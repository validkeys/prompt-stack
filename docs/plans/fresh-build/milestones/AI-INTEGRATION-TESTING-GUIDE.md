# AI Integration Testing Guide (Milestones 27-33)

**Milestone Group**: AI Integration  
**Milestones**: M27-M33  
**Focus**: API client, context selection, tokens, suggestions, diff

## Overview

This guide provides comprehensive testing strategies for the AI Integration milestone group, which implements Claude API integration with intelligent context selection, token budgeting, suggestion parsing, and diff application. Testing focuses on ensuring reliable API communication, efficient context management, accurate token estimation, robust suggestion parsing, and seamless diff operations.

## Integration Tests

### Test Suite: AI Integration

**Location**: `internal/ai/integration_test.go`

#### Test 1: Claude API Client
```go
func TestClaudeAPIClient(t *testing.T) {
    // Test API client functionality
    // 1. Send test request
    // 2. Receive valid response
    // 3. Handle authentication
    // 4. Handle rate limits
    // 5. Handle network errors
    // 6. Test retry logic
}
```

**Acceptance Criteria**:
- [ ] API requests sent successfully
- [ ] Valid responses received
- [ ] Authentication with API key works
- [ ] 401 auth errors handled
- [ ] 429 rate limits handled
- [ ] Network timeouts handled
- [ ] Retry on transient failures
- [ ] Stop after 3 failed retries
- [ ] Exponential backoff: 1s, 2s, 4s

#### Test 2: Context Selection Algorithm
```go
func TestContextSelectionAlgorithm(t *testing.T) {
    // Test intelligent context selection
    // 1. Select from library prompts
    // 2. Select from history
    // 3. Select from current file
    // 4. Apply token budget
    // 5. Verify relevance scoring
}
```

**Acceptance Criteria**:
- [ ] Library prompts selected (15% budget)
- [ ] History entries selected
- [ ] Current file content selected
- [ ] Token budget enforced (25% blocking threshold)
- [ ] Relevance scoring works
- [ ] Selection completes in <100ms
- [ ] No context exceeds budget

#### Test 3: Token Estimation & Budget
```go
func TestTokenEstimationAndBudget(t *testing.T) {
    // Test token estimation and budgeting
    // 1. Estimate tokens for text
    // 2. Estimate tokens for code
    // 3. Apply budget limits
    // 4. Track token usage
    // 5. Enforce blocking threshold
}
```

**Acceptance Criteria**:
- [ ] Token estimation accurate (±10%)
- [ ] Budget limits enforced
- [ ] Token usage tracked
- [ ] Blocking threshold enforced (25%)
- [ ] Estimation completes in <10ms
- [ ] Budget calculations correct
- [ ] No context exceeds limit

#### Test 4: Suggestion Parsing
```go
func TestSuggestionParsing(t *testing.T) {
    // Test suggestion parsing from API responses
    // 1. Parse valid suggestions
    // 2. Parse multiple suggestions
    // 3. Parse suggestions with code blocks
    // 4. Handle malformed responses
    // 5. Extract metadata
}
```

**Acceptance Criteria**:
- [ ] Valid suggestions parsed
- [ ] Multiple suggestions parsed
- [ ] Code blocks extracted
- [ ] Malformed responses handled
- [ ] Metadata extracted
- [ ] Parsing completes in <50ms
- [ ] No crashes on invalid data

#### Test 5: Suggestions Panel
```go
func TestSuggestionsPanel(t *testing.T) {
    // Test suggestions panel UI
    // 1. Display suggestions
    // 2. Navigate between suggestions
    // 3. View suggestion details
    // 4. Apply suggestion
    // 5. Dismiss suggestion
}
```

**Acceptance Criteria**:
- [ ] Suggestions displayed clearly
- [ ] Navigation works smoothly
- [ ] Details shown correctly
- [ ] Apply works
- [ ] Dismiss works
- [ ] No UI lag with 10+ suggestions
- [ ] Keyboard shortcuts work

