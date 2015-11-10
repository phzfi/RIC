package main

import (
	"errors"
	"github.com/joonazan/imagick/imagick"
	"path/filepath"
	"testing"
)

func TestResize1(t *testing.T) {
	doTest(t, "toresize.jpg", "resized.jpg", true)
}

func TestResize2(t *testing.T) {
	doTest(t, "toresize2.jpg", "resized2.jpg", true)
}

func TestResize3(t *testing.T) {
	doTest(t, "toresize.jpg", "resized_bad.jpg", false)
}

func doTest(t *testing.T, to_resize_fn, resized_fn string, should_pass bool) {
	const tolerance = 0.002

	distortion, err := GetDistortion(to_resize_fn, resized_fn)
	if err != nil {
		t.Fatal(err)
	}

	if distortion > tolerance == should_pass {
		t.Fatal("Resize failed. Distortion:", distortion, "Tolerance:", tolerance, "Should pass?:", should_pass)
	}
}

func GetDistortion(filename, filename_cmp string) (distortion float64, err error) {
	const image_folder = "testimages/resize/"
	imagick.Initialize()
	defer imagick.Terminate()

	mw_cmp := imagick.NewMagickWand()
	defer mw_cmp.Destroy()

	err = mw_cmp.ReadImage(filepath.FromSlash(image_folder + filename_cmp))
	if err != nil {
		err = errors.New("Could not load reference image:" + err.Error())
		return
	}

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	image, err := LoadImage(filepath.FromSlash(image_folder + filename))
	if err != nil {
		err = errors.New("LoadImage failed:" + err.Error())
		return
	}

	resized, err := image.Resized(100, 100)
	if err != nil {
		err = errors.New("Resize failed:" + err.Error())
		return
	}

	mw.ReadImageBlob(resized.ToBlob())

        trash, distortion := mw.CompareImages(mw_cmp, imagick.METRIC_MEAN_SQUARED_ERROR)
	trash.Destroy()

	err = mw.WriteImage(filepath.FromSlash("testresults/resize/" + filename))

	return
}
