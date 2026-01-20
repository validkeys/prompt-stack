package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocumentationCompleteness(t *testing.T) {
	repoDir := GetRepoDir(t)

	tests := []struct {
		name        string
		filePath    string
		checkFunc   func(content string) bool
		description string
	}{
		{
			name:     "README_exists",
			filePath: "README.md",
			checkFunc: func(content string) bool {
				return len(content) > 0
			},
			description: "README.md exists and is not empty",
		},
		{
			name:     "README_describes_tool_purpose",
			filePath: "README.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "prompt-stack") &&
					(strings.Contains(content, "AI-assisted") || strings.Contains(content, "development"))
			},
			description: "README describes tool purpose and core features",
		},
		{
			name:     "README_usage_examples_init",
			filePath: "README.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "init") &&
					strings.Contains(content, "./dist/prompt-stack init")
			},
			description: "README has usage examples for init command",
		},
		{
			name:     "README_usage_examples_plan",
			filePath: "README.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "plan") &&
					strings.Contains(content, "./dist/prompt-stack plan")
			},
			description: "README has usage examples for plan command",
		},
		{
			name:     "README_usage_examples_validate",
			filePath: "README.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "validate") &&
					strings.Contains(content, "./dist/prompt-stack validate")
			},
			description: "README has usage examples for validate command",
		},
		{
			name:     "README_usage_examples_review",
			filePath: "README.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "review") &&
					strings.Contains(content, "./dist/prompt-stack review")
			},
			description: "README has usage examples for review command",
		},
		{
			name:     "README_usage_examples_build",
			filePath: "README.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "build") &&
					strings.Contains(content, "./dist/prompt-stack build")
			},
			description: "README has usage examples for build command",
		},
		{
			name:     "CONTRIBUTING_exists",
			filePath: "CONTRIBUTING.md",
			checkFunc: func(content string) bool {
				return len(content) > 0
			},
			description: "CONTRIBUTING.md exists and is not empty",
		},
		{
			name:     "CONTRIBUTING_development_setup",
			filePath: "CONTRIBUTING.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "Development") &&
					(strings.Contains(content, "clone") || strings.Contains(content, "Clone")) &&
					strings.Contains(content, "go mod")
			},
			description: "CONTRIBUTING.md includes development setup instructions",
		},
		{
			name:     "CONTRIBUTING_code_style",
			filePath: "CONTRIBUTING.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "Code") &&
					(strings.Contains(content, "Style") || strings.Contains(content, "style"))
			},
			description: "CONTRIBUTING.md includes code style guidelines",
		},
		{
			name:     "CONTRIBUTING_testing",
			filePath: "CONTRIBUTING.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "Test") &&
					strings.Contains(content, "make test")
			},
			description: "CONTRIBUTING.md includes testing guidelines",
		},
		{
			name:     "commands_md_exists",
			filePath: "docs/commands.md",
			checkFunc: func(content string) bool {
				return len(content) > 0
			},
			description: "docs/commands.md exists and is not empty",
		},
		{
			name:     "commands_md_init_documentation",
			filePath: "docs/commands.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "init") &&
					strings.Contains(content, "Initialize")
			},
			description: "docs/commands.md documents init command",
		},
		{
			name:     "commands_md_plan_documentation",
			filePath: "docs/commands.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "plan") &&
					strings.Contains(content, "Generate")
			},
			description: "docs/commands.md documents plan command",
		},
		{
			name:     "commands_md_validate_documentation",
			filePath: "docs/commands.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "validate") &&
					strings.Contains(content, "Validate")
			},
			description: "docs/commands.md documents validate command",
		},
		{
			name:     "commands_md_review_documentation",
			filePath: "docs/commands.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "review") &&
					strings.Contains(content, "Review")
			},
			description: "docs/commands.md documents review command",
		},
		{
			name:     "commands_md_build_documentation",
			filePath: "docs/commands.md",
			checkFunc: func(content string) bool {
				return strings.Contains(content, "build") &&
					strings.Contains(content, "Build")
			},
			description: "docs/commands.md documents build command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath := filepath.Join(repoDir, tt.filePath)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Fatalf("failed to read %s: %v", tt.filePath, err)
			}

			if !tt.checkFunc(string(content)) {
				t.Errorf("%s: check failed", tt.description)
			}
		})
	}
}
