package cache

import (
	"errors"
	"github.com/phzfi/RIC/server/images"
	"path/filepath"
	"testing"
)

type DummyResizer bool

func (d DummyResizer) GetImage(fn string, w, h uint) (images.ImageBlob, error) {
	if d {
		if fn == "toserve.jpg" {
			return make(images.ImageBlob, 1001), nil
		} else {
			return nil, errors.New("That is not the name of my image!")
		}
	} else {
		return nil, errors.New("RemoveRoot was called last, so I assume I shouldn't find my image.")
	}
}

func (d DummyResizer) ImageSize(_ string) (uint, uint, error) {
	return 1000, 1000, nil
}

func (d *DummyResizer) AddRoot(_ string) error {
	*d = true
	return nil
}

func (d *DummyResizer) RemoveRoot(_ string) error {
	*d = false
	return nil
}

func makeCacheMaker(policyMaker func() Policy) CacheMaker {
	return func(mm uint64) *Cache {
		return NewCache(new(DummyResizer), policyMaker(), mm)
	}
}

func TestLRU(t *testing.T) {
	RunTests(makeCacheMaker(func() Policy { return NewLRUPolicy() }), t)
}

func TestFIFO(t *testing.T) {
	RunTests(makeCacheMaker(func() Policy { return new(FIFO) }), t)
}

type T struct {
	*testing.T
}

func (t T) FatalIfError(err error) {
	if err != nil {
		t.Fatal(err)
	}
}

var path string

func init() {
	path = filepath.FromSlash("../testimages/cache")
}

type CacheMaker func(uint64) *Cache

func RunTests(newResizer CacheMaker, t_in *testing.T) {

	t := T{t_in}

	normalOperation(t, newResizer)
	cache(t, newResizer)
	errorTest(t, newResizer)
	cacheExit(t, newResizer)
	noMemory(t, newResizer)
}

func normalOperation(t T, newResizer CacheMaker) {
	// 500MB cache
	cache := newResizer(500 * 1024 * 1024)

	t.FatalIfError(cache.AddRoot(path))

	_, err := cache.GetImage("toserve.jpg", 10, 10)
	t.FatalIfError(err)

	_, err = cache.GetImage("toserve.jpg", 104, 10)
	t.FatalIfError(err)

	_, err = cache.GetImage("toserve.jpg", 50, 200)
	t.FatalIfError(err)

	_, err = cache.GetImage("toserve.jpg", 100, 10)
	t.FatalIfError(err)

	_, err = cache.GetImage("toserve.jpg", 111, 111)
	t.FatalIfError(err)
}

func cache(t T, newResizer CacheMaker) {
	// 500MB cache
	cache := newResizer(500 * 1024 * 1024)

	t.FatalIfError(cache.AddRoot(path))

	_, err := cache.GetImage("toserve.jpg", 10, 10)
	t.FatalIfError(err)

	_, err = cache.GetImage("toserve.jpg", 104, 10)
	t.FatalIfError(err)

	// root is removed to verify that the image is not read from disk
	t.FatalIfError(cache.RemoveRoot(path))

	_, err = cache.GetImage("toserve.jpg", 10, 10)
	t.FatalIfError(err)
}

func errorTest(t T, newResizer CacheMaker) {
	// 500MB cache
	cache := newResizer(500 * 1024 * 1024)

	t.FatalIfError(cache.AddRoot(path))

	_, err := cache.GetImage("tosslntvgerve.jpg", 10, 10)
	if err == nil {
		t.Fatal("No error, although querying nonexistent image.")
	}
}

func cacheExit(t T, newResizer CacheMaker) {
	// 50kB cache
	cache := newResizer(4 * 1024)

	t.FatalIfError(cache.AddRoot(path))

	for i := 0; i < 6; i++ {
		_, err := cache.GetImage("toserve.jpg", uint(100-i), 100)
		t.FatalIfError(err)
	}

	t.FatalIfError(cache.RemoveRoot(path))

	_, err := cache.GetImage("toserve.jpg", 100, 100)
	if err == nil {
		t.Fatal("No error, although querying image that should not be in cache.")
	}
}

func noMemory(t T, newResizer CacheMaker) {
	cache := newResizer(0)

	t.FatalIfError(cache.AddRoot(path))

	_, err := cache.GetImage("toserve.jpg", 10, 10)
	t.FatalIfError(err)

	_, err = cache.GetImage("toserve.jpg", 104, 10)
	t.FatalIfError(err)

	_, err = cache.GetImage("toserve.jpg", 10, 10)
	t.FatalIfError(err)
}
