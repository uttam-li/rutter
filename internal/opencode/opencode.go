package opencode

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

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
		return "", fmt.Errorf("run %s debug paths: %w", c.binary(), err)
	}
	return ParseConfigDir(out)
}

// Reload asks the running OpenCode service to reload configuration.
func (c CLI) Reload(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, c.binary(), "reload")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("run %s reload: %w: %s", c.binary(), err, strings.TrimSpace(string(out)))
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
