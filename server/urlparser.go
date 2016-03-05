package main

import (
	"github.com/phzfi/RIC/server/ops"
	"github.com/phzfi/RIC/server/logging"
	"github.com/valyala/fasthttp"
	"path/filepath"
	"strings"
)

// Struct containing variable and methods needed when parsing url
type parseJob struct {
	args             *fasthttp.Args
	w, h             int
	ow, oh           int
	filename         string
	mode             string
	werr, herr, oerr error
	operations       []ops.Operation
}

// Adjust w according to the original aspect
func (p *parseJob) adjustWidth() {
	p.w = int(float32(p.h*p.ow)/float32(p.oh) + 0.5)
}

//Adjust h according to the original aspect
func (p *parseJob) adjustHeight() {
	p.h = int(float32(p.w*p.oh)/float32(p.ow) + 0.5)
}

// Adjust w, h so that the image will be correct aspect
func (p *parseJob) adjustSize() {
	if p.herr != nil && p.werr == nil {
		p.adjustHeight()
	} else if p.herr == nil && p.werr != nil {
		p.adjustWidth()
	} else if p.werr != nil && p.herr != nil {
		p.w, p.h = p.ow, p.oh
	}
}

// Adjust w, h so the image won't be upsaceld
func (p *parseJob) denyUpscale() {
	if p.w > p.ow {
		p.w = p.ow
		p.adjustHeight()
	}
	if p.h > p.oh {
		p.h = p.oh
		p.adjustWidth()
	}
}

// Generate resizes stack that halves the w, h until we are close to the requested
// size. The sub resized images will get cached and fasten future resizes.
func (p *parseJob) subResize() {
	hw := p.ow / 2
	hh := p.oh / 2
	for p.w < hw && p.h < hh && hw >= 600 && hh >= 600{
		logging.Debugf("Sub to %v %v", hw, hh)
		p.operations = append(p.operations, ops.Resize{hw, hh})
		hw = hw / 2
		hh = hh / 2
	}
}

// Generate resize operation stack (without loadImageOp)
func (p *parseJob) resize() {
	// Adjust w, h
	p.denyUpscale()
	p.adjustSize()
	// Generate sub-resizes
	p.subResize()
	// Add resize to the final size
	p.operations = append(p.operations, ops.Resize{p.w, p.h})
}

// Generate liquid resize stack (without loadImageOp)
func (p *parseJob) liquid() {
	p.denyUpscale()
	p.adjustSize()
	p.operations = append(p.operations, ops.LiquidRescale{p.w, p.h})
}

// Generate fit resize stack (without loadImageOp)
func (p *parseJob) fit() {
	if p.werr == nil && p.herr == nil {
		p.denyUpscale()
		if p.ow*p.h > p.w*p.oh {
			p.adjustHeight()
		} else {
			p.adjustWidth()
		}
		p.operations = append(p.operations, ops.Resize{p.w, p.h})
	} else {
		p.resize()
	}
}

// Convert file extension to format string
func extToFormat(ext string) string {
	ext = strings.ToUpper(strings.TrimLeft(ext, "."))
	if ext == "JPG" {
		return "JPEG"
	}
	return ext
}

// Generate []Operation thet the given URI represents
func ParseURI(uri *fasthttp.URI, source ops.ImageSourcer) (operations []ops.Operation, err error) {
	p := parseJob{}
	p.args = uri.QueryArgs()
	p.w, p.werr = p.args.GetUint("width")
	p.h, p.herr = p.args.GetUint("height")
	p.filename = string(uri.Path())
	p.ow, p.oh, p.oerr = source.ImageSize(p.filename)
	p.mode = string(p.args.Peek("mode"))
	p.operations = []ops.Operation{source.LoadImageOp(p.filename)}

	if p.oerr != nil {
		err = p.oerr
		return
	}

	switch p.mode {
	case "resize":
		p.resize()
	case "fit":
		p.fit()
	case "liquid":
		p.liquid()
	default:
		p.resize()
	}

	ext := filepath.Ext(p.filename)
	if ext != "" {
		p.operations = append(p.operations, ops.Convert{extToFormat(ext)})
	}

	operations = p.operations

	return
}
