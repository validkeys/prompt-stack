# History Testing Guide (Milestones 15-17)

**Milestone Group**: History  
**Milestones**: M15-M17  
**Focus**: SQLite setup, sync, browser

## Overview

This guide provides comprehensive testing strategies for the History milestone group, which implements a dual storage system with SQLite database and markdown files. Testing focuses on ensuring reliable database operations, accurate synchronization between storage systems, and efficient history browsing.

## Integration Tests

### Test Suite: History Integration

**Location**: `internal/history/integration_test.go`

#### Test 1: SQLite Database Setup
```go
func TestSQLiteDatabaseSetup(t *testing.T) {
    // Test database initialization and schema
    // 1. Create database file
    // 2. Initialize schema
    // 3. Verify tables created
    // 4. Verify indexes created
    // 5. Test database connection
}
```

**Acceptance Criteria**:
- [ ] Database file created at correct location
- [ ] All tables created (prompts, history, tags)
- [ ] All indexes created (title, date, tags)
- [ ] FTS5 full-text search enabled
- [ ] Connection pool configured
- [ ] Database is accessible
- [ ] Setup completes in <500ms

#### Test 2: History Synchronization
```go
func TestHistorySynchronization(t *testing.T) {
    // Test sync between markdown files and database
    // 1. Create markdown file
    // 2. Trigger sync
    // 3. Verify database updated
    // 4. Modify markdown file
    // 5. Trigger sync
    // 6. Verify database updated
    // 7. Delete markdown file
    // 8. Trigger sync
    // 9. Verify database updated
}
```

**Acceptance Criteria**:
- [ ] New files added to database
- [ ] Modified files updated in database
- [ ] Deleted files removed from database
- [ ] Metadata extracted correctly
- [ ] Content indexed for search
- [ ] Sync completes in <2s for 1000 files
- [ ] No duplicates in database
- [ ] No stale data in database

#### Test 3: History Browser Navigation
```go
func TestHistoryBrowserNavigation(t *testing.T) {
    // Test history browser UI navigation
    // 1. Open history browser
    // 2. Navigate through history
    // 3. Filter by date
    // 4. Filter by tags
    // 5. Search history
    // 6. View history details
}
```

**Acceptance Criteria**:
- [ ] Browser displays all history entries
- [ ] Navigation is smooth (<100ms response)
- [ ] Date filtering works correctly
- [ ] Tag filtering works correctly
- [ ] Search works with FTS5
- [ ] Details display correctly
- [ ] Keyboard navigation works
- [ ] No UI lag with 1000+ entries

#### Test 4: Concurrent Sync Operations
```go
func TestConcurrentSyncOperations(t *testing.T) {
    // Test concurrent sync operations
    // 1. Start multiple sync operations
    // 2. Verify no race conditions
    // 3. Verify database consistency
    // 4. Verify no data loss
}
```

**Acceptance Criteria**:
- [ ] Concurrent syncs work correctly
- [ ] No race conditions detected
- [ ] Database remains consistent
- [ ] No data loss
- [ ] No deadlocks
- [ ] Performance scales with concurrency

#### Test 5: Database Error Handling
```go
func TestDatabaseErrorHandling(t *testing.T) {
    // Test database error handling
    // 1. Simulate database lock
    // 2. Simulate disk full
    // 3. Simulate corrupted database
    // 4. Verify error handling
    // 5. Verify recovery
}
```

**Acceptance Criteria**:
- [ ] Database lock handled gracefully
- [ ] Disk full error caught and reported
- [ ] Corrupted database detected
- [ ] Error messages are clear
- [ ] Recovery works correctly
- [ ] No crashes on errors
- [ ] No data corruption

#### Test 6: History Search Performance
```go
func TestHistorySearchPerformance(t *testing.T) {
    // Test search performance with large history
    // 1. Load 1000 history entries
    // 2. Perform various searches
    // 3. Measure search latency
    // 4. Verify result relevance
}
```

**Acceptance Criteria**:
- [ ] Search completes in <100ms for 1000 entries
- [ ] Results are relevant and ranked correctly
- [ ] FTS5 full-text search works
- [ ] Partial matches work
- [ ] Case-insensitive search works
- [ ] Special characters handled correctly

## End-to-End Scenarios

### Scenario 1: Initial Database Setup

**Description**: Test initial database setup and first sync.

**Steps**:
1. User starts PromptStack for first time
2. Database is created automatically
3. Schema is initialized
4. User creates first prompt file
5. Sync triggers automatically
6. User opens history browser
7. User sees first entry
8. User searches for entry

