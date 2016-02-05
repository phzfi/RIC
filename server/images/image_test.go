package images

import (
	"fmt"
	"github.com/phzfi/RIC/server/logging"
	"testing"
)

func TestImageConvert(t *testing.T) {

	testfolder := "../testimages/convert/"
	testimage := testfolder + "toconvert.jpg"
	testimage2 := testfolder + "toconvert2.jpg"
	testimage3 := testfolder + "toconvert3.jpg"
	testimage4 := testfolder + "toconvert4.jpg"
	resfolder := "../testresults/images/"
	tolerance := 0.002

	cases := []FormatTestCase{
		{TestCase{testimage, testfolder + "converted.jpg", resfolder + "converted.jpg"}, "JPEG"},
		{TestCase{testimage, testfolder + "converted.webp", resfolder + "converted.webp"}, "WEBP"},
		{TestCase{testimage, testfolder + "converted.tiff", resfolder + "converted.tiff"}, "TIFF"},
		{TestCase{testimage, testfolder + "converted.gif", resfolder + "converted.gif"}, "GIF"},
		{TestCase{testimage, testfolder + "converted.png", resfolder + "converted.png"}, "PNG"},
		{TestCase{testimage, testfolder + "converted.bmp", resfolder + "converted.bmp"}, "BMP"},
		{TestCase{testimage2, testfolder + "converted2.jpg", resfolder + "converted2.jpg"}, "JPEG"},
		{TestCase{testimage2, testfolder + "converted2.webp", resfolder + "converted2.webp"}, "WEBP"},
		{TestCase{testimage3, testfolder + "converted3.jpg", resfolder + "converted3.jpg"}, "JPEG"},
		{TestCase{testimage3, testfolder + "converted3.webp", resfolder + "converted3.webp"}, "WEBP"},
		{TestCase{testimage4, testfolder + "converted4.jpg", resfolder + "converted4.jpg"}, "JPEG"},
		{TestCase{testimage4, testfolder + "converted4.webp", resfolder + "converted4.webp"}, "WEBP"},
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
	testimage2 := testfolder + "toresize2.jpg"
	testimage3 := testfolder + "toresize3.jpg"
	testimage4 := testfolder + "toresize4.jpg"
	resfolder := "../testresults/images/"
	tolerance := 0.002

	cases := []SizeTestCase{
		{TestCase{testimage, testfolder + "1_100x100.jpg", resfolder + "1_100x100.jpg"}, 100, 100},
		{TestCase{testimage, testfolder + "1_200x200.jpg", resfolder + "1_200x200.jpg"}, 200, 200},
		{TestCase{testimage, testfolder + "1_300x400.jpg", resfolder + "1_300x400.jpg"}, 300, 400},
		{TestCase{testimage, testfolder + "1_500x200.jpg", resfolder + "1_500x200.jpg"}, 500, 200},
		{TestCase{testimage, testfolder + "1_30x20.jpg", resfolder + "1_30x20.jpg"}, 30, 20},
		{TestCase{testimage, testfolder + "1_600x600.jpg", resfolder + "1_600x600.jpg"}, 600, 600},
		{TestCase{testimage2, testfolder + "2_100x100.jpg", resfolder + "2_100x100.jpg"}, 100, 100},
		{TestCase{testimage2, testfolder + "2_200x200.jpg", resfolder + "2_200x200.jpg"}, 200, 200},
		{TestCase{testimage3, testfolder + "3_100x100.jpg", resfolder + "3_100x100.jpg"}, 100, 100},
		{TestCase{testimage3, testfolder + "3_200x200.jpg", resfolder + "3_200x200.jpg"}, 200, 200},
		{TestCase{testimage4, testfolder + "4_100x100.jpg", resfolder + "4_100x100.jpg"}, 100, 100},
		{TestCase{testimage4, testfolder + "4_200x200.jpg", resfolder + "4_200x200.jpg"}, 200, 200},
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

func TestImageWatermark(t *testing.T) {

	testfolder := "../testimages/watermark/"
	testimage := testfolder + "towatermark.jpg"
	resfolder := "../testresults/images/"
	tolerance := 0.002

	wmimage := NewImage()
	defer wmimg.destroy
	err := img.FromFile(c.Testfn)
	if err != nil {
		t.Fatal(err)
	}

	horizontal := 0.0
	vertical := 0.0

	cases := []TestCase{
		{testimage, testfolder + "marked1.jpg", resfolder + "marked1.jpg"},
		{testimage, testfolder + "marked2.jpg", resfolder + "marked2.jpg"},
		{testimage, testfolder + "marked3.jpg", resfolder + "marked3.jpg"},
	}

	for _, c := range cases {

		logging.Debug(fmt.Sprintf("Testing watermark: %v, %v, %v", c.Testfn, c.Reffn, c.Resfn))

		img := NewImage()
		defer img.Destroy()
		err := img.FromFile(c.Testfn)
		if err != nil {
			t.Fatal(err)
		}

		err = img.Watermark(wmimg, horizontal, vertical)
		if err != nil {
			t.Fatal(err)
		}
		blob := img.Blob()
		horizontal = horizontal + 0.5
		vertical = vertical + 0.5
		err = FormatTest(c, blob, tolerance)
		if err != nil {
			t.Fatal(err)
		}
	}
}
