package common

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmationModel represents a confirmation dialog
type ConfirmationModel struct {
	title       string
	message     string
	confirmText string
	input       textinput.Model
	confirmed   bool
	cancelled   bool
	width       int
	height      int
	onConfirm   func() tea.Cmd
	onCancel    func() tea.Cmd
}

// NewConfirmation creates a new confirmation dialog
func NewConfirmation(title, message, confirmText string) ConfirmationModel {
	// Create input for typing confirmation
	input := textinput.New()
	input.Placeholder = "Type to confirm"
	input.Prompt = "> "
	input.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	input.Width = 30

	return ConfirmationModel{
		title:       title,
		message:     message,
		confirmText: confirmText,
		input:       input,
		confirmed:   false,
		cancelled:   false,
	}
}

// SetOnConfirm sets the callback for confirmation
func (m *ConfirmationModel) SetOnConfirm(fn func() tea.Cmd) {
	m.onConfirm = fn
}

// SetOnCancel sets the callback for cancellation
func (m *ConfirmationModel) SetOnCancel(fn func() tea.Cmd) {
	m.onCancel = fn
}

// SetSize sets the size of the dialog
func (m *ConfirmationModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.input.Width = width - 20 // Account for prompt and padding
}

// Init initializes the model
func (m ConfirmationModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m ConfirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			if m.onCancel != nil {
				return m, m.onCancel()
			}
			return m, tea.Quit

		case "enter":
			// Check if input matches confirmation text
			if strings.TrimSpace(m.input.Value()) == m.confirmText {
				m.confirmed = true
				if m.onConfirm != nil {
					return m, m.onConfirm()
				}
			}
			return m, nil
		}
	}

	// Update input
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View renders the confirmation dialog
func (m ConfirmationModel) View() string {
	if m.confirmed || m.cancelled {
		return ""
	}

	var b strings.Builder

	// Dialog box style
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("226")). // Yellow for warning
		Width(m.width).
		Height(m.height).
		Padding(1, 1)

	// Title style
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")). // Yellow
		Bold(true).
		Width(m.width - 4).
		Align(lipgloss.Center)

	// Message style
	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Width(m.width - 4).
		Align(lipgloss.Center)

	// Input style
	inputStyle := lipgloss.NewStyle().
		Width(m.width - 4).
		Align(lipgloss.Center)

	// Warning icon
	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")) // Yellow

	// Render title
	b.WriteString(titleStyle.Render(m.title))
	b.WriteString("\n\n")

	// Render warning icon
	b.WriteString(warningStyle.Render("⚠️  "))
	b.WriteString("\n\n")

	// Render message
	b.WriteString(messageStyle.Render(m.message))
	b.WriteString("\n\n")

	// Render confirmation instruction
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(m.width - 4).
		Align(lipgloss.Center)

	b.WriteString(instructionStyle.Render(fmt.Sprintf("Type '%s' to confirm:", m.confirmText)))
	b.WriteString("\n\n")

	// Render input
	b.WriteString(inputStyle.Render(m.input.View()))
	b.WriteString("\n\n")

	// Render help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(m.width - 4).
		Align(lipgloss.Center)

	helpText := "Enter: confirm | Esc: cancel"
	b.WriteString(helpStyle.Render(helpText))

	return dialogStyle.Render(b.String())
}

// IsConfirmed returns true if user confirmed
func (m ConfirmationModel) IsConfirmed() bool {
	return m.confirmed
}

// IsCancelled returns true if user cancelled
func (m ConfirmationModel) IsCancelled() bool {
	return m.cancelled
}

// SimpleConfirmation creates a simple yes/no confirmation dialog
func SimpleConfirmation(title, message string) ConfirmationModel {
	return NewConfirmation(title, message, "yes")
}

// DestructiveConfirmation creates a confirmation dialog for destructive operations
func DestructiveConfirmation(title, message string) ConfirmationModel {
	return NewConfirmation(title, message, "DELETE")
}

// Reset resets the confirmation dialog
func (m *ConfirmationModel) Reset() {
	m.input.Reset()
	m.confirmed = false
	m.cancelled = false
}
