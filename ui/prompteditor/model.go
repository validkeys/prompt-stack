package prompteditor

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/internal/editor"
	"github.com/kyledavis/prompt-stack/internal/prompt"
	"github.com/kyledavis/prompt-stack/ui/theme"
)

// Mode represents the editor mode
type Mode int

const (
	ModeEdit Mode = iota
	ModePreview
)

// Model represents the prompt editor
type Model struct {
	prompt          *prompt.Prompt
	filePath        string
	mode            Mode
	content         string
	cursor          cursor
	viewport        viewport.Model
	previewViewport viewport.Model
	width           int
	height          int
	isDirty         bool
	saveStatus      string // "saving", "saved", "error", ""
	statusMessage   string
	undoStack       *editor.UndoStack
	onSave          func(*prompt.Prompt) tea.Cmd
	onCancel        func() tea.Cmd
	storage         *prompt.Storage
	glamourRenderer *glamour.TermRenderer
}

type cursor struct {
	x int
	y int
}

// NewModel creates a new prompt editor model
func NewModel(p *prompt.Prompt, storage *prompt.Storage) Model {
	// Initialize viewport for edit mode
	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle()

	// Initialize viewport for preview mode
	previewVp := viewport.New(0, 0)
	previewVp.Style = lipgloss.NewStyle()

	// Initialize glamour renderer with minimal styling
	glamourRenderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		// Fallback to plain text if glamour fails to initialize
		glamourRenderer = nil
	}

	return Model{
		prompt:          p,
		filePath:        p.FilePath,
		mode:            ModeEdit,
		content:         p.Content,
		cursor:          cursor{x: 0, y: 0},
		viewport:        vp,
		previewViewport: previewVp,
		isDirty:         false,
		saveStatus:      "",
		statusMessage:   "",
		undoStack:       editor.NewUndoStack(),
		storage:         storage,
		glamourRenderer: glamourRenderer,
	}
}

// SetOnSave sets the callback for saving
func (m *Model) SetOnSave(fn func(*prompt.Prompt) tea.Cmd) {
	m.onSave = fn
}

// SetOnCancel sets the callback for cancellation
func (m *Model) SetOnCancel(fn func() tea.Cmd) {
	m.onCancel = fn
}

// SetSize sets the size of the model
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Update viewport sizes
	availableHeight := height - 2 // Leave room for header and status bar
	m.viewport.Width = width
	m.viewport.Height = availableHeight
	m.previewViewport.Width = width
	m.previewViewport.Height = availableHeight
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
		// Handle preview mode (read-only)
		if m.mode == ModePreview {
			return m.handlePreviewKeys(msg)
		}

		// Handle edit mode
		return m.handleEditKeys(msg)

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)

	case saveSuccessMsg:
		m.saveStatus = "saved"
		m.isDirty = false
		m.statusMessage = ""
		// Clear saved status after 2 seconds
		return m, tea.Tick(2*1000*1000*1000, func(t time.Time) tea.Msg {
			return clearSaveStatusMsg{}
		})

	case saveErrorMsg:
		m.saveStatus = "error"
		m.statusMessage = fmt.Sprintf("Save failed: %v", msg.err)
		return m, nil

	case clearSaveStatusMsg:
		m.saveStatus = ""
		m.statusMessage = ""
		return m, nil
	}

	return m, cmd
}

// handleEditKeys handles key events in edit mode
func (m Model) handleEditKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		// Check for unsaved changes before quitting
		if m.isDirty {
			m.statusMessage = "Unsaved changes. Press Ctrl+S to save or Esc to cancel."
			return m, nil
		}
		if m.onCancel != nil {
			return m, m.onCancel()
		}
		return m, tea.Quit

	case tea.KeyCtrlS:
		// Save prompt
		newModel, cmd := m.savePrompt()
		return newModel, cmd

	case tea.KeyCtrlP:
		// Toggle edit/preview mode
		m.toggleMode()
		return m, nil

	case tea.KeyEsc:
		// Cancel or exit
		if m.statusMessage != "" {
			m.statusMessage = ""
			return m, nil
		}
		if m.isDirty {
			m.statusMessage = "Unsaved changes. Press Ctrl+S to save or Esc to cancel."
			return m, nil
		}
		if m.onCancel != nil {
			return m, m.onCancel()
		}
		return m, tea.Quit

	case tea.KeyCtrlZ:
		// Undo
		if m.undoStack.CanUndo() {
			m.undo()
		}
		return m, nil

	case tea.KeyCtrlY:
		// Redo
		if m.undoStack.CanRedo() {
			m.redo()
		}
		return m, nil

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

	case tea.KeyEnter:
		m.insertNewline()
		m.markDirty()

	case tea.KeyTab:
		m.insertTab()
		m.markDirty()

	case tea.KeyRunes:
		m.insertRune(msg.Runes)
		m.markDirty()
	}

	// Update viewport to keep cursor in view
	m.adjustViewport()

	return m, nil
}

