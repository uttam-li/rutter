package link

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func copyPath(source, target string) error {
	info, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("inspect %s: %w", source, err)
	}

	if info.IsDir() {
		return copyDir(source, target)
	}
	return copyFile(source, target, info.Mode().Perm())
}

func copyDir(source, target string) error {
	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(target, rel)

		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			return os.MkdirAll(dest, info.Mode().Perm())
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		return copyFile(path, dest, info.Mode().Perm())
	})
}

func copyFile(source, target string, perm fs.FileMode) error {
	in, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open %s: %w", source, err)
	}
	defer func() { _ = in.Close() }()

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create parent for %s: %w", target, err)
	}

	out, err := os.CreateTemp(filepath.Dir(target), ".rutter-copy-*")
	if err != nil {
		return fmt.Errorf("create temporary copy: %w", err)
	}
	tmp := out.Name()
	defer func() { _ = os.Remove(tmp) }()

	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return fmt.Errorf("copy %s: %w", source, err)
	}
	if err := out.Chmod(perm); err != nil {
		_ = out.Close()
		return fmt.Errorf("set mode on %s: %w", target, err)
	}
	if err := out.Sync(); err != nil {
		_ = out.Close()
		return fmt.Errorf("sync %s: %w", target, err)
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("close %s: %w", target, err)
	}
	if err := os.Rename(tmp, target); err != nil {
		return fmt.Errorf("replace %s: %w", target, err)
	}
	return nil
}

func movePath(source, target string) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create parent for %s: %w", target, err)
	}
	if err := os.Rename(source, target); err == nil {
		return nil
	}
	if err := copyPath(source, target); err != nil {
		return err
	}
	return os.RemoveAll(source)
}

// Move moves a file or directory, falling back to copy and remove when the
// source and target are on different volumes.
func Move(source, target string) error {
	return movePath(source, target)
}

// sameContent compares a source and target file or directory by name, size,
// and SHA-256. It lets rutter avoid overwriting files a user has edited.
func sameContent(source, target string) bool {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return false
	}
	targetInfo, err := os.Stat(target)
	if err != nil {
		return false
	}

	if sourceInfo.IsDir() != targetInfo.IsDir() {
		return false
	}

	if !sourceInfo.IsDir() {
		return sameFile(source, target, sourceInfo.Size(), targetInfo.Size())
	}

	sourceFiles, ok := treeHashes(source)
	if !ok {
		return false
	}
	targetFiles, ok := treeHashes(target)
	if !ok {
		return false
	}
	if len(sourceFiles) != len(targetFiles) {
		return false
	}
	for rel, hash := range sourceFiles {
		if targetFiles[rel] != hash {
			return false
		}
	}
	return true
}

func sameFile(source, target string, sourceSize, targetSize int64) bool {
	if sourceSize != targetSize {
		return false
	}
	sourceHash, err := fileHash(source)
	if err != nil {
		return false
	}
	targetHash, err := fileHash(target)
	if err != nil {
		return false
	}
	return sourceHash == targetHash
}

func treeHashes(root string) (map[string]string, bool) {
	hashes := make(map[string]string)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		hash, err := fileHash(path)
		if err != nil {
			return err
		}
		hashes[rel] = hash
		return nil
	})
	return hashes, err == nil
}

func fileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
