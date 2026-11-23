package cli

import (
	"github.com/spf13/cobra"

	"github.com/uttam-li/rutter/internal/app"
)

func newStatus(service app.Service) *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the installed config and link state",
		RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := service.Status(cmd.Context())
			if err != nil {
				return err
			}
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), status)
			}
			printStatus(cmd.OutOrStdout(), status)
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "print JSON")
	return cmd
}
