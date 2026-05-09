package main

import (
	"errors"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/ops"
	"github.com/valyala/fasthttp"
	"strings"
)

func ParseURI(uri *fasthttp.URI, source ops.ImageSource, marker ops.Watermarker) (operations []ops.Operation, format string, watermark string, err, invalid error) {
	filename := string(uri.Path())

	w, h, cropx, cropy, mode, format, url, watermark, invalid := getParams(uri.QueryArgs())
	if invalid != nil {
		return
	}

	if url != "" {
		source.AddRoot(url)
	}

	ow, oh, err := source.ImageSize(filename)
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
		if h == 0 && w != 0 {
			adjustHeight()
		} else if h != 0 && w == 0 {
			adjustWidth()
		} else if w == 0 && h == 0 {
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
		operations = append(operations, ops.Resize{Width: w, Height: h})
	}

	liquid := func() {
		denyUpscale()
		adjustSize()
		operations = append(operations, ops.LiquidRescale{Width: w, Height: h})
	}

	crop := func() {
		if w == 0 {
			w = ow
		}
		if h == 0 {
			h = oh
		}
		operations = append(operations, ops.Crop{Width: w, Height: h, X: cropx, Y: cropy})
	}

	cropmid := func() {
		if w == 0 || w > ow {
			w = ow
		}
		if h == 0 || h > oh {
			h = oh
		}
		midW := roundedIntegerDivision(ow, 2)
		midH := roundedIntegerDivision(oh, 2)
		cropx := midW - roundedIntegerDivision(w, 2)
		cropy := midH - roundedIntegerDivision(h, 2)
		operations = append(operations, ops.Crop{Width: w, Height: h, X: cropx, Y: cropy})
	}

	fit := func() {
		if w > ow {
			w = ow
		}
		if h > oh {
			h = oh
		}
		if w != 0 && h != 0 {
			if ow*h > w*oh {
				adjustHeight()
			} else {
				adjustWidth()
			}
			operations = append(operations, ops.Resize{Width: w, Height: h})
		} else {
			resize()
		}
	}

	addWatermark := func() {
		effectiveWatermark := watermark
		shouldAddMark := false

		if marker.ForceMark != "" {
			if marker.ForceMark == "false" || marker.ForceMark == "0" {
				shouldAddMark = false
			} else {
				shouldAddMark = true
				effectiveWatermark = marker.ForceMark
			}
		} else if effectiveWatermark == "false" || effectiveWatermark == "0" {
			shouldAddMark = false
		} else if effectiveWatermark != "" {
			shouldAddMark = true
		} else {
			shouldAddMark = marker.AddMark
			if marker.Text != "" {
				effectiveWatermark = marker.Text
			}
		}

		heightOK := h > marker.MinHeight && h < marker.MaxHeight
		widthOK := w > marker.MinWidth && w < marker.MaxWidth
		if shouldAddMark && heightOK && widthOK {
			logging.Debug("Adding watermarkOp")
			if effectiveWatermark == "" || effectiveWatermark == "true" || effectiveWatermark == "1" {
				operations = append(operations, ops.WatermarkOp(marker.WatermarkImage, marker.Horizontal, marker.Vertical))
			} else {
				operations = append(operations, ops.TextWatermarkOp(effectiveWatermark, 0.96, 0.96, 0.04, 24, "rgba(255,255,255,0.7)"))
			}
		}
	}

	switch mode {
	case resizeMode:
		resize()
	case fitMode:
		fit()
	case liquidMode:
		liquid()
	case cropMode:
		crop()
	case cropmidMode:
		cropmid()
	}
	addWatermark()

	operations = append(operations, ops.Convert{Format: format})

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
	"":        resizeMode,
	"resize":  resizeMode,
	"fit":     fitMode,
	"crop":    cropMode,
	"cropmid": cropmidMode,
	"liquid":  liquidMode,
}

var supportedFormats = map[string]string{
	"":					"jpeg",
	"jpg":			"jpeg",
	"jpeg":			"jpeg",
	"gif":			"gif",
	"webp":			"webp",
	"bmp":			"bmp",
	"png":			"png",
	"tiff":			"tiff",
}


type mode int

const (
	fitMode = mode(1 + iota)
	cropMode
	cropmidMode
	liquidMode
	resizeMode

	widthParam  = "width"
	heightParam = "height"
	modeParam   = "mode"
	formatParam = "format"
	cropxParam  = "cropx"
	cropyParam  = "cropy"
	urlParam     = "url"
	watermarkParam = "watermark"
)

// returns validated parameters from request and error if invalid
func getParams(a *fasthttp.Args) (w, h, cropx, cropy int, mode mode, format, url, watermark string, err error) {

	if strings.Contains(a.String(), "%3F") { // %3F = ?
		err = errors.New("Invalid characters in request!")
		return
	}

	defer func() {
		if msg := recover(); msg != nil {
			err = msg.(error)
		}
	}()

	w = getUint(a, widthParam)
	h = getUint(a, heightParam)

	cropx = getUint(a, cropxParam)
	cropy = getUint(a, cropyParam)

	mode = stringToMode[string(a.Peek(modeParam))]
	if mode == 0 {
		err = errors.New("Invalid mode!")
		return
	}

	format, formatFound := supportedFormats[strings.ToLower(string(a.Peek(formatParam)))]
	if !formatFound {
		err = errors.New("Invalid format '" + string(a.Peek(formatParam)) + "'!")
		return
	}
	// TODO: verify that the format is one we support.
	// We do not want to support TXT, for instance

	url = string(a.Peek(urlParam))
	watermark = string(a.Peek(watermarkParam))

	a.Del(widthParam)
	a.Del(heightParam)
	a.Del(modeParam)
	a.Del(formatParam)
	a.Del(cropxParam)
	a.Del(cropyParam)
	a.Del(urlParam)
	a.Del(watermarkParam)

	if a.Len() != 0 {
		err = errors.New("Invalid parameter " + a.String())
		return
	}

	err = nil
	return
}

func getUint(a *fasthttp.Args, param string) int {
	v, err := a.GetUint(param)
	if isParseError(err) {
		panic(err)
	}
	if v == -1 {
		v = 0
	}
	return v
}

func isParseError(err error) bool {
	return err != nil && err != fasthttp.ErrNoArgValue
}
