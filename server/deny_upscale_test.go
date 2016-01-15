package main

import (
	"testing"
	)

func TestDenyUpscale(t *testing.T) {
	err := GetImageFromServer("toget.jpg", "?width=20000&height=20000", "toget.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
}
