# rutter

<p align="center">
  <img src="assets/logo.svg" alt="rutter logo" width="120" height="120">
</p>

`rutter` installs and syncs a team's shared OpenCode configuration on macOS, Windows, and Linux.

It keeps the CLI and the config in one repository:

- the Go CLI at the root
- the shared config under [`config/`](config/)

## What it does

- Clones the config repository into a local data directory.
- Merges `config/opencode.base.json` with a machine-local `opencode.personal.json`.
- Writes the generated config to OpenCode's config directory.
- Links `AGENTS.md`, `agents/`, `commands/`, and `skills/` into place.
- Backs up existing files instead of deleting them.
- Records a manifest so `rutter uninstall` can restore the machine.
- Asks OpenCode to reload after install and sync.
- Never prompts, so it works the same for people, scripts, and AI agents.

## Requirements

- OpenCode v2.3 or newer
- git
- Go 1.27.4 or newer when building from source

## Install

Download a binary from GitHub Releases, or build from source:

```bash
go install github.com/uttam-li/rutter/cmd/rutter@latest
```

Then:

```bash
rutter install
```

The default repository is `https://github.com/uttam-li/rutter.git`. Override it with `--repo`:

```bash
rutter install --repo git@github.com:acme/rutter.git
```

`rutter` asks `opencode debug paths` for the config directory. If OpenCode is not on `PATH`, pass it directly:

```bash
rutter install --config-dir "$HOME/.config/opencode"
```

## Update

```bash
rutter sync
```

This pulls the config repository, re-merges the config, refreshes the managed entries, and reloads OpenCode.

## For agents and scripts

`rutter` is built for unattended use: it never prompts, `install` and `sync` are idempotent, and every command has a stable JSON contract.

Start with the self-description:

```bash
rutter agent --json
```

It returns every command, flag, JSON object, exit code, environment variable, and the one interactive command to avoid. `rutter agent` without `--json` prints the same contract as a readable guide. For a one-line discovery call, `rutter --json` returns the same document.

Enable JSON once with `RUTTER_JSON=1`, or pass `--json` to any command:

```bash
RUTTER_JSON=1 rutter status
```

Output contract:

- Success: `{"ok":true,"data":...}` on stdout, exit `0`
- Failure: `{"ok":false,"error":{"code","exit_code","message","hint"}}` on stderr
- Human errors: `rutter: ...` plus a `hint:` line on stderr

A typical agent loop:

```bash
rutter status --json          # inspect without changing anything
rutter install --json         # first run
rutter sync --json            # receive the owner's config updates
rutter doctor --json          # diagnose the environment
rutter status --check         # non-zero when not installed or drifted
```

Exit codes:

| Code | Name | Meaning |
| --- | --- | --- |
| 0 | `ok` | success |
| 1 | `error` | unexpected failure |
| 2 | `usage` | unknown command, flag, or argument |
| 3 | `not_installed` | run `rutter install` |
| 4 | `environment` | git, OpenCode, or a doctor check failed |
| 5 | `drift` | a managed path changed outside rutter |
| 6 | `repository` | git clone, pull, or network failure |
| 7 | `config` | invalid base or personal config |

Environment variables:

| Variable | Purpose |
| --- | --- |
| `RUTTER_JSON` | `1` or `true` sets JSON output by default |
| `RUTTER_REPO` | repository URL for `install` |
| `RUTTER_CONFIG_DIR` | OpenCode config directory for `install` |
| `RUTTER_DATA_DIR` | data directory for the checkout, backups, and manifest |

Guarantees agents can rely on:

- No prompts and no TTY requirement.
- `install` and `sync` can be retried safely.
- Unmanaged files are moved to a backup, never deleted.
- `config edit` is the only interactive command, and `rutter agent` marks it.
- Nothing is committed or pushed. Config updates flow one way from the repository owner.

See [`docs/agent.md`](docs/agent.md) for worked examples and failure handling.

## Inspect

```bash
rutter agent
rutter status
rutter status --json
rutter status --check
rutter doctor
rutter config path
rutter config show
rutter config show --json
```

`rutter config edit` opens the machine-local personal config in `$VISUAL` or `$EDITOR`.

## Uninstall

```bash
rutter uninstall
```

This removes managed entries, restores backups, and deletes the generated config. Add `--purge` to remove the repository checkout, backups, and state:

```bash
rutter uninstall --purge
```

## How the config merge works

The team base and the personal file merge with these rules:

- Objects merge deeply.
- Scalars in the personal file win.
- Arrays replace, except `permissions`, `instructions`, `skills`, and `plugins`, which append.

The generated `opencode.jsonc` starts with a header. Treat it as machine-owned and edit `config/opencode.base.json` or `opencode.personal.json` instead.

## Platform behavior

On macOS and Linux, `rutter` creates symlinks. Changes in the config repository appear immediately.

On Windows, `rutter` copies files and directories. Run `rutter sync` after pulling changes. This avoids symlink privileges and Developer Mode requirements.

Treat linked entries as read-only. On Windows, `rutter sync` refreshes copies and does not merge edits made inside them.

Backups live under the rutter data directory:

- macOS: `~/Library/Application Support/rutter`
- Linux: `$XDG_DATA_HOME/rutter` or `~/.local/share/rutter`
- Windows: `%LOCALAPPDATA%\rutter`

## Repository layout

```text
rutter/
├── cmd/rutter/          # main package
├── internal/
│   ├── app/             # install, sync, uninstall, status, doctor
│   ├── cli/             # Cobra commands
│   ├── config/          # JSONC merge
│   ├── link/            # links, copies, backups, pruning
│   ├── opencode/        # config discovery and reload
│   ├── platform/        # data dir and atomic writes
│   ├── source/          # git clone and pull
│   ├── state/           # manifest
│   └── version/         # build metadata
├── config/              # the shared OpenCode config
├── docs/design.md       # module design
├── docs/agent.md        # agent and script contract
├── AGENTS.md            # instructions for agents working on this repo
├── .github/workflows/   # CI and release
└── .goreleaser.yaml
```

## Development

```bash
make build
make test
make lint
```

Run the test suite in the official Go image:

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.27.4 go test ./...
```

Cross-compile:

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build ./cmd/rutter
GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build ./cmd/rutter
```

Create a local release snapshot:

```bash
goreleaser release --snapshot --clean
```

## Release

Push a tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow builds macOS, Linux, and Windows binaries for amd64 and arm64, then publishes archives and checksums.

## Security

The shared config controls permissions, MCP servers, and agent behavior. Review changes to `config/` like code. Do not commit API keys, OAuth tokens, or machine-specific paths.

## License

MIT
