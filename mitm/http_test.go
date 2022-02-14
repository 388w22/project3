// This file has a lot of repetition, but we think that
// it makes it easier to linearly read what is happening
// in each test rather than passing control flow between functions.
// In an actual set of unit tests, you'd probably create
// helper functions.

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	// The header used across all tests for
	// client-to-server communication.
	ctsHeaderKey   = "X-388-Request-Header"
	ctsHeaderValue = "request header value"
	// The header used across all tests for
	// server-to-client communication.
	stcHeaderKey   = "X-388-Response-Header"
	stcHeaderValue = "response header value"

	uri = "/test/uri"
)

func TestPassthroughRequest(t *testing.T) {
	type requestResult struct {
		request *http.Request
		body    string
	}

	body := "test body"
	r := httptest.NewRequest("TEST", uri, strings.NewReader(body))
	r.Header.Add(ctsHeaderKey, ctsHeaderValue)

	w := httptest.NewRecorder()

	requests := make(chan requestResult, 1)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		requests <- requestResult{
			request: r,
			body:    string(b),
		}
		w.Header().Add(stcHeaderKey, stcHeaderValue)
		io.WriteString(w, "test response body")
	}))
	defer s.Close()

	PassthroughRequest(w, r, s.URL)

	var received requestResult
	// Wait up to 100 milliseconds for the response;
	// if we don't receive it by then, assume it's never coming.
	select {
	case received = <-requests:
	case <-time.After(100 * time.Millisecond):
		t.Error("request not received by real server")
		t.FailNow()
	}

	if received.request.Method != "TEST" {
		t.Errorf("real server expected method %q but got %q", "TEST", received.request.Method)
	}
	if received.request.RequestURI != uri {
		t.Errorf("real server expected URI %q but got %q", uri, received.request.RequestURI)
	}
	if len(received.request.Header.Values(ctsHeaderKey)) == 0 {
		t.Errorf("real server did not receive %q header sent in original request", ctsHeaderKey)
	} else if received.request.Header.Get(ctsHeaderKey) != ctsHeaderValue {
		t.Errorf("real server did not receive correct header value for key %s, expected %q but got %q", ctsHeaderKey, ctsHeaderValue, received.request.Header.Get(ctsHeaderKey))
	}
	if len(w.Result().Header.Values(stcHeaderKey)) == 0 {
		t.Errorf("client did not receive %q header sent in response", stcHeaderKey)
	} else if w.Result().Header.Get(stcHeaderKey) != stcHeaderValue {
		t.Errorf("client did not receive correct header value in response for key %s, expected %q but got %q", stcHeaderKey, stcHeaderValue, w.Result().Header.Get(stcHeaderKey))
	}
	if received.body != body {
		t.Errorf("real server expected body %q but got %q", body, received.body)
	}
	cl, _ := strconv.Atoi(w.Result().Header.Get("Content-Length"))
	if cl != w.Body.Len() {
		t.Errorf("client got response with declared Content-Length of %d bytes but actual body length of %d bytes", cl, w.Body.Len())
	}
	if w.Body.String() != "test response body" {
		t.Errorf("client expected response body %q but got %q", "test response body", w.Body.String())
	}
}

