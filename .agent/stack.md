# Technology Stack

## Language & Runtime
- **Go**: 1.22+
- **Module**: `github.com/Shaked/agentic-repo`

## Dependencies
- `github.com/spf13/cobra` — CLI framework
- `github.com/fatih/color` — Terminal colors

## Tooling
| Tool | Version | Purpose |
|------|---------|---------|
| golangci-lint | 1.62+ | Linting (aggregates multiple linters) |
| gofmt | (bundled) | Code formatting |
| goimports | latest | Import management |
| pre-commit | 4.0+ | Git hook management |

## Project Layout
```
cmd/agentic-repo/main.go    — CLI entry point
internal/cli/               — Command implementations
internal/detector/          — Stack detection logic
internal/generator/         — Template rendering
internal/templates/         — Embedded template files
```

## Build Output
- Binary: `bin/agentic-repo`
- Install location: `$GOPATH/bin/agentic-repo`
