package main

import (
	"errors"
	"fmt"
	"github.com/phzfi/RIC/server/config"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/operator"
	"github.com/phzfi/RIC/server/ops"
	"github.com/phzfi/RIC/server/testutils"
	"github.com/valyala/fasthttp"
	"log"
	"net"
	"time"
)

// Port will be incremented for each server created in the test
var port = 8022

// This is an utility function to launch a server.
func startServer() (server *fasthttp.Server, ln net.Listener, srverr chan error) {
	// Start the server
	conf, err := config.ReadConfig("config/testconfig.ini")

	if err != nil {
		log.Fatal("Error while reading config" + err.Error())
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
func SetupOperatorSource() (o operator.Operator, src ops.ImageSource) {
	o = operator.MakeDefault(512*1024*1024, "/tmp/RIC_testimagecache")
	src = ops.MakeImageSource()
	src.AddRoot("./")
	return
}

// Gets blob from server. package variable port is used as port and localhost as address
func getBlobFromServer(getname string) (blob []byte, err error) {
	statuscode, blob, err := fasthttp.Get(nil, fmt.Sprintf("http://localhost:%d/", port)+getname)
	if statuscode != 200 {
		return nil, errors.New(fmt.Sprintf("Server returned %d", statuscode))
	}
	if err != nil {
		return
	}

	return
}

// Tests getting images. c.Testfn is treated as GET string (filename?params).
// Executes the tests supported by TestCaseAll.
func testGetImages(cases []testutils.TestCaseAll) (err error) {

	s, ln, srverr := startServer()
	defer stopServer(s, ln, srverr)

	tolerance := 0.002

	// Todo: this threading is copied from common_test.go unify it to single implementation (DRY)
	sem := make(chan error, len(cases))
	for _, c := range cases {
		go func(tc testutils.TestCaseAll) {
			logging.Debug(fmt.Sprintf("Testing get: %v, %v, %v, %v, %v", tc.Testfn, tc.Reffn, tc.Resfn, tc.W, tc.H))
			blob, err := getBlobFromServer(tc.Testfn)
			if err != nil {
				sem <- err
				return
			}

			sem <- testutils.TestAll(tc, blob, tolerance)
		}(c)
	}

	for range cases {
		var verr = <-sem
		if verr != nil && err == nil {
			err = verr
		}
	}
	return
}
