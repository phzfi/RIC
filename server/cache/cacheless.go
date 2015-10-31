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
	Roots []string
}

// Get the image with the desired size. Size params can be nil if the
// original is desired.
func (self *Cacheless) GetImage(filename string, width *uint, height *uint) (blob images.ImageBlob, err error) {
	for _, root := range self.Roots {
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

