package cli

import (
	"github.com/spf13/cobra"

	"github.com/uttam-li/rutter/internal/app"
)

func newUninstall(service app.Service) *cobra.Command {
	var purge bool

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove managed links and restore backups",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := service.Uninstall(cmd.Context(), purge)
			if err != nil {
				return err
			}
			printResult(cmd.OutOrStdout(), result)
			return nil
		},
	}

	cmd.Flags().BoolVar(&purge, "purge", false, "also remove the repository checkout, backups, and state")
	return cmd
}
