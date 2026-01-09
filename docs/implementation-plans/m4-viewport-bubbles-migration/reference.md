# Task 5: Viewport with Scrolling - Reference Document

**Date:** 2026-01-09  
**Related:** task-list.md

---

## How to Use This Document

**Read this section when:**
- Before implementing Task 2 (Workspace Model Migration) - Read Architecture Context and API Mappings
- When implementing Task 4 (Cursor Visibility) - Read Middle-Third Scrolling Strategy section
- When writing tests (Task 5) - Read Testing Patterns and Test Examples sections
- When debugging viewport issues - Read Troubleshooting section

**Key sections:**
- Lines 30-120: Architecture Context - Read before Task 2
- Lines 121-220: API Reference and Mappings - Reference during Task 2 implementation
- Lines 221-320: Middle-Third Scrolling Strategy - Read before Task 4
- Lines 321-420: Code Examples - Reference during implementation
- Lines 421-520: Testing Patterns - Read before Task 5
- Lines 521-600: Troubleshooting - Consult when debugging

**Related documents:**
- See [`task-list.md`](task-list.md) for task breakdown and acceptance criteria
- Cross-reference with [`ui-domain.md`](../../learnings/ui-domain.md) for viewport patterns
- Cross-reference with [`editor-domain.md`](../../learnings/editor-domain.md) for editor integration

---

## Architecture Context

### Domain Overview

The viewport migration affects **two domains**:

1. **UI Domain** (`ui/workspace/`)
   - Workspace model owns viewport instance
   - Handles viewport initialization and lifecycle
   - Coordinates cursor visibility with viewport scrolling
   - Responsible for rendering viewport content

2. **Editor Domain** (`internal/editor/`)
   - Buffer provides content to viewport
   - Cursor position drives viewport scrolling
   - Buffer modifications trigger viewport content updates
   - No direct dependency on viewport (decoupled)

### Package Structure

Per [`project-structure.md`](../../project-structure.md):

```
ui/
└── workspace/
    ├── model.go              # Workspace model with viewport
    ├── view.go               # Rendering logic
    ├── viewport_cursor.go    # Cursor visibility logic (NEW)
    ├── model_test.go         # Model tests
    ├── view_test.go          # View tests
    ├── viewport_cursor_test.go        # Cursor tests (NEW)
    └── viewport_integration_test.go   # Integration tests (NEW)

internal/
└── editor/
    ├── buffer.go             # Text buffer (unchanged)
    ├── cursor.go             # Cursor management (unchanged)
    └── viewport.go           # TO BE DELETED

docs/
└── archived/
    └── viewport-custom-v1.go # Archived implementation
```

### Dependencies

**External Dependencies** (per [`DEPENDENCIES.md`](../../DEPENDENCIES.md)):
```
github.com/charmbracelet/bubbletea v0.25.0
github.com/charmbracelet/bubbles v0.18.0  # NEW
```

**Internal Dependencies**:
```
ui/workspace → internal/editor (buffer, cursor)
ui/workspace → github.com/charmbracelet/bubbles/viewport
```

**Dependency Flow**:
```
main.go
  → ui/app (TUI app)
    → ui/workspace (workspace model)
      → bubbles/viewport (viewport component)
      → internal/editor/buffer (text buffer)
      → internal/editor/cursor (cursor position)
```

### Key Architectural Decisions

1. **Why Bubbles Viewport?**
   - **Production-proven**: Used by Glow, Charm's official editor, and thousands of apps
   - **Automatic bounds checking**: Prevents negative scroll positions (bug in custom impl)
   - **Edge case handling**: Handles small viewports, empty content, resize events
   - **Performance**: Optimized for large documents (>10,000 lines)
   - **Maintenance**: Active development, community support, well-documented

2. **Why Keep Cursor Logic Custom?**
   - Bubbles viewport provides scrolling primitives, not cursor visibility strategies
   - Middle-third scrolling is app-specific (not all apps use this strategy)
   - Allows flexibility for future enhancements (e.g., configurable scroll margins)
   - Follows separation of concerns: viewport handles scrolling, workspace handles cursor

