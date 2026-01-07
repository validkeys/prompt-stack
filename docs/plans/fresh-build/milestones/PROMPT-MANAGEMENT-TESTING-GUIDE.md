# Prompt Management Testing Guide (Milestones 23-26)

**Milestone Group**: Prompt Management  
**Milestones**: M23-M26  
**Focus**: Validation, results, creator, editor

## Overview

This guide provides comprehensive testing strategies for the Prompt Management milestone group, which implements prompt validation, validation results display, prompt creation, and prompt editing. Testing focuses on ensuring accurate validation, clear error reporting, intuitive creation workflow, and seamless editing experience.

## Integration Tests

### Test Suite: Prompt Management Integration

**Location**: `internal/prompt/integration_test.go`

#### Test 1: Prompt Validation
```go
func TestPromptValidation(t *testing.T) {
    // Test comprehensive prompt validation
    // 1. Validate valid prompt
    // 2. Validate invalid YAML
    // 3. Validate missing required fields
    // 4. Validate invalid placeholder syntax
    // 5. Validate duplicate placeholder names
}
```

**Acceptance Criteria**:
- [ ] Valid prompts pass validation
- [ ] Invalid YAML detected
- [ ] Missing required fields detected
- [ ] Invalid placeholder syntax detected
- [ ] Duplicate placeholder names detected
- [ ] Validation completes in <100ms
- [ ] Clear error messages provided

#### Test 2: Validation Results Display
```go
func TestValidationResultsDisplay(t *testing.T) {
    // Test validation results UI
    // 1. Display validation errors
    // 2. Display validation warnings
    // 3. Navigate between errors
    // 4. Fix errors interactively
    // 5. Re-validate after fixes
}
```

**Acceptance Criteria**:
- [ ] Errors displayed clearly
- [ ] Warnings displayed clearly
- [ ] Navigation between errors works
- [ ] Interactive fixing works
- [ ] Re-validation works
- [ ] Error locations highlighted
- [ ] No UI lag with 100+ errors

#### Test 3: Prompt Creator
```go
func TestPromptCreator(t *testing.T) {
    // Test prompt creation workflow
    // 1. Open prompt creator
    // 2. Enter prompt metadata
    // 3. Write prompt content
    // 4. Add placeholders
    // 5. Validate prompt
    // 6. Save prompt
}
```

**Acceptance Criteria**:
- [ ] Creator opens quickly
- [ ] Metadata entry works
- [ ] Content editing works
- [ ] Placeholder insertion works
- [ ] Validation works
- [ ] Save works correctly
- [ ] File created in correct location

#### Test 4: Prompt Editor
```go
func TestPromptEditor(t *testing.T) {
    // Test prompt editing workflow
    // 1. Open existing prompt
    // 2. Edit metadata
    // 3. Edit content
    // 4. Add/remove placeholders
    // 5. Validate changes
    // 6. Save changes
}
```

**Acceptance Criteria**:
- [ ] Prompt opens correctly
- [ ] Metadata editing works
- [ ] Content editing works
- [ ] Placeholder editing works
- [ ] Validation works
- [ ] Save works correctly
- [ ] Original file preserved

#### Test 5: Validation with Placeholders
```go
func TestValidationWithPlaceholders(t *testing.T) {
    // Test validation with various placeholders
    // 1. Validate text placeholders
    // 2. Validate list placeholders
    // 3. Validate nested placeholders
    // 4. Validate invalid placeholders
    // 5. Validate duplicate names
}
```

**Acceptance Criteria**:
- [ ] Text placeholders validated
- [ ] List placeholders validated
- [ ] Nested placeholders validated
- [ ] Invalid placeholders detected
- [ ] Duplicate names detected
- [ ] Placeholder syntax checked
- [ ] Clear error messages

#### Test 6: Real-time Validation
```go
func TestRealtimeValidation(t *testing.T) {
    // Test real-time validation during editing
    // 1. Type content
    // 2. Validate in real-time
    // 3. Show errors immediately
    // 4. Update errors as typing
    // 5. Clear errors when fixed
}
```

**Acceptance Criteria**:
- [ ] Validation triggers on input
- [ ] Errors shown immediately
- [ ] Errors update as typing
- [ ] Errors clear when fixed
- [ ] No UI lag during typing
- [ ] Debouncing works correctly
- [ ] Performance remains responsive

## End-to-End Scenarios

### Scenario 1: Creating a New Prompt

**Description**: Test complete prompt creation workflow from start to finish.

