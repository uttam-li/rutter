# Agent and script contract

`rutter` is designed to be driven by an AI agent or a script. It never prompts, never needs a TTY, and reports state in stable JSON.

## Discover the contract

```bash
rutter agent --json
```

The response is the source of truth for:

- every command and flag
- the JSON object each command returns
- exit codes
- environment variables
- which command is interactive

The same data is available as a readable guide with `rutter agent`.

## Streams

| Stream | Content |
| --- | --- |
| stdout | The success envelope or the human report |
| stderr | The failure envelope or the human error |

A command that fails before producing its result writes no data to stdout. `status --check` is the exception: it returns the status object, then exits non-zero when the check fails.

## Envelopes

Success:

```json
{
  "ok": true,
  "data": {}
}
```

Failure:

```json
{
  "ok": false,
  "error": {
    "code": "drift",
    "exit_code": 5,
    "message": "managed path changed outside rutter: /home/me/.config/opencode/AGENTS.md",
    "hint": "move or delete the conflicting path, then run `rutter sync`"
  }
}
```

`code` is stable. `message` and `hint` are for humans and may change.

## Exit codes

| Code | Name | Recovery |
| --- | --- | --- |
| 0 | `ok` | none |
| 1 | `error` | read the message |
| 2 | `usage` | fix the command; run `rutter agent` |
| 3 | `not_installed` | run `rutter install` |
| 4 | `environment` | run `rutter doctor` and fix the errors |
| 5 | `drift` | move or delete the conflicting path, then run `rutter sync` |
| 6 | `repository` | check the URL, network, and git credentials |
| 7 | `config` | fix the JSON in the base or personal config |

## Recipes

### Check before acting

```bash
rutter status --json
```

```json
{
  "ok": true,
  "data": {
    "installed": true,
    "config_dir": "/home/me/.config/opencode",
    "healthy": true,
    "links": [
      {"target": "/home/me/.config/opencode/AGENTS.md", "state": "ok"}
    ]
  }
}
```

If `installed` is false, run `rutter install`. If `healthy` is false, run `rutter sync`.

For shell branching, `rutter status --check` is shorter: exit `3` when not installed, `5` when drifted, `0` when healthy. In JSON mode, stdout still contains the status object; the exit code carries the check result.

### First install

```bash
rutter install --json
```

In a sandbox, redirect all state:

```bash
RUTTER_DATA_DIR=/tmp/rutter-data \
RUTTER_CONFIG_DIR=/tmp/opencode \
rutter install --json --no-reload
```

`RUTTER_REPO` replaces `--repo`. The flag wins over the environment variable.

### Receive an update

```bash
rutter sync --json
```

This is fast-forward only. If the checkout has local changes, git fails and `rutter` exits `6`; it never discards work.

### Diagnose

```bash
rutter doctor --json
```

Each check has `status` of `ok`, `warn`, or `error`. Only `error` changes the exit code. A missing OpenCode binary is a warning because `--config-dir` can be passed explicitly. In JSON mode, stdout still contains the report when the command exits `4`; read the exit code for the result.

### Read the generated config

```bash
rutter config show --json
```

The file on disk is JSONC with a generated header. The JSON form has comments stripped, so it is safe to parse.

### Remove

```bash
rutter uninstall --json
```

`--purge` also removes the checkout, backups, and manifest. If a managed path was changed outside rutter, uninstall exits `5` and leaves that path in place. Run `rutter sync` to relink it, or move it aside manually.

## Environment

| Variable | Purpose |
| --- | --- |
| `RUTTER_JSON` | `1` or `true` makes JSON the default |
| `RUTTER_REPO` | Config repository URL for `install` |
| `RUTTER_CONFIG_DIR` | OpenCode config directory for `install` |
| `RUTTER_DATA_DIR` | Checkout, backups, and manifest location |
| `EDITOR` / `VISUAL` | Editor for the interactive `config edit` command |

`RUTTER_DATA_DIR` is the clean way to run several isolated installs on one machine or inside a container.

## Rules for agents

- Prefer `status --json` before mutating anything.
- Treat `install` and `sync` as retryable. Both are idempotent.
- Do not edit `opencode.jsonc`; it is regenerated. Edit `config/opencode.base.json` in the repository, or `opencode.personal.json` on the machine.
- Do not call `config edit`; it blocks on an editor.
- Do not run git commands in the checkout. `rutter sync` owns that checkout.
- Do not push config changes unless the user owns the repository and asked for it. `rutter` never commits or pushes.
- On Windows, entries are copies. Run `rutter sync` after the repository changes.

## Windows note

Paths in JSON use backslashes, and link mode is `copy` instead of `link`. Treat `mode` as data, not as a guarantee that a path is a symlink. Copies are managed by rutter: `sync` refreshes them, so do not edit them in place.

## Versioning

The envelope, exit codes, and command names are part of the public contract. Additive fields may appear in `data`; agents should ignore unknown fields.
