package detector

import (
	"os"
	"path/filepath"
)

// PythonDetector detects Python projects
type PythonDetector struct{}

// Detect checks for Python project indicators
func (d *PythonDetector) Detect(path string) bool {
	indicators := []string{
		"pyproject.toml",
		"requirements.txt",
		"setup.py",
		"setup.cfg",
		"uv.lock",
		"poetry.lock",
		"Pipfile",
	}

	for _, indicator := range indicators {
		if _, err := os.Stat(filepath.Join(path, indicator)); err == nil {
			return true
		}
	}

	return false
}

// Type returns the stack type
func (d *PythonDetector) Type() StackType {
	return StackPython
}
