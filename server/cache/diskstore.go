package cache

import (
	"encoding/base64"
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

		string := string(bytes)
		c.policy.Push(string)

		size := fileSize(fn)
		store.entries[string] = entry{fn, size}
		c.currentMemory += size
	}

	return c
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

var encoder = base64.RawURLEncoding

func stringToBase64(k string) string {
	return encoder.EncodeToString([]byte(k))
}

type DiskStore struct {
	sync.RWMutex
	entries stringToEntry

	folder string
}

type stringToEntry map[string]entry

type entry struct {
	Path string
	Size uint64
}

func NewDiskStore(folder string) *DiskStore {
	return &DiskStore{
		entries: make(stringToEntry),
		folder:  folder,
	}
}

func (d *DiskStore) Load(string string) (blob []byte, ok bool) {
	d.RLock()
	entry, ok := d.entries[string]
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

func (d *DiskStore) Store(string string, blob []byte) {
	filename := stringToBase64(string)
	path := filepath.Join(filepath.FromSlash(d.folder), filename)

	go func() {
		err := ioutil.WriteFile(path, blob, os.ModePerm)
		if err != nil {
			log.Println("Unable to write file into disk cache:", err)
		}
		d.Lock()
		d.entries[string] = entry{path, uint64(len(blob))}
		d.Unlock()
	}()
}

func (d *DiskStore) Delete(string string) uint64 {
	d.Lock()
	entry, ok := d.entries[string]

	// Pretty dirty solution, but this path is only used if an image is deleted just after being cached.
	for !ok {
		d.Unlock()
		time.Sleep(100)
		d.Lock()
		entry, ok = d.entries[string]
	}

	delete(d.entries, string)
	d.Unlock()

	go func() {
		err := os.Remove(entry.Path)
		if err != nil {
			log.Println("Error deleting file from disk cache:", err)
		}
	}()
	return entry.Size
}
