// Package validate_enforcement provides validation for multi-layer enforcement and commit/scope policies in Ralphy YAML inputs.
//
// # Purpose
//
// This tool validates that Ralphy YAML files include comprehensive multi-layer enforcement:
// 1. Prompt-level constraints in YAML
// 2. IDE/LSP integration considerations
// 3. Pre-commit hook specifications
// 4. CI check definitions
// 5. Runtime validation where applicable
// 6. Commit/revert and file-scope enforcement (atomic commits, files_in_scope enforcement)
//
// Multi-layer enforcement is critical for preventing architectural drift.
//
// Usage
//
//	go run validate.go --file <input.yaml>
//
//	go run validate.go -f final_ralphy_inputs.yaml
//
// Exit Codes
//
//	0 - Validation successful (≥3 verification layers, all requirements met)
//	1 - Validation failed (insufficient verification layers or missing requirements)
//	2 - Error in command execution (file not found, invalid YAML, etc.)
//
// Features
//
//   - YAML parsing and enforcement layer extraction
//   - Verification layer counting and validation
//   - Pre-commit hook specification validation
//   - CI check definition validation
//   - Runtime validation requirement checking
//   - Commit/scope policy validation
//   - Files_in_scope enforcement validation
//   - Detailed violation reporting with specific issues
//   - JSON output format for integration with other tools
//
// Integration Points
//
//   - CI/CD pipelines: Run as a quality gate before merge
//   - Pre-commit hooks: Validate enforcement layers before commits
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
//   - Time: O(n) where n = tasks + constraints
//   - Typical runtime: <100ms for 100 tasks on modern hardware
//
// Limitations
//
//   - Does not validate YAML syntax (relies on yaml.v3 parser)
//   - Focuses on enforcement layer validation only
//   - Does not validate other YAML properties
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
//	  - name: Validate multi-layer enforcement
//	    run: |
//	      cd tools/validate_enforcement
//	      go run validate.go \
//	        --file ../../final_ralphy_inputs.yaml
//
//	Makefile Target:
//	  validate-enforcement:
//	      cd tools/validate_enforcement && go run validate.go \
//	        --file $(FILE)
//
//	Pre-commit Hook:
//	  if grep -q "ralphy_inputs.yaml" .git/hooks/pre-commit.sample; then
//	    tools/validate_enforcement/validate.go \
//	      --file ralphy_inputs.yaml
//	  fi
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// Exit codes for predictable script behavior
const (
	ExitSuccess   = 0 // Validation passed
	ExitFailed    = 1 // Validation failed (enforcement violations)
	ExitExecution = 2 // Execution error (file I/O, invalid args, etc.)
)

// Minimum required verification layers
const minVerificationLayers = 3

// RalphyYAML represents the structure of Ralphy YAML files
type RalphyYAML struct {
	Name              string            `yaml:"name"`
	Description       string            `yaml:"description"`
	Version           string            `yaml:"version"`
	RulesFile         string            `yaml:"rules_file"`
	CI                CI                `yaml:"ci,omitempty"`
	DriftPolicyRef    string            `yaml:"drift_policy_ref,omitempty"`
	ValidationSchemas []string          `yaml:"validation_schemas,omitempty"`
	Outputs           Outputs           `yaml:"outputs"`
	GlobalConstraints GlobalConstraints `yaml:"global_constraints"`
	Tasks             []Task            `yaml:"tasks"`
}

// CI represents CI/CD configuration
type CI struct {
	Precommit []string `yaml:"precommit,omitempty"`
	CIChecks  []string `yaml:"ci_checks,omitempty"`
}

// Outputs represents output configuration
type Outputs struct {
	AllowedFileEdits    []string     `yaml:"allowed_file_edits,omitempty"`
	DisallowedFileEdits []string     `yaml:"disallowed_file_edits,omitempty"`
	CommitPolicy        CommitPolicy `yaml:"commit_policy,omitempty"`
}

