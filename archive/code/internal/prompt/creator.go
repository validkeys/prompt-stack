package prompt

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Creator handles guided prompt creation workflow
type Creator struct {
	logger  *zap.Logger
	storage *Storage
}

// NewCreator creates a new prompt creator
func NewCreator(logger *zap.Logger) *Creator {
	return &Creator{
		logger: logger,
	}
}

// NewCreatorWithStorage creates a new prompt creator with storage
func NewCreatorWithStorage(logger *zap.Logger, storage *Storage) *Creator {
	return &Creator{
		logger:  logger,
		storage: storage,
	}
}

// Step represents a step in the creation wizard
type Step int

const (
	StepTitle Step = iota
	StepDescription
	StepTags
	StepCategory
	StepContent
	StepReview
	StepComplete
)

// CreationState represents the state of the creation wizard
type CreationState struct {
	CurrentStep      Step
	Title            string
	Description      string
	Tags             []string
	Category         string
	Content          string
	ValidationErrors []string
}

// NewCreationState creates a new creation state
func NewCreationState() *CreationState {
	return &CreationState{
		CurrentStep:      StepTitle,
		Tags:             []string{},
		ValidationErrors: []string{},
	}
}

// StepInfo represents information about a wizard step
type StepInfo struct {
	Title       string
	Description string
	IsRequired  bool
	IsComplete  bool
}

// GetStepInfo returns information about the current step
func (c *Creator) GetStepInfo(state *CreationState) StepInfo {
	switch state.CurrentStep {
	case StepTitle:
		return StepInfo{
			Title:       "Title",
			Description: "Enter a title for your prompt (required)",
			IsRequired:  true,
			IsComplete:  state.Title != "",
		}

	case StepDescription:
		return StepInfo{
			Title:       "Description",
			Description: "Enter a description for your prompt (optional)",
			IsRequired:  false,
			IsComplete:  state.Description != "",
		}

	case StepTags:
		return StepInfo{
			Title:       "Tags",
			Description: "Enter tags separated by commas (optional)",
			IsRequired:  false,
			IsComplete:  len(state.Tags) > 0,
		}

	case StepCategory:
		return StepInfo{
			Title:       "Category",
			Description: "Select a category for your prompt (required)",
			IsRequired:  true,
			IsComplete:  state.Category != "",
		}

	case StepContent:
		return StepInfo{
			Title:       "Content",
			Description: "Enter the prompt content (markdown supported)",
			IsRequired:  true,
			IsComplete:  state.Content != "",
		}

	case StepReview:
		return StepInfo{
			Title:       "Review",
			Description: "Review your prompt before saving",
			IsRequired:  false,
			IsComplete:  true, // Review is always complete
		}

	case StepComplete:
		return StepInfo{
			Title:       "Complete",
			Description: "Prompt created successfully",
			IsRequired:  false,
			IsComplete:  true,
		}

	default:
		return StepInfo{
			Title:       "Unknown",
			Description: "Unknown step",
			IsRequired:  false,
			IsComplete:  false,
		}
	}
}

// ValidateStep validates the current step
func (c *Creator) ValidateStep(state *CreationState) []string {
	var errors []string

	switch state.CurrentStep {
	case StepTitle:
		if state.Title == "" {
			errors = append(errors, "Title is required")
		} else if len(state.Title) > 100 {
			errors = append(errors, "Title must be less than 100 characters")
		}

	case StepDescription:
		// Description is optional, no validation needed
		if len(state.Description) > 500 {
			errors = append(errors, "Description must be less than 500 characters")
		}

	case StepTags:
		// Tags are optional, but validate format if provided
		for _, tag := range state.Tags {
			tag = strings.TrimSpace(tag)
			if tag == "" {
				continue
			}
			if len(tag) > 50 {
				errors = append(errors, fmt.Sprintf("Tag '%s' must be less than 50 characters", tag))
			}
			if strings.Contains(tag, ",") {
				errors = append(errors, fmt.Sprintf("Tag '%s' contains invalid characters", tag))
			}
		}

	case StepCategory:
		if state.Category == "" {
			errors = append(errors, "Category is required")
		}

	case StepContent:
		if state.Content == "" {
			errors = append(errors, "Content is required")
		} else if len(state.Content) > 1000000 { // 1MB limit
			errors = append(errors, "Content must be less than 1MB")
		}

	case StepReview, StepComplete:
		// No validation needed for these steps
	}

	return errors
}

