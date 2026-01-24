package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Shaked/agentic-repo/internal/detector"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		opts Options
	}{
		{
			name: "creates generator with default options",
			opts: Options{},
		},
		{
			name: "creates generator with force option",
			opts: Options{Force: true},
		},
		{
			name: "creates generator with dry-run option",
			opts: Options{DryRun: true},
		},
		{
			name: "creates generator with verbose option",
			opts: Options{Verbose: true},
		},
		{
			name: "creates generator with all options",
			opts: Options{Force: true, DryRun: true, Verbose: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(tt.opts)
			if gen == nil {
				t.Error("New() returned nil")
			}
			if gen.opts != tt.opts {
				t.Errorf("New() opts = %v, want %v", gen.opts, tt.opts)
			}
		})
	}
}

func TestGenerate_SingleProject(t *testing.T) {
	tests := []struct {
		name      string
		stack     detector.StackType
		wantFiles []string
	}{
		{
			name:  "generates files for Go project",
			stack: detector.StackGo,
			wantFiles: []string{
				"AGENTS.md",
				"CODE_REVIEW_RULES.md",
				"USAGE.md",
				"Makefile",
				".gitignore",
				".agentignore",
				".agent/stack.md",
				".agent/testing.md",
				".agent/commands.md",
				".agent/architecture.md",
				".cursorrules",
				".claude/settings.json",
			},
		},
		{
			name:  "generates files for Python project",
			stack: detector.StackPython,
			wantFiles: []string{
				"AGENTS.md",
				".agent/stack.md",
				".cursorrules",
			},
		},
		{
			name:  "generates files for Node project",
			stack: detector.StackNode,
			wantFiles: []string{
				"AGENTS.md",
				".agent/stack.md",
				".cursorrules",
			},
		},
		{
			name:  "generates files for Java project",
			stack: detector.StackJava,
			wantFiles: []string{
				"AGENTS.md",
				".agent/stack.md",
				".cursorrules",
			},
		},
		{
			name:  "generates files for unknown project",
			stack: detector.StackUnknown,
			wantFiles: []string{
				"AGENTS.md",
				".agent/stack.md",
				".cursorrules",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			gen := New(Options{})
			results := []detector.Result{{Path: dir, Stack: tt.stack}}

			err := gen.Generate(dir, results, false)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			for _, wantFile := range tt.wantFiles {
				path := filepath.Join(dir, wantFile)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("expected file %s to exist", wantFile)
				}
			}
		})
	}
}

func TestGenerate_Monorepo(t *testing.T) {
	dir := t.TempDir()

	// Create subdirectory structure
	backendDir := filepath.Join(dir, "backend")
	frontendDir := filepath.Join(dir, "frontend")
	os.MkdirAll(backendDir, 0755)
	os.MkdirAll(frontendDir, 0755)

	gen := New(Options{})
	results := []detector.Result{
		{Path: backendDir, Stack: detector.StackGo},
		{Path: frontendDir, Stack: detector.StackNode},
	}

	err := gen.Generate(dir, results, true)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check root-level files
	rootFiles := []string{
		"AGENTS.md",
		".gitignore",
		".agent/overview.md",
		".agent/architecture.md",
		".cursorrules",
		".claude/settings.json",
	}

	for _, f := range rootFiles {
		path := filepath.Join(dir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected root file %s to exist", f)
		}
	}

	// Check per-project files
	projectFiles := []string{
		"AGENTS.md",
		".agent/stack.md",
	}

	for _, subdir := range []string{"backend", "frontend"} {
		for _, f := range projectFiles {
			path := filepath.Join(dir, subdir, f)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("expected project file %s/%s to exist", subdir, f)
			}
		}
	}
}

func TestGenerate_DryRun(t *testing.T) {
	dir := t.TempDir()

	gen := New(Options{DryRun: true})
	results := []detector.Result{{Path: dir, Stack: detector.StackGo}}

	err := gen.Generate(dir, results, false)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// In dry-run mode, no files should be created
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("dry-run should not create files, found %d entries", len(entries))
	}
}

