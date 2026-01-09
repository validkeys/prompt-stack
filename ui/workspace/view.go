package workspace

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/ui/theme"
)

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
	_, cursorY := m.buffer.CursorPosition()

	for i, line := range lines {
		visibleY := m.viewport.YOffset + i
		if visibleY == cursorY {
			// Cursor line - render with cursor
			renderedLines[i] = m.renderCursorLine(lines, i)
		} else {
			// Non-cursor line - render with placeholder highlighting
			renderedLines[i] = m.renderLineWithPlaceholders(line, visibleY)
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

// renderCursorLine renders the line with cursor position highlighted.
func (m Model) renderCursorLine(lines []string, visibleIndex int) string {
	cursorX, cursorY := m.buffer.CursorPosition()
	visibleY := m.viewport.YOffset + visibleIndex

	if visibleY != cursorY || visibleIndex >= len(lines) {
		if visibleIndex < len(lines) {
			return lines[visibleIndex]
		}
		return ""
	}

	line := lines[visibleIndex]

	// If in placeholder edit mode, show edit value instead of placeholder
	if m.placeholders.IsEditing() && m.placeholders.Active() != nil {
		ph := m.placeholders.Active()
		if ph.Type == "text" {
			// Calculate line start position
			lineStartPos := 0
			allLines := strings.Split(m.buffer.Content(), "\n")
			for i := 0; i < visibleY && i < len(allLines); i++ {
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
				cursorPosInEdit := cursorX - phStart
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

	if cursorX > len(line) {
		return line
	}

	before := line[:cursorX]
	after := line[cursorX:]

	// Style cursor
	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.CursorBackground)).
		Foreground(lipgloss.Color(theme.CursorForeground))

	if cursorX < len(line) {
		return before + cursorStyle.Render(string(line[cursorX])) + after
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
	lines := strings.Split(m.buffer.Content(), "\n")
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
