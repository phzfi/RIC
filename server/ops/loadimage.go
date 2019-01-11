package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type loadImageOp struct {
	is *ImageSource
	id string
}

func (i loadImageOp) Marshal() string {
	return string(loadID) + i.id + string(0)
}

func (i loadImageOp) Apply(img images.Image) error {
	logging.Debugf("Loading: %v", i.id)
	return i.is.searchRoots(i.id, img)
}
