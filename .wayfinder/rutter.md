# rutter v0.1

## Destination

A Go CLI that installs, syncs, and uninstalls the shared OpenCode configuration on macOS, Windows, and Linux, with one repository and a repeatable release pipeline.

## Notes

- Go 1.27.4.
- One repository: CLI at the root, config under `config/`.
- Module design is in `docs/design.md`.
- Docker and the GitHub Actions matrix are the test environments.

## Decisions so far

- Name: `rutter`.
- Language: Go.
- One repository for CLI and config.
- Config source: `config/` in the same repository.
- Link strategy: symlinks on Unix, copies on Windows.
- Config generation: base + personal merge with a generated-file header.
- Distribution: GitHub Releases built by GoReleaser.

## Not yet specified

- The `add` and `remove` command surface.
- Homebrew and Scoop channels.
- Config schema versioning.
- Multiple team profiles.

## Out of scope

- A TUI.
- An OpenCode plugin runtime.
- Running the Windows test suite on actual Windows ARM64 hardware. CI builds arm64 and tests amd64.
