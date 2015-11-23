package images

import (
	"github.com/joonazan/imagick/imagick"
)

func (img Image) Resized(w, h uint) (resized Image, err error) {

	resized = img.Clone()

	err = resized.ResizeImage(w, h, imagick.FILTER_LANCZOS, 1)

	return
}
