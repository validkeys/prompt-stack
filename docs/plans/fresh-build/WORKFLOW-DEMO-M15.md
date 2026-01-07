# Document Checking Workflow Demo: Milestone 15 (SQLite Setup)

**Purpose**: Demonstrate how the document checking system works when starting a milestone

---

## Step 1: Identify Milestone Context

**Milestone**: M15  
**Title**: SQLite Setup  
**Domain**: history  
**Key Features**: Database schema, CRUD operations, FTS5 search

---

## Step 2: Consult Document Reference Matrix

From [`DOCUMENT-REFERENCE-MATRIX.md`](DOCUMENT-REFERENCE-MATRIX.md):

| Milestone | Title | Required Documents | Additional Context |
|-----------|-------|-------------------|-------------------|
| **M15** | SQLite Setup | Core Planning, Implementation Guides, DATABASE-SCHEMA.md | DEPENDENCIES.md (modernc.org/sqlite) |

**Required Documents**:
1. Core Planning Documents (6 documents)
   - milestone-execution-prompt.md
   - milestones.md
   - requirements.md
   - project-structure.md
   - go-style-guide.md
   - go-testing-guide.md
2. Context-Specific Documents (2 documents)
   - DATABASE-SCHEMA.md
   - DEPENDENCIES.md

**Total**: 8 documents to read

---

## Step 3: Read Core Planning Documents (Batch 1)

```xml
<read_file>
<args>
  <file>
    <path>milestone-execution-prompt.md</path>
  </file>
  <file>
    <path>milestones.md</path>
  </file>
  <file>
    <path>requirements.md</path>
  </file>
  <file>
    <path>project-structure.md</path>
  </file>
  <file>
    <path>go-style-guide.md</path>
  </file>
</args>
</read_file>
```

**Extracted Information**:

### From milestones.md (M15 section):
- **Goal**: Initialize history database
- **Deliverables**:
  - Database schema with FTS5
  - Create at `~/.promptstack/data/history.db`
  - Basic CRUD operations
  - Index on created_at and working_directory
- **Test Criteria**:
  - [ ] Database file created on first run
  - [ ] Schema matches specification
  - [ ] Insert composition record
  - [ ] Query by ID
  - [ ] Full-text search works
  - [ ] Update existing record
  - [ ] Delete record
- **Files**: `internal/history/database.go`

### From requirements.md (History section):
- **Database Location**: `~/.promptstack/data/history.db`
- **Purpose**: Index composition history for fast searching and metadata queries
- **Data Flow**:
  - New composition created → Insert into SQLite + create markdown file
  - Auto-save triggered → Update `content`, `character_count`, `line_count`, `updated_at`
  - History browser opened → Query SQLite for list with metadata
  - Search in history → Query SQLite with full-text search
  - History cleanup → Delete from both SQLite and filesystem
- **Sync Strategy**:
  - Markdown files are source of truth
  - On startup, verify SQLite matches filesystem
  - If mismatch detected, rebuild SQLite from markdown files

### From project-structure.md (history domain):
- **Package**: `internal/history/`
- **Files**:
  - `database.go` - SQLite operations
  - `database_test.go`
  - `storage.go` - Markdown file operations
  - `storage_test.go`
  - `sync.go` - DB/file sync verification
  - `sync_test.go`
  - `search.go` - Full-text search
  - `search_test.go`
  - `cleanup.go` - History cleanup strategies
  - `cleanup_test.go`
- **Domain**: History Domain - Composition persistence and retrieval

### From go-style-guide.md (relevant patterns):
- **Error Handling**:
  ```go
  return fmt.Errorf("failed to open database: %w", err)
  ```
- **Type Design**:
  ```go
  func NewDatabase(path string) (*Database, error) {
      if path == "" {
          return nil, fmt.Errorf("database path required")
      }
      return &Database{path: path}, nil
  }
  ```
- **Dependency Injection**:
  ```go
  func NewManager(db *Database, storage *Storage) *Manager {
      return &Manager{db: db, storage: storage}
  }
  ```

---

## Step 4: Read Context-Specific Documents (Batch 2)

```xml
<read_file>
<args>
  <file>
    <path>go-testing-guide.md</path>
  </file>
  <file>
    <path>DATABASE-SCHEMA.md</path>
  </file>
  <file>
    <path>DEPENDENCIES.md</path>
  </file>
</args>
</read_file>
```

