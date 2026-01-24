package main

import (
	"os"

	"github.com/Shaked/agentic-repo/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
