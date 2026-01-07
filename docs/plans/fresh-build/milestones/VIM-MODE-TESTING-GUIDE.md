# Vim Mode Testing Guide (Milestones 34-35)

**Milestone Group**: Vim Mode  
**Milestones**: M34-M35  
**Focus**: State machine, keybindings

## Overview

This guide provides comprehensive testing strategies for the Vim Mode milestone group, which implements a context-aware Vim keybinding system across multiple components. Testing focuses on ensuring accurate state machine transitions, correct keybinding behavior, and seamless integration with existing editor functionality.

## Integration Tests

### Test Suite: Vim Mode Integration

**Location**: `internal/vim/integration_test.go`

#### Test 1: Vim State Machine
```go
func TestVimStateMachine(t *testing.T) {
    // Test state machine transitions
    // 1. Test normal mode
    // 2. Test insert mode
    // 3. Test visual mode
    // 4. Test command mode
    // 5. Test mode transitions
    // 6. Test state persistence
}
```

**Acceptance Criteria**:
- [ ] Normal mode works correctly
- [ ] Insert mode works correctly
- [ ] Visual mode works correctly
- [ ] Command mode works correctly
- [ ] Mode transitions work correctly
- [ ] State persists across operations
- [ ] No invalid states
- [ ] Transitions complete in <10ms

#### Test 2: Vim Keybindings in Editor
```go
func TestVimKeybindingsInEditor(t *testing.T) {
    // Test Vim keybindings in editor context
    // 1. Test movement keys (h, j, k, l)
    // 2. Test editing keys (i, a, o, dd)
    // 3. Test visual mode (v, V)
    // 4. Test copy/paste (y, p)
    // 5. Test undo/redo (u, Ctrl+r)
}
```

**Acceptance Criteria**:
- [ ] Movement keys work correctly
- [ ] Editing keys work correctly
- [ ] Visual mode works correctly
- [ ] Copy/paste works correctly
- [ ] Undo/redo works correctly
- [ ] Keybindings don't conflict with system
- [ ] Response time <50ms

#### Test 3: Vim Keybindings in Browser
```go
func TestVimKeybindingsInBrowser(t *testing.T) {
    // Test Vim keybindings in browser context
    // 1. Test navigation (j, k)
    // 2. Test selection (Enter)
    // 3. Test search (/)
    // 4. Test exit (Esc, q)
}
```

**Acceptance Criteria**:
- [ ] Navigation works correctly
- [ ] Selection works correctly
- [ ] Search works correctly
- [ ] Exit works correctly
- [ ] Context-specific behavior
- [ ] No conflicts with editor keybindings
- [ ] Response time <50ms

#### Test 4: Vim Keybindings in Palette
```go
func TestVimKeybindingsInPalette(t *testing.T) {
    // Test Vim keybindings in palette context
    // 1. Test navigation (j, k, Ctrl+j, Ctrl+k)
    // 2. Test selection (Enter)
    // 3. Test search (/)
    // 4. Test exit (Esc, Ctrl+c)
}
```

**Acceptance Criteria**:
- [ ] Navigation works correctly
- [ ] Selection works correctly
- [ ] Search works correctly
- [ ] Exit works correctly
- [ ] Context-specific behavior
- [ ] No conflicts with other contexts
- [ ] Response time <50ms

#### Test 5: Context-Aware Keybindings
```go
func TestContextAwareKeybindings(t *testing.T) {
    // Test context-aware keybinding behavior
    // 1. Test keybinding in editor
    // 2. Test same keybinding in browser
    // 3. Test same keybinding in palette
    // 4. Verify context-specific behavior
    // 5. Test context switching
}
```

**Acceptance Criteria**:
- [ ] Keybindings work in editor
- [ ] Keybindings work in browser
- [ ] Keybindings work in palette
- [ ] Context-specific behavior correct
- [ ] Context switching works
- [ ] No conflicts between contexts
- [ ] State preserved during context switch

#### Test 6: Vim Mode Toggle
```go
func TestVimModeToggle(t *testing.T) {
    // Test Vim mode enable/disable
    // 1. Enable Vim mode
    // 2. Verify Vim keybindings work
    // 3. Disable Vim mode
    // 4. Verify default keybindings work
    // 5. Toggle multiple times
}
```

**Acceptance Criteria**:
- [ ] Vim mode enables correctly
- [ ] Vim keybindings work when enabled
- [ ] Vim mode disables correctly
- [ ] Default keybindings work when disabled
- [ ] Toggle works multiple times
- [ ] No state corruption
- [ ] Settings persist

## End-to-End Scenarios

### Scenario 1: Basic Vim Editing Workflow

**Description**: Test basic Vim editing workflow from start to finish.

**Steps**:
1. User enables Vim mode
2. User opens file in editor
3. User is in normal mode
4. User navigates with h, j, k, l
5. User enters insert mode with i
6. User types content
7. User exits insert mode with Esc
8. User deletes line with dd
9. User undoes with u
10. User saves with :w