**Steps**:
1. User opens prompt creator
2. User enters title
3. User enters description
4. User adds tags
5. User writes prompt content
6. User inserts text placeholder
7. User inserts list placeholder
8. User validates prompt
9. User saves prompt
10. User verifies file created

**Expected Results**:
- [ ] Creator opens quickly
- [ ] All metadata entered
- [ ] Content written
- [ ] Placeholders inserted
- [ ] Validation passes
- [ ] File created
- [ ] File in correct location
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test creating a new prompt
promptstack create
# Enter title: "Code Review Prompt"
# Enter description: "A prompt for code reviews"
# Add tags: "review", "code"
# Write content with placeholders
# Validate
# Save
# Verify file created
```

### Scenario 2: Editing an Existing Prompt

**Description**: Test prompt editing workflow with various changes.

**Steps**:
1. User opens existing prompt
2. User edits title
3. User adds new tag
4. User modifies content
5. User adds new placeholder
6. User removes old placeholder
7. User validates changes
8. User saves changes
9. User verifies file updated

**Expected Results**:
- [ ] Prompt opens correctly
- [ ] Title edited
- [ ] Tag added
- [ ] Content modified
- [ ] Placeholder added
- [ ] Placeholder removed
- [ ] Validation passes
- [ ] File updated
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test editing an existing prompt
promptstack edit existing-prompt.md
# Edit title
# Add tag
# Modify content
# Add placeholder
# Remove placeholder
# Validate
# Save
# Verify file updated
```

### Scenario 3: Validation Error Resolution

**Description**: Test validation error detection and resolution workflow.

**Steps**:
1. User creates prompt with errors
2. User validates prompt
3. User sees validation errors
4. User navigates to first error
5. User fixes error
6. User navigates to next error
7. User fixes error
8. User re-validates
9. User saves prompt

**Expected Results**:
- [ ] Errors detected
- [ ] Errors displayed clearly
- [ ] Navigation works
- [ ] First error fixed
- [ ] Second error fixed
- [ ] Re-validation passes
- [ ] Prompt saved
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test validation error resolution
promptstack create
# Create prompt with errors
# Validate
# See errors
# Navigate to errors
# Fix errors
# Re-validate
# Save
```

### Scenario 4: Complex Prompt with Multiple Placeholders

**Description**: Test prompt creation with complex placeholder scenarios.

**Steps**:
1. User opens prompt creator
2. User writes complex content
3. User inserts 5 text placeholders
4. User inserts 3 list placeholders
5. User validates prompt
6. User sees no errors
7. User saves prompt
8. User reopens prompt
9. User verifies all placeholders

**Expected Results**:
- [ ] Content written
- [ ] 5 text placeholders inserted
- [ ] 3 list placeholders inserted
- [ ] Validation passes
- [ ] Prompt saved
- [ ] Prompt reopens
- [ ] All placeholders preserved
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test complex prompt with multiple placeholders
promptstack create
# Write complex content
# Insert 5 text placeholders
# Insert 3 list placeholders
# Validate
# Save
# Reopen
# Verify placeholders
```

### Scenario 5: Real-time Validation During Editing

**Description**: Test real-time validation feedback during editing.

**Steps**:
1. User opens prompt editor
2. User starts typing content
3. User sees validation errors appear
4. User continues typing
5. User sees errors update
6. User fixes errors
7. User sees errors clear
8. User completes prompt
9. User saves prompt

