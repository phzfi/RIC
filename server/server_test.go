package main

import (
	"bytes"
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/testutils"
	"github.com/valyala/fasthttp"
	"strings"
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

	cases := []testutils.TestCaseAll{
		{testutils.TestCase{testfolder + "01.jpg?width=100&height=100", testfolder + "01_100x100.jpg", resfolder + "01_100x100.jpg"}, "JPEG", 100, 100},
		{testutils.TestCase{testfolder + "01.jpg?width=200&height=100", testfolder + "01_200x100.jpg", resfolder + "01_200x100.jpg"}, "JPEG", 200, 100},
		{testutils.TestCase{testfolder + "01.jpg?width=300&height=100", testfolder + "01_300x100.jpg", resfolder + "01_300x100.jpg"}, "JPEG", 300, 100},
		{testutils.TestCase{testfolder + "01.jpg?width=300&height=50", testfolder + "01_300x50.jpg", resfolder + "01_300x50.jpg"}, "JPEG", 300, 50},
		{testutils.TestCase{testfolder + "01.jpg?format=webp&width=100&height=100", testfolder + "01_100x100.webp", resfolder + "01_100x100.webp"}, "WEBP", 100, 100},
		{testutils.TestCase{testfolder + "01.jpg?format=webp&width=200&height=100", testfolder + "01_200x100.webp", resfolder + "01_200x100.webp"}, "WEBP", 200, 100},
		{testutils.TestCase{testfolder + "01.jpg?format=webp&width=300&height=100", testfolder + "01_300x100.webp", resfolder + "01_300x100.webp"}, "WEBP", 300, 100},
		{testutils.TestCase{testfolder + "01.jpg?format=webp&width=300&height=50", testfolder + "01_300x50.webp", resfolder + "01_300x50.webp"}, "WEBP", 300, 50},
		{testutils.TestCase{testfolder + "02.jpg?width=100&height=100", testfolder + "02_100x100.jpg", resfolder + "02_100x100.jpg"}, "JPEG", 100, 100},
		{testutils.TestCase{testfolder + "02.jpg?width=200&height=100", testfolder + "02_200x100.jpg", resfolder + "02_200x100.jpg"}, "JPEG", 200, 100},
		{testutils.TestCase{testfolder + "02.jpg?width=300&height=100", testfolder + "02_300x100.jpg", resfolder + "02_300x100.jpg"}, "JPEG", 300, 100},
		{testutils.TestCase{testfolder + "02.jpg?width=300&height=50", testfolder + "02_300x50.jpg", resfolder + "02_300x50.jpg"}, "JPEG", 300, 50},
		{testutils.TestCase{testfolder + "02.jpg?format=webp&width=100&height=100", testfolder + "02_100x100.webp", resfolder + "02_100x100.webp"}, "WEBP", 100, 100},
		{testutils.TestCase{testfolder + "02.jpg?format=webp&width=200&height=100", testfolder + "02_200x100.webp", resfolder + "02_200x100.webp"}, "WEBP", 200, 100},
		{testutils.TestCase{testfolder + "02.jpg?format=webp&width=300&height=100", testfolder + "02_300x100.webp", resfolder + "02_300x100.webp"}, "WEBP", 300, 100},
		{testutils.TestCase{testfolder + "02.jpg?format=webp&width=300&height=50", testfolder + "02_300x50.webp", resfolder + "02_300x50.webp"}, "WEBP", 300, 50},
	}

	err := testGetImages(cases)
	if err != nil {
		t.Fatal(err)
	}
}

