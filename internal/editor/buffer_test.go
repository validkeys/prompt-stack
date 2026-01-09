package editor_test

import (
	"strings"
	"testing"

	"github.com/kyledavis/prompt-stack/internal/editor"
)

func TestNewBuffer(t *testing.T) {
	buf := editor.NewBuffer()

	if buf.Content() != "" {
		t.Errorf("NewBuffer() content = %q, want empty string", buf.Content())
	}

	x, y := buf.CursorPosition()
	if x != 0 || y != 0 {
		t.Errorf("NewBuffer() cursor position = (%d, %d), want (0, 0)", x, y)
	}
}

func TestNewBufferWithContent(t *testing.T) {
	content := "Hello, World!"
	buf := editor.NewBufferWithContent(content)

	if buf.Content() != content {
		t.Errorf("NewBufferWithContent() content = %q, want %q", buf.Content(), content)
	}

	x, y := buf.CursorPosition()
	if x != 13 || y != 0 {
		t.Errorf("NewBufferWithContent() cursor position = (%d, %d), want (13, 0)", x, y)
	}
}

func TestNewBufferWithContentMultiline(t *testing.T) {
	content := "Line 1\nLine 2\nLine 3"
	buf := editor.NewBufferWithContent(content)

	if buf.Content() != content {
		t.Errorf("NewBufferWithContent() content = %q, want %q", buf.Content(), content)
	}

	x, y := buf.CursorPosition()
	if x != 6 || y != 2 {
		t.Errorf("NewBufferWithContent() cursor position = (%d, %d), want (6, 2)", x, y)
	}
}

func TestBufferInsert(t *testing.T) {
	tests := []struct {
		name        string
		initial     string
		cursorX     int
		cursorY     int
		insert      rune
		wantContent string
		wantCursorX int
		wantCursorY int
		wantErr     bool
	}{
		{
			name:        "insert at beginning of empty buffer",
			initial:     "",
			cursorX:     0,
			cursorY:     0,
			insert:      'a',
			wantContent: "a",
			wantCursorX: 1,
			wantCursorY: 0,
			wantErr:     false,
		},
		{
			name:        "insert at beginning of non-empty buffer",
			initial:     "hello",
			cursorX:     0,
			cursorY:     0,
			insert:      'H',
			wantContent: "Hhello",
			wantCursorX: 1,
			wantCursorY: 0,
			wantErr:     false,
		},
		{
			name:        "insert in middle",
			initial:     "helo",
			cursorX:     2,
			cursorY:     0,
			insert:      'l',
			wantContent: "hello",
			wantCursorX: 3,
			wantCursorY: 0,
			wantErr:     false,
		},
		{
			name:        "insert at end",
			initial:     "hello",
			cursorX:     5,
			cursorY:     0,
			insert:      '!',
			wantContent: "hello!",
			wantCursorX: 6,
			wantCursorY: 0,
			wantErr:     false,
		},
		{
			name:        "insert newline",
			initial:     "hello",
			cursorX:     5,
			cursorY:     0,
			insert:      '\n',
			wantContent: "hello\n",
			wantCursorX: 0,
			wantCursorY: 1,
			wantErr:     false,
		},
		{
			name:        "insert unicode emoji",
			initial:     "Hello ",
			cursorX:     6,
			cursorY:     0,
			insert:      '😀',
			wantContent: "Hello 😀",
			wantCursorX: 7,
			wantCursorY: 0,
			wantErr:     false,
		},
		{
			name:        "insert multi-byte rune",
			initial:     "café",
			cursorX:     4,
			cursorY:     0,
			insert:      '🌍',
			wantContent: "café🌍",
			wantCursorX: 5,
			wantCursorY: 0,
			wantErr:     false,
		},
		{
			name:        "insert on second line",
			initial:     "line1\nline",
			cursorX:     4,
			cursorY:     1,
			insert:      '2',
			wantContent: "line1\nline2",
			wantCursorX: 5,
			wantCursorY: 1,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBufferWithContent(tt.initial)
			err := buf.SetCursorPosition(tt.cursorX, tt.cursorY)
			if err != nil {
				t.Fatalf("SetCursorPosition failed: %v", err)
			}

			err = buf.Insert(tt.insert)

			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if buf.Content() != tt.wantContent {
				t.Errorf("Insert() content = %q, want %q", buf.Content(), tt.wantContent)
			}

			x, y := buf.CursorPosition()
			if x != tt.wantCursorX || y != tt.wantCursorY {
				t.Errorf("Insert() cursor position = (%d, %d), want (%d, %d)", x, y, tt.wantCursorX, tt.wantCursorY)
			}
		})
	}
}

