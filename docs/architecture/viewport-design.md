# Viewport Architecture Design

**Date:** 2026-01-09
**Version:** 1.0
**Status:** Active

---

## Overview

This document describes the viewport architecture for PromptStack's text editor, including the middle-third scrolling strategy, bubbles viewport integration, and best practices for future development.

### Key Decisions

1. **Bubbles Viewport**: Production-proven `charmbracelet/bubbles/viewport` component
2. **Middle-Third Scrolling**: Cursor maintained in middle third for optimal context visibility
3. **Separation of Concerns**: Viewport handles scrolling, workspace handles cursor visibility

---

## Architecture

### Component Hierarchy

```
main.go
  → ui/app (TUI application)
    → ui/workspace (workspace model)
      → bubbles/viewport (viewport.Model)
      → internal/editor/buffer (text buffer)
      → internal/editor/cursor (cursor position)
```

### Package Structure

```
ui/workspace/
├── model.go                 # Workspace model with viewport instance
├── view.go                  # Rendering logic with viewport.View()
├── viewport_cursor_test.go   # Cursor visibility tests
└── model_test.go             # General workspace tests
```

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
- **Top third**: Lines 0-3 (lines 0 to YOffset + third)
- **Middle third**: Lines 3-8 (lines YOffset + third to YOffset + Height - third)
- **Bottom third**: Lines 8-12 (lines YOffset + Height - third to YOffset + Height)

### Scrolling Rules

1. **Cursor in middle third**: No scrolling needed ✅
2. **Cursor in top third**: Scroll viewport up (decrease YOffset)
3. **Cursor in bottom third**: Scroll viewport down (increase YOffset)

### Algorithm

```go
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

    // Ensure we don't scroll past document end
    totalLines := buffer.LineCount()
    maxOffset := max(0, totalLines - viewportHeight)
    newOffset = min(newOffset, maxOffset)

    viewport.SetYOffset(newOffset)

ELSE:
    // Cursor in middle third, no scroll needed
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

---

## API Reference

### Bubbles Viewport API

#### Initialization

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

#### Content Management

```go
// Set content (triggers internal layout recalculation)
vp.SetContent(content string)

// Get content
content := vp.View() // Returns visible portion
```

#### Scroll Position

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

#### Viewport Dimensions

```go
// Get viewport dimensions
width := vp.Width
height := vp.Height

// Update dimensions (typically on resize)
vp.Width = newWidth
vp.Height = newHeight
```

#### Scrolling Queries

```go
// Check scroll state
atTop := vp.AtTop()          // true if at document start
atBottom := vp.AtBottom()    // true if at document end
canScrollUp := !vp.AtTop()
canScrollDown := !vp.AtBottom()

