// Package validate_implementation_guidelines provides validation for implementation guidelines inclusion in Ralphy YAML inputs.
//
// # Purpose
//
// This tool validates that implementation-phase guidelines are present where applicable:
// 1. Identify tasks that will require test-first workflows at implementation time
// 2. Ensure task descriptions include 'testable' acceptance criteria when appropriate
// 3. Reference implementation guidance and style anchors for the implementation team
//
// Note: This planning-phase template does NOT generate tests or TDD artifacts; it only records guidance for later implementation phases.
//
// Usage
//
//	go run validate.go --file <input.yaml>
//
//	go run validate.go -f final_ralphy_inputs.yaml
//
// Exit Codes
//
//	0 - Validation successful (implementation guidelines present where needed)
//	1 - Validation failed (missing implementation guidelines)
//	2 - Error in command execution (file not found, invalid YAML, etc.)
//
// Features
//
//   - YAML parsing and task analysis
//   - Test-first workflow identification
//   - Acceptance criteria testability validation
//   - Implementation guidance reference validation
//   - Detailed violation reporting with specific issues
//   - JSON output format for integration with other tools
//
// Integration Points
//
//   - CI/CD pipelines: Run as a quality gate before merge
//   - Pre-commit hooks: Validate implementation guidelines before commits
//   - Build Mode: Validate generated Ralphy YAML before execution
//   - Local development: Quick validation during development
//
// Dependencies
//
//   - gopkg.in/yaml.v3: YAML parsing library
//
// Performance Characteristics
//
//   - Memory: O(file_size) - entire file loaded into memory
//   - Time: O(n) where n = tasks
//   - Typical runtime: <50ms for 100 tasks on modern hardware
//
// Limitations
//
//   - Does not validate YAML syntax (relies on yaml.v3 parser)
//   - Focuses on implementation guidelines only
//   - Does not validate other task properties
//
// Security Considerations
//
//   - Reads only local files specified via command-line arguments
//   - Does not execute external commands or evaluate arbitrary code
//   - No network access or external dependencies during validation
//
// Example Workflows
//
//	CI Pipeline:
//	  - name: Validate implementation guidelines
//	    run: |
//	      cd tools/validate_implementation_guidelines
//	      go run validate.go \
//	        --file ../../final_ralphy_inputs.yaml
//
//	Makefile Target:
//	  validate-implementation-guidelines:
//	      cd tools/validate_implementation_guidelines && go run validate.go \
//	        --file $(FILE)
//
//	Pre-commit Hook:
//	  if grep -q "ralphy_inputs.yaml" .git/hooks/pre-commit.sample; then
//	    tools/validate_implementation_guidelines/validate.go \
//	      --file ralphy_inputs.yaml
//	  fi
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Exit codes for predictable script behavior
const (
	ExitSuccess   = 0 // Validation passed
	ExitFailed    = 1 // Validation failed (guideline violations)
	ExitExecution = 2 // Execution error (file I/O, invalid args, etc.)
)

// Keywords that indicate test-first workflows
var testFirstKeywords = []string{
	"test", "testing", "unit test", "integration test", "tdd",
	"table-driven", "test coverage", "coverage", "assert",
	"verify", "validate", "check", "ensure",
}

// Keywords that indicate testable acceptance criteria
var testableCriteriaKeywords = []string{
	"passes", "fails", "returns", "outputs", "matches",
	"contains", "equals", "greater than", "less than",
	"successfully", "completes", "executes", "runs",
	"valid", "invalid", "correct", "incorrect",
	"pass", "test", "coverage", "expected",
}

// Keywords that indicate implementation guidance
var implementationGuidanceKeywords = []string{
	"guideline", "guide", "pattern", "example",
	"anchor", "style", "template", "specification", "standard",
	"best practice", "convention", "approach", "method",
	"technique", "strategy", "implementation", "development",
	"reference implementation", "coding standard", "design pattern",
}

// RalphyYAML represents the structure of Ralphy YAML files
type RalphyYAML struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	TDD         TDD    `yaml:"tdd"`
	Tasks       []Task `yaml:"tasks"`
}

// TDD represents the TDD configuration
type TDD struct {
	Required           bool   `yaml:"required"`
	TestCommand        string `yaml:"test_command"`
	FailureInstruction string `yaml:"failure_instruction"`
}

