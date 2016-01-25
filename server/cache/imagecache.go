package cache

import (
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type ImageCache interface {
	GetImage(filename string, width, height uint) (images.ImageBlob, error)

	GetOriginalSizedImage(filename string) (images.ImageBlob, error)
	GetImageByWidth(filename string, width uint) (images.ImageBlob, error)
	GetImageByHeight(filename string, height uint) (images.ImageBlob, error)
	GetImageFit(filename string, width uint, height uint) (images.ImageBlob, error)

	// TODO: These could be moved a separate filesystem handler
	// but these are ok, where they are at the moment (no need to bloat yet).
	AddRoot(string) error
	RemoveRoot(string) error
}

type AmbiguousSizeImageCache struct {
	ImageCache
}

func (a AmbiguousSizeImageCache) GetImage(filename string, width, height *uint, mode string) (blob images.ImageBlob, err error) {
	logging.Debug(fmt.Sprintf("Image request: filename=%v, width=%v, height=%v, mode=%v", filename, width, height, mode))
	scalefunc := ScaleFunc(mode)
	blob, err = scalefunc(a.ImageCache, filename, width, height)
	return
}
