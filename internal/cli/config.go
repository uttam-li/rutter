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
)

func newConfig(service app.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Inspect or edit config files",
	}
	cmd.AddCommand(
		newConfigPath(service),
		newConfigEdit(service),
		newConfigShow(service),
	)
	return cmd
}

func newConfigPath(service app.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print config paths",
		RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := installedStatus(cmd.Context(), service)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "config_dir\t%s\n", status.ConfigDir)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "config\t%s\n", status.ConfigPath)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "personal\t%s\n", status.PersonalPath)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "repo\t%s\n", status.RepoDir)
			return nil
		},
	}
}

func newConfigEdit(service app.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Edit the machine-local personal config",
		RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := installedStatus(cmd.Context(), service)
			if err != nil {
				return err
			}
			return openEditor(status.PersonalPath)
		},
	}
}

func newConfigShow(service app.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Print the generated OpenCode config",
		RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := installedStatus(cmd.Context(), service)
			if err != nil {
				return err
			}
			data, err := os.ReadFile(status.ConfigPath)
			if err != nil {
				return fmt.Errorf("read config: %w", err)
			}
			_, err = cmd.OutOrStdout().Write(data)
			return err
		},
	}
}

func installedStatus(ctx context.Context, service app.Service) (app.Status, error) {
	status, err := service.Status(ctx)
	if err != nil {
		return app.Status{}, err
	}
	if !status.Installed {
		return app.Status{}, errors.New("rutter is not installed")
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
