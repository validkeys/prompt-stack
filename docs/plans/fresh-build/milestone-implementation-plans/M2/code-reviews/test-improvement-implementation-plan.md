# Test Improvement Implementation Plan

**Created**: 2026-01-08  
**Based on**: [`test-review-against-best-practices.md`](./test-review-against-best-practices.md)  
**Target**: [`cmd/promptstack/main_test.go`](../../../../cmd/promptstack/main_test.go)  
**Goal**: Align test suite with [`bubble-tea-testing-best-practices.md`](../../bubble-tea-testing-best-practices.md)

---

## Executive Summary

This implementation plan addresses critical gaps in the test suite that allowed bug #002 (keyboard input not working) to slip through. The primary focus is on adding **effect verification** to ensure tests verify that input produces expected state changes, not just that code runs without panicking.

**Timeline**: 2-3 days  
**Priority**: High  
**Impact**: Prevents similar bugs, improves test reliability

---

## Phase 1: High Priority Fixes (Day 1)

### 1.1 Add Effect Verification to All Input Tests

**Objective**: Ensure all input tests verify the actual state changes produced by input handling.

**Tasks**:

#### Task 1.1.1: Update TestTUIHandlesKeyboardInput
- **File**: [`cmd/promptstack/main_test.go`](../../../../cmd/promptstack/main_test.go:169)
- **Current Issue**: Tests that input doesn't crash but doesn't verify character was added
- **Action**: Add assertions to verify content and character count
- **Expected Outcome**:
  ```go
  func TestTUIHandlesKeyboardInput(t *testing.T) {
      appModel := app.New()
      msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
      newModel, cmd := appModel.Update(msg)
      
      if cmd != nil {
          t.Error("Character input should not return a command")
      }
      
      updatedModel := newModel.(app.Model)
      
      // ✅ NEW: Verify effect
      if updatedModel.GetContent() != "a" {
          t.Errorf("Expected content 'a', got '%s'", updatedModel.GetContent())
      }
      if updatedModel.GetCharCount() != 1 {
          t.Errorf("Expected 1 character, got %d", updatedModel.GetCharCount())
      }
      if updatedModel.IsQuitting() {
          t.Error("Character input should not quit")
      }
  }
  ```

#### Task 1.1.2: Add Tests for Missing Input Types
- **File**: [`cmd/promptstack/main_test.go`](../../../../cmd/promptstack/main_test.go)
- **New Tests to Add**:
  - `TestTUIHandlesBackspace()` - Verify backspace removes characters
  - `TestTUIHandlesEnter()` - Verify Enter creates new lines
  - `TestTUIHandlesArrowKeys()` - Verify arrow keys move cursor
  - `TestTUIHandlesTab()` - Verify Tab behavior
  - `TestTUIHandlesEscape()` - Verify Esc behavior

**Implementation Template**:
```go
func TestTUIHandlesBackspace(t *testing.T) {
    model := app.New()
    
    // Type "hello"
    for _, r := range "hello" {
        msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
        model, _ = model.Update(msg)
    }
    
    // Press backspace
    msg := tea.KeyMsg{Type: tea.KeyBackspace}
    model, _ = model.Update(msg)
    
    // Verify effect: last character removed
    if model.GetContent() != "hell" {
        t.Errorf("Expected 'hell', got '%s'", model.GetContent())
    }
    if model.GetCharCount() != 4 {
        t.Errorf("Expected 4 characters, got %d", model.GetCharCount())
    }
}
```

**Acceptance Criteria**:
- [ ] All existing input tests verify state changes
- [ ] New tests added for backspace, enter, arrow keys, tab, escape
- [ ] All tests pass
- [ ] Test coverage for input handling increases by at least 20%

---

### 1.2 Remove View Output Tests

**Objective**: Eliminate fragile view output tests in favor of model state tests.

**Tasks**:

