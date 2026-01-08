package editor

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewCursor(t *testing.T) {
	cursor := NewCursor()

	if cursor.x != 0 {
		t.Errorf("expected x=0, got x=%d", cursor.x)
	}
	if cursor.y != 0 {
		t.Errorf("expected y=0, got y=%d", cursor.y)
	}
	if cursor.line != 0 {
		t.Errorf("expected line=0, got line=%d", cursor.line)
	}
}

func TestCursorMoveUp(t *testing.T) {
	tests := []struct {
		name     string
		initial  Cursor
		expected Cursor
	}{
		{
			name:     "move from y=5 to y=4",
			initial:  Cursor{x: 10, y: 5},
			expected: Cursor{x: 10, y: 4},
		},
		{
			name:     "move from y=1 to y=0",
			initial:  Cursor{x: 10, y: 1},
			expected: Cursor{x: 10, y: 0},
		},
		{
			name:     "stay at y=0",
			initial:  Cursor{x: 10, y: 0},
			expected: Cursor{x: 10, y: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.MoveUp()
			if result.x != tt.expected.x || result.y != tt.expected.y {
				t.Errorf("expected (%d, %d), got (%d, %d)", tt.expected.x, tt.expected.y, result.x, result.y)
			}
		})
	}
}

func TestCursorMoveDown(t *testing.T) {
	tests := []struct {
		name     string
		initial  Cursor
		expected Cursor
	}{
		{
			name:     "move from y=0 to y=1",
			initial:  Cursor{x: 10, y: 0},
			expected: Cursor{x: 10, y: 1},
		},
		{
			name:     "move from y=5 to y=6",
			initial:  Cursor{x: 10, y: 5},
			expected: Cursor{x: 10, y: 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.MoveDown()
			if result.x != tt.expected.x || result.y != tt.expected.y {
				t.Errorf("expected (%d, %d), got (%d, %d)", tt.expected.x, tt.expected.y, result.x, result.y)
			}
		})
	}
}

func TestCursorMoveLeft(t *testing.T) {
	tests := []struct {
		name     string
		initial  Cursor
		expected Cursor
	}{
		{
			name:     "move from x=10 to x=9",
			initial:  Cursor{x: 10, y: 5},
			expected: Cursor{x: 9, y: 5},
		},
		{
			name:     "move from x=1 to x=0",
			initial:  Cursor{x: 1, y: 5},
			expected: Cursor{x: 0, y: 5},
		},
		{
			name:     "stay at x=0",
			initial:  Cursor{x: 0, y: 5},
			expected: Cursor{x: 0, y: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.MoveLeft()
			if result.x != tt.expected.x || result.y != tt.expected.y {
				t.Errorf("expected (%d, %d), got (%d, %d)", tt.expected.x, tt.expected.y, result.x, result.y)
			}
		})
	}
}

func TestCursorMoveRight(t *testing.T) {
	tests := []struct {
		name     string
		initial  Cursor
		expected Cursor
	}{
		{
			name:     "move from x=0 to x=1",
			initial:  Cursor{x: 0, y: 5},
			expected: Cursor{x: 1, y: 5},
		},
		{
			name:     "move from x=10 to x=11",
			initial:  Cursor{x: 10, y: 5},
			expected: Cursor{x: 11, y: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.MoveRight()
			if result.x != tt.expected.x || result.y != tt.expected.y {
				t.Errorf("expected (%d, %d), got (%d, %d)", tt.expected.x, tt.expected.y, result.x, result.y)
			}
		})
	}
}

func TestCursorPosition(t *testing.T) {
	cursor := Cursor{x: 10, y: 5}

	x, y := cursor.Position()

	if x != 10 {
		t.Errorf("expected x=10, got x=%d", x)
	}
	if y != 5 {
		t.Errorf("expected y=5, got y=%d", y)
	}
}

func TestCursorSetPosition(t *testing.T) {
	cursor := NewCursor()
	result := cursor.SetPosition(15, 7)

	if result.x != 15 {
		t.Errorf("expected x=15, got x=%d", result.x)
	}
	if result.y != 7 {
		t.Errorf("expected y=7, got y=%d", result.y)
	}
}

func TestCursorSetX(t *testing.T) {
	cursor := NewCursor()
	result := cursor.SetX(20)

	if result.x != 20 {
		t.Errorf("expected x=20, got x=%d", result.x)
	}
	if result.y != 0 {
		t.Errorf("expected y=0, got y=%d", result.y)
	}
}

func TestCursorSetY(t *testing.T) {
	cursor := NewCursor()
	result := cursor.SetY(10)

	if result.x != 0 {
		t.Errorf("expected x=0, got x=%d", result.x)
	}
	if result.y != 10 {
		t.Errorf("expected y=10, got y=%d", result.y)
	}
}

func TestCursorX(t *testing.T) {
	cursor := Cursor{x: 15, y: 5}

	if cursor.X() != 15 {
		t.Errorf("expected X()=15, got X()=%d", cursor.X())
	}
}

