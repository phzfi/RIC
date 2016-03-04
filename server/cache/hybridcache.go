package cache

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
)

// Cache that looks in the first cache first etc.
// Images are stored in every cache
type HybridCache []*Cache

func (caches HybridCache) GetBlob(operations []ops.Operation) (images.ImageBlob, bool) {
	for _, cache := range caches {
		if blob, found := cache.GetBlob(operations); found {
			return blob, true
		}
	}
	return nil, false
}

func (caches HybridCache) AddBlob(operations []ops.Operation, blob images.ImageBlob) {
	for _, cache := range caches {
		cache.AddBlob(operations, blob)
	}
}
