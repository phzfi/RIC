package ops

import "github.com/phzfi/RIC/server/images"

type Operation interface {
	// Returns string representation of the operation used in cache keys.
	Marshal() string

	// Applies the operation to image
	Apply(images.Image) error
}
