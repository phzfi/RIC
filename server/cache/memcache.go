package cache

import (
	"github.com/phzfi/RIC/server/images"
)

// Takes the caching policy and the maximum size of the cache in bytes.
func NewCache(policy Policy, mm uint64) *Cache {
	logging.Debugf("Cache create: mem:%v", mm)
	return &Cache{
		maxMemory: mm,
		policy:    policy,
		storer:    make(MemoryStore),
		blobs:     make(map[cacheKey]images.ImageBlob),
	}
}

type MemoryStore map[cacheKey]images.ImageBlob

func (s MemoryStore) Load(key cacheKey) (images.ImageBlob, bool) {
	return s[key]
}

func (s MemoryStore) Store(key cacheKey, value images.ImageBlob) {
	s[key] = value
}

func (s MemoryStore) Delete(key cacheKey) {
	delete(s, key)
}