// handlePreviewKeys handles key events in preview mode
func (m Model) handlePreviewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		if m.onCancel != nil {
			return m, m.onCancel()
		}
		return m, tea.Quit

	case tea.KeyCtrlP:
		// Toggle back to edit mode
		m.toggleMode()
		return m, nil

	case tea.KeyEsc:
		// Toggle back to edit mode
		m.toggleMode()
		return m, nil

	case tea.KeyUp:
		m.previewViewport.LineUp(1)
		return m, nil

	case tea.KeyDown:
		m.previewViewport.LineDown(1)
		return m, nil

	case tea.KeyPgUp:
		m.previewViewport.HalfViewUp()
		return m, nil

	case tea.KeyPgDown:
		m.previewViewport.HalfViewDown()
		return m, nil
	}

	return m, nil
}

// toggleMode toggles between edit and preview modes
func (m *Model) toggleMode() {
	if m.mode == ModeEdit {
		m.mode = ModePreview
		// Update preview content
		m.updatePreview()
	} else {
		m.mode = ModeEdit
	}
}

// updatePreview updates the preview content
func (m *Model) updatePreview() {
	var previewContent string

	if m.glamourRenderer != nil {
		// Render markdown using glamour
		rendered, err := m.glamourRenderer.Render(m.content)
		if err != nil {
			// Fallback to plain text on error
			previewContent = m.content
		} else {
			previewContent = rendered
		}
	} else {
		// Fallback to plain text if glamour is not available
		previewContent = m.content
	}

	m.previewViewport.SetContent(previewContent)
	m.previewViewport.GotoTop()
}

// View renders the model
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// Render header
	header := m.renderHeader()

	// Render content based on mode
	var content string
	if m.mode == ModeEdit {
		content = m.renderEditMode()
	} else {
		content = m.renderPreviewMode()
	}

	// Render status bar
	statusBar := m.renderStatusBar()

	// Combine all parts
	return lipgloss.JoinVertical(lipgloss.Left, header, content, statusBar)
}

// renderHeader renders the header
func (m *Model) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Width(m.width).
		Padding(0, 1)

	// Build header text
	var parts []string
	parts = append(parts, fmt.Sprintf("Prompt: %s", m.prompt.Title))
	if m.prompt.Category != "" {
		parts = append(parts, fmt.Sprintf("Category: %s", m.prompt.Category))
	}

	headerText := strings.Join(parts, " | ")
	return headerStyle.Render(headerText)
}

// renderEditMode renders the editor in edit mode
func (m Model) renderEditMode() string {
	// Get visible lines
	lines := strings.Split(m.content, "\n")
	visibleLines := m.getVisibleLines(lines)

	// Render lines with cursor
	renderedLines := make([]string, len(visibleLines))
	for i, line := range visibleLines {
		if i == m.cursor.y {
			renderedLines[i] = m.renderCursorLine(line)
		} else {
			renderedLines[i] = line
		}
	}

	// Combine rendered lines
	content := strings.Join(renderedLines, "\n")

	// Style editor
	editorStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-2).
		Padding(0, 1)

	return editorStyle.Render(content)
}

// renderPreviewMode renders the editor in preview mode
func (m Model) renderPreviewMode() string {
	// Style preview
	previewStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-2).
		Padding(0, 1)

	return previewStyle.Render(m.previewViewport.View())
}

