# Database Schema Updates for Scalability

**Document**: [`DATABASE-SCHEMA.md`](../DATABASE-SCHEMA.md)  
**Priority**: **Medium**  
**Status**: Should update before implementation begins

---

## Overview

This document details all changes required to [`DATABASE-SCHEMA.md`](../DATABASE-SCHEMA.md) to incorporate Phase 1 scalability abstractions. These changes add a migration strategy section for database schema evolution.

---

## Migration Strategy

### Version Management

Track schema version in database:

```sql
CREATE TABLE schema_version (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP NOT NULL
);

-- Insert initial version
INSERT INTO schema_version (version, applied_at) VALUES (1, CURRENT_TIMESTAMP);
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
    id TEXT PRIMARY KEY,
    file_path TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    working_directory TEXT NOT NULL,
    content TEXT NOT NULL,
    character_count INTEGER NOT NULL,
    line_count INTEGER NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- 001_initial_schema_rollback.sql
DROP TABLE IF EXISTS compositions;
```

### Future Migrations

- PostgreSQL schema: Different dialect
- Neo4j schema: Graph-based
- Data migration: SQLite → PostgreSQL

---

## Migration Implementation

### Migration Runner

```go
// internal/storage/migrations/migrator.go
package migrations

import (
    "database/sql"
    "fmt"
    "sort"
    "strings"
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
            applied_at TIMESTAMP NOT NULL
        )
    `)
    return err
}

func (m *Migrator) getCurrentVersion() (int, error) {
    var version int
    err := m.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&version)
    return version, err
}

