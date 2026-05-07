package main

import (
	"testing"

	"github.com/phzfi/RIC/server/config"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
	"github.com/valyala/fasthttp"
)

func TestParseURI_WatermarkTextParam(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=500&height=500&mode=liquid&format=jpeg&watermark=PHZ.fi")
	uri.SetPath("/01.jpg")

	operations, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}

	var foundTextWatermark bool
	for _, op := range operations {
		if op.Marshal() != "" {
			m := op.Marshal()
			if len(m) > 1 && m[0] == '2' {
				foundTextWatermark = true
			}
		}
	}

	if !foundTextWatermark {
		t.Error("Should have found text watermark operation")
	}
}

func TestParseURI_WatermarkFalseParam(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=500&height=500&mode=liquid&format=jpeg&watermark=false")
	uri.SetPath("/01.jpg")

	operations, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}

	watermarkCount := 0
	for _, op := range operations {
		m := op.Marshal()
		if len(m) > 0 && m[0] == '4' {
			watermarkCount++
		}
	}

	if watermarkCount != 0 {
		t.Errorf("Should have no watermark operations, got %d", watermarkCount)
	}
}

func TestParseURI_WatermarkTrueParam(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=500&height=500&mode=liquid&format=jpeg&watermark=true")
	uri.SetPath("/01.jpg")

	operations, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}

	var foundImageWatermark bool
	for _, op := range operations {
		m := op.Marshal()
		if len(m) > 0 && m[0] == '4' {
			foundImageWatermark = true
		}
	}

	if !foundImageWatermark {
		t.Error("Should have found image watermark operation")
	}
}

func TestParseURI_InvalidMode(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=500&height=500&mode=invalid&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, _, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr == nil {
		t.Fatal("Should return error for invalid mode")
	}
}

func TestParseURI_InvalidFormat(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=500&height=500&mode=liquid&format=invalid")
	uri.SetPath("/01.jpg")

	_, _, _, _, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr == nil {
		t.Fatal("Should return error for invalid format")
	}
}

func TestParseURI_ImageNotFound(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=500&height=500&mode=liquid&format=jpeg")
	uri.SetPath("/nonexistent.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr == nil {
		t.Error("Should return error for non-existent image")
	}
}

func TestParseURI_FitMode(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=500&height=500&mode=fit&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_CropMode(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=200&height=200&mode=crop&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_CropmidMode(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=200&height=200&mode=cropmid&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_WithUrlParam(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=200&height=200&mode=fit&format=jpeg&url=http://example.com")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_FitModeWidthAdjust(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=100&height=100&mode=fit&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_ResizeOnlyWidth(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=200&mode=resize&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_ResizeOnlyHeight(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("height=200&mode=resize&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_ResizeNoDimensions(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("mode=resize&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_LiquidOnlyWidth(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=200&mode=liquid&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_LiquidOnlyHeight(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("height=200&mode=liquid&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_LiquidNoDimensions(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("mode=liquid&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_FitOnlyWidth(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=200&mode=fit&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_FitOnlyHeight(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("height=200&mode=fit&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_FitNoDimensions(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("mode=fit&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_CropOnlyWidth(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=200&mode=crop&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_CropOnlyHeight(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("height=200&mode=crop&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_CropNoDimensions(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("mode=crop&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_CropmidOnlyWidth(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=200&mode=cropmid&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_CropmidOnlyHeight(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("height=200&mode=cropmid&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestParseURI_CropmidNoDimensions(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("mode=cropmid&format=jpeg")
	uri.SetPath("/01.jpg")

	_, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}
}

func TestRoundedIntegerDivision_NegativeN(t *testing.T) {
	result := roundedIntegerDivision(-5, 6)
	if result != -1 {
		t.Errorf("Expected -1, got %d", result)
	}
}

func TestRoundedIntegerDivision_NegativeM(t *testing.T) {
	result := roundedIntegerDivision(5, -6)
	if result != -1 {
		t.Errorf("Expected -1, got %d", result)
	}
}

func TestRoundedIntegerDivision_BothNegative(t *testing.T) {
	result := roundedIntegerDivision(-5, -6)
	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}
}

func TestRoundedIntegerDivision_Positive(t *testing.T) {
	result := roundedIntegerDivision(5, 6)
	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}
}

func TestParseURI_ForceMarkOverridesURLParam(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
			ForceMark:  "FORCED",
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=500&height=500&mode=liquid&format=jpeg&watermark=URLTEXT")
	uri.SetPath("/01.jpg")

	operations, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}

	var foundForcedText bool
	for _, op := range operations {
		m := op.Marshal()
		if len(m) > 6 && m[0] == '4' && m[1:7] == "FORCED" {
			foundForcedText = true
		}
	}

	if !foundForcedText {
		t.Error("Should have found forced watermark text 'FORCED', not URL parameter 'URLTEXT'")
	}
}

func TestParseURI_ForceMarkImageWatermark(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
			ForceMark:  "true",
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=500&height=500&mode=liquid&format=jpeg&watermark=IGNORED")
	uri.SetPath("/01.jpg")

	operations, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}

	var foundImageWatermark bool
	for _, op := range operations {
		m := op.Marshal()
		if len(m) > 1 && m[0] == '4' && m[1] != 'F' && m[1] != 'T' {
			foundImageWatermark = true
		}
	}

	if !foundImageWatermark {
		t.Error("Should have found image watermark when ForceMark=true")
	}
}

func TestParseURI_ForceMarkDisablesWatermark(t *testing.T) {
	source := ops.MakeImageSource()
	err := source.AddRoot("/app/server/testimages/server/")
	if err != nil {
		t.Fatalf("AddRoot failed: %v", err)
	}

	wmImage := images.NewImage()
	defer wmImage.Destroy()
	err = wmImage.FromFile("/app/server/watermark.png")
	if err != nil {
		t.Fatalf("Load watermark failed: %v", err)
	}

	marker := ops.Watermarker{
		WatermarkImage: wmImage,
		Watermark: config.Watermark{
			AddMark:    true,
			MinHeight:  200,
			MinWidth:   200,
			MaxHeight:  5000,
			MaxWidth:   5000,
			Horizontal: 0.5,
			Vertical:   0.5,
			ForceMark:  "false",
		},
	}

	uri := fasthttp.URI{}
	uri.SetQueryString("width=500&height=500&mode=liquid&format=jpeg&watermark=true")
	uri.SetPath("/01.jpg")

	operations, _, _, opErr, invalidErr := ParseURI(&uri, source, marker)
	if invalidErr != nil {
		t.Fatalf("Invalid error: %v", invalidErr)
	}
	if opErr != nil {
		t.Fatalf("Operation error: %v", opErr)
	}

	watermarkCount := 0
	for _, op := range operations {
		m := op.Marshal()
		if len(m) > 0 && m[0] == '4' {
			watermarkCount++
		}
	}

	if watermarkCount != 0 {
		t.Errorf("Should have no watermark operations when ForceMark=false, got %d", watermarkCount)
	}
}
