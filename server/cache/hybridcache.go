package cache

import (
	"github.com/phzfi/RIC/server/images"
)

// Cache that looks in the first cache first etc.
// Images are stored in every cache
type HybridCache []*Cache

func (caches HybridCache) GetBlob(string string) (images.ImageBlob, bool) {
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

func (caches HybridCache) AddBlob(string string, blob images.ImageBlob) {
	for _, cache := range caches {
		cache.AddBlob(string, blob)
	}
}
