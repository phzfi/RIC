package main

import (
	"fmt"
	"github.com/phzfi/RIC/server/cache"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/ops"
	"github.com/phzfi/RIC/server/configuration"
	"github.com/valyala/fasthttp"
	"net"
	"time"
)

// Port will be incremented for each server created in the test
var port = 8022

// This is an utility function to launch a server.
func startServer() (server *fasthttp.Server, ln net.Listener, srverr chan error) {
	// Start the server
	conf, err := configuration.ReadConfig("testconfig.ini")

	if err != nil {
		logging.Debug("Error while reading config" + err.Error())
		return
	}

	port++
	server, _, ln = NewServer(port, 500000, conf)
	srverr = make(chan error)
	go func() {
		srverr <- server.Serve(ln)
	}()
	time.Sleep(100 * time.Millisecond)
	return
}

// Stop server and block until stopped
func stopServer(server *fasthttp.Server, ln net.Listener, srverr chan error) error {
	_ = ln.Close()
	select {
	case err := <-srverr:
		return err
	default:
	}
	return nil
}

// A setup function for simple Operator tests. Creates a cache, operator that
// uses the cache and ImageSource operation with current working dir as root.
// Returns the operator and the source operation.
func SetupOperatorSource() (operator cache.Operator, src ops.ImageSource) {
	operator = cache.MakeOperator(512 * 1024 * 1024)
	src = ops.MakeImageSource()
	src.AddRoot("./")
	return
}

// Gets blob from server. package variable port is used as port and localhost as address
func getBlobFromServer(getname string) (blob images.ImageBlob, err error) {
	_, blob, err = fasthttp.Get(nil, fmt.Sprintf("http://localhost:%d/", port)+getname)
	if err != nil {
		return
	}

	return
}

// Tests getting images. c.Testfn is treated as GET string (filename?params).
// Executes the tests supported by TestCaseAll.
func testGetImages(cases []images.TestCaseAll) (err error) {

	s, ln, srverr := startServer()
	defer stopServer(s, ln, srverr)

	tolerance := 0.002

	for _, c := range cases {
		logging.Debug(fmt.Sprintf("Testing get: %v, %v, %v, %v, %v", c.Testfn, c.Reffn, c.Resfn, c.W, c.H))
		blob, err := getBlobFromServer(c.Testfn)
		if err != nil {
			return err
		}

		err = images.TestAll(c, blob, tolerance)
		if err != nil {
			return err
		}
	}
	return
}
