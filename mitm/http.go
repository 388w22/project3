package main

import (
	"net/http"
)

// PassthroughRequest should take the incoming request r
// and send it to the HTTP server located at endpoint,
// then mirror the response back to w.
// It should make no changes to the incoming request.
func PassthroughRequest(w http.ResponseWriter, r *http.Request, endpoint string) {
	panic("PassthroughRequest unimplemented!")
}

// InterceptAndRelayRequest should take the incoming request r,
// and if it has a `to` parameter in the body, change it to spoofed.
// It should then relay this request to the HTTP server located at endpoint,
// feeding the response back to the client, replacing any occurrences
// of spoofed with original.
//
// You may assume that the request is a POST request
// that contains a valid application/x-www-form-urlencoded body
// (and should thus only call this function
// with requests that fit these requirements).
func InterceptAndRelayRequest(w http.ResponseWriter, r *http.Request, endpoint, spoofed string) {
	panic("InterceptAndRelayRequest unimplemented!")
}