#### Task 1.2.1: Refactor TestTUIIntegration
- **File**: [`cmd/promptstack/main_test.go`](../../../../cmd/promptstack/main_test.go:318)
- **Current Issue**: Tests view content (lines 350-357)
- **Action**: Replace view assertions with model state assertions
- **Expected Outcome**:
  ```go
  // ❌ REMOVE:
  view := appModel.View()
  if view == "" {
      t.Error("View should render content")
  }
  if !strings.Contains(view, "PromptStack TUI") {
      t.Error("View should contain 'PromptStack TUI'")
  }
  
  // ✅ ADD:
  if appModel.GetTitle() != "PromptStack TUI" {
      t.Errorf("Expected title 'PromptStack TUI', got '%s'", appModel.GetTitle())
  }
  if !appModel.IsInitialized() {
      t.Error("Model should be initialized")
  }
  ```

#### Task 1.2.2: Audit All Tests for View Assertions
- **Action**: Search for all `.View()` calls in test files
- **Action**: Replace with appropriate model state assertions
- **Files to Check**:
  - [`cmd/promptstack/main_test.go`](../../../../cmd/promptstack/main_test.go)
  - Any other test files in the project

**Acceptance Criteria**:
- [ ] No tests assert on view output strings
- [ ] All view tests replaced with model state tests
- [ ] All tests pass after refactoring

---

### 1.3 Add Input Sequence Tests

**Objective**: Test realistic user workflows with multiple input operations.

**Tasks**:

#### Task 1.3.1: Create TestTypingWorkflow
- **File**: [`cmd/promptstack/main_test.go`](../../../../cmd/promptstack/main_test.go)
- **Action**: Add test for typing, editing, and navigation workflow
- **Expected Outcome**:
  ```go
  func TestTypingWorkflow(t *testing.T) {
      model := app.New()
      
      // Type "hello"
      for _, r := range "hello" {
          msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
          model, _ = model.Update(msg)
      }
      
      // Verify intermediate state
      if model.GetContent() != "hello" {
          t.Errorf("Expected 'hello', got '%s'", model.GetContent())
      }
      
      // Press Enter
      msg := tea.KeyMsg{Type: tea.KeyEnter}
      model, _ = model.Update(msg)
      
      // Type "world"
      for _, r := range "world" {
          msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
          model, _ = model.Update(msg)
      }
      
      // Verify final state
      if model.GetContent() != "hello\nworld" {
          t.Errorf("Expected 'hello\\nworld', got '%s'", model.GetContent())
      }
      if model.GetLineCount() != 2 {
          t.Errorf("Expected 2 lines, got %d", model.GetLineCount())
      }
  }
  ```

#### Task 1.3.2: Create TestEditWorkflow
- **Action**: Add test for typing, backspacing, and retyping
- **Expected Outcome**:
  ```go
  func TestEditWorkflow(t *testing.T) {
      model := app.New()
      
      // Type "hello world"
      model = testutil.TypeText(model, "hello world")
      
      // Move cursor left 5 times
      for i := 0; i < 5; i++ {
          model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
      }
      
      // Delete "world"
      for i := 0; i < 5; i++ {
          model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
      }
      
      // Type "there"
      model = testutil.TypeText(model, "there")
      
      // Verify final content
      if model.GetContent() != "hello there" {
          t.Errorf("Expected 'hello there', got '%s'", model.GetContent())
      }
  }
  ```

**Acceptance Criteria**:
- [ ] At least 3 input sequence tests added
- [ ] Tests cover typing, editing, and navigation workflows
- [ ] All tests verify state changes at each step
- [ ] All tests pass

---

## Phase 2: Medium Priority Improvements (Day 2)

### 2.1 Create Test Helpers

**Objective**: Reduce boilerplate and improve test maintainability.

**Tasks**:

#### Task 2.1.1: Create testutil Package
- **New File**: `testutil/bubbletea.go`
- **Action**: Create reusable test helpers
- **Expected Outcome**:
  ```go
  // testutil/bubbletea.go
  package testutil
  
  import (
      "testing"
      tea "github.com/charmbracelet/bubbletea"
  )
  
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
  
  // AssertContent checks model content
  func AssertContent(t *testing.T, model tea.Model, want string) {
      t.Helper()
      if m, ok := model.(interface{ GetContent() string }); ok {
          if m.GetContent() != want {
              t.Errorf("got %q, want %q", m.GetContent(), want)
          }
      } else {
          t.Error("model does not implement GetContent()")
      }
  }
  
  // AssertCharCount checks character count
  func AssertCharCount(t *testing.T, model tea.Model, want int) {
      t.Helper()
      if m, ok := model.(interface{ GetCharCount() int }); ok {
          if m.GetCharCount() != want {
              t.Errorf("got %d, want %d", m.GetCharCount(), want)
          }
      } else {
          t.Error("model does not implement GetCharCount()")
      }
  }
  ```

