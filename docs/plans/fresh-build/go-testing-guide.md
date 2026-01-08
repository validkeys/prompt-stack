# Effect Testing Guide for PromptStack

**Target**: AI-assisted test development for TUI applications
**Focus**: Bubble Tea testing, user input simulation, Effect patterns
**Optimization**: Context-window efficient, actionable patterns

---

## Core Testing Philosophy

### Effect-First Testing
Test **effects** (observable outcomes), not implementation details.

```go
// ✅ GOOD: Test effect
func TestInsertText(t *testing.T) {
    buf := editor.New("hello")
    buf.Insert(" world")
    if buf.Content() != "hello world" {
        t.Errorf("got %q, want %q", buf.Content(), "hello world")
    }
}

// ❌ BAD: Test implementation
func TestInsertText(t *testing.T) {
    buf := editor.New("hello")
    buf.Insert(" world")
    if buf.cursor != 11 { // implementation detail
        t.Errorf("cursor at %d, want 11", buf.cursor)
    }
}
```

### Test Pyramid
```
        E2E (5%)
       /      \
    Integration (15%)
   /              \
Unit Tests (80%)
```

**Priority**: Unit > Integration > E2E

---

## Bubble Tea Testing Patterns

### Model Testing

Test [`Model.Update()`](archive/code/ui/app/model.go:91) and [`Model.View()`](archive/code/ui/app/model.go:320) independently.

```go
func TestWorkspaceUpdate(t *testing.T) {
    tests := []struct {
        name     string
        initial  workspace.Model
        msg      tea.Msg
        want     workspace.Model
        wantCmd  tea.Cmd
    }{
        {
            name:    "insert character",
            initial: workspace.New(""),
            msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
            want:    workspace.New("a"),
            wantCmd: nil,
        },
        {
            name:    "backspace",
            initial: workspace.New("ab"),
            msg:     tea.KeyMsg{Type: tea.KeyBackspace},
            want:    workspace.New("a"),
            wantCmd: nil,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, cmd := tt.initial.Update(tt.msg)
            if got.Content() != tt.want.Content() {
                t.Errorf("content = %q, want %q", got.Content(), tt.want.Content())
            }
            if (cmd == nil) != (tt.wantCmd == nil) {
                t.Errorf("cmd = %v, want %v", cmd, tt.wantCmd)
            }
        })
    }
}
```

### Command Testing

Test async operations via [`tea.Cmd`](archive/code/ui/app/model.go:91).

```go
func TestSaveCommand(t *testing.T) {
    // Create temp file
    tmp := t.TempDir()
    path := filepath.Join(tmp, "test.md")
    
    // Create save command
    cmd := saveCommand(path, "content")
    
    // Execute command
    msg := cmd()
    
    // Verify message type
    saveMsg, ok := msg.(saveSuccessMsg)
    if !ok {
        t.Fatalf("got %T, want saveSuccessMsg", msg)
    }
    
    // Verify file was saved
    content, err := os.ReadFile(path)
    if err != nil {
        t.Fatal(err)
    }
    if string(content) != "content" {
        t.Errorf("got %q, want %q", string(content), "content")
    }
}
```

### Message Testing

Test message handling in [`Update()`](archive/code/ui/app/model.go:91).

```go
func TestHandleWindowSize(t *testing.T) {
    model := workspace.New("")
    msg := tea.WindowSizeMsg{Width: 80, Height: 24}
    
    updated, _ := model.Update(msg)
    
    if updated.width != 80 {
        t.Errorf("width = %d, want 80", updated.width)
    }
    if updated.height != 24 {
        t.Errorf("height = %d, want 24", updated.height)
    }
}
```

---

## User Input Simulation

### Key Message Construction

```go
// Single key
msg := tea.KeyMsg{Type: tea.KeyEnter}

// Character input
msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}

// Multiple characters
msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h', 'e', 'l', 'l', 'o'}}

// Ctrl key
msg := tea.KeyMsg{Type: tea.KeyCtrlC}

// Alt key
msg := tea.KeyMsg{Type: tea.KeyAltLeft}
```

### Input Sequences

