package palette

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/internal/commands"
	"github.com/kyledavis/prompt-stack/ui/theme"
)

// Model represents the command palette modal
type Model struct {
	registry    *commands.Registry
	filtered    []*commands.Command
	selected    int
	searchInput string
	width       int
	height      int
	visible     bool
	vimMode     bool
}

// New creates a new command palette model
func New(registry *commands.Registry, vimMode bool) Model {
	return Model{
		registry:    registry,
		filtered:    registry.GetAll(),
		selected:    0,
		searchInput: "",
		visible:     false,
		vimMode:     vimMode,
	}
}

// Init initializes the palette model
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
			// Execute selected command
			if len(m.filtered) > 0 {
				cmd := m.filtered[m.selected]
				m.visible = false
				return m, func() tea.Msg {
					err := cmd.Handler()
					if err != nil {
						return ExecuteErrorMsg{
							CommandID: cmd.ID,
							Error:     err,
						}
					}
					return ExecuteSuccessMsg{
						CommandID: cmd.ID,
					}
				}
			}

		case tea.KeyUp:
			m.moveSelection(-1)

		case tea.KeyDown:
			m.moveSelection(1)

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

// View renders the command palette
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	// Calculate modal dimensions
	modalWidth := min(m.width-4, 80)
	modalHeight := min(m.height-4, 25)

	// Create modal style
	modalStyle := theme.BaseModal().
		Width(modalWidth).
		Height(modalHeight)

	// Render header
	header := m.renderHeader()

	// Render search input
	searchInput := m.renderSearchInput(modalWidth)

	// Render command list
	commandList := m.renderCommandList(modalWidth, modalHeight-3)

	// Combine sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		searchInput,
		commandList,
	)

	return modalStyle.Render(content)
}

// renderHeader renders the palette header
func (m Model) renderHeader() string {
	return theme.HeaderStyle().Render("Command Palette")
}

// renderSearchInput renders the search input field
func (m Model) renderSearchInput(width int) string {
	promptText := ">"
	if m.searchInput != "" {
		promptText = "> " + m.searchInput
	}

	return theme.SearchInputStyle().
		Width(width - 2).
		Render(promptText)
}

// renderCommandList renders the filtered command list
func (m Model) renderCommandList(width, height int) string {
	listStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(0, 1)

	// Render each command
	var items []string
	for i, cmd := range m.filtered {
		// Style based on selection
		var itemStyle lipgloss.Style
		if i == m.selected {
			itemStyle = theme.ListItemSelectedStyle()
		} else {
			itemStyle = theme.ListItemStyle()
		}

		// Render with category label if present
		var commandText string
		if cmd.Category != "" {
			categoryLabel := theme.ListCategoryStyle().Render(fmt.Sprintf("[%s]", cmd.Category))
			commandText = fmt.Sprintf("%s %s", categoryLabel, cmd.Name)
		} else {
			commandText = cmd.Name
		}

		// Add description if available
		if cmd.Description != "" {
			commandText += " " + theme.ListDescriptionStyle().Render(cmd.Description)
		}

		items = append(items, itemStyle.Render(commandText))
	}

	// Add empty state if no results
	if len(items) == 0 {
		items = append(items, theme.ListEmptyStyle().Render("No commands found"))
	}

	// Truncate to fit height
	if len(items) > height {
		items = items[:height]
	}

	return listStyle.Render(strings.Join(items, "\n"))
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

// applyFilter applies fuzzy filtering to commands
func (m *Model) applyFilter() {
	m.filtered = m.registry.Search(m.searchInput)
	m.selected = 0
}

// Show makes the palette visible
func (m *Model) Show() {
	m.visible = true
	m.searchInput = ""
	m.applyFilter()
}

// Hide makes the palette invisible
func (m *Model) Hide() {
	m.visible = false
}

// IsVisible returns whether the palette is visible
func (m Model) IsVisible() bool {
	return m.visible
}

// Messages
type ExecuteSuccessMsg struct {
	CommandID string
}

type ExecuteErrorMsg struct {
	CommandID string
	Error     error
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
