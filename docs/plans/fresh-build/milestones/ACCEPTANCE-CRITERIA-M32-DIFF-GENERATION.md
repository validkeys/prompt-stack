# Acceptance Criteria: M32 - Diff Generation

**Milestone**: M32 - Diff Generation  
**Milestone Group**: AI Integration (M27-M33)  
**Complexity**: High  
**Status**: Draft

## Overview

Diff Generation implements unified diff generation between original and modified content, enabling users to review AI-suggested changes before applying them. This milestone is critical for providing clear, accurate, and actionable diff representations.

## Functional Requirements

### FR-1: Basic Diff Generation

**Description**: System must generate unified diff between original and modified content.

**Acceptance Criteria**:
- [ ] Diff is generated for text changes
- [ ] Diff is generated for code changes
- [ ] Diff format is unified (standard format)
- [ ] Diff includes line numbers
- [ ] Diff includes change indicators (+, -, @@)
- [ ] Diff generation completes within 100ms for 1KB content
- [ ] Diff is accurate (no false positives/negatives)
- [ ] Diff is complete (all changes included)

**Test Cases**:
```go
func TestBasicDiffGeneration(t *testing.T) {
    // Setup: Original and modified content
    // Action: Generate diff
    // Expected: Unified diff generated
    // Verify: Format is unified
    // Verify: Line numbers included
    // Verify: Change indicators correct
}
```

### FR-2: Multiple Change Detection

**Description**: System must detect and represent multiple changes in a single diff.

**Acceptance Criteria**:
- [ ] Multiple changes are detected
- [ ] Each change is represented as a separate hunk
- [ ] Hunks are ordered by line number
- [ ] Hunks are separated by blank lines
- [ ] Each hunk has correct header (@@ line, count @@)
- [ ] Diff generation completes within 200ms for 10 changes
- [ ] No changes are missed
- [ ] No duplicate changes are included

**Test Cases**:
```go
func TestMultipleChangeDetection(t *testing.T) {
    // Setup: Original and modified with 10 changes
    // Action: Generate diff
    // Expected: All changes detected
    // Verify: 10 hunks generated
    // Verify: Hunks ordered correctly
    // Verify: No changes missed
}
```

### FR-3: Insertion Detection

**Description**: System must detect and represent insertions correctly.

**Acceptance Criteria**:
- [ ] Insertions are marked with + prefix
- [ ] Inserted lines are shown in green (if colored)
- [ ] Insertion context is included (3 lines before/after)
- [ ] Insertion position is accurate
- [ ] Multiple consecutive insertions are grouped
- [ ] Empty insertions are handled
- [ ] Insertion at start/end of file works
- [ ] Insertion detection is accurate

**Test Cases**:
```go
func TestInsertionDetection(t *testing.T) {
    // Setup: Original and modified with insertions
    // Action: Generate diff
    // Expected: Insertions marked correctly
    // Verify: + prefix used
    // Verify: Context included
    // Verify: Position accurate
}
```

### FR-4: Deletion Detection

**Description**: System must detect and represent deletions correctly.

**Acceptance Criteria**:
- [ ] Deletions are marked with - prefix
- [ ] Deleted lines are shown in red (if colored)
- [ ] Deletion context is included (3 lines before/after)
- [ ] Deletion position is accurate
- [ ] Multiple consecutive deletions are grouped
- [ ] Deletion at start/end of file works
- [ ] Deletion detection is accurate
- [ ] Deleted content is preserved in diff

**Test Cases**:
```go
func TestDeletionDetection(t *testing.T) {
    // Setup: Original and modified with deletions
    // Action: Generate diff
    // Expected: Deletions marked correctly
    // Verify: - prefix used
    // Verify: Context included
    // Verify: Position accurate
}
```

### FR-5: Modification Detection

**Description**: System must detect and represent modifications (replace) correctly.

**Acceptance Criteria**:
- [ ] Modifications are represented as deletion + insertion
- [ ] Original line is shown with - prefix
- [ ] Modified line is shown with + prefix
- [ ] Modification context is included
- [ ] Modification position is accurate
- [ ] Multiple consecutive modifications are grouped
- [ ] Modification detection is accurate
- [ ] Both original and modified content are preserved

