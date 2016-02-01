package images

import (
	"errors"
	"github.com/joonazan/imagick/imagick"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/config"
	"strings"
)


func init() {
	imagick.Initialize()
}

// ImageBlob is just an image file dumped, byte by byte to an byte array.
type ImageBlob []byte

// Image is an uncompressed image that must be convertd to blob before serving to a client.
type Image struct {
	*imagick.MagickWand
}

func NewImage() Image {
	return Image{imagick.NewMagickWand()}
}

// Clone an image. Remember images and made clones need to be destroyed using Destroy().
func (img Image) Clone() Image {
	return Image{img.MagickWand.Clone()}
}

// Converts the image to different format. Takes extension as parameter.
func (img Image) Convert(ext string) (err error) {
	err = img.SetImageFormat(ext)
	if err != nil {
		err = errors.New("Could not convert image: " + err.Error())
	}
	return
}

// Returns image width
func (img Image) GetWidth() (width int) {
	return int(img.GetImageWidth())
}

// Returns image height
func (img Image) GetHeight() (height int) {
	return int(img.GetImageHeight())
}

// Returns filename extension of the image e.g. jpg, gif, webp
func (img Image) GetExtension() (ext string) {
	format := img.GetImageFormat()
	ext = strings.ToLower(format)
	if strings.EqualFold(ext, "jpeg") {
		ext = "jpg"
	}
	return
}

// Method for converting Image to ImageBlob.
func (img Image) Blob() ImageBlob {
	return img.GetImageBlob()
}

// Watermark watermarks image.
func (img Image) Watermark() (err error) {
	logging.Debug("Watermarking")
  logging.Debug(*config.Config())
	minHeight, err := config.WatermarkInt("minheight")
	minWidth, err := config.WatermarkInt("minwidth")
	maxHeight, err := config.WatermarkInt("maxheight")
	maxWidth, err := config.WatermarkInt("maxwidth")

	if err != nil {
		logging.Debug("Error reading config size restrictions." + err.Error())
		return
	}

	heightOK := img.GetHeight() > uint(minHeight) && img.GetHeight() < uint(maxHeight)
	widthOK := img.GetWidth() > uint(minWidth) && img.GetWidth() < uint(maxWidth)
	if (!heightOK && !widthOK) {
		return
	}

	watermark, err := LoadImage(config.Watermark("path"))
	if err != nil {
		logging.Debug("Error loading watermark image." + err.Error())
		return
	}

	horizontal, err := config.WatermarkFloat64("horizontal")
	vertical, err := config.WatermarkFloat64("vertical")

	if err != nil {
		logging.Debug("Error loading config alignment." + err.Error())
		return
	}

	x := int(float64((img.GetWidth() - watermark.GetWidth())) * horizontal)
	y := int(float64((img.GetHeight() - watermark.GetHeight())) * vertical)
	return img.CompositeImage(watermark.MagickWand, imagick.COMPOSITE_OP_OVER, x, y)
}
