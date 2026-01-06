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
}

// Initialize creates and initializes the SQLite database
func Initialize(dbPath string, logger *zap.Logger) (*Database, error) {
	logger.Info("Initializing database", zap.String("path", dbPath))

	// Ensure directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

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

// Close closes the database connection
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// GetDB returns the underlying database connection
func (d *Database) GetDB() *sql.DB {
	return d.db
}