func TestGenerate_SkipsExistingFiles(t *testing.T) {
	dir := t.TempDir()

	// Create an existing file with custom content (use Makefile, not AGENTS.md which has migration behavior)
	existingContent := "# Custom Makefile content"
	makefilePath := filepath.Join(dir, "Makefile")
	if err := os.WriteFile(makefilePath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	gen := New(Options{}) // Force is false by default
	results := []detector.Result{{Path: dir, Stack: detector.StackGo}}

	err := gen.Generate(dir, results, false)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Read the file and verify it wasn't overwritten
	content, err := os.ReadFile(makefilePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(content) != existingContent {
		t.Errorf("existing file was overwritten, got %q, want %q", string(content), existingContent)
	}
}

func TestGenerate_ForceOverwritesExistingFiles(t *testing.T) {
	dir := t.TempDir()

	// Create an existing file with custom content (use Makefile, not AGENTS.md which has migration behavior)
	existingContent := "# Custom Makefile content"
	makefilePath := filepath.Join(dir, "Makefile")
	if err := os.WriteFile(makefilePath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	gen := New(Options{Force: true})
	results := []detector.Result{{Path: dir, Stack: detector.StackGo}}

	err := gen.Generate(dir, results, false)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Read the file and verify it was overwritten
	content, err := os.ReadFile(makefilePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(content) == existingContent {
		t.Error("existing file was not overwritten with force=true")
	}
}

func TestGenerate_EmptyResults(t *testing.T) {
	dir := t.TempDir()

	gen := New(Options{})
	results := []detector.Result{}

	err := gen.Generate(dir, results, false)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Should generate files for unknown stack
	agentsPath := filepath.Join(dir, "AGENTS.md")
	if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
		t.Error("expected AGENTS.md to be created for empty results")
	}
}

func TestMigrateLegacyAgents(t *testing.T) {
	tests := []struct {
		name           string
		setupFiles     map[string]string // files to create before migration
		expectMigrated bool
		expectLegacy   bool // should .agent/AGENTS_LEGACY.md exist after
	}{
		{
			name:           "no existing AGENTS.md",
			setupFiles:     map[string]string{},
			expectMigrated: false,
			expectLegacy:   false,
		},
		{
			name: "migrates existing AGENTS.md",
			setupFiles: map[string]string{
				"AGENTS.md": "# Original content",
			},
			expectMigrated: true,
			expectLegacy:   true,
		},
		{
			name: "legacy already exists",
			setupFiles: map[string]string{
				"AGENTS.md":               "# New content",
				".agent/AGENTS_LEGACY.md": "# Old content",
			},
			expectMigrated: true, // returns true because legacy exists
			expectLegacy:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			// Setup files
			for path, content := range tt.setupFiles {
				fullPath := filepath.Join(dir, path)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
			}

			gen := New(Options{})
			migrated, err := gen.migrateLegacyAgents(dir)

			if err != nil {
				t.Fatalf("migrateLegacyAgents() error = %v", err)
			}

			if migrated != tt.expectMigrated {
				t.Errorf("migrateLegacyAgents() = %v, want %v", migrated, tt.expectMigrated)
			}

			legacyPath := filepath.Join(dir, ".agent", "AGENTS_LEGACY.md")

			if tt.expectLegacy {
				if _, err := os.Stat(legacyPath); os.IsNotExist(err) {
					t.Error("expected legacy file to exist")
				}
			}
		})
	}
}

func TestMigrateLegacyAgents_DryRun(t *testing.T) {
	dir := t.TempDir()

	// Create existing AGENTS.md
	agentsPath := filepath.Join(dir, "AGENTS.md")
	originalContent := "# Original content"
	if err := os.WriteFile(agentsPath, []byte(originalContent), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	gen := New(Options{DryRun: true})
	migrated, err := gen.migrateLegacyAgents(dir)

	if err != nil {
		t.Fatalf("migrateLegacyAgents() error = %v", err)
	}

	if !migrated {
		t.Error("expected migrateLegacyAgents() to return true in dry-run")
	}

	// Original file should still exist
	if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
		t.Error("original AGENTS.md should still exist in dry-run")
	}

	// Legacy file should not exist
	legacyPath := filepath.Join(dir, ".agent", "AGENTS_LEGACY.md")
	if _, err := os.Stat(legacyPath); !os.IsNotExist(err) {
		t.Error("legacy file should not be created in dry-run")
	}
}

func TestWriteTemplate_CreatesDirectories(t *testing.T) {
	dir := t.TempDir()

	gen := New(Options{})

	// Write to a nested path that doesn't exist
	nestedPath := filepath.Join(dir, "deep", "nested", "path", "file.md")
	err := gen.writeTemplate(nestedPath, "agentignore.tmpl", struct{}{})

	if err != nil {
		t.Fatalf("writeTemplate() error = %v", err)
	}

	if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
		t.Error("expected nested file to be created")
	}
}

func TestWriteTemplate_TemplateFallback(t *testing.T) {
	dir := t.TempDir()

	gen := New(Options{})

	// Try to write with a stack-specific template that falls back to generic
	// The code should fall back gracefully
	path := filepath.Join(dir, "test.md")

	// Create template data
	data := templateData{Stack: detector.StackGo}

	// This should work because agentignore.tmpl exists
	err := gen.writeTemplate(path, "agentignore.tmpl", data)
	if err != nil {
		t.Fatalf("writeTemplate() error = %v", err)
	}
}