**Extracted Information**:

### From DATABASE-SCHEMA.md:
- **Database Location**: `~/.promptstack/data/history.db`
- **Tables**:
  - `compositions` - Main table with metadata
  - `composition_tags` - Tag associations
  - `compositions_fts` - FTS5 virtual table for search
- **Schema**:
  ```sql
  CREATE TABLE compositions (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      uuid TEXT UNIQUE NOT NULL,
      content TEXT NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      working_directory TEXT,
      prompt_count INTEGER DEFAULT 0,
      token_count INTEGER DEFAULT 0
  );
  ```
- **Indexes**:
  ```sql
  CREATE INDEX idx_compositions_created_at ON compositions(created_at DESC);
  CREATE INDEX idx_compositions_working_directory ON compositions(working_directory);
  CREATE INDEX idx_compositions_uuid ON compositions(uuid);
  ```
- **FTS5 Table**:
  ```sql
  CREATE VIRTUAL TABLE compositions_fts USING fts5(
      content,
      content=compositions,
      content_rowid=rowid
  );
  ```
- **Triggers**: Auto-update timestamps, FTS5 sync triggers
- **Query Patterns**: CRUD operations, search, tag operations

### From DEPENDENCIES.md:
- **SQLite Package**: `modernc.org/sqlite` v1.42.2
- **License**: BSD-3-Clause
- **Rationale**: Pure Go SQLite implementation (no CGO) for easier cross-compilation
- **Usage**: Database operations, FTS5 full-text search

### From go-testing-guide.md (database testing patterns):
- **Test Setup**:
  ```go
  func TestDatabase(t *testing.T) {
      tmp := t.TempDir()
      db, err := NewDatabase(filepath.Join(tmp, "test.db"))
      if err != nil {
          t.Fatal(err)
      }
      defer db.Close()
      // tests...
  }
  ```
- **Table-Driven Tests**:
  ```go
  tests := []struct {
      name    string
      input   Composition
      wantErr bool
  }{
      // test cases...
  }
  ```
- **Error Handling in Tests**:
  ```go
  if err != nil {
      t.Fatalf("failed to create database: %v", err)
  }
  ```

---

## Step 5: Extract and Organize Information

### Milestone Context
- **Number**: M15
- **Title**: SQLite Setup
- **Domain**: history
- **Goal**: Initialize history database

### Requirements Summary
- Create database at `~/.promptstack/data/history.db`
- Implement schema with FTS5 support
- Provide basic CRUD operations
- Create indexes for performance
- Implement triggers for auto-updates

### Architecture Context
- **Package**: `internal/history/`
- **File**: `database.go`
- **Dependencies**: `modernc.org/sqlite`
- **Integration**: Will be used by `internal/history/manager.go`

### Style Guide References
- **Error Handling**: Wrap errors with context using `%w`
- **Type Design**: Constructor pattern with validation
- **Dependency Injection**: Pass database to manager
- **Comments**: Document exported functions

### Testing Guide References
- **Test Patterns**: Table-driven tests for CRUD operations
- **Test Setup**: Use `t.TempDir()` for test databases
- **Error Handling**: Check errors in tests, use `t.Fatal()` for setup failures
- **Cleanup**: Use `defer db.Close()` for resource cleanup

### Technical Specifications
- **Schema**: compositions table with FTS5 virtual table
- **Indexes**: created_at, working_directory, uuid
- **Triggers**: Auto-update timestamps, FTS5 sync
- **Queries**: CRUD, search, tag operations

### Dependencies
- **Package**: `modernc.org/sqlite` v1.42.2
- **Purpose**: Pure Go SQLite implementation
- **License**: BSD-3-Clause

---

## Step 6: Create Implementation Plan

Based on extracted information, create task list and reference document following [`milestone-execution-prompt.md`](milestone-execution-prompt.md) format.

### Task List (M15-task-list.md)

