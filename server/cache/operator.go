package cache

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
	"sync"
)

type Operator struct {
	cache Cacher

	sync.Mutex
	inProgress map[cacheKey]*sync.RWMutex
	tokens     TokenPool
}

type Cacher interface {
	GetBlob(cacheKey) (images.ImageBlob, bool)
	AddBlob(cacheKey, images.ImageBlob)
}

func MakeOperator(mm uint64, cacheFolder string) Operator {
	return Operator{
		cache: HybridCache{
			NewLRU(mm),
			NewDiskCache(cacheFolder, 1024*1024*1024*4, NewLRUPolicy()),
		},
		inProgress: make(map[cacheKey]*sync.RWMutex),
		tokens:     MakeTokenPool(2),
	}
}

func (o Operator) GetBlob(operations ...ops.Operation) (blob images.ImageBlob, err error) {

	key := toKey(operations)

	var startimage images.ImageBlob
	var start int

	for start = len(operations); start > 0; start-- {
		var found bool
		startimage, found = o.cache.GetBlob(toKey(operations[:start]))
		if found {
			break
		}
	}

	if start == len(operations) {
		return startimage, nil
	} else {
		o.Lock()
		cond, inProgress := o.inProgress[key]
		if !inProgress {
			// image may have entered cache while this goroutine moved to this place in code
			var found bool
			blob, found = o.cache.GetBlob(key)
			if found {
				o.Unlock()
				return
			}
			cond = o.addInProgress(key)
		}
		o.Unlock()

		if inProgress {
			// Blocks until image has been processed
			cond.RLock()

			var found bool
			blob, found = o.cache.GetBlob(key)
			if found {
				return
			}

			// This only happens if the freshly resized image is dropped from cache too quickly
			o.Lock()
			o.addInProgress(key)
			o.Unlock()
		}

		o.tokens.Borrow()
		defer o.tokens.Return()

		img := images.NewImage()
		defer img.Destroy()

		if start != 0 {
			img.FromBlob(startimage)
		}

		// TODO: do not ignore error
		o.applyOpsToImage(operations[start:], img)
		blob = img.Blob()

		o.cache.AddBlob(key, blob)

		cond.Unlock()
		o.Lock()
		delete(o.inProgress, key)
		o.Unlock()
	}

	return
}

func (o Operator) addInProgress(key cacheKey) *sync.RWMutex {
	m := &sync.RWMutex{}
	m.Lock()
	o.inProgress[key] = m
	return m
}

func (o Operator) applyOpsToImage(operations []ops.Operation, img images.Image) (err error) {
	for _, op := range operations {
		err = op.Apply(img)
		if err != nil {
			return
		}
	}
	return
}
