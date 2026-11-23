package cli

import "github.com/spf13/cobra"

const (
	annotationJSON        = "rutter.io/json"
	annotationInteractive = "rutter.io/interactive"
)

// annotate records the JSON object a command returns. `rutter agent` reads it
// so the documented contract cannot drift from the command tree.
func annotate(cmd *cobra.Command, object string) *cobra.Command {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[annotationJSON] = object
	return cmd
}

// markInteractive flags a command an agent must not call unattended.
func markInteractive(cmd *cobra.Command) *cobra.Command {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[annotationInteractive] = "true"
	return cmd
}
