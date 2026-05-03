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

func TextWatermarkOp(text string, hor, ver, margin float64, fontSize float64, color string) Operation {
	return watermark{
		text:       text,
		horizontal: hor,
		vertical:   ver,
		margin:     margin,
		fontSize:   fontSize,
		color:      color,
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