#### Task 2.1.2: Refactor Existing Tests to Use Helpers
- **Action**: Update tests in [`cmd/promptstack/main_test.go`](../../../../cmd/promptstack/main_test.go) to use new helpers
- **Example Refactoring**:
  ```go
  // Before:
  func TestTUIHandlesKeyboardInput(t *testing.T) {
      appModel := app.New()
      for _, r := range "hello" {
          msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
          appModel, _ = appModel.Update(msg)
      }
      if appModel.GetContent() != "hello" {
          t.Errorf("got %q, want 'hello'", appModel.GetContent())
      }
  }
  
  // After:
  func TestTUIHandlesKeyboardInput(t *testing.T) {
      model := app.New()
      model = testutil.TypeText(model, "hello")
      testutil.AssertContent(t, model, "hello")
      testutil.AssertCharCount(t, model, 5)
  }
  ```

**Acceptance Criteria**:
- [ ] testutil package created with at least 5 helper functions
- [ ] All existing tests refactored to use helpers
- [ ] Test code reduced by at least 30%
- [ ] All tests pass after refactoring

---

### 2.2 Add State Transition Tests

**Objective**: Verify that models transition through expected states during user interaction.

**Tasks**:

#### Task 2.2.1: Create TestStateTransitions
- **File**: [`cmd/promptstack/main_test.go`](../../../../cmd/promptstack/main_test.go)
- **Action**: Add test for mode transitions
- **Expected Outcome**:
  ```go
  func TestStateTransitions(t *testing.T) {
      model := app.New()
      
      // Initial state
      if model.GetMode() != "normal" {
          t.Errorf("Expected initial mode 'normal', got '%s'", model.GetMode())
      }
      
      // Transition to insert mode
      model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
      if model.GetMode() != "insert" {
          t.Errorf("Expected mode 'insert', got '%s'", model.GetMode())
      }
      
      // Transition back to normal
      model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
      if model.GetMode() != "normal" {
          t.Errorf("Expected mode 'normal', got '%s'", model.GetMode())
      }
  }
  ```

#### Task 2.2.2: Create TestLifecycleStates
- **Action**: Add test for application lifecycle states
- **Expected Outcome**:
  ```go
  func TestLifecycleStates(t *testing.T) {
      model := app.New()
      
      // Initial state
      if !model.IsInitialized() {
          t.Error("Model should be initialized")
      }
      
      // Active state
      model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
      if !model.IsActive() {
          t.Error("Model should be active")
      }
      
      // Quitting state
      model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
      if !model.IsQuitting() {
          t.Error("Model should be quitting")
      }
  }
  ```

**Acceptance Criteria**:
- [ ] At least 2 state transition tests added
- [ ] Tests cover mode transitions and lifecycle states
- [ ] All tests verify state changes at each transition
- [ ] All tests pass

---

### 2.3 Add Unicode and Special Character Tests

**Objective**: Ensure the application handles unicode and special characters correctly.

**Tasks**:

#### Task 2.3.1: Create TestUnicodeInput
- **File**: [`cmd/promptstack/main_test.go`](../../../../cmd/promptstack/main_test.go)
- **Action**: Add test for emoji and non-ASCII characters
- **Expected Outcome**:
  ```go
  func TestUnicodeInput(t *testing.T) {
      model := app.New()
      
      // Test emoji
      msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'😀'}}
      model, _ = model.Update(msg)
      
      if model.GetCharCount() != 1 {
          t.Errorf("Emoji should count as 1 character, got %d", model.GetCharCount())
      }
      
      // Test multi-byte character
      msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'世'}}
      model, _ = model.Update(msg)
      
      if model.GetCharCount() != 2 {
          t.Errorf("Expected 2 characters, got %d", model.GetCharCount())
      }
  }
  ```

