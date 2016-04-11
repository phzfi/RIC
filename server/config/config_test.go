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
