// Package workspace provides text editor workspace model for PromptStack TUI.
// It integrates editor components (cursor, viewport, placeholder, fileio) to provide
// a complete text editing experience with Bubble Tea.
package workspace

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/internal/editor"
	"github.com/kyledavis/prompt-stack/ui/theme"
)

// Model represents text editor workspace.
// It integrates editor components to provide text editing functionality.
type Model struct {
	cursor       editor.Cursor
	viewport     editor.Viewport
	placeholders editor.Manager
	fileManager  editor.FileManager
	content      string
	width        int
	height       int
	statusBar    statusBar
	isReadOnly   bool
}

type statusBar struct {
	charCount int
	lineCount int
	message   string
}

// New creates a new workspace model.
func New() Model {
	return Model{
		cursor:       editor.NewCursor(),
		viewport:     editor.NewViewport(24),
		placeholders: editor.New(),
		fileManager:  editor.NewFileManager(""),
		content:      "",
		width:        0,
		height:       0,
		isReadOnly:   false,
	}
}

// Init initializes the workspace model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Block all editing when in read-only mode
		if m.isReadOnly {
			// Only allow cursor navigation in read-only mode
			switch msg.Type {
			case tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
				// Allow cursor navigation
			default:
				// Block all other keys
				return m, nil
			}
		}

		// Handle text placeholder edit mode
		if m.placeholders.IsEditing() {
			return m.handlePlaceholderEdit(msg)
		}

		switch msg.Type {
		case tea.KeyCtrlC:
			// Check for unsaved changes before quitting
			if m.fileManager.IsModified() {
				// In a real implementation, we'd show a confirmation dialog
				// For now, we'll just save before quitting
				newModel := m
				if err := newModel.saveToFile(); err != nil {
					// Show error in status bar
					newModel = newModel.SetStatus(fmt.Sprintf("Save failed: %v", err))
					return newModel, nil
				}
				return newModel, tea.Quit
			}
			return m, tea.Quit

		case tea.KeyUp:
			return m.moveCursorUp().adjustViewport(), nil

		case tea.KeyDown:
			return m.moveCursorDown().adjustViewport(), nil

		case tea.KeyLeft:
			return m.moveCursorLeft().adjustViewport(), nil

		case tea.KeyRight:
			return m.moveCursorRight().adjustViewport(), nil

		case tea.KeyBackspace:
			newModel := m.backspace()
			newModel = newModel.markDirty()
			newModel = newModel.updateStatusBar()
			return newModel, nil

		case tea.KeyEnter:
			newModel := m.insertNewline()
			newModel = newModel.markDirty()
			newModel = newModel.updateStatusBar()
			return newModel, nil

		case tea.KeyTab:
			// Navigate to next placeholder
			newModel := m.navigateToNextPlaceholder()
			// Check if we successfully navigated to a placeholder
			if newModel.placeholders.Active() != nil {
				// Successfully navigated to placeholder
				return newModel, nil
			} else {
				// No placeholder found, insert tab
				newModel = newModel.insertTab()
				newModel = newModel.markDirty()
				newModel = newModel.updateStatusBar()
				return newModel, nil
			}

		case tea.KeyShiftTab:
			// Navigate to previous placeholder
			return m.navigateToPreviousPlaceholder(), nil

		case tea.KeySpace:
			// Handle space bar explicitly
			newModel := m.insertRune([]rune{' '})
			newModel = newModel.markDirty()
			newModel = newModel.updateStatusBar()
			return newModel, nil

		case tea.KeyRunes:
			// Check for 'i' to enter placeholder edit mode
			if m.placeholders.Active() != nil && len(msg.Runes) == 1 {
				r := msg.Runes[0]
				if r == 'i' {
					return m.enterPlaceholderEditMode(), nil
				}
			}
			// Normal typing
			newModel := m.insertRune(msg.Runes)
			newModel = newModel.markDirty()
			newModel = newModel.updateStatusBar()
			return newModel, nil
		}

	case tea.WindowSizeMsg:
		newModel := m
		newModel.width = msg.Width
		newModel.height = msg.Height
		newModel.viewport = newModel.viewport.SetHeight(msg.Height - 1) // Leave room for status bar
		return newModel, nil
	}

	// Update viewport to keep cursor in view
	m = m.adjustViewport()

	// Update status bar counts
	m = m.updateStatusBar()

	return m, nil
}

