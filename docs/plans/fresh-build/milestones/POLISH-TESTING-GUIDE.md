# Polish Testing Guide (Milestones 36-38)

**Milestone Group**: Polish  
**Milestones**: M36-M38  
**Focus**: Settings, responsive layout, error handling

## Overview

This guide provides comprehensive testing strategies for the Polish milestone group, which implements settings management, responsive layout, and comprehensive error handling with log viewer. Testing focuses on ensuring reliable settings persistence, smooth responsive behavior, and robust error recovery.

## Integration Tests

### Test Suite: Polish Integration

**Location**: `internal/polish/integration_test.go`

#### Test 1: Settings Panel
```go
func TestSettingsPanel(t *testing.T) {
    // Test settings panel functionality
    // 1. Open settings panel
    // 2. Navigate settings categories
    // 3. Modify settings
    // 4. Save settings
    // 5. Reload settings
    // 6. Reset to defaults
}
```

**Acceptance Criteria**:
- [ ] Settings panel opens quickly
- [ ] All categories accessible
- [ ] Settings modify correctly
- [ ] Settings save correctly
- [ ] Settings reload correctly
- [ ] Reset to defaults works
- [ ] Invalid values rejected
- [ ] Settings persist across restarts

#### Test 2: Responsive Layout
```go
func TestResponsiveLayout(t *testing.T) {
    // Test responsive layout behavior
    // 1. Test minimum terminal size
    // 2. Test medium terminal size
    // 3. Test large terminal size
    // 4. Test dynamic resizing
    // 5. Test layout transitions
}
```

**Acceptance Criteria**:
- [ ] Minimum size (80x24) works
- [ ] Medium size (120x40) works
- [ ] Large size (200x60) works
- [ ] Dynamic resizing works
- [ ] Layout transitions smooth
- [ ] No content clipping
- [ ] No UI corruption
- [ ] Resizing completes in <100ms

#### Test 3: Error Handling
```go
func TestErrorHandling(t *testing.T) {
    // Test comprehensive error handling
    // 1. Test file not found errors
    // 2. Test permission errors
    // 3. Test network errors
    // 4. Test API errors
    // 5. Test validation errors
}
```

**Acceptance Criteria**:
- [ ] File not found handled
- [ ] Permission errors handled
- [ ] Network errors handled
- [ ] API errors handled
- [ ] Validation errors handled
- [ ] Error messages clear
- [ ] No crashes on errors
- [ ] Recovery works

#### Test 4: Log Viewer
```go
func TestLogViewer(t *testing.T) {
    // Test log viewer functionality
    // 1. Open log viewer
    // 2. View log entries
    // 3. Filter logs by level
    // 4. Search logs
    // 5. Export logs
}
```

**Acceptance Criteria**:
- [ ] Log viewer opens quickly
- [ ] All entries displayed
- [ ] Filtering works (INFO, WARN, ERROR)
- [ ] Search works
- [ ] Export works
- [ ] No UI lag with 1000+ entries
- [ ] Real-time updates work

#### Test 5: Settings Validation
```go
func TestSettingsValidation(t *testing.T) {
    // Test settings validation
    // 1. Validate numeric ranges
    // 2. Validate string formats
    // 3. Validate file paths
    // 4. Validate boolean values
    // 5. Test invalid values
}
```

**Acceptance Criteria**:
- [ ] Numeric ranges validated
- [ ] String formats validated
- [ ] File paths validated
- [ ] Boolean values validated
- [ ] Invalid values rejected
- [ ] Validation errors clear
- [ ] Default values used for invalid

#### Test 6: Error Recovery
```go
func TestErrorRecovery(t *testing.T) {
    // Test error recovery mechanisms
    // 1. Simulate crash
    // 2. Verify auto-recovery
    // 3. Verify data preserved
    // 4. Test graceful degradation
    // 5. Test fallback behavior
}
```

**Acceptance Criteria**:
- [ ] Crash detected
- [ ] Auto-recovery works
- [ ] Data preserved
- [ ] Graceful degradation works
- [ ] Fallback behavior works
- [ ] No data loss
- [ ] Recovery completes in <5s

## End-to-End Scenarios

### Scenario 1: Settings Configuration Workflow

**Description**: Test complete settings configuration workflow.

**Steps**:
1. User opens settings panel
2. User navigates to editor settings
3. User modifies tab width
4. User navigates to theme settings
5. User changes theme
6. User navigates to keybinding settings
7. User enables Vim mode
8. User saves settings
9. User restarts PromptStack
10. User verifies settings persisted

