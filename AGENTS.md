# Working in this repository

`rutter` is a Go CLI that installs a shared OpenCode configuration. The Go code lives at the root; the shared config lives under `config/`.

## Common commands

```bash
make build
make test
make lint
go vet ./...
```

Run the tests in the official image when changing platform or link behavior:

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.27.4 go test ./...
```

Smoke-test without touching real state:

```bash
RUTTER_DATA_DIR=/tmp/rutter-data go run ./cmd/rutter status --json
```

## Layout

- `cmd/rutter` — entry point
- `internal/app` — orchestration; the deep module
- `internal/cli` — commands, JSON envelope, exit codes, self-description
- `internal/config` — JSONC merge and document parsing
- `internal/link` — symlinks, Windows copies, backups, drift
- `internal/opencode` — config discovery and reload
- `internal/source` — git clone and pull
- `internal/state` — manifest
- `config/` — the shared OpenCode config

Design details are in [docs/design.md](docs/design.md). The agent-facing contract is [docs/agent.md](docs/agent.md).

## Rules

- Do not edit the generated `opencode.jsonc`. Edit `config/opencode.base.json` or the machine-local `opencode.personal.json`.
- Treat `config/` changes as code review: they change permissions, MCP servers, and agent behavior for the whole team.
- Config updates flow one way: the owner edits and pushes, teammates run `rutter sync`. The CLI never commits or pushes.
- Do not create or push a release tag without approval.
- Keep the agent contract accurate. When adding a command, flag, exit code, or JSON field, update `rutter agent`, [docs/agent.md](docs/agent.md), and the README table.
- Add a test for new behavior. CLI tests live in `internal/cli/cli_test.go`; orchestration tests in `internal/app/app_test.go`.
- Keep human output short and put machine data behind `--json`.
