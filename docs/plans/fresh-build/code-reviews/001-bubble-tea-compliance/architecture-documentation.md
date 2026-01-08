# PromptStack Architecture Documentation

**Created:** 2026-01-08  
**Version:** 2.0 (Post-Refactoring)  
**Status:** Current

## Overview

PromptStack is a Bubble Tea-based TUI application for AI-assisted prompt composition. This document describes the current architecture after the Bubble Tea compliance refactoring, which achieved 95% architectural compliance, 85.8% test coverage, and eliminated all state mutation violations.

## Architecture Principles

### 1. Bubble Tea Compliance

The application strictly follows [Bubble Tea best practices](../../../bubble-tea-best-practices.md):

- **Immutable Update Pattern:** All `Update()` methods return new model instances, never mutating receivers
- **Model-View-Update Separation:** Clear separation between state, rendering, and logic
- **Component Decomposition:** Complex models broken into focused, testable components
- **Message-Based Communication:** Async operations use Bubble Tea's command system

### 2. Component-Based Design

The editor functionality is decomposed into four focused components:

| Component | Location | Lines | Responsibility |
|-----------|----------|-------|----------------|
| Cursor | `internal/editor/cursor.go` | 189 | Position tracking, navigation |
| Viewport | `internal/editor/viewport.go` | 123 | Scroll management, visible lines |
| Placeholder | `internal/editor/placeholder.go` | 385 | Template variable parsing, editing |
| FileIO | `internal/editor/fileio.go` | 87 | File save/load operations |

### 3. Test-Driven Development

- **Test Coverage:** 85.8% for workspace model (up from 52.5%)
- **Test Utilities:** Centralized helpers in `internal/testutil/bubbletea.go`
- **Effect Tests:** Verify state changes from user input
- **Command Tests:** Verify async operations return correct commands
- **Message Coverage:** Tests for all message types

## Component Architecture

### Editor Components (`internal/editor/`)

#### Cursor Component

**Responsibility:** Track cursor position and handle movement

```go
type Cursor struct {
    x    int // Column position
    y    int // Line position
    line int // Current line index (for tracking)
}
```

**Key Methods:**
- `Update(msg tea.Msg) Cursor` - Handle keyboard messages
- `MoveUp() / MoveDown() / MoveLeft() / MoveRight()` - Movement operations
- `Position() (x, y int)` - Get current position
- `MoveToPosition(pos int, content string) Cursor` - Jump to absolute position
- `AdjustToLineLength(content string) Cursor` - Clamp to line boundaries

**Testing:** 100% coverage for all movement methods and edge cases

#### Viewport Component

**Responsibility:** Manage scrolling and visible content area

```go
type Viewport struct {
    topLine    int // First visible line
    height     int // Viewport height in lines
    totalLines int // Total lines in document
}
```

**Key Methods:**
- `Update(msg tea.Msg) Viewport` - Handle resize and scroll messages
- `ScrollTo(line int) Viewport` - Move viewport to specific line
- `ScrollUp() / ScrollDown()` - Incremental scrolling
- `EnsureVisible(line int) Viewport` - Ensure cursor stays visible
- `VisibleLines() (start, end int)` - Get visible range

**Testing:** 100% coverage including edge cases (empty doc, single line, etc.)

#### Placeholder Manager

**Responsibility:** Parse, validate, and edit template variables

```go
type Placeholder struct {
    Type         string   // "text" or "list"
    Name         string   // placeholder name
    StartPos     int      // position in content
    EndPos       int      // position in content
    CurrentValue string   // current filled value (for text)
    ListValues   []string // current filled values (for list)
    IsValid      bool     // whether syntax is valid
    IsActive     bool     // whether currently selected
}

type Manager struct {
    placeholders []Placeholder
    activeIndex  int
    isEditing    bool
    editValue    string
}
```

**Key Methods:**
- `Update(msg tea.Msg) Manager` - Handle all placeholder-related messages
- `parsePlaceholders(content string) Manager` - Extract placeholders from content
- `nextPlaceholder() / previousPlaceholder()` - Navigate between placeholders
- `editPlaceholder(value string) Manager` - Update edit value
- `saveEdit() Manager` - Commit edit to placeholder

