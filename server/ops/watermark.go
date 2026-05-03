package ops

import (
	"math"
	"strconv"

	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type watermark struct {
	stamp      images.Image
	text       string
	horizontal float64
	vertical   float64
	margin     float64
	fontSize   float64
	color      string
}

func (w watermark) Marshal() string {
	if w.text != "" {
		return strconv.Itoa(int(watermarkID)) + w.text + float64ToString(w.vertical) + float64ToString(w.horizontal)
	}
	return strconv.Itoa(int(watermarkID)) + float64ToString(w.vertical) + float64ToString(w.horizontal)
}

func (w watermark) Apply(img images.Image) (err error) {
	logging.Debug("Adding watermark")
	if w.text != "" {
		return img.WatermarkText(w.text, w.fontSize, w.color, w.horizontal, w.vertical, w.margin)
	}
	return img.Watermark(w.stamp, w.horizontal, w.vertical)
}

func float64ToString(x float64) string {
	return int64ToString(math.Float64bits(x))
}

func int64ToString(x uint64) string {
	return int32ToString(uint32(x>>32)) + int32ToString(uint32(x))
}
