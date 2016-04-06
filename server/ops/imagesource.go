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
func (i ImageSource) searchRoots(fn string, img images.Image) (err error) {
	if len(i.roots) == 0 && len(i.webroots) == 0 {
		logging.Debug("No roots")
		err = os.ErrNotExist
		return
	}
	// Extract requested type/extension and id from filename
	ext := strings.TrimLeft(filepath.Ext(fn), ".")
	id := strings.TrimRight(fn[0:len(fn)-len(ext)], ".")
	// Assume image is stored as .jpg -> change extension to .jpg
	filename := id + ".jpg"
	// Search requested image from all roots by trial and error
	for _, root := range i.roots {
		// TODO: Fix escape vulnerability (sanitize filename from at least ".." etc)
		err = img.FromFile(filepath.Join(root, filename))
		if err == nil {
			break
		}
	}

	for _, root := range i.webroots {
		logging.Debugf("Attempting to load %s", root+filename)
		err = img.FromWeb(root + filename)
		if err == nil {
			break
		}
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
	return i.searchRoots(fn, img)
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

var (
	ErrRootNotFound     = errors.New("Root not found")
	ErrRootAlreadyAdded = errors.New("Root is already served")
)

type roots []string

func (roots *roots) Add(n string) error {
	logging.Debug("Adding root: " + n)
	for _, path := range *roots {
		if path == n {
			return ErrRootAlreadyAdded
		}
	}

	*roots = append(*roots, n)
	return nil
}

func (roots *roots) Remove(r string) error {
	for i, path := range *roots {
		if path == r {
			*roots = append((*roots)[:i], (*roots)[i+1:]...)
			return nil
		}
	}

	return ErrRootNotFound
}
