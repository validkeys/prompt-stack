package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/internal/editor"
	"github.com/kyledavis/prompt-stack/ui/theme"
)

// Model represents the composition workspace
type Model struct {
	content              string
	cursor               cursor
	filePath             string
	workingDir           string
	isDirty              bool
	lastSave             time.Time
	saveStatus           string // "saving", "saved", ""
	saveTimer            *time.Timer
	viewport             viewport
	statusBar            statusBar
	width                int
	height               int
	placeholders         []editor.Placeholder
	activePlaceholder    int    // -1 if none active
	placeholderEditMode  bool   // true when editing a placeholder
	placeholderEditValue string // current value being edited
}

type cursor struct {
	x int
	y int
}

type viewport struct {
	x int
	y int
}

type statusBar struct {
	charCount int
	lineCount int
	message   string
}

// New creates a new workspace model
func New(workingDir string) Model {
	return Model{
		content:    "",
		cursor:     cursor{x: 0, y: 0},
		viewport:   viewport{x: 0, y: 0},
		workingDir: workingDir,
		isDirty:    false,
		saveStatus: "",
	}
}

// Init initializes the workspace model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle placeholder edit mode
		if m.placeholderEditMode {
			return m.handlePlaceholderEdit(msg)
		}

		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyUp:
			m.moveCursorUp()

		case tea.KeyDown:
			m.moveCursorDown()

		case tea.KeyLeft:
			m.moveCursorLeft()

		case tea.KeyRight:
			m.moveCursorRight()

		case tea.KeyBackspace:
			m.backspace()
			m.markDirty()
			m.scheduleAutoSave()

		case tea.KeyEnter:
			m.insertNewline()
			m.markDirty()
			m.scheduleAutoSave()

		case tea.KeyTab:
			// Navigate to next placeholder
			if m.navigateToNextPlaceholder() {
				// Successfully navigated to placeholder
			} else {
				m.insertTab()
				m.markDirty()
				m.scheduleAutoSave()
			}

		case tea.KeyShiftTab:
			// Navigate to previous placeholder
			m.navigateToPreviousPlaceholder()

		case tea.KeyRunes:
			// Check for 'i' or 'a' to enter placeholder edit mode
			if m.activePlaceholder >= 0 && len(msg.Runes) == 1 {
				r := msg.Runes[0]
				if r == 'i' || r == 'a' {
					m.enterPlaceholderEditMode()
					return m, nil
				}
			}
			// Normal typing
			m.insertRune(msg.Runes)
			m.markDirty()
			m.scheduleAutoSave()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case autoSaveMsg:
		m.saveStatus = "saving"
		return m, tea.Cmd(func() tea.Msg {
			err := m.saveToFile()
			if err != nil {
				return saveErrorMsg{err}
			}
			return saveSuccessMsg{}
		})

	case saveSuccessMsg:
		m.saveStatus = "saved"
		m.lastSave = time.Now()
		m.isDirty = false
		// Clear saved status after 2 seconds
		return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
			return clearSaveStatusMsg{}
		})

	case saveErrorMsg:
		m.saveStatus = "error"
		m.statusBar.message = fmt.Sprintf("Save failed: %v", msg.err)
		// Keep error status visible
		return m, nil

	case clearSaveStatusMsg:
		m.saveStatus = ""
		m.statusBar.message = ""
		return m, nil
	}

	// Update viewport to keep cursor in view
	m.adjustViewport()

	// Update status bar counts
	m.updateStatusBar()

	return m, cmd
}

