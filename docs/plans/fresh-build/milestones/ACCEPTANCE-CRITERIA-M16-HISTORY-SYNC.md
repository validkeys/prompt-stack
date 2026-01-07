# Acceptance Criteria: M16 - History Sync

**Milestone**: M16 - History Sync  
**Milestone Group**: History (M15-M17)  
**Complexity**: High  
**Status**: Draft

## Overview

History Sync implements bidirectional synchronization between markdown files and SQLite database, ensuring data consistency across both storage systems. This milestone is critical for maintaining accurate history records and enabling fast searches.

## Functional Requirements

### FR-1: Initial Sync on Startup

**Description**: When PromptStack starts, it must automatically sync all markdown files to the database.

**Acceptance Criteria**:
- [ ] On first startup, all markdown files in library directory are scanned
- [ ] Each file's YAML frontmatter is parsed
- [ ] Each file's content is indexed
- [ ] Database records are created for all files
- [ ] Sync completes within 2 seconds for 100 files
- [ ] Progress indicator shows sync status
- [ ] User can cancel sync during startup
- [ ] Cancelled sync can be resumed

**Test Cases**:
```go
func TestInitialSyncOnStartup(t *testing.T) {
    // Setup: Create 100 markdown files with YAML frontmatter
    // Action: Start PromptStack
    // Expected: All files synced to database
    // Verify: Database has 100 records
    // Verify: All metadata extracted correctly
    // Verify: All content indexed
}
```

### FR-2: Incremental Sync on File Changes

**Description**: When a markdown file is modified, the database must be updated automatically.

**Acceptance Criteria**:
- [ ] File modifications are detected within 1 second
- [ ] Modified file's metadata is updated in database
- [ ] Modified file's content is re-indexed
- [ ] Sync completes within 100ms for single file
- [ ] Multiple file changes are batched
- [ ] Batch sync completes within 500ms for 10 files
- [ ] No duplicate records created
- [ ] Original record is updated, not replaced

**Test Cases**:
```go
func TestIncrementalSyncOnFileChanges(t *testing.T) {
    // Setup: Sync initial file
    // Action: Modify file content
    // Expected: Database updated within 1 second
    // Verify: Metadata updated
    // Verify: Content re-indexed
    // Verify: No duplicate records
}
```

### FR-3: Sync on New Files

**Description**: When a new markdown file is created, it must be added to the database.

**Acceptance Criteria**:
- [ ] New files are detected within 1 second
- [ ] New file's YAML frontmatter is parsed
- [ ] New file's content is indexed
- [ ] Database record is created
- [ ] Sync completes within 100ms for single file
- [ ] Multiple new files are batched
- [ ] Batch sync completes within 500ms for 10 files
- [ ] File is added to history browser

**Test Cases**:
```go
func TestSyncOnNewFiles(t *testing.T) {
    // Setup: Start PromptStack
    // Action: Create new markdown file
    // Expected: File synced to database within 1 second
    // Verify: Database has new record
    // Verify: Metadata extracted
    // Verify: Content indexed
    // Verify: File appears in history browser
}
```

### FR-4: Sync on File Deletion

**Description**: When a markdown file is deleted, it must be removed from the database.

**Acceptance Criteria**:
- [ ] File deletions are detected within 1 second
- [ ] Database record is marked as deleted
- [ ] Record is removed from search index
- [ ] File is removed from history browser
- [ ] Sync completes within 100ms for single file
- [ ] Multiple deletions are batched
- [ ] Batch sync completes within 500ms for 10 files
- [ ] Soft delete is used (record preserved for recovery)

**Test Cases**:
```go
func TestSyncOnFileDeletion(t *testing.T) {
    // Setup: Sync initial file
    // Action: Delete file
    // Expected: Record marked as deleted within 1 second
    // Verify: Record marked as deleted
    // Verify: Removed from search index
    // Verify: Removed from history browser
}
```

### FR-5: Manual Sync Trigger

**Description**: User must be able to manually trigger a sync operation.

**Acceptance Criteria**:
- [ ] Manual sync command is available
- [ ] Sync can be triggered via keyboard shortcut
- [ ] Sync can be triggered via command palette
- [ ] Sync progress is displayed
- [ ] Sync can be cancelled
- [ ] Sync results are shown (files synced, errors)
- [ ] Sync completes within 2 seconds for 100 files
- [ ] User is notified of sync completion