```markdown
# Milestone 15: SQLite Setup

## Overview
- **Goal**: Initialize history database
- **Deliverables**: Database schema with FTS5, CRUD operations, indexes
- **Dependencies**: None (first milestone in history domain)

## Tasks

### Task 1: Create Database Package Structure
- **Dependencies**: None
- **Files**: `internal/history/database.go`, `internal/history/database_test.go`
- **Integration Points**: None
- **Estimated Complexity**: Low

**Description**: Create the database package with basic structure and constructor.

**Acceptance Criteria**:
- [ ] Package created with proper imports
- [ ] Database struct defined
- [ ] NewDatabase constructor implemented
- [ ] Tests verify constructor behavior

---

### Task 2: Implement Schema Creation
- **Dependencies**: Task 1
- **Files**: `internal/history/database.go`, `internal/history/database_test.go`
- **Integration Points**: None
- **Estimated Complexity**: Medium

**Description**: Create database schema with compositions table, indexes, and FTS5 virtual table.

**Acceptance Criteria**:
- [ ] compositions table created with correct schema
- [ ] composition_tags table created
- [ ] compositions_fts virtual table created
- [ ] Indexes created on created_at, working_directory, uuid
- [ ] Tests verify schema structure

---

### Task 3: Implement CRUD Operations
- **Dependencies**: Task 2
- **Files**: `internal/history/database.go`, `internal/history/database_test.go`
- **Integration Points**: None
- **Estimated Complexity**: Medium

**Description**: Implement Create, Read, Update, Delete operations for compositions.

**Acceptance Criteria**:
- [ ] Insert composition with UUID
- [ ] Query composition by ID
- [ ] Query composition by UUID
- [ ] Update composition content and metadata
- [ ] Delete composition
- [ ] Tests verify all CRUD operations

---

### Task 4: Implement Triggers
- **Dependencies**: Task 2
- **Files**: `internal/history/database.go`, `internal/history/database_test.go`
- **Integration Points**: None
- **Estimated Complexity**: Low

**Description**: Create triggers for auto-updating timestamps and syncing FTS5.

**Acceptance Criteria**:
- [ ] Trigger updates updated_at on UPDATE
- [ ] Trigger inserts into FTS5 on INSERT
- [ ] Trigger updates FTS5 on UPDATE
- [ ] Trigger deletes from FTS5 on DELETE
- [ ] Tests verify trigger behavior

---

### Task 5: Implement Full-Text Search
- **Dependencies**: Task 2
- **Files**: `internal/history/database.go`, `internal/history/database_test.go`
- **Integration Points**: None
- **Estimated Complexity**: Medium

**Description**: Implement search functionality using FTS5 virtual table.

**Acceptance Criteria**:
- [ ] Search by content using MATCH operator
- [ ] Search results ordered by created_at DESC
- [ ] Tests verify search functionality
- [ ] Tests verify search ranking

---

### Task 6: Implement Tag Operations
- **Dependencies**: Task 2
- **Files**: `internal/history/database.go`, `internal/history/database_test.go`
- **Integration Points**: None
- **Estimated Complexity**: Low

**Description**: Implement tag management operations for compositions.

**Acceptance Criteria**:
- [ ] Add tag to composition
- [ ] Remove tag from composition
- [ ] Get all tags for composition
- [ ] Get all unique tags
- [ ] Tests verify tag operations
```

### Reference Document (M15-reference.md)

```markdown
# Milestone 15 Reference Guide

## Architecture Context

### Domain Overview
History Domain - Composition persistence and retrieval using SQLite database with full-text search capabilities.

### Package Structure
```
internal/history/
├── database.go          # SQLite operations (this milestone)
├── database_test.go     # Database tests
├── storage.go           # Markdown file operations (future)
├── sync.go             # DB/file sync verification (future)
└── search.go           # Full-text search (future)
```

### Dependencies
- `modernc.org/sqlite` v1.42.2 - Pure Go SQLite implementation
- Standard library: `database/sql`, `path/filepath`, `time`

## Style Guide References

### Relevant Patterns

**Error Handling** (from go-style-guide.md):
```go
// Wrap errors with context
return fmt.Errorf("failed to open database: %w", err)

// Include context in error messages
return fmt.Errorf("failed to insert composition %s: %w", uuid, err)
```

**Type Design** (from go-style-guide.md):
```go
// Constructor pattern with validation
func NewDatabase(path string) (*Database, error) {
    if path == "" {
        return nil, fmt.Errorf("database path required")
    }
    db, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    return &Database{db: db, path: path}, nil
}
```

**Method Receivers** (from go-style-guide.md):
```go
// Use pointer receivers for mutating methods
func (d *Database) Insert(comp Composition) error { ... }

