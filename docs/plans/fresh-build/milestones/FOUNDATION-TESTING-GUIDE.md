# Foundation Testing Guide (Milestones 1-6)

**Milestone Group**: Foundation  
**Milestones**: M1-M6  
**Focus**: Bootstrap, TUI shell, file I/O, editor, auto-save, undo/redo

## Overview

This guide provides comprehensive testing strategies for the Foundation milestone group, which establishes the core infrastructure for the PromptStack application. Testing focuses on ensuring robust file operations, reliable text editing, and seamless user interactions.

## Integration Tests

### Test Suite: Foundation Integration

**Location**: `internal/foundation/integration_test.go`

#### Test 1: Bootstrap to Editor Workflow
```go
func TestBootstrapToEditorWorkflow(t *testing.T) {
    // Test complete flow from bootstrap to text editor
    // 1. Initialize config
    // 2. Launch TUI shell
    // 3. Create new file
    // 4. Write text
    // 5. Verify file saved
}
```

**Acceptance Criteria**:
- [ ] Config loads successfully from default location
- [ ] TUI shell initializes without errors
- [ ] New file creation works
- [ ] Text input is captured correctly
- [ ] File is saved to disk with correct content

#### Test 2: File I/O with Auto-save
```go
func TestFileIOWithAutosave(t *testing.T) {
    // Test file operations with auto-save enabled
    // 1. Open existing file
    // 2. Make modifications
    // 3. Wait for auto-save interval
    // 4. Verify file updated on disk
    // 5. Verify no data loss
}
```

**Acceptance Criteria**:
- [ ] File opens successfully
- [ ] Modifications are tracked
- [ ] Auto-save triggers at correct interval
- [ ] File content matches editor state
- [ ] No data loss during auto-save

#### Test 3: Undo/Redo with File Operations
```go
func TestUndoRedoWithFileOperations(t *testing.T) {
    // Test undo/redo across file operations
    // 1. Open file
    // 2. Make multiple edits
    [ ] Undo several changes
    [ ] Redo changes
    [ ] Verify state consistency
}
```

**Acceptance Criteria**:
- [ ] Each edit is recorded in history
- [ ] Undo restores previous state correctly
- [ ] Redo reapplies changes correctly
- [ ] History persists across file operations
- [ ] Memory usage remains bounded

#### Test 4: Config Changes Runtime
```go
func TestConfigChangesRuntime(t *testing.T) {
    // Test config changes while app is running
    // 1. Start app with initial config
    // 2. Modify config file
    // 3. Trigger config reload
    // 4. Verify changes applied
}
```

**Acceptance Criteria**:
- [ ] Config reloads without restart
- [ ] Changes are applied immediately
- [ ] Invalid config is rejected gracefully
- [ ] Default values used for missing keys
- [ ] No crashes on config errors

#### Test 5: TUI State Management
```go
func TestTUIStateManagement(t *testing.T) {
    // Test TUI state across multiple operations
    // 1. Initialize TUI
    // 2. Navigate between views
    // 3. Perform operations
    // 4. Verify state consistency
}
```

**Acceptance Criteria**:
- [ ] State transitions are correct
- [ ] No state corruption
- [ ] Memory leaks are prevented
- [ ] UI updates are timely
- [ ] Error states are handled

## End-to-End Scenarios

### Scenario 1: First-Time User Setup

**Description**: Simulate a new user installing and using PromptStack for the first time.

**Steps**:
1. User runs `promptstack` for the first time
2. Bootstrap wizard creates default config
3. User creates first prompt file
4. User writes content
5. User saves and exits
6. User reopens file
7. User makes edits and uses undo/redo
8. User saves changes

**Expected Results**:
- [ ] Config created at `~/.promptstack/config.yaml`
- [ ] Default library directory created
- [ ] File created successfully
- [ ] Content saved correctly
- [ ] File reopens with correct content
- [ ] Undo/redo works as expected
- [ ] No data loss

**Test Script**:
```bash
#!/bin/bash
# Test first-time user setup
rm -rf ~/.promptstack
promptstack
# Follow wizard prompts
# Create test file
# Write content
# Save and exit
# Verify file exists
# Reopen and verify content
```

### Scenario 2: Editing Session with Auto-save

**Description**: Simulate a typical editing session with auto-save enabled.

