// Package workspace provides text editor workspace model for PromptStack TUI.
// It integrates editor components (cursor, viewport, placeholder, fileio) to provide
// a complete text editing experience with Bubble Tea.
package workspace

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyledavis/prompt-stack/internal/editor"
)

// Model represents text editor workspace.
// It integrates editor components to provide text editing functionality.
type Model struct {
	buffer       *editor.Buffer
	viewport     viewport.Model
	placeholders editor.Manager
	fileManager  editor.FileManager
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
		buffer:       editor.NewBuffer(),
		viewport:     viewport.New(0, 24),
		placeholders: editor.New(),
		fileManager:  editor.NewFileManager(""),
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
					// Show error in status bar - can't use fmt here
					newModel = newModel.markDirty()
					return newModel, nil
				}
				return newModel, tea.Quit
			}
			return m, tea.Quit

		case tea.KeyUp:
			newModel := m
			newModel.buffer.MoveUp()
			return newModel.adjustViewport(), nil

		case tea.KeyDown:
			newModel := m
			newModel.buffer.MoveDown()
			return newModel.adjustViewport(), nil

		case tea.KeyLeft:
			newModel := m
			newModel.buffer.MoveLeft()
			return newModel.adjustViewport(), nil

		case tea.KeyRight:
			newModel := m
			newModel.buffer.MoveRight()
			return newModel.adjustViewport(), nil

		case tea.KeyBackspace:
			newModel := m
			newModel.buffer.Delete()
			newModel = newModel.markDirty()
			newModel = newModel.updateStatusBar()
			return newModel, nil

		case tea.KeyEnter:
			newModel := m
			newModel.buffer.Insert('\n')
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
				newModel = m.insertTab()
				newModel = newModel.markDirty()
				newModel = newModel.updateStatusBar()
				return newModel, nil
			}

		case tea.KeyShiftTab:
			// Navigate to previous placeholder
			return m.navigateToPreviousPlaceholder(), nil

		case tea.KeySpace:
			// Handle space bar explicitly
			newModel := m.insertRune(' ')
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
			newModel := m.insertRunes(msg.Runes)
			newModel = newModel.markDirty()
			newModel = newModel.updateStatusBar()
			return newModel, nil
		}

	case tea.WindowSizeMsg:
		newModel := m
		newModel.width = msg.Width
		newModel.height = msg.Height
		newModel.viewport.Width = msg.Width
		newModel.viewport.Height = msg.Height - 1 // Leave room for status bar
		return newModel, nil
	}

	// Update viewport to keep cursor in view
	m = m.adjustViewport()

	// Update status bar counts
	m = m.updateStatusBar()

	// Sync viewport with buffer line count
	m = m.syncViewportLines()

	return m, nil
}

// insertTab inserts a tab at the cursor position.
func (m Model) insertTab() Model {
	newModel := m
	newModel.buffer.Insert(' ')
	newModel.buffer.Insert(' ')
	newModel.buffer.Insert(' ')
	newModel.buffer.Insert(' ')
	return newModel
}

// insertRune inserts a rune at the cursor position.
func (m Model) insertRune(r rune) Model {
	newModel := m
	newModel.buffer.Insert(r)
	return newModel
}

// insertRunes inserts runes at the cursor position.
func (m Model) insertRunes(runes []rune) Model {
	newModel := m
	for _, r := range runes {
		newModel.buffer.Insert(r)
	}
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
	return m.fileManager.Save(m.buffer.Content())
}

