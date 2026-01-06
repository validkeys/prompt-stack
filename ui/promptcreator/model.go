package promptcreator

import (
	"fmt"
	"strings"

	"github.com/kyledavis/prompt-stack/internal/library"
	"github.com/kyledavis/prompt-stack/internal/prompt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the prompt creation wizard UI
type Model struct {
	state            *prompt.CreationState
	titleInput       textinput.Model
	descriptionInput textinput.Model
	tagsInput        textinput.Model
	categoryInput    textinput.Model
	contentInput     textinput.Model
	quitting         bool
	width            int
	height           int
	onSave           func(*prompt.PromptData) tea.Cmd
	onCancel         func() tea.Cmd
	storage          *prompt.Storage
	lib              *library.Library
}

// NewModel creates a new prompt creator model
func NewModel(storage *prompt.Storage, lib *library.Library) Model {
	// Create title input
	title := textinput.New()
	title.Placeholder = "Enter prompt title..."
	title.Prompt = "Title: "
	title.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	title.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	title.Width = 60
	title.Focus()

	// Create description input
	description := textinput.New()
	description.Placeholder = "Enter description (optional)..."
	description.Prompt = "Description: "
	description.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	description.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	description.Width = 60

	// Create tags input
	tags := textinput.New()
	tags.Placeholder = "Enter tags (comma-separated, optional)..."
	tags.Prompt = "Tags: "
	tags.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	tags.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	tags.Width = 60

	// Create category input
	category := textinput.New()
	category.Placeholder = "Select category..."
	category.Prompt = "Category: "
	category.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	category.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	category.Width = 60

	// Create content input
	content := textinput.New()
	content.Placeholder = "Enter prompt content (markdown)..."
	content.Prompt = "Content: "
	content.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	content.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	content.Width = 60

	return Model{
		state:            prompt.NewCreationState(),
		titleInput:       title,
		descriptionInput: description,
		tagsInput:        tags,
		categoryInput:    category,
		contentInput:     content,
		storage:          storage,
		lib:              lib,
	}
}

// SetOnSave sets the callback for saving
func (m *Model) SetOnSave(fn func(*prompt.PromptData) tea.Cmd) {
	m.onSave = fn
}

// SetOnCancel sets the callback for cancellation
func (m *Model) SetOnCancel(fn func() tea.Cmd) {
	m.onCancel = fn
}

// SetSize sets the size of the model
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Adjust input widths
	inputWidth := width - 20
	m.titleInput.Width = inputWidth
	m.descriptionInput.Width = inputWidth
	m.tagsInput.Width = inputWidth
	m.categoryInput.Width = inputWidth
	m.contentInput.Width = inputWidth
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			if m.onCancel != nil {
				return m, m.onCancel()
			}
			return m, tea.Quit

		case "enter":
			// Validate current step
			creator := prompt.NewCreator(nil)
			errors := creator.ValidateStep(m.state)

			if len(errors) > 0 {
				// Show errors, don't proceed
				m.state.ValidationErrors = errors
				return m, nil
			}

			// Save current input to state
			m.saveCurrentInput()

			// Move to next step or complete
			if creator.CanProceed(m.state) {
				if creator.IsComplete(m.state) {
					// Complete wizard - save prompt
					promptData := creator.GetPromptData(m.state)

					// Save prompt to storage
					creatorWithStorage := prompt.NewCreatorWithStorage(nil, m.storage)
					filePath, err := creatorWithStorage.SavePrompt(promptData)

					if err != nil {
						// Show error and stay in complete step
						m.state.ValidationErrors = []string{fmt.Sprintf("Failed to save prompt: %v", err)}
						return m, nil
					}

					// Add to library index
					if m.lib != nil {
						newPrompt := &prompt.Prompt{
							ID:           filePath,
							Title:        promptData.Title,
							Description:  promptData.Description,
							Tags:         promptData.Tags,
							Category:     promptData.Category,
							FilePath:     filePath,
							Content:      promptData.Content,
							Metadata:     make(map[string]string),
							Placeholders: prompt.ParsePlaceholders(promptData.Content),
							ValidationStatus: prompt.ValidationResult{
								IsValid: true,
							},
							UsageStats: prompt.UsageMetadata{
								UseCount: 0,
							},
						}
						m.lib.AddPrompt(newPrompt)
					}

					// Call save callback if provided
					if m.onSave != nil {
						return m, m.onSave(promptData)
					}
				} else {
					// Go to next step
					creator.GoToNext(m.state)
					m.focusCurrentInput()
				}
			}
			return m, nil

		case "tab":
			// Move to next step
			creator := prompt.NewCreator(nil)
			if creator.CanProceed(m.state) {
				creator.GoToNext(m.state)
				m.focusCurrentInput()
			}
			return m, nil

		case "shift+tab":
			// Move to previous step
			creator := prompt.NewCreator(nil)
			if creator.CanGoBack(m.state) {
				creator.GoToPrevious(m.state)
				m.focusCurrentInput()
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	// Update current input based on step
	switch m.state.CurrentStep {
	case prompt.StepTitle:
		m.titleInput, cmd = m.titleInput.Update(msg)
	case prompt.StepDescription:
		m.descriptionInput, cmd = m.descriptionInput.Update(msg)
	case prompt.StepTags:
		m.tagsInput, cmd = m.tagsInput.Update(msg)
	case prompt.StepCategory:
		m.categoryInput, cmd = m.categoryInput.Update(msg)
	case prompt.StepContent:
		m.contentInput, cmd = m.contentInput.Update(msg)
	case prompt.StepReview, prompt.StepComplete:
		// No input in review/complete steps
	}

	return m, cmd
}

// saveCurrentInput saves current input to state
func (m *Model) saveCurrentInput() {
	m.state.Title = strings.TrimSpace(m.titleInput.Value())
	m.state.Description = strings.TrimSpace(m.descriptionInput.Value())

	// Parse tags
	tagsStr := strings.TrimSpace(m.tagsInput.Value())
	if tagsStr != "" {
		tags := strings.Split(tagsStr, ",")
		m.state.Tags = make([]string, 0, len(tags))
		for i, tag := range tags {
			m.state.Tags[i] = strings.TrimSpace(tag)
		}
	}

	m.state.Category = strings.TrimSpace(m.categoryInput.Value())
	m.state.Content = strings.TrimSpace(m.contentInput.Value())
}

// focusCurrentInput focuses the appropriate input for current step
func (m *Model) focusCurrentInput() {
	// Blur all inputs first
	m.titleInput.Blur()
	m.descriptionInput.Blur()
	m.tagsInput.Blur()
	m.categoryInput.Blur()
	m.contentInput.Blur()

	// Focus current input
	switch m.state.CurrentStep {
	case prompt.StepTitle:
		m.titleInput.Focus()
	case prompt.StepDescription:
		m.descriptionInput.Focus()
	case prompt.StepTags:
		m.tagsInput.Focus()
	case prompt.StepCategory:
		m.categoryInput.Focus()
	case prompt.StepContent:
		m.contentInput.Focus()
	case prompt.StepReview, prompt.StepComplete:
		// No input to focus
	}
}

// View renders the model
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Render header
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Width(m.width).
		Align(lipgloss.Center)

	b.WriteString(headerStyle.Render("Create New Prompt"))
	b.WriteString("\n\n")

	// Render progress
	m.renderProgress(&b)

	// Render current step
	m.renderCurrentStep(&b)

	// Render validation errors
	if len(m.state.ValidationErrors) > 0 {
		m.renderValidationErrors(&b)
	}

	// Render help text
	m.renderHelpText(&b)

	return b.String()
}

