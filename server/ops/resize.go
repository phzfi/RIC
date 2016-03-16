package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

type Resize struct {
	Width, Height int
}

func (r Resize) Marshal() string {
	return string(resize) + int32ToString(int32(r.Width)) + int32ToString(int32(r.Height))
}

func (r Resize) Apply(img images.Image) error {
	logging.Debug("Resizing image to: %v, %v", r.Width, r.Height)
	return img.Resize(r.Width, r.Height)
}

func int32ToString(x int32) string {
	return string([]byte{
		byte((x >> 24) & 0xff),
		byte((x >> 16) & 0xff),
		byte((x >> 8) & 0xff),
		byte(x & 0xff),
	})
}
