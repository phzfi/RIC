package cache

import (
	"github.com/phzfi/RIC/server/images"
)

// Cache that looks in the first cache first etc.
// Images are stored in every cache
type HybridCache []*Cache

func (caches HybridCache) GetBlob(key cacheKey) (images.ImageBlob, bool) {
	for i, cache := range caches {
		if blob, found := cache.GetBlob(key); found {
			for j := 0; j < i; j++ {
				caches[j].AddBlob(key, blob)
			}
			return blob, true
		}
	}
	return nil, false
}

func (caches HybridCache) AddBlob(key cacheKey, blob images.ImageBlob) {
	for _, cache := range caches {
		cache.AddBlob(key, blob)
	}
}
