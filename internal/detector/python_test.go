package detector

import (
	"testing"
)

func TestPythonDetector_Detect(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected bool
	}{
		{
			name:     "detects pyproject.toml",
			files:    []string{"pyproject.toml"},
			expected: true,
		},
		{
			name:     "detects requirements.txt",
			files:    []string{"requirements.txt"},
			expected: true,
		},
		{
			name:     "detects setup.py",
			files:    []string{"setup.py"},
			expected: true,
		},
		{
			name:     "detects setup.cfg",
			files:    []string{"setup.cfg"},
			expected: true,
		},
		{
			name:     "detects uv.lock",
			files:    []string{"uv.lock"},
			expected: true,
		},
		{
			name:     "detects poetry.lock",
			files:    []string{"poetry.lock"},
			expected: true,
		},
		{
			name:     "detects Pipfile",
			files:    []string{"Pipfile"},
			expected: true,
		},
		{
			name:     "detects multiple Python indicators",
			files:    []string{"pyproject.toml", "requirements.txt", "poetry.lock"},
			expected: true,
		},
		{
			name:     "no python files returns false",
			files:    []string{"main.go", "go.mod"},
			expected: false,
		},
		{
			name:     "empty directory returns false",
			files:    []string{},
			expected: false,
		},
		{
			name:     "python files in subdirectory not detected at root",
			files:    []string{"subdir/pyproject.toml"},
			expected: false,
		},
	}

	detector := &PythonDetector{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := createTempProject(t, tt.files)
			result := detector.Detect(dir)
			if result != tt.expected {
				t.Errorf("PythonDetector.Detect() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPythonDetector_Type(t *testing.T) {
	detector := &PythonDetector{}
	if detector.Type() != StackPython {
		t.Errorf("PythonDetector.Type() = %v, want %v", detector.Type(), StackPython)
	}
}
