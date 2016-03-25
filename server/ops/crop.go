package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type Crop struct {
	Width, Height, X, Y int
  Mid bool
}

func (c Crop) Marshal() string {
	return string(cropID) + int32ToString(uint32(c.Width)) + int32ToString(uint32(c.Height)) +
          int32ToString(uint32(c.X)) + int32ToString(uint32(c.Y))
}

func (c Crop) Apply(img images.Image) error {
  if c.X < 0 {
    c.X = 0
  }
  if c.Y < 0 {
    c.Y = 0
  }
  if c.Width < 0 {
    c.Width = img.GetWidth()
  }
  if c.Height < 0 {
    c.Height = img.GetHeight()
  }
  if c.Mid {
    logging.Debugf("Crop midimage to: %d, %d", c.Width, c.Height)
    midW := roundedIntegerDivision(img.GetWidth(), 2)
    midH := roundedIntegerDivision(img.GetWidth(), 2)
    c.X = midW - roundedIntegerDivision(c.Width, 2)
    c.Y = midH - roundedIntegerDivision(c.Height, 2)
  	return img.Crop(c.Width, c.Height, c.X, c.Y)
  }
	logging.Debugf("Crop image to: %d, %d with offset: %d, %d", c.Width, c.Height, c.X, c.Y)
	return img.Crop(c.Width, c.Height, c.X, c.Y)
}

func roundedIntegerDivision(n, m int) int {
	if (n < 0) == (m < 0) {
		return (n + m/2) / m
	} else { // -5 / 6 should round to -1
		return (n - m/2) / m
	}
}
