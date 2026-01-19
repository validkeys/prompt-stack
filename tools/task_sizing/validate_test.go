package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestLoadYAML tests loading a valid YAML file
func TestLoadYAML(t *testing.T) {
	// Create a temporary YAML file
	yamlContent := `name: test-project
description: Test project for validation
version: 1.0.0
task_sizing:
  min_minutes: 30
  max_minutes: 150
  max_files: 5
tasks:
  - id: task-1
    title: Test Task 1
    estimated_duration_minutes: 45
    files_in_scope:
      - file1.go
      - file2.go
  - id: task-2
    title: Test Task 2
    estimated_duration_minutes: 60
    dependencies:
      - task-1`

	tmpFile, err := os.CreateTemp("", "test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Test loading the YAML
	config, err := loadYAML(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load YAML: %v", err)
	}

	// Verify the loaded configuration
	if config.Name != "test-project" {
		t.Errorf("Expected name 'test-project', got %q", config.Name)
	}
	if config.TaskSizing.MinMinutes != 30 {
		t.Errorf("Expected min_minutes 30, got %d", config.TaskSizing.MinMinutes)
	}
	if config.TaskSizing.MaxMinutes != 150 {
		t.Errorf("Expected max_minutes 150, got %d", config.TaskSizing.MaxMinutes)
	}
	if len(config.Tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(config.Tasks))
	}
}

