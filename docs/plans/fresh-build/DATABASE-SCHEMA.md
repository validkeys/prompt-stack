# Database Schema

**Version**: 1.0  
**Database Engine**: SQLite 3 with FTS5 extension  
**Location**: `~/.promptstack/data/history.db`

---

## Overview

The PromptStack history database stores all compositions (user-created prompts) with full-text search capabilities, metadata tracking, and efficient querying. The database uses SQLite with FTS5 (Full-Text Search) extension for fast content search.

**Key Features:**
- Dual storage: SQLite for search/metadata + Markdown files for content
- Full-text search with FTS5
- Automatic timestamp tracking
- Tag-based organization
- Working directory tracking
- Token count tracking for AI context management

---

## Database Location

```
~/.promptstack/data/history.db
```

The database is created automatically on first run (Milestone 15). The parent directory structure is created if it doesn't exist.

---

## Tables

### compositions

Stores metadata and references to composition content.

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

**Columns:**
- `id` - Auto-incrementing primary key
- `uuid` - Unique identifier (UUID v4 format) for cross-referencing with markdown files
- `content` - Full composition text (mirrored from markdown file)
- `created_at` - Timestamp when composition was first created
- `updated_at` - Timestamp of last modification
- `working_directory` - Current working directory when composition was created
- `prompt_count` - Number of prompts used in composition
- `token_count` - Estimated token count for AI context management

**Indexes:**
```sql
CREATE INDEX idx_compositions_created_at ON compositions(created_at DESC);
CREATE INDEX idx_compositions_working_directory ON compositions(working_directory);
CREATE INDEX idx_compositions_uuid ON compositions(uuid);
```

**Constraints:**
- `uuid` must be unique (enforced by UNIQUE constraint)
- `content` cannot be NULL
- `prompt_count` and `token_count` default to 0

---

### composition_tags

Stores tags associated with compositions for categorization and filtering.

```sql
CREATE TABLE composition_tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    composition_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
    FOREIGN KEY (composition_id) REFERENCES compositions(id) ON DELETE CASCADE
);
```

**Columns:**
- `id` - Auto-incrementing primary key
- `composition_id` - Foreign key reference to compositions table
- `tag` - Tag name (e.g., "bug-report", "code-review", "documentation")

**Indexes:**
```sql
CREATE INDEX idx_composition_tags_composition_id ON composition_tags(composition_id);
CREATE INDEX idx_composition_tags_tag ON composition_tags(tag);
```

**Constraints:**
- `composition_id` must reference valid composition (enforced by FOREIGN KEY)
- Cascade delete: removing a composition removes all its tags
- `tag` cannot be NULL

**Relationship:**
- One-to-many: One composition can have multiple tags
- Many-to-many: Multiple compositions can share the same tag

---

## Full-Text Search (FTS5)

### compositions_fts

Virtual table for full-text search across composition content.

```sql
CREATE VIRTUAL TABLE compositions_fts USING fts5(
    content,
    content=compositions,
    content_rowid=rowid
);
```

**Configuration:**
- External content table: `compositions`
- Rowid mapping: `content_rowid=rowid`
- Searchable column: `content`

**Usage:**
```sql
-- Simple search
SELECT c.* FROM compositions c
JOIN compositions_fts fts ON c.id = fts.rowid
WHERE compositions_fts MATCH 'search query';

-- Boolean search
SELECT c.* FROM compositions c
JOIN compositions_fts fts ON c.id = fts.rowid
WHERE compositions_fts MATCH 'bug AND fix';

-- Phrase search
SELECT c.* FROM compositions c
JOIN compositions_fts fts ON c.id = fts.rowid
WHERE compositions_fts MATCH '"exact phrase"';
```

---

## Triggers

### Auto-update updated_at

Automatically updates the `updated_at` timestamp when a composition is modified.

```sql
CREATE TRIGGER update_compositions_timestamp
AFTER UPDATE ON compositions
FOR EACH ROW
BEGIN
    UPDATE compositions
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
```

