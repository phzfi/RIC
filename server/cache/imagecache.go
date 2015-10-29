package cache

import "server/images"

type ImageCache interface {
	GetImage(filename string, width *int, height *int) (images.ImageBlob, error)
}