// TestLoadYAMLInvalid tests loading an invalid YAML file
func TestLoadYAMLInvalid(t *testing.T) {
	// Create a temporary invalid YAML file
	invalidYAML := `name: test-project
invalid: [`

	tmpFile, err := os.CreateTemp("", "test-invalid-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(invalidYAML)); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Test loading the invalid YAML
	_, err = loadYAML(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

// TestLoadYAMLNotFound tests loading a non-existent file
func TestLoadYAMLNotFound(t *testing.T) {
	_, err := loadYAML("/non/existent/file.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

// TestValidateTaskSizingValid tests validation of valid tasks
func TestValidateTaskSizingValid(t *testing.T) {
	config := &RalphyYAML{
		Name:        "test-project",
		Description: "Test project",
		Version:     "1.0.0",
		TaskSizing: TaskSizingConfig{
			MinMinutes: 30,
			MaxMinutes: 150,
			MaxFiles:   5,
		},
		Tasks: []Task{
			{
				ID:                       "task-1",
				Title:                    "Valid Task 1",
				EstimatedDurationMinutes: 45,
				FilesInScope:             []string{"file1.go", "file2.go"},
			},
			{
				ID:                       "task-2",
				Title:                    "Valid Task 2",
				EstimatedDurationMinutes: 90,
				Dependencies:             []string{"task-1"},
				FilesInScope:             []string{"file3.go"},
			},
		},
	}

	result := validateTaskSizing(config)

	if !result.Valid {
		t.Errorf("Expected valid result, got invalid. Violations: %v", result.Violations)
	}
	if result.TotalTasks != 2 {
		t.Errorf("Expected 2 total tasks, got %d", result.TotalTasks)
	}
	if result.Summary.TasksWithinRange != 2 {
		t.Errorf("Expected 2 tasks within range, got %d", result.Summary.TasksWithinRange)
	}
	if result.Summary.TasksOutsideRange != 0 {
		t.Errorf("Expected 0 tasks outside range, got %d", result.Summary.TasksOutsideRange)
	}
	if result.Summary.MinDuration != 45 {
		t.Errorf("Expected min duration 45, got %d", result.Summary.MinDuration)
	}
	if result.Summary.MaxDuration != 90 {
		t.Errorf("Expected max duration 90, got %d", result.Summary.MaxDuration)
	}
	if result.Summary.AverageDuration != 67.5 {
		t.Errorf("Expected average duration 67.5, got %f", result.Summary.AverageDuration)
	}
}

// TestValidateTaskSizingBelowMinimum tests tasks below minimum duration
func TestValidateTaskSizingBelowMinimum(t *testing.T) {
	config := &RalphyYAML{
		Name:        "test-project",
		Description: "Test project",
		Version:     "1.0.0",
		TaskSizing: TaskSizingConfig{
			MinMinutes: 30,
			MaxMinutes: 150,
			MaxFiles:   5,
		},
		Tasks: []Task{
			{
				ID:                       "task-1",
				Title:                    "Too Short Task",
				EstimatedDurationMinutes: 20, // Below minimum
				FilesInScope:             []string{"file1.go"},
			},
		},
	}

	result := validateTaskSizing(config)

	if result.Valid {
		t.Error("Expected invalid result for task below minimum, got valid")
	}
	if len(result.Violations) != 1 {
		t.Errorf("Expected 1 violation, got %d", len(result.Violations))
	}
	if result.Violations[0].Issue != "duration_below_minimum" {
		t.Errorf("Expected violation 'duration_below_minimum', got %q", result.Violations[0].Issue)
	}
	if result.Summary.TasksOutsideRange != 1 {
		t.Errorf("Expected 1 task outside range, got %d", result.Summary.TasksOutsideRange)
	}
}

// TestValidateTaskSizingAboveMaximum tests tasks above maximum duration
func TestValidateTaskSizingAboveMaximum(t *testing.T) {
	config := &RalphyYAML{
		Name:        "test-project",
		Description: "Test project",
		Version:     "1.0.0",
		TaskSizing: TaskSizingConfig{
			MinMinutes: 30,
			MaxMinutes: 150,
			MaxFiles:   5,
		},
		Tasks: []Task{
			{
				ID:                       "task-1",
				Title:                    "Too Long Task",
				EstimatedDurationMinutes: 200, // Above maximum
				FilesInScope:             []string{"file1.go"},
			},
		},
	}

	result := validateTaskSizing(config)

	if result.Valid {
		t.Error("Expected invalid result for task above maximum, got valid")
	}
	if len(result.Violations) != 1 {
		t.Errorf("Expected 1 violation, got %d", len(result.Violations))
	}
	if result.Violations[0].Issue != "duration_above_maximum" {
		t.Errorf("Expected violation 'duration_above_maximum', got %q", result.Violations[0].Issue)
	}
	if result.Summary.TasksOutsideRange != 1 {
		t.Errorf("Expected 1 task outside range, got %d", result.Summary.TasksOutsideRange)
	}
}

// TestValidateTaskSizingTooManyFiles tests tasks with too many files
func TestValidateTaskSizingTooManyFiles(t *testing.T) {
	config := &RalphyYAML{
		Name:        "test-project",
		Description: "Test project",
		Version:     "1.0.0",
		TaskSizing: TaskSizingConfig{
			MinMinutes: 30,
			MaxMinutes: 150,
			MaxFiles:   5,
		},
		Tasks: []Task{
			{
				ID:                       "task-1",
				Title:                    "Too Many Files Task",
				EstimatedDurationMinutes: 60,
				FilesInScope:             []string{"f1.go", "f2.go", "f3.go", "f4.go", "f5.go", "f6.go"}, // 6 files > max 5
			},
		},
	}

	result := validateTaskSizing(config)

	if result.Valid {
		t.Error("Expected invalid result for too many files, got valid")
	}
	if len(result.Violations) != 1 {
		t.Errorf("Expected 1 violation, got %d", len(result.Violations))
	}
	if result.Violations[0].Issue != "too_many_files" {
		t.Errorf("Expected violation 'too_many_files', got %q", result.Violations[0].Issue)
	}
}

// TestValidateTaskSizingMultipleViolations tests multiple violations in one task
func TestValidateTaskSizingMultipleViolations(t *testing.T) {
	config := &RalphyYAML{
		Name:        "test-project",
		Description: "Test project",
		Version:     "1.0.0",
		TaskSizing: TaskSizingConfig{
			MinMinutes: 30,
			MaxMinutes: 150,
			MaxFiles:   5,
		},
		Tasks: []Task{
			{
				ID:                       "task-1",
				Title:                    "Multiple Violations Task",
				EstimatedDurationMinutes: 10,                                                                      // Below minimum
				FilesInScope:             []string{"f1.go", "f2.go", "f3.go", "f4.go", "f5.go", "f6.go", "f7.go"}, // 7 files > max 5
			},
		},
	}

	result := validateTaskSizing(config)

	if result.Valid {
		t.Error("Expected invalid result for multiple violations, got valid")
	}
	if len(result.Violations) != 2 {
		t.Errorf("Expected 2 violations, got %d", len(result.Violations))
	}

	// Check that both violations are present
	durationViolation := false
	filesViolation := false
	for _, v := range result.Violations {
		if v.Issue == "duration_below_minimum" {
			durationViolation = true
		}
		if v.Issue == "too_many_files" {
			filesViolation = true
		}
	}

	if !durationViolation {
		t.Error("Expected duration_below_minimum violation")
	}
	if !filesViolation {
		t.Error("Expected too_many_files violation")
	}
}

// TestValidateTaskSizingDependencies tests dependency tracking
func TestValidateTaskSizingDependencies(t *testing.T) {
	config := &RalphyYAML{
		Name:        "test-project",
		Description: "Test project",
		Version:     "1.0.0",
		TaskSizing: TaskSizingConfig{
			MinMinutes: 30,
			MaxMinutes: 150,
			MaxFiles:   5,
		},
		Tasks: []Task{
			{
				ID:                       "task-1",
				Title:                    "Independent Task",
				EstimatedDurationMinutes: 45,
			},
			{
				ID:                       "task-2",
				Title:                    "Dependent Task",
				EstimatedDurationMinutes: 60,
				Dependencies:             []string{"task-1"},
			},
			{
				ID:                       "task-3",
				Title:                    "Another Independent Task",
				EstimatedDurationMinutes: 75,
			},
		},
	}

	result := validateTaskSizing(config)

	if !result.Valid {
		t.Errorf("Expected valid result, got invalid. Violations: %v", result.Violations)
	}
	if result.Summary.TasksWithDependencies != 1 {
		t.Errorf("Expected 1 task with dependencies, got %d", result.Summary.TasksWithDependencies)
	}
	if result.Summary.TasksWithoutDependencies != 2 {
		t.Errorf("Expected 2 tasks without dependencies, got %d", result.Summary.TasksWithoutDependencies)
	}
	if result.ParallelOpportunity.TotalIndependentTasks != 2 {
		t.Errorf("Expected 2 independent tasks, got %d", result.ParallelOpportunity.TotalIndependentTasks)
	}
	if len(result.ParallelOpportunity.IndependentTaskIDs) != 2 {
		t.Errorf("Expected 2 independent task IDs, got %d", len(result.ParallelOpportunity.IndependentTaskIDs))
	}
}

// TestIdentifyParallelGroups tests parallel group identification
func TestIdentifyParallelGroups(t *testing.T) {
	tasks := []Task{
		{ID: "task-1", Dependencies: []string{}},
		{ID: "task-2", Dependencies: []string{"task-1"}},
		{ID: "task-3", Dependencies: []string{}},
		{ID: "task-4", Dependencies: []string{"task-2"}},
	}

	dependencyMap := map[string][]string{
		"task-2": {"task-1"},
		"task-4": {"task-2"},
	}

	groups := identifyParallelGroups(tasks, dependencyMap)

	if len(groups) != 1 {
		t.Errorf("Expected 1 parallel group, got %d", len(groups))
	}
	if len(groups[0]) != 2 {
		t.Errorf("Expected 2 tasks in parallel group, got %d", len(groups[0]))
	}

	// Check that task-1 and task-3 are in the group
	task1Found := false
	task3Found := false
	for _, taskID := range groups[0] {
		if taskID == "task-1" {
			task1Found = true
		}
		if taskID == "task-3" {
			task3Found = true
		}
	}

	if !task1Found {
		t.Error("Expected task-1 in parallel group")
	}
	if !task3Found {
		t.Error("Expected task-3 in parallel group")
	}
}

// TestValidateTaskSizingEmptyTasks tests validation with empty tasks
func TestValidateTaskSizingEmptyTasks(t *testing.T) {
	config := &RalphyYAML{
		Name:        "test-project",
		Description: "Test project",
		Version:     "1.0.0",
		TaskSizing: TaskSizingConfig{
			MinMinutes: 30,
			MaxMinutes: 150,
			MaxFiles:   5,
		},
		Tasks: []Task{},
	}

	result := validateTaskSizing(config)

	if !result.Valid {
		t.Error("Expected valid result for empty tasks, got invalid")
	}
	if result.TotalTasks != 0 {
		t.Errorf("Expected 0 total tasks, got %d", result.TotalTasks)
	}
	if result.Summary.AverageDuration != 0 {
		t.Errorf("Expected average duration 0, got %f", result.Summary.AverageDuration)
	}
}

// TestValidateTaskSizingJSONOutput tests JSON output format
func TestValidateTaskSizingJSONOutput(t *testing.T) {
	config := &RalphyYAML{
		Name:        "test-project",
		Description: "Test project",
		Version:     "1.0.0",
		TaskSizing: TaskSizingConfig{
			MinMinutes: 30,
			MaxMinutes: 150,
			MaxFiles:   5,
		},
		Tasks: []Task{
			{
				ID:                       "task-1",
				Title:                    "Test Task",
				EstimatedDurationMinutes: 45,
			},
		},
	}

	result := validateTaskSizing(config)

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal result to JSON: %v", err)
	}

	// Unmarshal back to verify structure
	var decodedResult ValidationResult
	if err := json.Unmarshal(jsonData, &decodedResult); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify the decoded result matches the original
	if decodedResult.Valid != result.Valid {
		t.Errorf("JSON Valid mismatch: expected %v, got %v", result.Valid, decodedResult.Valid)
	}
	if decodedResult.TotalTasks != result.TotalTasks {
		t.Errorf("JSON TotalTasks mismatch: expected %d, got %d", result.TotalTasks, decodedResult.TotalTasks)
	}
}

// TestIntegration tests the full integration with a real YAML file
func TestIntegration(t *testing.T) {
	// Get the path to the test data directory
	testDataDir := filepath.Join("..", "..", "docs", "implementation-plan", "m0")
	yamlPath := filepath.Join(testDataDir, "ralphy_inputs.yaml")

	// Check if the file exists
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		t.Skip("Test YAML file not found, skipping integration test")
	}

	// Load and validate the YAML
	config, err := loadYAML(yamlPath)
	if err != nil {
		t.Fatalf("Failed to load YAML file %q: %v", yamlPath, err)
	}

	result := validateTaskSizing(config)

	// Basic validation of the result
	if result.TotalTasks == 0 {
		t.Error("Expected at least one task in the YAML file")
	}

	// The actual validation result depends on the content of the YAML file
	// We just verify that the validation completed without panicking
	t.Logf("Validation completed: %d tasks, valid=%v", result.TotalTasks, result.Valid)
	if len(result.Violations) > 0 {
		t.Logf("Violations found: %v", result.Violations)
	}
}