func TestBufferDelete(t *testing.T) {
	tests := []struct {
		name        string
		initial     string
		cursorX     int
		cursorY     int
		wantContent string
		wantCursorX int
		wantCursorY int
		wantErr     bool
	}{
		{
			name:        "delete from end",
			initial:     "hello",
			cursorX:     5,
			cursorY:     0,
			wantContent: "hell",
			wantCursorX: 4,
			wantCursorY: 0,
			wantErr:     false,
		},
		{
			name:        "delete from middle",
			initial:     "hello",
			cursorX:     2,
			cursorY:     0,
			wantContent: "hllo",
			wantCursorX: 1,
			wantCursorY: 0,
			wantErr:     false,
		},
		{
			name:        "delete from beginning - should error",
			initial:     "hello",
			cursorX:     0,
			cursorY:     0,
			wantContent: "hello",
			wantCursorX: 0,
			wantCursorY: 0,
			wantErr:     true,
		},
		{
			name:        "delete unicode character",
			initial:     "café🌍",
			cursorX:     5,
			cursorY:     0,
			wantContent: "café",
			wantCursorX: 4,
			wantCursorY: 0,
			wantErr:     false,
		},
		{
			name:        "delete from second line",
			initial:     "line1\nline2",
			cursorX:     5,
			cursorY:     1,
			wantContent: "line1\nline",
			wantCursorX: 4,
			wantCursorY: 1,
			wantErr:     false,
		},
		{
			name:        "delete newline character",
			initial:     "line1\nline2",
			cursorX:     0,
			cursorY:     1,
			wantContent: "line1line2",
			wantCursorX: 5,
			wantCursorY: 0,
			wantErr:     false,
		},
		{
			name:        "delete single character buffer",
			initial:     "a",
			cursorX:     1,
			cursorY:     0,
			wantContent: "",
			wantCursorX: 0,
			wantCursorY: 0,
			wantErr:     false,
		},
		{
			name:        "delete emoji",
			initial:     "Hello 😀",
			cursorX:     7,
			cursorY:     0,
			wantContent: "Hello ",
			wantCursorX: 6,
			wantCursorY: 0,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBufferWithContent(tt.initial)
			err := buf.SetCursorPosition(tt.cursorX, tt.cursorY)
			if err != nil {
				t.Fatalf("SetCursorPosition failed: %v", err)
			}

			err = buf.Delete()

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if buf.Content() != tt.wantContent {
				t.Errorf("Delete() content = %q, want %q", buf.Content(), tt.wantContent)
			}

			x, y := buf.CursorPosition()
			if x != tt.wantCursorX || y != tt.wantCursorY {
				t.Errorf("Delete() cursor position = (%d, %d), want (%d, %d)", x, y, tt.wantCursorX, tt.wantCursorY)
			}
		})
	}
}

