package main

import (
	"bytes"
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/testutils"
	"github.com/valyala/fasthttp"
	"testing"
)

func TestImageWatermark(t *testing.T) {

	testfolder := "testimages/watermark/"
	testimage := testfolder + "towatermark.jpg"
	resfolder := "testresults/images/"
	tolerance := 0.002

	wmimage := images.NewImage()
	defer wmimage.Destroy()
	err := wmimage.FromFile(testfolder + "watermark.png")
	if err != nil {
		t.Fatal(err)
	}

	horizontal := 0.0
	vertical := 0.0

	cases := []testutils.TestCase{
		{Testfn: testimage, Reffn: testfolder + "marked1.jpg", Resfn: resfolder + "marked1.jpg"},
		{Testfn: testimage, Reffn: testfolder + "marked2.jpg", Resfn: resfolder + "marked2.jpg"},
		{Testfn: testimage, Reffn: testfolder + "marked3.jpg", Resfn: resfolder + "marked3.jpg"},
	}

	for _, c := range cases {

		img := images.NewImage()
		defer img.Destroy()
		err := img.FromFile(c.Testfn)
		if err != nil {
			t.Fatal(err)
		}

		err = img.Watermark(wmimage, horizontal, vertical)
		if err != nil {
			t.Fatal(err)
		}
		blob := img.Blob()
		horizontal = horizontal + 0.5
		vertical = vertical + 0.5
		err = testutils.CheckDistortion(blob, c.Reffn, tolerance, c.Resfn)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestTextWatermarkServer(t *testing.T) {
	s, ln, srverr := startServer()
	defer stopServer(s, ln, srverr)

	response := fasthttp.AcquireResponse()
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.SetRequestURI(fmt.Sprintf("http://localhost:%d/testimages/server/01.jpg?width=500&height=500&mode=liquid&watermark=PHZ.fi", port))
	fasthttp.Do(request, response)

	if response.Header.StatusCode() != 200 {
		t.Fatalf("Expected 200, got %d", response.Header.StatusCode())
	}

	body := response.Body()
	if len(body) == 0 {
		t.Fatal("Response body should not be empty")
	}

	img := images.NewImage()
	defer img.Destroy()
	err := img.FromBlob(body)
	if err != nil {
		t.Fatalf("Failed to parse response as image: %v", err)
	}

	if img.GetWidth() != 500 || img.GetHeight() != 500 {
		t.Errorf("Image dimensions mismatch: got %dx%d, want 500x500", img.GetWidth(), img.GetHeight())
	}
}

func TestWatermarkDisabledParam(t *testing.T) {
	s, ln, srverr := startServer()
	defer stopServer(s, ln, srverr)

	response := fasthttp.AcquireResponse()
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.SetRequestURI(fmt.Sprintf("http://localhost:%d/testimages/server/01.jpg?width=500&height=500&mode=liquid&watermark=false", port))
	fasthttp.Do(request, response)

	if response.Header.StatusCode() != 200 {
		t.Fatalf("Expected 200, got %d", response.Header.StatusCode())
	}

	bodyWithFalse := response.Body()

	request.Reset()
	response.Reset()
	request.SetRequestURI(fmt.Sprintf("http://localhost:%d/testimages/server/01.jpg?width=500&height=500&mode=liquid", port))
	fasthttp.Do(request, response)

	if response.Header.StatusCode() != 200 {
		t.Fatalf("Expected 200, got %d", response.Header.StatusCode())
	}

	bodyWithout := response.Body()

	if !bytes.Equal(bodyWithFalse, bodyWithout) {
		t.Fatal("watermark=false should produce same result as no watermark param (both disabled)")
	}
}