func TestCursorY(t *testing.T) {
	cursor := Cursor{x: 15, y: 5}

	if cursor.Y() != 5 {
		t.Errorf("expected Y()=5, got Y()=%d", cursor.Y())
	}
}

func TestCursorUpdate(t *testing.T) {
	tests := []struct {
		name     string
		initial  Cursor
		msg      tea.Msg
		expected Cursor
	}{
		{
			name:     "key up",
			initial:  Cursor{x: 10, y: 5},
			msg:      tea.KeyMsg{Type: tea.KeyUp},
			expected: Cursor{x: 10, y: 4},
		},
		{
			name:     "key down",
			initial:  Cursor{x: 10, y: 5},
			msg:      tea.KeyMsg{Type: tea.KeyDown},
			expected: Cursor{x: 10, y: 6},
		},
		{
			name:     "key left",
			initial:  Cursor{x: 10, y: 5},
			msg:      tea.KeyMsg{Type: tea.KeyLeft},
			expected: Cursor{x: 9, y: 5},
		},
		{
			name:     "key right",
			initial:  Cursor{x: 10, y: 5},
			msg:      tea.KeyMsg{Type: tea.KeyRight},
			expected: Cursor{x: 11, y: 5},
		},
		{
			name:     "other key",
			initial:  Cursor{x: 10, y: 5},
			msg:      tea.KeyMsg{Type: tea.KeyEnter},
			expected: Cursor{x: 10, y: 5},
		},
		{
			name:     "non-key message",
			initial:  Cursor{x: 10, y: 5},
			msg:      tea.WindowSizeMsg{Width: 80, Height: 24},
			expected: Cursor{x: 10, y: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.Update(tt.msg)
			if result.x != tt.expected.x || result.y != tt.expected.y {
				t.Errorf("expected (%d, %d), got (%d, %d)", tt.expected.x, tt.expected.y, result.x, result.y)
			}
		})
	}
}

func TestCursorMoveToPosition(t *testing.T) {
	tests := []struct {
		name     string
		initial  Cursor
		pos      int
		content  string
		expected Cursor
	}{
		{
			name:     "move to position 5 in single line",
			initial:  Cursor{x: 0, y: 0},
			pos:      5,
			content:  "Hello World",
			expected: Cursor{x: 5, y: 0},
		},
		{
			name:     "move to position 0",
			initial:  Cursor{x: 10, y: 5},
			pos:      0,
			content:  "Hello\nWorld",
			expected: Cursor{x: 0, y: 0},
		},
		{
			name:     "move to position 7 (after newline)",
			initial:  Cursor{x: 0, y: 0},
			pos:      7,
			content:  "Hello\nWorld",
			expected: Cursor{x: 1, y: 1},
		},
		{
			name:     "move to end of content",
			initial:  Cursor{x: 0, y: 0},
			pos:      11,
			content:  "Hello\nWorld",
			expected: Cursor{x: 5, y: 1},
		},
		{
			name:     "move beyond content",
			initial:  Cursor{x: 0, y: 0},
			pos:      100,
			content:  "Hello\nWorld",
			expected: Cursor{x: 5, y: 1},
		},
		{
			name:     "empty content",
			initial:  Cursor{x: 0, y: 0},
			pos:      0,
			content:  "",
			expected: Cursor{x: 0, y: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.MoveToPosition(tt.pos, tt.content)
			if result.x != tt.expected.x || result.y != tt.expected.y {
				t.Errorf("expected (%d, %d), got (%d, %d)", tt.expected.x, tt.expected.y, result.x, result.y)
			}
		})
	}
}

