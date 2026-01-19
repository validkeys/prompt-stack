package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIsNegativePhrasing tests negative phrasing detection
func TestIsNegativePhrasing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		phrase   string
	}{
		{
			name:     "don't pattern",
			input:    "Don't use global variables",
			expected: true,
			phrase:   "(?i)\\bdon'?t\\b",
		},
		{
			name:     "do not pattern",
			input:    "Do not ignore errors",
			expected: true,
			phrase:   "(?i)\\bdo not\\b",
		},
		{
			name:     "never pattern",
			input:    "Never use magic numbers",
			expected: true,
			phrase:   "(?i)\\bnever\\b",
		},
		{
			name:     "avoid pattern",
			input:    "Avoid using any type",
			expected: true,
			phrase:   "(?i)\\bavoid\\b",
		},
		{
			name:     "prevent pattern",
			input:    "Prevent race conditions",
			expected: true,
			phrase:   "(?i)\\bprevent\\b",
		},
		{
			name:     "no using pattern",
			input:    "No using deprecated APIs",
			expected: true,
			phrase:   "(?i)\\bno\\s+[a-z]+ing\\b",
		},
		{
			name:     "not writing pattern",
			input:    "Not writing to stdout",
			expected: true,
			phrase:   "(?i)\\bnot\\s+[a-z]+ing\\b",
		},
		{
			name:     "affirmative phrasing",
			input:    "Always validate inputs",
			expected: false,
			phrase:   "",
		},
		{
			name:     "neutral phrasing",
			input:    "Use proper error handling",
			expected: false,
			phrase:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isNegative, phrase := isNegativePhrasing(tt.input)
			if isNegative != tt.expected {
				t.Errorf("isNegativePhrasing(%q) = %v, want %v", tt.input, isNegative, tt.expected)
			}
			if isNegative && phrase != tt.phrase {
				t.Errorf("isNegativePhrasing(%q) phrase = %q, want %q", tt.input, phrase, tt.phrase)
			}
		})
	}
}

