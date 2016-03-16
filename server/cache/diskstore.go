package cache

import (
	"encoding/base64"
	"github.com/phzfi/RIC/server/images"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// initializes cache with images found in given folder
func NewDiskCache(folder string, mm uint64, policy Policy) *Cache {
	store := NewDiskStore(folder)
	c := &Cache{
		maxMemory: mm,
		policy:    policy,
		storer:    store,
	}

	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		log.Println("Unable to create folder for disk-caching:", err)
	}

	files, err := filepath.Glob(folder + "/*")
	if err != nil {
		log.Println("Error reading previously cached files from disk:", err)
	}
	for _, fn := range files {
		bytes, err := encoder.DecodeString(filepath.Base(fn))
		if err != nil {
			log.Println("Malformed filename", fn, "in previously cached files:", err)
			continue
		}

		key := cacheKey(bytes)
		c.policy.Push(key)

		size := fileSize(fn)
		store.entries[key] = entry{fn, size}
		c.currentMemory += size
	}

	return c
}

var encoder = base64.RawURLEncoding

func keyToBase64(k cacheKey) string {
	return encoder.EncodeToString([]byte(k))
}

type keyToEntry map[cacheKey]entry

type entry struct {
	Path string
	Size uint64
}

type DiskStore struct {
	sync.RWMutex
	entries keyToEntry

	folder string
}

func NewDiskStore(folder string) *DiskStore {
	return &DiskStore{
		entries: make(keyToEntry),
		folder:  folder,
	}
}

func (d *DiskStore) Load(key cacheKey) (blob images.ImageBlob, ok bool) {
	d.RLock()
	entry, ok := d.entries[key]
	d.RUnlock()

	if ok {
		var err error
		blob, err = ioutil.ReadFile(entry.Path)
		if err != nil {
			log.Println("Error reading file from disk cache:", err)
			ok = false
		}
	}
	return
}

func (d *DiskStore) Store(key cacheKey, blob images.ImageBlob) {
	filename := keyToBase64(key)
	path := filepath.Join(filepath.FromSlash(d.folder), filename)

	go func() {
		err := ioutil.WriteFile(path, blob, os.ModePerm)
		if err != nil {
			log.Println("Unable to write file into disk cache:", err)
		}
		d.Lock()
		d.entries[key] = entry{path, uint64(len(blob))}
		d.Unlock()
	}()
}

func (d *DiskStore) Delete(key cacheKey) uint64 {
	d.Lock()
	entry, ok := d.entries[key]

	// Pretty dirty solution, but this path is only used if an image is deleted just after being cached.
	for !ok {
		d.Unlock()
		time.Sleep(100)
		d.Lock()
		entry, ok = d.entries[key]
	}

	delete(d.entries, key)
	d.Unlock()

	go func() {
		err := os.Remove(entry.Path)
		if err != nil {
			log.Println("Error deleting file from disk cache:", err)
		}
	}()
	return entry.Size
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
