# Acceptance Criteria: M33 - Diff Application

**Milestone**: M33 - Diff Application  
**Milestone Group**: AI Integration (M27-M33)  
**Complexity**: High  
**Status**: Draft

## Overview

Diff Application implements applying unified diffs to editor content, enabling users to accept AI-suggested changes. This milestone is critical for providing reliable, safe, and reversible diff application with conflict detection.

## Functional Requirements

### FR-1: Single Change Application

**Description**: System must apply single diff changes to editor content.

**Acceptance Criteria**:
- [ ] Single insertion is applied correctly
- [ ] Single deletion is applied correctly
- [ ] Single modification is applied correctly
- [ ] Line numbers are preserved
- [ ] Content after application matches diff
- [ ] Application completes within 50ms
- [ ] No corruption occurs
- [ ] Original content is backed up

**Test Cases**:
```go
func TestSingleChangeApplication(t *testing.T) {
    // Setup: Generate diff with single change
    // Action: Apply diff to editor
    // Expected: Change applied correctly
    // Verify: Content matches diff
    // Verify: Line numbers preserved
    // Verify: No corruption
}
```

### FR-2: Multiple Change Application

**Description**: System must apply multiple diff changes to editor content.

**Acceptance Criteria**:
- [ ] Multiple insertions are applied correctly
- [ ] Multiple deletions are applied correctly
- [ ] Multiple modifications are applied correctly
- [ ] Changes are applied in order
- [ ] Line numbers are adjusted correctly
- [ ] Content after application matches diff
- [ ] Application completes within 200ms for 10 changes
- [ ] No corruption occurs
- [ ] Original content is backed up

**Test Cases**:
```go
func TestMultipleChangeApplication(t *testing.T) {
    // Setup: Generate diff with 10 changes
    // Action: Apply diff to editor
    // Expected: All changes applied correctly
    // Verify: Content matches diff
    // Verify: Order preserved
    // Verify: Line numbers adjusted
}
```

### FR-3: Partial Application

**Description**: System must allow applying only selected diff hunks.

**Acceptance Criteria**:
- [ ] User can select specific hunks to apply
- [ ] Selected hunks are applied correctly
- [ ] Unselected hunks are not applied
- [ ] Line numbers are adjusted correctly
- [ ] Content after application matches selected hunks
- [ ] Application completes within 100ms for 5 hunks
- [ ] No corruption occurs
- [ ] Original content is backed up

**Test Cases**:
```go
func TestPartialApplication(t *testing.T) {
    // Setup: Generate diff with 10 hunks
    // Action: Apply selected 5 hunks
    // Expected: Selected hunks applied correctly
    // Verify: Selected hunks applied
    // Verify: Unselected hunks not applied
    // Verify: Content matches selection
}
```

### FR-4: Conflict Detection

**Description**: System must detect conflicts between diff and current content.

**Acceptance Criteria**:
- [ ] Conflicts are detected automatically
- [ ] Conflict type is identified (insertion, deletion, modification)
- [ ] Conflict location is identified (line number)
- [ ] Conflict details are shown
- [ ] User is notified of conflicts
- [ ] Application is paused on conflict
- [ ] Conflict resolution options are presented
- [ ] No corruption occurs

**Test Cases**:
```go
func TestConflictDetection(t *testing.T) {
    // Setup: Generate diff with conflict
    // Action: Apply diff to modified content
    // Expected: Conflict detected
    // Verify: Conflict type identified
    // Verify: Conflict location identified
    // Verify: User notified
}
```

### FR-5: Conflict Resolution

**Description**: System must provide options for resolving conflicts.

**Acceptance Criteria**:
- [ ] Resolution options are presented (use diff, use current, merge)
- [ ] Use diff option applies diff change
- [ ] Use current option keeps current content
- [ ] Merge option shows both versions
- [ ] User can edit merged content
- [ ] Resolution is applied correctly
- [ ] Resolution completes within 100ms
- [ ] No corruption occurs