// Test GETting different sized and formats with mode=fit
func TestGetImageFit(t *testing.T) {

	testfolder := "testimages/server/"
	resfolder := "testresults/server/"

	cases := []testutils.TestCaseAll{
		{testutils.TestCase{testfolder + "03.jpg?width=500&height=100&mode=fit", testfolder + "03_h100.jpg", resfolder + "03_h100.jpg"}, "JPEG", 143, 100},
		{testutils.TestCase{testfolder + "03.jpg?width=200&height=500&mode=fit", testfolder + "03_w200.jpg", resfolder + "03_w200.jpg"}, "JPEG", 200, 140},
		{testutils.TestCase{testfolder + "03.jpg?width=300&height=500&mode=fit", testfolder + "03_w300.jpg", resfolder + "03_w300.jpg"}, "JPEG", 300, 210},
		{testutils.TestCase{testfolder + "03.jpg?width=500&height=50&mode=fit", testfolder + "03_h50.jpg", resfolder + "03_h50.jpg"}, "JPEG", 71, 50},
		{testutils.TestCase{testfolder + "03.jpg?format=webp&width=500&height=100&mode=fit", testfolder + "03_h100.webp", resfolder + "03_h100.webp"}, "WEBP", 143, 100},
		{testutils.TestCase{testfolder + "03.jpg?format=webp&width=200&height=500&mode=fit", testfolder + "03_w200.webp", resfolder + "03_w200.webp"}, "WEBP", 200, 140},
		{testutils.TestCase{testfolder + "03.jpg?format=webp&width=300&height=500&mode=fit", testfolder + "03_w300.webp", resfolder + "03_w300.webp"}, "WEBP", 300, 210},
		{testutils.TestCase{testfolder + "03.jpg?format=webp&width=500&height=50&mode=fit", testfolder + "03_h50.webp", resfolder + "03_h50.webp"}, "WEBP", 71, 50},
		{testutils.TestCase{testfolder + "04.jpg?width=500&height=100&mode=fit", testfolder + "04_h100.jpg", resfolder + "04_h100.jpg"}, "JPEG", 143, 100},
		{testutils.TestCase{testfolder + "04.jpg?width=200&height=500&mode=fit", testfolder + "04_w200.jpg", resfolder + "04_w200.jpg"}, "JPEG", 200, 140},
		{testutils.TestCase{testfolder + "04.jpg?width=300&height=500&mode=fit", testfolder + "04_w300.jpg", resfolder + "04_w300.jpg"}, "JPEG", 300, 210},
		{testutils.TestCase{testfolder + "04.jpg?width=500&height=50&mode=fit", testfolder + "04_h50.jpg", resfolder + "04_h50.jpg"}, "JPEG", 71, 50},
		{testutils.TestCase{testfolder + "04.jpg?format=webp&width=500&height=100&mode=fit", testfolder + "04_h100.webp", resfolder + "04_h100.webp"}, "WEBP", 143, 100},
		{testutils.TestCase{testfolder + "04.jpg?format=webp&width=200&height=500&mode=fit", testfolder + "04_w200.webp", resfolder + "04_w200.webp"}, "WEBP", 200, 140},
		{testutils.TestCase{testfolder + "04.jpg?format=webp&width=300&height=500&mode=fit", testfolder + "04_w300.webp", resfolder + "04_w300.webp"}, "WEBP", 300, 210},
		{testutils.TestCase{testfolder + "04.jpg?format=webp&width=500&height=50&mode=fit", testfolder + "04_h50.webp", resfolder + "04_h50.webp"}, "WEBP", 71, 50},
	}

	err := testGetImages(cases)
	if err != nil {
		t.Fatal(err)
	}
}

