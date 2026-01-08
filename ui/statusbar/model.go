// Package statusbar provides a status bar component for the PromptStack TUI.
// It displays application status information including character and line counts,
// and integrates with the theme system for consistent styling.
package statusbar

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kyledavis/prompt-stack/ui/theme"
)

// Model represents the status bar component state.
type Model struct {
	charCount int
	lineCount int
	width     int
}

// New creates a new status bar model with default values.
func New() Model {
	return Model{
		charCount: 0,
		lineCount: 0,
		width:     80, // Default width
	}
}

// Init initializes the status bar model.
// Returns nil command as no initialization is needed.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages for the status bar.
// It processes window resize messages to adjust the status bar width.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	return m, nil
}

// View renders the status bar as a string.
// It displays character and line counts using the theme's status style.
func (m Model) View() string {
	statusText := fmt.Sprintf("Chars: %d | Lines: %d", m.charCount, m.lineCount)
	return theme.StatusStyle().
		Width(m.width).
		Render(statusText)
}

// SetCharCount updates the character count displayed in the status bar.
// Negative values are clamped to 0.
func (m *Model) SetCharCount(count int) {
	if count < 0 {
		count = 0
	}
	m.charCount = count
}

// SetLineCount updates the line count displayed in the status bar.
// Negative values are clamped to 0.
func (m *Model) SetLineCount(count int) {
	if count < 0 {
		count = 0
	}
	m.lineCount = count
}

// GetCharCount returns the current character count.
func (m Model) GetCharCount() int {
	return m.charCount
}

// GetLineCount returns the current line count.
func (m Model) GetLineCount() int {
	return m.lineCount
}

// GetWidth returns the current status bar width.
func (m Model) GetWidth() int {
	return m.width
}