3. **Why Middle-Third Scrolling?**
   - Provides optimal context: user sees 1-2 lines above and below cursor
   - Reduces jarring jumps: cursor doesn't snap to top/bottom edges
   - Matches user expectations: similar to vim's `scrolloff` behavior
   - Proven pattern: documented in [`ui-domain.md`](../../learnings/ui-domain.md) Category 2

---

## API Reference and Mappings

### Bubbles Viewport API

**Initialization**:
```go
import "github.com/charmbracelet/bubbles/viewport"

// Create viewport with dimensions
vp := viewport.New(width int, height int) viewport.Model

// Alternative: create and configure
vp := viewport.Model{
    Width:  80,
    Height: 24,
}
```

**Content Management**:
```go
// Set content (triggers internal layout recalculation)
vp.SetContent(content string)

// Get content
content := vp.View() // Returns visible portion
```

**Scroll Position**:
```go
// Get/set scroll offset (0-based line number)
offset := vp.YOffset         // Read current offset
vp.SetYOffset(10)            // Set offset to line 10
vp.YOffset = 15              // Direct assignment also works

// Relative scrolling
vp.LineUp(5)                 // Scroll up 5 lines
vp.LineDown(3)               // Scroll down 3 lines
vp.HalfViewUp()              // Scroll up half viewport height
vp.HalfViewDown()            // Scroll down half viewport height
```

**Viewport Dimensions**:
```go
// Get viewport dimensions
width := vp.Width
height := vp.Height

// Update dimensions (typically on resize)
vp.Width = newWidth
vp.Height = newHeight
```

**Scrolling Queries**:
```go
// Check scroll state
atTop := vp.AtTop()          // true if at document start
atBottom := vp.AtBottom()    // true if at document end
canScrollUp := !vp.AtTop()
canScrollDown := !vp.AtBottom()

// Get total height
totalHeight := vp.TotalLineCount() // Total lines in content
```

### Custom Viewport → Bubbles Viewport Mapping

| Custom Viewport API | Bubbles Viewport API | Notes |
|-------------------|---------------------|-------|
| `viewport.TopLine()` | `viewport.YOffset` | Read property directly |
| `viewport.ScrollTo(line int)` | `viewport.SetYOffset(line)` | Method call |
| `viewport.EnsureVisible(line int)` | **Custom logic** | Implement with YOffset + middle-third |
| `viewport.SetTotalLines(n int)` | `viewport.SetContent(content)` | Handled by content |
| `viewport.VisibleLines()` | `viewport.YOffset` to `YOffset + Height` | Calculate range |
| `viewport.Render()` | `viewport.View()` | Returns string |
| `viewport.Update(msg)` | `viewport.Update(msg)` | Same signature |
| `viewport.Init()` | `viewport.Init()` | Same signature |

### Migration Gotchas

1. **Content vs. Line Count**
   - **Custom**: Tracked line count separately with `SetTotalLines(n)`
   - **Bubbles**: Infers line count from `SetContent(content string)`
   - **Implication**: Must call `SetContent()` whenever buffer changes

2. **YOffset Bounds**
   - **Custom**: Manual bounds checking (had bug with negative values)
   - **Bubbles**: Automatic bounds checking in `SetYOffset()`
   - **Implication**: No need to check bounds manually, but still check in cursor logic

3. **Cursor Visibility**
   - **Custom**: `EnsureVisible(line)` method
   - **Bubbles**: No built-in cursor visibility
   - **Implication**: Must implement custom cursor visibility logic (Task 4)

4. **Update() Return Value**
   - **Custom**: May have returned `(viewport.Viewport, tea.Cmd)`
   - **Bubbles**: Returns `(viewport.Model, tea.Cmd)`
   - **Implication**: Type mismatch if using wrong type

---

## Middle-Third Scrolling Strategy

### Visual Representation

```
┌─────────────────────────────────┐
│ Line 1                          │
│ Line 2                          │ ← Top third starts
│ Line 3                          │
├─────────────────────────────────┤
│ Line 4    ← Middle third starts │
│ Line 5                          │
│ Line 6    ← CURSOR HERE         │ ← Cursor in middle third (optimal)
│ Line 7                          │
│ Line 8    ← Middle third ends   │
├─────────────────────────────────┤
│ Line 9                          │
│ Line 10   ← Bottom third starts │
│ Line 11                         │
│ Line 12                         │
└─────────────────────────────────┘
```

