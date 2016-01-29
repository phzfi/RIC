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
	
	testfolder := "testimages/convert/"
	testimage := testfolder + "toconvert"
	testimage2 := testfolder + "toconvert2"
	testimage3 := testfolder + "toconvert3"
	testimage4 := testfolder + "toconvert4"
	resfolder := "testresults/common/"
	tolerance := 0.002
	
	cases := []images.FormatTestCase {
		{images.TestCase{testimage, testfolder + "converted.jpg",  resfolder + "converted.jpg"},  "JPEG"},
		{images.TestCase{testimage, testfolder + "converted.webp", resfolder + "converted.webp"}, "WEBP"},
		{images.TestCase{testimage, testfolder + "converted.tiff", resfolder + "converted.tiff"}, "TIFF"},
		{images.TestCase{testimage, testfolder + "converted.gif",  resfolder + "converted.gif"},  "GIF"},
		{images.TestCase{testimage, testfolder + "converted.png",  resfolder + "converted.png"},  "PNG"},
		{images.TestCase{testimage, testfolder + "converted.bmp",  resfolder + "converted.bmp"},  "BMP"},
		{images.TestCase{testimage2, testfolder + "converted2.jpg",  resfolder + "converted2.jpg"},  "JPEG"},
		{images.TestCase{testimage2, testfolder + "converted2.webp", resfolder + "converted2.webp"}, "WEBP"},
		{images.TestCase{testimage3, testfolder + "converted3.jpg",  resfolder + "converted3.jpg"},  "JPEG"},
		{images.TestCase{testimage3, testfolder + "converted3.webp", resfolder + "converted3.webp"}, "WEBP"},
		{images.TestCase{testimage4, testfolder + "converted4.jpg",  resfolder + "converted4.jpg"},  "JPEG"},
		{images.TestCase{testimage4, testfolder + "converted4.webp", resfolder + "converted4.webp"}, "WEBP"},
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
	
	testfolder := "testimages/resize/"
	testimage := testfolder + "toresize"
	testimage2 := testfolder + "toresize2"
	testimage3 := testfolder + "toresize3"
	testimage4 := testfolder + "toresize4"
	resfolder := "testresults/common/"
	tolerance := 0.002
	
	cases := []images.SizeTestCase {
		{images.TestCase{testimage, testfolder + "1_100x100.jpg", resfolder + "1_100x100.jpg"}, 100, 100},
		{images.TestCase{testimage, testfolder + "1_200x200.jpg", resfolder + "1_200x200.jpg"}, 200, 200},
		{images.TestCase{testimage, testfolder + "1_300x400.jpg", resfolder + "1_300x400.jpg"}, 300, 400},
		{images.TestCase{testimage, testfolder + "1_500x200.jpg", resfolder + "1_500x200.jpg"}, 500, 200},
		{images.TestCase{testimage, testfolder + "1_30x20.jpg",   resfolder + "1_30x20.jpg"},   30, 20},
		{images.TestCase{testimage, testfolder + "1_600x600.jpg", resfolder + "1_600x600.jpg"}, 600, 600},
		{images.TestCase{testimage2, testfolder + "2_100x100.jpg", resfolder + "2_100x100.jpg"}, 100, 100},
		{images.TestCase{testimage2, testfolder + "2_200x200.jpg", resfolder + "2_200x200.jpg"}, 200, 200},
		{images.TestCase{testimage3, testfolder + "3_100x100.jpg", resfolder + "3_100x100.jpg"}, 100, 100},
		{images.TestCase{testimage3, testfolder + "3_200x200.jpg", resfolder + "3_200x200.jpg"}, 200, 200},
		{images.TestCase{testimage4, testfolder + "4_100x100.jpg", resfolder + "4_100x100.jpg"}, 100, 100},
		{images.TestCase{testimage4, testfolder + "4_200x200.jpg", resfolder + "4_200x200.jpg"}, 200, 200},
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


func TestOperatorLiquidRescale(t *testing.T) {

	operator, src := SetupOperatorSource()
	
	testfolder := "testimages/resize/"
	testimage := testfolder + "toresize"
	testimage2 := testfolder + "toresize2"
	resfolder := "testresults/common/"
	tolerance := 0.002
	
	cases := []images.SizeTestCase {
		{images.TestCase{testimage, testfolder + "liquid1_100x100.jpg", resfolder + "liquid1_100x100.jpg"}, 100, 100},
		{images.TestCase{testimage, testfolder + "liquid1_500x200.jpg", resfolder + "liquid1_500x200.jpg"}, 500, 200},
		{images.TestCase{testimage2, testfolder + "liquid2_200x200.jpg", resfolder + "liquid2_200x200.jpg"}, 200, 200},
	}
	
	for _, c := range cases {

		logging.Debug(fmt.Sprintf("Testing resize: %v, %v, %v, %v, %v", c.Testfn, c.Reffn, c.W, c.H, c.Resfn))
		
		blob, err := operator.GetBlob(src.LoadImageOp(c.Testfn), ops.LiquidRescale{c.W, c.H})
		if err != nil {
			t.Fatal(err)
		}
		
		err = images.SizeTest(c, blob, tolerance)
		if err != nil {
			t.Fatal(err)
		}
	}
}
