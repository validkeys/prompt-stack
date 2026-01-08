# Milestone 2: Basic TUI Shell - Reference Document (Part 2)

**Milestone Number**: 2  
**Title**: Basic TUI Shell  
**Goal**: Render functional TUI with quit handling

---

## How to Use This Document

**Read this section when:**
- [Before implementing any task] - Understanding style guide patterns
- [When writing code] - Referencing coding standards and patterns
- [When writing tests] - Referencing testing patterns and examples
- [When debugging] - Checking for common anti-patterns

**Key sections:**
- Lines 15-100: Style Guide References - Reference during implementation
- Lines 102-250: Testing Guide References - Consult when writing tests

**Related documents:**
- See [`reference.md`](reference.md) for Architecture Context
- See [`reference-part-3.md`](reference-part-3.md) for Key Learnings and Implementation Notes (Tasks 1-2)
- See [`reference-part-4.md`](reference-part-4.md) for Implementation Notes (Tasks 3-4), Common Patterns, Performance, Testing Checklist

---

## Style Guide References

### Bubble Tea Model Pattern

**Reference**: [`go-style-guide.md`](../../go-style-guide.md) - Section: Type Design (lines 52-95)

**Pattern**: Implement Bubble Tea Model interface with Init(), Update(), View()

```go
// ✅ GOOD: Standard Bubble Tea model
type Model struct {
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

**Common Pitfalls**:
- ❌ Forgetting to return updated model from Update()
- ❌ Forgetting to return command from Update()
- ❌ Not handling tea.WindowSizeMsg
- ❌ Using value receiver when pointer receiver needed

### Constructor Pattern

**Reference**: [`go-style-guide.md`](../../go-style-guide.md) - Section: Type Design (lines 52-95)

**Pattern**: Use New() or NewType() for constructors

```go
// ✅ GOOD: Simple constructor
func New() Model {
    return Model{
        content: "",
        width:   80,
        height:  24,
    }
}

// ✅ GOOD: Constructor with dependencies
func New(logger *zap.Logger, config *config.Config) Model {
    return Model{
        logger: logger,
        config: config,
        content: "",
    }
}

// ❌ BAD: Stuttering name
func NewModel() Model { // Don't repeat type name
    return Model{}
}
```

### Error Handling Pattern

**Reference**: [`go-style-guide.md`](../../go-style-guide.md) - Section: Error Handling (lines 113-201)

**Pattern**: Use fmt.Errorf with %w for error wrapping

```go
// ✅ GOOD: Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create model: %w", err)
}

// ✅ GOOD: Check for specific errors
if errors.Is(err, os.ErrNotExist) {
    // Handle file not found
}

// ❌ BAD: Discarding original error
if err != nil {
    return fmt.Errorf("failed to create model") // Lost original error
}
```

### Package Documentation

**Reference**: [`go-style-guide.md`](../../go-style-guide.md) - Section: Package Organization (lines 18-48)

**Pattern**: Package-level documentation comment

```go
// ✅ GOOD: Package documentation
// Package app provides the root Bubble Tea model for the PromptStack TUI.
//
// The app model coordinates all UI components and handles keyboard input,
// window resizing, and quit functionality.
package app