### FTS5 Sync Triggers

Keep the FTS5 virtual table in sync with the compositions table.

```sql
-- Insert trigger
CREATE TRIGGER compositions_fts_insert
AFTER INSERT ON compositions
BEGIN
    INSERT INTO compositions_fts(rowid, content)
    VALUES (NEW.id, NEW.content);
END;

-- Update trigger
CREATE TRIGGER compositions_fts_update
AFTER UPDATE ON compositions
BEGIN
    UPDATE compositions_fts
    SET content = NEW.content
    WHERE rowid = NEW.id;
END;

-- Delete trigger
CREATE TRIGGER compositions_fts_delete
AFTER DELETE ON compositions
BEGIN
    DELETE FROM compositions_fts WHERE rowid = OLD.id;
END;
```

---

## Query Patterns

### Basic CRUD Operations

#### Create Composition
```sql
INSERT INTO compositions (uuid, content, working_directory, prompt_count, token_count)
VALUES (?, ?, ?, ?, ?);
```

#### Read Composition by ID
```sql
SELECT * FROM compositions WHERE id = ?;
```

#### Read Composition by UUID
```sql
SELECT * FROM compositions WHERE uuid = ?;
```

#### Update Composition
```sql
UPDATE compositions
SET content = ?, prompt_count = ?, token_count = ?
WHERE id = ?;
```

#### Delete Composition
```sql
DELETE FROM compositions WHERE id = ?;
-- Tags automatically deleted via CASCADE
```

---

### Listing and Filtering

#### List All Compositions (Newest First)
```sql
SELECT * FROM compositions
ORDER BY created_at DESC;
```

#### List by Working Directory
```sql
SELECT * FROM compositions
WHERE working_directory = ?
ORDER BY created_at DESC;
```

#### List by Date Range
```sql
SELECT * FROM compositions
WHERE created_at BETWEEN ? AND ?
ORDER BY created_at DESC;
```

#### List with Tags
```sql
SELECT c.*, GROUP_CONCAT(ct.tag, ',') as tags
FROM compositions c
LEFT JOIN composition_tags ct ON c.id = ct.composition_id
GROUP BY c.id
ORDER BY c.created_at DESC;
```

---

### Search Operations

#### Full-Text Search
```sql
SELECT c.* FROM compositions c
JOIN compositions_fts fts ON c.id = fts.rowid
WHERE compositions_fts MATCH ?
ORDER BY c.created_at DESC;
```

#### Search with Tag Filter
```sql
SELECT DISTINCT c.* FROM compositions c
JOIN compositions_fts fts ON c.id = fts.rowid
JOIN composition_tags ct ON c.id = ct.composition_id
WHERE compositions_fts MATCH ? AND ct.tag = ?
ORDER BY c.created_at DESC;
```

#### Search by Multiple Tags
```sql
SELECT c.* FROM compositions c
WHERE c.id IN (
    SELECT composition_id FROM composition_tags
    WHERE tag IN (?, ?, ?)
    GROUP BY composition_id
    HAVING COUNT(DISTINCT tag) = ?
)
ORDER BY c.created_at DESC;
```

---

### Tag Operations

#### Add Tag to Composition
```sql
INSERT INTO composition_tags (composition_id, tag)
VALUES (?, ?);
```

#### Remove Tag from Composition
```sql
DELETE FROM composition_tags
WHERE composition_id = ? AND tag = ?;
```

#### Get All Tags for Composition
```sql
SELECT tag FROM composition_tags
WHERE composition_id = ?;
```

#### Get All Tags (Unique)
```sql
SELECT DISTINCT tag FROM composition_tags
ORDER BY tag;
```

#### Get Popular Tags
```sql
SELECT tag, COUNT(*) as count
FROM composition_tags
GROUP BY tag
ORDER BY count DESC
LIMIT 20;
```

---

### Statistics and Analytics

