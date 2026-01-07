package history

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// Sync handles synchronization between markdown files and SQLite database
type Sync struct {
	db      *Database
	storage *Storage
	logger  *zap.Logger
}

// NewSync creates a new sync handler
func NewSync(db *Database, storage *Storage, logger *zap.Logger) *Sync {
	return &Sync{
		db:      db,
		storage: storage,
		logger:  logger,
	}
}

// SyncStatus represents the synchronization status
type SyncStatus struct {
	IsSynced       bool
	TotalFiles     int
	TotalDBEntries int
	MissingFiles   []string
	MissingEntries []string
	OutOfSync      []SyncMismatch
}

// SyncMismatch represents a mismatch between file and database
type SyncMismatch struct {
	FilePath      string
	FileSize      int64
	DBSize        int
	FileModTime   time.Time
	DBUpdatedTime time.Time
}

// VerifySync checks if markdown files and database are in sync
func (s *Sync) VerifySync() (*SyncStatus, error) {
	s.logger.Info("Verifying sync between markdown files and database")

	status := &SyncStatus{
		IsSynced: true,
	}

	// Get all markdown files
	files, err := s.storage.ListHistoryFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to list history files: %w", err)
	}
	status.TotalFiles = len(files)

	// Get all database entries
	compositions, err := s.db.GetAllCompositions()
	if err != nil {
		return nil, fmt.Errorf("failed to get database entries: %w", err)
	}
	status.TotalDBEntries = len(compositions)

	// Create maps for comparison
	fileMap := make(map[string]bool)
	for _, file := range files {
		fileMap[file] = true
	}

	dbMap := make(map[string]*Composition)
	for i := range compositions {
		dbMap[compositions[i].FilePath] = &compositions[i]
	}

	// Check for files not in database
	for file := range fileMap {
		if _, exists := dbMap[file]; !exists {
			status.MissingEntries = append(status.MissingEntries, file)
			status.IsSynced = false
		}
	}

	// Check for database entries without files
	for filePath := range dbMap {
		if _, exists := fileMap[filePath]; !exists {
			status.MissingFiles = append(status.MissingFiles, filePath)
			status.IsSynced = false
		}
	}

	// Check for mismatches between file and database
	for filePath := range fileMap {
		if comp, exists := dbMap[filePath]; exists {
			mismatch, err := s.checkMismatch(filePath, comp)
			if err != nil {
				s.logger.Warn("Failed to check mismatch",
					zap.String("file_path", filePath),
					zap.Error(err),
				)
				continue
			}

			if mismatch != nil {
				status.OutOfSync = append(status.OutOfSync, *mismatch)
				status.IsSynced = false
			}
		}
	}

	if status.IsSynced {
		s.logger.Info("Sync verification passed",
			zap.Int("files", status.TotalFiles),
			zap.Int("db_entries", status.TotalDBEntries),
		)
	} else {
		s.logger.Warn("Sync verification failed",
			zap.Int("missing_files", len(status.MissingFiles)),
			zap.Int("missing_entries", len(status.MissingEntries)),
			zap.Int("out_of_sync", len(status.OutOfSync)),
		)
	}

	return status, nil
}

// checkMismatch checks if a file and database entry are out of sync
func (s *Sync) checkMismatch(filePath string, comp *Composition) (*SyncMismatch, error) {
	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Get file size
	fileSize := fileInfo.Size()

	// Parse database update time
	dbUpdatedTime, err := time.Parse(time.RFC3339, comp.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database update time: %w", err)
	}

	// Compare sizes (allow small tolerance for metadata differences)
	sizeDiff := abs(fileSize - int64(comp.CharacterCount))
	if sizeDiff > 100 { // More than 100 characters difference
		return &SyncMismatch{
			FilePath:      filePath,
			FileSize:      fileSize,
			DBSize:        comp.CharacterCount,
			FileModTime:   fileInfo.ModTime(),
			DBUpdatedTime: dbUpdatedTime,
		}, nil
	}

	// Compare modification times (allow 1 second tolerance)
	timeDiff := fileInfo.ModTime().Sub(dbUpdatedTime)
	if timeDiff.Abs() > time.Second {
		return &SyncMismatch{
			FilePath:      filePath,
			FileSize:      fileSize,
			DBSize:        comp.CharacterCount,
			FileModTime:   fileInfo.ModTime(),
			DBUpdatedTime: dbUpdatedTime,
		}, nil
	}

	return nil, nil
}

