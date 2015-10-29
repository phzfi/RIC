package cacheless

import (
	"server/images"
)

type Cacheless struct {
}

func (self *Cacheless) GetImage(filename string, width *uint, height *uint) (images.ImageBlob, error) {
	blob, err := images.LoadImage(filename)
	if err != nil {
		return nil, err
	}
	// TODO: If only one is set?
	if width != nil && height != nil {
		blob, err := blob.Resized(*width, *height)
		if err != nil {
			return nil, err
		}
	}
	return blob, nil
}
