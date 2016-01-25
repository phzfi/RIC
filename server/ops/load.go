package ops

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/logging"
	"os"
	"strings"
	"path/filepath"
	"errors"
)

type ImageSource struct{
	roots []string
}

func (i *ImageSource) LoadImageOp(id string) Operation{
	return OperationFunc(func(img images.Image) error {
		return i.searchRoots(id, img)
	})
}

// Search root for an image. Returned image should be destroyed by image.Destroy, image.Resized or image.ToBlob or other.
func (i *ImageSource) searchRoots(filename string, img images.Image) (err error) {
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
