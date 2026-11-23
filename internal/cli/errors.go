package cli

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/uttam-li/rutter/internal/config"
	"github.com/uttam-li/rutter/internal/link"
	"github.com/uttam-li/rutter/internal/opencode"
	"github.com/uttam-li/rutter/internal/source"
	"github.com/uttam-li/rutter/internal/state"
)

// Exit codes are part of rutter's contract with scripts and agents. They are
// documented by `rutter agent` and must stay stable once released.
type exitCode int

const (
	exitOK           exitCode = 0
	exitFailure      exitCode = 1
	exitUsage        exitCode = 2
	exitNotInstalled exitCode = 3
	exitEnvironment  exitCode = 4
	exitDrift        exitCode = 5
	exitRepository   exitCode = 6
	exitConfig       exitCode = 7
)

// errEnvironment marks a failed environment or health check.
var errEnvironment = errors.New("environment check failed")

// usageError marks invalid commands and flags.
type usageError struct {
	err error
}

func (e usageError) Error() string { return e.err.Error() }
func (e usageError) Unwrap() error { return e.err }

// classify maps an error to an exit code and a stable machine label.
func classify(err error) (exitCode, string) {
	switch {
	case err == nil:
		return exitOK, "ok"
	case errors.Is(err, state.ErrNotInstalled):
		return exitNotInstalled, "not_installed"
	case errors.Is(err, link.ErrDrift):
		return exitDrift, "drift"
	case errors.Is(err, source.ErrRepository):
		return exitRepository, "repository"
	case errors.Is(err, config.ErrInvalid):
		return exitConfig, "config"
	case errors.Is(err, opencode.ErrUnavailable), errors.Is(err, errEnvironment):
		return exitEnvironment, "environment"
	}

	var usage usageError
	if errors.As(err, &usage) || isUsageMessage(err.Error()) {
		return exitUsage, "usage"
	}
	return exitFailure, "error"
}

func isUsageMessage(message string) bool {
	prefixes := []string{
		"unknown command",
		"unknown flag",
		"unknown shorthand flag",
		"invalid argument",
		"required flag",
		"accepts ",
		"flag needs an argument",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(message, prefix) {
			return true
		}
	}
	return false
}

// hintFor returns the next action for a failure, so both humans and agents can
// recover without reading the source.
func hintFor(label string) string {
	switch label {
	case "not_installed":
		return "run `rutter install` first"
	case "environment":
		return "run `rutter doctor` and fix the failing checks"
	case "drift":
		return "move or delete the conflicting path, then run `rutter sync`"
	case "repository":
		return "check the repository URL, network access, and git credentials"
	case "config":
		return "fix the JSON in config/opencode.base.json or opencode.personal.json"
	case "usage":
		return "run `rutter agent` for the command contract"
	default:
		return ""
	}
}

type errorDocument struct {
	Code     string `json:"code"`
	ExitCode int    `json:"exit_code"`
	Message  string `json:"message"`
	Hint     string `json:"hint,omitempty"`
}

type errorEnvelope struct {
	OK    bool          `json:"ok"`
	Error errorDocument `json:"error"`
}

// writeError renders a failure. In JSON mode it writes one object to stderr so
// stdout stays reserved for data.
func writeError(w io.Writer, code exitCode, label string, err error, asJSON bool) {
	hint := hintFor(label)
	if asJSON {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(errorEnvelope{
			OK: false,
			Error: errorDocument{
				Code:     label,
				ExitCode: int(code),
				Message:  err.Error(),
				Hint:     hint,
			},
		})
		return
	}

	line(w, "rutter: %v\n", err)
	if hint != "" {
		line(w, "hint: %s\n", hint)
	}
}
