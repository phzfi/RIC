package ops

import "github.com/phzfi/RIC/server/images"

type Operation interface {
	Marshal() string
	Apply(images.Image) error
}
