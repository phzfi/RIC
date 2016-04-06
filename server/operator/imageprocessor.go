package operator

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
)

type token struct{}

type imageProcessor chan images.Image

func makeImageProcessor(size int) (t imageProcessor) {
	t = make(imageProcessor, size)

	for i := 0; i < size; i++ {
		t <- images.NewImage()
	}

	return
}

func (p imageProcessor) makeBlob(startBlob []byte, operations []ops.Operation) ([]byte, error) {
	img := p.borrow()
	defer p.giveBack(img)

	if startBlob != nil {
		img.FromBlob(startBlob)
	}

	for _, op := range operations {
		err := op.Apply(img)
		if err != nil {
			return nil, err
		}
	}

	return img.Blob(), nil
}

func (p imageProcessor) borrow() images.Image {
	return <-p
}

func (p imageProcessor) giveBack(img images.Image) {
	p <- img
}