// CommitPolicy represents commit policy configuration
type CommitPolicy struct {
	PrefixRules                []string `yaml:"prefix_rules,omitempty"`
	RequireScope               bool     `yaml:"require_scope,omitempty"`
	RequireConventionalCommits bool     `yaml:"require_conventional_commits,omitempty"`
}

// GlobalConstraints represents the constraints section
type GlobalConstraints struct {
	ForbiddenPatterns      []PatternConstraint `yaml:"forbidden_patterns,omitempty"`
	RequiredPatterns       []PatternConstraint `yaml:"required_patterns,omitempty"`
	AffirmativeConstraints []string            `yaml:"affirmative_constraints,omitempty"`
}

// PatternConstraint represents a pattern-based constraint
type PatternConstraint struct {
	Pattern string `yaml:"pattern"`
	Message string `yaml:"message"`
	When    string `yaml:"when,omitempty"`
}

// Task represents a single task definition
type Task struct {
	ID                   string       `yaml:"id"`
	Title                string       `yaml:"title"`
	Description          string       `yaml:"description"`
	FilesInScope         []string     `yaml:"files_in_scope,omitempty"`
	Verification         Verification `yaml:"verification,omitempty"`
	SingleResponsibility string       `yaml:"single_responsibility,omitempty"`
}

// Verification represents task verification configuration
type Verification struct {
	PreCommit  []string `yaml:"pre_commit,omitempty"`
	PostCommit []string `yaml:"post_commit,omitempty"`
	Runtime    []string `yaml:"runtime,omitempty"`
}

// ValidationResult represents the result of enforcement validation
type ValidationResult struct {
	Valid                 bool               `json:"valid"`
	TotalTasks            int                `json:"total_tasks"`
	TasksWithFilesInScope int                `json:"tasks_with_files_in_scope"`
	TasksWithVerification int                `json:"tasks_with_verification"`
	VerificationLayers    VerificationLayers `json:"verification_layers"`
	CommitPolicy          CommitPolicyStatus `json:"commit_policy"`
	ScopeEnforcement      ScopeEnforcement   `json:"scope_enforcement"`
	Violations            []Violation        `json:"violations,omitempty"`
	Recommendations       []string           `json:"recommendations,omitempty"`
}

// VerificationLayers counts verification layers
type VerificationLayers struct {
	PromptLevel    bool `json:"prompt_level"`
	IDEIntegration bool `json:"ide_integration"`
	PreCommit      bool `json:"pre_commit"`
	CIChecks       bool `json:"ci_checks"`
	Runtime        bool `json:"runtime"`
	TotalLayers    int  `json:"total_layers"`
}

// CommitPolicyStatus tracks commit policy compliance
type CommitPolicyStatus struct {
	HasPrefixRules         bool `json:"has_prefix_rules"`
	HasScopeRequirement    bool `json:"has_scope_requirement"`
	HasConventionalCommits bool `json:"has_conventional_commits"`
	Complete               bool `json:"complete"`
}

// ScopeEnforcement tracks scope enforcement
type ScopeEnforcement struct {
	HasAllowedFileEdits      bool `json:"has_allowed_file_edits"`
	HasDisallowedFileEdits   bool `json:"has_disallowed_file_edits"`
	AllTasksHaveFilesInScope bool `json:"all_tasks_have_files_in_scope"`
	Complete                 bool `json:"complete"`
}

// Violation represents a single enforcement violation
type Violation struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	TaskID      string `json:"task_id,omitempty"`
	Suggestion  string `json:"suggestion,omitempty"`
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

