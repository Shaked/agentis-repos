package detector

import (
	"os"
	"path/filepath"
)

// StackType represents a detected project stack
type StackType string

const (
	StackGo      StackType = "go"
	StackPython  StackType = "python"
	StackNode    StackType = "node"
	StackJava    StackType = "java"
	StackUnknown StackType = "unknown"
)

func (s StackType) String() string {
	return string(s)
}

// Result represents a detected project at a specific path
type Result struct {
	Path  string
	Stack StackType
}

// Detector interface for stack-specific detection
type Detector interface {
	// Detect checks if the given directory contains this stack type
	Detect(path string) bool
	// Type returns the stack type this detector identifies
	Type() StackType
}

// All registered detectors in priority order
var detectors = []Detector{
	&GoDetector{},
	&PythonDetector{},
	&NodeDetector{},
	&JavaDetector{},
}

// Scan recursively scans a directory for project types
func Scan(root string) ([]Result, error) {
	var results []Result

	// Check root directory first
	if stack := detectStack(root); stack != StackUnknown {
		results = append(results, Result{Path: root, Stack: stack})
	}

	// Scan subdirectories (depth 1-2)
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() || isIgnoredDir(entry.Name()) {
			continue
		}

		subPath := filepath.Join(root, entry.Name())

		// Check this subdirectory
		if stack := detectStack(subPath); stack != StackUnknown {
			results = append(results, Result{Path: subPath, Stack: stack})
		}

		// Check one level deeper for monorepo structures like services/api/
		subEntries, err := os.ReadDir(subPath)
		if err != nil {
			continue
		}

		for _, subEntry := range subEntries {
			if !subEntry.IsDir() || isIgnoredDir(subEntry.Name()) {
				continue
			}

			deepPath := filepath.Join(subPath, subEntry.Name())
			if stack := detectStack(deepPath); stack != StackUnknown {
				results = append(results, Result{Path: deepPath, Stack: stack})
			}
		}
	}

	// Deduplicate - if root is detected, remove subdirectory matches of same type
	results = deduplicateResults(results, root)

	return results, nil
}

// detectStack checks a single directory for any known stack
func detectStack(path string) StackType {
	for _, d := range detectors {
		if d.Detect(path) {
			return d.Type()
		}
	}
	return StackUnknown
}

// IsMonorepo determines if the results indicate a monorepo structure
func IsMonorepo(results []Result) bool {
	if len(results) <= 1 {
		return false
	}

	// Multiple results = monorepo
	// Or single result not at root = monorepo
	return true
}

// isIgnoredDir returns true for directories that should be skipped
func isIgnoredDir(name string) bool {
	ignored := map[string]bool{
		".git":         true,
		".agent":       true,
		".claude":      true,
		"node_modules": true,
		"vendor":       true,
		"venv":         true,
		".venv":        true,
		"__pycache__":  true,
		"target":       true, // Java/Rust
		"bin":          true,
		"dist":         true,
		"build":        true,
	}
	return ignored[name]
}

// deduplicateResults removes redundant detections
func deduplicateResults(results []Result, root string) []Result {
	if len(results) <= 1 {
		return results
	}

	// If we have a root detection and subdirectory detections of the same type,
	// only keep the subdirectory ones (more specific)
	var rootStack StackType
	hasRootDetection := false

	for _, r := range results {
		if r.Path == root {
			rootStack = r.Stack
			hasRootDetection = true
			break
		}
	}

	if !hasRootDetection {
		return results
	}

	// Check if all subdirectory detections are the same type as root
	allSameAsRoot := true
	for _, r := range results {
		if r.Path != root && r.Stack != rootStack {
			allSameAsRoot = false
			break
		}
	}

	// If all same type, just keep root
	if allSameAsRoot {
		for _, r := range results {
			if r.Path == root {
				return []Result{r}
			}
		}
	}

	// Otherwise, keep all results but the root if there are subdirectories
	// with the same type
	var filtered []Result
	for _, r := range results {
		if r.Path == root {
			// Check if any subdirectory has same stack
			hasSameSubdir := false
			for _, sub := range results {
				if sub.Path != root && sub.Stack == rootStack {
					hasSameSubdir = true
					break
				}
			}
			if hasSameSubdir {
				continue // Skip root, keep subdirs
			}
		}
		filtered = append(filtered, r)
	}

	return filtered
}