```go
func TestTypingSequence(t *testing.T) {
    model := workspace.New("")
    
    // Type "hello"
    for _, r := range "hello" {
        msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
        model, _ = model.Update(msg)
    }
    
    if model.Content() != "hello" {
        t.Errorf("got %q, want %q", model.Content(), "hello")
    }
}
```

### Complex Interactions

```go
func TestEditWorkflow(t *testing.T) {
    model := workspace.New("hello world")
    
    // Move cursor left 5 times
    for i := 0; i < 5; i++ {
        msg := tea.KeyMsg{Type: tea.KeyLeft}
        model, _ = model.Update(msg)
    }
    
    // Delete "world"
    for i := 0; i < 5; i++ {
        msg := tea.KeyMsg{Type: tea.KeyBackspace}
        model, _ = model.Update(msg)
    }
    
    // Type "there"
    for _, r := range "there" {
        msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
        model, _ = model.Update(msg)
    }
    
    if model.Content() != "hello there" {
        t.Errorf("got %q, want %q", model.Content(), "hello there")
    }
}
```

---

## Domain Testing

### Editor Domain

Test [`internal/editor/`](archive/code/internal/editor/) logic independently.

```go
func TestParsePlaceholders(t *testing.T) {
    tests := []struct {
        input string
        want  []editor.Placeholder
    }{
        {
            input: "Hello {{text:name}}",
            want: []editor.Placeholder{
                {Type: "text", Name: "name"},
            },
        },
        {
            input: "{{list:items}}",
            want: []editor.Placeholder{
                {Type: "list", Name: "items"},
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            got := editor.ParsePlaceholders(tt.input)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Prompt Domain

Test [`internal/prompt/`](archive/code/internal/prompt/) parsing and validation.

```go
func TestExtractKeywords(t *testing.T) {
    content := "Write a function to sort an array of integers"
    keywords := prompt.ExtractKeywords(content)
    
    if keywords["write"] == 0 {
        t.Error("expected 'write' keyword")
    }
    if keywords["function"] == 0 {
        t.Error("expected 'function' keyword")
    }
    if keywords["the"] != 0 { // stop word
        t.Error("stop word should be filtered")
    }
}
```

### AI Domain

Mock [`internal/ai/`](archive/code/internal/ai/) client for testing.

```go
type mockAIClient struct {
    response string
    err      error
}

func (m *mockAIClient) SendMessage(ctx context.Context, req ai.MessageRequest) (*ai.MessageResponse, error) {
    if m.err != nil {
        return nil, m.err
    }
    return &ai.MessageResponse{
        Content: m.response,
    }, nil
}

func TestAISuggestions(t *testing.T) {
    mock := &mockAIClient{
        response: `{"suggestions": [{"title": "Test", "description": "Test suggestion"}]}`,
    }
    
    model := app.NewWithDependencies("", nil, nil, nil, false, mock, nil)
    
    // Trigger AI suggestions
    msg := app.TriggerAISuggestionsMsg{}
    model, cmd := model.Update(msg)
    
    // Execute command
    result := cmd()
    
    // Verify suggestions generated
    suggestionsMsg, ok := result.(app.AISuggestionsGeneratedMsg)
    if !ok {
        t.Fatalf("got %T, want AISuggestionsGeneratedMsg", result)
    }
    
    if len(suggestionsMsg.Suggestions) == 0 {
        t.Error("expected suggestions")
    }
}
```

---

## Integration Testing

### Component Integration

Test interactions between UI components.

```go
func TestBrowserInsertPrompt(t *testing.T) {
    // Setup
    lib := &library.Library{
        Prompts: map[string]*prompt.Prompt{
            "test.md": {
                Title:   "Test Prompt",
                Content: "Hello {{text:name}}",
            },
        },
    }
    
    model := app.NewWithDependencies("", lib, nil, nil, false, nil, nil)
    
    // Show browser
    msg := tea.KeyMsg{Type: tea.KeyCtrlB}
    model, _ = model.Update(msg)
    
    // Select prompt
    msg = tea.KeyMsg{Type: tea.KeyEnter}
    model, _ = model.Update(msg)
    
    // Verify prompt inserted
    if !strings.Contains(model.workspace.GetContent(), "Hello {{text:name}}") {
        t.Error("prompt not inserted into workspace")
    }
}
```

### End-to-End Workflows

```go
func TestCompleteWorkflow(t *testing.T) {
    // Setup
    tmp := t.TempDir()
    lib := library.Load(tmp, nil)
    model := app.NewWithDependencies(tmp, lib, nil, nil, false, nil, nil)
    
    // 1. Type composition
    for _, r := range "Write a function" {
        msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
        model, _ = model.Update(msg)
    }
    
    // 2. Open browser
    msg := tea.KeyMsg{Type: tea.KeyCtrlB}
    model, _ = model.Update(msg)
    
    // 3. Select prompt
    msg = tea.KeyMsg{Type: tea.KeyEnter}
    model, _ = model.Update(msg)
    
    // 4. Save
    msg = tea.KeyMsg{Type: tea.KeyCtrlS}
    model, cmd := model.Update(msg)
    
    // Execute save command
    cmd()
    
    // Verify saved
    files, _ := os.ReadDir(filepath.Join(tmp, ".history"))
    if len(files) == 0 {
        t.Error("composition not saved")
    }
}
```

---

## Test Utilities

### Test Helpers

```go
// testutil/helpers.go
package testutil

