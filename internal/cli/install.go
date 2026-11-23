package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/uttam-li/rutter/internal/app"
	"github.com/uttam-li/rutter/internal/version"
)

func newInstall(service app.Service) *cobra.Command {
	var (
		repoURL   string
		configDir string
		noReload  bool
	)

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the shared config on this machine",
		Long: `Clone the config repository, write the merged OpenCode config, link the
shared entries, and record what was installed.

Existing files are moved to a backup and restored by "rutter uninstall".
Running install again is safe.`,
		Example: "  rutter install\n  rutter install --json\n  rutter install --repo git@github.com:acme/rutter.git",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !cmd.Flags().Changed("repo") {
				if value := os.Getenv("RUTTER_REPO"); value != "" {
					repoURL = value
				}
			}
			if !cmd.Flags().Changed("config-dir") {
				if value := os.Getenv("RUTTER_CONFIG_DIR"); value != "" {
					configDir = value
				}
			}

			result, err := service.Install(cmd.Context(), app.InstallOptions{
				RepoURL:   repoURL,
				ConfigDir: configDir,
				Reload:    !noReload,
			})
			if err != nil {
				return err
			}
			if jsonEnabled(cmd) {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			printInstallResult(cmd.OutOrStdout(), result)
			return nil
		},
	}

	cmd.Flags().StringVar(&repoURL, "repo", version.DefaultRepo, "config repository URL (env: RUTTER_REPO)")
	cmd.Flags().StringVar(&configDir, "config-dir", "", "OpenCode config directory (default: detected from opencode; env: RUTTER_CONFIG_DIR)")
	cmd.Flags().BoolVar(&noReload, "no-reload", false, "do not ask OpenCode to reload")
	return cmd
}
