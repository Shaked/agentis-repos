# Testing Standards

## Framework
- Standard `testing` package
- **Table-driven tests** are mandatory

## Test File Naming
- Test files: `*_test.go`
- Same package as code being tested

## Table-Driven Test Pattern

```go
func TestDetectStack(t *testing.T) {
    tests := []struct {
        name     string
        files    []string
        expected StackType
    }{
        {
            name:     "detects Go project",
            files:    []string{"go.mod", "main.go"},
            expected: StackGo,
        },
        {
            name:     "detects Python project",
            files:    []string{"pyproject.toml", "main.py"},
            expected: StackPython,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := DetectStack(tt.files)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Test Commands
```bash
# Run all tests
go test ./...

# With race detection
go test -race ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Verbose output
go test -v ./...
```

## Assertions
- Use standard `t.Errorf` / `t.Fatalf`
- Or `github.com/stretchr/testify/assert` for convenience

## Test Fixtures
- Place in `testdata/` directory
- Use `os.MkdirTemp` for temp directories
- Clean up with `t.Cleanup()`

## What to Test
- All public functions
- Edge cases (empty input, nil, errors)
- Detector logic for each stack type
- Template rendering output
- Monorepo detection scenarios
