package cache

import (
	"encoding/json"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/ops"
	"md5"
	"sync"
)

type cacheKey string

// Returns a unique representation of an ops chain. This unique representation can be used as a map key unlike the original ops chain (slice cannot be a key).
func toKey(operations []ops.Operation) cacheKey {
	// Todo: Instead of JSON, use compact binary encoding
	bytes, err := json.Marshal(operations)
	if err != nil {
		panic(err)
	}
	return cacheKey(string(md5.Sum(bytes)))
}

type Cache struct {
	sync.RWMutex

	policy Policy
	storer LoadStorer

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
	Delete(cacheKey)
}

// Gets an image blob of requested dimensions
func (c *Cache) GetBlob(operations []ops.Operation) (blob images.ImageBlob, found bool) {
	key := toKey(operations)
	logging.Debugf("Cache get with key: %v", key)

	// TODO: GetBlob calls policy.Visit(), AddBlob calls policy.Push().
	// Figure out how thread safety should be handled. Is this current
	// solution ok?
	c.RLock()
	defer c.RUnlock()
	blob, found = c.storer.Load(key)

	if found {
		logging.Debugf("Cache found: %v", key)
		c.policy.Visit(key)
	} else {
		logging.Debugf("Cache not found: %v", key)
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
	logging.Debugf("Cache add: %v", key)

	c.Lock()
	defer c.Unlock()
	for c.currentMemory+size > c.maxMemory {
		c.deleteOne()
	}
	c.policy.Push(key)
	c.currentMemory += uint64(len(blob))
	c.storer.Store(key, blob)
}

func (c *Cache) deleteOne() {
	to_delete := c.policy.Pop()
	logging.Debugf("Cache delete: %v", to_delete)
	c.currentMemory -= uint64(len(c.blobs[to_delete]))
	c.storer.Delete(to_delete)
}
