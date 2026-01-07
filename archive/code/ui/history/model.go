package history

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the history browser UI
type Model struct {
	list        list.Model
	searchInput textinput.Model
	showSearch  bool
	quitting    bool
	width       int
	height      int
	onSelect    func(string) tea.Cmd
	onDelete    func(string) tea.Cmd
}

// NewModel creates a new history browser model
func NewModel() Model {
	// Create list delegate
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	delegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	// Create list
	l := list.New(nil, delegate, 0, 0)
	l.Title = "History"
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)
	l.Styles.FilterPrompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))
	l.Styles.FilterCursor = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	// Create search input
	search := textinput.New()
	search.Placeholder = "Search history..."
	search.Prompt = "/"
	search.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))
	search.TextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))
	search.Width = 50

	return Model{
		list:        l,
		searchInput: search,
		showSearch:  false,
	}
}

// SetItems sets the list items
func (m *Model) SetItems(items []list.Item) {
	m.list.SetItems(items)
}

// SetOnSelect sets the callback for when an item is selected
func (m *Model) SetOnSelect(fn func(string) tea.Cmd) {
	m.onSelect = fn
}

// SetOnDelete sets the callback for when an item is deleted
func (m *Model) SetOnDelete(fn func(string) tea.Cmd) {
	m.onDelete = fn
}

// SetSize sets the size of the model
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Adjust list size based on whether search is shown
	listHeight := height
	if m.showSearch {
		listHeight = height - 3 // Reserve space for search input
	}

	m.list.SetSize(width, listHeight)
	m.searchInput.Width = width - 2 // Account for prompt
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.showSearch {
				// Exit search mode
				m.showSearch = false
				m.searchInput.Reset()
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit

		case "/":
			// Toggle search mode
			m.showSearch = !m.showSearch
			if m.showSearch {
				m.searchInput.Focus()
				return m, textinput.Blink
			}
			m.searchInput.Blur()
			m.searchInput.Reset()
			return m, nil

		case "enter":
			if m.showSearch {
				// Apply search filter
				m.showSearch = false
				m.searchInput.Blur()
				query := strings.TrimSpace(m.searchInput.Value())
				if query != "" {
					// Filter items manually
					var filteredItems []list.Item
					for _, item := range m.list.Items() {
						if strings.Contains(strings.ToLower(item.FilterValue()), strings.ToLower(query)) {
							filteredItems = append(filteredItems, item)
						}
					}
					m.list.SetItems(filteredItems)
				}
				return m, nil
			}
			// Select item
			if m.list.SelectedItem() != nil && m.onSelect != nil {
				if item, ok := m.list.SelectedItem().(Item); ok {
					return m, m.onSelect(item.FilePath())
				}
			}

		case "delete", "backspace":
			if !m.showSearch && m.list.SelectedItem() != nil && m.onDelete != nil {
				if item, ok := m.list.SelectedItem().(Item); ok {
					return m, m.onDelete(item.FilePath())
				}
			}
		}

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	// Update search input if in search mode
	if m.showSearch {
		m.searchInput, cmd = m.searchInput.Update(msg)
		return m, cmd
	}

	// Update list
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the model
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Render search input if active
	if m.showSearch {
		searchStyle := lipgloss.NewStyle().
			Width(m.width).
			Height(3).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205"))

		b.WriteString(searchStyle.Render(m.searchInput.View()))
		b.WriteString("\n")
	}

	// Render list
	b.WriteString(m.list.View())

	// Render help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(m.width)

	helpText := "↑/↓: navigate | Enter: load | /: search | Delete: delete | Esc: close"
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}

// Item represents a history composition in the list
type Item struct {
	filePath         string
	timestamp        string
	workingDirectory string
	preview          string
	charCount        int
	lineCount        int
}

// NewItem creates a new history item
func NewItem(filePath, timestamp, workingDir, preview string, charCount, lineCount int) Item {
	return Item{
		filePath:         filePath,
		timestamp:        timestamp,
		workingDirectory: workingDir,
		preview:          preview,
		charCount:        charCount,
		lineCount:        lineCount,
	}
}

// FilterValue implements list.Item
func (i Item) FilterValue() string {
	return i.timestamp + " " + i.workingDirectory + " " + i.preview
}

// Title implements list.Item
func (i Item) Title() string {
	return fmt.Sprintf("%s - %s", i.timestamp, i.workingDirectory)
}

// Description implements list.Item
func (i Item) Description() string {
	// Truncate preview if too long
	maxLen := 80
	preview := i.preview
	if len(preview) > maxLen {
		preview = preview[:maxLen] + "..."
	}
	return fmt.Sprintf("%s | %d chars, %d lines", preview, i.charCount, i.lineCount)
}

// FilePath returns the file path
func (i Item) FilePath() string {
	return i.filePath
}

// Timestamp returns the timestamp
func (i Item) Timestamp() string {
	return i.timestamp
}

// WorkingDirectory returns the working directory
func (i Item) WorkingDirectory() string {
	return i.workingDirectory
}

// Preview returns the preview text
func (i Item) Preview() string {
	return i.preview
}

// CharCount returns the character count
func (i Item) CharCount() int {
	return i.charCount
}

// LineCount returns the line count
func (i Item) LineCount() int {
	return i.lineCount
}
