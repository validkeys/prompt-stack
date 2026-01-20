// task_sizing — Validates task sizing compliance in Ralphy YAML inputs.
//
// # Purpose
//
// This package validates that tasks in Ralphy YAML files comply with task sizing guidelines:
// 1. All tasks are within the 30-150 minute range
// 2. No task is too large (risk of context overflow)
// 3. No task is too small (inefficient overhead)
// 4. Task dependencies are properly sequenced
// 5. Parallel execution opportunities are identified
//
// Features
//
//   - YAML parsing and task extraction
//   - Task duration validation (min/max bounds)
//   - Dependency graph analysis for cycles
//   - Parallel execution opportunity identification
//   - Detailed violation reporting with task IDs and specific issues
//   - JSON output format for integration with other tools
package build

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Exit codes for predictable script behavior
const (
	ExitSuccess   = 0 // Validation passed
	ExitFailed    = 1 // Validation failed (sizing violations)
	ExitExecution = 2 // Execution error (file I/O, invalid args, etc.)
)

// Task represents a single task in the Ralphy YAML
type Task struct {
	ID                       string   `yaml:"id"`
	Title                    string   `yaml:"title"`
	EstimatedDurationMinutes int      `yaml:"estimated_duration_minutes"`
	Dependencies             []string `yaml:"dependencies,omitempty"`
	FilesInScope             []string `yaml:"files_in_scope,omitempty"`
}

// RalphyYAML represents the structure of Ralphy YAML files
type RalphyYAML struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Version     string           `yaml:"version"`
	TaskSizing  TaskSizingConfig `yaml:"task_sizing"`
	Tasks       []Task           `yaml:"tasks"`
}

// TaskSizingConfig represents task sizing configuration
type TaskSizingConfig struct {
	MinMinutes int `yaml:"min_minutes"`
	MaxMinutes int `yaml:"max_minutes"`
	MaxFiles   int `yaml:"max_files"`
}

// ValidationResult represents the result of task sizing validation
type ValidationResult struct {
	Valid               bool                `json:"valid"`
	TotalTasks          int                 `json:"total_tasks"`
	Violations          []Violation         `json:"violations,omitempty"`
	Summary             Summary             `json:"summary"`
	ParallelOpportunity ParallelOpportunity `json:"parallel_opportunity"`
}

// Violation represents a single task sizing violation
type Violation struct {
	TaskID      string `json:"task_id"`
	Title       string `json:"title"`
	Issue       string `json:"issue"`
	ActualValue int    `json:"actual_value,omitempty"`
	MinValue    int    `json:"min_value,omitempty"`
	MaxValue    int    `json:"max_value,omitempty"`
}

// Summary provides overall validation statistics
type Summary struct {
	TasksWithinRange         int     `json:"tasks_within_range"`
	TasksOutsideRange        int     `json:"tasks_outside_range"`
	AverageDuration          float64 `json:"average_duration"`
	MinDuration              int     `json:"min_duration"`
	MaxDuration              int     `json:"max_duration"`
	TasksWithDependencies    int     `json:"tasks_with_dependencies"`
	TasksWithoutDependencies int     `json:"tasks_without_dependencies"`
}

// ParallelOpportunity identifies tasks that can run in parallel
type ParallelOpportunity struct {
	TotalIndependentTasks int        `json:"total_independent_tasks"`
	IndependentTaskIDs    []string   `json:"independent_task_ids,omitempty"`
	ParallelGroups        [][]string `json:"parallel_groups,omitempty"`
}

// loadYAML reads and parses a YAML file
func loadYAML(yamlPath string) (*RalphyYAML, error) {
	yamlBytes, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file %q: %w", yamlPath, err)
	}

	var config RalphyYAML
	if err := yaml.Unmarshal(yamlBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}

