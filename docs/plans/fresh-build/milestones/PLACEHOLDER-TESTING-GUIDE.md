# Placeholder Testing Guide (Milestones 11-14)

**Milestone Group**: Placeholders  
**Milestones**: M11-M14  
**Focus**: Parse, navigate, edit text/list placeholders

## Overview

This guide provides comprehensive testing strategies for the Placeholder milestone group, which enables dynamic content in prompts through template syntax. Testing focuses on ensuring accurate parsing, smooth navigation, and reliable editing of both text and list placeholders.

## Integration Tests

### Test Suite: Placeholder Integration

**Location**: `internal/placeholder/integration_test.go`

#### Test 1: Placeholder Parsing
```go
func TestPlaceholderParsing(t *testing.T) {
    // Test comprehensive placeholder parsing
    // 1. Parse text placeholders {{text:name}}
    // 2. Parse list placeholders {{list:name}}
    // 3. Parse multiple placeholders
    // 4. Parse nested placeholders
    // 5. Parse invalid placeholders
}
```

**Acceptance Criteria**:
- [ ] Text placeholders parsed correctly
- [ ] List placeholders parsed correctly
- [ ] Multiple placeholders in same document parsed
- [ ] Placeholder names extracted correctly
- [ ] Invalid placeholders handled gracefully
- [ ] Malformed syntax detected and reported
- [ ] Parsing completes in <100ms for 100 placeholders

#### Test 2: Placeholder Navigation
```go
func TestPlaceholderNavigation(t *testing.T) {
    // Test navigation between placeholders
    // 1. Navigate to next placeholder
    // 2. Navigate to previous placeholder
    // 3. Jump to specific placeholder
    // 4. Navigate with keyboard shortcuts
    // 5. Navigate with mouse clicks
}
```

**Acceptance Criteria**:
- [ ] Next placeholder navigation works
- [ ] Previous placeholder navigation works
- [ ] Jump to specific placeholder works
- [ ] Keyboard shortcuts work (Tab, Shift+Tab)
- [ ] Mouse clicks work
- [ ] Cursor positioned correctly at placeholder
- [ ] Navigation is smooth (<50ms response)

#### Test 3: Text Placeholder Editing
```go
func TestTextPlaceholderEditing(t *testing.T) {
    // Test editing text placeholders
    // 1. Enter text placeholder
    // 2. Type content
    // 3. Exit placeholder
    // 4. Verify content preserved
    // 5. Edit existing content
}
```

**Acceptance Criteria**:
- [ ] Enter placeholder works
- [ ] Typing works normally
- [ ] Exit placeholder works (Tab, Enter, Esc)
- [ ] Content preserved after exit
- [ ] Placeholder syntax preserved
- [ ] Cursor positioned correctly
- [ ] Undo/redo works for edits

#### Test 4: List Placeholder Editing
```go
func TestListPlaceholderEditing(t *testing.T) {
    // Test editing list placeholders
    // 1. Enter list placeholder
    // 2. Add items
    // 3. Remove items
    // 4. Edit items
    // 5. Reorder items
    // 6. Exit placeholder
}
```

**Acceptance Criteria**:
- [ ] Enter placeholder works
- [ ] Add item works (Enter, Ctrl+Enter)
- [ ] Remove item works (Ctrl+D, Backspace)
- [ ] Edit item works
- [ ] Reorder items works (Ctrl+Up/Down)
- [ ] Exit placeholder works
- [ ] List syntax preserved

#### Test 5: Placeholder with Undo/Redo
```go
func TestPlaceholderWithUndoRedo(t *testing.T) {
    // Test undo/redo with placeholder operations
    // 1. Edit text placeholder
    // 2. Edit list placeholder
    // 3. Undo edits
    // 4. Redo edits
    // 5. Verify state consistency
}
```

**Acceptance Criteria**:
- [ ] Text placeholder edits recorded
- [ ] List placeholder edits recorded
- [ ] Undo restores previous state
- [ ] Redo reapplies changes
- [ ] Placeholder syntax preserved
- [ ] No corruption after undo/redo

#### Test 6: Placeholder Validation
```go
func TestPlaceholderValidation(t *testing.T) {
    // Test placeholder validation
    // 1. Validate text placeholder names
    // 2. Validate list placeholder names
    // 3. Detect duplicate names
    // 4. Detect invalid characters
    // 5. Report validation errors
}
```

**Acceptance Criteria**:
- [ ] Valid names accepted
- [ ] Invalid names rejected
- [ ] Duplicate names detected
- [ ] Invalid characters detected
- [ ] Validation errors reported clearly
- [ ] Validation completes in <50ms

## End-to-End Scenarios

### Scenario 1: Simple Text Placeholder Workflow

**Description**: Test basic text placeholder usage from insertion to editing.

**Steps**:
1. User opens editor
2. User inserts prompt with text placeholder
3. User navigates to placeholder
4. User enters placeholder
5. User types content
6. User exits placeholder
7. User saves file
8. User reopens file
9. User verifies content

