# Milestone 2: Basic TUI Shell - Reference Document (Part 4)

**Milestone Number**: 2  
**Title**: Basic TUI Shell  
**Goal**: Render functional TUI with quit handling

---

## How to Use This Document

**Read this section when:**
- [Before implementing Task 3] - Understanding root app model structure
- [Before implementing Task 4] - Understanding main application integration
- [When debugging] - Checking for common anti-patterns
- [When optimizing] - Referencing performance considerations
- [Before completing tasks] - Verifying testing checklist

**Key sections:**
- Lines 15-250: Implementation Notes - Task-specific code examples and patterns (Tasks 3-4)
- Lines 252-350: Common Patterns and Anti-Patterns
- Lines 352-400: Performance Considerations
- Lines 402-450: Testing Checklist

**Related documents:**
- See [`reference.md`](reference.md) for Architecture Context
- See [`reference-part-2.md`](reference-part-2.md) for Style Guide and Testing Guide references
- See [`reference-part-3.md`](reference-part-3.md) for Key Learnings and Implementation Notes (Tasks 1-2)

---

## Implementation Notes

### Task 3: Create Root App Model

**File**: [`ui/app/model.go`](../../project-structure.md:227)

**Code Example**:
```go
// Package app provides a root Bubble Tea model for the PromptStack TUI.
//
// The app model coordinates all UI components and handles keyboard input,
// window resizing, and quit functionality.
package app

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/kyledavis/prompt-stack/ui/statusbar"
    "github.com/kyledavis/prompt-stack/ui/theme"
)

// Model represents the root application model.
type Model struct {
    statusBar statusbar.Model
    width     int
    height    int
    quitting  bool
}

// New creates a new app model.
func New() Model {
    return Model{
        statusBar: statusbar.New(),
        width:     80,
        height:    24,
        quitting:  false,
    }
}

// Init initializes the app model.
func (m Model) Init() tea.Cmd {
    return nil
}

// Update handles messages for the app.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC:
            m.quitting = true
            return m, tea.Quit
            
        case tea.KeyRunes:
            // Check for 'q' key
            if len(msg.Runes) == 1 && msg.Runes[0] == 'q' {
                m.quitting = true
                return m, tea.Quit
            }
        }
        
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        
        // Update status bar with new width
        var cmd tea.Cmd
        m.statusBar, cmd = m.statusBar.Update(msg)
        return m, cmd
    }
    
    return m, nil
}

// View renders the app.
func (m Model) View() string {
    if m.quitting {
        return ""
    }
    
    // Render status bar at bottom
    statusBar := m.statusBar.View()
    
    // Combine views
    return statusBar
}
```

**Test Example**:
```go
// Package app provides tests for the root app model.
package app

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

func TestModelUpdateCtrlC(t *testing.T) {
    model := New()
    msg := tea.KeyMsg{Type: tea.KeyCtrlC}
    
    updated, cmd := model.Update(msg)
    
    if cmd != tea.Quit {
        t.Errorf("Update() returned %v, want tea.Quit", cmd)
    }
    
    updatedModel, ok := updated.(Model)
    if !ok {
        t.Fatalf("Update() returned non-Model type")
    }
    
    if !updatedModel.quitting {
        t.Error("quitting flag not set")
    }
}

func TestModelUpdateQKey(t *testing.T) {
    model := New()
    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
    
    updated, cmd := model.Update(msg)
    
    if cmd != tea.Quit {
        t.Errorf("Update() returned %v, want tea.Quit", cmd)
    }
    
    updatedModel, ok := updated.(Model)
    if !ok {
        t.Fatalf("Update() returned non-Model type")
    }
    
    if !updatedModel.quitting {
        t.Error("quitting flag not set")
    }
}

func TestModelUpdateCharacterInput(t *testing.T) {
    model := New()
    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
    
    updated, cmd := model.Update(msg)
    
    if cmd != nil {
        t.Errorf("Update() returned non-nil command, want nil")
    }
    
    updatedModel, ok := updated.(Model)
    if !ok {
        t.Fatalf("Update() returned non-Model type")
    }
    
    if updatedModel.quitting {
        t.Error("quitting flag should not be set for character input")
    }
}

func TestModelUpdateWindowSize(t *testing.T) {
    model := New()
    msg := tea.WindowSizeMsg{Width: 100, Height: 30}
    
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
    
    if updatedModel.height != 30 {
        t.Errorf("height = %d, want 30", updatedModel.height)
    }
}

func TestView(t *testing.T) {
    model := New()
    model.width = 80
    model.height = 24
    
    view := model.View()
    
    if view == "" {
        t.Error("View() returned empty string")
    }
}

func TestViewWhenQuitting(t *testing.T) {
    model := New()
    model.quitting = true
    
    view := model.View()
    
    if view != "" {
        t.Errorf("View() returned %q, want empty string when quitting", view)
    }
}

func TestRapidKeyboardInput(t *testing.T) {
    model := New()
    
    // Simulate rapid input
    for i := 0; i < 100; i++ {
        msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
        model, _ = model.Update(msg)
    }
    
    // Should not panic
    if model.quitting {
        t.Error("quitting flag should not be set for character input")
    }
}
```

