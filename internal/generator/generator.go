package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Shaked/agentic-repo/internal/detector"
	"github.com/Shaked/agentic-repo/internal/templates"
	"github.com/fatih/color"
)

// Options configures the generator behavior
type Options struct {
	Force   bool
	DryRun  bool
	Verbose bool
}

// Generator creates agent context files
type Generator struct {
	opts Options
}

// New creates a new Generator with the given options
func New(opts Options) *Generator {
	return &Generator{opts: opts}
}

// Generate creates all necessary files for the detected stacks
func (g *Generator) Generate(root string, results []detector.Result, isMonorepo bool) error {
	if isMonorepo {
		return g.generateMonorepo(root, results)
	}

	// Single project - use first result or unknown
	stack := detector.StackUnknown
	if len(results) > 0 {
		stack = results[0].Stack
	}

	return g.generateSingleProject(root, stack)
}

// generateSingleProject generates files for a single-stack project
func (g *Generator) generateSingleProject(root string, stack detector.StackType) error {
	// Migrate existing AGENTS.md to legacy location
	hasLegacy, err := g.migrateLegacyAgents(root)
	if err != nil {
		return fmt.Errorf("failed to migrate legacy AGENTS.md: %w", err)
	}

	// Template data with legacy flag
	data := templateData{Stack: stack, IsMonorepo: false, HasLegacy: hasLegacy}

	// Generate root files
	files := []struct {
		path     string
		template string
		data     any
	}{
		{"AGENTS.md", "agents.md.tmpl", data},
		{"CODE_REVIEW_RULES.md", fmt.Sprintf("%s/code-review-rules.md.tmpl", stack), data},
		{"repo-best-practices.md", "repo-best-practices.md.tmpl", data},
		{"USAGE.md", "usage.md.tmpl", data},
		{"Makefile", fmt.Sprintf("%s/Makefile.tmpl", stack), data},
		{".gitignore", "gitignore.tmpl", data},
		{".agentignore", "agentignore.tmpl", data},
		{".pre-commit-config.yaml", fmt.Sprintf("%s/pre-commit-config.yaml.tmpl", stack), data},
		{".agent/stack.md", fmt.Sprintf("%s/stack.md.tmpl", stack), data},
		{".agent/testing.md", fmt.Sprintf("%s/testing.md.tmpl", stack), data},
		{".agent/commands.md", fmt.Sprintf("%s/commands.md.tmpl", stack), data},
		{".agent/architecture.md", "architecture.md.tmpl", data},
		{".cursorrules", "cursorrules.tmpl", data},
		{".claude/settings.json", "claude-settings.json.tmpl", data},
	}

	for _, f := range files {
		fullPath := filepath.Join(root, f.path)
		if err := g.writeTemplate(fullPath, f.template, f.data); err != nil {
			return fmt.Errorf("failed to write %s: %w", f.path, err)
		}
	}

	return nil
}

// generateMonorepo generates files for a monorepo structure
func (g *Generator) generateMonorepo(root string, results []detector.Result) error {
	// Migrate existing AGENTS.md to legacy location
	hasLegacy, err := g.migrateLegacyAgents(root)
	if err != nil {
		return fmt.Errorf("failed to migrate legacy AGENTS.md: %w", err)
	}

	// Monorepo data with legacy flag
	monoData := monorepoData{Results: results, Root: root, HasLegacy: hasLegacy}

	// Generate root-level files
	rootFiles := []struct {
		path     string
		template string
		data     any
	}{
		{"AGENTS.md", "agents-monorepo.md.tmpl", monoData},
		{"CODE_REVIEW_RULES.md", "code-review-rules.md.tmpl", monoData},
		{"repo-best-practices.md", "repo-best-practices.md.tmpl", monoData},
		{"USAGE.md", "usage-monorepo.md.tmpl", monoData},
		{"Makefile", "Makefile-monorepo.tmpl", monoData},
		{".gitignore", "gitignore.tmpl", templateData{Stack: detector.StackUnknown}},
		{".agentignore", "agentignore.tmpl", templateData{Stack: detector.StackUnknown}},
		{".agent/overview.md", "overview.md.tmpl", monoData},
		{".agent/architecture.md", "architecture.md.tmpl", monoData},
		{".cursorrules", "cursorrules-monorepo.tmpl", monoData},
		{".claude/settings.json", "claude-settings-monorepo.json.tmpl", monoData},
	}

	for _, f := range rootFiles {
		fullPath := filepath.Join(root, f.path)
		if err := g.writeTemplate(fullPath, f.template, f.data); err != nil {
			return fmt.Errorf("failed to write %s: %w", f.path, err)
		}
	}

	// Generate per-project files
	for _, result := range results {
		if result.Path == root {
			continue // Skip root, already handled
		}

		// Check for legacy in subproject
		subHasLegacy, err := g.migrateLegacyAgents(result.Path)
		if err != nil {
			return fmt.Errorf("failed to migrate legacy AGENTS.md in %s: %w", result.Path, err)
		}

		relPath, _ := filepath.Rel(root, result.Path)
		subData := templateData{Stack: result.Stack, IsMonorepo: true, RelPath: relPath, HasLegacy: subHasLegacy}

		projectFiles := []struct {
			path     string
			template string
			data     any
		}{
			{"AGENTS.md", "agents.md.tmpl", subData},
			{"CODE_REVIEW_RULES.md", fmt.Sprintf("%s/code-review-rules.md.tmpl", result.Stack), subData},
			{"repo-best-practices.md", "repo-best-practices.md.tmpl", subData},
			{"USAGE.md", "usage.md.tmpl", subData},
			{".pre-commit-config.yaml", fmt.Sprintf("%s/pre-commit-config.yaml.tmpl", result.Stack), subData},
			{".agent/stack.md", fmt.Sprintf("%s/stack.md.tmpl", result.Stack), subData},
			{".agent/testing.md", fmt.Sprintf("%s/testing.md.tmpl", result.Stack), subData},
			{".agent/commands.md", fmt.Sprintf("%s/commands.md.tmpl", result.Stack), subData},
		}

		for _, f := range projectFiles {
			fullPath := filepath.Join(result.Path, f.path)
			if err := g.writeTemplate(fullPath, f.template, f.data); err != nil {
				return fmt.Errorf("failed to write %s: %w", fullPath, err)
			}
		}
	}

	return nil
}

