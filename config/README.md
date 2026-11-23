# Shared config

`rutter` clones this repository and links the entries in this directory into OpenCode's global config directory.

```text
config/
├── AGENTS.md                    # global instructions
├── agents/                      # subagents
├── commands/                    # slash commands
├── skills/                      # skills
├── opencode.base.json           # team settings
├── opencode.personal.example.json
└── cli.example.json             # optional TUI template
```

## What gets linked

- `AGENTS.md`
- `agents/`
- `commands/`
- `skills/`

`opencode.base.json` is not linked. `rutter` merges it with the machine-local `opencode.personal.json` and writes `opencode.jsonc`.

## Merge rules

- Objects merge deeply.
- Personal scalars win.
- Arrays replace, except `permissions`, `instructions`, `skills`, and `plugins`, which append.

## Add a skill

Create `skills/<id>/SKILL.md` with frontmatter:

```markdown
---
name: Example
description: When the model should use this skill
---

Instructions go here.
```

Copy real files into the repository. Do not add symlinks that point outside the repository.

## Add an agent or command

- Agents: `agents/<name>.md`
- Commands: `commands/<name>.md`

Both support YAML frontmatter.

## Rules

- Do not put API keys, OAuth tokens, or machine-specific paths in these files.
- Treat config changes like code. They change what every teammate's agent can do.
- Do not edit the generated `opencode.jsonc`. Edit `opencode.base.json` or `opencode.personal.json`.
- Run `rutter sync` after pulling changes.
