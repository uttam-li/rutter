package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/uttam-li/rutter/internal/app"
	"github.com/uttam-li/rutter/internal/link"
)

// line writes human output. Callers deliberately ignore write errors: the
// terminal is not a durability boundary, and stdout failures surface as a
// broken pipe from the process anyway.
func line(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

// jsonEnabled reports whether the command should emit machine-readable output.
func jsonEnabled(cmd *cobra.Command) bool {
	enabled, err := cmd.Flags().GetBool("json")
	return jsonDefault || (err == nil && enabled)
}

// writeJSON writes the success envelope. Every command uses the same shape so
// an agent can branch on "ok" before reading "data".
func writeJSON(w io.Writer, data any) error {
	payload := struct {
		OK   bool `json:"ok"`
		Data any  `json:"data"`
	}{OK: true, Data: data}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(payload)
}

func printInstallResult(w io.Writer, result app.Result) {
	printResult(w, "rutter is installed.", result, "LINKED")
}

func printSyncResult(w io.Writer, result app.Result) {
	printResult(w, "rutter sync complete.", result, "LINKED")
}

func printUninstallResult(w io.Writer, result app.Result) {
	line(w, "rutter has been uninstalled.\n\n")
	line(w, "Config dir: %s\n", result.ConfigDir)
	line(w, "Repository: %s\n", result.RepoDir)
	printLinks(w, "REMOVED", result.Links)
	printReload(w, result.Reloaded)
	printWarnings(w, result.Warnings)
}

func printResult(w io.Writer, heading string, result app.Result, linkHeading string) {
	line(w, "%s\n\n", heading)
	line(w, "Config dir: %s\n", result.ConfigDir)
	line(w, "Repository: %s\n", result.RepoDir)
	line(w, "Config:     %s\n", result.ConfigPath)
	line(w, "Personal:   %s\n", result.PersonalPath)
	if result.PersonalCreated {
		line(w, "A personal config was created from the example. Edit it for this machine.\n")
	}
	printLinks(w, linkHeading, result.Links)
	printReload(w, result.Reloaded)
	printWarnings(w, result.Warnings)
}

func printLinks(w io.Writer, heading string, links []link.Link) {
	if len(links) == 0 {
		return
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	line(tw, "\n%s\tMODE\n", heading)
	for _, item := range links {
		line(tw, "%s\t%s\n", filepath.Base(item.Target), item.Mode)
	}
	_ = tw.Flush()
}

func printReload(w io.Writer, reloaded bool) {
	if reloaded {
		line(w, "\nOpenCode reloaded.\n")
	}
}

func printWarnings(w io.Writer, warnings []string) {
	for _, warning := range warnings {
		line(w, "Warning: %s\n", warning)
	}
}

func printStatus(w io.Writer, status app.Status) {
	if !status.Installed {
		line(w, "rutter is not installed.\n")
		line(w, "Run `rutter install` to set it up.\n")
		return
	}

	line(w, "rutter is installed.\n\n")
	line(w, "Config dir: %s\n", status.ConfigDir)
	line(w, "Repository: %s\n", status.RepoURL)
	line(w, "Checkout:   %s\n", status.RepoDir)
	line(w, "Config:     %s\n", status.ConfigPath)
	line(w, "Personal:   %s\n", status.PersonalPath)
	if !status.ConfigGenerated {
		line(w, "Config is not the generated rutter file.\n")
	}

	if len(status.Links) == 0 {
		return
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	line(tw, "\nLINK\tSTATE\tDETAIL\n")
	needsSync := false
	for _, check := range status.Links {
		line(tw, "%s\t%s\t%s\n", filepath.Base(check.Target), check.State, check.Detail)
		if check.State != "ok" {
			needsSync = true
		}
	}
	_ = tw.Flush()

	if needsSync {
		line(w, "\nSome entries need attention. Run `rutter sync` to repair them.\n")
	}
}

func printDoctor(w io.Writer, report app.DoctorReport) {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	line(tw, "CHECK\tSTATUS\tDETAIL\n")
	warnings := 0
	for _, check := range report.Checks {
		line(tw, "%s\t%s\t%s\n", check.Name, check.Status, check.Detail)
		if check.Status == "warn" {
			warnings++
		}
	}
	_ = tw.Flush()

	switch {
	case !report.OK:
		line(w, "\nSome checks failed. Fix the errors above and run `rutter doctor` again.\n")
	case warnings > 0:
		line(w, "\nNo blocking problems found, but %d warning(s) need attention.\n", warnings)
	default:
		line(w, "\nAll checks passed.\n")
	}
}
