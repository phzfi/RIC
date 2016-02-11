package images

import (
	"github.com/gographics/imagick/imagick"
)

func (img Image) Resize(w, h int) error {
	return img.ResizeImage(uint(w), uint(h), imagick.FILTER_LANCZOS, 1)
}
