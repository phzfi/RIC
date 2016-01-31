package ops

import (
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type Resize struct {
	Width, Height int
}

func (r Resize) Apply(img images.Image) error {
	logging.Debug(fmt.Sprintf("Resizing image to: %v, %v", r.Width, r.Height))
	return img.Resize(r.Width, r.Height)
}
