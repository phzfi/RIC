package cache

import (
	"testing"
	"path/filepath"
	"github.com/joonazan/imagick/imagick"
	"errors"
	"strings"
)


func TestCacheGetPNG(t *testing.T) {
	err := testCacheFormat("png")
	if err != nil {
		t.Fatal(err)
	}
}


func TestCacheGetJPG(t *testing.T) {
	err := testCacheFormat("jpg")
	if err != nil {
		t.Fatal(err)
	}
}


func TestCacheGetWEBP(t *testing.T) {
	err := testCacheFormat("webp")
	if err != nil {
		t.Fatal(err)
	}
}


func TestCacheGetTIFF(t *testing.T) {
	err := testCacheFormat("tiff")
	if err != nil {
		t.Fatal(err)
	}
}


func TestCacheGetGIF(t *testing.T) {
	err := testCacheFormat("gif")
	if err != nil {
		t.Fatal(err)
	}
}


func testCacheFormat(ext string) (err error) {
	var cache Cacheless
	path := filepath.FromSlash("../testimages/formats")
	err = cache.AddRoot(path)
	if err != nil {
		return
	}
	filename := "formats." + ext
	blob, err := cache.GetImage(filename, 480, 240)
	if err != nil {
		return
	}
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	mw.ReadImageBlob(blob)
	if !strings.EqualFold(mw.GetImageFormat(), ext) && !(strings.EqualFold(mw.GetImageFormat(), "JPEG") && strings.EqualFold(ext, "jpg")) {
		err = errors.New("Cache returned wrong blob format. Requested: " + ext + ", Got: " + mw.GetImageFormat())
		return
	}
	err = mw.WriteImage(filepath.FromSlash("../testresults/formats/" + filename))
	if err != nil {
		return
	}
	return
}