import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
)

// KeyPress creates a KeyMsg for a single character
func KeyPress(r rune) tea.Msg {
    return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

// KeySequence creates multiple KeyMsgs
func KeySequence(s string) []tea.Msg {
    var msgs []tea.Msg
    for _, r := range s {
        msgs = append(msgs, KeyPress(r))
    }
    return msgs
}

// ApplyKeys applies a sequence of keys to a model
func ApplyKeys(model tea.Model, msgs []tea.Msg) tea.Model {
    for _, msg := range msgs {
        model, _ = model.Update(msg)
    }
    return model
}

// AssertEqual marks function as test helper and checks equality
func AssertEqual(t *testing.T, got, want interface{}) {
    t.Helper() // Marks this as helper for better error reporting
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

### Mock Factories

```go
// testutil/mocks.go
package testutil

import (
    "github.com/kyledavis/prompt-stack/internal/ai"
    "github.com/kyledavis/prompt-stack/internal/library"
    "github.com/kyledavis/prompt-stack/internal/prompt"
)

// MockLibrary creates a test library with sample prompts
func MockLibrary() *library.Library {
    return &library.Library{
        Prompts: map[string]*prompt.Prompt{
            "test.md": {
                Title:   "Test",
                Content: "Test content",
            },
        },
    }
}

// MockAIClient creates a mock AI client
func MockAIClient(response string, err error) *mockAIClient {
    return &mockAIClient{response: response, err: err}
}
```

### Fixtures

```go
// test/fixtures/prompts/sample.md
---
title: Sample Prompt
description: A sample prompt for testing
tags: test,example
---

This is a sample prompt with {{text:placeholder}}.

// test/fixtures/compositions/simple.md
Write a function to sort an array.
```

---

## Testing Anti-Patterns

### Don't Test View Output

```go
// ❌ BAD: Testing view string
func TestView(t *testing.T) {
    model := workspace.New("hello")
    view := model.View()
    if view != "hello" {
        t.Errorf("view = %q, want %q", view, "hello")
    }
}

// ✅ GOOD: Test model state
func TestContent(t *testing.T) {
    model := workspace.New("hello")
    if model.Content() != "hello" {
        t.Errorf("content = %q, want %q", model.Content(), "hello")
    }
}
```

### Don't Test Private Methods

```go
// ❌ BAD: Testing internal implementation
func TestInternalHelper(t *testing.T) {
    result := internalHelper("input") // unexported
    if result != "output" {
        t.Error("wrong result")
    }
}

// ✅ GOOD: Test public API
func TestPublicMethod(t *testing.T) {
    model := New("input")
    result := model.PublicMethod()
    if result != "output" {
        t.Error("wrong result")
    }
}
```

### Don't Ignore Errors

```go
// ❌ BAD: Ignoring errors
func TestLoad(t *testing.T) {
    lib, _ := library.Load(path) // error ignored
    if lib == nil {
        t.Error("library is nil")
    }
}

// ✅ GOOD: Handling errors
func TestLoad(t *testing.T) {
    lib, err := library.Load(path)
    if err != nil {
        t.Fatalf("failed to load library: %v", err)
    }
    if lib == nil {
        t.Error("library is nil")
    }
}
```

---

## Performance Testing

### Benchmark Tests

```go
func BenchmarkParsePlaceholders(b *testing.B) {
    content := strings.Repeat("{{text:name}} ", 100)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        editor.ParsePlaceholders(content)
    }
}

func BenchmarkUpdate(b *testing.B) {
    model := workspace.New(strings.Repeat("a", 1000))
    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        model.Update(msg)
    }
}
```

### Memory Testing

```go
func TestMemoryUsage(t *testing.T) {
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    
    // Perform operation
    lib := library.Load("large/path", nil)
    
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    
    allocated := m2.Alloc - m1.Alloc
    if allocated > 10<<20 { // 10MB
        t.Errorf("allocated %d bytes, want < 10MB", allocated)
    }
}
```

---

## Advanced Testing Patterns

### Parallel Tests

```go
func TestEditor(t *testing.T) {
    t.Run("Insert", func(t *testing.T) {
        t.Parallel() // Run concurrently with other subtests
        buf := editor.New("")
        buf.Insert("text")
        if buf.Content() != "text" {
            t.Errorf("got %q, want %q", buf.Content(), "text")
        }
    })
    
    t.Run("Delete", func(t *testing.T) {
        t.Parallel()
        buf := editor.New("text")
        buf.Delete(2)
        if buf.Content() != "te" {
            t.Errorf("got %q, want %q", buf.Content(), "te")
        }
    })
}
```

### Cleanup Handlers

```go
func TestWithCleanup(t *testing.T) {
    // Setup
    db := openTestDB()
    t.Cleanup(func() {
        db.Close() // Runs after test completes
    })
    
    // Test operations
    db.Insert("test")
    
    // No need for defer - t.Cleanup handles it
}
```

### Context and Timeouts

```go
func TestWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()
    
    cmd := loadDataCommand(ctx)
    msg := cmd()
    
    if _, ok := msg.(timeoutMsg); !ok {
        t.Error("expected timeout")
    }
}

func TestWithDeadline(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping slow test in short mode")
    }
    
    ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
    defer cancel()
    
    // Long-running operation
}
```

---

## Test Organization

### File Structure

```
internal/editor/
├── buffer.go
├── buffer_test.go       # Unit tests
├── placeholder.go
├── placeholder_test.go # Unit tests
└── editor_test.go      # Integration tests

ui/workspace/
├── model.go
├── model_test.go       # Model tests
├── update.go
└── update_test.go      # Update tests

test/
├── integration/
│   ├── workflow_test.go # E2E tests
│   └── ai_test.go     # AI integration tests
├── fixtures/
│   ├── prompts/
│   └── compositions/
└── testutil/
    ├── helpers.go
    └── mocks.go
```

### Test Naming

```go
// Unit tests
func TestBuffer_Insert(t *testing.T) { ... }
func TestBuffer_Delete(t *testing.T) { ... }

// Integration tests
func TestWorkspace_SaveWorkflow(t *testing.T) { ... }
func TestApp_AISuggestions(t *testing.T) { ... }

// Benchmark tests
func BenchmarkParsePlaceholders(t *testing.B) { ... }
func BenchmarkUpdate(t *testing.B) { ... }
```

---

## Running Tests

### Unit Tests

```bash
# Run all tests
go test ./...

# Run specific package
go test ./internal/editor

# Run with coverage
go test -cover ./...

# Run with race detector
go test -race ./...

# Run parallel tests
go test -parallel 4 ./...

# Skip slow tests
go test -short ./...
```

### Integration Tests

```bash
# Run integration tests
go test ./test/integration

# Run with verbose output
go test -v ./test/integration

# Run specific test
go test -run TestCompleteWorkflow ./test/integration

# Run with timeout
go test -timeout 30s ./test/integration
```

### Benchmarks

```bash
# Run benchmarks
go test -bench=. ./...

# Run with memory stats
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkParsePlaceholders ./internal/editor

# Compare benchmarks
go test -bench=. -count 5 ./... | tee old.txt
# (make changes)
go test -bench=. -count 5 ./... | tee new.txt
benchstat old.txt new.txt
```

---

## Quick Reference

### Key Testing Patterns

| Pattern | Use Case | Example |
|----------|-----------|---------|
| Table-driven | Multiple test cases | [`TestParsePlaceholders`](#domain-testing) |
| Mock interfaces | External dependencies | [`mockAIClient`](#domain-testing) |
| Message testing | Bubble Tea updates | [`TestHandleWindowSize`](#message-testing) |
| Command testing | Async operations | [`TestSaveCommand`](#command-testing) |
| Input simulation | User interactions | [`TestTypingSequence`](#input-sequences) |
| Parallel tests | Concurrent execution | [`TestEditor`](#parallel-tests) |
| Cleanup handlers | Resource management | [`TestWithCleanup`](#cleanup-handlers) |
| Context/timeout | Time-bound operations | [`TestWithTimeout`](#context-and-timeouts) |

### Common Assertions

```go
// Equality
if got != want { t.Errorf("got %v, want %v", got, want) }

// Deep equality
if !reflect.DeepEqual(got, want) { ... }

// Error checking
if err != nil { t.Fatalf("unexpected error: %v", err) }

// String contains
if !strings.Contains(got, want) { ... }

// Length
if len(got) != want { ... }
```

### Test Helpers

```go
// Create temp directory
tmp := t.TempDir()

// Create temp file
path := filepath.Join(tmp, "test.txt")
os.WriteFile(path, []byte("content"), 0644)

// Cleanup (automatic with t.TempDir)
// No manual cleanup needed

// Helper function pattern
func assertEqual(t *testing.T, got, want string) {
    t.Helper() // Better error line numbers
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```

---

## Best Practices Summary

1. **Test effects, not implementation**
2. **Use table-driven tests for multiple cases**
3. **Mock external dependencies**
4. **Test public API, not private methods**
5. **Handle errors explicitly**
6. **Use descriptive test names**
7. **Keep tests focused and independent**
8. **Use `t.Helper()` in test utilities**
9. **Use `t.Cleanup()` for resource cleanup**
10. **Use `t.Parallel()` for concurrent tests**
11. **Use `testing.Short()` to skip slow tests**
12. **Run tests with race detector (`-race`)**
13. **Maintain high test coverage (>80%)**

---

## Related Documentation

### Bubble Tea Testing Best Practices

For comprehensive guidance on testing Bubble Tea TUI applications, see [`bubble-tea-testing-best-practices.md`](./bubble-tea-testing-best-practices.md).

**Key Topics Covered**:
- Effect-first testing philosophy
- Essential test patterns for Bubble Tea
- Message simulation and input sequences
- Test helpers for AI development
- Anti-patterns to avoid
- AI development checklist

**Why This Matters**: The best practices document was created in response to bug #002 (keyboard input not working), which occurred because tests didn't verify input handling effects. Following these practices prevents similar bugs.

**Quick Reference**:
| Pattern | Use Case | Example |
|---------|----------|---------|
| Table-driven | Multiple test cases | `TestUpdateHandlesAllMessages` |
| Input simulation | User workflows | `TestTypingWorkflow` |
| Edge cases | Boundary conditions | `TestEdgeCases` |
| Test helpers | Reduce boilerplate | `TypeText()`, `PressKey()` |

**See Also**:
- [Test Review Analysis](./milestone-implementation-plans/M2/code-reviews/test-review-against-best-practices.md) - Detailed review of existing tests against best practices
- [Test Improvement Implementation Plan](./milestone-implementation-plans/M2/code-reviews/test-improvement-implementation-plan.md) - Step-by-step plan to align tests with best practices

---

**Remember**: Tests are documentation. Make them clear, maintainable, and focused on behavior.