package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

// Operation for resizing image to Width and Height with liquid rescale method
type LiquidRescale struct {
	Width, Height int
}

func (r LiquidRescale) Marshal() string {
	return string(liquidRescaleID) + int32ToString(uint32(r.Width)) + int32ToString(uint32(r.Height))
}

func (r LiquidRescale) Apply(img images.Image) error {
	logging.Debugf("Liquid rescaling image to: %v, %v", r.Width, r.Height)
	// The third argument to LiquidRescaleImage is the maximum transversal step, or how many pixels a seam can move sideways.
	// The fourth is the rigidity, which makes seams less steep. Recommended if transversal step is greater than one.
	return img.LiquidRescaleImage(uint(r.Width), uint(r.Height), 1, 0)
}