// validateEnforcement validates multi-layer enforcement
func validateEnforcement(config *RalphyYAML) ValidationResult {
	result := ValidationResult{
		Valid:              true,
		TotalTasks:         len(config.Tasks),
		Violations:         []Violation{},
		Recommendations:    []string{},
		VerificationLayers: VerificationLayers{},
		CommitPolicy:       CommitPolicyStatus{},
		ScopeEnforcement:   ScopeEnforcement{},
	}

	// Check verification layers
	result.VerificationLayers = checkVerificationLayers(config)

	// Check commit policy
	result.CommitPolicy = checkCommitPolicy(config)

	// Check scope enforcement
	result.ScopeEnforcement = checkScopeEnforcement(config)

	// Check tasks
	result = checkTasks(config, result)

	// Validate overall requirements
	result = validateRequirements(config, result)

	return result
}

// checkVerificationLayers checks for verification layers
func checkVerificationLayers(config *RalphyYAML) VerificationLayers {
	layers := VerificationLayers{}

	// 1. Prompt-level constraints
	if len(config.GlobalConstraints.ForbiddenPatterns) > 0 ||
		len(config.GlobalConstraints.RequiredPatterns) > 0 ||
		len(config.GlobalConstraints.AffirmativeConstraints) > 0 {
		layers.PromptLevel = true
	}

	// 2. IDE/LSP integration considerations
	// Check for rules_file and validation_schemas
	if config.RulesFile != "" || len(config.ValidationSchemas) > 0 {
		layers.IDEIntegration = true
	}

	// 3. Pre-commit hook specifications
	if len(config.CI.Precommit) > 0 {
		layers.PreCommit = true
	}

	// 4. CI check definitions
	if len(config.CI.CIChecks) > 0 {
		layers.CIChecks = true
	}

	// 5. Runtime validation where applicable
	// Check for validation_schemas or drift_policy_ref
	if len(config.ValidationSchemas) > 0 || config.DriftPolicyRef != "" {
		layers.Runtime = true
	}

	// Count total layers
	total := 0
	if layers.PromptLevel {
		total++
	}
	if layers.IDEIntegration {
		total++
	}
	if layers.PreCommit {
		total++
	}
	if layers.CIChecks {
		total++
	}
	if layers.Runtime {
		total++
	}
	layers.TotalLayers = total

	return layers
}

// checkCommitPolicy checks commit policy configuration
func checkCommitPolicy(config *RalphyYAML) CommitPolicyStatus {
	status := CommitPolicyStatus{}

	if len(config.Outputs.CommitPolicy.PrefixRules) > 0 {
		status.HasPrefixRules = true
	}

	if config.Outputs.CommitPolicy.RequireScope {
		status.HasScopeRequirement = true
	}

	if config.Outputs.CommitPolicy.RequireConventionalCommits {
		status.HasConventionalCommits = true
	}

	// Consider complete if at least prefix rules are present
	status.Complete = status.HasPrefixRules

	return status
}

// checkScopeEnforcement checks scope enforcement configuration
func checkScopeEnforcement(config *RalphyYAML) ScopeEnforcement {
	enforcement := ScopeEnforcement{}

	if len(config.Outputs.AllowedFileEdits) > 0 {
		enforcement.HasAllowedFileEdits = true
	}

	if len(config.Outputs.DisallowedFileEdits) > 0 {
		enforcement.HasDisallowedFileEdits = true
	}

	// This will be updated when checking tasks
	enforcement.AllTasksHaveFilesInScope = true // Assume true, will be checked

	// Consider complete if both allowed and disallowed file edits are specified
	enforcement.Complete = enforcement.HasAllowedFileEdits && enforcement.HasDisallowedFileEdits

	return enforcement
}

