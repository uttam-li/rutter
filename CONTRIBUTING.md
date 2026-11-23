# Contributing

## Requirements

- Go 1.27.4 or newer
- git
- OpenCode v2.3 or newer for manual testing
- golangci-lint v2 for linting
- GoReleaser v2 for release snapshots

## Before opening a pull request

```bash
make test
make lint
go vet ./...
```

Config changes live under `config/` and need the same review as code. They can change permissions, MCP servers, and agent behavior for the whole team.

## Commit style

Use short imperative subjects. Keep CLI changes and config changes in separate commits when they are unrelated.

## Adding a skill

Put real files under `config/skills/<id>/`. Do not add symlinks that point outside the repository.

## Platform notes

- macOS and Linux use symlinks.
- Windows uses copies and refreshes them on `rutter sync`.
- CI runs the test suite on all three platforms.
