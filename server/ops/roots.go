package ops

import (
	"errors"
	"github.com/phzfi/RIC/server/logging"
)

var (
	ErrRootNotFound     = errors.New("Root not found")
	ErrRootAlreadyAdded = errors.New("Root is already served")
)

type roots []string

func (roots *roots) Add(r string) error {
	logging.Debug("Adding root: " + r)
	if roots.HasRoot(r) {
		return ErrRootAlreadyAdded
	}

	*roots = append(*roots, r)
	return nil
}

func (roots *roots) HasRoot(r string) bool {
	for _, path := range *roots {
		if path == r {
			return true
		}
	}
	return false
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
