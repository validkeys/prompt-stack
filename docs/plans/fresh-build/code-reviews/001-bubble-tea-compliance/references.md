# Bubble Tea Compliance - Implementation References

**Purpose:** Detailed code examples, templates, and specifications for implementing the compliance plan.

## Code Examples

### State Mutation Fixes

#### Example 1: App Model Active Panel
**File:** `archive/code/ui/app/model.go:103`

```go
// BEFORE (incorrect):
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.Type == tea.KeyCtrlP {
            m.activePanel = "palette"  // Direct mutation
            return m, nil
        }
    }
    return m, nil
}

// AFTER (correct):
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.Type == tea.KeyCtrlP {
            newModel := m
            newModel.activePanel = "palette"
            return newModel, nil
        }
    }
    return m, nil
}
```

#### Example 2: Workspace Model Cursor Movement
**File:** `archive/code/ui/workspace/model.go`

```go
// BEFORE (incorrect):
func (m *Model) moveCursorUp() {
    m.cursor.y--
    if m.cursor.y < 0 {
        m.cursor.y = 0
    }
}

// AFTER (correct):
func (m *Model) moveCursorUp() Model {
    newModel := m
    newModel.cursor.y--
    if newModel.cursor.y < 0 {
        newModel.cursor.y = 0
    }
    return newModel
}
```

#### Example 3: Palette Model State Update
**File:** `archive/code/ui/palette/model.go`

```go
// BEFORE (incorrect):
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.Type == tea.KeyEnter {
            m.visible = false  // Direct mutation
            return m, nil
        }
    }
    return m, nil
}

// AFTER (correct):
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.Type == tea.KeyEnter {
            newModel := m
            newModel.visible = false
            return newModel, nil
        }
    }
    return m, nil
}
```

### Copy Helper Pattern

```go
// For complex models with many fields
func (m *Model) copy() Model {
    return Model{
        field1:      m.field1,
        field2:      m.field2,
        field3:      m.field3,
        // ... all fields
    }
}

// Usage:
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    newModel := m.copy()
    newModel.field1 = newValue
    return newModel, nil
}
```

## Test Templates

### Effect Testing Template

```go
func TestCharacterInputUpdatesContent(t *testing.T) {
    model := workspace.New("")
    
    // Type text using helper
    model = testutil.TypeText(model, "hello")
    
    // Verify effect - content changed
    if model.GetContent() != "hello" {
        t.Errorf("expected 'hello', got '%s'", model.GetContent())
    }
}

func TestCursorMovement(t *testing.T) {
    model := workspace.New("test")
    
    // Move cursor right
    model = testutil.PressKey(model, tea.KeyRight)
    
    // Verify effect - cursor position changed
    if model.GetCursorX() != 1 {
        t.Errorf("expected cursor at x=1, got x=%d", model.GetCursorX())
    }
}
```

### Command Verification Template

```go
func TestQuitReturnsCommand(t *testing.T) {
    model := app.New()
    
    // Send quit message
    _, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
    
    // Verify command is returned
    if cmd == nil {
        t.Fatal("expected command, got nil")
    }
    
    // Verify it's a quit command
    if cmd != tea.Quit {
        t.Errorf("expected tea.Quit, got %v", cmd)
    }
}

func TestAISuggestionReturnsCommand(t *testing.T) {
    model := app.New()
    
    // Trigger AI suggestion
    _, cmd := model.Update(AISuggestionRequestMsg{})
    
    // Verify command is returned
    if cmd == nil {
        t.Fatal("expected command, got nil")
    }
    
    // Verify command type (implementation-specific)
    // cmd should be a tea.Cmd that triggers AI generation
}
```

### Message Coverage Template

```go
func TestMessageCoverage(t *testing.T) {
    tests := []struct {
        name     string
        msg      tea.Msg
        wantCmd  bool
    }{
        {"resize", tea.WindowSizeMsg{Width: 80, Height: 24}, false},
        {"key enter", tea.KeyMsg{Type: tea.KeyEnter}, false},
        {"key escape", tea.KeyMsg{Type: tea.KeyEsc}, false},
        {"mouse click", tea.MouseMsg{Type: tea.MouseLeft}, false},
        {"custom message", CustomMsg{}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            model := workspace.New("test")
            _, cmd := model.Update(tt.msg)
            
            if tt.wantCmd && cmd == nil {
                t.Errorf("expected command for %s", tt.name)
            }
        })
    }
}
```

### Input Sequence Template

