package history

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// Storage handles markdown file operations for history
type Storage struct {
	historyDir string
	logger     *zap.Logger
}

// NewStorage creates a new storage instance for history files
func NewStorage(historyDir string, logger *zap.Logger) (*Storage, error) {
	// Ensure history directory exists
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create history directory: %w", err)
	}

	return &Storage{
		historyDir: historyDir,
		logger:     logger,
	}, nil
}

// CreateHistoryFile creates a new timestamped markdown file for a composition
func (s *Storage) CreateHistoryFile(workingDir string, content string) (string, error) {
	// Generate timestamp-based filename
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s.md", timestamp)
	filePath := filepath.Join(s.historyDir, filename)

	s.logger.Debug("Creating history file",
		zap.String("file_path", filePath),
		zap.String("working_directory", workingDir),
		zap.Int("content_length", len(content)),
	)

	// Create file with content
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write history file: %w", err)
	}

	s.logger.Info("History file created", zap.String("file_path", filePath))
	return filePath, nil
}

// ReadHistoryFile reads content from a history file
func (s *Storage) ReadHistoryFile(filePath string) (string, error) {
	s.logger.Debug("Reading history file", zap.String("file_path", filePath))

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read history file: %w", err)
	}

	return string(content), nil
}

// UpdateHistoryFile updates an existing history file with new content
func (s *Storage) UpdateHistoryFile(filePath string, content string) error {
	s.logger.Debug("Updating history file",
		zap.String("file_path", filePath),
		zap.Int("content_length", len(content)),
	)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to update history file: %w", err)
	}

	return nil
}

// DeleteHistoryFile deletes a history file
func (s *Storage) DeleteHistoryFile(filePath string) error {
	s.logger.Debug("Deleting history file", zap.String("file_path", filePath))

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete history file: %w", err)
	}

	s.logger.Info("History file deleted", zap.String("file_path", filePath))
	return nil
}

// ListHistoryFiles returns all history files sorted by modification time (newest first)
func (s *Storage) ListHistoryFiles() ([]string, error) {
	s.logger.Debug("Listing history files")

	entries, err := os.ReadDir(s.historyDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read history directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			filePath := filepath.Join(s.historyDir, entry.Name())
			files = append(files, filePath)
		}
	}

	// Sort files by modification time (newest first)
	if len(files) > 0 {
		sortedFiles, err := s.sortFilesByModTime(files)
		if err != nil {
			return nil, fmt.Errorf("failed to sort history files: %w", err)
		}
		files = sortedFiles
	}

	s.logger.Debug("Found history files", zap.Int("count", len(files)))
	return files, nil
}

// sortFilesByModTime sorts files by modification time (newest first)
func (s *Storage) sortFilesByModTime(files []string) ([]string, error) {
	type fileInfo struct {
		path    string
		modTime time.Time
	}

	var fileInfos []fileInfo
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return nil, fmt.Errorf("failed to stat file %s: %w", file, err)
		}
		fileInfos = append(fileInfos, fileInfo{
			path:    file,
			modTime: info.ModTime(),
		})
	}

	// Sort by modification time (newest first)
	for i := 0; i < len(fileInfos); i++ {
		for j := i + 1; j < len(fileInfos); j++ {
			if fileInfos[i].modTime.Before(fileInfos[j].modTime) {
				fileInfos[i], fileInfos[j] = fileInfos[j], fileInfos[i]
			}
		}
	}

	var sortedFiles []string
	for _, fi := range fileInfos {
		sortedFiles = append(sortedFiles, fi.path)
	}

	return sortedFiles, nil
}

// GetHistoryDir returns the history directory path
func (s *Storage) GetHistoryDir() string {
	return s.historyDir
}

// FileExists checks if a history file exists
func (s *Storage) FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// GetFileModTime returns the modification time of a history file
func (s *Storage) GetFileModTime(filePath string) (time.Time, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get file modification time: %w", err)
	}
	return info.ModTime(), nil
}

// GetFileSize returns the size of a history file in bytes
func (s *Storage) GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file size: %w", err)
	}
	return info.Size(), nil
}
