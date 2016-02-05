package ops

import (
  "github.com/phzfi/RIC/server/logging"
  "github.com/phzfi/RIC/server/config"
  "github.com/phzfi/RIC/server/images"
)

struct Watermarker{
  watermarkImage Image
}

func WatermarkOp() Operation {
  return watermark{}
}

func MakeWatermarker() Watermarker{
  image, err := loadimage.FromFile(config.GetString("watermark", "path"))
  if err != nil {
    logging.Debug("Error loading watermark image." + err.Error())
  }
  return Watermarker {
    watermarkImage: image
  }
}
