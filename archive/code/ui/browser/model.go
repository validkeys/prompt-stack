package browser

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/internal/prompt"
	"github.com/kyledavis/prompt-stack/ui/theme"
	"github.com/sahilm/fuzzy"
)

// Model represents the library browser modal
type Model struct {
	prompts     map[string]*prompt.Prompt // keyed by file path
	filtered    []string                  // filtered prompt file paths
	selected    int                       // index of selected prompt
	searchInput string                    // current search query
	width       int                       // modal width
	height      int                       // modal height
	visible     bool                      // whether browser is visible
	insertMode  InsertMode                // insert at cursor or new line
	vimMode     bool                      // vim mode enabled
}

// InsertMode determines how prompts are inserted
type InsertMode int

const (
	InsertAtCursor InsertMode = iota
	InsertOnNewLine
)

// New creates a new library browser model
func New(prompts map[string]*prompt.Prompt, vimMode bool) Model {
	// Convert prompts to slice for initial display
	var allPaths []string
	for path := range prompts {
		allPaths = append(allPaths, path)
	}

	return Model{
		prompts:     prompts,
		filtered:    allPaths,
		selected:    0,
		searchInput: "",
		visible:     false,
		insertMode:  InsertAtCursor,
		vimMode:     vimMode,
	}
}

// Init initializes the browser model
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
			return m, nil

		case tea.KeyEnter:
			// Insert selected prompt
			if len(m.filtered) > 0 {
				p := m.prompts[m.filtered[m.selected]]
				if p == nil {
					return m, nil
				}

				// Check validation status
				if !p.ValidationStatus.IsValid {
					// Block insertion of error-level prompts
					return m, func() tea.Msg {
						return ValidationErrorMsg{
							FilePath: m.filtered[m.selected],
							Errors:   p.ValidationStatus.Errors,
						}
					}
				}

				// Allow insertion of warning-level prompts with indicator
				insertMode := InsertAtCursor
				if msg.Alt {
					insertMode = InsertOnNewLine
				}

				return m, func() tea.Msg {
					return InsertPromptMsg{
						FilePath:    m.filtered[m.selected],
						InsertMode:  insertMode,
						HasWarnings: len(p.ValidationStatus.Warnings) > 0,
					}
				}
			}

		case tea.KeyUp:
			m.moveSelection(-1)

		case tea.KeyDown:
			m.moveSelection(1)

		case tea.KeyLeft:
			// Move to start of search input
			// Not implemented for now

		case tea.KeyRight:
			// Move to end of search input
			// Not implemented for now

		case tea.KeyBackspace:
			if len(m.searchInput) > 0 {
				m.searchInput = m.searchInput[:len(m.searchInput)-1]
				m.applyFilter()
			}

		default:
			if msg.Type == tea.KeyRunes {
				// Add character to search input
				m.searchInput += string(msg.Runes)
				m.applyFilter()
			}
		}

		// Vim mode keybindings
		if m.vimMode {
			switch msg.String() {
			case "j":
				m.moveSelection(1)
			case "k":
				m.moveSelection(-1)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the library browser
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

	// Render search input
	searchInput := m.renderSearchInput(modalWidth)

	// Render prompt list
	promptList := m.renderPromptList(modalWidth, modalHeight-4)

	// Render preview
	preview := m.renderPreview(modalWidth, modalHeight-4)

	// Combine sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		searchInput,
		lipgloss.JoinHorizontal(lipgloss.Top, promptList, preview),
	)

	return modalStyle.Render(content)
}

// renderHeader renders the browser header
func (m Model) renderHeader() string {
	return theme.HeaderStyle().Render("Library Browser")
}

// renderSearchInput renders the search input field
func (m Model) renderSearchInput(width int) string {
	promptText := "Search: "
	if m.searchInput != "" {
		promptText = "Search: " + m.searchInput
	}

	return theme.SearchInputStyle().
		Width(width - 2).
		Render(promptText)
}