// Get total height
totalHeight := vp.TotalLineCount() // Total lines in content
```

### Workspace Viewport Integration

#### Sync Content

```go
// ui/workspace/model.go
func (m Model) syncViewportLines() Model {
    newModel := m
    newModel.viewport.SetContent(newModel.buffer.Content())
    return newModel
}
```

#### Adjust Viewport for Cursor

```go
// ui/workspace/model.go
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
```

---

## Integration Points

### M6: Undo/Redo System

**Status**: Future consideration (M6 not yet implemented)

**Integration Point**: Scroll position should be part of undo state

**Details**:
- When implementing M6, consider adding `viewport.YOffset` to undo state
- Undoing a text edit SHOULD restore scroll position to where user was
- Redoing SHOULD restore scroll position forward

**Testing Strategy**:
- Add integration test in M6: Undo edit → verify viewport position restored
- Test: Redo edit → verify viewport position restored forward

**How to Test**:
1. Make edit at bottom of document (scroll down)
2. Undo edit
3. Verify viewport scrolls back to previous position

**Implementation Note**:
```go
type UndoState struct {
    Content       string
    CursorX       int
    CursorY       int
    ViewportOffset int  // Add this for M6
}
```

### M34-M35: Vim Mode

**Status**: Future consideration (M34-M35 not yet implemented)

**Integration Point**: j/k keys MUST work consistently in normal mode

**Details**:
- j/k keys in vim normal mode should trigger same cursor movement as arrow keys
- Cursor visibility logic MUST work with vim navigation commands
- Middle-third scrolling MUST apply to vim navigation

**Testing Strategy**:
- Add integration test in M34-M35: Press j/k 50 times → verify smooth scrolling
- Test: Vim motions (10j, 5k) → verify correct cursor/viewport coordination

**How to Test**:
1. Enter vim normal mode
2. Press j 100 times rapidly
3. Verify cursor stays visible
4. Verify middle-third scrolling applies

**Implementation Note**:
Vim mode will call the same `adjustViewport()` method that arrow keys use, so no special handling needed.

### M37: Responsive Layout

**Status**: Future consideration (M37 not yet implemented)

**Integration Point**: Viewport MUST resize correctly on terminal resize

**Details**:
- Terminal resize SHOULD preserve scroll position when possible
- Viewport MUST handle resize from large to small gracefully
- Status bar MUST update after resize

**Testing Strategy**:
- Add integration test in M37: Resize terminal → verify viewport adjusts
- Test: Resize to minimum size (80x24) → verify no crashes
- Test: Resize to very small (20x10) → verify graceful degradation

**How to Test**:
1. Open document with 100 lines
2. Scroll to line 50
3. Resize terminal from 120x40 to 80x24
4. Verify viewport adjusts and cursor stays visible
5. Verify scroll position preserved proportionally

**Implementation Note**:
```go
// ui/workspace/model.go - Update() method
case tea.WindowSizeMsg:
    newModel := m
    newModel.width = msg.Width
    newModel.height = msg.Height
    newModel.viewport.Width = msg.Width
    newModel.viewport.Height = msg.Height - 1 // Reserve status bar
    newModel = newModel.syncViewportLines()
    newModel = newModel.adjustViewport() // Re-adjust scroll
    return newModel, nil
```

---

## Best Practices

### When to Call `syncViewportLines()`

**Always call after**:
- Buffer content changes (insert, delete, paste)
- File load operations
- Terminal resize (to recalculate layout)

**Don't call after**:
- Cursor movement only (no content change)
- Scroll-only operations

### When to Call `adjustViewport()`

**Always call after**:
- Cursor movement (arrow keys, vim j/k)
- Text insertion/deletion that moves cursor
- File navigation (goto line)
- Terminal resize

**Don't call when**:
- Content-only changes without cursor movement (rare)

### Performance Considerations

1. **Avoid redundant `SetContent()` calls**: Only call when buffer actually changes
2. **Batch viewport updates**: If multiple operations, sync once at the end
3. **Debounce resize events**: If rapidly resizing, use debounce to avoid excessive recalculations

### Bounds Checking

**Always verify**:
- YOffset is never negative (handled by bubbles, but double-check in cursor logic)
- YOffset never exceeds `max(0, totalLines - viewportHeight)`
- Cursor position is within document bounds

---

## Troubleshooting

### Issue 1: Cursor Not Visible After Scrolling

**Symptoms**:
- Cursor disappears when navigating
- Cursor is off-screen but should be visible

**Possible Causes**:
1. `adjustViewport()` not called after cursor movement
2. Middle-third calculation incorrect
3. YOffset not updating correctly

**Debug Steps**:
1. Add logging to `adjustViewport()`:
   ```go
   log.Printf("CursorY=%d, ViewportOffset=%d, MiddleTop=%d, MiddleBottom=%d",
       cursorY, m.viewport.YOffset, middleTop, middleBottom)
   ```
2. Verify cursor position with `m.buffer.CursorPosition()`
3. Check viewport bounds: `m.viewport.YOffset`, `m.viewport.Height`

**Solution**:
- Ensure `adjustViewport()` called after every cursor move
- Verify third calculation: `third = max(1, height/3)`
- Check for integer division issues (use float, then convert)

### Issue 2: Negative YOffset Values

**Symptoms**:
- Panic: "negative YOffset"
- Content scrolls incorrectly upward

**Possible Causes**:
1. Bounds checking missing in custom cursor logic
2. Integer underflow in offset calculation

**Debug Steps**:
1. Add assertion: `if m.viewport.YOffset < 0 { panic() }`
2. Check offset calculation logic
3. Verify `SetYOffset()` clamps to >= 0

**Solution**:
- Always use `max(0, calculated_offset)` when setting YOffset
- Bubbles viewport should handle this automatically, but check custom logic

### Issue 3: Content Not Updating

**Symptoms**:
- Buffer changes don't appear in viewport
- Old content still visible after edits

**Possible Causes**:
1. `SetContent()` not called after buffer modification
2. `View()` called before `SyncViewportContent()`

**Debug Steps**:
1. Add logging: `log.Printf("Setting content: %d lines", buffer.LineCount())`
2. Verify `SyncViewportContent()` called in Update() and View()
3. Check if buffer reference is stale

**Solution**:
- Call `syncViewportLines()` in `View()` before rendering
- Consider calling `SetContent()` in Update() when buffer changes
- Ensure buffer reference is current (not nil)

### Issue 4: Viewport Doesn't Resize

**Symptoms**:
- Terminal resize doesn't adjust viewport
- Content cut off after resize

**Possible Causes**:
1. `tea.WindowSizeMsg` not handled
2. Viewport dimensions not updated
3. Content not re-synced after resize

**Debug Steps**:
1. Verify WindowSizeMsg handler exists in Update()
2. Check if Width/Height updated: `m.viewport.Width = msg.Width`
3. Verify `syncViewportContent()` called after resize

**Solution**:
```go
case tea.WindowSizeMsg:
    m.viewport.Width = msg.Width
    m.viewport.Height = msg.Height - 1 // Status bar
    m.syncViewportLines()
    m.adjustViewport() // Re-adjust scroll
