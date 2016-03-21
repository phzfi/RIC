package operator

import (
	"github.com/phzfi/RIC/server/cache"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
	"sync"
)

type Operator struct {
	cache Cacher

	sync.Mutex
	inProgress map[string]*sync.RWMutex
	tokens     TokenPool
}

type Cacher interface {
	GetBlob(string) (images.ImageBlob, bool)
	AddBlob(string, images.ImageBlob)
}

func Make(cache Cacher) Operator {
	return Operator{
		cache:      cache,
		inProgress: make(map[string]*sync.RWMutex),
		tokens:     MakeTokenPool(2),
	}
}

func MakeDefault(mm uint64, cacheFolder string) Operator {
	return Make(cache.HybridCache{
		cache.NewLRU(mm),
		cache.NewDiskCache(cacheFolder, 1024*1024*1024*4, cache.NewLRUPolicy()),
	})
}

func (o *Operator) GetBlob(operations ...ops.Operation) (blob images.ImageBlob, err error) {

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
		// Only one goroutine gets through this with inProgress == false
		o.Lock()
		isReady, inProgress := o.inProgress[key]
		if !inProgress {
			isReady = o.addInProgress(key)
		}
		o.Unlock()

		if inProgress {
			// Blocks until image has been processed
			isReady.RLock()

			var found bool
			blob, found = o.cache.GetBlob(key)
			if found {
				return
			}

			// This only happens if the freshly resized image is dropped from cache too quickly
			//o.addInProgress(key)
			// TODO: do something about it
		}

		// image may have entered cache while this goroutine moved to this place in code
		var found bool
		blob, found = o.cache.GetBlob(key)
		if found {
			return
		}

		o.tokens.Borrow()
		defer o.tokens.Return()

		img := images.NewImage()
		defer img.Destroy()

		if start != 0 {
			img.FromBlob(startimage)
		}

		// TODO: do not ignore error
		applyOpsToImage(operations[start:], img)
		blob = img.Blob()

		o.cache.AddBlob(key, blob)

		isReady.Unlock()
		o.Lock()
		delete(o.inProgress, key)
		o.Unlock()
	}

	return
}

func (o *Operator) addInProgress(key string) *sync.RWMutex {
	m := &sync.RWMutex{}
	m.Lock()

	o.inProgress[key] = m
	return m
}

func applyOpsToImage(operations []ops.Operation, img images.Image) (err error) {
	for _, op := range operations {
		err = op.Apply(img)
		if err != nil {
			return
		}
	}
	return
}
