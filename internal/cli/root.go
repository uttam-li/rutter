package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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

// Execute builds the service and runs the CLI.
func Execute() error {
	service, err := buildService()
	if err != nil {
		return err
	}
	return newRoot(service).ExecuteContext(context.Background())
}

func buildService() (app.Service, error) {
	dataDir, err := platform.DataDir()
	if err != nil {
		return app.Service{}, err
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
	root := &cobra.Command{
		Use:           "rutter",
		Short:         "Install and sync the team OpenCode configuration",
		Long:          "rutter installs the shared OpenCode config, keeps it current, and removes it cleanly.",
		Version:       version.Version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)

	root.AddCommand(
		newInstall(service),
		newSync(service),
		newUninstall(service),
		newStatus(service),
		newDoctor(service),
		newConfig(service),
		newVersion(),
	)
	return root
}

func newVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), version.String())
			return nil
		},
	}
}