#### Count Total Compositions
```sql
SELECT COUNT(*) FROM compositions;
```

#### Count by Working Directory
```sql
SELECT working_directory, COUNT(*) as count
FROM compositions
GROUP BY working_directory
ORDER BY count DESC;
```

#### Average Token Count
```sql
SELECT AVG(token_count) as avg_tokens
FROM compositions;
```

#### Recent Activity
```sql
SELECT 
    DATE(created_at) as date,
    COUNT(*) as count
FROM compositions
WHERE created_at >= date('now', '-7 days')
GROUP BY DATE(created_at)
ORDER BY date DESC;
```

---

## Migration Strategy

### Version Management

Track schema version in database:

```sql
CREATE TABLE schema_version (
    version INTEGER PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

-- Insert initial version
INSERT INTO schema_version (version, applied_at, description)
VALUES (1, CURRENT_TIMESTAMP, 'Initial schema with compositions and tags');
```

### Migration Process

1. Check current schema version
2. Apply migrations in order
3. Update schema version
4. Verify migration success

### Migration Files

Store migrations in `internal/storage/migrations/`:

```
migrations/
├── 001_initial_schema.sql
├── 001_initial_schema_rollback.sql
├── 002_add_indexes.sql
├── 002_add_indexes_rollback.sql
├── 003_add_fts5.sql
├── 003_add_fts5_rollback.sql
└── 004_add_composition_tags.sql
```

### Rollback Support

Each migration should have rollback:

```sql
-- 001_initial_schema.sql
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

-- 001_initial_schema_rollback.sql
DROP TABLE IF EXISTS compositions;
```

### Migration Implementation

The migration system uses a Go migrator to apply changes:

```go
// internal/storage/migrations/migrator.go
package migrations

import (
    "database/sql"
    "fmt"
)

type Migrator struct {
    db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
    return &Migrator{db: db}
}

func (m *Migrator) Migrate() error {
    // Create schema_version table if not exists
    if err := m.createSchemaVersionTable(); err != nil {
        return fmt.Errorf("failed to create schema_version table: %w", err)
    }
    
    // Get current version
    currentVersion, err := m.getCurrentVersion()
    if err != nil {
        return fmt.Errorf("failed to get current version: %w", err)
    }
    
    // Get available migrations
    migrations, err := m.getAvailableMigrations()
    if err != nil {
        return fmt.Errorf("failed to get available migrations: %w", err)
    }
    
    // Apply pending migrations
    for _, migration := range migrations {
        if migration.Version > currentVersion {
            if err := m.applyMigration(migration); err != nil {
                return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
            }
        }
    }
    
    return nil
}

func (m *Migrator) createSchemaVersionTable() error {
    _, err := m.db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_version (
            version INTEGER PRIMARY KEY,
            applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            description TEXT
        )
    `)
    return err
}

func (m *Migrator) getCurrentVersion() (int, error) {
    var version int
    err := m.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&version)
    return version, err
}

func (m *Migrator) applyMigration(migration Migration) error {
    // Start transaction
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Execute migration SQL
    if _, err := tx.Exec(migration.SQL); err != nil {
        return fmt.Errorf("failed to execute migration: %w", err)
    }
    
    // Update schema version
    if _, err := tx.Exec(
        "INSERT INTO schema_version (version, applied_at, description) VALUES (?, CURRENT_TIMESTAMP, ?)",
        migration.Version,
        migration.Description,
    ); err != nil {
        return fmt.Errorf("failed to update schema version: %w", err)
    }
    
    // Commit transaction
    if err := tx.Commit(); err != nil {
        return err
    }
    
    return nil
}

type Migration struct {
    Version     int
    SQL         string
    Name        string
    Description string
}
```

### Migration Files

#### 001_initial_schema.sql

```sql
-- Create compositions table
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

-- Create indexes
CREATE INDEX idx_compositions_created_at ON compositions(created_at DESC);
CREATE INDEX idx_compositions_working_directory ON compositions(working_directory);
CREATE INDEX idx_compositions_uuid ON compositions(uuid);

