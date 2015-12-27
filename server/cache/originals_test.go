package cache

import (
	"errors"
	"fmt"
	"github.com/joonazan/imagick/imagick"
	"github.com/phzfi/RIC/server/images"
	"path/filepath"
	"testing"
)

func TestCachelessOriginals(t *testing.T) {
	err := testCacheOriginals(NewCacheless())
	if err != nil {
		t.Fatal(err)
	}
}

func TestCacheOriginals(t *testing.T) {
	err := testCacheOriginals(New(&DummyPolicy{}, 512*1024*1024))
	if err != nil {
		t.Fatal(err)
	}
}

func testCacheOriginals(c ImageCache) (err error) {

	path := filepath.FromSlash("../testimages/originals/")
	filename := "original.jpg"

	err = c.AddRoot(path)
	if err != nil {
		return
	}
	blob, err := c.GetOriginalSizedImage(filename)
	if err != nil {
		return
	}

	err = CompareBlobToImage(blob, path+filename)
	return
}

func CompareBlobToImage(blob []byte, filename string) (err error) {
	img_cmp, err := images.LoadImage(filepath.FromSlash(filename))

	if err != nil {
		return errors.New("LoadImage failed: " + err.Error())
	}

	img := imagick.NewMagickWand()
	img.ReadImageBlob(blob)
	trash, distortion := img.CompareImages(img_cmp.MagickWand, imagick.METRIC_MEAN_SQUARED_ERROR)
	trash.Destroy()

	const tolerance = 0.000001
	if distortion > tolerance {
		return errors.New(fmt.Sprintf("Image load failed. Distortion: %f, Tolerance: %f", distortion, tolerance))
	}

	return nil
}
