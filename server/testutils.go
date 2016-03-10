package main

import (
	"fmt"
	"github.com/phzfi/RIC/server/cache"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/ops"
	"github.com/valyala/fasthttp"
	"net"
	"time"
)

// Port will be incremented for each server created in the test
var port = 8022

// This is an utility function to launch a server.
func startServer() (server *fasthttp.Server, ln net.Listener, srverr chan error) {
	// Start the server
	port++
	server, _, ln = NewServer(port, 500000)
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
	operator = cache.MakeOperator(512*1024*1024, "/tmp/RIC_testimagecache")
	src = ops.MakeImageSource()
	src.AddRoot("./")
	return
}

// Gets blob from server. package variable port is used as port and lovalhost as address
func getBlobFromServer(getname string) (blob images.ImageBlob, err error) {
	_, blob, err = fasthttp.Get(nil, fmt.Sprintf("http://localhost:%d/", port)+getname)
	if err != nil {
		return
	}

	return
}

// Tests getting images. c.Testfn is treated as GET string (filename?params).
// Executes the tests supported by TestCaseAll.
func testGetImages(cases []images.TestCaseAll, tolerance float64) (err error) {

	s, ln, srverr := startServer()
	defer stopServer(s, ln, srverr)

	// Todo: this threading is copied from common_test.go unify it to single implementation (DRY)
	sem := make(chan error, len(cases))
	for _, c := range cases {
		go func (tc images.TestCaseAll) {
			logging.Debug(fmt.Sprintf("Testing get: %v, %v, %v, %v, %v", tc.Testfn, tc.Reffn, tc.Resfn, tc.W, tc.H))
			blob, err := getBlobFromServer(tc.Testfn)
			if err != nil {
				sem <- err
				return
			}

			sem <- images.TestAll(tc, blob, tolerance)
		} (c)
	}

	for range cases {
		var verr = <- sem
		if verr != nil && err == nil {
			err = verr
		}
	}
	return
}