-- Insert initial schema version
INSERT INTO schema_version (version, applied_at, description)
VALUES (1, CURRENT_TIMESTAMP, 'Initial schema with compositions table');
```

#### 001_initial_schema_rollback.sql

```sql
DROP INDEX IF EXISTS idx_compositions_uuid;
DROP INDEX IF EXISTS idx_compositions_working_directory;
DROP INDEX IF EXISTS idx_compositions_created_at;
DROP TABLE IF EXISTS compositions;
```

#### 002_add_fts5.sql

```sql
-- Create FTS5 virtual table for full-text search
CREATE VIRTUAL TABLE compositions_fts USING fts5(
    content,
    content=compositions,
    content_rowid=rowid
);

-- Populate FTS5 table
INSERT INTO compositions_fts(rowid, content)
SELECT rowid, content FROM compositions;

-- Create triggers to keep FTS5 in sync
CREATE TRIGGER compositions_fts_insert AFTER INSERT ON compositions BEGIN
    INSERT INTO compositions_fts(rowid, content)
    VALUES (NEW.id, NEW.content);
END;

CREATE TRIGGER compositions_fts_delete AFTER DELETE ON compositions BEGIN
    DELETE FROM compositions_fts WHERE rowid = OLD.id;
END;

CREATE TRIGGER compositions_fts_update AFTER UPDATE ON compositions BEGIN
    UPDATE compositions_fts SET content = NEW.content WHERE rowid = NEW.id;
END;