**Integration Considerations**:
- Depends on Task 1 (Theme System) and Task 2 (Status Bar)
- Integrates with main application in Task 4
- Must handle quit messages correctly (Ctrl+C, 'q')

**Rollback Scenarios**:

### Decision Criteria: Fix Forward vs Rollback

**Fix Forward If** (Preferred):
- Issue is minor (< 30 min fix)
- Only affects current task
- No architectural violations
- Test coverage remains acceptable (>80%)
- No integration breaks
- Single point of failure
- Examples: Missing import, test assertion error, wrong error type, simple refactoring mistakes

**Rollback If** (Required):
- Fundamental design issue
- Multiple cascading failures
- Architecture violations (e.g., circular dependencies, global state)
- Significant test coverage drop (>20%)
- Integration breaks existing functionality
- Time to fix > 60 minutes
- Examples: Circular dependency, breaks 3+ existing tests, violates core architectural principles

### Task 3 Rollback Procedure

**Scenario**: App doesn't start or crashes on initialization

**Fix Forward Steps** (< 30 min):
```bash
# 1. Check Init method returns tea.Cmd
grep "func (m Model) Init() tea.Cmd" internal/ui/app/model.go

# 2. Verify Update method signature
grep "func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)" internal/ui/app/model.go

# 3. Check View method returns string
grep "func (m Model) View() string" internal/ui/app/model.go

# 4. Run app with debug output
go run ./cmd/promptstack/main.go --debug 2>&1 | head -50
```

**Rollback Procedure** (if fix forward fails):
```bash
# 1. Identify affected files
git status

# 2. Revert app package changes
git checkout -- internal/ui/app/

# 3. Verify clean state
git status
go test ./... -v

# 4. Document rollback in checkpoint
# Add note: "Rolled back Task 3 due to [reason]. Time spent: X minutes."
```

**Decision Point**: If root model architecture is fundamentally flawed or violates Bubble Tea patterns, rollback immediately.

---

### Task 4: Integrate TUI with Main Application

**File**: [`cmd/promptstack/main.go`](../../project-structure.md:227)

**Code Example**:
```go
// Package main provides the entry point for the PromptStack CLI.
package main

import (
    "os"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/kyledavis/prompt-stack/internal/config"
    "github.com/kyledavis/prompt-stack/internal/logging"
    "github.com/kyledavis/prompt-stack/ui/app"
    "go.uber.org/zap"
)

func main() {
    // Initialize logger
    logger, err := logging.New()
    if err != nil {
        os.Exit(1)
    }
    defer logger.Sync()
    
    logger.Info("starting promptstack")
    
    // Load config
    cfg, err := config.Load()
    if err != nil {
        logger.Error("failed to load config", zap.Error(err))
        os.Exit(1)
    }
    
    logger.Info("config loaded",
        zap.String("path", cfg.Path),
        zap.String("version", cfg.Version),
    )
    
    // Create app model
    model := app.New()
    
    // Create Bubble Tea program
    p := tea.NewProgram(
        model,
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )
    
    logger.Info("starting tui")
    
    // Run program
    if _, err := p.Run(); err != nil {
        logger.Error("failed to run tui", zap.Error(err))
        os.Exit(1)
    }
    
    logger.Info("tui exited")
}
```