**Test Cases**:
```go
func TestModificationDetection(t *testing.T) {
    // Setup: Original and modified with replacements
    // Action: Generate diff
    // Expected: Modifications represented correctly
    // Verify: - prefix for original
    // Verify: + prefix for modified
    // Verify: Both preserved
}
```

### FR-6: Whitespace Handling

**Description**: System must handle whitespace changes correctly.

**Acceptance Criteria**:
- [ ] Whitespace-only changes are detected
- [ ] Trailing whitespace changes are detected
- [ ] Leading whitespace changes are detected
- [ ] Tab/space conversions are detected
- [ ] Empty line changes are detected
- [ ] Whitespace changes can be ignored (configurable)
- [ ] Whitespace is preserved in diff
- [ ] Whitespace handling is consistent

**Test Cases**:
```go
func TestWhitespaceHandling(t *testing.T) {
    // Setup: Original and modified with whitespace changes
    // Action: Generate diff
    // Expected: Whitespace changes detected
    // Verify: Trailing changes detected
    // Verify: Leading changes detected
    // Verify: Tab/space conversions detected
}
```

## Integration Requirements

### IR-1: Integration with AI Suggestions

**Description**: Diff generation must integrate with AI suggestion parsing.

**Acceptance Criteria**:
- [ ] AI suggestions are parsed for changes
- [ ] Original content is extracted from suggestion
- [ ] Modified content is extracted from suggestion
- [ ] Diff is generated from suggestion
- [ ] Multiple suggestions can be diffed
- [ ] Suggestion metadata is preserved
- [ ] Parsing errors are handled gracefully
- [ ] Integration completes within 200ms

**Test Cases**:
```go
func TestAISuggestionIntegration(t *testing.T) {
    // Setup: Parse AI suggestion
    // Action: Generate diff from suggestion
    // Expected: Diff generated correctly
    // Verify: Original extracted
    // Verify: Modified extracted
    // Verify: Diff accurate
}
```

### IR-2: Integration with Diff Viewer

**Description**: Diff generation must integrate with diff viewer UI.

**Acceptance Criteria**:
- [ ] Diff is passed to diff viewer
- [ ] Diff format is compatible with viewer
- [ ] Diff metadata is passed (line numbers, change types)
- [ ] Multiple diffs can be displayed
- [ ] Diff updates are reflected in viewer
- [ ] Viewer errors are handled gracefully
- [ ] Integration completes within 50ms
- [ ] Viewer can navigate diff hunks

**Test Cases**:
```go
func TestDiffViewerIntegration(t *testing.T) {
    // Setup: Generate diff
    // Action: Pass to diff viewer
    // Expected: Diff displayed correctly
    // Verify: Format compatible
    // Verify: Metadata passed
    // Verify: Navigation works
}
```

### IR-3: Integration with Editor

**Description**: Diff generation must integrate with editor for content retrieval.

**Acceptance Criteria**:
- [ ] Original content is retrieved from editor
- [ ] Modified content is retrieved from editor
- [ ] Editor state is preserved during diff generation
- [ ] Editor cursor position is considered
- [ ] Editor selection is considered
- [ ] Editor errors are handled gracefully
- [ ] Integration completes within 100ms
- [ ] Editor is not blocked during diff generation

**Test Cases**:
```go
func TestEditorIntegration(t *testing.T) {
    // Setup: Open editor with content
    // Action: Generate diff
    // Expected: Content retrieved correctly
    // Verify: Original retrieved
    // Verify: Modified retrieved
    // Verify: Editor state preserved
}
```

## Edge Cases & Error Handling

### EC-1: Empty Original Content

**Description**: System must handle empty original content gracefully.

**Acceptance Criteria**:
- [ ] Empty original is detected
- [ ] Diff shows all modified content as insertions
- [ ] No errors are raised
- [ ] Diff is generated successfully
- [ ] Diff format is correct
- [ ] Diff is complete
- [ ] Performance is not degraded
- [ ] User is notified of empty original

**Test Cases**:
```go
func TestEmptyOriginalContent(t *testing.T) {
    // Setup: Empty original, modified content
    // Action: Generate diff
    // Expected: Diff shows insertions
    // Verify: All content marked as +
    // Verify: No errors
    // Verify: Format correct
}
```