// View renders the workspace.
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
		if i == m.cursor.Y() {
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

// moveCursorUp moves the cursor up one line.
func (m Model) moveCursorUp() Model {
	newModel := m
	newModel.cursor = newModel.cursor.MoveUp()
	newModel.cursor = newModel.cursor.AdjustToLineLength(newModel.content)
	return newModel
}

// moveCursorDown moves the cursor down one line.
func (m Model) moveCursorDown() Model {
	newModel := m
	newModel.cursor = newModel.cursor.MoveDown()
	newModel.cursor = newModel.cursor.AdjustToLineLength(newModel.content)
	return newModel
}

// moveCursorLeft moves the cursor left one character.
func (m Model) moveCursorLeft() Model {
	newModel := m
	newModel.cursor = newModel.cursor.MoveLeft()
	return newModel
}

// moveCursorRight moves the cursor right one character.
func (m Model) moveCursorRight() Model {
	newModel := m
	newModel.cursor = newModel.cursor.MoveRight()
	return newModel
}

// backspace deletes the character before the cursor.
func (m Model) backspace() Model {
	newModel := m
	if newModel.cursor.X() == 0 && newModel.cursor.Y() == 0 {
		return newModel
	}

	lines := strings.Split(newModel.content, "\n")

	if newModel.cursor.X() == 0 {
		// Delete newline, merge with previous line
		if newModel.cursor.Y() > 0 {
			prevLine := lines[newModel.cursor.Y()-1]
			currentLine := lines[newModel.cursor.Y()]
			newModel.cursor = newModel.cursor.SetX(len(prevLine))
			newModel.cursor = newModel.cursor.SetY(newModel.cursor.Y() - 1)
			lines[newModel.cursor.Y()] = prevLine + currentLine
			lines = append(lines[:newModel.cursor.Y()+1], lines[newModel.cursor.Y()+2:]...)
		}
	} else {
		// Delete character in current line
		if newModel.cursor.Y() < len(lines) {
			line := lines[newModel.cursor.Y()]
			if newModel.cursor.X() <= len(line) {
				lines[newModel.cursor.Y()] = line[:newModel.cursor.X()-1] + line[newModel.cursor.X():]
				newModel.cursor = newModel.cursor.SetX(newModel.cursor.X() - 1)
			}
		}
	}

	newModel.content = strings.Join(lines, "\n")

	// Re-parse placeholders after content change
	newModel = newModel.updatePlaceholders()

	return newModel
}

// insertNewline inserts a newline at the cursor position.
func (m Model) insertNewline() Model {
	newModel := m
	lines := strings.Split(newModel.content, "\n")

	if newModel.cursor.Y() < len(lines) {
		line := lines[newModel.cursor.Y()]
		before := line[:newModel.cursor.X()]
		after := line[newModel.cursor.X():]

		lines[newModel.cursor.Y()] = before
		lines = append(lines[:newModel.cursor.Y()+1], lines[newModel.cursor.Y():]...)
		lines[newModel.cursor.Y()+1] = after

		newModel.cursor = newModel.cursor.SetY(newModel.cursor.Y() + 1)
		newModel.cursor = newModel.cursor.SetX(0)
	}

	newModel.content = strings.Join(lines, "\n")

	// Re-parse placeholders after content change
	newModel = newModel.updatePlaceholders()

	return newModel
}

// insertTab inserts a tab at the cursor position.
func (m Model) insertTab() Model {
	return m.insertRune([]rune{' ', ' ', ' ', ' '}) // 4 spaces
}

// insertRune inserts runes at the cursor position.
func (m Model) insertRune(runes []rune) Model {
	newModel := m
	lines := strings.Split(newModel.content, "\n")

	if newModel.cursor.Y() < len(lines) {
		line := lines[newModel.cursor.Y()]
		before := line[:newModel.cursor.X()]
		after := line[newModel.cursor.X():]

		lines[newModel.cursor.Y()] = before + string(runes) + after
		newModel.cursor = newModel.cursor.SetX(newModel.cursor.X() + len(runes))
	}

	newModel.content = strings.Join(lines, "\n")

	// Re-parse placeholders after content change
	newModel = newModel.updatePlaceholders()

	return newModel
}

// markDirty marks the composition as having unsaved changes.
func (m Model) markDirty() Model {
	newModel := m
	newModel.fileManager = newModel.fileManager.MarkModified()
	return newModel
}

// saveToFile saves the composition to a markdown file.
func (m Model) saveToFile() error {
	return m.fileManager.Save(m.content)
}

// adjustViewport adjusts the viewport to keep the cursor visible.
func (m Model) adjustViewport() Model {
	newModel := m
	newModel.viewport = newModel.viewport.EnsureVisible(newModel.cursor.Y())
	return newModel
}

// getVisibleLines returns the lines visible in the viewport.
func (m Model) getVisibleLines(height int) []string {
	lines := strings.Split(m.content, "\n")

	start, end := m.viewport.VisibleLines()
	if start >= len(lines) {
		return []string{}
	}

	if end > len(lines) {
		end = len(lines)
	}

	return lines[start:end]
}

// renderCursorLine renders the line with cursor position highlighted.
func (m Model) renderCursorLine(lines []string) string {
	if m.cursor.Y() >= len(lines) {
		return ""
	}

	line := lines[m.cursor.Y()]

	// If in placeholder edit mode, show edit value instead of placeholder
	if m.placeholders.IsEditing() && m.placeholders.Active() != nil {
		ph := m.placeholders.Active()
		if ph.Type == "text" {
			// Calculate line start position
			lineStartPos := 0
			allLines := strings.Split(m.content, "\n")
			for i := 0; i < m.cursor.Y() && i < len(allLines); i++ {
				lineStartPos += len(allLines[i]) + 1
			}

			// Check if placeholder is on this line
			if ph.StartPos >= lineStartPos && ph.EndPos <= lineStartPos+len(line) {
				phStart := ph.StartPos - lineStartPos
				phEnd := ph.EndPos - lineStartPos

				// Replace placeholder with edit value
				before := line[:phStart]
				after := line[phEnd:]
				editValue := m.placeholders.EditValue()

				// Adjust cursor position to be within edit value
				cursorPosInEdit := m.cursor.X() - phStart
				if cursorPosInEdit < 0 {
					cursorPosInEdit = 0
				} else if cursorPosInEdit > len(editValue) {
					cursorPosInEdit = len(editValue)
				}

				// Style cursor
				cursorStyle := lipgloss.NewStyle().
					Background(lipgloss.Color(theme.CursorBackground)).
					Foreground(lipgloss.Color(theme.CursorForeground))

				if cursorPosInEdit < len(editValue) {
					return before + cursorStyle.Render(string(editValue[cursorPosInEdit])) + editValue[cursorPosInEdit+1:] + after
				}

				return before + editValue + cursorStyle.Render(" ") + after
			}
		}
	}

	if m.cursor.X() > len(line) {
		return line
	}

	before := line[:m.cursor.X()]
	after := line[m.cursor.X():]

	// Style cursor
	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.CursorBackground)).
		Foreground(lipgloss.Color(theme.CursorForeground))

	if m.cursor.X() < len(line) {
		return before + cursorStyle.Render(string(line[m.cursor.X()])) + after
	}

	return before + cursorStyle.Render(" ")
}

