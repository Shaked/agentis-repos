package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Shaked/agentic-repo/internal/detector"
)

func TestInitCmd_Flags(t *testing.T) {
	// Verify init command has expected flags
	flags := initCmd.Flags()

	tests := []struct {
		name      string
		shorthand string
	}{
		{"force", "f"},
		{"dry-run", "n"},
		{"verbose", "v"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := flags.Lookup(tt.name)
			if flag == nil {
				t.Errorf("flag --%s not found", tt.name)
				return
			}
			if flag.Shorthand != tt.shorthand {
				t.Errorf("flag --%s shorthand = %q, want %q", tt.name, flag.Shorthand, tt.shorthand)
			}
		})
	}
}

func TestRunInit(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		setupDir  func(t *testing.T) string
		wantErr   bool
		wantFiles []string
	}{
		{
			name: "initializes current directory",
			args: []string{},
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				// Create a Go project indicator
				os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0644)
				return dir
			},
			wantErr: false,
			wantFiles: []string{
				"AGENTS.md",
				".agent/stack.md",
			},
		},
		{
			name: "initializes specified directory",
			args: []string{}, // Will be set to temp dir
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0644)
				return dir
			},
			wantErr: false,
			wantFiles: []string{
				"AGENTS.md",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)

			// Reset flags
			flagForce = false
			flagDryRun = false
			flagVerbose = false

			// Set args to target the temp directory
			args := append(tt.args, dir)

			err := runInit(initCmd, args)

			if (err != nil) != tt.wantErr {
				t.Errorf("runInit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				for _, wantFile := range tt.wantFiles {
					path := filepath.Join(dir, wantFile)
					if _, err := os.Stat(path); os.IsNotExist(err) {
						t.Errorf("expected file %s to exist", wantFile)
					}
				}
			}
		})
	}
}

func TestRunInit_NonExistentDirectory(t *testing.T) {
	err := runInit(initCmd, []string{"/nonexistent/path/that/does/not/exist"})
	if err == nil {
		t.Error("expected error for non-existent directory")
	}
}

func TestRunInit_NotADirectory(t *testing.T) {
	// Create a temporary file (not a directory)
	dir := t.TempDir()
	filePath := filepath.Join(dir, "file.txt")
	os.WriteFile(filePath, []byte("test"), 0644)

	err := runInit(initCmd, []string{filePath})
	if err == nil {
		t.Error("expected error for file path (not directory)")
	}
}

func TestRunInit_DryRun(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0644)

	// Set dry-run flag
	flagDryRun = true
	defer func() { flagDryRun = false }()

	err := runInit(initCmd, []string{dir})
	if err != nil {
		t.Fatalf("runInit() error = %v", err)
	}

	// In dry-run mode, only the go.mod we created should exist
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}

	// Should only have the go.mod file we created
	if len(entries) != 1 {
		t.Errorf("dry-run should not create files, found %d entries (expected 1 for go.mod)", len(entries))
	}
}

func TestRunInit_Verbose(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0644)

	// Set verbose flag
	flagVerbose = true
	defer func() { flagVerbose = false }()

	err := runInit(initCmd, []string{dir})
	if err != nil {
		t.Fatalf("runInit() error = %v", err)
	}

	// Just verify it runs without error in verbose mode
	// The actual verbose output goes to stdout
}

func TestPrintDetectionResults(t *testing.T) {
	tests := []struct {
		name       string
		results    []detector.Result
		isMonorepo bool
	}{
		{
			name: "prints single project",
			results: []detector.Result{
				{Path: "/project", Stack: detector.StackGo},
			},
			isMonorepo: false,
		},
		{
			name: "prints monorepo",
			results: []detector.Result{
				{Path: "/project/backend", Stack: detector.StackGo},
				{Path: "/project/frontend", Stack: detector.StackNode},
			},
			isMonorepo: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify it doesn't panic
			// The actual output goes to stdout
			printDetectionResults(tt.results, tt.isMonorepo)
		})
	}
}

func TestInitCmd_Help(t *testing.T) {
	// Verify init command is properly configured
	if initCmd.Use != "init [directory]" {
		t.Errorf("initCmd.Use = %q, want %q", initCmd.Use, "init [directory]")
	}

	if initCmd.Short == "" {
		t.Error("initCmd.Short should not be empty")
	}

	if initCmd.Long == "" {
		t.Error("initCmd.Long should not be empty")
	}

	// Verify help contains expected keywords
	if !strings.Contains(initCmd.Long, "Initialize") {
		t.Error("initCmd.Long should contain 'Initialize'")
	}
}

func TestRunInit_EmptyDirectory(t *testing.T) {
	// Test with an empty directory (no project indicators)
	dir := t.TempDir()

	err := runInit(initCmd, []string{dir})
	if err != nil {
		t.Fatalf("runInit() error = %v", err)
	}

	// Should generate files for unknown stack
	agentsPath := filepath.Join(dir, "AGENTS.md")
	if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
		t.Error("expected AGENTS.md to be created for empty directory")
	}
}
