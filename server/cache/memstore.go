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

type MemoryStore map[string]images.ImageBlob

func (s MemoryStore) Load(string string) (b images.ImageBlob, ok bool) {
	b, ok = s[string]
	return
}

func (s MemoryStore) Store(string string, value images.ImageBlob) {
	s[string] = value
}

func (s MemoryStore) Delete(string string) (size uint64) {
	size = uint64(len(s[string]))
	delete(s, string)
	return
}
