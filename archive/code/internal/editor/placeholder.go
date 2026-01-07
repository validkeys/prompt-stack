package editor

import (
	"regexp"
	"strings"
)

// Placeholder represents a template variable in a composition
type Placeholder struct {
	Type         string   `json:"type"`          // "text" or "list"
	Name         string   `json:"name"`          // placeholder name
	StartPos     int      `json:"start_pos"`     // position in content
	EndPos       int      `json:"end_pos"`       // position in content
	CurrentValue string   `json:"current_value"` // current filled value (for text)
	ListValues   []string `json:"list_values"`   // current filled values (for list)
	IsValid      bool     `json:"is_valid"`      // whether syntax is valid
	IsActive     bool     `json:"is_active"`     // whether currently selected
}

// PlaceholderState tracks all placeholders in a composition
type PlaceholderState struct {
	Placeholders []Placeholder `json:"placeholders"`
	ActiveIndex  int           `json:"active_index"` // -1 if none active
}

// ParsePlaceholders extracts all placeholders from content
func ParsePlaceholders(content string) []Placeholder {
	// Regex to match {{type:name}} pattern
	re := regexp.MustCompile(`\{\{(\w+):(\w+)\}\}`)
	matches := re.FindAllStringSubmatchIndex(content, -1)

	var placeholders []Placeholder
	for _, match := range matches {
		if len(match) >= 4 {
			fullMatchStart := match[0]
			fullMatchEnd := match[1]
			typeStart := match[2]
			typeEnd := match[3]
			nameStart := match[4]
			nameEnd := match[5]

			// Extract type and name
			placeholderType := content[typeStart:typeEnd]
			placeholderName := content[nameStart:nameEnd]

			placeholder := Placeholder{
				Type:     placeholderType,
				Name:     placeholderName,
				StartPos: fullMatchStart,
				EndPos:   fullMatchEnd,
				IsValid:  isValidPlaceholderType(placeholderType) && isValidPlaceholderName(placeholderName),
				IsActive: false,
			}

			// Initialize current value based on type
			if placeholderType == "list" {
				placeholder.ListValues = []string{}
			} else {
				placeholder.CurrentValue = ""
			}

			placeholders = append(placeholders, placeholder)
		}
	}

	return placeholders
}

// ValidatePlaceholders checks for duplicate names and validates all placeholders
func ValidatePlaceholders(placeholders []Placeholder) []ValidationError {
	var errors []ValidationError
	nameMap := make(map[string]int)

	// Check for duplicate names
	for i, ph := range placeholders {
		if !ph.IsValid {
			continue // Skip invalid placeholders for duplicate check
		}

		if prevIndex, exists := nameMap[ph.Name]; exists {
			errors = append(errors, ValidationError{
				Type:    "error",
				Message: "Duplicate placeholder name: " + ph.Name,
				Line:    getLineNumber(placeholders, i),
				Column:  ph.StartPos,
			})
			errors = append(errors, ValidationError{
				Type:    "error",
				Message: "Duplicate placeholder name: " + ph.Name,
				Line:    getLineNumber(placeholders, prevIndex),
				Column:  placeholders[prevIndex].StartPos,
			})
		} else {
			nameMap[ph.Name] = i
		}
	}

	// Validate each placeholder
	for i, ph := range placeholders {
		if !ph.IsValid {
			if !isValidPlaceholderType(ph.Type) {
				errors = append(errors, ValidationError{
					Type:    "error",
					Message: "Invalid placeholder type: " + ph.Type + " (must be 'text' or 'list')",
					Line:    getLineNumber(placeholders, i),
					Column:  ph.StartPos,
				})
			}
			if !isValidPlaceholderName(ph.Name) {
				errors = append(errors, ValidationError{
					Type:    "error",
					Message: "Invalid placeholder name: " + ph.Name + " (must be alphanumeric and underscores only)",
					Line:    getLineNumber(placeholders, i),
					Column:  ph.StartPos,
				})
			}
		}
	}

	return errors
}

// FindPlaceholderAtPosition finds the placeholder at the given cursor position
func FindPlaceholderAtPosition(placeholders []Placeholder, pos int) int {
	for i, ph := range placeholders {
		if pos >= ph.StartPos && pos <= ph.EndPos {
			return i
		}
	}
	return -1
}

// GetNextPlaceholder finds the next placeholder after the current position
func GetNextPlaceholder(placeholders []Placeholder, currentPos int) int {
	for i, ph := range placeholders {
		if ph.StartPos > currentPos {
			return i
		}
	}
	return -1
}

// GetPreviousPlaceholder finds the previous placeholder before the current position
func GetPreviousPlaceholder(placeholders []Placeholder, currentPos int) int {
	for i := len(placeholders) - 1; i >= 0; i-- {
		if placeholders[i].EndPos < currentPos {
			return i
		}
	}
	return -1
}

// ReplacePlaceholder replaces a placeholder with its filled value in content
func ReplacePlaceholder(content string, ph Placeholder) string {
	var replacement string
	if ph.Type == "list" {
		// Format as markdown bullet list
		if len(ph.ListValues) > 0 {
			items := make([]string, len(ph.ListValues))
			for i, val := range ph.ListValues {
				items[i] = "- " + val
			}
			replacement = strings.Join(items, "\n")
		} else {
			replacement = ""
		}
	} else {
		replacement = ph.CurrentValue
	}

	// Replace placeholder with value
	before := content[:ph.StartPos]
	after := content[ph.EndPos:]
	return before + replacement + after
}

// isValidPlaceholderType checks if placeholder type is valid
func isValidPlaceholderType(typ string) bool {
	return typ == "text" || typ == "list"
}

// isValidPlaceholderName checks if placeholder name is valid
func isValidPlaceholderName(name string) bool {
	if name == "" {
		return false
	}
	// Must be alphanumeric and underscores only
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return true
}

// getLineNumber calculates the line number for a placeholder
func getLineNumber(placeholders []Placeholder, index int) int {
	if index < 0 || index >= len(placeholders) {
		return 0
	}
	// This is a simplified calculation - in practice, we'd need the full content
	// to accurately determine line numbers. For now, return 1-based index.
	return index + 1
}

// ValidationError represents a validation error
type ValidationError struct {
	Type    string `json:"type"`    // "error" or "warning"
	Message string `json:"message"` // human-readable message
	Line    int    `json:"line"`    // line number
	Column  int    `json:"column"`  // column number
}
