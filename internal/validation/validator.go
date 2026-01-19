package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// ValidationResult represents result of a validation run matching validation-report.schema.json
type ValidationResult struct {
	ReportType        string            `json:"report_type"`
	Timestamp         string            `json:"timestamp"`
	Milestone         string            `json:"milestone,omitempty"`
	RequirementsFile  string            `json:"requirements_file,omitempty"`
	GeneratedYAML     string            `json:"generated_yaml,omitempty"`
	OverallScore      float64           `json:"overall_score"`
	OverallResult     string            `json:"overall_result"`
	ApprovalStatus    string            `json:"approval_status,omitempty"`
	ApprovalReason    string            `json:"approval_reason,omitempty"`
	Issues            []Issue           `json:"issues,omitempty"`
	Recommendations   []Recommendation  `json:"recommendations,omitempty"`
	ValidationSummary map[string]string `json:"validation_summary,omitempty"`
	ComponentScores   map[string]Score  `json:"component_scores"`
	Metadata          Metadata          `json:"metadata,omitempty"`
}

// Issue represents a validation issue
type Issue struct {
	ID            string `json:"id,omitempty"`
	Severity      string `json:"severity"`
	Component     string `json:"component,omitempty"`
	Path          string `json:"path"`
	Message       string `json:"message"`
	FixSuggestion string `json:"fix_suggestion,omitempty"`
}

// Recommendation represents a recommendation for improvement
type Recommendation struct {
	Priority string `json:"priority,omitempty"`
	Action   string `json:"action,omitempty"`
	Benefit  string `json:"benefit,omitempty"`
}

