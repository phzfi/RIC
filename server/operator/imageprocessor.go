package operator

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
)

type token struct{}

type imageProcessor chan token

func makeImageProcessor(size int) (t imageProcessor) {
	t = make(imageProcessor, size)

	for i := 0; i < size; i++ {
		t <- token{}
	}

	return
}

func (p imageProcessor) MakeBlob(startBlob []byte, operations []ops.Operation) ([]byte, error) {
	p.borrow()
	defer p.giveBack()

	img := images.NewImage()
	defer img.Destroy()

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

func (p imageProcessor) borrow() {
	<-p
}

func (p imageProcessor) giveBack() {
	p <- token{}
}
