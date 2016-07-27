package main

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/testutils"
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
		{testimage, testfolder + "marked1.jpg", resfolder + "marked1.jpg"},
		{testimage, testfolder + "marked2.jpg", resfolder + "marked2.jpg"},
		{testimage, testfolder + "marked3.jpg", resfolder + "marked3.jpg"},
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