func TestBufferContent(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		inserts []rune
		want    string
	}{
		{
			name:    "empty buffer",
			initial: "",
			inserts: nil,
			want:    "",
		},
		{
			name:    "single character",
			initial: "",
			inserts: []rune{'a'},
			want:    "a",
		},
		{
			name:    "multiple characters",
			initial: "",
			inserts: []rune{'h', 'e', 'l', 'l', 'o'},
			want:    "hello",
		},
		{
			name:    "with newlines",
			initial: "",
			inserts: []rune{'a', '\n', 'b'},
			want:    "a\nb",
		},
		{
			name:    "unicode characters",
			initial: "",
			inserts: []rune{'H', 'e', 'l', 'l', 'o', ' ', '😀', '🌍'},
			want:    "Hello 😀🌍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBuffer()

			for _, r := range tt.inserts {
				if err := buf.Insert(r); err != nil {
					t.Fatalf("Insert(%c) failed: %v", r, err)
				}
			}

			if got := buf.Content(); got != tt.want {
				t.Errorf("Content() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBufferCursorPosition(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		moveX   int
		moveY   int
		wantX   int
		wantY   int
	}{
		{
			name:    "initial position",
			initial: "",
			moveX:   0,
			moveY:   0,
			wantX:   0,
			wantY:   0,
		},
		{
			name:    "after insert",
			initial: "hello",
			moveX:   5,
			moveY:   0,
			wantX:   5,
			wantY:   0,
		},
		{
			name:    "after newline insert",
			initial: "hello\nworld",
			moveX:   5,
			moveY:   1,
			wantX:   5,
			wantY:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBufferWithContent(tt.initial)
			err := buf.SetCursorPosition(tt.moveX, tt.moveY)
			if err != nil {
				t.Fatalf("SetCursorPosition failed: %v", err)
			}

			x, y := buf.CursorPosition()
			if x != tt.wantX || y != tt.wantY {
				t.Errorf("CursorPosition() = (%d, %d), want (%d, %d)", x, y, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestBufferSetContent(t *testing.T) {
	tests := []struct {
		name        string
		initial     string
		newContent  string
		wantContent string
	}{
		{
			name:        "set empty content",
			initial:     "hello",
			newContent:  "",
			wantContent: "",
		},
		{
			name:        "set new content",
			initial:     "",
			newContent:  "hello world",
			wantContent: "hello world",
		},
		{
			name:        "set multiline content",
			initial:     "old",
			newContent:  "line1\nline2\nline3",
			wantContent: "line1\nline2\nline3",
		},
		{
			name:        "set unicode content",
			initial:     "",
			newContent:  "Hello 😀🌍",
			wantContent: "Hello 😀🌍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBufferWithContent(tt.initial)
			buf.SetContent(tt.newContent)

			if buf.Content() != tt.wantContent {
				t.Errorf("SetContent() result = %q, want %q", buf.Content(), tt.wantContent)
			}
		})
	}
}

func TestBufferSetCursorPosition(t *testing.T) {
	tests := []struct {
		name    string
		content string
		x       int
		y       int
		wantX   int
		wantY   int
		wantErr bool
	}{
		{
			name:    "valid position",
			content: "hello",
			x:       2,
			y:       0,
			wantX:   2,
			wantY:   0,
			wantErr: false,
		},
		{
			name:    "position at end of line",
			content: "hello",
			x:       5,
			y:       0,
			wantX:   5,
			wantY:   0,
			wantErr: false,
		},
		{
			name:    "position on second line",
			content: "line1\nline2",
			x:       3,
			y:       1,
			wantX:   3,
			wantY:   1,
			wantErr: false,
		},
		{
			name:    "invalid line too high",
			content: "line1\nline2",
			x:       0,
			y:       5,
			wantX:   0,
			wantY:   0,
			wantErr: true,
		},
		{
			name:    "invalid line negative",
			content: "line1",
			x:       0,
			y:       -1,
			wantX:   0,
			wantY:   0,
			wantErr: true,
		},
		{
			name:    "invalid column too high",
			content: "hello",
			x:       10,
			y:       0,
			wantX:   0,
			wantY:   0,
			wantErr: true,
		},
		{
			name:    "invalid column negative",
			content: "hello",
			x:       -1,
			y:       0,
			wantX:   0,
			wantY:   0,
			wantErr: true,
		},
		{
			name:    "empty buffer position 0,0",
			content: "",
			x:       0,
			y:       0,
			wantX:   0,
			wantY:   0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBufferWithContent(tt.content)
			err := buf.SetCursorPosition(tt.x, tt.y)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetCursorPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				x, y := buf.CursorPosition()
				if x != tt.wantX || y != tt.wantY {
					t.Errorf("SetCursorPosition() result = (%d, %d), want (%d, %d)", x, y, tt.wantX, tt.wantY)
				}
			}
		})
	}
}

func TestBufferEdgeCases(t *testing.T) {
	t.Run("empty buffer operations", func(t *testing.T) {
		buf := editor.NewBuffer()

		if buf.Content() != "" {
			t.Errorf("empty buffer content = %q, want empty", buf.Content())
		}

		x, y := buf.CursorPosition()
		if x != 0 || y != 0 {
			t.Errorf("empty buffer cursor = (%d, %d), want (0, 0)", x, y)
		}
	})

	t.Run("single character buffer", func(t *testing.T) {
		buf := editor.NewBufferWithContent("a")

		if buf.Content() != "a" {
			t.Errorf("single char buffer content = %q, want 'a'", buf.Content())
		}

		err := buf.Delete()
		if err != nil {
			t.Errorf("Delete on single char buffer failed: %v", err)
		}

		if buf.Content() != "" {
			t.Errorf("after delete content = %q, want empty", buf.Content())
		}
	})

	t.Run("very long line", func(t *testing.T) {
		longLine := strings.Repeat("a", 1000)
		buf := editor.NewBufferWithContent(longLine)

		if buf.Content() != longLine {
			t.Errorf("long line content length = %d, want %d", len(buf.Content()), len(longLine))
		}

		err := buf.Insert('b')
		if err != nil {
			t.Errorf("Insert on long line failed: %v", err)
		}

		if buf.Content() != longLine+"b" {
			t.Errorf("after insert content length = %d, want %d", len(buf.Content()), len(longLine)+1)
		}
	})

	t.Run("rapid insert operations", func(t *testing.T) {
		buf := editor.NewBuffer()

		for i := 0; i < 100; i++ {
			err := buf.Insert('a')
			if err != nil {
				t.Errorf("Insert iteration %d failed: %v", i, err)
			}
		}

		if buf.Content() != strings.Repeat("a", 100) {
			t.Errorf("rapid insert result = %q, want %q", buf.Content(), strings.Repeat("a", 100))
		}
	})

	t.Run("rapid delete operations", func(t *testing.T) {
		buf := editor.NewBufferWithContent(strings.Repeat("a", 100))

		for i := 0; i < 100; i++ {
			err := buf.Delete()
			if err != nil {
				t.Logf("Delete iteration %d failed (expected at start): %v", i, err)
				break
			}
		}

		if buf.Content() != "" {
			t.Errorf("rapid delete result = %q, want empty", buf.Content())
		}
	})

	t.Run("multibyte unicode handling", func(t *testing.T) {
		buf := editor.NewBuffer()

		unicodes := []rune{'世', '界', '🌍', '😀', '🎉'}
		for _, r := range unicodes {
			if err := buf.Insert(r); err != nil {
				t.Errorf("Insert unicode %c failed: %v", r, err)
			}
		}

		want := "世界🌍😀🎉"
		if buf.Content() != want {
			t.Errorf("unicode content = %q, want %q", buf.Content(), want)
		}
	})
}

func TestBufferMultilineContent(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		inserts   []rune
		cursorX   int
		cursorY   int
		want      string
		wantLines int
	}{
		{
			name:      "insert newlines",
			content:   "",
			inserts:   []rune{'a', '\n', 'b', '\n', 'c'},
			want:      "a\nb\nc",
			wantLines: 3,
		},
		{
			name:      "multiline with text",
			content:   "line1\nline2",
			inserts:   []rune{'!'},
			cursorX:   5,
			cursorY:   0,
			want:      "line1!\nline2",
			wantLines: 2,
		},
		{
			name:      "empty lines",
			content:   "",
			inserts:   []rune{'\n', '\n', '\n'},
			want:      "\n\n\n",
			wantLines: 4,
		},
		{
			name:      "insert in middle of multiline",
			content:   "first\nsecond\nthird",
			inserts:   []rune{'-'},
			cursorX:   6,
			cursorY:   1,
			want:      "first\nsecond-\nthird",
			wantLines: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBufferWithContent(tt.content)

			if tt.cursorX > 0 || tt.cursorY > 0 {
				err := buf.SetCursorPosition(tt.cursorX, tt.cursorY)
				if err != nil {
					t.Fatalf("SetCursorPosition failed: %v", err)
				}
			}

			for _, r := range tt.inserts {
				if err := buf.Insert(r); err != nil {
					t.Fatalf("Insert(%c) failed: %v", r, err)
				}
			}

			if got := buf.Content(); got != tt.want {
				t.Errorf("content = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBufferMoveUp(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		cursorX     int
		cursorY     int
		wantCursorX int
		wantCursorY int
	}{
		{
			name:        "move up from line 1 to line 0",
			content:     "line1\nline2",
			cursorX:     2,
			cursorY:     1,
			wantCursorX: 2,
			wantCursorY: 0,
		},
		{
			name:        "move up preserves column when possible",
			content:     "short\nlonger line",
			cursorX:     5,
			cursorY:     1,
			wantCursorX: 5,
			wantCursorY: 0,
		},
		{
			name:        "move up adjusts to shorter line",
			content:     "short\nlonger line",
			cursorX:     10,
			cursorY:     1,
			wantCursorX: 5,
			wantCursorY: 0,
		},
		{
			name:        "move up at first line stays at line 0",
			content:     "line1\nline2",
			cursorX:     2,
			cursorY:     0,
			wantCursorX: 2,
			wantCursorY: 0,
		},
		{
			name:        "move up on single line stays",
			content:     "single line",
			cursorX:     5,
			cursorY:     0,
			wantCursorX: 5,
			wantCursorY: 0,
		},
		{
			name:        "move up between multiple lines",
			content:     "line1\nline2\nline3",
			cursorX:     3,
			cursorY:     2,
			wantCursorX: 3,
			wantCursorY: 1,
		},
		{
			name:        "move up on empty line",
			content:     "line1\n\nline3",
			cursorX:     0,
			cursorY:     2,
			wantCursorX: 0,
			wantCursorY: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBufferWithContent(tt.content)
			err := buf.SetCursorPosition(tt.cursorX, tt.cursorY)
			if err != nil {
				t.Fatalf("SetCursorPosition failed: %v", err)
			}

			buf.MoveUp()

			x, y := buf.CursorPosition()
			if x != tt.wantCursorX || y != tt.wantCursorY {
				t.Errorf("MoveUp() cursor = (%d, %d), want (%d, %d)", x, y, tt.wantCursorX, tt.wantCursorY)
			}
		})
	}
}

func TestBufferMoveDown(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		cursorX     int
		cursorY     int
		wantCursorX int
		wantCursorY int
	}{
		{
			name:        "move down from line 0 to line 1",
			content:     "line1\nline2",
			cursorX:     2,
			cursorY:     0,
			wantCursorX: 2,
			wantCursorY: 1,
		},
		{
			name:        "move down preserves column when possible",
			content:     "longer line\nshort",
			cursorX:     5,
			cursorY:     0,
			wantCursorX: 5,
			wantCursorY: 1,
		},
		{
			name:        "move down adjusts to shorter line",
			content:     "longer line\nshort",
			cursorX:     10,
			cursorY:     0,
			wantCursorX: 5,
			wantCursorY: 1,
		},
		{
			name:        "move down at last line stays",
			content:     "line1\nline2",
			cursorX:     2,
			cursorY:     1,
			wantCursorX: 2,
			wantCursorY: 1,
		},
		{
			name:        "move down on single line stays",
			content:     "single line",
			cursorX:     5,
			cursorY:     0,
			wantCursorX: 5,
			wantCursorY: 0,
		},
		{
			name:        "move down between multiple lines",
			content:     "line1\nline2\nline3",
			cursorX:     3,
			cursorY:     0,
			wantCursorX: 3,
			wantCursorY: 1,
		},
		{
			name:        "move down on empty line",
			content:     "line1\n\nline3",
			cursorX:     0,
			cursorY:     0,
			wantCursorX: 0,
			wantCursorY: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBufferWithContent(tt.content)
			err := buf.SetCursorPosition(tt.cursorX, tt.cursorY)
			if err != nil {
				t.Fatalf("SetCursorPosition failed: %v", err)
			}

			buf.MoveDown()

			x, y := buf.CursorPosition()
			if x != tt.wantCursorX || y != tt.wantCursorY {
				t.Errorf("MoveDown() cursor = (%d, %d), want (%d, %d)", x, y, tt.wantCursorX, tt.wantCursorY)
			}
		})
	}
}

func TestBufferMoveLeft(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		cursorX     int
		cursorY     int
		wantCursorX int
		wantCursorY int
	}{
		{
			name:        "move left within line",
			content:     "hello",
			cursorX:     3,
			cursorY:     0,
			wantCursorX: 2,
			wantCursorY: 0,
		},
		{
			name:        "move left at column 0 wraps to previous line end",
			content:     "line1\nline2",
			cursorX:     0,
			cursorY:     1,
			wantCursorX: 5,
			wantCursorY: 0,
		},
		{
			name:        "move left at start of buffer stays",
			content:     "hello",
			cursorX:     0,
			cursorY:     0,
			wantCursorX: 0,
			wantCursorY: 0,
		},
		{
			name:        "move left on empty line wraps to previous line end",
			content:     "line1\n\nline3",
			cursorX:     0,
			cursorY:     1,
			wantCursorX: 5,
			wantCursorY: 0,
		},
		{
			name:        "move left handles unicode correctly",
			content:     "Hello 😀",
			cursorX:     7,
			cursorY:     0,
			wantCursorX: 6,
			wantCursorY: 0,
		},
		{
			name:        "move left at end of line",
			content:     "hello",
			cursorX:     5,
			cursorY:     0,
			wantCursorX: 4,
			wantCursorY: 0,
		},
		{
			name:        "move left on single character line",
			content:     "a",
			cursorX:     1,
			cursorY:     0,
			wantCursorX: 0,
			wantCursorY: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBufferWithContent(tt.content)
			err := buf.SetCursorPosition(tt.cursorX, tt.cursorY)
			if err != nil {
				t.Fatalf("SetCursorPosition failed: %v", err)
			}

			buf.MoveLeft()

			x, y := buf.CursorPosition()
			if x != tt.wantCursorX || y != tt.wantCursorY {
				t.Errorf("MoveLeft() cursor = (%d, %d), want (%d, %d)", x, y, tt.wantCursorX, tt.wantCursorY)
			}
		})
	}
}

func TestBufferMoveRight(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		cursorX     int
		cursorY     int
		wantCursorX int
		wantCursorY int
	}{
		{
			name:        "move right within line",
			content:     "hello",
			cursorX:     2,
			cursorY:     0,
			wantCursorX: 3,
			wantCursorY: 0,
		},
		{
			name:        "move right at line end wraps to next line start",
			content:     "line1\nline2",
			cursorX:     5,
			cursorY:     0,
			wantCursorX: 0,
			wantCursorY: 1,
		},
		{
			name:        "move right at end of buffer stays",
			content:     "hello",
			cursorX:     5,
			cursorY:     0,
			wantCursorX: 5,
			wantCursorY: 0,
		},
		{
			name:        "move right on empty line wraps to next line",
			content:     "line1\n\nline3",
			cursorX:     0,
			cursorY:     0,
			wantCursorX: 1,
			wantCursorY: 0,
		},
		{
			name:        "move right handles unicode correctly",
			content:     "Hello 😀",
			cursorX:     6,
			cursorY:     0,
			wantCursorX: 7,
			wantCursorY: 0,
		},
		{
			name:        "move right at beginning of line",
			content:     "hello",
			cursorX:     0,
			cursorY:     0,
			wantCursorX: 1,
			wantCursorY: 0,
		},
		{
			name:        "move right wraps through empty line",
			content:     "line1\n\nline3",
			cursorX:     5,
			cursorY:     0,
			wantCursorX: 0,
			wantCursorY: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBufferWithContent(tt.content)
			err := buf.SetCursorPosition(tt.cursorX, tt.cursorY)
			if err != nil {
				t.Fatalf("SetCursorPosition failed: %v", err)
			}

			buf.MoveRight()

			x, y := buf.CursorPosition()
			if x != tt.wantCursorX || y != tt.wantCursorY {
				t.Errorf("MoveRight() cursor = (%d, %d), want (%d, %d)", x, y, tt.wantCursorX, tt.wantCursorY)
			}
		})
	}
}