```

### Issue 5: Performance Issues with Large Documents

**Symptoms**:
- Typing or scrolling feels laggy
- Frame rate drops

**Possible Causes**:
1. Excessive `SetContent()` calls
2. Inefficient cursor visibility calculations

**Debug Steps**:
1. Run benchmark tests
2. Profile with `go test -cpuprofile=cpu.out -bench=.`
3. Check for hot loops

**Solution**:
- Optimize `SetContent()` calls (cache content, only update on change)
- Batch viewport updates (debounce)
- Profile and optimize hot paths

---

## Migration History

### v1.0 (2026-01-09): Bubbles Viewport Integration

**Changes**:
- Replaced custom viewport implementation with `charmbracelet/bubbles/viewport` v0.21.0
- Implemented middle-third scrolling strategy in `adjustViewport()`
- Updated view rendering to use `viewport.View()` for content display
- Added comprehensive test suite with >95% coverage for cursor visibility

**Benefits**:
- Production-proven implementation (used by Glow, Charm's editor)
- Automatic bounds checking (prevents negative YOffset)
- Better edge case handling (small viewports, empty content, resize)
- Performance optimization for large documents (>10,000 lines)
- Active maintenance and community support

**Archived**: `docs/archived/viewport-custom-v1.go`

---

## References

### Internal Documentation
- [`ui-domain.md`](../plans/fresh-build/learnings/ui-domain.md) - Category 2: Cursor and Viewport Management
- [`editor-domain.md`](../plans/fresh-build/learnings/editor-domain.md) - Editor Integration Patterns
- [`task-list.md`](../implementation-plans/m4-viewport-bubbles-migration/task-list.md) - Task 5 Implementation Plan
- [`reference.md`](../implementation-plans/m4-viewport-bubbles-migration/reference.md) - Detailed Reference Documentation

### External Documentation
- [Bubbles Viewport GoDoc](https://pkg.go.dev/github.com/charmbracelet/bubbles/viewport)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Charm Examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)

### Related Design Documents
- [`opencode-design-system.md`](../plans/fresh-build/opencode-design-system.md) - UI Styling Guidelines
- [`bubble-tea-best-practices.md`](../patterns/bubble-tea-best-practices.md) - Bubble Tea Architecture

---

**Maintainer**: PromptStack Development Team
**Last Updated**: 2026-01-09
**Next Review**: After M6, M34-M35, M37 implementation
