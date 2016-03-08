package cache

import (
	"github.com/phzfi/RIC/server/ops"
	"github.com/phzfi/RIC/server/images"
	"sync"
)

// Operator processes image operations and uses a cache to cache the results.
type Operator struct {
	cache  *Cache
	tokens chan bool
	ops_mutex_map map[cacheKey]*sync.Mutex
	map_mutex *sync.Mutex
}

// Returns an operations with mm bytes of cache
func MakeOperator(mm uint64) Operator {
	concurrent := 2

	lru := NewLRU(mm)
	c := make(chan bool, concurrent)
	mutex_map := make(map[cacheKey]*sync.Mutex)
	map_mutex := &sync.Mutex{}
	o := Operator{lru, c, mutex_map, map_mutex}
	// TODO: Currently only 2 simult. operations allowed. Increate tokens and make them configurable.
	for i := 0; i < concurrent; i++ {
		o.tokens <- true
	}
	return o
}

// Returns an ImageBlob corresponding the given operations
func (o Operator) GetBlob(operations ...ops.Operation) (blob images.ImageBlob, err error) {

	// Only one identical operation at time
	mUnlock := o.lockOps(operations)
	defer mUnlock()
	// Only n operations at time
	t := <-o.tokens
	defer func() {
		o.tokens <- t
	}()
	
	// Check if requested blob can be found from cache and return if so
	blob, found := o.checkCache(operations)
	if found {
		return
	}

	// Begin operation. Check if any sub result is cached and begin operation from the last cached sub result. i.e if result of operations[:i] is cached we can use that. 
	img, i, err := o.beginOp(operations)
	if err != nil {
		return
	}
	defer img.Destroy()
	
	// Apply operations[i:] to the image
	o.applyOpsToImage(operations[i:], img)
	blob = img.Blob()

	// Cache the result
	o.cache.AddBlob(operations, blob)

	return
}

// Locks a mutex specific to given operations. Used to block same operations being executed simultaneously. Returns unlock function.
func (o Operator) lockOps(operations []ops.Operation) (unlock func ()) {
	key := toKey(operations)
	o.map_mutex.Lock()
	ops_mutex, found := o.ops_mutex_map[key]
	
	if found {
		o.map_mutex.Unlock()
		ops_mutex.Lock()
		return func() {
			ops_mutex.Unlock()
		}
	}
	// Create the mutex if not found.
	ops_mutex = &sync.Mutex{}
	o.ops_mutex_map[key] = ops_mutex
	// Cannot get blocked here.
	ops_mutex.Lock()
	o.map_mutex.Unlock()
	// Creator is also responsible for destroying the mutex - We do not want to keep old mutexes in the map. Return function that Unlocks and removes the mutex from the map
	return func() {
		delete(o.ops_mutex_map, key)
		ops_mutex.Unlock()
	}
}


// Tries to find blob corresponging the given Operation[] from cache. Returns Imageblob and boolean depening on whether or not the blob was found.
func (o Operator) checkCache(operations []ops.Operation) (blob images.ImageBlob, found bool) {
	blob, found = o.cache.GetBlob(operations)
	return
}

// Returns the image operations[i:] can be applied to. Checks cache for any sub result of operations[:i] and returns the result Image and i. If no sub results are found an empty image and i=0 will be returned. The returned image needs to be destroyed.
func (o Operator) beginOp(operations []ops.Operation) (img images.Image, i int, err error) {
	img = images.NewImage()
	for i = len(operations); i > 0; i-- {
		var found bool
		blob, found := o.cache.GetBlob(operations[:i])
		if found {
			err = img.FromBlob(blob)
			return
		}
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
