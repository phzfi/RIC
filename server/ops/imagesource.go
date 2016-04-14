package ops

import (
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
	roots, webroots roots
	sizes           idToSize
	mutex           *sync.RWMutex
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
func (i ImageSource) searchRoots(fn string, img images.Image) error {
	return i.searchRootsInternal(fn, img.FromFile, img.FromWeb)
}

func (i ImageSource) searchRootsInternal(filename string, visitPath, visitURL func(string) error) (err error) {
	if len(i.roots) == 0 && len(i.webroots) == 0 {
		logging.Debug("No roots")
		err = os.ErrNotExist
		return
	}

	// Search requested image from all roots by trial and error
	for _, root := range i.roots {
		// TODO: Fix escape vulnerability (sanitize filename from at least ".." etc)
		err = visitPath(filepath.Join(root, filename))
		if err == nil {
			return
		}
	}

	for _, root := range i.webroots {
		logging.Debugf("Attempting to load %s", root+filename)
		err = visitURL(root + filename)
		if err == nil {
			return
		}
	}
	return
}

// Searches root for an image. If found, loads only the image metadata to img. Otherwise does nothing and returns an error.
func (i ImageSource) pingRoots(fn string, img images.Image) (err error) {
	return i.searchRootsInternal(fn, img.PingImage, img.FromWeb)
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

func isWebroot(root string) bool {
	return strings.HasPrefix(root, "http:") || strings.HasPrefix(root, "https:")
}

func (i *ImageSource) AddRoot(root string) error {

	if isWebroot(root) {
		return i.webroots.Add(root)
	}

	abspath, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	return i.roots.Add(abspath)
}

func (is *ImageSource) RemoveRoot(root string) error {

	if isWebroot(root) {
		return is.webroots.Remove(root)
	}

	abspath, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	return is.roots.Remove(abspath)
}
