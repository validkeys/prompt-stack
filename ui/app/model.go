// Package app provides the root Bubble Tea model for the PromptStack TUI.
// It serves as the main application entry point, managing overall TUI state
// and coordinating child components like the status bar.
package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/kyledavis/prompt-stack/ui/statusbar"
	"github.com/kyledavis/prompt-stack/ui/theme"
)

const (
	defaultTitle   = "PromptStack TUI"
	defaultMessage = "Press 'q' or Ctrl+C to quit"
)

// Model represents the root application model state.
// It manages the overall TUI state and coordinates child components.
type Model struct {
	statusBar statusbar.Model
	width     int
	height    int
	quitting  bool
}

// New creates a new root app model with default values.
func New() Model {
	return Model{
		statusBar: statusbar.New(),
		width:     80, // Default width
		height:    24, // Default height
		quitting:  false,
	}
}

// Init initializes the root app model.
// Returns nil command as no initialization is needed.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages for the root app model.
// It processes keyboard input, window resize events, and quit commands.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle quit commands
		switch msg.Type {
		case tea.KeyCtrlC:
			newModel := m
			newModel.quitting = true
			return newModel, tea.Quit
		case tea.KeyRunes:
			// Check for 'q' key to quit
			if len(msg.Runes) == 1 && msg.Runes[0] == 'q' {
				newModel := m
				newModel.quitting = true
				return newModel, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		// Update window dimensions
		newModel := m
		newModel.width = msg.Width
		newModel.height = msg.Height
		// Update status bar with new dimensions
		var statusCmd tea.Cmd
		statusModel, statusCmd := newModel.statusBar.Update(msg)
		newModel.statusBar = statusModel.(statusbar.Model)
		return newModel, statusCmd
	}

	// Update status bar with any message
	var statusCmd tea.Cmd
	statusModel, statusCmd := m.statusBar.Update(msg)
	newModel := m
	newModel.statusBar = statusModel.(statusbar.Model)
	return newModel, statusCmd
}

// View renders the root app model as a string.
// It displays the main content area with the status bar at the bottom.
func (m Model) View() string {
	// Calculate available height for main content (minus status bar)
	mainHeight := m.height - 1
	if mainHeight < 0 {
		mainHeight = 0
	}

	// Render main content area with theme styling
	mainContent := theme.ModalStyle().
		Width(m.width).
		Height(mainHeight).
		Render(defaultTitle + "\n\n" + defaultMessage)

	// Render status bar at the bottom
	statusBar := m.statusBar.View()

	// Combine main content and status bar
	return mainContent + "\n" + statusBar
}

// IsQuitting returns whether the application is in the process of quitting.
func (m Model) IsQuitting() bool {
	return m.quitting
}

// GetWidth returns the current window width.
func (m Model) GetWidth() int {
	return m.width
}

// GetHeight returns the current window height.
func (m Model) GetHeight() int {
	return m.height
}

// GetStatusBar returns the status bar model.
func (m Model) GetStatusBar() statusbar.Model {
	return m.statusBar
}
