package logging;

import (
	"bytes"
	"testing"
)

func TestDebug(t *testing.T) {
	// Test that Debug doesn't panic with various inputs
	var buf bytes.Buffer
	// We can't easily test the output since Debug uses log.Println
	// Just verify it doesn't panic
	Debug("test")
	Debug("test", 123, 45.6)
	Debug()
	
	// If we're here, the function didn't panic
	_ = buf
}

func TestDebugf(t *testing.T) {
	// Test that Debugf doesn't panic with various inputs
	Debugf("test %s", "value")
	Debugf("test %d %f", 123, 45.6)
	Debugf("test")
	
	// If we're here, the function didn't panic
}