func (m *Migrator) getAvailableMigrations() ([]Migration, error) {
    // Read migration files from migrations directory
    // Parse version numbers from filenames
    // Return sorted list of migrations
    return []Migration{}, nil
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
        "INSERT INTO schema_version (version, applied_at) VALUES (?, CURRENT_TIMESTAMP)",
        migration.Version,
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
    Version int
    SQL     string
    Name    string
}
```

---

## Migration Files

### 001_initial_schema.sql

```sql
-- Create compositions table
CREATE TABLE compositions (
    id TEXT PRIMARY KEY,
    file_path TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    working_directory TEXT NOT NULL,
    content TEXT NOT NULL,
    character_count INTEGER NOT NULL,
    line_count INTEGER NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create indexes
CREATE INDEX idx_compositions_created_at ON compositions(created_at);
CREATE INDEX idx_compositions_working_directory ON compositions(working_directory);
```

### 001_initial_schema_rollback.sql

```sql
DROP INDEX IF EXISTS idx_compositions_working_directory;
DROP INDEX IF EXISTS idx_compositions_created_at;
DROP TABLE IF EXISTS compositions;
```

### 002_add_fts5.sql

```sql
-- Create FTS5 virtual table for full-text search
CREATE VIRTUAL TABLE compositions_fts USING fts5(
    id,
    content,
    content=compositions,
    content_rowid=rowid
);

-- Populate FTS5 table
INSERT INTO compositions_fts(rowid, id, content)
SELECT rowid, id, content FROM compositions;

-- Create triggers to keep FTS5 in sync
CREATE TRIGGER compositions_ai AFTER INSERT ON compositions BEGIN
    INSERT INTO compositions_fts(rowid, id, content)
    VALUES (new.rowid, new.id, new.content);
END;

CREATE TRIGGER compositions_ad AFTER DELETE ON compositions BEGIN
    DELETE FROM compositions_fts WHERE rowid = old.rowid;
END;

CREATE TRIGGER compositions_au AFTER UPDATE ON compositions BEGIN
    UPDATE compositions_fts SET content = new.content WHERE rowid = new.rowid;
END;
```

### 002_add_fts5_rollback.sql

```sql
DROP TRIGGER IF EXISTS compositions_au;
DROP TRIGGER IF EXISTS compositions_ad;
DROP TRIGGER IF EXISTS compositions_ai;
DROP TABLE IF EXISTS compositions_fts;
```

### 003_add_composition_tags.sql

```sql
-- Create tags table
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

-- Create composition_tags junction table
CREATE TABLE composition_tags (
    composition_id TEXT NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (composition_id, tag_id),
    FOREIGN KEY (composition_id) REFERENCES compositions(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_composition_tags_composition_id ON composition_tags(composition_id);
CREATE INDEX idx_composition_tags_tag_id ON composition_tags(tag_id);
```

### 003_add_composition_tags_rollback.sql

```sql
DROP INDEX IF EXISTS idx_composition_tags_tag_id;
DROP INDEX IF EXISTS idx_composition_tags_composition_id;
DROP TABLE IF EXISTS composition_tags;
DROP TABLE IF EXISTS tags;
```

---

## Cross-Database Migration

### SQLite to PostgreSQL

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
            id TEXT PRIMARY KEY,
            file_path TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL,
            working_directory TEXT NOT NULL,
            content TEXT NOT NULL,
            character_count INTEGER NOT NULL,
            line_count INTEGER NOT NULL,
            updated_at TIMESTAMP NOT NULL
        );
        
        CREATE INDEX IF NOT EXISTS idx_compositions_created_at ON compositions(created_at);
        CREATE INDEX IF NOT EXISTS idx_compositions_working_directory ON compositions(working_directory);
        
        CREATE TABLE IF NOT EXISTS tags (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL UNIQUE
        );
        
        CREATE TABLE IF NOT EXISTS composition_tags (
            composition_id TEXT NOT NULL,
            tag_id INTEGER NOT NULL,
            PRIMARY KEY (composition_id, tag_id),
            FOREIGN KEY (composition_id) REFERENCES compositions(id) ON DELETE CASCADE,
            FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
        );
    `
    
    _, err := m.postgresDB.Exec(schema)
    return err
}

func (m *SQLiteToPostgresMigrator) migrateCompositions() error {
    rows, err := m.sqliteDB.Query(`
        SELECT id, file_path, created_at, working_directory, 
               content, character_count, line_count, updated_at
        FROM compositions
    `)
    if err != nil {
        return err
    }
    defer rows.Close()
    
    stmt, err := m.postgresDB.Prepare(`
        INSERT INTO compositions (id, file_path, created_at, working_directory,
                               content, character_count, line_count, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for rows.Next() {
        var comp struct {
            ID               string
            FilePath         string
            CreatedAt        string
            WorkingDirectory string
            Content          string
            CharacterCount   int
            LineCount        int
            UpdatedAt        string
        }
        
        if err := rows.Scan(
            &comp.ID, &comp.FilePath, &comp.CreatedAt, &comp.WorkingDirectory,
            &comp.Content, &comp.CharacterCount, &comp.LineCount, &comp.UpdatedAt,
        ); err != nil {
            return err
        }
        
        if _, err := stmt.Exec(
            comp.ID, comp.FilePath, comp.CreatedAt, comp.WorkingDirectory,
            comp.Content, comp.CharacterCount, comp.LineCount, comp.UpdatedAt,
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

---

## Migration Testing

### Unit Tests

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

### Integration Tests

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

---

## Migration Best Practices

### Guidelines

1. **Always use transactions**: Wrap migrations in transactions for atomicity
2. **Test rollbacks**: Ensure rollback scripts work correctly
3. **Version sequentially**: Use sequential version numbers
4. **Document changes**: Add comments explaining what each migration does
5. **Test thoroughly**: Test migrations on sample data before production
6. **Backup first**: Always backup database before migration
7. **Monitor performance**: Watch for slow migrations on large datasets
8. **Handle errors gracefully**: Provide clear error messages for failures

### Checklist

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

---

**Last Updated**: 2026-01-07  
**Related Documents**: 
- [`scalability-implementation-summary.md`](./scalability-implementation-summary.md)
- [`scalability-project-structure-updates.md`](./scalability-project-structure-updates.md)
- [`scalability-dependencies-updates.md`](./scalability-dependencies-updates.md)