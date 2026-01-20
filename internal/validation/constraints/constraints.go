package constraints

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	ExitSuccess   = 0
	ExitFailed    = 1
	ExitExecution = 2
)

var negativePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bdon'?t\b`),
	regexp.MustCompile(`(?i)\bdo not\b`),
	regexp.MustCompile(`(?i)\bnever\b`),
	regexp.MustCompile(`(?i)\bavoid\b`),
	regexp.MustCompile(`(?i)\bprevent\b`),
	regexp.MustCompile(`(?i)\bprohibit\b`),
	regexp.MustCompile(`(?i)\bforbid\b`),
	regexp.MustCompile(`(?i)\bno\s+[a-z]+ing\b`),
	regexp.MustCompile(`(?i)\bnot\s+[a-z]+ing\b`),
}

var affirmativePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\balways\b`),
	regexp.MustCompile(`(?i)\bmust\b`),
	regexp.MustCompile(`(?i)\bshould\b`),
	regexp.MustCompile(`(?i)\bshall\b`),
	regexp.MustCompile(`(?i)\buse\b`),
	regexp.MustCompile(`(?i)\bfollow\b`),
	regexp.MustCompile(`(?i)\bwrite\b`),
	regexp.MustCompile(`(?i)\bcreate\b`),
	regexp.MustCompile(`(?i)\bimplement\b`),
	regexp.MustCompile(`(?i)\bvalidate\b`),
	regexp.MustCompile(`(?i)\bhandle\b`),
	regexp.MustCompile(`(?i)\baddress\b`),
	regexp.MustCompile(`(?i)\bfix\b`),
	regexp.MustCompile(`(?i)\bonly\b`),
}

type RalphyYAML struct {
	Name              string            `yaml:"name"`
	Description       string            `yaml:"description"`
	Version           string            `yaml:"version"`
	GlobalConstraints GlobalConstraints `yaml:"global_constraints"`
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

type ValidationResult struct {
	Valid            bool        `json:"valid"`
	TotalConstraints int         `json:"total_constraints"`
	AffirmativeCount int         `json:"affirmative_count"`
	NegativeCount    int         `json:"negative_count"`
	Violations       []Violation `json:"violations,omitempty"`
	Summary          Summary     `json:"summary"`
	Recommendations  []string    `json:"recommendations,omitempty"`
}

type Violation struct {
	ConstraintType string `json:"constraint_type"`
	ConstraintText string `json:"constraint_text"`
	Issue          string `json:"issue"`
	NegativePhrase string `json:"negative_phrase,omitempty"`
	Suggestion     string `json:"suggestion,omitempty"`
}

type Summary struct {
	AffirmativePercentage float64 `json:"affirmative_percentage"`
	SpecificConstraints   int     `json:"specific_constraints"`
	VagueConstraints      int     `json:"vague_constraints"`
	PatternReferences     int     `json:"pattern_references"`
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

func ValidateConstraints(config *RalphyYAML) ValidationResult {
	result := ValidationResult{
		Valid:           true,
		Violations:      []Violation{},
		Summary:         Summary{},
		Recommendations: []string{},
	}

	for _, constraint := range config.GlobalConstraints.ForbiddenPatterns {
		result.TotalConstraints++

		if isNegative, phrase := isNegativePhrasing(constraint.Message); isNegative {
			result.Valid = false
			result.NegativeCount++
			result.Violations = append(result.Violations, Violation{
				ConstraintType: "forbidden_pattern",
				ConstraintText: constraint.Message,
				Issue:          "negative_phrasing",
				NegativePhrase: phrase,
				Suggestion:     "Rephrase to use affirmative language (e.g., 'Use X instead of Y')",
			})
		} else if isAffirmativePhrasing(constraint.Message) {
			result.AffirmativeCount++
		}

		if isSpecificConstraint(constraint.Message) {
			result.Summary.SpecificConstraints++
		} else {
			result.Summary.VagueConstraints++
		}

		if containsPatternReference(constraint.Message) {
			result.Summary.PatternReferences++
		}
	}

	for _, constraint := range config.GlobalConstraints.RequiredPatterns {
		result.TotalConstraints++

		if isNegative, phrase := isNegativePhrasing(constraint.Message); isNegative {
			result.Valid = false
			result.NegativeCount++
			result.Violations = append(result.Violations, Violation{
				ConstraintType: "required_pattern",
				ConstraintText: constraint.Message,
				Issue:          "negative_phrasing",
				NegativePhrase: phrase,
				Suggestion:     "Rephrase to use affirmative language (e.g., 'Always use X when Y')",
			})
		} else if isAffirmativePhrasing(constraint.Message) {
			result.AffirmativeCount++
		}

		if isSpecificConstraint(constraint.Message) {
			result.Summary.SpecificConstraints++
		} else {
			result.Summary.VagueConstraints++
		}

		if containsPatternReference(constraint.Message) {
			result.Summary.PatternReferences++
		}
	}

	for _, constraint := range config.GlobalConstraints.AffirmativeConstraints {
		result.TotalConstraints++

		if isNegative, phrase := isNegativePhrasing(constraint); isNegative {
			result.Valid = false
			result.NegativeCount++
			result.Violations = append(result.Violations, Violation{
				ConstraintType: "affirmative_constraint",
				ConstraintText: constraint,
				Issue:          "negative_phrasing",
				NegativePhrase: phrase,
				Suggestion:     "Rephrase to use affirmative language (e.g., 'Always do X' instead of 'Don't do Y')",
			})
		} else if isAffirmativePhrasing(constraint) {
			result.AffirmativeCount++
		} else {
			result.Summary.VagueConstraints++
			result.Violations = append(result.Violations, Violation{
				ConstraintType: "affirmative_constraint",
				ConstraintText: constraint,
				Issue:          "vague_phrasing",
				Suggestion:     "Make constraint more specific and actionable (e.g., 'Always validate inputs' instead of 'Be careful with inputs')",
			})
		}

		if isSpecificConstraint(constraint) {
			result.Summary.SpecificConstraints++
		} else {
			result.Summary.VagueConstraints++
		}

		if containsPatternReference(constraint) {
			result.Summary.PatternReferences++
		}
	}

	if result.TotalConstraints > 0 {
		result.Summary.AffirmativePercentage = float64(result.AffirmativeCount) / float64(result.TotalConstraints) * 100
	}

	if result.Summary.AffirmativePercentage < 100 {
		result.Recommendations = append(result.Recommendations,
			fmt.Sprintf("Increase affirmative phrasing from %.1f%% to 100%%", result.Summary.AffirmativePercentage))
	}

	if result.Summary.VagueConstraints > 0 {
		result.Recommendations = append(result.Recommendations,
			fmt.Sprintf("Make %d vague constraints more specific and actionable", result.Summary.VagueConstraints))
	}

	if result.Summary.PatternReferences < result.TotalConstraints/2 {
		result.Recommendations = append(result.Recommendations,
			"Reference existing patterns and examples in more constraints")
	}

	return result
}

func isNegativePhrasing(text string) (bool, string) {
	for _, pattern := range negativePatterns {
		if pattern.MatchString(text) {
			return true, pattern.String()
		}
	}
	return false, ""
}

func isAffirmativePhrasing(text string) bool {
	for _, pattern := range affirmativePatterns {
		if pattern.MatchString(text) {
			return true
		}
	}
	return false
}

func isSpecificConstraint(text string) bool {
	vagueTerms := []string{
		"properly", "correctly", "appropriately", "well",
		"good", "bad", "better", "best", "worse", "worst",
		"some", "any", "all", "every", "none",
	}

	lowerText := strings.ToLower(text)
	for _, term := range vagueTerms {
		if strings.Contains(lowerText, term) {
			return false
		}
	}

	actionableVerbs := []string{
		"use", "follow", "write", "create", "implement",
		"validate", "handle", "test", "run", "check",
		"verify", "ensure", "require", "include",
	}

	for _, verb := range actionableVerbs {
		if strings.Contains(lowerText, verb) {
			return true
		}
	}

	return false
}

func containsPatternReference(text string) bool {
	patternReferences := []string{
		"pattern", "example", "reference", "anchor",
		"style", "template", "specification", "standard",
		"guideline", "best practice", "convention",
	}

	lowerText := strings.ToLower(text)
	for _, ref := range patternReferences {
		if strings.Contains(lowerText, ref) {
			return true
		}
	}

	return false
}

func ValidateConstraintsFromFile(yamlPath string) (int, *ValidationResult, error) {
	config, err := LoadYAML(yamlPath)
	if err != nil {
		return ExitExecution, nil, err
	}

	result := ValidateConstraints(config)

	if !result.Valid {
		return ExitFailed, &result, nil
	}
	return ExitSuccess, &result, nil
}
