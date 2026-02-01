# Installing agentic-repo

## Prerequisites

### Go 1.22+

```bash
go version
```

Expected output: `go version go1.22.x` or higher.

**If Go is not installed:**
- macOS: `brew install go`
- Linux: See https://go.dev/doc/install
- Windows: Download from https://go.dev/dl/

### Git

```bash
git --version
```

Expected output: `git version 2.x.x`

---

## Method 1: Quick Install (Recommended)

```bash
go install github.com/Shaked/agentic-repo/cmd/agentic-repo@latest
```

### Verify Installation

```bash
agentic-repo --help
```

Expected: Help output showing available commands.

If you get `command not found`, ensure `$(go env GOPATH)/bin` is in your PATH:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

---

## Method 2: Build from Source

### Clone the Repository

```bash
git clone https://github.com/Shaked/agentic-repo.git
cd agentic-repo
```

### Build

```bash
make build
```

### Verify Build

```bash
./bin/agentic-repo --help
```

Expected: Help output showing available commands.

### Optional: Add to PATH

```bash
# Add to your shell profile (~/.zshrc, ~/.bashrc, etc.)
export PATH="$PATH:/path/to/agentic-repo/bin"
```

Or copy to GOPATH:

```bash
make install
```

---

## Quick Test

After installation, test by initializing a sample project:

```bash
mkdir /tmp/test-agentic && cd /tmp/test-agentic
git init
agentic-repo init --dry-run
```

Expected: Preview of files that would be generated.

---

## Troubleshooting

### `command not found: agentic-repo`

Your Go bin directory is not in PATH. Fix:

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

### `go install` fails with network error

Check your Go proxy settings:

```bash
go env GOPROXY
```

If needed, reset to default:

```bash
go env -w GOPROXY=https://proxy.golang.org,direct
```

### Build fails with missing dependencies

```bash
go mod download
make build
```
