package main

import (
	"github.com/joonazan/imagick/imagick"
)

func (img ImageBlob) Resized(w, h uint) (resized ImageBlob, err error) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	mw.ReadImageBlob(img)
	err = mw.ResizeImage(w, h, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		return
	}

	resized = mw.GetImageBlob()
	return

}