**Viewport height**: 12 lines  
**Top third**: Lines 1-3 (lines 0-3 in 0-indexed)  
**Middle third**: Lines 4-8 (lines 3-8 in 0-indexed)  
**Bottom third**: Lines 9-12 (lines 8-12 in 0-indexed)

### Scrolling Rules

1. **Cursor in middle third**: No scrolling needed ✅
2. **Cursor in top third**: Scroll viewport up (decrease YOffset)
3. **Cursor in bottom third**: Scroll viewport down (increase YOffset)

### Algorithm

```
Given:
- viewportHeight = viewport.Height
- cursorY = cursor line (0-based)
- currentOffset = viewport.YOffset

Calculate:
- third = viewportHeight / 3
- middleTop = currentOffset + third
- middleBottom = currentOffset + viewportHeight - third

Rules:
IF cursorY < middleTop:
    newOffset = max(0, cursorY - third)
    viewport.SetYOffset(newOffset)
    
ELSE IF cursorY >= middleBottom:
    newOffset = cursorY - viewportHeight + third
    viewport.SetYOffset(newOffset)
    
ELSE:
    # Cursor in middle third, no scroll needed
```

### Edge Cases

**Case 1: Viewport taller than document**
```
Document: 10 lines
Viewport: 20 lines
Solution: YOffset = 0, cursor can be anywhere
```

**Case 2: Cursor at document start**
```
Cursor at line 0
Solution: YOffset = 0 (always)
```

**Case 3: Cursor at document end**
```
Document: 100 lines
Viewport: 20 lines
Cursor at line 99
Solution: YOffset = max(0, 99 - 20 + third) = max(0, 79 + 6) = 85
But document only has 100 lines, so YOffset = min(85, 100 - 20) = 80
```

**Case 4: One-line viewport**
```
Viewport height = 1
third = 1/3 = 0
Middle region = entire viewport
Solution: Always keep cursor at top (YOffset = cursorY)
```

### Key Learning Applied

From [`ui-domain.md`](../../learnings/ui-domain.md) Category 2:

> **Learning**: Always keep cursor visible by adjusting viewport. Use a "middle third" strategy to provide smooth scrolling. This prevents cursor from getting stuck at edges.

This pattern is proven and well-documented from previous implementation. We apply the exact same logic but integrate with bubbles viewport instead of custom viewport.

---

## Code Examples

### Example 1: Workspace Model with Bubbles Viewport

**File**: `ui/workspace/model.go`

```go
package workspace

import (
    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/viewport"
    
    "promptstack/internal/editor"
)

// Model represents the workspace with text editor and viewport.
type Model struct {
    buffer   *editor.Buffer
    cursor   *editor.Cursor
    viewport viewport.Model  // Bubbles viewport
    width    int
    height   int
}

// New creates a new workspace model with initialized viewport.
func New(width, height int) Model {
    // Create viewport with dimensions
    vp := viewport.New(width, height-1) // Reserve 1 line for status bar
    
    return Model{
        buffer:   editor.NewBuffer(),
        cursor:   editor.NewCursor(),
        viewport: vp,
        width:    width,
        height:   height,
    }
}

// Init initializes the workspace model.
func (m Model) Init() tea.Cmd {
    return m.viewport.Init()
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        // Handle terminal resize
        m.width = msg.Width
        m.height = msg.Height
        
        // Update viewport dimensions
        m.viewport.Width = msg.Width
        m.viewport.Height = msg.Height - 1 // Reserve status bar
        
        // Sync content after resize
        m.syncViewportContent()
        
        return m, nil
        
    case tea.KeyMsg:
        // Handle keyboard input
        switch msg.String() {
        case "up":
            m.cursor.MoveUp()
            m.ensureCursorVisible()
            return m, nil
            
        case "down":
            m.cursor.MoveDown()
            m.ensureCursorVisible()
            return m, nil
            
        // ... other keys
        }
    }
    
    // Forward message to viewport
    m.viewport, cmd = m.viewport.Update(msg)
    return m, cmd
}

// syncViewportContent synchronizes buffer content with viewport.
func (m *Model) syncViewportContent() {
    content := m.buffer.String()
    m.viewport.SetContent(content)
}

// View renders the workspace.
func (m Model) View() string {
    // Sync content before rendering
    m.syncViewportContent()
    
    // Render viewport
    return m.viewport.View()
}
```