// View renders the workspace
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// Calculate available space (leave room for status bar)
	availableHeight := m.height - 1

	// Get visible lines
	lines := m.getVisibleLines(availableHeight)

	// Render lines with placeholder highlighting
	renderedLines := make([]string, len(lines))
	for i, line := range lines {
		if i == m.cursor.y {
			// Cursor line - render with cursor
			renderedLines[i] = m.renderCursorLine(lines)
		} else {
			// Non-cursor line - render with placeholder highlighting
			renderedLines[i] = m.renderLineWithPlaceholders(line, i)
		}
	}

	// Combine rendered lines
	content := strings.Join(renderedLines, "\n")

	// Style the editor
	editorStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(availableHeight).
		Padding(0, 1)

	editorView := editorStyle.Render(content)

	// Render status bar
	statusBarView := m.renderStatusBar()

	// Combine
	return lipgloss.JoinVertical(lipgloss.Left, editorView, statusBarView)
}

// moveCursorUp moves the cursor up one line
func (m *Model) moveCursorUp() {
	if m.cursor.y > 0 {
		m.cursor.y--
		// Adjust x to line length
		lines := strings.Split(m.content, "\n")
		if m.cursor.y < len(lines) {
			lineLen := lines[m.cursor.y]
			if m.cursor.x > len(lineLen) {
				m.cursor.x = len(lineLen)
			}
		}
	}
}

// moveCursorDown moves the cursor down one line
func (m *Model) moveCursorDown() {
	lines := strings.Split(m.content, "\n")
	if m.cursor.y < len(lines)-1 {
		m.cursor.y++
		// Adjust x to line length
		if m.cursor.y < len(lines) {
			lineLen := lines[m.cursor.y]
			if m.cursor.x > len(lineLen) {
				m.cursor.x = len(lineLen)
			}
		}
	}
}

// moveCursorLeft moves the cursor left one character
func (m *Model) moveCursorLeft() {
	if m.cursor.x > 0 {
		m.cursor.x--
	} else if m.cursor.y > 0 {
		// Move to end of previous line
		m.cursor.y--
		lines := strings.Split(m.content, "\n")
		if m.cursor.y < len(lines) {
			lineLen := lines[m.cursor.y]
			m.cursor.x = len(lineLen)
		}
	}
}

// moveCursorRight moves the cursor right one character
func (m *Model) moveCursorRight() {
	lines := strings.Split(m.content, "\n")
	if m.cursor.y < len(lines) {
		lineLen := lines[m.cursor.y]
		if m.cursor.x < len(lineLen) {
			m.cursor.x++
		} else if m.cursor.y < len(lines)-1 {
			// Move to start of next line
			m.cursor.y++
			m.cursor.x = 0
		}
	}
}

// backspace deletes the character before the cursor
func (m *Model) backspace() {
	if m.cursor.x == 0 && m.cursor.y == 0 {
		return
	}

	lines := strings.Split(m.content, "\n")

	if m.cursor.x == 0 {
		// Delete newline, merge with previous line
		if m.cursor.y > 0 {
			prevLine := lines[m.cursor.y-1]
			currentLine := lines[m.cursor.y]
			m.cursor.x = len(prevLine)
			lines[m.cursor.y-1] = prevLine + currentLine
			lines = append(lines[:m.cursor.y], lines[m.cursor.y+1:]...)
			m.cursor.y--
		}
	} else {
		// Delete character in current line
		if m.cursor.y < len(lines) {
			line := lines[m.cursor.y]
			if m.cursor.x <= len(line) {
				lines[m.cursor.y] = line[:m.cursor.x-1] + line[m.cursor.x:]
				m.cursor.x--
			}
		}
	}

	m.content = strings.Join(lines, "\n")

	// Re-parse placeholders after content change
	m.updatePlaceholders()
}

// insertNewline inserts a newline at the cursor position
func (m *Model) insertNewline() {
	lines := strings.Split(m.content, "\n")

	if m.cursor.y < len(lines) {
		line := lines[m.cursor.y]
		before := line[:m.cursor.x]
		after := line[m.cursor.x:]

		lines[m.cursor.y] = before
		lines = append(lines[:m.cursor.y+1], lines[m.cursor.y:]...)
		lines[m.cursor.y+1] = after

		m.cursor.y++
		m.cursor.x = 0
	}

	m.content = strings.Join(lines, "\n")

	// Re-parse placeholders after content change
	m.updatePlaceholders()
}