**Test Cases**:
```go
func TestConflictResolution(t *testing.T) {
    // Setup: Detect conflict
    // Action: Resolve conflict
    // Expected: Resolution applied correctly
    // Verify: Options presented
    // Verify: Use diff works
    // Verify: Use current works
    // Verify: Merge works
}
```

### FR-6: Undo Support

**Description**: System must support undoing applied diffs.

**Acceptance Criteria**:
- [ ] Applied diff can be undone
- [ ] Undo restores original content
- [ ] Undo works for single changes
- [ ] Undo works for multiple changes
- [ ] Undo works for partial applications
- [ ] Undo completes within 50ms
- [ ] No corruption occurs
- [ ] Undo history is maintained

**Test Cases**:
```go
func TestUndoSupport(t *testing.T) {
    // Setup: Apply diff
    // Action: Undo application
    // Expected: Original content restored
    // Verify: Content restored
    // Verify: No corruption
    // Verify: History maintained
}
```

## Integration Requirements

### IR-1: Integration with Diff Viewer

**Description**: Diff application must integrate with diff viewer UI.

**Acceptance Criteria**:
- [ ] Diff viewer provides apply button
- [ ] Diff viewer provides apply all button
- [ ] Diff viewer provides apply selected button
- [ ] Diff viewer shows application status
- [ ] Diff viewer updates after application
- [ ] Diff viewer shows conflicts
- [ ] Integration completes within 50ms
- [ ] Viewer remains responsive

**Test Cases**:
```go
func TestDiffViewerIntegration(t *testing.T) {
    // Setup: Open diff viewer
    // Action: Apply diff from viewer
    // Expected: Application succeeds
    // Verify: Apply button works
    // Verify: Apply all works
    // Verify: Apply selected works
    // Verify: Status shown
}
```

### IR-2: Integration with Editor

**Description**: Diff application must integrate with editor for content modification.

**Acceptance Criteria**:
- [ ] Editor content is retrieved
- [ ] Editor cursor position is considered
- [ ] Editor selection is considered
- [ ] Editor content is modified
- [ ] Editor cursor is repositioned
- [ ] Editor state is preserved
- [ ] Integration completes within 100ms
- [ ] Editor is not blocked

**Test Cases**:
```go
func TestEditorIntegration(t *testing.T) {
    // Setup: Open editor with content
    // Action: Apply diff to editor
    // Expected: Content modified correctly
    // Verify: Content retrieved
    // Verify: Cursor considered
    // Verify: Content modified
    // Verify: Cursor repositioned
}
```

### IR-3: Integration with Undo System

**Description**: Diff application must integrate with undo system.

**Acceptance Criteria**:
- [ ] Undo operation is recorded
- [ ] Undo history is updated
- [ ] Undo stack is maintained
- [ ] Undo can be triggered
- [ ] Undo restores correct state
- [ ] Redo works after undo
- [ ] Integration completes within 50ms
- [ ] No conflicts with other undo operations

**Test Cases**:
```go
func TestUndoSystemIntegration(t *testing.T) {
    // Setup: Apply diff
    // Action: Undo application
    // Expected: Undo works correctly
    // Verify: Operation recorded
    // Verify: History updated
    // Verify: Stack maintained
    // Verify: Redo works
}
```

## Edge Cases & Error Handling

### EC-1: Empty Diff

**Description**: System must handle empty diff gracefully.

**Acceptance Criteria**:
- [ ] Empty diff is detected
- [ ] No changes are applied
- [ ] User is notified of empty diff
- [ ] No error is raised
- [ ] No corruption occurs
- [ ] Performance is not degraded
- [ ] Editor state is preserved
- [ ] User can continue working

**Test Cases**:
```go
func TestEmptyDiff(t *testing.T) {
    // Setup: Generate empty diff
    // Action: Apply diff
    // Expected: No changes applied
    // Verify: Empty diff detected
    // Verify: User notified
    // Verify: No error raised
}
```

### EC-2: Invalid Diff Format

**Description**: System must handle invalid diff format gracefully.

