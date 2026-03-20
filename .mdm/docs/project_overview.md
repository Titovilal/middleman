# Project Overview

## What It Does
MDM (The Middleman) is a CLI tool that brings AI-powered documentation and agent orchestration to any software project. It initializes a `.mdm/` directory containing guides, templates, and generated docs, then uses AI coding assistants (Claude, Gemini, Copilot, Codex, OpenCode) to automatically generate and maintain project documentation or act as a Middleman agent that delegates tasks to subagents.

## Main Files
- `main.go` - Entry point; embeds defaults and starts the CLI
- `cmd/` - All CLI commands: init, sync-docs, open, update, version
- `cmd/connector.go` - AI CLI integrations (Claude, Copilot, Gemini, Codex, OpenCode)
- `defaults/` - Embedded default files copied into projects on init
- `.mdm/guides/` - Operational guides for the Middleman and doc syncing
- `.mdm/templates/` - Templates for generating documentation
- `.mdm/docs/` - Generated project documentation
- `AGENTS.md` - Instructions for AI subagents working in the codebase
- `CLAUDE.md` - Entry point that routes Claude to agent or middleman behavior
- `install.sh` / `install.ps1` - Platform installers
- `.github/workflows/release.yml` - Multi-platform build and release pipeline

## Flow
1. `mdm init` scaffolds the `.mdm/` directory with guides, templates, and agent config files in any project
2. `mdm sync-docs` generates or updates `.mdm/docs/` by sending the codebase to an AI connector with the doc guide and templates
3. `mdm open <request>` launches a Middleman agent session that reads the guide, delegates work to AI subagents, and returns results

## Documentation available in `.mdm/docs/`
- **`cli_and_connectors.md`** — CLI commands, AI connectors, update system, styling, installers, and CI pipeline
