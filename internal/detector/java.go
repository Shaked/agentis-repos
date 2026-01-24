package detector

import (
	"os"
	"path/filepath"
)

// JavaDetector detects Java projects
type JavaDetector struct{}

// Detect checks for Java project indicators
func (d *JavaDetector) Detect(path string) bool {
	indicators := []string{
		"pom.xml",
		"build.gradle",
		"build.gradle.kts",
		"settings.gradle",
		"settings.gradle.kts",
	}

	for _, indicator := range indicators {
		if _, err := os.Stat(filepath.Join(path, indicator)); err == nil {
			return true
		}
	}

	// Also check for .mvn directory (Maven Wrapper)
	if _, err := os.Stat(filepath.Join(path, ".mvn")); err == nil {
		return true
	}

	return false
}

// Type returns the stack type
func (d *JavaDetector) Type() StackType {
	return StackJava
}
