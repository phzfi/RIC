package images

import (
	"github.com/joonazan/imagick/imagick"
)

func init() {
	imagick.Initialize()
}

// Imageblob is just an image file dumped, byte by byte to an byte array.
type ImageBlob []byte

//Image is an uncompressed image that must be convertd to blob before serving to a client.
type Image struct {
	*imagick.MagickWand
}

func NewImage() Image {
	return Image{imagick.NewMagickWand()}
}

func (img Image) Clone() Image {
	return Image{img.MagickWand.Clone()}
}

//Method for converting Image to ImageBlob. Note: Method Destroys the used Image and frees the memory used.
func (img Image) ToBlob() (blob ImageBlob) {
	blob = img.GetImageBlob()
	img.Destroy()
	return
}
