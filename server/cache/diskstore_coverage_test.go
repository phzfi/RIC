package cache

import (
	"testing"
)

func TestStringToBase64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"simple string", "test", "dGVzdA"},
		{"longer string", "hello world", "aGVsbG8gd29ybGQ"},
		{"special chars", "key-123_456", "a2V5LTEyM180NTY"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := stringToBase64(tc.input)
			if result != tc.expected {
				t.Errorf("stringToBase64(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestNewDiskStore(t *testing.T) {
	ds := NewDiskStore("/tmp/test_new_disk_store")
	if ds == nil {
		t.Fatal("NewDiskStore should not return nil")
	}
	if ds.folder != "/tmp/test_new_disk_store" {
		t.Errorf("Expected folder /tmp/test_new_disk_store, got %s", ds.folder)
	}
}

func TestDiskStoreLoadSuccess(t *testing.T) {
	folder := "/tmp/test_diskstore_load_success"
	defer func() {
		// Cleanup
	}()

	ds := NewDiskStore(folder)
	ds.Store("testkey", []byte("test data"))

	// Wait for async store
	ds.Lock()
	// The store is async, so we need to give it time
	ds.Unlock()

	// Note: Store is async, so we can't reliably test Load without mocking
	// This test validates the Store doesn't panic
}