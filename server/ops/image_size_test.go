package ops

import (
	"testing"
)

func TestImageSize(t *testing.T) {
	s := MakeImageSource()
	s.AddRoot("/app/server/testimages/server/")

	// Test with local image - 01.jpg is 900x1200 (portrait)
	w, h, err := s.ImageSize("01.jpg")
	if err != nil {
		t.Fatal(err)
	}
	if w != 900 || h != 1200 {
		t.Fatalf("Wrong image size returned. got w: %v, h:%v. expected w:%v, h:%v", w, h, 900, 1200)
	}

	// Test cache hit - second call should return from cache
	w2, h2, err := s.ImageSize("01.jpg")
	if err != nil {
		t.Fatal(err)
	}
	if w2 != 900 || h2 != 1200 {
		t.Fatalf("Wrong image size from cache. got w: %v, h:%v. expected w:%v, h:%v", w2, h2, 900, 1200)
	}

	// Test with non-existent image
	w, h, err = s.ImageSize("IMAGETHATDOESNOTEXIST")
	if err == nil {
		t.Fatal("No error returned when trying to get size of non existing image from fs.")
	}
}

func TestImageSize_ExternalSource(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping external network test in short mode")
	}

	s := MakeImageSource()
	s.AddRoot("https://upload.wikimedia.org/wikipedia/commons/")

	// Test loading image from external HTTP source
	// Note: Wikimedia may rate-limit, so this test may fail intermittently
	w, h, err := s.ImageSize("b/b4/JPEG_example_JPG_RIP_100.jpg")
	if err != nil {
		// Skip if rate-limited or network error
		if err.Error() == "Couldn't load image. Server returned 429" ||
			err.Error() == "Couldn't load image. Server returned 403" {
			t.Skip("External server rate-limited or blocked request")
		}
		t.Fatal(err)
	}
	if w != 313 || h != 234 {
		t.Fatalf("Wrong image size returned. got w: %v, h:%v. expected w:%v, h:%v", w, h, 313, 234)
	}

	// Test error handling for non-existent external image
	_, _, err = s.ImageSize("https://IMAGETHATDOESNOTEXIST")
	if err == nil {
		t.Fatal("No error returned when trying to get size of non existing image from fs.")
	}
}
