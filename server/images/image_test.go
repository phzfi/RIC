package images

import (
	"testing"
)

func TestImageClone(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	clone := img.Clone()
	defer clone.Destroy()

	if clone.GetWidth() != img.GetWidth() {
		t.Errorf("Clone width mismatch: got %d, want %d", clone.GetWidth(), img.GetWidth())
	}

	if clone.GetHeight() != img.GetHeight() {
		t.Errorf("Clone height mismatch: got %d, want %d", clone.GetHeight(), img.GetHeight())
	}
}

func TestImageConvert(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	err = img.Convert("png")
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	ext := img.GetExtension()
	if ext != "png" {
		t.Errorf("GetExtension() = %q, want %q", ext, "png")
	}
}

func TestImageConvertInvalidFormat(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	err = img.Convert("invalidformatthatdoesnotexist")
	if err != nil {
		t.Logf("Convert with invalid format returned error as expected: %v", err)
	}
}

func TestImageResize(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	origW := img.GetWidth()
	origH := img.GetHeight()

	err = img.Resize(origW/2, origH/2)
	if err != nil {
		t.Fatalf("Resize failed: %v", err)
	}

	if img.GetWidth() != origW/2 {
		t.Errorf("Resize width: got %d, want %d", img.GetWidth(), origW/2)
	}
}

func TestImageCrop(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	origW := img.GetWidth()
	origH := img.GetHeight()

	err = img.Crop(origW/2, origH/2, 0, 0)
	if err != nil {
		t.Fatalf("Crop failed: %v", err)
	}

	if img.GetWidth() != origW/2 {
		t.Errorf("Crop width: got %d, want %d", img.GetWidth(), origW/2)
	}
}

func TestImageWatermark(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	wm := NewImage()
	defer wm.Destroy()

	err = wm.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("FromFile watermark failed: %v", err)
	}

	err = img.Watermark(wm, 0.5, 0.5)
	if err != nil {
		t.Fatalf("Watermark failed: %v", err)
	}
}

func TestImageWatermarkCorners(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	wm := NewImage()
	defer wm.Destroy()

	err = wm.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("FromFile watermark failed: %v", err)
	}

	tests := []struct {
		h, v float64
	}{
		{0.0, 0.0},
		{1.0, 0.0},
		{0.0, 1.0},
		{1.0, 1.0},
	}

	for _, tc := range tests {
		img := NewImage()
		defer img.Destroy()

		err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
		if err != nil {
			t.Fatalf("FromFile failed: %v", err)
		}

		err = img.Watermark(wm, tc.h, tc.v)
		if err != nil {
			t.Fatalf("Watermark failed at %f,%f: %v", tc.h, tc.v, err)
		}
	}
}

func TestImageBlob(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	blob := img.Blob()
	if len(blob) == 0 {
		t.Error("Blob should not be empty")
	}
}

func TestFromBlob(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	orig := NewImage()
	defer orig.Destroy()

	err := orig.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	blob := orig.Blob()

	err = img.FromBlob(blob)
	if err != nil {
		t.Fatalf("FromBlob failed: %v", err)
	}

	if img.GetWidth() != orig.GetWidth() {
		t.Errorf("Width mismatch: got %d, want %d", img.GetWidth(), orig.GetWidth())
	}
}

func TestFromBlobInvalid(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromBlob([]byte("invalid image data"))
	if err == nil {
		t.Error("FromBlob with invalid data should return error")
	}
}

func TestFromFileNotExist(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("nonexistent_file.jpg")
	if err == nil {
		t.Error("FromFile with non-existent file should return error")
	}
}

func TestNewImage(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	if img.MagickWand == nil {
		t.Error("NewImage should create non-nil MagickWand")
	}
}

func TestImageGetExtension(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	ext := img.GetExtension()
	if ext != "jpg" {
		t.Errorf("GetExtension() = %q, want %q", ext, "jpg")
	}

	err = img.Convert("png")
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	ext = img.GetExtension()
	if ext != "png" {
		t.Errorf("GetExtension() = %q, want %q", ext, "png")
	}

	err = img.Convert("jpeg")
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	ext = img.GetExtension()
	if ext != "jpg" {
		t.Errorf("GetExtension() = %q, want %q", ext, "jpg")
	}
}