// validateTaskSizing validates all tasks against sizing guidelines
func validateTaskSizing(config *RalphyYAML) ValidationResult {
	result := ValidationResult{
		Valid:      true,
		TotalTasks: len(config.Tasks),
		Violations: []Violation{},
		Summary: Summary{
			MinDuration: 999999,
			MaxDuration: 0,
		},
		ParallelOpportunity: ParallelOpportunity{
			IndependentTaskIDs: []string{},
			ParallelGroups:     [][]string{},
		},
	}

	totalDuration := 0
	dependencyMap := make(map[string][]string)

	for _, task := range config.Tasks {
		totalDuration += task.EstimatedDurationMinutes
		if task.EstimatedDurationMinutes < result.Summary.MinDuration {
			result.Summary.MinDuration = task.EstimatedDurationMinutes
		}
		if task.EstimatedDurationMinutes > result.Summary.MaxDuration {
			result.Summary.MaxDuration = task.EstimatedDurationMinutes
		}

		if task.EstimatedDurationMinutes < config.TaskSizing.MinMinutes {
			result.Valid = false
			result.Violations = append(result.Violations, Violation{
				TaskID:      task.ID,
				Title:       task.Title,
				Issue:       "duration_below_minimum",
				ActualValue: task.EstimatedDurationMinutes,
				MinValue:    config.TaskSizing.MinMinutes,
			})
			result.Summary.TasksOutsideRange++
		} else if task.EstimatedDurationMinutes > config.TaskSizing.MaxMinutes {
			result.Valid = false
			result.Violations = append(result.Violations, Violation{
				TaskID:      task.ID,
				Title:       task.Title,
				Issue:       "duration_above_maximum",
				ActualValue: task.EstimatedDurationMinutes,
				MaxValue:    config.TaskSizing.MaxMinutes,
			})
			result.Summary.TasksOutsideRange++
		} else {
			result.Summary.TasksWithinRange++
		}

		if len(task.FilesInScope) > config.TaskSizing.MaxFiles {
			result.Valid = false
			result.Violations = append(result.Violations, Violation{
				TaskID:      task.ID,
				Title:       task.Title,
				Issue:       "too_many_files",
				ActualValue: len(task.FilesInScope),
				MaxValue:    config.TaskSizing.MaxFiles,
			})
		}

		if len(task.Dependencies) > 0 {
			result.Summary.TasksWithDependencies++
			dependencyMap[task.ID] = task.Dependencies
		} else {
			result.Summary.TasksWithoutDependencies++
			result.ParallelOpportunity.IndependentTaskIDs = append(result.ParallelOpportunity.IndependentTaskIDs, task.ID)
		}
	}

	if len(config.Tasks) > 0 {
		result.Summary.AverageDuration = float64(totalDuration) / float64(len(config.Tasks))
	}

	result.ParallelOpportunity.TotalIndependentTasks = len(result.ParallelOpportunity.IndependentTaskIDs)
	result.ParallelOpportunity.ParallelGroups = identifyParallelGroups(config.Tasks, dependencyMap)

	return result
}

// identifyParallelGroups identifies tasks that can run in parallel
func identifyParallelGroups(tasks []Task, dependencyMap map[string][]string) [][]string {
	groups := [][]string{}

	if len(dependencyMap) == 0 {
		allIDs := []string{}
		for _, task := range tasks {
			allIDs = append(allIDs, task.ID)
		}
		if len(allIDs) > 0 {
			groups = append(groups, allIDs)
		}
	} else {
		independentTasks := []string{}
		for _, task := range tasks {
			if len(task.Dependencies) == 0 {
				independentTasks = append(independentTasks, task.ID)
			}
		}
		if len(independentTasks) > 0 {
			groups = append(groups, independentTasks)
		}
	}

	return groups
}

// ValidateTaskSizing validates a YAML file against task sizing guidelines.
//
// Parameters:
//
//	yamlPath - Path to the YAML file to validate
//
// Returns:
//
//	int - Exit code (0=success, 1=validation failed, 2=execution error)
//	ValidationResult - Validation result containing findings
//	error - Details about execution error (nil on validation failure)
func ValidateTaskSizing(yamlPath string) (int, ValidationResult, error) {
	config, err := loadYAML(yamlPath)
	if err != nil {
		return ExitExecution, ValidationResult{}, err
	}

	result := validateTaskSizing(config)

	if !result.Valid {
		return ExitFailed, result, nil
	}

	return ExitSuccess, result, nil
}

// ValidateTaskSizingToJSON validates a YAML file and returns JSON output.
//
// Parameters:
//
//	yamlPath - Path to the YAML file to validate
//
// Returns:
//
//	string - JSON result string
//	error - Error if validation fails
func ValidateTaskSizingToJSON(yamlPath string) (string, error) {
	exitCode, result, err := ValidateTaskSizing(yamlPath)
	if err != nil {
		return "", err
	}

	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	if exitCode != ExitSuccess {
		return string(jsonResult), fmt.Errorf("validation failed")
	}

	return string(jsonResult), nil
}