**Test Example**:
```go
// Package main provides tests for the main application.
package main

import (
    "os"
    "testing"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/kyledavis/prompt-stack/ui/app"
)

func TestMainIntegration(t *testing.T) {
    // This is an integration test that verifies the main function
    // can be called without panicking.
    // Note: This test doesn't actually run main() as it would
    // start the TUI and block.
    
    // Instead, we test that components can be created
    model := app.New()
    
    if model == (app.Model{}) {
        t.Error("failed to create app model")
    }
    
    // Test that model can be updated
    msg := tea.KeyMsg{Type: tea.KeyCtrlC}
    updated, cmd := model.Update(msg)
    
    if updated == (app.Model{}) {
        t.Error("Update() returned zero value")
    }
    
    if cmd != tea.Quit {
        t.Errorf("Update() returned %v, want tea.Quit", cmd)
    }
}

func TestProgramCreation(t *testing.T) {
    model := app.New()
    
    // Test that program can be created
    p := tea.NewProgram(
        model,
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )
    
    if p == nil {
        t.Error("failed to create program")
    }
}

func TestEnvironmentVariables(t *testing.T) {
    // Test that required environment variables are set
    // (if any)
    
    // For now, we just verify the test environment is valid
    if os.Getenv("PATH") == "" {
        t.Error("PATH environment variable not set")
    }
}
```

**Integration Considerations**:
- Depends on Task 3 (Root App Model)
- Integrates with config system (M1)
- Integrates with logging system (M1)
- Must handle errors gracefully

**Rollback Scenarios**:

### Task 4 Rollback Procedure

**Scenario**: Build fails or app crashes on startup

**Fix Forward Steps** (< 30 min):
```bash
# 1. Check build errors
go build ./cmd/promptstack

# 2. Verify imports are correct
head -20 cmd/promptstack/main.go | grep import -A 10

# 3. Check main function signature
grep "func main()" cmd/promptstack/main.go

# 4. Verify tea.NewProgram call
grep "tea.NewProgram" cmd/promptstack/main.go

# 5. Run with verbose logging
go run ./cmd/promptstack/main.go 2>&1 | head -50
```

**Rollback Procedure** (if fix forward fails):
```bash
# 1. Identify affected files
git status

# 2. Revert main.go changes
git checkout -- cmd/promptstack/main.go

# 3. Verify clean state
git status
go test ./... -v

# 4. Document rollback in checkpoint
# Add note: "Rolled back Task 4 due to [reason]. Time spent: X minutes."
```

**Decision Point**: If integration requires architectural changes to M1 components or creates circular dependencies, rollback immediately.

### Rollback Documentation Requirements

**Every rollback MUST include in checkpoint**:
1. **Decision Made**: Fix forward or rollback (with justification)
2. **Reasoning**: Why this decision was made
3. **Time Spent**: Total time on analysis and fix attempt
4. **Lessons Learned**: What was learned from the failure
5. **Changes to Approach**: How the next attempt will differ
6. **Root Cause**: What caused the issue (if known)

**Example Rollback Documentation**:
```markdown
## Rollback: Task 3 (Root App Model)

**Decision**: Rollback
**Reasoning**: Circular dependency detected between ui/app and ui/statusbar packages
**Time Spent**: 45 minutes (30 min analysis, 15 min fix attempt)
**Lessons Learned**: Status bar should be initialized in app.New(), not passed as dependency
**Changes to Approach**: Will restructure to have app create statusbar internally
**Root Cause**: Misunderstanding of Bubble Tea component composition pattern
```