// CanProceed checks if user can proceed to next step
func (c *Creator) CanProceed(state *CreationState) bool {
	errors := c.ValidateStep(state)
	return len(errors) == 0
}

// CanGoBack checks if user can go to previous step
func (c *Creator) CanGoBack(state *CreationState) bool {
	return state.CurrentStep > StepTitle
}

// GoToNext advances to the next step
func (c *Creator) GoToNext(state *CreationState) {
	if state.CurrentStep < StepComplete {
		state.CurrentStep++
	}
}

// GoToPrevious goes to the previous step
func (c *Creator) GoToPrevious(state *CreationState) {
	if state.CurrentStep > StepTitle {
		state.CurrentStep--
	}
}

// GoToStep jumps to a specific step
func (c *Creator) GoToStep(state *CreationState, step Step) {
	if step >= StepTitle && step <= StepComplete {
		state.CurrentStep = step
	}
}

// GetProgress returns the progress percentage
func (c *Creator) GetProgress(state *CreationState) int {
	totalSteps := int(StepComplete) + 1
	currentStep := int(state.CurrentStep) + 1
	return (currentStep * 100) / totalSteps
}

// GetPromptData returns the final prompt data
func (c *Creator) GetPromptData(state *CreationState) *PromptData {
	return &PromptData{
		Title:       state.Title,
		Description: state.Description,
		Tags:        state.Tags,
		Category:    state.Category,
		Content:     state.Content,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
}

// PromptData represents the final prompt data
type PromptData struct {
	Title       string
	Description string
	Tags        []string
	Category    string
	Content     string
	CreatedAt   string
}

// GenerateFilename generates a filename from the title
func (c *Creator) GenerateFilename(title string) string {
	// Convert to kebab-case
	filename := strings.ToLower(title)
	filename = strings.ReplaceAll(filename, " ", "-")

	// Remove special characters
	var result []rune
	for _, r := range filename {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result = append(result, r)
		}
	}
	filename = string(result)

	// Remove consecutive hyphens
	for strings.Contains(filename, "--") {
		filename = strings.ReplaceAll(filename, "--", "-")
	}

	// Trim hyphens from start/end
	filename = strings.Trim(filename, "-")

	return fmt.Sprintf("%s.md", filename)
}

// ValidateFinal validates the complete prompt before saving
func (c *Creator) ValidateFinal(state *CreationState) []string {
	var errors []string

	// Validate all steps
	for step := StepTitle; step <= StepContent; step++ {
		// Temporarily set current step to validate
		originalStep := state.CurrentStep
		state.CurrentStep = step
		stepErrors := c.ValidateStep(state)
		state.CurrentStep = originalStep

		errors = append(errors, stepErrors...)
	}

	return errors
}

// IsComplete checks if the wizard is complete
func (c *Creator) IsComplete(state *CreationState) bool {
	return state.CurrentStep == StepComplete
}

// Reset resets the wizard to initial state
func (c *Creator) Reset(state *CreationState) {
	state.CurrentStep = StepTitle
	state.Title = ""
	state.Description = ""
	state.Tags = []string{}
	state.Category = ""
	state.Content = ""
	state.ValidationErrors = []string{}
}

// GetStepCount returns the total number of steps
func (c *Creator) GetStepCount() int {
	return int(StepComplete) + 1
}

// GetCurrentStepNumber returns the current step number (1-indexed)
func (c *Creator) GetCurrentStepNumber(state *CreationState) int {
	return int(state.CurrentStep) + 1
}

// SavePrompt saves the prompt to storage
func (c *Creator) SavePrompt(data *PromptData) (string, error) {
	if c.storage == nil {
		return "", fmt.Errorf("storage not initialized")
	}

	// Save prompt to file
	filePath, err := c.storage.SavePrompt(data)
	if err != nil {
		return "", fmt.Errorf("failed to save prompt: %w", err)
	}

	return filePath, nil
}