// Use value receivers for read-only methods
func (d Database) Path() string { ... }
```

### Common Pitfalls
- Don't ignore errors from database operations
- Don't forget to close database connections
- Don't use string concatenation for SQL queries (use prepared statements)
- Don't forget to handle foreign key constraints

## Testing Guide References

### Test Patterns

**Test Setup** (from go-testing-guide.md):
```go
func TestDatabase(t *testing.T) {
    tmp := t.TempDir()
    db, err := NewDatabase(filepath.Join(tmp, "test.db"))
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()
    // tests...
}
```

**Table-Driven Tests** (from go-testing-guide.md):
```go
func TestInsertComposition(t *testing.T) {
    tests := []struct {
        name    string
        input   Composition
        wantErr bool
    }{
        {
            name:  "valid composition",
            input: Composition{UUID: "test-uuid", Content: "test"},
            wantErr: false,
        },
        {
            name:    "empty uuid",
            input:   Composition{UUID: "", Content: "test"},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Test Examples

**CRUD Operations Test**:
```go
func TestCRUDOperations(t *testing.T) {
    tmp := t.TempDir()
    db, err := NewDatabase(filepath.Join(tmp, "test.db"))
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()
    
    // Create
    comp := Composition{
        UUID: "test-uuid-1",
        Content: "Test content",
        WorkingDirectory: "/test/dir",
    }
    id, err := db.Insert(comp)
    if err != nil {
        t.Fatalf("failed to insert: %v", err)
    }
    
    // Read
    got, err := db.GetByID(id)
    if err != nil {
        t.Fatalf("failed to get by ID: %v", err)
    }
    if got.UUID != comp.UUID {
        t.Errorf("got UUID %q, want %q", got.UUID, comp.UUID)
    }
    
    // Update
    comp.Content = "Updated content"
    err = db.Update(id, comp)
    if err != nil {
        t.Fatalf("failed to update: %v", err)
    }
    
    // Delete
    err = db.Delete(id)
    if err != nil {
        t.Fatalf("failed to delete: %v", err)
    }
}
```

## Implementation Notes

### Task 1: Create Database Package Structure

**Code Examples**:
```go
package history

import (
    "database/sql"
    "fmt"
    "path/filepath"
    _ "modernc.org/sqlite"
)

type Database struct {
    db   *sql.DB
    path string
}

func NewDatabase(path string) (*Database, error) {
    if path == "" {
        return nil, fmt.Errorf("database path required")
    }
    
    db, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    return &Database{db: db, path: path}, nil
}

func (d *Database) Close() error {
    return d.db.Close()
}
```

**Test Examples**:
```go
func TestNewDatabase(t *testing.T) {
    tests := []struct {
        name    string
        path    string
        wantErr bool
    }{
        {
            name:    "valid path",
            path:    filepath.Join(t.TempDir(), "test.db"),
            wantErr: false,
        },
        {
            name:    "empty path",
            path:    "",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db, err := NewDatabase(tt.path)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewDatabase() error = %v, wantErr %v", err, tt.wantErr)
            }
            if db != nil {
                db.Close()
            }
        })
    }
}
```

**Integration Considerations**:
- This is the foundation for all history domain operations
- Will be used by manager.go in future milestones
- Must handle concurrent access safely

### Task 2: Implement Schema Creation

**Code Examples**:
```go
func (d *Database) CreateSchema() error {
    // Create compositions table
    _, err := d.db.Exec(`
        CREATE TABLE IF NOT EXISTS compositions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            uuid TEXT UNIQUE NOT NULL,
            content TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            working_directory TEXT,
            prompt_count INTEGER DEFAULT 0,
            token_count INTEGER DEFAULT 0
        );
    `)
    if err != nil {
        return fmt.Errorf("failed to create compositions table: %w", err)
    }
    
    // Create indexes
    _, err = d.db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_compositions_created_at 
        ON compositions(created_at DESC);
    `)
    if err != nil {
        return fmt.Errorf("failed to create created_at index: %w", err)
    }
    
    // Create FTS5 virtual table
    _, err = d.db.Exec(`
        CREATE VIRTUAL TABLE IF NOT EXISTS compositions_fts 
        USING fts5(content, content=compositions, content_rowid=rowid);
    `)
    if err != nil {
        return fmt.Errorf("failed to create FTS5 table: %w", err)
    }
    
    return nil
}
```

**Test Examples**:
```go
func TestCreateSchema(t *testing.T) {
    tmp := t.TempDir()
    db, err := NewDatabase(filepath.Join(tmp, "test.db"))
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()
    
    err = db.CreateSchema()
    if err != nil {
        t.Fatalf("CreateSchema() error = %v", err)
    }
    
    // Verify tables exist
    var tableName string
    err = db.db.QueryRow(
        "SELECT name FROM sqlite_master WHERE type='table' AND name='compositions'",
    ).Scan(&tableName)
    if err != nil {
        t.Errorf("compositions table not found: %v", err)
    }
}
```

**Integration Considerations**:
- Schema must match DATABASE-SCHEMA.md exactly
- FTS5 table must be synced with compositions table
- Triggers will be added in Task 4

### Task 3: Implement CRUD Operations

**Code Examples**:
```go
type Composition struct {
    ID               int
    UUID             string
    Content          string
    CreatedAt        time.Time
    UpdatedAt        time.Time
    WorkingDirectory string
    PromptCount      int
    TokenCount       int
}

