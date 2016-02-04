// +build vips

package images

import (
	"github.com/valyala/fasthttp"
	"github.com/h2non/bimg"
	"errors"
	"fmt"
	"strings"
)

func init() { }

// Imageblob is just an image file dumped, byte by byte to an byte array.
type ImageBlob []byte

// Image is an uncompressed image that must be convertd to blob before serving to a client.
type image struct{
	*bimg.Image
}

type Image struct {
	*image
}

// Create a new Image.
func NewImage() Image {
	img := Image{&image{bimg.NewImage([]byte{})}}
	return img
}

// Just for compability. Images don't need to be destroyed with vips.
func (img Image) Destroy() {
}

// Clone an image. Remember images and made clones need to be destroyed using Destroy().
func (img Image) Clone() Image {
	i := Image{&image{bimg.NewImage(img.Image.Image())}}
	return i
}

// Helper func for converting extension to format string
func extToFormat(ext string) string {
	if strings.EqualFold(ext, "jpg") {
		return "JPEG"
	}
	return strings.ToUpper(ext)
}

// Converts the image to different format. Takes extension as parameter
func (img Image) Convert(ext string) (err error) {
	formatToType := map[string]bimg.ImageType {
		"JPEG": bimg.JPEG,
		"PNG": bimg.PNG,
		"WEBP": bimg.WEBP,
		"TIFF": bimg.TIFF,
		"MAGICK": bimg.MAGICK,
	}
	_, err = img.Image.Convert(formatToType[extToFormat(ext)])
	if err != nil {
		err = errors.New("Could not convert image: " + err.Error())
	}
	return
}

// Returns image width
func (img Image) GetWidth() (width int) {
	s, _ := img.Size()
	return s.Width
}

// Returns image height
func (img Image) GetHeight() (height int) {
	s, _ := img.Size()
	return s.Height
}

// Returns filename extension of the image e.g. jpg, gif, webp
func (img Image) GetExtension() (ext string) {
	format := img.Type()
	if strings.EqualFold(format, "jpeg") {
		ext = "jpg"
	}
	return
}

func (img Image) GetImageFormat() (format string) {
	return strings.ToUpper(img.Type())
}

// Method for converting Image to ImageBlob.
func (img Image) Blob() ImageBlob {
	// The vips go binding uses blobs natively so a simple cast is enough
	return img.Image.Image()
}

// Returns Image from file.
func (img Image) FromFile(filename string) error {
	blob, err := bimg.Read(filename)
	if err != nil {
		return err
	}
	img.image.Image = bimg.NewImage(blob)
	return nil
}

// Return binary ImageBlob of an image from web.
func (img Image) FromWeb(url string) error {
	//resp, err := http.Get(url)
	statuscode, body, err := fasthttp.Get(nil, url)
	if err != nil {
		return err
	}
	//defer resp.Body.Close()

	if statuscode != 200 {
		return errors.New(fmt.Sprintf("Couldn't load image. Server returned %i", statuscode))
	}

	return img.FromBlob(body)
}

func (img Image) FromBlob(blob ImageBlob) error {
	img.image.Image = bimg.NewImage(blob)
	return nil
}

// Resize image
func (img Image) Resize(w, h int) error {
	_, err := img.ForceResize(w, h)
	return err
}

// At the moment this does just a regular resize
func (img Image) LiquidRescaleImage(w, h uint, d, r float64) error {
	_, err := img.ForceResize(int(w), int(h))
	return err
}