**Expected Results**:
- [ ] Database created at `~/.promptstack/history.db`
- [ ] All tables and indexes created
- [ ] FTS5 enabled
- [ ] First prompt synced to database
- [ ] Entry appears in history browser
- [ ] Search finds the entry
- [ ] No errors

**Test Script**:
```bash
#!/bin/bash
# Test initial database setup
rm -rf ~/.promptstack
promptstack
# Verify database created
# Create first prompt
# Verify sync
# Open history browser
# Verify entry
# Search for entry
```

### Scenario 2: Large History Management

**Description**: Test performance and usability with large history (1000+ entries).

**Steps**:
1. User has 1000 history entries
2. User starts PromptStack
3. Database loads
4. User opens history browser
5. User performs various searches
6. User filters by date
7. User filters by tags
8. User views details
9. User adds new entry
10. User verifies sync

**Expected Results**:
- [ ] Database loads in <2s
- [ ] Browser is responsive
- [ ] Search completes in <100ms
- [ ] Filtering works correctly
- [ ] Details display correctly
- [ ] New entry synced
- [ ] No performance degradation

**Test Script**:
```bash
#!/bin/bash
# Test large history management
# Generate 1000 test entries
./scripts/test/generate-history.sh 1000
promptstack
# Measure load time
# Test searches
# Test filters
# Test details
# Add new entry
# Verify sync
```

### Scenario 3: Sync Conflict Resolution

**Description**: Test sync when files are modified externally.

**Steps**:
1. User has history entries
2. User modifies file externally
3. User triggers sync
4. Database updates
5. User modifies file again externally
6. User triggers sync
7. Database updates again
8. User verifies consistency

**Expected Results**:
- [ ] External changes detected
- [ ] Database updated correctly
- [ ] No conflicts
- [ ] No data loss
- [ ] Consistency maintained
- [ ] No duplicates

**Test Script**:
```bash
#!/bin/bash
# Test sync conflict resolution
promptstack
# Modify file externally
# Trigger sync
# Verify database
# Modify file again
# Trigger sync
# Verify database
```

### Scenario 4: Database Recovery

**Description**: Test recovery from database corruption.

**Steps**:
1. User has history entries
2. Database becomes corrupted
3. User starts PromptStack
4. Corruption detected
5. Database rebuilt from markdown files
6. User verifies history
7. User performs searches

**Expected Results**:
- [ ] Corruption detected
- [ ] Error message displayed
- [ ] Database rebuilt successfully
- [ ] All entries restored
- [ ] Search works correctly
- [ ] No data loss

**Test Script**:
```bash
#!/bin/bash
# Test database recovery
promptstack
# Corrupt database
# Restart PromptStack
# Verify recovery
# Verify history
# Test searches
```

### Scenario 5: History Search Scenarios

**Description**: Test various search scenarios and edge cases.

**Steps**:
1. User has diverse history
2. User searches for exact match
3. User searches for partial match
4. User searches with typos
5. User searches for tags
6. User searches for combinations
7. User searches for non-existent term
8. User performs empty search

**Expected Results**:
- [ ] Exact match returns correct results
- [ ] Partial match returns relevant results
- [ ] Typos handled gracefully (FTS5)
- [ ] Tag search works
- [ ] Combination search works
- [ ] Non-existent search shows "no results"
- [ ] Empty search returns all entries

**Test Script**:
```bash
#!/bin/bash
# Test history search scenarios
promptstack
# Search: "code review" (exact)
# Search: "code rev" (partial)
# Search: "code revew" (typo)
# Search: "tag:testing" (tag)
# Search: "code review tag:testing" (combo)
# Search: "nonexistent" (no results)
# Search: "" (empty)
# Verify results
```

## Performance Benchmarks

### Benchmark 1: Database Setup

**Test**: Measure performance of database initialization

```go
func BenchmarkDatabaseSetup(b *testing.B) {
    for i := 0; i < b.N; i++ {
        SetupDatabase()
    }
}
```

**Thresholds**:
- [ ] Empty database: <100ms
- [ ] 100 entries: <200ms
- [ ] 1000 entries: <500ms
- [ ] 10000 entries: <2s

### Benchmark 2: History Sync

**Test**: Measure performance of sync operations

```go
func BenchmarkHistorySync(b *testing.B) {
    db := SetupDatabase()
    for i := 0; i < b.N; i++ {
        SyncHistory(db)
    }
}
```

**Thresholds**:
- [ ] 10 files: <100ms
- [ ] 100 files: <500ms
- [ ] 1000 files: <2s
- [ ] 10000 files: <10s

