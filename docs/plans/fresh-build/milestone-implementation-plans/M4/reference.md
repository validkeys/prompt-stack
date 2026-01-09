# Milestone 4: Basic Text Editor - Reference Document

## How to Use This Document

**Read this section when:**
- Before implementing any task
- When writing buffer operations (Lines 30-150)
- When implementing cursor movement (Lines 151-250)
- When building Bubble Tea UI (Lines 251-400)
- When writing tests (Lines 401-550)

**Key sections:**
- Lines 30-150: Buffer Architecture - Read before Task 1 (Buffer Package)
- Lines 151-250: Cursor Movement Patterns - Read before Task 2 (Cursor Movement)
- Lines 251-400: Bubble Tea Integration - Read before Task 3 (Workspace TUI)
- Lines 401-550: Testing Patterns - Reference during all tasks
- Lines 551-600: OpenCode Design Integration - Read before Task 6 (Theme Styles)

**Related documents:**
- See [task-list.md](task-list.md) for detailed task breakdown
- Cross-reference with [`project-structure.md`](../../project-structure.md) for file paths
- Cross-reference with [`go-style-guide.md`](../../go-style-guide.md) for coding patterns

---

## Architecture Context

### Domain Overview

Milestone 4 implements the **Editor Domain** and its **UI integration**. The editor domain provides core text editing functionality, while the UI layer presents it using Bubble Tea's TUI framework.

**Key Architectural Principles**:
1. **Separation of Concerns**: Editor logic is independent of UI
2. **Rune-Based Operations**: All text operations use runes for proper unicode support
3. **Immutable Updates**: Bubble Tea models return new instances, not mutating state
4. **Effect-Based Testing**: Test observable outcomes, not implementation details

### Package Structure

```
internal/editor/          # Editor Domain (business logic)
├── buffer.go             # Core text buffer with cursor
├── buffer_test.go        # Buffer unit tests
└── ...                   # Future: undo.go, autosave.go

ui/workspace/             # Workspace UI (Bubble Tea component)
├── model.go              # Model state
├── update.go             # Update logic (event handling)
├── view.go               # View rendering
├── messages.go           # Custom message types
└── model_test.go         # Bubble Tea tests
```

**Dependencies Flow**:
```
ui/workspace → internal/editor → (no dependencies)
```

### Key Types

```go
// internal/editor/buffer.go
type Buffer struct {
    content string  // Full text content
    cursorX int     // Cursor column (0-based)
    cursorY int     // Cursor line (0-based)
}

// ui/workspace/model.go
type Model struct {
    buffer   *editor.Buffer  // Editor buffer
    width    int             // Terminal width
    height   int             // Terminal height
    viewportX int            // Viewport horizontal offset
    viewportY int            // Viewport vertical offset
}
```

---

## Buffer Architecture (Lines 30-150)

### Core Buffer Implementation

**Purpose**: Store text content with cursor position, provide text manipulation methods.

**Key Design Decisions**:
1. Use `string` for content (Go strings are UTF-8 by default)
2. Use `[]rune` for indexing (handles multi-byte characters correctly)
3. Track cursor with (x, y) coordinates, not absolute position
4. Split on `\n` for line-based operations

**Implementation Pattern**:

```go
package editor

import (
    "fmt"
    "strings"
)

// Buffer represents a text buffer with cursor position.
type Buffer struct {
    content string
    cursorX int
    cursorY int
}

// New creates a new Buffer with the given content.
func New(content string) *Buffer {
    return &Buffer{
        content: content,
        cursorX: 0,
        cursorY: 0,
    }
}

// Content returns the current buffer content.
func (b *Buffer) Content() string {
    return b.content
}

// CursorPosition returns the current cursor position (x, y).
func (b *Buffer) CursorPosition() (int, int) {
    return b.cursorX, b.cursorY
}

// SetCursorPosition sets the cursor position to the given coordinates.
// This method is primarily used for testing.
// Returns an error if the position is out of bounds.
func (b *Buffer) SetCursorPosition(x, y int) error {
    lines := strings.Split(b.content, "\n")
    
    if y < 0 || y >= len(lines) {
        return fmt.Errorf("cursor y %d out of bounds [0, %d]", y, len(lines)-1)
    }
    
    line := lines[y]
    runes := []rune(line)
    
    if x < 0 || x > len(runes) {
        return fmt.Errorf("cursor x %d out of bounds [0, %d] for line %d", x, len(runes), y)
    }
    
    b.cursorX = x
    b.cursorY = y
    return nil
}

// Insert inserts a rune at the current cursor position.
// The cursor moves to after the inserted rune.
func (b *Buffer) Insert(r rune) error {
    lines := strings.Split(b.content, "\n")
    
    // Validate cursor position
    if b.cursorY >= len(lines) {
        return fmt.Errorf("cursor y %d exceeds line count %d", b.cursorY, len(lines))
    }
    
    line := lines[b.cursorY]
    runes := []rune(line)
    
    // Validate cursor x position
    if b.cursorX > len(runes) {
        return fmt.Errorf("cursor x %d exceeds line length %d", b.cursorX, len(runes))
    }
    
    // Insert rune at cursor position
    newRunes := make([]rune, len(runes)+1)
    copy(newRunes[:b.cursorX], runes[:b.cursorX])
    newRunes[b.cursorX] = r
    copy(newRunes[b.cursorX+1:], runes[b.cursorX:])
    
    lines[b.cursorY] = string(newRunes)
    b.content = strings.Join(lines, "\n")
    b.cursorX++
    
    return nil
}

// Delete deletes the character at the current cursor position (Delete key).
func (b *Buffer) Delete() error {
    lines := strings.Split(b.content, "\n")
    
    if b.cursorY >= len(lines) {
        return fmt.Errorf("cursor y %d exceeds line count %d", b.cursorY, len(lines))
    }
    
    line := lines[b.cursorY]
    runes := []rune(line)
    
    // Check if at end of line
    if b.cursorX >= len(runes) {
        // At end of line, join with next line
        if b.cursorY < len(lines)-1 {
            lines[b.cursorY] = line + lines[b.cursorY+1]
            lines = append(lines[:b.cursorY+1], lines[b.cursorY+2:]...)
            b.content = strings.Join(lines, "\n")
        }
        return nil
    }
    
    // Delete character at cursor
    newRunes := make([]rune, len(runes)-1)
    copy(newRunes[:b.cursorX], runes[:b.cursorX])
    copy(newRunes[b.cursorX:], runes[b.cursorX+1:])
    
    lines[b.cursorY] = string(newRunes)
    b.content = strings.Join(lines, "\n")
    
    return nil
}

// Backspace deletes the character before the cursor (Backspace key).
func (b *Buffer) Backspace() error {
    if b.cursorX == 0 && b.cursorY == 0 {
        // At start of buffer, nothing to delete
        return nil
    }
    
    if b.cursorX == 0 {
        // At start of line, join with previous line
        lines := strings.Split(b.content, "\n")
        prevLine := lines[b.cursorY-1]
        b.cursorX = len([]rune(prevLine))
        lines[b.cursorY-1] = prevLine + lines[b.cursorY]
        lines = append(lines[:b.cursorY], lines[b.cursorY+1:]...)
        b.content = strings.Join(lines, "\n")
        b.cursorY--
        return nil
    }
    
    // Move cursor left and delete
    b.cursorX--
    return b.Delete()
}

// InsertNewline inserts a newline at the current cursor position (Enter key).
func (b *Buffer) InsertNewline() error {
    lines := strings.Split(b.content, "\n")
    
    if b.cursorY >= len(lines) {
        return fmt.Errorf("cursor y %d exceeds line count %d", b.cursorY, len(lines))
    }
    
    line := lines[b.cursorY]
    runes := []rune(line)
    
    // Split line at cursor
    beforeCursor := string(runes[:b.cursorX])
    afterCursor := string(runes[b.cursorX:])
    
    lines[b.cursorY] = beforeCursor
    newLines := make([]string, len(lines)+1)
    copy(newLines[:b.cursorY+1], lines[:b.cursorY+1])
    newLines[b.cursorY+1] = afterCursor
    copy(newLines[b.cursorY+2:], lines[b.cursorY+1:])
    
    b.content = strings.Join(newLines, "\n")
    b.cursorY++
    b.cursorX = 0
    
    return nil
}

// CharCount returns the total number of characters in the buffer.
func (b *Buffer) CharCount() int {
    return len([]rune(b.content))
}

// LineCount returns the total number of lines in the buffer.
func (b *Buffer) LineCount() int {
    if b.content == "" {
        return 0
    }
    return strings.Count(b.content, "\n") + 1
}
```

