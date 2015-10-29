package cache

import (
	"github.com/phzfi/RIC/server/images"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Cacheless struct {
	roots []string
}

// Get the image with the desired size. Size params can be nil if the
// original is desired.
func (self *Cacheless) GetImage(filename string, width *uint, height *uint) (images.ImageBlob, error) {
	var blob images.ImageBlob = nil
	var err error = nil

	for _, root := range self.roots {
		// TODO: Fix escape vulnerability (sanitize filename from at least ".." etc)
		trial := root + filename

		blob, err = images.LoadImage(trial)
		if err == nil {
			log.Println("Found: " + trial)
			break
		}
		if !os.IsNotExist(err) {
			return nil, err
		}
		log.Println("Not found: " + trial)
	}
	if err != nil {
		return nil, err
	}

	// TODO: If only one is set?
	if width != nil && height != nil {
		wx := strconv.FormatUint(uint64(*width), 10)
		wy := strconv.FormatUint(uint64(*height), 10)
		log.Println("Resize to: " + wx + "x" + wy)
		blob, err = blob.Resized(*width, *height)
		if err != nil {
			return nil, err
		}
	}
	return blob, nil
}

// A very trivial (and inefficient way to handle roots)
// Can be used for development work, however.
func (self *Cacheless) AddRoot(fileroot string) error {
	log.Println("Adding root: " + fileroot)
	for i := range self.roots {
		path := self.roots[i]
		if path == fileroot {
			return errors.New("Root is already served")
		}
	}

	abspath, err := filepath.Abs(fileroot)
	if err != nil {
		return err
	}

	// TODO: Check if the path exists?
	log.Println("Resolved root: " + abspath)
	self.roots = append(self.roots, abspath)
	return nil
}

// A very trivial (and inefficient way to handle roots)
// Can be used for development work, however.
func (self *Cacheless) RemoveRoot(fileroot string) error {
	abspath, err := filepath.Abs(fileroot)
	if err != nil {
		return err
	}
	for i := range self.roots {
		path := self.roots[i]
		if path == abspath {
			// TODO: Fix possible memory leak
			// https://github.com/golang/go/wiki/SliceTricks
			self.roots = append(self.roots[:i], self.roots[i+1:]...)
			return nil
		}
	}
	return errors.New("Root not found")
}