```go
func TestTypingWorkflow(t *testing.T) {
    model := workspace.New("")
    
    // Simulate typing workflow
    model = testutil.TypeText(model, "Hello")
    model = testutil.PressKey(model, tea.KeySpace)
    model = testutil.TypeText(model, "World")
    model = testutil.PressKey(model, tea.KeyEnter)
    
    // Verify final state
    expected := "Hello World\n"
    if model.GetContent() != expected {
        t.Errorf("expected '%s', got '%s'", expected, model.GetContent())
    }
}

func TestEditWorkflow(t *testing.T) {
    model := workspace.New("Hello World")
    
    // Move to end, delete, type new text
    model = testutil.PressKey(model, tea.KeyRight)  // Move right 11 times
    model = testutil.PressKey(model, tea.KeyRight)
    model = testutil.PressKey(model, tea.KeyRight)
    model = testutil.PressKey(model, tea.KeyRight)
    model = testutil.PressKey(model, tea.KeyRight)
    model = testutil.PressKey(model, tea.KeyRight)
    model = testutil.PressKey(model, tea.KeyRight)
    model = testutil.PressKey(model, tea.KeyRight)
    model = testutil.PressKey(model, tea.KeyRight)
    model = testutil.PressKey(model, tea.KeyRight)
    model = testutil.PressKey(model, tea.KeyRight)
    
    model = testutil.PressKey(model, tea.KeyBackspace)  // Delete 'd'
    model = testutil.TypeText(model, "k")  // Type 'k'
    
    // Verify
    if model.GetContent() != "Hello Work" {
        t.Errorf("expected 'Hello Work', got '%s'", model.GetContent())
    }
}
```

### Edge Case Template

```go
func TestEdgeCases(t *testing.T) {
    tests := []struct {
        name     string
        setup    func() workspace.Model
        action   func(workspace.Model) workspace.Model
        verify   func(workspace.Model) error
    }{
        {
            name: "empty content cursor movement",
            setup: func() workspace.Model { return workspace.New("") },
            action: func(m workspace.Model) workspace.Model {
                return testutil.PressKey(m, tea.KeyRight)
            },
            verify: func(m workspace.Model) error {
                if m.GetCursorX() != 0 {
                    return fmt.Errorf("cursor should stay at 0")
                }
                return nil
            },
        },
        {
            name: "unicode characters",
            setup: func() workspace.Model { return workspace.New("") },
            action: func(m workspace.Model) workspace.Model {
                return testutil.TypeText(m, "Hello 世界 🌍")
            },
            verify: func(m workspace.Model) error {
                if m.GetContent() != "Hello 世界 🌍" {
                    return fmt.Errorf("unicode not handled correctly")
                }
                return nil
            },
        },
        {
            name: "very long line",
            setup: func() workspace.Model { return workspace.New("") },
            action: func(m workspace.Model) workspace.Model {
                longText := strings.Repeat("a", 2000)
                return testutil.TypeText(m, longText)
            },
            verify: func(m workspace.Model) error {
                if len(m.GetContent()) != 2000 {
                    return fmt.Errorf("long line not handled correctly")
                }
                return nil
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            model := tt.setup()
            model = tt.action(model)
            if err := tt.verify(model); err != nil {
                t.Error(err)
            }
        })
    }
}
```

## Component Interface Specifications

### Cursor Component
**Location:** `internal/editor/cursor.go`

```go
package cursor

import tea "github.com/charmbracelet/bubbletea"

// Model represents cursor position and movement
type Model struct {
    x    int  // Column position
    y    int  // Line position
    line int  // Current line index
}

// New creates a new cursor at origin
func New() Model {
    return Model{x: 0, y: 0, line: 0}
}

// Update handles cursor-related messages
func (m Model) Update(msg tea.Msg) Model {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyUp:
            return m.moveUp()
        case tea.KeyDown:
            return m.moveDown()
        case tea.KeyLeft:
            return m.moveLeft()
        case tea.KeyRight:
            return m.moveRight()
        }
    }
    return m
}

// moveUp moves cursor up one line
func (m Model) moveUp() Model {
    newModel := m
    newModel.y--
    if newModel.y < 0 {
        newModel.y = 0
    }
    return newModel
}

// moveDown moves cursor down one line
func (m Model) moveDown() Model {
    newModel := m
    newModel.y++
    return newModel
}

// moveLeft moves cursor left one character
func (m Model) moveLeft() Model {
    newModel := m
    newModel.x--
    if newModel.x < 0 {
        newModel.x = 0
    }
    return newModel
}

// moveRight moves cursor right one character
func (m Model) moveRight() Model {
    newModel := m
    newModel.x++
    return newModel
}

// Position returns current cursor position
func (m Model) Position() (x, y int) {
    return m.x, m.y
}
```

### Viewport Component
**Location:** `internal/editor/viewport.go`

