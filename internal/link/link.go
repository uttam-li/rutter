package link

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrDrift reports a managed path that was changed outside rutter.
var ErrDrift = errors.New("managed path changed outside rutter")

// Kind describes whether a linked item is a file or a directory.
type Kind string

const (
	KindFile Kind = "file"
	KindDir  Kind = "dir"
)

// Mode records how rutter materialised a link.
type Mode string

const (
	ModeLink Mode = "link"
	ModeCopy Mode = "copy"
)

// Link is one item rutter places in the OpenCode config directory.
type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Kind   Kind   `json:"kind"`
	Mode   Mode   `json:"mode,omitempty"`
	Backup string `json:"backup,omitempty"`
}

// Check reports the state of one installed link.
type Check struct {
	Target string `json:"target"`
	Source string `json:"source"`
	State  string `json:"state"`
	Detail string `json:"detail,omitempty"`
}

// Manager installs and removes links. It owns backups, platform differences,
// and the guarantee that user files are moved aside rather than deleted.
type Manager struct {
	BackupRoot string
	Now        func() time.Time
}

// Plan builds links for items found under root.
func (m Manager) Plan(root, dest string, items []string) ([]Link, error) {
	links := make([]Link, 0, len(items))
	for _, item := range items {
		source := filepath.Join(root, item)
		info, err := os.Stat(source)
		if err != nil {
			return nil, fmt.Errorf("inspect %s: %w", source, err)
		}

		kind := KindFile
		if info.IsDir() {
			kind = KindDir
		}
		links = append(links, Link{
			Source: source,
			Target: filepath.Join(dest, item),
			Kind:   kind,
		})
	}
	return links, nil
}

// Apply creates or refreshes links. Existing files not owned by rutter are
// moved into the backup directory and recorded on the link.
func (m Manager) Apply(ctx context.Context, links, previous []Link) ([]Link, error) {
	previousByTarget := make(map[string]Link, len(previous))
	for _, old := range previous {
		previousByTarget[old.Target] = old
	}

	backupDir := filepath.Join(m.BackupRoot, m.timestamp())
	applied := make([]Link, 0, len(links))

	for _, item := range links {
		if err := ctx.Err(); err != nil {
			return applied, err
		}

		old, wasOurs := previousByTarget[item.Target]
		exists, err := pathExists(item.Target)
		if err != nil {
			return applied, err
		}
		if wasOurs && old.Source == item.Source {
			item.Backup = old.Backup
		}

		if exists {
			owned := wasOurs && old.Source == item.Source && old.Mode != ""
			if owned && old.Mode == ModeLink && pointsTo(item.Target, item.Source) {
				item.Mode = old.Mode
				applied = append(applied, item)
				continue
			}
			if owned && old.Mode == ModeCopy {
				if sameContent(item.Source, item.Target) {
					item.Mode = old.Mode
					applied = append(applied, item)
					continue
				}
				// A copy rutter owns is disposable. Replacing it keeps the
				// original pre-install backup that uninstall must restore.
				if err := removeLink(item.Target, old.Mode); err != nil {
					return applied, fmt.Errorf("replace copy %s: %w", item.Target, err)
				}
			} else {
				// Anything else is user data or drift. Preserve it.
				backup, err := m.backupPath(backupDir, item.Target)
				if err != nil {
					return applied, err
				}
				if err := movePath(item.Target, backup); err != nil {
					return applied, fmt.Errorf("back up %s: %w", item.Target, err)
				}
				item.Backup = backup
			}
		}

		if err := os.MkdirAll(filepath.Dir(item.Target), 0o755); err != nil {
			return applied, fmt.Errorf("create parent for %s: %w", item.Target, err)
		}
		mode, err := createLink(item.Source, item.Target, item.Kind)
		if err != nil {
			return applied, fmt.Errorf("install %s: %w", item.Target, err)
		}
		item.Mode = mode
		applied = append(applied, item)
	}

	return applied, nil
}

// Prune removes links that existed in the previous manifest but are gone from
// the new plan. Backups are restored when possible.
func (m Manager) Prune(ctx context.Context, previous, current []Link) error {
	currentTargets := make(map[string]bool, len(current))
	for _, item := range current {
		currentTargets[item.Target] = true
	}

	for _, old := range previous {
		if currentTargets[old.Target] {
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := m.removeOne(old, true); err != nil {
			return err
		}
	}
	return nil
}

// Remove uninstalls every link and restores backups.
func (m Manager) Remove(ctx context.Context, links []Link) error {
	for _, item := range links {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := m.removeOne(item, false); err != nil {
			return err
		}
	}
	return nil
}

func (m Manager) removeOne(item Link, allowMissing bool) error {
	exists, err := pathExists(item.Target)
	if err != nil {
		return err
	}

	if exists {
		owned := false
		switch item.Mode {
		case ModeLink:
			owned = pointsTo(item.Target, item.Source)
		case ModeCopy:
			owned = sameContent(item.Source, item.Target)
		default:
			owned = false
		}
		if owned {
			if err := removeLink(item.Target, item.Mode); err != nil {
				return fmt.Errorf("remove %s: %w", item.Target, err)
			}
		} else if !allowMissing {
			return fmt.Errorf("%w: %s", ErrDrift, item.Target)
		}
	}

	if item.Backup != "" {
		if _, err := os.Stat(item.Backup); err == nil {
			if err := os.MkdirAll(filepath.Dir(item.Target), 0o755); err != nil {
				return fmt.Errorf("create parent for restore: %w", err)
			}
			if err := movePath(item.Backup, item.Target); err != nil {
				return fmt.Errorf("restore %s: %w", item.Target, err)
			}
		}
	}
	return nil
}

// Check reports whether installed links still point at their sources.
func (m Manager) Check(links []Link) []Check {
	checks := make([]Check, 0, len(links))
	for _, item := range links {
		check := Check{Target: item.Target, Source: item.Source}
		exists, err := pathExists(item.Target)
		switch {
		case err != nil:
			check.State = "error"
			check.Detail = err.Error()
		case !exists:
			check.State = "missing"
		case item.Mode == ModeLink:
			if pointsTo(item.Target, item.Source) {
				check.State = "ok"
			} else {
				check.State = "drift"
				check.Detail = "link points elsewhere"
			}
		case item.Mode == ModeCopy:
			if sameContent(item.Source, item.Target) {
				check.State = "ok"
			} else {
				check.State = "drift"
				check.Detail = "copy differs from source"
			}
		default:
			check.State = "untracked"
		}
		checks = append(checks, check)
	}
	return checks
}

func (m Manager) timestamp() string {
	now := time.Now
	if m.Now != nil {
		now = m.Now
	}
	return now().Format("20060102-150405")
}

func (m Manager) backupPath(backupDir, target string) (string, error) {
	if m.BackupRoot == "" {
		return "", errors.New("backup directory is not configured")
	}
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return "", fmt.Errorf("create backup directory: %w", err)
	}
	return filepath.Join(backupDir, filepath.Base(target)), nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("inspect %s: %w", path, err)
}
