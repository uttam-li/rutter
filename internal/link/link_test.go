package link

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestApplyKeepsSourceAndRestoresBackup(t *testing.T) {
	root := t.TempDir()
	sourceDir := filepath.Join(root, "config")
	destDir := filepath.Join(root, "opencode")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "AGENTS.md"), []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "AGENTS.md"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := Manager{
		BackupRoot: filepath.Join(root, "backups"),
		Now:        func() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) },
	}
	planned, err := manager.Plan(sourceDir, destDir, []string{"AGENTS.md"})
	if err != nil {
		t.Fatalf("plan: %v", err)
	}
	applied, err := manager.Apply(context.Background(), planned, nil)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(applied) != 1 || applied[0].Backup == "" {
		t.Fatalf("applied = %+v, want one backed-up link", applied)
	}

	assertContent(t, filepath.Join(destDir, "AGENTS.md"), "new")

	checks := manager.Check(applied)
	if len(checks) != 1 || checks[0].State != "ok" {
		t.Fatalf("checks = %+v, want one ok check", checks)
	}

	if err := manager.Remove(context.Background(), applied); err != nil {
		t.Fatalf("remove: %v", err)
	}
	assertContent(t, filepath.Join(destDir, "AGENTS.md"), "old")
}

func TestApplyIsIdempotent(t *testing.T) {
	root := t.TempDir()
	sourceDir := filepath.Join(root, "config")
	destDir := filepath.Join(root, "opencode")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(sourceDir, "skills", "one"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "skills", "one", "SKILL.md"), []byte("skill"), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := Manager{BackupRoot: filepath.Join(root, "backups")}
	planned, err := manager.Plan(sourceDir, destDir, []string{"skills"})
	if err != nil {
		t.Fatalf("plan: %v", err)
	}
	first, err := manager.Apply(context.Background(), planned, nil)
	if err != nil {
		t.Fatalf("first apply: %v", err)
	}
	second, err := manager.Apply(context.Background(), planned, first)
	if err != nil {
		t.Fatalf("second apply: %v", err)
	}
	if second[0].Backup != "" {
		t.Fatalf("second apply created a backup: %+v", second[0])
	}
}

func TestPruneRemovesLinkGoneFromPlan(t *testing.T) {
	root := t.TempDir()
	sourceDir := filepath.Join(root, "config")
	destDir := filepath.Join(root, "opencode")
	if err := os.MkdirAll(filepath.Join(sourceDir, "skills", "one"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "skills", "one", "SKILL.md"), []byte("skill"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "AGENTS.md"), []byte("rules"), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := Manager{BackupRoot: filepath.Join(root, "backups")}
	planned, err := manager.Plan(sourceDir, destDir, []string{"AGENTS.md", "skills"})
	if err != nil {
		t.Fatalf("plan: %v", err)
	}
	applied, err := manager.Apply(context.Background(), planned, nil)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	current, err := manager.Plan(sourceDir, destDir, []string{"AGENTS.md"})
	if err != nil {
		t.Fatalf("plan current: %v", err)
	}
	if err := manager.Prune(context.Background(), applied, current); err != nil {
		t.Fatalf("prune: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(destDir, "skills")); !os.IsNotExist(err) {
		t.Fatalf("skills still exists after prune: %v", err)
	}
}

func assertContent(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(data) != want {
		t.Fatalf("%s = %q, want %q", path, string(data), want)
	}
}
