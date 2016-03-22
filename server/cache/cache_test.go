package cache

import (
	"bytes"
	"github.com/phzfi/RIC/server/testutils"
	"testing"
	"time"
)

const (
	cachefolder = "/tmp/cachetests"
	cacheSize   = 100
)

func TestMemCache(t *testing.T) {
	allTests(t, setupMemcache)
}

func TestDiskCache(t *testing.T) {
	allTests(t, func() (*DummyPolicy, *Cache) {
		testutils.RemoveContents(cachefolder)
		return setupDiskCache()
	})
}

func TestDiskCachePersistence(t *testing.T) {
	id := "testdiskpersist"
	data := []byte{1, 2, 3, 4, 7}

	_, cache := setupDiskCache()
	cache.AddBlob(id, data)

	time.Sleep(100 * time.Millisecond)

	_, cache = setupDiskCache()
	recovered, ok := cache.GetBlob(id)

	if !ok {
		t.Fatal("The new cache instance did not find the image previously saved on disk.")
	}

	if !bytes.Equal(data, recovered) {
		t.Fatal("The cache returned different data than what was cached.")
	}
}

type setupFunc func() (dp *DummyPolicy, cache *Cache)

func setupDiskCache() (dp *DummyPolicy, cache *Cache) {
	dp = NewDummyPolicy(make(Log))
	cache = NewDiskCache(cachefolder, cacheSize, dp)
	return
}

func setupMemcache() (dp *DummyPolicy, cache *Cache) {
	dp = NewDummyPolicy(make(Log))
	cache = NewCache(dp, cacheSize)
	return
}

func allTests(t *testing.T, f setupFunc) {
	testCache(t, f)
	testCacheExit(t, f)
}

func testCache(t *testing.T, setup setupFunc) {
	id := "testcache"
	dp, cache := setup()

	found := func() bool {
		_, ok := cache.GetBlob(id)
		return ok
	}

	if found() {
		t.Fatal("Cache claimed to contain a blob that was never added")
	}

	cache.AddBlob(id, make([]byte, 10))

	time.Sleep(100 * time.Millisecond) // only necessary for pure disk cache

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
		id1 = "cacheexit1"
		id2 = "cacheexit2"
		id3 = "cacheexit3"
	)
	dp, cache := setup()

	cache.AddBlob(id1, make([]byte, 50))
	cache.AddBlob(id2, make([]byte, 40))
	cache.AddBlob(id3, make([]byte, 20))

	if dp.pops != 1 {
		t.Fatal("Wrong amount of blobs removed from cache")
	}
}

func TestTooBig(t *testing.T) {
	dp, cache := setupMemcache()
	const id = "string"
	cache.AddBlob(id, make([]byte, cacheSize+1))

	if len(dp.loki[id]) != 0 {
		t.Fatalf("Despite being too big, resource was cached. %#v", dp.loki[id])
	}
}
