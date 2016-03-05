package main

import (
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/valyala/fasthttp"
	"github.com/phzfi/RIC/server/ops"
	"testing"
)


type MockImageSource struct {
}

type MockLoadImageOp struct {
}

func (MockLoadImageOp) Apply (images.Image) error {
	return nil
}

func (i MockImageSource) LoadImageOp(id string) ops.Operation {
	return MockLoadImageOp{}
}

func (i MockImageSource) ImageSize(id string) (w int, h int, err error) {
	return 10000, 16000, nil
}

func TestURLParser(t* testing.T) {
	uri := fasthttp.AcquireURI()
	uri.Parse("dummyhost.com", "dummyhost.com/dummyimage.webp?width=20&height=20")
	operations, err := ParseURI(uri, MockImageSource{})
	if err != nil {
		t.Fatal(err)
		return
	}
	resize1 := operations[1]
	resize2 := operations[2]
	resize3 := operations[3]
	resize4 := operations[4]
	
	r := resize1
	w := r.Width
	h := r.Height
	ew := 5000
	eh := 8000
	if w != 5000 || h != 8000 {
		t.Fatal(fmt.Sprintf("Wrong sub-resize. Expected: %v, %v. Got: %v, %v.", ew, eh, w, h))
	}

}