### Example 2: Cursor Visibility with Middle-Third Scrolling

**File**: `ui/workspace/viewport_cursor.go`

```go
package workspace

import (
    "math"
)

// ensureCursorVisible adjusts viewport to keep cursor visible using middle-third strategy.
// 
// This implements the middle-third scrolling strategy documented in learnings/ui-domain.md.
// The cursor is kept in the middle third of the viewport for optimal context visibility.
func (m *Model) ensureCursorVisible() {
    if m.buffer == nil || m.cursor == nil {
        return // Safety check
    }
    
    // Get cursor position (0-based line number)
    _, cursorY := m.cursor.Position()
    
    // Get current viewport state
    viewportTop := m.viewport.YOffset
    viewportHeight := m.viewport.Height
    viewportBottom := viewportTop + viewportHeight
    
    // Calculate middle third boundaries
    third := int(math.Max(1, float64(viewportHeight)/3.0))
    middleTop := viewportTop + third
    middleBottom := viewportBottom - third
    
    // Check if cursor is outside middle third
    if cursorY < middleTop && cursorY >= 0 {
        // Cursor in top third - scroll up
        newOffset := int(math.Max(0, float64(cursorY-third)))
        m.viewport.SetYOffset(newOffset)
        
    } else if cursorY >= middleBottom {
        // Cursor in bottom third - scroll down
        newOffset := cursorY - viewportHeight + third
        
        // Ensure we don't scroll past document end
        totalLines := m.buffer.LineCount()
        maxOffset := int(math.Max(0, float64(totalLines-viewportHeight)))
        newOffset = int(math.Min(float64(newOffset), float64(maxOffset)))
        
        m.viewport.SetYOffset(newOffset)
    }
    
    // If cursor in middle third, no action needed
}

// scrollToCursor scrolls viewport to make cursor visible at top of viewport.
// This is used for explicit scroll operations (e.g., "zz" in vim mode).
func (m *Model) scrollToCursor() {
    if m.cursor == nil {
        return
    }
    
    _, cursorY := m.cursor.Position()
    m.viewport.SetYOffset(cursorY)
}

// centerCursor centers the cursor in the viewport.
// This is used for "zz" command in vim mode (future M34-M35).
func (m *Model) centerCursor() {
    if m.cursor == nil || m.buffer == nil {
        return
    }
    
    _, cursorY := m.cursor.Position()
    viewportHeight := m.viewport.Height
    
    // Center cursor in viewport
    newOffset := cursorY - (viewportHeight / 2)
    
    // Clamp to valid range
    totalLines := m.buffer.LineCount()
    maxOffset := int(math.Max(0, float64(totalLines-viewportHeight)))
    newOffset = int(math.Max(0, math.Min(float64(newOffset), float64(maxOffset))))
    
    m.viewport.SetYOffset(newOffset)
}
```

### Example 3: View Rendering

**File**: `ui/workspace/view.go`

```go
package workspace

import (
    "fmt"
    "strings"
    
    "github.com/charmbracelet/lipgloss"
)

// View renders the workspace with viewport and status bar.
func (m Model) View() string {
    // Sync buffer content to viewport before rendering
    m.syncViewportContent()
    
    // Render viewport (returns visible content)
    viewportContent := m.viewport.View()
    
    // Render status bar
    statusBar := m.renderStatusBar()
    
    // Combine viewport and status bar
    return fmt.Sprintf("%s\n%s", viewportContent, statusBar)
}

// renderStatusBar renders the status bar with cursor position and line count.
func (m Model) renderStatusBar() string {
    if m.cursor == nil || m.buffer == nil {
        return ""
    }
    
    // Get cursor position (1-based for display)
    cursorX, cursorY := m.cursor.Position()
    line := cursorY + 1
    col := cursorX + 1
    
    // Get buffer stats
    totalLines := m.buffer.LineCount()
    totalChars := m.buffer.CharCount()
    
    // Format status bar
    leftStatus := fmt.Sprintf(" Ln %d, Col %d ", line, col)
    rightStatus := fmt.Sprintf(" %d lines, %d chars ", totalLines, totalChars)
    
    // Calculate padding
    padding := m.width - len(leftStatus) - len(rightStatus)
    if padding < 0 {
        padding = 0
    }
    
    // Apply styling (per opencode-design-system.md)
    statusStyle := lipgloss.NewStyle().
        Background(lipgloss.Color("236")).
        Foreground(lipgloss.Color("250")).
        Bold(true)
    
    statusText := leftStatus + strings.Repeat(" ", padding) + rightStatus
    return statusStyle.Render(statusText)
}
```

