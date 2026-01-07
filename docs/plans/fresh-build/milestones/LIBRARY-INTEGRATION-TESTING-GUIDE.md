# Library Integration Testing Guide (Milestones 7-10)

**Milestone Group**: Library Integration  
**Milestones**: M7-M10  
**Focus**: Load library, browse, search, insert prompts

## Overview

This guide provides comprehensive testing strategies for the Library Integration milestone group, which enables users to manage and interact with their prompt library. Testing focuses on ensuring reliable library loading, efficient browsing, fast fuzzy search, and seamless prompt insertion.

## Integration Tests

### Test Suite: Library Integration

**Location**: `internal/library/integration_test.go`

#### Test 1: Library Loading and Indexing
```go
func TestLibraryLoadingAndIndexing(t *testing.T) {
    // Test complete library loading workflow
    // 1. Load library from directory
    // 2. Parse all prompt files
    // 3. Build search index
    // 4. Verify all prompts indexed
    // 5. Verify metadata extracted
}
```

**Acceptance Criteria**:
- [ ] Library loads from configured directory
- [ ] All .md files are parsed
- [ ] YAML frontmatter extracted correctly
- [ ] Search index built successfully
- [ ] Metadata (title, tags, description) extracted
- [ ] Invalid files are skipped gracefully
- [ ] Loading completes in <2s for 1000 prompts

#### Test 2: Library Browser Navigation
```go
func TestLibraryBrowserNavigation(t *testing.T) {
    // Test library browser UI navigation
    // 1. Open library browser
    // 2. Navigate through prompts
    // 3. Filter by tags
    // 4. Sort by different criteria
    // 5. View prompt details
}
```

**Acceptance Criteria**:
- [ ] Browser displays all prompts
- [ ] Navigation is smooth (<100ms response)
- [ ] Tag filtering works correctly
- [ ] Sorting works (title, date, usage)
- [ ] Prompt details display correctly
- [ ] Keyboard navigation works
- [ ] No UI lag with 1000+ prompts

#### Test 3: Fuzzy Search Performance
```go
func TestFuzzySearchPerformance(t *testing.T) {
    // Test fuzzy search with large library
    // 1. Load 1000 prompts
    // 2. Perform various searches
    // 3. Measure search latency
    // 4. Verify result relevance
    // 5. Test partial matches
}
```

**Acceptance Criteria**:
- [ ] Search completes in <100ms for 1000 prompts
- [ ] Results are relevant and ranked correctly
- [ ] Partial matches work (e.g., "code rev" matches "code review")
- [ ] Case-insensitive search works
- [ ] Special characters handled correctly
- [ ] Empty search returns all prompts
- [ ] No results found handled gracefully

#### Test 4: Prompt Insertion Workflow
```go
func TestPromptInsertionWorkflow(t *testing.T) {
    // Test prompt insertion into editor
    // 1. Open editor with existing content
    // 2. Select prompt from library
    // 3. Insert prompt at cursor
    // 4. Verify content inserted correctly
    // 5. Verify placeholders preserved
}
```

**Acceptance Criteria**:
- [ ] Prompt inserted at correct position
- [ ] Existing content preserved
- [ ] Placeholders preserved ({{text:name}}, {{list:name}})
- [ ] Cursor positioned after insertion
- [ ] Undo works for insertion
- [ ] Multiple insertions work correctly
- [ ] Insertion completes in <50ms

#### Test 5: Library Reload and Sync
```go
func TestLibraryReloadAndSync(t *testing.T) {
    // Test library reload and synchronization
    // 1. Load initial library
    // 2. Add new prompt file
    // 3. Modify existing prompt
    // 4. Delete prompt file
    // 5. Trigger library reload
    // 6. Verify changes reflected
}
```

