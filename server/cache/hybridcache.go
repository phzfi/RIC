package cache

// Cache that looks in the first cache first etc.
// Images are stored in every cache
type HybridCache []Cacher

type Cacher interface {
	GetBlob(string) ([]byte, bool)
	AddBlob(string, []byte)
}

func (caches HybridCache) GetBlob(string string) ([]byte, bool) {
	for i, cache := range caches {
		if blob, found := cache.GetBlob(string); found {
			for j := 0; j < i; j++ {
				caches[j].AddBlob(string, blob)
			}
			return blob, true
		}
	}
	return nil, false
}

func (caches HybridCache) AddBlob(string string, blob []byte) {
	for _, cache := range caches {
		cache.AddBlob(string, blob)
	}
}