### Example 4: Error Handling

**File**: `ui/workspace/model.go` (error handling additions)

```go
package workspace

import (
    "fmt"
    
    "promptstack/internal/platform/errors"
)

// validateViewportState checks for invalid viewport state and returns error.
func (m *Model) validateViewportState() error {
    // Check viewport dimensions
    if m.viewport.Width <= 0 || m.viewport.Height <= 0 {
        return errors.New(errors.ErrorTypeValidation, "invalid viewport dimensions").
            WithDetails(fmt.Sprintf("width=%d, height=%d", m.viewport.Width, m.viewport.Height)).
            WithSeverity(errors.SeverityError)
    }
    
    // Check YOffset bounds
    if m.viewport.YOffset < 0 {
        return errors.New(errors.ErrorTypeValidation, "negative viewport offset").
            WithDetails(fmt.Sprintf("offset=%d", m.viewport.YOffset)).
            WithSeverity(errors.SeverityWarning). // Warning, will auto-correct
            WithRetryable(true)
    }
    
    return nil
}

// syncViewportContent synchronizes buffer content with viewport and handles errors.
func (m *Model) syncViewportContent() {
    if m.buffer == nil {
        // Log warning but don't crash
        // TODO: Add logger reference when available
        return
    }
    
    // Get buffer content
    content := m.buffer.String()
    
    // Update viewport content
    m.viewport.SetContent(content)
    
    // Validate state
    if err := m.validateViewportState(); err != nil {
        // Handle validation error
        if appErr, ok := err.(*errors.AppError); ok {
            if appErr.Retryable {
                // Auto-correct negative offset
                if m.viewport.YOffset < 0 {
                    m.viewport.SetYOffset(0)
                }
            }
        }
    }
}
```

---

## Testing Patterns

### Pattern 1: Table-Driven Tests for Cursor Visibility

**File**: `ui/workspace/viewport_cursor_test.go`

Per [`go-testing-guide.md`](../../go-testing-guide.md) and [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md):

