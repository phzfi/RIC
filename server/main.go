package main

import (
	"net/http"
)

type MyServer struct {
}

func (*MyServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	result := "Hello world!"
	omg := ([]byte)(result)
	i, err := writer.Write(omg)
}

func main() {
	server := &http.Server{
		Addr:           ":8000",
		Handler:        null,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