### Benchmark 3: History Search

**Test**: Measure performance of search operations

```go
func BenchmarkHistorySearch(b *testing.B) {
    db := SetupDatabase()
    LoadHistory(db, 1000)
    for i := 0; i < b.N; i++ {
        db.Search("code review")
    }
}
```

**Thresholds**:
- [ ] 100 entries: <10ms
- [ ] 1000 entries: <100ms
- [ ] 10000 entries: <500ms
- [ ] FTS5 full-text search enabled

### Benchmark 4: History Browser Rendering

**Test**: Measure performance of history browser UI rendering

```go
func BenchmarkHistoryBrowserRender(b *testing.B) {
    db := SetupDatabase()
    LoadHistory(db, 1000)
    browser := NewHistoryBrowser(db)
    for i := 0; i < b.N; i++ {
        browser.Render()
    }
}
```

**Thresholds**:
- [ ] 100 entries: <16ms (60 FPS)
- [ ] 1000 entries: <16ms (60 FPS)
- [ ] 10000 entries: <16ms (60 FPS)
- [ ] Virtual scrolling enabled for 1000+ entries

### Benchmark 5: Concurrent Operations

**Test**: Measure performance of concurrent operations

```go
func BenchmarkConcurrentOperations(b *testing.B) {
    db := SetupDatabase()
    LoadHistory(db, 1000)
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            db.Search("test")
        }
    })
}
```

**Thresholds**:
- [ ] 10 concurrent operations: <100ms
- [ ] 100 concurrent operations: <500ms
- [ ] No race conditions
- [ ] No deadlocks

## Test Execution

### Running Integration Tests

```bash
# Run all history integration tests
go test ./internal/history -v -tags=integration

# Run specific test
go test ./internal/history -v -run TestHistorySynchronization

# Run with coverage
go test ./internal/history -cover -coverprofile=coverage.out
```

### Running End-to-End Tests

```bash
# Run all E2E scenarios
./scripts/test/history-e2e.sh

# Run specific scenario
./scripts/test/history-e2e.sh scenario1

# Run with performance monitoring
./scripts/test/history-e2e.sh --perf
```

### Running Benchmarks

```bash
# Run all benchmarks
go test ./internal/history -bench=. -benchmem

# Run specific benchmark
go test ./internal/history -bench=BenchmarkHistorySync

# Run with CPU profiling
go test ./internal/history -bench=. -cpuprofile=cpu.prof
```

## Test Data

### Sample History

**Location**: `test/data/history/`

- `small/` - 10 history entries
- `medium/` - 100 history entries
- `large/` - 1000 history entries
- `with-tags/` - Entries with diverse tags
- `with-dates/` - Entries spanning different time periods

### Sample Database

**Location**: `test/data/history/databases/`

- `empty.db` - Empty database
- `small.db` - Database with 10 entries
- `medium.db` - Database with 100 entries
- `large.db` - Database with 1000 entries
- `corrupted.db` - Corrupted database (for error testing)

## Success Criteria

### Integration Tests
- [ ] All integration tests pass
- [ ] Code coverage >80% for history components
- [ ] No memory leaks detected
- [ ] No race conditions detected

### End-to-End Scenarios
- [ ] All scenarios complete successfully
- [ ] Database setup works correctly
- [ ] Sync works correctly in all scenarios
- [ ] Browser works correctly
- [ ] Performance meets thresholds

### Performance Benchmarks
- [ ] All benchmarks meet thresholds
- [ ] No performance regression from baseline
- [ ] Memory usage remains bounded
- [ ] CPU usage is reasonable

## Known Issues and Limitations

### Current Limitations
- Very large databases (>100,000 entries) may cause performance issues
- Concurrent sync operations may have edge cases
- Database recovery may be slow for very large histories
- FTS5 search may not handle all typos correctly

### Future Improvements
- Add incremental sync for better performance
- Implement database sharding for very large histories
- Add database backup and restore
- Optimize search index for very large databases
- Add parallel sync operations

## References

- [Enhanced Test Criteria Template](ENHANCED-TEST-CRITERIA-TEMPLATE.md)
- [Milestones Documentation](../milestones.md)
- [Go Testing Guide](../go-testing-guide.md)
- [Project Structure](../project-structure.md)
- [Foundation Testing Guide](FOUNDATION-TESTING-GUIDE.md)
- [Library Integration Testing Guide](LIBRARY-INTEGRATION-TESTING-GUIDE.md)
- [Database Schema](../DATABASE-SCHEMA.md)