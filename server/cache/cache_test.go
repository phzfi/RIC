package cache

import (
	"bytes"
	"github.com/phzfi/RIC/server/testutils"
	"testing"
	"time"
)

func TestMemCache(t *testing.T) {
	allTests(t, setupMemcache)
}

func TestDiskCache(t *testing.T) {
	allTests(t, func() (*DummyPolicy, *Cache) {
		testutils.RemoveContents(testutils.CacheFolder)
		return setupDiskCache()
	})
}

func TestDiskCachePersistence(t *testing.T) {
	id := string("testdiskpersist")
	data := []byte{1, 2, 3, 4, 7}

	_, cache := setupDiskCache()
	cache.AddBlob(id, data)

	time.Sleep(100)

	_, cache = setupDiskCache()
	recovered, ok := cache.GetBlob(id)

	if !ok {
		t.Fatal("The new cache instance did not find the image previously saved on disk.")
	}

	if !bytes.Equal(data, recovered) {
		t.Fatal("The cache returned different data than what was cached.")
	}
}

func setupDiskCache() (dp *DummyPolicy, cache *Cache) {
	dp = NewDummyPolicy(make(Log))
	cache = NewDiskCache(testutils.CacheFolder, 100, dp)
	return
}

func allTests(t *testing.T, f setupFunc) {
	testCache(t, f)
	testCacheExit(t, f)
}

const (
	Visit = iota
	Push
	Pop
)

type Log map[string][]uint

type DummyPolicy struct {
	fifo Policy

	loki Log
	pops int
}

func (d DummyPolicy) Visit(k string) {
	d.log(k, Visit)
	d.fifo.Visit(k)
}

func (d DummyPolicy) log(k string, t uint) {
	d.loki[k] = append(d.loki[k], t)
}

func (d DummyPolicy) Push(k string) {
	d.log(k, Push)
	d.fifo.Push(k)
}

func (d *DummyPolicy) Pop() string {
	d.pops += 1
	return d.fifo.Pop()
}

func NewDummyPolicy(log Log) *DummyPolicy {
	return &DummyPolicy{fifo: &FIFO{}, loki: log}
}

func setupMemcache() (dp *DummyPolicy, cache *Cache) {
	dp = NewDummyPolicy(make(Log))
	cache = NewCache(dp, 100)
	return
}

type setupFunc func() (dp *DummyPolicy, cache *Cache)

func testCache(t *testing.T, setup setupFunc) {
	id := string("testcache")
	dp, cache := setup()

	found := func() bool {
		_, ok := cache.GetBlob(id)
		return ok
	}

	if found() {
		t.Fatal("Cache claimed to contain a blob that was never added")
	}

	cache.AddBlob(id, make([]byte, 10))

	time.Sleep(100) // only necessary for pure disk cache

	if tx := dp.loki[id]; len(tx) != 1 || tx[0] != Push {
		t.Fatal("Cache did not use policy properly")
	}

	if !found() {
		t.Fatal("Not found after adding to cache")
	}

	if tx := dp.loki[id]; len(tx) != 2 || tx[1] != Visit {
		t.Fatal("Cache did not use policy properly")
	}
}

func testCacheExit(t *testing.T, setup setupFunc) {
	var (
		id1 = string("cacheexit1")
		id2 = string("cacheexit2")
		id3 = string("cacheexit3")
	)
	dp, cache := setup()

	cache.AddBlob(id1, make([]byte, 50))
	cache.AddBlob(id2, make([]byte, 40))
	cache.AddBlob(id3, make([]byte, 20))

	if dp.pops != 1 {
		t.Fatal("Wrong amount of blobs removed from cache")
	}
}
