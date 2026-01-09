package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"

	"go.uber.org/zap"
)

// Storage handles reading and writing prompt files from disk.
type Storage struct {
	logger *zap.Logger
}

// NewStorage creates a new prompt storage instance.
func NewStorage(logger *zap.Logger) *Storage {
	return &Storage{
		logger: logger,
	}
}

// Load reads a prompt file from disk, parses frontmatter, and returns a Prompt structure.
func (s *Storage) Load(path string) (*Prompt, error) {
	if !validateMarkdownExtension(path) {
		return nil, fmt.Errorf("file %s must have .md extension", path)
	}

	// Check file size before reading
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	if info.Size() > MaxFileSize {
		return nil, fmt.Errorf("file %s exceeds maximum size of %d bytes", path, MaxFileSize)
	}

	// Read file content
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// Validate UTF-8 encoding
	if !isValidUTF8(data) {
		return nil, fmt.Errorf("file %s contains invalid UTF-8 encoding", path)
	}

	// Parse frontmatter
	metadata, content, err := ParseFrontmatter(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter from %s: %w", path, err)
	}

	prompt := &Prompt{
		Metadata: metadata,
		Content:  content,
	}

	s.logger.Debug("loaded prompt file", zap.String("path", path))

	return prompt, nil
}

// Save writes a prompt file to disk with frontmatter.
// Creates parent directories if they don't exist.
// Uses atomic write pattern (write to temp file then rename).
func (s *Storage) Save(path string, prompt *Prompt) error {
	if !validateMarkdownExtension(path) {
		return fmt.Errorf("file %s must have .md extension", path)
	}

	if err := ValidateMetadata(prompt.Metadata); err != nil {
		return fmt.Errorf("invalid metadata: %w", err)
	}

	// Create parent directories if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Format content with frontmatter
	content := FormatFrontmatter(prompt.Metadata, prompt.Content)

	// Atomic write: write to temp file, then rename
	tempPath := path + ".tmp"
	if err := os.WriteFile(tempPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write temp file %s: %w", tempPath, err)
	}

	// Rename temp file to final path (atomic on most systems)
	if err := os.Rename(tempPath, path); err != nil {
		// Clean up temp file if rename fails
		os.Remove(tempPath)
		return fmt.Errorf("failed to rename temp file to %s: %w", path, err)
	}

	s.logger.Debug("saved prompt file", zap.String("path", path))

	return nil
}

// isValidUTF8 checks if a byte slice contains valid UTF-8 encoding.
func isValidUTF8(data []byte) bool {
	return utf8.Valid(data)
}

// validateMarkdownExtension checks if the file path has a .md extension.
func validateMarkdownExtension(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".md"
}