// TestIsAffirmativePhrasing tests affirmative phrasing detection
func TestIsAffirmativePhrasing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "always pattern",
			input:    "Always use table-driven tests",
			expected: true,
		},
		{
			name:     "must pattern",
			input:    "Must validate external inputs",
			expected: true,
		},
		{
			name:     "should pattern",
			input:    "Should write documentation",
			expected: true,
		},
		{
			name:     "use pattern",
			input:    "Use Go formatting conventions",
			expected: true,
		},
		{
			name:     "follow pattern",
			input:    "Follow best practices",
			expected: true,
		},
		{
			name:     "write pattern",
			input:    "Write unit tests",
			expected: true,
		},
		{
			name:     "negative phrasing",
			input:    "Don't ignore errors",
			expected: false,
		},
		{
			name:     "neutral phrasing",
			input:    "Consider performance implications",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAffirmativePhrasing(tt.input)
			if result != tt.expected {
				t.Errorf("isAffirmativePhrasing(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsSpecificConstraint tests constraint specificity detection
func TestIsSpecificConstraint(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "specific with actionable verb",
			input:    "Always validate external inputs at entry points",
			expected: true,
		},
		{
			name:     "specific with use verb",
			input:    "Use Go formatting conventions",
			expected: true,
		},
		{
			name:     "vague with properly",
			input:    "Handle errors properly",
			expected: false,
		},
		{
			name:     "vague with correctly",
			input:    "Implement tests correctly",
			expected: false,
		},
		{
			name:     "vague with some",
			input:    "Add some validation",
			expected: false,
		},
		{
			name:     "specific with check verb",
			input:    "Check for nil pointers",
			expected: true,
		},
		{
			name:     "specific with ensure verb",
			input:    "Ensure code quality",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSpecificConstraint(tt.input)
			if result != tt.expected {
				t.Errorf("isSpecificConstraint(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestContainsPatternReference tests pattern reference detection
func TestContainsPatternReference(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "contains pattern",
			input:    "Follow the error handling pattern",
			expected: true,
		},
		{
			name:     "contains example",
			input:    "See the validation example",
			expected: true,
		},
		{
			name:     "contains reference",
			input:    "Reference the style guide",
			expected: true,
		},
		{
			name:     "contains anchor",
			input:    "Use style anchors",
			expected: true,
		},
		{
			name:     "contains best practice",
			input:    "Follow best practices",
			expected: true,
		},
		{
			name:     "no pattern reference",
			input:    "Always validate inputs",
			expected: false,
		},
		{
			name:     "contains convention",
			input:    "Use naming conventions",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsPatternReference(tt.input)
			if result != tt.expected {
				t.Errorf("containsPatternReference(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestValidateConstraints tests the main validation function
func TestValidateConstraints(t *testing.T) {
	tests := []struct {
		name       string
		config     *RalphyYAML
		expected   bool
		violations int
	}{
		{
			name: "all affirmative constraints",
			config: &RalphyYAML{
				GlobalConstraints: GlobalConstraints{
					AffirmativeConstraints: []string{
						"Always validate inputs",
						"Use table-driven tests",
						"Write documentation",
					},
				},
			},
			expected:   true,
			violations: 0,
		},
		{
			name: "mixed constraints with negative phrasing",
			config: &RalphyYAML{
				GlobalConstraints: GlobalConstraints{
					ForbiddenPatterns: []PatternConstraint{
						{
							Pattern: "\\bany\\b",
							Message: "Don't use any type",
						},
					},
					AffirmativeConstraints: []string{
						"Always validate inputs",
					},
				},
			},
			expected:   false,
			violations: 1,
		},
		{
			name: "multiple negative constraints",
			config: &RalphyYAML{
				GlobalConstraints: GlobalConstraints{
					ForbiddenPatterns: []PatternConstraint{
						{
							Pattern: "@ts-ignore",
							Message: "Never use ts-ignore",
						},
					},
					RequiredPatterns: []PatternConstraint{
						{
							Pattern: "func Test",
							Message: "Do not skip tests",
						},
					},
					AffirmativeConstraints: []string{
						"Don't ignore errors",
					},
				},
			},
			expected:   false,
			violations: 3,
		},
		{
			name: "vague constraints",
			config: &RalphyYAML{
				GlobalConstraints: GlobalConstraints{
					AffirmativeConstraints: []string{
						"Do things properly",
						"Be careful with inputs",
					},
				},
			},
			expected:   true, // Vague constraints don't fail validation, just get flagged
			violations: 2,    // Both are vague and get violation entries
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateConstraints(tt.config)
			if result.Valid != tt.expected {
				t.Errorf("validateConstraints() valid = %v, want %v", result.Valid, tt.expected)
			}
			if len(result.Violations) != tt.violations {
				t.Errorf("validateConstraints() violations = %d, want %d", len(result.Violations), tt.violations)
			}
		})
	}
}

// TestLoadYAML tests YAML loading
func TestLoadYAML(t *testing.T) {
	// Create a temporary YAML file
	yamlContent := `
name: test-project
description: Test project for validation
version: 1.0.0
global_constraints:
  forbidden_patterns:
    - pattern: "\\bany\\b"
      message: "Use unknown with type guards instead"
  affirmative_constraints:
    - "Always validate inputs"
`

	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	config, err := loadYAML(yamlPath)
	if err != nil {
		t.Fatalf("loadYAML() error = %v", err)
	}

	if config.Name != "test-project" {
		t.Errorf("loadYAML() name = %q, want %q", config.Name, "test-project")
	}
	if config.Description != "Test project for validation" {
		t.Errorf("loadYAML() description = %q, want %q", config.Description, "Test project for validation")
	}
	if config.Version != "1.0.0" {
		t.Errorf("loadYAML() version = %q, want %q", config.Version, "1.0.0")
	}
	if len(config.GlobalConstraints.ForbiddenPatterns) != 1 {
		t.Errorf("loadYAML() forbidden_patterns count = %d, want %d", len(config.GlobalConstraints.ForbiddenPatterns), 1)
	}
	if len(config.GlobalConstraints.AffirmativeConstraints) != 1 {
		t.Errorf("loadYAML() affirmative_constraints count = %d, want %d", len(config.GlobalConstraints.AffirmativeConstraints), 1)
	}
}

// TestValidateConstraintsIntegration tests integration with actual YAML
func TestValidateConstraintsIntegration(t *testing.T) {
	// Test with the actual ralphy_inputs.yaml
	yamlPath := "../../docs/implementation-plan/m0/ralphy_inputs.yaml"

	// Check if file exists
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		t.Skip("ralphy_inputs.yaml not found, skipping integration test")
	}

	config, err := loadYAML(yamlPath)
	if err != nil {
		t.Fatalf("loadYAML(%q) error = %v", yamlPath, err)
	}

	result := validateConstraints(config)

	// The actual file should have affirmative constraints
	if result.TotalConstraints == 0 {
		t.Error("validateConstraints() found 0 constraints in actual YAML file")
	}

	// Check that we have some affirmative constraints
	if result.AffirmativeCount == 0 && result.TotalConstraints > 0 {
		t.Error("validateConstraints() found 0 affirmative constraints in actual YAML file")
	}

	// Log the results for debugging
	t.Logf("Validation result: valid=%v, total=%d, affirmative=%d, negative=%d",
		result.Valid, result.TotalConstraints, result.AffirmativeCount, result.NegativeCount)
}