**Test Cases**:
```go
func TestManualSyncTrigger(t *testing.T) {
    // Setup: Start PromptStack
    // Action: Trigger manual sync
    // Expected: Sync completes and results shown
    // Verify: All files synced
    // Verify: Progress displayed
    // Verify: Results shown
    // Verify: Completion notification
}
```

### FR-6: Conflict Resolution

**Description**: When conflicts occur between file and database, they must be resolved gracefully.

**Acceptance Criteria**:
- [ ] Conflicts are detected automatically
- [ ] User is notified of conflicts
- [ ] Conflict resolution options are presented
- [ ] Options: Use file version, Use database version, Merge
- [ ] User's choice is applied
- [ ] Conflict is logged
- [ ] No data loss occurs
- [ ] Resolution completes within 1 second

**Test Cases**:
```go
func TestConflictResolution(t *testing.T) {
    // Setup: Create conflict (file modified externally)
    // Action: Trigger sync
    // Expected: Conflict detected and resolved
    // Verify: User notified
    // Verify: Options presented
    // Verify: Resolution applied
    // Verify: No data loss
}
```

## Integration Requirements

### IR-1: Integration with File System

**Description**: History sync must integrate with the file system for monitoring changes.

**Acceptance Criteria**:
- [ ] File system watcher is initialized
- [ ] Watcher monitors library directory
- [ ] Watcher detects file modifications
- [ ] Watcher detects file creations
- [ ] Watcher detects file deletions
- [ ] Watcher handles directory changes
- [ ] Watcher is resilient to temporary unavailability
- [ ] Watcher uses minimal system resources (<1% CPU)

**Test Cases**:
```go
func TestFileSystemWatcherIntegration(t *testing.T) {
    // Setup: Initialize file system watcher
    // Action: Create, modify, delete files
    // Expected: All changes detected
    // Verify: Modifications detected
    // Verify: Creations detected
    // Verify: Deletions detected
    // Verify: Resource usage minimal
}
```

### IR-2: Integration with SQLite Database

**Description**: History sync must integrate with SQLite database for storing history records.

**Acceptance Criteria**:
- [ ] Database connection is established
- [ ] Database schema is compatible
- [ ] Records are inserted correctly
- [ ] Records are updated correctly
- [ ] Records are deleted correctly
- [ ] Transactions are used for batch operations
- [ ] Connection pooling is configured
- [ ] Database errors are handled gracefully

**Test Cases**:
```go
func TestSQLiteDatabaseIntegration(t *testing.T) {
    // Setup: Initialize database connection
    // Action: Perform sync operations
    // Expected: Database operations succeed
    // Verify: Records inserted
    // Verify: Records updated
    // Verify: Records deleted
    // Verify: Transactions used
    // Verify: Errors handled
}
```

### IR-3: Integration with History Browser

**Description**: History sync must integrate with history browser for displaying records.

**Acceptance Criteria**:
- [ ] History browser receives sync updates
- [ ] New records appear in browser
- [ ] Updated records refresh in browser
- [ ] Deleted records disappear from browser
- [ ] Browser updates are debounced (100ms)
- [ ] Browser remains responsive during sync
- [ ] Browser shows sync status
- [ ] Browser filters work with synced data

**Test Cases**:
```go
func TestHistoryBrowserIntegration(t *testing.T) {
    // Setup: Open history browser
    // Action: Trigger sync
    // Expected: Browser updates correctly
    // Verify: New records appear
    // Verify: Updated records refresh
    // Verify: Deleted records disappear
    // Verify: Browser responsive
}
```

## Edge Cases & Error Handling

### EC-1: Corrupted YAML Frontmatter

**Description**: System must handle files with corrupted YAML frontmatter.

**Acceptance Criteria**:
- [ ] Corrupted YAML is detected
- [ ] Error is logged
- [ ] File is skipped gracefully
- [ ] Sync continues with other files
- [ ] User is notified of error
- [ ] Error message includes file path
- [ ] Error message includes parse error
- [ ] No crash occurs

