package cache

import (
	"strings"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"fmt"
)

type ScaleMode struct {
	ScaleImage func(c ImageCache, filename string, width, height *uint) (img images.ImageBlob, err error)
}

var SCALEMODE_RESIZE  = ScaleMode{resize}
var SCALEMODE_FIT     = ScaleMode{fit}
//var SCALEMODE_FILL  = ScaleMode{fill}
// var SCALEMODE_CROP  = ScaleMode{crop}
var SCALEMODE_DEFAULT = SCALEMODE_RESIZE

func (m *ScaleMode) FromString(s string) {
	s = strings.ToLower(s)
	switch s {
	case "resize":
		*m = SCALEMODE_RESIZE
		return
	case "fit":
		*m = SCALEMODE_FIT
		return
	}
	*m = SCALEMODE_DEFAULT
	return
}

func resize(c ImageCache, filename string, width, height *uint) (blob images.ImageBlob, err error){
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
