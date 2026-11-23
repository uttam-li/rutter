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
		Long: `Remove the entries rutter installed, restore backed-up files, delete the
generated config, and remove the manifest.

Files changed outside rutter are left in place and reported as drift.`,
		Example: "  rutter uninstall\n  rutter uninstall --purge\n  rutter uninstall --json",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := service.Uninstall(cmd.Context(), purge)
			if err != nil {
				return err
			}
			if jsonEnabled(cmd) {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			printUninstallResult(cmd.OutOrStdout(), result)
			return nil
		},
	}

	cmd.Flags().BoolVar(&purge, "purge", false, "also remove the repository checkout, backups, and state")
	return cmd
}
