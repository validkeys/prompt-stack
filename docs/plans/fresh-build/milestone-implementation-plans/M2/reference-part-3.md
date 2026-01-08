# Milestone 2: Basic TUI Shell - Reference Document (Part 3)

**Milestone Number**: 2  
**Title**: Basic TUI Shell  
**Goal**: Render functional TUI with quit handling

---

## How to Use This Document

**Read this section when:**
- [Before implementing Task 1] - Understanding theme system requirements and patterns
- [Before implementing Task 2] - Understanding status bar component patterns
- [When applying learnings] - Referencing patterns from previous attempts
- [When writing code] - Applying key learnings to implementation

**Key sections:**
- Lines 15-180: Key Learnings References - Apply patterns from previous attempts
- Lines 182-450: Implementation Notes - Task-specific code examples and patterns (Tasks 1-2)

**Related documents:**
- See [`reference.md`](reference.md) for Architecture Context
- See [`reference-part-2.md`](reference-part-2.md) for Style Guide and Testing Guide references
- See [`reference-part-4.md`](reference-part-4.md) for Implementation Notes (Tasks 3-4), Common Patterns, Performance, Testing Checklist

---

## Key Learnings References

### Bubble Tea Model Implementation

**Reference**: [`learnings/ui-domain.md`](../learnings/ui-domain.md) - Category 1: Bubble Tea Model Implementation (lines 16-63)

**Learning**: Follow Bubble Tea's model-view-update pattern strictly

**Implementation Pattern**:
```go
type Model struct {
    // State fields
    content string
    width   int
    height  int
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle keyboard input
        return m, nil
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil
    }
    return m, nil
}

func (m Model) View() string {
    // Render UI
    return styledContent
}
```

**Lesson**: Always return the updated model and any commands from Update(). This ensures proper state management and message flow.

### Centralized Theme System

**Reference**: [`learnings/ui-domain.md`](../learnings/ui-domain.md) - Category 6: Centralized Theme System (lines 285-364)

**Learning**: Create a centralized theme package with color constants and style helper functions

**Implementation Pattern**:
```go
// ui/theme/theme.go
package theme

import "github.com/charmbracelet/lipgloss"

// Color Constants
const (
    BackgroundPrimary   = "#1e1e2e"
    BackgroundSecondary = "#181825"
    ForegroundPrimary   = "#cdd6f4"
    ForegroundMuted    = "#a6adc8"
    AccentBlue   = "#89b4fa"
    AccentGreen  = "#a6e3a1"
    AccentYellow = "#f9e2af"
    AccentRed    = "#f38ba8"
    BorderPrimary = "#45475a"
)

// Style Helper Functions
func ModalStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color(BorderPrimary)).
        Padding(1, 2).
        Background(lipgloss.Color(BackgroundPrimary))
}

func StatusStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color(ForegroundMuted)).
        Background(lipgloss.Color(BackgroundSecondary)).
        Padding(0, 1)
}
```

**Benefits**:
- Single source of truth for all colors
- Easy to update entire UI theme
- Consistent color palette across components
- Type-safe constants prevent typos

**Lesson**: Create a centralized theme package with color constants and style helper functions. Replace hard-coded lipgloss styles with theme helpers.

### Error Handling Patterns

**Reference**: [`learnings/go-fundamentals.md`](../learnings/go-fundamentals.md) - Category 7: Error Handling Patterns (lines 209-253)

**Learning**: Use fmt.Errorf with %w for error wrapping

**Implementation Pattern**:
```go
import (
    "errors"
    "fmt"
)

// Wrap errors with context
if err != nil {
    return nil, fmt.Errorf("failed to load config: %w", err)
}

// Check for specific errors
if errors.Is(err, os.ErrNotExist) {
    // Handle file not found
}

// Extract wrapped error
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    // Handle path error
}
```

**Benefits**:
- Preserves original error for unwrapping
- Adds context at each layer
- Enables errors.Is() and errors.As() checks
- Clear error messages for debugging

**Lesson**: Always wrap errors with context using %w. Never discard the original error.

---

## Implementation Notes

### Task 1: Create Theme System

**File**: [`ui/theme/theme.go`](../../project-structure.md:267)

