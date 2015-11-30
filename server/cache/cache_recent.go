package cache

import (
	"github.com/phzfi/RIC/server/images"
	"sync"
)

type imageInfo struct {
	name          string
	width, height uint
}

type CacheRecent struct {
	Cacheless

	sync.Mutex

	blobs map[imageInfo]images.ImageBlob

	deletionQueue            imageInfoQueue
	maxMemory, currentMemory uint
}

// Takes the maximum size of the cache in bytes
func NewCacherecent(mm uint) *CacheRecent {
	return &CacheRecent{maxMemory: mm, blobs: make(map[imageInfo]images.ImageBlob)}
}

func (c *CacheRecent) GetImage(filename string, width, height uint) (images.ImageBlob, error) {

	info := imageInfo{filename, width, height}

	if blob, ok := c.blobs[info]; ok {
		return blob, nil
	}

	// TODO: Prevent scenario where requesting the same imageinfo simultaneously leads to the image being resized many times.
	blob, err := c.Cacheless.GetImage(filename, width, height)
	if err == nil {
		c.addBlob(info, blob)
	}

	return blob, err
}

func (c *CacheRecent) addBlob(info imageInfo, blob images.ImageBlob) {

	// This is the only point where the cache is mutated, and therefore can't run in parallel.
	// GetImage can be run in parallel even during this operation due map being thread safe.
	c.Lock()
	defer c.Unlock()

	if _, ok := c.blobs[info]; ok {
		return
	}

	size := uint(len(blob))

	if size > c.maxMemory {
		return
	}

	for c.currentMemory+size > c.maxMemory {
		c.deleteOldest()
	}

	c.blobs[info] = blob
	c.currentMemory += uint(len(blob))
	c.deletionQueue.Push(info)
}

func (c *CacheRecent) deleteOldest() {

	to_delete := c.deletionQueue.Pop()

	c.currentMemory -= uint(len(c.blobs[to_delete]))
	delete(c.blobs, to_delete)
}