**Test Cases**:
```go
func TestCorruptedYAMLFrontmatter(t *testing.T) {
    // Setup: Create file with corrupted YAML
    // Action: Trigger sync
    // Expected: Error handled gracefully
    // Verify: Error logged
    // Verify: File skipped
    // Verify: Sync continues
    // Verify: User notified
}
```

### EC-2: Missing YAML Frontmatter

**Description**: System must handle files without YAML frontmatter.

**Acceptance Criteria**:
- [ ] Missing YAML is detected
- [ ] Default metadata is used
- [ ] File is synced with defaults
- [ ] Title is extracted from first heading
- [ ] Tags are empty
- [ ] Description is empty
- [ ] No error is logged
- [ ] Sync continues normally

**Test Cases**:
```go
func TestMissingYAMLFrontmatter(t *testing.T) {
    // Setup: Create file without YAML
    // Action: Trigger sync
    // Expected: File synced with defaults
    // Verify: Default metadata used
    // Verify: Title extracted from heading
    // Verify: Tags empty
    // Verify: No error logged
}
```

### EC-3: Large File Sync

**Description**: System must handle large files (>10MB) during sync.

**Acceptance Criteria**:
- [ ] Large files are detected
- [ ] Sync is performed in chunks
- [ ] Progress is shown during sync
- [ ] Sync completes within 10 seconds for 10MB file
- [ ] Memory usage remains bounded (<100MB)
- [ ] No timeout occurs
- [ ] No crash occurs
- [ ] User can cancel sync

**Test Cases**:
```go
func TestLargeFileSync(t *testing.T) {
    // Setup: Create 10MB markdown file
    // Action: Trigger sync
    // Expected: Sync completes successfully
    // Verify: Progress shown
    // Verify: Sync completes in <10s
    // Verify: Memory bounded
    // Verify: No timeout
}
```

### EC-4: Concurrent Sync Operations

**Description**: System must handle concurrent sync operations safely.

**Acceptance Criteria**:
- [ ] Concurrent syncs are detected
- [ ] Only one sync runs at a time
- [ ] Concurrent requests are queued
- [ ] Queue is processed in order
- [ ] No race conditions occur
- [ ] No deadlocks occur
- [ ] No data corruption occurs
- [ ] All queued syncs complete

**Test Cases**:
```go
func TestConcurrentSyncOperations(t *testing.T) {
    // Setup: Start PromptStack
    // Action: Trigger 10 concurrent syncs
    // Expected: All syncs complete safely
    // Verify: Only one sync at a time
    // Verify: Queue processed in order
    // Verify: No race conditions
    // Verify: No data corruption
}
```

### EC-5: Database Lock Errors

**Description**: System must handle database lock errors gracefully.

**Acceptance Criteria**:
- [ ] Database lock is detected
- [ ] Error is logged
- [ ] Sync is retried with exponential backoff
- [ ] Retry attempts: 3 with 1s, 2s, 4s delays
- [ ] User is notified of retry
- [ ] If all retries fail, user is notified
- [ ] No crash occurs
- [ ] No data corruption occurs

**Test Cases**:
```go
func TestDatabaseLockErrors(t *testing.T) {
    // Setup: Simulate database lock
    // Action: Trigger sync
    // Expected: Retry with backoff
    // Verify: Lock detected
    // Verify: Retry with 1s, 2s, 4s
    // Verify: User notified
    // Verify: No crash
}
```

## Performance Requirements

### PR-1: Sync Latency

**Description**: Sync operations must complete within specified time limits.

**Acceptance Criteria**:
- [ ] Single file sync: <100ms
- [ ] 10 files sync: <500ms
- [ ] 100 files sync: <2s
- [ ] 1000 files sync: <10s
- [ ] Incremental sync: <100ms
- [ ] Manual sync: <2s for 100 files
- [ ] Startup sync: <2s for 100 files
- [ ] Latency measured from start to completion

**Test Cases**:
```go
func TestSyncLatency(t *testing.T) {
    // Setup: Create test files
    // Action: Measure sync time
    // Expected: Sync completes within limits
    // Verify: Single file <100ms
    // Verify: 10 files <500ms
    // Verify: 100 files <2s
}
```

### PR-2: Memory Usage

