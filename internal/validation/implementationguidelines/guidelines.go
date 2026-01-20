package implementationguidelines

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	ExitSuccess   = 0
	ExitFailed    = 1
	ExitExecution = 2
)

var testFirstKeywords = []string{
	"test", "testing", "unit test", "integration test", "tdd",
	"table-driven", "test coverage", "coverage", "assert",
	"verify", "validate", "check", "ensure",
}

var testableCriteriaKeywords = []string{
	"passes", "fails", "returns", "outputs", "matches",
	"contains", "equals", "greater than", "less than",
	"successfully", "completes", "executes", "runs",
	"valid", "invalid", "correct", "incorrect",
	"pass", "test", "coverage", "expected",
}

var implementationGuidanceKeywords = []string{
	"guideline", "guide", "pattern", "example",
	"anchor", "style", "template", "specification", "standard",
	"best practice", "convention", "approach", "method",
	"technique", "strategy", "implementation", "development",
	"reference implementation", "coding standard", "design pattern",
}

type RalphyYAML struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	TDD         TDD    `yaml:"tdd"`
	Tasks       []Task `yaml:"tasks"`
}

type TDD struct {
	Required           bool   `yaml:"required"`
	TestCommand        string `yaml:"test_command"`
	FailureInstruction string `yaml:"failure_instruction"`
}

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

type StyleAnchor struct {
	File   string `yaml:"file"`
	Reason string `yaml:"reason"`
}

type Verification struct {
	PreCommit []string `yaml:"pre_commit"`
}

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

type Violation struct {
	TaskID     string `json:"task_id"`
	TaskTitle  string `json:"task_title"`
	Issue      string `json:"issue"`
	Details    string `json:"details"`
	Suggestion string `json:"suggestion,omitempty"`
}

type Summary struct {
	TestFirstCoverage              float64 `json:"test_first_coverage"`
	TestableCriteriaCoverage       float64 `json:"testable_criteria_coverage"`
	ImplementationGuidanceCoverage float64 `json:"implementation_guidance_coverage"`
	OverallScore                   float64 `json:"overall_score"`
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

func ValidateImplementationGuidelines(config *RalphyYAML) ValidationResult {
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

	if result.TotalTasks > 0 {
		result.Summary.TestFirstCoverage = float64(result.TasksNeedingTests) / float64(result.TotalTasks) * 100
		result.Summary.TestableCriteriaCoverage = float64(result.TasksWithTestableCriteria) / float64(result.TotalTasks) * 100
		result.Summary.ImplementationGuidanceCoverage = float64(result.TasksWithImplementationGuidance) / float64(result.TotalTasks) * 100
		result.Summary.OverallScore = (result.Summary.TestFirstCoverage*0.4 + result.Summary.TestableCriteriaCoverage*0.3 + result.Summary.ImplementationGuidanceCoverage*0.3) / 100
	}

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

func requiresTestFirstWorkflow(task Task) bool {
	lowerDesc := strings.ToLower(task.Description)
	for _, keyword := range testFirstKeywords {
		if strings.Contains(lowerDesc, keyword) {
			return true
		}
	}

	lowerTitle := strings.ToLower(task.Title)
	for _, keyword := range testFirstKeywords {
		if strings.Contains(lowerTitle, keyword) {
			return true
		}
	}

	for _, file := range task.FilesInScope {
		if strings.Contains(file, "_test.go") || strings.Contains(file, "_test.") {
			return true
		}
	}

	return false
}

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

func hasImplementationGuidance(task Task) bool {
	for _, anchor := range task.StyleAnchors {
		lowerReason := strings.ToLower(anchor.Reason)
		for _, keyword := range implementationGuidanceKeywords {
			if strings.Contains(lowerReason, keyword) {
				return true
			}
		}
	}

	lowerDesc := strings.ToLower(task.Description)
	for _, keyword := range implementationGuidanceKeywords {
		if strings.Contains(lowerDesc, keyword) {
			return true
		}
	}

	return false
}

func ValidateImplementationGuidelinesFromFile(yamlPath string) (int, *ValidationResult, error) {
	config, err := LoadYAML(yamlPath)
	if err != nil {
		return ExitExecution, nil, err
	}

	result := ValidateImplementationGuidelines(config)

	if !result.Valid {
		return ExitFailed, &result, nil
	}
	return ExitSuccess, &result, nil
}