// renderProgress renders the progress indicator
func (m *Model) renderProgress(b *strings.Builder) {
	creator := prompt.NewCreator(nil)
	progress := creator.GetProgress(m.state)
	stepCount := creator.GetStepCount()
	currentStep := creator.GetCurrentStepNumber(m.state)

	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(m.width)

	progressText := fmt.Sprintf("Step %d of %d (%d%%)", currentStep, stepCount, progress)
	b.WriteString(progressStyle.Render(progressText))
	b.WriteString("\n\n")
}

// renderCurrentStep renders the current step
func (m *Model) renderCurrentStep(b *strings.Builder) {
	creator := prompt.NewCreator(nil)
	stepInfo := creator.GetStepInfo(m.state)

	// Step title style
	stepTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Width(m.width)

	b.WriteString(stepTitleStyle.Render(fmt.Sprintf("%s %s", stepInfo.Title, m.getRequiredIndicator(stepInfo.IsRequired))))
	b.WriteString("\n")

	// Step description
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(m.width)

	b.WriteString(descStyle.Render(stepInfo.Description))
	b.WriteString("\n\n")

	// Render input for current step
	m.renderStepInput(b)
}

// renderStepInput renders the input for the current step
func (m *Model) renderStepInput(b *strings.Builder) {
	switch m.state.CurrentStep {
	case prompt.StepTitle:
		b.WriteString(m.titleInput.View())
		b.WriteString("\n")
	case prompt.StepDescription:
		b.WriteString(m.descriptionInput.View())
		b.WriteString("\n")
	case prompt.StepTags:
		b.WriteString(m.tagsInput.View())
		b.WriteString("\n")
	case prompt.StepCategory:
		b.WriteString(m.categoryInput.View())
		b.WriteString("\n")
	case prompt.StepContent:
		b.WriteString(m.contentInput.View())
		b.WriteString("\n")
	case prompt.StepReview:
		m.renderReview(b)
	case prompt.StepComplete:
		completeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Width(m.width).
			Align(lipgloss.Center)

		b.WriteString(completeStyle.Render("✓ Prompt created successfully!"))
		b.WriteString("\n\n")
	}
}