**Expected Results**:
- [ ] Prompt inserted correctly
- [ ] Placeholder syntax preserved
- [ ] Navigation works
- [ ] Content typed correctly
- [ ] Placeholder syntax preserved after exit
- [ ] File saved correctly
- [ ] Content preserved on reopen

**Test Script**:
```bash
#!/bin/bash
# Test simple text placeholder workflow
promptstack edit test.md
# Insert prompt with {{text:name}}
# Navigate to placeholder
# Type "John Doe"
# Exit placeholder
# Save file
# Reopen file
# Verify content
```

### Scenario 2: Complex List Placeholder Workflow

**Description**: Test list placeholder with multiple items and operations.

**Steps**:
1. User opens editor
2. User inserts prompt with list placeholder
3. User navigates to placeholder
4. User enters placeholder
5. User adds 5 items
6. User edits 2 items
7. User removes 1 item
8. User reorders items
9. User exits placeholder
10. User saves file

**Expected Results**:
- [ ] Prompt inserted correctly
- [ ] List placeholder syntax preserved
- [ ] All 5 items added
- [ ] 2 items edited correctly
- [ ] 1 item removed
- [ ] Items reordered correctly
- [ ] List syntax preserved after exit
- [ ] File saved correctly

**Test Script**:
```bash
#!/bin/bash
# Test complex list placeholder workflow
promptstack edit test.md
# Insert prompt with {{list:items}}
# Navigate to placeholder
# Add 5 items
# Edit 2 items
# Remove 1 item
# Reorder items
# Exit placeholder
# Save file
# Verify content
```

### Scenario 3: Multiple Placeholders in Document

**Description**: Test document with multiple placeholders of different types.

**Steps**:
1. User opens editor
2. User inserts prompt with 3 text placeholders
3. User inserts prompt with 2 list placeholders
4. User navigates between all placeholders
5. User fills all text placeholders
6. User fills all list placeholders
7. User edits some placeholders
8. User uses undo/redo
9. User saves file

**Expected Results**:
- [ ] All placeholders inserted correctly
- [ ] Navigation works between all placeholders
- [ ] All text placeholders filled
- [ ] All list placeholders filled
- [ ] Edits work correctly
- [ ] Undo/redo works correctly
- [ ] All placeholders preserved
- [ ] File saved correctly

**Test Script**:
```bash
#!/bin/bash
# Test multiple placeholders in document
promptstack edit test.md
# Insert prompt with 3 text placeholders
# Insert prompt with 2 list placeholders
# Navigate between all placeholders
# Fill all text placeholders
# Fill all list placeholders
# Edit some placeholders
# Use undo/redo
# Save file
# Verify content
```

### Scenario 4: Placeholder with Special Characters

**Description**: Test placeholders with special characters and edge cases.

**Steps**:
1. User opens editor
2. User inserts prompt with text placeholder
3. User types content with special characters
4. User types content with newlines
5. User types content with unicode
6. User exits placeholder
7. User verifies content preserved

**Expected Results**:
- [ ] Special characters preserved
- [ ] Newlines preserved
- [ ] Unicode characters preserved
- [ ] Placeholder syntax preserved
- [ ] No corruption
- [ ] Content displays correctly

**Test Script**:
```bash
#!/bin/bash
# Test placeholder with special characters
promptstack edit test.md
# Insert prompt with {{text:content}}
# Type: "Hello, World! @#$%^&*()"
# Type: "Line 1\nLine 2\nLine 3"
# Type: "Hello 世界 🌍"
# Exit placeholder
# Verify content
```

### Scenario 5: Placeholder Error Recovery

**Description**: Test error handling and recovery with placeholders.

**Steps**:
1. User opens editor
2. User inserts prompt with invalid placeholder syntax
3. User sees error message
4. User fixes syntax
5. User edits placeholder
6. User types very long content
7. User exits placeholder
8. User verifies no errors

**Expected Results**:
- [ ] Invalid syntax detected
- [ ] Error message displayed clearly
- [ ] Fix works correctly
- [ ] Long content handled
- [ ] No crashes
- [ ] No corruption

**Test Script**:
```bash
#!/bin/bash
# Test placeholder error recovery
promptstack edit test.md
# Insert: "{{text:name" (missing closing brace)
# Verify error message
# Fix: "{{text:name}}"
# Edit placeholder
# Type 1000 characters
# Exit placeholder
# Verify no errors
```

## Performance Benchmarks

### Benchmark 1: Placeholder Parsing

**Test**: Measure performance of placeholder parsing

```go
func BenchmarkPlaceholderParsing(b *testing.B) {
    content := LoadTestContent("with-placeholders.md")
    for i := 0; i < b.N; i++ {
        ParsePlaceholders(content)
    }
}
```

**Thresholds**:
- [ ] 10 placeholders: <10ms
- [ ] 50 placeholders: <50ms
- [ ] 100 placeholders: <100ms
- [ ] 500 placeholders: <500ms

### Benchmark 2: Placeholder Navigation

**Test**: Measure performance of placeholder navigation

```go
func BenchmarkPlaceholderNavigation(b *testing.B) {
    doc := ParseDocument("with-placeholders.md")
    for i := 0; i < b.N; i++ {
        doc.NavigateToNextPlaceholder()
    }
}
```

