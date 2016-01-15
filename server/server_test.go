package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/joonazan/imagick/imagick"
	"github.com/phzfi/RIC/server/images"
	"gopkg.in/tylerb/graceful.v1"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
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
	server.Stop(10 * time.Millisecond)
	<-server.StopChan()
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
	case err := <-srverr:
		t.Fatal(err)
	default:
	}
}

const TOLERANCE = 0.0005

// Test that the web server returns requested JPG image at right size
func TestGetJPGFromServer(t *testing.T) {
	err := testGetImageFromServer("toget.jpg", "?width=100&height=100", "getref.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns requested JPG image at right size
func TestGetPNGFromServer(t *testing.T) {
	err := testGetImageFromServer("toget.png", "?width=100&height=100", "getref.png")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns requested JPG image at right size
func TestGetGIFFromServer(t *testing.T) {
	err := testGetImageFromServer("toget.gif", "?width=100&height=100", "getref.gif")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns requested JPG image at right size
func TestGetTIFFFromServer(t *testing.T) {
	err := testGetImageFromServer("toget.tiff", "?width=100&height=100", "getref.tiff")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns requested JPG image at right size
func TestGetWEBPFromServer(t *testing.T) {
	err := testGetImageFromServer("toget.webp", "?width=100&height=100", "getref.webp")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns JPG requested by defining width only
func TestGetJPGByWidth(t *testing.T) {
	err := testGetImageFromServer("toget.jpg", "?width=200", "jpgbywidth.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns PNG requested by defining width only
func TestGetPNGByWidth(t *testing.T) {
	err := testGetImageFromServer("toget.png", "?width=200", "pngbywidth.png")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns JPG requested by defining width only
func TestGetJPGByHeight(t *testing.T) {
	err := testGetImageFromServer("toget.jpg", "?height=200", "jpgbyheight.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns PNG requested by defining width only
func TestGetPNGByHeight(t *testing.T) {
	err := testGetImageFromServer("toget.png", "?height=200", "pngbyheight.png")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns JPG fit to given dimensions
func TestGetJPGFitByWidth(t *testing.T) {
	err := testGetImageFromServer("toget.jpg", "?width=200&height=1000&mode=fit", "jpgfitbywidth.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns JPG fit to given dimensions
func TestGetJPGFitByHeight(t *testing.T) {
	err := testGetImageFromServer("toget.jpg", "?width=1000&height=200&mode=fit", "jpgfitbyheight.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns PNG fit to given dimensions
func TestGetPNGFitByWidth(t *testing.T) {
	err := testGetImageFromServer("toget.png", "?width=200&height=1000&mode=fit", "pngfitbywidth.png")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns PNG fit to given dimensions
func TestGetPNGFitByHeight(t *testing.T) {
	err := testGetImageFromServer("toget.png", "?width=1000&height=200&mode=fit", "pngfitbyheight.png")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the test fails with wrong sized PNG image (testing that tests work)
func TestGetPNGFitByHeightFail(t *testing.T) {
	err := testGetImageFromServer("toget.png", "?width=180&height=200", "pngbyheight.png")
	if err == nil {
		t.Fatal("Test passed even with bad image - there is something seriously wrong with these tests")
		return
	}
}

// Test that the web server returns JPG at original size
func TestGetJPGOriginalSize(t *testing.T) {
	err := testGetImageFromServer("toget.jpg", "", "origref.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
}

// Test that the web server returns original image
func TestGetOriginal(t *testing.T) {
	err := testGetImageFromServer("toget", "", "origref")
	if err != nil {
		t.Fatal(err)
		return
	}
}

func testGetImageFromServer(getname string, params string, refname string) (err error) {

	// Start the server
	server, _ := NewServer(500000)
	srverr := make(chan error)
	go startServer(srverr, server)
	defer stopServer(server)
	time.Sleep(100 * time.Millisecond)

	// Get and read requested image (blob) of size 100x100 from the server
	resp, err := http.Get("http://localhost:8005/testimages/server/" + getname + params)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	imagick.Initialize()
	defer imagick.Terminate()

	// Save retrieved image to testresults
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err = mw.ReadImageBlob(body)
	if err != nil {
		return
	}
	err = mw.WriteImage(filepath.FromSlash("testresults/server/" + refname))
	if err != nil {
		return
	}

	// Get distortion compared to refrence image and check it is inside tolerance
	distortion, err := getDistortion(body, refname)
	if err != nil {
		return
	}
	if distortion > TOLERANCE {
		return errors.New(fmt.Sprintf("Bad image returned. Distortion: %v, Tolerance: %v", distortion, TOLERANCE))
	}

	// Check for server errors
	select {
	case err = <-srverr:
		return
	default:
	}

	return
}

func getDistortion(imageblob images.ImageBlob, filename_cmp string) (distortion float64, err error) {
	const image_folder = "testimages/server/"

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
