package validation

import "github.com/kyledavis/prompt-stack/internal/prompt"

// ShowValidationMsg is a message to show the validation results modal
type ShowValidationMsg struct {
	Results map[string]prompt.ValidationResult
}
