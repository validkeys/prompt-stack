package history

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Manager handles history management with auto-save functionality
type Manager struct {
	db      *Database
	storage *Storage
	logger  *zap.Logger

	// Auto-save state
	mu               sync.Mutex
	currentFilePath  string
	workingDirectory string
	pendingSave      bool
	saveTimer        *time.Timer
	saveInterval     time.Duration
}

// NewManager creates a new history manager
func NewManager(db *Database, storage *Storage, logger *zap.Logger) *Manager {
	return &Manager{
		db:           db,
		storage:      storage,
		logger:       logger,
		saveInterval: 750 * time.Millisecond, // 750ms debounce interval
	}
}

// NewComposition creates a new composition with history file and database entry
func (m *Manager) NewComposition(workingDir, content string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("Creating new composition",
		zap.String("working_directory", workingDir),
		zap.Int("content_length", len(content)),
	)

	// Create timestamped history file
	filePath, err := m.storage.CreateHistoryFile(workingDir, content)
	if err != nil {
		return "", fmt.Errorf("failed to create history file: %w", err)
	}

	// Calculate metadata
	charCount := len(content)
	lineCount := strings.Count(content, "\n") + 1
	now := time.Now().Format(time.RFC3339)

	// Insert into database
	comp := Composition{
		FilePath:         filePath,
		CreatedAt:        now,
		WorkingDirectory: workingDir,
		Content:          content,
		CharacterCount:   charCount,
		LineCount:        lineCount,
		UpdatedAt:        now,
	}

	if err := m.db.InsertComposition(comp); err != nil {
		// Rollback: delete the file if database insert fails
		_ = m.storage.DeleteHistoryFile(filePath)
		return "", fmt.Errorf("failed to insert composition into database: %w", err)
	}

	// Update current composition state
	m.currentFilePath = filePath
	m.workingDirectory = workingDir
	m.pendingSave = false

	m.logger.Info("New composition created",
		zap.String("file_path", filePath),
		zap.Int("character_count", charCount),
		zap.Int("line_count", lineCount),
	)

	return filePath, nil
}

// LoadComposition loads an existing composition from history
func (m *Manager) LoadComposition(filePath string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Debug("Loading composition", zap.String("file_path", filePath))

	// Read from markdown file (source of truth)
	content, err := m.storage.ReadHistoryFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read history file: %w", err)
	}

	// Verify composition exists in database
	comp, err := m.db.GetCompositionByPath(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get composition from database: %w", err)
	}

	if comp == nil {
		m.logger.Warn("Composition not found in database, will sync on save",
			zap.String("file_path", filePath),
		)
	} else {
		// Update current composition state
		m.workingDirectory = comp.WorkingDirectory
	}

	m.currentFilePath = filePath
	m.pendingSave = false

	m.logger.Info("Composition loaded",
		zap.String("file_path", filePath),
		zap.Int("content_length", len(content)),
	)

	return content, nil
}

// SaveComposition saves composition changes to both markdown and SQLite
func (m *Manager) SaveComposition(content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.currentFilePath == "" {
		return fmt.Errorf("no active composition to save")
	}

	m.logger.Debug("Saving composition",
		zap.String("file_path", m.currentFilePath),
		zap.Int("content_length", len(content)),
	)

	// Calculate metadata
	charCount := len(content)
	lineCount := strings.Count(content, "\n") + 1
	now := time.Now().Format(time.RFC3339)

	// Update markdown file (source of truth)
	if err := m.storage.UpdateHistoryFile(m.currentFilePath, content); err != nil {
		return fmt.Errorf("failed to update history file: %w", err)
	}

	// Update database
	comp := Composition{
		FilePath:       m.currentFilePath,
		Content:        content,
		CharacterCount: charCount,
		LineCount:      lineCount,
		UpdatedAt:      now,
	}

	if err := m.db.UpdateComposition(comp); err != nil {
		// Log error but don't fail - markdown is source of truth
		m.logger.Error("Failed to update database, but markdown was saved",
			zap.String("file_path", m.currentFilePath),
			zap.Error(err),
		)
	}

	m.pendingSave = false

	m.logger.Debug("Composition saved successfully",
		zap.String("file_path", m.currentFilePath),
		zap.Int("character_count", charCount),
		zap.Int("line_count", lineCount),
	)

	return nil
}

// TriggerAutoSave triggers a debounced auto-save
func (m *Manager) TriggerAutoSave(content string) {
	m.mu.Lock()

	if m.currentFilePath == "" {
		m.mu.Unlock()
		return
	}

	// Cancel any pending save timer
	if m.saveTimer != nil {
		m.saveTimer.Stop()
	}

	// Mark as pending save
	m.pendingSave = true
	m.mu.Unlock()

	m.logger.Debug("Auto-save triggered", zap.Int("content_length", len(content)))

	// Start new timer for debounced save
	m.saveTimer = time.AfterFunc(m.saveInterval, func() {
		if err := m.SaveComposition(content); err != nil {
			m.logger.Error("Auto-save failed", zap.Error(err))
		}
	})
}

// HasPendingSave returns true if there's a pending auto-save
func (m *Manager) HasPendingSave() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.pendingSave
}

// GetCurrentFilePath returns the current composition file path
func (m *Manager) GetCurrentFilePath() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.currentFilePath
}

// GetWorkingDirectory returns the current working directory
func (m *Manager) GetWorkingDirectory() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.workingDirectory
}

// Close closes the manager and performs final save if needed
func (m *Manager) Close(content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("Closing history manager")

	// Cancel any pending save timer
	if m.saveTimer != nil {
		m.saveTimer.Stop()
		m.saveTimer = nil
	}

	// Perform final save if there's a pending save
	if m.pendingSave && m.currentFilePath != "" {
		m.logger.Info("Performing final save before close")
		if err := m.SaveComposition(content); err != nil {
			m.logger.Error("Final save failed", zap.Error(err))
			return err
		}
	}

	return nil
}

// GetAllCompositions returns all compositions from database
func (m *Manager) GetAllCompositions() ([]Composition, error) {
	return m.db.GetAllCompositions()
}

// GetCompositionsByDirectory returns compositions for a specific working directory
func (m *Manager) GetCompositionsByDirectory(workingDir string) ([]Composition, error) {
	return m.db.GetCompositionsByDirectory(workingDir)
}

// GetCompositionsByDateRange returns compositions within a date range
func (m *Manager) GetCompositionsByDateRange(startDate, endDate string) ([]Composition, error) {
	return m.db.GetCompositionsByDateRange(startDate, endDate)
}

// SearchCompositions searches compositions using full-text search
func (m *Manager) SearchCompositions(query string) ([]Composition, error) {
	return m.db.SearchCompositions(query)
}

// DeleteComposition deletes a composition from both markdown and database
func (m *Manager) DeleteComposition(filePath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("Deleting composition", zap.String("file_path", filePath))

	// Delete from database
	if err := m.db.DeleteComposition(filePath); err != nil {
		return fmt.Errorf("failed to delete from database: %w", err)
	}

	// Delete markdown file
	if err := m.storage.DeleteHistoryFile(filePath); err != nil {
		return fmt.Errorf("failed to delete history file: %w", err)
	}

	// Clear current composition if it was the deleted one
	if m.currentFilePath == filePath {
		m.currentFilePath = ""
		m.workingDirectory = ""
		m.pendingSave = false
	}

	return nil
}

// GetCompositionMetadata returns metadata for a composition
func (m *Manager) GetCompositionMetadata(filePath string) (*Composition, error) {
	return m.db.GetCompositionByPath(filePath)
}
