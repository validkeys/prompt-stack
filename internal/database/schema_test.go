package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestInit(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	t.Run("creates database file", func(t *testing.T) {
		err := Init(dbPath)
		if err != nil {
			t.Fatalf("Init() error = %v", err)
		}

		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Fatal("database file was not created")
		}
	})

	t.Run("creates required tables", func(t *testing.T) {
		dbPath2 := filepath.Join(tmpDir, "test2.db")
		if err := Init(dbPath2); err != nil {
			t.Fatalf("Init() error = %v", err)
		}

		db, err := sql.Open("sqlite3", dbPath2)
		if err != nil {
			t.Fatalf("failed to open database: %v", err)
		}
		defer db.Close()

		tables := []string{"patterns", "requirements", "tasks", "validation_reports"}
		for _, table := range tables {
			var exists bool
			err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name=?)", table).Scan(&exists)
			if err != nil {
				t.Fatalf("failed to check table existence: %v", err)
			}
			if !exists {
				t.Errorf("table %s was not created", table)
			}
		}
	})

	t.Run("creates parent directories", func(t *testing.T) {
		nestedPath := filepath.Join(tmpDir, "nested", "dir", "test.db")
		err := Init(nestedPath)
		if err != nil {
			t.Fatalf("Init() error = %v", err)
		}

		if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
			t.Fatal("database file was not created in nested directory")
		}
	})
}

func TestCreateTables(t *testing.T) {
	t.Run("creates all required tables", func(t *testing.T) {
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("failed to open in-memory database: %v", err)
		}
		defer db.Close()

		if err := createTables(db); err != nil {
			t.Fatalf("createTables() error = %v", err)
		}

		tables := []string{"patterns", "requirements", "tasks", "validation_reports"}
		for _, table := range tables {
			var exists bool
			err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name=?)", table).Scan(&exists)
			if err != nil {
				t.Fatalf("failed to check table existence: %v", err)
			}
			if !exists {
				t.Errorf("table %s was not created", table)
			}
		}
	})

	t.Run("tables have expected columns", func(t *testing.T) {
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("failed to open in-memory database: %v", err)
		}
		defer db.Close()

		if err := createTables(db); err != nil {
			t.Fatalf("createTables() error = %v", err)
		}

		tests := []struct {
			table   string
			columns []string
		}{
			{"patterns", []string{"id", "name", "description", "confidence", "category"}},
			{"requirements", []string{"id", "milestone_id", "requirement_type", "content", "priority", "status"}},
			{"tasks", []string{"id", "milestone_id", "task_id", "title", "description", "dependencies", "status"}},
			{"validation_reports", []string{"id", "plan_path", "overall_score", "approval_status", "report_json"}},
		}

		for _, tt := range tests {
			for _, column := range tt.columns {
				var exists bool
				err := db.QueryRow(`
					SELECT EXISTS (
						SELECT 1 
						FROM pragma_table_info(?) 
						WHERE name=?
					)
				`, tt.table, column).Scan(&exists)
				if err != nil {
					t.Errorf("failed to check column %s.%s: %v", tt.table, column, err)
				}
				if !exists {
					t.Errorf("column %s.%s was not created", tt.table, column)
				}
			}
		}
	})
}

func TestSetSchemaVersion(t *testing.T) {
	t.Run("sets schema version", func(t *testing.T) {
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("failed to open in-memory database: %v", err)
		}
		defer db.Close()

		if err := setSchemaVersion(db); err != nil {
			t.Fatalf("setSchemaVersion() error = %v", err)
		}

		var version int
		err = db.QueryRow("SELECT version FROM schema_version").Scan(&version)
		if err != nil {
			t.Fatalf("failed to query schema version: %v", err)
		}

		if version != schemaVersion {
			t.Errorf("schema version = %d, want %d", version, schemaVersion)
		}
	})

	t.Run("updates existing schema version", func(t *testing.T) {
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("failed to open in-memory database: %v", err)
		}
		defer db.Close()

		_, err = db.Exec("CREATE TABLE schema_version (version INTEGER PRIMARY KEY)")
		if err != nil {
			t.Fatalf("failed to create schema_version table: %v", err)
		}

		_, err = db.Exec("INSERT INTO schema_version (version) VALUES (0)")
		if err != nil {
			t.Fatalf("failed to insert initial version: %v", err)
		}

		if err := setSchemaVersion(db); err != nil {
			t.Fatalf("setSchemaVersion() error = %v", err)
		}

		var version int
		err = db.QueryRow("SELECT version FROM schema_version").Scan(&version)
		if err != nil {
			t.Fatalf("failed to query schema version: %v", err)
		}

		if version != schemaVersion {
			t.Errorf("schema version = %d, want %d", version, schemaVersion)
		}
	})
}
