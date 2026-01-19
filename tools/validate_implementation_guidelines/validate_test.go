package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRequiresTestFirstWorkflow tests the test-first workflow detection
func TestRequiresTestFirstWorkflow(t *testing.T) {
	tests := []struct {
		name     string
		task     Task
		expected bool
	}{
		{
			name: "task with test in description",
			task: Task{
				Title:       "Add unit tests",
				Description: "Create unit tests for the core package",
			},
			expected: true,
		},
		{
			name: "task with test in title",
			task: Task{
				Title:       "Test infrastructure setup",
				Description: "Set up testing framework",
			},
			expected: true,
		},
		{
			name: "task with test files in scope",
			task: Task{
				Title:        "Implement feature",
				Description:  "Add new functionality",
				FilesInScope: []string{"pkg/feature/feature.go", "pkg/feature/feature_test.go"},
			},
			expected: true,
		},
		{
			name: "task without test references",
			task: Task{
				Title:       "Documentation update",
				Description: "Update README with new features",
			},
			expected: false,
		},
		{
			name: "task with TDD in description",
			task: Task{
				Title:       "Implement validation",
				Description: "Use TDD approach to implement input validation",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := requiresTestFirstWorkflow(tt.task)
			if result != tt.expected {
				t.Errorf("requiresTestFirstWorkflow() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestHasTestableAcceptanceCriteria tests the testable acceptance criteria detection
func TestHasTestableAcceptanceCriteria(t *testing.T) {
	tests := []struct {
		name     string
		task     Task
		expected bool
	}{
		{
			name: "task with testable criteria",
			task: Task{
				AcceptanceCriteria: []string{
					"Function returns expected value",
					"Test passes with valid input",
				},
			},
			expected: true,
		},
		{
			name: "task with non-testable criteria",
			task: Task{
				AcceptanceCriteria: []string{
					"Code is well-structured",
					"Documentation is clear",
				},
			},
			expected: false,
		},
		{
			name: "task with mixed criteria",
			task: Task{
				AcceptanceCriteria: []string{
					"Code is well-structured",
					"Function outputs correct result",
					"Documentation is clear",
				},
			},
			expected: true,
		},
		{
			name: "task with multiple violations",
			task: Task{
				ID:          "test-001",
				Title:       "Add comprehensive tests",
				Description: "Implement thorough unit testing",
				AcceptanceCriteria: []string{
					"Implementation is complete",
					"All functionality works",
				},
				StyleAnchors: []StyleAnchor{
					{
						File:   "some/file.go",
						Reason: "File reference",
					},
				},
			},
			expected: false,
		},
		{
			name: "task with 'contains' criteria",
			task: Task{
				AcceptanceCriteria: []string{
					"Output contains expected message",
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasTestableAcceptanceCriteria(tt.task)
			if result != tt.expected {
				t.Errorf("hasTestableAcceptanceCriteria() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestHasImplementationGuidance tests the implementation guidance detection
func TestHasImplementationGuidance(t *testing.T) {
	tests := []struct {
		name     string
		task     Task
		expected bool
	}{
		{
			name: "task with style anchor containing guidance",
			task: Task{
				Description: "Implement feature",
				StyleAnchors: []StyleAnchor{
					{
						File:   "examples/style-anchor/pkg/greeter/greeter.go",
						Reason: "Reference implementation pattern for simple packages",
					},
				},
			},
			expected: true,
		},
		{
			name: "task with guidance in description",
			task: Task{
				Description:  "Follow best practices for error handling",
				StyleAnchors: []StyleAnchor{},
			},
			expected: true,
		},
		{
			name: "task with no guidance",
			task: Task{
				Description: "Do something",
				StyleAnchors: []StyleAnchor{
					{
						File:   "some/file.go",
						Reason: "Just a file",
					},
				},
			},
			expected: false,
		},
		{
			name: "task with pattern reference",
			task: Task{
				Description:  "Use the design pattern for CLI commands",
				StyleAnchors: []StyleAnchor{},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasImplementationGuidance(tt.task)
			if result != tt.expected {
				t.Errorf("hasImplementationGuidance() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestValidateImplementationGuidelines tests the complete validation function
func TestValidateImplementationGuidelines(t *testing.T) {
	tests := []struct {
		name     string
		config   *RalphyYAML
		expected bool
	}{
		{
			name: "valid configuration with all guidelines",
			config: &RalphyYAML{
				TDD: TDD{
					Required: false,
				},
				Tasks: []Task{
					{
						ID:          "test-001",
						Title:       "Add unit tests",
						Description: "Create unit tests for core package",
						AcceptanceCriteria: []string{
							"Tests pass with valid input",
							"Coverage exceeds 80%",
						},
						StyleAnchors: []StyleAnchor{
							{
								File:   "examples/style-anchor/pkg/greeter/greeter_test.go",
								Reason: "Reference test pattern for table-driven tests",
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "invalid configuration missing testable criteria",
			config: &RalphyYAML{
				TDD: TDD{
					Required: false,
				},
				Tasks: []Task{
					{
						ID:          "test-001",
						Title:       "Add unit tests",
						Description: "Create unit tests for core package",
						AcceptanceCriteria: []string{
							"Code quality is good",
							"Implementation is thorough",
						},
						StyleAnchors: []StyleAnchor{
							{
								File:   "examples/style-anchor/pkg/greeter/greeter_test.go",
								Reason: "Reference test pattern",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "invalid configuration missing implementation guidance",
			config: &RalphyYAML{
				TDD: TDD{
					Required: false,
				},
				Tasks: []Task{
					{
						ID:          "test-001",
						Title:       "Add unit tests",
						Description: "Create unit tests for core package",
						AcceptanceCriteria: []string{
							"Tests pass with valid input",
						},
						StyleAnchors: []StyleAnchor{
							{
								File:   "some/file.go",
								Reason: "Just a file",
							},
						},
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateImplementationGuidelines(tt.config)
			if result.Valid != tt.expected {
				t.Errorf("validateImplementationGuidelines().Valid = %v, expected %v", result.Valid, tt.expected)
			}
		})
	}
}

// TestValidationResultStructure tests that validation result has correct structure
func TestValidationResultStructure(t *testing.T) {
	config := &RalphyYAML{
		TDD: TDD{
			Required: false,
		},
		Tasks: []Task{
			{
				ID:          "test-001",
				Title:       "Test task",
				Description: "Task with unit tests",
				AcceptanceCriteria: []string{
					"Tests pass",
				},
				StyleAnchors: []StyleAnchor{
					{
						File:   "examples/test.go",
						Reason: "Test pattern reference",
					},
				},
			},
			{
				ID:          "test-002",
				Title:       "Documentation task",
				Description: "Update documentation",
				AcceptanceCriteria: []string{
					"Documentation is clear",
				},
				StyleAnchors: []StyleAnchor{},
			},
		},
	}

	result := validateImplementationGuidelines(config)

	// Check basic structure
	if result.TotalTasks != 2 {
		t.Errorf("TotalTasks = %d, expected 2", result.TotalTasks)
	}

	// First task needs tests, has testable criteria, has guidance
	if result.TasksNeedingTests != 1 {
		t.Errorf("TasksNeedingTests = %d, expected 1", result.TasksNeedingTests)
	}

	if result.TasksWithTestableCriteria != 1 {
		t.Errorf("TasksWithTestableCriteria = %d, expected 1", result.TasksWithTestableCriteria)
	}

	if result.TasksWithImplementationGuidance != 1 {
		t.Errorf("TasksWithImplementationGuidance = %d, expected 1", result.TasksWithImplementationGuidance)
	}

	// Check summary calculations
	if result.Summary.TestFirstCoverage != 50.0 {
		t.Errorf("TestFirstCoverage = %.1f, expected 50.0", result.Summary.TestFirstCoverage)
	}

	if result.Summary.TestableCriteriaCoverage != 50.0 {
		t.Errorf("TestableCriteriaCoverage = %.1f, expected 50.0", result.Summary.TestableCriteriaCoverage)
	}

	if result.Summary.ImplementationGuidanceCoverage != 50.0 {
		t.Errorf("ImplementationGuidanceCoverage = %.1f, expected 50.0", result.Summary.ImplementationGuidanceCoverage)
	}
}

// TestCommandLineInterface tests the command-line interface
func TestCommandLineInterface(t *testing.T) {
	// Create a temporary test YAML file
	tempDir := t.TempDir()
	testYAML := `name: test-project
description: Test project
version: 1.0.0
tdd:
  required: false
  test_command: go test ./...
  failure_instruction: Fix tests
tasks:
  - id: test-001
    title: Test task
    description: Task with unit tests
    completed: false
    files_in_scope: []
    style_anchors:
      - file: examples/test.go
        reason: Test pattern reference
    estimated_duration_minutes: 30
    single_responsibility: Test implementation
    acceptance_criteria:
      - Tests pass
    verification:
      pre_commit: []
`

	yamlPath := filepath.Join(tempDir, "test.yaml")
	if err := os.WriteFile(yamlPath, []byte(testYAML), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Test with valid file
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"validate", "--file", yamlPath}

	// We can't easily test the main function exit codes in a unit test
	// without complex mocking, so we'll just verify the file can be loaded
	config, err := loadYAML(yamlPath)
	if err != nil {
		t.Fatalf("Failed to load test YAML: %v", err)
	}

	if config.Name != "test-project" {
		t.Errorf("Config.Name = %s, expected test-project", config.Name)
	}

	if len(config.Tasks) != 1 {
		t.Errorf("Config.Tasks length = %d, expected 1", len(config.Tasks))
	}
}

// TestEmptyConfig tests validation with empty configuration
func TestEmptyConfig(t *testing.T) {
	config := &RalphyYAML{
		TDD:   TDD{Required: false},
		Tasks: []Task{},
	}

	result := validateImplementationGuidelines(config)

	if result.TotalTasks != 0 {
		t.Errorf("TotalTasks = %d, expected 0", result.TotalTasks)
	}

	if !result.Valid {
		t.Error("Empty config should be valid")
	}

	if len(result.Violations) != 0 {
		t.Errorf("Violations count = %d, expected 0", len(result.Violations))
	}
}

// TestTaskWithMultipleViolations tests detection of multiple violations in one task
func TestTaskWithMultipleViolations(t *testing.T) {
	config := &RalphyYAML{
		TDD: TDD{Required: false},
		Tasks: []Task{
			{
				ID:          "test-001",
				Title:       "Add comprehensive tests",
				Description: "Implement thorough unit testing",
				AcceptanceCriteria: []string{
					"Implementation is complete",
					"All functionality works",
				},
				StyleAnchors: []StyleAnchor{
					{
						File:   "some/file.go",
						Reason: "File reference",
					},
				},
			},
		},
	}

	result := validateImplementationGuidelines(config)

	if result.Valid {
		t.Error("Config with violations should not be valid")
	}

	if len(result.Violations) < 2 {
		t.Errorf("Expected at least 2 violations, got %d", len(result.Violations))
	}

	// Check for specific violation types
	hasMissingTestableCriteria := false
	hasMissingImplementationGuidance := false

	for _, violation := range result.Violations {
		if violation.Issue == "missing_testable_criteria" {
			hasMissingTestableCriteria = true
		}
		if violation.Issue == "missing_implementation_guidance" {
			hasMissingImplementationGuidance = true
		}
	}

	if !hasMissingTestableCriteria {
		t.Error("Missing testable criteria violation not found")
	}

	if !hasMissingImplementationGuidance {
		t.Error("Missing implementation guidance violation not found")
	}
}
