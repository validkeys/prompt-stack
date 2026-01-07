package cleanup

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/ui/common"
)

// Model represents cleanup UI
type Model struct {
	list             list.Model
	strategyInput    textinput.Model
	ageInput         textinput.Model
	countInput       textinput.Model
	directoryInput   textinput.Model
	showPreview      bool
	previewResult    *PreviewResult
	statistics       *Statistics
	quitting         bool
	width            int
	height           int
	onExecute        func() tea.Cmd
	onCancel         func() tea.Cmd
	onPreview        func() tea.Cmd
	showConfirmation bool
	confirmation     common.ConfirmationModel
}

// PreviewResult represents cleanup preview result
type PreviewResult struct {
	Strategy      string
	FilesToDelete []FileItem
	TotalSize     string
	FileCount     int
}

// Statistics represents history statistics
type Statistics struct {
	TotalCompositions int
	TotalSize         string
	OldestDate        string
	NewestDate        string
	AgeDays           int
}

// FileItem represents a file in the cleanup list
type FileItem struct {
	filePath  string
	date      string
	directory string
	size      string
	preview   string
}

// NewFileItem creates a new file item
func NewFileItem(filePath, date, directory, size, preview string) FileItem {
	return FileItem{
		filePath:  filePath,
		date:      date,
		directory: directory,
		size:      size,
		preview:   preview,
	}
}

// FilterValue implements list.Item interface
func (f FileItem) FilterValue() string {
	return f.filePath
}

// Title implements list.Item interface
func (f FileItem) Title() string {
	return f.filePath
}

// Description implements list.Item interface
func (f FileItem) Description() string {
	return fmt.Sprintf("%s | %s | %s", f.date, f.directory, f.size)
}

// NewModel creates a new cleanup model
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
	l.Title = "History Cleanup"
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	// Create strategy input
	strategy := textinput.New()
	strategy.Placeholder = "Select strategy..."
	strategy.Prompt = "Strategy: "
	strategy.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))
	strategy.TextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))
	strategy.Width = 30

	// Create age input
	age := textinput.New()
	age.Placeholder = "Days"
	age.Prompt = "Age (days): "
	age.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))
	age.TextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))
	age.Width = 10

	// Create count input
	count := textinput.New()
	count.Placeholder = "Count"
	count.Prompt = "Count: "
	count.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))
	count.TextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))
	count.Width = 10

	// Create directory input
	directory := textinput.New()
	directory.Placeholder = "Directory path"
	directory.Prompt = "Directory: "
	directory.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))
	directory.TextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))
	directory.Width = 50

	return Model{
		list:           l,
		strategyInput:  strategy,
		ageInput:       age,
		countInput:     count,
		directoryInput: directory,
		showPreview:    false,
	}
}

// SetPreviewResult sets the preview result
func (m *Model) SetPreviewResult(result *PreviewResult) {
	m.previewResult = result
	m.showPreview = true
}

// SetOnExecute sets the callback for execute action
func (m *Model) SetOnExecute(fn func() tea.Cmd) {
	m.onExecute = fn
}

// SetOnCancel sets the callback for cancel action
func (m *Model) SetOnCancel(fn func() tea.Cmd) {
	m.onCancel = fn
}

// SetOnPreview sets the callback for preview action
func (m *Model) SetOnPreview(fn func() tea.Cmd) {
	m.onPreview = fn
}

// SetStatistics sets the statistics display
func (m *Model) SetStatistics(stats *Statistics) {
	m.statistics = stats
}

// SetFiles sets the list of files to display
func (m *Model) SetFiles(files []FileItem) {
	items := make([]list.Item, len(files))
	for i, file := range files {
		items[i] = file
	}
	m.list.SetItems(items)
}

// SetSize sets the size of the model
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Adjust list size based on preview
	listHeight := height
	if m.showPreview {
		listHeight = height - 10 // Reserve space for preview
	}

	m.list.SetSize(width, listHeight)
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
			if m.showPreview {
				// Exit preview mode
				m.showPreview = false
				m.previewResult = nil
				return m, nil
			}
			m.quitting = true
			if m.onCancel != nil {
				return m, m.onCancel()
			}
			return m, tea.Quit

		case "enter":
			if m.showConfirmation {
				// Let confirmation dialog handle enter
				break
			}

			if m.showPreview {
				// Show confirmation dialog before executing
				m.showConfirmation = true
				m.confirmation = common.DestructiveConfirmation(
					"Confirm Cleanup",
					fmt.Sprintf("Are you sure you want to delete %d files? This action cannot be undone.", m.previewResult.FileCount),
				)
				m.confirmation.SetOnConfirm(func() tea.Cmd {
					if m.onExecute != nil {
						return m.onExecute()
					}
					return nil
				})
				m.confirmation.SetOnCancel(func() tea.Cmd {
					m.showConfirmation = false
					return nil
				})
				m.confirmation.SetSize(m.width, m.height)
				return m, nil
			} else {
				// Show preview
				if m.onPreview != nil {
					return m, m.onPreview()
				}
				m.showPreview = true
			}
			return m, nil

		case "p":
			if !m.showConfirmation {
				// Toggle preview
				m.showPreview = !m.showPreview
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	// Update confirmation dialog if active
	if m.showConfirmation {
		var confirmCmd tea.Cmd
		model, confirmCmd := m.confirmation.Update(msg)
		m.confirmation = model.(common.ConfirmationModel)

		// Check if confirmation was handled
		if m.confirmation.IsConfirmed() || m.confirmation.IsCancelled() {
			m.showConfirmation = false
		}

		return m, confirmCmd
	}

	// Update inputs
	m.strategyInput, cmd = m.strategyInput.Update(msg)
	m.ageInput, cmd = m.ageInput.Update(msg)
	m.countInput, cmd = m.countInput.Update(msg)
	m.directoryInput, cmd = m.directoryInput.Update(msg)

	// Update list
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the model
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// Show confirmation dialog if active
	if m.showConfirmation {
		return m.confirmation.View()
	}

	var b strings.Builder

	// Render title
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Width(m.width).
		Align(lipgloss.Center)

	b.WriteString(titleStyle.Render("History Cleanup"))
	b.WriteString("\n\n")

	// Render statistics if available
	if m.statistics != nil {
		b.WriteString(m.renderStatistics())
		b.WriteString("\n")
	}

	// Render strategy selection
	b.WriteString(m.renderStrategySelection())
	b.WriteString("\n")

	// Render preview if active
	if m.showPreview && m.previewResult != nil {
		b.WriteString(m.renderPreview())
		b.WriteString("\n")
	}

	// Render list
	b.WriteString(m.list.View())
	b.WriteString("\n")

	// Render help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(m.width)

	helpText := "↑/↓: navigate | Enter: preview/execute | P: toggle preview | Esc: cancel"
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}

