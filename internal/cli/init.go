package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Shaked/agentic-repo/internal/detector"
	"github.com/Shaked/agentic-repo/internal/generator"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	flagForce   bool
	flagDryRun  bool
	flagVerbose bool
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a repository with agent context files",
	Long: `Initialize a repository with the Agent-Native Repository Standard.

This command will:
1. Detect your project type (Go, Python, Node/TS, Java)
2. Detect if it's a monorepo with multiple project types
3. Generate appropriate context files (AGENTS.md, .agent/, etc.)
4. Create integration stubs for AI tools (.cursorrules, .claude/)`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVarP(&flagForce, "force", "f", false, "Overwrite existing files")
	initCmd.Flags().BoolVarP(&flagDryRun, "dry-run", "n", false, "Preview changes without writing")
	initCmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "Show detailed output")
}

func runInit(cmd *cobra.Command, args []string) error {
	// Determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	// Resolve to absolute path
	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if directory exists
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("directory does not exist: %s", absPath)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", absPath)
	}

	// Print header
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Printf("🔍 Scanning %s\n", absPath)

	// Detect project stacks
	results, err := detector.Scan(absPath)
	if err != nil {
		return fmt.Errorf("detection failed: %w", err)
	}

	if len(results) == 0 {
		yellow := color.New(color.FgYellow)
		yellow.Println("⚠️  No recognized project types found")
		yellow.Println("   Generating generic context files...")
		results = []detector.Result{{
			Path:  absPath,
			Stack: detector.StackUnknown,
		}}
	}

	// Determine if monorepo
	isMonorepo := detector.IsMonorepo(results)

	if flagVerbose {
		printDetectionResults(results, isMonorepo)
	}

	// Generate files
	gen := generator.New(generator.Options{
		Force:   flagForce,
		DryRun:  flagDryRun,
		Verbose: flagVerbose,
	})

	if err := gen.Generate(absPath, results, isMonorepo); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	// Print success
	green := color.New(color.FgGreen, color.Bold)
	if flagDryRun {
		green.Println("\n✓ Dry run complete (no files written)")
	} else {
		green.Println("\n✓ Repository initialized with Agent-Native Standard")
	}

	return nil
}

func printDetectionResults(results []detector.Result, isMonorepo bool) {
	fmt.Println()
	if isMonorepo {
		color.New(color.FgMagenta).Println("📦 Detected: Monorepo")
	} else {
		color.New(color.FgMagenta).Println("📦 Detected: Single project")
	}

	for _, r := range results {
		fmt.Printf("   • %s: %s\n", r.Path, r.Stack)
	}
	fmt.Println()
}
