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
		
		img := NewImage()
		defer img.Destroy()
		err := img.FromFile(c.testfn)
		if err != nil {
			t.Fatal(err)
		}

		err = img.Convert(c.format)
		if err != nil {
			t.Fatal(err)
		}
		blob := img.Blob()

		img = NewImage()
		defer img.Destroy()
		img.FromBlob(blob)

		f := img.GetImageFormat()
		if f != c.format {
			t.Fatal(fmt.Sprintf("Bad image format. Requested %v, Got %v", c.format, f))
		}

		err = CheckDistortion(blob, c.reffn, tolerance, c.resfn)
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
		
		img := NewImage()
		defer img.Destroy()
		err := img.FromFile(c.testfn)
		if err != nil {
			t.Fatal(err)
		}
		
		err = img.Resize(c.w, c.h)
		if err != nil {
			t.Fatal(err)
		}
		blob := img.Blob()

		img = NewImage()
		defer img.Destroy()
		img.FromBlob(blob)
		
		w := img.GetWidth()
		h := img.GetHeight()
		if w != c.w || h != c.h {
			t.Fatal(fmt.Sprintf("Bad image size. Requested (%v, %v) , Got (%v, %v)", c.w, c.h, w, h))
		}

		err = CheckDistortion(blob, c.reffn, tolerance, c.resfn)
		if err != nil {
			t.Fatal(err)
		}
	}
}
