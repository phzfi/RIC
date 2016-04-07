package operator

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
)

type token struct{}

type ImageProcessor chan token

func MakeImageProcessor(size int) (t imageProcessor) {
	t = make(imageProcessor, size)

	for i := 0; i < size; i++ {
		t <- token{}
	}

	return
}

// Takes an image as a blob and applies the given operations to it
// startBlob can be nil, in which case operations should start with an image loading operation
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
