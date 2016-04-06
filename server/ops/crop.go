package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type Crop struct {
	Width, Height, X, Y int
}

func (c Crop) Marshal() string {
	return string(cropID) + int32ToString(uint32(c.Width)) + int32ToString(uint32(c.Height)) +
          int32ToString(uint32(c.X)) + int32ToString(uint32(c.Y))
}

func (c Crop) Apply(img images.Image) error {
	logging.Debugf("Crop image to: %d, %d with offset: %d, %d", c.Width, c.Height, c.X, c.Y)
	return img.Crop(c.Width, c.Height, c.X, c.Y)
}
