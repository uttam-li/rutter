package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/uttam-li/rutter/internal/config"
	"github.com/uttam-li/rutter/internal/link"
	"github.com/uttam-li/rutter/internal/platform"
	"github.com/uttam-li/rutter/internal/state"
)

// ConfigItems are the entries linked from the config repository into the
// OpenCode config directory.
var ConfigItems = []string{"AGENTS.md", "agents", "commands", "skills"}

// Source keeps the config repository checkout current.
type Source interface {
	Ensure(ctx context.Context, url, dir string) error
	Sync(ctx context.Context, dir string) error
}

// Merger joins the team base and personal config files.
type Merger interface {
	MergeFiles(basePath, personalPath string) ([]byte, error)
}

// Linker installs and removes config entries.
type Linker interface {
	Plan(root, dest string, items []string) ([]link.Link, error)
	Apply(ctx context.Context, links, previous []link.Link) ([]link.Link, error)
	Prune(ctx context.Context, previous, current []link.Link) error
	Remove(ctx context.Context, links []link.Link) error
	Check(links []link.Link) []link.Check
}

// OpenCode discovers the config directory and asks the service to reload.
type OpenCode interface {
	ConfigDir(ctx context.Context) (string, error)
	Reload(ctx context.Context) error
}

// Store persists the installation manifest.
type Store interface {
	Load() (state.Manifest, error)
	Save(state.Manifest) error
	Remove() error
}

// Service is the deep module the CLI talks to. It hides repository syncing,
// config merging, file linking, backups, and OpenCode integration.
type Service struct {
	DataDir  string
	RepoDir  string
	Manifest Store
	Source   Source
	Merger   Merger
	Links    Linker
	OpenCode OpenCode
	Now      func() time.Time
}

// InstallOptions controls a fresh install.
type InstallOptions struct {
	RepoURL   string
	ConfigDir string
	Reload    bool
}

// Result describes a completed install or sync.
type Result struct {
	ConfigDir       string      `json:"config_dir"`
	RepoDir         string      `json:"repo_dir"`
	ConfigPath      string      `json:"config_path"`
	PersonalPath    string      `json:"personal_path"`
	PersonalCreated bool        `json:"personal_created"`
	Links           []link.Link `json:"links"`
	Reloaded        bool        `json:"reloaded"`
	Warnings        []string    `json:"warnings,omitempty"`
}

// Status is the read-only view used by `rutter status`.
type Status struct {
	Installed       bool         `json:"installed"`
	RepoURL         string       `json:"repo_url,omitempty"`
	RepoDir         string       `json:"repo_dir,omitempty"`
	ConfigDir       string       `json:"config_dir,omitempty"`
	ConfigPath      string       `json:"config_path,omitempty"`
	PersonalPath    string       `json:"personal_path,omitempty"`
	ConfigGenerated bool         `json:"config_generated,omitempty"`
	Links           []link.Check `json:"links,omitempty"`
}

// DoctorCheck is one diagnostic result.
type DoctorCheck struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

// DoctorReport is the result of `rutter doctor`.
type DoctorReport struct {
	OK     bool          `json:"ok"`
	Checks []DoctorCheck `json:"checks"`
}

// Install clones or updates the config repo, writes the merged config, links
// the shared entries, and records a manifest.
func (s Service) Install(ctx context.Context, opts InstallOptions) (Result, error) {
	if opts.RepoURL == "" {
		return Result{}, errors.New("config repository URL is required")
	}
	if err := s.Source.Ensure(ctx, opts.RepoURL, s.RepoDir); err != nil {
		return Result{}, err
	}

	configDir := opts.ConfigDir
	if configDir == "" {
		discovered, err := s.OpenCode.ConfigDir(ctx)
		if err != nil {
			return Result{}, fmt.Errorf("find OpenCode config directory: %w; pass --config-dir to override", err)
		}
		configDir = discovered
	}

	now := s.now()
	manifest := state.Manifest{
		Version:      1,
		InstalledAt:  now,
		UpdatedAt:    now,
		RepoURL:      opts.RepoURL,
		RepoDir:      s.RepoDir,
		ConfigDir:    configDir,
		PersonalPath: filepath.Join(configDir, "opencode.personal.json"),
		ConfigPath:   filepath.Join(configDir, "opencode.jsonc"),
	}
	if previous, err := s.Manifest.Load(); err == nil {
		manifest.InstalledAt = previous.InstalledAt
		manifest.ConfigBackup = previous.ConfigBackup
		manifest.Links = previous.Links
	} else if !errors.Is(err, state.ErrNotInstalled) {
		return Result{}, err
	}

	return s.apply(ctx, manifest, opts.Reload)
}

