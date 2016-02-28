package ops

import (
  "github.com/phzfi/RIC/server/logging"
  "github.com/phzfi/RIC/server/configuration"
  "github.com/phzfi/RIC/server/images"
)

type Watermarker struct {
  WatermarkImage images.Image
}

func WatermarkOp(stamp images.Image) Operation {
  return watermark{
    stamp: stamp,
  }
}

func MakeWatermarker() Watermarker {
  image := images.NewImage()
  err := image.FromFile(configuration.GetString("watermark", "path"))
  if err != nil {
    logging.Debug("Error loading watermark image." + err.Error())
  }
  return Watermarker {
    WatermarkImage: image,
  }
}
