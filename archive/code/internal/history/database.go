package history

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

// Database handles SQLite database operations
type Database struct {
	db     *sql.DB
	logger *zap.Logger

	// Prepared statements
	stmtInsert    *sql.Stmt
	stmtUpdate    *sql.Stmt
	stmtGetByPath *sql.Stmt
	stmtGetAll    *sql.Stmt
	stmtGetByDir  *sql.Stmt
	stmtGetByDate *sql.Stmt
	stmtDelete    *sql.Stmt
	stmtSearch    *sql.Stmt
	stmtExists    *sql.Stmt
}

// Initialize creates and initializes the SQLite database
func Initialize(dbPath string, logger *zap.Logger) (*Database, error) {
	logger.Info("Initializing database", zap.String("path", dbPath))

	// Ensure directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection with connection pooling
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for optimal performance
	// SetMaxOpenConns: Maximum number of open connections to the database
	// SQLite typically works well with 1-2 connections due to file locking
	db.SetMaxOpenConns(2)

	// SetMaxIdleConns: Maximum number of idle connections in the pool
	// Keep 1 idle connection ready for reuse
	db.SetMaxIdleConns(1)

	// SetConnMaxLifetime: Maximum amount of time a connection may be reused
	// 5 minutes to prevent long-lived connections from accumulating issues
	db.SetConnMaxLifetime(5 * 60 * 1000000000) // 5 minutes in nanoseconds

	// SetConnMaxIdleTime: Maximum amount of time a connection may be idle
	// 2 minutes before closing idle connections
	db.SetConnMaxIdleTime(2 * 60 * 1000000000) // 2 minutes in nanoseconds

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{
		db:     db,
		logger: logger,
	}

	// Create schema
	if err := database.createSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	// Enable WAL mode for better concurrency
	if err := database.enableWALMode(); err != nil {
		logger.Warn("Failed to enable WAL mode", zap.Error(err))
	}

	// Prepare statements for better performance
	if err := database.prepareStatements(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to prepare statements: %w", err)
	}

	logger.Info("Database initialized successfully")
	return database, nil
}

// createSchema creates the database schema
func (d *Database) createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS compositions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL,
		working_directory TEXT NOT NULL,
		content TEXT NOT NULL,
		character_count INTEGER NOT NULL,
		line_count INTEGER NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_created_at ON compositions(created_at);
	CREATE INDEX IF NOT EXISTS idx_working_directory ON compositions(working_directory);

	CREATE VIRTUAL TABLE IF NOT EXISTS compositions_fts USING fts5(
		content,
		working_directory,
		content='compositions',
		content_rowid='id'
	);
	`

	_, err := d.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// enableWALMode enables Write-Ahead Logging mode for better concurrency
func (d *Database) enableWALMode() error {
	_, err := d.db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return fmt.Errorf("failed to enable WAL mode: %w", err)
	}
	return nil
}

// prepareStatements prepares all SQL statements for reuse
func (d *Database) prepareStatements() error {
	var err error

	// Insert new composition
	d.stmtInsert, err = d.db.Prepare(`
		INSERT INTO compositions (file_path, created_at, working_directory, content, character_count, line_count, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}

	// Update existing composition
	d.stmtUpdate, err = d.db.Prepare(`
		UPDATE compositions
		SET content = ?, character_count = ?, line_count = ?, updated_at = ?
		WHERE file_path = ?
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}

	// Get composition by file path
	d.stmtGetByPath, err = d.db.Prepare(`
		SELECT id, file_path, created_at, working_directory, content, character_count, line_count, updated_at
		FROM compositions
		WHERE file_path = ?
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare get by path statement: %w", err)
	}

	// Get all compositions ordered by created_at desc
	d.stmtGetAll, err = d.db.Prepare(`
		SELECT id, file_path, created_at, working_directory, content, character_count, line_count, updated_at
		FROM compositions
		ORDER BY created_at DESC
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare get all statement: %w", err)
	}

	// Get compositions by working directory
	d.stmtGetByDir, err = d.db.Prepare(`
		SELECT id, file_path, created_at, working_directory, content, character_count, line_count, updated_at
		FROM compositions
		WHERE working_directory = ?
		ORDER BY created_at DESC
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare get by directory statement: %w", err)
	}

	// Get compositions by date range
	d.stmtGetByDate, err = d.db.Prepare(`
		SELECT id, file_path, created_at, working_directory, content, character_count, line_count, updated_at
		FROM compositions
		WHERE created_at >= ? AND created_at <= ?
		ORDER BY created_at DESC
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare get by date statement: %w", err)
	}

	// Delete composition by file path
	d.stmtDelete, err = d.db.Prepare(`
		DELETE FROM compositions WHERE file_path = ?
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}

	// Search compositions using FTS5
	d.stmtSearch, err = d.db.Prepare(`
		SELECT c.id, c.file_path, c.created_at, c.working_directory, c.content, c.character_count, c.line_count, c.updated_at
		FROM compositions c
		INNER JOIN compositions_fts fts ON c.id = fts.rowid
		WHERE compositions_fts MATCH ?
		ORDER BY c.created_at DESC
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare search statement: %w", err)
	}

	// Check if composition exists
	d.stmtExists, err = d.db.Prepare(`
		SELECT COUNT(*) FROM compositions WHERE file_path = ?
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare exists statement: %w", err)
	}

	return nil
}

