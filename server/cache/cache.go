package cache

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/ops"
	"sync"
)

type Cache struct {
	sync.RWMutex

	policy Policy
	storer Storer

	maxMemory, currentMemory uint64
}

type Policy interface {
	// Push and Pop do not need to be thread safe
	Push(cacheKey)
	Pop() cacheKey

	// Image is requested and found in cache. Needs to be thread safe.
	Visit(cacheKey)
}

type Storer interface {
	Load(cacheKey) (images.ImageBlob, bool)
	Store(cacheKey, images.ImageBlob)
	Delete(cacheKey) uint64
}

// Gets an image blob of requested dimensions
func (c *Cache) GetBlob(operations []ops.Operation) (blob images.ImageBlob, found bool) {
	key := toKey(operations)

	b64 := keyToBase64(key)
	logging.Debugf("Cache get with key: %v", b64)

	// TODO: GetBlob calls policy.Visit(), AddBlob calls policy.Push().
	// Figure out how thread safety should be handled. Is this current
	// solution ok?
	c.RLock()
	defer c.RUnlock()
	blob, found = c.storer.Load(key)

	if found {
		logging.Debugf("Cache found: %v", b64)
		c.policy.Visit(key)
	} else {
		logging.Debugf("Cache not found: %v", b64)
	}

	return
}

func (c *Cache) AddBlob(operations []ops.Operation, blob images.ImageBlob) {

	// This is the only point where the cache is mutated.
	// While this runs the there can be no reads from "blobs".
	size := uint64(len(blob))

	if size > c.maxMemory {
		return
	}

	key := toKey(operations)
	logging.Debugf("Cache add: %v", keyToBase64(key))

	c.Lock()
	defer c.Unlock()
	for c.currentMemory+size > c.maxMemory {
		c.deleteOne()
	}
	c.policy.Push(key)
	c.currentMemory += uint64(len(blob))
	logging.Debugf("New cache size: %v", c.currentMemory)
	c.storer.Store(key, blob)
}

func (c *Cache) deleteOne() {
	to_delete := c.policy.Pop()
	logging.Debugf("Cache delete: %v", to_delete)
	c.currentMemory -= c.storer.Delete(to_delete)
	logging.Debugf("New cache size: %v", c.currentMemory)
}