func (d *Database) Insert(comp Composition) (int64, error) {
    result, err := d.db.Exec(`
        INSERT INTO compositions (uuid, content, working_directory, prompt_count, token_count)
        VALUES (?, ?, ?, ?, ?);
    `, comp.UUID, comp.Content, comp.WorkingDirectory, comp.PromptCount, comp.TokenCount)
    if err != nil {
        return 0, fmt.Errorf("failed to insert composition: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to get last insert ID: %w", err)
    }
    
    return id, nil
}

func (d *Database) GetByID(id int64) (*Composition, error) {
    var comp Composition
    err := d.db.QueryRow(`
        SELECT id, uuid, content, created_at, updated_at, working_directory, prompt_count, token_count
        FROM compositions WHERE id = ?;
    `, id).Scan(&comp.ID, &comp.UUID, &comp.Content, &comp.CreatedAt, &comp.UpdatedAt, 
        &comp.WorkingDirectory, &comp.PromptCount, &comp.TokenCount)
    if err != nil {
        return nil, fmt.Errorf("failed to get composition by ID: %w", err)
    }
    
    return &comp, nil
}
```

**Test Examples**:
```go
func TestInsertAndGet(t *testing.T) {
    tmp := t.TempDir()
    db, err := NewDatabase(filepath.Join(tmp, "test.db"))
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()
    
    err = db.CreateSchema()
    if err != nil {
        t.Fatal(err)
    }
    
    comp := Composition{
        UUID: "test-uuid-1",
        Content: "Test content",
        WorkingDirectory: "/test/dir",
    }
    
    id, err := db.Insert(comp)
    if err != nil {
        t.Fatalf("Insert() error = %v", err)
    }
    
    got, err := db.GetByID(id)
    if err != nil {
        t.Fatalf("GetByID() error = %v", err)
    }
    
    if got.UUID != comp.UUID {
        t.Errorf("got UUID %q, want %q", got.UUID, comp.UUID)
    }
    if got.Content != comp.Content {
        t.Errorf("got Content %q, want %q", got.Content, comp.Content)
    }
}
```

**Integration Considerations**:
- Use prepared statements to prevent SQL injection
- Handle foreign key constraints when composition_tags is added
- Return descriptive errors for debugging

### Task 4: Implement Triggers

**Code Examples**:
```go
func (d *Database) CreateTriggers() error {
    // Auto-update updated_at trigger
    _, err := d.db.Exec(`
        CREATE TRIGGER IF NOT EXISTS update_compositions_timestamp
        AFTER UPDATE ON compositions
        FOR EACH ROW
        BEGIN
            UPDATE compositions
            SET updated_at = CURRENT_TIMESTAMP
            WHERE id = OLD.id;
        END;
    `)
    if err != nil {
        return fmt.Errorf("failed to create timestamp trigger: %w", err)
    }
    
    // FTS5 sync triggers
    _, err = d.db.Exec(`
        CREATE TRIGGER IF NOT EXISTS compositions_fts_insert
        AFTER INSERT ON compositions
        BEGIN
            INSERT INTO compositions_fts(rowid, content)
            VALUES (NEW.id, NEW.content);
        END;
    `)
    if err != nil {
        return fmt.Errorf("failed to create FTS5 insert trigger: %w", err)
    }
    
    return nil
}
```

**Test Examples**:
```go
func TestTriggers(t *testing.T) {
    tmp := t.TempDir()
    db, err := NewDatabase(filepath.Join(tmp, "test.db"))
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()
    
    err = db.CreateSchema()
    if err != nil {
        t.Fatal(err)
    }
    
    err = db.CreateTriggers()
    if err != nil {
        t.Fatal(err)
    }
    
    // Test auto-update timestamp
    comp := Composition{UUID: "test-uuid", Content: "test"}
    id, err := db.Insert(comp)
    if err != nil {
        t.Fatal(err)
    }
    
    original, err := db.GetByID(id)
    if err != nil {
        t.Fatal(err)
    }
    
    time.Sleep(10 * time.Millisecond) // Ensure time difference
    
    err = db.Update(id, Composition{Content: "updated"})
    if err != nil {
        t.Fatal(err)
    }
    
    updated, err := db.GetByID(id)
    if err != nil {
        t.Fatal(err)
    }
    
    if !updated.UpdatedAt.After(original.UpdatedAt) {
        t.Error("updated_at was not auto-updated")
    }
}
```

**Integration Considerations**:
- Triggers must be created after schema
- FTS5 triggers ensure search index stays in sync
- Test trigger behavior with time delays

### Task 5: Implement Full-Text Search

**Code Examples**:
```go
func (d *Database) Search(query string) ([]Composition, error) {
    rows, err := d.db.Query(`
        SELECT c.id, c.uuid, c.content, c.created_at, c.updated_at, c.working_directory, c.prompt_count, c.token_count
        FROM compositions c
        JOIN compositions_fts fts ON c.id = fts.rowid
        WHERE compositions_fts MATCH ?
        ORDER BY c.created_at DESC;
    `, query)
    if err != nil {
        return nil, fmt.Errorf("failed to search compositions: %w", err)
    }
    defer rows.Close()
    
    var results []Composition
    for rows.Next() {
        var comp Composition
        err := rows.Scan(&comp.ID, &comp.UUID, &comp.Content, &comp.CreatedAt, &comp.UpdatedAt,
            &comp.WorkingDirectory, &comp.PromptCount, &comp.TokenCount)
        if err != nil {
            return nil, fmt.Errorf("failed to scan composition: %w", err)
        }
        results = append(results, comp)
    }
    
    return results, nil
}
```

**Test Examples**:
```go
func TestSearch(t *testing.T) {
    tmp := t.TempDir()
    db, err := NewDatabase(filepath.Join(tmp, "test.db"))
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()
    
    err = db.CreateSchema()
    if err != nil {
        t.Fatal(err)
    }
    
    err = db.CreateTriggers()
    if err != nil {
        t.Fatal(err)
    }
    
    // Insert test data
    comp1 := Composition{UUID: "1", Content: "Write a function to sort array"}
    comp2 := Composition{UUID: "2", Content: "Create a database schema"}
    comp3 := Composition{UUID: "3", Content: "Sort the array by date"}
    
    _, err = db.Insert(comp1)
    if err != nil {
        t.Fatal(err)
    }
    _, err = db.Insert(comp2)
    if err != nil {
        t.Fatal(err)
    }
    _, err = db.Insert(comp3)
    if err != nil {
        t.Fatal(err)
    }
    
    // Search for "sort"
    results, err := db.Search("sort")
    if err != nil {
        t.Fatalf("Search() error = %v", err)
    }
    
    if len(results) != 2 {
        t.Errorf("got %d results, want 2", len(results))
    }
}
```

**Integration Considerations**:
- FTS5 MATCH operator is case-insensitive by default
- Results ordered by created_at DESC (newest first)
- Empty query should return all results

### Task 6: Implement Tag Operations

**Code Examples**:
```go
func (d *Database) AddTag(compositionID int64, tag string) error {
    _, err := d.db.Exec(`
        INSERT INTO composition_tags (composition_id, tag)
        VALUES (?, ?);
    `, compositionID, tag)
    if err != nil {
        return fmt.Errorf("failed to add tag: %w", err)
    }
    
    return nil
}