**Acceptance Criteria**:
- [ ] New prompts detected and loaded
- [ ] Modified prompts updated
- [ ] Deleted prompts removed
- [ ] Search index updated
- [ ] No duplicates in library
- [ ] Reload completes in <2s for 1000 prompts
- [ ] No crashes during reload

#### Test 6: Library Browser with Placeholders
```go
func TestLibraryBrowserWithPlaceholders(t *testing.T) {
    // Test library browser handles placeholders correctly
    // 1. Browse prompts with placeholders
    // 2. View placeholder details
    // 3. Insert prompt with placeholders
    // 4. Verify placeholders preserved
}
```

**Acceptance Criteria**:
- [ ] Placeholders displayed in preview
- [ ] Placeholder types identified (text/list)
- [ ] Placeholder names extracted
- [ ] Placeholders preserved during insertion
- [ ] Invalid placeholders handled gracefully
- [ ] Complex placeholders work correctly

## End-to-End Scenarios

### Scenario 1: First-Time Library Setup

**Description**: Simulate a new user setting up their prompt library for the first time.

**Steps**:
1. User creates library directory
2. User adds first prompt file with YAML frontmatter
3. User starts PromptStack
4. Library loads automatically
5. User opens library browser
6. User sees their prompt
7. User searches for prompt
8. User inserts prompt into editor
9. User verifies content

**Expected Results**:
- [ ] Library directory created
- [ ] Prompt file created with valid YAML
- [ ] Library loads on startup
- [ ] Prompt appears in browser
- [ ] Search finds the prompt
- [ ] Insertion works correctly
- [ ] Placeholders preserved

**Test Script**:
```bash
#!/bin/bash
# Test first-time library setup
mkdir -p ~/prompts
cat > ~/prompts/test.md << 'EOF'
---
title: Test Prompt
tags: [test, example]
description: A test prompt
---
This is a test prompt with {{text:name}} placeholder.
EOF
promptstack
# Verify library loads
# Open browser
# Search for "test"
# Insert prompt
# Verify content
```

### Scenario 2: Large Library Management

**Description**: Test performance and usability with a large prompt library (1000+ prompts).

**Steps**:
1. User has 1000 prompts in library
2. User starts PromptStack
3. Library loads and indexes
4. User opens library browser
5. User performs various searches
6. User filters by tags
7. User sorts by different criteria
8. User inserts multiple prompts
9. User adds new prompt
10. User reloads library

**Expected Results**:
- [ ] Library loads in <2s
- [ ] Browser is responsive
- [ ] Search completes in <100ms
- [ ] Filtering works correctly
- [ ] Sorting works correctly
- [ ] Insertions work correctly
- [ ] New prompt detected on reload
- [ ] No performance degradation

**Test Script**:
```bash
#!/bin/bash
# Test large library management
# Generate 1000 test prompts
./scripts/test/generate-prompts.sh 1000
promptstack
# Measure load time
# Test searches
# Test filters
# Test sorts
# Test insertions
# Add new prompt
# Reload library
# Verify performance
```

### Scenario 3: Complex Search Queries

**Description**: Test various search query patterns and edge cases.

**Steps**:
1. User has diverse prompt library
2. User searches for exact match
3. User searches for partial match
4. User searches with typos
5. User searches with special characters
6. User searches for tags
7. User searches for combinations
8. User searches for non-existent term
9. User performs empty search

**Expected Results**:
- [ ] Exact match returns correct results
- [ ] Partial match returns relevant results
- [ ] Typos handled gracefully (fuzzy matching)
- [ ] Special characters handled correctly
- [ ] Tag search works
- [ ] Combination search works
- [ ] Non-existent search shows "no results"
- [ ] Empty search returns all prompts

**Test Script**:
```bash
#!/bin/bash
# Test complex search queries
promptstack
# Search: "code review" (exact)
# Search: "code rev" (partial)
# Search: "code revew" (typo)
# Search: "code-review" (special char)
# Search: "tag:testing" (tag)
# Search: "code review tag:testing" (combo)
# Search: "nonexistent" (no results)
# Search: "" (empty)
# Verify results
```

