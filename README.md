# 🤖 agentic-repo

**Transform any Git repository into a Hybrid Human/Agent Environment.**

Stop your AI assistant from hallucinating commands, wasting tokens on irrelevant files, or getting confused by inconsistent patterns.

---

## 🤖 Agentic Installation

Tell your AI coding assistant:

> Read https://raw.githubusercontent.com/Shaked/agentic-repo/main/INSTALL.md and install agentic-repo

Or see [INSTALL.md](./INSTALL.md) for manual installation.

---

## ⚡ Quick Start

```bash
# Install
go install github.com/Shaked/agentic-repo/cmd/agentic-repo@latest

# Initialize your repo
cd your-project
agentic-repo init

# Preview first (optional)
agentic-repo init --dry-run
```

---

## 🎯 The Problem

AI agents struggle with repositories because they:

- 🔥 **Waste tokens** on irrelevant files
- 😵 **Get confused** by inconsistent patterns
- 🤷 **Hallucinate** commands and configurations

## 💡 The Solution

**Context Routing** — Your agent loads a lightweight router file (`AGENTS.md`) first, then pulls specific context on-demand.

```
Agent reads AGENTS.md (< 100 tokens)
         ↓
    Needs to run tests?  →  Load .agent/testing.md
    Needs CLI commands?  →  Load .agent/commands.md
    Doing code review?   →  Load CODE_REVIEW_RULES.md
```

No more token waste. No more guessing.

---

## 📁 What Gets Generated

```
your-project/
├── AGENTS.md                 # 🗺️  Router (< 100 tokens)
├── CODE_REVIEW_RULES.md      # ✅ CI review requirements
├── repo-best-practices.md    # 📚 Team patterns
├── USAGE.md                  # 👤 Human-readable guide
├── Makefile                  # 🔧 Standard targets
├── .agentignore              # 🚫 Files to skip
├── .pre-commit-config.yaml   # 🔒 Enforcement hooks
├── .agent/
│   ├── stack.md              # 🛠️  Tech stack & versions
│   ├── testing.md            # 🧪 Testing patterns
│   └── commands.md           # 💻 CLI cheat sheet
├── .cursorrules              # Cursor AI integration
└── .claude/
    └── settings.json         # Claude integration
```

---

## 🔄 The Agent Workflow

Agents follow an iterative loop that ensures quality:

```
Make Changes → Run pre-commit → Run tests
                                    ↓
                              Passing?
                             /        \
                           No          Yes
                           ↓            ↓
                    Fix issues    Ready for review
                         ↓              ↓
                    (loop back)        CI
```

---

## 🧩 Supported Stacks

| Language | Package Manager | Linter | Formatter | Testing |
|----------|----------------|--------|-----------|---------|
| **Go** | go mod | golangci-lint | gofmt | go test |
| **Python** | uv | ruff | ruff | pytest |
| **Node/TS** | pnpm | eslint | prettier | vitest |
| **Java** | Maven | Checkstyle | Spotless | JUnit 5 |

---

## 📦 Monorepo Support

Auto-detects monorepos and creates hierarchical context:

```
monorepo/
├── AGENTS.md              # Root router
├── backend/
│   ├── AGENTS.md          # Python context
│   └── .agent/
├── frontend/
│   ├── AGENTS.md          # TypeScript context
│   └── .agent/
└── services/api/
    ├── AGENTS.md          # Go context
    └── .agent/
```

---

## 🚀 CLI Options

| Flag | Description |
|------|-------------|
| `--dry-run`, `-n` | Preview changes without writing |
| `--force`, `-f` | Overwrite existing files |
| `--verbose`, `-v` | Show detailed output |

---

## 🛠️ Development

```bash
make build    # Build binary
make test     # Run tests
make lint     # Run linter
```

---

## 📄 License

Apache 2.0
