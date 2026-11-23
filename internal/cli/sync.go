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
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := service.Sync(cmd.Context(), !noReload)
			if err != nil {
				return err
			}
			printResult(cmd.OutOrStdout(), result)
			return nil
		},
	}

	cmd.Flags().BoolVar(&noReload, "no-reload", false, "do not ask OpenCode to reload")
	return cmd
}
