package editor

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Buffer manages text content and cursor position for editing operations.
// It provides rune-based storage for proper unicode support and tracks
// cursor position in (x, y) coordinates representing (column, line).
type Buffer struct {
	content string
	cursor  Cursor
}

// New creates a new empty buffer with cursor at origin.
func NewBuffer() *Buffer {
	return &Buffer{
		content: "",
		cursor:  NewCursor(),
	}
}

// NewWithContent creates a new buffer with initial content.
func NewBufferWithContent(content string) *Buffer {
	buf := &Buffer{
		content: content,
		cursor:  NewCursor(),
	}
	buf.cursor = buf.cursor.AdjustToLineLength(content)
	lines := buf.getLines()
	if len(lines) > 0 {
		lastLineRunes := []rune(lines[len(lines)-1])
		buf.cursor = buf.cursor.SetPosition(len(lastLineRunes), len(lines)-1)
	}
	return buf
}

// Insert inserts a rune at the current cursor position and moves cursor forward.
func (b *Buffer) Insert(r rune) error {
	if !utf8.ValidRune(r) {
		return fmt.Errorf("invalid rune: %c", r)
	}

	lines := b.getLines()
	if b.cursor.Y() < len(lines) {
		line := lines[b.cursor.Y()]
		runes := []rune(line)
		if b.cursor.X() > len(runes) {
			b.cursor = b.cursor.SetX(len(runes))
		}

		before := string(runes[:b.cursor.X()])
		after := string(runes[b.cursor.X():])
		newLine := before + string(r) + after
		lines[b.cursor.Y()] = newLine
		b.content = strings.Join(lines, "\n")

		if r == Newline {
			b.cursor = b.cursor.SetPosition(0, b.cursor.Y()+1)
		} else {
			b.cursor = b.cursor.SetX(b.cursor.X() + 1)
		}
	}
	return nil
}

// Delete removes the character before the cursor position and moves cursor backward.
// Returns error if cursor is at the beginning of the buffer.
func (b *Buffer) Delete() error {
	lines := b.getLines()

	if b.cursor.Y() < 0 || b.cursor.Y() >= len(lines) {
		return fmt.Errorf("cursor at invalid position")
	}

	if b.cursor.Y() == 0 && b.cursor.X() == 0 {
		return fmt.Errorf("cursor at start of buffer, nothing to delete")
	}

	line := lines[b.cursor.Y()]
	runes := []rune(line)

	if b.cursor.X() == 0 && b.cursor.Y() > 0 {
		prevLine := []rune(lines[b.cursor.Y()-1])
		mergedRunes := append(prevLine, runes...)
		mergedLine := string(mergedRunes)

		newLines := make([]string, 0, len(lines)-1)
		newLines = append(newLines, lines[:b.cursor.Y()-1]...)
		newLines = append(newLines, mergedLine)
		if b.cursor.Y()+1 < len(lines) {
			newLines = append(newLines, lines[b.cursor.Y()+1:]...)
		}

		b.cursor = b.cursor.SetX(len(prevLine))
		b.cursor = b.cursor.SetY(b.cursor.Y() - 1)
		b.content = strings.Join(newLines, "\n")
	} else {
		newRunes := append(runes[:b.cursor.X()-1], runes[b.cursor.X():]...)
		lines[b.cursor.Y()] = string(newRunes)
		b.cursor = b.cursor.SetX(b.cursor.X() - 1)
		b.content = strings.Join(lines, "\n")
	}

	return nil
}

// Content returns the current buffer content as a string.
func (b *Buffer) Content() string {
	return b.content
}

// CursorPosition returns the current cursor position as (x, y) coordinates.
func (b *Buffer) CursorPosition() (x, y int) {
	return b.cursor.Position()
}

// SetContent sets the buffer content and adjusts cursor position.
func (b *Buffer) SetContent(content string) {
	b.content = content
	b.cursor = b.cursor.AdjustToLineLength(content)
}

// SetCursorPosition sets the cursor to a specific position.
// Returns error if position is out of bounds.
func (b *Buffer) SetCursorPosition(x, y int) error {
	lines := b.getLines()
	if y < 0 || y >= len(lines) {
		return fmt.Errorf("line %d out of bounds (0-%d)", y, len(lines)-1)
	}

	if x < 0 || x > len(lines[y]) {
		return fmt.Errorf("column %d out of bounds for line %d (0-%d)", x, y, len(lines[y]))
	}

	b.cursor = b.cursor.SetPosition(x, y)
	return nil
}

