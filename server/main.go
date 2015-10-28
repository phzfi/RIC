package main

import (
	"gopkg.in/tylerb/graceful.v1"
	"log"
	"net/http"
	"time"
)

type MyHandler struct {
}

func (*MyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	result := "Hello world!"
	omg := ([]byte)(result)
	_, err := writer.Write(omg)
	if err != nil {
		log.Println(err)
	}
}

func NewServer() *graceful.Server {
	server := &graceful.Server{
		Timeout: 5 * time.Second,
		Server: &http.Server{
			Addr:           ":8005",
			Handler:        &MyHandler{},
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
	return server
}

func main() {
	server := NewServer()
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
