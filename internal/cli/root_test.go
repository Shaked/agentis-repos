package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	// Test that Execute runs without panic
	// We can't easily test the full execution since it depends on os.Args
	// but we can verify the root command is properly configured

	if rootCmd == nil {
		t.Error("rootCmd is nil")
	}

	if rootCmd.Use != "agentic-repo" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "agentic-repo")
	}
}

func TestVersionCmd(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	versionCmd.SetOut(buf)

	// Execute version command
	versionCmd.Run(versionCmd, []string{})

	output := buf.String()
	if !strings.Contains(output, "agentic-repo") {
		t.Errorf("version output should contain 'agentic-repo', got %q", output)
	}

	if !strings.Contains(output, "v0.1.0") {
		t.Errorf("version output should contain version number, got %q", output)
	}
}

func TestRootCmd_HasSubcommands(t *testing.T) {
	// Verify subcommands are registered
	subcommands := rootCmd.Commands()

	expectedCommands := []string{"init", "version"}
	for _, expected := range expectedCommands {
		found := false
		for _, cmd := range subcommands {
			if cmd.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand %q not found", expected)
		}
	}
}

func TestRootCmd_Help(t *testing.T) {
	// Verify root command is properly configured
	if rootCmd.Use != "agentic-repo" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "agentic-repo")
	}

	if rootCmd.Short == "" {
		t.Error("rootCmd.Short should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}

	// Verify description contains expected keywords
	expectedStrings := []string{
		"agentic-repo",
		"AI",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(rootCmd.Long, expected) {
			t.Errorf("rootCmd.Long should contain %q", expected)
		}
	}
}
