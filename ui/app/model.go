package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/internal/ai"
	"github.com/kyledavis/prompt-stack/ui/diffviewer"
	"github.com/kyledavis/prompt-stack/ui/suggestions"
	"github.com/kyledavis/prompt-stack/ui/workspace"
)

// Model represents the root application model
type Model struct {
	workspace       workspace.Model
	suggestions     suggestions.Model
	diffViewer      diffviewer.Model
	activePanel     string // "workspace", "suggestions", "diffviewer"
	width           int
	height          int
	pendingEdits    []ai.Edit // Store edits for diff viewer
	pendingOriginal string    // Store original content for diff viewer
}

// New creates a new application model
func New(workingDir string) Model {
	return Model{
		workspace:   workspace.New(workingDir),
		suggestions: suggestions.NewModel(),
		diffViewer:  diffviewer.NewModel(),
		activePanel: "workspace",
		width:       80,
		height:      24,
	}
}

// Init initializes the application model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keybindings
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.workspace.SetSize(msg.Width, msg.Height)
		m.suggestions.SetSize(msg.Width, msg.Height)
		m.diffViewer.SetSize(msg.Width, msg.Height)

	case diffviewer.ShowDiffMsg:
		// Show diff viewer modal
		m.activePanel = "diffviewer"
		m.diffViewer.SetDiff(msg.Diff)
		m.pendingOriginal = msg.Original
		m.pendingEdits = msg.Edits

		m.diffViewer.SetOnAccept(func() tea.Cmd {
			return func() tea.Msg {
				return diffviewer.AcceptDiffMsg{}
			}
		})
		m.diffViewer.SetOnReject(func() tea.Cmd {
			return func() tea.Msg {
				return diffviewer.RejectDiffMsg{}
			}
		})

	case diffviewer.HideDiffMsg:
		// Hide diff viewer modal
		m.activePanel = "workspace"

	case diffviewer.AcceptDiffMsg:
		// User accepted the diff - apply edits
		if m.pendingOriginal != "" && len(m.pendingEdits) > 0 {
			err := m.workspace.ApplyEditsAsSingleAction(m.pendingEdits)
			if err != nil {
				// Handle error - show in status bar
				m.workspace.SetStatus(fmt.Sprintf("Failed to apply edits: %v", err))
			}
		}
		// Unlock editor
		m.workspace.SetAIApplying(false)
		m.workspace.SetReadOnly(false)
		// Hide diff viewer
		m.activePanel = "workspace"
		// Clear pending data
		m.pendingEdits = nil
		m.pendingOriginal = ""

	case diffviewer.RejectDiffMsg:
		// User rejected the diff - discard changes
		// Unlock editor
		m.workspace.SetAIApplying(false)
		m.workspace.SetReadOnly(false)
		// Hide diff viewer
		m.activePanel = "workspace"
		// Clear pending data
		m.pendingEdits = nil
		m.pendingOriginal = ""
	}

	// Update active panel
	switch m.activePanel {
	case "workspace":
		updatedModel, updateCmd := m.workspace.Update(msg)
		m.workspace = updatedModel.(workspace.Model)
		cmd = updateCmd
	case "suggestions":
		m.suggestions, cmd = m.suggestions.Update(msg)
	case "diffviewer":
		m.diffViewer, cmd = m.diffViewer.Update(msg)
	}

	return m, cmd
}

// View renders the application
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	switch m.activePanel {
	case "workspace":
		return m.renderWorkspace()
	case "suggestions":
		return m.renderSuggestions()
	case "diffviewer":
		return m.renderDiffViewer()
	default:
		return m.renderWorkspace()
	}
}

// renderWorkspace renders the workspace view
func (m Model) renderWorkspace() string {
	// Render workspace
	workspaceView := m.workspace.View()

	// If suggestions panel is active, render split view
	if m.activePanel == "workspace" && len(m.suggestions.GetSuggestions()) > 0 {
		// Calculate split (70% workspace, 30% suggestions)
		workspaceWidth := m.width * 7 / 10
		suggestionsWidth := m.width - workspaceWidth

		// Style workspace
		workspaceStyle := lipgloss.NewStyle().
			Width(workspaceWidth).
			Height(m.height)

		// Style suggestions
		suggestionsStyle := lipgloss.NewStyle().
			Width(suggestionsWidth).
			Height(m.height)

		// Render suggestions panel
		m.suggestions.SetSize(suggestionsWidth, m.height)
		suggestionsView := suggestionsStyle.Render(m.suggestions.View())

		// Combine views
		return lipgloss.JoinHorizontal(lipgloss.Top,
			workspaceStyle.Render(workspaceView),
			suggestionsView,
		)
	}

	return workspaceView
}

// renderSuggestions renders the suggestions panel view
func (m Model) renderSuggestions() string {
	return m.suggestions.View()
}

// renderDiffViewer renders the diff viewer modal
func (m Model) renderDiffViewer() string {
	// Render diff viewer as modal overlay
	diffView := m.diffViewer.View()

	// Create modal style
	modalStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	return modalStyle.Render(diffView)
}

// SetSuggestions sets the suggestions to display
func (m *Model) SetSuggestions(suggestionList []ai.Suggestion) {
	m.suggestions.SetSuggestions(suggestionList)

	// Set up callbacks for applying/dismissing suggestions
	m.suggestions.SetOnApply(func(suggestion *ai.Suggestion) tea.Cmd {
		return func() tea.Msg {
			// When user presses 'a' to apply suggestion:
			// 1. Lock editor (read-only mode)
			m.workspace.SetAIApplying(true)
			m.workspace.SetReadOnly(true)

			// 2. Generate diff from suggestion edits
			diff, err := ai.NewDiffGenerator().GenerateUnifiedDiff(m.workspace.GetContent(), suggestion.ProposedChanges)
			if err != nil {
				// Handle error - show in status bar and unlock editor
				m.workspace.SetStatus(fmt.Sprintf("Failed to generate diff: %v", err))
				m.workspace.SetAIApplying(false)
				m.workspace.SetReadOnly(false)
				return nil
			}

			// 3. Store pending data
			m.pendingOriginal = m.workspace.GetContent()
			m.pendingEdits = suggestion.ProposedChanges

			// 4. Show diff viewer modal
			return diffviewer.ShowDiffMsg{
				Diff:     diff,
				Original: m.workspace.GetContent(),
				Edits:    suggestion.ProposedChanges,
			}
		}
	})

	m.suggestions.SetOnDismiss(func(suggestion *ai.Suggestion) tea.Cmd {
		return func() tea.Msg {
			// Mark suggestion as dismissed
			suggestion.MarkAsDismissed()
			return nil
		}
	})
}

// GetSuggestions returns the current suggestions
func (m *Model) GetSuggestions() []ai.Suggestion {
	return m.suggestions.GetSuggestions()
}

// ShowSuggestions shows the suggestions panel
func (m *Model) ShowSuggestions() {
	m.activePanel = "workspace" // Show split view
}

// HideSuggestions hides the suggestions panel
func (m *Model) HideSuggestions() {
	m.activePanel = "workspace"
}
