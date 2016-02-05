package ops

import (
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type watermark struct {
}

func (w Watermark) Apply(img images.Image) error {
	logging.Debug("Adding watermark")
	return img.Watermark()
}
