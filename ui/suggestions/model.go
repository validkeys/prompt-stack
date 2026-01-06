package suggestions

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/internal/ai"
	"github.com/kyledavis/prompt-stack/ui/theme"
)

// Model represents the suggestions panel Bubble Tea model
type Model struct {
	list          list.Model
	viewport      viewport.Model
	suggestions   []ai.Suggestion
	selectedIndex int
	width         int
	height        int
	onApply       func(suggestion *ai.Suggestion) tea.Cmd
	onDismiss     func(suggestion *ai.Suggestion) tea.Cmd
}

// NewModel creates a new suggestions panel model
func NewModel() Model {
	// Initialize list with custom delegate
	delegate := suggestionDelegate{}

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.SetShowTitle(false)

	// Initialize viewport
	vp := viewport.New(0, 0)

	return Model{
		list:        l,
		viewport:    vp,
		suggestions: []ai.Suggestion{},
		width:       80,
		height:      20,
	}
}

// SetSuggestions sets the suggestions to display
func (m *Model) SetSuggestions(suggestions []ai.Suggestion) {
	m.suggestions = suggestions

	// Convert suggestions to list items
	items := make([]list.Item, len(suggestions))
	for i := range suggestions {
		items[i] = suggestionItem{
			suggestion: &suggestions[i],
			index:      i,
		}
	}

	m.list.SetItems(items)

	// Reset selection
	if len(items) > 0 {
		m.list.Select(0)
	}
}

// GetSelectedSuggestion returns the currently selected suggestion
func (m *Model) GetSelectedSuggestion() *ai.Suggestion {
	if len(m.suggestions) == 0 {
		return nil
	}

	selected := m.list.Index()
	if selected < 0 || selected >= len(m.suggestions) {
		return nil
	}

	return &m.suggestions[selected]
}

// GetSuggestions returns all suggestions
func (m *Model) GetSuggestions() []ai.Suggestion {
	return m.suggestions
}

// SetOnApply sets the callback for applying a suggestion
func (m *Model) SetOnApply(callback func(suggestion *ai.Suggestion) tea.Cmd) {
	m.onApply = callback
}

// SetOnDismiss sets the callback for dismissing a suggestion
func (m *Model) SetOnDismiss(callback func(suggestion *ai.Suggestion) tea.Cmd) {
	m.onDismiss = callback
}

// SetSize sets the size of the suggestions panel
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Update list size
	m.list.SetWidth(width)
	m.list.SetHeight(height - 2) // Reserve space for header/footer

	// Update viewport size
	m.viewport.Width = width
	m.viewport.Height = height - 2
}

// Update handles messages and updates the model state
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a", "A":
			// Apply selected suggestion
			if suggestion := m.GetSelectedSuggestion(); suggestion != nil && suggestion.IsApplicable() {
				if m.onApply != nil {
					return m, m.onApply(suggestion)
				}
			}

		case "d", "D":
			// Dismiss selected suggestion
			if suggestion := m.GetSelectedSuggestion(); suggestion != nil {
				if m.onDismiss != nil {
					return m, m.onDismiss(suggestion)
				}
			}

		case "up", "k":
			// Navigate up
			m.list.CursorUp()

		case "down", "j":
			// Navigate down
			m.list.CursorDown()

		case "home", "g":
			// Go to top
			m.list.Select(0)

		case "end", "G":
			// Go to bottom
			if len(m.suggestions) > 0 {
				m.list.Select(len(m.suggestions) - 1)
			}

		case "pgup":
			// Page up
			m.list.PrevPage()

		case "pgdown":
			// Page down
			m.list.NextPage()
		}

	case tea.WindowSizeMsg:
		// Handle window resize
		m.SetSize(msg.Width, msg.Height)
	}

	// Update list
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

// View renders the suggestions panel
func (m Model) View() string {
	if len(m.suggestions) == 0 {
		return m.renderEmpty()
	}

	// Render header
	header := m.renderHeader()

	// Render list
	listView := m.list.View()

	// Render footer
	footer := m.renderFooter()

	// Combine all parts
	return lipgloss.JoinVertical(lipgloss.Left, header, listView, footer)
}

// renderHeader renders the panel header
func (m Model) renderHeader() string {
	title := theme.HeaderStyle().Render("✨ AI Suggestions")

	count := fmt.Sprintf("%d suggestion%s", len(m.suggestions), pluralize(len(m.suggestions)))
	countStyle := theme.InfoStyle().Render(count)

	return lipgloss.JoinHorizontal(lipgloss.Left, title, "  ", countStyle)
}

