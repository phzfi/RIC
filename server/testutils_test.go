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
	"net"
	"time"
	"encoding/base64"

	"log"
	"os"
	"io/ioutil"
	"crypto/md5"
	"io"
	"bufio"
)

// Port will be incremented for each server created in the test
var increment = 0
var port int
var tokens = 3
var configPath = "config/testconfig.ini"

// This is an utility function to launch a server.
func startServer() (server *fasthttp.Server, listener net.Listener, serverErr chan error) {
	// Start the server
	conf := config.ReadConfig(configPath)
	port = conf.Server.Port + increment
	increment++
	logging.Debugf("Starting a new server on port %v", port)

	server, _, listener = NewServer(port, 500000, conf)
	serverErr = make(chan error)
	go func() {
		serverErr <- server.Serve(listener)
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
	conf := config.ReadConfig(configPath)
	o = operator.MakeWithDefaultCacheSet(512*1024*1024, conf.Server.CacheFolder, tokens)
	src = ops.MakeImageSource()
	src.AddRoot("./")
	return
}

// Gets blob from server. package variable port is used as port and localhost as address
func getBlobFromServer(getname string) (blob []byte, err error)  {
	requestURL := fmt.Sprintf("http://localhost:%d/", port)+getname
	logging.Debugf("Requesting URL: %v", requestURL)
	statusCode, blob, httpErr := fasthttp.Get(nil, requestURL)
	logging.Debugf("STATUS CODE: %v; ERR: %v", statusCode, httpErr)

	err = httpErr
	if statusCode != 200 {
		err = errors.New(fmt.Sprintf("Server returned %d", statusCode))
	}

	return
}

// Tests getting images. c.TestFilename is treated as GET string (filename?params).
// Executes the tests supported by TestCaseAll.
func testGetImages(cases []testutils.TestCaseAll) (err error) {

	server, listener, serverErr := startServer()
	defer stopServer(server, listener, serverErr)

	tolerance := 0.002

	// Todo: this threading is copied from common_test.go unify it to single implementation (DRY)
	sem := make(chan error, len(cases))
	for _, testCase := range cases {
		go func(testCase testutils.TestCaseAll) {
			logging.Debug(fmt.Sprintf("Testing get: %v%v, %v, %v, %v, %v", testCase.TestFilename, testCase.TestParameters, testCase.ReferenceFilename, testCase.ResultFilename, testCase.W, testCase.H))
			requestFilename := base64.StdEncoding.EncodeToString([]byte(testCase.TestFilename)) + testCase.TestParameters
			blob, err := getBlobFromServer(requestFilename)

			//fmt.Printf("--->BLOB BEGINNING: %x", len(blob))
			if err != nil {
				logging.Debugf("%v", err)
			}

			if err != nil {
				sem <- err
				return
			}

			sem <- testutils.TestAll(testCase, blob, tolerance)
		}(testCase)
	}

	// HACKY(?) sync for go routines
	for range cases {
		var verr = <-sem
		if verr != nil && err == nil {
			err = verr
		}
	}
	return
}


func createTestImageFolderStructure() {

	conf := config.ReadConfig(configPath)

	emptyDirectory(conf.Server.CacheFolder)
	//emptyDirectory(conf.Server.ImageFolder)

	testFolder := "testimages/server/"
	files, err := ioutil.ReadDir(testFolder)
	if err != nil {
		logging.Debug("Could not read test image source files")
		log.Fatal(err)
	}

	for _, filename := range files {
		sourcePath := testFolder + filename.Name()
		md5Hash := md5.New()
		io.WriteString(md5Hash, sourcePath)
		md5Filename := fmt.Sprintf("%x", md5Hash.Sum(nil))

		file, copyErr := os.Create(conf.Server.ImageFolder + "/" + md5Filename)
		f, _ := os.Open(sourcePath)
		_, copyErr = io.Copy(file, bufio.NewReader(f))


		if copyErr != nil {
			fmt.Printf("Failed to copy file: %s: %s", sourcePath, copyErr)
		}
	}
}

func emptyDirectory(directoryPath string) {
	_ = os.RemoveAll(directoryPath)
	err := os.MkdirAll(directoryPath, 777)
	if err != nil {
		log.Fatal(err)
	}
}