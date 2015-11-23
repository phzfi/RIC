package cache

import "github.com/phzfi/RIC/server/images"

type ImageCache interface {
	GetImage(filename string, width *uint, height *uint) (images.ImageBlob, error)

	// TODO: These could be moved a separate filesystem handler
	// but these are ok, where they are at the moment (no need to bloat yet).
	AddRoot(string) error
	RemoveRoot(string) error
}
