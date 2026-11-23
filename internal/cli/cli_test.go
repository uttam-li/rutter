package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/uttam-li/rutter/internal/app"
	"github.com/uttam-li/rutter/internal/link"
	"github.com/uttam-li/rutter/internal/state"
)

func TestAgentSpecJSON(t *testing.T) {
	out, stderr, code := runCLI(t, app.Service{}, "agent", "--json")
	if code != int(exitOK) {
		t.Fatalf("exit = %d, stderr = %s", code, stderr)
	}

	var envelope struct {
		OK   bool      `json:"ok"`
		Data agentSpec `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &envelope); err != nil {
		t.Fatalf("decode agent spec: %v\n%s", err, out)
	}
	if !envelope.OK {
		t.Fatal("agent spec is not marked ok")
	}
	if !envelope.Data.NonInteractive {
		t.Fatal("agent spec must declare non_interactive")
	}

	paths := map[string]bool{}
	for _, command := range envelope.Data.Commands {
		paths[command.Path] = true
	}
	for _, want := range []string{"install", "sync", "uninstall", "status", "doctor", "config path", "config show", "version", "agent"} {
		if !paths[want] {
			t.Errorf("agent spec is missing command %q", want)
		}
	}

	codes := map[int]bool{}
	for _, item := range envelope.Data.ExitCodes {
		codes[item.Code] = true
	}
	for _, want := range []int{0, 2, 3, 5, 6, 7} {
		if !codes[want] {
			t.Errorf("agent spec is missing exit code %d", want)
		}
	}
}

func TestAgentGuide(t *testing.T) {
	out, _, code := runCLI(t, app.Service{}, "agent")
	if code != int(exitOK) {
		t.Fatalf("exit = %d", code)
	}
	for _, want := range []string{"EXIT CODES", "rutter install", "RUTTER_JSON", "not_installed"} {
		if !strings.Contains(out, want) {
			t.Errorf("agent guide is missing %q", want)
		}
	}
}

func TestRootJSONIsAgentSpec(t *testing.T) {
	out, stderr, code := runCLI(t, app.Service{}, "--json")
	if code != int(exitOK) {
		t.Fatalf("exit = %d, stderr = %s", code, stderr)
	}
	if !strings.Contains(out, `"non_interactive": true`) || !strings.Contains(out, `"commands"`) {
		t.Fatalf("stdout = %s, want the agent spec", out)
	}
}

func TestUnknownCommandIsUsageJSON(t *testing.T) {
	_, stderr, code := runCLI(t, app.Service{}, "bogus", "--json")
	if code != int(exitUsage) {
		t.Fatalf("exit = %d, want %d", code, int(exitUsage))
	}
	for _, want := range []string{`"ok": false`, `"code": "usage"`, `"exit_code": 2`} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr = %s, want %q", stderr, want)
		}
	}
}

func TestUnknownFlagIsUsageJSON(t *testing.T) {
	_, stderr, code := runCLI(t, app.Service{}, "status", "--nope", "--json")
	if code != int(exitUsage) {
		t.Fatalf("exit = %d, want %d", code, int(exitUsage))
	}
	if !strings.Contains(stderr, `"code": "usage"`) {
		t.Fatalf("stderr = %s", stderr)
	}
}

func TestVersionJSON(t *testing.T) {
	out, stderr, code := runCLI(t, app.Service{}, "version", "--json")
	if code != int(exitOK) {
		t.Fatalf("exit = %d, stderr = %s", code, stderr)
	}
	if !strings.Contains(out, `"name": "rutter"`) {
		t.Fatalf("version output = %s", out)
	}
}

func TestStatusCheckNotInstalled(t *testing.T) {
	service := app.Service{Manifest: notInstalledStore{}}

	out, _, code := runCLI(t, service, "status", "--check")
	if code != int(exitNotInstalled) {
		t.Fatalf("exit = %d, want %d", code, int(exitNotInstalled))
	}
	if !strings.Contains(out, "not installed") {
		t.Fatalf("stdout = %s", out)
	}
}

func TestStatusNotInstalledJSONHealthyFalse(t *testing.T) {
	service := app.Service{Manifest: notInstalledStore{}}

	out, _, code := runCLI(t, service, "status", "--json")
	if code != int(exitOK) {
		t.Fatalf("exit = %d", code)
	}
	if !strings.Contains(out, `"healthy": false`) {
		t.Fatalf("stdout = %s", out)
	}
}

func TestJSONEnv(t *testing.T) {
	t.Setenv("RUTTER_JSON", "1")

	out, _, code := runCLI(t, app.Service{}, "version")
	if code != int(exitOK) {
		t.Fatalf("exit = %d", code)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Fatalf("stdout = %s, want JSON", out)
	}
}

func TestStatusHealthy(t *testing.T) {
	ok := app.Status{Installed: true, ConfigGenerated: true, Links: []link.Check{{State: "ok"}}}
	if !statusHealthy(ok) {
		t.Fatal("healthy status reported unhealthy")
	}

	for name, status := range map[string]app.Status{
		"not installed": {},
		"config drift":  {Installed: true, Links: []link.Check{{State: "ok"}}},
		"missing link":  {Installed: true, ConfigGenerated: true, Links: []link.Check{{State: "missing"}}},
		"no links":      {Installed: true, ConfigGenerated: true},
	} {
		if statusHealthy(status) {
			t.Errorf("%s reported healthy", name)
		}
	}
}

func TestClassify(t *testing.T) {
	tests := []struct {
		err   error
		code  exitCode
		label string
	}{
		{nil, exitOK, "ok"},
		{state.ErrNotInstalled, exitNotInstalled, "not_installed"},
		{link.ErrDrift, exitDrift, "drift"},
		{errEnvironment, exitEnvironment, "environment"},
		{usageError{err: errors.New("unknown flag")}, exitUsage, "usage"},
		{errors.New("boom"), exitFailure, "error"},
	}
	for _, test := range tests {
		code, label := classify(test.err)
		if code != test.code || label != test.label {
			t.Errorf("classify(%v) = (%d, %s), want (%d, %s)", test.err, code, label, test.code, test.label)
		}
	}
}

type notInstalledStore struct{}

func (notInstalledStore) Load() (state.Manifest, error) {
	return state.Manifest{}, state.ErrNotInstalled
}

func (notInstalledStore) Save(state.Manifest) error { return nil }
func (notInstalledStore) Remove() error             { return nil }

func runCLI(t *testing.T, service app.Service, args ...string) (string, string, int) {
	t.Helper()
	jsonDefault = false
	t.Cleanup(func() { jsonDefault = false })

	var out, errOut bytes.Buffer
	code := runRoot(newRoot(service), args, &out, &errOut)
	return out.String(), errOut.String(), code
}
