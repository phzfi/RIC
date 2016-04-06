package operator

import (
	"github.com/phzfi/RIC/server/cache"
	"github.com/phzfi/RIC/server/ops"
	"sync"
)

type Operator struct {
	cache Cacher

	sync.Mutex
	inProgress map[string]*Progress

	processor imageProcessor
}

type Progress struct {
	sync.RWMutex
	blob []byte
}

type Cacher interface {
	GetBlob(string) ([]byte, bool)
	AddBlob(string, []byte)
}

func Make(cache Cacher) Operator {
	return Operator{
		cache:      cache,
		inProgress: make(map[string]*Progress),
		processor:  makeImageProcessor(2),
	}
}

func MakeDefault(mm uint64, cacheFolder string) Operator {
	return Make(cache.HybridCache{
		cache.NewCache(cache.NewLRU(), mm),
		cache.NewDiskCache(cacheFolder, 1024*1024*1024*4, cache.NewLRU()),
	})
}

func (o *Operator) GetBlob(operations ...ops.Operation) (blob []byte, err error) {

	key := toKey(operations)

	var startimage []byte
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
		// Only one goroutine gets through this with inProgress being false
		o.Lock()
		isReady, inProgress := o.inProgress[key]
		if !inProgress {
			isReady = o.addInProgress(key)
		}
		o.Unlock()

		if inProgress {
			// Blocks until image has been processed
			isReady.RLock()
			return isReady.blob, nil
		}

		// image may have entered cache while this goroutine moved to this place in code
		var found bool
		blob, found = o.cache.GetBlob(key)
		if !found {
			blob, err = o.processor.makeBlob(startimage, operations[start:])
			if err != nil {
				return nil, err
			}
		}

		isReady.blob = blob
		isReady.Unlock()

		o.cache.AddBlob(key, blob)

		o.Lock()
		delete(o.inProgress, key)
		o.Unlock()
	}

	return
}

func (o *Operator) addInProgress(key string) *Progress {
	p := new(Progress)
	p.Lock()

	o.inProgress[key] = p
	return p
}