// getPreviousRunePosition returns the byte position of the rune before pos.
// This handles multi-byte UTF-8 characters correctly.
func (b *Buffer) getPreviousRunePosition(pos int) int {
	if pos == 0 {
		return 0
	}

	for i := pos - 1; i >= 0; i-- {
		if utf8.RuneStart(b.content[i]) {
			return i
		}
	}
	return 0
}

// MoveUp moves the cursor up one line, preserving column position if possible.
// At first line, cursor position does not change.
func (b *Buffer) MoveUp() {
	lines := b.getLines()
	if b.cursor.Y() == 0 {
		return
	}

	newY := b.cursor.Y() - 1
	targetLine := []rune(lines[newY])
	newX := b.cursor.X()
	if newX > len(targetLine) {
		newX = len(targetLine)
	}
	b.cursor = b.cursor.SetPosition(newX, newY)
}

// MoveDown moves the cursor down one line, preserving column position if possible.
// At last line, cursor position does not change.
func (b *Buffer) MoveDown() {
	lines := b.getLines()
	if b.cursor.Y() == len(lines)-1 {
		return
	}

	newY := b.cursor.Y() + 1
	targetLine := []rune(lines[newY])
	newX := b.cursor.X()
	if newX > len(targetLine) {
		newX = len(targetLine)
	}
	b.cursor = b.cursor.SetPosition(newX, newY)
}

// MoveLeft moves the cursor left one character.
// At column 0, wraps to the end of the previous line.
// At start of buffer, cursor position does not change.
func (b *Buffer) MoveLeft() {
	lines := b.getLines()
	if b.cursor.Y() == 0 && b.cursor.X() == 0 {
		return
	}

	if b.cursor.X() == 0 {
		prevLine := []rune(lines[b.cursor.Y()-1])
		b.cursor = b.cursor.SetPosition(len(prevLine), b.cursor.Y()-1)
	} else {
		b.cursor = b.cursor.SetX(b.cursor.X() - 1)
	}
}

// MoveRight moves the cursor right one character.
// At end of line, wraps to the start of the next line.
// At end of buffer, cursor position does not change.
func (b *Buffer) MoveRight() {
	lines := b.getLines()
	currentLine := []rune(lines[b.cursor.Y()])

	if b.cursor.X() == len(currentLine) {
		if b.cursor.Y() == len(lines)-1 {
			return
		}
		b.cursor = b.cursor.SetPosition(0, b.cursor.Y()+1)
	} else {
		b.cursor = b.cursor.SetX(b.cursor.X() + 1)
	}
}

// MoveToLineStart moves the cursor to column 0 of the current line.
func (b *Buffer) MoveToLineStart() {
	b.cursor = b.cursor.SetX(0)
}

// MoveToLineEnd moves the cursor to the end of the current line.
func (b *Buffer) MoveToLineEnd() {
	lines := b.getLines()
	currentLine := []rune(lines[b.cursor.Y()])
	b.cursor = b.cursor.SetX(len(currentLine))
}

// getLines splits the buffer content into lines.
func (b *Buffer) getLines() []string {
	if b.content == "" {
		return []string{""}
	}
	return strings.Split(b.content, "\n")
}

// CharCount returns the total number of characters in the buffer.
// Counts runes, not bytes, to correctly handle multi-byte unicode characters.
func (b *Buffer) CharCount() int {
	return utf8.RuneCountInString(b.content)
}

// LineCount returns the total number of lines in the buffer.
// Empty buffer returns 1 (representing a single empty line).
func (b *Buffer) LineCount() int {
	if b.content == "" {
		return 1
	}
	return strings.Count(b.content, "\n") + 1
}

// SetCursorPositionAbsolute sets the cursor to an absolute position in content.
func (b *Buffer) SetCursorPositionAbsolute(pos int, content string) {
	lines := b.getLines()
	currentPos := 0

	for i, line := range lines {
		lineEnd := currentPos + len(line)

		if pos <= lineEnd {
			b.cursor = b.cursor.SetPosition(pos-currentPos, i)
			return
		}

		currentPos = lineEnd + 1 // +1 for newline
	}

	// If position is beyond content, set to end
	if len(lines) > 0 {
		b.cursor = b.cursor.SetPosition(len(lines[len(lines)-1]), len(lines)-1)
	}
}
