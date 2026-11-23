//go:build !windows

package link

import (
	"os"
	"path/filepath"
)

func createLink(source, target string, _ Kind) (Mode, error) {
	if err := os.Symlink(source, target); err != nil {
		return "", err
	}
	return ModeLink, nil
}

func removeLink(target string, mode Mode) error {
	if mode == ModeCopy {
		return os.RemoveAll(target)
	}
	return os.Remove(target)
}

func pointsTo(target, source string) bool {
	destination, err := os.Readlink(target)
	if err != nil {
		return false
	}
	if !filepath.IsAbs(destination) {
		destination = filepath.Join(filepath.Dir(target), destination)
	}
	absoluteDestination, err := filepath.Abs(destination)
	if err != nil {
		return false
	}
	absoluteSource, err := filepath.Abs(source)
	if err != nil {
		return false
	}
	return absoluteDestination == absoluteSource
}