### EC-2: Empty Modified Content

**Description**: System must handle empty modified content gracefully.

**Acceptance Criteria**:
- [ ] Empty modified is detected
- [ ] Diff shows all original content as deletions
- [ ] No errors are raised
- [ ] Diff is generated successfully
- [ ] Diff format is correct
- [ ] Diff is complete
- [ ] Performance is not degraded
- [ ] User is notified of empty modified

**Test Cases**:
```go
func TestEmptyModifiedContent(t *testing.T) {
    // Setup: Original content, empty modified
    // Action: Generate diff
    // Expected: Diff shows deletions
    // Verify: All content marked as -
    // Verify: No errors
    // Verify: Format correct
}
```

### EC-3: Identical Content

**Description**: System must handle identical original and modified content.

**Acceptance Criteria**:
- [ ] Identical content is detected
- [ ] Diff is empty (no changes)
- [ ] No errors are raised
- [ ] Diff is generated successfully
- [ ] Diff format is correct
- [ ] User is notified of no changes
- [ ] Performance is optimized (early exit)
- [ ] No unnecessary processing

**Test Cases**:
```go
func TestIdenticalContent(t *testing.T) {
    // Setup: Identical original and modified
    // Action: Generate diff
    // Expected: Empty diff
    // Verify: No changes detected
    // Verify: No errors
    // Verify: User notified
}
```

### EC-4: Very Large Content

**Description**: System must handle very large content (>100KB) efficiently.

**Acceptance Criteria**:
- [ ] Large content is detected
- [ ] Diff is generated in chunks
- [ ] Progress is shown during generation
- [ ] Diff generation completes within 5 seconds for 100KB
- [ ] Memory usage remains bounded (<200MB)
- [ ] No timeout occurs
- [ ] No crash occurs
- [ ] User can cancel generation

**Test Cases**:
```go
func TestVeryLargeContent(t *testing.T) {
    // Setup: 100KB original and modified
    // Action: Generate diff
    // Expected: Diff generated successfully
    // Verify: Completes in <5s
    // Verify: Memory bounded
    // Verify: No timeout
}
```

### EC-5: Binary Content

**Description**: System must handle binary content gracefully.

**Acceptance Criteria**:
- [ ] Binary content is detected
- [ ] Error is raised for binary content
- [ ] Error message is clear
- [ ] No crash occurs
- [ ] No corruption occurs
- [ ] User is notified of binary content
- [ ] Diff is not generated for binary content
- [ ] System remains stable

**Test Cases**:
```go
func TestBinaryContent(t *testing.T) {
    // Setup: Binary original and modified
    // Action: Generate diff
    // Expected: Error raised
    // Verify: Binary detected
    // Verify: Error message clear
    // Verify: No crash
}
```

## Performance Requirements

### PR-1: Diff Generation Latency

**Description**: Diff generation must complete within specified time limits.

**Acceptance Criteria**:
- [ ] 1KB content: <100ms
- [ ] 10KB content: <500ms
- [ ] 100KB content: <5s
- [ ] 1MB content: <30s
- [ ] Latency is measured from start to completion
- [ ] Latency is consistent (±20%)
- [ ] Latency is monitored
- [ ] Latency warnings are logged

**Test Cases**:
```go
func TestDiffGenerationLatency(t *testing.T) {
    // Setup: Create test content
    // Action: Measure diff generation time
    // Expected: Completes within limits
    // Verify: 1KB <100ms
    // Verify: 10KB <500ms
    // Verify: 100KB <5s
}
```

### PR-2: Memory Usage

**Description**: Diff generation must use bounded memory.

**Acceptance Criteria**:
- [ ] Memory usage <50MB for 10KB content
- [ ] Memory usage <200MB for 100KB content
- [ ] Memory usage <500MB for 1MB content
- [ ] No memory leaks detected
- [ ] Memory is released after generation
- [ ] Memory usage is monitored
- [ ] Memory warnings are logged
- [ ] Memory limits are enforced

**Test Cases**:
```go
func TestMemoryUsage(t *testing.T) {
    // Setup: Create test content
    // Action: Monitor memory during generation
    // Expected: Memory usage bounded
    // Verify: Memory <50MB for 10KB
    // Verify: Memory <200MB for 100KB
    // Verify: No leaks
}
```

