package images

import (
	"errors"
	"gopkg.in/gographics/imagick.v2/imagick"
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

// Clone an image. Remember images and made clones need to be destroyed using Destroy(), ToBlob() or Resize().
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

// Method for converting Image to ImageBlob. Note: Method Destroys the used Image and frees the memory used.
func (img Image) Blob() (blob ImageBlob) {
	blob = img.GetImageBlob()
	img.Destroy()
	return
}

// Watermark watermarks image.
func (img Image) Watermark(watermark Image, horizontal, vertical float64) (err error) {
	x := int(float64((img.GetWidth() - watermark.GetWidth())) * horizontal)
	y := int(float64((img.GetHeight() - watermark.GetHeight())) * vertical)
	return img.CompositeImage(watermark.MagickWand, imagick.COMPOSITE_OP_OVER, x, y)
}
