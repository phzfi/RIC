package cache

import (
	"errors"
	"github.com/phzfi/RIC/server/images"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Cacheless struct {
	Roots []string
}

// Get image whose dimensions are known.
func (c *Cacheless) GetImage(filename string, width, height uint) (blob images.ImageBlob, err error) {
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
	log.Println("Resize to: " + wx + "x" + wy)
	image, err = image.Resized(width, height)
	if err != nil {
		return
	}
	// Extract requested format/extension from the filename
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	//Convert the resized image to requested format
	// TODO: Unsupported format -> Convert does nothing and we get the default format (jpg). Is this behaviour fine?
	log.Println("Converting: " + ext)
	err = image.Convert(ext)
	if err != nil {
		return
	}
	// Get blob from the image. This destroys the image object.
	blob = image.ToBlob()
	return
}


// Get original image as Image type (for resizing/converting)
func (c *Cacheless) GetOriginal(filename string) (blob images.ImageBlob, err error) {
	// Search roots for the Image
	image, err := c.searchRoots(filename)
	if err != nil {
		return
	}
	// Extract requested format/extension from the filename
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	//Convert image to requested format
	// TODO: Unsupported format -> Convert does nothing and we get the default format (jpg). Is this behaviour fine?
	log.Println("Converting: " + ext)
	err = image.Convert(ext)
	if err != nil {
		return
	}
	// Get blob from the Image. This also destroys the image -object.
	blob = image.ToBlob()
	return
}


// Search root for an image. Returned image should be destroyed by image.Destroy, image.Resized or image.ToBlob or other.
func (c *Cacheless) searchRoots(filename string) (image images.Image, err error) {
	if len(c.Roots) == 0 {
		err = os.ErrNotExist
		return
	}
	// Extract requested type/extension and id from filename
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	id := filename[0:len(filename)-len(ext) - 1]
	// Search requested image from all roots by trial and error
	for _, root := range c.Roots {
		// TODO: Fix escape vulnerability (sanitize filename from at least ".." etc)
		// Assume image is stored as .jpg -> change extension to .jpg
		trial := filepath.Join(root, id) + ".jpg"
		image, err = images.LoadImage(trial)
		if err == nil {
			log.Println("Found: " + trial)
			break
		}
		if !os.IsNotExist(err) {
			return
		}
		log.Println("Not found: " + trial)
	}
	return
}


// A very trivial (and inefficient way to handle roots)
// Can be used for development work, however.
func (self *Cacheless) AddRoot(fileroot string) error {
	abspath, err := filepath.Abs(fileroot)
	if err != nil {
		return err
	}

	log.Println("Adding root: " + fileroot + " -> " + abspath)
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
func (self *Cacheless) RemoveRoot(fileroot string) error {
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
