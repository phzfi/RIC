package cache

import "github.com/phzfi/RIC/server/images"

type imageinfo struct {
	name          string
	width, height uint
}

type CacheRecent struct {
	Cacheless

	blobs map[imageinfo]images.ImageBlob
}

func (c *CacheRecent) GetImage(filename string, width, height uint) (images.ImageBlob, error) {

	info := imageinfo{filename, width, height}

	if blob, ok := c.blobs[info]; ok {
		return blob, nil
	}

	blob, err := c.Cacheless.GetImage(filename, width, height)

	if err == nil {
		c.blobs[info] = blob
	}

	return blob, err
}