**Expected Results**:
- [ ] Vim mode enabled
- [ ] File opens in normal mode
- [ ] Navigation works
- [ ] Insert mode works
- [ ] Content typed
- [ ] Exit insert mode works
- [ ] Delete line works
- [ ] Undo works
- [ ] Save works
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test basic Vim editing workflow
promptstack edit test.md
# Enable Vim mode
# Navigate with h, j, k, l
# Enter insert mode with i
# Type content
# Exit with Esc
# Delete line with dd
# Undo with u
# Save with :w
# Verify result
```

### Scenario 2: Visual Mode Selection

**Description**: Test visual mode selection and operations.

**Steps**:
1. User opens file with content
2. User enters visual mode with v
3. User selects text with movement keys
4. User copies selection with y
5. User moves cursor
6. User pastes with p
7. User enters visual line mode with V
8. User selects lines
9. User deletes selection with d

**Expected Results**:
- [ ] Visual mode entered
- [ ] Text selected
- [ ] Copy works
- [ ] Cursor moved
- [ ] Paste works
- [ ] Visual line mode entered
- [ ] Lines selected
- [ ] Delete works
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test visual mode selection
promptstack edit test.md
# Enter visual mode with v
# Select text
# Copy with y
# Move cursor
# Paste with p
# Enter visual line mode with V
# Select lines
# Delete with d
# Verify result
```

### Scenario 3: Context Switching with Vim

**Description**: Test Vim keybindings across different contexts.

**Steps**:
1. User is in editor with Vim mode
2. User opens command palette
3. User navigates palette with j, k
4. User selects command with Enter
5. User opens library browser
6. User navigates browser with j, k
7. User selects prompt with Enter
8. User returns to editor
9. User continues editing with Vim keybindings

**Expected Results**:
- [ ] Editor Vim keybindings work
- [ ] Palette Vim keybindings work
- [ ] Command selected
- [ ] Browser Vim keybindings work
- [ ] Prompt selected
- [ ] Editor Vim keybindings still work
- [ ] No conflicts
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test context switching with Vim
promptstack edit test.md
# Use Vim keybindings in editor
# Open command palette
# Navigate with j, k
# Select with Enter
# Open library browser
# Navigate with j, k
# Select with Enter
# Return to editor
# Continue editing
# Verify no conflicts
```

### Scenario 4: Complex Vim Operations

**Description**: Test complex Vim operations and combinations.

**Steps**:
1. User opens file with content
2. User performs 5dd (delete 5 lines)
3. User performs 3yy (yank 3 lines)
4. User performs p (paste)
5. User performs /search (search)
6. User performs n (next match)
7. User performs cw (change word)
8. User performs . (repeat)
9. User performs :%s/old/new/g (substitute)

**Expected Results**:
- [ ] 5 lines deleted
- [ ] 3 lines yanked
- [ ] Paste works
- [ ] Search works
- [ ] Next match works
- [ ] Change word works
- [ ] Repeat works
- [ ] Substitute works
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test complex Vim operations
promptstack edit test.md
# Delete 5 lines: 5dd
# Yank 3 lines: 3yy
# Paste: p
# Search: /test
# Next match: n
# Change word: cw
# Repeat: .
# Substitute: :%s/old/new/g
# Verify result
```

### Scenario 5: Vim Mode Toggle and Persistence

**Description**: Test Vim mode toggle and settings persistence.

**Steps**:
1. User starts with Vim mode disabled
2. User uses default keybindings
3. User enables Vim mode
4. User uses Vim keybindings
5. User saves file
6. User exits PromptStack
7. User restarts PromptStack
8. User verifies Vim mode still enabled
9. User disables Vim mode
10. User verifies default keybindings work

