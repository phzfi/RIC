package cache;

import (
	"os"
	"testing"
	"time"
)

// TestNewDiskCache tests the NewDiskCache function thoroughly
func TestNewDiskCache(t *testing.T) {
	// Test with new directory
	dp := NewDummyPolicy()
	cache := NewDiskCache("/tmp/test_diskcache_newdir", 1024*1024*10, dp)
	if cache == nil {
		t.Fatal("NewDiskCache returned nil")
	}

	// Add blob and verify
	cache.AddBlob("key1", []byte("data1"))
	time.Sleep(100 * time.Millisecond)
	
	data, found := cache.GetBlob("key1")
	if !found {
		t.Error("Blob should be found")
	}
	if string(data) != "data1" {
		t.Errorf("Expected 'data1', got '%s'", string(data))
	}

	// Test loading existing cache from disk
	dp2 := NewDummyPolicy()
	cache2 := NewDiskCache("/tmp/test_diskcache_newdir", 1024*1024*10, dp2)
	data2, found2 := cache2.GetBlob("key1")
	if !found2 {
		t.Error("Blob should be found in reinitialized cache")
	}
	if string(data2) != "data1" {
		t.Errorf("Expected 'data1', got '%s'", string(data2))
	}
}

// TestFileSize tests the fileSize function
func TestFileSize(t *testing.T) {
	// Create a test file
	testData := []byte("test data for filesize")
	err := os.WriteFile("/tmp/test_filesize.txt", testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	size := fileSize("/tmp/test_filesize.txt")
	if size != uint64(len(testData)) {
		t.Errorf("Expected file size %d, got %d", len(testData), size)
	}

	// Test with non-existent file
	size = fileSize("/tmp/nonexistent_file_12345")
	if size != 0 {
		t.Errorf("Expected 0 for non-existent file, got %d", size)
	}
}

// TestFIFOPolicy_Visit tests the Visit method
func TestFIFOPolicy_Visit(t *testing.T) {
	dp := NewDummyPolicy()

	// Add some entries
	dp.Push("key1")
	dp.Push("key2")
	dp.Push("key3")

	// Visit should not panic
	dp.Visit("key1")
	dp.Visit("key2")
	dp.Visit("key3")

	// Check that visits were recorded
	if len(dp.loki["key1"]) == 0 {
		t.Error("Visit should record visits")
	}
	if len(dp.loki["key1"]) > 0 && dp.loki["key1"][len(dp.loki["key1"])-1] != Visit {
		t.Error("Last operation should be Visit")
	}
}
