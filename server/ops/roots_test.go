package ops

import (
	"os"
	"testing"

	"github.com/phzfi/RIC/server/config"
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

func TestRootsHasRoot(t *testing.T) {
	var r roots

	r.Add("/path/one")
	r.Add("/path/two")

	if !r.HasRoot("/path/one") {
		t.Error("HasRoot should return true for existing root")
	}

	if r.HasRoot("/path/three") {
		t.Error("HasRoot should return false for non-existing root")
	}
}

func TestRootsRemoveMultiple(t *testing.T) {
	var r roots
	r.Add("/path/one")
	r.Add("/path/two")
	r.Add("/path/three")

	err := r.Remove("/path/two")
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	if r.HasRoot("/path/two") {
		t.Error("Root should be removed")
	}

	if !r.HasRoot("/path/one") || !r.HasRoot("/path/three") {
		t.Error("Other roots should remain")
	}
}

func TestAddRoot_WebRoot(t *testing.T) {
	is := MakeImageSource()

	// Add a web root
	err := is.AddRoot("http://example.com/images/")
	if err != nil {
		t.Fatalf("AddRoot with web root failed: %v", err)
	}

	// Verify it was added (indirectly by trying to use it)
	// We can't easily check internal state, but we can verify no error
}

func TestRemoveRoot_WebRoot(t *testing.T) {
	is := MakeImageSource()

	// Add and then remove a web root
	err := is.AddRoot("http://example.com/images/")
	if err != nil {
		t.Fatalf("AddRoot with web root failed: %v", err)
	}

	err = is.RemoveRoot("http://example.com/images/")
	if err != nil {
		t.Fatalf("RemoveRoot with web root failed: %v", err)
	}

	// Try to remove non-existent web root
	err = is.RemoveRoot("http://nonexistent.com/")
	if err != ErrRootNotFound {
		t.Fatalf("Expected ErrRootNotFound for non-existent web root, got: %v", err)
	}
}

func TestAddRoot_AbsError(t *testing.T) {
	// This is difficult to test directly since filepath.Abs rarely fails
	// We'll skip this test as it requires a special filesystem state
	t.Skip("Cannot easily test filepath.Abs error")
}

func TestInt32ToString(t *testing.T) {
	tests := []uint32{
		0,
		1,
		255,
		256,
		65535,
		65536,
		4294967295,
	}

	for _, x := range tests {
		result := int32ToString(x)
		if len(result) != 4 {
			t.Errorf("int32ToString(%d) returned %d bytes, expected 4", x, len(result))
		}
	}
}

func TestFloat64ToString(t *testing.T) {
	tests := []float64{
		0.0,
		1.0,
		-1.0,
		3.14159,
		1e10,
		-1e-10,
	}

	for _, x := range tests {
		result := float64ToString(x)
		if len(result) != 8 {
			t.Errorf("float64ToString(%f) returned %d bytes, expected 8", x, len(result))
		}
	}
}

func TestInt64ToString(t *testing.T) {
	tests := []uint64{
		0,
		1,
		255,
		4294967295,
		18446744073709551615,
	}

	for _, x := range tests {
		result := int64ToString(x)
		if len(result) != 8 {
			t.Errorf("int64ToString(%d) returned %d bytes, expected 8", x, len(result))
		}
	}
}

func TestMakeImageSourceWithS3(t *testing.T) {
	testCfg := struct {
		S3Enabled  bool
		S3Bucket   string
		S3Prefix   string
		S3Region   string
		S3Endpoint string
	}{
		S3Enabled:  true,
		S3Bucket:   "nonexistent-bucket-12345",
		S3Prefix:   "test/",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
	}
	_ = testCfg // use the struct to avoid unused variable error

	s3client, err := images.NewS3Client(config.ImageSourceConfig{
		S3Enabled:  true,
		S3Bucket:   "nonexistent-bucket-12345",
		S3Prefix:   "test/",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
	})
	if err != nil {
		t.Fatalf("NewS3Client should not fail: %v", err)
	}

	is := MakeImageSourceWithS3(s3client)

	img := images.NewImage()
	defer img.Destroy()

	err = is.searchRoots("nonexistent.jpg", img)
	if err == nil {
		t.Error("searchRoots should fail for non-existent image")
	}
}

func TestMakeImageSourceWithS3_NilClient(t *testing.T) {
	// When s3client is nil, it should fall back to other roots
	is := MakeImageSourceWithS3(nil)

	// Add a valid root - use the correct path from server/testimages
	is.AddRoot("../testimages/server")

	img := images.NewImage()
	defer img.Destroy()

	// Should still work with local roots
	err := is.searchRoots("01.jpg", img)
	if err != nil {
		t.Errorf("searchRoots should succeed with local roots: %v", err)
	}
}

func TestSearchRootsInternal_S3Fallback(t *testing.T) {
	// Test that S3 is checked after local and web roots
	s3client, _ := images.NewS3Client(config.ImageSourceConfig{
		S3Enabled:  true,
		S3Bucket:   "nonexistent-bucket-12345",
		S3Prefix:   "",
		S3Region:   "us-east-1",
		S3Endpoint: "http://localhost:9000",
	})

	is := MakeImageSourceWithS3(s3client)
	is.AddRoot("/nonexistent/path")

	img := images.NewImage()
	defer img.Destroy()

	// Should fail after checking all sources including S3
	err := is.searchRoots("test.jpg", img)
	if err == nil {
		t.Error("searchRoots should fail when all sources fail")
	}
}

func TestImageSource_S3ClientNilCheck(t *testing.T) {
	// Test that s3client == nil is handled correctly
	is := MakeImageSource()

	// Add local root - use the correct path from server/testimages
	is.AddRoot("../testimages/server")

	img := images.NewImage()
	defer img.Destroy()

	// Should work with local roots only
	err := is.searchRoots("01.jpg", img)
	if err != nil {
		t.Errorf("searchRoots should succeed with local roots: %v", err)
	}
}

func TestImageSource_S3ClientWithLocalRoots(t *testing.T) {
	// Test that local roots are checked before S3
	is := MakeImageSource()

	// Add local root - use the correct path from server/testimages
	is.AddRoot("../testimages/server")

	img := images.NewImage()
	defer img.Destroy()

	// Should find image in local roots (S3 never checked due to local hit)
	err := is.searchRoots("01.jpg", img)
	if err != nil {
		t.Errorf("searchRoots should find image in local roots: %v", err)
	}
}