package cache

import (
	"encoding/base64"
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

func NewDiskCache(policy Policy, mm uint64) *Cache {
	return &Cache{
		maxMemory: mm,
		policy:    policy,
		storer: &DiskStorer{
			keyToPath: make(map[cacheKey]string),
			folder:    "/tmp/kuvia",
		},
	}
	// TODO: load keyToPath by doing a directory listing
}

func (d *DiskStorer) Load(key cacheKey) (blob images.ImageBlob, ok bool) {
	path, ok := d.keyToPath[key]
	if ok {
		var err error
		blob, err = ioutil.ReadFile(path)
		if err != nil {
			log.Println("Error reading file from disk cache:", err)
			ok = false
		}
	}
	return
}

func (d *DiskStorer) Store(key cacheKey, blob images.ImageBlob) {
	filename := base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(key))
	path := filepath.Join(filepath.FromSlash(d.folder), filename)
	d.keyToPath[key] = path
	err := ioutil.WriteFile(path, blob, os.ModePerm)
	if err != nil {
		log.Println("Unable to write file into disk cache:", err)
	}
}

func (d *DiskStorer) Delete(key cacheKey) (size uint64) {
	path := d.keyToPath[key]

	size = fileSize(path)

	err := os.Remove(path)
	if err != nil {
		log.Println("Error deleting file in from disk cache:", err)
	}
	delete(d.keyToPath, key)

	return
}

func fileSize(path string) uint64 {
	f, err := os.Open(path)
	if err != nil {
		log.Println("Unable to open file to get its size:", err)
	}
	stat, err := f.Stat()
	if err != nil {
		log.Println("Unable to get file stats:", err)
	}
	return uint64(stat.Size())
}
