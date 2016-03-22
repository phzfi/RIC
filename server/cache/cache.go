package cache

import (
	"github.com/phzfi/RIC/server/logging"
	"sync"
)

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
	sync.RWMutex

	policy Policy
	storer Storer

	maxMemory, currentMemory uint64
}

// Gets an image blob of requested dimensions
func (c *Cache) GetBlob(string string) (blob []byte, found bool) {

	b64 := stringToBase64(string)
	logging.Debugf("Cache get with string: %v", b64)

	c.RLock()
	blob, found = c.storer.Load(string)
	c.RUnlock()

	if found {
		logging.Debugf("Cache found: %v", b64)
		c.policy.Visit(string)
	} else {
		logging.Debugf("Cache not found: %v", b64)
	}

	return
}

func (c *Cache) AddBlob(string string, blob []byte) {

	// This is the only point where the cache is mutated.
	// While this runs the there can be no reads from the storer.
	size := uint64(len(blob))

	if size > c.maxMemory {
		return
	}

	logging.Debugf("Cache add: %v", stringToBase64(string))

	c.Lock()
	defer c.Unlock()
	for c.currentMemory+size > c.maxMemory {
		c.deleteOne()
	}
	c.policy.Push(string)
	c.currentMemory += uint64(len(blob))
	c.storer.Store(string, blob)
}

func (c *Cache) deleteOne() {
	to_delete := c.policy.Pop()
	logging.Debugf("Cache delete: %v", stringToBase64(to_delete))
	c.currentMemory -= c.storer.Delete(to_delete)
}