**Expected Results**:
- [ ] Default keybindings work initially
- [ ] Vim mode enables
- [ ] Vim keybindings work
- [ ] File saves
- [ ] Settings persist
- [ ] Vim mode enabled on restart
- [ ] Vim mode disables
- [ ] Default keybindings work
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test Vim mode toggle and persistence
promptstack edit test.md
# Use default keybindings
# Enable Vim mode
# Use Vim keybindings
# Save file
# Exit
# Restart PromptStack
# Verify Vim mode enabled
# Disable Vim mode
# Verify default keybindings
```

## Performance Benchmarks

### Benchmark 1: State Machine Transitions

**Test**: Measure performance of state machine transitions

```go
func BenchmarkStateMachineTransitions(b *testing.B) {
    stateMachine := NewVimStateMachine()
    for i := 0; i < b.N; i++ {
        stateMachine.Transition("normal", "insert")
        stateMachine.Transition("insert", "normal")
    }
}
```

**Thresholds**:
- [ ] Single transition: <1ms
- [ ] 100 transitions: <10ms
- [ ] 1000 transitions: <100ms
- [ ] No state corruption

### Benchmark 2: Keybinding Processing

**Test**: Measure performance of keybinding processing

```go
func BenchmarkKeybindingProcessing(b *testing.B) {
    processor := NewKeybindingProcessor()
    for i := 0; i < b.N; i++ {
        processor.ProcessKey("j", "editor")
    }
}
```

**Thresholds**:
- [ ] Single key: <1ms
- [ ] 100 keys: <10ms
- [ ] 1000 keys: <100ms
- [ ] No lag

### Benchmark 3: Context Switching

**Test**: Measure performance of context switching

```go
func BenchmarkContextSwitching(b *testing.B) {
    vim := NewVimMode()
    for i := 0; i < b.N; i++ {
        vim.SwitchContext("editor", "browser")
        vim.SwitchContext("browser", "editor")
    }
}
```

**Thresholds**:
- [ ] Single switch: <10ms
- [ ] 100 switches: <100ms
- [ ] 1000 switches: <1s
- [ ] No state loss

### Benchmark 4: Visual Mode Operations

**Test**: Measure performance of visual mode operations

```go
func BenchmarkVisualModeOperations(b *testing.B) {
    editor := NewEditor()
    for i := 0; i < b.N; i++ {
        editor.EnterVisualMode()
        editor.SelectText(100)
        editor.ExitVisualMode()
    }
}
```

**Thresholds**:
- [ ] Enter/exit: <10ms
- [ ] Selection: <10ms
- [ ] 100 operations: <100ms
- [ ] No UI lag

### Benchmark 5: Complex Vim Commands

**Test**: Measure performance of complex Vim commands

```go
func BenchmarkComplexVimCommands(b *testing.B) {
    editor := NewEditor()
    for i := 0; i < b.N; i++ {
        editor.ExecuteCommand("5dd")
        editor.ExecuteCommand("3yy")
        editor.ExecuteCommand("p")
    }
}
```

**Thresholds**:
- [ ] Single command: <10ms
- [ ] 10 commands: <100ms
- [ ] 100 commands: <1s
- [ ] No lag

## Test Execution

### Running Integration Tests

```bash
# Run all Vim mode integration tests
go test ./internal/vim -v -tags=integration

# Run specific test
go test ./internal/vim -v -run TestVimStateMachine

# Run with coverage
go test ./internal/vim -cover -coverprofile=coverage.out
```

### Running End-to-End Tests

```bash
# Run all E2E scenarios
./scripts/test/vim-mode-e2e.sh

# Run specific scenario
./scripts/test/vim-mode-e2e.sh scenario1

# Run with performance monitoring
./scripts/test/vim-mode-e2e.sh --perf
```

### Running Benchmarks

```bash
# Run all benchmarks
go test ./internal/vim -bench=. -benchmem

# Run specific benchmark
go test ./internal/vim -bench=BenchmarkStateMachineTransitions

# Run with CPU profiling
go test ./internal/vim -bench=. -cpuprofile=cpu.prof
```

## Test Data

### Sample Keybindings

**Location**: `test/data/vim/keybindings/`

- `editor.yaml` - Editor keybindings
- `browser.yaml` - Browser keybindings
- `palette.yaml` - Palette keybindings
- `conflicts.yaml` - Conflicting keybindings (for testing)

### Sample Scenarios

**Location**: `test/data/vim/scenarios/`

- `basic-editing.md` - Basic editing scenario
- `visual-mode.md` - Visual mode scenario
- `complex-operations.md` - Complex operations scenario
- `context-switching.md` - Context switching scenario

## Success Criteria

### Integration Tests
- [ ] All integration tests pass
- [ ] Code coverage >80% for Vim mode components
- [ ] No memory leaks detected
- [ ] No race conditions detected

### End-to-End Scenarios
- [ ] All scenarios complete successfully
- [ ] State machine works correctly
- [ ] Keybindings work correctly in all contexts
- [ ] Context switching works correctly
- [ ] Performance meets thresholds

### Performance Benchmarks
- [ ] All benchmarks meet thresholds
- [ ] No performance regression from baseline
- [ ] Memory usage remains bounded
- [ ] CPU usage is reasonable

## Known Issues and Limitations

### Current Limitations
- Some advanced Vim features not implemented (macros, registers)
- Complex keybinding combinations may have edge cases
- Context switching may have latency in some scenarios
- Visual mode may have limitations with very large selections

### Future Improvements
- Add more Vim features (macros, registers, marks)
- Implement custom keybinding configuration
- Add Vim command mode with ex commands
- Optimize context switching
- Add Vim tutorial and help

## References

- [Enhanced Test Criteria Template](ENHANCED-TEST-CRITERIA-TEMPLATE.md)
- [Milestones Documentation](../milestones.md)
- [Go Testing Guide](../go-testing-guide.md)
- [Project Structure](../project-structure.md)
- [Foundation Testing Guide](FOUNDATION-TESTING-GUIDE.md)
- [Keybinding System](../keybinding-system.md)