package cache

import (
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type AspectPreserver struct {
	Resizer
}

func New(policy Policy, mm uint64) *AspectPreserver {
	return &AspectPreserver{Resizer: NewCache(NewBasicResizer(), policy, mm)}
}

func NewCacheless() *AspectPreserver {
	return &AspectPreserver{Resizer: NewBasicResizer()}
}

// Gets image blob by width only. Image is scaled and aspect ratio preserved.
func (a *AspectPreserver) GetImageByWidth(filename string, width uint) (blob images.ImageBlob, err error) {
	logging.Debug(fmt.Sprintf("Get image by width: %v, %v", filename, width))

	originalWidth, originalHeight, err := a.ImageSize(filename)
	if err != nil {
		return
	}

	height := width * originalHeight / originalWidth
	blob, err = a.GetImage(filename, width, height)
	return

}

// Gets image blob by height only. Image is scaled and aspect ratio preserved.
func (a *AspectPreserver) GetImageByHeight(filename string, height uint) (blob images.ImageBlob, err error) {
	logging.Debug(fmt.Sprintf("Get image by height: %v, %v", filename, height))

	originalWidth, originalHeight, err := a.ImageSize(filename)
	if err != nil {
		return
	}

	width := height * originalWidth / originalHeight
	blob, err = a.GetImage(filename, width, height)
	return
}

// Get image fitted to w, h preserving aspect ratio
func (a *AspectPreserver) GetImageFit(filename string, width uint, height uint) (blob images.ImageBlob, err error) {
	logging.Debug(fmt.Sprintf("Get image fitted to: %v %v %v", filename, width, height))

	originalWidth, originalHeight, err := a.ImageSize(filename)
	if err != nil {
		return
	}

	aspect := width / height
	originalAspect := originalWidth / originalHeight
	
	if aspect <= originalAspect {
		blob, err = a.GetImageByWidth(filename, width)
	} else {
		blob, err = a.GetImageByHeight(filename, height)
	}

	return
}

// Gets original sized image blob.
func (a *AspectPreserver) GetOriginalSizedImage(filename string) (blob images.ImageBlob, err error) {
	logging.Debug(fmt.Sprintf("Get original sized image: %v", filename))

	width, height, err := a.ImageSize(filename)
	if err != nil {
		return
	}

	blob, err = a.GetImage(filename, width, height)
	return
}
