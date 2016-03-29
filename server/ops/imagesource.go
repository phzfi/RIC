package ops

import (
	"errors"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type dim [2]int
type idToSize map[string]dim

type ImageSource struct {
	roots []string
	sizes idToSize
	mutex *sync.RWMutex
}

func MakeImageSource() ImageSource {
	return ImageSource{
		sizes: make(idToSize),
		mutex: new(sync.RWMutex),
	}
}

func (i ImageSource) LoadImageOp(id string) Operation {
	return loadImageOp{&i, id}
}

// Searches root for an image. If found loads the image to img. Otherwise does nothing and returns an error.
func (i ImageSource) searchRoots(fn string, img images.Image) (err error) {
	return i.searchRootsCustomTrialFunc(fn, img.FromFile)
}

// Searches root for an image. Calls the given trialFunc with the given fn for every root until trialFunc does not return an error. Returns if trialFunc succeeds. returns with error if no trialFunc succeeds.
func (i ImageSource) searchRootsCustomTrialFunc(fn string, trialFunc func(fn string) (err error)) (err error) {
	if len(i.roots) == 0 {
		logging.Debug("No roots")
		err = os.ErrNotExist
		return
	}
	// Extract requested type/extension and id from filename
	ext := strings.TrimLeft(filepath.Ext(fn), ".")
	id := strings.TrimRight(fn[0:len(fn)-len(ext)], ".")
	// Search requested image from all roots by trial and error
	for _, root := range i.roots {
		// TODO: Fix escape vulnerability (sanitize filename from at least ".." etc)
		// Assume image is stored as .jpg -> change extension to .jpg
		trial := filepath.Join(root, id) + ".jpg"
		err = trialFunc(trial)
		if err == nil {
			logging.Debug("Found: " + trial)
			break
		}
		logging.Debug("Not found: " + trial)
	}
	return
}

// Get image size
func (i ImageSource) ImageSize(fn string) (w int, h int, err error) {
	i.mutex.RLock()
	s, ok := i.sizes[fn]
	i.mutex.RUnlock()

	if ok {
		return s[0], s[1], nil
	}

	image := images.NewImage()
	defer image.Destroy()

	err = i.pingRoots(fn, image)
	if err != nil {
		return
	}

	w = image.GetWidth()
	h = image.GetHeight()

	i.mutex.Lock()
	i.sizes[fn] = dim{w, h}
	i.mutex.Unlock()

	return
}

// Searches root for an image. If found, loads only the image metadata to img. Otherwise does nothing and returns an error.
func (i ImageSource) pingRoots(fn string, img images.Image) (err error) {
	return i.searchRootsCustomTrialFunc(fn, img.PingImage)
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