```go
package viewport

import tea "github.com/charmbracelet/bubbletea"

// Model manages viewport scrolling and visible area
type Model struct {
    topLine    int  // First visible line
    height     int  // Viewport height in lines
    totalLines int  // Total lines in document
}

// New creates a new viewport
func New(height int) Model {
    return Model{
        topLine:    0,
        height:     height,
        totalLines: 0,
    }
}

// Update handles viewport-related messages
func (m Model) Update(msg tea.Msg) Model {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        newModel := m
        newModel.height = msg.Height
        return newModel
    case ScrollMsg:
        return m.scrollTo(msg.Line)
    }
    return m
}

// scrollTo moves viewport to show specified line
func (m Model) scrollTo(line int) Model {
    newModel := m
    newModel.topLine = line
    return newModel
}

// VisibleLines returns range of visible line numbers
func (m Model) VisibleLines() (start, end int) {
    return m.topLine, m.topLine + m.height
}

// EnsureVisible ensures line is visible in viewport
func (m Model) EnsureVisible(line int) Model {
    start, end := m.VisibleLines()
    
    if line < start {
        return m.scrollTo(line)
    }
    if line >= end {
        return m.scrollTo(line - m.height + 1)
    }
    return m
}
```

### Placeholder Component
**Location:** `internal/editor/placeholder.go`

```go
package placeholder

import tea "github.com/charmbracelet/bubbletea"

// Manager handles placeholder detection and editing
type Manager struct {
    placeholders []Placeholder
    activeIndex  int
    isEditing    bool
}

type Placeholder struct {
    ID      string
    Content string
    Start   int
    End     int
}

// New creates a new placeholder manager
func New() Manager {
    return Manager{
        placeholders: []Placeholder{},
        activeIndex:  -1,
        isEditing:    false,
    }
}

// Update handles placeholder-related messages
func (m Manager) Update(msg tea.Msg) Manager {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.Type == tea.KeyTab {
            return m.nextPlaceholder()
        }
        if msg.Type == tea.KeyEnter && m.isEditing {
            return m.exitEditMode()
        }
    case EnterPlaceholderMsg:
        return m.enterEditMode(msg.ID)
    }
    return m
}

// nextPlaceholder moves to next placeholder
func (m Manager) nextPlaceholder() Manager {
    newModel := m
    newModel.activeIndex = (newModel.activeIndex + 1) % len(newModel.placeholders)
    return newModel
}

// enterEditMode enters edit mode for placeholder
func (m Manager) enterEditMode(id string) Manager {
    newModel := m
    for i, p := range newModel.placeholders {
        if p.ID == id {
            newModel.activeIndex = i
            newModel.isEditing = true
            break
        }
    }
    return newModel
}

// exitEditMode exits placeholder edit mode
func (m Manager) exitEditMode() Manager {
    newModel := m
    newModel.isEditing = false
    return newModel
}

// Active returns currently active placeholder
func (m Manager) Active() *Placeholder {
    if m.activeIndex < 0 || m.activeIndex >= len(m.placeholders) {
        return nil
    }
    return &m.placeholders[m.activeIndex]
}

// IsEditing returns true if editing a placeholder
func (m Manager) IsEditing() bool {
    return m.isEditing
}
```

### File I/O Component
**Location:** `internal/editor/fileio.go`

```go
package fileio

import (
    "os"
    "path/filepath"
)

// Manager handles file save/load operations
type Manager struct {
    filepath string
    modified bool
}

// New creates a new file I/O manager
func New(path string) Manager {
    return Manager{
        filepath: path,
        modified: false,
    }
}

// Load reads file content
func (m Manager) Load() (string, error) {
    if m.filepath == "" {
        return "", nil
    }
    
    content, err := os.ReadFile(m.filepath)
    if err != nil {
        return "", err
    }
    
    return string(content), nil
}

// Save writes content to file
func (m Manager) Save(content string) error {
    if m.filepath == "" {
        return nil
    }
    
    // Ensure directory exists
    dir := filepath.Dir(m.filepath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    
    return os.WriteFile(m.filepath, []byte(content), 0644)
}

// MarkModified marks file as modified
func (m Manager) MarkModified() Manager {
    newModel := m
    newModel.modified = true
    return newModel
}

// IsModified returns true if file has unsaved changes
func (m Manager) IsModified() bool {
    return m.modified
}

// Path returns the file path
func (m Manager) Path() string {
    return m.filepath
}
```

## Integration Points

### Workspace Model After Refactoring
**File:** `archive/code/ui/workspace/model.go`

