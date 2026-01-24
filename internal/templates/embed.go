package templates

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed files/*
var templateFS embed.FS

// Get retrieves a template by name
func Get(name string) (string, error) {
	// Try direct path first
	content, err := fs.ReadFile(templateFS, "files/"+name)
	if err == nil {
		return string(content), nil
	}

	// Try without stack prefix (for common templates)
	return "", fmt.Errorf("template not found: %s", name)
}

// MustGet retrieves a template or panics
func MustGet(name string) string {
	content, err := Get(name)
	if err != nil {
		panic(err)
	}
	return content
}
