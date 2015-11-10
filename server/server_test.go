package main

import (
	"bytes"
	"gopkg.in/tylerb/graceful.v1"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// This is an utility function to launch a graceful server.
// This is intended to be run as a goroutine as it takes an
// error channel as the first parameter.
func startServer(errors chan<- error, server *graceful.Server) {
	errors <- server.ListenAndServe()
}

// Test that the web server return "Hello world" and does not
// raise any exceptions or errors. This also starts and stops
// a web server instance for the duration of the test.
func TestHello(test *testing.T) {
	server, _ := NewServer()
	errors := make(chan error)

	go startServer(errors, server)
	defer server.Stop(3 * time.Second)

	resp, err := http.Post("http://localhost:8005", "text/plain", nil)
	if err != nil {
		test.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		test.Fatal(err)
	}

	expected := ([]byte)("Hello world!")

	ok := bytes.Equal(expected, body)
	if !ok {
		test.Fatal("Server did not greet us properly!")
	}

	if len(errors) > 0 {
		err, ok := <-errors
		if !ok {
			// TODO: You should do this properly
			test.Fatal("There was an error, but we missed it (too soon or too late")
		}
		if err != nil {
			test.Fatal(err)
		}
	}
}
