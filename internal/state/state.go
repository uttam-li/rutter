package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/uttam-li/rutter/internal/link"
	"github.com/uttam-li/rutter/internal/platform"
)

var ErrNotInstalled = errors.New("rutter is not installed")

// Manifest records what rutter installed so uninstall can restore the machine.
type Manifest struct {
	Version      int         `json:"version"`
	InstalledAt  time.Time   `json:"installed_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	RepoURL      string      `json:"repo_url"`
	RepoDir      string      `json:"repo_dir"`
	ConfigDir    string      `json:"config_dir"`
	PersonalPath string      `json:"personal_path"`
	ConfigPath   string      `json:"config_path"`
	ConfigBackup string      `json:"config_backup,omitempty"`
	Links        []link.Link `json:"links"`
}

// Store persists the manifest as JSON.
type Store struct {
	Path string
}

func (s Store) Load() (Manifest, error) {
	data, err := os.ReadFile(s.Path)
	if errors.Is(err, os.ErrNotExist) {
		return Manifest{}, ErrNotInstalled
	}
	if err != nil {
		return Manifest{}, fmt.Errorf("read manifest: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, fmt.Errorf("parse manifest: %w", err)
	}
	if manifest.Version == 0 {
		manifest.Version = 1
	}
	return manifest, nil
}

func (s Store) Save(manifest Manifest) error {
	manifest.Version = 1
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("encode manifest: %w", err)
	}
	if err := platform.WriteFileAtomic(s.Path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	return nil
}

func (s Store) Remove() error {
	err := os.Remove(s.Path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("remove manifest: %w", err)
	}
	return nil
}
