package main

import (
	"errors"
	"github.com/phzfi/RIC/server/config"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/ops"
	"github.com/valyala/fasthttp"
	"strings"
)

func ParseURI(uri *fasthttp.URI, source ops.ImageSource, marker ops.Watermarker, conf config.Conf) (operations []ops.Operation, format string, err, invalid error) {
	filename := string(uri.Path())
	w, h, mode, format, invalid := getParams(uri.QueryArgs())
	ow, oh, err := source.ImageSize(filename)
	if invalid != nil {
		return
	}
	if err != nil {
		return
	}

	operations = []ops.Operation{source.LoadImageOp(filename)}

	adjustWidth := func() {
		w = roundedIntegerDivision(h*ow, oh)
	}

	adjustHeight := func() {
		h = roundedIntegerDivision(w*oh, ow)
	}

	adjustSize := func() {
		if h == -1 && w != -1 {
			adjustHeight()
		} else if h != -1 && w == -1 {
			adjustWidth()
		} else if w == -1 && h == -1 {
			w, h = ow, oh
		}
	}

	denyUpscale := func() {
		h0 := h
		w0 := w
		if w > ow {
			h = roundedIntegerDivision(ow*h0, w0)
			w = ow
		}
		if h > oh || h > h0 {
			w = roundedIntegerDivision(oh*w0, h0)
			h = oh
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
		if w > ow {
			w = ow
		}
		if h > oh {
			h = oh
		}
		if w != -1 && h != -1 {
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
	case resizeMode:
		resize()
	case fitMode:
		fit()
	case liquidMode:
		liquid()
	}
	watermark()

	operations = append(operations, ops.Convert{format})

	return
}

func roundedIntegerDivision(n, m int) int {
	if (n < 0) == (m < 0) {
		return (n + m/2) / m
	} else { // -5 / 6 should round to -1
		return (n - m/2) / m
	}
}

var stringToMode = map[string]mode{
	"":       resizeMode,
	"resize": resizeMode,
	"fit":    fitMode,
	"crop":   cropMode,
	"liquid": liquidMode,
}

type mode int

const (
	fitMode = mode(1 + iota)
	cropMode
	liquidMode
	resizeMode

	widthParam  = "width"
	heightParam = "height"
	modeParam   = "mode"
	formatParam = "format"
)

// returns validated parameters from request and error if invalid
func getParams(a *fasthttp.Args) (w int, h int, mode mode, format string, err error) {
	if strings.Contains(a.String(), "%") {
		err = errors.New("Invalid characters in request!")
		return
	}

	w, err = a.GetUint(widthParam)
	if isParseError(err) {
		return
	}

	h, err = a.GetUint(heightParam)
	if isParseError(err) {
		return
	}

	mode = stringToMode[string(a.Peek(modeParam))]
	if mode == 0 {
		err = errors.New("Invalid mode!")
		return
	}

	format = strings.ToLower(string(a.Peek(formatParam)))
	if format == "" {
		format = "jpeg"
	}
	// TODO: verify that the format is one we support.
	// We do not want to support TXT, for instance

	a.Del(widthParam)
	a.Del(heightParam)
	a.Del(modeParam)
	a.Del(formatParam)
	if a.Len() != 0 {
		err = errors.New("Invalid parameter " + a.String())
		return
	}

	err = nil
	return
}

func isParseError(err error) bool {
	return err != nil && err != fasthttp.ErrNoArgValue
}