// Rebuild rebuilds the database from markdown files
func (s *Sync) Rebuild() error {
	s.logger.Info("Rebuilding database from markdown files")

	// Get all markdown files
	files, err := s.storage.ListHistoryFiles()
	if err != nil {
		return fmt.Errorf("failed to list history files: %w", err)
	}

	s.logger.Info("Found history files for rebuild", zap.Int("count", len(files)))

	// Process each file
	for i, filePath := range files {
		s.logger.Debug("Processing file for rebuild",
			zap.String("file_path", filePath),
			zap.Int("progress", i+1),
			zap.Int("total", len(files)),
		)

		// Read file content
		content, err := s.storage.ReadHistoryFile(filePath)
		if err != nil {
			s.logger.Warn("Failed to read file during rebuild",
				zap.String("file_path", filePath),
				zap.Error(err),
			)
			continue
		}

		// Get file info
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			s.logger.Warn("Failed to stat file during rebuild",
				zap.String("file_path", filePath),
				zap.Error(err),
			)
			continue
		}

		// Check if entry exists in database
		exists, err := s.db.CompositionExists(filePath)
		if err != nil {
			s.logger.Warn("Failed to check composition existence",
				zap.String("file_path", filePath),
				zap.Error(err),
			)
			continue
		}

		// Calculate metadata
		charCount := len(content)
		lineCount := countLines(content)
		modTime := fileInfo.ModTime().Format(time.RFC3339)

		comp := Composition{
			FilePath:         filePath,
			CreatedAt:        modTime,                    // Use mod time as created time for rebuild
			WorkingDirectory: extractWorkingDir(content), // Try to extract from content
			Content:          content,
			CharacterCount:   charCount,
			LineCount:        lineCount,
			UpdatedAt:        modTime,
		}

		// Insert or update
		if exists {
			if err := s.db.UpdateComposition(comp); err != nil {
				s.logger.Warn("Failed to update composition during rebuild",
					zap.String("file_path", filePath),
					zap.Error(err),
				)
			}
		} else {
			if err := s.db.InsertComposition(comp); err != nil {
				s.logger.Warn("Failed to insert composition during rebuild",
					zap.String("file_path", filePath),
					zap.Error(err),
				)
			}
		}
	}

	s.logger.Info("Database rebuild completed", zap.Int("files_processed", len(files)))
	return nil
}

// BackupDatabase creates a backup of the database
func (s *Sync) BackupDatabase() (string, error) {
	s.logger.Info("Creating database backup")

	// Get database path from storage
	dbPath := filepath.Join(s.storage.GetHistoryDir(), "history.db")

	// Create backup filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := filepath.Join(s.storage.GetHistoryDir(), fmt.Sprintf("history.db.backup-%s", timestamp))

	// Read database file
	dbData, err := os.ReadFile(dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to read database file: %w", err)
	}

	// Write backup
	if err := os.WriteFile(backupPath, dbData, 0644); err != nil {
		return "", fmt.Errorf("failed to write backup file: %w", err)
	}

	s.logger.Info("Database backup created", zap.String("backup_path", backupPath))
	return backupPath, nil
}

// Helper functions

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func countLines(content string) int {
	count := 0
	for _, c := range content {
		if c == '\n' {
			count++
		}
	}
	if len(content) > 0 && content[len(content)-1] != '\n' {
		count++
	}
	return count
}

func extractWorkingDir(content string) string {
	// Try to extract working directory from content
	// This is a simple heuristic - in production, you might store this in metadata
	lines := splitLines(content, 10) // Check first 10 lines
	for _, line := range lines {
		if len(line) > 20 && (line[0:20] == "Working Directory: " || line[0:20] == "working_directory:") {
			// Extract directory from line
			parts := splitString(line, ":")
			if len(parts) > 1 {
				return trimSpace(parts[1])
			}
		}
	}
	return "unknown"
}

func splitLines(content string, maxLines int) []string {
	lines := []string{}
	count := 0
	start := 0

	for i, c := range content {
		if c == '\n' {
			lines = append(lines, content[start:i])
			start = i + 1
			count++
			if count >= maxLines {
				break
			}
		}
	}

	if start < len(content) && count < maxLines {
		lines = append(lines, content[start:])
	}

	return lines
}

func splitString(s, sep string) []string {
	var parts []string
	start := 0

	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			parts = append(parts, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}

	if start < len(s) {
		parts = append(parts, s[start:])
	}

	return parts
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}

	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}

	return s[start:end]
}
