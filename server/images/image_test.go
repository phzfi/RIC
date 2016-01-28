package images

import (
	"github.com/phzfi/RIC/server/logging"
	"fmt"
	"testing"
)

func TestImageConvert(t *testing.T) {
	
	testfolder := "../testimages/convert/"
	testimage := testfolder + "toconvert.jpg"
	resfolder := "../testresults/convert/"
	tolerance := 0.002
	
	cases := [...]FormatTestCase {
		{TestCase{testimage, testfolder + "converted.jpg",  resfolder + "converted.jpg"},  "JPEG"},
		{TestCase{testimage, testfolder + "converted.webp", resfolder + "converted.webp"}, "WEBP"},
		{TestCase{testimage, testfolder + "converted.tiff", resfolder + "converted.tiff"}, "TIFF"},
		// {TestCase{testimage, testfolder + "converted.gif",  resfolder + "converted.gif"},  "GIF"}, This test takes too much time
		{TestCase{testimage, testfolder + "converted.png",  resfolder + "converted.png"},  "PNG"},
		{TestCase{testimage, testfolder + "converted.bmp",  resfolder + "converted.bmp"},  "BMP"},
	}
	
	for _, c := range cases {
		
		logging.Debug(fmt.Sprintf("Testing convert: %v, %v, %v, %v", c.Testfn, c.Reffn, c.Format, c.Resfn))
		
		img := NewImage()
		defer img.Destroy()
		err := img.FromFile(c.Testfn)
		if err != nil {
			t.Fatal(err)
		}

		err = img.Convert(c.Format)
		if err != nil {
			t.Fatal(err)
		}
		blob := img.Blob()
		
		err = FormatTest(c, blob, tolerance)
		if err != nil {
			t.Fatal(err)
		}
	}
}


func TestImageResize(t *testing.T) {
	
	testfolder := "../testimages/resize/"
	testimage := testfolder + "toresize.jpg"
	resfolder := "../testresults/operator/"
	tolerance := 0.002

	cases := [...]SizeTestCase {
		{TestCase{testimage, testfolder + "100x100.jpg", resfolder + "100x100.jpg"}, 100, 100},
		{TestCase{testimage, testfolder + "200x200.jpg", resfolder + "200x200.jpg"}, 200, 200},
		{TestCase{testimage, testfolder + "300x400.jpg", resfolder + "300x400.jpg"}, 300, 400},
		{TestCase{testimage, testfolder + "500x200.jpg", resfolder + "500x200.jpg"}, 500, 200},
		{TestCase{testimage, testfolder + "30x20.jpg",   resfolder + "30x20.jpg"},  30, 20},
		{TestCase{testimage, testfolder + "600x600.jpg", resfolder + "600x600.jpg"}, 600, 600},
	}
	
	for _, c := range cases {

		logging.Debug(fmt.Sprintf("Testing resize: %v, %v, %v, %v, %v", c.Testfn, c.Reffn, c.W, c.H, c.Resfn))
		
		img := NewImage()
		defer img.Destroy()
		err := img.FromFile(c.Testfn)
		if err != nil {
			t.Fatal(err)
		}
		
		err = img.Resize(c.W, c.H)
		if err != nil {
			t.Fatal(err)
		}
		blob := img.Blob()

		err = SizeTest(c, blob, tolerance)
		if err != nil {
			t.Fatal(err)
		}
	}
}
