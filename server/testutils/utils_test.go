package testutils

import (
	"os"
	"testing"
)

func TestRemoveContents(t *testing.T) {
	// Create a test directory with some files
	testDir := "/tmp/testremovecontents"
	os.MkdirAll(testDir, os.ModePerm)
	os.WriteFile(testDir+"/file1.txt", []byte("test"), 0644)
	os.WriteFile(testDir+"/file2.txt", []byte("test"), 0644)

	err := RemoveContents(testDir)
	if err != nil {
		t.Fatalf("RemoveContents failed: %v", err)
	}

	// Verify directory is empty
	d, _ := os.Open(testDir)
	names, _ := d.Readdirnames(-1)
	if len(names) != 0 {
		t.Error("Directory should be empty after RemoveContents")
	}
	d.Close()
	os.RemoveAll(testDir)
}

func TestRemoveContentsError(t *testing.T) {
	err := RemoveContents("/nonexistent/dir")
	if err == nil {
		t.Error("RemoveContents should return error for nonexistent dir")
	}
}