**Custom Messages:**
- `ParsePlaceholdersMsg` - Trigger placeholder parsing
- `ActivatePlaceholderMsg` - Activate placeholder by index
- `EditPlaceholderMsg` - Update edit value
- `ExitEditModeMsg` - Exit edit mode without saving

**Testing:** 95% coverage including regex parsing, validation, and navigation

#### File Manager

**Responsibility:** Handle file I/O operations with modification tracking

```go
type FileManager struct {
    filepath string
    modified bool
}
```

**Key Methods:**
- `Load() (string, error)` - Read file content
- `Save(content string) error` - Write content to file
- `MarkModified() FileManager` - Mark as having unsaved changes
- `ClearModified() FileManager` - Clear modification flag
- `SetPath(path string) FileManager` - Update file path

**Testing:** 100% coverage including directory creation and error handling

### Workspace Model (`ui/workspace/model.go`)

**Responsibility:** Integrate editor components and provide main editor interface

**Structure:**
```go
type Model struct {
    content      string
    cursor       editor.Cursor
    viewport     editor.Viewport
    placeholders editor.Manager
    fileManager  editor.FileManager
    width        int
    height       int
    isDirty      bool
    saveStatus   string
}
```

**Update Pattern:**
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Update child components
    m.cursor = m.cursor.Update(msg)
    m.viewport = m.viewport.Update(msg)
    m.placeholders = m.placeholders.Update(msg)
    
    // Handle workspace-specific messages
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKey(msg)
    case tea.WindowSizeMsg:
        return m.handleResize(msg)
    }
    
    return m, nil
}
```

**Testing:** 85.8% coverage (47 tests, all passing)

### App Model (`ui/app/model.go`)

**Responsibility:** Top-level model coordinating workspace, palette, and other UI components

**Structure:**
```go
type Model struct {
    workspace    workspace.Model
    palette      palette.Model
    activePanel  string
    width        int
    height       int
    quitting     bool
}
```

**Update Pattern:** Immutable updates returning new model instances

**Testing:** Comprehensive tests for all message types and state transitions

## Immutable Update Pattern

All models follow the immutable update pattern:

```go
// BEFORE (incorrect - violates Bubble Tea):
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.field = value  // Direct mutation
    return m, nil
}

// AFTER (correct - immutable):
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    newModel := m
    newModel.field = value
    return newModel, nil
}
```

**Benefits:**
- No state mutation violations
- Easier to reason about state changes
- Better testability (pure functions)
- Aligns with Bubble Tea's Elm Architecture

## Test Infrastructure

### Test Utilities (`internal/testutil/bubbletea.go`)

Centralized helpers for testing Bubble Tea models:

```go
// TypeText simulates typing text character by character
func TypeText(model tea.Model, text string) tea.Model

// PressKey simulates pressing a key
func PressKey(model tea.Model, keyType tea.KeyType) tea.Model

// AssertState compares model states
func AssertState(t *testing.T, got, want tea.Model)

// SimulateWorkflow runs a sequence of operations
func SimulateWorkflow(model tea.Model, steps []Step) tea.Model
```

### Test Coverage

| Component | Coverage | Tests | Status |
|-----------|----------|-------|--------|
| workspace | 85.8% | 47 | ✅ All passing |
| cursor | 100% | 15 | ✅ All passing |
| viewport | 100% | 12 | ✅ All passing |
| placeholder | 95% | 28 | ✅ All passing |
| fileio | 100% | 10 | ✅ All passing |
| app | ~80% | 20 | ✅ All passing |

### Test Types

1. **Effect Tests:** Verify state changes from user input
   - Example: Typing text updates content
   - Example: Cursor movement changes position

2. **Command Tests:** Verify async operations return correct commands
   - Example: Quit command returns tea.Quit
   - Example: AI suggestion returns generation command

3. **Message Coverage:** Tests for all message types
   - tea.KeyMsg
   - tea.WindowSizeMsg
   - tea.MouseMsg
   - Custom messages (ParsePlaceholdersMsg, etc.)

4. **Edge Cases:** Boundary conditions and error cases
   - Empty content
   - Unicode characters
   - Very long lines
   - File I/O errors

## Success Metrics

| Metric | Before | After | Target | Status |
|--------|--------|-------|--------|--------|
| Test Coverage | 52.5% | 85.8% | 80% | ✅ Exceeded |
| State Mutations | 15+ | 0 | 0 | ✅ Achieved |
| Workspace Lines | 1292 | 387 | <400 | ✅ Achieved |
| Effect Tests | 0% | 100% | 100% | ✅ Achieved |
| Arch Compliance | 70% | 95% | 95% | ✅ Achieved |

## Message Flow

### User Input Flow

```
User Types Character
    ↓
