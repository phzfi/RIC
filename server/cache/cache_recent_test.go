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

	// 550MB cache
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

	// 550MB cache
	cache := NewCacherecent(500 * 1024 * 1024)

	t.FatalIfError(cache.AddRoot(path))

	_, err := cache.GetImage("toserve.jpg", 10, 10)
	t.FatalIfError(err)

	_, err = cache.GetImage("toserve.jpg", 104, 10)
	t.FatalIfError(err)

	t.FatalIfError(cache.RemoveRoot(path))

	_, err = cache.GetImage("toserve.jpg", 10, 10)
	t.FatalIfError(err)
}

func TestError(t_in *testing.T) {
	t := T{t_in}

	// 550MB cache
	cache := NewCacherecent(500 * 1024 * 1024)

	t.FatalIfError(cache.AddRoot(path))

	_, err := cache.GetImage("tosslntvgerve.jpg", 10, 10)
	if err == nil {
		t.Fatal("No error, although querying nonexistent image.")
	}
}