// renderCursorLine renders a line with cursor position highlighted
func (m Model) renderCursorLine(line string) string {
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

// renderStatusBar renders the status bar
func (m Model) renderStatusBar() string {
	statusStyle := theme.StatusStyle().
		Width(m.width).
		Height(1)

	// Build status message
	var parts []string

	// Mode indicator
	if m.mode == ModeEdit {
		parts = append(parts, "EDIT")
	} else {
		parts = append(parts, "PREVIEW")
	}

	// Save status
	if m.saveStatus == "saving" {
		parts = append(parts, "Saving...")
	} else if m.saveStatus == "saved" {
		parts = append(parts, "Saved ✓")
	} else if m.saveStatus == "error" {
		parts = append(parts, "Save failed")
	} else if m.isDirty {
		parts = append(parts, "Unsaved changes")
	}

	// Undo/Redo indicators
	if m.undoStack.CanUndo() {
		parts = append(parts, "Undo: Ctrl+Z")
	}
	if m.undoStack.CanRedo() {
		parts = append(parts, "Redo: Ctrl+Y")
	}

	// Character and line counts
	charCount := len(m.content)
	lineCount := strings.Count(m.content, "\n") + 1
	parts = append(parts, fmt.Sprintf("%d chars, %d lines", charCount, lineCount))

	// Custom message
	if m.statusMessage != "" {
		parts = append(parts, m.statusMessage)
	}

	// Help text
	if m.mode == ModeEdit {
		parts = append(parts, "Ctrl+S: Save | Ctrl+P: Preview | Esc: Exit")
	} else {
		parts = append(parts, "Ctrl+P: Edit | Esc: Exit")
	}

	// Join with separator
	statusText := strings.Join(parts, " | ")

	return statusStyle.Render(statusText)
}

// Cursor movement methods
func (m *Model) moveCursorUp() {
	if m.cursor.y > 0 {
		m.cursor.y--
		lines := strings.Split(m.content, "\n")
		if m.cursor.y < len(lines) {
			lineLen := len(lines[m.cursor.y])
			if m.cursor.x > lineLen {
				m.cursor.x = lineLen
			}
		}
	}
}

func (m *Model) moveCursorDown() {
	lines := strings.Split(m.content, "\n")
	if m.cursor.y < len(lines)-1 {
		m.cursor.y++
		if m.cursor.y < len(lines) {
			lineLen := len(lines[m.cursor.y])
			if m.cursor.x > lineLen {
				m.cursor.x = lineLen
			}
		}
	}
}

func (m *Model) moveCursorLeft() {
	if m.cursor.x > 0 {
		m.cursor.x--
	} else if m.cursor.y > 0 {
		m.cursor.y--
		lines := strings.Split(m.content, "\n")
		if m.cursor.y < len(lines) {
			m.cursor.x = len(lines[m.cursor.y])
		}
	}
}

func (m *Model) moveCursorRight() {
	lines := strings.Split(m.content, "\n")
	if m.cursor.y < len(lines) {
		lineLen := len(lines[m.cursor.y])
		if m.cursor.x < lineLen {
			m.cursor.x++
		} else if m.cursor.y < len(lines)-1 {
			m.cursor.y++
			m.cursor.x = 0
		}
	}
}

// Editing methods
func (m *Model) backspace() {
	if m.cursor.x == 0 && m.cursor.y == 0 {
		return
	}

	cursorPos := m.getCursorPosition()
	lines := strings.Split(m.content, "\n")

	var deletedContent string

	if m.cursor.x == 0 {
		if m.cursor.y > 0 {
			prevLine := lines[m.cursor.y-1]
			currentLine := lines[m.cursor.y]
			deletedContent = "\n"
			m.cursor.x = len(prevLine)
			lines[m.cursor.y-1] = prevLine + currentLine
			lines = append(lines[:m.cursor.y], lines[m.cursor.y+1:]...)
			m.cursor.y--
		}
	} else {
		if m.cursor.y < len(lines) {
			line := lines[m.cursor.y]
			if m.cursor.x <= len(line) {
				deletedContent = string(line[m.cursor.x-1])
				lines[m.cursor.y] = line[:m.cursor.x-1] + line[m.cursor.x:]
				m.cursor.x--
			}
		}
	}

	m.content = strings.Join(lines, "\n")

	if deletedContent != "" {
		action := editor.CreateBackspaceAction(deletedContent, cursorPos, m.cursor.x, m.cursor.y)
		m.undoStack.PushAction(action)
	}
}

func (m *Model) insertNewline() {
	cursorPos := m.getCursorPosition()
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

	action := editor.CreateNewlineAction(cursorPos, m.cursor.x, m.cursor.y)
	m.undoStack.PushAction(action)
}

func (m *Model) insertTab() {
	m.insertRune([]rune{' ', ' ', ' ', ' '})
}

func (m *Model) insertRune(runes []rune) {
	cursorPos := m.getCursorPosition()
	lines := strings.Split(m.content, "\n")

	if m.cursor.y < len(lines) {
		line := lines[m.cursor.y]
		before := line[:m.cursor.x]
		after := line[m.cursor.x:]

		lines[m.cursor.y] = before + string(runes) + after
		m.cursor.x += len(runes)
	}

	m.content = strings.Join(lines, "\n")

	action := editor.CreateInsertAction(string(runes), cursorPos, m.cursor.x, m.cursor.y)
	m.undoStack.PushAction(action)
}

// Helper methods
func (m *Model) markDirty() {
	m.isDirty = true
}

func (m *Model) adjustViewport() {
	// Simple viewport adjustment
	availableHeight := m.height - 2
	if availableHeight <= 0 {
		return
	}

	third := availableHeight / 3

	if m.cursor.y < m.viewport.YOffset+third {
		m.viewport.GotoTop()
		m.viewport.LineDown(max(0, m.cursor.y-third))
	} else if m.cursor.y > m.viewport.YOffset+availableHeight-third {
		m.viewport.GotoTop()
		m.viewport.LineDown(max(0, m.cursor.y-availableHeight+third))
	}
}

func (m *Model) getVisibleLines(lines []string) []string {
	start := m.viewport.YOffset
	end := min(start+m.viewport.Height, len(lines))

	if start >= len(lines) {
		return []string{}
	}

	return lines[start:end]
}

func (m *Model) getCursorPosition() int {
	lines := strings.Split(m.content, "\n")
	pos := 0

	for i := 0; i < m.cursor.y && i < len(lines); i++ {
		pos += len(lines[i]) + 1
	}

	if m.cursor.y < len(lines) {
		pos += m.cursor.x
	}

	return pos
}

func (m *Model) undo() {
	action, ok := m.undoStack.Undo()
	if !ok {
		return
	}

	switch action.Type {
	case editor.ActionTypeInsert, editor.ActionTypePaste, editor.ActionTypePromptInsert, editor.ActionTypePlaceholderFill, editor.ActionTypeNewline:
		m.deleteContent(action.Position, len(action.Content))
	case editor.ActionTypeDelete, editor.ActionTypeBackspace:
		m.insertContent(action.Position, action.Content)
	}

	m.cursor.x = action.CursorPos.X
	m.cursor.y = action.CursorPos.Y
	m.markDirty()
}

func (m *Model) redo() {
	action, ok := m.undoStack.Redo()
	if !ok {
		return
	}

	switch action.Type {
	case editor.ActionTypeInsert, editor.ActionTypePaste, editor.ActionTypePromptInsert, editor.ActionTypePlaceholderFill, editor.ActionTypeNewline:
		m.insertContent(action.Position, action.Content)
	case editor.ActionTypeDelete, editor.ActionTypeBackspace:
		m.deleteContent(action.Position, len(action.Content))
	}

	m.cursor.x = action.CursorPos.X
	m.cursor.y = action.CursorPos.Y
	m.markDirty()
}

func (m *Model) insertContent(position int, content string) {
	lines := strings.Split(m.content, "\n")
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
		currentPos = lineEnd + 1
	}

	if targetLine < len(lines) {
		line := lines[targetLine]
		before := line[:targetCol]
		after := line[targetCol:]
		lines[targetLine] = before + content + after
	}

	m.content = strings.Join(lines, "\n")
}