// Test GETting few liquid rescaled images
func TestGetImageSingleParam(t *testing.T) {
	testfolder := "testimages/server/"
	resfolder := "testresults/server/"

	cases := []testutils.TestCaseAll{
		{testutils.TestCase{testfolder + "03.jpg?height=100", testfolder + "03_h100.jpg", resfolder + "03_h100.jpg"}, "JPEG", 143, 100},
		{testutils.TestCase{testfolder + "03.jpg?width=200", testfolder + "03_w200.jpg", resfolder + "03_w200.jpg"}, "JPEG", 200, 140},
		{testutils.TestCase{testfolder + "03.jpg?width=300", testfolder + "03_w300.jpg", resfolder + "03_w300.jpg"}, "JPEG", 300, 210},
		{testutils.TestCase{testfolder + "03.jpg?height=50", testfolder + "03_h50.jpg", resfolder + "03_h50.jpg"}, "JPEG", 71, 50},
		{testutils.TestCase{testfolder + "03.jpg?format=webp&height=100", testfolder + "03_h100.webp", resfolder + "03_h100.webp"}, "WEBP", 143, 100},
		{testutils.TestCase{testfolder + "03.jpg?format=webp&width=200", testfolder + "03_w200.webp", resfolder + "03_w200.webp"}, "WEBP", 200, 140},
		{testutils.TestCase{testfolder + "03.jpg?format=webp&width=300", testfolder + "03_w300.webp", resfolder + "03_w300.webp"}, "WEBP", 300, 210},
		{testutils.TestCase{testfolder + "03.jpg?format=webp&height=50", testfolder + "03_h50.webp", resfolder + "03_h50.webp"}, "WEBP", 71, 50},
		{testutils.TestCase{testfolder + "04.jpg?height=100", testfolder + "04_h100.jpg", resfolder + "04_h100.jpg"}, "JPEG", 143, 100},
		{testutils.TestCase{testfolder + "04.jpg?width=200", testfolder + "04_w200.jpg", resfolder + "04_w200.jpg"}, "JPEG", 200, 140},
		{testutils.TestCase{testfolder + "04.jpg?width=300", testfolder + "04_w300.jpg", resfolder + "04_w300.jpg"}, "JPEG", 300, 210},
		{testutils.TestCase{testfolder + "04.jpg?height=50", testfolder + "04_h50.jpg", resfolder + "04_h50.jpg"}, "JPEG", 71, 50},
		{testutils.TestCase{testfolder + "04.jpg?format=webp&height=100", testfolder + "04_h100.webp", resfolder + "04_h100.webp"}, "WEBP", 143, 100},
		{testutils.TestCase{testfolder + "04.jpg?format=webp&width=200", testfolder + "04_w200.webp", resfolder + "04_w200.webp"}, "WEBP", 200, 140},
		{testutils.TestCase{testfolder + "04.jpg?format=webp&width=300", testfolder + "04_w300.webp", resfolder + "04_w300.webp"}, "WEBP", 300, 210},
		{testutils.TestCase{testfolder + "04.jpg?format=webp&height=50", testfolder + "04_h50.webp", resfolder + "04_h50.webp"}, "WEBP", 71, 50},
	}

	err := testGetImages(cases)
	if err != nil {
		t.Fatal(err)
	}
}

// Test GETting different sized and formats with mode=liquid
func TestGetLiquid(t *testing.T) {
	testfolder := "testimages/server/"
	resfolder := "testresults/server/"

	cases := []testutils.TestCaseAll{
		{testutils.TestCase{testfolder + "01.jpg?width=143&height=100&mode=liquid", testfolder + "liquid_01_143x100.jpg", resfolder + "liquid_01_143x100.jpg"}, "JPEG", 143, 100},
		{testutils.TestCase{testfolder + "02.jpg?width=200&height=140&mode=liquid", testfolder + "liquid_02_200x140.jpg", resfolder + "liquid_02_200x140.jpg"}, "JPEG", 200, 140},
		{testutils.TestCase{testfolder + "03.jpg?width=300&mode=liquid", testfolder + "liquid_03_w300.jpg", resfolder + "liquid_03_w300.jpg"}, "JPEG", 300, 210},
		{testutils.TestCase{testfolder + "03.jpg?height=300&mode=liquid", testfolder + "liquid_03_h300.jpg", resfolder + "liquid_03_h300.jpg"}, "JPEG", 429, 300},
	}

	err := testGetImages(cases)
	if err != nil {
		t.Fatal(err)
	}
}

