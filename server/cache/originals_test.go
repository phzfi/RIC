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

func TestCachlessOriginals(t *testing.T) {
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
	id := "original"
	filename := id + ".jpg"

	err = c.AddRoot(path)
	if err != nil {
		return
	}
	blob, err := c.GetOriginalSizedImage(id)
	if err != nil {
		return
	}

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err = mw.ReadImageBlob(blob)
	if err != nil {
		return
	}
	err = mw.WriteImage(filepath.FromSlash(path + filename))
	if err != nil {
		return
	}
	if !strings.EqualFold(mw.GetImageFormat(), "jpeg") {
		err = errors.New("Cache returned wrong blob format. Original is jpeg, got: " + mw.GetImageFormat())
		return
	}

	err = CompareBlobToImage(blob, path+id+".jpg")
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
