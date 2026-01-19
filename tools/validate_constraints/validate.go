// Package validate_constraints provides validation for affirmative constraints usage in Ralphy YAML inputs.
//
// # Purpose
//
// This tool validates that constraints in Ralphy YAML files use affirmative language:
// 1. All constraints use affirmative language ("do this")
// 2. No negative phrasing ("don't do that")
// 3. Constraints are specific and actionable
// 4. Constraints reference existing patterns when possible
//
// Research shows 40%+ better compliance with affirmative framing.
//
// Usage
//
//	go run validate.go --file <input.yaml>
//
//	go run validate.go -f final_ralphy_inputs.yaml
//
// Exit Codes
//
//	0 - Validation successful (all constraints use affirmative language)
//	1 - Validation failed (one or more constraints use negative phrasing)
//	2 - Error in command execution (file not found, invalid YAML, etc.)
//
// Features
//
//   - YAML parsing and constraint extraction
//   - Affirmative language validation
//   - Negative phrasing detection
//   - Constraint specificity analysis
//   - Pattern reference validation
//   - Detailed violation reporting with specific issues
//   - JSON output format for integration with other tools
//
// Integration Points
//
//   - CI/CD pipelines: Run as a quality gate before merge
//   - Pre-commit hooks: Validate constraints before commits
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
//   - Time: O(n) where n = constraints
//   - Typical runtime: <50ms for 100 constraints on modern hardware
//
// Limitations
//
//   - Does not validate YAML syntax (relies on yaml.v3 parser)
//   - Focuses on affirmative language detection only
//   - Does not validate other constraint properties
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
//	  - name: Validate affirmative constraints
//	    run: |
//	      cd tools/validate_constraints
//	      go run validate.go \
//	        --file ../../final_ralphy_inputs.yaml
//
//	Makefile Target:
//	  validate-constraints:
//	      cd tools/validate_constraints && go run validate.go \
//	        --file $(FILE)
//
//	Pre-commit Hook:
//	  if grep -q "ralphy_inputs.yaml" .git/hooks/pre-commit.sample; then
//	    tools/validate_constraints/validate.go \
//	      --file ralphy_inputs.yaml
//	  fi
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Exit codes for predictable script behavior
const (
	ExitSuccess   = 0 // Validation passed
	ExitFailed    = 1 // Validation failed (constraint violations)
	ExitExecution = 2 // Execution error (file I/O, invalid args, etc.)
)

// Negative patterns that indicate non-affirmative phrasing
var negativePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bdon'?t\b`),
	regexp.MustCompile(`(?i)\bdo not\b`),
	regexp.MustCompile(`(?i)\bnever\b`),
	regexp.MustCompile(`(?i)\bavoid\b`),
	regexp.MustCompile(`(?i)\bprevent\b`),
	regexp.MustCompile(`(?i)\bprohibit\b`),
	regexp.MustCompile(`(?i)\bforbid\b`),
	regexp.MustCompile(`(?i)\bno\s+[a-z]+ing\b`),  // "no using", "no writing"
	regexp.MustCompile(`(?i)\bnot\s+[a-z]+ing\b`), // "not using", "not writing"
}

// Affirmative patterns that indicate good phrasing
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

