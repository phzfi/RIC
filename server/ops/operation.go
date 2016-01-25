package ops

import "github.com/phzfi/RIC/server/images"

type Operation interface {
	Apply(images.Image) error
}

type OperationFunc func(images.Image) error

func (o OperationFunc) Apply (img images.Image) error {
	return o(img)
}
