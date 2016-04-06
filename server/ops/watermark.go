package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"math"
)

type watermark struct {
	stamp      images.Image
	horizontal float64
	vertical   float64
}

func (w watermark) Marshal() string {
	return string(watermarkID) + float64ToString(w.vertical) + float64ToString(w.horizontal)
}

func (w watermark) Apply(img images.Image) (err error) {
	logging.Debug("Adding watermark")
	return img.Watermark(w.stamp, w.horizontal, w.vertical)
}

func float64ToString(x float64) string {
	return int64ToString(math.Float64bits(x))
}

func int64ToString(x uint64) string {
	return int32ToString(uint32(x>>32)) + int32ToString(uint32(x))
}
