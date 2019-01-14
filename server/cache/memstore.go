package cache

import (
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

type MemoryStore map[string]map[string][]byte

func (s MemoryStore) Load(namespace string, identifier string) (b []byte, ok bool) {
	b, ok = s[namespace][identifier]
	if !ok {
		mm := make(map[string][]byte)
		s[namespace] = mm
	}
	return
}

func (s MemoryStore) Store(namespace string, identifier string, value []byte) {
	s[namespace][identifier] = value
}

func (s MemoryStore) Delete(namespace string, identifier string) (size uint64) {
	//key := createKey(namespace, identifier)
	size = uint64(len(s[namespace][identifier]))
	delete(s[namespace], identifier)
	return
}

func (s MemoryStore) DeleteNamespace(namespace string) (err error) {
	delete(s, namespace)
	return
}