**Code Example**:
```go
// Package theme provides a centralized theme system for the PromptStack TUI.
//
// This package defines the Catppuccin Mocha color palette and provides
// style helper functions for consistent styling across all UI components.
package theme

import "github.com/charmbracelet/lipgloss"

// Color Constants - Catppuccin Mocha Palette
const (
    // Background colors
    BackgroundPrimary   = "#1e1e2e"
    BackgroundSecondary = "#181825"
    
    // Foreground colors
    ForegroundPrimary = "#cdd6f4"
    ForegroundMuted  = "#a6adc8"
    
    // Accent colors
    AccentBlue   = "#89b4fa"
    AccentGreen  = "#a6e3a1"
    AccentYellow = "#f9e2af"
    AccentRed    = "#f38ba8"
    
    // Border colors
    BorderPrimary = "#45475a"
)

// StatusStyle returns a lipgloss style for the status bar.
func StatusStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color(ForegroundMuted)).
        Background(lipgloss.Color(BackgroundSecondary)).
        Padding(0, 1)
}

// ModalStyle returns a lipgloss style for modal overlays.
func ModalStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color(BorderPrimary)).
        Padding(1, 2).
        Background(lipgloss.Color(BackgroundPrimary))
}

// ActivePlaceholderStyle returns a lipgloss style for active placeholders.
func ActivePlaceholderStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Background(lipgloss.Color(AccentYellow)).
        Foreground(lipgloss.Color(BackgroundPrimary)).
        Bold(true)
}
```

**Test Example**:
```go
// Package theme provides tests for the theme system.
package theme

import (
    "testing"
    
    "github.com/charmbracelet/lipgloss"
)

func TestColorConstants(t *testing.T) {
    tests := []struct {
        name  string
        value string
    }{
        {"BackgroundPrimary", BackgroundPrimary},
        {"BackgroundSecondary", BackgroundSecondary},
        {"ForegroundPrimary", ForegroundPrimary},
        {"ForegroundMuted", ForegroundMuted},
        {"AccentBlue", AccentBlue},
        {"AccentGreen", AccentGreen},
        {"AccentYellow", AccentYellow},
        {"AccentRed", AccentRed},
        {"BorderPrimary", BorderPrimary},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.value == "" {
                t.Errorf("color constant is empty")
            }
            if len(tt.value) != 7 { // #RRGGBB format
                t.Errorf("color constant %q has invalid length", tt.value)
            }
            if tt.value[0] != '#' {
                t.Errorf("color constant %q missing # prefix", tt.value)
            }
        })
    }
}

func TestStatusStyle(t *testing.T) {
    style := StatusStyle()
    
    if style == (lipgloss.Style{}) {
        t.Error("StatusStyle returned zero value")
    }
    
    // Test that style renders
    rendered := style.Render("test")
    if rendered == "" {
        t.Error("style rendered empty string")
    }
}

func TestModalStyle(t *testing.T) {
    style := ModalStyle()
    
    if style == (lipgloss.Style{}) {
        t.Error("ModalStyle returned zero value")
    }
    
    // Test that style renders
    rendered := style.Render("test")
    if rendered == "" {
        t.Error("style rendered empty string")
    }
}

func TestActivePlaceholderStyle(t *testing.T) {
    style := ActivePlaceholderStyle()
    
    if style == (lipgloss.Style{}) {
        t.Error("ActivePlaceholderStyle returned zero value")
    }
    
    // Test that style renders
    rendered := style.Render("test")
    if rendered == "" {
        t.Error("style rendered empty string")
    }
}
```

**Integration Considerations**:
- This is a foundation task with no dependencies
- All other UI tasks depend on this task
- Colors must match Catppuccin Mocha specification exactly

**Rollback Scenarios**:
- If color constants are incorrect, update hex values
- If style functions don't compile, check lipgloss API

---

### Task 2: Create Status Bar Component

**File**: [`ui/statusbar/model.go`](../../project-structure.md:247)