#### Test 6: Diff Generation
```go
func TestDiffGeneration(t *testing.T) {
    // Test unified diff generation
    // 1. Generate diff for text changes
    // 2. Generate diff for code changes
    // 3. Generate diff for multiple changes
    // 4. Verify diff format
    // 5. Handle edge cases
}
```

**Acceptance Criteria**:
- [ ] Text diffs generated correctly
- [ ] Code diffs generated correctly
- [ ] Multiple changes handled
- [ ] Unified diff format correct
- [ ] Edge cases handled
- [ ] Generation completes in <100ms
- [ ] Diff is accurate

#### Test 7: Diff Application
```go
func TestDiffApplication(t *testing.T) {
    // Test diff application to editor
    // 1. Apply single change
    // 2. Apply multiple changes
    // 3. Apply with conflicts
    // 4. Undo applied diff
    // 5. Verify content updated
}
```

**Acceptance Criteria**:
- [ ] Single change applied
- [ ] Multiple changes applied
- [ ] Conflicts detected
- [ ] Undo works
- [ ] Content updated correctly
- [ ] Application completes in <50ms
- [ ] No data corruption

## End-to-End Scenarios

### Scenario 1: Basic AI Suggestion Workflow

**Description**: Test complete AI suggestion workflow from request to application.

**Steps**:
1. User opens editor with code
2. User requests AI suggestions
3. Context selected from library
4. Context selected from history
5. Context selected from current file
6. API request sent
7. Suggestions received
8. Suggestions displayed in panel
9. User reviews suggestions
10. User applies suggestion
11. Diff generated
12. Diff applied to editor

**Expected Results**:
- [ ] Context selected correctly
- [ ] API request successful
- [ ] Suggestions received
- [ ] Suggestions displayed
- [ ] Review works
- [ ] Suggestion applied
- [ ] Diff generated
- [ ] Diff applied
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test basic AI suggestion workflow
promptstack edit code.go
# Request suggestions
# Verify context selection
# Verify API request
# Verify suggestions
# Review suggestions
# Apply suggestion
# Verify diff
# Verify application
```

### Scenario 2: Token Budget Management

**Description**: Test token budget enforcement with large context.

**Steps**:
1. User has large library (1000 prompts)
2. User has large history (1000 entries)
3. User has large current file (10KB)
4. User requests AI suggestions
5. Context selection algorithm runs
6. Token budget enforced
7. Context trimmed to fit budget
8. API request sent
9. Suggestions received

**Expected Results**:
- [ ] Context selected from all sources
- [ ] Token budget enforced
- [ ] Context trimmed correctly
- [ ] No context exceeds limit
- [ ] API request successful
- [ ] Suggestions received
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test token budget management
# Generate large library
./scripts/test/generate-library.sh 1000
# Generate large history
./scripts/test/generate-history.sh 1000
# Create large file
./scripts/test/create-large-file.sh 10KB
promptstack edit large-file.go
# Request suggestions
# Verify budget enforcement
# Verify API request
# Verify suggestions
```

### Scenario 3: Error Handling and Retry Logic

**Description**: Test error handling and retry logic for API failures.

**Steps**:
1. User requests AI suggestions
2. Simulate network timeout
3. Retry with exponential backoff
4. Simulate rate limit (429)
5. Wait and retry
6. Simulate auth error (401)
7. Display error message
8. User fixes API key
9. Retry request
10. Suggestions received

**Expected Results**:
- [ ] Timeout handled
- [ ] Retry with backoff (1s, 2s, 4s)
- [ ] Rate limit handled
- [ ] Wait and retry works
- [ ] Auth error displayed
- [ ] Error message clear
- [ ] Retry after fix works
- [ ] Suggestions received
- [ ] No crashes

**Test Script**:
```bash
#!/bin/bash
# Test error handling and retry logic
promptstack edit code.go
# Request suggestions
# Simulate timeout
# Verify retry
# Simulate rate limit
# Verify wait and retry
# Simulate auth error
# Verify error message
# Fix API key
# Retry
# Verify suggestions
```

