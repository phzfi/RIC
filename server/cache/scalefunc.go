package cache

import (
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"strings"
)

var SCALEFUNC_RESIZE = resize
var SCALEFUNC_FIT = fit

//var SCALEFUNC_FILL  = fill
// var SCALEFUNC_CROP  = crop
var SCALEFUNC_DEFAULT = SCALEFUNC_RESIZE

// Get ScaleFunction according to given mode string (eg. "resize", "fit", "crop", ...). If given string is nil, default scale function is returned.
func ScaleFunc(s *string) (f func(c ImageCache, filename string, width, height *uint) (blob images.ImageBlob, err error)) {
	if s == nil {
		f = SCALEFUNC_DEFAULT
		return
	}
	t := strings.ToLower(*s)
	switch t {
	case "resize":
		f = SCALEFUNC_RESIZE
		return
	case "fit":
		f = SCALEFUNC_FIT
		return
	}
	f = SCALEFUNC_DEFAULT
	return
}

// ScaleFunction for resize mode
func resize(c ImageCache, filename string, width, height *uint) (blob images.ImageBlob, err error) {
	logging.Debug(fmt.Sprintf("Scale mode: resize: filename=%v, width=%v, height=%v", filename, width, height))
	if width == nil && height != nil {
		blob, err = c.GetImageByHeight(filename, *height)
		return
	}
	if width != nil && height == nil {
		blob, err = c.GetImageByWidth(filename, *width)
		return
	}
	if width == nil && height == nil {
		blob, err = c.GetOriginalSizedImage(filename)
		return
	}
	blob, err = c.GetImage(filename, *width, *height)
	return
}

// ScaleFunction for fit mode
func fit(c ImageCache, filename string, width, height *uint) (blob images.ImageBlob, err error) {
	logging.Debug(fmt.Sprintf("Scale mode: fit: filename=%v, width=%v, height=%v", filename, width, height))
	if width == nil && height != nil {
		blob, err = c.GetImageByHeight(filename, *height)
		return
	}
	if width != nil && height == nil {
		blob, err = c.GetImageByWidth(filename, *width)
		return
	}
	if width == nil && height == nil {
		blob, err = c.GetOriginalSizedImage(filename)
		return
	}
	blob, err = c.GetImageFit(filename, *width, *height)
	return
}
