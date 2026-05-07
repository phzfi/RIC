package ops

import (
	"testing"

	"github.com/phzfi/RIC/server/config"
	"github.com/phzfi/RIC/server/images"
)

func TestTextWatermarkOp(t *testing.T) {
	text := "Test Watermark"
	hor := 0.96
	ver := 0.96
	margin := 0.04
	fontSize := 24.0
	color := "rgba(255,255,255,0.7)"

	op := TextWatermarkOp(text, hor, ver, margin, fontSize, color)
	if op == nil {
		t.Fatal("TextWatermarkOp should return non-nil operation")
	}

	marshaled := op.Marshal()
	if marshaled == "" {
		t.Error("Marshal should return non-empty string")
	}
}

func TestTextWatermarkOp_Marshal(t *testing.T) {
	op := TextWatermarkOp("PHZ.fi", 0.96, 0.96, 0.04, 24, "rgba(255,255,255,0.7)")

	marshaled := op.Marshal()
	if marshaled == "" {
		t.Fatal("Marshal should return non-empty string")
	}
}

func TestTextWatermarkOp_Apply(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	op := TextWatermarkOp("PHZ.fi", 0.96, 0.96, 0.04, 24, "rgba(255,255,255,0.7)")

	err = op.Apply(img)
	if err != nil {
		t.Fatalf("TextWatermarkOp Apply failed: %v", err)
	}

	blob := img.Blob()
	if len(blob) == 0 {
		t.Error("Blob should not be empty after applying text watermark")
	}
}

func TestMakeWatermarker_Success(t *testing.T) {
	settings := config.Watermark{
		ImagePath:  "/app/server/watermark.png",
		MinHeight:  200,
		MinWidth:   200,
		MaxHeight:  5000,
		MaxWidth:   5000,
		AddMark:    true,
		Vertical:   0.5,
		Horizontal: 0.5,
	}

	wm, err := MakeWatermarker(settings)
	if err != nil {
		t.Fatalf("MakeWatermarker failed: %v", err)
	}

	if wm.WatermarkImage.MagickWand == nil {
		t.Error("WatermarkImage should be loaded")
	}

	if !wm.AddMark {
		t.Error("AddMark should be true")
	}
}

func TestMakeWatermarker_InvalidPath(t *testing.T) {
	settings := config.Watermark{
		ImagePath: "/nonexistent/path/watermark.png",
	}

	_, err := MakeWatermarker(settings)
	if err == nil {
		t.Error("MakeWatermarker should return error for invalid path")
	}
}

func TestWatermarker_EmbeddedConfig(t *testing.T) {
	settings := config.Watermark{
		ImagePath:  "/app/server/watermark.png",
		MinHeight:  100,
		MinWidth:   100,
		MaxHeight:  1000,
		MaxWidth:   1000,
		AddMark:    true,
		Vertical:   0.8,
		Horizontal: 0.2,
	}

	wm, err := MakeWatermarker(settings)
	if err != nil {
		t.Fatalf("MakeWatermarker failed: %v", err)
	}

	if wm.MinHeight != 100 {
		t.Errorf("MinHeight mismatch: got %d, want 100", wm.MinHeight)
	}
	if wm.MinWidth != 100 {
		t.Errorf("MinWidth mismatch: got %d, want 100", wm.MinWidth)
	}
	if wm.MaxHeight != 1000 {
		t.Errorf("MaxHeight mismatch: got %d, want 1000", wm.MaxHeight)
	}
	if wm.MaxWidth != 1000 {
		t.Errorf("MaxWidth mismatch: got %d, want 1000", wm.MaxWidth)
	}
	if wm.Vertical != 0.8 {
		t.Errorf("Vertical mismatch: got %f, want 0.8", wm.Vertical)
	}
	if wm.Horizontal != 0.2 {
		t.Errorf("Horizontal mismatch: got %f, want 0.2", wm.Horizontal)
	}
}
