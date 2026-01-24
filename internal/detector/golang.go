package detector

import (
	"os"
	"path/filepath"
)

// GoDetector detects Go projects
type GoDetector struct{}

// Detect checks for Go project indicators
func (d *GoDetector) Detect(path string) bool {
	indicators := []string{
		"go.mod",
		"go.sum",
	}

	for _, indicator := range indicators {
		if _, err := os.Stat(filepath.Join(path, indicator)); err == nil {
			return true
		}
	}

	return false
}

// Type returns the stack type
func (d *GoDetector) Type() StackType {
	return StackGo
}
