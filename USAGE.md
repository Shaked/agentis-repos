# agentic-repo CLI Usage Guide

A CLI tool that initializes repositories with the **Agent-Native Repository Standard** — optimizing codebases for AI agent interaction through structured context files.

## Installation

```bash
# Build from source
make build

# Binary will be at ./bin/agentic-repo
```

## Quick Start

```bash
# Initialize current directory (auto-detects project type)
agentic-repo init

# Initialize a specific directory
agentic-repo init ./my-project

# Preview what would be generated (dry run)
agentic-repo init --dry-run

# Force overwrite existing files
agentic-repo init --force
```

## What It Does

Running `agentic-repo init` will:

1. **Detect your project type** — Scans for `go.mod`, `package.json`, `pyproject.toml`, `pom.xml`, etc.
2. **Detect monorepos** — If multiple project types exist in subdirectories, creates a hierarchical structure
3. **Generate context files**:
   - `AGENTS.md` — Lightweight router (<100 tokens) for AI agents
   - `USAGE.md` — Human-readable usage guide
   - `.agent/stack.md` — Technology stack and versions
   - `.agent/testing.md` — Testing patterns and requirements
   - `.agent/commands.md` — CLI commands cheat sheet
   - `Makefile` — Standard build/test/lint targets
   - `.agentignore` — Files AI agents should skip
   - `.pre-commit-config.yaml` — Linting enforcement hooks
4. **Create integration stubs** — `.cursorrules`, `.claude/` for AI tool compatibility

## Supported Stacks

| Language | Package Manager | Linter | Formatter | Testing |
|----------|----------------|--------|-----------|---------|
| **Go** | go mod | golangci-lint | gofmt | go test (table-driven) |
| **Python** | uv | ruff | ruff | pytest |
| **Node/TS** | pnpm | eslint | prettier | vitest/jest |
| **Java** | Maven (mvnw) | Checkstyle | Spotless | JUnit 5 |

## Monorepo Support

For monorepos with multiple languages, the tool creates:

```
monorepo/
├── AGENTS.md              # Root router → references subdirs
├── USAGE.md
├── .agent/overview.md     # High-level architecture
├── backend/               # Python service
│   ├── AGENTS.md
│   └── .agent/
├── frontend/              # TypeScript app
│   ├── AGENTS.md
│   └── .agent/
└── services/api/          # Go service
    ├── AGENTS.md
    └── .agent/
```

## CLI Flags

| Flag | Description |
|------|-------------|
| `--dry-run` | Preview generated files without writing |
| `--force` | Overwrite existing files |
| `--verbose` | Show detailed detection and generation logs |

## The Agent Workflow

1. **Agent reads `AGENTS.md`** — Gets the map of the repository
2. **Agent loads context on demand**:
   - Bug fix? → Load `.agent/stack.md`
   - Write test? → Load `.agent/testing.md`
   - Run something? → Load `.agent/commands.md`
3. **Agent executes with confidence** — Exact commands, no guessing
4. **Pre-commit hooks validate** — Linting failures block commits, agent self-corrects

## Development

```bash
# Run tests
make test

# Lint code
make lint

# Format code
make fmt

# Clean build artifacts
make clean
```
