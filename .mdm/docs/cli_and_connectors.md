# CLI and Connectors

## What It Does
MDM is a Cobra-based CLI that manages AI-powered project documentation. It initializes a `.mdm/` directory in any project, generates documentation by delegating to AI CLI tools (connectors), and can launch a Middleman agent session. It also handles self-updating from GitHub releases.

## Main Files
- `main.go` - Entry point; embeds the `defaults/` directory and starts Cobra
- `cmd/root.go` - Root command, `init` command, CLI integration selection, config management, and banner
- `cmd/sync_docs.go` - `sync-docs` command; builds a prompt from guides/templates and sends it to a connector
- `cmd/open.go` - `open` command; launches a Middleman agent session through a connector
- `cmd/connector.go` - Connector definitions for Claude, Copilot, Gemini, Codex, and OpenCode
- `cmd/update.go` - `update` and `version` commands; downloads latest binary from GitHub, handles cross-platform install
- `cmd/style.go` - Shared ANSI styling helpers used across all commands
- `defaults/` - Embedded default files (guides, templates, AGENTS.md, CLAUDE.md, GEMINI.md) copied on `mdm init`
- `install.sh` - Shell installer that downloads the latest release for Linux/macOS
- `install.ps1` - PowerShell installer for Windows
- `.github/workflows/release.yml` - CI workflow that builds multi-platform binaries and creates GitHub releases

## Flow
1. User runs `mdm init` in a project root, which creates `.mdm/` with guides, templates, and docs directory, copies agent instruction files (AGENTS.md, CLAUDE.md, etc.), and saves a config with the default CLI
2. User runs `mdm sync-docs` to generate documentation — the command reads the sync guide and templates, assembles a prompt, and sends it to the chosen AI connector which reads the codebase and writes docs into `.mdm/docs/`
3. User can run `mdm open <request>` to start a Middleman agent session — the command injects the Middleman guide into a prompt and delegates the user's request to the chosen AI connector
