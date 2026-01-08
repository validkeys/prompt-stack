package editor

import tea "github.com/charmbracelet/bubbletea"

// Model represents cursor position and movement
type Cursor struct {
	x    int // Column position
	y    int // Line position
	line int // Current line index (for tracking)
}

// NewCursor creates a new cursor at origin
func NewCursor() Cursor {
	return Cursor{x: 0, y: 0, line: 0}
}

// Update handles cursor-related messages
func (m Cursor) Update(msg tea.Msg) Cursor {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			return m.MoveUp()
		case tea.KeyDown:
			return m.MoveDown()
		case tea.KeyLeft:
			return m.MoveLeft()
		case tea.KeyRight:
			return m.MoveRight()
		}
	}
	return m
}

// MoveUp moves cursor up one line
func (m Cursor) MoveUp() Cursor {
	newModel := m
	newModel.y--
	if newModel.y < 0 {
		newModel.y = 0
	}
	return newModel
}

// MoveDown moves cursor down one line
func (m Cursor) MoveDown() Cursor {
	newModel := m
	newModel.y++
	return newModel
}

// MoveLeft moves cursor left one character
func (m Cursor) MoveLeft() Cursor {
	newModel := m
	newModel.x--
	if newModel.x < 0 {
		newModel.x = 0
	}
	return newModel
}

// MoveRight moves cursor right one character
func (m Cursor) MoveRight() Cursor {
	newModel := m
	newModel.x++
	return newModel
}

// Position returns current cursor position
func (m Cursor) Position() (x, y int) {
	return m.x, m.y
}

// SetPosition sets the cursor to a specific position
func (m Cursor) SetPosition(x, y int) Cursor {
	newModel := m
	newModel.x = x
	newModel.y = y
	return newModel
}

// SetX sets the cursor X position
func (m Cursor) SetX(x int) Cursor {
	newModel := m
	newModel.x = x
	return newModel
}

// SetY sets the cursor Y position
func (m Cursor) SetY(y int) Cursor {
	newModel := m
	newModel.y = y
	return newModel
}

// X returns the cursor X position
func (m Cursor) X() int {
	return m.x
}

// Y returns the cursor Y position
func (m Cursor) Y() int {
	return m.y
}

// MoveToPosition moves cursor to an absolute position in content
// This is a helper that converts absolute position to line/column
func (m Cursor) MoveToPosition(pos int, content string) Cursor {
	newModel := m
	lines := splitLines(content)
	currentPos := 0

	for i, line := range lines {
		lineEnd := currentPos + len(line)

		if pos <= lineEnd {
			newModel.y = i
			newModel.x = pos - currentPos
			return newModel
		}

		currentPos = lineEnd + 1 // +1 for newline
	}

	// If position is beyond content, set to end
	if len(lines) > 0 {
		newModel.y = len(lines) - 1
		newModel.x = len(lines[len(lines)-1])
	}
	return newModel
}

// GetAbsolutePosition returns the absolute cursor position in content
func (m Cursor) GetAbsolutePosition(content string) int {
	lines := splitLines(content)
	pos := 0

	for i := 0; i < m.y && i < len(lines); i++ {
		pos += len(lines[i]) + 1 // +1 for newline
	}

	if m.y < len(lines) {
		pos += m.x
	}

	return pos
}

// AdjustToLineLength adjusts cursor X position to fit within line length
func (m Cursor) AdjustToLineLength(content string) Cursor {
	newModel := m
	lines := splitLines(content)

	if newModel.y < len(lines) {
		lineLen := len(lines[newModel.y])
		if newModel.x > lineLen {
			newModel.x = lineLen
		}
	}

	return newModel
}

// splitLines is a helper to split content into lines
func splitLines(content string) []string {
	if content == "" {
		return []string{""}
	}
	return splitLinesHelper(content)
}

// splitLinesHelper splits content by newlines
func splitLinesHelper(content string) []string {
	var lines []string
	start := 0

	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			lines = append(lines, content[start:i])
			start = i + 1
		}
	}

	// Add the last line
	lines = append(lines, content[start:])

	return lines
}