// renderLineWithPlaceholders renders a line with placeholder highlighting.
func (m Model) renderLineWithPlaceholders(line string, lineIndex int) string {
	placeholders := m.placeholders.Placeholders()
	if len(placeholders) == 0 {
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

	for _, ph := range placeholders {
		// Check if placeholder is on this line
		if ph.StartPos >= lineStartPos && ph.EndPos <= lineStartPos+len(line) {
			// Calculate position within line
			phStart := ph.StartPos - lineStartPos
			phEnd := ph.EndPos - lineStartPos

			// Apply highlighting if active
			if ph.IsActive {
				placeholderStyle := theme.ActivePlaceholderStyle()
				placeholderText := line[phStart+offset : phEnd+offset]
				result = result[:phStart+offset] + placeholderStyle.Render(placeholderText) + result[phEnd+offset:]
				offset += len(placeholderStyle.Render(placeholderText)) - (phEnd - phStart)
			}
		}
	}

	return result
}

// updateStatusBar updates the status bar information.
func (m Model) updateStatusBar() Model {
	newModel := m
	newModel.statusBar.charCount = len(newModel.content)
	newModel.statusBar.lineCount = strings.Count(newModel.content, "\n") + 1
	return newModel
}

// renderStatusBar renders the status bar.
func (m Model) renderStatusBar() string {
	statusStyle := theme.StatusStyle().
		Width(m.width).
		Height(1)

	// Build status message
	var parts []string

	// Placeholder edit mode indicator
	if m.placeholders.IsEditing() {
		parts = append(parts, "[PLACEHOLDER EDIT MODE]")
	}

	// Auto-save indicator
	if m.fileManager.IsModified() {
		parts = append(parts, "Modified")
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

// updatePlaceholders re-parses placeholders from content.
func (m Model) updatePlaceholders() Model {
	newModel := m
	newModel.placeholders = newModel.placeholders.Update(editor.ParsePlaceholdersMsg{Content: newModel.content})
	return newModel
}

// navigateToNextPlaceholder moves to the next placeholder.
func (m Model) navigateToNextPlaceholder() Model {
	newModel := m
	placeholders := newModel.placeholders.Placeholders()

	if len(placeholders) == 0 {
		return newModel
	}

	cursorPos := newModel.cursor.GetAbsolutePosition(newModel.content)

	// Find next placeholder after current cursor position
	nextIndex := -1
	for i, ph := range placeholders {
		if ph.StartPos > cursorPos {
			nextIndex = i
			break
		}
	}

	// If no placeholder found after cursor, wrap to first
	if nextIndex == -1 {
		nextIndex = 0
	}

	if nextIndex >= 0 {
		newModel.placeholders = newModel.placeholders.Update(editor.ActivatePlaceholderMsg{Index: nextIndex})
		ph := placeholders[nextIndex]
		newModel.cursor = newModel.cursor.MoveToPosition(ph.StartPos, newModel.content)
	}

	return newModel
}

// navigateToPreviousPlaceholder moves to the previous placeholder.
func (m Model) navigateToPreviousPlaceholder() Model {
	newModel := m
	placeholders := newModel.placeholders.Placeholders()

	if len(placeholders) == 0 {
		return newModel
	}

	cursorPos := newModel.cursor.GetAbsolutePosition(newModel.content)

	// Find previous placeholder before current cursor position
	prevIndex := -1
	for i := len(placeholders) - 1; i >= 0; i-- {
		if placeholders[i].EndPos < cursorPos {
			prevIndex = i
			break
		}
	}

	// If no placeholder found before cursor, wrap to last
	if prevIndex == -1 {
		prevIndex = len(placeholders) - 1
	}

	if prevIndex >= 0 {
		newModel.placeholders = newModel.placeholders.Update(editor.ActivatePlaceholderMsg{Index: prevIndex})
		ph := placeholders[prevIndex]
		newModel.cursor = newModel.cursor.MoveToPosition(ph.StartPos, newModel.content)
	}

	return newModel
}

// enterPlaceholderEditMode enters placeholder edit mode.
func (m Model) enterPlaceholderEditMode() Model {
	newModel := m
	ph := m.placeholders.Active()
	if ph == nil {
		return newModel
	}

	// Handle list placeholders
	if ph.Type == "list" {
		// List editing not yet implemented
		return newModel
	}

	// Handle text placeholders - already handled by placeholder manager
	return newModel
}

// handlePlaceholderEdit handles key events when in placeholder edit mode.
func (m Model) handlePlaceholderEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		// Exit placeholder edit mode and save the value
		newModel := m
		ph := newModel.placeholders.Active()
		if ph != nil {
			// Replace placeholder in content with the filled value
			newModel.content = editor.ReplacePlaceholder(newModel.content, *ph)
			newModel = newModel.updatePlaceholders()
			newModel = newModel.markDirty()
			newModel = newModel.updateStatusBar()
		}
		newModel.placeholders = newModel.placeholders.Update(editor.ExitEditModeMsg{})
		return newModel, nil

	case tea.KeyBackspace:
		// Delete character from edit value
		newModel := m
		editValue := newModel.placeholders.EditValue()
		if len(editValue) > 0 {
			newModel.placeholders = newModel.placeholders.Update(editor.EditPlaceholderMsg{Value: editValue[:len(editValue)-1]})
		}
		return newModel, nil

	case tea.KeyEnter:
		// Exit placeholder edit mode and save the value
		newModel := m
		ph := newModel.placeholders.Active()
		if ph != nil {
			// Replace placeholder in content with the filled value
			newModel.content = editor.ReplacePlaceholder(newModel.content, *ph)
			newModel = newModel.updatePlaceholders()
			newModel = newModel.markDirty()
			newModel = newModel.updateStatusBar()
		}
		newModel.placeholders = newModel.placeholders.Update(editor.ExitEditModeMsg{})
		return newModel, nil

	case tea.KeyRunes:
		// Append characters to edit value
		newModel := m
		editValue := newModel.placeholders.EditValue()
		newModel.placeholders = newModel.placeholders.Update(editor.EditPlaceholderMsg{Value: editValue + string(msg.Runes)})
		newModel = newModel.updateStatusBar()
		return newModel, nil
	}

	return m, nil
}

// SetReadOnly sets the read-only state of the workspace.
func (m Model) SetReadOnly(readOnly bool) Model {
	newModel := m
	newModel.isReadOnly = readOnly
	return newModel
}

// IsReadOnly returns whether the workspace is in read-only mode.
func (m Model) IsReadOnly() bool {
	return m.isReadOnly
}

// SetSize sets the size of the workspace.
func (m Model) SetSize(width, height int) Model {
	newModel := m
	newModel.width = width
	newModel.height = height
	newModel.viewport = newModel.viewport.SetHeight(height - 1)
	return newModel
}

// GetContent returns the current content of the workspace.
func (m Model) GetContent() string {
	return m.content
}

// SetContent sets the content of the workspace.
func (m Model) SetContent(content string) Model {
	newModel := m
	newModel.content = content
	newModel = newModel.updatePlaceholders()
	newModel = newModel.updateStatusBar()
	newModel = newModel.markDirty()
	return newModel
}

// GetCursorPosition returns the absolute cursor position in content.
func (m Model) GetCursorPosition() int {
	return m.cursor.GetAbsolutePosition(m.content)
}

// InsertContent inserts content at a specific position.
func (m Model) InsertContent(position int, content string) Model {
	newModel := m
	// Convert position to line/column
	lines := strings.Split(newModel.content, "\n")
	currentPos := 0
	targetLine := 0
	targetCol := 0

	for i, line := range lines {
		lineEnd := currentPos + len(line)
		if position <= lineEnd {
			targetLine = i
			targetCol = position - currentPos
			break
		}
		currentPos = lineEnd + 1 // +1 for newline
	}

	// Insert content at target position
	if targetLine < len(lines) {
		line := lines[targetLine]
		before := line[:targetCol]
		after := line[targetCol:]
		lines[targetLine] = before + content + after
	}

	newModel.content = strings.Join(lines, "\n")
	newModel = newModel.updatePlaceholders()
	newModel = newModel.markDirty()
	return newModel
}

// MarkDirty marks the composition as having unsaved changes.
func (m Model) MarkDirty() Model {
	return m.markDirty()
}

// SetStatus sets a status message in the status bar.
func (m Model) SetStatus(message string) Model {
	newModel := m
	newModel.statusBar.message = message
	return newModel
}