### Scenario 4: Prompt Insertion with Placeholders

**Description**: Test prompt insertion with various placeholder types and scenarios.

**Steps**:
1. User has prompts with placeholders
2. User opens editor with content
3. User inserts prompt with text placeholder
4. User inserts prompt with list placeholder
5. User inserts prompt with multiple placeholders
6. User inserts prompt with nested placeholders
7. User inserts prompt at different positions
8. User undoes insertions
9. User redoes insertions

**Expected Results**:
- [ ] Text placeholders preserved
- [ ] List placeholders preserved
- [ ] Multiple placeholders preserved
- [ ] Nested placeholders preserved
- [ ] Insertion at different positions works
- [ ] Undo works correctly
- [ ] Redo works correctly
- [ ] No placeholder corruption

**Test Script**:
```bash
#!/bin/bash
# Test prompt insertion with placeholders
promptstack edit test.md
# Insert prompt with {{text:name}}
# Insert prompt with {{list:items}}
# Insert prompt with multiple placeholders
# Insert at beginning
# Insert at end
# Insert in middle
# Undo insertions
# Redo insertions
# Verify placeholders
```

### Scenario 5: Library Synchronization

**Description**: Test library synchronization when files change externally.

**Steps**:
1. User has library with prompts
2. User adds new prompt externally
3. User modifies existing prompt externally
4. User deletes prompt externally
5. User triggers library reload
6. User verifies changes reflected
7. User performs searches
8. User inserts prompts

**Expected Results**:
- [ ] New prompt detected
- [ ] Modified prompt updated
- [ ] Deleted prompt removed
- [ ] Search index updated
- [ ] Browser shows correct state
- [ ] Insertions work with updated library
- [ ] No duplicates
- [ ] No stale data

**Test Script**:
```bash
#!/bin/bash
# Test library synchronization
promptstack
# Add new prompt externally
# Modify existing prompt externally
# Delete prompt externally
# Trigger reload
# Verify changes
# Test searches
# Test insertions
```

## Performance Benchmarks

### Benchmark 1: Library Loading

**Test**: Measure performance of library loading and indexing

```go
func BenchmarkLibraryLoad(b *testing.B) {
    for i := 0; i < b.N; i++ {
        LoadLibrary("~/prompts")
    }
}
```

**Thresholds**:
- [ ] 100 prompts: <200ms
- [ ] 500 prompts: <1s
- [ ] 1000 prompts: <2s
- [ ] 5000 prompts: <10s

### Benchmark 2: Fuzzy Search

**Test**: Measure performance of fuzzy search operations

```go
func BenchmarkFuzzySearch(b *testing.B) {
    library := LoadLibrary("~/prompts")
    for i := 0; i < b.N; i++ {
        library.Search("code review")
    }
}
```

**Thresholds**:
- [ ] 100 prompts: <10ms
- [ ] 500 prompts: <50ms
- [ ] 1000 prompts: <100ms
- [ ] 5000 prompts: <500ms

### Benchmark 3: Library Browser Rendering

**Test**: Measure performance of library browser UI rendering

```go
func BenchmarkLibraryBrowserRender(b *testing.B) {
    library := LoadLibrary("~/prompts")
    browser := NewLibraryBrowser(library)
    for i := 0; i < b.N; i++ {
        browser.Render()
    }
}
```

**Thresholds**:
- [ ] 100 prompts: <16ms (60 FPS)
- [ ] 500 prompts: <16ms (60 FPS)
- [ ] 1000 prompts: <16ms (60 FPS)
- [ ] Virtual scrolling enabled for 1000+ prompts

### Benchmark 4: Prompt Insertion

**Test**: Measure performance of prompt insertion operations