// insertTab inserts a tab at the cursor position
func (m *Model) insertTab() {
	m.insertRune([]rune{' ', ' ', ' ', ' '}) // 4 spaces
}

// insertRune inserts runes at the cursor position
func (m *Model) insertRune(runes []rune) {
	lines := strings.Split(m.content, "\n")

	if m.cursor.y < len(lines) {
		line := lines[m.cursor.y]
		before := line[:m.cursor.x]
		after := line[m.cursor.x:]

		lines[m.cursor.y] = before + string(runes) + after
		m.cursor.x += len(runes)
	}

	m.content = strings.Join(lines, "\n")

	// Re-parse placeholders after content change
	m.updatePlaceholders()
}

// markDirty marks the composition as having unsaved changes
func (m *Model) markDirty() {
	m.isDirty = true
}

// scheduleAutoSave schedules an auto-save with debouncing
func (m *Model) scheduleAutoSave() {
	if m.saveTimer != nil {
		m.saveTimer.Stop()
	}
	m.saveTimer = time.AfterFunc(750*time.Millisecond, func() {
		// This will be handled by the tea.Cmd in Update
	})
}

// saveToFile saves the composition to a markdown file
func (m *Model) saveToFile() error {
	// Create file path if not exists
	if m.filePath == "" {
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		m.filePath = filepath.Join(m.workingDir, ".promptstack", ".history", timestamp+".md")
	}

	// Ensure directory exists
	dir := filepath.Dir(m.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write content
	if err := os.WriteFile(m.filePath, []byte(m.content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// adjustViewport adjusts the viewport to keep the cursor visible
func (m *Model) adjustViewport() {
	// Simple viewport adjustment - keep cursor in middle third
	availableHeight := m.height - 1
	if availableHeight <= 0 {
		return
	}

	third := availableHeight / 3

	if m.cursor.y < m.viewport.y+third {
		m.viewport.y = max(0, m.cursor.y-third)
	} else if m.cursor.y > m.viewport.y+availableHeight-third {
		m.viewport.y = max(0, m.cursor.y-availableHeight+third)
	}
}

// getVisibleLines returns the lines visible in the viewport
func (m *Model) getVisibleLines(height int) []string {
	lines := strings.Split(m.content, "\n")

	start := m.viewport.y
	end := min(start+height, len(lines))

	if start >= len(lines) {
		return []string{}
	}

	return lines[start:end]
}

// renderCursorLine renders the line with cursor position highlighted
func (m *Model) renderCursorLine(lines []string) string {
	if m.cursor.y >= len(lines) {
		return ""
	}

	line := lines[m.cursor.y]

	// If in placeholder edit mode, show edit value instead of placeholder
	if m.placeholderEditMode && m.activePlaceholder >= 0 {
		ph := m.placeholders[m.activePlaceholder]
		if ph.Type == "text" {
			// Calculate line start position
			lineStartPos := 0
			allLines := strings.Split(m.content, "\n")
			for i := 0; i < m.cursor.y && i < len(allLines); i++ {
				lineStartPos += len(allLines[i]) + 1
			}

			// Check if placeholder is on this line
			if ph.StartPos >= lineStartPos && ph.EndPos <= lineStartPos+len(line) {
				phStart := ph.StartPos - lineStartPos
				phEnd := ph.EndPos - lineStartPos

				// Replace placeholder with edit value
				before := line[:phStart]
				after := line[phEnd:]
				editValue := m.placeholderEditValue

				// Adjust cursor position to be within edit value
				cursorPosInEdit := m.cursor.x - phStart
				if cursorPosInEdit < 0 {
					cursorPosInEdit = 0
				} else if cursorPosInEdit > len(editValue) {
					cursorPosInEdit = len(editValue)
				}

				// Style cursor
				cursorStyle := theme.CursorStyle()

				if cursorPosInEdit < len(editValue) {
					return before + cursorStyle.Render(string(editValue[cursorPosInEdit])) + editValue[cursorPosInEdit+1:] + after
				}

				return before + editValue + cursorStyle.Render(" ") + after
			}
		}
	}

	if m.cursor.x > len(line) {
		return line
	}

	before := line[:m.cursor.x]
	after := line[m.cursor.x:]

	// Style cursor
	cursorStyle := theme.CursorStyle()

	if m.cursor.x < len(line) {
		return before + cursorStyle.Render(string(line[m.cursor.x])) + after
	}

	return before + cursorStyle.Render(" ")
}

// renderLineWithPlaceholders renders a line with placeholder highlighting
func (m *Model) renderLineWithPlaceholders(line string, lineIndex int) string {
	if len(m.placeholders) == 0 {
		return line
	}

	// Calculate line start position in content
	lineStartPos := 0
	lines := strings.Split(m.content, "\n")
	for i := 0; i < lineIndex && i < len(lines); i++ {
		lineStartPos += len(lines[i]) + 1 // +1 for newline
	}

	// Find placeholders on this line
	result := line
	offset := 0

	for _, ph := range m.placeholders {
		// Check if placeholder is on this line
		if ph.StartPos >= lineStartPos && ph.EndPos <= lineStartPos+len(line) {
			// Calculate position within line
			phStart := ph.StartPos - lineStartPos
			phEnd := ph.EndPos - lineStartPos

			// Apply highlighting if active
			if m.activePlaceholder >= 0 && m.placeholders[m.activePlaceholder].Name == ph.Name {
				placeholderStyle := theme.ActivePlaceholderStyle()
				placeholderText := line[phStart+offset : phEnd+offset]
				result = result[:phStart+offset] + placeholderStyle.Render(placeholderText) + result[phEnd+offset:]
				offset += len(placeholderStyle.Render(placeholderText)) - (phEnd - phStart)
			}
		}
	}

	return result
}

// updateStatusBar updates the status bar information
func (m *Model) updateStatusBar() {
	m.statusBar.charCount = len(m.content)
	m.statusBar.lineCount = strings.Count(m.content, "\n") + 1
}

// renderStatusBar renders the status bar
func (m *Model) renderStatusBar() string {
	statusStyle := theme.StatusStyle().
		Width(m.width).
		Height(1)

	// Build status message
	var parts []string

	// Placeholder edit mode indicator
	if m.placeholderEditMode {
		parts = append(parts, "[PLACEHOLDER EDIT MODE]")
	}

	// Auto-save indicator
	if m.saveStatus == "saving" {
		parts = append(parts, "Saving...")
	} else if m.saveStatus == "saved" {
		parts = append(parts, "Saved ✓")
	} else if m.saveStatus == "error" {
		parts = append(parts, "Save failed")
	}

	// Character and line counts
	if m.statusBar.charCount > 0 {
		parts = append(parts, fmt.Sprintf("%d chars, %d lines", m.statusBar.charCount, m.statusBar.lineCount))
	}

	// Custom message
	if m.statusBar.message != "" {
		parts = append(parts, m.statusBar.message)
	}

	// Join with separator
	statusText := strings.Join(parts, " | ")

	return statusStyle.Render(statusText)
}

// Messages
type autoSaveMsg struct{}
type saveSuccessMsg struct{}
type saveErrorMsg struct {
	err error
}
type clearSaveStatusMsg struct{}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// updatePlaceholders re-parses placeholders from content
func (m *Model) updatePlaceholders() {
	m.placeholders = editor.ParsePlaceholders(m.content)

	// If we had an active placeholder, try to find it again
	if m.activePlaceholder >= 0 && m.activePlaceholder < len(m.placeholders) {
		// Keep the same index if it's still valid
	} else {
		m.activePlaceholder = -1
	}
}

// navigateToNextPlaceholder moves to the next placeholder
func (m *Model) navigateToNextPlaceholder() bool {
	cursorPos := m.getCursorPosition()
	nextIndex := editor.GetNextPlaceholder(m.placeholders, cursorPos)

	if nextIndex >= 0 {
		m.activePlaceholder = nextIndex
		ph := m.placeholders[nextIndex]
		m.setCursorToPosition(ph.StartPos)
		return true
	}

	return false
}

// navigateToPreviousPlaceholder moves to the previous placeholder
func (m *Model) navigateToPreviousPlaceholder() bool {
	cursorPos := m.getCursorPosition()
	prevIndex := editor.GetPreviousPlaceholder(m.placeholders, cursorPos)

	if prevIndex >= 0 {
		m.activePlaceholder = prevIndex
		ph := m.placeholders[prevIndex]
		m.setCursorToPosition(ph.StartPos)
		return true
	}

	return false
}

// getCursorPosition returns the absolute cursor position in content
func (m *Model) getCursorPosition() int {
	lines := strings.Split(m.content, "\n")
	pos := 0

	for i := 0; i < m.cursor.y && i < len(lines); i++ {
		pos += len(lines[i]) + 1 // +1 for newline
	}

	if m.cursor.y < len(lines) {
		pos += m.cursor.x
	}

	return pos
}

// setCursorToPosition sets the cursor to an absolute position in content
func (m *Model) setCursorToPosition(pos int) {
	lines := strings.Split(m.content, "\n")
	currentPos := 0

	for i, line := range lines {
		lineEnd := currentPos + len(line)

		if pos <= lineEnd {
			m.cursor.y = i
			m.cursor.x = pos - currentPos
			return
		}

		currentPos = lineEnd + 1 // +1 for newline
	}

	// If position is beyond content, set to end
	m.cursor.y = len(lines) - 1
	m.cursor.x = len(lines[len(lines)-1])
}

// enterPlaceholderEditMode enters placeholder edit mode
func (m *Model) enterPlaceholderEditMode() {
	if m.activePlaceholder < 0 || m.activePlaceholder >= len(m.placeholders) {
		return
	}

	ph := &m.placeholders[m.activePlaceholder]

	// Only text placeholders can be edited in this mode
	if ph.Type != "text" {
		return
	}

	// Initialize edit value with current value or empty string
	m.placeholderEditMode = true
	m.placeholderEditValue = ph.CurrentValue

	// Move cursor to placeholder position
	m.setCursorToPosition(ph.StartPos)
}

// handlePlaceholderEdit handles key events when in placeholder edit mode
func (m Model) handlePlaceholderEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		// Exit placeholder edit mode and save the value
		m.exitPlaceholderEditMode()
		return m, nil

	case tea.KeyBackspace:
		// Delete character from edit value
		if len(m.placeholderEditValue) > 0 {
			m.placeholderEditValue = m.placeholderEditValue[:len(m.placeholderEditValue)-1]
		}
		return m, nil

	case tea.KeyEnter:
		// Exit placeholder edit mode and save the value
		m.exitPlaceholderEditMode()
		return m, nil

	case tea.KeyRunes:
		// Append characters to edit value
		m.placeholderEditValue += string(msg.Runes)
		return m, nil
	}

	return m, nil
}

// exitPlaceholderEditMode exits placeholder edit mode and applies the value
func (m *Model) exitPlaceholderEditMode() {
	if m.activePlaceholder < 0 || m.activePlaceholder >= len(m.placeholders) {
		m.placeholderEditMode = false
		m.placeholderEditValue = ""
		return
	}

	ph := &m.placeholders[m.activePlaceholder]

	// Update placeholder's current value
	ph.CurrentValue = m.placeholderEditValue

	// Replace placeholder in content with the filled value
	m.content = editor.ReplacePlaceholder(m.content, *ph)

	// Re-parse placeholders after replacement
	m.updatePlaceholders()

	// Mark as dirty and schedule auto-save
	m.markDirty()
	m.scheduleAutoSave()

	// Exit edit mode
	m.placeholderEditMode = false
	m.placeholderEditValue = ""
}
