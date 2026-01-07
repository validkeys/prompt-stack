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

The database schema includes a `schema_version` table to track migrations:

```sql
CREATE TABLE schema_version (
    version INTEGER PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);
```

### Migration Process

1. **Check Current Version**
   ```sql
   SELECT MAX(version) FROM schema_version;
   ```

2. **Apply Migrations Sequentially**
   - Each migration is a SQL file named `migration_001.sql`, `migration_002.sql`, etc.
   - Migrations are applied in order
   - Each migration updates the `schema_version` table

3. **Migration Example**
   ```sql
   -- migration_002_add_index.sql
   BEGIN TRANSACTION;
   
   CREATE INDEX idx_compositions_token_count ON compositions(token_count);
   
   INSERT INTO schema_version (version, description)
   VALUES (2, 'Add index on token_count');
   
   COMMIT;
   ```

### Rollback Strategy

- Keep rollback scripts for each migration
- Test rollback procedures before deployment
- Document breaking changes clearly

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