package main

import (
	"bytes"
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/valyala/fasthttp"
	"testing"
)

// Test that the web server return "Hello world" and does not
// raise any exceptions or errors. This also starts and stops
// a web server instance for the duration of the t.
func TestHello(t *testing.T) {
	s, ln, srverr := startServer()
	defer stopServer(s, ln, srverr)

	_, body, err := fasthttp.Post(nil, fmt.Sprintf("http://localhost:%d", port), nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := ([]byte)("Hello world!")
	if !bytes.Equal(expected, body) {
		t.Fatal("Server did not greet us properly!")
	}
}

// Test GETting different sizes and formats
func TestGetImageBySize(t *testing.T) {

	testfolder := "testimages/server/"
	resfolder := "testresults/server/"
	tolerance := 0.002

	cases := []images.TestCaseAll{
		{images.TestCase{testfolder + "01.jpg?width=100&height=100", testfolder + "01_100x100.jpg", resfolder + "01_100x100.jpg"}, "JPEG", 100, 100},
		{images.TestCase{testfolder + "01.jpg?width=200&height=100", testfolder + "01_200x100.jpg", resfolder + "01_200x100.jpg"}, "JPEG", 200, 100},
		{images.TestCase{testfolder + "01.jpg?width=300&height=100", testfolder + "01_300x100.jpg", resfolder + "01_300x100.jpg"}, "JPEG", 300, 100},
		{images.TestCase{testfolder + "01.jpg?width=300&height=50", testfolder + "01_300x50.jpg", resfolder + "01_300x50.jpg"}, "JPEG", 300, 50},
		{images.TestCase{testfolder + "01.webp?width=100&height=100", testfolder + "01_100x100.webp", resfolder + "01_100x100.webp"}, "WEBP", 100, 100},
		{images.TestCase{testfolder + "01.webp?width=200&height=100", testfolder + "01_200x100.webp", resfolder + "01_200x100.webp"}, "WEBP", 200, 100},
		{images.TestCase{testfolder + "01.webp?width=300&height=100", testfolder + "01_300x100.webp", resfolder + "01_300x100.webp"}, "WEBP", 300, 100},
		{images.TestCase{testfolder + "01.webp?width=300&height=50", testfolder + "01_300x50.webp", resfolder + "01_300x50.webp"}, "WEBP", 300, 50},
		{images.TestCase{testfolder + "02.jpg?width=100&height=100", testfolder + "02_100x100.jpg", resfolder + "02_100x100.jpg"}, "JPEG", 100, 100},
		{images.TestCase{testfolder + "02.jpg?width=200&height=100", testfolder + "02_200x100.jpg", resfolder + "02_200x100.jpg"}, "JPEG", 200, 100},
		{images.TestCase{testfolder + "02.jpg?width=300&height=100", testfolder + "02_300x100.jpg", resfolder + "02_300x100.jpg"}, "JPEG", 300, 100},
		{images.TestCase{testfolder + "02.jpg?width=300&height=50", testfolder + "02_300x50.jpg", resfolder + "02_300x50.jpg"}, "JPEG", 300, 50},
		{images.TestCase{testfolder + "02.webp?width=100&height=100", testfolder + "02_100x100.webp", resfolder + "02_100x100.webp"}, "WEBP", 100, 100},
		{images.TestCase{testfolder + "02.webp?width=200&height=100", testfolder + "02_200x100.webp", resfolder + "02_200x100.webp"}, "WEBP", 200, 100},
		{images.TestCase{testfolder + "02.webp?width=300&height=100", testfolder + "02_300x100.webp", resfolder + "02_300x100.webp"}, "WEBP", 300, 100},
		{images.TestCase{testfolder + "02.webp?width=300&height=50", testfolder + "02_300x50.webp", resfolder + "02_300x50.webp"}, "WEBP", 300, 50},
	}

	err := testGetImages(cases, tolerance)
	if err != nil {
		t.Fatal(err)
	}
}

// Test GETting different sized and formats with mode=fit
func TestGetImageFit(t *testing.T) {

	testfolder := "testimages/server/"
	resfolder := "testresults/server/"
	tolerance := 0.002

	cases := []images.TestCaseAll{
		{images.TestCase{testfolder + "03.jpg?width=500&height=100&mode=fit", testfolder + "03_h100.jpg", resfolder + "03_h100.jpg"}, "JPEG", 143, 100},
		{images.TestCase{testfolder + "03.jpg?width=200&height=500&mode=fit", testfolder + "03_w200.jpg", resfolder + "03_w200.jpg"}, "JPEG", 200, 140},
		{images.TestCase{testfolder + "03.jpg?width=300&height=500&mode=fit", testfolder + "03_w300.jpg", resfolder + "03_w300.jpg"}, "JPEG", 300, 210},
		{images.TestCase{testfolder + "03.jpg?width=500&height=50&mode=fit", testfolder + "03_h50.jpg", resfolder + "03_h50.jpg"}, "JPEG", 71, 50},
		{images.TestCase{testfolder + "03.webp?width=500&height=100&mode=fit", testfolder + "03_h100.webp", resfolder + "03_h100.webp"}, "WEBP", 143, 100},
		{images.TestCase{testfolder + "03.webp?width=200&height=500&mode=fit", testfolder + "03_w200.webp", resfolder + "03_w200.webp"}, "WEBP", 200, 140},
		{images.TestCase{testfolder + "03.webp?width=300&height=500&mode=fit", testfolder + "03_w300.webp", resfolder + "03_w300.webp"}, "WEBP", 300, 210},
		{images.TestCase{testfolder + "03.webp?width=500&height=50&mode=fit", testfolder + "03_h50.webp", resfolder + "03_h50.webp"}, "WEBP", 71, 50},
		{images.TestCase{testfolder + "04.jpg?width=500&height=100&mode=fit", testfolder + "04_h100.jpg", resfolder + "04_h100.jpg"}, "JPEG", 143, 100},
		{images.TestCase{testfolder + "04.jpg?width=200&height=500&mode=fit", testfolder + "04_w200.jpg", resfolder + "04_w200.jpg"}, "JPEG", 200, 140},
		{images.TestCase{testfolder + "04.jpg?width=300&height=500&mode=fit", testfolder + "04_w300.jpg", resfolder + "04_w300.jpg"}, "JPEG", 300, 210},
		{images.TestCase{testfolder + "04.jpg?width=500&height=50&mode=fit", testfolder + "04_h50.jpg", resfolder + "04_h50.jpg"}, "JPEG", 71, 50},
		{images.TestCase{testfolder + "04.webp?width=500&height=100&mode=fit", testfolder + "04_h100.webp", resfolder + "04_h100.webp"}, "WEBP", 143, 100},
		{images.TestCase{testfolder + "04.webp?width=200&height=500&mode=fit", testfolder + "04_w200.webp", resfolder + "04_w200.webp"}, "WEBP", 200, 140},
		{images.TestCase{testfolder + "04.webp?width=300&height=500&mode=fit", testfolder + "04_w300.webp", resfolder + "04_w300.webp"}, "WEBP", 300, 210},
		{images.TestCase{testfolder + "04.webp?width=500&height=50&mode=fit", testfolder + "04_h50.webp", resfolder + "04_h50.webp"}, "WEBP", 71, 50},
	}

	err := testGetImages(cases, tolerance)
	if err != nil {
		t.Fatal(err)
	}
}

// Test GETting few liquid rescaled images
func TestGetImageSingleParam(t *testing.T) {
	testfolder := "testimages/server/"
	resfolder := "testresults/server/"
	tolerance := 0.002

	cases := []images.TestCaseAll{
		{images.TestCase{testfolder + "03.jpg?height=100", testfolder + "03_h100.jpg", resfolder + "03_h100.jpg"}, "JPEG", 143, 100},
		{images.TestCase{testfolder + "03.jpg?width=200", testfolder + "03_w200.jpg", resfolder + "03_w200.jpg"}, "JPEG", 200, 140},
		{images.TestCase{testfolder + "03.jpg?width=300", testfolder + "03_w300.jpg", resfolder + "03_w300.jpg"}, "JPEG", 300, 210},
		{images.TestCase{testfolder + "03.jpg?height=50", testfolder + "03_h50.jpg", resfolder + "03_h50.jpg"}, "JPEG", 71, 50},
		{images.TestCase{testfolder + "03.webp?height=100", testfolder + "03_h100.webp", resfolder + "03_h100.webp"}, "WEBP", 143, 100},
		{images.TestCase{testfolder + "03.webp?width=200", testfolder + "03_w200.webp", resfolder + "03_w200.webp"}, "WEBP", 200, 140},
		{images.TestCase{testfolder + "03.webp?width=300", testfolder + "03_w300.webp", resfolder + "03_w300.webp"}, "WEBP", 300, 210},
		{images.TestCase{testfolder + "03.webp?height=50", testfolder + "03_h50.webp", resfolder + "03_h50.webp"}, "WEBP", 71, 50},
		{images.TestCase{testfolder + "04.jpg?height=100", testfolder + "04_h100.jpg", resfolder + "04_h100.jpg"}, "JPEG", 143, 100},
		{images.TestCase{testfolder + "04.jpg?width=200", testfolder + "04_w200.jpg", resfolder + "04_w200.jpg"}, "JPEG", 200, 140},
		{images.TestCase{testfolder + "04.jpg?width=300", testfolder + "04_w300.jpg", resfolder + "04_w300.jpg"}, "JPEG", 300, 210},
		{images.TestCase{testfolder + "04.jpg?height=50", testfolder + "04_h50.jpg", resfolder + "04_h50.jpg"}, "JPEG", 71, 50},
		{images.TestCase{testfolder + "04.webp?height=100", testfolder + "04_h100.webp", resfolder + "04_h100.webp"}, "WEBP", 143, 100},
		{images.TestCase{testfolder + "04.webp?width=200", testfolder + "04_w200.webp", resfolder + "04_w200.webp"}, "WEBP", 200, 140},
		{images.TestCase{testfolder + "04.webp?width=300", testfolder + "04_w300.webp", resfolder + "04_w300.webp"}, "WEBP", 300, 210},
		{images.TestCase{testfolder + "04.webp?height=50", testfolder + "04_h50.webp", resfolder + "04_h50.webp"}, "WEBP", 71, 50},
	}

	err := testGetImages(cases, tolerance)
	if err != nil {
		t.Fatal(err)
	}
}

// Test GETting different sized and formats with mode=fit
func TestGetLiquid(t *testing.T) {
	testfolder := "testimages/server/"
	resfolder := "testresults/server/"
	tolerance := 0.06

	cases := []images.TestCaseAll{
		{images.TestCase{testfolder + "01.jpg?width=143&height=100&mode=liquid", testfolder + "liquid_01_143x100.jpg", resfolder + "liquid_01_143x100.jpg"}, "JPEG", 143, 100},
		{images.TestCase{testfolder + "02.jpg?width=200&height=140&mode=liquid", testfolder + "liquid_02_200x140.jpg", resfolder + "liquid_02_200x140.jpg"}, "JPEG", 200, 140},
		{images.TestCase{testfolder + "03.jpg?width=300&mode=liquid", testfolder + "liquid_03_w300.jpg", resfolder + "liquid_03_w300.jpg"}, "JPEG", 300, 210},
		{images.TestCase{testfolder + "03.jpg?height=300&mode=liquid", testfolder + "liquid_03_h300.jpg", resfolder + "liquid_03_h300.jpg"}, "JPEG", 429, 300},
	}

	err := testGetImages(cases, tolerance)
	if err != nil {
		t.Fatal(err)
	}
}

// Test for right content-type request header
func TestMIMEtype(t *testing.T) {
	s, ln, srverr := startServer()
	defer stopServer(s, ln, srverr)
	response := fasthttp.AcquireResponse()
	request := fasthttp.AcquireRequest()
	cases := []string{"01.jpg", "01.png",	"01.webp", "01.tiff", "01.bmp",	"01.gif"}
	folder := "testimages/server/"

	for _, c := range cases {
		request.SetRequestURI(fmt.Sprintf("http://localhost:%d/", port) + folder + c)
		fasthttp.Do(request, response)
		MIME := string(response.Header.ContentType())

		img := images.NewImage()
		img.FromBlob(response.Body())
		expected := "image/" + img.GetImageFormat()

		if MIME != expected {
			t.Fatal("Server returned: " + MIME)
		}
		request.Reset()
		response.Reset()
		img.Destroy()
	}
	fasthttp.ReleaseRequest(request)
	fasthttp.ReleaseResponse(response)
}
