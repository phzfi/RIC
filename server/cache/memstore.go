package cache

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
)

// Takes the caching policy and the maximum size of the cache in bytes.
func NewCache(policy Policy, mm uint64) *Cache {
	logging.Debugf("Cache create: mem:%v", mm)
	return &Cache{
		maxMemory: mm,
		policy:    policy,
		storer:    make(MemoryStore),
	}
}

type MemoryStore map[cacheKey]images.ImageBlob

func (s MemoryStore) Load(key cacheKey) (b images.ImageBlob, ok bool) {
	b, ok = s[key]
	return
}

func (s MemoryStore) Store(key cacheKey, value images.ImageBlob) {
	s[key] = value
}

func (s MemoryStore) Delete(key cacheKey) (size uint64) {
	size = uint64(len(s[key]))
	delete(s, key)
	return
}
