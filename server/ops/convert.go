package ops

import (
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type Convert struct {
	Format string
}

func (c Convert) Marshal() string {
	return string(convertID) + c.Format + string(0)
}

func (c Convert) Apply(img images.Image) error {
	logging.Debug(fmt.Sprintf("Converting image to: %v", c.Format))
	return img.Convert(c.Format)
}
