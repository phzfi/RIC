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
  return Watermarker {
    watermarkImage: image.LoadImage(config.GetString("watermark", "path"))
  }
}
