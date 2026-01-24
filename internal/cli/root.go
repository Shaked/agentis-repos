package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agentic-repo",
	Short: "Initialize repositories with the Agent-Native Repository Standard",
	Long: `agentic-repo is a CLI tool that transforms standard Git repositories
into Hybrid Human/Agent Environments, optimizing for AI context windows,
reduced hallucinations, and strict code governance.

It generates structured context files that help AI agents understand
your codebase efficiently without wasting tokens.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("agentic-repo v0.1.0")
	},
}
