package images

import (
	"github.com/fubla/imagick/imagick"
)

func (img Image) Resize(w, h int) error {
	
    if isOpenCL{
        return img.AccelerateResizeImage(uint(w), uint(h),
        imagick.FILTER_LANCZOS)
    } else {
        return img.ResizeImage(uint(w), uint(h), imagick.FILTER_LANCZOS, 1)
    }
}


