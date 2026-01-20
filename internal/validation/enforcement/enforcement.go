package enforcement

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ExitSuccess   = 0
	ExitFailed    = 1
	ExitExecution = 2
)

const minVerificationLayers = 3

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

type CI struct {
	Precommit []string `yaml:"precommit,omitempty"`
	CIChecks  []string `yaml:"ci_checks,omitempty"`
}

type Outputs struct {
	AllowedFileEdits    []string     `yaml:"allowed_file_edits,omitempty"`
	DisallowedFileEdits []string     `yaml:"disallowed_file_edits,omitempty"`
	CommitPolicy        CommitPolicy `yaml:"commit_policy,omitempty"`
}

type CommitPolicy struct {
	PrefixRules                []string `yaml:"prefix_rules,omitempty"`
	RequireScope               bool     `yaml:"require_scope,omitempty"`
	RequireConventionalCommits bool     `yaml:"require_conventional_commits,omitempty"`
}

type GlobalConstraints struct {
	ForbiddenPatterns      []PatternConstraint `yaml:"forbidden_patterns,omitempty"`
	RequiredPatterns       []PatternConstraint `yaml:"required_patterns,omitempty"`
	AffirmativeConstraints []string            `yaml:"affirmative_constraints,omitempty"`
}

type PatternConstraint struct {
	Pattern string `yaml:"pattern"`
	Message string `yaml:"message"`
	When    string `yaml:"when,omitempty"`
}

type Task struct {
	ID                   string       `yaml:"id"`
	Title                string       `yaml:"title"`
	Description          string       `yaml:"description"`
	FilesInScope         []string     `yaml:"files_in_scope,omitempty"`
	Verification         Verification `yaml:"verification,omitempty"`
	SingleResponsibility string       `yaml:"single_responsibility,omitempty"`
}

type Verification struct {
	PreCommit  []string `yaml:"pre_commit,omitempty"`
	PostCommit []string `yaml:"post_commit,omitempty"`
	Runtime    []string `yaml:"runtime,omitempty"`
}

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

type VerificationLayers struct {
	PromptLevel    bool `json:"prompt_level"`
	IDEIntegration bool `json:"ide_integration"`
	PreCommit      bool `json:"pre_commit"`
	CIChecks       bool `json:"ci_checks"`
	Runtime        bool `json:"runtime"`
	TotalLayers    int  `json:"total_layers"`
}

type CommitPolicyStatus struct {
	HasPrefixRules         bool `json:"has_prefix_rules"`
	HasScopeRequirement    bool `json:"has_scope_requirement"`
	HasConventionalCommits bool `json:"has_conventional_commits"`
	Complete               bool `json:"complete"`
}

type ScopeEnforcement struct {
	HasAllowedFileEdits      bool `json:"has_allowed_file_edits"`
	HasDisallowedFileEdits   bool `json:"has_disallowed_file_edits"`
	AllTasksHaveFilesInScope bool `json:"all_tasks_have_files_in_scope"`
	Complete                 bool `json:"complete"`
}

type Violation struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	TaskID      string `json:"task_id,omitempty"`
	Suggestion  string `json:"suggestion,omitempty"`
}

func LoadYAML(yamlPath string) (*RalphyYAML, error) {
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

func ValidateEnforcement(config *RalphyYAML) ValidationResult {
	result := ValidationResult{
		Valid:              true,
		TotalTasks:         len(config.Tasks),
		Violations:         []Violation{},
		Recommendations:    []string{},
		VerificationLayers: VerificationLayers{},
		CommitPolicy:       CommitPolicyStatus{},
		ScopeEnforcement:   ScopeEnforcement{},
	}

	result.VerificationLayers = checkVerificationLayers(config)
	result.CommitPolicy = checkCommitPolicy(config)
	result.ScopeEnforcement = checkScopeEnforcement(config)
	result = checkTasks(config, result)
	result = validateRequirements(config, result)

	return result
}

func checkVerificationLayers(config *RalphyYAML) VerificationLayers {
	layers := VerificationLayers{}

	if len(config.GlobalConstraints.ForbiddenPatterns) > 0 ||
		len(config.GlobalConstraints.RequiredPatterns) > 0 ||
		len(config.GlobalConstraints.AffirmativeConstraints) > 0 {
		layers.PromptLevel = true
	}

	if config.RulesFile != "" || len(config.ValidationSchemas) > 0 {
		layers.IDEIntegration = true
	}

	if len(config.CI.Precommit) > 0 {
		layers.PreCommit = true
	}

	if len(config.CI.CIChecks) > 0 {
		layers.CIChecks = true
	}

	if len(config.ValidationSchemas) > 0 || config.DriftPolicyRef != "" {
		layers.Runtime = true
	}

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

	status.Complete = status.HasPrefixRules

	return status
}

func checkScopeEnforcement(config *RalphyYAML) ScopeEnforcement {
	enforcement := ScopeEnforcement{}

	if len(config.Outputs.AllowedFileEdits) > 0 {
		enforcement.HasAllowedFileEdits = true
	}

	if len(config.Outputs.DisallowedFileEdits) > 0 {
		enforcement.HasDisallowedFileEdits = true
	}

	enforcement.AllTasksHaveFilesInScope = true
	enforcement.Complete = enforcement.HasAllowedFileEdits && enforcement.HasDisallowedFileEdits

	return enforcement
}

func checkTasks(config *RalphyYAML, result ValidationResult) ValidationResult {
	for _, task := range config.Tasks {
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

func validateRequirements(config *RalphyYAML, result ValidationResult) ValidationResult {
	if result.VerificationLayers.TotalLayers < minVerificationLayers {
		result.Valid = false
		result.Violations = append(result.Violations, Violation{
			Type:        "insufficient_verification_layers",
			Description: fmt.Sprintf("Only %d verification layers found (minimum %d required)", result.VerificationLayers.TotalLayers, minVerificationLayers),
			Suggestion:  "Add more verification layers (prompt-level, IDE integration, pre-commit, CI checks, runtime validation)",
		})
	}

	if !result.CommitPolicy.Complete {
		result.Valid = false
		result.Violations = append(result.Violations, Violation{
			Type:        "incomplete_commit_policy",
			Description: "Commit policy is incomplete or missing",
			Suggestion:  "Add commit_policy.prefix_rules to define allowed commit message prefixes",
		})
	}

	if !result.ScopeEnforcement.Complete {
		result.Valid = false
		result.Violations = append(result.Violations, Violation{
			Type:        "incomplete_scope_enforcement",
			Description: "Scope enforcement is incomplete",
			Suggestion:  "Add both outputs.allowed_file_edits and outputs.disallowed_file_edits to define file scope boundaries",
		})
	}

	if !result.ScopeEnforcement.AllTasksHaveFilesInScope {
		result.Valid = false
		result.Violations = append(result.Violations, Violation{
			Type:        "incomplete_task_scope",
			Description: fmt.Sprintf("%d/%d tasks have files_in_scope defined", result.TasksWithFilesInScope, result.TotalTasks),
			Suggestion:  "Add files_in_scope to all tasks to enforce scope boundaries",
		})
	}

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

func ValidateEnforcementFromFile(yamlPath string) (int, *ValidationResult, error) {
	config, err := LoadYAML(yamlPath)
	if err != nil {
		return ExitExecution, nil, err
	}

	result := ValidateEnforcement(config)

	if !result.Valid {
		return ExitFailed, &result, nil
	}
	return ExitSuccess, &result, nil
}