func TestInterceptAndRelayNoChanges(t *testing.T) {
	type requestResult struct {
		request *http.Request
		body    url.Values
	}

	body := "real=test&loc=body"
	r := httptest.NewRequest("POST", uri, strings.NewReader(body))
	r.Header.Add(ctsHeaderKey, ctsHeaderValue)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	requests := make(chan requestResult, 1)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		v, _ := url.ParseQuery(string(b))
		requests <- requestResult{
			request: r,
			body:    v,
		}
		w.Header().Add(stcHeaderKey, stcHeaderValue)
		io.WriteString(w, "test response body")
	}))
	defer s.Close()

	InterceptAndRelayRequest(w, r, s.URL, "fake")

	var received requestResult
	select {
	case received = <-requests:
	case <-time.After(100 * time.Millisecond):
		t.Error("request not received by real server")
		t.FailNow()
	}

	if received.request.Method != "POST" {
		t.Errorf("real server expected method %q but got %q", "POST", received.request.Method)
	}
	if received.request.RequestURI != uri {
		t.Errorf("real server expected URI %q but got %q", uri, received.request.RequestURI)
	}
	if len(received.request.Header.Values(ctsHeaderKey)) == 0 {
		t.Errorf("real server did not receive %q header sent in original request", ctsHeaderKey)
	} else if received.request.Header.Get(ctsHeaderKey) != ctsHeaderValue {
		t.Errorf("real server did not receive correct header value for key %s, expected %q but got %q", ctsHeaderKey, ctsHeaderValue, received.request.Header.Get(ctsHeaderKey))
	}
	if len(w.Result().Header.Values(stcHeaderKey)) == 0 {
		t.Errorf("client did not receive %q header sent in response", stcHeaderKey)
	} else if w.Result().Header.Get(stcHeaderKey) != stcHeaderValue {
		t.Errorf("client did not receive correct header value in response for key %s, expected %q but got %q", stcHeaderKey, stcHeaderValue, w.Result().Header.Get(stcHeaderKey))
	}
	ex, _ := url.ParseQuery(body)
	if !reflect.DeepEqual(received.body, ex) {
		t.Errorf("real server expected body %q but got %q", ex, received.body)
	}
	cl, _ := strconv.Atoi(w.Result().Header.Get("Content-Length"))
	if cl != w.Body.Len() {
		t.Errorf("client got response with declared Content-Length of %d bytes but actual body length of %d bytes", cl, w.Body.Len())
	}
	if w.Body.String() != "test response body" {
		t.Errorf("client expected response body %q but got %q", "test response body", w.Body.String())
	}
}

func TestInterceptAndRelayChangeBoth(t *testing.T) {
	type requestResult struct {
		request *http.Request
		body    url.Values
	}

	body := "test=real&to=real"
	expectedAtServer := "test=real&to=not"
	expectedAtClient := "sabrina sent $1000 to real"
	r := httptest.NewRequest("POST", uri, strings.NewReader(body))
	r.Header.Add(ctsHeaderKey, ctsHeaderValue)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	requests := make(chan requestResult, 1)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		v, _ := url.ParseQuery(string(b))
		requests <- requestResult{
			request: r,
			body:    v,
		}
		w.Header().Add(stcHeaderKey, stcHeaderValue)
		io.WriteString(w, "sabrina sent $1000 to "+v.Get("to"))
	}))
	defer s.Close()

	InterceptAndRelayRequest(w, r, s.URL, "not")

	var received requestResult
	select {
	case received = <-requests:
	case <-time.After(100 * time.Millisecond):
		t.Error("request not received by real server")
		t.FailNow()
	}

	if received.request.Method != "POST" {
		t.Errorf("real server expected method %q but got %q", "POST", received.request.Method)
	}
	if received.request.RequestURI != uri {
		t.Errorf("real server expected URI %q but got %q", uri, received.request.RequestURI)
	}
	if len(received.request.Header.Values(ctsHeaderKey)) == 0 {
		t.Errorf("real server did not receive %q header sent in original request", ctsHeaderKey)
	} else if received.request.Header.Get(ctsHeaderKey) != ctsHeaderValue {
		t.Errorf("real server did not receive correct header value for key %s, expected %q but got %q", ctsHeaderKey, ctsHeaderValue, received.request.Header.Get(ctsHeaderKey))
	}
	if len(w.Result().Header.Values(stcHeaderKey)) == 0 {
		t.Errorf("client did not receive %q header sent in response", stcHeaderKey)
	} else if w.Result().Header.Get(stcHeaderKey) != stcHeaderValue {
		t.Errorf("client did not receive correct header value in response for key %s, expected %q but got %q", stcHeaderKey, stcHeaderValue, w.Result().Header.Get(stcHeaderKey))
	}
	ex, _ := url.ParseQuery(expectedAtServer)
	if !reflect.DeepEqual(received.body, ex) {
		t.Errorf("real server expected body %q but got %q", ex, received.body)
	}
	cl, _ := strconv.Atoi(w.Result().Header.Get("Content-Length"))
	if cl != w.Body.Len() {
		t.Errorf("client got response with declared Content-Length of %d bytes but actual body length of %d bytes", cl, w.Body.Len())
	}
	if w.Body.String() != expectedAtClient {
		t.Errorf("client expected response body %q but got %q", expectedAtClient, w.Body.String())
	}
}
