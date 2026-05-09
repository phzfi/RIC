//go:build debug
// +build debug

package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestDebug_DebugBuild(t *testing.T) {
	// This test is only run when the debug tag is set
	// The Debug function should actually log something
	
	// We can't easily capture log output, but we can verify it doesn't panic
	Debug("test message")
	Debug("test", 123, 45.6)
	Debug()
	
	// If we're here, the function didn't panic
}

func TestDebugf_DebugBuild(t *testing.T) {
	// This test is only run when the debug tag is set
	// The Debugf function should actually log something
	
	// We can't easily capture log output, but we can verify it doesn't panic
	Debugf("test %s", "value")
	Debugf("test %d %f", 123, 45.6)
	Debugf("test")
	
	// If we're here, the function didn't panic
}
