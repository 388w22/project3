package thievery

import (
	"fmt"
	"log"
	"os"
)

// StealClientCookies helps steal cookies sent by the client.
// The actual stealing logic is up to you to deduce;
// this simply marks them as stolen for grading.
// This function may change in grading; do not hardcode its functionality.
func StealClientCookie(name string, value string) {
	fmt.Println("MITM:   Intercepted Cookie Sent By Client")
	fmt.Println("        Name: ", name)
	fmt.Println("        Value:", value)
}

// StealClientCookies helps steal cookies set by the server.
// The actual stealing logic is up to you to deduce;
// this simply marks them as stolen for grading.
// This function may change in grading; do not hardcode its functionality.
func StealServerCookie(name string, value string) {
	fmt.Println("MITM:   Intercepted Cookie Set By Server")
	fmt.Println("        Name: ", name)
	fmt.Println("        Value:", value)
}

// StealClientCookies helps steal login credentials.
// The actual stealing logic is up to you to deduce;
// this simply marks them as stolen for grading.
// This function may change in grading; do not hardcode its functionality.
func StealCredentials(username string, password string) {
	fmt.Println("MITM:   Intercepted Credentials")
	fmt.Println("        Username:", username)
	fmt.Println("        Password:", password)
}

// StealClientCookies helps steal files sent by the server.
// Write the contents of the file to the returned *os.File,
// then call Close on it.
// This function may change in grading; do not hardcode its functionality.
func StealFile(name string) *os.File {
	f, err := os.Create("/files/" + name)
	if err != nil {
		log.Panic(err)
	}
	return f
}