```go
package workspace_test

import (
    "testing"
    
    "promptstack/ui/workspace"
    "promptstack/internal/editor"
)

func TestEnsureCursorVisible_MiddleThirdScrolling(t *testing.T) {
    tests := []struct {
        name            string
        viewportHeight  int
        totalLines      int
        cursorY         int
        initialOffset   int
        expectedOffset  int
        expectedScroll  bool
    }{
        {
            name:           "cursor in middle third, no scroll",
            viewportHeight: 12,
            totalLines:     100,
            cursorY:        6, // Middle of viewport
            initialOffset:  0,
            expectedOffset: 0,
            expectedScroll: false,
        },
        {
            name:           "cursor in top third, scroll up",
            viewportHeight: 12,
            totalLines:     100,
            cursorY:        2, // Top third (< third=4)
            initialOffset:  5,
            expectedOffset: 0, // max(0, 2-4) = 0
            expectedScroll: true,
        },
        {
            name:           "cursor in bottom third, scroll down",
            viewportHeight: 12,
            totalLines:     100,
            cursorY:        10, // Bottom third (>= 8)
            initialOffset:  0,
            expectedOffset: 2, // 10 - 12 + 4 = 2
            expectedScroll: true,
        },
        {
            name:           "cursor at document start",
            viewportHeight: 12,
            totalLines:     100,
            cursorY:        0,
            initialOffset:  5,
            expectedOffset: 0, // Always 0 at start
            expectedScroll: true,
        },
        {
            name:           "cursor at document end",
            viewportHeight: 12,
            totalLines:     100,
            cursorY:        99,
            initialOffset:  0,
            expectedOffset: 88, // 100 - 12 = 88
            expectedScroll: true,
        },
        {
            name:           "viewport taller than document",
            viewportHeight: 50,
            totalLines:     10,
            cursorY:        5,
            initialOffset:  0,
            expectedOffset: 0, // No scroll needed
            expectedScroll: false,
        },
        {
            name:           "one line viewport",
            viewportHeight: 1,
            totalLines:     100,
            cursorY:        50,
            initialOffset:  0,
            expectedOffset: 50, // Cursor always at top
            expectedScroll: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            m := workspace.New(80, tt.viewportHeight)
            
            // Create buffer with content
            buffer := editor.NewBuffer()
            for i := 0; i < tt.totalLines; i++ {
                buffer.InsertLine(fmt.Sprintf("Line %d", i+1))
            }
            m.SetBuffer(buffer)
            
            // Set initial scroll position
            m.Viewport().SetYOffset(tt.initialOffset)
            
            // Set cursor position
            m.Cursor().MoveTo(0, tt.cursorY)
            
            // Execute
            m.EnsureCursorVisible()
            
            // Assert
            actualOffset := m.Viewport().YOffset
            if actualOffset != tt.expectedOffset {
                t.Errorf("Expected offset %d, got %d", tt.expectedOffset, actualOffset)
            }
            
            // Verify scroll occurred (or didn't)
            scrolled := (actualOffset != tt.initialOffset)
            if scrolled != tt.expectedScroll {
                t.Errorf("Expected scroll=%v, got scroll=%v", tt.expectedScroll, scrolled)
            }
        })
    }
}
```

### Pattern 2: Integration Tests

**File**: `ui/workspace/viewport_integration_test.go`

```go
package workspace_test

import (
    "testing"
    
    "github.com/charmbracelet/bubbletea"
    
    "promptstack/ui/workspace"
    "promptstack/internal/editor"
)

func TestViewportIntegration_BufferSync(t *testing.T) {
    // Setup
    m := workspace.New(80, 24)
    buffer := editor.NewBuffer()
    buffer.Insert("Hello\nWorld\nTest")
    m.SetBuffer(buffer)
    
    // Execute: Sync content
    m.SyncViewportContent()
    
    // Assert: Viewport has correct content
    view := m.Viewport().View()
    expected := "Hello\nWorld\nTest"
    if view != expected {
        t.Errorf("Expected viewport content %q, got %q", expected, view)
    }
}

func TestViewportIntegration_CursorMovement(t *testing.T) {
    // Setup: 100-line document, viewport height 12
    m := workspace.New(80, 12)
    buffer := editor.NewBuffer()
    for i := 0; i < 100; i++ {
        buffer.InsertLine(fmt.Sprintf("Line %d", i+1))
    }
    m.SetBuffer(buffer)
    
    // Initial state: cursor at line 0, viewport at offset 0
    m.Cursor().MoveTo(0, 0)
    m.SyncViewportContent()
    
    // Execute: Move cursor down 50 times
    for i := 0; i < 50; i++ {
        msg := bubbletea.KeyMsg{Type: bubbletea.KeyDown}
        m, _ = m.Update(msg)
    }
    
    // Assert: Cursor at line 50
    _, cursorY := m.Cursor().Position()
    if cursorY != 50 {
        t.Errorf("Expected cursor at line 50, got %d", cursorY)
    }
    
    // Assert: Viewport scrolled (middle-third strategy)
    // Expected: offset = 50 - 12 + 4 = 42
    offset := m.Viewport().YOffset
    if offset < 38 || offset > 46 {
        t.Errorf("Expected viewport offset ~42, got %d", offset)
    }
}

func TestViewportIntegration_TerminalResize(t *testing.T) {
    // Setup
    m := workspace.New(120, 40)
    buffer := editor.NewBuffer()
    for i := 0; i < 100; i++ {
        buffer.InsertLine(fmt.Sprintf("Line %d", i+1))
    }
    m.SetBuffer(buffer)
    m.Cursor().MoveTo(0, 50)
    m.SyncViewportContent()
    
    initialOffset := m.Viewport().YOffset
    
    // Execute: Resize to smaller dimensions
    resizeMsg := bubbletea.WindowSizeMsg{Width: 80, Height: 24}
    m, _ = m.Update(resizeMsg)
    
    // Assert: Viewport dimensions updated
    if m.Viewport().Width != 80 {
        t.Errorf("Expected viewport width 80, got %d", m.Viewport().Width)
    }
    if m.Viewport().Height != 23 { // -1 for status bar
        t.Errorf("Expected viewport height 23, got %d", m.Viewport().Height)
    }
    
    // Assert: Cursor still visible after resize
    _, cursorY := m.Cursor().Position()
    offset := m.Viewport().YOffset
    viewportBottom := offset + m.Viewport().Height
    
    if cursorY < offset || cursorY >= viewportBottom {
        t.Errorf("Cursor not visible after resize: cursor=%d, offset=%d, bottom=%d",
            cursorY, offset, viewportBottom)
    }
}
```

