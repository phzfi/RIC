package cache

import (
	"github.com/phzfi/RIC/server/images"
)

type DenyUpscale struct {
	Resizer
}

// TODO: Should we "cap" each dimension individually to at most original or use
// original image size when either of requested dimensions exceed originals?
// Currently handles each dimension individually.
func (d DenyUpscale) GetImage(filename string, xsize uint, ysize uint) (images.ImageBlob, error) {
	x, y, _ := d.ImageSize(filename)
	if x < xsize {
		xsize = x
	}
	if y < ysize {
		ysize = y
	}
	return d.Resizer.GetImage(filename, xsize, ysize)
}
