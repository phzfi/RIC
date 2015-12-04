package cache

import "testing"

func TestLRU(t *testing.T) {
	RunTests(NewLRU, t)
}

func TestFIFO(t *testing.T) {
	RunTests(NewFIFO, t)
}
