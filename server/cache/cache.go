package cache

import (
	"github.com/phzfi/RIC/server/logging"
	"sync"
)

var rwmutex = &sync.RWMutex{}

type Policy interface {
	// Push and Pop do not need to be thread safe
	Push(string)
	Pop() string

	// Image is requested and found in cache. Needs to be thread safe.
	Visit(string)
}

type Storer interface {
	Load(string) ([]byte, bool)
	Store(string, []byte)
	Delete(string) uint64
}

type Cache struct {

	policy Policy
	storer Storer

	maxMemory, currentMemory uint64
}

// Gets an image blob of requested dimensions
func (c *Cache) GetBlob(key string) (blob []byte, found bool) {

	b64 := stringToBase64(key)
	logging.Debugf("Cache get with key: %v", b64)

	rwmutex.RLock()
	blob, found = c.storer.Load(key)
	rwmutex.RUnlock()

	if found {
		logging.Debugf("Cache found: %v", b64)
		c.policy.Visit(key)
	} else {
		logging.Debugf("Cache not found: %v", b64)
	}

	return
}

func (c *Cache) AddBlob(key string, blob []byte) {

	size := uint64(len(blob))

	if size > c.maxMemory {
		return
	}

	logging.Debugf("Cache add: %v", stringToBase64(key))

	// This is the only point where the cache is mutated.
	// While this runs the there can be no reads from the storer.
	rwmutex.Lock()
	defer rwmutex.Unlock()
	for c.currentMemory+size > c.maxMemory {
		c.deleteOne()
	}
	c.policy.Push(key)
	c.currentMemory += uint64(len(blob))
	c.storer.Store(key, blob)
}

func (c *Cache) deleteOne() {
	to_delete := c.policy.Pop()
	logging.Debugf("Cache delete: %v", stringToBase64(to_delete))
	c.currentMemory -= c.storer.Delete(to_delete)
}
