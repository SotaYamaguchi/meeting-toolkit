# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`meeting-toolkit` is a CLI tool for organizing meeting documents before and after customer project meetings. It automates file renaming and collection based on project-specific prefixes.

## Development Commands

All commands are run from the repository root:

```bash
# Build
make build          # Creates ./mtg binary

# Testing
make test           # Run tests with coverage
go test -v          # Run tests verbose
go test -run TestFunctionName  # Run single test

# Linting
make lint           # Run golangci-lint
golangci-lint run   # Direct invocation

# Installation (local)
make install        # Install to ~/bin/mtg with config
make uninstall      # Remove binary and config
```

## Architecture

### Package Structure (Red Hat style)

Code is organized following Red Hat upstream project patterns (e.g., StackRox, OpenShift CLI):

```
.
├── main.go              # Entry point (cmd.Execute() only)
├── cmd/                 # Cobra subcommand definitions
│   ├── root.go         # Root command and subcommand registration
│   ├── prep.go         # prep subcommand
│   ├── memo.go         # memo subcommand
│   ├── mail.go         # mail / mail init subcommands
│   ├── list.go         # list subcommand
│   └── completion.go   # completion subcommand (cobra auto-generated)
└── pkg/                # Reusable business logic
    ├── config/         # Configuration management
    │   └── config.go   # Load, Save, ResolvePrefix
    ├── file/           # File operations
    │   └── operations.go  # Rename, Collect, ProcessPrep/Memo
    └── mail/           # Mail template handling
        └── template.go # Get, Parse, Format, CreateFile
```

**Subcommands:**
- `prep` - Rename files (main→date) and collect for pre-meeting
- `memo` - Rename files (main→date_MTG後) and collect for post-meeting
- `mail prep` / `mail memo` - Display mail template for project
- `mail init prep` / `mail init memo` - Create mail template file
- `list` - Show configured projects from config.json
- `completion` - Generate shell completion script (bash/zsh/fish/powershell)

**Core flow:**
1. `main.go` calls `cmd.Execute()` which runs the cobra root command
2. Cobra dispatches to the appropriate subcommand's `RunE` function
3. Subcommand resolves prefix via `pkg/config.ResolvePrefix()` and executes operations via `pkg/file.*()` or `pkg/mail.*()`

**Key packages:**
- `pkg/config` - Config struct, Load/Save, prefix resolution
- `pkg/file` - Rename/Collect files, ProcessPrep/ProcessMemo
- `pkg/mail` - Template parsing, formatting, file creation

**Configuration:**
- User config: `~/.config/mtg/config.json` (maps project names to file prefixes)
- Sample: `config.sample.json` (committed to repo)
- **Important:** Actual `config.json` contains customer information and is `.gitignore`d

### Testing Strategy

Tests in `main_test.go` use temporary directories and files:
- `config.Load` - JSON parsing and validation
- `config.ResolvePrefix` - Project name to prefix resolution
- `file.Rename` - File renaming with date/suffix
- `file.Collect` - File moving to destination folder
- `mail.Parse` - Email template parsing
- `mail.Format` - Email output formatting

Test files must handle cleanup with `defer func() { _ = os.RemoveAll(tmpDir) }()` pattern to satisfy errcheck linter.

## CI/CD

GitHub Actions (`.github/workflows/test.yml`) runs three jobs:
1. **Lint** - golangci-lint v2.11 (config: `.golangci.yml`)
2. **Test** - Unit tests with race detector, coverage displayed in Step Summary
3. **Build** - Binary compilation check

## pre-commit Hooks

`.pre-commit-config.yaml` runs on commit:
- golangci-lint (only on changed Go files)
- Standard checks (trailing whitespace, EOF, YAML syntax, file size)

Setup: `pre-commit install`

## Key Constraints

- **Minimal dependencies** - Uses cobra for CLI framework, otherwise Go standard library
- **Stateless** - No database, all config from JSON file
- **File-based** - Operates on filesystem directly using glob patterns
- **Date format** - Always YYYYMMDD (time.Now().Format("20060102"))
- **Japanese output** - All user messages and help text in Japanese
