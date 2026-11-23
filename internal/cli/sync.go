package cli

import (
	"github.com/spf13/cobra"

	"github.com/uttam-li/rutter/internal/app"
)

func newSync(service app.Service) *cobra.Command {
	var noReload bool

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Pull the config repo and reapply it",
		Long: `Pull the config repository with a fast-forward-only merge, re-merge the
config, refresh links, and reload OpenCode.

Local changes in the config checkout are never discarded.`,
		Example: "  rutter sync\n  rutter sync --no-reload\n  rutter sync --json",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := service.Sync(cmd.Context(), !noReload)
			if err != nil {
				return err
			}
			if jsonEnabled(cmd) {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			printSyncResult(cmd.OutOrStdout(), result)
			return nil
		},
	}

	cmd.Flags().BoolVar(&noReload, "no-reload", false, "do not ask OpenCode to reload")
	return cmd
}
