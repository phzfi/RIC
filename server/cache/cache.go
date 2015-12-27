package cache

import (
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"sync"
)

type ImageInfo struct {
	name          string
	width, height uint
	original      bool
}

type OriginalInfo struct {
	filename      string
	width, height uint
}

type Cache struct {
	Resizer

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
func NewCache(policy Policy, mm uint64) *Cache {
	return &Cache{
		Resizer:   new(BasicResizer),
		maxMemory: mm,
		policy:    policy,
		blobs:     make(map[ImageInfo]images.ImageBlob),
	}
}

// Gets an image blob of requested dimensions
func (c *Cache) GetImage(filename string, width, height uint) (blob images.ImageBlob, err error) {
	logging.Debug(fmt.Sprintf("Get image: %v, %v, %v", filename, width, height))
	info := ImageInfo{filename, width, height, false}
	if blob, ok := c.blobs[info]; ok {
		c.policy.Visit(info)
		return blob, nil
	}
	// TODO: Requesting nonexistent images causes roots to be accessed unneccessarily. Could it be avoided?
	// TODO: Prevent scenario where requesting the same ImageInfo simultaneously leads to the image being loaded/resized many times.
	blob, err = c.Resizer.GetImage(filename, width, height)
	if err == nil {
		c.addBlob(info, blob)
	}
	return
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