-- Update schema version
INSERT INTO schema_version (version, applied_at, description)
VALUES (2, CURRENT_TIMESTAMP, 'Add FTS5 full-text search support');
```

#### 002_add_fts5_rollback.sql

```sql
DROP TRIGGER IF EXISTS compositions_fts_update;
DROP TRIGGER IF EXISTS compositions_fts_delete;
DROP TRIGGER IF EXISTS compositions_fts_insert;
DROP TABLE IF EXISTS compositions_fts;
```

#### 003_add_composition_tags.sql

```sql
-- Create composition_tags table
CREATE TABLE composition_tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    composition_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
    FOREIGN KEY (composition_id) REFERENCES compositions(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_composition_tags_composition_id ON composition_tags(composition_id);
CREATE INDEX idx_composition_tags_tag ON composition_tags(tag);

-- Update schema version
INSERT INTO schema_version (version, applied_at, description)
VALUES (3, CURRENT_TIMESTAMP, 'Add composition tags support');
```

#### 003_add_composition_tags_rollback.sql

```sql
DROP INDEX IF EXISTS idx_composition_tags_tag;
DROP INDEX IF EXISTS idx_composition_tags_composition_id;
DROP TABLE IF EXISTS composition_tags;
```

### Cross-Database Migration

#### SQLite to PostgreSQL

For future scalability, support migration from SQLite to PostgreSQL:

```go
// internal/storage/migrations/sqlite_to_postgres.go
package migrations

import (
    "database/sql"
    "fmt"
)

type SQLiteToPostgresMigrator struct {
    sqliteDB    *sql.DB
    postgresDB  *sql.DB
}

func NewSQLiteToPostgresMigrator(sqliteDB, postgresDB *sql.DB) *SQLiteToPostgresMigrator {
    return &SQLiteToPostgresMigrator{
        sqliteDB:   sqliteDB,
        postgresDB: postgresDB,
    }
}

func (m *SQLiteToPostgresMigrator) Migrate() error {
    // Create PostgreSQL schema
    if err := m.createPostgreSQLSchema(); err != nil {
        return fmt.Errorf("failed to create PostgreSQL schema: %w", err)
    }
    
    // Migrate compositions
    if err := m.migrateCompositions(); err != nil {
        return fmt.Errorf("failed to migrate compositions: %w", err)
    }
    
    // Migrate tags
    if err := m.migrateTags(); err != nil {
        return fmt.Errorf("failed to migrate tags: %w", err)
    }
    
    return nil
}

func (m *SQLiteToPostgresMigrator) createPostgreSQLSchema() error {
    // Create PostgreSQL schema
    schema := `
        CREATE TABLE IF NOT EXISTS compositions (
            id SERIAL PRIMARY KEY,
            uuid TEXT UNIQUE NOT NULL,
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            working_directory TEXT,
            prompt_count INTEGER DEFAULT 0,
            token_count INTEGER DEFAULT 0
        );
        
        CREATE INDEX IF NOT EXISTS idx_compositions_created_at ON compositions(created_at DESC);
        CREATE INDEX IF NOT EXISTS idx_compositions_working_directory ON compositions(working_directory);
        CREATE INDEX IF NOT EXISTS idx_compositions_uuid ON compositions(uuid);
        
        CREATE TABLE IF NOT EXISTS composition_tags (
            id SERIAL PRIMARY KEY,
            composition_id INTEGER NOT NULL,
            tag TEXT NOT NULL,
            FOREIGN KEY (composition_id) REFERENCES compositions(id) ON DELETE CASCADE
        );
        
        CREATE INDEX IF NOT EXISTS idx_composition_tags_composition_id ON composition_tags(composition_id);
        CREATE INDEX IF NOT EXISTS idx_composition_tags_tag ON composition_tags(tag);
    `
    
    _, err := m.postgresDB.Exec(schema)
    return err
}

func (m *SQLiteToPostgresMigrator) migrateCompositions() error {
    rows, err := m.sqliteDB.Query(`
        SELECT id, uuid, content, created_at, updated_at, working_directory,
               prompt_count, token_count
        FROM compositions
    `)
    if err != nil {
        return err
    }
    defer rows.Close()
    
    stmt, err := m.postgresDB.Prepare(`
        INSERT INTO compositions (id, uuid, content, created_at, updated_at,
                                working_directory, prompt_count, token_count)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for rows.Next() {
        var comp struct {
            ID               int
            UUID             string
            Content          string
            CreatedAt        string
            UpdatedAt        string
            WorkingDirectory string
            PromptCount      int
            TokenCount       int
        }
        
        if err := rows.Scan(
            &comp.ID, &comp.UUID, &comp.Content, &comp.CreatedAt, &comp.UpdatedAt,
            &comp.WorkingDirectory, &comp.PromptCount, &comp.TokenCount,
        ); err != nil {
            return err
        }
        
        if _, err := stmt.Exec(
            comp.ID, comp.UUID, comp.Content, comp.CreatedAt, comp.UpdatedAt,
            comp.WorkingDirectory, comp.PromptCount, comp.TokenCount,
        ); err != nil {
            return err
        }
    }
    
    return nil
}

func (m *SQLiteToPostgresMigrator) migrateTags() error {
    // Similar implementation for tags
    return nil
}
```

### Migration Testing

#### Unit Tests

```go
func TestMigrator(t *testing.T) {
    db := setupTestDB(t)
    migrator := NewMigrator(db)
    
    // Test initial migration
    if err := migrator.Migrate(); err != nil {
        t.Fatalf("Migrate() error = %v", err)
    }
    
    // Verify schema version
    var version int
    err := db.QueryRow("SELECT MAX(version) FROM schema_version").Scan(&version)
    if err != nil {
        t.Fatalf("Failed to get schema version: %v", err)
    }
    
    if version != 1 {
        t.Errorf("Expected schema version 1, got %d", version)
    }
    
    // Verify tables exist
    var tableName string
    err = db.QueryRow(`
        SELECT name FROM sqlite_master
        WHERE type='table' AND name='compositions'
    `).Scan(&tableName)
    if err != nil {
        t.Errorf("compositions table not created: %v", err)
    }
}
```

#### Integration Tests

```go
func TestSQLiteToPostgresMigrator(t *testing.T) {
    sqliteDB := setupSQLiteTestDB(t)
    postgresDB := setupPostgreSQLTestDB(t)
    
    // Insert test data into SQLite
    insertTestData(t, sqliteDB)
    
    // Migrate to PostgreSQL
    migrator := NewSQLiteToPostgresMigrator(sqliteDB, postgresDB)
    if err := migrator.Migrate(); err != nil {
        t.Fatalf("Migrate() error = %v", err)
    }
    
    // Verify data in PostgreSQL
    var count int
    err := postgresDB.QueryRow("SELECT COUNT(*) FROM compositions").Scan(&count)
    if err != nil {
        t.Fatalf("Failed to count compositions: %v", err)
    }
    
    if count != 10 {
        t.Errorf("Expected 10 compositions, got %d", count)
    }
}
```

### Migration Best Practices

#### Guidelines

1. **Always use transactions**: Wrap migrations in transactions for atomicity
2. **Test rollbacks**: Ensure rollback scripts work correctly
3. **Version sequentially**: Use sequential version numbers
4. **Document changes**: Add comments explaining what each migration does
5. **Test thoroughly**: Test migrations on sample data before production
6. **Backup first**: Always backup database before migration
7. **Monitor performance**: Watch for slow migrations on large datasets
8. **Handle errors gracefully**: Provide clear error messages for failures

#### Checklist

- [ ] Migration file follows naming convention (XXX_description.sql)
- [ ] Rollback file exists for each migration
- [ ] Migration is tested on sample data
- [ ] Migration is tested on empty database
- [ ] Migration is tested on populated database
- [ ] Rollback is tested
- [ ] Performance is acceptable on expected data size
- [ ] Documentation is updated
- [ ] Schema version is updated
- [ ] Backup is created before migration

### Schema Changes

**Adding Columns:**
```sql
ALTER TABLE compositions ADD COLUMN new_column TEXT DEFAULT '';
```

**Adding Indexes:**
```sql
CREATE INDEX idx_new_index ON compositions(new_column);
```

**Dropping Columns:**
```sql
-- SQLite doesn't support DROP COLUMN directly
-- Must recreate table:
CREATE TABLE compositions_new (...);
INSERT INTO compositions_new SELECT ... FROM compositions;
DROP TABLE compositions;
ALTER TABLE compositions_new RENAME TO compositions;
```

---

## Performance Considerations

### Indexing Strategy

- **Created at index**: Essential for listing compositions by date
- **UUID index**: Fast lookups by UUID
- **Working directory index**: Filter by project context
- **Tag index**: Fast tag-based filtering

### Query Optimization

**Use EXPLAIN QUERY PLAN:**
```sql
EXPLAIN QUERY PLAN
SELECT * FROM compositions
WHERE working_directory = ?
ORDER BY created_at DESC;
```

**Best Practices:**
- Use indexed columns in WHERE clauses
- Avoid SELECT * when possible
- Use LIMIT for large result sets
- Consider prepared statements for repeated queries

### FTS5 Performance

- FTS5 is optimized for full-text search
- Use MATCH operator for best performance
- Consider FTS5 ranking for relevance sorting
- Rebuild FTS5 index if performance degrades:
  ```sql
  INSERT INTO compositions_fts(compositions_fts) VALUES('rebuild');
  ```

---

## Data Integrity

### Foreign Key Constraints

Foreign keys are enforced to maintain referential integrity:

```sql
PRAGMA foreign_keys = ON;
```

### Transactions

Use transactions for multi-step operations:

```sql
BEGIN TRANSACTION;
INSERT INTO compositions (...) VALUES (...);
INSERT INTO composition_tags (...) VALUES (...);
COMMIT;
```

### Validation

- UUID format validation in application layer
- Tag name validation (alphanumeric, hyphens, underscores)
- Working directory path validation
- Token count non-negative constraint

---

## Backup and Recovery

### Backup Strategy

**Online Backup:**
```bash
sqlite3 ~/.promptstack/data/history.db ".backup ~/.promptstack/data/history.db.backup"
```

**Copy Method:**
```bash
cp ~/.promptstack/data/history.db ~/.promptstack/data/history.db.backup
```

### Recovery

**From Backup:**
```bash
cp ~/.promptstack/data/history.db.backup ~/.promptstack/data/history.db
```

**Corruption Recovery:**
```bash
sqlite3 ~/.promptstack/data/history.db "PRAGMA integrity_check;"
```

---

## Maintenance

### Vacuum

Reclaim unused space and defragment:
```sql
VACUUM;
```

### Analyze

Update query planner statistics:
```sql
ANALYZE;
```

### Cleanup Old Records

Delete compositions older than specified date:
```sql
DELETE FROM compositions
WHERE created_at < date('now', '-90 days');
```

---

## Testing

### Test Data Setup

```sql
-- Insert test composition
INSERT INTO compositions (uuid, content, working_directory, prompt_count, token_count)
VALUES ('test-uuid-1', 'Test content', '/test/dir', 1, 100);

-- Add test tags
INSERT INTO composition_tags (composition_id, tag)
VALUES (1, 'test'), (1, 'example');
```

### Test Queries

```sql
-- Verify FTS5 search
SELECT * FROM compositions_fts WHERE compositions_fts MATCH 'Test';

-- Verify triggers
UPDATE compositions SET content = 'Updated content' WHERE id = 1;
SELECT updated_at FROM compositions WHERE id = 1;

-- Verify cascade delete
DELETE FROM compositions WHERE id = 1;
SELECT * FROM composition_tags WHERE composition_id = 1; -- Should be empty
```

---

## Security Considerations

### File Permissions

Ensure database file has appropriate permissions:
```bash
chmod 600 ~/.promptstack/data/history.db
```

### SQL Injection Prevention

- Always use parameterized queries
- Never concatenate user input into SQL
- Validate and sanitize all inputs

### Data Privacy

- Database contains user's composition content
- Consider encryption for sensitive data
- Respect user's privacy expectations

---

## Integration with Markdown Files

### Dual Storage Strategy

The database works in tandem with markdown files:

1. **Markdown Files**: Primary storage in `~/.promptstack/data/.history/`
   - Filename format: `YYYY-MM-DD_HH-MM-SS.md`
   - Contains full composition with YAML frontmatter
   - Human-readable and editable

2. **SQLite Database**: Search and metadata
   - Mirrors content from markdown files
   - Enables fast search and filtering
   - Tracks metadata and statistics

### Sync Verification

On startup, verify database matches markdown files:

```sql
-- Check for missing database entries
SELECT uuid FROM compositions
WHERE uuid NOT IN (
    SELECT substr(filename, 1, 36) FROM markdown_files
);

-- Check for missing markdown files
SELECT uuid FROM compositions
WHERE uuid NOT IN (
    SELECT uuid FROM markdown_files
);
```

### Rebuild Option

If sync mismatch detected, offer to rebuild database from markdown files.

---

## Future Enhancements

### Potential Schema Changes

- Add `author` field for multi-user scenarios
- Add `is_favorite` flag for quick access
- Add `template_id` for template-based compositions
- Add `session_id` for session grouping
- Add `ai_model_used` for tracking AI interactions

### Performance Improvements

- Consider partitioning by date for large datasets
- Implement query result caching
- Add materialized views for common queries
- Optimize FTS5 configuration for specific use cases

### Features

- Add composition versioning
- Implement composition sharing/collaboration
- Add composition templates
- Implement composition analytics

---

## References

- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [SQLite FTS5 Extension](https://www.sqlite.org/fts5.html)
- [SQLite Foreign Keys](https://www.sqlite.org/foreignkeys.html)
- [SQLite Transaction Management](https://www.sqlite.org/lang_transaction.html)

---

**Last Updated**: 2026-01-07  
**Schema Version**: 1.0  
**Status**: Ready for implementation (Milestone 15)