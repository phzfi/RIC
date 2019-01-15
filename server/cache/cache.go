package cache

import (
	"github.com/phzfi/RIC/server/logging"
	"sync"
)

type Policy interface {
	// Push and Pop do not need to be thread safe
	Push(string)
	Pop() (string, error)

	// Image is requested and found in cache. Needs to be thread safe.
	Visit(string)
}

type Storer interface {
	Load(string, string) ([]byte, bool)
	Store(string, string, []byte)
	Delete(string, string) uint64
	DeleteNamespace(string) (error)
}

type Cache struct {
	sync.RWMutex

	policy Policy
	storer Storer

	maxMemory, currentMemory uint64
}

// Gets an image blob of requested dimensions
func (c *Cache) GetBlob(namespace string, key string) (blob []byte, found bool) {

	b64 := stringToBase64(key)
	logging.Debugf("Cache get with key: %v:%v", namespace, b64)

	// TODO: Change back to RLock after memstore.go 2d structure refactor
	//c.RLock()
	c.Lock()
	blob, found = c.storer.Load(namespace, key)
	//c.RUnlock()
	c.Unlock()

	if found {
		logging.Debugf("Cache found from %T: %v:%v", c.storer, namespace, b64)
		c.policy.Visit(createKey(namespace, key))
	} else {
		logging.Debugf("Cache not found in %T: %v:%v", c.storer, namespace, b64)
	}

	return
}

func (c *Cache) AddBlob(namespace string, key string, blob []byte) {

	size := uint64(len(blob))

	if size > c.maxMemory {
		return
	}

	logging.Debugf("Cache add to %T: %v:%v", c.storer, namespace, stringToBase64(key))

	// This is the only point where the cache is mutated.
	// While this runs the there can be no reads from the storer.
	c.Lock()
	defer c.Unlock()
	for c.currentMemory+size > c.maxMemory {
		c.deleteOne()
	}
	c.policy.Push(createKey(namespace, key))
	c.currentMemory += uint64(len(blob))
	c.storer.Store(namespace, key, blob)
}

func (c *Cache) deleteOne() {
	toDelete, err := c.policy.Pop()
	if err != nil {
		logging.Debug("Cache delete failed, not items in cache")
		return
	}
	namespace, identifier := splitKey(toDelete)
	logging.Debugf("Cache delete: %v:%v", namespace, stringToBase64(identifier))
	c.currentMemory -= c.storer.Delete(namespace, identifier)
}


func (c *Cache) DeleteNamespace(namespace string) {
	c.storer.DeleteNamespace(namespace)
}