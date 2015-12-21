package cache


import (
	"github.com/phzfi/RIC/server/images"
	"sync"
)


type ImageInfo struct {
	name          string
	width, height uint
	original      bool
}


type Cache struct {
	Cacheless

	sync.Mutex

	blobs map[ImageInfo]images.ImageBlob

	policy                   Policy
	maxMemory, currentMemory uint64
}


type Policy interface {
	Push(ImageInfo)
	Pop() ImageInfo

	// Image is requested and found in cache. Needs to be thread safe.
	Visit(ImageInfo)
}


// Takes the caching policy and the maximum size of the cache in bytes.
func New(policy Policy, mm uint64) *Cache {
	return &Cache{
		maxMemory: mm,
		policy:    policy,
		blobs:     make(map[ImageInfo]images.ImageBlob),
	}
}


// Gets an image blob of requested dimensions
func (c *Cache) GetImage(filename string, width, height uint) (images.ImageBlob, error) {
	info := ImageInfo{filename, width, height, false}
	if blob, ok := c.blobs[info]; ok {
		c.policy.Visit(info)
		return blob, nil
	}
	// TODO: Prevent scenario where requesting the same ImageInfo simultaneously leads to the image being loaded/resized many times.
	blob, err := c.Cacheless.GetImage(filename, width, height)
	if err == nil {
		c.addBlob(info, blob)
	}
	return blob, err
}


// Gets and image blob of original image dimensions
func (c *Cache) GetOriginal(filename string) (images.ImageBlob, error) {
	info := ImageInfo{filename, 0, 0, true}
	if blob, ok := c.blobs[info]; ok {
		c.policy.Visit(info)
		return blob, nil
	}
	// TODO: Prevent scenario where requesting the same ImageInfo simultaneously leads to the image being loaded/resized
	blob, err := c.Cacheless.GetOriginal(filename)
	if err == nil {
		c.addBlob(info, blob)
	}
	return blob, err
}


func (c *Cache) addBlob(info ImageInfo, blob images.ImageBlob) {

	// This is the only point where the cache is mutated, and therefore can't run in parallel.
	// GetImage can be run in parallel even during this operation due map being thread safe.
	c.Lock()
	defer c.Unlock()

	if _, ok := c.blobs[info]; ok {
		return
	}

	size := uint64(len(blob))

	if size > c.maxMemory {
		return
	}

	for c.currentMemory+size > c.maxMemory {
		c.deleteOldest()
	}

	c.policy.Push(info)

	c.currentMemory += uint64(len(blob))
	c.blobs[info] = blob
}


func (c *Cache) deleteOldest() {

	to_delete := c.policy.Pop()

	c.currentMemory -= uint64(len(c.blobs[to_delete]))
	delete(c.blobs, to_delete)
}
