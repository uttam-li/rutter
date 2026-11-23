# Design

## Destination

A Go CLI that installs, syncs, and uninstalls the shared OpenCode configuration on macOS, Windows, and Linux, with one repository and a repeatable release pipeline.

## Modules

| Module | Interface | Hides |
| --- | --- | --- |
| `internal/app` | `Install`, `Sync`, `Uninstall`, `Status`, `Doctor` | Ordering, repository sync, config generation, links, backups, manifest |
| `internal/config` | `MergeFiles`, `Merge` | JSONC parsing, merge rules, generated header |
| `internal/link` | `Plan`, `Apply`, `Prune`, `Remove`, `Check` | Symlinks, Windows copies, backups, drift detection |
| `internal/opencode` | `ConfigDir`, `Reload` | `opencode debug paths`, `opencode reload`, path parsing |
| `internal/source` | `Ensure`, `Sync` | `git clone`, `git pull --ff-only`, credential handling |
| `internal/state` | `Load`, `Save`, `Remove` | Manifest shape and persistence |
| `internal/platform` | `DataDir`, `WriteFileAtomic` | OS data directories and atomic replacement |
| `internal/cli` | Cobra commands | Flag parsing, JSON envelope, exit codes, self-description |

## The deep module

`app.Service` is the module the CLI calls. Its interface is five verbs, but each verb hides a sequence of filesystem and process work:

```go
Install(ctx, InstallOptions) (Result, error)
Sync(ctx, reload bool) (Result, error)
Uninstall(ctx, purge bool) (Result, error)
Status(ctx) (Status, error)
Doctor(ctx) DoctorReport
```

The CLI never merges JSON, creates a link, moves a backup, or shells out to git. Tests exercise the same seam with fakes for `Source`, `OpenCode`, and `Store`, and the real `config.Merger` and `link.Manager`.

Deletion test: remove `app.Service` and every command would have to rebuild the sequence and the manifest rules. The complexity reappears in five call sites, so the module earns its keep.

## Internal seams

- `Source` has one production adapter (`source.Git`) and test fakes. It could grow a pure-Go adapter later.
- `OpenCode` has one production adapter (`opencode.CLI`). Tests replace it, so install and sync never require a running OpenCode.
- `Store` has one production adapter (`state.Store`) and an in-memory test adapter.
- `link.Manager` hides the OS seam behind build-tagged `link_unix.go` and `link_windows.go`.

No seam exists without a second adapter or a demonstrated variation. The one exception is the platform split in `internal/link`, where the variation is the operating system itself.

## Invariants

- A file that rutter did not create is moved to a backup, never deleted.
- A managed entry is removed only when it still matches its source.
- The generated config carries a header and is replaced atomically.
- The manifest is saved only after the config and links succeed.
- A failed sync leaves the previous generated config in place until the new one is written.
- On Windows, copies are refreshed by `sync`; on Unix, symlinks follow the repository.
- The CLI never prompts and never requires a TTY.
- Success JSON goes to stdout; failure JSON goes to stderr.
- Exit codes and the `ok`/`data`/`error` envelope are stable.
- `rutter agent` is generated from the live command tree, so documented flags do not drift.

## Test surface

- `internal/config`: merge rules, JSONC, missing personal file, invalid roots.
- `internal/link`: link creation, copy creation, backup and restore, idempotent apply, prune.
- `internal/opencode`: path parsing for Unix and Windows output.
- `internal/app`: full install, status, sync, and uninstall with fakes.
- `internal/cli`: JSON envelopes, exit-code mapping, usage errors, `rutter agent` coverage.
- CI: Linux, macOS, and Windows matrix, plus the official `golang:1.27.4` container.

## Non-goals for v0.1

- A TUI.
- An OpenCode plugin runtime.
- A built-in git implementation.
- Config profiles beyond one team base and one personal file.
