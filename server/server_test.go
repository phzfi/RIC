package main

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"github.com/valyala/fasthttp"
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"
)

var port = 8022

// This is an utility function to launch a server.
// This is intended to be run as a goroutine as it takes an
// error channel as the first parameter.
func startServer(srverr chan<- error, server *fasthttp.Server, ln net.Listener) {
	srverr <- server.Serve(ln)
}

// Stop server and block until stopped
func stopServer(ln net.Listener) {
	ln.Close()
}

// Test that the web server return "Hello world" and does not
// raise any exceptions or errors. This also starts and stops
// a web server instance for the duration of the t.
func TestHello(t *testing.T) {

	// Start the server
	port++
	server, _, ln := NewServer(port, 500000)
	srverr := make(chan error)
	go startServer(srverr, server, ln)
	defer stopServer(ln)
	time.Sleep(100 * time.Millisecond)

	_, body, err := fasthttp.Post(nil, fmt.Sprintf("http://localhost:%d",port), nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := ([]byte)("Hello world!")
	if !bytes.Equal(expected, body) {
		t.Fatal("Server did not greet us properly!")
	}

	// Yes we do this properly
	select {
	case err := <-srverr:
		t.Fatal(err)
	default:
	}
}

func TestGetImageBySize(t *testing.T) {

	testfolder := "testimages/server/"
	resfolder := "testresults/server/"
	tolerance := 0.002
	
	cases := [...]images.TestCaseAll {
		{images.TestCase{testfolder + "01.jpg", testfolder + "01_100x100.jpg", resfolder + "01_100x100.jpg"}, "JPEG", 100, 100},
		{images.TestCase{testfolder + "01.jpg", testfolder + "01_200x100.jpg", resfolder + "01_200x100.jpg"}, "JPEG", 200, 100},
		{images.TestCase{testfolder + "01.jpg", testfolder + "01_300x100.jpg", resfolder + "01_300x100.jpg"}, "JPEG", 300, 100},
		{images.TestCase{testfolder + "01.jpg", testfolder + "01_300x50.jpg", resfolder + "01_300x50.jpg"}, "JPEG", 300, 50},
		{images.TestCase{testfolder + "01.webp", testfolder + "01_100x100.webp", resfolder + "01_100x100.webp"}, "WEBP", 100, 100},
		{images.TestCase{testfolder + "01.webp", testfolder + "01_200x100.webp", resfolder + "01_200x100.webp"}, "WEBP", 200, 100},
		{images.TestCase{testfolder + "01.webp", testfolder + "01_300x100.webp", resfolder + "01_300x100.webp"}, "WEBP", 300, 100},
		{images.TestCase{testfolder + "01.webp", testfolder + "01_300x50.webp", resfolder + "01_300x50.webp"}, "WEBP", 300, 50},
		{images.TestCase{testfolder + "02.jpg", testfolder + "02_100x100.jpg", resfolder + "02_100x100.jpg"}, "JPEG", 100, 100},
		{images.TestCase{testfolder + "02.jpg", testfolder + "02_200x100.jpg", resfolder + "02_200x100.jpg"}, "JPEG", 200, 100},
		{images.TestCase{testfolder + "02.jpg", testfolder + "02_300x100.jpg", resfolder + "02_300x100.jpg"}, "JPEG", 300, 100},
		{images.TestCase{testfolder + "02.jpg", testfolder + "02_300x50.jpg", resfolder + "02_300x50.jpg"}, "JPEG", 300, 50},
		{images.TestCase{testfolder + "02.webp", testfolder + "02_100x100.webp", resfolder + "02_100x100.webp"}, "WEBP", 100, 100},
		{images.TestCase{testfolder + "02.webp", testfolder + "02_200x100.webp", resfolder + "02_200x100.webp"}, "WEBP", 200, 100},
		{images.TestCase{testfolder + "02.webp", testfolder + "02_300x100.webp", resfolder + "02_300x100.webp"}, "WEBP", 300, 100},
		{images.TestCase{testfolder + "02.webp", testfolder + "02_300x50.webp", resfolder + "02_300x50.webp"}, "WEBP", 300, 50},
	}
	
	
	for _, c := range cases {	
		logging.Debug(fmt.Sprintf("Testing get: %v, %v, %v, %v, %v, %v", c.Testfn, c.Reffn, c.W, c.H, c.Format, c.Resfn))
		blob, err := GetBlobFromServer(fmt.Sprintf("%v?width=%v&height=%v", c.Testfn, c.W, c.H))
		if err != nil {
			t.Fatal(err)
		}

		err = images.TestAll(c, blob, tolerance)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func GetBlobFromServer(getname string) (blob images.ImageBlob, err error) {
	// Start the server
	port++
	server, _, ln := NewServer(port, 500000)
	srverr := make(chan error)
	go startServer(srverr, server, ln)
	defer stopServer(ln)
	time.Sleep(100 * time.Millisecond)

	_, blob, err = fasthttp.Get(nil, fmt.Sprintf("http://localhost:%d/", port) + getname)
	if err != nil {
		return
	}

	// Check for server errors
	select {
	case err = <-srverr:
		return
	default:
	}

	return
}

// TODO: Tests for different modes and parameters (get by aspect etc.)
