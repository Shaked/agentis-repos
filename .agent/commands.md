# CLI Commands Cheat Sheet

## Build & Run
```bash
# Build binary
make build

# Run directly
go run ./cmd/agentic-repo init

# Install to GOPATH
make install
```

## Testing
```bash
# Run all tests
make test

# Run specific test
go test -v -run TestDetectStack ./internal/detector/

# With coverage
go test -coverprofile=coverage.out ./...
```

## Linting & Formatting
```bash
# Run linter
make lint

# Auto-fix lint issues
golangci-lint run --fix

# Format code
make fmt
```

## Dependencies
```bash
# Add dependency
go get github.com/some/package

# Tidy modules
go mod tidy

# Vendor dependencies (optional)
go mod vendor
```

## Pre-commit Hooks
```bash
# Install hooks
pre-commit install

# Run manually
pre-commit run --all-files

# Update hook versions
pre-commit autoupdate
```

## Debugging
```bash
# Build with debug symbols
go build -gcflags="all=-N -l" -o bin/agentic-repo ./cmd/agentic-repo

# Run with delve
dlv exec ./bin/agentic-repo -- init
```

## Release
```bash
# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o bin/agentic-repo-linux ./cmd/agentic-repo
GOOS=darwin GOARCH=arm64 go build -o bin/agentic-repo-darwin ./cmd/agentic-repo
GOOS=windows GOARCH=amd64 go build -o bin/agentic-repo.exe ./cmd/agentic-repo
```
