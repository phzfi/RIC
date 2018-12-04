package cache

// Cache that looks in the first cache first etc.
// Images are stored in every cache
type HybridCache []Cacher

type Cacher interface {
	GetBlob(string, string) ([]byte, bool)
	AddBlob(string, string, []byte)
	DeleteNamespace(string)
}

func (caches HybridCache) GetBlob(namespace string, identifier string) ([]byte, bool) {
	for i, cache := range caches {
		if blob, found := cache.GetBlob(namespace, identifier); found {
			for j := 0; j < i; j++ {
				caches[j].AddBlob(namespace, identifier, blob)
			}
			return blob, true
		}
	}
	return nil, false
}

func (caches HybridCache) AddBlob(namespace string, identifier string, blob []byte) {
	for _, cache := range caches {
		cache.AddBlob(namespace, identifier, blob)
	}
}


func (caches HybridCache) DeleteNamespace(namespace string)  {
	for _, cache := range caches {
		cache.DeleteNamespace(namespace)
	}
}