#### Task 2.3.2: Create TestSpecialKeyCombinations
- **Action**: Add test for Ctrl+ and Alt+ key combinations
- **Expected Outcome**:
  ```go
  func TestSpecialKeyCombinations(t *testing.T) {
      model := app.New()
      
      // Test Ctrl+S (save)
      model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
      if cmd == nil {
          t.Error("Ctrl+S should return a command")
      }
      
      // Test Ctrl+C (quit)
      model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
      if !model.IsQuitting() {
          t.Error("Ctrl+C should quit")
      }
  }
  ```

**Acceptance Criteria**:
- [ ] At least 2 unicode/special character tests added
- [ ] Tests cover emoji, multi-byte characters, and key combinations
- [ ] All tests verify correct handling
- [ ] All tests pass

---

## Phase 3: Low Priority Enhancements (Day 3)

### 3.1 Improve Command Verification

**Objective**: Ensure commands returned from Update() are properly tested.

**Tasks**:

#### Task 3.1.1: Audit All Update Calls
- **Action**: Find all `model.Update(msg)` calls that ignore the command
- **Action**: Add command verification where appropriate
- **Example**:
  ```go
  // Before:
  newModel, _ := appModel.Update(msg)
  
  // After:
  newModel, cmd := appModel.Update(msg)
  if cmd != nil {
      // Execute command
      resultMsg := cmd()
      // Verify command result
      if _, ok := resultMsg.(expectedType); !ok {
          t.Errorf("Expected expectedType, got %T", resultMsg)
      }
  }
  ```

#### Task 3.1.2: Create TestCommandChains
- **Action**: Add test for sequences of commands
- **Expected Outcome**:
  ```go
  func TestCommandChains(t *testing.T) {
      model := app.New()
      
      // Trigger save command
      model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
      
      // Execute command
      msg := cmd()
      
      // Verify command returns expected message
      saveMsg, ok := msg.(saveSuccessMsg)
      if !ok {
          t.Fatalf("Expected saveSuccessMsg, got %T", msg)
      }
      
      // Verify model updated with command result
      model, _ = model.Update(saveMsg)
      if !model.IsSaved() {
          t.Error("Model should be saved")
      }
  }
  ```

**Acceptance Criteria**:
- [ ] All Update calls verify commands where appropriate
- [ ] At least 1 command chain test added
- [ ] All tests pass

---

### 3.2 Add Performance Benchmarks

**Objective**: Ensure the application performs well under various conditions.

**Tasks**:

#### Task 3.2.1: Create Benchmark Tests
- **New File**: `cmd/promptstack/main_bench_test.go`
- **Action**: Add benchmarks for common operations
- **Expected Outcome**:
  ```go
  func BenchmarkTyping(b *testing.B) {
      model := app.New()
      msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
      
      b.ResetTimer()
      for i := 0; i < b.N; i++ {
          model, _ = model.Update(msg)
      }
  }
  
  func BenchmarkLargeBuffer(b *testing.B) {
      model := app.New()
      model = testutil.TypeText(model, strings.Repeat("a", 10000))
      
      msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}}
      
      b.ResetTimer()
      for i := 0; i < b.N; i++ {
          model, _ = model.Update(msg)
      }
  }
  
  func BenchmarkRendering(b *testing.B) {
      model := app.New()
      model = testutil.TypeText(model, strings.Repeat("a", 1000))
      
      b.ResetTimer()
      for i := 0; i < b.N; i++ {
          _ = model.View()
      }
  }
  ```

#### Task 3.2.2: Establish Performance Baselines
- **Action**: Run benchmarks and record baseline results
- **Action**: Document expected performance thresholds
- **Expected Outcome**:
  ```
  BenchmarkTyping-8              1000000    1234 ns/op
  BenchmarkLargeBuffer-8         100000    12345 ns/op
  BenchmarkRendering-8           100000     5678 ns/op
  ```

