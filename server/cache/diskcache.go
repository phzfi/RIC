package cache

import (
	"github.com/phzfi/RIC/server/images"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type DiskStorer struct {
	keyToPath map[cacheKey]string

	folder string
}

func NewDiskCache(mm uint64) *Cache {
	return &Cache{
		maxMemory: mm,
		policy:    NewLRUPolicy(),
		storer: &DiskStorer{
			keyToPath: make(map[cacheKey]string),
			folder:    "/tmp/kuvia",
		},
	}
	// TODO: load opsToPath by doing a directory listing
}

func (d *DiskStorer) Load(key cacheKey) (blob images.ImageBlob, ok bool) {
	path, ok := d.keyToPath[key]
	if ok {
		var err error
		blob, err = ioutil.ReadFile(path)
		if err != nil {
			log.Println("Error reading file in DiskStorer:", err)
		}
	}
	return
}

func (d *DiskStorer) Store(key cacheKey, blob images.ImageBlob) {
	path := filepath.Join(filepath.FromSlash(d.folder), string(key))
	d.keyToPath[key] = path
	ioutil.WriteFile(path, blob, os.ModePerm)
}

func (d *DiskStorer) Delete(key cacheKey) (size uint64) {
	path := d.keyToPath[key]

	f, err := os.Open(path)
	if err != nil {
		log.Println("Unable to open file in DiskStorer:", err)
	}
	stat, err := f.Stat()
	if err != nil {
		log.Println("Unable to get file stats in Diskstorer:", err)
	}
	size = uint64(stat.Size())

	err = os.Remove(path)
	if err != nil {
		log.Println("Error deleting file in DiskStorer:", err)
	}
	delete(d.keyToPath, key)

	return
}
