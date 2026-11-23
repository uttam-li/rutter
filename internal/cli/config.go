package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/uttam-li/rutter/internal/app"
	"github.com/uttam-li/rutter/internal/config"
	"github.com/uttam-li/rutter/internal/state"
)

type configPaths struct {
	ConfigDir string `json:"config_dir"`
	Config    string `json:"config"`
	Personal  string `json:"personal"`
	Repo      string `json:"repo"`
}

func newConfig(service app.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Inspect or edit config files",
	}
	cmd.AddCommand(
		newConfigPath(service),
		newConfigShow(service),
		newConfigEdit(service),
	)
	return cmd
}

func newConfigPath(service app.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "path",
		Short:   "Print config paths",
		Example: "  rutter config path\n  rutter config path --json",
		RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := installedStatus(cmd.Context(), service)
			if err != nil {
				return err
			}

			paths := configPaths{
				ConfigDir: status.ConfigDir,
				Config:    status.ConfigPath,
				Personal:  status.PersonalPath,
				Repo:      status.RepoDir,
			}
			if jsonEnabled(cmd) {
				return writeJSON(cmd.OutOrStdout(), paths)
			}
			line(cmd.OutOrStdout(), "config_dir\t%s\n", paths.ConfigDir)
			line(cmd.OutOrStdout(), "config\t%s\n", paths.Config)
			line(cmd.OutOrStdout(), "personal\t%s\n", paths.Personal)
			line(cmd.OutOrStdout(), "repo\t%s\n", paths.Repo)
			return nil
		},
	}
	return annotate(cmd, "ConfigPaths")
}

func newConfigShow(service app.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Short:   "Print the generated OpenCode config",
		Long:    "Print the generated opencode.jsonc. With --json, print the parsed document with comments removed.",
		Example: "  rutter config show\n  rutter config show --json",
		RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := installedStatus(cmd.Context(), service)
			if err != nil {
				return err
			}

			if jsonEnabled(cmd) {
				document, err := config.ReadDocument(status.ConfigPath)
				if err != nil {
					return err
				}
				return writeJSON(cmd.OutOrStdout(), document)
			}

			data, err := os.ReadFile(status.ConfigPath)
			if err != nil {
				return fmt.Errorf("read config: %w", err)
			}
			_, err = cmd.OutOrStdout().Write(data)
			return err
		},
	}
	return annotate(cmd, "ConfigDocument")
}

func newConfigEdit(service app.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit the machine-local personal config",
		Long: `Open the machine-local opencode.personal.json in $VISUAL or $EDITOR.

This command is interactive and is not meant for unattended agents.`,
		Example: "  rutter config edit",
		RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := installedStatus(cmd.Context(), service)
			if err != nil {
				return err
			}
			return openEditor(status.PersonalPath)
		},
	}
	return markInteractive(cmd)
}

func installedStatus(ctx context.Context, service app.Service) (app.Status, error) {
	status, err := service.Status(ctx)
	if err != nil {
		return app.Status{}, err
	}
	if !status.Installed {
		return app.Status{}, state.ErrNotInstalled
	}
	return status, nil
}

func openEditor(path string) error {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "vi"
		}
	}

	fields := strings.Fields(editor)
	if len(fields) == 0 {
		return errors.New("editor command is empty")
	}
	cmd := exec.Command(fields[0], append(fields[1:], path)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
