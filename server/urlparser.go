package main

import (
	"github.com/phzfi/RIC/server/ops"
	"github.com/phzfi/RIC/server/config"
	"github.com/phzfi/RIC/server/logging"
	"github.com/valyala/fasthttp"
	"path/filepath"
	"strings"
)

func ExtToFormat(ext string) string {
	ext = strings.ToUpper(strings.TrimLeft(ext, "."))
	if ext == "JPG" {
		return "JPEG"
	}
	return ext
}

func ParseURI(uri *fasthttp.URI, source ops.ImageSource, marker ops.Watermarker, conf config.Conf) (operations []ops.Operation, err error) {
	args := uri.QueryArgs()
	filename := string(uri.Path())
	w, werr := args.GetUint("width")
	h, herr := args.GetUint("height")
	ow, oh, err := source.ImageSize(filename)
	if err != nil {
		return
	}
	mode := string(args.Peek("mode"))

	operations = []ops.Operation{source.LoadImageOp(filename)}

	adjustWidth := func() {
		w = int(float32(h*ow)/float32(oh) + 0.5)
	}

	adjustHeight := func() {
		h = int(float32(w*oh)/float32(ow) + 0.5)
	}

	adjustSize := func() {
		if herr != nil && werr == nil {
			adjustHeight()
		} else if herr == nil && werr != nil {
			adjustWidth()
		} else if werr != nil && herr != nil {
			w, h = ow, oh
		}
	}

	denyUpscale := func() {
		f := float32(w) / float32(h)
		if w > ow {
			w = ow
			h = int(float32(w) / f + 0.5)
		}
		if h > oh {
			h = oh
			w = int(f * float32(h) + 0.5)
		}
	}

	resize := func() {
		adjustSize()
		denyUpscale()
		operations = append(operations, ops.Resize{w, h})
	}

	liquid := func() {
		adjustSize()
		denyUpscale()
		operations = append(operations, ops.LiquidRescale{w, h})
	}

	fit := func() {
		if werr == nil && herr == nil {
			denyUpscale()
			if ow*h > w*oh {
				adjustHeight()
			} else {
				adjustWidth()
			}
			operations = append(operations, ops.Resize{w, h})
		} else {
			resize()
		}
	}

	watermark := func() {
		heightOK := h > marker.Minheight && h < marker.Maxheight
		widthOK := w > marker.Minwidth && w < marker.Maxwidth
		if marker.AddMark && heightOK && widthOK {
			logging.Debug("Adding watermarkOp")
			operations = append(operations, ops.WatermarkOp(marker.WatermarkImage, marker.Horizontal, marker.Vertical))
		}
	}

	switch mode {
	case "resize":
		resize()
	case "fit":
		fit()
	case "liquid":
		liquid()
	default:
		resize()
	}
	watermark()

	ext := filepath.Ext(filename)
	if ext != "" {
		operations = append(operations, ops.Convert{ExtToFormat(ext)})
	}

	return
}