**Acceptance Criteria**:
- [ ] Invalid format is detected
- [ ] Error is raised
- [ ] Error message is clear
- [ ] No changes are applied
- [ ] No corruption occurs
- [ ] User is notified of error
- [ ] Error includes line number
- [ ] Error includes expected format

**Test Cases**:
```go
func TestInvalidDiffFormat(t *testing.T) {
    // Setup: Generate invalid diff
    // Action: Apply diff
    // Expected: Error raised
    // Verify: Invalid format detected
    // Verify: Error message clear
    // Verify: No changes applied
}
```

### EC-3: Out of Range Line Numbers

**Description**: System must handle out of range line numbers gracefully.

**Acceptance Criteria**:
- [ ] Out of range is detected
- [ ] Error is raised
- [ ] Error message is clear
- [ ] No changes are applied
- [ ] No corruption occurs
- [ ] User is notified of error
- [ ] Error includes line number
- [ ] Error includes valid range

**Test Cases**:
```go
func TestOutOfRangeLineNumbers(t *testing.T) {
    // Setup: Generate diff with out of range lines
    // Action: Apply diff
    // Expected: Error raised
    // Verify: Out of range detected
    // Verify: Error message clear
    // Verify: No changes applied
}
```

### EC-4: Concurrent Applications

**Description**: System must handle concurrent diff applications safely.

**Acceptance Criteria**:
- [ ] Concurrent applications are detected
- [ ] Only one application runs at a time
- [ ] Concurrent requests are queued
- [ ] Queue is processed in order
- [ ] No race conditions occur
- [ ] No deadlocks occur
- [ ] No data corruption occurs
- [ ] All queued applications complete

**Test Cases**:
```go
func TestConcurrentApplications(t *testing.T) {
    // Setup: Start editor
    // Action: Apply 10 concurrent diffs
    // Expected: All complete safely
    // Verify: Only one at a time
    // Verify: Queue processed in order
    // Verify: No race conditions
}
```

### EC-5: Very Large Diff

**Description**: System must handle very large diffs (>1000 hunks) efficiently.

**Acceptance Criteria**:
- [ ] Large diff is detected
- [ ] Application is performed in chunks
- [ ] Progress is shown during application
- [ ] Application completes within 10 seconds for 1000 hunks
- [ ] Memory usage remains bounded (<200MB)
- [ ] No timeout occurs
- [ ] No crash occurs
- [ ] User can cancel application

**Test Cases**:
```go
func TestVeryLargeDiff(t *testing.T) {
    // Setup: Generate diff with 1000 hunks
    // Action: Apply diff
    // Expected: Application completes successfully
    // Verify: Progress shown
    // Verify: Completes in <10s
    // Verify: Memory bounded
}
```

## Performance Requirements

### PR-1: Application Latency

**Description**: Diff application must complete within specified time limits.

**Acceptance Criteria**:
- [ ] Single change: <50ms
- [ ] 10 changes: <200ms
- [ ] 100 changes: <1s
- [ ] 1000 changes: <10s
- [ ] Latency is measured from start to completion
- [ ] Latency is consistent (±20%)
- [ ] Latency is monitored
- [ ] Latency warnings are logged

**Test Cases**:
```go
func TestApplicationLatency(t *testing.T) {
    // Setup: Generate test diffs
    // Action: Measure application time
    // Expected: Completes within limits
    // Verify: Single change <50ms
    // Verify: 10 changes <200ms
    // Verify: 100 changes <1s
}
```

### PR-2: Memory Usage

**Description**: Diff application must use bounded memory.

**Acceptance Criteria**:
- [ ] Memory usage <50MB for 10 changes
- [ ] Memory usage <100MB for 100 changes
- [ ] Memory usage <200MB for 1000 changes
- [ ] No memory leaks detected
- [ ] Memory is released after application
- [ ] Memory usage is monitored
- [ ] Memory warnings are logged
- [ ] Memory limits are enforced

**Test Cases**:
```go
func TestMemoryUsage(t *testing.T) {
    // Setup: Generate test diffs
    // Action: Monitor memory during application
    // Expected: Memory usage bounded
    // Verify: Memory <50MB for 10 changes
    // Verify: Memory <100MB for 100 changes
    // Verify: No leaks
}
```

