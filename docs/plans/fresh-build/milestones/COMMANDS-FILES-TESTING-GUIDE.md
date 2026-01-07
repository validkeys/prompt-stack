# Commands & Files Testing Guide (Milestones 18-22)

**Milestone Group**: Commands & Files  
**Milestones**: M18-M22  
**Focus**: Command system, palette, file references

## Overview

This guide provides comprehensive testing strategies for the Commands & Files milestone group, which implements a command system with palette UI and file reference management. Testing focuses on ensuring reliable command registration, efficient command palette, fast file finding, accurate title extraction, and seamless batch operations.

## Integration Tests

### Test Suite: Commands & Files Integration

**Location**: `internal/commands/integration_test.go`

#### Test 1: Command Registry
```go
func TestCommandRegistry(t *testing.T) {
    // Test command registration and execution
    // 1. Register commands
    // 2. List all commands
    // 3. Execute command by name
    // 4. Execute command by keybinding
    // 5. Test command aliases
}
```

**Acceptance Criteria**:
- [ ] Commands register successfully
- [ ] All commands listed correctly
- [ ] Command execution works by name
- [ ] Command execution works by keybinding
- [ ] Command aliases work
- [ ] Invalid commands handled gracefully
- [ ] Command metadata preserved

#### Test 2: Command Palette UI
```go
func TestCommandPaletteUI(t *testing.T) {
    // Test command palette interface
    // 1. Open command palette
    // 2. Search commands
    // 3. Navigate through results
    // 4. Execute selected command
    // 5. Test keyboard shortcuts
}
```

**Acceptance Criteria**:
- [ ] Palette opens quickly (<100ms)
- [ ] Search filters commands correctly
- [ ] Navigation is smooth
- [ ] Execution works
- [ ] Keyboard shortcuts work
- [ ] No UI lag with 100+ commands
- [ ] Results ranked correctly

#### Test 3: File Finder
```go
func TestFileFinder(t *testing.T) {
    // Test file finding functionality
    // 1. Search for files
    // 2. Filter by extension
    // 3. Filter by path
    // 4. Navigate through results
    // 5. Open selected file
}
```

**Acceptance Criteria**:
- [ ] File search works
- [ ] Extension filtering works
- [ ] Path filtering works
- [ ] Navigation is smooth
- [ ] File opens correctly
- [ ] Search completes in <100ms for 1000 files
- [ ] Results ranked correctly

#### Test 4: Title Extraction
```go
func TestTitleExtraction(t *testing.T) {
    // Test title extraction from files
    // 1. Extract from markdown files
    // 2. Extract from YAML frontmatter
    // 3. Extract from first heading
    // 4. Handle missing titles
    // 5. Handle special characters
}
```

**Acceptance Criteria**:
- [ ] Titles extracted from markdown
- [ ] Titles extracted from YAML frontmatter
- [ ] Titles extracted from first heading
- [ ] Missing titles handled gracefully
- [ ] Special characters handled correctly
- [ ] Extraction completes in <10ms per file
- [ ] Titles cached for performance

#### Test 5: Batch Title Editor & Link Insertion
```go
func TestBatchTitleEditorAndLinkInsertion(t *testing.T) {
    // Test batch operations
    // 1. Select multiple files
    // 2. Batch edit titles
    // 3. Insert links to files
    // 4. Verify all files updated
    // 5. Test undo/redo
}
```

**Acceptance Criteria**:
- [ ] Multiple files selected
- [ ] Batch title edit works
- [ ] Link insertion works
- [ ] All files updated correctly
- [ ] No data loss
- [ ] Undo/redo works
- [ ] Operations complete in <1s for 100 files

#### Test 6: Command Context Awareness
```go
func TestCommandContextAwareness(t *testing.T) {
    // Test commands adapt to context
    // 1. Test commands in editor
    // 2. Test commands in browser
    // 3. Test commands in palette
    // 4. Verify context-specific behavior
}
```

**Acceptance Criteria**:
- [ ] Commands work in editor context
- [ ] Commands work in browser context
- [ ] Commands work in palette context
- [ ] Context-specific behavior correct
- [ ] Invalid context handled gracefully
- [ ] Context switches work correctly

## End-to-End Scenarios

### Scenario 1: Command Registration and Execution

**Description**: Test complete command lifecycle from registration to execution.