func TestBufferMoveToLineStart(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		cursorX     int
		cursorY     int
		wantCursorX int
		wantCursorY int
	}{
		{
			name:        "move to start from middle of line",
			content:     "hello",
			cursorX:     3,
			cursorY:     0,
			wantCursorX: 0,
			wantCursorY: 0,
		},
		{
			name:        "move to start from end of line",
			content:     "hello",
			cursorX:     5,
			cursorY:     0,
			wantCursorX: 0,
			wantCursorY: 0,
		},
		{
			name:        "move to start on multiline document",
			content:     "line1\nline2\nline3",
			cursorX:     5,
			cursorY:     2,
			wantCursorX: 0,
			wantCursorY: 2,
		},
		{
			name:        "move to start on empty line",
			content:     "line1\n\nline3",
			cursorX:     0,
			cursorY:     1,
			wantCursorX: 0,
			wantCursorY: 1,
		},
		{
			name:        "move to start when already at start",
			content:     "hello",
			cursorX:     0,
			cursorY:     0,
			wantCursorX: 0,
			wantCursorY: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBufferWithContent(tt.content)
			err := buf.SetCursorPosition(tt.cursorX, tt.cursorY)
			if err != nil {
				t.Fatalf("SetCursorPosition failed: %v", err)
			}

			buf.MoveToLineStart()

			x, y := buf.CursorPosition()
			if x != tt.wantCursorX || y != tt.wantCursorY {
				t.Errorf("MoveToLineStart() cursor = (%d, %d), want (%d, %d)", x, y, tt.wantCursorX, tt.wantCursorY)
			}
		})
	}
}

