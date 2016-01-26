package main

import (
	"errors"
	"fmt"
	"github.com/joonazan/imagick/imagick"
	"github.com/phzfi/RIC/server/cache"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
	"path/filepath"
	"testing"
)

const (
	testfolder    = "testimages/"
	resultsfolder = "testresults/"
	testgroup     = "resize/"
)

const TOLERANCE = 0.002

func TestResize(t *testing.T) {
	operator := cache.MakeOperator(512 * 1024 * 1024)
	src := ops.ImageSource{}
	src.AddRoot(testfolder + testgroup)
	blob, err := operator.GetBlob(
		src.LoadImageOp("toresize.jpg"),
		ops.Resize{100, 100},
	)
	if err != nil {
		return
	}
	d, err := getDistortion(blob, "resized.jpg")
	if err != nil {
		return
	}
	if d > TOLERANCE {
		t.Fatal(fmt.Sprintf("Bad image returned. Distortion: %v, Tolerance: %v", d, TOLERANCE))
	}
}

func getDistortion(blob images.ImageBlob, filename_cmp string) (distortion float64, err error) {
	const image_folder = testfolder + testgroup

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
	err = mw.ReadImageBlob(blob)
	if err != nil {
		return
	}

	// Save image, just in case someone wants to look at it
	err = mw.WriteImage(filepath.FromSlash(resultsfolder + testgroup + filename_cmp))
	if err != nil {
		return
	}

	trash, distortion := mw.CompareImages(mw_cmp, imagick.METRIC_MEAN_SQUARED_ERROR)
	trash.Destroy()

	return
}
