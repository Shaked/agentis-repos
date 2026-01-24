package detector

import (
	"os"
	"path/filepath"
	"testing"
)

// createTempProject creates a temp directory with the specified files
func createTempProject(t *testing.T, files []string) string {
	t.Helper()
	dir := t.TempDir()
	for _, f := range files {
		path := filepath.Join(dir, f)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("failed to create directory for %s: %v", f, err)
		}
		if err := os.WriteFile(path, []byte{}, 0644); err != nil {
			t.Fatalf("failed to create file %s: %v", f, err)
		}
	}
	return dir
}

// createTempDir creates a temp directory with the specified subdirectories
func createTempDir(t *testing.T, dirs []string) string {
	t.Helper()
	root := t.TempDir()
	for _, d := range dirs {
		path := filepath.Join(root, d)
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("failed to create directory %s: %v", d, err)
		}
	}
	return root
}