### Pattern 3: Benchmark Tests

**File**: `ui/workspace/viewport_benchmark_test.go`

```go
package workspace_test

import (
    "fmt"
    "testing"
    
    "promptstack/ui/workspace"
    "promptstack/internal/editor"
)

func BenchmarkEnsureCursorVisible(b *testing.B) {
    m := workspace.New(80, 24)
    buffer := editor.NewBuffer()
    for i := 0; i < 10000; i++ {
        buffer.InsertLine(fmt.Sprintf("Line %d", i+1))
    }
    m.SetBuffer(buffer)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        m.Cursor().MoveTo(0, i%10000)
        m.EnsureCursorVisible()
    }
}

func BenchmarkViewportRender(b *testing.B) {
    m := workspace.New(80, 24)
    buffer := editor.NewBuffer()
    for i := 0; i < 1000; i++ {
        buffer.InsertLine(fmt.Sprintf("Line %d with some content", i+1))
    }
    m.SetBuffer(buffer)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = m.View()
    }
}

func BenchmarkBufferSync(b *testing.B) {
    m := workspace.New(80, 24)
    buffer := editor.NewBuffer()
    for i := 0; i < 5000; i++ {
        buffer.InsertLine(fmt.Sprintf("Line %d", i+1))
    }
    m.SetBuffer(buffer)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        m.SyncViewportContent()
    }
}
```

---

## Rollback Scenarios

### Scenario 1: Bubbles Viewport Behavior Differs

**Symptom:** Scrolling feels different or cursor visibility is inconsistent

**Diagnosis:**
1. Check YOffset values during scrolling (add debug logging)
2. Verify middle-third calculation is correct
3. Compare with custom viewport behavior (archived implementation)

**Rollback Decision:**
- If issue is fixable with < 30 min adjustment → Fix forward
- If fundamental incompatibility → Rollback

**Rollback Steps:**
1. `cp docs/archived/viewport-custom-v1.go internal/editor/viewport.go`
2. `git checkout ui/workspace/model.go ui/workspace/view.go`
3. Remove bubbles dependency from go.mod
4. `go mod tidy`
5. `go test ./...`

### Scenario 2: Performance Regression

**Symptom:** Typing or scrolling feels laggy, frame rate drops

**Diagnosis:**
1. Run benchmark tests (see Pattern 3 above)
2. Profile with `go test -cpuprofile=cpu.out -bench=.`
3. Compare with custom viewport benchmarks (if available)

**Rollback Decision:**
- If <10ms regression → Acceptable, fix forward
- If >10ms regression → Investigate cause
- If >50ms regression → Rollback

**Fix Forward Options:**
- Optimize `SetContent()` calls (cache content, only update on change)
- Batch viewport updates (debounce)
- Profile and optimize hot paths

### Scenario 3: Integration Issues

**Symptom:** Viewport doesn't work with other components (status bar, etc.)

**Diagnosis:**
1. Check integration points (status bar, cursor, buffer)
2. Verify message passing in Update()
3. Test isolation: Does viewport work standalone?

