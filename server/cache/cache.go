package cache


import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"sync"
	"strings"
	"path/filepath"
	"fmt"
)


type ImageInfo struct {
	name                string
	width, height       uint
	original            bool
}

type OriginalInfo struct {
	filename string
	width, height uint
}


type Cache struct {
	Cacheless

	sync.Mutex

	blobs map[ImageInfo]images.ImageBlob
	originals map[string]OriginalInfo

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
		originals: make(map[string]OriginalInfo),
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
	blob, err = c.Cacheless.GetImage(filename, width, height)
	if err == nil {
		c.addBlob(info, blob)
	}
	return
}


// Gets image blob by width only. Image is scaled and aspect ratio preserved.
func (c *Cache) GetImageByWidth(filename string, width uint) (blob images.ImageBlob, err error) {
	logging.Debug(fmt.Sprintf("Get image by width: %v, %v", filename, width))
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	id := strings.TrimRight(filename[0:len(filename)-len(ext)], ".")
	info, ok := c.originals[id]
	if ok {
		height := uint(float64(info.height) / float64(info.width) * float64(width))
		blob, err = c.GetImage(filename, width, height)
		return
	}
	// Image probably not accessed before. Use GetOriginal to access the original image and to store it's info to c.originals
	// TODO: Requesting nonexistent images causes roots to be accessed unneccessarily. Could it be avoided?
	// TODO: Prevent scenario where requesting the same ImageInfo simultaneously leads to the image being loaded/resized
	info, _, err = c.GetOriginal(filename)
	if err != nil {
		return
	}
	c.originals[id] = info
	height := uint(float64(info.height) / float64(info.width) * float64(width))
	blob, err = c.GetImage(filename, width, height)
	return
}



// Gets image blob by height only. Image is scaled and aspect ratio preserved.
func (c *Cache) GetImageByHeight(filename string, height uint) (blob images.ImageBlob, err error) {
	logging.Debug(fmt.Sprintf("Get image by height: %v, %v", filename, height))
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	id := strings.TrimRight(filename[0:len(filename)-len(ext)], ".")
	info, ok := c.originals[id]
	if ok {
		width := uint(float64(info.width) / float64(info.height) * float64(height))
		blob, err = c.GetImage(filename, width, height)
		return
	}
	// Image probably not accessed before. Use GetOriginal to access the original image and to store it's info to c.originals
	// TODO: Requesting nonexistent images causes roots to be accessed unneccessarily. Could it be avoided?
	// TODO: Prevent scenario where requesting the same ImageInfo simultaneously leads to the image being loaded/resized
	info, _, err = c.GetOriginal(filename)
	if err != nil {
		return
	}
	c.originals[id] = info
	width := uint(float64(info.width) / float64(info.height) * float64(height))
	blob, err = c.GetImage(filename, width, height)
	return
}


// Gets image blob by height only. Image is scaled and aspect ratio preserved.
func (c *Cache) GetOriginalSized(filename string) (blob images.ImageBlob, err error) {
	logging.Debug(fmt.Sprintf("Get original sized image: %v", filename))
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	id := strings.TrimRight(filename[0:len(filename)-len(ext)], ".")
	info, ok := c.originals[id]
	if ok {
		width := info.width
		height := info.height
		blob, err = c.GetImage(filename, width, height)
		return
	}
	// Image probably not accessed before. Use GetOriginal to access the original image and to store it's info to c.originals
	// TODO: Requesting nonexistent images causes roots to be accessed unneccessarily. Could it be avoided?
	// TODO: Prevent scenario where requesting the same ImageInfo simultaneously leads to the image being loaded/resized
	info, _, err = c.GetOriginal(filename)
	if err != nil {
		return
	}
	c.originals[id] = info
	width := info.width
	height := info.height
	blob, err = c.GetImage(filename, width, height)
	return
}


// Gets image blob in original format. Also returns original image info.
func (c *Cache) GetOriginal(id string) (originalinfo OriginalInfo, blob images.ImageBlob, err error) {
	logging.Debug(fmt.Sprintf("Get original image: %v", id))
	originalinfo, ok := c.originals[id]
	if ok {
		info := ImageInfo{originalinfo.filename, 0, 0, true}
		blob, ok = c.blobs[info]
		if ok {
			c.policy.Visit(info)
			return
		}
	}
	// TODO: Requesting nonexistent images causes roots to be accessed unneccessarily. Could it be avoided?
	// TODO: Prevent scenario where requesting the same ImageInfo simultaneously leads to the image being loaded/resized
	originalinfo, blob, err = c.Cacheless.GetOriginal(id)
	c.originals[id] = originalinfo
	if err != nil {
		return
	}
	info := ImageInfo{originalinfo.filename, 0, 0, true}
	c.addBlob(info, blob)
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
