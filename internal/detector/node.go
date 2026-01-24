package detector

import (
	"os"
	"path/filepath"
)

// NodeDetector detects Node.js/TypeScript projects
type NodeDetector struct{}

// Detect checks for Node.js/TypeScript project indicators
func (d *NodeDetector) Detect(path string) bool {
	indicators := []string{
		"package.json",
		"pnpm-lock.yaml",
		"package-lock.json",
		"yarn.lock",
		"tsconfig.json",
	}

	for _, indicator := range indicators {
		if _, err := os.Stat(filepath.Join(path, indicator)); err == nil {
			return true
		}
	}

	return false
}

// Type returns the stack type
func (d *NodeDetector) Type() StackType {
	return StackNode
}
