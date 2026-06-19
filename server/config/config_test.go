package config

import (
	"os"
	"testing"
)

func TestReadConfigDefaults(t *testing.T) {
	// Ensure no env vars interfere
	os.Unsetenv("WATERMARK_PATH")
	os.Unsetenv("WATERMARK_HORIZONTAL")
	os.Unsetenv("WATERMARK_VERTICAL")
	os.Unsetenv("WATERMARK_MAXWIDTH")
	os.Unsetenv("WATERMARK_MINWIDTH")
	os.Unsetenv("WATERMARK_MAXHEIGHT")
	os.Unsetenv("WATERMARK_MINHEIGHT")
	os.Unsetenv("WATERMARK_ADDMARK")
	os.Unsetenv("WATERMARK_FORCEMARK")
	os.Unsetenv("WATERMARK_TEXT")
	os.Unsetenv("SERVER_TOKENS")
	os.Unsetenv("SERVER_MEMORY")
	os.Unsetenv("CACHE_DISK_PATH")
	os.Unsetenv("CACHE_DISK_MAX_MB")
	os.Unsetenv("CACHE_S3_ENABLED")
	os.Unsetenv("CACHE_S3_BUCKET")
	os.Unsetenv("CACHE_S3_PREFIX")
	os.Unsetenv("CACHE_S3_REGION")
	os.Unsetenv("CACHE_S3_ENDPOINT")
	os.Unsetenv("CACHE_S3_MAX_MB")
	os.Unsetenv("SOURCE_S3_ENABLED")
	os.Unsetenv("SOURCE_S3_BUCKET")
	os.Unsetenv("SOURCE_S3_PREFIX")
	os.Unsetenv("SOURCE_S3_REGION")
	os.Unsetenv("SOURCE_S3_ENDPOINT")

	conf := ReadConfig()
	if *conf != defaults {
		t.Fatal("Expected default config")
	}
}

func TestReadConfigEnvVars(t *testing.T) {
	os.Setenv("WATERMARK_PATH", "/path/to/wm.png")
	os.Setenv("WATERMARK_HORIZONTAL", "0.3")
	os.Setenv("WATERMARK_VERTICAL", "0.7")
	os.Setenv("WATERMARK_MAXWIDTH", "3000")
	os.Setenv("WATERMARK_MINWIDTH", "100")
	os.Setenv("WATERMARK_MAXHEIGHT", "4000")
	os.Setenv("WATERMARK_MINHEIGHT", "150")
	os.Setenv("WATERMARK_ADDMARK", "true")
	os.Setenv("WATERMARK_FORCEMARK", "forced")
	os.Setenv("WATERMARK_TEXT", "Sample Text")
	os.Setenv("SERVER_TOKENS", "5")
	os.Setenv("SERVER_MEMORY", "1073741824")

	defer func() {
		os.Unsetenv("WATERMARK_PATH")
		os.Unsetenv("WATERMARK_HORIZONTAL")
		os.Unsetenv("WATERMARK_VERTICAL")
		os.Unsetenv("WATERMARK_MAXWIDTH")
		os.Unsetenv("WATERMARK_MINWIDTH")
		os.Unsetenv("WATERMARK_MAXHEIGHT")
		os.Unsetenv("WATERMARK_MINHEIGHT")
		os.Unsetenv("WATERMARK_ADDMARK")
		os.Unsetenv("WATERMARK_FORCEMARK")
		os.Unsetenv("WATERMARK_TEXT")
		os.Unsetenv("SERVER_TOKENS")
		os.Unsetenv("SERVER_MEMORY")
	}()

	conf := ReadConfig()
	if conf.Watermark.ImagePath != "/path/to/wm.png" {
		t.Fatalf("Expected ImagePath /path/to/wm.png, got %s", conf.Watermark.ImagePath)
	}
	if conf.Watermark.Horizontal != 0.3 {
		t.Fatalf("Expected Horizontal 0.3, got %f", conf.Watermark.Horizontal)
	}
	if conf.Watermark.Vertical != 0.7 {
		t.Fatalf("Expected Vertical 0.7, got %f", conf.Watermark.Vertical)
	}
	if conf.Watermark.MaxWidth != 3000 {
		t.Fatalf("Expected MaxWidth 3000, got %d", conf.Watermark.MaxWidth)
	}
	if conf.Watermark.MinWidth != 100 {
		t.Fatalf("Expected MinWidth 100, got %d", conf.Watermark.MinWidth)
	}
	if conf.Watermark.MaxHeight != 4000 {
		t.Fatalf("Expected MaxHeight 4000, got %d", conf.Watermark.MaxHeight)
	}
	if conf.Watermark.MinHeight != 150 {
		t.Fatalf("Expected MinHeight 150, got %d", conf.Watermark.MinHeight)
	}
	if !conf.Watermark.AddMark {
		t.Fatal("Expected AddMark true")
	}
	if conf.Watermark.ForceMark != "forced" {
		t.Fatalf("Expected ForceMark 'forced', got %s", conf.Watermark.ForceMark)
	}
	if conf.Watermark.Text != "Sample Text" {
		t.Fatalf("Expected Text 'Sample Text', got %s", conf.Watermark.Text)
	}
	if conf.Server.Tokens != 5 {
		t.Fatalf("Expected Tokens 5, got %d", conf.Server.Tokens)
	}
	if conf.Server.Memory != 1073741824 {
		t.Fatalf("Expected Memory 1073741824, got %d", conf.Server.Memory)
	}
}

