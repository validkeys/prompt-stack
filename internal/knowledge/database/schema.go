package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const schemaVersion = 1

func Init(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := createTables(db); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	if err := setSchemaVersion(db); err != nil {
		// attempt to close the DB and return the original error
		_ = db.Close()
		return fmt.Errorf("failed to set schema version: %w", err)
	}

	if cerr := db.Close(); cerr != nil {
		return fmt.Errorf("failed to close database: %w", cerr)
	}

	return nil
}

func createTables(db *sql.DB) error {
	patternsTable := `
	CREATE TABLE IF NOT EXISTS patterns (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		confidence REAL,
		category TEXT,
		source TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	requirementsTable := `
	CREATE TABLE IF NOT EXISTS requirements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		milestone_id TEXT NOT NULL,
		requirement_type TEXT,
		content TEXT,
		priority INTEGER,
		status TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	tasksTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		milestone_id TEXT NOT NULL,
		task_id TEXT NOT NULL,
		title TEXT,
		description TEXT,
		dependencies TEXT,
		status TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	validationTable := `
	CREATE TABLE IF NOT EXISTS validation_reports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		plan_path TEXT NOT NULL,
		overall_score REAL,
		approval_status TEXT,
		report_json TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	tables := []string{patternsTable, requirementsTable, tasksTable, validationTable}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

func setSchemaVersion(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER PRIMARY KEY
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		DELETE FROM schema_version
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO schema_version (version) VALUES (?)
	`, schemaVersion)
	return err
}