// Sync pulls the config repo and reapplies it using the saved manifest.
func (s Service) Sync(ctx context.Context, reload bool) (Result, error) {
	manifest, err := s.Manifest.Load()
	if err != nil {
		if errors.Is(err, state.ErrNotInstalled) {
			return Result{}, errors.New("rutter is not installed; run `rutter install` first")
		}
		return Result{}, err
	}

	if err := s.Source.Sync(ctx, manifest.RepoDir); err != nil {
		return Result{}, err
	}
	return s.apply(ctx, manifest, reload)
}

func (s Service) apply(ctx context.Context, manifest state.Manifest, reload bool) (Result, error) {
	repoConfigDir := filepath.Join(manifest.RepoDir, "config")
	basePath := filepath.Join(repoConfigDir, "opencode.base.json")
	examplePath := filepath.Join(repoConfigDir, "opencode.personal.example.json")

	if err := os.MkdirAll(manifest.ConfigDir, 0o755); err != nil {
		return Result{}, fmt.Errorf("create config directory %s: %w", manifest.ConfigDir, err)
	}

	result := Result{
		ConfigDir:    manifest.ConfigDir,
		RepoDir:      manifest.RepoDir,
		ConfigPath:   manifest.ConfigPath,
		PersonalPath: manifest.PersonalPath,
	}

	if _, err := os.Stat(manifest.PersonalPath); errors.Is(err, os.ErrNotExist) {
		data, readErr := os.ReadFile(examplePath)
		if readErr != nil {
			return Result{}, fmt.Errorf("read personal config example: %w", readErr)
		}
		if err := platform.WriteFileAtomic(manifest.PersonalPath, data, 0o644); err != nil {
			return Result{}, fmt.Errorf("create personal config: %w", err)
		}
		result.PersonalCreated = true
	} else if err != nil {
		return Result{}, fmt.Errorf("inspect personal config: %w", err)
	}

	merged, err := s.Merger.MergeFiles(basePath, manifest.PersonalPath)
	if err != nil {
		return Result{}, err
	}

	if manifest.ConfigBackup == "" {
		current, readErr := os.ReadFile(manifest.ConfigPath)
		if readErr == nil && !bytes.HasPrefix(current, []byte(config.GeneratedHeader)) {
			backupDir := filepath.Join(s.DataDir, "backups", s.timestamp())
			backupPath := filepath.Join(backupDir, "opencode.jsonc")
			if err := link.Move(manifest.ConfigPath, backupPath); err != nil {
				return Result{}, fmt.Errorf("back up existing config: %w", err)
			}
			manifest.ConfigBackup = backupPath
		}
	}

	if err := platform.WriteFileAtomic(manifest.ConfigPath, merged, 0o644); err != nil {
		return Result{}, fmt.Errorf("write merged config: %w", err)
	}

	planned, err := s.Links.Plan(repoConfigDir, manifest.ConfigDir, ConfigItems)
	if err != nil {
		return Result{}, err
	}
	applied, err := s.Links.Apply(ctx, planned, manifest.Links)
	if err != nil {
		return Result{}, err
	}
	if err := s.Links.Prune(ctx, manifest.Links, applied); err != nil {
		return Result{}, err
	}

	manifest.Links = applied
	manifest.UpdatedAt = s.now()
	if err := s.Manifest.Save(manifest); err != nil {
		return Result{}, err
	}

	result.Links = applied
	if reload {
		if err := s.OpenCode.Reload(ctx); err != nil {
			result.Warnings = append(result.Warnings, err.Error())
		} else {
			result.Reloaded = true
		}
	}
	return result, nil
}

