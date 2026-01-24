package detector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStackType_String(t *testing.T) {
	tests := []struct {
		name     string
		stack    StackType
		expected string
	}{
		{"go stack", StackGo, "go"},
		{"python stack", StackPython, "python"},
		{"node stack", StackNode, "node"},
		{"java stack", StackJava, "java"},
		{"unknown stack", StackUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.stack.String()
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestScan(t *testing.T) {
	tests := []struct {
		name           string
		files          []string
		expectedStacks []StackType
		expectedCount  int
	}{
		{
			name:           "detects Go project at root",
			files:          []string{"go.mod", "main.go"},
			expectedStacks: []StackType{StackGo},
			expectedCount:  1,
		},
		{
			name:           "detects Python project at root",
			files:          []string{"pyproject.toml", "main.py"},
			expectedStacks: []StackType{StackPython},
			expectedCount:  1,
		},
		{
			name:           "detects Node project at root",
			files:          []string{"package.json", "index.js"},
			expectedStacks: []StackType{StackNode},
			expectedCount:  1,
		},
		{
			name:           "detects Java project at root",
			files:          []string{"pom.xml", "src/Main.java"},
			expectedStacks: []StackType{StackJava},
			expectedCount:  1,
		},
		{
			name:           "empty directory returns no results",
			files:          []string{},
			expectedStacks: []StackType{},
			expectedCount:  0,
		},
		{
			name: "detects monorepo with multiple projects",
			files: []string{
				"services/api/go.mod",
				"services/web/package.json",
			},
			expectedStacks: []StackType{StackGo, StackNode},
			expectedCount:  2,
		},
		{
			name: "detects nested project in subdir",
			files: []string{
				"backend/go.mod",
			},
			expectedStacks: []StackType{StackGo},
			expectedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := createTempProject(t, tt.files)

			results, err := Scan(dir)
			if err != nil {
				t.Fatalf("Scan() error = %v", err)
			}

			if len(results) != tt.expectedCount {
				t.Errorf("got %d results, want %d", len(results), tt.expectedCount)
			}

			// Check that all expected stacks are found
			for _, expectedStack := range tt.expectedStacks {
				found := false
				for _, r := range results {
					if r.Stack == expectedStack {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected stack %s not found in results", expectedStack)
				}
			}
		})
	}
}

func TestScan_ErrorOnNonExistentDir(t *testing.T) {
	_, err := Scan("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("expected error for non-existent directory, got nil")
	}
}

func TestIsMonorepo(t *testing.T) {
	tests := []struct {
		name     string
		results  []Result
		expected bool
	}{
		{
			name:     "empty results is not monorepo",
			results:  []Result{},
			expected: false,
		},
		{
			name: "single result is not monorepo",
			results: []Result{
				{Path: "/project", Stack: StackGo},
			},
			expected: false,
		},
		{
			name: "multiple results is monorepo",
			results: []Result{
				{Path: "/project", Stack: StackGo},
				{Path: "/project/frontend", Stack: StackNode},
			},
			expected: true,
		},
		{
			name: "multiple same-type results is monorepo",
			results: []Result{
				{Path: "/project/service1", Stack: StackGo},
				{Path: "/project/service2", Stack: StackGo},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsMonorepo(tt.results)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsIgnoredDir(t *testing.T) {
	tests := []struct {
		name     string
		dirName  string
		expected bool
	}{
		{"git directory", ".git", true},
		{"agent directory", ".agent", true},
		{"claude directory", ".claude", true},
		{"node_modules", "node_modules", true},
		{"vendor", "vendor", true},
		{"venv", "venv", true},
		{"dot venv", ".venv", true},
		{"pycache", "__pycache__", true},
		{"target", "target", true},
		{"bin", "bin", true},
		{"dist", "dist", true},
		{"build", "build", true},
		{"src directory", "src", false},
		{"cmd directory", "cmd", false},
		{"internal directory", "internal", false},
		{"services directory", "services", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIgnoredDir(tt.dirName)
			if result != tt.expected {
				t.Errorf("isIgnoredDir(%q) = %v, want %v", tt.dirName, result, tt.expected)
			}
		})
	}
}

func TestDeduplicateResults(t *testing.T) {
	root := "/project"

	tests := []struct {
		name          string
		results       []Result
		expectedCount int
		description   string
	}{
		{
			name:          "empty results",
			results:       []Result{},
			expectedCount: 0,
			description:   "empty input returns empty output",
		},
		{
			name: "single result unchanged",
			results: []Result{
				{Path: root, Stack: StackGo},
			},
			expectedCount: 1,
			description:   "single result is not deduplicated",
		},
		{
			name: "root and subdirs same type keeps only root",
			results: []Result{
				{Path: root, Stack: StackGo},
				{Path: filepath.Join(root, "cmd"), Stack: StackGo},
				{Path: filepath.Join(root, "internal"), Stack: StackGo},
			},
			expectedCount: 1,
			description:   "all same type as root, keep only root",
		},
		{
			name: "root and subdirs different types keeps all except root same-type subdirs",
			results: []Result{
				{Path: root, Stack: StackGo},
				{Path: filepath.Join(root, "frontend"), Stack: StackNode},
			},
			expectedCount: 2,
			description:   "different types are preserved",
		},
		{
			name: "no root detection keeps all",
			results: []Result{
				{Path: filepath.Join(root, "service1"), Stack: StackGo},
				{Path: filepath.Join(root, "service2"), Stack: StackPython},
			},
			expectedCount: 2,
			description:   "no root detection means all results kept",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deduplicateResults(tt.results, root)
			if len(result) != tt.expectedCount {
				t.Errorf("%s: got %d results, want %d", tt.description, len(result), tt.expectedCount)
			}
		})
	}
}

func TestDetectStack(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected StackType
	}{
		{
			name:     "detects Go",
			files:    []string{"go.mod"},
			expected: StackGo,
		},
		{
			name:     "detects Python",
			files:    []string{"pyproject.toml"},
			expected: StackPython,
		},
		{
			name:     "detects Node",
			files:    []string{"package.json"},
			expected: StackNode,
		},
		{
			name:     "detects Java",
			files:    []string{"pom.xml"},
			expected: StackJava,
		},
		{
			name:     "returns unknown for empty dir",
			files:    []string{},
			expected: StackUnknown,
		},
		{
			name:     "returns unknown for unrecognized files",
			files:    []string{"random.txt", "data.csv"},
			expected: StackUnknown,
		},
		{
			name:     "Go takes priority over Python",
			files:    []string{"go.mod", "pyproject.toml"},
			expected: StackGo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := createTempProject(t, tt.files)

			result := detectStack(dir)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestScan_IgnoresSpecialDirectories(t *testing.T) {
	// Create a project with ignored directories that contain project markers
	dir := t.TempDir()

	// Create root go.mod
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte{}, 0644)

	// Create ignored directories with project markers
	ignoredDirs := []string{"node_modules", "vendor", ".git", "build"}
	for _, ignored := range ignoredDirs {
		ignoredPath := filepath.Join(dir, ignored)
		os.MkdirAll(ignoredPath, 0755)
		// Put a project marker in each ignored directory
		os.WriteFile(filepath.Join(ignoredPath, "package.json"), []byte{}, 0644)
	}

	results, err := Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	// Should only find the root Go project, not the projects in ignored directories
	if len(results) != 1 {
		t.Errorf("expected 1 result (root Go project), got %d", len(results))
	}

	if len(results) > 0 && results[0].Stack != StackGo {
		t.Errorf("expected Go stack, got %v", results[0].Stack)
	}
}
