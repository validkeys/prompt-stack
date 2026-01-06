package ai

import (
	"encoding/json"
	"fmt"
	"time"
)

// SuggestionType represents the type of AI suggestion
type SuggestionType string

const (
	// SuggestionTypeRecommendation suggests relevant library prompts
	SuggestionTypeRecommendation SuggestionType = "recommendation"

	// SuggestionTypeGap identifies missing context or information
	SuggestionTypeGap SuggestionType = "gap"

	// SuggestionTypeFormatting suggests better structure or organization
	SuggestionTypeFormatting SuggestionType = "formatting"

	// SuggestionTypeContradiction identifies conflicting instructions
	SuggestionTypeContradiction SuggestionType = "contradiction"

	// SuggestionTypeClarity identifies unclear or ambiguous instructions
	SuggestionTypeClarity SuggestionType = "clarity"

	// SuggestionTypeReformatting suggests alternative ways to structure content
	SuggestionTypeReformatting SuggestionType = "reformatting"
)

// SuggestionStatus represents the current status of a suggestion
type SuggestionStatus string

const (
	// SuggestionStatusPending is waiting for user action
	SuggestionStatusPending SuggestionStatus = "pending"

	// SuggestionStatusApplying is currently being applied
	SuggestionStatusApplying SuggestionStatus = "applying"

	// SuggestionStatusApplied has been successfully applied
	SuggestionStatusApplied SuggestionStatus = "applied"

	// SuggestionStatusDismissed has been dismissed by user
	SuggestionStatusDismissed SuggestionStatus = "dismissed"

	// SuggestionStatusError failed to apply
	SuggestionStatusError SuggestionStatus = "error"
)

// Edit represents a single edit operation
type Edit struct {
	// Line is the line number where the edit starts (1-indexed)
	Line int `json:"line"`

	// Column is the column number where the edit starts (1-indexed)
	Column int `json:"column"`

	// OldContent is the content to be replaced
	OldContent string `json:"old_content"`

	// NewContent is the new content to insert
	NewContent string `json:"new_content"`

	// Length is the number of characters to replace
	Length int `json:"length"`
}

// Suggestion represents an AI suggestion for improving the composition
type Suggestion struct {
	// ID is a unique identifier for this suggestion
	ID string `json:"id"`

	// Type is the type of suggestion
	Type SuggestionType `json:"type"`

	// Title is a short, descriptive title
	Title string `json:"title"`

	// Description explains the suggestion in detail
	Description string `json:"description"`

	// ProposedChanges are the edits to apply
	ProposedChanges []Edit `json:"proposed_changes"`

	// Status is the current status of the suggestion
	Status SuggestionStatus `json:"status"`

	// Error contains error message if status is "error"
	Error string `json:"error,omitempty"`

	// CreatedAt is when the suggestion was created
	CreatedAt time.Time `json:"created_at"`

	// AppliedAt is when the suggestion was applied (if applicable)
	AppliedAt *time.Time `json:"applied_at,omitempty"`
}

// NewSuggestion creates a new suggestion with the given parameters
func NewSuggestion(suggestionType SuggestionType, title, description string, changes []Edit) *Suggestion {
	return &Suggestion{
		ID:              generateSuggestionID(),
		Type:            suggestionType,
		Title:           title,
		Description:     description,
		ProposedChanges: changes,
		Status:          SuggestionStatusPending,
		CreatedAt:       time.Now(),
	}
}

// generateSuggestionID generates a unique ID for a suggestion
func generateSuggestionID() string {
	return fmt.Sprintf("suggestion-%d", time.Now().UnixNano())
}

// GetDisplayTitle returns a formatted title for display
func (s *Suggestion) GetDisplayTitle() string {
	icon := s.GetTypeIcon()
	return fmt.Sprintf("%s %s", icon, s.Title)
}

// GetTypeIcon returns an icon representing the suggestion type
func (s *Suggestion) GetTypeIcon() string {
	switch s.Type {
	case SuggestionTypeRecommendation:
		return "💡"
	case SuggestionTypeGap:
		return "🔍"
	case SuggestionTypeFormatting:
		return "📝"
	case SuggestionTypeContradiction:
		return "⚠️"
	case SuggestionTypeClarity:
		return "🔍"
	case SuggestionTypeReformatting:
		return "🔄"
	default:
		return "💡"
	}
}