**Key Patterns from [`learnings/editor-domain.md`](../../learnings/editor-domain.md)**:
- Use rune-based indexing for unicode correctness
- Validate cursor position before operations
- Return descriptive errors with context
- Keep operations simple and focused

---

## Cursor Movement Patterns (Lines 151-250)

### Cursor Navigation Methods

**Purpose**: Move cursor through buffer with proper boundary handling.

**Implementation Pattern**:

```go
// MoveUp moves the cursor up one line.
func (b *Buffer) MoveUp() {
    if b.cursorY > 0 {
        b.cursorY--
        b.constrainCursorX()
    }
}

// MoveDown moves the cursor down one line.
func (b *Buffer) MoveDown() {
    lines := strings.Split(b.content, "\n")
    if b.cursorY < len(lines)-1 {
        b.cursorY++
        b.constrainCursorX()
    }
}

// MoveLeft moves the cursor left one character.
func (b *Buffer) MoveLeft() {
    if b.cursorX > 0 {
        b.cursorX--
    } else if b.cursorY > 0 {
        // Wrap to end of previous line
        b.cursorY--
        lines := strings.Split(b.content, "\n")
        if b.cursorY < len(lines) {
            b.cursorX = len([]rune(lines[b.cursorY]))
        }
    }
}

// MoveRight moves the cursor right one character.
func (b *Buffer) MoveRight() {
    lines := strings.Split(b.content, "\n")
    if b.cursorY < len(lines) {
        lineLen := len([]rune(lines[b.cursorY]))
        if b.cursorX < lineLen {
            b.cursorX++
        } else if b.cursorY < len(lines)-1 {
            // Wrap to start of next line
            b.cursorY++
            b.cursorX = 0
        }
    }
}

// MoveToLineStart moves the cursor to the start of the current line (Home key).
func (b *Buffer) MoveToLineStart() {
    b.cursorX = 0
}

// MoveToLineEnd moves the cursor to the end of the current line (End key).
func (b *Buffer) MoveToLineEnd() {
    lines := strings.Split(b.content, "\n")
    if b.cursorY < len(lines) {
        b.cursorX = len([]rune(lines[b.cursorY]))
    }
}

// constrainCursorX ensures cursor X is within the current line bounds.
func (b *Buffer) constrainCursorX() {
    lines := strings.Split(b.content, "\n")
    if b.cursorY < len(lines) {
        lineLen := len([]rune(lines[b.cursorY]))
        if b.cursorX > lineLen {
            b.cursorX = lineLen
        }
    }
}
```

**Key Patterns from [`learnings/editor-domain.md`](../../learnings/editor-domain.md)**:
- Handle wrapping at line boundaries (left/right)
- Preserve cursor column when possible (up/down)
- Constrain cursor X when moving to shorter lines
- No-op at document boundaries

---

## Bubble Tea Integration (Lines 251-400)

### Workspace Model Implementation

**Purpose**: Connect editor buffer to Bubble Tea TUI framework.

**Implementation Pattern**:

```go
package workspace

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/kyledavis/prompt-stack/internal/editor"
    "github.com/kyledavis/prompt-stack/ui/theme"
    "github.com/charmbracelet/lipgloss"
    "strings"
)

// Model represents the workspace editor state.
type Model struct {
    buffer    *editor.Buffer
    width     int
    height    int
    viewportX int
    viewportY int
}

// New creates a new workspace model.
func New() Model {
    return Model{
        buffer:    editor.New(""),
        width:     80,
        height:    24,
        viewportX: 0,
        viewportY: 0,
    }
}

// Init initializes the workspace model.
func (m Model) Init() tea.Cmd {
    return nil
}

// Update handles messages and updates model state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyMsg(msg)
        
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil
    }
    
    return m, nil
}

// handleKeyMsg handles keyboard input.
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.Type {
    case tea.KeyRunes:
        // Insert characters
        for _, r := range msg.Runes {
            if err := m.buffer.Insert(r); err != nil {
                // Log error but don't crash
            }
        }
        m.adjustViewport()
        return m, nil
        
    case tea.KeyBackspace:
        if err := m.buffer.Backspace(); err != nil {
            // Log error but don't crash
        }
        m.adjustViewport()
        return m, nil
        
    case tea.KeyDelete:
        if err := m.buffer.Delete(); err != nil {
            // Log error but don't crash
        }
        m.adjustViewport()
        return m, nil
        
    case tea.KeyEnter:
        if err := m.buffer.InsertNewline(); err != nil {
            // Log error but don't crash
        }
        m.adjustViewport()
        return m, nil
        
    case tea.KeyUp:
        m.buffer.MoveUp()
        m.adjustViewport()
        return m, nil
        
    case tea.KeyDown:
        m.buffer.MoveDown()
        m.adjustViewport()
        return m, nil
        
    case tea.KeyLeft:
        m.buffer.MoveLeft()
        m.adjustViewport()
        return m, nil
        
    case tea.KeyRight:
        m.buffer.MoveRight()
        m.adjustViewport()
        return m, nil
        
    case tea.KeyHome:
        m.buffer.MoveToLineStart()
        m.adjustViewport()
        return m, nil
        
    case tea.KeyEnd:
        m.buffer.MoveToLineEnd()
        m.adjustViewport()
        return m, nil
    }
    
    return m, nil
}

// adjustViewport adjusts viewport to keep cursor visible.
func (m *Model) adjustViewport() {
    _, cursorY := m.buffer.CursorPosition()
    availableHeight := m.height - 1 // Reserve 1 line for status bar
    
    // "Middle third" scrolling strategy
    third := availableHeight / 3
    
    if cursorY < m.viewportY+third {
        // Cursor too close to top, scroll up
        m.viewportY = max(0, cursorY-third)
    } else if cursorY > m.viewportY+availableHeight-third {
        // Cursor too close to bottom, scroll down
        m.viewportY = max(0, cursorY-availableHeight+third)
    }
}

// View renders the workspace.
func (m Model) View() string {
    // Render editor content
    editorView := m.renderEditor()
    
    // Render status bar
    statusBar := m.renderStatusBar()
    
    // Combine vertically
    return lipgloss.JoinVertical(lipgloss.Left, editorView, statusBar)
}

// renderEditor renders the editor content with viewport.
func (m Model) renderEditor() string {
    lines := strings.Split(m.buffer.Content(), "\n")
    cursorX, cursorY := m.buffer.CursorPosition()
    availableHeight := m.height - 1
    
    // Calculate visible line range
    startLine := m.viewportY
    endLine := min(len(lines), m.viewportY+availableHeight)
    
    var displayLines []string
    for i := startLine; i < endLine; i++ {
        line := lines[i]
        
        // Highlight cursor if on this line
        if i == cursorY {
            runes := []rune(line)
            if cursorX >= 0 && cursorX <= len(runes) {
                // Insert cursor marker
                before := string(runes[:cursorX])
                cursor := " " // Cursor character
                if cursorX < len(runes) {
                    cursor = string(runes[cursorX])
                }
                after := ""
                if cursorX < len(runes) {
                    after = string(runes[cursorX+1:])
                }
                
                // Apply cursor style
                cursorStyled := theme.CursorStyle().Render(cursor)
                line = before + cursorStyled + after
            }
        }
        
        displayLines = append(displayLines, line)
    }
    
    // Apply editor style
    editorStyle := theme.EditorStyle().
        Width(m.width).
        Height(availableHeight)
    
    return editorStyle.Render(strings.Join(displayLines, "\n"))
}

// renderStatusBar renders the status bar.
func (m Model) renderStatusBar() string {
    charCount := m.buffer.CharCount()
    lineCount := m.buffer.LineCount()
    
    statusText := fmt.Sprintf("%d chars | %d lines", charCount, lineCount)
    
    statusStyle := theme.StatusStyle().
        Width(m.width).
        Height(1)
    
    return statusStyle.Render(statusText)
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

**Key Patterns from [`learnings/ui-domain.md`](../../learnings/ui-domain.md)**:
- Follow Bubble Tea model-view-update pattern
- Return updated model and commands from Update()
- Adjust viewport after cursor movements
- Use "middle third" scrolling strategy

---

## Testing Patterns (Lines 401-550)

### Buffer Unit Tests

**Test Pattern from [`bubble-tea-testing-best-practices.md`](../../bubble-tea-testing-best-practices.md)**:

```go
package editor_test