**Description**: Sync operations must use bounded memory.

**Acceptance Criteria**:
- [ ] Memory usage <50MB for 100 files
- [ ] Memory usage <100MB for 1000 files
- [ ] Memory usage <200MB for 10000 files
- [ ] No memory leaks detected
- [ ] Memory is released after sync
- [ ] Memory usage is monitored
- [ ] Memory warnings are logged
- [ ] Memory limits are enforced

**Test Cases**:
```go
func TestMemoryUsage(t *testing.T) {
    // Setup: Create test files
    // Action: Monitor memory during sync
    // Expected: Memory usage bounded
    // Verify: Memory <50MB for 100 files
    // Verify: Memory <100MB for 1000 files
    // Verify: No leaks
}
```

### PR-3: CPU Usage

**Description**: Sync operations must use minimal CPU resources.

**Acceptance Criteria**:
- [ ] CPU usage <10% during sync
- [ ] CPU usage <5% during idle
- [ ] CPU usage spikes are <50% and <1s duration
- [ ] CPU usage is monitored
- [ ] CPU warnings are logged
- [ ] CPU limits are enforced
- [ ] Sync is throttled if CPU >80%
- [ ] Throttling is logged

**Test Cases**:
```go
func TestCPUUsage(t *testing.T) {
    // Setup: Create test files
    // Action: Monitor CPU during sync
    // Expected: CPU usage minimal
    // Verify: CPU <10% during sync
    // Verify: CPU <5% during idle
    // Verify: Spikes <50% and <1s
}
```

## User Experience Requirements

### UX-1: Progress Indication

**Description**: Users must see clear progress indication during sync operations.

**Acceptance Criteria**:
- [ ] Progress bar shows sync progress
- [ ] Progress percentage is accurate
- [ ] Current file being synced is shown
- [ ] Estimated time remaining is shown
- [ ] Progress updates every 100ms
- [ ] Progress is visible in status bar
- [ ] Progress is visible in modal
- [ ] Progress is clear and understandable

**Test Cases**:
```go
func TestProgressIndication(t *testing.T) {
    // Setup: Create test files
    // Action: Trigger sync
    // Expected: Progress shown clearly
    // Verify: Progress bar visible
    // Verify: Percentage accurate
    // Verify: Current file shown
    // Verify: ETA shown
}
```

### UX-2: Error Messages

**Description**: Error messages must be clear and actionable.

**Acceptance Criteria**:
- [ ] Error messages are in plain language
- [ ] Error messages include file path
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
    // Action: Trigger sync
    // Expected: Clear error message shown
    // Verify: Message in plain language
    // Verify: File path included
    // Verify: Error type included
    // Verify: Resolution suggested
}
```

### UX-3: Cancel Operation

**Description**: Users must be able to cancel sync operations.

**Acceptance Criteria**:
- [ ] Cancel button is visible during sync
- [ ] Cancel button is accessible via keyboard
- [ ] Cancel stops sync immediately
- [ ] Partial sync is preserved
- [ ] Cancel is confirmed
- [ ] Cancel is logged
- [ ] User is notified of cancellation
- [ ] Cancel can be retried

**Test Cases**:
```go
func TestCancelOperation(t *testing.T) {
    // Setup: Start sync
    // Action: Cancel sync
    // Expected: Sync stops cleanly
    // Verify: Cancel button visible
    // Verify: Sync stops immediately
    // Verify: Partial sync preserved
    // Verify: User notified
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
- [ ] Sync completes within 2s for 100 files
- [ ] Memory usage <100MB for 1000 files
- [ ] CPU usage <10% during sync
- [ ] Progress indication is clear
- [ ] Error messages are actionable
- [ ] Cancel operation works smoothly

### Nice to Have (P2)
- [ ] Sync can be scheduled
- [ ] Sync can be configured
- [ ] Sync history is maintained
- [ ] Sync statistics are available
- [ ] Sync can be debugged

## References

- [Enhanced Test Criteria Template](ENHANCED-TEST-CRITERIA-TEMPLATE.md)
- [Milestones Documentation](../milestones.md)
- [History Testing Guide](HISTORY-TESTING-GUIDE.md)
- [Database Schema](../DATABASE-SCHEMA.md)