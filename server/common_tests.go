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
	
	cases := [...]images.FormatTestCase {
		{images.TestCase{testimage, testfolder + "converted.jpg",  resfolder + "converted.jpg"},  "JPEG"},
		{images.TestCase{testimage, testfolder + "converted.webp", resfolder + "converted.webp"}, "WEBP"},
		{images.TestCase{testimage, testfolder + "converted.tiff", resfolder + "converted.tiff"}, "TIFF"},
		{images.TestCase{testimage, testfolder + "converted.gif",  resfolder + "converted.gif"},  "GIF"},
		{images.TestCase{testimage, testfolder + "converted.png",  resfolder + "converted.png"},  "PNG"},
		{images.TestCase{testimage, testfolder + "converted.bmp",  resfolder + "converted.bmp"},  "BMP"},
	}
	
	for _, c := range cases {
		logging.Debug(fmt.Sprintf("Testing convert: %v, %v, %v, %v", c.Testfn, c.Reffn, c.Format, c.Resfn))

		blob, err := operator.GetBlob(src.LoadImageOp(c.Testfn), ops.Convert{c.Format})
		if err != nil {
			t.Fatal(err)
		}
		
		err = images.FormatTest(c, blob, tolerance)
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
	
	cases := [...]images.SizeTestCase {
		{images.TestCase{testimage, testfolder + "100x100.jpg", resfolder + "100x100.jpg"}, 100, 100},
		{images.TestCase{testimage, testfolder + "200x200.jpg", resfolder + "200x200.jpg"}, 200, 200},
		{images.TestCase{testimage, testfolder + "300x400.jpg", resfolder + "300x400.jpg"}, 300, 400},
		{images.TestCase{testimage, testfolder + "500x200.jpg", resfolder + "500x200.jpg"}, 500, 200},
		{images.TestCase{testimage, testfolder + "30x20.jpg",   resfolder + "30x20.jpg"},   30, 20},
		{images.TestCase{testimage, testfolder + "600x600.jpg", resfolder + "600x600.jpg"}, 600, 600},
	}
	
	for _, c := range cases {

		logging.Debug(fmt.Sprintf("Testing resize: %v, %v, %v, %v, %v", c.Testfn, c.Reffn, c.W, c.H, c.Resfn))
		
		blob, err := operator.GetBlob(src.LoadImageOp(c.Testfn), ops.Resize{c.W, c.H})
		if err != nil {
			t.Fatal(err)
		}
		
		err = images.SizeTest(c, blob, tolerance)
		if err != nil {
			t.Fatal(err)
		}
	}
}