func TestReadConfigInvalidEnvVars(t *testing.T) {
	os.Setenv("WATERMARK_HORIZONTAL", "not_a_number")
	os.Setenv("WATERMARK_VERTICAL", "not_a_number")
	os.Setenv("WATERMARK_MAXWIDTH", "not_a_number")
	os.Setenv("WATERMARK_MINWIDTH", "not_a_number")
	os.Setenv("WATERMARK_MAXHEIGHT", "not_a_number")
	os.Setenv("WATERMARK_MINHEIGHT", "not_a_number")
	os.Setenv("WATERMARK_ADDMARK", "not_a_bool")
	os.Setenv("SERVER_TOKENS", "not_a_number")
	os.Setenv("SERVER_MEMORY", "not_a_number")

	defer func() {
		os.Unsetenv("WATERMARK_HORIZONTAL")
		os.Unsetenv("WATERMARK_VERTICAL")
		os.Unsetenv("WATERMARK_MAXWIDTH")
		os.Unsetenv("WATERMARK_MINWIDTH")
		os.Unsetenv("WATERMARK_MAXHEIGHT")
		os.Unsetenv("WATERMARK_MINHEIGHT")
		os.Unsetenv("WATERMARK_ADDMARK")
		os.Unsetenv("SERVER_TOKENS")
		os.Unsetenv("SERVER_MEMORY")
	}()

	conf := ReadConfig()
	// Should fall back to defaults
	if conf.Watermark.Horizontal != defaults.Watermark.Horizontal {
		t.Fatalf("Expected Horizontal %f (default), got %f", defaults.Watermark.Horizontal, conf.Watermark.Horizontal)
	}
	if conf.Watermark.Vertical != defaults.Watermark.Vertical {
		t.Fatalf("Expected Vertical %f (default), got %f", defaults.Watermark.Vertical, conf.Watermark.Vertical)
	}
	if conf.Watermark.MaxWidth != defaults.Watermark.MaxWidth {
		t.Fatalf("Expected MaxWidth %d (default), got %d", defaults.Watermark.MaxWidth, conf.Watermark.MaxWidth)
	}
	if conf.Server.Tokens != defaults.Server.Tokens {
		t.Fatalf("Expected Tokens %d (default), got %d", defaults.Server.Tokens, conf.Server.Tokens)
	}
	if conf.Server.Memory != defaults.Server.Memory {
		t.Fatalf("Expected Memory %d (default), got %d", defaults.Server.Memory, conf.Server.Memory)
	}
}

func TestReadConfigCacheEnvVars(t *testing.T) {
	os.Setenv("CACHE_DISK_PATH", "/custom/cache")
	os.Setenv("CACHE_DISK_MAX_MB", "8192")
	os.Setenv("CACHE_S3_ENABLED", "true")
	os.Setenv("CACHE_S3_BUCKET", "my-cache-bucket")
	os.Setenv("CACHE_S3_PREFIX", "cache/")
	os.Setenv("CACHE_S3_REGION", "us-east-1")
	os.Setenv("CACHE_S3_ENDPOINT", "http://minio:9000")
	os.Setenv("CACHE_S3_MAX_MB", "500")

	defer func() {
		os.Unsetenv("CACHE_DISK_PATH")
		os.Unsetenv("CACHE_DISK_MAX_MB")
		os.Unsetenv("CACHE_S3_ENABLED")
		os.Unsetenv("CACHE_S3_BUCKET")
		os.Unsetenv("CACHE_S3_PREFIX")
		os.Unsetenv("CACHE_S3_REGION")
		os.Unsetenv("CACHE_S3_ENDPOINT")
		os.Unsetenv("CACHE_S3_MAX_MB")
	}()

	conf := ReadConfig()
	if conf.Cache.DiskPath != "/custom/cache" {
		t.Fatalf("Expected DiskPath /custom/cache, got %s", conf.Cache.DiskPath)
	}
	if conf.Cache.DiskMaxMB != 8192 {
		t.Fatalf("Expected DiskMaxMB 8192, got %d", conf.Cache.DiskMaxMB)
	}
	if !conf.Cache.S3Enabled {
		t.Fatal("Expected S3Enabled true")
	}
	if conf.Cache.S3Bucket != "my-cache-bucket" {
		t.Fatalf("Expected S3Bucket my-cache-bucket, got %s", conf.Cache.S3Bucket)
	}
	if conf.Cache.S3Prefix != "cache/" {
		t.Fatalf("Expected S3Prefix cache/, got %s", conf.Cache.S3Prefix)
	}
	if conf.Cache.S3Region != "us-east-1" {
		t.Fatalf("Expected S3Region us-east-1, got %s", conf.Cache.S3Region)
	}
	if conf.Cache.S3Endpoint != "http://minio:9000" {
		t.Fatalf("Expected S3Endpoint http://minio:9000, got %s", conf.Cache.S3Endpoint)
	}
	if conf.Cache.S3MaxMB != 500 {
		t.Fatalf("Expected S3MaxMB 500, got %d", conf.Cache.S3MaxMB)
	}
}