// Close closes the database connection and all prepared statements
func (d *Database) Close() error {
	// Close all prepared statements
	if d.stmtInsert != nil {
		d.stmtInsert.Close()
	}
	if d.stmtUpdate != nil {
		d.stmtUpdate.Close()
	}
	if d.stmtGetByPath != nil {
		d.stmtGetByPath.Close()
	}
	if d.stmtGetAll != nil {
		d.stmtGetAll.Close()
	}
	if d.stmtGetByDir != nil {
		d.stmtGetByDir.Close()
	}
	if d.stmtGetByDate != nil {
		d.stmtGetByDate.Close()
	}
	if d.stmtDelete != nil {
		d.stmtDelete.Close()
	}
	if d.stmtSearch != nil {
		d.stmtSearch.Close()
	}
	if d.stmtExists != nil {
		d.stmtExists.Close()
	}

	// Close database connection
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// GetDB returns the underlying database connection
func (d *Database) GetDB() *sql.DB {
	return d.db
}

// Composition represents a composition record in the database
type Composition struct {
	ID               int
	FilePath         string
	CreatedAt        string
	WorkingDirectory string
	Content          string
	CharacterCount   int
	LineCount        int
	UpdatedAt        string
}

// InsertComposition inserts a new composition into the database
func (d *Database) InsertComposition(comp Composition) error {
	d.logger.Debug("Inserting composition",
		zap.String("file_path", comp.FilePath),
		zap.Int("character_count", comp.CharacterCount),
		zap.Int("line_count", comp.LineCount),
	)

	_, err := d.stmtInsert.Exec(
		comp.FilePath,
		comp.CreatedAt,
		comp.WorkingDirectory,
		comp.Content,
		comp.CharacterCount,
		comp.LineCount,
		comp.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert composition: %w", err)
	}

	return nil
}

// UpdateComposition updates an existing composition in the database
func (d *Database) UpdateComposition(comp Composition) error {
	d.logger.Debug("Updating composition",
		zap.String("file_path", comp.FilePath),
		zap.Int("character_count", comp.CharacterCount),
		zap.Int("line_count", comp.LineCount),
	)

	result, err := d.stmtUpdate.Exec(
		comp.Content,
		comp.CharacterCount,
		comp.LineCount,
		comp.UpdatedAt,
		comp.FilePath,
	)
	if err != nil {
		return fmt.Errorf("failed to update composition: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("composition not found: %s", comp.FilePath)
	}

	return nil
}

// GetCompositionByPath retrieves a composition by its file path
func (d *Database) GetCompositionByPath(filePath string) (*Composition, error) {
	d.logger.Debug("Getting composition by path", zap.String("file_path", filePath))

	var comp Composition
	err := d.stmtGetByPath.QueryRow(filePath).Scan(
		&comp.ID,
		&comp.FilePath,
		&comp.CreatedAt,
		&comp.WorkingDirectory,
		&comp.Content,
		&comp.CharacterCount,
		&comp.LineCount,
		&comp.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get composition: %w", err)
	}

	return &comp, nil
}

// GetAllCompositions retrieves all compositions ordered by creation date
func (d *Database) GetAllCompositions() ([]Composition, error) {
	d.logger.Debug("Getting all compositions")

	rows, err := d.stmtGetAll.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to query compositions: %w", err)
	}
	defer rows.Close()

	return d.scanCompositions(rows)
}

// GetCompositionsByDirectory retrieves compositions for a specific working directory
func (d *Database) GetCompositionsByDirectory(workingDir string) ([]Composition, error) {
	d.logger.Debug("Getting compositions by directory", zap.String("working_directory", workingDir))

	rows, err := d.stmtGetByDir.Query(workingDir)
	if err != nil {
		return nil, fmt.Errorf("failed to query compositions by directory: %w", err)
	}
	defer rows.Close()

	return d.scanCompositions(rows)
}

// GetCompositionsByDateRange retrieves compositions within a date range
func (d *Database) GetCompositionsByDateRange(startDate, endDate string) ([]Composition, error) {
	d.logger.Debug("Getting compositions by date range",
		zap.String("start_date", startDate),
		zap.String("end_date", endDate),
	)

	rows, err := d.stmtGetByDate.Query(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query compositions by date range: %w", err)
	}
	defer rows.Close()

	return d.scanCompositions(rows)
}

// DeleteComposition deletes a composition by its file path
func (d *Database) DeleteComposition(filePath string) error {
	d.logger.Debug("Deleting composition", zap.String("file_path", filePath))

	result, err := d.stmtDelete.Exec(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete composition: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("composition not found: %s", filePath)
	}

	return nil
}

// SearchCompositions searches compositions using FTS5 full-text search
func (d *Database) SearchCompositions(query string) ([]Composition, error) {
	d.logger.Debug("Searching compositions", zap.String("query", query))

	rows, err := d.stmtSearch.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search compositions: %w", err)
	}
	defer rows.Close()

	return d.scanCompositions(rows)
}

// CompositionExists checks if a composition exists by its file path
func (d *Database) CompositionExists(filePath string) (bool, error) {
	d.logger.Debug("Checking if composition exists", zap.String("file_path", filePath))

	var count int
	err := d.stmtExists.QueryRow(filePath).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check composition existence: %w", err)
	}

	return count > 0, nil
}

// scanCompositions scans rows and returns a slice of compositions
func (d *Database) scanCompositions(rows *sql.Rows) ([]Composition, error) {
	var compositions []Composition

	for rows.Next() {
		var comp Composition
		err := rows.Scan(
			&comp.ID,
			&comp.FilePath,
			&comp.CreatedAt,
			&comp.WorkingDirectory,
			&comp.Content,
			&comp.CharacterCount,
			&comp.LineCount,
			&comp.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan composition: %w", err)
		}
		compositions = append(compositions, comp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating compositions: %w", err)
	}

	return compositions, nil
}
