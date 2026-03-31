# Project Overview

## What It Does
Context0 is a CLI tool that brings AI-powered documentation to any software project. It initializes a `.ctx/` directory containing guides, templates, and generated docs, then uses AI coding assistants (Claude, Gemini, Copilot, Codex, OpenCode) to automatically generate and maintain project documentation.

## Main Files
- `main.go` - Entry point; embeds defaults and starts the CLI
- `cmd/` - All CLI commands: init, sync-docs, update, version
- `cmd/connector.go` - AI CLI integrations (Claude, Copilot, Gemini, Codex, OpenCode)
- `defaults/` - Embedded default files copied into projects on init
- `.ctx/guides/` - Operational guides for doc syncing and the Middleman concept
- `.ctx/templates/` - Templates for generating documentation
- `.ctx/docs/` - Generated project documentation
- `AGENTS.md` - Instructions for AI agents working in the codebase
- `CLAUDE.md` - Entry point that routes Claude to agent behavior
- `install.sh` / `install.ps1` - Platform installers
- `.github/workflows/release.yml` - Multi-platform build and release pipeline

## Flow
1. `ctx init` scaffolds the `.ctx/` directory with guides, templates, and agent config files in any project
2. `ctx sync-docs` generates or updates `.ctx/docs/` by sending the codebase to an AI connector with the doc guide and templates

## Documentation available in `.ctx/docs/`
- **`cli_and_connectors.md`** — CLI commands, AI connectors, update system, styling, installers, and CI pipeline