**Steps**:
1. User opens existing file
2. User makes multiple edits over 5 minutes
3. Auto-save triggers every 30 seconds
4. User continues editing
5. User manually saves
6. User exits
7. User reopens file
8. Verify all changes preserved

**Expected Results**:
- [ ] File opens successfully
- [ ] All edits are captured
- [ ] Auto-save works at correct intervals
- [ ] Manual save works
- [ ] All changes preserved on reopen
- [ ] No data loss

**Test Script**:
```bash
#!/bin/bash
# Test editing session with auto-save
promptstack edit test.md
# Make edits
# Wait for auto-save
# Make more edits
# Manually save
# Exit
# Verify file content
```

### Scenario 3: Complex Undo/Redo Workflow

**Description**: Test complex undo/redo scenarios with multiple file operations.

**Steps**:
1. User opens file A
2. User makes 10 edits
3. User undoes 5 edits
4. User makes 3 new edits
5. User redoes 3 edits
6. User opens file B
7. User makes edits in file B
8. User switches back to file A
9. User continues undo/redo
10. Verify both files have correct content

**Expected Results**:
- [ ] Undo/redo works correctly in file A
- [ ] New edits after undo work correctly
- [ ] Redo works correctly
- [ ] File B operations don't affect file A
- [ ] Switching between files preserves history
- [ ] Both files have correct final content

**Test Script**:
```bash
#!/bin/bash
# Test complex undo/redo workflow
promptstack edit fileA.md
# Make 10 edits
# Undo 5 edits
# Make 3 new edits
# Redo 3 edits
# Open fileB.md
# Make edits
# Switch back to fileA.md
# Continue undo/redo
# Verify both files
```

### Scenario 4: Error Recovery

**Description**: Test recovery from various error conditions.

**Steps**:
1. User opens file
2. Simulate disk full error during save
3. Verify error message displayed
4. User fixes disk space
5. User retries save
6. Verify save succeeds
7. Simulate file permission error
8. Verify error message displayed
9. User fixes permissions
10. User retries save
11. Verify save succeeds

**Expected Results**:
- [ ] Disk full error caught and reported
- [ ] App doesn't crash
- [ ] Retry works after fixing issue
- [ ] Permission error caught and reported
- [ ] App doesn't crash
- [ ] Retry works after fixing issue
- [ ] No data loss

**Test Script**:
```bash
#!/bin/bash
# Test error recovery
# Fill disk
promptstack edit test.md
# Try to save
# Verify error
# Free disk space
# Retry save
# Verify success
# Change permissions
# Try to save
# Verify error
# Fix permissions
# Retry save
# Verify success
```

### Scenario 5: Performance Under Load

**Description**: Test performance with large files and rapid operations.

**Steps**:
1. User opens 10MB file
2. User makes rapid edits (100 edits in 10 seconds)
3. User performs undo/redo operations
4. User saves file
5. Measure performance metrics
6. Verify no lag or crashes

**Expected Results**:
- [ ] Large file opens in <2 seconds
- [ ] Edits are responsive (<100ms latency)
- [ ] Undo/redo is responsive (<100ms latency)
- [ ] Save completes in <1 second
- [ ] No memory leaks
- [ ] No crashes

**Test Script**:
```bash
#!/bin/bash
# Test performance under load
# Create 10MB test file
promptstack edit largefile.md
# Make rapid edits
# Perform undo/redo
# Save
# Measure metrics
```

## Performance Benchmarks

### Benchmark 1: File Operations

**Test**: Measure performance of file I/O operations

```go
func BenchmarkFileOpen(b *testing.B) {
    for i := 0; i < b.N; i++ {
        OpenFile("test.md")
    }
}

func BenchmarkFileSave(b *testing.B) {
    for i := 0; i < b.N; i++ {
        SaveFile("test.md", content)
    }
}
```

**Thresholds**:
- [ ] File open: <100ms for 1MB file
- [ ] File save: <200ms for 1MB file
- [ ] File open: <500ms for 10MB file
- [ ] File save: <1s for 10MB file

### Benchmark 2: Text Editing

**Test**: Measure performance of text editing operations

