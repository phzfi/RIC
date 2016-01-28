package main

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/ops"
	"fmt"
	"testing"
)

func TestOperatorConvert(t *testing.T) {

	operator, src := SetupOperatorSource()
	
	testfolder := "../testimages/convert/"
	testimage := testfolder + "toconvert"
	resfolder := "../testresults/convert/"
	tolerance := 0.002
	
	type testCase struct {
		testfn, reffn, resfn string
		format string
	}
	
	cases := [...]testCase {
		{testimage, testfolder + "converted.jpg",  resfolder + "converted.jpg",  "JPEG"},
		{testimage, testfolder + "converted.webp", resfolder + "converted.webp", "WEBP"},
		{testimage, testfolder + "converted.tiff", resfolder + "converted.tiff", "TIFF"},
		{testimage, testfolder + "converted.gif",  resfolder + "converted.gif",  "GIF"},
		{testimage, testfolder + "converted.png",  resfolder + "converted.png",  "PNG"},
		{testimage, testfolder + "converted.bmp",  resfolder + "converted.bmp",  "BMP"},
	}
	
	for _, c := range cases {
		
		logging.Debug(fmt.Sprintf("Testing convert: %v, %v, %v, %v", c.testfn, c.reffn, c.format, c.resfn))
		
		blob, err := operator.GetBlob(src.LoadImageOp(c.testfn), ops.Convert{c.format})
		if err != nil {
			t.Fatal(err)
		}

		img := images.NewImage()
		defer img.Destroy()
		img.FromBlob(blob)

		f := img.GetImageFormat()
		if f != c.format {
			t.Fatal(fmt.Sprintf("Bad image format. Requested %v, Got %v", c.format, f))
		}

		err = images.CheckDistortion(blob, c.reffn, tolerance, c.resfn)
		if err != nil {
			t.Fatal(err)
		}
	}
}


func TestOperatorResize(t *testing.T) {

	operator, src := SetupOperatorSource()
	
	testfolder := "../testimages/resize/"
	testimage := testfolder + "toresize"
	resfolder := "../testresults/operator/"
	tolerance := 0.002

	type testCase struct {
		testfn, reffn, resfn string
		w, h int
	}
	
	cases := [...]testCase {
		{testimage, testfolder + "100x100.jpg", resfolder + "100x100.jpg", 100, 100},
		{testimage, testfolder + "200x200.jpg", resfolder + "200x200.jpg", 200, 200},
		{testimage, testfolder + "300x400.jpg", resfolder + "300x400.jpg", 300, 400},
		{testimage, testfolder + "500x200.jpg", resfolder + "500x200.jpg", 500, 200},
		{testimage, testfolder + "30x20.jpg",   resfolder + "30x20.jpg",  30, 20},
		{testimage, testfolder + "600x600.jpg", resfolder + "600x600.jpg", 600, 600},
	}
	
	for _, c := range cases {

		logging.Debug(fmt.Sprintf("Testing resize: %v, %v, %v, %v, %v", c.testfn, c.reffn, c.w, c.h, c.resfn))
		
		blob, err := operator.GetBlob(src.LoadImageOp(c.testfn), ops.Resize{c.w, c.h})
		if err != nil {
			t.Fatal(err)
		}

		img := images.NewImage()
		defer img.Destroy()
		img.FromBlob(blob)
		
		w := img.GetWidth()
		h := img.GetHeight()
		if w != c.w || h != c.h {
			t.Fatal(fmt.Sprintf("Bad image size. Requested (%v, %v) , Got (%v, %v)", c.w, c.h, w, h))
		}

		err = images.CheckDistortion(blob, c.reffn, tolerance, c.resfn)
		if err != nil {
			t.Fatal(err)
		}
	}
}
