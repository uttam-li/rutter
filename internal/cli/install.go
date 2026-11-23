package cli

import (
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := service.Install(cmd.Context(), app.InstallOptions{
				RepoURL:   repoURL,
				ConfigDir: configDir,
				Reload:    !noReload,
			})
			if err != nil {
				return err
			}
			printResult(cmd.OutOrStdout(), result)
			return nil
		},
	}

	cmd.Flags().StringVar(&repoURL, "repo", version.DefaultRepo, "config repository URL")
	cmd.Flags().StringVar(&configDir, "config-dir", "", "OpenCode config directory (default: detected from opencode)")
	cmd.Flags().BoolVar(&noReload, "no-reload", false, "do not ask OpenCode to reload")
	return cmd
}
