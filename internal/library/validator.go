package library

import (
	"fmt"
	"os"
	"strings"

	"github.com/kyledavis/prompt-stack/internal/prompt"
	"go.uber.org/zap"
)

const (
	// MaxFileSize is the maximum allowed file size (1MB)
	MaxFileSize = 1 * 1024 * 1024
)

// Validate validates a prompt and returns validation results
func Validate(p *prompt.Prompt, logger *zap.Logger) prompt.ValidationResult {
	result := prompt.ValidationResult{
		Errors:   []prompt.ValidationError{},
		Warnings: []prompt.ValidationError{},
		IsValid:  true,
	}

	// Check file size
	validateFileSize(p, &result, logger)

	// Check required fields
	validateRequiredFields(p, &result, logger)

	// Check placeholders
	validatePlaceholders(p, &result, logger)

	// Check YAML frontmatter
	validateFrontmatter(p, &result, logger)

	// Set overall validity
	result.IsValid = len(result.Errors) == 0

	return result
}

// validateFileSize checks if file size exceeds limit
func validateFileSize(p *prompt.Prompt, result *prompt.ValidationResult, logger *zap.Logger) {
	fileInfo, err := os.Stat(p.FilePath)
	if err != nil {
		result.Errors = append(result.Errors, prompt.ValidationError{
			Type:    "error",
			Message: fmt.Sprintf("Cannot read file: %v", err),
		})
		return
	}

	if fileInfo.Size() > MaxFileSize {
		result.Errors = append(result.Errors, prompt.ValidationError{
			Type:    "error",
			Message: fmt.Sprintf("File size exceeds %d bytes limit (current: %d bytes)", MaxFileSize, fileInfo.Size()),
		})
	}
}

// validateRequiredFields checks for required fields
func validateRequiredFields(p *prompt.Prompt, result *prompt.ValidationResult, logger *zap.Logger) {
	if p.Title == "" {
		result.Errors = append(result.Errors, prompt.ValidationError{
			Type:    "error",
			Message: "Missing required field: title",
		})
	}

	// Optional fields warnings
	if p.Description == "" {
		result.Warnings = append(result.Warnings, prompt.ValidationError{
			Type:    "warning",
			Message: "Missing optional field: description",
		})
	}

	if len(p.Tags) == 0 {
		result.Warnings = append(result.Warnings, prompt.ValidationError{
			Type:    "warning",
			Message: "Missing optional field: tags",
		})
	}
}

// validatePlaceholders checks placeholder syntax and duplicates
func validatePlaceholders(p *prompt.Prompt, result *prompt.ValidationResult, logger *zap.Logger) {
	// Track placeholder names to detect duplicates
	seenNames := make(map[string]bool)

	for _, ph := range p.Placeholders {
		// Check if placeholder is valid
		if !ph.IsValid {
			result.Warnings = append(result.Warnings, prompt.ValidationError{
				Type:    "warning",
				Message: fmt.Sprintf("Invalid placeholder syntax: {{%s:%s}}", ph.Type, ph.Name),
			})
			continue
		}

		// Check for duplicate names
		if seenNames[ph.Name] {
			result.Errors = append(result.Errors, prompt.ValidationError{
				Type:    "error",
				Message: fmt.Sprintf("Duplicate placeholder name: %s", ph.Name),
			})
		} else {
			seenNames[ph.Name] = true
		}
	}
}

// validateFrontmatter checks YAML frontmatter syntax
func validateFrontmatter(p *prompt.Prompt, result *prompt.ValidationResult, logger *zap.Logger) {
	// Check if frontmatter exists
	if len(p.Metadata) == 0 {
		result.Warnings = append(result.Warnings, prompt.ValidationError{
			Type:    "warning",
			Message: "Missing YAML frontmatter",
		})
		return
	}

	// Check for invalid YAML syntax
	// This is a basic check - could be enhanced with proper YAML parser
	for key, value := range p.Metadata {
		if strings.Contains(value, "\t") {
			result.Warnings = append(result.Warnings, prompt.ValidationError{
				Type:    "warning",
				Message: fmt.Sprintf("YAML frontmatter may contain tabs (use spaces instead): %s", key),
			})
		}
	}
}

// ValidateLibrary validates all prompts in the library
func (l *Library) ValidateLibrary() map[string]prompt.ValidationResult {
	results := make(map[string]prompt.ValidationResult)

	for filePath, p := range l.Prompts {
		results[filePath] = Validate(p, l.logger)
	}

	return results
}

// GetInvalidPrompts returns prompts with validation errors
func (l *Library) GetInvalidPrompts() []*prompt.Prompt {
	var invalid []*prompt.Prompt

	for _, p := range l.Prompts {
		if !p.ValidationStatus.IsValid {
			invalid = append(invalid, p)
		}
	}

	return invalid
}

// GetPromptsWithWarnings returns prompts with validation warnings
func (l *Library) GetPromptsWithWarnings() []*prompt.Prompt {
	var withWarnings []*prompt.Prompt

	for _, p := range l.Prompts {
		if len(p.ValidationStatus.Warnings) > 0 {
			withWarnings = append(withWarnings, p)
		}
	}

	return withWarnings
}