### PR-3: CPU Usage

**Description**: Diff application must use minimal CPU resources.

**Acceptance Criteria**:
- [ ] CPU usage <20% during application
- [ ] CPU usage <5% during idle
- [ ] CPU usage spikes are <50% and <1s duration
- [ ] CPU usage is monitored
- [ ] CPU warnings are logged
- [ ] CPU limits are enforced
- [ ] Application is throttled if CPU >80%
- [ ] Throttling is logged

**Test Cases**:
```go
func TestCPUUsage(t *testing.T) {
    // Setup: Generate test diffs
    // Action: Monitor CPU during application
    // Expected: CPU usage minimal
    // Verify: CPU <20% during application
    // Verify: CPU <5% during idle
    // Verify: Spikes <50% and <1s
}
```

## User Experience Requirements

### UX-1: Application Confirmation

**Description**: Users must confirm before applying diffs.

**Acceptance Criteria**:
- [ ] Confirmation dialog is shown
- [ ] Dialog shows change summary
- [ ] Dialog shows number of changes
- [ ] Dialog shows affected lines
- [ ] User can confirm or cancel
- [ ] Confirmation is required for apply all
- [ ] Confirmation is optional for apply selected
- [ ] Dialog is clear and understandable

**Test Cases**:
```go
func TestApplicationConfirmation(t *testing.T) {
    // Setup: Generate diff
    // Action: Show confirmation dialog
    // Expected: Dialog shown correctly
    // Verify: Summary shown
    // Verify: Change count shown
    // Verify: Affected lines shown
    // Verify: Confirm/cancel works
}
```

### UX-2: Progress Indication

**Description**: Users must see progress during diff application.

**Acceptance Criteria**:
- [ ] Progress bar shows application progress
- [ ] Progress percentage is accurate
- [ ] Current hunk being applied is shown
- [ ] Estimated time remaining is shown
- [ ] Progress updates every 50ms
- [ ] Progress is visible in status bar
- [ ] Progress is visible in modal
- [ ] Progress is clear and understandable

**Test Cases**:
```go
func TestProgressIndication(t *testing.T) {
    // Setup: Generate large diff
    // Action: Apply diff
    // Expected: Progress shown clearly
    // Verify: Progress bar visible
    // Verify: Percentage accurate
    // Verify: Current hunk shown
    // Verify: ETA shown
}
```

### UX-3: Error Messages

**Description**: Error messages must be clear and actionable.

**Acceptance Criteria**:
- [ ] Error messages are in plain language
- [ ] Error messages include line number
- [ ] Error messages include error type
- [ ] Error messages suggest resolution
- [ ] Error messages are logged
- [ ] Error messages are shown in modal
- [ ] Error messages can be dismissed
- [ ] Error messages are consistent

**Test Cases**:
```go
func TestErrorMessages(t *testing.T) {
    // Setup: Create error condition
    // Action: Apply diff
    // Expected: Clear error message shown
    // Verify: Message in plain language
    // Verify: Line number included
    // Verify: Error type included
    // Verify: Resolution suggested
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
- [ ] Application completes in <10s for 1000 hunks
- [ ] Memory usage <200MB for 1000 hunks
- [ ] CPU usage <20% during application
- [ ] Progress indication is clear
- [ ] Error messages are actionable
- [ ] Conflict resolution works smoothly

### Nice to Have (P2)
- [ ] Diff can be previewed before application
- [ ] Diff can be applied in background
- [ ] Diff application can be scheduled
- [ ] Diff application history is maintained
- [ ] Diff application can be automated

## References

- [Enhanced Test Criteria Template](ENHANCED-TEST-CRITERIA-TEMPLATE.md)
- [Milestones Documentation](../milestones.md)
- [AI Integration Testing Guide](AI-INTEGRATION-TESTING-GUIDE.md)
- [Diff Generation](ACCEPTANCE-CRITERIA-M32-DIFF-GENERATION.md)