// GetTypeLabel returns a human-readable label for the suggestion type
func (s *Suggestion) GetTypeLabel() string {
	switch s.Type {
	case SuggestionTypeRecommendation:
		return "Recommendation"
	case SuggestionTypeGap:
		return "Gap Analysis"
	case SuggestionTypeFormatting:
		return "Formatting"
	case SuggestionTypeContradiction:
		return "Contradiction"
	case SuggestionTypeClarity:
		return "Clarity"
	case SuggestionTypeReformatting:
		return "Reformatting"
	default:
		return "Suggestion"
	}
}

// IsApplicable returns true if the suggestion can be applied
func (s *Suggestion) IsApplicable() bool {
	return s.Status == SuggestionStatusPending && len(s.ProposedChanges) > 0
}

// IsDismissed returns true if the suggestion has been dismissed
func (s *Suggestion) IsDismissed() bool {
	return s.Status == SuggestionStatusDismissed
}

// HasError returns true if the suggestion has an error
func (s *Suggestion) HasError() bool {
	return s.Status == SuggestionStatusError && s.Error != ""
}

// MarkAsApplying marks the suggestion as being applied
func (s *Suggestion) MarkAsApplying() {
	s.Status = SuggestionStatusApplying
}

// MarkAsApplied marks the suggestion as successfully applied
func (s *Suggestion) MarkAsApplied() {
	s.Status = SuggestionStatusApplied
	now := time.Now()
	s.AppliedAt = &now
}

// MarkAsDismissed marks the suggestion as dismissed
func (s *Suggestion) MarkAsDismissed() {
	s.Status = SuggestionStatusDismissed
}

// MarkAsError marks the suggestion as having an error
func (s *Suggestion) MarkAsError(err error) {
	s.Status = SuggestionStatusError
	if err != nil {
		s.Error = err.Error()
	}
}

// SuggestionsResponse represents the structured response from Claude containing suggestions
type SuggestionsResponse struct {
	// Suggestions is the list of suggestions
	Suggestions []Suggestion `json:"suggestions"`

	// Summary is a brief summary of the analysis
	Summary string `json:"summary"`
}

// ParseSuggestionsResponse parses a JSON response from Claude into suggestions
func ParseSuggestionsResponse(jsonResponse string) (*SuggestionsResponse, error) {
	var response SuggestionsResponse

	err := json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse suggestions response: %w", err)
	}

	// Validate response
	if len(response.Suggestions) == 0 {
		return nil, fmt.Errorf("no suggestions found in response")
	}

	// Validate each suggestion
	for i, suggestion := range response.Suggestions {
		if suggestion.ID == "" {
			suggestion.ID = generateSuggestionID()
			response.Suggestions[i] = suggestion
		}

		if suggestion.Type == "" {
			return nil, fmt.Errorf("suggestion %d is missing type", i)
		}

		if suggestion.Title == "" {
			return nil, fmt.Errorf("suggestion %d is missing title", i)
		}
	}

	return &response, nil
}

// SuggestionRequest represents a request to Claude for suggestions
type SuggestionRequest struct {
	// Composition is the current composition content
	Composition string `json:"composition"`

	// ContextPrompts are the library prompts included as context
	ContextPrompts []IndexedPrompt `json:"context_prompts,omitempty"`

	// MaxSuggestions is the maximum number of suggestions to return
	MaxSuggestions int `json:"max_suggestions,omitempty"`

	// SuggestionTypes are the types of suggestions to generate
	// If empty, all types are generated
	SuggestionTypes []SuggestionType `json:"suggestion_types,omitempty"`
}

// GetSystemPrompt returns the system prompt for generating suggestions
func GetSystemPrompt() string {
	return `You are an AI assistant that helps improve prompt compositions. Analyze the given composition and provide structured suggestions for improvement.

Your suggestions should fall into these categories:
1. **recommendation**: Suggest relevant library prompts that could enhance the composition
2. **gap**: Identify missing context or information that would make the prompt more effective
3. **formatting**: Suggest better structure or organization of the content
4. **contradiction**: Identify conflicting instructions or requirements
5. **clarity**: Point out unclear or ambiguous instructions that need clarification
6. **reformatting**: Suggest alternative ways to structure the content for better flow

For each suggestion, provide:
- A clear, concise title
- A detailed description explaining why the suggestion is helpful
- Specific edits to apply (with line numbers, old content, and new content)

Format your response as JSON:
{
  "suggestions": [
    {
      "type": "suggestion_type",
      "title": "Short title",
      "description": "Detailed explanation",
      "proposed_changes": [
        {
          "line": 1,
          "column": 1,
          "old_content": "text to replace",
          "new_content": "new text",
          "length": 10
        }
      ]
    }
  ],
  "summary": "Brief summary of your analysis"
}

Be conservative and practical. Only suggest changes that will genuinely improve the prompt's effectiveness.`
}