import (
    "testing"
    "github.com/kyledavis/prompt-stack/internal/editor"
)

func TestBuffer_Insert(t *testing.T) {
    tests := []struct {
        name        string
        initial     string
        cursorX     int
        cursorY     int
        insertRune  rune
        wantContent string
        wantCursorX int
        wantCursorY int
    }{
        {
            name:        "insert at start",
            initial:     "hello",
            cursorX:     0,
            cursorY:     0,
            insertRune:  'x',
            wantContent: "xhello",
            wantCursorX: 1,
            wantCursorY: 0,
        },
        {
            name:        "insert in middle",
            initial:     "hello",
            cursorX:     3,
            cursorY:     0,
            insertRune:  'x',
            wantContent: "helxlo",
            wantCursorX: 4,
            wantCursorY: 0,
        },
        {
            name:        "insert unicode",
            initial:     "hello",
            cursorX:     5,
            cursorY:     0,
            insertRune:  '😀',
            wantContent: "hello😀",
            wantCursorX: 6,
            wantCursorY: 0,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            buf := editor.New(tt.initial)
            // Set cursor position
            buf.SetCursorPosition(tt.cursorX, tt.cursorY)
            
            err := buf.Insert(tt.insertRune)
            if err != nil {
                t.Fatalf("Insert() error = %v", err)
            }
            
            if got := buf.Content(); got != tt.wantContent {
                t.Errorf("Content() = %q, want %q", got, tt.wantContent)
            }
            
            gotX, gotY := buf.CursorPosition()
            if gotX != tt.wantCursorX || gotY != tt.wantCursorY {
                t.Errorf("CursorPosition() = (%d, %d), want (%d, %d)", 
                    gotX, gotY, tt.wantCursorX, tt.wantCursorY)
            }
        })
    }
}

func TestBuffer_MoveCursor(t *testing.T) {
    tests := []struct {
        name        string
        initial     string
        startX      int
        startY      int
        operation   func(*editor.Buffer)
        wantX       int
        wantY       int
    }{
        {
            name:      "move right",
            initial:   "hello\nworld",
            startX:    0,
            startY:    0,
            operation: func(b *editor.Buffer) { b.MoveRight() },
            wantX:     1,
            wantY:     0,
        },
        {
            name:      "move right at line end wraps",
            initial:   "hello\nworld",
            startX:    5,
            startY:    0,
            operation: func(b *editor.Buffer) { b.MoveRight() },
            wantX:     0,
            wantY:     1,
        },
        {
            name:      "move up",
            initial:   "hello\nworld",
            startX:    2,
            startY:    1,
            operation: func(b *editor.Buffer) { b.MoveUp() },
            wantX:     2,
            wantY:     0,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            buf := editor.New(tt.initial)
            buf.SetCursorPosition(tt.startX, tt.startY)
            
            tt.operation(buf)
            
            gotX, gotY := buf.CursorPosition()
            if gotX != tt.wantX || gotY != tt.wantY {
                t.Errorf("CursorPosition() = (%d, %d), want (%d, %d)", 
                    gotX, gotY, tt.wantX, tt.wantY)
            }
        })
    }
}
```

### Bubble Tea Integration Tests

**Test Pattern from [`bubble-tea-testing-best-practices.md`](../../bubble-tea-testing-best-practices.md)**:

```go
package workspace_test

import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/kyledavis/prompt-stack/ui/workspace"
)

func TestWorkspace_TypeText(t *testing.T) {
    model := workspace.New()
    
    // Type "hello"
    for _, r := range "hello" {
        msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
        var cmd tea.Cmd
        model, cmd = model.Update(msg)
        
        if cmd != nil {
            // Handle any commands returned
        }
    }
    
    // Verify content
    if got := model.GetContent(); got != "hello" {
        t.Errorf("Content() = %q, want %q", got, "hello")
    }
    
    // Verify character count
    if got := model.GetCharCount(); got != 5 {
        t.Errorf("CharCount() = %d, want 5", got)
    }
}