// Test GETting different sized images with mode=crop and mode=cropmid
func TestGetCrop(t *testing.T) {
	testfolder := "testimages/server/"
	resfolder := "testresults/server/"

	cases := []testutils.TestCaseAll{
		{testutils.TestCase{testfolder + "02.jpg?width=200&height=100&mode=crop", testfolder + "crop_200x100.jpg", resfolder + "crop_200x100.jpg"}, "JPEG", 200, 100},
		{testutils.TestCase{testfolder + "02.jpg?width=200&height=200&mode=crop&cropx=100", testfolder + "crop_200x200_offset100x0y.jpg", resfolder + "crop_200x200_offset100x0y.jpg"}, "JPEG", 200, 200},
		{testutils.TestCase{testfolder + "02.jpg?width=200&height=200&mode=crop&cropx=100&cropy=100", testfolder + "crop_200x200_offset100x100y.jpg", resfolder + "crop_200x200_offset100x100y.jpg"}, "JPEG", 200, 200},
		{testutils.TestCase{testfolder + "02.jpg?width=300&mode=crop", testfolder + "crop_w300.jpg", resfolder + "crop_w300.jpg"}, "JPEG", 300, 900},
		{testutils.TestCase{testfolder + "02.jpg?height=300&mode=crop", testfolder + "crop_h300.jpg", resfolder + "crop_h300.jpg"}, "JPEG", 1200, 300},
		{testutils.TestCase{testfolder + "02.jpg?width=200&height=200&mode=cropmid", testfolder + "cropmid_200x200.jpg", resfolder + "cropmid_200x200.jpg"}, "JPEG", 200, 200},
		{testutils.TestCase{testfolder + "02.jpg?width=300&mode=cropmid", testfolder + "cropmid_w300.jpg", resfolder + "cropmid_w300.jpg"}, "JPEG", 300, 900},
		{testutils.TestCase{testfolder + "02.jpg?height=300&mode=cropmid", testfolder + "cropmid_h300.jpg", resfolder + "cropmid_h300.jpg"}, "JPEG", 1200, 300},
	}

	err := testGetImages(cases)
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
	cases := []string{"", "?format=jpeg", "?format=png", "?format=webp", "?format=tiff", "?format=bmp", "?format=gif"}
	folder := "testimages/server/"

	for _, c := range cases {
		request.SetRequestURI(fmt.Sprintf("http://localhost:%d/", port) + folder + "01.jpg" + c)
		fasthttp.Do(request, response)
		MIME := string(response.Header.ContentType())

		img := images.NewImage()
		img.FromBlob(response.Body())
		expected := "image/" + strings.ToLower(img.GetImageFormat())

		if MIME != expected {
			t.Fatal(fmt.Sprintf("Server returned: %s, image is %s", MIME, expected))
		}
		request.Reset()
		response.Reset()
		img.Destroy()
	}
	fasthttp.ReleaseRequest(request)
	fasthttp.ReleaseResponse(response)
}

func TestInvalidParams(t *testing.T) {
	s, ln, srverr := startServer()
	defer stopServer(s, ln, srverr)
	response := fasthttp.AcquireResponse()
	request := fasthttp.AcquireRequest()
	base := fmt.Sprintf("http://localhost:%d/", port) + "testimages/server/01.jpg"
	cases := []string{
		"?width=abc&height=200",
		"?width=200&height=abc",
		"?width=200?height=200",
		"?width=200&height=200&mode=fit&something=x",
		"?width=200&height=200&mode=y",
		"?format=abc"}

	for _, c := range cases {
		uri := base + c
		request.SetRequestURI(uri)
		fasthttp.Do(request, response)
		if response.Header.StatusCode() != 400 {
			t.Fatal("Server did not return 400 to request: " + uri)
		}
		request.Reset()
		response.Reset()
	}
	fasthttp.ReleaseRequest(request)
	fasthttp.ReleaseResponse(response)
}