func (d *Database) GetTags(compositionID int64) ([]string, error) {
    rows, err := d.db.Query(`
        SELECT tag FROM composition_tags WHERE composition_id = ?;
    `, compositionID)
    if err != nil {
        return nil, fmt.Errorf("failed to get tags: %w", err)
    }
    defer rows.Close()
    
    var tags []string
    for rows.Next() {
        var tag string
        err := rows.Scan(&tag)
        if err != nil {
            return nil, fmt.Errorf("failed to scan tag: %w", err)
        }
        tags = append(tags, tag)
    }
    
    return tags, nil
}
```

**Test Examples**:
```go
func TestTagOperations(t *testing.T) {
    tmp := t.TempDir()
    db, err := NewDatabase(filepath.Join(tmp, "test.db"))
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()
    
    err = db.CreateSchema()
    if err != nil {
        t.Fatal(err)
    }
    
    // Insert composition
    comp := Composition{UUID: "test-uuid", Content: "test"}
    id, err := db.Insert(comp)
    if err != nil {
        t.Fatal(err)
    }
    
    // Add tags
    err = db.AddTag(id, "bug-report")
    if err != nil {
        t.Fatalf("AddTag() error = %v", err)
    }
    
    err = db.AddTag(id, "code-review")
    if err != nil {
        t.Fatalf("AddTag() error = %v", err)
    }
    
    // Get tags
    tags, err := db.GetTags(id)
    if err != nil {
        t.Fatalf("GetTags() error = %v", err)
    }
    
    if len(tags) != 2 {
        t.Errorf("got %d tags, want 2", len(tags))
    }
}
```

**Integration Considerations**:
- Tags will be used for filtering in history browser
- Foreign key constraints ensure referential integrity
- Cascade delete removes tags when composition is deleted
```

