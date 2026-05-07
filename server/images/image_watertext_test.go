package images

import (
	"testing"
)

func TestImageWatermarkText(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	err = img.WatermarkText("Test", 24, "rgba(255,255,255,0.7)", 1.0, 1.0, 0.04)
	if err != nil {
		t.Fatalf("WatermarkText failed: %v", err)
	}

	blob := img.Blob()
	if len(blob) == 0 {
		t.Error("Blob should not be empty after watermarking")
	}
}

func TestImageWatermarkText_DefaultColor(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	err = img.WatermarkText("Test", 24, "", 1.0, 1.0, 0.04)
	if err != nil {
		t.Fatalf("WatermarkText with empty color failed: %v", err)
	}

	blob := img.Blob()
	if len(blob) == 0 {
		t.Error("Blob should not be empty after watermarking with default color")
	}
}

func TestImageWatermarkText_Positions(t *testing.T) {
	tests := []struct {
		horizontal float64
		vertical   float64
		margin     float64
		name       string
	}{
		{1.0, 1.0, 0.04, "bottom-right"},
		{0.0, 1.0, 0.04, "bottom-left"},
		{1.0, 0.0, 0.04, "top-right"},
		{0.0, 0.0, 0.04, "top-left"},
		{0.5, 0.5, 0.10, "center-large-margin"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			img := NewImage()
			defer img.Destroy()

			err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
			if err != nil {
				t.Fatalf("FromFile failed: %v", err)
			}

			err = img.WatermarkText("Test", 24, "rgba(255,255,255,0.7)", tc.horizontal, tc.vertical, tc.margin)
			if err != nil {
				t.Fatalf("WatermarkText failed for %s: %v", tc.name, err)
			}

			blob := img.Blob()
			if len(blob) == 0 {
				t.Errorf("Blob should not be empty after watermarking at %s", tc.name)
			}
		})
	}
}

func TestImageWatermarkText_FontSizes(t *testing.T) {
	sizes := []float64{12, 18, 24, 36, 48, 72}

	for _, size := range sizes {
		t.Run("fontSize", func(t *testing.T) {
			img := NewImage()
			defer img.Destroy()

			err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
			if err != nil {
				t.Fatalf("FromFile failed: %v", err)
			}

			err = img.WatermarkText("Test", size, "rgba(255,255,255,0.7)", 1.0, 1.0, 0.04)
			if err != nil {
				t.Fatalf("WatermarkText with fontSize %v failed: %v", size, err)
			}

			blob := img.Blob()
			if len(blob) == 0 {
				t.Errorf("Blob should not be empty after watermarking with fontSize %v", size)
			}
		})
	}
}

func TestImageWatermarkText_LongText(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	longText := "This is a very long watermark text that should still fit on the image"
	err = img.WatermarkText(longText, 24, "rgba(255,255,255,0.7)", 1.0, 1.0, 0.04)
	if err != nil {
		t.Fatalf("WatermarkText with long text failed: %v", err)
	}

	blob := img.Blob()
	if len(blob) == 0 {
		t.Error("Blob should not be empty after watermarking with long text")
	}
}

func TestImageWatermarkText_EmptyText(t *testing.T) {
	img := NewImage()
	defer img.Destroy()

	err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
	if err != nil {
		t.Fatalf("FromFile failed: %v", err)
	}

	origBlob := img.Blob()

	err = img.WatermarkText("", 24, "rgba(255,255,255,0.7)", 1.0, 1.0, 0.04)
	if err != nil {
		t.Fatalf("WatermarkText with empty text failed: %v", err)
	}

	newBlob := img.Blob()
	if len(newBlob) == 0 {
		t.Error("Blob should not be empty after watermarking with empty text")
	}
	if len(newBlob) != len(origBlob) {
		t.Log("Blob size changed with empty text (expected due to image processing)")
	}
}

func TestImageWatermarkText_DifferentMargins(t *testing.T) {
	margins := []float64{0.01, 0.04, 0.10, 0.20}

	for _, margin := range margins {
		t.Run("margin", func(t *testing.T) {
			img := NewImage()
			defer img.Destroy()

			err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
			if err != nil {
				t.Fatalf("FromFile failed: %v", err)
			}

			err = img.WatermarkText("Test", 24, "rgba(255,255,255,0.7)", 1.0, 1.0, margin)
			if err != nil {
				t.Fatalf("WatermarkText with margin %v failed: %v", margin, err)
			}

			blob := img.Blob()
			if len(blob) == 0 {
				t.Errorf("Blob should not be empty after watermarking with margin %v", margin)
			}
		})
	}
}

func TestImageWatermarkText_DifferentColors(t *testing.T) {
	colors := []string{
		"rgba(255,255,255,0.7)",
		"rgba(0,0,0,0.5)",
		"rgba(255,0,0,0.8)",
		"white",
		"black",
		"",
	}

	for _, color := range colors {
		t.Run("color", func(t *testing.T) {
			img := NewImage()
			defer img.Destroy()

			err := img.FromFile("/app/server/testimages/loadimage/test.jpg")
			if err != nil {
				t.Fatalf("FromFile failed: %v", err)
			}

			err = img.WatermarkText("Test", 24, color, 1.0, 1.0, 0.04)
			if err != nil {
				t.Fatalf("WatermarkText with color %q failed: %v", color, err)
			}

			blob := img.Blob()
			if len(blob) == 0 {
				t.Errorf("Blob should not be empty after watermarking with color %q", color)
			}
		})
	}
}
