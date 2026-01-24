package detector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestJavaDetector_Detect(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		dirs     []string // directories to create
		expected bool
	}{
		{
			name:     "detects pom.xml",
			files:    []string{"pom.xml"},
			expected: true,
		},
		{
			name:     "detects build.gradle",
			files:    []string{"build.gradle"},
			expected: true,
		},
		{
			name:     "detects build.gradle.kts",
			files:    []string{"build.gradle.kts"},
			expected: true,
		},
		{
			name:     "detects settings.gradle",
			files:    []string{"settings.gradle"},
			expected: true,
		},
		{
			name:     "detects settings.gradle.kts",
			files:    []string{"settings.gradle.kts"},
			expected: true,
		},
		{
			name:     "detects .mvn directory",
			dirs:     []string{".mvn"},
			expected: true,
		},
		{
			name:     "detects multiple Java indicators",
			files:    []string{"pom.xml", "settings.gradle"},
			dirs:     []string{".mvn"},
			expected: true,
		},
		{
			name:     "no java files returns false",
			files:    []string{"main.py", "requirements.txt"},
			expected: false,
		},
		{
			name:     "empty directory returns false",
			files:    []string{},
			expected: false,
		},
		{
			name:     "java files in subdirectory not detected at root",
			files:    []string{"subdir/pom.xml"},
			expected: false,
		},
	}

	detector := &JavaDetector{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			// Create files
			for _, f := range tt.files {
				path := filepath.Join(dir, f)
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
				if err := os.WriteFile(path, []byte{}, 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
			}

			// Create directories
			for _, d := range tt.dirs {
				path := filepath.Join(dir, d)
				if err := os.MkdirAll(path, 0755); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
			}

			result := detector.Detect(dir)
			if result != tt.expected {
				t.Errorf("JavaDetector.Detect() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestJavaDetector_Type(t *testing.T) {
	detector := &JavaDetector{}
	if detector.Type() != StackJava {
		t.Errorf("JavaDetector.Type() = %v, want %v", detector.Type(), StackJava)
	}
}
