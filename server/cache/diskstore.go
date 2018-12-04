package cache

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
	"fmt"
	"github.com/phzfi/RIC/server/logging"
	"errors"
	"strings"
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

	folders, err := filepath.Glob(folder + "/*")
	if err != nil {
		log.Println("Error reading previously cached files from disk:", err)
	}


	for _, folderPath := range folders {
		folder := filepath.Base(folderPath)
		files, err := filepath.Glob(folderPath + "/*")
		if err != nil {
			log.Println("Error reading previously cached files from disk:", err)
		}
		for _, filePath := range files {
			filename := filepath.Base(filePath)
			bytes, err := encoder.DecodeString(filepath.Base(filePath))
			if err != nil {
				log.Println("Malformed filename", filename, "in previously cached files:", err)
				continue
			}

			string := string(bytes)
			c.policy.Push(string)

			size := fileSize(filePath)
			key := createKey(folder, filename)
			store.entries[key] = entry{filePath, size}
			c.currentMemory += size
		}
	}

	return c
}

func fileSize(path string) uint64 {
	f, err := os.Open(path)
	if err != nil {
		log.Println("Unable to open file to get its size:", err)
		return 0
	}
	defer f.Close()

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

func (d *DiskStore) Load(namespace string, identifier string) (blob []byte, ok bool) {
	d.RLock()
	filename := stringToBase64(identifier)
	key := createKey(namespace, filename)
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

func (d *DiskStore) Store(namespace string, identifier string, blob []byte) {
	filename := stringToBase64(identifier)
	folder := filepath.FromSlash(d.folder + "/" + namespace)
	folderErr := assertFolder(folder)
	if folderErr != nil {
		logging.Debug(folderErr)
		return
	}

	path := filepath.Join(folder, filename)

	go func() {
		err := ioutil.WriteFile(path, blob, os.ModePerm)
		if err != nil {
			log.Println("Unable to write file into disk cache:", err)
		}
		d.Lock()
		key := createKey(namespace, filename)
		d.entries[key] = entry{path, uint64(len(blob))}
		d.Unlock()
	}()
}

func (d *DiskStore) Delete(namespace string, identifier string) uint64 {
	d.Lock()
	entry, ok := d.entries[identifier]

	// Pretty dirty solution, but this path is only used if an image is deleted just after being cached.
	for !ok {
		d.Unlock()
		time.Sleep(100)
		d.Lock()
		entry, ok = d.entries[identifier]
	}

	delete(d.entries, identifier)
	d.Unlock()

	go func() {
		err := os.Remove(entry.Path)
		if err != nil {
			log.Println("Error deleting file from disk cache:", err)
		}
	}()
	return entry.Size
}

func (d *DiskStore) DeleteNamespace(namespace string) (err error) {
	if len(namespace) == 0 {
		return errors.New("invalid namespace given")
	}
	namespaceFolder := filepath.FromSlash(d.folder + "/" + namespace)
	err = os.RemoveAll(namespaceFolder)
	if err != nil {
		logging.Debug(err)
	}
	return
}

func assertFolder(path string) (err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		createErr := os.MkdirAll(path, os.ModePerm)
		if createErr != nil {
			logging.Debug(createErr)
			return createErr
		}
	}

	return
}

func createKey(namespace string, identifier string) string {
	return fmt.Sprintf("%s:%s", namespace, identifier)
}

func splitKey(key string) (namespace string, identifier string) {
	keyParts := strings.Split(key, ":")
	return keyParts[0], keyParts[1]
}
