package detector

import (
	"testing"
)

func TestGoDetector_Detect(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected bool
	}{
		{
			name:     "detects go.mod",
			files:    []string{"go.mod"},
			expected: true,
		},
		{
			name:     "detects go.sum",
			files:    []string{"go.sum"},
			expected: true,
		},
		{
			name:     "detects both go.mod and go.sum",
			files:    []string{"go.mod", "go.sum"},
			expected: true,
		},
		{
			name:     "no go files returns false",
			files:    []string{"main.py", "requirements.txt"},
			expected: false,
		},
		{
			name:     "empty directory returns false",
			files:    []string{},
			expected: false,
		},
		{
			name:     "go files in subdirectory not detected at root",
			files:    []string{"subdir/go.mod"},
			expected: false,
		},
	}

	detector := &GoDetector{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := createTempProject(t, tt.files)
			result := detector.Detect(dir)
			if result != tt.expected {
				t.Errorf("GoDetector.Detect() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGoDetector_Type(t *testing.T) {
	detector := &GoDetector{}
	if detector.Type() != StackGo {
		t.Errorf("GoDetector.Type() = %v, want %v", detector.Type(), StackGo)
	}
}
