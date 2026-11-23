package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/uttam-li/rutter/internal/app"
	"github.com/uttam-li/rutter/internal/config"
	"github.com/uttam-li/rutter/internal/link"
	"github.com/uttam-li/rutter/internal/opencode"
	"github.com/uttam-li/rutter/internal/platform"
	"github.com/uttam-li/rutter/internal/source"
	"github.com/uttam-li/rutter/internal/state"
	"github.com/uttam-li/rutter/internal/version"
)

// jsonDefault is set from RUTTER_JSON or an explicit --json before Cobra
// parses anything, so even flag errors can be reported as JSON.
var jsonDefault bool

// Execute runs the CLI and returns the process exit code.
func Execute() int {
	args := os.Args[1:]
	jsonDefault = envJSON() || argsContainJSON(args)

	service, err := buildService()
	if err != nil {
		wrapped := fmt.Errorf("%w: %w", errEnvironment, err)
		code, label := classify(wrapped)
		writeError(os.Stderr, code, label, wrapped, jsonDefault)
		return int(code)
	}

	return runRoot(newRoot(service), args, os.Stdout, os.Stderr)
}

func runRoot(root *cobra.Command, args []string, out, errOut io.Writer) int {
	jsonDefault = jsonDefault || envJSON()
	root.SetArgs(args)
	root.SetOut(out)
	root.SetErr(errOut)

	err := root.ExecuteContext(context.Background())
	if err == nil {
		return int(exitOK)
	}

	code, label := classify(err)
	writeError(errOut, code, label, err, jsonDefault || argsContainJSON(args))
	return int(code)
}

func buildService() (app.Service, error) {
	dataDir := os.Getenv("RUTTER_DATA_DIR")
	if dataDir == "" {
		discovered, err := platform.DataDir()
		if err != nil {
			return app.Service{}, err
		}
		dataDir = discovered
	}

	return app.Service{
		DataDir:  dataDir,
		RepoDir:  filepath.Join(dataDir, "repo"),
		Manifest: state.Store{Path: filepath.Join(dataDir, "manifest.json")},
		Source:   source.Git{},
		Merger:   config.Merger{},
		Links:    link.Manager{BackupRoot: filepath.Join(dataDir, "backups")},
		OpenCode: opencode.CLI{},
	}, nil
}

func newRoot(service app.Service) *cobra.Command {
	var root *cobra.Command
	root = &cobra.Command{
		Use:   "rutter",
		Short: "Install and sync the team OpenCode configuration",
		Long: `rutter installs the shared OpenCode config, keeps it current, and removes it cleanly.

It never prompts. Pass --json or set RUTTER_JSON=1 for machine-readable output,
and run "rutter agent" for the complete contract used by scripts and agents.`,
		Example: `  rutter install
  rutter status --json
  rutter sync --no-reload
  rutter agent --json`,
		Version:       version.Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if jsonEnabled(cmd) {
				return writeJSON(cmd.OutOrStdout(), buildAgentSpec(root))
			}
			return cmd.Help()
		},
	}
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)
	root.PersistentFlags().Bool("json", false, "print machine-readable JSON (env: RUTTER_JSON=1)")
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return usageError{err: err}
	})
	root.CompletionOptions.HiddenDefaultCmd = true

	root.AddCommand(
		annotate(newInstall(service), "InstallResult"),
		annotate(newSync(service), "InstallResult"),
		annotate(newUninstall(service), "InstallResult"),
		annotate(newStatus(service), "Status"),
		annotate(newDoctor(service), "DoctorReport"),
		newConfig(service),
		annotate(newVersion(), "Version"),
		newAgent(root),
	)
	return root
}

func newVersion() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Print version information",
		Example: "  rutter version\n  rutter version --json",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if jsonEnabled(cmd) {
				return writeJSON(cmd.OutOrStdout(), version.BuildInfo())
			}
			line(cmd.OutOrStdout(), "%s\n", version.String())
			return nil
		},
	}
}

func envJSON() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("RUTTER_JSON"))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func argsContainJSON(args []string) bool {
	for _, arg := range args {
		if arg == "--json" || strings.HasPrefix(arg, "--json=") {
			return true
		}
	}
	return false
}