**Acceptance Criteria**:
- [ ] At least 3 benchmark tests added
- [ ] Baseline performance documented
- [ ] Benchmarks run successfully

---

## Testing Strategy

### Test Execution Plan

1. **Before Implementation**:
   ```bash
   # Run existing tests to establish baseline
   go test ./cmd/promptstack -v -cover
   ```

2. **During Implementation**:
   ```bash
   # Run tests after each phase
   go test ./cmd/promptstack -v -cover
   go test ./testutil -v -cover
   ```

3. **After Implementation**:
   ```bash
   # Run all tests with coverage
   go test ./... -v -cover
   
   # Run with race detector
   go test ./... -race
   
   # Run benchmarks
   go test ./... -bench=. -benchmem
   ```

### Coverage Goals

- **Current Coverage**: ~60% (estimated)
- **Target Coverage**: >80%
- **Critical Path Coverage**: 100%

### Success Metrics

- [ ] All existing tests pass after refactoring
- [ ] Test coverage increases by at least 20%
- [ ] No tests assert on view output
- [ ] All input tests verify state changes
- [ ] Test code reduced by at least 30% (due to helpers)
- [ ] All new tests pass
- [ ] Benchmarks establish performance baselines

---

## Risk Mitigation

### Potential Issues

1. **Breaking Changes**: Refactoring tests may break existing functionality
   - **Mitigation**: Run tests after each change, use git to track progress

2. **Missing Model Methods**: Tests may require model methods that don't exist
   - **Mitigation**: Add necessary methods to model interface, document changes

3. **Performance Regression**: New tests may slow down test suite
   - **Mitigation**: Use table-driven tests, parallel tests where appropriate

4. **Incomplete Coverage**: Some edge cases may be missed
   - **Mitigation**: Use coverage reports, manual code review

### Rollback Plan

If implementation causes issues:
1. Revert to previous commit
2. Implement changes incrementally
3. Test each change before proceeding

---

## Documentation Updates

### Required Documentation Changes

1. **Update [`go-testing-guide.md`](../../go-testing-guide.md)**:
   - Add reference to [`bubble-tea-testing-best-practices.md`](../../bubble-tea-testing-best-practices.md)
   - Document new test helpers
   - Add examples of effect verification

2. **Update Test File Headers**:
   - Add comments explaining test purpose
   - Reference best practices document

3. **Create Test Checklist**:
   - Document test requirements for new features
   - Include effect verification checklist

---

## Timeline Summary

| Phase | Tasks | Duration | Priority |
|-------|-------|----------|----------|
| Phase 1 | Effect verification, remove view tests, input sequences | Day 1 | High |
| Phase 2 | Test helpers, state transitions, unicode tests | Day 2 | Medium |
| Phase 3 | Command verification, benchmarks | Day 3 | Low |

**Total Duration**: 2-3 days

---

## Next Steps

1. **Immediate**:
   - Review this plan with team
   - Get approval to proceed
   - Set up branch for implementation

2. **Day 1**:
   - Implement Phase 1 tasks
   - Run tests after each task
   - Commit changes frequently

3. **Day 2**:
   - Implement Phase 2 tasks
   - Refactor existing tests
   - Create test helpers

4. **Day 3**:
   - Implement Phase 3 tasks
   - Run full test suite
   - Document results

5. **Completion**:
   - Update documentation
   - Create pull request
   - Present results to team

---

## References

- **Test Review**: [`test-review-against-best-practices.md`](./test-review-against-best-practices.md)
- **Best Practices**: [`bubble-tea-testing-best-practices.md`](../../bubble-tea-testing-best-practices.md)
- **Bug Report**: [`issues/002-tui-input-not-working/report.md`](../../../../issues/002-tui-input-not-working/report.md)
- **Test File**: [`cmd/promptstack/main_test.go`](../../../../cmd/promptstack/main_test.go)
- **Official Guide**: https://github.com/charmbracelet/bubbletea/tree/master/tutorials#testing

---

**Remember**: The goal is to prevent bugs like #002 by ensuring tests verify that input produces expected state changes, not just that code runs without panicking.