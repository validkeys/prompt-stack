package editor

import (
	"os"
	"path/filepath"
)

// FileManager handles file save/load operations
type FileManager struct {
	filepath string
	modified bool
}

// NewFileManager creates a new file I/O manager
func NewFileManager(path string) FileManager {
	return FileManager{
		filepath: path,
		modified: false,
	}
}

// Load reads file content
func (m FileManager) Load() (string, error) {
	if m.filepath == "" {
		return "", nil
	}

	content, err := os.ReadFile(m.filepath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// Save writes content to file
func (m FileManager) Save(content string) error {
	if m.filepath == "" {
		return nil
	}

	// Ensure directory exists
	dir := filepath.Dir(m.filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(m.filepath, []byte(content), 0644)
}

// MarkModified marks file as modified
func (m FileManager) MarkModified() FileManager {
	newModel := m
	newModel.modified = true
	return newModel
}

// ClearModified clears the modified flag
func (m FileManager) ClearModified() FileManager {
	newModel := m
	newModel.modified = false
	return newModel
}

// IsModified returns true if file has unsaved changes
func (m FileManager) IsModified() bool {
	return m.modified
}

// Path returns the file path
func (m FileManager) Path() string {
	return m.filepath
}

// SetPath sets a new file path
func (m FileManager) SetPath(path string) FileManager {
	newModel := m
	newModel.filepath = path
	newModel.modified = false
	return newModel
}

// HasPath returns true if a file path is set
func (m FileManager) HasPath() bool {
	return m.filepath != ""
}
