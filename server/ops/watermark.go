package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type watermark struct {
	stamp images.Image
	horizontal float64
	vertical float64
}

func (w watermark) Apply(img images.Image) (err error) {
	logging.Debug("Adding watermark")
	return img.Watermark(w.stamp, w.horizontal, w.vertical)
}