// renderReview renders the review step
func (m *Model) renderReview(b *strings.Builder) {
	creator := prompt.NewCreator(nil)
	promptData := creator.GetPromptData(m.state)

	// Review style
	reviewStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Width(m.width)

	// Title
	b.WriteString(reviewStyle.Render(fmt.Sprintf("Title: %s", promptData.Title)))
	b.WriteString("\n")

	// Description
	if promptData.Description != "" {
		b.WriteString(reviewStyle.Render(fmt.Sprintf("Description: %s", promptData.Description)))
		b.WriteString("\n")
	}

	// Tags
	if len(promptData.Tags) > 0 {
		b.WriteString(reviewStyle.Render(fmt.Sprintf("Tags: %s", strings.Join(promptData.Tags, ", "))))
		b.WriteString("\n")
	}

	// Category
	b.WriteString(reviewStyle.Render(fmt.Sprintf("Category: %s", promptData.Category)))
	b.WriteString("\n")

	// Content preview
	contentPreview := promptData.Content
	if len(contentPreview) > 200 {
		contentPreview = contentPreview[:200] + "..."
	}
	b.WriteString(reviewStyle.Render(fmt.Sprintf("Content: %s", contentPreview)))
	b.WriteString("\n\n")

	// Filename
	filename := creator.GenerateFilename(promptData.Title)
	b.WriteString(reviewStyle.Render(fmt.Sprintf("Filename: %s", filename)))
	b.WriteString("\n")
}

// renderValidationErrors renders validation errors
func (m *Model) renderValidationErrors(b *strings.Builder) {
	if len(m.state.ValidationErrors) == 0 {
		return
	}

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")). // Yellow
		Width(m.width)

	b.WriteString(errorStyle.Render("⚠️  Validation Errors:"))
	b.WriteString("\n\n")

	for _, err := range m.state.ValidationErrors {
		b.WriteString(errorStyle.Render(fmt.Sprintf("  • %s", err)))
		b.WriteString("\n")
	}
	b.WriteString("\n")
}

// renderHelpText renders the help text
func (m *Model) renderHelpText(b *strings.Builder) {
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(m.width)

	helpText := "Enter: next | Tab: next | Shift+Tab: previous | Esc: cancel"
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(helpText))
}

// getRequiredIndicator returns a required indicator
func (m *Model) getRequiredIndicator(isRequired bool) string {
	if isRequired {
		return "*"
	}
	return ""
}

// GetState returns the current creation state
func (m *Model) GetState() *prompt.CreationState {
	return m.state
}

// Reset resets the model to initial state
func (m *Model) Reset() {
	m.state = prompt.NewCreationState()
	m.titleInput.Reset()
	m.descriptionInput.Reset()
	m.tagsInput.Reset()
	m.categoryInput.Reset()
	m.contentInput.Reset()
	m.titleInput.Focus()
	m.quitting = false
}