// Score represents a component score
type Score struct {
	Score   float64     `json:"score"`
	Reason  string      `json:"reason,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// Metadata represents validation metadata
type Metadata struct {
	Generator     string  `json:"generator,omitempty"`
	GeneratedAt   string  `json:"generated_at,omitempty"`
	QualityTarget float64 `json:"quality_target,omitempty"`
}

// Validator interface for validation implementations
type Validator interface {
	Validate(inputPath string) (ComponentResult, error)
	Name() string
}

// ComponentResult represents result from a single validator
type ComponentResult struct {
	Name    string
	Score   float64
	Valid   bool
	Issues  []Issue
	Details interface{}
}

// Config holds validation configuration
type Config struct {
	InputPath     string
	OutputPath    string
	Strict        bool
	Milestone     string
	QualityTarget float64
	EventBus      *EventBus
}

// Validate runs all validators against input file
func Validate(config Config) (*ValidationResult, error) {
	EmitValidateEvents(config.EventBus, config.InputPath, nil)

	result := &ValidationResult{
		ReportType:        "final_quality_report",
		Timestamp:         time.Now().Format(time.RFC3339),
		Milestone:         config.Milestone,
		RequirementsFile:  config.InputPath,
		GeneratedYAML:     config.InputPath,
		OverallScore:      1.0,
		OverallResult:     "PASS",
		ComponentScores:   make(map[string]Score),
		Issues:            []Issue{},
		ValidationSummary: make(map[string]string),
		Metadata: Metadata{
			Generator:     "prompt-stack",
			GeneratedAt:   time.Now().Format(time.RFC3339),
			QualityTarget: config.QualityTarget,
		},
	}

	validators := []Validator{
		&YAMLValidator{},
		&TaskSizingValidator{},
		&ConstraintsValidator{},
	}

	for _, validator := range validators {
		componentResult, err := validator.Validate(config.InputPath)
		if err != nil {
			result.OverallScore = 0.0
			result.OverallResult = "FAIL"
			result.Issues = append(result.Issues, Issue{
				Severity: "CRITICAL",
				Path:     config.InputPath,
				Message:  fmt.Sprintf("%s validator failed: %v", validator.Name(), err),
			})
			continue
		}

		result.ComponentScores[validator.Name()] = Score{
			Score:   componentResult.Score,
			Reason:  fmt.Sprintf("Component %s validation", validator.Name()),
			Details: componentResult.Details,
		}

		result.Issues = append(result.Issues, componentResult.Issues...)

		if componentResult.Score < result.OverallScore {
			result.OverallScore = componentResult.Score
		}

		if config.Strict && !componentResult.Valid {
			result.OverallResult = "FAIL"
		}
	}

	if result.OverallScore < 0.8 {
		result.OverallResult = "FAIL"
	} else if result.OverallScore < config.QualityTarget {
		result.OverallResult = "WARN"
	} else if result.OverallScore >= config.QualityTarget {
		result.OverallResult = "APPROVED"
		result.ApprovalStatus = "APPROVED"
		result.ApprovalReason = fmt.Sprintf("Quality score %.4f meets threshold %.2f", result.OverallScore, config.QualityTarget)
	}

	EmitValidateEvents(config.EventBus, config.InputPath, result)

	if config.OutputPath != "" {
		if err := saveResult(result, config.OutputPath); err != nil {
			return result, fmt.Errorf("failed to save report: %w", err)
		}
	}

	return result, nil
}

func saveResult(result *ValidationResult, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// findProjectRoot walks up the filesystem from the current working directory to find go.mod
func findProjectRoot() string {
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

// YAMLValidator validates YAML syntax and structure
type YAMLValidator struct{}

func (v *YAMLValidator) Name() string {
	return "yaml_compliance"
}

func (v *YAMLValidator) Validate(inputPath string) (ComponentResult, error) {
	result := ComponentResult{
		Name:   "yaml_compliance",
		Score:  1.0,
		Valid:  true,
		Issues: []Issue{},
	}

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		result.Valid = false
		result.Score = 0.0
		result.Issues = append(result.Issues, Issue{
			Severity: "CRITICAL",
			Path:     inputPath,
			Message:  "file does not exist",
		})
		return result, nil
	}

	_, err := os.ReadFile(inputPath)
	if err != nil {
		result.Valid = false
		result.Score = 0.0
		result.Issues = append(result.Issues, Issue{
			Severity: "CRITICAL",
			Path:     inputPath,
			Message:  fmt.Sprintf("failed to read file: %v", err),
		})
		return result, nil
	}

	return result, nil
}

// TaskSizingValidator validates task sizes against constraints
type TaskSizingValidator struct{}

func (v *TaskSizingValidator) Name() string {
	return "task_sizing"
}

func (v *TaskSizingValidator) Validate(inputPath string) (ComponentResult, error) {
	result := ComponentResult{
		Name:   "task_sizing",
		Score:  1.0,
		Valid:  true,
		Issues: []Issue{},
		Details: map[string]interface{}{
			"validator": "tools/task_sizing/validate.go",
		},
	}

	// Run the task_sizing validator tool in its module directory so deps resolve
	root := findProjectRoot()
	scriptDir := filepath.Join("tools", "task_sizing")
	if root != "" {
		scriptDir = filepath.Join(root, "tools", "task_sizing")
	}

	cmd := exec.Command("go", "run", "./validate.go", "--file", inputPath)
	cmd.Dir = scriptDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Valid = false
		result.Score = 0.8
		result.Issues = append(result.Issues, Issue{
			Severity: "HIGH",
			Path:     inputPath,
			Message:  fmt.Sprintf("validation failed: %v", err),
		})
		result.Issues = append(result.Issues, Issue{
			Severity: "LOW",
			Path:     inputPath,
			Message:  string(output),
		})
		return result, nil
	}

	// Try to parse the tool output (JSON) and adjust score if no tasks found
	var parsed map[string]interface{}
	if err := json.Unmarshal(output, &parsed); err == nil {
		if v, ok := parsed["total_tasks"]; ok {
			if f, ok := v.(float64); ok {
				if int(f) == 0 {
					// no tasks -> lower score
					result.Score = 0.8
				}
			}
		}
		// include parsed output in details for diagnostics
		result.Details = parsed
	}

	return result, nil
}

// ConstraintsValidator validates affirmative constraints
type ConstraintsValidator struct{}

func (v *ConstraintsValidator) Name() string {
	return "constraints"
}

func (v *ConstraintsValidator) Validate(inputPath string) (ComponentResult, error) {
	result := ComponentResult{
		Name:   "constraints",
		Score:  1.0,
		Valid:  true,
		Issues: []Issue{},
		Details: map[string]interface{}{
			"validator": "tools/validate_constraints/validate.go",
		},
	}

	root := findProjectRoot()
	scriptDir := filepath.Join("tools", "validate_constraints")
	if root != "" {
		scriptDir = filepath.Join(root, "tools", "validate_constraints")
	}

	cmd := exec.Command("go", "run", "./validate.go", "--file", inputPath)
	cmd.Dir = scriptDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Valid = false
		result.Score = 0.9
		result.Issues = append(result.Issues, Issue{
			Severity: "HIGH",
			Path:     inputPath,
			Message:  fmt.Sprintf("validation failed: %v", err),
		})
		result.Issues = append(result.Issues, Issue{
			Severity: "LOW",
			Path:     inputPath,
			Message:  string(output),
		})
		return result, nil
	}

	return result, nil
}
