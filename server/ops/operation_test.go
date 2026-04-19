package ops

import (
	"testing"

	"github.com/phzfi/RIC/server/images"
)

func TestResizeOperation(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	origW := img.GetWidth()
	origH := img.GetHeight()

	op := Resize{Width: origW / 2, Height: origH / 2}

	marshaled := op.Marshal()
	if marshaled == "" {
		t.Error("Marshal should return non-empty string")
	}

	err = op.Apply(img)
	if err != nil {
		t.Fatalf("Resize Apply failed: %v", err)
	}

	if img.GetWidth() != origW/2 {
		t.Errorf("Width mismatch: got %d, want %d", img.GetWidth(), origW/2)
	}
}

func TestCropOperation(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	origW := img.GetWidth()
	origH := img.GetHeight()

	op := Crop{Width: origW / 2, Height: origH / 2, X: 10, Y: 10}

	marshaled := op.Marshal()
	if marshaled == "" {
		t.Error("Marshal should return non-empty string")
	}

	err = op.Apply(img)
	if err != nil {
		t.Fatalf("Crop Apply failed: %v", err)
	}
}

func TestConvertOperation(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	op := Convert{Format: "png"}

	marshaled := op.Marshal()
	if marshaled == "" {
		t.Error("Marshal should return non-empty string")
	}

	err = op.Apply(img)
	if err != nil {
		t.Fatalf("Convert Apply failed: %v", err)
	}

	ext := img.GetExtension()
	if ext != "png" {
		t.Errorf("Extension mismatch: got %s, want png", ext)
	}
}

func TestLiquidRescaleOperation(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	origW := img.GetWidth()
	origH := img.GetHeight()

	op := LiquidRescale{Width: origW / 2, Height: origH / 2}

	marshaled := op.Marshal()
	if marshaled == "" {
		t.Error("Marshal should return non-empty string")
	}

	err = op.Apply(img)
	if err != nil {
		t.Fatalf("LiquidRescale Apply failed: %v", err)
	}
}

func TestLoadImageOperation(t *testing.T) {
	is := MakeImageSource()
	err := is.AddRoot("/app/server/testimages/loadimage/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	op := is.LoadImageOp("test.jpg")

	marshaled := op.Marshal()
	if marshaled == "" {
		t.Error("Marshal should return non-empty string")
	}

	img := images.NewImage()
	defer img.Destroy()

	err = op.Apply(img)
	if err != nil {
		t.Fatalf("LoadImage Apply failed: %v", err)
	}

	if img.GetWidth() == 0 || img.GetHeight() == 0 {
		t.Error("Image should have dimensions")
	}
}

func TestLoadImageOperationNotFound(t *testing.T) {
	is := MakeImageSource()
	err := is.AddRoot("/app/server/testimages/loadimage/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	op := is.LoadImageOp("nonexistent.jpg")

	img := images.NewImage()
	defer img.Destroy()

	err = op.Apply(img)
	if err == nil {
		t.Error("LoadImage with non-existent file should return error")
	}
}

func TestWatermarkOperation(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	wm := images.NewImage()
	defer wm.Destroy()

	err = wm.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("FromFile watermark failed: %v", err)
	}

	op := watermark{
		stamp:      wm,
		horizontal: 0.5,
		vertical:   0.5,
	}

	marshaled := op.Marshal()
	if marshaled == "" {
		t.Error("Marshal should return non-empty string")
	}

	err = op.Apply(img)
	if err != nil {
		t.Fatalf("Watermark Apply failed: %v", err)
	}
}

func TestWatermarkOp(t *testing.T) {
	wm := images.NewImage()
	defer wm.Destroy()

	err := wm.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	op := WatermarkOp(wm, 0.5, 0.5)
	if op == nil {
		t.Error("WatermarkOp should return non-nil operation")
	}

	marshaled := op.Marshal()
	if marshaled == "" {
		t.Error("Marshal should return non-empty string")
	}
}