package ops

import (
	"os"
	"testing"

	"github.com/phzfi/RIC/server/images"
)

func TestIsWebroot(t *testing.T) {
	tests := []struct {
		root     string
		expected bool
	}{
		{"http://example.com/", true},
		{"https://example.com/", true},
		{"http:", true},
		{"https:", true},
		{"/path/to/root", false},
		{"../relative", false},
		{"/absolute/path", false},
		{"", false},
	}

	for _, tc := range tests {
		result := isWebroot(tc.root)
		if result != tc.expected {
			t.Errorf("isWebroot(%q) = %v, expected %v", tc.root, result, tc.expected)
		}
	}
}

func TestAddRoot(t *testing.T) {
	is := MakeImageSource()

	err := is.AddRoot("testimages/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	err = is.AddRoot("testimages/")
	if err != ErrRootAlreadyAdded {
		t.Fatalf("Expected ErrRootAlreadyAdded, got: %v", err)
	}
}

func TestRemoveRoot(t *testing.T) {
	is := MakeImageSource()

	err := is.AddRoot("testimages/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	err = is.RemoveRoot("testimages/")
	if err != nil {
		t.Fatalf("RemoveRoot failed: %v", err)
	}

	err = is.RemoveRoot("testimages/")
	if err != ErrRootNotFound {
		t.Fatalf("Expected ErrRootNotFound, got: %v", err)
	}
}

func TestSearchRootsEmpty(t *testing.T) {
	is := MakeImageSource()

	img := images.NewImage()
	defer img.Destroy()

	err := is.searchRoots("test.jpg", img)
	if err != os.ErrNotExist {
		t.Fatalf("Expected os.ErrNotExist for empty roots, got: %v", err)
	}
}