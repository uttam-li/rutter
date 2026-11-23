//go:build windows

package link

import "os"

func createLink(source, target string, _ Kind) (Mode, error) {
	if err := copyPath(source, target); err != nil {
		return "", err
	}
	return ModeCopy, nil
}

func removeLink(target string, _ Mode) error {
	return os.RemoveAll(target)
}

func pointsTo(_, _ string) bool {
	return false
}