func TestBufferMoveToLineEnd(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		cursorX     int
		cursorY     int
		wantCursorX int
		wantCursorY int
	}{
		{
			name:        "move to end from start of line",
			content:     "hello",
			cursorX:     0,
			cursorY:     0,
			wantCursorX: 5,
			wantCursorY: 0,
		},
		{
			name:        "move to end from middle of line",
			content:     "hello",
			cursorX:     2,
			cursorY:     0,
			wantCursorX: 5,
			wantCursorY: 0,
		},
		{
			name:        "move to end on multiline document",
			content:     "line1\nline2\nline3",
			cursorX:     0,
			cursorY:     1,
			wantCursorX: 5,
			wantCursorY: 1,
		},
		{
			name:        "move to end on empty line",
			content:     "line1\n\nline3",
			cursorX:     0,
			cursorY:     1,
			wantCursorX: 0,
			wantCursorY: 1,
		},
		{
			name:        "move to end when already at end",
			content:     "hello",
			cursorX:     5,
			cursorY:     0,
			wantCursorX: 5,
			wantCursorY: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBufferWithContent(tt.content)
			err := buf.SetCursorPosition(tt.cursorX, tt.cursorY)
			if err != nil {
				t.Fatalf("SetCursorPosition failed: %v", err)
			}

			buf.MoveToLineEnd()

			x, y := buf.CursorPosition()
			if x != tt.wantCursorX || y != tt.wantCursorY {
				t.Errorf("MoveToLineEnd() cursor = (%d, %d), want (%d, %d)", x, y, tt.wantCursorX, tt.wantCursorY)
			}
		})
	}
}

