package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

// Storage handles prompt file operations
type Storage struct {
	dataPath string
	logger   *zap.Logger
}

// NewStorage creates a new prompt storage instance
func NewStorage(dataPath string, logger *zap.Logger) *Storage {
	return &Storage{
		dataPath: dataPath,
		logger:   logger,
	}
}

// SavePrompt saves a prompt to the library with YAML frontmatter
func (s *Storage) SavePrompt(data *PromptData) (string, error) {
	// Validate required fields
	if data.Title == "" {
		return "", fmt.Errorf("title is required")
	}
	if data.Category == "" {
		return "", fmt.Errorf("category is required")
	}
	if data.Content == "" {
		return "", fmt.Errorf("content is required")
	}

	// Generate filename from title
	filename := generateFilenameFromTitle(data.Title)

	// Build full file path
	categoryPath := filepath.Join(s.dataPath, data.Category)
	filePath := filepath.Join(categoryPath, filename)

	// Ensure category directory exists
	if err := os.MkdirAll(categoryPath, 0755); err != nil {
		s.logger.Error("Failed to create category directory",
			zap.String("category", data.Category),
			zap.Error(err))
		return "", fmt.Errorf("failed to create category directory: %w", err)
	}

	// Generate YAML frontmatter
	frontmatter := s.generateFrontmatter(data)

	// Build complete file content
	fileContent := fmt.Sprintf("%s\n%s", frontmatter, data.Content)

	// Write file
	if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
		s.logger.Error("Failed to save prompt file",
			zap.String("path", filePath),
			zap.Error(err))
		return "", fmt.Errorf("failed to save prompt file: %w", err)
	}

	s.logger.Info("Prompt saved successfully",
		zap.String("title", data.Title),
		zap.String("category", data.Category),
		zap.String("path", filePath))

	return filePath, nil
}

// generateFrontmatter generates YAML frontmatter from prompt data
func (s *Storage) generateFrontmatter(data *PromptData) string {
	var builder strings.Builder

	builder.WriteString("---\n")
	builder.WriteString(fmt.Sprintf("title: \"%s\"\n", escapeYAMLString(data.Title)))

	if data.Description != "" {
		builder.WriteString(fmt.Sprintf("description: \"%s\"\n", escapeYAMLString(data.Description)))
	}

	if len(data.Tags) > 0 {
		tagsStr := strings.Join(data.Tags, "\", \"")
		builder.WriteString(fmt.Sprintf("tags: [\"%s\"]\n", tagsStr))
	}

	builder.WriteString("---\n")

	return builder.String()
}

// escapeYAMLString escapes special characters in YAML strings
func escapeYAMLString(s string) string {
	// Escape backslashes and quotes
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// generateFilenameFromTitle generates a kebab-case filename from title
func generateFilenameFromTitle(title string) string {
	// Convert to lowercase
	filename := strings.ToLower(title)

	// Replace spaces with hyphens
	filename = strings.ReplaceAll(filename, " ", "-")

	// Remove special characters (keep only alphanumeric and hyphens)
	var result []rune
	for _, r := range filename {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result = append(result, r)
		}
	}
	filename = string(result)

	// Remove consecutive hyphens
	for strings.Contains(filename, "--") {
		filename = strings.ReplaceAll(filename, "--", "-")
	}

	// Trim hyphens from start/end
	filename = strings.Trim(filename, "-")

	// Add .md extension
	if filename == "" {
		filename = "untitled"
	}

	return fmt.Sprintf("%s.md", filename)
}

// DeletePrompt deletes a prompt file
func (s *Storage) DeletePrompt(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		s.logger.Error("Failed to delete prompt file",
			zap.String("path", filePath),
			zap.Error(err))
		return fmt.Errorf("failed to delete prompt file: %w", err)
	}

	s.logger.Info("Prompt deleted successfully", zap.String("path", filePath))
	return nil
}

// PromptExists checks if a prompt file exists
func (s *Storage) PromptExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// GetDataPath returns the data path
func (s *Storage) GetDataPath() string {
	return s.dataPath
}
