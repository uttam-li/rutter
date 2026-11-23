package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/uttam-li/rutter/internal/app"
)

func newDoctor(service app.Service) *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check OpenCode, git, config, and links",
		RunE: func(cmd *cobra.Command, _ []string) error {
			report := service.Doctor(cmd.Context())
			if jsonOutput {
				if err := writeJSON(cmd.OutOrStdout(), report); err != nil {
					return err
				}
			} else {
				printDoctor(cmd.OutOrStdout(), report)
			}
			if !report.OK {
				return errors.New("doctor found problems")
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "print JSON")
	return cmd
}