### PR-3: CPU Usage

**Description**: Diff generation must use minimal CPU resources.

**Acceptance Criteria**:
- [ ] CPU usage <30% during generation
- [ ] CPU usage <5% during idle
- [ ] CPU usage spikes are <70% and <2s duration
- [ ] CPU usage is monitored
- [ ] CPU warnings are logged
- [ ] CPU limits are enforced
- [ ] Generation is throttled if CPU >90%
- [ ] Throttling is logged

**Test Cases**:
```go
func TestCPUUsage(t *testing.T) {
    // Setup: Create test content
    // Action: Monitor CPU during generation
    // Expected: CPU usage minimal
    // Verify: CPU <30% during generation
    // Verify: CPU <5% during idle
    // Verify: Spikes <70% and <2s
}
```

## User Experience Requirements

### UX-1: Diff Visualization

**Description**: Users must see clear diff visualization.

**Acceptance Criteria**:
- [ ] Diff is displayed with color coding
- [ ] Insertions are shown in green
- [ ] Deletions are shown in red
- [ ] Modifications show both colors
- [ ] Line numbers are visible
- [ ] Change indicators are visible
- [ ] Diff is readable
- [ ] Diff is scrollable

**Test Cases**:
```go
func TestDiffVisualization(t *testing.T) {
    // Setup: Generate diff
    // Action: Display diff
    // Expected: Clear visualization
    // Verify: Color coding used
    // Verify: Insertions green
    // Verify: Deletions red
    // Verify: Line numbers visible
}
```

### UX-2: Diff Navigation

**Description**: Users must be able to navigate through diff hunks.

**Acceptance Criteria**:
- [ ] Next hunk navigation works
- [ ] Previous hunk navigation works
- [ ] Jump to specific hunk works
- [ ] Keyboard shortcuts work (n, p, j, k)
- [ ] Navigation is smooth (<50ms response)
- [ ] Current hunk is highlighted
- [ ] Hunk count is displayed
- [ ] Navigation wraps around

**Test Cases**:
```go
func TestDiffNavigation(t *testing.T) {
    // Setup: Generate diff with multiple hunks
    // Action: Navigate through hunks
    // Expected: Navigation works smoothly
    // Verify: Next hunk works
    // Verify: Previous hunk works
    // Verify: Current hunk highlighted
}
```

### UX-3: Diff Statistics

**Description**: Users must see diff statistics.

**Acceptance Criteria**:
- [ ] Total changes count is displayed
- [ ] Insertions count is displayed
- [ ] Deletions count is displayed
- [ ] Modifications count is displayed
- [ ] Statistics are accurate
- [ ] Statistics update in real-time
- [ ] Statistics are visible in header
- [ ] Statistics are clear and understandable

**Test Cases**:
```go
func TestDiffStatistics(t *testing.T) {
    // Setup: Generate diff
    // Action: Display statistics
    // Expected: Accurate statistics
    // Verify: Total changes correct
    // Verify: Insertions correct
    // Verify: Deletions correct
    // Verify: Modifications correct
}
```

## Success Criteria

### Must Have (P0)
- [ ] All functional requirements met
- [ ] All integration requirements met
- [ ] All edge cases handled
- [ ] All performance requirements met
- [ ] All user experience requirements met
- [ ] No critical bugs
- [ ] No data loss scenarios
- [ ] Code coverage >80%

### Should Have (P1)
- [ ] Diff generation completes in <5s for 100KB
- [ ] Memory usage <200MB for 100KB
- [ ] CPU usage <30% during generation
- [ ] Diff visualization is clear
- [ ] Diff navigation works smoothly
- [ ] Diff statistics are accurate

### Nice to Have (P2)
- [ ] Diff can be exported
- [ ] Diff can be shared
- [ ] Diff can be compared
- [ ] Diff can be annotated
- [ ] Diff can be filtered

## References

- [Enhanced Test Criteria Template](ENHANCED-TEST-CRITERIA-TEMPLATE.md)
- [Milestones Documentation](../milestones.md)
- [AI Integration Testing Guide](AI-INTEGRATION-TESTING-GUIDE.md)
- [Diff Application](ACCEPTANCE-CRITERIA-M33-DIFF-APPLICATION.md)