// renderFooter renders the panel footer with help text
func (m Model) renderFooter() string {
	help := []string{
		"↑/↓ or j/k: Navigate",
		"a: Apply",
		"d: Dismiss",
		"Home/End: Jump",
	}

	helpText := strings.Join(help, "  |  ")

	return theme.StatusStyle().Render(helpText)
}

// renderEmpty renders the empty state
func (m Model) renderEmpty() string {
	emptyText := "No suggestions available"

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(theme.ListEmptyStyle().Render(emptyText))
}

// suggestionItem implements list.Item for suggestions
type suggestionItem struct {
	suggestion *ai.Suggestion
	index      int
}

// FilterValue implements list.Item
func (i suggestionItem) FilterValue() string {
	return i.suggestion.Title
}

// suggestionDelegate implements list.ItemDelegate for custom rendering
type suggestionDelegate struct{}

// Height implements list.ItemDelegate
func (d suggestionDelegate) Height() int {
	return 4
}

// Spacing implements list.ItemDelegate
func (d suggestionDelegate) Spacing() int {
	return 0
}

// Update implements list.ItemDelegate
func (d suggestionDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

// Render implements list.ItemDelegate
func (d suggestionDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	s, ok := item.(suggestionItem)
	if !ok {
		return
	}

	suggestion := s.suggestion
	isSelected := index == m.Index()

	// Build suggestion display
	var lines []string

	// Line 1: Icon + Title + Status
	titleLine := d.renderTitleLine(suggestion, isSelected)
	lines = append(lines, titleLine)

	// Line 2: Description (truncated)
	description := d.renderDescription(suggestion, isSelected)
	lines = append(lines, description)

	// Line 3: Type + Changes count
	metaLine := d.renderMetaLine(suggestion, isSelected)
	lines = append(lines, metaLine)

	// Line 4: Empty for spacing
	lines = append(lines, "")

	// Render all lines
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
}

// renderTitleLine renders the title line with icon and status
func (d suggestionDelegate) renderTitleLine(suggestion *ai.Suggestion, isSelected bool) string {
	icon := suggestion.GetTypeIcon()
	title := suggestion.Title

	// Get style based on selection
	var style lipgloss.Style
	if isSelected {
		style = theme.ListItemSelectedStyle()
	} else {
		style = theme.ListItemStyle()
	}

	// Add status indicator
	status := d.renderStatus(suggestion)

	return style.Render(fmt.Sprintf("%s %s %s", icon, title, status))
}

// renderDescription renders the description (truncated)
func (d suggestionDelegate) renderDescription(suggestion *ai.Suggestion, isSelected bool) string {
	description := suggestion.Description

	// Truncate if too long
	maxLen := 60
	if len(description) > maxLen {
		description = description[:maxLen] + "..."
	}

	// Get style based on selection
	var style lipgloss.Style
	if isSelected {
		style = theme.ListItemSelectedStyle()
	} else {
		style = theme.ListDescriptionStyle()
	}

	return style.Render(description)
}

// renderMetaLine renders metadata (type, changes count)
func (d suggestionDelegate) renderMetaLine(suggestion *ai.Suggestion, isSelected bool) string {
	typeLabel := suggestion.GetTypeLabel()
	changesCount := len(suggestion.ProposedChanges)

	meta := fmt.Sprintf("%s • %d change%s", typeLabel, changesCount, pluralize(changesCount))

	// Get style based on selection
	var style lipgloss.Style
	if isSelected {
		style = theme.ListItemSelectedStyle()
	} else {
		style = theme.ListCategoryStyle()
	}

	return style.Render(meta)
}

// renderStatus renders the status indicator
func (d suggestionDelegate) renderStatus(suggestion *ai.Suggestion) string {
	switch suggestion.Status {
	case ai.SuggestionStatusPending:
		return ""
	case ai.SuggestionStatusApplying:
		return theme.InfoStyle().Render("[Applying...]")
	case ai.SuggestionStatusApplied:
		return theme.SuccessStyle().Render("[Applied ✓]")
	case ai.SuggestionStatusDismissed:
		return theme.StatusStyle().Render("[Dismissed]")
	case ai.SuggestionStatusError:
		return theme.ErrorStyle().Render("[Error]")
	default:
		return ""
	}
}

// pluralize returns the plural form of a word based on count
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
