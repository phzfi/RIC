package ops

import (
  "github.com/phzfi/RIC/server/logging"
  "github.com/phzfi/RIC/server/configuration"
  "github.com/phzfi/RIC/server/images"
)

type Watermarker struct {
  watermarkImage images.Image
}

func WatermarkOp() Operation {
  return watermark{}
}

func MakeWatermarker() Watermarker {
  image := images.NewImage()
  err := image.FromFile(configuration.GetString("watermark", "path"))
  if err != nil {
    logging.Debug("Error loading watermark image." + err.Error())
  }
  return Watermarker {
    watermarkImage: image,
  }
}