// checkTasks validates task-level enforcement
func checkTasks(config *RalphyYAML, result ValidationResult) ValidationResult {
	for _, task := range config.Tasks {
		// Check files_in_scope
		if len(task.FilesInScope) == 0 {
			result.ScopeEnforcement.AllTasksHaveFilesInScope = false
			result.Valid = false
			result.Violations = append(result.Violations, Violation{
				Type:        "missing_files_in_scope",
				Description: fmt.Sprintf("Task %q does not have files_in_scope defined", task.ID),
				TaskID:      task.ID,
				Suggestion:  "Add files_in_scope to define which files this task can modify",
			})
		} else {
			result.TasksWithFilesInScope++
		}

		// Check verification
		if len(task.Verification.PreCommit) == 0 && len(task.Verification.PostCommit) == 0 && len(task.Verification.Runtime) == 0 {
			result.Violations = append(result.Violations, Violation{
				Type:        "missing_verification",
				Description: fmt.Sprintf("Task %q does not have verification steps defined", task.ID),
				TaskID:      task.ID,
				Suggestion:  "Add verification.pre_commit, verification.post_commit, or verification.runtime steps",
			})
		} else {
			result.TasksWithVerification++
		}

		// Check single responsibility
		if task.SingleResponsibility == "" {
			result.Violations = append(result.Violations, Violation{
				Type:        "missing_single_responsibility",
				Description: fmt.Sprintf("Task %q does not have single_responsibility defined", task.ID),
				TaskID:      task.ID,
				Suggestion:  "Add single_responsibility to clearly define the task's purpose",
			})
		}
	}

	return result
}

// validateRequirements validates overall requirements
func validateRequirements(config *RalphyYAML, result ValidationResult) ValidationResult {
	// Check minimum verification layers
	if result.VerificationLayers.TotalLayers < minVerificationLayers {
		result.Valid = false
		result.Violations = append(result.Violations, Violation{
			Type:        "insufficient_verification_layers",
			Description: fmt.Sprintf("Only %d verification layers found (minimum %d required)", result.VerificationLayers.TotalLayers, minVerificationLayers),
			Suggestion:  fmt.Sprintf("Add more verification layers (prompt-level, IDE integration, pre-commit, CI checks, runtime validation)"),
		})
	}

	// Check commit policy completeness
	if !result.CommitPolicy.Complete {
		result.Valid = false
		result.Violations = append(result.Violations, Violation{
			Type:        "incomplete_commit_policy",
			Description: "Commit policy is incomplete or missing",
			Suggestion:  "Add commit_policy.prefix_rules to define allowed commit message prefixes",
		})
	}

	// Check scope enforcement completeness
	if !result.ScopeEnforcement.Complete {
		result.Valid = false
		result.Violations = append(result.Violations, Violation{
			Type:        "incomplete_scope_enforcement",
			Description: "Scope enforcement is incomplete",
			Suggestion:  "Add both outputs.allowed_file_edits and outputs.disallowed_file_edits to define file scope boundaries",
		})
	}

	// Check if all tasks have files_in_scope
	if !result.ScopeEnforcement.AllTasksHaveFilesInScope {
		result.Valid = false
		result.Violations = append(result.Violations, Violation{
			Type:        "incomplete_task_scope",
			Description: fmt.Sprintf("%d/%d tasks have files_in_scope defined", result.TasksWithFilesInScope, result.TotalTasks),
			Suggestion:  "Add files_in_scope to all tasks to enforce scope boundaries",
		})
	}

	// Generate recommendations
	if result.VerificationLayers.TotalLayers < 5 {
		result.Recommendations = append(result.Recommendations,
			fmt.Sprintf("Add more verification layers (currently %d/5)", result.VerificationLayers.TotalLayers))
	}

	if !result.CommitPolicy.HasScopeRequirement {
		result.Recommendations = append(result.Recommendations,
			"Consider adding commit_policy.require_scope for better commit organization")
	}

	if !result.CommitPolicy.HasConventionalCommits {
		result.Recommendations = append(result.Recommendations,
			"Consider adding commit_policy.require_conventional_commits for standardized commit messages")
	}

	if result.TasksWithVerification < result.TotalTasks {
		result.Recommendations = append(result.Recommendations,
			fmt.Sprintf("Add verification steps to %d tasks without verification", result.TotalTasks-result.TasksWithVerification))
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

// main is the entry point for the validate_enforcement tool
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

	result := validateEnforcement(config)

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
