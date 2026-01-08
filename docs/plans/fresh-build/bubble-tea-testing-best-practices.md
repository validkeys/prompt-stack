# Bubble Tea Testing Best Practices

**Purpose**: AI-focused guide for testing interactive Bubble Tea TUI systems
**Context**: Prevent bugs like [issue #002](../../../issues/002-tui-input-not-working/report.md) through automated testing

## Core Principle: Test Effects, Not Implementation

**Critical**: Bug #002 occurred because tests verified model creation but not input handling effects.

```go
// ✅ GOOD: Test the effect of input
func TestCharacterInputUpdatesCount(t *testing.T) {
    model := app.New()
    
    // Simulate typing "hello"
    for _, r := range "hello" {
        msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
        model, _ = model.Update(msg)
    }
    
    // Verify effect: character count updated
    if model.GetCharCount() != 5 {
        t.Errorf("got %d, want 5", model.GetCharCount())
    }
}

// ❌ BAD: Test implementation details
func TestCharacterInput(t *testing.T) {
    model := app.New()
    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
    model.Update(msg)
    // No assertion - test passes even if input is ignored!
}
```

## Essential Test Patterns

### 1. Message Simulation

Test all message types your model handles:

```go
func TestUpdateHandlesAllMessages(t *testing.T) {
    tests := []struct {
        name     string
        msg      tea.Msg
        wantChar int
        wantLine int
    }{
        {"character", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}, 1, 0},
        {"enter", tea.KeyMsg{Type: tea.KeyEnter}, 0, 1},
        {"backspace", tea.KeyMsg{Type: tea.KeyBackspace}, -1, 0},
        {"resize", tea.WindowSizeMsg{Width: 100, Height: 50}, 0, 0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            model := app.New()
            model, _ = model.Update(tt.msg)
            // Verify state changed as expected
        })
    }
}
```

### 2. Input Sequences

Test realistic user workflows:

```go
func TestTypingWorkflow(t *testing.T) {
    model := app.New()
    
    // Type "hello", press Enter, type "world"
    for _, r := range "hello" {
        model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
    }
    model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
    for _, r := range "world" {
        model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
    }
    
    // Verify final state
    if model.GetCharCount() != 10 {
        t.Errorf("got %d, want 10", model.GetCharCount())
    }
    if model.GetLineCount() != 1 {
        t.Errorf("got %d, want 1", model.GetLineCount())
    }
}
```

### 3. Edge Cases

Test boundary conditions:

```go
func TestEdgeCases(t *testing.T) {
    tests := []struct {
        name string
        test func(*testing.T)
    }{
        {"empty buffer backspace", func(t *testing.T) {
            model := app.New()
            model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
            if model.GetCharCount() != 0 {
                t.Error("should remain 0")
            }
        }},
        {"unicode input", func(t *testing.T) {
            model := app.New()
            model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'😀'}})
            if model.GetCharCount() != 1 {
                t.Error("emoji should count as 1 character")
            }
        }},
        {"rapid input", func(t *testing.T) {
            model := app.New()
            for i := 0; i < 100; i++ {
                model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
            }
            if model.GetCharCount() != 100 {
                t.Error("should handle rapid input")
            }
        }},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, tt.test)
    }
}
```

## Test Helpers

Create reusable helpers to reduce boilerplate:

```go
// testutil/bubbletea.go
package testutil

import tea "github.com/charmbracelet/bubbletea"

// TypeText simulates typing a string
func TypeText(model tea.Model, text string) tea.Model {
    for _, r := range text {
        msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
        model, _ = model.Update(msg)
    }
    return model
}

// PressKey simulates a single key press
func PressKey(model tea.Model, keyType tea.KeyType) tea.Model {
    msg := tea.KeyMsg{Type: keyType}
    model, _ = model.Update(msg)
    return model
}

// AssertState checks model state
func AssertState(t *testing.T, got, want interface{}) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

## Anti-Patterns

### ❌ Don't Test View Output

```go
// BAD: Fragile, tests implementation
func TestView(t *testing.T) {
    model := app.New()
    view := model.View()
    if view != "expected string" { ... }
}

// GOOD: Test model state
func TestContent(t *testing.T) {
    model := app.New()
    if model.GetContent() != "expected" { ... }
}
```

### ❌ Don't Ignore Commands

```go
// BAD: Commands may contain important side effects
model, _ = model.Update(msg)

// GOOD: Check commands
model, cmd := model.Update(msg)
if cmd != nil {
    // Execute and verify command
    msg := cmd()
    // Assert on message type/content
}
```

### ❌ Don't Test Private Methods

```go
// BAD: Tests implementation details
func TestInternalHelper(t *testing.T) { ... }

// GOOD: Test public API
func TestPublicMethod(t *testing.T) { ... }
```

## Test Coverage Checklist

Required coverage for new features:
- [ ] All [`Update()`](cmd/promptstack/main.go:1) message types
- [ ] Input sequences (typing, editing, navigation)
- [ ] Edge cases (empty state, overflow, unicode)
- [ ] State transitions (initial → intermediate → final)
- [ ] Commands returned from [`Update()`](cmd/promptstack/main.go:1)
- [ ] Error conditions
- [ ] Performance (rapid input, large buffers)

## Pattern Reference

| Pattern | Use Case | Example |
|---------|----------|---------|
| Table-driven | Multiple test cases | `TestUpdateHandlesAllMessages` |
| Input simulation | User workflows | `TestTypingWorkflow` |
| Edge cases | Boundary conditions | `TestEdgeCases` |
| Test helpers | Reduce boilerplate | `TypeText()`, `PressKey()` |

## Resources

- [Bubble Tea Testing Guide](https://github.com/charmbracelet/bubbletea/tree/master/tutorials#testing)
- [Bubble Tea Examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling library

## Key Principle

**Rule**: If a user can do it manually, test it automatically. Always verify messages produce expected state changes, not just that code executes without panicking.