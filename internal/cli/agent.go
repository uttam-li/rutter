package cli

import (
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/uttam-li/rutter/internal/version"
)

// agentSpec is the self-description returned by `rutter agent --json`. It is
// generated from the live command tree, so flags and summaries cannot drift.
type agentSpec struct {
	Name           string            `json:"name"`
	Version        string            `json:"version"`
	Summary        string            `json:"summary"`
	NonInteractive bool              `json:"non_interactive"`
	Invocation     string            `json:"invocation"`
	Output         string            `json:"output"`
	Environment    []agentEnvVar     `json:"environment"`
	ExitCodes      []agentExitCode   `json:"exit_codes"`
	Commands       []agentCommand    `json:"commands"`
	Examples       []string          `json:"examples"`
	JSONObjects    map[string]string `json:"json_objects"`
	Notes          []string          `json:"notes"`
}

type agentEnvVar struct {
	Name    string `json:"name"`
	Purpose string `json:"purpose"`
	Default string `json:"default,omitempty"`
}

type agentExitCode struct {
	Code    int    `json:"code"`
	Name    string `json:"name"`
	Meaning string `json:"meaning"`
}

type agentCommand struct {
	Path        string      `json:"path"`
	Summary     string      `json:"summary"`
	Usage       string      `json:"usage"`
	Interactive bool        `json:"interactive,omitempty"`
	JSONObject  string      `json:"json_object,omitempty"`
	Flags       []agentFlag `json:"flags,omitempty"`
}

type agentFlag struct {
	Name      string `json:"name"`
	Shorthand string `json:"shorthand,omitempty"`
	Type      string `json:"type"`
	Default   string `json:"default,omitempty"`
	Usage     string `json:"usage"`
}

var jsonObjectDocs = map[string]string{
	"InstallResult":  "config_dir, repo_dir, config_path, personal_path, personal_created, links[{source,target,kind,mode,backup}], reloaded, warnings[]",
	"Status":         "installed, repo_url, repo_dir, config_dir, config_path, personal_path, config_generated, links[{target,source,state,detail}], healthy",
	"DoctorReport":   "ok, checks[{name,status,detail}]; status is ok, warn, or error",
	"ConfigPaths":    "config_dir, config, personal, repo",
	"ConfigDocument": "the parsed opencode.jsonc document with comments removed",
	"Version":        "name, version, commit, date, repo",
	"AgentSpec":      "this document",
}

func newAgent(root *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Print the non-interactive contract for agents and scripts",
		Long: `Print everything an agent needs to call rutter safely: commands, flags,
JSON objects, exit codes, and environment variables.

Add --json for the machine-readable version.`,
		Example: "  rutter agent\n  rutter agent --json",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			spec := buildAgentSpec(root)
			if jsonEnabled(cmd) {
				return writeJSON(cmd.OutOrStdout(), spec)
			}
			writeAgentGuide(cmd.OutOrStdout(), spec)
			return nil
		},
	}
	return annotate(cmd, "AgentSpec")
}

func buildAgentSpec(root *cobra.Command) agentSpec {
	return agentSpec{
		Name:           "rutter",
		Version:        version.Version,
		Summary:        "Installs, syncs, and uninstalls a shared OpenCode configuration on macOS, Windows, and Linux.",
		NonInteractive: true,
		Invocation:     "rutter <command> [flags]",
		Output:         `Success: {"ok":true,"data":...} on stdout. Failure: {"ok":false,"error":{code,exit_code,message,hint}} on stderr.`,
		Environment: []agentEnvVar{
			{Name: "RUTTER_JSON", Purpose: "Set to 1 or true to default every command to JSON output"},
			{Name: "RUTTER_REPO", Purpose: "Config repository URL used by install when --repo is not passed", Default: version.DefaultRepo},
			{Name: "RUTTER_CONFIG_DIR", Purpose: "OpenCode config directory used by install when --config-dir is not passed", Default: "detected from `opencode debug paths`"},
			{Name: "RUTTER_DATA_DIR", Purpose: "Override the rutter data directory for checkouts, backups, and the manifest", Default: "platform data directory"},
			{Name: "EDITOR / VISUAL", Purpose: "Editor opened by `rutter config edit`; agents should not call it"},
		},
		ExitCodes: []agentExitCode{
			{Code: int(exitOK), Name: "ok", Meaning: "command succeeded"},
			{Code: int(exitFailure), Name: "error", Meaning: "unexpected failure; read the message"},
			{Code: int(exitUsage), Name: "usage", Meaning: "unknown command, flag, or argument"},
			{Code: int(exitNotInstalled), Name: "not_installed", Meaning: "rutter is not installed on this machine"},
			{Code: int(exitEnvironment), Name: "environment", Meaning: "git, OpenCode, or a health check is failing"},
			{Code: int(exitDrift), Name: "drift", Meaning: "a managed path was changed outside rutter"},
			{Code: int(exitRepository), Name: "repository", Meaning: "git clone, pull, or network failure"},
			{Code: int(exitConfig), Name: "config", Meaning: "the shared or personal config is not valid JSON"},
		},
		Commands:    collectCommands(root),
		Examples:    []string{"rutter install --json", "rutter status --json", "rutter status --check", "rutter sync --json --no-reload", "rutter doctor --json", "rutter config show --json", "rutter uninstall --json"},
		JSONObjects: jsonObjectDocs,
		Notes: []string{
			"rutter never prompts, so it is safe to run unattended.",
			"The repository owner publishes config changes; teammates run `rutter sync` to receive them.",
			"Existing files that rutter does not own are moved to a backup, never deleted.",
			"Install and sync are idempotent; running them again is safe.",
			"On Windows, config entries are copies, so `rutter sync` is required after changes.",
		},
	}
}

