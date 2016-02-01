package cache

import (
	
	"errors"
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Resizer interface {
	GetImage(string, uint, uint) (images.ImageBlob, error)
	ImageSize(string) (uint, uint, error)

	AddRoot(string) error
	RemoveRoot(string) error
}

type BasicResizer struct {
	Roots []string
	sizes map[string][2]uint
}

func NewBasicResizer() *BasicResizer {
	return &BasicResizer{sizes: make(map[string][2]uint)}
}

// Get image whose dimensions are known.
func (c *BasicResizer) GetImage(filename string, width, height uint) (blob images.ImageBlob, err error) {
	logging.Debug(fmt.Sprintf("Get image: %v, %v, %v", filename, width, height))
	// Search requested image from all roots by trial and error
	image, err := c.searchRoots(filename)
	if err != nil {
		return
	}
	// Unlike in GetOriginal Destroy is needed since Image.Resized creates a copy of this original image
	defer image.Destroy()
	// Get resized image
	wx := strconv.FormatUint(uint64(width), 10)
	wy := strconv.FormatUint(uint64(height), 10)
	logging.Debug("Resize to: " + wx + "x" + wy)
	resized, err := image.Resized(width, height)
	if err != nil {
		return
	}
	// Add watermark

	err = resized.Watermark()
	if err != nil {
		return
	}
	// Extract requested format/extension from the filename
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	// If no extension was given, return image as is.
	if len(ext) == 0 {
		logging.Debug("No extension provided, return image as is without convesion")
		blob = resized.ToBlob()
		return
	}
	// Convert the resized image to requested format
	logging.Debug("Converting: " + ext)
	err = resized.Convert(ext)
	if err != nil {
		return
	}
	// Get blob from the image. This destroys the image object.
	blob = resized.ToBlob()
	return
}

func (c *BasicResizer) ImageSize(fn string) (w uint, h uint, err error) {
	logging.Debug(fmt.Sprintf("Get original image: %v", fn))

	if s, ok := c.sizes[fn]; ok {
		return s[0], s[1], nil
	}

	image, err := c.searchRoots(fn)
	if err != nil {
		return
	}
	defer image.Destroy()

	w = image.GetWidth()
	h = image.GetHeight()
	c.sizes[fn] = [2]uint{w, h}
	return
}

// Search root for an image. Returned image should be destroyed by image.Destroy, image.Resized or image.ToBlob or other.
func (c *BasicResizer) searchRoots(filename string) (image images.Image, err error) {
	if len(c.Roots) == 0 {
		logging.Debug("No roots")
		err = os.ErrNotExist
		return
	}
	// Extract requested type/extension and id from filename
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	id := strings.TrimRight(filename[0:len(filename)-len(ext)], ".")
	// Search requested image from all roots by trial and error
	for _, root := range c.Roots {
		// TODO: Fix escape vulnerability (sanitize filename from at least ".." etc)
		// Assume image is stored as .jpg -> change extension to .jpg
		trial := filepath.Join(root, id) + ".jpg"
		image, err = images.LoadImage(trial)
		if err == nil {
			logging.Debug("Found: " + trial)
			break
		}
		logging.Debug("Not found: " + trial)
	}
	return
}

// A very trivial (and inefficient way to handle roots)
// Can be used for development work, however.
func (self *BasicResizer) AddRoot(fileroot string) error {
	abspath, err := filepath.Abs(fileroot)
	if err != nil {
		return err
	}

	logging.Debug("Adding root: " + fileroot + " -> " + abspath)
	for _, path := range self.Roots {
		if path == abspath {
			return errors.New("Root is already served")
		}
	}

	self.Roots = append(self.Roots, abspath)
	return nil
}

// A very trivial (and inefficient way to handle roots)
// Can be used for development work, however.
func (self *BasicResizer) RemoveRoot(fileroot string) error {
	abspath, err := filepath.Abs(fileroot)
	if err != nil {
		return err
	}

	for i, path := range self.Roots {
		if path == abspath {
			// TODO: Fix possible memory leak
			// https://github.com/golang/go/wiki/SliceTricks
			self.Roots = append(self.Roots[:i], self.Roots[i+1:]...)
			return nil
		}
	}

	return errors.New("Root not found")
}
