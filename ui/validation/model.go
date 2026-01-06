package validation

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/internal/prompt"
	"github.com/kyledavis/prompt-stack/ui/theme"
)

// Model represents the validation results modal
type Model struct {
	results  map[string]prompt.ValidationResult // keyed by file path
	prompts  map[string]*prompt.Prompt          // keyed by file path
	selected int                                // index of selected prompt
	scroll   int                                // scroll position
	width    int                                // modal width
	height   int                                // modal height
	visible  bool                               // whether modal is visible
	vimMode  bool                               // vim mode enabled
	onClose  func()                             // callback when modal closes
}

// New creates a new validation results model
func New(prompts map[string]*prompt.Prompt, vimMode bool) Model {
	return Model{
		results:  make(map[string]prompt.ValidationResult),
		prompts:  prompts,
		selected: 0,
		scroll:   0,
		visible:  false,
		vimMode:  vimMode,
	}
}

// Init initializes the validation model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.visible = false
			if m.onClose != nil {
				m.onClose()
			}
			return m, nil

		case tea.KeyUp:
			m.moveSelection(-1)

		case tea.KeyDown:
			m.moveSelection(1)

		case tea.KeyPgUp:
			m.moveSelection(-10)

		case tea.KeyPgDown:
			m.moveSelection(10)

		case tea.KeyHome:
			m.selected = 0
			m.scroll = 0

		case tea.KeyEnd:
			m.selected = len(m.prompts) - 1
			m.adjustScroll()
		}

		// Vim mode keybindings
		if m.vimMode {
			switch msg.String() {
			case "j":
				m.moveSelection(1)
			case "k":
				m.moveSelection(-1)
			case "G":
				m.selected = len(m.prompts) - 1
				m.adjustScroll()
			case "g":
				m.selected = 0
				m.scroll = 0
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the validation results modal
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	// Calculate modal dimensions
	modalWidth := min(m.width-4, 100)
	modalHeight := min(m.height-4, 30)

	// Create modal style
	modalStyle := theme.BaseModal().
		Width(modalWidth).
		Height(modalHeight)

	// Render header
	header := m.renderHeader()

	// Render summary
	summary := m.renderSummary(modalWidth)

	// Render validation list
	list := m.renderValidationList(modalWidth, modalHeight-4)

	// Combine sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		summary,
		list,
	)

	return modalStyle.Render(content)
}

// renderHeader renders the modal header
func (m Model) renderHeader() string {
	return theme.HeaderStyle().Render("Validation Results")
}

// renderSummary renders the validation summary
func (m Model) renderSummary(width int) string {
	errorCount := 0
	warningCount := 0
	totalPrompts := len(m.prompts)

	for _, result := range m.results {
		errorCount += len(result.Errors)
		warningCount += len(result.Warnings)
	}

	var summaryParts []string
	summaryParts = append(summaryParts, fmt.Sprintf("Total: %d", totalPrompts))

	if errorCount > 0 {
		summaryParts = append(summaryParts,
			theme.ValidationErrorStyle().Render(fmt.Sprintf("Errors: %d", errorCount)))
	}

	if warningCount > 0 {
		summaryParts = append(summaryParts,
			theme.ValidationWarningStyle().Render(fmt.Sprintf("Warnings: %d", warningCount)))
	}

	if errorCount == 0 && warningCount == 0 {
		summaryParts = append(summaryParts,
			theme.SuccessStyle().Render("All prompts valid ✓"))
	}

	return theme.ValidationSummaryStyle().Render(strings.Join(summaryParts, " | "))
}

// renderValidationList renders the list of validation results
func (m Model) renderValidationList(width, height int) string {
	listStyle := lipgloss.NewStyle().
		Width(width - 2).
		Height(height)

	// Get sorted prompt paths
	var paths []string
	for path := range m.prompts {
		paths = append(paths, path)
	}

	// Filter to only show prompts with issues
	var issuePaths []string
	for _, path := range paths {
		result := m.results[path]
		if len(result.Errors) > 0 || len(result.Warnings) > 0 {
			issuePaths = append(issuePaths, path)
		}
	}

	// If no issues, show success message
	if len(issuePaths) == 0 {
		return listStyle.Render(
			theme.ListEmptyStyle().Render("No validation issues found"))
	}

	// Render each prompt with issues
	var items []string
	for i, path := range issuePaths {
		if i < m.scroll || i >= m.scroll+height {
			continue
		}

		p := m.prompts[path]
		result := m.results[path]

		// Style based on selection
		var itemStyle lipgloss.Style
		if i == m.selected {
			itemStyle = theme.ListItemSelectedStyle()
		} else {
			itemStyle = theme.ListItemStyle()
		}

		// Determine icon and style based on validation status
		var icon string
		var promptStyle lipgloss.Style
		if len(result.Errors) > 0 {
			icon = theme.ValidationErrorStyle().Render("✗")
			promptStyle = theme.ValidationPromptErrorStyle()
		} else {
			icon = theme.ValidationWarningStyle().Render("⚠")
			promptStyle = theme.ValidationPromptWarningStyle()
		}

		// Render prompt title
		titleText := fmt.Sprintf("%s %s", icon, promptStyle.Render(p.Title))

		// Render validation details
		var details []string
		for _, err := range result.Errors {
			details = append(details,
				theme.ValidationErrorStyle().Render(fmt.Sprintf("  ✗ %s", err.Message)))
		}
		for _, warn := range result.Warnings {
			details = append(details,
				theme.ValidationWarningStyle().Render(fmt.Sprintf("  ⚠ %s", warn.Message)))
		}

		// Combine title and details
		itemContent := titleText
		if len(details) > 0 {
			itemContent += "\n" + strings.Join(details, "\n")
		}

		items = append(items, itemStyle.Render(itemContent))
	}

	// Add help text at bottom
	helpText := theme.ListDescriptionStyle().Render(
		"↑/↓ or j/k: Navigate | Esc: Close")
	items = append(items, helpText)

	return listStyle.Render(strings.Join(items, "\n"))
}

// moveSelection moves the selection up or down
func (m *Model) moveSelection(delta int) {
	// Get issue paths count
	issueCount := 0
	for _, result := range m.results {
		if len(result.Errors) > 0 || len(result.Warnings) > 0 {
			issueCount++
		}
	}

	if issueCount == 0 {
		return
	}

	m.selected += delta

	// Clamp selection
	if m.selected < 0 {
		m.selected = 0
	}
	if m.selected >= issueCount {
		m.selected = issueCount - 1
	}

	m.adjustScroll()
}

// adjustScroll adjusts the scroll position to keep selection visible
func (m *Model) adjustScroll() {
	// Calculate visible height (approximate)
	visibleHeight := m.height - 6 // Account for header, summary, padding

	if m.selected < m.scroll {
		m.scroll = m.selected
	} else if m.selected >= m.scroll+visibleHeight {
		m.scroll = m.selected - visibleHeight + 1
	}
}

// SetResults sets the validation results to display
func (m *Model) SetResults(results map[string]prompt.ValidationResult) {
	m.results = results
	m.selected = 0
	m.scroll = 0
}

// Show makes the modal visible
func (m *Model) Show() {
	m.visible = true
}

// Hide makes the modal invisible
func (m *Model) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is visible
func (m Model) IsVisible() bool {
	return m.visible
}

// SetOnClose sets the callback for when the modal closes
func (m *Model) SetOnClose(onClose func()) {
	m.onClose = onClose
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