tea.KeyMsg (KeyRunes)
    ↓
Workspace.Update()
    ↓
1. Update cursor component
2. Insert character into content
3. Update viewport (ensure cursor visible)
4. Parse placeholders
5. Mark as dirty
6. Schedule auto-save
    ↓
Return (newModel, autoSaveCmd)
```

### Placeholder Editing Flow

```
User Presses Tab
    ↓
tea.KeyMsg (KeyTab)
    ↓
PlaceholderManager.Update()
    ↓
1. Save current edit if editing
2. Move to next placeholder
3. Activate new placeholder
4. Enter edit mode for text placeholders
    ↓
Return (newManager, nil)
    ↓
Workspace renders active placeholder with highlight
```

### Auto-Save Flow

```
User stops typing (750ms delay)
    ↓
autoSaveMsg triggers
    ↓
Workspace.Update()
    ↓
1. Set saveStatus = "saving"
2. Return save command
    ↓
File save executes in background
    ↓
saveSuccessMsg or saveErrorMsg
    ↓
Update saveStatus
    ↓
Clear status after 2 seconds
```

## Integration Points

### Component Communication

Components communicate through:

1. **Direct Updates:** Parent calls child's Update() method
   ```go
   m.cursor = m.cursor.Update(msg)
   ```

2. **Custom Messages:** Bubble Tea message system for async operations
   ```go
   return m, tea.Tick(750*time.Millisecond, func(t time.Time) tea.Msg {
       return autoSaveMsg{}
   })
   ```

3. **State Access:** Parent reads child state via getter methods
   ```go
   x, y := m.cursor.Position()
   ```

### Data Flow

```
User Input → Messages → Model Updates → State Changes → View Rendering
                ↑                    ↓
            Commands          Components
```

## Performance Considerations

### String Operations

- Use `strings.Builder` for complex view rendering
- Avoid repeated string concatenation in loops
- Cache line splitting where possible

### Placeholder Parsing

- Re-parse only on content changes
- Regex with position tracking for efficiency
- Maintain active placeholder across re-parses

### Viewport Scrolling

- "Middle third" strategy for smooth scrolling
- Ensure cursor visible only when needed
- Minimal calculations per frame

## Compliance Checklist

✅ All Update() methods return new model instances  
✅ No direct mutations in Update() methods  
✅ All models implement Init(), Update(), View()  
✅ Component decomposition (<400 lines per component)  
✅ 80%+ test coverage achieved  
✅ Effect tests for all state changes  
✅ Command tests for async operations  
✅ Message coverage for all message types  
✅ No race conditions detected  
✅ All integration tests passing  

## Related Documentation

- [Bubble Tea Best Practices](../../../bubble-tea-best-practices.md) - Framework guidelines
- [Architecture Patterns Key Learnings](../learnings/architecture-patterns.md) - Patterns used
- [Editor Domain Key Learnings](../learnings/editor-domain.md) - Editor-specific patterns
- [UI Domain Key Learnings](../learnings/ui-domain.md) - UI-specific patterns
- [Go Testing Guide](../go-testing-guide.md) - Testing strategies
- [Go Style Guide](../go-style-guide.md) - Code standards

## Future Enhancements

### Planned Improvements

1. **Performance Optimization** (Phase 3)
   - Profile view rendering for large documents
   - Optimize string operations in Update
   - Implement efficient diff algorithms
   - Add rendering caching

2. **Testing Infrastructure** (Phase 3)
   - CI pipeline with coverage reporting
   - Benchmark tests
   - Golden file tests for views
   - Property-based testing

3. **Developer Experience** (Phase 3)
   - Model template generator
   - Linting rules for Bubble Tea
   - Comprehensive test helpers
   - Interactive documentation

---

**Last Updated:** 2026-01-08  
**Audit:** 001-bubble-tea-compliance  
**Status:** Complete