// renderStatistics renders history statistics
func (m Model) renderStatistics() string {
	if m.statistics == nil {
		return ""
	}

	var b strings.Builder

	// Statistics box style
	statsStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(m.width).
		Padding(0, 1)

	// Header style
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	// Value style
	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))

	b.WriteString(headerStyle.Render("History Statistics"))
	b.WriteString("\n")
	b.WriteString(valueStyle.Render(fmt.Sprintf("Total compositions: %d", m.statistics.TotalCompositions)))
	b.WriteString("\n")
	b.WriteString(valueStyle.Render(fmt.Sprintf("Total size: %s", m.statistics.TotalSize)))
	b.WriteString("\n")
	b.WriteString(valueStyle.Render(fmt.Sprintf("Date range: %s to %s", m.statistics.OldestDate, m.statistics.NewestDate)))
	b.WriteString("\n")
	b.WriteString(valueStyle.Render(fmt.Sprintf("Age: %d days", m.statistics.AgeDays)))

	return statsStyle.Render(b.String())
}

// renderStrategySelection renders strategy selection inputs
func (m Model) renderStrategySelection() string {
	var b strings.Builder

	// Strategy input
	b.WriteString(m.strategyInput.View())
	b.WriteString("\n")

	// Age input (for age-based strategy)
	b.WriteString(m.ageInput.View())
	b.WriteString("\n")

	// Count input (for count-based strategies)
	b.WriteString(m.countInput.View())
	b.WriteString("\n")

	// Directory input (for directory-based strategy)
	b.WriteString(m.directoryInput.View())
	b.WriteString("\n")

	return b.String()
}

// renderPreview renders cleanup preview
func (m Model) renderPreview() string {
	if m.previewResult == nil {
		return ""
	}

	var b strings.Builder

	// Preview box style
	previewStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Width(m.width).
		Padding(0, 1)

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	b.WriteString(headerStyle.Render(fmt.Sprintf("Preview: %s", m.previewResult.Strategy)))
	b.WriteString("\n\n")

	// Statistics
	statStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))

	b.WriteString(statStyle.Render(fmt.Sprintf("Files to delete: %d", m.previewResult.FileCount)))
	b.WriteString("\n")
	b.WriteString(statStyle.Render(fmt.Sprintf("Total size: %s", m.previewResult.TotalSize)))
	b.WriteString("\n\n")

	// Show first few files as preview
	if len(m.previewResult.FilesToDelete) > 0 {
		previewCount := 5
		if len(m.previewResult.FilesToDelete) < previewCount {
			previewCount = len(m.previewResult.FilesToDelete)
		}

		fileStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

		b.WriteString(headerStyle.Render("Files to be deleted:"))
		b.WriteString("\n")
		for i := 0; i < previewCount; i++ {
			file := m.previewResult.FilesToDelete[i]
			b.WriteString(fileStyle.Render(fmt.Sprintf("  • %s", file.filePath)))
			b.WriteString("\n")
		}

		if len(m.previewResult.FilesToDelete) > previewCount {
			b.WriteString(fileStyle.Render(fmt.Sprintf("  ... and %d more files", len(m.previewResult.FilesToDelete)-previewCount)))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Warning
	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")) // Yellow

	b.WriteString(warningStyle.Render("⚠️  This action cannot be undone!"))
	b.WriteString("\n")

	return previewStyle.Render(b.String())
}

// GetStrategy returns the selected strategy
func (m Model) GetStrategy() string {
	return strings.TrimSpace(m.strategyInput.Value())
}

// GetAgeDays returns the age in days
func (m Model) GetAgeDays() int {
	age := strings.TrimSpace(m.ageInput.Value())
	if age == "" {
		return 0
	}
	var days int
	fmt.Sscanf(age, "%d", &days)
	return days
}

// GetCount returns the count
func (m Model) GetCount() int {
	count := strings.TrimSpace(m.countInput.Value())
	if count == "" {
		return 0
	}
	var n int
	fmt.Sscanf(count, "%d", &n)
	return n
}

// GetDirectory returns the directory path
func (m Model) GetDirectory() string {
	return strings.TrimSpace(m.directoryInput.Value())
}

// FormatDate formats a time.Time to a readable string
func FormatDate(t time.Time) string {
	return t.Format("Jan 2, 2006 3:04 PM")
}

// TruncateString truncates a string to max length with ellipsis
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
