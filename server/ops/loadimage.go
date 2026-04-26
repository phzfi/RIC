package ops

import (
	"strconv"

	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type loadImageOp struct {
	is *ImageSource
	id string
}

func (i loadImageOp) Marshal() string {
	return strconv.Itoa(int(loadID)) + i.id + "0"
}

func (i loadImageOp) Apply(img images.Image) error {
	logging.Debugf("Loading: %v", i.id)
	return i.is.searchRoots(i.id, img)
}
