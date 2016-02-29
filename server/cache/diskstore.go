package cache

import (
	"encoding/base64"
	"github.com/phzfi/RIC/server/images"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var encoder = base64.RawURLEncoding

type DiskStore struct {
	keyToPath map[cacheKey]string

	folder string
}

func NewDiskStore(folder string) *DiskStore {
	d := DiskStore{
		keyToPath: make(map[cacheKey]string),
		folder:    folder,
	}

	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		log.Println("Unable to create folder for disk-caching:", err)
	}

	files, err := filepath.Glob(d.folder + "/*")
	if err != nil {
		log.Println("Error reading previously cached files from disk:", err)
	}
	for _, fn := range files {
		bytes, err := encoder.DecodeString(filepath.Base(fn))
		if err != nil {
			log.Println("Malformed filename", fn, "in previously cached files:", err)
			continue
		}
		d.keyToPath[cacheKey(bytes)] = fn
	}
	return &d
}

func (d *DiskStore) Load(key cacheKey) (blob images.ImageBlob, ok bool) {
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

func (d *DiskStore) Store(key cacheKey, blob images.ImageBlob) {
	filename := encoder.EncodeToString([]byte(key))
	path := filepath.Join(filepath.FromSlash(d.folder), filename)
	d.keyToPath[key] = path
	err := ioutil.WriteFile(path, blob, os.ModePerm)
	if err != nil {
		log.Println("Unable to write file into disk cache:", err)
	}
}

func (d *DiskStore) Delete(key cacheKey) (size uint64) {
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
		return 0
	}
	stat, err := f.Stat()
	if err != nil {
		log.Println("Unable to get file stats:", err)
		return 0
	}
	return uint64(stat.Size())
}
