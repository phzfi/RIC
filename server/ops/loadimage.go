package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type loadImageOp struct {
	is ImageSource
	Id string
}

func (i loadImageOp) Apply(img images.Image) error {
	logging.Debug("Loading: %v", i.Id)
	return i.is.searchRoots(i.Id, img)
}
