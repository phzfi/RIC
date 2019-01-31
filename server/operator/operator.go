package operator

import (
	"github.com/phzfi/RIC/server/cache"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/ops"
	"github.com/phzfi/RIC/server/ric_file"
	"github.com/valyala/fasthttp"
	"sync"
)

type Operator struct {
	cache Cacher

	sync.Mutex
	inProgress map[string]*Progress

	processor ImageProcessor
}

type Progress struct {
	sync.RWMutex
	blob []byte
}

type Cacher interface {
	GetBlob(string, string) ([]byte, bool)
	AddBlob(string, string, []byte)
	DeleteNamespace(namespace string)
}

func Make(cache Cacher, tokens int) Operator {
	return Operator{
		cache:      cache,
		inProgress: make(map[string]*Progress),
		processor:  MakeImageProcessor(tokens),
	}
}

func MakeWithDefaultCacheSet(mm uint64, cacheFolder string, tokens int) Operator {
	return Make(cache.HybridCache{
		cache.NewCache(cache.NewLRU(), mm),
		cache.NewDiskCache(cacheFolder, 1024*1024*1024*4, cache.NewLRU()),
	}, tokens)
}

func (o *Operator) GetBlob(namespace string, operations ...ops.Operation) (blob []byte, err error) {

	key := ToKey(operations)

	var startimage []byte
	var start int

	for start = len(operations); start > 0; start-- {
		var found bool
		startimage, found = o.cache.GetBlob(namespace, ToKey(operations[:start]))
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
		blob, found = o.cache.GetBlob(namespace, key)
		if !found {
			blob, err = o.processor.MakeBlob(startimage, operations[start:])
			if err != nil {
				return nil, err
			}
		}

		isReady.blob = blob
		isReady.Unlock()

		o.cache.AddBlob(namespace, key, blob)

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

func (o *Operator) DeleteCacheNamespace(uri *fasthttp.URI, source ops.ImageSource) error {
	filename := string(uri.Path())
	_, md5Filename, decodeErr := ric_file.DecodeFilename(filename)
	if decodeErr != nil {
		logging.Debug(decodeErr)
		return decodeErr
	}

	o.cache.DeleteNamespace(md5Filename)

	return nil
}
