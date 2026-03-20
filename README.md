# MDM — The Middleman

One agent to rule them all.

A CLI that turns any AI coding assistant into a Middleman — an orchestrator that delegates work to subagents without writing code itself.

## The problem

When working on complex projects with AI coding assistants, you end up managing multiple instances manually. MDM externalizes that cognitive load: you talk to one agent, and it decides how to split the work across subagents.

## The Middleman

The Middleman is an orchestration layer. It sits between you and the AI coding agents. You describe what you want, and the Middleman breaks it down, delegates each piece to a subagent in parallel, and returns control to you immediately.

**What it does:**
- Understands your request and splits it into independent concerns
- Spawns one subagent per concern, all in parallel
- Can run commands (build, test, git) but never writes application code itself

**What it doesn't do:**
- Write code — that's what subagents are for
- Block or poll — after delegating, control returns to you
- Chat — it only speaks when a decision is needed or a result matters

## Subagents

Subagents are the ones that actually write code. Every subagent follows the same three rules defined in `AGENTS.md`:

1. **Orient** — Read `.mdm/docs/project_overview.md` and list `.mdm/docs/` to understand the project
2. **Read before writing** — Read the specific doc(s) for the area they'll modify
3. **Keep docs in sync** — After making changes, update the affected doc(s)

This means subagents don't start from zero. The `.mdm/docs/` directory gives them the context they need to make informed changes without reading the entire codebase.

## The `.mdm/` directory

MDM initializes a `.mdm/` directory in your project. This is the knowledge base that both the Middleman and subagents rely on.

```
.mdm/
├── config.json        — default CLI and settings
├── guides/
│   ├── the_middleman.md    — how the Middleman operates
│   └── how_to_sync_docs.md — guide for doc generation
├── templates/
│   ├── doc_template.md            — template for feature docs
│   └── project_overview_template.md — template for the overview
└── docs/
    ├── project_overview.md   — high-level project summary
    └── *.md                  — one doc per feature/area
```

### Docs

Each doc in `.mdm/docs/` groups ~8-16 related files and describes what they do, how they connect, and the flow of data through them. Docs follow the template in `.mdm/templates/doc_template.md`: what the feature does, its main files, and its flow. No code snippets, no implementation details — just the big picture in plain language.

The `project_overview.md` is the entry point. It describes what the project does, its main files/directories, and links to all available docs.

### Templates

Templates define the structure that docs must follow. When you run `mdm sync-docs`, the AI CLI reads these templates and generates or updates the docs accordingly. You can edit the templates to change what gets documented and how.

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
curl -sL https://raw.githubusercontent.com/Titovilal/middleman/main/install.sh | sh
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/Titovilal/middleman/main/install.ps1 | iex
```

### From source

```bash
go install github.com/Titovilal/middleman@latest
mv $(go env GOPATH)/bin/middleman $(go env GOPATH)/bin/mdm
```

## Quick start

```bash
# Initialize .mdm/ in your project — asks which CLIs to integrate
mdm init

# Generate project documentation
mdm sync-docs

# Open a Middleman session with a request
mdm open "refactor the auth module and add tests"
```

## Commands

| Command | Description |
|---|---|
| `mdm init` | Initialize `.mdm/` and copy instruction files to project root |
| `mdm sync-docs` | Generate/update `.mdm/docs/` using an AI CLI |
| `mdm open [request]` | Open an interactive Middleman session |
| `mdm update` | Self-update to the latest version |
| `mdm version` | Print current version |

### Global flags

| Flag | Description |
|---|---|
| `--workdir`, `-w` | Project directory (default: current dir) |

### `mdm init` flags

| Flag | Description |
|---|---|
| `--mode`, `-m` | Init mode: `overwrite`, `fresh`, or `keep` (when `.mdm/` already exists) |
| `--clis` | Comma-separated CLIs to integrate (e.g. `claude,gemini,codex`) |
| `--default` | Default CLI for sync-docs and open |
| `--sync` | Run sync-docs after init |

### `mdm open` / `mdm sync-docs` flags

| Flag | Description |
|---|---|
| `--connector`, `-c` | AI CLI to use (overrides default from config) |

## Project structure

```
your-project/
├── AGENTS.md              — agent behavior: subagent + middleman roles
├── CLAUDE.md              — Claude instructions (if selected)
├── GEMINI.md              — Gemini instructions (if selected)
└── .mdm/
    ├── config.json
    ├── guides/
    ├── templates/
    └── docs/
```

## Architecture

```
mdm/
├── main.go          # Entry point, embeds defaults/
├── cmd/
│   ├── root.go      # CLI setup, init command, banner
│   ├── open.go      # mdm open — launch interactive Middleman session
│   ├── sync_docs.go # mdm sync-docs — doc generation
│   ├── connector.go # AI CLI connectors (claude, codex, gemini, copilot, opencode)
│   └── update.go    # Self-update
└── defaults/        # Embedded files copied on mdm init
    ├── AGENTS.md
    ├── CLAUDE.md
    ├── GEMINI.md
    ├── guides/
    └── templates/
```
