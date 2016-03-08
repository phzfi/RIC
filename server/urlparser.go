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
		if w > ow {
			w = ow
			adjustHeight()
		}
		if h > oh {
			h = oh
			adjustWidth()
		}
	}

	resize := func() {
		denyUpscale()
		adjustSize()
		operations = append(operations, ops.Resize{w, h})
	}

	liquid := func() {
		denyUpscale()
		adjustSize()
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
		logging.Debug("Reading size restrictions")
		minHeight, err := conf.GetInt("watermark", "minheight")
		minWidth, err := conf.GetInt("watermark", "minwidth")
		maxHeight, err := conf.GetInt("watermark", "maxheight")
		maxWidth, err := conf.GetInt("watermark", "maxwidth")
		addMark, err := conf.GetBool("watermark", "addmark")

		if err != nil {
			logging.Debug("Error reading config size restrictions." + err.Error())
			return
		}

		ver, err := conf.GetFloat64("watermark", "vertical")
		hor, err := conf.GetFloat64("watermark", "horizontal")

		if err != nil {
			logging.Debug("Error loading config alignment." + err.Error())
			return
		}

		heightOK := h > minHeight && h < maxHeight
		widthOK := w > minWidth && w < maxWidth
		if addMark && heightOK && widthOK {
			logging.Debug("Adding watermarkOp")
			operations = append(operations, ops.WatermarkOp(marker.WatermarkImage, hor, ver))
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