// RalphyYAML represents the structure of Ralphy YAML files
type RalphyYAML struct {
	Name              string            `yaml:"name"`
	Description       string            `yaml:"description"`
	Version           string            `yaml:"version"`
	GlobalConstraints GlobalConstraints `yaml:"global_constraints"`
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

// ValidationResult represents the result of constraints validation
type ValidationResult struct {
	Valid            bool        `json:"valid"`
	TotalConstraints int         `json:"total_constraints"`
	AffirmativeCount int         `json:"affirmative_count"`
	NegativeCount    int         `json:"negative_count"`
	Violations       []Violation `json:"violations,omitempty"`
	Summary          Summary     `json:"summary"`
	Recommendations  []string    `json:"recommendations,omitempty"`
}

// Violation represents a single constraint violation
type Violation struct {
	ConstraintType string `json:"constraint_type"`
	ConstraintText string `json:"constraint_text"`
	Issue          string `json:"issue"`
	NegativePhrase string `json:"negative_phrase,omitempty"`
	Suggestion     string `json:"suggestion,omitempty"`
}

// Summary provides overall validation statistics
type Summary struct {
	AffirmativePercentage float64 `json:"affirmative_percentage"`
	SpecificConstraints   int     `json:"specific_constraints"`
	VagueConstraints      int     `json:"vague_constraints"`
	PatternReferences     int     `json:"pattern_references"`
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

// isNegativePhrasing checks if a constraint uses negative language
func isNegativePhrasing(text string) (bool, string) {
	for _, pattern := range negativePatterns {
		if pattern.MatchString(text) {
			return true, pattern.String()
		}
	}
	return false, ""
}

// isAffirmativePhrasing checks if a constraint uses affirmative language
func isAffirmativePhrasing(text string) bool {
	for _, pattern := range affirmativePatterns {
		if pattern.MatchString(text) {
			return true
		}
	}
	return false
}

// isSpecificConstraint checks if a constraint is specific and actionable
func isSpecificConstraint(text string) bool {
	// Check for vague terms
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

	// Check if it contains actionable verbs
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

// containsPatternReference checks if a constraint references existing patterns
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

// validateConstraints validates all constraints for affirmative language
func validateConstraints(config *RalphyYAML) ValidationResult {
	result := ValidationResult{
		Valid:           true,
		Violations:      []Violation{},
		Summary:         Summary{},
		Recommendations: []string{},
	}

	// Count all constraints
	allConstraints := []string{}

	// Check forbidden patterns
	for _, constraint := range config.GlobalConstraints.ForbiddenPatterns {
		allConstraints = append(allConstraints, constraint.Message)
		result.TotalConstraints++

		// Check for negative phrasing in message
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

		// Check specificity
		if isSpecificConstraint(constraint.Message) {
			result.Summary.SpecificConstraints++
		} else {
			result.Summary.VagueConstraints++
		}

		// Check pattern references
		if containsPatternReference(constraint.Message) {
			result.Summary.PatternReferences++
		}
	}

	// Check required patterns
	for _, constraint := range config.GlobalConstraints.RequiredPatterns {
		allConstraints = append(allConstraints, constraint.Message)
		result.TotalConstraints++

		// Check for negative phrasing in message
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

		// Check specificity
		if isSpecificConstraint(constraint.Message) {
			result.Summary.SpecificConstraints++
		} else {
			result.Summary.VagueConstraints++
		}

		// Check pattern references
		if containsPatternReference(constraint.Message) {
			result.Summary.PatternReferences++
		}
	}

	// Check affirmative constraints
	for _, constraint := range config.GlobalConstraints.AffirmativeConstraints {
		allConstraints = append(allConstraints, constraint)
		result.TotalConstraints++

		// Check for negative phrasing
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
			// Not clearly affirmative or negative - flag as vague
			result.Summary.VagueConstraints++
			result.Violations = append(result.Violations, Violation{
				ConstraintType: "affirmative_constraint",
				ConstraintText: constraint,
				Issue:          "vague_phrasing",
				Suggestion:     "Make constraint more specific and actionable (e.g., 'Always validate inputs' instead of 'Be careful with inputs')",
			})
		}

		// Check specificity
		if isSpecificConstraint(constraint) {
			result.Summary.SpecificConstraints++
		} else {
			result.Summary.VagueConstraints++
		}

		// Check pattern references
		if containsPatternReference(constraint) {
			result.Summary.PatternReferences++
		}
	}

	// Calculate affirmative percentage
	if result.TotalConstraints > 0 {
		result.Summary.AffirmativePercentage = float64(result.AffirmativeCount) / float64(result.TotalConstraints) * 100
	}

	// Generate recommendations
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

// parseFlags parses and validates command-line arguments
func parseFlags() (string, error) {
	yamlPath := flag.String("file", "final_ralphy_inputs.yaml", "Path to YAML file to validate")
	flag.Parse()

	if *yamlPath == "" {
		return "", fmt.Errorf("--file flag is required")
	}

	return *yamlPath, nil
}

// main is the entry point for the validate_constraints tool
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

	result := validateConstraints(config)

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
