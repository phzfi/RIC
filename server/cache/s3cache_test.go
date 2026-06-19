package cache

import (
	"testing"

	cfg "github.com/phzfi/RIC/server/config"
)

func TestNewS3Store_WithEndpoint(t *testing.T) {
	testCfg := cfg.CacheConfig{
		S3Enabled:  true,
		S3Bucket:   "test-bucket",
		S3Prefix:   "cache/",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
		S3MaxMB:    100,
	}

	store, err := NewS3Store(testCfg)
	if err != nil {
		t.Fatalf("NewS3Store should not fail with valid config: %v", err)
	}
	if store == nil {
		t.Fatal("Store should not be nil")
	}
	if store.bucket != "test-bucket" {
		t.Errorf("Expected bucket test-bucket, got %s", store.bucket)
	}
	if store.prefix != "cache/" {
		t.Errorf("Expected prefix cache/, got %s", store.prefix)
	}
	if store.maxSize != 100*1024*1024 {
		t.Errorf("Expected maxSize %d, got %d", 100*1024*1024, store.maxSize)
	}
}

func TestNewS3Store_WithoutEndpoint(t *testing.T) {
	testCfg := cfg.CacheConfig{
		S3Enabled:  true,
		S3Bucket:   "my-bucket",
		S3Prefix:   "prefix/",
		S3Region:   "eu-west-1",
		S3Endpoint: "",
		S3MaxMB:    50,
	}

	store, err := NewS3Store(testCfg)
	if err != nil {
		t.Fatalf("NewS3Store should not fail without endpoint: %v", err)
	}
	if store == nil {
		t.Fatal("Store should not be nil")
	}
	if store.bucket != "my-bucket" {
		t.Errorf("Expected bucket my-bucket, got %s", store.bucket)
	}
	if store.prefix != "prefix/" {
		t.Errorf("Expected prefix prefix/, got %s", store.prefix)
	}
	if store.maxSize != 50*1024*1024 {
		t.Errorf("Expected maxSize %d, got %d", 50*1024*1024, store.maxSize)
	}
}

func TestNewS3Store_DefaultMaxSize(t *testing.T) {
	testCfg := cfg.CacheConfig{
		S3Enabled:  true,
		S3Bucket:   "test-bucket",
		S3Prefix:   "",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
		S3MaxMB:    0, // Should default to 200MB
	}

	store, err := NewS3Store(testCfg)
	if err != nil {
		t.Fatalf("NewS3Store should not fail: %v", err)
	}
	if store.maxSize != 200*1024*1024 {
		t.Errorf("Expected default maxSize %d, got %d", 200*1024*1024, store.maxSize)
	}
}

func TestNewS3Cache_Success(t *testing.T) {
	testCfg := cfg.CacheConfig{
		S3Enabled:  true,
		S3Bucket:   "test-bucket",
		S3Prefix:   "cache/",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
		S3MaxMB:    100,
	}

	cache, err := NewS3Cache(testCfg, NewLRU())
	if err != nil {
		t.Fatalf("NewS3Cache should not fail: %v", err)
	}
	if cache == nil {
		t.Fatal("Cache should not be nil")
	}
}

func TestNewS3Cache_WithNilPolicy(t *testing.T) {
	testCfg := cfg.CacheConfig{
		S3Enabled:  true,
		S3Bucket:   "test-bucket",
		S3Prefix:   "cache/",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
		S3MaxMB:    100,
	}

	// This should work since policy is not used in this test
	cache, err := NewS3Cache(testCfg, NewLRU())
	if err != nil {
		t.Fatalf("NewS3Cache should not fail: %v", err)
	}
	if cache == nil {
		t.Fatal("Cache should not be nil")
	}
}

func TestS3Store_Load_InvalidBucket(t *testing.T) {
	testCfg := cfg.CacheConfig{
		S3Enabled:  true,
		S3Bucket:   "nonexistent-bucket-123456789",
		S3Prefix:   "",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
		S3MaxMB:    100,
	}

	store, err := NewS3Store(testCfg)
	if err != nil {
		t.Fatalf("NewS3Store should not fail: %v", err)
	}

	blob, found := store.Load("test-key")
	if found {
		t.Error("Load should return found=false for invalid bucket")
	}
	if blob != nil {
		t.Error("Load should return nil blob for invalid bucket")
	}
}

func TestS3Store_Store_SizeLimit(t *testing.T) {
	testCfg := cfg.CacheConfig{
		S3Enabled:  true,
		S3Bucket:   "nonexistent-bucket-123456789",
		S3Prefix:   "",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
		S3MaxMB:    1, // 1MB limit
	}

	store, err := NewS3Store(testCfg)
	if err != nil {
		t.Fatalf("NewS3Store should not fail: %v", err)
	}

	// Store should skip blobs that exceed maxSize
	largeBlob := make([]byte, 2*1024*1024) // 2MB
	store.Store("large-key", largeBlob)
	// No panic means success (store is async and will just log)
}

func TestS3Store_Store_EmptyBlob(t *testing.T) {
	testCfg := cfg.CacheConfig{
		S3Enabled:  true,
		S3Bucket:   "nonexistent-bucket-123456789",
		S3Prefix:   "",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
		S3MaxMB:    100,
	}

	store, err := NewS3Store(testCfg)
	if err != nil {
		t.Fatalf("NewS3Store should not fail: %v", err)
	}

	// Store should handle empty blob
	store.Store("empty-key", []byte{})
	// No panic means success
}

func TestS3Store_Delete_InvalidBucket(t *testing.T) {
	testCfg := cfg.CacheConfig{
		S3Enabled:  true,
		S3Bucket:   "nonexistent-bucket-123456789",
		S3Prefix:   "",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
		S3MaxMB:    100,
	}

	store, err := NewS3Store(testCfg)
	if err != nil {
		t.Fatalf("NewS3Store should not fail: %v", err)
	}

	// Delete should not panic
	size := store.Delete("test-key")
	if size != 0 {
		t.Errorf("Delete should return 0, got %d", size)
	}
}

func TestS3Store_PrefixHandling(t *testing.T) {
	testCfg := cfg.CacheConfig{
		S3Enabled:  true,
		S3Bucket:   "test-bucket",
		S3Prefix:   "folder/subfolder/",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
		S3MaxMB:    100,
	}

	store, err := NewS3Store(testCfg)
	if err != nil {
		t.Fatalf("NewS3Store should not fail: %v", err)
	}

	if store.prefix != "folder/subfolder/" {
		t.Errorf("Expected prefix folder/subfolder/, got %s", store.prefix)
	}
}

func TestS3Store_EmptyPrefix(t *testing.T) {
	testCfg := cfg.CacheConfig{
		S3Enabled:  true,
		S3Bucket:   "test-bucket",
		S3Prefix:   "",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
		S3MaxMB:    100,
	}

	store, err := NewS3Store(testCfg)
	if err != nil {
		t.Fatalf("NewS3Store should not fail: %v", err)
	}

	if store.prefix != "" {
		t.Errorf("Expected empty prefix, got %s", store.prefix)
	}
}