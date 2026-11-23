# Contributing

## Requirements

- Go 1.27.4 or newer
- git
- make
- jq
- OpenCode v2.3 or newer for manual testing
- golangci-lint v2.15.0 for linting
- GoReleaser v2 for release snapshots and release-config validation

## Local gate (required)

Run this before every commit and push:

```bash
make check
```

The gate runs formatting, `go vet`, tests, a build, `golangci-lint`, config validation, and the agent contract check. Do not commit or push a change that fails it, and do not use GitHub Actions to discover failures that can run locally.

Install the pre-push hook once:

```bash
make hooks
```

The hook runs `make check` automatically. Use `git push --no-verify` only in an emergency, never as a shortcut.

When `.goreleaser.yaml` or the release assets change, also run:

```bash
make check-release
```

GitHub Actions is reserved for what a laptop cannot prove: the macOS and Windows test matrix, the Docker image, and release validation. Pull requests run the Ubuntu gate; the full matrix runs on pushes to `main` and on PRs labeled `ci:full`.

## Before opening a pull request

```bash
make check
```

Config changes live under `config/` and need the same review as code. They can change permissions, MCP servers, and agent behavior for the whole team.

The CLI is part of a contract with scripts and AI agents. If you add a command, flag, exit code, or JSON field, update `rutter agent` and [docs/agent.md](docs/agent.md) in the same change.

## Commit style

Use short imperative subjects. Keep CLI changes and config changes in separate commits when they are unrelated.

## Adding a skill

Put real files under `config/skills/<id>/`. Do not add symlinks that point outside the repository.

## Platform notes

- macOS and Linux use symlinks.
- Windows uses copies and refreshes them on `rutter sync`.
- CI runs the Ubuntu gate on every pull request, the macOS and Windows matrix on pushes to `main`, and the full matrix on PRs labeled `ci:full`.
