package cache

import (
	"path/filepath"
	"testing"
)

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

type CacheMaker func(uint64) ImageCache

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
	cache := newResizer(50 * 1024)

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
