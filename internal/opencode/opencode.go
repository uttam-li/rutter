package opencode

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// ErrUnavailable reports that the OpenCode CLI could not be used.
var ErrUnavailable = errors.New("opencode unavailable")

// Client talks to the installed OpenCode CLI.
type Client interface {
	ConfigDir(ctx context.Context) (string, error)
	Reload(ctx context.Context) error
}

// CLI is the production OpenCode client.
type CLI struct {
	Binary string
}

func (c CLI) binary() string {
	if c.Binary != "" {
		return c.Binary
	}
	return "opencode"
}

// ConfigDir asks OpenCode where its global config directory lives. Deriving
// this from the installed binary avoids hardcoding Unix paths on Windows.
func (c CLI) ConfigDir(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, c.binary(), "debug", "paths")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%w: run %s debug paths: %w", ErrUnavailable, c.binary(), err)
	}
	dir, err := ParseConfigDir(out)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrUnavailable, err)
	}
	return dir, nil
}

// Reload asks the running OpenCode service to reload configuration.
func (c CLI) Reload(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, c.binary(), "reload")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: run %s reload: %w (%s)", ErrUnavailable, c.binary(), err, strings.TrimSpace(string(out)))
	}
	return nil
}

// ParseConfigDir extracts the config path from `opencode debug paths`.
func ParseConfigDir(out []byte) (string, error) {
	for _, raw := range strings.Split(string(out), "\n") {
		line := strings.TrimSpace(raw)
		fields := strings.Fields(line)
		if len(fields) < 2 || fields[0] != "config" {
			continue
		}

		index := strings.Index(line, "config")
		path := strings.TrimSpace(line[index+len("config"):])
		if path == "" {
			continue
		}
		return path, nil
	}

	return "", fmt.Errorf("config path not found in opencode debug paths output")
}
