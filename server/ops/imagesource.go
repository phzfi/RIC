package ops

import (
	"errors"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"os"
	"sync"
	"path/filepath"
	"strings"
)

type dim [2]int
type idToSize map[string]dim

type ImageSource struct {
	roots []string
	sizes idToSize
    mutex sync.Mutex
}

func MakeImageSource() ImageSource {
	return ImageSource{
		sizes: make(idToSize),
	}
}

func (i ImageSource) LoadImageOp(id string) Operation {
	return loadImageOp{i, id}
}

// Search root for an image. Returned image should be destroyed by image.Destroy, image.Resized or image.ToBlob or other.
func (i ImageSource) searchRoots(filename string, img images.Image) (err error) {
	if len(i.roots) == 0 {
		logging.Debug("No roots")
		err = os.ErrNotExist
		return
	}
	// Extract requested type/extension and id from filename
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	id := strings.TrimRight(filename[0:len(filename)-len(ext)], ".")
	// Search requested image from all roots by trial and error
	for _, root := range i.roots {
		// TODO: Fix escape vulnerability (sanitize filename from at least ".." etc)
		// Assume image is stored as .jpg -> change extension to .jpg
		trial := filepath.Join(root, id) + ".jpg"
		err = img.FromFile(trial)
		if err == nil {
			logging.Debug("Found: " + trial)
			break
		}
		logging.Debug("Not found: " + trial)
	}
	return
}


// TODO: This is a temp solution for ImageSize creating too many Images.
// Limit to creating only one at time for finding the image size

func (i ImageSource) ImageSize(fn string) (w int, h int, err error) {
	i.mutex.Lock()

	if s, ok := i.sizes[fn]; ok {
		i.mutex.Unlock()
		return s[0], s[1], nil
	}

	image := images.NewImage()
	defer func () {
		i.mutex.Unlock()
		image.Destroy()
	}()

	err = i.searchRoots(fn, image)
	if err != nil {
		return
	}

	w = image.GetWidth()
	h = image.GetHeight()
	i.sizes[fn] = dim{w, h}
	return
}

// A very trivial (and inefficient way to handle roots)
// Can be used for development work, however.
func (i *ImageSource) AddRoot(root string) error {
	abspath, err := filepath.Abs(root)
	if err != nil {
		return err
	}

	logging.Debug("Adding root: " + root + " -> " + abspath)
	for _, path := range i.roots {
		if path == abspath {
			return errors.New("Root is already served")
		}
	}

	i.roots = append(i.roots, abspath)
	return nil
}

// A very trivial (and inefficient way to handle roots)
// Can be used for development work, however.
func (is *ImageSource) RemoveRoot(root string) error {
	abspath, err := filepath.Abs(root)
	if err != nil {
		return err
	}

	for i, path := range is.roots {
		if path == abspath {
			is.roots = append(is.roots[:i], is.roots[i+1:]...)
			return nil
		}
	}

	return errors.New("Root not found")
}
