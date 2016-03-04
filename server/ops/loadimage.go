package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type loadImageOp struct {
	is ImageSource
	Id string
}

func (i loadImageOp) Apply(img images.Image) (err error) {
	logging.Debug("Loading: %v", i.id)
	err = i.is.searchRoots(i.id, img)
	return
}
