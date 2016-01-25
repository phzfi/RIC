package images

import (
	"errors"
	"github.com/joonazan/imagick/imagick"
	"github.com/phzfi/RIC/server/logging"
	"gopkg.in/gcfg.v1"
	"strings"
)

var conf = struct {
	watermark struct {
		path string
	}

	min struct {
		width uint
		height uint
	}

	max struct {
		width uint
		height uint
	}

	alignment struct {
		horizontal float64
		vertical float64
	}
}{}

func init() {
	imagick.Initialize()
	logging.Debug("Reading config")
	err := gcfg.ReadFileInto(&conf, "watermark-config.gcfg")
	if err != nil {
		logging.Debug("Couldn't read watermark config." + err.Error())
	}
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
	if img.GetHeight() < conf.min.height && img.GetWidth() < conf.min.width && img.GetHeight() > conf.max.height && img.GetWidth() > conf.max.width {
		return
	}

	logging.Debug(conf)

	watermark, err := LoadImage(conf.watermark.path)
	if err != nil {
		logging.Debug("Error loading watermark image." + err.Error())
		return
	}
	x := int(float64((img.GetWidth() - watermark.GetWidth())) * conf.alignment.horizontal)
	y := int(float64((img.GetHeight() - watermark.GetHeight())) * conf.alignment.vertical)
	return img.CompositeImage(watermark.MagickWand, imagick.COMPOSITE_OP_OVER, x, y)
}