### Scenario 4: Multiple Suggestions and Diff Application

**Description**: Test handling multiple suggestions and applying diffs.

**Steps**:
1. User requests AI suggestions
2. API returns 5 suggestions
3. Suggestions displayed in panel
4. User navigates between suggestions
5. User applies suggestion 1
6. Diff generated and applied
7. User applies suggestion 3
8. Diff generated and applied
9. User undoes both changes
10. User reapplies suggestion 2

**Expected Results**:
- [ ] 5 suggestions displayed
- [ ] Navigation works
- [ ] Suggestion 1 applied
- [ ] Diff generated
- [ ] Diff applied
- [ ] Suggestion 3 applied
- [ ] Diff generated
- [ ] Diff applied
- [ ] Undo works
- [ ] Reapply works
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test multiple suggestions and diff application
promptstack edit code.go
# Request suggestions
# Verify 5 suggestions
# Navigate between suggestions
# Apply suggestion 1
# Verify diff
# Apply suggestion 3
# Verify diff
# Undo both
# Reapply suggestion 2
# Verify result
```

### Scenario 5: Complex Code Changes with Diff

**Description**: Test complex code changes with multiple diff hunks.

**Steps**:
1. User has complex code file
2. User requests AI suggestions
3. API returns suggestion with multiple changes
4. Diff generated with multiple hunks
5. User reviews each hunk
6. User applies all hunks
7. User verifies code compiles
8. User undoes changes
9. User reapplies changes

**Expected Results**:
- [ ] Complex suggestion received
- [ ] Diff with multiple hunks generated
- [ ] Each hunk displayed
- [ ] Review works
- [ ] All hunks applied
- [ ] Code compiles
- [ ] Undo works
- [ ] Reapply works
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test complex code changes with diff
promptstack edit complex-code.go
# Request suggestions
# Verify complex suggestion
# Verify diff with multiple hunks
# Review each hunk
# Apply all hunks
# Verify compilation
# Undo changes
# Reapply changes
# Verify result
```

## Performance Benchmarks

### Benchmark 1: API Request Latency

**Test**: Measure API request and response latency

```go
func BenchmarkAPIRequestLatency(b *testing.B) {
    client := NewAPIClient()
    for i := 0; i < b.N; i++ {
        client.SendRequest("test prompt")
    }
}
```

**Thresholds**:
- [ ] Simple request: <2s
- [ ] Complex request: <5s
- [ ] With retry: <10s
- [ ] Timeout: 30s

### Benchmark 2: Context Selection

**Test**: Measure performance of context selection algorithm

```go
func BenchmarkContextSelection(b *testing.B) {
    library := LoadLibrary(1000)
    history := LoadHistory(1000)
    for i := 0; i < b.N; i++ {
        SelectContext(library, history, "current file")
    }
}
```

**Thresholds**:
- [ ] Small context: <50ms
- [ ] Medium context: <100ms
- [ ] Large context: <200ms
- [ ] Very large context: <500ms

### Benchmark 3: Token Estimation

**Test**: Measure performance of token estimation

```go
func BenchmarkTokenEstimation(b *testing.B) {
    content := LoadTestContent("large-file.md")
    for i := 0; i < b.N; i++ {
        EstimateTokens(content)
    }
}
```

**Thresholds**:
- [ ] Small content: <1ms
- [ ] Medium content: <5ms
- [ ] Large content: <10ms
- [ ] Very large content: <50ms

### Benchmark 4: Suggestion Parsing

**Test**: Measure performance of suggestion parsing

```go
func BenchmarkSuggestionParsing(b *testing.B) {
    response := LoadTestResponse("complex-suggestion.json")
    for i := 0; i < b.N; i++ {
        ParseSuggestions(response)
    }
}
```

**Thresholds**:
- [ ] Simple suggestion: <10ms
- [ ] Medium suggestion: <50ms
- [ ] Complex suggestion: <100ms
- [ ] Multiple suggestions: <200ms