**Expected Results**:
- [ ] Settings panel opens
- [ ] Navigation works
- [ ] Tab width modified
- [ ] Theme changed
- [ ] Vim mode enabled
- [ ] Settings saved
- [ ] Settings persisted
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test settings configuration workflow
promptstack
# Open settings panel
# Navigate to editor settings
# Modify tab width
# Navigate to theme settings
# Change theme
# Navigate to keybinding settings
# Enable Vim mode
# Save settings
# Restart PromptStack
# Verify settings persisted
```

### Scenario 2: Responsive Layout Resizing

**Description**: Test responsive layout with dynamic terminal resizing.

**Steps**:
1. User starts PromptStack in 80x24 terminal
2. User opens editor with content
3. User resizes to 120x40
4. User resizes to 200x60
5. User resizes back to 80x24
6. User opens library browser
7. User resizes to 100x30
8. User opens command palette
9. User resizes to 150x50

**Expected Results**:
- [ ] Layout works at 80x24
- [ ] Layout adapts to 120x40
- [ ] Layout adapts to 200x60
- [ ] Layout adapts back to 80x24
- [ ] Browser works at 100x30
- [ ] Palette works at 150x50
- [ ] No content clipping
- [ ] No UI corruption
- [ ] Smooth transitions

**Test Script**:
```bash
#!/bin/bash
# Test responsive layout resizing
promptstack edit test.md
# Resize to 80x24
# Resize to 120x40
# Resize to 200x60
# Resize to 80x24
# Open library browser
# Resize to 100x30
# Open command palette
# Resize to 150x50
# Verify layout
```

### Scenario 3: Error Handling and Recovery

**Description**: Test error handling and recovery from various error conditions.

**Steps**:
1. User tries to open non-existent file
2. User sees error message
3. User tries to save to read-only location
4. User sees error message
5. User fixes permissions
6. User retries save
7. User experiences network error
8. User sees error message
9. User retries operation
10. User verifies recovery

**Expected Results**:
- [ ] File not found error displayed
- [ ] Error message clear
- [ ] Permission error displayed
- [ ] Error message clear
- [ ] Fix works
- [ ] Retry works
- [ ] Network error displayed
- [ ] Error message clear
- [ ] Retry works
- [ ] No crashes

**Test Script**:
```bash
#!/bin/bash
# Test error handling and recovery
promptstack edit nonexistent.md
# Verify error message
# Try to save to read-only location
# Verify error message
# Fix permissions
# Retry save
# Simulate network error
# Verify error message
# Retry operation
# Verify recovery
```

### Scenario 4: Log Viewer and Debugging

**Description**: Test log viewer for debugging and troubleshooting.

**Steps**:
1. User performs various operations
2. User opens log viewer
3. User sees all log entries
4. User filters by ERROR level
5. User searches for specific error
6. User exports logs to file
7. User shares logs for debugging
8. User clears logs
9. User verifies logs cleared

**Expected Results**:
- [ ] Operations logged
- [ ] Log viewer opens
- [ ] All entries displayed
- [ ] Filtering works
- [ ] Search works
- [ ] Export works
- [ ] Logs exported
- [ ] Logs cleared
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test log viewer and debugging
promptstack
# Perform various operations
# Open log viewer
# Verify all entries
# Filter by ERROR
# Search for error
# Export logs
# Clear logs
# Verify cleared
```

### Scenario 5: Settings Reset and Defaults

**Description**: Test settings reset to defaults and recovery.

**Steps**:
1. User has custom settings
2. User opens settings panel
3. User resets to defaults
4. User verifies defaults applied
5. User modifies some settings
6. User saves settings
7. User corrupts settings file
8. User restarts PromptStack
9. User verifies defaults restored
10. User reconfigures settings

**Expected Results**:
- [ ] Custom settings loaded
- [ ] Reset to defaults works
- [ ] Defaults applied
- [ ] Settings modified
- [ ] Settings saved
- [ ] Corruption detected
- [ ] Defaults restored
- [ ] No crashes
- [ ] Reconfiguration works

**Test Script**:
```bash
#!/bin/bash
# Test settings reset and defaults
promptstack
# Configure custom settings
# Open settings panel
# Reset to defaults
# Verify defaults
# Modify settings
# Save settings
# Corrupt settings file
# Restart PromptStack
# Verify defaults restored
# Reconfigure settings
```

## Performance Benchmarks

### Benchmark 1: Settings Operations

**Test**: Measure performance of settings operations

```go
func BenchmarkSettingsOperations(b *testing.B) {
    settings := NewSettings()
    for i := 0; i < b.N; i++ {
        settings.Set("key", "value")
        settings.Save()
        settings.Load()
    }
}
```

**Thresholds**:
- [ ] Set value: <1ms
- [ ] Save settings: <10ms
- [ ] Load settings: <10ms
- [ ] Reset to defaults: <10ms

### Benchmark 2: Layout Resizing

**Test**: Measure performance of layout resizing

```go
func BenchmarkLayoutResizing(b *testing.B) {
    layout := NewResponsiveLayout()
    for i := 0; i < b.N; i++ {
        layout.Resize(100, 30)
        layout.Resize(150, 50)
    }
}
```