// renderPromptList renders the filtered prompt list
func (m Model) renderPromptList(width, height int) string {
	listWidth := width / 2
	listHeight := height

	listStyle := lipgloss.NewStyle().
		Width(listWidth).
		Height(listHeight).
		Padding(0, 1)

	// Render each prompt
	var items []string
	for i, filePath := range m.filtered {
		p := m.prompts[filePath]
		if p == nil {
			continue
		}

		// Style based on selection
		var itemStyle lipgloss.Style
		if i == m.selected {
			itemStyle = theme.ListItemSelectedStyle()
		} else {
			itemStyle = theme.ListItemStyle()
		}

		// Render with category label and validation icon
		categoryLabel := theme.ListCategoryStyle().Render(fmt.Sprintf("[%s]", p.Category))

		// Add validation icon if prompt has issues
		var validationIcon string
		if !p.ValidationStatus.IsValid {
			validationIcon = theme.ValidationErrorStyle().Render("✗ ")
		} else if len(p.ValidationStatus.Warnings) > 0 {
			validationIcon = theme.ValidationWarningStyle().Render("⚠ ")
		}

		titleText := fmt.Sprintf("%s%s%s", validationIcon, categoryLabel, p.Title)

		items = append(items, itemStyle.Render(titleText))
	}

	// Add empty state if no results
	if len(items) == 0 {
		items = append(items, theme.ListEmptyStyle().Render("No prompts found"))
	}

	// Truncate to fit height
	if len(items) > listHeight {
		items = items[:listHeight]
	}

	return listStyle.Render(strings.Join(items, "\n"))
}

// renderPreview renders the preview pane
func (m Model) renderPreview(width, height int) string {
	previewWidth := width / 2
	previewHeight := height

	previewStyle := theme.PreviewStyle().
		Width(previewWidth).
		Height(previewHeight)

	// Get selected prompt
	if len(m.filtered) == 0 {
		return previewStyle.Render(theme.ListEmptyStyle().Render("Select a prompt to preview"))
	}

	p := m.prompts[m.filtered[m.selected]]
	if p == nil {
		return previewStyle.Render("")
	}

	// Render preview content
	var previewLines []string

	// Title
	previewLines = append(previewLines, theme.PreviewTitleStyle().Render(p.Title))

	// Description
	if p.Description != "" {
		previewLines = append(previewLines, "")
		previewLines = append(previewLines, theme.PreviewDescriptionStyle().Render(p.Description))
	}

	// Tags
	if len(p.Tags) > 0 {
		previewLines = append(previewLines, "")
		tagsText := fmt.Sprintf("Tags: %s", strings.Join(p.Tags, ", "))
		previewLines = append(previewLines, theme.PreviewTagsStyle().Render(tagsText))
	}

	// Content preview
	previewLines = append(previewLines, "")
	previewLines = append(previewLines, "---")
	contentLines := strings.Split(p.Content, "\n")
	maxContentLines := previewHeight - len(previewLines) - 2

	if len(contentLines) > maxContentLines {
		contentLines = contentLines[:maxContentLines]
	}

	for _, line := range contentLines {
		// Truncate long lines
		if len(line) > previewWidth-4 {
			line = line[:previewWidth-4] + "..."
		}
		previewLines = append(previewLines, theme.PreviewContentStyle().Render(line))
	}

	return previewStyle.Render(strings.Join(previewLines, "\n"))
}

// moveSelection moves the selection up or down
func (m *Model) moveSelection(delta int) {
	m.selected += delta

	// Clamp selection
	if m.selected < 0 {
		m.selected = 0
	}
	if m.selected >= len(m.filtered) {
		m.selected = len(m.filtered) - 1
	}
}

// applyFilter applies fuzzy filtering to prompts
func (m *Model) applyFilter() {
	if m.searchInput == "" {
		// Show all prompts
		m.filtered = make([]string, 0, len(m.prompts))
		for path := range m.prompts {
			m.filtered = append(m.filtered, path)
		}
		m.selected = 0
		return
	}

	// Build searchable strings
	var stringsToMatch []string
	var paths []string
	for path, p := range m.prompts {
		// Combine title, tags, and category for search
		searchable := fmt.Sprintf("%s %s %s", p.Title, strings.Join(p.Tags, " "), p.Category)
		stringsToMatch = append(stringsToMatch, searchable)
		paths = append(paths, path)
	}

	// Apply fuzzy matching
	matches := fuzzy.Find(m.searchInput, stringsToMatch)

	// Update filtered list
	m.filtered = make([]string, 0, len(matches))
	for _, match := range matches {
		m.filtered = append(m.filtered, paths[match.Index])
	}

	m.selected = 0
}

// Show makes the browser visible
func (m *Model) Show() {
	m.visible = true
	m.searchInput = ""
	m.applyFilter()
}

// Hide makes the browser invisible
func (m *Model) Hide() {
	m.visible = false
}

// IsVisible returns whether the browser is visible
func (m Model) IsVisible() bool {
	return m.visible
}

// Messages
type InsertPromptMsg struct {
	FilePath    string
	InsertMode  InsertMode
	HasWarnings bool // true if prompt has validation warnings
}

type ValidationErrorMsg struct {
	FilePath string
	Errors   []prompt.ValidationError
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
