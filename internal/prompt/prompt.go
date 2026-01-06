package prompt

import (
	"regexp"
	"strings"
	"time"
)

// Prompt represents a reusable prompt template
type Prompt struct {
	ID               string            `json:"id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Tags             []string          `json:"tags"`
	Category         string            `json:"category"`
	FilePath         string            `json:"file_path"`
	Content          string            `json:"content"`
	Metadata         map[string]string `json:"metadata"`
	Placeholders     []Placeholder     `json:"placeholders"`
	ValidationStatus ValidationResult  `json:"validation_status"`
	UsageStats       UsageMetadata     `json:"usage_stats"`
}

// Placeholder represents a template variable in a prompt
type Placeholder struct {
	Type         string `json:"type"`          // "text" or "list"
	Name         string `json:"name"`          // placeholder name
	StartPos     int    `json:"start_pos"`     // position in content
	EndPos       int    `json:"end_pos"`       // position in content
	CurrentValue string `json:"current_value"` // current filled value
	IsValid      bool   `json:"is_valid"`      // whether syntax is valid
}

// ValidationResult represents validation results for a prompt
type ValidationResult struct {
	Errors   []ValidationError `json:"errors"`
	Warnings []ValidationError `json:"warnings"`
	IsValid  bool              `json:"is_valid"`
}

// ValidationError represents a single validation issue
type ValidationError struct {
	Type    string `json:"type"`    // "error" or "warning"
	Message string `json:"message"` // human-readable message
	Line    int    `json:"line"`    // line number if applicable
	Column  int    `json:"column"`  // column number if applicable
}

// UsageMetadata tracks usage statistics for a prompt
type UsageMetadata struct {
	LastUsed time.Time `json:"last_used"`
	UseCount int       `json:"use_count"`
}

// IndexedPrompt represents a prompt in the library index
type IndexedPrompt struct {
	PromptID      string         `json:"prompt_id"`
	Title         string         `json:"title"`
	Tags          []string       `json:"tags"`
	Category      string         `json:"category"`
	WordFrequency map[string]int `json:"word_frequency"`
	LastUsed      time.Time      `json:"last_used"`
	UseCount      int            `json:"use_count"`
}

// LibraryIndex represents the in-memory index of all prompts
type LibraryIndex struct {
	Prompts   map[string]*IndexedPrompt `json:"prompts"` // keyed by prompt ID
	LastBuilt time.Time                 `json:"last_built"`
	Version   string                    `json:"version"`
}

// ParsePlaceholders extracts placeholders from prompt content
func ParsePlaceholders(content string) []Placeholder {
	// Regex to match {{type:name}} pattern
	re := regexp.MustCompile(`\{\{(\w+):(\w+)\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	var placeholders []Placeholder
	for _, match := range matches {
		if len(match) >= 3 {
			fullMatch := match[0]
			placeholderType := match[1]
			placeholderName := match[2]

			// Find positions in content
			startPos := strings.Index(content, fullMatch)
			if startPos == -1 {
				continue
			}
			endPos := startPos + len(fullMatch)

			placeholder := Placeholder{
				Type:     placeholderType,
				Name:     placeholderName,
				StartPos: startPos,
				EndPos:   endPos,
				IsValid:  isValidPlaceholderType(placeholderType) && isValidPlaceholderName(placeholderName),
			}
			placeholders = append(placeholders, placeholder)
		}
	}

	return placeholders
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

// ExtractKeywords extracts keywords from content for indexing
func ExtractKeywords(content string) map[string]int {
	// Remove frontmatter if present
	content = removeFrontmatter(content)

	// Split into words
	words := strings.Fields(content)

	// Count word frequency
	frequency := make(map[string]int)
	for _, word := range words {
		// Convert to lowercase
		word = strings.ToLower(word)

		// Remove punctuation
		word = strings.Trim(word, ".,!?;:\"'()[]{}")

		// Skip short words and common stop words
		if len(word) < 3 || isStopWord(word) {
			continue
		}

		frequency[word]++
	}

	return frequency
}

// removeFrontmatter removes YAML frontmatter from content
func removeFrontmatter(content string) string {
	if strings.HasPrefix(content, "---") {
		if idx := strings.Index(content[3:], "---"); idx != -1 {
			return content[idx+6:] // Skip both --- markers
		}
	}
	return content
}

// isStopWord checks if a word is a common stop word
func isStopWord(word string) bool {
	stopWords := map[string]bool{
		"the": true, "and": true, "for": true, "are": true,
		"but": true, "not": true, "you": true, "all": true,
		"can": true, "had": true, "her": true, "was": true,
		"one": true, "our": true, "out": true, "has": true,
		"have": true, "been": true, "this": true, "that": true,
		"with": true, "they": true, "from": true, "what": true,
		"when": true, "make": true, "more": true, "will": true,
		"just": true, "know": true, "take": true, "into": true,
		"your": true, "some": true, "could": true, "them": true,
		"than": true, "then": true, "look": true, "only": true,
		"come": true, "over": true, "also": true, "back": true,
		"use": true, "two": true, "how": true, "like": true,
		"first": true, "want": true, "any": true, "work": true,
		"now": true, "such": true, "give": true, "find": true,
	}
	return stopWords[word]
}
