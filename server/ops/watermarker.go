package ops

import (
  "github.com/phzfi/RIC/server/images"
)

type Watermarker struct {
  WatermarkImage images.Image
  Horizontal float64
  Vertical float64
  Maxwidth int
  Minwidth int
  Maxheight int
  Minheight int
  AddMark bool
}

func WatermarkOp(stamp images.Image, hor, ver float64) Operation {
  return watermark{
    stamp: stamp,
    horizontal: hor,
    vertical: ver,
  }
}

func MakeWatermarker(path string, hor, ver float64, maxwidth, minwidth, maxheight, minheight int, addmark bool) (wm Watermarker, err error) {
  image := images.NewImage()
  err = image.FromFile(path)
  if err != nil {
    return
  }
  wm = Watermarker {
    WatermarkImage: image,
    Horizontal: hor,
    Vertical: ver,
    Maxwidth: maxwidth,
    Minwidth: minwidth,
    Maxheight: maxheight,
    Minheight: minheight,
    AddMark: addmark,
  }
  return
}
