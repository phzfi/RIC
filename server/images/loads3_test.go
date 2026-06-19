package images

import (
	"testing"

	cfg "github.com/phzfi/RIC/server/config"
)

func TestNewS3Client_WithEndpoint(t *testing.T) {
	testCfg := cfg.ImageSourceConfig{
		S3Enabled:  true,
		S3Bucket:   "test-bucket",
		S3Prefix:   "images/",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
	}

	client, err := NewS3Client(testCfg)
	if err != nil {
		t.Fatalf("NewS3Client should not fail with valid config: %v", err)
	}
	if client == nil {
		t.Fatal("Client should not be nil")
	}
	if client.bucket != "test-bucket" {
		t.Errorf("Expected bucket test-bucket, got %s", client.bucket)
	}
	if client.prefix != "images/" {
		t.Errorf("Expected prefix images/, got %s", client.prefix)
	}
}

func TestNewS3Client_WithoutEndpoint(t *testing.T) {
	testCfg := cfg.ImageSourceConfig{
		S3Enabled:  true,
		S3Bucket:   "my-bucket",
		S3Prefix:   "prefix/",
		S3Region:   "eu-west-1",
		S3Endpoint: "",
	}

	client, err := NewS3Client(testCfg)
	if err != nil {
		t.Fatalf("NewS3Client should not fail without endpoint: %v", err)
	}
	if client == nil {
		t.Fatal("Client should not be nil")
	}
	if client.bucket != "my-bucket" {
		t.Errorf("Expected bucket my-bucket, got %s", client.bucket)
	}
	if client.prefix != "prefix/" {
		t.Errorf("Expected prefix prefix/, got %s", client.prefix)
	}
}

func TestS3Client_GetObject_InvalidBucket(t *testing.T) {
	testCfg := cfg.ImageSourceConfig{
		S3Enabled:  true,
		S3Bucket:   "nonexistent-bucket-12345",
		S3Prefix:   "",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
	}

	client, err := NewS3Client(testCfg)
	if err != nil {
		t.Fatalf("NewS3Client should not fail: %v", err)
	}

	// GetObject should return an error for invalid bucket
	_, err = client.GetObject("test-key")
	if err == nil {
		t.Error("GetObject should return error for invalid bucket")
	}
}

func TestS3Client_HeadObject_InvalidBucket(t *testing.T) {
	testCfg := cfg.ImageSourceConfig{
		S3Enabled:  true,
		S3Bucket:   "nonexistent-bucket-12345",
		S3Prefix:   "",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
	}

	client, err := NewS3Client(testCfg)
	if err != nil {
		t.Fatalf("NewS3Client should not fail: %v", err)
	}

	// HeadObject should return (0, false, nil) for invalid bucket
	size, exists, err := client.HeadObject("test-key")
	if err != nil {
		t.Errorf("HeadObject should not return error: %v", err)
	}
	if exists {
		t.Error("HeadObject should return exists=false for invalid bucket")
	}
	if size != 0 {
		t.Error("HeadObject should return size=0")
	}
}

func TestS3Client_EmptyKey(t *testing.T) {
	testCfg := cfg.ImageSourceConfig{
		S3Enabled:  true,
		S3Bucket:   "nonexistent-bucket-12345",
		S3Prefix:   "",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
	}

	client, err := NewS3Client(testCfg)
	if err != nil {
		t.Fatalf("NewS3Client should not fail: %v", err)
	}

	// GetObject with empty key
	_, err = client.GetObject("")
	if err == nil {
		t.Error("GetObject should return error for empty key with invalid bucket")
	}
}

func TestS3Client_PrefixHandling(t *testing.T) {
	testCfg := cfg.ImageSourceConfig{
		S3Enabled:  true,
		S3Bucket:   "test-bucket",
		S3Prefix:   "folder/subfolder/",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
	}

	client, err := NewS3Client(testCfg)
	if err != nil {
		t.Fatalf("NewS3Client should not fail: %v", err)
	}

	// Verify prefix is stored correctly
	if client.prefix != "folder/subfolder/" {
		t.Errorf("Expected prefix folder/subfolder/, got %s", client.prefix)
	}
}