func TestBufferMovementEdgeCases(t *testing.T) {
	t.Run("empty buffer movement", func(t *testing.T) {
		buf := editor.NewBuffer()

		buf.MoveUp()
		x, y := buf.CursorPosition()
		if x != 0 || y != 0 {
			t.Errorf("MoveUp on empty buffer: got (%d, %d), want (0, 0)", x, y)
		}

		buf.MoveDown()
		x, y = buf.CursorPosition()
		if x != 0 || y != 0 {
			t.Errorf("MoveDown on empty buffer: got (%d, %d), want (0, 0)", x, y)
		}

		buf.MoveLeft()
		x, y = buf.CursorPosition()
		if x != 0 || y != 0 {
			t.Errorf("MoveLeft on empty buffer: got (%d, %d), want (0, 0)", x, y)
		}

		buf.MoveRight()
		x, y = buf.CursorPosition()
		if x != 0 || y != 0 {
			t.Errorf("MoveRight on empty buffer: got (%d, %d), want (0, 0)", x, y)
		}

		buf.MoveToLineStart()
		x, y = buf.CursorPosition()
		if x != 0 || y != 0 {
			t.Errorf("MoveToLineStart on empty buffer: got (%d, %d), want (0, 0)", x, y)
		}

		buf.MoveToLineEnd()
		x, y = buf.CursorPosition()
		if x != 0 || y != 0 {
			t.Errorf("MoveToLineEnd on empty buffer: got (%d, %d), want (0, 0)", x, y)
		}
	})

	t.Run("single character buffer", func(t *testing.T) {
		buf := editor.NewBufferWithContent("a")

		buf.MoveRight()
		x, y := buf.CursorPosition()
		if x != 1 || y != 0 {
			t.Errorf("MoveRight on single char: got (%d, %d), want (1, 0)", x, y)
		}

		buf.MoveLeft()
		x, y = buf.CursorPosition()
		if x != 0 || y != 0 {
			t.Errorf("MoveLeft on single char: got (%d, %d), want (0, 0)", x, y)
		}
	})

	t.Run("very long line movement", func(t *testing.T) {
		longLine := "a" + strings.Repeat("b", 1000)
		buf := editor.NewBufferWithContent(longLine)

		buf.MoveToLineEnd()
		x, y := buf.CursorPosition()
		if x != 1001 || y != 0 {
			t.Errorf("MoveToLineEnd on long line: got (%d, %d), want (1001, 0)", x, y)
		}

		for i := 0; i < 100; i++ {
			buf.MoveLeft()
		}
		x, y = buf.CursorPosition()
		if x != 901 || y != 0 {
			t.Errorf("MoveLeft 100 times: got (%d, %d), want (901, 0)", x, y)
		}
	})

	t.Run("rapid cursor movements", func(t *testing.T) {
		buf := editor.NewBufferWithContent("line1\nline2\nline3")

		for i := 0; i < 100; i++ {
			buf.MoveUp()
			buf.MoveDown()
			buf.MoveLeft()
			buf.MoveRight()
		}

		x, y := buf.CursorPosition()
		if x != 5 || y != 2 {
			t.Errorf("After 100 rapid movements: got (%d, %d), want (5, 2)", x, y)
		}
	})

	t.Run("column preservation across lines", func(t *testing.T) {
		buf := editor.NewBufferWithContent("short\nmuch longer line here\nmedium")

		buf.SetCursorPosition(15, 1)
		buf.MoveUp()
		x, y := buf.CursorPosition()
		if x != 5 || y != 0 {
			t.Errorf("MoveUp to shorter line: got (%d, %d), want (5, 0)", x, y)
		}

		buf.MoveDown()
		x, y = buf.CursorPosition()
		if x != 5 || y != 1 {
			t.Errorf("MoveDown back to longer line: got (%d, %d), want (5, 1)", x, y)
		}
	})

	t.Run("unicode cursor movement", func(t *testing.T) {
		buf := editor.NewBufferWithContent("Hello 世界 😀 World")

		buf.MoveToLineEnd()
		x, y := buf.CursorPosition()
		if x != 16 || y != 0 {
			t.Errorf("MoveToLineEnd with unicode: got (%d, %d), want (16, 0)", x, y)
		}

		for i := 0; i < 6; i++ {
			buf.MoveLeft()
		}
		x, y = buf.CursorPosition()
		if x != 10 || y != 0 {
			t.Errorf("MoveLeft 6 times: got (%d, %d), want (10, 0)", x, y)
		}
	})

	t.Run("home and end movements", func(t *testing.T) {
		buf := editor.NewBufferWithContent("line1\nline2\nline3")

		buf.SetCursorPosition(2, 1)
		buf.MoveToLineStart()
		x, y := buf.CursorPosition()
		if x != 0 || y != 1 {
			t.Errorf("MoveToLineStart: got (%d, %d), want (0, 1)", x, y)
		}

		buf.MoveToLineEnd()
		x, y = buf.CursorPosition()
		if x != 5 || y != 1 {
			t.Errorf("MoveToLineEnd: got (%d, %d), want (5, 1)", x, y)
		}
	})
}

