package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"fmt"
)


type Convert struct {
	Format string
}

func (c Convert) Apply(img images.Image) error {
	logging.Debug(fmt.Sprintf("Converting image to: %v", c.Format))
	return img.Convert(c.Format)
}