// ✅ GOOD: Exported function documentation
// New creates a new app model with the given dependencies.
func New(logger *zap.Logger, config *config.Config) Model {
    return Model{
        logger: logger,
        config: config,
    }
}
```

---

## Testing Guide References

### Bubble Tea Model Testing

**Reference**: [`go-testing-guide.md`](../../go-testing-guide.md) - Section: Bubble Tea Testing Patterns (lines 47-143)

**Pattern**: Test Update() and View() independently

```go
func TestModelUpdate(t *testing.T) {
    tests := []struct {
        name     string
        initial  Model
        msg      tea.Msg
        want     Model
        wantCmd  tea.Cmd
    }{
        {
            name:    "handle window size",
            initial: New(),
            msg:     tea.WindowSizeMsg{Width: 80, Height: 24},
            want:    Model{width: 80, height: 24},
            wantCmd: nil,
        },
        {
            name:    "handle quit",
            initial: New(),
            msg:     tea.KeyMsg{Type: tea.KeyCtrlC},
            want:    New(),
            wantCmd: tea.Quit,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, cmd := tt.initial.Update(tt.msg)
            if got.width != tt.want.width {
                t.Errorf("width = %d, want %d", got.width, tt.want.width)
            }
            if (cmd == nil) != (tt.wantCmd == nil) {
                t.Errorf("cmd = %v, want %v", cmd, tt.wantCmd)
            }
        })
    }
}
```

### User Input Simulation

**Reference**: [`go-testing-guide.md`](../../go-testing-guide.md) - Section: User Input Simulation (lines 147-214)

**Pattern**: Construct tea.KeyMsg for keyboard input

```go
// Single key
msg := tea.KeyMsg{Type: tea.KeyEnter}

// Character input
msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}

// Ctrl key
msg := tea.KeyMsg{Type: tea.KeyCtrlC}

// Alt key
msg := tea.KeyMsg{Type: tea.KeyAltLeft}

// Input sequence
func TestTypingSequence(t *testing.T) {
    model := New()
    
    // Type "hello"
    for _, r := range "hello" {
        msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
        model, _ = model.Update(msg)
    }
    
    if model.content != "hello" {
        t.Errorf("got %q, want %q", model.content, "hello")
    }
}
```

### Test Organization

**Reference**: [`go-testing-guide.md`](../../go-testing-guide.md) - Section: Test Organization (lines 678-723)

**Pattern**: Test files alongside implementation files

```
ui/app/
├── model.go
└── model_test.go       # Unit tests

ui/statusbar/
├── model.go
└── model_test.go       # Unit tests

ui/theme/
├── theme.go
└── theme_test.go       # Unit tests
```

### Testing Anti-Patterns

**Reference**: [`go-testing-guide.md`](../../go-testing-guide.md) - Section: Testing Anti-Patterns (lines 490-555)

**Don't Test View Output**:
```go
// ❌ BAD: Testing view string
func TestView(t *testing.T) {
    model := New()
    view := model.View()
    if view != "expected" {
        t.Errorf("view = %q, want %q", view, "expected")
    }
}

// ✅ GOOD: Test model state
func TestContent(t *testing.T) {
    model := New()
    if model.content != "expected" {
        t.Errorf("content = %q, want %q", model.content, "expected")
    }
}
```

**Don't Ignore Errors**:
```go
// ❌ BAD: Ignoring errors
func TestLoad(t *testing.T) {
    model, _ := New() // error ignored
    if model == nil {
        t.Error("model is nil")
    }
}

// ✅ GOOD: Handling errors
func TestLoad(t *testing.T) {
    model, err := New()
    if err != nil {
        t.Fatalf("failed to create model: %v", err)
    }
    if model == nil {
        t.Error("model is nil")
    }
}
```

---

## Document Navigation

This reference document is split into multiple parts to comply with the 600-line limit:

- **Part 1**: [`reference.md`](reference.md) - Architecture Context, Navigation Guide
- **Part 2 (this file)**: Style Guide References, Testing Guide References
- **Part 3**: [`reference-part-3.md`](reference-part-3.md) - Key Learnings References, Implementation Notes (Tasks 1-2)
- **Part 4**: [`reference-part-4.md`](reference-part-4.md) - Implementation Notes (Tasks 3-4), Common Patterns, Performance, Testing Checklist

**Continue reading**: See [`reference-part-3.md`](reference-part-3.md) for Key Learnings References and Implementation Notes (Tasks 1-2).

---

**Last Updated**: 2026-01-07  
**Milestone Group**: Foundation (M1-M6)  
**Testing Guide**: [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md)