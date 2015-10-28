package main

import (
	"bytes"
	"gopkg.in/tylerb/graceful.v1"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func startServer(errors chan<- error, server *graceful.Server) {
	errors <- server.ListenAndServe()
}

func TestHello(test *testing.T) {
	server := NewServer()
	errors := make(chan error)

	go startServer(errors, server)
	defer server.Stop(5 * time.Second)

	resp, err := http.Get("http://localhost:8005")
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
			// TODO: You could do this *properly*
			test.Fatal("There was an error, but we missed it (too soon or too late")
		}
		if err != nil {
			test.Fatal(err)
		}
	}
}
