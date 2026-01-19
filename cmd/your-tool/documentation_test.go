package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func TestDocumentationExists(t *testing.T) {
	root := getProjectRoot()
	if root == "" {
		t.Fatal("Could not find project root")
	}

	tests := []struct {
		name string
		path string
	}{
		{"README.md exists", "README.md"},
		{"CONTRIBUTING.md exists", "CONTRIBUTING.md"},
		{"architecture.md exists", "docs/architecture.md"},
		{"commands.md exists", "docs/commands.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath := filepath.Join(root, tt.path)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Errorf("File %q does not exist at %q", tt.path, fullPath)
			}
		})
	}
}

func TestREADMEContent(t *testing.T) {
	root := getProjectRoot()
	if root == "" {
		t.Fatal("Could not find project root")
	}

	content, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readme := string(content)

	requiredSections := []string{
		"# prompt-stack",
		"## Quickstart",
		"## Features",
		"## Usage",
		"## Project Structure",
		"## Development",
		"## Documentation",
		"## Contributing",
	}

	for _, section := range requiredSections {
		if !strings.Contains(readme, section) {
			t.Errorf("README.md missing required section: %q", section)
		}
	}

	requiredFeatures := []string{
		"Plan mode",
		"Build mode",
		"Validate",
		"Review",
		"Init",
	}

	for _, feature := range requiredFeatures {
		if !strings.Contains(readme, feature) {
			t.Errorf("README.md missing required feature description: %q", feature)
		}
	}
}

func TestREADMEUsageExamples(t *testing.T) {
	root := getProjectRoot()
	if root == "" {
		t.Fatal("Could not find project root")
	}

	content, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readme := string(content)

	requiredCommands := []string{
		"your-tool init",
		"your-tool plan",
		"your-tool validate",
		"your-tool build",
		"your-tool review",
	}

	for _, cmd := range requiredCommands {
		if !strings.Contains(readme, cmd) {
			t.Errorf("README.md missing usage example for command: %q", cmd)
		}
	}
}

func TestContributingContent(t *testing.T) {
	root := getProjectRoot()
	if root == "" {
		t.Fatal("Could not find project root")
	}

	content, err := os.ReadFile(filepath.Join(root, "CONTRIBUTING.md"))
	if err != nil {
		t.Fatalf("Failed to read CONTRIBUTING.md: %v", err)
	}

	contrib := string(content)

	requiredSections := []string{
		"# Contributing",
		"## Development Setup",
		"## Code Style",
		"## Testing",
		"## Commit Messages",
		"## Pull Request Process",
	}

	for _, section := range requiredSections {
		if !strings.Contains(contrib, section) {
			t.Errorf("CONTRIBUTING.md missing required section: %q", section)
		}
	}

	requiredPatterns := []string{
		"conventional commits",
		"table-driven",
		"Pull Request",
		"make test",
	}

	for _, pattern := range requiredPatterns {
		if !strings.Contains(contrib, pattern) {
			t.Errorf("CONTRIBUTING.md missing required pattern: %q", pattern)
		}
	}
}

func TestArchitectureContent(t *testing.T) {
	root := getProjectRoot()
	if root == "" {
		t.Fatal("Could not find project root")
	}

	content, err := os.ReadFile(filepath.Join(root, "docs/architecture.md"))
	if err != nil {
		t.Fatalf("Failed to read docs/architecture.md: %v", err)
	}

	arch := string(content)

	requiredSections := []string{
		"# Architecture",
		"## Overview",
		"## Components",
		"## Data Flow",
		"## Design Principles",
	}

	for _, section := range requiredSections {
		if !strings.Contains(arch, section) {
			t.Errorf("docs/architecture.md missing required section: %q", section)
		}
	}

	requiredComponents := []string{
		"CLI Layer",
		"Executor Package",
		"Prompt Package",
	}

	for _, component := range requiredComponents {
		if !strings.Contains(arch, component) {
			t.Errorf("docs/architecture.md missing required component description: %q", component)
		}
	}
}

func TestCommandsDocumentation(t *testing.T) {
	root := getProjectRoot()
	if root == "" {
		t.Fatal("Could not find project root")
	}

	content, err := os.ReadFile(filepath.Join(root, "docs/commands.md"))
	if err != nil {
		t.Fatalf("Failed to read docs/commands.md: %v", err)
	}

	cmds := string(content)

	requiredCommands := []string{
		"### init",
		"### plan",
		"### validate",
		"### build",
		"### review",
	}

	for _, cmd := range requiredCommands {
		if !strings.Contains(cmds, cmd) {
			t.Errorf("docs/commands.md missing documentation for command: %q", cmd)
		}
	}

	requiredPatterns := []string{
		"## Overview",
		"## Commands",
		"## Global Flags",
		"## Exit Codes",
	}

	for _, pattern := range requiredPatterns {
		if !strings.Contains(cmds, pattern) {
			t.Errorf("docs/commands.md missing required pattern: %q", pattern)
		}
	}
}
