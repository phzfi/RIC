package ops

import "github.com/phzfi/RIC/server/images"

type Resize struct{
	Width, Height int
}

func (r Resize) Apply(img images.Image) error {
	return img.Resize(r.Width, r.Height)
}
