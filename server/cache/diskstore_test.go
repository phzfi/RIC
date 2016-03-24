package cache

import (
	"os"
	"testing"
	"time"
)

// The value of these tests is questionable, but they raise coverage.

func TestBreakPath(t *testing.T) {
	const folder = "/tmp/tobebroken"
	ds := Storer(NewDiskStore(folder))
	ds.Store("abc", []byte{1, 1, 2, 0})

	time.Sleep(100 * time.Millisecond)

	os.RemoveAll(folder)

	_, found := ds.Load("abc")
	if found {
		t.Fatal("Found in cache, although file was deleted.")
	}

	ds.Delete("abc")
}