**Code Example**:
```go
// Package statusbar provides a status bar component for the PromptStack TUI.
//
// The status bar displays application status information including character count,
// line count, and status messages.
package statusbar

import (
    "fmt"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/kyledavis/prompt-stack/ui/theme"
)

// Model represents the status bar component.
type Model struct {
    charCount int
    lineCount int
    width     int
}

// New creates a new status bar model.
func New() Model {
    return Model{
        charCount: 0,
        lineCount: 0,
        width:     80,
    }
}

// Init initializes the status bar model.
func (m Model) Init() tea.Cmd {
    return nil
}

// Update handles messages for the status bar.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        return m, nil
    }
    return m, nil
}

// View renders the status bar.
func (m Model) View() string {
    statusStyle := theme.StatusStyle().
        Width(m.width).
        Height(1)
    
    statusText := fmt.Sprintf("%d chars, %d lines", m.charCount, m.lineCount)
    return statusStyle.Render(statusText)
}

// SetCharCount sets the character count.
func (m *Model) SetCharCount(count int) {
    m.charCount = count
}

// SetLineCount sets the line count.
func (m *Model) SetLineCount(count int) {
    m.lineCount = count
}
```

**Test Example**:
```go
// Package statusbar provides tests for the status bar component.
package statusbar

import (
    "testing"
    
    tea "github.com/charmbracelet/bubbletea"
)

func TestModelInit(t *testing.T) {
    model := New()
    cmd := model.Init()
    
    if cmd != nil {
        t.Errorf("Init() returned non-nil command, want nil")
    }
}

func TestModelUpdateWindowSize(t *testing.T) {
    model := New()
    msg := tea.WindowSizeMsg{Width: 100, Height: 24}
    
    updated, cmd := model.Update(msg)
    
    if cmd != nil {
        t.Errorf("Update() returned non-nil command, want nil")
    }
    
    updatedModel, ok := updated.(Model)
    if !ok {
        t.Fatalf("Update() returned non-Model type")
    }
    
    if updatedModel.width != 100 {
        t.Errorf("width = %d, want 100", updatedModel.width)
    }
}

func TestModelView(t *testing.T) {
    model := New()
    model.SetCharCount(42)
    model.SetLineCount(3)
    model.width = 80
    
    view := model.View()
    
    if view == "" {
        t.Error("View() returned empty string")
    }
    
    // Check that view contains expected content
    expected := "42 chars, 3 lines"
    if view != expected {
        t.Errorf("View() = %q, want %q", view, expected)
    }
}

func TestSetCharCount(t *testing.T) {
    model := New()
    model.SetCharCount(100)
    
    if model.charCount != 100 {
        t.Errorf("charCount = %d, want 100", model.charCount)
    }
}

func TestSetLineCount(t *testing.T) {
    model := New()
    model.SetLineCount(5)
    
    if model.lineCount != 5 {
        t.Errorf("lineCount = %d, want 5", model.lineCount)
    }
}

func TestViewWithZeroWidth(t *testing.T) {
    model := New()
    model.width = 0
    
    view := model.View()
    
    // Should not panic with zero width
    if view == "" {
        t.Error("View() returned empty string")
    }
}

func TestViewWithLargeWidth(t *testing.T) {
    model := New()
    model.width = 1000
    
    view := model.View()
    
    // Should not panic with large width
    if view == "" {
        t.Error("View() returned empty string")
    }
}
```

**Integration Considerations**:
- Depends on Task 1 (Theme System)
- Will be integrated into root app model in Task 3
- Must handle window resize messages correctly

**Rollback Scenarios**:
- If theme integration fails, check theme package imports
- If window resize doesn't work, verify Update() handles tea.WindowSizeMsg

---

## Document Navigation

This reference document is split into multiple parts to comply with the 600-line limit:

- **Part 1**: [`reference.md`](reference.md) - Architecture Context, Navigation Guide
- **Part 2**: [`reference-part-2.md`](reference-part-2.md) - Style Guide References, Testing Guide References
- **Part 3 (this file)**: Key Learnings References, Implementation Notes (Tasks 1-2)
- **Part 4**: [`reference-part-4.md`](reference-part-4.md) - Implementation Notes (Tasks 3-4), Common Patterns, Performance, Testing Checklist

**Continue reading**: See [`reference-part-4.md`](reference-part-4.md) for Implementation Notes (Tasks 3-4), Common Patterns, Performance, and Testing Checklist.

---

**Last Updated**: 2026-01-07  
**Milestone Group**: Foundation (M1-M6)  
**Testing Guide**: [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md)