---

## Step 7: Verify Completeness

### Document Coverage
- [x] Core planning documents read (6 documents)
  - [x] milestone-execution-prompt.md
  - [x] milestones.md
  - [x] requirements.md
  - [x] project-structure.md
  - [x] go-style-guide.md
  - [x] go-testing-guide.md
- [x] Context-specific documents read (2 documents)
  - [x] DATABASE-SCHEMA.md
  - [x] DEPENDENCIES.md
- [x] All documents from reference matrix consulted
- [x] Document sections properly referenced

### Plan Quality
- [x] All deliverables from milestones.md addressed
  - [x] Database schema with FTS5
  - [x] Create at `~/.promptstack/data/history.db`
  - [x] Basic CRUD operations
  - [x] Index on created_at and working_directory
- [x] File paths match project-structure.md
  - [x] `internal/history/database.go`
  - [x] `internal/history/database_test.go`
- [x] Code follows go-style-guide.md patterns
  - [x] Error handling with %w
  - [x] Constructor pattern
  - [x] Pointer receivers for mutation
- [x] Tests follow go-testing-guide.md patterns
  - [x] Table-driven tests
  - [x] Test setup with t.TempDir()
  - [x] Error handling in tests
- [x] Technical specifications correctly applied
  - [x] Schema matches DATABASE-SCHEMA.md
  - [x] FTS5 virtual table created
  - [x] Triggers implemented
- [x] Dependencies correctly identified
  - [x] modernc.org/sqlite v1.42.2
- [x] Integration points noted
  - [x] Will be used by manager.go
  - [x] Foundation for history domain

### Plan Completeness
- [x] Task list created with clear dependencies
- [x] Reference document created with examples
- [x] Acceptance criteria are testable
- [x] File paths are explicit
- [x] Integration points are documented

---

## Summary

The document checking workflow successfully:

1. **Identified all required documents** from the reference matrix
2. **Read documents efficiently** in two batches (8 documents total)
3. **Extracted relevant information** from each document
4. **Organized information** by category (requirements, architecture, style, testing, technical specs)
5. **Created comprehensive implementation plan** with task list and reference document
6. **Verified completeness** against all checklist items

**Result**: The implementation plan for M15 is complete, well-documented, and ready for execution. All required documents were consulted, all deliverables are addressed, and all patterns from style/testing guides are applied.

---

**Key Benefits of Document Checking Workflow**:

✅ **Comprehensive Coverage**: No critical information missed  
✅ **Efficient Reading**: Batch reading minimizes context window usage  
✅ **Organized Information**: Structured extraction makes planning easier  
✅ **Pattern Application**: Style and testing patterns consistently applied  
✅ **Verification Checklist**: Ensures nothing is overlooked  
✅ **Reference Citations**: Clear links to source documents  

**This workflow ensures AI creates high-quality, well-documented implementation plans for every milestone.**