### Rollback Verification Checklist

After any rollback, verify:
- [ ] All affected files reverted
- [ ] `git status` shows clean state (or only expected changes)
- [ ] `go build ./...` succeeds
- [ ] `go test ./... -v` passes
- [ ] No race conditions: `go test ./... -race`
- [ ] Coverage remains acceptable: `go test ./... -cover`
- [ ] Rollback documented in checkpoint
- [ ] Lessons learned captured
- [ ] New approach planned before retry

---

## Common Patterns and Anti-Patterns

### Common Patterns

**1. Bubble Tea Model Pattern**:
```go
type Model struct {
    // State fields
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle messages
    return m, nil
}

func (m Model) View() string {
    // Render UI
    return styledContent
}
```

**2. Message Handling Pattern**:
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle keyboard input
        return m, nil
    case tea.WindowSizeMsg:
        // Handle window resize
        m.width = msg.Width
        m.height = msg.Height
        return m, nil
    }
    return m, nil
}
```

**3. Style Helper Pattern**:
```go
func StatusStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color(ForegroundMuted)).
        Background(lipgloss.Color(BackgroundSecondary)).
        Padding(0, 1)
}
```

### Common Anti-Patterns

**1. Forgetting to Return Updated Model**:
```go
// ❌ BAD: Not returning updated model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        // Forgot to return m
    }
    return m, nil
}

// ✅ GOOD: Returning updated model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        return m, nil
    }
    return m, nil
}
```

**2. Not Handling Window Size Messages**:
```go
// ❌ BAD: Not handling window resize
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle keyboard input
        return m, nil
    }
    return m, nil
}

// ✅ GOOD: Handling window resize
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
```

**3. Hard-Coding Styles**:
```go
// ❌ BAD: Hard-coded styles
modalStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#45475a")).
    Background(lipgloss.Color("#1e1e2e"))

// ✅ GOOD: Using theme helpers
modalStyle := theme.ModalStyle()
```

---

## Performance Considerations

### Rendering Performance

**Target**: < 16ms per render (60 FPS)

**Optimizations**:
- Minimize string allocations in View()
- Use lipgloss style caching
- Avoid complex calculations in View()

### Input Handling Performance

**Target**: < 16ms per input event

**Optimizations**:
- Keep Update() logic simple
- Avoid expensive operations in Update()
- Use efficient data structures

### Memory Usage

**Target**: < 10MB for TUI

**Optimizations**:
- Reuse model instances
- Avoid unnecessary allocations
- Clean up resources on quit

---

## Testing Checklist

Before marking a task complete, verify:

### Unit Tests
- [ ] All tests pass
- [ ] Coverage > 80%
- [ ] No race conditions
- [ ] No ignored errors

### Integration Tests
- [ ] Components integrate correctly
- [ ] Messages flow properly
- [ ] State is consistent

### Manual Tests
- [ ] TUI launches without errors
- [ ] Keyboard input works
- [ ] Quit works (Ctrl+C, 'q')
- [ ] Window resize works
- [ ] No visual glitches

---

## Document Navigation

This reference document is split into multiple parts to comply with the 600-line limit:

- **Part 1**: [`reference.md`](reference.md) - Architecture Context, Navigation Guide
- **Part 2**: [`reference-part-2.md`](reference-part-2.md) - Style Guide References, Testing Guide References
- **Part 3**: [`reference-part-3.md`](reference-part-3.md) - Key Learnings References, Implementation Notes (Tasks 1-2)
- **Part 4 (this file)**: Implementation Notes (Tasks 3-4), Common Patterns, Performance, Testing Checklist

**End of Reference Document**

---

**Last Updated**: 2026-01-07  
**Milestone Group**: Foundation (M1-M6)  
**Testing Guide**: [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md)