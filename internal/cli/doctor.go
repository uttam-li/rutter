package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/uttam-li/rutter/internal/app"
)

func newDoctor(service app.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check OpenCode, git, config, and links",
		Long: `Run local diagnostics: git and OpenCode availability, the manifest, the
config merge, and every managed link.

Exits non-zero only when a check reports an error. Warnings do not fail.`,
		Example: "  rutter doctor\n  rutter doctor --json",
		RunE: func(cmd *cobra.Command, _ []string) error {
			report := service.Doctor(cmd.Context())
			if jsonEnabled(cmd) {
				if err := writeJSON(cmd.OutOrStdout(), report); err != nil {
					return err
				}
			} else {
				printDoctor(cmd.OutOrStdout(), report)
			}
			if !report.OK {
				return fmt.Errorf("%w: one or more checks failed", errEnvironment)
			}
			return nil
		},
	}

	return cmd
}
