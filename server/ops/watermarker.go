package ops

import (
  "github.com/phzfi/RIC/server/logging"
  "github.com/phzfi/RIC/server/images"
)

type Watermarker struct {
  WatermarkImage images.Image
}

func WatermarkOp(stamp images.Image, hor, ver float64) Operation {
  return watermark{
    stamp: stamp,
    horizontal: hor,
    vertical: ver,
  }
}

func MakeWatermarker(path string) Watermarker {
  image := images.NewImage()
  err := image.FromFile(path)
  if err != nil {
    logging.Debug("Error loading watermark image." + err.Error())
  }
  return Watermarker {
    WatermarkImage: image,
  }
}
