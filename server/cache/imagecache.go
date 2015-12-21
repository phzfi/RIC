package cache

import (
	"github.com/phzfi/RIC/server/images"
)

type ImageCache interface {
	GetImage(filename string, width, height uint) (images.ImageBlob, error)
	GetOriginal(filename string) (images.ImageBlob, error)

	// TODO: These could be moved a separate filesystem handler
	// but these are ok, where they are at the moment (no need to bloat yet).
	AddRoot(string) error
	RemoveRoot(string) error
}

type AmbiguousSizeImageCache struct {
	ImageCache
}

func (a AmbiguousSizeImageCache) GetImage(filename string, width, height *uint) (blob images.ImageBlob, err error) {

	if width == nil || height == nil {
		blob, err = a.ImageCache.GetOriginal(filename)
		return
	}

	blob, err = a.ImageCache.GetImage(filename, *width, *height)
	return
}
