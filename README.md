# Context0 — AI Documentation Manager

AI-powered documentation for any project.

## The problem

AI coding assistants work better when they understand your project. Context0 generates and maintains structured documentation in a `.ctx/docs/` directory that agents can read before touching code — so they don't start from zero.

## The Middleman Concept

The Middleman is an orchestration pattern where one agent sits between you and AI coding subagents. You describe what you want, and the Middleman breaks it down, delegates each piece to a subagent in parallel, and returns control immediately. See `.ctx/guides/the_middleman.md` for details.

## The `.ctx/` directory

Context0 initializes a `.ctx/` directory in your project:

```
.ctx/
├── config.json        — default CLI and settings
├── guides/
│   ├── the_middleman.md    — the Middleman concept explained
│   └── how_to_sync_docs.md — guide for doc generation
├── templates/
│   ├── doc_template.md            — template for feature docs
│   └── project_overview_template.md — template for the overview
└── docs/
    ├── project_overview.md   — high-level project summary
    └── *.md                  — one doc per feature/area
```

### Docs

Each doc in `.ctx/docs/` groups ~8-16 related files and describes what they do, how they connect, and the flow of data through them. Docs follow the template in `.ctx/templates/doc_template.md`: what the feature does, its main files, and its flow. No code snippets, no implementation details — just the big picture in plain language.

The `project_overview.md` is the entry point. It describes what the project does, its main files/directories, and links to all available docs.

### Templates

Templates define the structure that docs must follow. When you run `ctx sync-docs`, the AI CLI reads these templates and generates or updates the docs accordingly. You can edit the templates to change what gets documented and how.

## Supported CLIs

| CLI | Instruction file | Status |
|---|---|---|
| Claude Code | `CLAUDE.md` + `AGENTS.md` | Tested |
| Codex | `AGENTS.md` | Untested |
| Copilot | `AGENTS.md` | Untested |
| Gemini CLI | `GEMINI.md` + `AGENTS.md` | Untested |
| OpenCode | `AGENTS.md` | Untested |

## Installation

### Linux / macOS

```bash
curl -sL https://raw.githubusercontent.com/Titovilal/context0/main/install.sh | sh
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/Titovilal/context0/main/install.ps1 | iex
```

### From source

```bash
go install github.com/Titovilal/context0@latest
mv $(go env GOPATH)/bin/context0 $(go env GOPATH)/bin/ctx
```

## Quick start

```bash
# Initialize .ctx/ in your project — asks which CLIs to integrate
ctx init

# Generate project documentation
ctx sync-docs
```

## Commands

| Command | Description |
|---|---|
| `ctx init` | Initialize `.ctx/` and copy instruction files to project root |
| `ctx sync-docs` | Generate/update `.ctx/docs/` using an AI CLI |
| `ctx update` | Self-update to the latest version |
| `ctx version` | Print current version |

### Global flags

| Flag | Description |
|---|---|
| `--workdir`, `-w` | Project directory (default: current dir) |

### `ctx init` flags

| Flag | Description |
|---|---|
| `--mode`, `-m` | Init mode: `overwrite`, `fresh`, or `keep` (when `.ctx/` already exists) |
| `--clis` | Comma-separated CLIs to integrate (e.g. `claude,gemini,codex`) |
| `--default` | Default CLI for sync-docs |
| `--sync` | Run sync-docs after init |

### `ctx sync-docs` flags

| Flag | Description |
|---|---|
| `--connector`, `-c` | AI CLI to use (overrides default from config) |

## Project structure

```
your-project/
├── AGENTS.md              — agent behavior rules
├── CLAUDE.md              — Claude instructions (if selected)
├── GEMINI.md              — Gemini instructions (if selected)
└── .ctx/
    ├── config.json
    ├── guides/
    ├── templates/
    └── docs/
```

## Architecture

```
ctx/
├── main.go          # Entry point, embeds defaults/
├── cmd/
│   ├── root.go      # CLI setup, init command, banner
│   ├── sync_docs.go # ctx sync-docs — doc generation
│   ├── connector.go # AI CLI connectors (claude, codex, gemini, copilot, opencode)
│   └── update.go    # Self-update
└── defaults/        # Embedded files copied on ctx init
    ├── AGENTS.md
    ├── CLAUDE.md
    ├── GEMINI.md
    ├── guides/
    └── templates/
```
