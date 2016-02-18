package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/configuration"
)

type watermark struct {
	img images.Image
}

func (w watermark) Apply(img images.Image) (err error) {
	logging.Debug("Adding watermark")
	horizontal, err := configuration.GetFloat64("watermark", "horizontal")
	vertical, err := configuration.GetFloat64("watermark", "vertical")

	if err != nil {
		logging.Debug("Error loading config alignment." + err.Error())
		return
	}
	return img.Watermark(img, horizontal, vertical)
}
