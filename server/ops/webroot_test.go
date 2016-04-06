package ops

import (
	"github.com/phzfi/RIC/server/images"
	"testing"
)

func TestWebroot(t *testing.T) {
	s := MakeImageSource()
	s.AddRoot("https://upload.wikimedia.org/wikipedia/commons/")

	i := images.NewImage()
	defer i.Destroy()

	err := s.LoadImageOp("b/b4/JPEG_example_JPG_RIP_100.jpg").Apply(i)

	if err != nil {
		t.Fatalf("Error loading image from web: %s", err)
	}

	if len(i.Blob()) == 0 {
		t.Fatal("No error, but image is empty.")
	}
}