func TestWorkspace_CursorMovement(t *testing.T) {
    model := workspace.New()
    
    // Type "hello\nworld"
    for _, r := range "hello" {
        msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
        model, _ = model.Update(msg)
    }
    model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
    for _, r := range "world" {
        msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
        model, _ = model.Update(msg)
    }
    
    // Move cursor up
    model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
    
    // Verify cursor position (should be on first line)
    x, y := model.GetCursorPosition()
    if y != 0 {
        t.Errorf("CursorY = %d, want 0", y)
    }
}
```

**Key Testing Principles**:
- Test effects, not implementation
- Use table-driven tests for multiple cases
- Simulate real user input sequences
- Verify observable state changes
- Handle edge cases (empty, boundaries, unicode)

---

## OpenCode Design Integration (Lines 551-600)

### Theme Style Application

**Purpose**: Apply OpenCode design system to workspace following centralized theme package.

**Style Pattern from [`opencode-design-system.md`](../../opencode-design-system.md)**:

```go
// ui/theme/theme.go additions for M4

// EditorStyle returns the style for the editor content area.
func EditorStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Background(lipgloss.Color(BackgroundPrimary)).
        Foreground(lipgloss.Color(ForegroundPrimary)).
        Padding(0, 1)  // 1-character unit padding
}

// CursorStyle returns the style for cursor highlighting.
func CursorStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Background(lipgloss.Color(CursorBackground)).
        Foreground(lipgloss.Color(CursorForeground))
}
```

**Usage in Workspace**:

```go
// ui/workspace/view.go

import "github.com/kyledavis/prompt-stack/ui/theme"

func (m Model) renderEditor() string {
    // ... content rendering ...
    
    editorStyle := theme.EditorStyle().
        Width(m.width).
        Height(m.height - 1)  // Reserve 1 line for status bar
    
    return editorStyle.Render(content)
}

func (m Model) renderStatusBar() string {
    statusText := fmt.Sprintf("%d chars | %d lines", charCount, lineCount)
    
    statusStyle := theme.StatusStyle().
        Width(m.width).
        Height(1)
    
    return statusStyle.Render(statusText)
}
```

**Key Design Principles**:
- Use theme helpers, never hard-code colors
- Follow 1-character unit spacing system
- Ensure high contrast for cursor visibility
- Maintain visual hierarchy (editor vs status bar)
- Apply OpenCode color palette consistently

---

## Common Pitfalls

### ❌ Don't Use Byte Indexing

```go
// BAD: Byte indexing breaks unicode
content[5] = 'x'

// GOOD: Rune-based indexing
runes := []rune(content)
runes[5] = 'x'
content = string(runes)
```

### ❌ Don't Forget Viewport Adjustment

```go
// BAD: Cursor moves off screen
m.buffer.MoveDown()
return m, nil

// GOOD: Adjust viewport after movement
m.buffer.MoveDown()
m.adjustViewport()
return m, nil
```

### ❌ Don't Hard-Code Colors

```go
// BAD: Hard-coded color
style := lipgloss.NewStyle().Background(lipgloss.Color("#1e1e2e"))

// GOOD: Use theme
style := theme.EditorStyle()
```

### ❌ Don't Ignore Test Coverage

```go
// BAD: Only test happy path
func TestInsert(t *testing.T) {
    buf := editor.New("hello")
    buf.Insert('x')
    // No assertion!
}

// GOOD: Test effects with assertions
func TestInsert(t *testing.T) {
    buf := editor.New("hello")
    buf.Insert('x')
    if got := buf.Content(); got != "xhello" {
        t.Errorf("got %q, want %q", got, "xhello")
    }
}
```

---

## Implementation Checklist

Before starting each task:

- [ ] Read relevant section of this reference document
- [ ] Review [`go-style-guide.md`](../../go-style-guide.md) patterns
- [ ] Review [`go-testing-guide.md`](../../go-testing-guide.md) test patterns
- [ ] Check [`project-structure.md`](../../project-structure.md) for correct file paths
- [ ] Review relevant learnings from editor-domain.md or ui-domain.md
- [ ] Write tests FIRST (TDD)
- [ ] Implement code to pass tests
- [ ] Refactor for clarity
- [ ] Verify all tests pass
- [ ] Check style compliance
- [ ] Create checkpoint document

---

## Next Steps

After completing M4:
- **M5**: Auto-save (builds on buffer and workspace)
- **M6**: Undo/Redo (extends buffer with history)
- **M11**: Placeholder Parser (uses buffer content)
