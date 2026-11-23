package source

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ErrRepository reports a git or network failure.
var ErrRepository = errors.New("repository operation failed")

// Source keeps the config repository checkout up to date.
type Source interface {
	Ensure(ctx context.Context, url, dir string) error
	Sync(ctx context.Context, dir string) error
}

// Git uses the local git binary so existing SSH keys, credential helpers, and
// repository permissions keep working.
type Git struct {
	Binary string
}

func (g Git) binary() string {
	if g.Binary != "" {
		return g.Binary
	}
	return "git"
}

// Ensure clones the repository when dir is missing or empty, and pulls when it
// already holds a checkout.
func (g Git) Ensure(ctx context.Context, url, dir string) error {
	info, err := os.Stat(dir)
	switch {
	case err == nil && !info.IsDir():
		return fmt.Errorf("%s exists and is not a directory", dir)
	case err == nil:
		entries, readErr := os.ReadDir(dir)
		if readErr != nil {
			return fmt.Errorf("read %s: %w", dir, readErr)
		}
		if len(entries) == 0 {
			return g.clone(ctx, url, dir)
		}
		if _, statErr := os.Stat(filepath.Join(dir, ".git")); statErr != nil {
			return fmt.Errorf("%s exists but is not a git checkout", dir)
		}
		return g.Sync(ctx, dir)
	case !os.IsNotExist(err):
		return fmt.Errorf("inspect %s: %w", dir, err)
	}

	if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
		return fmt.Errorf("create parent directory: %w", err)
	}
	return g.clone(ctx, url, dir)
}

func (g Git) clone(ctx context.Context, url, dir string) error {
	if strings.TrimSpace(url) == "" {
		return fmt.Errorf("config repository URL is empty")
	}
	cmd := exec.CommandContext(ctx, g.binary(), "clone", "--depth", "1", url, dir)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: git clone %s: %w (%s)", ErrRepository, url, err, strings.TrimSpace(string(out)))
	}
	return nil
}

// Sync pulls fast-forward only. Local changes are never discarded.
func (g Git) Sync(ctx context.Context, dir string) error {
	cmd := exec.CommandContext(ctx, g.binary(), "-C", dir, "pull", "--ff-only")
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: git pull in %s: %w (%s)", ErrRepository, dir, err, strings.TrimSpace(string(out)))
	}
	return nil
}