```go
package workspace

import (
    "github.com/charmbracelet/bubbletea"
    "promptstack/internal/editor/cursor"
    "promptstack/internal/editor/viewport"
    "promptstack/internal/editor/placeholder"
    "promptstack/internal/editor/fileio"
)

type Model struct {
    cursor      *cursor.Model
    viewport    *viewport.Model
    placeholders *placeholder.Manager
    fileio      *fileio.Manager
    content     string
    lines       []string
}

func New(content string) Model {
    return Model{
        cursor:      cursor.New(),
        viewport:    viewport.New(24),
        placeholders: placeholder.New(),
        fileio:      fileio.New(""),
        content:     content,
        lines:       splitLines(content),
    }
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // Update cursor
    m.cursor = m.cursor.Update(msg)
    
    // Update viewport
    m.viewport = m.viewport.Update(msg)
    
    // Update placeholders
    m.placeholders = m.placeholders.Update(msg)
    
    // Handle other messages
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKey(msg)
    case tea.WindowSizeMsg:
        return m.handleResize(msg)
    }
    
    return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
    switch msg.Type {
    case tea.KeyRunes:
        return m.insertRunes(msg.Runes)
    case tea.KeyBackspace:
        return m.deleteChar()
    }
    return m, nil
}

func (m Model) handleResize(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
    newModel := m
    newModel.viewport = viewport.New(msg.Height)
    return newModel, nil
}
```

## Task Dependencies

```
Phase 1 Dependencies:
1.1 Fix Mutations (app, palette) → 1.2 Create Test Utils → 1.3 Add Tests
1.1 Fix Mutations (workspace) → 1.3 Add Tests

Phase 2 Dependencies:
1.3 Add Tests → 2.1 Extract Cursor → 2.1 Extract Viewport → 2.1 Extract Placeholder → 2.1 Extract FileIO → 2.2 Improve Coverage

Phase 3 Dependencies:
2.2 Improve Coverage → 3.1 Performance → 3.2 Infrastructure → 3.3 Developer Experience
```

## Verification Commands

### After State Mutation Fixes
```bash
# Run tests for specific model
go test ./ui/app/... -v
go test ./ui/workspace/... -v
go test ./ui/palette/... -v

# Check for direct mutations (manual review)
grep -n "m\." archive/code/ui/app/model.go | grep "="
```

### After Test Utilities Creation
```bash
# Test the test utilities
go test ./internal/testutil/... -v

# Verify helpers work
go test -run TestTypeText ./internal/testutil/...
```

### After Adding Tests
```bash
# Run all tests
go test ./... -v

# Check coverage
go test ./... -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### After Workspace Refactoring
```bash
# Run workspace tests
go test ./ui/workspace/... -v

# Run integration tests
go test ./cmd/promptstack/... -v

# Verify no regressions
go test ./... -race
```

## Specific Test Cases

### Edge Cases to Test
1. **Empty Content**
   - Cursor movement on empty document
   - Backspace on empty document
   - Enter on empty document

2. **Unicode Characters**
   - Multi-byte characters (emoji, CJK)
   - Cursor movement through multi-byte chars
   - Deletion of multi-byte chars

3. **Very Long Lines**
   - Lines > 1000 characters
   - Cursor movement on long lines
   - Viewport scrolling with long lines

4. **Window Resize**
   - Resize to 0x0
   - Resize to very large dimensions
   - Rapid resize events

5. **Rapid Input**
   - Type 100+ characters quickly
   - Rapid key presses
   - Concurrent message handling

6. **Placeholder Edge Cases**
   - Empty placeholder list
   - Navigate past last placeholder
   - Edit placeholder with special chars

7. **File I/O Edge Cases**
   - Save to non-existent directory
   - Load non-existent file
   - Save empty file

8. **Undo/Redo Edge Cases**
   - Undo with empty stack
   - Redo with empty stack
   - Undo after file save

9. **Viewport Edge Cases**
   - Document shorter than viewport
   - Cursor at document boundaries
   - Scroll past document end

10. **Memory Edge Cases**
    - Very large document (10MB+)
    - Many rapid edits
    - Long undo/redo history

## Error Handling Strategy

### When Tests Fail After Refactoring

1. **Identify Failure**
   ```bash
   go test ./ui/workspace/... -v -run TestName
   ```

2. **Check Integration Points**
   - Verify message passing between components
   - Check component initialization
   - Validate state updates

3. **Debug Step-by-Step**
   - Add logging to Update methods
   - Verify state before/after each message
   - Check command execution

4. **Rollback Strategy**
   - If stuck, revert last change
   - Re-apply with more tests
   - Use git bisect if needed

### Performance Regression Detection

```bash
# Run benchmarks before refactoring
go test -bench=. ./ui/workspace/... > before.txt

# Run benchmarks after refactoring
go test -bench=. ./ui/workspace/... > after.txt

# Compare results
benchcmp before.txt after.txt
```

## References

- [`audit-tracking.md`](./audit-tracking.md:1) - Full audit findings
- [`plan.md`](./plan.md:1) - Implementation plan
- [`bubble-tea-testing-best-practices.md`](../../bubble-tea-testing-best-practices.md:1) - Testing guidelines
- [`bubble-tea-best-practices.md`](../../../../bubble-tea-best-practices.md:1) - Architecture guidelines
- [`go-testing-guide.md`](../../go-testing-guide.md:1) - Go testing patterns