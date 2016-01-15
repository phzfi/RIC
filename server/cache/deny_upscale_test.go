package cache

import (
	"testing"
	)

func TestDenyUpscale(t *testing.T){

	err := GetImageFromServer("toget.jpg", "?width=2000000&height=2000000", "toget.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
}

