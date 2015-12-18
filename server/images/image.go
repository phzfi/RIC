package images

import (
	"github.com/joonazan/imagick/imagick"
	"errors"
)

func init() {
	imagick.Initialize()
}

// Imageblob is just an image file dumped, byte by byte to an byte array.
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


// Method for converting Image to ImageBlob. Note: Method Destroys the used Image and frees the memory used.
func (img Image) ToBlob() (blob ImageBlob) {
	blob = img.GetImageBlob()
	img.Destroy()
	return
}
