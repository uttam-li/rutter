package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/uttam-li/rutter/internal/app"
	"github.com/uttam-li/rutter/internal/link"
	"github.com/uttam-li/rutter/internal/state"
)

type statusOutput struct {
	app.Status
	Healthy bool `json:"healthy"`
}

func newStatus(service app.Service) *cobra.Command {
	var check bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the installed config and link state",
		Long: `Show the manifest, the generated config, and whether managed links still
match their sources.

Use --check to exit with a non-zero code when rutter is not installed or any
managed entry needs attention.`,
		Example: "  rutter status\n  rutter status --json\n  rutter status --check",
		RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := service.Status(cmd.Context())
			if err != nil {
				return err
			}

			healthy := statusHealthy(status)
			if jsonEnabled(cmd) {
				if err := writeJSON(cmd.OutOrStdout(), statusOutput{Status: status, Healthy: healthy}); err != nil {
					return err
				}
			} else {
				printStatus(cmd.OutOrStdout(), status)
			}

			if !check || healthy {
				return nil
			}
			return unhealthyError(status)
		},
	}

	cmd.Flags().BoolVar(&check, "check", false, "exit non-zero when not installed or links have drifted")
	return cmd
}

func statusHealthy(status app.Status) bool {
	if !status.Installed || !status.ConfigGenerated || len(status.Links) == 0 {
		return false
	}
	for _, item := range status.Links {
		if item.State != "ok" {
			return false
		}
	}
	return true
}

func unhealthyError(status app.Status) error {
	if !status.Installed {
		return state.ErrNotInstalled
	}
	for _, item := range status.Links {
		if item.State != "ok" {
			return fmt.Errorf("%w: %s (%s)", link.ErrDrift, item.Target, item.State)
		}
	}
	return fmt.Errorf("%w: %s is not managed by rutter", link.ErrDrift, status.ConfigPath)
}
