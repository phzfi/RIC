package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type loadImageOp struct {
	is ImageSource
	id string
}

func (i loadImageOp) GetKey() string {
	return i.id
}

func (i loadImageOp) Apply(img images.Image) error {
	logging.Debug("Loading: %v", i.id)
	return i.is.searchRoots(i.id, img)
}