func TestBufferMovementBoundaryHandling(t *testing.T) {
	t.Run("boundary at start of buffer", func(t *testing.T) {
		buf := editor.NewBufferWithContent("hello")
		buf.SetCursorPosition(0, 0)

		buf.MoveLeft()
		x, y := buf.CursorPosition()
		if x != 0 || y != 0 {
			t.Errorf("MoveLeft at start: got (%d, %d), want (0, 0)", x, y)
		}

		buf.MoveUp()
		x, y = buf.CursorPosition()
		if x != 0 || y != 0 {
			t.Errorf("MoveUp at start: got (%d, %d), want (0, 0)", x, y)
		}
	})

	t.Run("boundary at end of buffer", func(t *testing.T) {
		buf := editor.NewBufferWithContent("hello")

		buf.MoveRight()
		x, y := buf.CursorPosition()
		if x != 5 || y != 0 {
			t.Errorf("MoveRight at end: got (%d, %d), want (5, 0)", x, y)
		}

		buf.MoveDown()
		x, y = buf.CursorPosition()
		if x != 5 || y != 0 {
			t.Errorf("MoveDown at end: got (%d, %d), want (5, 0)", x, y)
		}
	})

	t.Run("boundary at line transitions", func(t *testing.T) {
		buf := editor.NewBufferWithContent("line1\nline2")

		buf.SetCursorPosition(0, 1)
		buf.MoveLeft()
		x, y := buf.CursorPosition()
		if x != 5 || y != 0 {
			t.Errorf("MoveLeft at line transition: got (%d, %d), want (5, 0)", x, y)
		}

		buf.SetCursorPosition(5, 0)
		buf.MoveRight()
		x, y = buf.CursorPosition()
		if x != 0 || y != 1 {
			t.Errorf("MoveRight at line transition: got (%d, %d), want (0, 1)", x, y)
		}
	})
}

func TestBufferCharCount(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantChar int
		wantLine int
	}{
		{
			name:     "empty buffer",
			content:  "",
			wantChar: 0,
			wantLine: 1,
		},
		{
			name:     "single character",
			content:  "a",
			wantChar: 1,
			wantLine: 1,
		},
		{
			name:     "multiple characters single line",
			content:  "Hello, World!",
			wantChar: 13,
			wantLine: 1,
		},
		{
			name:     "multiple lines",
			content:  "line1\nline2\nline3",
			wantChar: 17,
			wantLine: 3,
		},
		{
			name:     "only newlines",
			content:  "\n\n\n",
			wantChar: 3,
			wantLine: 4,
		},
		{
			name:     "empty lines",
			content:  "line1\n\nline3",
			wantChar: 12,
			wantLine: 3,
		},
		{
			name:     "unicode characters",
			content:  "Hello 世界 😀 World",
			wantChar: 16,
			wantLine: 1,
		},
		{
			name:     "emoji only",
			content:  "😀🎉🌍",
			wantChar: 3,
			wantLine: 1,
		},
		{
			name:     "mixed unicode and ascii",
			content:  "a世b😀c🌍d",
			wantChar: 7,
			wantLine: 1,
		},
		{
			name:     "large document",
			content:  strings.Repeat("a", 10000) + "\n" + strings.Repeat("b", 10000),
			wantChar: 20001,
			wantLine: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := editor.NewBufferWithContent(tt.content)

			gotChar := buf.CharCount()
			if gotChar != tt.wantChar {
				t.Errorf("CharCount() = %d, want %d", gotChar, tt.wantChar)
			}

			gotLine := buf.LineCount()
			if gotLine != tt.wantLine {
				t.Errorf("LineCount() = %d, want %d", gotLine, tt.wantLine)
			}
		})
	}
}

