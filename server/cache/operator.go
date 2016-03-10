package cache

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
	"sync"
)

type Operator struct {
	cache *Cache

	sync.Mutex
	inProgress map[cacheKey]*sync.Cond
	tokens     TokenPool
}

func MakeOperator(mm uint64) Operator {
	return Operator{
		cache:      NewLRU(mm),
		inProgress: make(map[cacheKey]*sync.Cond),
		tokens:     MakeTokenPool(2),
	}
}

func (o Operator) GetBlob(operations ...ops.Operation) (blob images.ImageBlob, err error) {

	key := toKey(operations)

	var startimage images.ImageBlob
	var start int

	for start = len(operations); start > 0; start-- {
		var found bool
		startimage, found = o.cache.GetBlob(operations[:start])
		if found {
			break
		}
	}

	if start == len(operations) {
		return startimage, nil
	} else {
		o.Lock()
		cond, ok := o.inProgress[key]
		if !ok {
			o.inProgress[key] = sync.NewCond(&sync.Mutex{})
		}
		o.Unlock()

		if ok {
			cond.Wait()
			var found bool
			blob, found = o.cache.GetBlob(operations)
			if found {
				return
			}

			// This only happens if the freshly resized image is dropped from cache too quickly
			o.Lock()
			o.inProgress[key] = sync.NewCond(&sync.Mutex{})
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

		o.cache.AddBlob(operations, blob)
	}

	return
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