// writeTemplate renders a template and writes it to disk
func (g *Generator) writeTemplate(path, tmplName string, data any) error {
	// Check if file exists
	if !g.opts.Force {
		if _, err := os.Stat(path); err == nil {
			if g.opts.Verbose {
				color.Yellow("   ⏭  Skipping %s (exists)", path)
			}
			return nil
		}
	}

	// Get template content
	content, err := templates.Get(tmplName)
	if err != nil {
		// Fall back to generic template if stack-specific doesn't exist
		genericName := filepath.Base(tmplName)
		content, err = templates.Get(genericName)
		if err != nil {
			return fmt.Errorf("template not found: %s", tmplName)
		}
	}

	// Parse and execute template
	tmpl, err := template.New(tmplName).Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output
	var output string
	if g.opts.DryRun {
		// In dry-run mode, just print what would be created
		color.Cyan("   📄 Would create: %s", path)
		return nil
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Execute template
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	color.Green("   ✓ Created: %s", path)
	_ = output // silence unused warning

	return nil
}

// templateData holds data for single-project templates
type templateData struct {
	Stack      detector.StackType
	IsMonorepo bool
	RelPath    string
	HasLegacy  bool
}

// monorepoData holds data for monorepo templates
type monorepoData struct {
	Results   []detector.Result
	Root      string
	HasLegacy bool
}

// migrateLegacyAgents moves existing AGENTS.md to .agent/AGENTS_LEGACY.md
// Returns true if migration occurred, false otherwise
func (g *Generator) migrateLegacyAgents(root string) (bool, error) {
	agentsPath := filepath.Join(root, "AGENTS.md")
	legacyPath := filepath.Join(root, ".agent", "AGENTS_LEGACY.md")

	// Check if AGENTS.md exists
	if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
		return false, nil
	}

	// Check if legacy already exists (don't overwrite)
	if _, err := os.Stat(legacyPath); err == nil {
		if g.opts.Verbose {
			color.Yellow("   ⏭  Legacy file already exists: %s", legacyPath)
		}
		return true, nil // Legacy exists, consider it migrated
	}

	if g.opts.DryRun {
		color.Cyan("   📦 Would migrate: AGENTS.md → .agent/AGENTS_LEGACY.md")
		return true, nil
	}

	// Ensure .agent directory exists
	agentDir := filepath.Join(root, ".agent")
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		return false, fmt.Errorf("failed to create .agent directory: %w", err)
	}

	// Read existing AGENTS.md
	content, err := os.ReadFile(agentsPath)
	if err != nil {
		return false, fmt.Errorf("failed to read existing AGENTS.md: %w", err)
	}

	// Write to legacy location
	if err := os.WriteFile(legacyPath, content, 0644); err != nil {
		return false, fmt.Errorf("failed to write legacy file: %w", err)
	}

	// Remove original AGENTS.md
	if err := os.Remove(agentsPath); err != nil {
		return false, fmt.Errorf("failed to remove original AGENTS.md: %w", err)
	}

	color.Magenta("   📦 Migrated: AGENTS.md → .agent/AGENTS_LEGACY.md")
	return true, nil
}
