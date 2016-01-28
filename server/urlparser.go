package main

import (
	"github.com/phzfi/RIC/server/ops"
	"github.com/valyala/fasthttp"
)


func ParseURI(uri *fasthttp.URI, source ops.ImageSource) (operations []ops.Operation, err error){
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
		w = h * ow / oh
	}

	adjustHeight := func() {
		h = w * oh / ow
	}

	resize := func(){
		if herr != nil && werr == nil {
			adjustHeight()
		} else if herr == nil && werr != nil {
			adjustWidth()
		} else if werr != nil && herr != nil {
			w, h = ow, oh
		}
		operations = append(operations, ops.Resize{w, h})	
	}
	
	switch mode {
	case "resize":
		resize()
	case "fit":
		if werr == nil && herr == nil {
			if ow * h > w * oh {
				adjustHeight()
			} else {
				adjustWidth()
			}
			operations = append(operations, ops.Resize{w, h})
		} else {
			resize()
		}
	default:
		resize()
	}
	return
}


