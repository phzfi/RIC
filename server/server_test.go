package main

import (
	"bytes"
	"errors"
	"github.com/joonazan/imagick/imagick"
	"github.com/phzfi/RIC/server/images"
	"gopkg.in/tylerb/graceful.v1"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"fmt"
	"time"
)

// This is an utility function to launch a graceful server.
// This is intended to be run as a goroutine as it takes an
// error channel as the first parameter.
func startServer(srverr chan<- error, server *graceful.Server) {
	srverr <- server.ListenAndServe()
}


// Stop server and block until stopped
func stopServer(server *graceful.Server) {
	server.Stop(0)
	<- server.StopChan()
}


// Test that the web server return "Hello world" and does not
// raise any exceptions or errors. This also starts and stops
// a web server instance for the duration of the t.
func TestHello(t *testing.T) {
	
	// Start the server
	server, _ := NewServer(500000)
	srverr := make(chan error)
	go startServer(srverr, server)
	defer stopServer(server)
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Post("http://localhost:8005", "text/plain", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := ([]byte)("Hello world!")

	if !bytes.Equal(expected, body) {
		t.Fatal("Server did not greet us properly!")
	}
	
	// Yes we do this properly
	select {
		case err := <- srverr:
			t.Fatal(err)
		default:
	}
}

const TOLERANCE = 0.0005

// Test that the web server returns requested JPG image at right size
func TestGetJPGFromServer(t *testing.T) {
	err := testGetImageFromServer("toget.jpg", "getref.jpg", TOLERANCE)
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns requested JPG image at right size
func TestGetPNGFromServer(t *testing.T) {
	err := testGetImageFromServer("toget.png", "getref.png", TOLERANCE)
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns requested JPG image at right size
func TestGetGIFFromServer(t *testing.T) {
	err := testGetImageFromServer("toget.gif", "getref.gif", TOLERANCE)
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns requested JPG image at right size
func TestGetTIFFFromServer(t *testing.T) {
	err := testGetImageFromServer("toget.tiff", "getref.tiff", TOLERANCE)
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns requested JPG image at right size
func TestGetWEBPFromServer(t *testing.T) {
	err := testGetImageFromServer("toget.webp", "getref.webp", TOLERANCE)
	if err != nil {
		t.Fatal(err)
		return
	}
}

func testGetImageFromServer(getname string, refname string, tolerance float64) (err error) {

	// Start the server
	server, _ := NewServer(500000)
	srverr := make(chan error)
	go startServer(srverr, server)
	defer stopServer(server)
	time.Sleep(100 * time.Millisecond)

	// Get and read requested image (blob) of size 100x100 from the server
	resp, err := http.Get("http://localhost:8005/testimages/server/" + getname + "?width=100&height=100")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	
	// Save retrieved image to testresults

	mw := imagick.NewMagickWand()
	//defer mw.Destroy()
	err = mw.ReadImageBlob(body)
	if err != nil {
		return
	}
	err = mw.WriteImage(filepath.FromSlash("testresults/server/" + getname))
	if err != nil {
		return
	}
	
	// Get distortion compared to refrence image and check it is inside tolerance
	distortion, err := getDistortion(body, refname)
	if err != nil {
		return
	}
	if distortion > tolerance {
		return errors.New(fmt.Sprintf("Bad image returned. Distortion: %v, Tolerance: %v", distortion, tolerance))
	}

	// Check for server errors
	select {
		case err = <- srverr:
			return
		default:
	}
	
	return
}

func getDistortion(imageblob images.ImageBlob, filename_cmp string) (distortion float64, err error) {
	const image_folder = "testimages/server/"
	imagick.Initialize()
	defer imagick.Terminate()

	mw_cmp := imagick.NewMagickWand()
	defer mw_cmp.Destroy()

	err = mw_cmp.ReadImage(filepath.FromSlash(image_folder + filename_cmp))
	if err != nil {
		err = errors.New("Could not load reference image:" + err.Error())
		return
	}

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err = mw.ReadImageBlob(imageblob)
	if err != nil {
		return
	}

	trash, distortion := mw.CompareImages(mw_cmp, imagick.METRIC_MEAN_SQUARED_ERROR)
	trash.Destroy()

	return
}