// adjustViewport adjusts the viewport to keep the cursor visible.
func (m Model) adjustViewport() Model {
	newModel := m
	_, cursorY := newModel.buffer.CursorPosition()

	viewportHeight := newModel.viewport.Height
	if viewportHeight <= 0 {
		return newModel
	}

	third := viewportHeight / 3
	middleTop := newModel.viewport.YOffset + third
	middleBottom := newModel.viewport.YOffset + viewportHeight - third

	if cursorY < middleTop {
		newOffset := cursorY - third
		if newOffset < 0 {
			newOffset = 0
		}
		newModel.viewport.SetYOffset(newOffset)
	} else if cursorY >= middleBottom {
		newOffset := cursorY - viewportHeight + third
		maxOffset := newModel.buffer.LineCount() - viewportHeight
		if maxOffset < 0 {
			maxOffset = 0
		}
		if newOffset > maxOffset {
			newOffset = maxOffset
		}
		newModel.viewport.SetYOffset(newOffset)
	}

	return newModel
}

// syncViewportLines syncs viewport content with buffer.
func (m Model) syncViewportLines() Model {
	newModel := m
	newModel.viewport.SetContent(newModel.buffer.Content())
	return newModel
}

// getVisibleLines returns the lines visible in the viewport.
func (m Model) getVisibleLines(height int) []string {
	lines := strings.Split(m.buffer.Content(), "\n")

	start := m.viewport.YOffset
	end := start + height
	if start >= len(lines) {
		return []string{}
	}

	if end > len(lines) {
		end = len(lines)
	}

	return lines[start:end]
}

// updateStatusBar updates the status bar information.
func (m Model) updateStatusBar() Model {
	newModel := m
	newModel.statusBar.charCount = newModel.buffer.CharCount()
	newModel.statusBar.lineCount = newModel.buffer.LineCount()
	return newModel
}

// updatePlaceholders re-parses placeholders from content.
func (m Model) updatePlaceholders() Model {
	newModel := m
	newModel.placeholders = newModel.placeholders.Update(editor.ParsePlaceholdersMsg{Content: newModel.buffer.Content()})
	return newModel
}

// navigateToNextPlaceholder moves to the next placeholder.
func (m Model) navigateToNextPlaceholder() Model {
	newModel := m
	placeholders := newModel.placeholders.Placeholders()

	if len(placeholders) == 0 {
		return newModel
	}

	cursorX, cursorY := newModel.buffer.CursorPosition()
	cursorPos := newModel.getAbsolutePosition(cursorX, cursorY)

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
		newModel.buffer.SetCursorPositionAbsolute(ph.StartPos, newModel.buffer.Content())
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

	cursorX, cursorY := newModel.buffer.CursorPosition()
	cursorPos := newModel.getAbsolutePosition(cursorX, cursorY)

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
		newModel.buffer.SetCursorPositionAbsolute(ph.StartPos, newModel.buffer.Content())
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
	newModel.viewport.Width = width
	newModel.viewport.Height = height - 1
	return newModel
}

// GetContent returns the current content of the workspace.
func (m Model) GetContent() string {
	return m.buffer.Content()
}

// SetContent sets the content of the workspace.
func (m Model) SetContent(content string) Model {
	newModel := m
	newModel.buffer.SetContent(content)
	newModel = newModel.updatePlaceholders()
	newModel = newModel.updateStatusBar()
	newModel = newModel.markDirty()
	return newModel
}

// GetCursorPosition returns the absolute cursor position in content.
func (m Model) GetCursorPosition() int {
	cursorX, cursorY := m.buffer.CursorPosition()
	return m.getAbsolutePosition(cursorX, cursorY)
}

// getAbsolutePosition converts (x, y) coordinates to absolute position in content.
func (m Model) getAbsolutePosition(x, y int) int {
	lines := strings.Split(m.buffer.Content(), "\n")
	pos := 0
	for i := 0; i < y && i < len(lines); i++ {
		pos += len(lines[i]) + 1 // +1 for newline
	}
	pos += x
	return pos
}

// InsertContent inserts content at a specific position.
func (m Model) InsertContent(position int, content string) Model {
	newModel := m
	// Convert position to line/column
	lines := strings.Split(newModel.buffer.Content(), "\n")
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

	newModel.buffer.SetContent(strings.Join(lines, "\n"))
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