func (m *Model) deleteContent(position int, length int) {
	lines := strings.Split(m.content, "\n")
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
		currentPos = lineEnd + 1
	}

	if targetLine < len(lines) {
		line := lines[targetLine]
		if targetCol+length <= len(line) {
			before := line[:targetCol]
			after := line[targetCol+length:]
			lines[targetLine] = before + after
		}
	}

	m.content = strings.Join(lines, "\n")
}

// savePrompt saves the prompt to storage
func (m Model) savePrompt() (tea.Model, tea.Cmd) {
	m.saveStatus = "saving"
	m.statusMessage = ""

	// Update prompt content
	m.prompt.Content = m.content

	// Re-parse placeholders
	m.prompt.Placeholders = prompt.ParsePlaceholders(m.content)

	// Save to storage
	return m, tea.Cmd(func() tea.Msg {
		// Create prompt data for saving
		promptData := &prompt.PromptData{
			Title:       m.prompt.Title,
			Description: m.prompt.Description,
			Tags:        m.prompt.Tags,
			Category:    m.prompt.Category,
			Content:     m.content,
		}

		// Save using storage
		filePath, err := m.storage.SavePrompt(promptData)
		if err != nil {
			return saveErrorMsg{err}
		}

		// Update file path
		m.prompt.FilePath = filePath

		// Call save callback if provided
		if m.onSave != nil {
			return m.onSave(m.prompt)
		}

		return saveSuccessMsg{}
	})
}

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

// Messages
type saveSuccessMsg struct{}
type saveErrorMsg struct {
	err error
}
type clearSaveStatusMsg struct{}
