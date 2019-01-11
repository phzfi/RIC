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
	Load(string, string) ([]byte, bool)
	Store(string, []byte, string)
	Delete(string, string) uint64
	DeleteNamespace(string) (error)
}

type Cache struct {
	sync.RWMutex

	policy Policy
	storer Storer

	maxMemory, currentMemory uint64
	namespace string
}

// Gets an image blob of requested dimensions
func (c *Cache) GetBlob(namespace string, key string) (blob []byte, found bool) {

	b64 := stringToBase64(key)
	logging.Debugf("Cache get with key: %v:%v", namespace, b64)

	c.RLock()
	blob, found = c.storer.Load(key, namespace)
	c.RUnlock()

	if found {
		logging.Debugf("Cache found: %v:%v", namespace, b64)
		c.policy.Visit(key)
	} else {
		logging.Debugf("Cache not found: %v:%v", namespace, b64)
	}

	return
}

func (c *Cache) AddBlob(namespace string, key string, blob []byte) {

	size := uint64(len(blob))

	if size > c.maxMemory {
		return
	}

	logging.Debugf("Cache add: %v:%v", namespace, stringToBase64(key))

	// This is the only point where the cache is mutated.
	// While this runs the there can be no reads from the storer.
	c.Lock()
	defer c.Unlock()
	for c.currentMemory+size > c.maxMemory {
		c.deleteOne()
	}
	c.policy.Push(key)
	c.currentMemory += uint64(len(blob))
	c.storer.Store(key, blob, namespace)
}

func (c *Cache) deleteOne() {
	// TODO: FIX namespace
	//to_delete := c.policy.Pop()
	//logging.Debugf("Cache delete: %v:%v", c.namespace, stringToBase64(to_delete))
	//c.currentMemory -= c.storer.Delete(to_delete, c.namespace)
}


func (c *Cache) DeleteNamespace(namespace string) {
	c.storer.DeleteNamespace(namespace)
}