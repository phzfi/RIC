package ops

import "github.com/phzfi/RIC/server/images"

type Operation interface {
	GetKey() string
	Apply(images.Image) error
}

// At the moment OperationFunc is not used anywhere.
// OperationFunc is a closure implementing Operation but it does not
// convert to string uniquely and therefore is bad when creating cache
// keys from [Operation]. (Different image id:s might generate same
// cache keys
/*
type OperationFunc func(images.Image) error

func (o OperationFunc) Apply(img images.Image) error {
	return o(img)
}
*/