func collectCommands(root *cobra.Command) []agentCommand {
	var commands []agentCommand

	var walk func(parent *cobra.Command, prefix string)
	walk = func(parent *cobra.Command, prefix string) {
		for _, child := range parent.Commands() {
			if child.Hidden || child.Name() == "help" || child.Name() == "completion" {
				continue
			}

			path := strings.TrimSpace(prefix + " " + child.Name())
			var flags []agentFlag
			if child.Run != nil || child.RunE != nil {
				flags = collectFlags(child)
			}
			commands = append(commands, agentCommand{
				Path:        path,
				Summary:     child.Short,
				Usage:       child.UseLine(),
				Interactive: child.Annotations[annotationInteractive] == "true",
				JSONObject:  child.Annotations[annotationJSON],
				Flags:       flags,
			})
			walk(child, path)
		}
	}
	walk(root, "")

	sort.Slice(commands, func(i, j int) bool { return commands[i].Path < commands[j].Path })
	return commands
}

func collectFlags(cmd *cobra.Command) []agentFlag {
	seen := map[string]bool{}
	var flags []agentFlag

	add := func(flagSet *pflag.FlagSet) {
		flagSet.VisitAll(func(flag *pflag.Flag) {
			if flag.Name == "help" || seen[flag.Name] {
				return
			}
			seen[flag.Name] = true
			flags = append(flags, agentFlag{
				Name:      flag.Name,
				Shorthand: flag.Shorthand,
				Type:      flag.Value.Type(),
				Default:   flag.DefValue,
				Usage:     flag.Usage,
			})
		})
	}

	add(cmd.LocalFlags())
	add(cmd.InheritedFlags())

	sort.Slice(flags, func(i, j int) bool { return flags[i].Name < flags[j].Name })
	return flags
}

func writeAgentGuide(w io.Writer, spec agentSpec) {
	line(w, "%s %s - guide for agents and scripts\n\n", spec.Name, spec.Version)
	line(w, "%s\n\n", spec.Summary)
	line(w, "rutter never prompts. Set RUTTER_JSON=1 or pass --json for machine-readable output.\n\n")

	line(w, "COMMANDS\n")
	for _, command := range spec.Commands {
		line(w, "  rutter %s\n", command.Path)
		line(w, "      %s\n", command.Summary)
		if command.Interactive {
			line(w, "      interactive: do not call unattended\n")
		}
		if command.JSONObject != "" {
			line(w, "      json: %s\n", command.JSONObject)
		}
		for _, flag := range command.Flags {
			line(w, "      --%s %s\n", flag.Name, flag.Usage)
		}
		line(w, "\n")
	}

	line(w, "EXIT CODES\n")
	for _, code := range spec.ExitCodes {
		line(w, "  %d  %-13s %s\n", code.Code, code.Name, code.Meaning)
	}

	line(w, "\nOUTPUT\n")
	line(w, "  success: {\"ok\":true,\"data\":...} on stdout\n")
	line(w, "  failure: {\"ok\":false,\"error\":{...}} on stderr\n")

	line(w, "\nENVIRONMENT\n")
	for _, variable := range spec.Environment {
		line(w, "  %s - %s\n", variable.Name, variable.Purpose)
	}

	line(w, "\nEXAMPLES\n")
	for _, example := range spec.Examples {
		line(w, "  %s\n", example)
	}

	line(w, "\nJSON OBJECTS\n")
	names := make([]string, 0, len(spec.JSONObjects))
	for name := range spec.JSONObjects {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		line(w, "  %-14s %s\n", name, spec.JSONObjects[name])
	}

	line(w, "\nNOTES\n")
	for _, note := range spec.Notes {
		line(w, "  - %s\n", note)
	}
}
