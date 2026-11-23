package app

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/uttam-li/rutter/internal/config"
	"github.com/uttam-li/rutter/internal/link"
	"github.com/uttam-li/rutter/internal/state"
)

func TestInstallSyncStatusUninstall(t *testing.T) {
	root := t.TempDir()
	repoDir := filepath.Join(root, "repo")
	configDir := filepath.Join(root, "opencode")
	dataDir := filepath.Join(root, "data")
	writeRepo(t, repoDir)

	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "opencode.jsonc"), []byte("user config"), 0o644); err != nil {
		t.Fatal(err)
	}

	source := &fakeSource{}
	opencode := &fakeOpenCode{configDir: configDir}
	store := &memoryStore{}
	service := Service{
		DataDir:  dataDir,
		RepoDir:  repoDir,
		Manifest: store,
		Source:   source,
		Merger:   config.Merger{},
		Links:    link.Manager{BackupRoot: filepath.Join(dataDir, "backups")},
		OpenCode: opencode,
		Now:      func() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) },
	}

	result, err := service.Install(context.Background(), InstallOptions{
		RepoURL:   "https://example.com/team/rutter.git",
		ConfigDir: configDir,
		Reload:    true,
	})
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if source.ensureCalls != 1 {
		t.Fatalf("ensure calls = %d, want 1", source.ensureCalls)
	}
	if !result.PersonalCreated {
		t.Fatal("expected a personal config to be created")
	}
	if opencode.reloads != 1 {
		t.Fatalf("reloads = %d, want 1", opencode.reloads)
	}
	if store.manifest.ConfigBackup == "" {
		t.Fatal("expected the existing config to be backed up")
	}

	assertGeneratedConfig(t, filepath.Join(configDir, "opencode.jsonc"), "me/model", 2)
	assertFile(t, filepath.Join(configDir, "AGENTS.md"), "# rules")
	assertFile(t, filepath.Join(configDir, "agents", "design.md"), "design")
	assertFile(t, filepath.Join(configDir, "skills", "foo", "SKILL.md"), "skill")

	status, err := service.Status(context.Background())
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if !status.Installed || len(status.Links) != len(ConfigItems) {
		t.Fatalf("status = %+v", status)
	}
	for _, check := range status.Links {
		if check.State != "ok" {
			t.Fatalf("link %s state = %s, want ok", check.Target, check.State)
		}
	}

	if err := os.WriteFile(filepath.Join(repoDir, "config", "opencode.base.json"), []byte(`{
		"$schema": "https://opencode.ai/config.json",
		"permissions": [{"action": "read", "resource": "*", "effect": "allow"}],
		"mcp": {"servers": {"mobbin": {"type": "remote", "url": "https://example.com/mcp"}}}
	}`), 0o644); err != nil {
		t.Fatal(err)
	}

	syncResult, err := service.Sync(context.Background(), true)
	if err != nil {
		t.Fatalf("sync: %v", err)
	}
	if source.syncCalls != 1 {
		t.Fatalf("sync calls = %d, want 1", source.syncCalls)
	}
	if !syncResult.Reloaded {
		t.Fatal("expected sync to reload OpenCode")
	}

	_, err = service.Uninstall(context.Background(), true)
	if err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	assertFile(t, filepath.Join(configDir, "opencode.jsonc"), "user config")
	if _, err := os.Lstat(filepath.Join(configDir, "AGENTS.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("AGENTS.md still exists: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(configDir, "skills")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("skills still exists: %v", err)
	}

	status, err = service.Status(context.Background())
	if err != nil {
		t.Fatalf("status after uninstall: %v", err)
	}
	if status.Installed {
		t.Fatal("status reports installed after uninstall")
	}
}

func TestInstallWithoutOpenCode(t *testing.T) {
	root := t.TempDir()
	repoDir := filepath.Join(root, "repo")
	writeRepo(t, repoDir)

	service := Service{
		DataDir:  filepath.Join(root, "data"),
		RepoDir:  repoDir,
		Manifest: &memoryStore{},
		Source:   &fakeSource{},
		Merger:   config.Merger{},
		Links:    link.Manager{BackupRoot: filepath.Join(root, "data", "backups")},
		OpenCode: &fakeOpenCode{configDirErr: errors.New("opencode not found")},
	}

	_, err := service.Install(context.Background(), InstallOptions{
		RepoURL: "https://example.com/team/rutter.git",
		Reload:  true,
	})
	if err == nil || !strings.Contains(err.Error(), "config-dir") {
		t.Fatalf("install error = %v, want a hint to pass --config-dir", err)
	}
}

type fakeSource struct {
	ensureCalls int
	syncCalls   int
	ensureErr   error
	syncErr     error
}

func (f *fakeSource) Ensure(_ context.Context, _, _ string) error {
	f.ensureCalls++
	return f.ensureErr
}

func (f *fakeSource) Sync(_ context.Context, _ string) error {
	f.syncCalls++
	return f.syncErr
}

type fakeOpenCode struct {
	configDir    string
	configDirErr error
	reloads      int
	reloadErr    error
}

func (f *fakeOpenCode) ConfigDir(_ context.Context) (string, error) {
	return f.configDir, f.configDirErr
}

func (f *fakeOpenCode) Reload(_ context.Context) error {
	f.reloads++
	return f.reloadErr
}

type memoryStore struct {
	manifest state.Manifest
	exists   bool
}

func (m *memoryStore) Load() (state.Manifest, error) {
	if !m.exists {
		return state.Manifest{}, state.ErrNotInstalled
	}
	return m.manifest, nil
}

func (m *memoryStore) Save(manifest state.Manifest) error {
	m.manifest = manifest
	m.exists = true
	return nil
}

func (m *memoryStore) Remove() error {
	m.exists = false
	return nil
}

func writeRepo(t *testing.T, repoDir string) {
	t.Helper()
	files := map[string]string{
		"config/opencode.base.json": `{
			"$schema": "https://opencode.ai/config.json",
			"permissions": [{"action": "read", "resource": "*", "effect": "allow"}],
			"mcp": {"servers": {"mobbin": {"type": "remote", "url": "https://example.com/mcp"}}}
		}`,
		"config/opencode.personal.example.json": `{
			"model": "me/model",
			"permissions": [{"action": "shell", "resource": "*git push*", "effect": "ask"}]
		}`,
		"config/AGENTS.md":           "# rules",
		"config/agents/design.md":    "design",
		"config/commands/test.md":    "command",
		"config/skills/foo/SKILL.md": "skill",
	}
	for name, body := range files {
		path := filepath.Join(repoDir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func assertGeneratedConfig(t *testing.T, path, model string, permissions int) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.HasPrefix(string(data), config.GeneratedHeader) {
		t.Fatalf("%s is missing the generated header", path)
	}
	var value map[string]any
	if err := json.Unmarshal([]byte(strings.TrimPrefix(string(data), config.GeneratedHeader)), &value); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	if value["model"] != model {
		t.Fatalf("model = %v, want %s", value["model"], model)
	}
	if got := len(value["permissions"].([]any)); got != permissions {
		t.Fatalf("permissions = %d, want %d", got, permissions)
	}
}

func assertFile(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(data) != want {
		t.Fatalf("%s = %q, want %q", path, string(data), want)
	}
}