func TestBufferCharCountUpdates(t *testing.T) {
	t.Run("count updates after insert", func(t *testing.T) {
		buf := editor.NewBuffer()

		if buf.CharCount() != 0 {
			t.Errorf("initial CharCount = %d, want 0", buf.CharCount())
		}
		if buf.LineCount() != 1 {
			t.Errorf("initial LineCount = %d, want 1", buf.LineCount())
		}

		buf.Insert('a')
		if buf.CharCount() != 1 {
			t.Errorf("after insert 'a': CharCount = %d, want 1", buf.CharCount())
		}
		if buf.LineCount() != 1 {
			t.Errorf("after insert 'a': LineCount = %d, want 1", buf.LineCount())
		}

		buf.Insert('\n')
		if buf.CharCount() != 2 {
			t.Errorf("after insert newline: CharCount = %d, want 2", buf.CharCount())
		}
		if buf.LineCount() != 2 {
			t.Errorf("after insert newline: LineCount = %d, want 2", buf.LineCount())
		}

		buf.Insert('b')
		if buf.CharCount() != 3 {
			t.Errorf("after insert 'b': CharCount = %d, want 3", buf.CharCount())
		}
		if buf.LineCount() != 2 {
			t.Errorf("after insert 'b': LineCount = %d, want 2", buf.LineCount())
		}
	})

	t.Run("count updates after delete", func(t *testing.T) {
		buf := editor.NewBufferWithContent("abc\nxyz")

		if buf.CharCount() != 7 {
			t.Errorf("initial CharCount = %d, want 7", buf.CharCount())
		}
		if buf.LineCount() != 2 {
			t.Errorf("initial LineCount = %d, want 2", buf.LineCount())
		}

		buf.Delete()
		if buf.CharCount() != 6 {
			t.Errorf("after delete: CharCount = %d, want 6", buf.CharCount())
		}
		if buf.LineCount() != 2 {
			t.Errorf("after delete: LineCount = %d, want 2", buf.LineCount())
		}

		buf.Delete()
		buf.Delete()
		buf.Delete()
		if buf.CharCount() != 3 {
			t.Errorf("after 3 more deletes: CharCount = %d, want 3", buf.CharCount())
		}
		if buf.LineCount() != 1 {
			t.Errorf("after 3 more deletes: LineCount = %d, want 1", buf.LineCount())
		}
	})

	t.Run("count with rapid inserts", func(t *testing.T) {
		buf := editor.NewBuffer()

		for i := 0; i < 100; i++ {
			buf.Insert('a')
		}

		if buf.CharCount() != 100 {
			t.Errorf("after 100 inserts: CharCount = %d, want 100", buf.CharCount())
		}
		if buf.LineCount() != 1 {
			t.Errorf("after 100 inserts: LineCount = %d, want 1", buf.LineCount())
		}
	})
}

func TestBufferLineCountEdgeCases(t *testing.T) {
	t.Run("empty buffer has 1 line", func(t *testing.T) {
		buf := editor.NewBuffer()
		if buf.LineCount() != 1 {
			t.Errorf("empty buffer LineCount = %d, want 1", buf.LineCount())
		}
	})

	t.Run("single newline creates 2 lines", func(t *testing.T) {
		buf := editor.NewBufferWithContent("\n")
		if buf.LineCount() != 2 {
			t.Errorf("'\\n' LineCount = %d, want 2", buf.LineCount())
		}
		if buf.CharCount() != 1 {
			t.Errorf("'\\n' CharCount = %d, want 1", buf.CharCount())
		}
	})

	t.Run("consecutive newlines", func(t *testing.T) {
		buf := editor.NewBufferWithContent("\n\n\n")
		if buf.LineCount() != 4 {
			t.Errorf("'\\n\\n\\n' LineCount = %d, want 4", buf.LineCount())
		}
		if buf.CharCount() != 3 {
			t.Errorf("'\\n\\n\\n' CharCount = %d, want 3", buf.CharCount())
		}
	})

	t.Run("text with trailing newline", func(t *testing.T) {
		buf := editor.NewBufferWithContent("hello\n")
		if buf.LineCount() != 2 {
			t.Errorf("'hello\\n' LineCount = %d, want 2", buf.LineCount())
		}
		if buf.CharCount() != 6 {
			t.Errorf("'hello\\n' CharCount = %d, want 6", buf.CharCount())
		}
	})

	t.Run("text without trailing newline", func(t *testing.T) {
		buf := editor.NewBufferWithContent("hello")
		if buf.LineCount() != 1 {
			t.Errorf("'hello' LineCount = %d, want 1", buf.LineCount())
		}
		if buf.CharCount() != 5 {
			t.Errorf("'hello' CharCount = %d, want 5", buf.CharCount())
		}
	})
}

func TestBufferCharCountUnicode(t *testing.T) {
	t.Run("multibyte unicode characters", func(t *testing.T) {
		tests := []struct {
			content  string
			wantChar int
		}{
			{"世", 1},
			{"世界", 2},
			{"😀", 1},
			{"🌍", 1},
			{"🎉", 1},
			{"😀🎉🌍", 3},
			{"Hello 世界 😀", 10},
			{"a世b😀c🌍d", 7},
		}

		for _, tt := range tests {
			t.Run(tt.content, func(t *testing.T) {
				buf := editor.NewBufferWithContent(tt.content)
				got := buf.CharCount()
				if got != tt.wantChar {
					t.Errorf("CharCount() for %q = %d, want %d", tt.content, got, tt.wantChar)
				}
			})
		}
	})

	t.Run("insert unicode characters", func(t *testing.T) {
		buf := editor.NewBuffer()

		unicodes := []rune{'世', '界', '😀', '🎉', '🌍'}
		for _, r := range unicodes {
			buf.Insert(r)
		}

		if buf.CharCount() != 5 {
			t.Errorf("after inserting 5 unicode chars: CharCount = %d, want 5", buf.CharCount())
		}
		if buf.LineCount() != 1 {
			t.Errorf("after inserting 5 unicode chars: LineCount = %d, want 1", buf.LineCount())
		}
	})
}

func TestBufferCountingPerformance(t *testing.T) {
	t.Run("large document char count", func(t *testing.T) {
		content := strings.Repeat("a", 10000)
		buf := editor.NewBufferWithContent(content)

		got := buf.CharCount()
		if got != 10000 {
			t.Errorf("CharCount() = %d, want 10000", got)
		}
	})

	t.Run("large document line count", func(t *testing.T) {
		lines := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			lines[i] = "line " + string(rune('0'+i%10))
		}
		content := strings.Join(lines, "\n")
		buf := editor.NewBufferWithContent(content)

		got := buf.LineCount()
		if got != 1000 {
			t.Errorf("LineCount() = %d, want 1000", got)
		}
	})
}
