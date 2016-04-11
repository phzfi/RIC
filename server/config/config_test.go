package config

import (
	"testing"
)

// Test reading config from not existing file.
func TestReadConfigErr(t *testing.T) {
	conf := ReadConfig("nosuchfile.ini")
	if *conf != defaults {
		t.Fatal("Expected default config")
	}
}

// If there is a bug in the config reading, it will panic when reading a file
func TestNoReadConfigPanic(t *testing.T) {
	_ = ReadConfig("testconfig.ini")
}
