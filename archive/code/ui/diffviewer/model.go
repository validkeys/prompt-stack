package diffviewer

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

// Model represents the diff viewer modal Bubble Tea model
type Model struct {
	viewport     viewport.Model
	diff         *ai.UnifiedDiff
	original     string
	edits        []ai.Edit
	width        int
	height       int
	onAccept     func() tea.Cmd
	onReject     func() tea.Cmd
	scrollOffset int
}

// NewModel creates a new diff viewer model
func NewModel() Model {
	vp := viewport.New(0, 0)

	return Model{
		viewport:     vp,
		diff:         nil,
		original:     "",
		edits:        []ai.Edit{},
		width:        80,
		height:       20,
		scrollOffset: 0,
	}
}

// SetDiff sets the diff to display
func (m *Model) SetDiff(diff *ai.UnifiedDiff) {
	m.diff = diff
	m.renderDiff()
}

// SetOriginal sets the original content
func (m *Model) SetOriginal(original string) {
	m.original = original
}

// SetEdits sets the edits that generated the diff
func (m *Model) SetEdits(edits []ai.Edit) {
	m.edits = edits
}

// SetOnAccept sets the callback for accepting the diff
func (m *Model) SetOnAccept(callback func() tea.Cmd) {
	m.onAccept = callback
}

// SetOnReject sets the callback for rejecting the diff
func (m *Model) SetOnReject(callback func() tea.Cmd) {
	m.onReject = callback
}

// SetSize sets the size of the diff viewer
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Update viewport size
	m.viewport.Width = width
	m.viewport.Height = height - 4 // Reserve space for header and footer
}

// Update handles messages and updates the model state
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "y":
			// Accept diff
			if m.onAccept != nil {
				return m, m.onAccept()
			}

		case "esc", "n", "q":
			// Reject diff
			if m.onReject != nil {
				return m, m.onReject()
			}

		case "up", "k":
			// Scroll up
			m.viewport.LineUp(1)

		case "down", "j":
			// Scroll down
			m.viewport.LineDown(1)

		case "pgup":
			// Page up
			m.viewport.HalfViewUp()

		case "pgdown":
			// Page down
			m.viewport.HalfViewDown()

		case "home", "g":
			// Go to top
			m.viewport.GotoTop()

		case "end", "G":
			// Go to bottom
			m.viewport.GotoBottom()
		}

	case tea.WindowSizeMsg:
		// Handle window resize
		m.SetSize(msg.Width, msg.Height)
	}

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}

// View renders the diff viewer
func (m Model) View() string {
	if m.diff == nil {
		return m.renderEmpty()
	}

	// Render header
	header := m.renderHeader()

	// Render diff content
	diffContent := m.viewport.View()

	// Render footer
	footer := m.renderFooter()

	// Combine all parts
	return lipgloss.JoinVertical(lipgloss.Left, header, diffContent, footer)
}

// renderHeader renders the modal header
func (m Model) renderHeader() string {
	title := theme.HeaderStyle().Render("📝 Diff Viewer")

	// Calculate stats
	additions := 0
	deletions := 0
	for _, hunk := range m.diff.Hunks {
		for _, line := range hunk.Lines {
			if line.Type == ai.DiffLineAddition {
				additions++
			} else if line.Type == ai.DiffLineDeletion {
				deletions++
			}
		}
	}

	stats := fmt.Sprintf("+%d/-%d", additions, deletions)
	statsStyle := theme.InfoStyle().Render(stats)

	return lipgloss.JoinHorizontal(lipgloss.Left, title, "  ", statsStyle)
}

// renderFooter renders the modal footer with help text
func (m Model) renderFooter() string {
	help := []string{
		"↑/↓ or j/k: Scroll",
		"Enter/y: Accept",
		"Esc/n/q: Reject",
		"Home/End: Jump",
	}

	helpText := strings.Join(help, "  |  ")

	return theme.StatusStyle().Render(helpText)
}

// renderEmpty renders the empty state
func (m Model) renderEmpty() string {
	emptyText := "No diff to display"

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(theme.ListEmptyStyle().Render(emptyText))
}

// renderDiff renders the diff content into the viewport
func (m *Model) renderDiff() {
	if m.diff == nil {
		m.viewport.SetContent("")
		return
	}

	var builder strings.Builder

	// Write header
	builder.WriteString(m.diff.Header)
	builder.WriteString("\n\n")

	// Write hunks
	for _, hunk := range m.diff.Hunks {
		// Write hunk header
		hunkHeader := fmt.Sprintf("@@ -%d,%d +%d,%d @@",
			hunk.OldStart, hunk.OldLines, hunk.NewStart, hunk.NewLines)
		builder.WriteString(theme.DiffHunkHeaderStyle().Render(hunkHeader))
		builder.WriteString("\n")

		// Write lines
		for _, line := range hunk.Lines {
			var styledLine string
			switch line.Type {
			case ai.DiffLineContext:
				styledLine = theme.DiffContextStyle().Render(" " + line.Content)
			case ai.DiffLineAddition:
				styledLine = theme.DiffAdditionStyle().Render("+" + line.Content)
			case ai.DiffLineDeletion:
				styledLine = theme.DiffDeletionStyle().Render("-" + line.Content)
			}
			builder.WriteString(styledLine)
			builder.WriteString("\n")
		}

		builder.WriteString("\n")
	}

	m.viewport.SetContent(builder.String())
	m.viewport.GotoTop()
}

// diffItem implements list.Item for diff hunks (alternative implementation)
type diffItem struct {
	hunk ai.DiffHunk
}

// FilterValue implements list.Item
func (i diffItem) FilterValue() string {
	return fmt.Sprintf("@@ -%d,%d +%d,%d @@",
		i.hunk.OldStart, i.hunk.OldLines, i.hunk.NewStart, i.hunk.NewLines)
}

// diffDelegate implements list.ItemDelegate for custom rendering
type diffDelegate struct{}

// Height implements list.ItemDelegate
func (d diffDelegate) Height() int {
	return 1
}

// Spacing implements list.ItemDelegate
func (d diffDelegate) Spacing() int {
	return 0
}

// Update implements list.ItemDelegate
func (d diffDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

// Render implements list.ItemDelegate
func (d diffDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	h, ok := item.(diffItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	// Render hunk header
	hunkHeader := fmt.Sprintf("@@ -%d,%d +%d,%d @@",
		h.hunk.OldStart, h.hunk.OldLines, h.hunk.NewStart, h.hunk.NewLines)

	var style lipgloss.Style
	if isSelected {
		style = theme.ListItemSelectedStyle()
	} else {
		style = theme.ListItemStyle()
	}

	fmt.Fprintln(w, style.Render(hunkHeader))
}