### Benchmark 5: Diff Generation

**Test**: Measure performance of diff generation

```go
func BenchmarkDiffGeneration(b *testing.B) {
    original := LoadTestContent("original.md")
    modified := LoadTestContent("modified.md")
    for i := 0; i < b.N; i++ {
        GenerateDiff(original, modified)
    }
}
```

**Thresholds**:
- [ ] Small changes: <10ms
- [ ] Medium changes: <50ms
- [ ] Large changes: <100ms
- [ ] Very large changes: <500ms

## Test Execution

### Running Integration Tests

```bash
# Run all AI integration tests
go test ./internal/ai -v -tags=integration

# Run specific test
go test ./internal/ai -v -run TestClaudeAPIClient

# Run with coverage
go test ./internal/ai -cover -coverprofile=coverage.out
```

### Running End-to-End Tests

```bash
# Run all E2E scenarios
./scripts/test/ai-integration-e2e.sh

# Run specific scenario
./scripts/test/ai-integration-e2e.sh scenario1

# Run with performance monitoring
./scripts/test/ai-integration-e2e.sh --perf
```

### Running Benchmarks

```bash
# Run all benchmarks
go test ./internal/ai -bench=. -benchmem

# Run specific benchmark
go test ./internal/ai -bench=BenchmarkAPIRequestLatency

# Run with CPU profiling
go test ./internal/ai -bench=. -cpuprofile=cpu.prof
```

## Test Data

### Sample API Responses

**Location**: `test/data/ai/responses/`

- `simple-suggestion.json` - Simple suggestion
- `multiple-suggestions.json` - Multiple suggestions
- `complex-suggestion.json` - Complex suggestion with code
- `error-response.json` - Error response
- `rate-limit-response.json` - Rate limit response

### Sample Contexts

**Location**: `test/data/ai/contexts/`

- `small/` - Small context
- `medium/` - Medium context
- `large/` - Large context
- `with-code/` - Context with code blocks
- `with-placeholders/` - Context with placeholders

### Sample Diffs

**Location**: `test/data/ai/diffs/`

- `simple-diff.txt` - Simple diff
- `multiple-hunks.txt` - Diff with multiple hunks
- `complex-diff.txt` - Complex diff
- `conflict-diff.txt` - Diff with conflicts

## Success Criteria

### Integration Tests
- [ ] All integration tests pass
- [ ] Code coverage >80% for AI integration components
- [ ] No memory leaks detected
- [ ] No race conditions detected

### End-to-End Scenarios
- [ ] All scenarios complete successfully
- [ ] API communication works correctly
- [ ] Context selection works correctly
- [ ] Token budgeting works correctly
- [ ] Suggestion parsing works correctly
- [ ] Diff operations work correctly
- [ ] Performance meets thresholds

### Performance Benchmarks
- [ ] All benchmarks meet thresholds
- [ ] No performance regression from baseline
- [ ] Memory usage remains bounded
- [ ] CPU usage is reasonable

## Known Issues and Limitations

### Current Limitations
- Very large contexts (>100K tokens) may cause performance issues
- Complex diffs may have edge cases
- API rate limits may affect user experience
- Token estimation may not be 100% accurate
- Retry logic may not handle all error scenarios

### Future Improvements
- Add streaming API responses
- Implement more sophisticated context selection
- Add caching for API responses
- Optimize diff generation for very large files
- Add suggestion ranking and filtering

## References

- [Enhanced Test Criteria Template](ENHANCED-TEST-CRITERIA-TEMPLATE.md)
- [Milestones Documentation](../milestones.md)
- [Go Testing Guide](../go-testing-guide.md)
- [Project Structure](../project-structure.md)
- [Foundation Testing Guide](FOUNDATION-TESTING-GUIDE.md)
- [Library Integration Testing Guide](LIBRARY-INTEGRATION-TESTING-GUIDE.md)
- [History Testing Guide](HISTORY-TESTING-GUIDE.md)