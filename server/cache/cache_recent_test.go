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

func TestNormalOperation(t_in *testing.T) {
	t := T{t_in}

	// 500MB cache
	cache := NewCacherecent(500 * 1024 * 1024)

	t.FatalIfError(cache.AddRoot(path))

	_, err := cache.GetImage("toserve.jpg", 10, 10)
	t.FatalIfError(err)

	_, err = cache.GetImage("toserve.jpg", 104, 10)
	t.FatalIfError(err)

	_, err = cache.GetImage("toserve.jpg", 10, 10)
	t.FatalIfError(err)
}

func TestCache(t_in *testing.T) {
	t := T{t_in}

	// 500MB cache
	cache := NewCacherecent(500 * 1024 * 1024)

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

func TestError(t_in *testing.T) {
	t := T{t_in}

	// 500MB cache
	cache := NewCacherecent(500 * 1024 * 1024)

	t.FatalIfError(cache.AddRoot(path))

	_, err := cache.GetImage("tosslntvgerve.jpg", 10, 10)
	if err == nil {
		t.Fatal("No error, although querying nonexistent image.")
	}
}

func TestCacheExit(t_in *testing.T) {
	t := T{t_in}

	// 50kB cache
	cache := NewCacherecent(50 * 1024)

	t.FatalIfError(cache.AddRoot(path))

	for i := 0; i < 6; i++{
		_, err := cache.GetImage("toserve.jpg", uint(100 - i), 100)
		t.FatalIfError(err)
		println(cache.currentMemory)
	}

	t.FatalIfError(cache.RemoveRoot(path))

	_, err := cache.GetImage("toserve.jpg", 100, 100)
	if err == nil {
		t.Fatal("No error, although querying image that should not be in cache.")
	}
}

func TestNoMemory(t_in *testing.T) {
	t := T{t_in}

	cache := NewCacherecent(0)

	t.FatalIfError(cache.AddRoot(path))

	_, err := cache.GetImage("toserve.jpg", 10, 10)
	t.FatalIfError(err)

	_, err = cache.GetImage("toserve.jpg", 104, 10)
	t.FatalIfError(err)

	_, err = cache.GetImage("toserve.jpg", 10, 10)
	t.FatalIfError(err)
}