package ops

import (
	"github.com/phzfi/RIC/server/config"
	"github.com/phzfi/RIC/server/images"
)

type Watermarker struct {
	WatermarkImage images.Image
	config.Watermark
}

func WatermarkOp(stamp images.Image, hor, ver float64) Operation {
	return watermark{
		stamp:      stamp,
		horizontal: hor,
		vertical:   ver,
	}
}

func MakeWatermarker(settings config.Watermark) (wm Watermarker, err error) {
	image := images.NewImage()
	err = image.FromFile(settings.ImagePath)
	if err != nil {
		return
	}
	wm = Watermarker{
		WatermarkImage: image,
		Watermark:      settings,
	}
	return
}
