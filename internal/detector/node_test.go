package detector

import (
	"testing"
)

func TestNodeDetector_Detect(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected bool
	}{
		{
			name:     "detects package.json",
			files:    []string{"package.json"},
			expected: true,
		},
		{
			name:     "detects pnpm-lock.yaml",
			files:    []string{"pnpm-lock.yaml"},
			expected: true,
		},
		{
			name:     "detects package-lock.json",
			files:    []string{"package-lock.json"},
			expected: true,
		},
		{
			name:     "detects yarn.lock",
			files:    []string{"yarn.lock"},
			expected: true,
		},
		{
			name:     "detects tsconfig.json",
			files:    []string{"tsconfig.json"},
			expected: true,
		},
		{
			name:     "detects multiple Node indicators",
			files:    []string{"package.json", "package-lock.json", "tsconfig.json"},
			expected: true,
		},
		{
			name:     "no node files returns false",
			files:    []string{"main.py", "requirements.txt"},
			expected: false,
		},
		{
			name:     "empty directory returns false",
			files:    []string{},
			expected: false,
		},
		{
			name:     "node files in subdirectory not detected at root",
			files:    []string{"subdir/package.json"},
			expected: false,
		},
	}

	detector := &NodeDetector{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := createTempProject(t, tt.files)
			result := detector.Detect(dir)
			if result != tt.expected {
				t.Errorf("NodeDetector.Detect() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNodeDetector_Type(t *testing.T) {
	detector := &NodeDetector{}
	if detector.Type() != StackNode {
		t.Errorf("NodeDetector.Type() = %v, want %v", detector.Type(), StackNode)
	}
}
