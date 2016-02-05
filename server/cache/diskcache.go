package cache

import (
	"crypto/md5"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type DiskStorer struct {
	opsToPath map[cacheKey]string

	folder string
}

func NewDiskCache() *DiskCache {
	return &DiskCache{
		opsToPath: make(map[cacheKey]string),
		policy:    NewLRUPolicy(),
		folder:    "/tmp/kuvia",
	}
	// TODO: load opsToPath by doing a directory listing
}

func (d *DiskStorer) Load(key cacheKey) (blob images.ImageBlob, ok bool) {
	return nil, false
}

func (d *DiskStorer) Store(key cacheKey, blob images.ImageBlob) {
	path := filepath.Join(filepath.FromSlash(folder), key)
	opsToPath[key] = path
	ioutil.WriteFile(path, blob, os.ModePerm)
}