**Rollback Decision:**
- If integration issue is localized → Fix forward
- If cascading failures → Rollback

---

## Troubleshooting

### Issue 1: Cursor Not Visible After Scrolling

**Symptoms:**
- Cursor disappears when navigating
- Cursor is off-screen but should be visible

**Possible Causes:**
1. `ensureCursorVisible()` not called after cursor movement
2. Middle-third calculation incorrect
3. YOffset not updating correctly

**Debug Steps:**
1. Add logging to `ensureCursorVisible()`:
   ```go
   log.Printf("CursorY=%d, ViewportOffset=%d, MiddleTop=%d, MiddleBottom=%d",
       cursorY, m.viewport.YOffset, middleTop, middleBottom)
   ```
2. Verify cursor position with `m.cursor.Position()`
3. Check viewport bounds: `m.viewport.YOffset`, `m.viewport.Height`

**Solution:**
- Ensure `ensureCursorVisible()` called after every cursor move
- Verify third calculation: `third = max(1, height/3)`
- Check for integer division issues (use float, then convert)

### Issue 2: Negative YOffset Values

**Symptoms:**
- Panic: "negative YOffset"
- Content scrolls incorrectly upward

**Possible Causes:**
1. Bounds checking missing in custom cursor logic
2. Integer underflow in offset calculation

**Debug Steps:**
1. Add assertion: `if m.viewport.YOffset < 0 { panic() }`
2. Check offset calculation logic
3. Verify `SetYOffset()` clamps to >= 0

**Solution:**
- Always use `max(0, calculated_offset)` when setting YOffset
- Bubbles viewport should handle this automatically, but check custom logic

### Issue 3: Content Not Updating

**Symptoms:**
- Buffer changes don't appear in viewport
- Old content still visible after edits

**Possible Causes:**
1. `SetContent()` not called after buffer modification
2. `View()` called before `SyncViewportContent()`

**Debug Steps:**
1. Add logging: `log.Printf("Setting content: %d lines", buffer.LineCount())`
2. Verify `SyncViewportContent()` called in Update() and View()
3. Check if buffer reference is stale

**Solution:**
- Call `syncViewportContent()` in `View()` before rendering
- Consider calling `SetContent()` in Update() when buffer changes
- Ensure buffer reference is current (not nil)

### Issue 4: Viewport Doesn't Resize

**Symptoms:**
- Terminal resize doesn't adjust viewport
- Content cut off after resize

**Possible Causes:**
1. `tea.WindowSizeMsg` not handled
2. Viewport dimensions not updated
3. Content not re-synced after resize

**Debug Steps:**
1. Verify WindowSizeMsg handler exists in Update()
2. Check if Width/Height updated: `m.viewport.Width = msg.Width`
3. Verify `syncViewportContent()` called after resize

**Solution:**
```go
case tea.WindowSizeMsg:
    m.viewport.Width = msg.Width
    m.viewport.Height = msg.Height - 1 // Status bar
    m.syncViewportContent()
    m.ensureCursorVisible() // Re-adjust scroll
```

---

## Additional References

### External Documentation
- [Bubbles Viewport GoDoc](https://pkg.go.dev/github.com/charmbracelet/bubbles/viewport)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Charm Examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)

### Internal Documentation
- [`ui-domain.md`](../../learnings/ui-domain.md) - Category 2: Cursor and Viewport Management
- [`editor-domain.md`](../../learnings/editor-domain.md) - Editor Integration Patterns
- [`error-handling.md`](../../learnings/error-handling.md) - Error Handling Patterns
- [`bubble-tea-best-practices.md`](../../../../patterns/bubble-tea-best-practices.md) - Bubble Tea Architecture
- [`opencode-design-system.md`](../../opencode-design-system.md) - UI Styling Guidelines

### Related Tasks
- M6 (Undo/Redo): Consider scroll position in undo state
- M34-M35 (Vim Mode): Ensure j/k work with middle-third scrolling
- M37 (Responsive Layout): Viewport resize integration

---

**Document Version:** 1.0  
**Last Updated:** 2026-01-09  
**Status:** Ready for implementation