**Steps**:
1. Developer registers new command
2. Command appears in registry
3. User opens command palette
4. User searches for command
5. User executes command
6. Command executes successfully
7. User verifies result

**Expected Results**:
- [ ] Command registered successfully
- [ ] Command appears in list
- [ ] Palette opens quickly
- [ ] Search finds command
- [ ] Execution works
- [ ] Result is correct
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test command registration and execution
promptstack
# Register test command
# Open command palette
# Search for command
# Execute command
# Verify result
```

### Scenario 2: File Finding and Opening

**Description**: Test file finding workflow with various search patterns.

**Steps**:
1. User has 1000 files in project
2. User opens file finder
3. User searches for file by name
4. User filters by extension
5. User filters by path
6. User navigates through results
7. User opens selected file
8. User verifies file content

**Expected Results**:
- [ ] Finder opens quickly
- [ ] Search finds correct files
- [ ] Extension filtering works
- [ ] Path filtering works
- [ ] Navigation is smooth
- [ ] File opens correctly
- [ ] Content is correct

**Test Script**:
```bash
#!/bin/bash
# Test file finding and opening
# Generate 1000 test files
./scripts/test/generate-files.sh 1000
promptstack
# Open file finder
# Search for file
# Filter by extension
# Filter by path
# Navigate results
# Open file
# Verify content
```

### Scenario 3: Title Extraction from Various Sources

**Description**: Test title extraction from different file types and formats.

**Steps**:
1. User has files with YAML frontmatter
2. User has files with markdown headings
3. User has files without titles
4. User has files with special characters
5. User opens file browser
6. User verifies titles displayed
7. User searches by title

**Expected Results**:
- [ ] YAML titles extracted
- [ ] Heading titles extracted
- [ ] Missing titles handled
- [ ] Special characters preserved
- [ ] Titles displayed correctly
- [ ] Search by title works
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test title extraction
# Create test files with different title formats
./scripts/test/create-title-files.sh
promptstack
# Open file browser
# Verify titles
# Search by title
```

### Scenario 4: Batch Title Editing

**Description**: Test batch title editing across multiple files.

**Steps**:
1. User has 50 files with titles
2. User selects 10 files
3. User opens batch title editor
4. User applies title transformation
5. User confirms changes
6. User verifies all files updated
7. User uses undo to revert

**Expected Results**:
- [ ] Multiple files selected
- [ ] Batch editor opens
- [ ] Transformation applied
- [ ] All files updated
- [ ] No data loss
- [ ] Undo works
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test batch title editing
# Create 50 test files
./scripts/test/generate-titled-files.sh 50
promptstack
# Select 10 files
# Open batch editor
# Apply transformation
# Confirm changes
# Verify updates
# Test undo
```

### Scenario 5: Link Insertion Workflow

**Description**: Test link insertion to reference files.

**Steps**:
1. User has editor open with content
2. User opens file finder
3. User selects file to reference
4. User inserts link at cursor
5. User verifies link format
6. User clicks link to open file
7. User verifies file opens

**Expected Results**:
- [ ] File finder opens
- [ ] File selected
- [ ] Link inserted at cursor
- [ ] Link format correct
- [ ] Link is clickable
- [ ] File opens on click
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test link insertion workflow
promptstack edit test.md
# Type some content
# Open file finder
# Select file
# Insert link
# Verify format
# Click link
# Verify file opens
```

## Performance Benchmarks

### Benchmark 1: Command Registry Operations

**Test**: Measure performance of command registry operations

```go
func BenchmarkCommandRegistry(b *testing.B) {
    registry := NewCommandRegistry()
    for i := 0; i < b.N; i++ {
        registry.RegisterCommand("test", func() {})
    }
}
```

**Thresholds**:
- [ ] Register command: <1ms
- [ ] List commands: <10ms
- [ ] Execute command: <1ms
- [ ] Search commands: <10ms

### Benchmark 2: Command Palette Search

**Test**: Measure performance of command palette search

```go
func BenchmarkCommandPaletteSearch(b *testing.B) {
    palette := NewCommandPalette()
    LoadCommands(palette, 100)
    for i := 0; i < b.N; i++ {
        palette.Search("test")
    }
}
```

**Thresholds**:
- [ ] 10 commands: <10ms
- [ ] 50 commands: <50ms
- [ ] 100 commands: <100ms
- [ ] 500 commands: <500ms

