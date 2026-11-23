package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"text/tabwriter"

	"github.com/uttam-li/rutter/internal/app"
)

func writeJSON(w io.Writer, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}

func printResult(w io.Writer, result app.Result) {
	_, _ = fmt.Fprintln(w, "rutter is installed.")
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintf(w, "Config dir: %s\n", result.ConfigDir)
	_, _ = fmt.Fprintf(w, "Repository: %s\n", result.RepoDir)
	_, _ = fmt.Fprintf(w, "Config:     %s\n", result.ConfigPath)
	_, _ = fmt.Fprintf(w, "Personal:   %s\n", result.PersonalPath)
	if result.PersonalCreated {
		_, _ = fmt.Fprintln(w, "A personal config was created from the example. Edit it for this machine.")
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "\nLINKED\tMODE")
	for _, item := range result.Links {
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", filepath.Base(item.Target), item.Mode)
	}
	_ = tw.Flush()

	if result.Reloaded {
		_, _ = fmt.Fprintln(w, "\nOpenCode reloaded.")
	}
	for _, warning := range result.Warnings {
		_, _ = fmt.Fprintf(w, "Warning: %s\n", warning)
	}
}

func printStatus(w io.Writer, status app.Status) {
	if !status.Installed {
		_, _ = fmt.Fprintln(w, "rutter is not installed.")
		return
	}

	_, _ = fmt.Fprintln(w, "rutter is installed.")
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintf(w, "Config dir: %s\n", status.ConfigDir)
	_, _ = fmt.Fprintf(w, "Repository: %s\n", status.RepoURL)
	_, _ = fmt.Fprintf(w, "Checkout:   %s\n", status.RepoDir)
	_, _ = fmt.Fprintf(w, "Config:     %s\n", status.ConfigPath)
	_, _ = fmt.Fprintf(w, "Personal:   %s\n", status.PersonalPath)
	if !status.ConfigGenerated {
		_, _ = fmt.Fprintln(w, "Config is not the generated rutter file.")
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "\nLINK\tSTATE\tDETAIL")
	for _, check := range status.Links {
		_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\n", filepath.Base(check.Target), check.State, check.Detail)
	}
	_ = tw.Flush()
}

func printDoctor(w io.Writer, report app.DoctorReport) {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "CHECK\tSTATUS\tDETAIL")
	for _, check := range report.Checks {
		_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\n", check.Name, check.Status, check.Detail)
	}
	_ = tw.Flush()
	if report.OK {
		_, _ = fmt.Fprintln(w, "\nAll checks passed.")
	} else {
		_, _ = fmt.Fprintln(w, "\nSome checks failed.")
	}
}