func TestCursorGetAbsolutePosition(t *testing.T) {
	tests := []struct {
		name     string
		cursor   Cursor
		content  string
		expected int
	}{
		{
			name:     "position at start",
			cursor:   Cursor{x: 0, y: 0},
			content:  "Hello\nWorld",
			expected: 0,
		},
		{
			name:     "position in first line",
			cursor:   Cursor{x: 3, y: 0},
			content:  "Hello\nWorld",
			expected: 3,
		},
		{
			name:     "position at end of first line",
			cursor:   Cursor{x: 5, y: 0},
			content:  "Hello\nWorld",
			expected: 5,
		},
		{
			name:     "position at start of second line",
			cursor:   Cursor{x: 0, y: 1},
			content:  "Hello\nWorld",
			expected: 6,
		},
		{
			name:     "position in second line",
			cursor:   Cursor{x: 3, y: 1},
			content:  "Hello\nWorld",
			expected: 9,
		},
		{
			name:     "position at end of content",
			cursor:   Cursor{x: 5, y: 1},
			content:  "Hello\nWorld",
			expected: 11,
		},
		{
			name:     "empty content",
			cursor:   Cursor{x: 0, y: 0},
			content:  "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cursor.GetAbsolutePosition(tt.content)
			if result != tt.expected {
				t.Errorf("expected position %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestCursorAdjustToLineLength(t *testing.T) {
	tests := []struct {
		name     string
		initial  Cursor
		content  string
		expected Cursor
	}{
		{
			name:     "cursor within line length",
			initial:  Cursor{x: 3, y: 0},
			content:  "Hello\nWorld",
			expected: Cursor{x: 3, y: 0},
		},
		{
			name:     "cursor at line length",
			initial:  Cursor{x: 5, y: 0},
			content:  "Hello\nWorld",
			expected: Cursor{x: 5, y: 0},
		},
		{
			name:     "cursor beyond line length",
			initial:  Cursor{x: 10, y: 0},
			content:  "Hello\nWorld",
			expected: Cursor{x: 5, y: 0},
		},
		{
			name:     "cursor on second line within length",
			initial:  Cursor{x: 3, y: 1},
			content:  "Hello\nWorld",
			expected: Cursor{x: 3, y: 1},
		},
		{
			name:     "cursor on second line beyond length",
			initial:  Cursor{x: 10, y: 1},
			content:  "Hello\nWorld",
			expected: Cursor{x: 5, y: 1},
		},
		{
			name:     "cursor beyond line count",
			initial:  Cursor{x: 10, y: 10},
			content:  "Hello\nWorld",
			expected: Cursor{x: 10, y: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.AdjustToLineLength(tt.content)
			if result.x != tt.expected.x || result.y != tt.expected.y {
				t.Errorf("expected (%d, %d), got (%d, %d)", tt.expected.x, tt.expected.y, result.x, result.y)
			}
		})
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "empty string",
			content:  "",
			expected: []string{""},
		},
		{
			name:     "single line",
			content:  "Hello",
			expected: []string{"Hello"},
		},
		{
			name:     "two lines",
			content:  "Hello\nWorld",
			expected: []string{"Hello", "World"},
		},
		{
			name:     "multiple lines",
			content:  "Line1\nLine2\nLine3",
			expected: []string{"Line1", "Line2", "Line3"},
		},
		{
			name:     "trailing newline",
			content:  "Hello\n",
			expected: []string{"Hello", ""},
		},
		{
			name:     "multiple trailing newlines",
			content:  "Hello\n\n\n",
			expected: []string{"Hello", "", "", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitLines(tt.content)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d lines, got %d", len(tt.expected), len(result))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("line %d: expected '%s', got '%s'", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

func TestCursorImmutability(t *testing.T) {
	// Test that cursor operations return new instances
	cursor := Cursor{x: 10, y: 5}

	// MoveUp should not modify original
	result := cursor.MoveUp()
	if cursor.x != 10 || cursor.y != 5 {
		t.Error("MoveUp modified original cursor")
	}
	if result.x != 10 || result.y != 4 {
		t.Error("MoveUp did not return correct result")
	}

	// MoveDown should not modify original
	result = cursor.MoveDown()
	if cursor.x != 10 || cursor.y != 5 {
		t.Error("MoveDown modified original cursor")
	}
	if result.x != 10 || result.y != 6 {
		t.Error("MoveDown did not return correct result")
	}

	// SetPosition should not modify original
	result = cursor.SetPosition(20, 10)
	if cursor.x != 10 || cursor.y != 5 {
		t.Error("SetPosition modified original cursor")
	}
	if result.x != 20 || result.y != 10 {
		t.Error("SetPosition did not return correct result")
	}
}

func TestCursorEdgeCases(t *testing.T) {
	t.Run("move up from y=0 stays at y=0", func(t *testing.T) {
		cursor := Cursor{x: 10, y: 0}
		result := cursor.MoveUp()
		if result.y != 0 {
			t.Errorf("expected y=0, got y=%d", result.y)
		}
	})

	t.Run("move left from x=0 stays at x=0", func(t *testing.T) {
		cursor := Cursor{x: 0, y: 5}
		result := cursor.MoveLeft()
		if result.x != 0 {
			t.Errorf("expected x=0, got x=%d", result.x)
		}
	})

	t.Run("move to position in empty content", func(t *testing.T) {
		cursor := Cursor{x: 10, y: 5}
		result := cursor.MoveToPosition(0, "")
		if result.x != 0 || result.y != 0 {
			t.Errorf("expected (0, 0), got (%d, %d)", result.x, result.y)
		}
	})

	t.Run("get absolute position in empty content", func(t *testing.T) {
		cursor := Cursor{x: 0, y: 0}
		pos := cursor.GetAbsolutePosition("")
		if pos != 0 {
			t.Errorf("expected position 0, got %d", pos)
		}
	})

	t.Run("adjust to line length in empty content", func(t *testing.T) {
		cursor := Cursor{x: 10, y: 0}
		result := cursor.AdjustToLineLength("")
		// Empty content has one empty line with length 0
		if result.x != 0 || result.y != 0 {
			t.Errorf("expected (0, 0), got (%d, %d)", result.x, result.y)
		}
	})
}
