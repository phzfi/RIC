package ops

import (
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type watermark struct {
	img Image
}

func (w Watermark) Apply(img images.Image) error {
	logging.Debug("Adding watermark")
	horizontal, err := config.GetFloat64("watermark", "horizontal")
	vertical, err := config.GetFloat64("watermark", "vertical")

	if err != nil {
		logging.Debug("Error loading config alignment." + err.Error())
		return
	}
	return img.Watermark(img, horizontal, vertical)
}