**Expected Results**:
- [ ] Errors appear as typing
- [ ] Errors update as typing
- [ ] No UI lag
- [ ] Errors clear when fixed
- [ ] Validation is responsive
- [ ] Prompt saved
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test real-time validation during editing
promptstack edit test.md
# Start typing
# Watch for errors
# Continue typing
# Watch errors update
# Fix errors
# Watch errors clear
# Save
```

## Performance Benchmarks

### Benchmark 1: Prompt Validation

**Test**: Measure performance of prompt validation

```go
func BenchmarkPromptValidation(b *testing.B) {
    prompt := LoadTestPrompt("complex.md")
    for i := 0; i < b.N; i++ {
        ValidatePrompt(prompt)
    }
}
```

**Thresholds**:
- [ ] Simple prompt: <10ms
- [ ] Medium prompt: <50ms
- [ ] Complex prompt: <100ms
- [ ] Very complex prompt: <200ms

### Benchmark 2: Validation Results Display

**Test**: Measure performance of validation results display

```go
func BenchmarkValidationResultsDisplay(b *testing.B) {
    results := ValidatePrompt("complex.md")
    for i := 0; i < b.N; i++ {
        DisplayValidationResults(results)
    }
}
```

**Thresholds**:
- [ ] 10 errors: <10ms
- [ ] 50 errors: <50ms
- [ ] 100 errors: <100ms
- [ ] No UI lag

### Benchmark 3: Prompt Creation

**Test**: Measure performance of prompt creation

```go
func BenchmarkPromptCreation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        CreatePrompt("test.md", content)
    }
}
```

**Thresholds**:
- [ ] Small prompt: <50ms
- [ ] Medium prompt: <100ms
- [ ] Large prompt: <200ms
- [ ] File creation included

### Benchmark 4: Prompt Editing

**Test**: Measure performance of prompt editing

```go
func BenchmarkPromptEditing(b *testing.B) {
    prompt := LoadPrompt("test.md")
    for i := 0; i < b.N; i++ {
        EditPrompt(prompt, "new content")
    }
}
```

**Thresholds**:
- [ ] Small edit: <10ms
- [ ] Medium edit: <50ms
- [ ] Large edit: <100ms
- [ ] No UI lag

### Benchmark 5: Real-time Validation

**Test**: Measure performance of real-time validation

```go
func BenchmarkRealtimeValidation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ValidateRealtime("test content")
    }
}
```

**Thresholds**:
- [ ] Small content: <10ms
- [ ] Medium content: <50ms
- [ ] Large content: <100ms
- [ ] Debouncing works

## Test Execution

### Running Integration Tests

```bash
# Run all prompt management integration tests
go test ./internal/prompt -v -tags=integration

# Run specific test
go test ./internal/prompt -v -run TestPromptValidation

# Run with coverage
go test ./internal/prompt -cover -coverprofile=coverage.out
```

### Running End-to-End Tests

```bash
# Run all E2E scenarios
./scripts/test/prompt-management-e2e.sh

# Run specific scenario
./scripts/test/prompt-management-e2e.sh scenario1

# Run with performance monitoring
./scripts/test/prompt-management-e2e.sh --perf
```

### Running Benchmarks

```bash
# Run all benchmarks
go test ./internal/prompt -bench=. -benchmem

# Run specific benchmark
go test ./internal/prompt -bench=BenchmarkPromptValidation

# Run with CPU profiling
go test ./internal/prompt -bench=. -cpuprofile=cpu.prof
```

## Test Data

### Sample Prompts

**Location**: `test/data/prompts/`

- `valid/` - Valid prompts
- `invalid-yaml/` - Invalid YAML prompts
- `missing-fields/` - Prompts with missing fields
- `invalid-placeholders/` - Prompts with invalid placeholders
- `duplicate-placeholders/` - Prompts with duplicate placeholders
- `complex/` - Complex prompts with many placeholders

### Sample Validation Results

**Location**: `test/data/validation/`

- `simple-errors.md` - Prompt with simple errors
- `complex-errors.md` - Prompt with complex errors
- `warnings.md` - Prompt with warnings
- `mixed.md` - Prompt with errors and warnings

## Success Criteria

### Integration Tests
- [ ] All integration tests pass
- [ ] Code coverage >80% for prompt management components
- [ ] No memory leaks detected
- [ ] No race conditions detected

### End-to-End Scenarios
- [ ] All scenarios complete successfully
- [ ] Validation works correctly in all scenarios
- [ ] Creation works correctly
- [ ] Editing works correctly
- [ ] Performance meets thresholds

### Performance Benchmarks
- [ ] All benchmarks meet thresholds
- [ ] No performance regression from baseline
- [ ] Memory usage remains bounded
- [ ] CPU usage is reasonable

## Known Issues and Limitations

### Current Limitations
- Very complex prompts (>100 placeholders) may cause performance issues
- Real-time validation may have latency with very large content
- Complex YAML frontmatter may have parsing edge cases
- Nested placeholders may have validation edge cases

### Future Improvements
- Add prompt templates
- Implement more sophisticated validation rules
- Add prompt autocomplete
- Optimize validation for very large prompts
- Add prompt versioning

## References

- [Enhanced Test Criteria Template](ENHANCED-TEST-CRITERIA-TEMPLATE.md)
- [Milestones Documentation](../milestones.md)
- [Go Testing Guide](../go-testing-guide.md)
- [Project Structure](../project-structure.md)
- [Foundation Testing Guide](FOUNDATION-TESTING-GUIDE.md)
- [Library Integration Testing Guide](LIBRARY-INTEGRATION-TESTING-GUIDE.md)
- [Placeholder Testing Guide](PLACEHOLDER-TESTING-GUIDE.md)