### Benchmark 3: File Finder Search

**Test**: Measure performance of file finder search

```go
func BenchmarkFileFinderSearch(b *testing.B) {
    finder := NewFileFinder()
    LoadFiles(finder, 1000)
    for i := 0; i < b.N; i++ {
        finder.Search("test")
    }
}
```

**Thresholds**:
- [ ] 100 files: <10ms
- [ ] 500 files: <50ms
- [ ] 1000 files: <100ms
- [ ] 5000 files: <500ms

### Benchmark 4: Title Extraction

**Test**: Measure performance of title extraction

```go
func BenchmarkTitleExtraction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ExtractTitle("test.md")
    }
}
```

**Thresholds**:
- [ ] Small file (<1KB): <1ms
- [ ] Medium file (<10KB): <5ms
- [ ] Large file (<100KB): <10ms
- [ ] Titles cached for performance

### Benchmark 5: Batch Operations

**Test**: Measure performance of batch operations

```go
func BenchmarkBatchOperations(b *testing.B) {
    files := LoadFiles(100)
    for i := 0; i < b.N; i++ {
        BatchEditTitles(files, "prefix")
    }
}
```

**Thresholds**:
- [ ] 10 files: <100ms
- [ ] 50 files: <500ms
- [ ] 100 files: <1s
- [ ] 500 files: <5s

## Test Execution

### Running Integration Tests

```bash
# Run all commands & files integration tests
go test ./internal/commands -v -tags=integration

# Run specific test
go test ./internal/commands -v -run TestCommandRegistry

# Run with coverage
go test ./internal/commands -cover -coverprofile=coverage.out
```

### Running End-to-End Tests

```bash
# Run all E2E scenarios
./scripts/test/commands-files-e2e.sh

# Run specific scenario
./scripts/test/commands-files-e2e.sh scenario1

# Run with performance monitoring
./scripts/test/commands-files-e2e.sh --perf
```

### Running Benchmarks

```bash
# Run all benchmarks
go test ./internal/commands -bench=. -benchmem

# Run specific benchmark
go test ./internal/commands -bench=BenchmarkCommandPaletteSearch

# Run with CPU profiling
go test ./internal/commands -bench=. -cpuprofile=cpu.prof
```

## Test Data

### Sample Commands

**Location**: `test/data/commands/`

- `basic/` - Basic commands
- `with-keybindings/` - Commands with keybindings
- `with-aliases/` - Commands with aliases
- `context-aware/` - Context-aware commands

### Sample Files

**Location**: `test/data/files/`

- `small/` - 10 files
- `medium/` - 100 files
- `large/` - 1000 files
- `with-titles/` - Files with various title formats
- `with-special-chars/` - Files with special characters

## Success Criteria

### Integration Tests
- [ ] All integration tests pass
- [ ] Code coverage >80% for commands & files components
- [ ] No memory leaks detected
- [ ] No race conditions detected

### End-to-End Scenarios
- [ ] All scenarios complete successfully
- [ ] Commands work correctly in all scenarios
- [ ] File finding works correctly
- [ ] Title extraction works correctly
- [ ] Batch operations work correctly
- [ ] Performance meets thresholds

### Performance Benchmarks
- [ ] All benchmarks meet thresholds
- [ ] No performance regression from baseline
- [ ] Memory usage remains bounded
- [ ] CPU usage is reasonable

## Known Issues and Limitations

### Current Limitations
- Very large file counts (>10,000) may cause performance issues
- Complex title extraction may have edge cases
- Batch operations may be slow for very large file sets
- Command palette may have latency with 500+ commands

### Future Improvements
- Add incremental file indexing
- Implement more sophisticated search algorithms
- Add command autocomplete
- Optimize batch operations for very large file sets
- Add command history and favorites

## References

- [Enhanced Test Criteria Template](ENHANCED-TEST-CRITERIA-TEMPLATE.md)
- [Milestones Documentation](../milestones.md)
- [Go Testing Guide](../go-testing-guide.md)
- [Project Structure](../project-structure.md)
- [Foundation Testing Guide](FOUNDATION-TESTING-GUIDE.md)
- [Library Integration Testing Guide](LIBRARY-INTEGRATION-TESTING-GUIDE.md)