```go
func BenchmarkPromptInsertion(b *testing.B) {
    library := LoadLibrary("~/prompts")
    editor := NewEditor()
    for i := 0; i < b.N; i++ {
        prompt := library.GetPrompt(0)
        editor.InsertPrompt(prompt)
    }
}
```

**Thresholds**:
- [ ] Small prompt (<1KB): <10ms
- [ ] Medium prompt (<10KB): <50ms
- [ ] Large prompt (<100KB): <200ms
- [ ] No UI blocking during insertion

### Benchmark 5: Library Reload

**Test**: Measure performance of library reload operations

```go
func BenchmarkLibraryReload(b *testing.B) {
    library := LoadLibrary("~/prompts")
    for i := 0; i < b.N; i++ {
        library.Reload()
    }
}
```

**Thresholds**:
- [ ] 100 prompts: <100ms
- [ ] 500 prompts: <500ms
- [ ] 1000 prompts: <2s
- [ ] 5000 prompts: <10s

## Test Execution

### Running Integration Tests

```bash
# Run all library integration tests
go test ./internal/library -v -tags=integration

# Run specific test
go test ./internal/library -v -run TestLibraryLoadingAndIndexing

# Run with coverage
go test ./internal/library -cover -coverprofile=coverage.out
```

### Running End-to-End Tests

```bash
# Run all E2E scenarios
./scripts/test/library-e2e.sh

# Run specific scenario
./scripts/test/library-e2e.sh scenario1

# Run with performance monitoring
./scripts/test/library-e2e.sh --perf
```

### Running Benchmarks

```bash
# Run all benchmarks
go test ./internal/library -bench=. -benchmem

# Run specific benchmark
go test ./internal/library -bench=BenchmarkFuzzySearch

# Run with CPU profiling
go test ./internal/library -bench=. -cpuprofile=cpu.prof
```

## Test Data

### Sample Library

**Location**: `test/data/library/`

- `small/` - 10 prompts
- `medium/` - 100 prompts
- `large/` - 1000 prompts
- `with-placeholders/` - Prompts with various placeholders
- `with-tags/` - Prompts with diverse tags

### Sample Prompts

**Location**: `test/data/library/prompts/`

- `simple.md` - Simple prompt without placeholders
- `with-text-placeholder.md` - Prompt with text placeholder
- `with-list-placeholder.md` - Prompt with list placeholder
- `with-multiple-placeholders.md` - Prompt with multiple placeholders
- `with-complex-frontmatter.md` - Prompt with complex YAML
- `invalid-frontmatter.md` - Invalid YAML (for error testing)

## Success Criteria

### Integration Tests
- [ ] All integration tests pass
- [ ] Code coverage >80% for library components
- [ ] No memory leaks detected
- [ ] No race conditions detected

### End-to-End Scenarios
- [ ] All scenarios complete successfully
- [ ] Library loads correctly in all scenarios
- [ ] Search works correctly in all scenarios
- [ ] Insertion works correctly in all scenarios
- [ ] Performance meets thresholds

### Performance Benchmarks
- [ ] All benchmarks meet thresholds
- [ ] No performance regression from baseline
- [ ] Memory usage remains bounded
- [ ] CPU usage is reasonable

## Known Issues and Limitations

### Current Limitations
- Very large libraries (>10,000 prompts) may cause performance issues
- Complex YAML frontmatter may have parsing edge cases
- Fuzzy search may not handle all typos correctly
- Library reload may be slow with many changes

### Future Improvements
- Add incremental library loading
- Implement more sophisticated fuzzy search algorithms
- Add library caching for faster reloads
- Optimize search index for very large libraries
- Add parallel library loading

## References

- [Enhanced Test Criteria Template](ENHANCED-TEST-CRITERIA-TEMPLATE.md)
- [Milestones Documentation](../milestones.md)
- [Go Testing Guide](../go-testing-guide.md)
- [Project Structure](../project-structure.md)
- [Foundation Testing Guide](FOUNDATION-TESTING-GUIDE.md)