// Task represents a single task in the Ralphy YAML
type Task struct {
	ID                       string        `yaml:"id"`
	Title                    string        `yaml:"title"`
	Description              string        `yaml:"description"`
	Completed                bool          `yaml:"completed"`
	FilesInScope             []string      `yaml:"files_in_scope"`
	StyleAnchors             []StyleAnchor `yaml:"style_anchors"`
	Dependencies             []string      `yaml:"dependencies,omitempty"`
	EstimatedDurationMinutes int           `yaml:"estimated_duration_minutes"`
	EstimatedContextTokens   int           `yaml:"estimated_context_tokens,omitempty"`
	SingleResponsibility     string        `yaml:"single_responsibility"`
	AcceptanceCriteria       []string      `yaml:"acceptance_criteria"`
	Verification             Verification  `yaml:"verification"`
}

// StyleAnchor represents a style anchor reference
type StyleAnchor struct {
	File   string `yaml:"file"`
	Reason string `yaml:"reason"`
}

// Verification represents verification steps for a task
type Verification struct {
	PreCommit []string `yaml:"pre_commit"`
}

// ValidationResult represents the result of implementation guidelines validation
type ValidationResult struct {
	Valid                           bool        `json:"valid"`
	TotalTasks                      int         `json:"total_tasks"`
	TasksNeedingTests               int         `json:"tasks_needing_tests"`
	TasksWithTestableCriteria       int         `json:"tasks_with_testable_criteria"`
	TasksWithImplementationGuidance int         `json:"tasks_with_implementation_guidance"`
	Violations                      []Violation `json:"violations,omitempty"`
	Summary                         Summary     `json:"summary"`
	Recommendations                 []string    `json:"recommendations,omitempty"`
}

// Violation represents a single guideline violation
type Violation struct {
	TaskID     string `json:"task_id"`
	TaskTitle  string `json:"task_title"`
	Issue      string `json:"issue"`
	Details    string `json:"details"`
	Suggestion string `json:"suggestion,omitempty"`
}

// Summary provides overall validation statistics
type Summary struct {
	TestFirstCoverage              float64 `json:"test_first_coverage"`
	TestableCriteriaCoverage       float64 `json:"testable_criteria_coverage"`
	ImplementationGuidanceCoverage float64 `json:"implementation_guidance_coverage"`
	OverallScore                   float64 `json:"overall_score"`
}

// loadYAML reads and parses a YAML file
func loadYAML(yamlPath string) (*RalphyYAML, error) {
	yamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file %q: %w", yamlPath, err)
	}

	var config RalphyYAML
	if err := yaml.Unmarshal(yamlBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}

// requiresTestFirstWorkflow checks if a task requires test-first workflow
func requiresTestFirstWorkflow(task Task) bool {
	// Check task description for test-related keywords
	lowerDesc := strings.ToLower(task.Description)
	for _, keyword := range testFirstKeywords {
		if strings.Contains(lowerDesc, keyword) {
			return true
		}
	}

	// Check title for test-related keywords
	lowerTitle := strings.ToLower(task.Title)
	for _, keyword := range testFirstKeywords {
		if strings.Contains(lowerTitle, keyword) {
			return true
		}
	}

	// Check files in scope for test files
	for _, file := range task.FilesInScope {
		if strings.Contains(file, "_test.go") || strings.Contains(file, "_test.") {
			return true
		}
	}

	return false
}

// hasTestableAcceptanceCriteria checks if task has testable acceptance criteria
func hasTestableAcceptanceCriteria(task Task) bool {
	if len(task.AcceptanceCriteria) == 0 {
		return false
	}

	for _, criteria := range task.AcceptanceCriteria {
		lowerCriteria := strings.ToLower(criteria)
		for _, keyword := range testableCriteriaKeywords {
			if strings.Contains(lowerCriteria, keyword) {
				return true
			}
		}
	}

	return false
}

// hasImplementationGuidance checks if task references implementation guidance
func hasImplementationGuidance(task Task) bool {
	// Check style anchors for implementation guidance
	for _, anchor := range task.StyleAnchors {
		lowerReason := strings.ToLower(anchor.Reason)
		for _, keyword := range implementationGuidanceKeywords {
			if strings.Contains(lowerReason, keyword) {
				return true
			}
		}
	}

	// Check description for implementation guidance
	lowerDesc := strings.ToLower(task.Description)
	for _, keyword := range implementationGuidanceKeywords {
		if strings.Contains(lowerDesc, keyword) {
			return true
		}
	}

	return false
}

