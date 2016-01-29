package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type LiquidRescale struct {
	Width, Height int
}

func (r LiquidRescale) Apply(img images.Image) error {
	logging.Debugf("Liquid rescaling image to: %v, %v", r.Width, r.Height)
	// The third argument to LiquidRescaleImage is the maximum transversal step, or how many pixels a seam can move sideways.
	// The fourth is the rigidity, which makes seams less steep. Recommended if transversal step is greater than one.
	return img.LiquidRescaleImage(uint(r.Width), uint(r.Height), 1, 0)
}