func TestReadConfigCacheInvalidEnvVars(t *testing.T) {
	os.Setenv("CACHE_DISK_MAX_MB", "not_a_number")
	os.Setenv("CACHE_S3_ENABLED", "not_a_bool")
	os.Setenv("CACHE_S3_MAX_MB", "also_not_number")

	defer func() {
		os.Unsetenv("CACHE_DISK_MAX_MB")
		os.Unsetenv("CACHE_S3_ENABLED")
		os.Unsetenv("CACHE_S3_MAX_MB")
	}()

	conf := ReadConfig()
	if conf.Cache.DiskMaxMB != defaults.Cache.DiskMaxMB {
		t.Fatalf("Expected DiskMaxMB %d (default), got %d", defaults.Cache.DiskMaxMB, conf.Cache.DiskMaxMB)
	}
	if conf.Cache.S3Enabled != defaults.Cache.S3Enabled {
		t.Fatalf("Expected S3Enabled %v (default), got %v", defaults.Cache.S3Enabled, conf.Cache.S3Enabled)
	}
	if conf.Cache.S3MaxMB != defaults.Cache.S3MaxMB {
		t.Fatalf("Expected S3MaxMB %d (default), got %d", defaults.Cache.S3MaxMB, conf.Cache.S3MaxMB)
	}
}

func TestReadConfigImageSourceEnvVars(t *testing.T) {
	os.Setenv("SOURCE_S3_ENABLED", "true")
	os.Setenv("SOURCE_S3_BUCKET", "my-source-bucket")
	os.Setenv("SOURCE_S3_PREFIX", "images/")
	os.Setenv("SOURCE_S3_REGION", "eu-west-1")
	os.Setenv("SOURCE_S3_ENDPOINT", "http://s3.example.com")

	defer func() {
		os.Unsetenv("SOURCE_S3_ENABLED")
		os.Unsetenv("SOURCE_S3_BUCKET")
		os.Unsetenv("SOURCE_S3_PREFIX")
		os.Unsetenv("SOURCE_S3_REGION")
		os.Unsetenv("SOURCE_S3_ENDPOINT")
	}()

	conf := ReadConfig()
	if !conf.ImageSource.S3Enabled {
		t.Fatal("Expected S3Enabled true")
	}
	if conf.ImageSource.S3Bucket != "my-source-bucket" {
		t.Fatalf("Expected S3Bucket my-source-bucket, got %s", conf.ImageSource.S3Bucket)
	}
	if conf.ImageSource.S3Prefix != "images/" {
		t.Fatalf("Expected S3Prefix images/, got %s", conf.ImageSource.S3Prefix)
	}
	if conf.ImageSource.S3Region != "eu-west-1" {
		t.Fatalf("Expected S3Region eu-west-1, got %s", conf.ImageSource.S3Region)
	}
	if conf.ImageSource.S3Endpoint != "http://s3.example.com" {
		t.Fatalf("Expected S3Endpoint http://s3.example.com, got %s", conf.ImageSource.S3Endpoint)
	}
}

func TestReadConfigImageSourceInvalidEnvVars(t *testing.T) {
	os.Setenv("SOURCE_S3_ENABLED", "not_a_bool")

	defer func() {
		os.Unsetenv("SOURCE_S3_ENABLED")
	}()

	conf := ReadConfig()
	if conf.ImageSource.S3Enabled != defaults.ImageSource.S3Enabled {
		t.Fatalf("Expected S3Enabled %v (default), got %v", defaults.ImageSource.S3Enabled, conf.ImageSource.S3Enabled)
	}
}