// validateImplementationGuidelines validates all tasks for implementation guidelines
func validateImplementationGuidelines(config *RalphyYAML) ValidationResult {
	result := ValidationResult{
		Valid:           true,
		Violations:      []Violation{},
		Summary:         Summary{},
		Recommendations: []string{},
	}

	result.TotalTasks = len(config.Tasks)

	for _, task := range config.Tasks {
		needsTests := requiresTestFirstWorkflow(task)
		hasTestableCriteria := hasTestableAcceptanceCriteria(task)
		hasGuidance := hasImplementationGuidance(task)

		if needsTests {
			result.TasksNeedingTests++
		}

		if hasTestableCriteria {
			result.TasksWithTestableCriteria++
		}

		if hasGuidance {
			result.TasksWithImplementationGuidance++
		}

		// Check for violations
		if needsTests && !hasTestableCriteria {
			result.Valid = false
			result.Violations = append(result.Violations, Violation{
				TaskID:     task.ID,
				TaskTitle:  task.Title,
				Issue:      "missing_testable_criteria",
				Details:    "Task appears to require test-first workflow but acceptance criteria are not testable",
				Suggestion: "Add testable acceptance criteria (e.g., 'test passes', 'function returns expected value', 'output matches pattern')",
			})
		}

		if needsTests && !hasGuidance {
			result.Valid = false
			result.Violations = append(result.Violations, Violation{
				TaskID:     task.ID,
				TaskTitle:  task.Title,
				Issue:      "missing_implementation_guidance",
				Details:    "Task requires test-first workflow but lacks implementation guidance references",
				Suggestion: "Add style anchors or references to testing patterns and implementation examples",
			})
		}
	}

	// Calculate coverage percentages
	if result.TotalTasks > 0 {
		result.Summary.TestFirstCoverage = float64(result.TasksNeedingTests) / float64(result.TotalTasks) * 100
		result.Summary.TestableCriteriaCoverage = float64(result.TasksWithTestableCriteria) / float64(result.TotalTasks) * 100
		result.Summary.ImplementationGuidanceCoverage = float64(result.TasksWithImplementationGuidance) / float64(result.TotalTasks) * 100

		// Calculate overall score (weighted average)
		result.Summary.OverallScore = (result.Summary.TestFirstCoverage*0.4 + result.Summary.TestableCriteriaCoverage*0.3 + result.Summary.ImplementationGuidanceCoverage*0.3) / 100
	}

	// Generate recommendations
	if result.Summary.TestFirstCoverage < 50 {
		result.Recommendations = append(result.Recommendations,
			fmt.Sprintf("Increase test-first workflow coverage from %.1f%% to at least 50%%", result.Summary.TestFirstCoverage))
	}

	if result.Summary.TestableCriteriaCoverage < 80 {
		result.Recommendations = append(result.Recommendations,
			fmt.Sprintf("Increase testable acceptance criteria coverage from %.1f%% to at least 80%%", result.Summary.TestableCriteriaCoverage))
	}

	if result.Summary.ImplementationGuidanceCoverage < 70 {
		result.Recommendations = append(result.Recommendations,
			fmt.Sprintf("Increase implementation guidance coverage from %.1f%% to at least 70%%", result.Summary.ImplementationGuidanceCoverage))
	}

	if result.Summary.OverallScore < 0.7 {
		result.Recommendations = append(result.Recommendations,
			fmt.Sprintf("Improve overall implementation guidelines score from %.2f to at least 0.7", result.Summary.OverallScore))
	}

	return result
}

// parseFlags parses and validates command-line arguments
func parseFlags() (string, error) {
	yamlPath := flag.String("file", "final_ralphy_inputs.yaml", "Path to YAML file to validate")
	flag.Parse()

	if *yamlPath == "" {
		return "", fmt.Errorf("--file flag is required")
	}

	return *yamlPath, nil
}

// main is the entry point for the validate_implementation_guidelines tool
func main() {
	yamlPath, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		fmt.Fprintf(os.Stderr, "Usage: go run validate.go --file <input.yaml>\n")
		os.Exit(ExitExecution)
	}

	config, err := loadYAML(yamlPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ExitExecution)
	}

	result := validateImplementationGuidelines(config)

	// Output JSON result
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal result: %v\n", err)
		os.Exit(ExitExecution)
	}

	fmt.Println(string(jsonResult))

	if !result.Valid {
		os.Exit(ExitFailed)
	}
	os.Exit(ExitSuccess)
}