```go
func BenchmarkTextInsert(b *testing.B) {
    for i := 0; i < b.N; i++ {
        InsertText("test content")
    }
}

func BenchmarkTextDelete(b *testing.B) {
    for i := 0; i < b.N; i++ {
        DeleteText(10)
    }
}
```

**Thresholds**:
- [ ] Text insert: <10ms per operation
- [ ] Text delete: <10ms per operation
- [ ] Batch insert (100 chars): <50ms
- [ ] Batch delete (100 chars): <50ms

### Benchmark 3: Undo/Redo

**Test**: Measure performance of undo/redo operations

```go
func BenchmarkUndo(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Undo()
    }
}

func BenchmarkRedo(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Redo()
    }
}
```

**Thresholds**:
- [ ] Undo: <50ms per operation
- [ ] Redo: <50ms per operation
- [ ] History size: <100MB for 1000 operations
- [ ] Memory growth: Linear with operation count

### Benchmark 4: Auto-save

**Test**: Measure performance of auto-save operations

```go
func BenchmarkAutosave(b *testing.B) {
    for i := 0; i < b.N; i++ {
        TriggerAutosave()
    }
}
```

**Thresholds**:
- [ ] Auto-save: <500ms for 1MB file
- [ ] Auto-save: <2s for 10MB file
- [ ] No UI blocking during auto-save
- [ ] Minimal CPU usage during auto-save

### Benchmark 5: TUI Rendering

**Test**: Measure performance of TUI rendering

```go
func BenchmarkTUIRender(b *testing.B) {
    for i := 0; i < b.N; i++ {
        RenderTUI()
    }
}
```

**Thresholds**:
- [ ] Render: <16ms (60 FPS)
- [ ] Input response: <16ms
- [ ] Screen updates: <16ms
- [ ] No frame drops during rapid input

## Test Execution

### Running Integration Tests

```bash
# Run all foundation integration tests
go test ./internal/foundation -v -tags=integration

# Run specific test
go test ./internal/foundation -v -run TestBootstrapToEditorWorkflow

# Run with coverage
go test ./internal/foundation -cover -coverprofile=coverage.out
```

### Running End-to-End Tests

```bash
# Run all E2E scenarios
./scripts/test/foundation-e2e.sh

# Run specific scenario
./scripts/test/foundation-e2e.sh scenario1

# Run with performance monitoring
./scripts/test/foundation-e2e.sh --perf
```

### Running Benchmarks

```bash
# Run all benchmarks
go test ./internal/foundation -bench=. -benchmem

# Run specific benchmark
go test ./internal/foundation -bench=BenchmarkFileOpen

# Run with CPU profiling
go test ./internal/foundation -bench=. -cpuprofile=cpu.prof
```

## Test Data

### Sample Files

**Location**: `test/data/foundation/`

- `small.md` (1KB)
- `medium.md` (100KB)
- `large.md` (10MB)
- `unicode.md` (UTF-8 content)
- `special-chars.md` (special characters)

### Sample Configs

**Location**: `test/configs/`

- `default.yaml` - Default configuration
- `minimal.yaml` - Minimal configuration
- `custom.yaml` - Custom configuration
- `invalid.yaml` - Invalid configuration (for error testing)

## Success Criteria

### Integration Tests
- [ ] All integration tests pass
- [ ] Code coverage >80% for foundation components
- [ ] No memory leaks detected
- [ ] No race conditions detected

### End-to-End Scenarios
- [ ] All scenarios complete successfully
- [ ] No data loss in any scenario
- [ ] Error recovery works correctly
- [ ] Performance meets thresholds

### Performance Benchmarks
- [ ] All benchmarks meet thresholds
- [ ] No performance regression from baseline
- [ ] Memory usage remains bounded
- [ ] CPU usage is reasonable

## Known Issues and Limitations

### Current Limitations
- Large files (>100MB) may cause performance issues
- Unicode handling may have edge cases
- Auto-save may conflict with manual save in rare cases

### Future Improvements
- Add support for incremental file loading
- Implement more efficient undo/redo for large files
- Add parallel auto-save for multiple files
- Optimize TUI rendering for large content

## References

- [Enhanced Test Criteria Template](ENHANCED-TEST-CRITERIA-TEMPLATE.md)
- [Milestones Documentation](../milestones.md)
- [Go Testing Guide](../go-testing-guide.md)
- [Project Structure](../project-structure.md)