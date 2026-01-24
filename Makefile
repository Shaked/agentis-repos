.PHONY: build test lint fmt clean install

# Build the CLI binary
build:
	go build -o bin/agentic-repo ./cmd/agentic-repo

# Run all tests with race detection
test:
	go test -v -race ./...

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	gofmt -w .
	goimports -w .

# Clean build artifacts
clean:
	rm -rf bin/

# Install binary to GOPATH
install: build
	cp bin/agentic-repo $(GOPATH)/bin/
