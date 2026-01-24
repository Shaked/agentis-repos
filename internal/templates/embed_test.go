package templates

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name         string
		templateName string
		wantErr      bool
		contains     string // optional: check if content contains this string
	}{
		{
			name:         "loads agentignore template",
			templateName: "agentignore.tmpl",
			wantErr:      false,
		},
		{
			name:         "loads gitignore template",
			templateName: "gitignore.tmpl",
			wantErr:      false,
		},
		{
			name:         "loads architecture template",
			templateName: "architecture.md.tmpl",
			wantErr:      false,
			contains:     "Architecture",
		},
		{
			name:         "loads agents.md template",
			templateName: "agents.md.tmpl",
			wantErr:      false,
		},
		{
			name:         "loads cursorrules template",
			templateName: "cursorrules.tmpl",
			wantErr:      false,
		},
		{
			name:         "loads go stack template",
			templateName: "go/stack.md.tmpl",
			wantErr:      false,
			contains:     "Go",
		},
		{
			name:         "loads go testing template",
			templateName: "go/testing.md.tmpl",
			wantErr:      false,
		},
		{
			name:         "loads go commands template",
			templateName: "go/commands.md.tmpl",
			wantErr:      false,
		},
		{
			name:         "loads python stack template",
			templateName: "python/stack.md.tmpl",
			wantErr:      false,
			contains:     "Python",
		},
		{
			name:         "loads node stack template",
			templateName: "node/stack.md.tmpl",
			wantErr:      false,
		},
		{
			name:         "loads java stack template",
			templateName: "java/stack.md.tmpl",
			wantErr:      false,
			contains:     "Java",
		},
		{
			name:         "loads unknown stack template",
			templateName: "unknown/stack.md.tmpl",
			wantErr:      false,
		},
		{
			name:         "loads overview template",
			templateName: "overview.md.tmpl",
			wantErr:      false,
		},
		{
			name:         "loads monorepo agents template",
			templateName: "agents-monorepo.md.tmpl",
			wantErr:      false,
		},
		{
			name:         "returns error for nonexistent template",
			templateName: "nonexistent.tmpl",
			wantErr:      true,
		},
		{
			name:         "returns error for invalid path",
			templateName: "../../../etc/passwd",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := Get(tt.templateName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Get(%q) expected error, got nil", tt.templateName)
				}
				return
			}

			if err != nil {
				t.Errorf("Get(%q) unexpected error: %v", tt.templateName, err)
				return
			}

			if content == "" {
				t.Errorf("Get(%q) returned empty content", tt.templateName)
			}

			if tt.contains != "" && !strings.Contains(content, tt.contains) {
				t.Errorf("Get(%q) content does not contain %q", tt.templateName, tt.contains)
			}
		})
	}
}

func TestGet_NotFound(t *testing.T) {
	_, err := Get("this-template-does-not-exist.tmpl")
	if err == nil {
		t.Error("expected error for nonexistent template, got nil")
	}
}

func TestMustGet(t *testing.T) {
	tests := []struct {
		name         string
		templateName string
		shouldPanic  bool
	}{
		{
			name:         "returns content for valid template",
			templateName: "agentignore.tmpl",
			shouldPanic:  false,
		},
		{
			name:         "panics for nonexistent template",
			templateName: "nonexistent.tmpl",
			shouldPanic:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("MustGet() should have panicked but didn't")
					}
				}()
			}

			content := MustGet(tt.templateName)

			if !tt.shouldPanic && content == "" {
				t.Error("MustGet() returned empty content")
			}
		})
	}
}

func TestGet_AllStackTemplates(t *testing.T) {
	stacks := []string{"go", "python", "node", "java", "unknown"}
	templateTypes := []string{"stack.md.tmpl", "testing.md.tmpl", "commands.md.tmpl", "code-review-rules.md.tmpl", "Makefile.tmpl", "pre-commit-config.yaml.tmpl"}

	for _, stack := range stacks {
		for _, tmplType := range templateTypes {
			templateName := stack + "/" + tmplType
			t.Run(templateName, func(t *testing.T) {
				content, err := Get(templateName)
				if err != nil {
					t.Errorf("Get(%q) failed: %v", templateName, err)
					return
				}
				if content == "" {
					t.Errorf("Get(%q) returned empty content", templateName)
				}
			})
		}
	}
}
