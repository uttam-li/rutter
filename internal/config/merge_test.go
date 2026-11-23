package config

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMergeAppendsConfiguredArrays(t *testing.T) {
	base := []byte(`{
		"model": "team/model",
		"permissions": [{"action": "read"}],
		"agents": {"plan": {"permissions": [{"action": "edit"}]}}
	}`)
	personal := []byte(`{
		"model": "me/model",
		"permissions": [{"action": "shell"}],
		"agents": {"plan": {"permissions": [{"action": "shell"}]}}
	}`)

	out, err := (Merger{}).Merge(base, personal)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}

	got := decodeGenerated(t, out)
	if got["model"] != "me/model" {
		t.Fatalf("model = %v, want me/model", got["model"])
	}

	permissions := got["permissions"].([]any)
	if len(permissions) != 2 {
		t.Fatalf("permissions length = %d, want 2", len(permissions))
	}

	agents := got["agents"].(map[string]any)
	plan := agents["plan"].(map[string]any)
	planPermissions := plan["permissions"].([]any)
	if len(planPermissions) != 2 {
		t.Fatalf("plan permissions length = %d, want 2", len(planPermissions))
	}
}

func TestMergeReplacesNormalArrays(t *testing.T) {
	base := []byte(`{"mcp":{"servers":{"browse":{"command":["old","mcp"]}}}}`)
	personal := []byte(`{"mcp":{"servers":{"browse":{"command":["new","mcp"]}}}}`)

	out, err := (Merger{}).Merge(base, personal)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}

	got := decodeGenerated(t, out)
	command := got["mcp"].(map[string]any)["servers"].(map[string]any)["browse"].(map[string]any)["command"].([]any)
	if len(command) != 2 || command[0] != "new" {
		t.Fatalf("command = %v, want [new mcp]", command)
	}
}

func TestMergeAcceptsJSONC(t *testing.T) {
	base := []byte(`{
		// team default
		"model": "team/model",
	}`)
	personal := []byte(`{
		"model": "me/model", // local
	}`)

	out, err := (Merger{}).Merge(base, personal)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	if !bytes.HasPrefix(out, []byte(GeneratedHeader)) {
		t.Fatal("merged output is missing the generated header")
	}
	if got := decodeGenerated(t, out); got["model"] != "me/model" {
		t.Fatalf("model = %v, want me/model", got["model"])
	}
}

func TestMergeFilesMissingPersonal(t *testing.T) {
	dir := t.TempDir()
	basePath := filepath.Join(dir, "base.json")
	if err := os.WriteFile(basePath, []byte(`{"model":"team/model"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := (Merger{}).MergeFiles(basePath, filepath.Join(dir, "missing.json"))
	if err != nil {
		t.Fatalf("merge files: %v", err)
	}
	if got := decodeGenerated(t, out); got["model"] != "team/model" {
		t.Fatalf("model = %v, want team/model", got["model"])
	}
}

func TestMergeRejectsArrayRoot(t *testing.T) {
	if _, err := (Merger{}).Merge([]byte(`[]`), []byte(`{}`)); err == nil {
		t.Fatal("expected an error for an array root")
	}
}

func decodeGenerated(t *testing.T, out []byte) map[string]any {
	t.Helper()
	body := bytes.TrimPrefix(out, []byte(GeneratedHeader))
	var value map[string]any
	if err := json.Unmarshal(body, &value); err != nil {
		t.Fatalf("decode generated config: %v", err)
	}
	return value
}