**Thresholds**:
- [ ] Navigation: <10ms per operation
- [ ] Jump to placeholder: <10ms
- [ ] No UI lag during navigation

### Benchmark 3: Text Placeholder Editing

**Test**: Measure performance of text placeholder editing

```go
func BenchmarkTextPlaceholderEditing(b *testing.B) {
    doc := ParseDocument("with-placeholders.md")
    placeholder := doc.GetPlaceholder(0)
    for i := 0; i < b.N; i++ {
        placeholder.SetText("test content")
    }
}
```

**Thresholds**:
- [ ] Set text: <10ms
- [ ] Get text: <10ms
- [ ] Typing latency: <50ms
- [ ] No UI blocking

### Benchmark 4: List Placeholder Editing

**Test**: Measure performance of list placeholder editing

```go
func BenchmarkListPlaceholderEditing(b *testing.B) {
    doc := ParseDocument("with-placeholders.md")
    placeholder := doc.GetListPlaceholder(0)
    for i := 0; i < b.N; i++ {
        placeholder.AddItem("item")
    }
}
```

**Thresholds**:
- [ ] Add item: <10ms
- [ ] Remove item: <10ms
- [ ] Edit item: <10ms
- [ ] Reorder items: <10ms
- [ ] No UI blocking

### Benchmark 5: Placeholder Validation

**Test**: Measure performance of placeholder validation

```go
func BenchmarkPlaceholderValidation(b *testing.B) {
    doc := ParseDocument("with-placeholders.md")
    for i := 0; i < b.N; i++ {
        doc.ValidatePlaceholders()
    }
}
```

**Thresholds**:
- [ ] 10 placeholders: <10ms
- [ ] 50 placeholders: <50ms
- [ ] 100 placeholders: <100ms
- [ ] 500 placeholders: <500ms

## Test Execution

### Running Integration Tests

```bash
# Run all placeholder integration tests
go test ./internal/placeholder -v -tags=integration

# Run specific test
go test ./internal/placeholder -v -run TestPlaceholderParsing

# Run with coverage
go test ./internal/placeholder -cover -coverprofile=coverage.out
```

### Running End-to-End Tests

```bash
# Run all E2E scenarios
./scripts/test/placeholder-e2e.sh

# Run specific scenario
./scripts/test/placeholder-e2e.sh scenario1

# Run with performance monitoring
./scripts/test/placeholder-e2e.sh --perf
```

### Running Benchmarks

```bash
# Run all benchmarks
go test ./internal/placeholder -bench=. -benchmem

# Run specific benchmark
go test ./internal/placeholder -bench=BenchmarkPlaceholderParsing

# Run with CPU profiling
go test ./internal/placeholder -bench=. -cpuprofile=cpu.prof
```

## Test Data

### Sample Documents

**Location**: `test/data/placeholders/`

- `simple-text.md` - Document with simple text placeholder
- `simple-list.md` - Document with simple list placeholder
- `multiple-placeholders.md` - Document with multiple placeholders
- `nested-placeholders.md` - Document with nested placeholders
- `invalid-syntax.md` - Document with invalid placeholder syntax
- `special-chars.md` - Document with special characters
- `large-document.md` - Document with 100 placeholders

### Sample Placeholders

**Location**: `test/data/placeholders/examples/`

- `text-basic.md` - Basic text placeholder
- `text-multiline.md` - Multiline text placeholder
- `list-basic.md` - Basic list placeholder
- `list-large.md` - List with 100 items
- `mixed.md` - Mixed text and list placeholders

## Success Criteria

### Integration Tests
- [ ] All integration tests pass
- [ ] Code coverage >80% for placeholder components
- [ ] No memory leaks detected
- [ ] No race conditions detected

### End-to-End Scenarios
- [ ] All scenarios complete successfully
- [ ] Placeholders work correctly in all scenarios
- [ ] Navigation works smoothly
- [ ] Editing works correctly
- [ ] Performance meets thresholds

### Performance Benchmarks
- [ ] All benchmarks meet thresholds
- [ ] No performance regression from baseline
- [ ] Memory usage remains bounded
- [ ] CPU usage is reasonable

## Known Issues and Limitations

### Current Limitations
- Very long placeholder content (>10KB) may cause performance issues
- Complex nested placeholders may have parsing edge cases
- Unicode handling may have edge cases in some scenarios
- Placeholder names with special characters may cause issues

### Future Improvements
- Add placeholder templates with default values
- Implement placeholder validation rules
- Add placeholder autocomplete
- Optimize parsing for very large documents
- Add placeholder snippets

## References

- [Enhanced Test Criteria Template](ENHANCED-TEST-CRITERIA-TEMPLATE.md)
- [Milestones Documentation](../milestones.md)
- [Go Testing Guide](../go-testing-guide.md)
- [Project Structure](../project-structure.md)
- [Foundation Testing Guide](FOUNDATION-TESTING-GUIDE.md)
- [Library Integration Testing Guide](LIBRARY-INTEGRATION-TESTING-GUIDE.md)