package cache

import (
	"errors"
	"fmt"
	"github.com/joonazan/imagick/imagick"
	"github.com/phzfi/RIC/server/images"
	"path/filepath"
	"strings"
	"testing"
)

const imagespath = "../testimages/aspect/"
const resultspath = "../testresults/aspect/"
const tolerance = 0.0005

func TestCachlessAspect(t *testing.T) {
	err := testCacheAspect(NewCacheless())
	if err != nil {
		t.Fatal(err)
	}
}

func TestCacheAspect(t *testing.T) {
	err := testCacheAspect(New(&DummyPolicy{}, 512*1024*1024))
	if err != nil {
		t.Fatal(err)
	}
}

func testCacheAspect(c ImageCache) (err error) {
	err = c.AddRoot(imagespath)
	if err != nil {
		return
	}
	err = testByWidth(c, "aspect.jpg", 200, "bywidth.jpg")
	if err != nil {
		return
	}
	err = testByWidth(c, "aspect.webp", 200, "bywidth.jpg")
	if err != nil {
		return
	}
	err = testByHeight(c, "aspect.jpg", 200, "byheight.jpg")
	if err != nil {
		return
	}
	err = testByHeight(c, "aspect.webp", 200, "byheight.jpg")
	return
}

func testByWidth(c ImageCache, filename string, width uint, reffile string) (err error) {
	blob, err := c.GetImageByWidth(filename, width)
	if err != nil {
		return
	}
	img, err := images.ImageFromBlob(blob)
	defer img.Destroy()
	if err != nil {
		return
	}
	err = img.WriteImage(filepath.Join(resultspath, filename))
	if err != nil {
		return
	}
	err = testFormat(img, filename)
	if err != nil {
		return
	}
	err = compareImage(img, reffile)
	return
}

func testByHeight(c ImageCache, filename string, height uint, reffile string) (err error) {
	blob, err := c.GetImageByHeight(filename, height)
	if err != nil {
		return
	}
	img, err := images.ImageFromBlob(blob)
	defer img.Destroy()
	if err != nil {
		return
	}
	err = img.WriteImage(filepath.Join(resultspath, filename))
	if err != nil {
		return
	}
	err = testFormat(img, filename)
	if err != nil {
		return
	}
	err = compareImage(img, reffile)
	return
}

func testFormat(img images.Image, filename string) (err error) {
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	if !strings.EqualFold(img.GetImageFormat(), ext) && !(strings.EqualFold(img.GetImageFormat(), "JPEG") && strings.EqualFold(ext, "jpg")) {
		err = errors.New("Cache returned wrong blob format. Requested: " + ext + ", Got: " + img.GetImageFormat())
	}
	return
}

func compareImage(img images.Image, filename string) (err error) {
	img_cmp, err := images.LoadImage(filepath.Join(imagespath, filename))
	if err != nil {
		return errors.New("LoadImage failed: " + err.Error())
	}
	defer img_cmp.Destroy()
	trash, distortion := img.CompareImages(img_cmp.MagickWand, imagick.METRIC_MEAN_SQUARED_ERROR)
	trash.Destroy()
	if distortion > tolerance {
		err = errors.New(fmt.Sprintf("Cache returned bad image. Distortion: %f, Tolerance: %f", distortion, tolerance))
	}
	return
}