**Thresholds**:
- [ ] Single resize: <50ms
- [ ] 10 resizes: <100ms
- [ ] 100 resizes: <1s
- [ ] No UI lag

### Benchmark 3: Error Handling

**Test**: Measure performance of error handling

```go
func BenchmarkErrorHandling(b *testing.B) {
    handler := NewErrorHandler()
    for i := 0; i < b.N; i++ {
        handler.HandleError(errors.New("test error"))
    }
}
```

**Thresholds**:
- [ ] Single error: <10ms
- [ ] 10 errors: <50ms
- [ ] 100 errors: <500ms
- [ ] No performance degradation

### Benchmark 4: Log Viewer

**Test**: Measure performance of log viewer operations

```go
func BenchmarkLogViewer(b *testing.B) {
    viewer := NewLogViewer()
    LoadLogs(viewer, 1000)
    for i := 0; i < b.N; i++ {
        viewer.Filter("ERROR")
        viewer.Search("test")
    }
}
```

**Thresholds**:
- [ ] Filter: <50ms
- [ ] Search: <50ms
- [ ] Display 1000 entries: <100ms
- [ ] No UI lag

### Benchmark 5: Settings Validation

**Test**: Measure performance of settings validation

```go
func BenchmarkSettingsValidation(b *testing.B) {
    validator := NewSettingsValidator()
    for i := 0; i < b.N; i++ {
        validator.Validate(settings)
    }
}
```

**Thresholds**:
- [ ] Single validation: <10ms
- [ ] 10 validations: <50ms
- [ ] 100 validations: <500ms
- [ ] No blocking

## Test Execution

### Running Integration Tests

```bash
# Run all polish integration tests
go test ./internal/polish -v -tags=integration

# Run specific test
go test ./internal/polish -v -run TestSettingsPanel

# Run with coverage
go test ./internal/polish -cover -coverprofile=coverage.out
```

### Running End-to-End Tests

```bash
# Run all E2E scenarios
./scripts/test/polish-e2e.sh

# Run specific scenario
./scripts/test/polish-e2e.sh scenario1

# Run with performance monitoring
./scripts/test/polish-e2e.sh --perf
```

### Running Benchmarks

```bash
# Run all benchmarks
go test ./internal/polish -bench=. -benchmem

# Run specific benchmark
go test ./internal/polish -bench=BenchmarkSettingsOperations

# Run with CPU profiling
go test ./internal/polish -bench=. -cpuprofile=cpu.prof
```

## Test Data

### Sample Settings

**Location**: `test/data/polish/settings/`

- `default.yaml` - Default settings
- `custom.yaml` - Custom settings
- `invalid.yaml` - Invalid settings (for testing)
- `corrupted.yaml` - Corrupted settings (for testing)

### Sample Logs

**Location**: `test/data/polish/logs/`

- `small.log` - Small log file (100 entries)
- `medium.log` - Medium log file (1000 entries)
- `large.log` - Large log file (10000 entries)
- `with-errors.log` - Log with errors
- `with-warnings.log` - Log with warnings

### Sample Terminal Sizes

**Location**: `test/data/polish/terminals/`

- `minimum.txt` - 80x24 terminal
- `medium.txt` - 120x40 terminal
- `large.txt` - 200x60 terminal
- `ultrawide.txt` - 300x80 terminal

## Success Criteria

### Integration Tests
- [ ] All integration tests pass
- [ ] Code coverage >80% for polish components
- [ ] No memory leaks detected
- [ ] No race conditions detected

### End-to-End Scenarios
- [ ] All scenarios complete successfully
- [ ] Settings work correctly in all scenarios
- [ ] Responsive layout works correctly
- [ ] Error handling works correctly
- [ ] Log viewer works correctly
- [ ] Performance meets thresholds

### Performance Benchmarks
- [ ] All benchmarks meet thresholds
- [ ] No performance regression from baseline
- [ ] Memory usage remains bounded
- [ ] CPU usage is reasonable

## Known Issues and Limitations

### Current Limitations
- Very small terminals (<80x24) may have layout issues
- Very large terminals (>300x100) may have performance issues
- Settings validation may not catch all edge cases
- Log viewer may have latency with 10000+ entries

### Future Improvements
- Add more settings options
- Implement more sophisticated responsive layouts
- Add log analysis and insights
- Implement settings profiles
- Add settings import/export

## References

- [Enhanced Test Criteria Template](ENHANCED-TEST-CRITERIA-TEMPLATE.md)
- [Milestones Documentation](../milestones.md)
- [Go Testing Guide](../go-testing-guide.md)
- [Project Structure](../project-structure.md)
- [Foundation Testing Guide](FOUNDATION-TESTING-GUIDE.md)
- [Config Schema](../CONFIG-SCHEMA.md)