// Uninstall removes managed links, restores backups, and deletes the manifest.
func (s Service) Uninstall(ctx context.Context, purge bool) (Result, error) {
	manifest, err := s.Manifest.Load()
	if err != nil {
		if errors.Is(err, state.ErrNotInstalled) {
			return Result{}, errors.New("rutter is not installed")
		}
		return Result{}, err
	}

	result := Result{
		ConfigDir:    manifest.ConfigDir,
		RepoDir:      manifest.RepoDir,
		ConfigPath:   manifest.ConfigPath,
		PersonalPath: manifest.PersonalPath,
	}

	if err := s.Links.Remove(ctx, manifest.Links); err != nil {
		return result, err
	}

	current, readErr := os.ReadFile(manifest.ConfigPath)
	if readErr == nil && bytes.HasPrefix(current, []byte(config.GeneratedHeader)) {
		if err := os.Remove(manifest.ConfigPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return result, fmt.Errorf("remove generated config: %w", err)
		}
	}

	if manifest.ConfigBackup != "" {
		if _, err := os.Stat(manifest.ConfigBackup); err == nil {
			if err := os.MkdirAll(filepath.Dir(manifest.ConfigPath), 0o755); err != nil {
				return result, fmt.Errorf("create config directory: %w", err)
			}
			if err := link.Move(manifest.ConfigBackup, manifest.ConfigPath); err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("restore config backup: %v", err))
			}
		}
	}

	if err := s.Manifest.Remove(); err != nil {
		return result, err
	}

	if purge {
		if err := os.RemoveAll(manifest.RepoDir); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("remove repository: %v", err))
		}
		if err := os.RemoveAll(s.DataDir); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("remove data directory: %v", err))
		}
	}

	if err := s.OpenCode.Reload(ctx); err != nil {
		result.Warnings = append(result.Warnings, err.Error())
	} else {
		result.Reloaded = true
	}
	return result, nil
}

// Status reports the installed manifest and whether links still match.
func (s Service) Status(ctx context.Context) (Status, error) {
	manifest, err := s.Manifest.Load()
	if err != nil {
		if errors.Is(err, state.ErrNotInstalled) {
			return Status{Installed: false}, nil
		}
		return Status{}, err
	}

	status := Status{
		Installed:    true,
		RepoURL:      manifest.RepoURL,
		RepoDir:      manifest.RepoDir,
		ConfigDir:    manifest.ConfigDir,
		ConfigPath:   manifest.ConfigPath,
		PersonalPath: manifest.PersonalPath,
	}
	if data, readErr := os.ReadFile(manifest.ConfigPath); readErr == nil {
		status.ConfigGenerated = bytes.HasPrefix(data, []byte(config.GeneratedHeader))
	}
	status.Links = s.Links.Check(manifest.Links)
	return status, nil
}

// Doctor checks the local environment and installation.
func (s Service) Doctor(ctx context.Context) DoctorReport {
	report := DoctorReport{OK: true}
	add := func(name, status, detail string) {
		if status == "error" {
			report.OK = false
		}
		report.Checks = append(report.Checks, DoctorCheck{Name: name, Status: status, Detail: detail})
	}

	if path, err := exec.LookPath("git"); err != nil {
		add("git", "error", "git not found in PATH")
	} else {
		add("git", "ok", path)
	}

	if path, err := exec.LookPath("opencode"); err != nil {
		add("opencode", "warn", "opencode not found in PATH; pass --config-dir during install")
	} else {
		add("opencode", "ok", path)
	}

	manifest, err := s.Manifest.Load()
	if err != nil {
		if errors.Is(err, state.ErrNotInstalled) {
			add("installation", "warn", "rutter is not installed")
			return report
		}
		add("installation", "error", err.Error())
		return report
	}
	add("installation", "ok", manifest.ConfigDir)

	if _, err := os.Stat(manifest.RepoDir); err != nil {
		add("repository", "error", err.Error())
	} else {
		add("repository", "ok", manifest.RepoDir)
	}

	if _, err := os.Stat(manifest.ConfigPath); err != nil {
		add("config", "error", err.Error())
	} else {
		add("config", "ok", manifest.ConfigPath)
	}

	if _, err := s.Merger.MergeFiles(filepath.Join(manifest.RepoDir, "config", "opencode.base.json"), manifest.PersonalPath); err != nil {
		add("config parse", "error", err.Error())
	} else {
		add("config parse", "ok", "base and personal configs merge")
	}

	for _, check := range s.Links.Check(manifest.Links) {
		status := "ok"
		if check.State != "ok" {
			status = "error"
		}
		add("link "+filepath.Base(check.Target), status, check.State+" "+check.Detail)
	}

	return report
}

func (s Service) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}

func (s Service) timestamp() string {
	return s.now().Format("20060102-150405")
}
