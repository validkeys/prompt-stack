package library

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	apperrors "github.com/kyledavis/prompt-stack/internal/errors"
	"github.com/kyledavis/prompt-stack/internal/prompt"
	"go.uber.org/zap"
)

// LoadError represents an error that occurred during library loading
type LoadError struct {
	FilePath string
	Error    error
	Severity string
}

// Library manages the collection of prompts
type Library struct {
	Prompts    map[string]*prompt.Prompt // keyed by file path
	Index      *prompt.LibraryIndex
	LoadErrors []LoadError
	logger     *zap.Logger
}

// Load scans the data directory and loads all prompts
func Load(dataPath string, logger *zap.Logger) (*Library, error) {
	logger.Info("Loading library", zap.String("data_path", dataPath))

	lib := &Library{
		Prompts: make(map[string]*prompt.Prompt),
		Index: &prompt.LibraryIndex{
			Prompts: make(map[string]*prompt.IndexedPrompt),
		},
		logger:     logger,
		LoadErrors: []LoadError{},
	}

	// Walk through data directory
	err := filepath.Walk(dataPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Skip history directory
		if strings.Contains(path, ".history") {
			return nil
		}

		// Load prompt
		p, loadErr := loadPrompt(path, dataPath, logger)
		if loadErr != nil {
			// Log the error and track it
			logger.Error("Failed to load prompt",
				zap.String("path", path),
				zap.Error(loadErr))

			// Determine severity based on error type
			severity := "error"
			if apperrors.IsRetryable(loadErr) {
				severity = "warning"
			}

			lib.LoadErrors = append(lib.LoadErrors, LoadError{
				FilePath: path,
				Error:    loadErr,
				Severity: severity,
			})

			// Continue loading other prompts
			return nil
		}

		// Add to library
		lib.Prompts[path] = p

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan library: %w", err)
	}

	// Log summary of loading results
	logger.Info("Library loaded",
		zap.Int("prompt_count", len(lib.Prompts)),
		zap.Int("error_count", len(lib.LoadErrors)))

	// If there were errors, log them for debugging
	if len(lib.LoadErrors) > 0 {
		for _, loadErr := range lib.LoadErrors {
			logger.Debug("Load error details",
				zap.String("path", loadErr.FilePath),
				zap.String("severity", loadErr.Severity),
				zap.Error(loadErr.Error))
		}
	}

	// Run validation summary and log results
	validationSummary := lib.GetValidationSummary()
	if validationSummary.ErrorCount > 0 {
		logger.Warn("Library validation found errors",
			zap.Int("error_count", validationSummary.ErrorCount),
			zap.Int("warning_count", validationSummary.WarningCount))
	} else if validationSummary.WarningCount > 0 {
		logger.Info("Library validation found warnings",
			zap.Int("warning_count", validationSummary.WarningCount))
	} else {
		logger.Info("Library validation passed")
	}

	return lib, nil
}

// loadPrompt loads a single prompt from file with graceful error handling
func loadPrompt(filePath, dataPath string, logger *zap.Logger) (*prompt.Prompt, error) {
	// Read file with error handling
	content, err := readFileGracefully(filePath, logger)
	if err != nil {
		return nil, err
	}

	// Parse frontmatter with error handling
	metadata, contentWithoutFrontmatter, parseErr := parseFrontmatter(string(content))
	if parseErr != nil {
		// If frontmatter parsing fails, treat entire content as prompt
		logger.Warn("Failed to parse frontmatter, treating as plain content",
			zap.String("path", filePath),
			zap.Error(parseErr))
		contentWithoutFrontmatter = string(content)
		metadata = make(map[string]string)
	}

	// Extract category from path
	category, err := filepath.Rel(dataPath, filePath)
	if err != nil {
		category = "unknown"
	} else {
		// Get parent directory as category
		category = filepath.Dir(category)
		if category == "." {
			category = "root"
		}
	}

	// Parse placeholders
	placeholders := prompt.ParsePlaceholders(contentWithoutFrontmatter)

	// Create prompt
	p := &prompt.Prompt{
		ID:           generateID(filePath),
		Title:        metadata["title"],
		Description:  metadata["description"],
		Tags:         parseTags(metadata["tags"]),
		Category:     category,
		FilePath:     filePath,
		Content:      contentWithoutFrontmatter,
		Metadata:     metadata,
		Placeholders: placeholders,
		ValidationStatus: prompt.ValidationResult{
			IsValid: true,
		},
		UsageStats: prompt.UsageMetadata{
			UseCount: 0,
		},
	}

	// Validate prompt (silent during loading)
	validationResult := Validate(p, logger)
	p.ValidationStatus = validationResult

	// Log validation errors/warnings for debugging
	if !validationResult.IsValid {
		logger.Warn("Prompt validation failed",
			zap.String("path", filePath),
			zap.Int("error_count", len(validationResult.Errors)),
			zap.Int("warning_count", len(validationResult.Warnings)))
	} else if len(validationResult.Warnings) > 0 {
		logger.Debug("Prompt has validation warnings",
			zap.String("path", filePath),
			zap.Int("warning_count", len(validationResult.Warnings)))
	}

	return p, nil
}

// readFileGracefully reads a file with comprehensive error handling
func readFileGracefully(filePath string, logger *zap.Logger) ([]byte, error) {
	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, apperrors.FileError("File not found", err).
				WithDetails(fmt.Sprintf("The file %s does not exist", filePath))
		}
		return nil, apperrors.FileError("Failed to access file", err).
			WithDetails(fmt.Sprintf("Could not access file: %s", filePath))
	}

	// Check file size (1MB limit)
	const maxFileSize = 1 << 20 // 1MB
	if fileInfo.Size() > maxFileSize {
		err := apperrors.FileError("File size exceeds limit", nil).
			WithDetails(fmt.Sprintf("File %s is %.2f MB (max: 1MB)",
				filePath, float64(fileInfo.Size())/(1024*1024)))
		logger.Warn("File size exceeds limit",
			zap.String("path", filePath),
			zap.Int64("size", fileInfo.Size()))
		return nil, err
	}

	// Check file permissions
	if fileInfo.Mode().Perm()&0400 == 0 {
		err := apperrors.FileError("Permission denied", nil).
			WithDetails(fmt.Sprintf("No read permission for file: %s", filePath))
		logger.Warn("Permission denied", zap.String("path", filePath))
		return nil, err
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		// Handle specific error types
		if os.IsPermission(err) {
			return nil, apperrors.FileError("Permission denied", err).
				WithDetails(fmt.Sprintf("Cannot read file: %s", filePath))
		}
		if errors.Is(err, os.ErrClosed) {
			return nil, apperrors.FileError("File closed unexpectedly", err).
				WithDetails(fmt.Sprintf("File handle closed: %s", filePath))
		}

		// Generic file read error
		return nil, apperrors.FileError("Failed to read file", err).
			WithDetails(fmt.Sprintf("Error reading file: %s", filePath))
	}

	// Check if content is empty
	if len(content) == 0 {
		logger.Debug("Empty file loaded", zap.String("path", filePath))
	}

	return content, nil
}

// GetLoadErrors returns all errors that occurred during library loading
func (l *Library) GetLoadErrors() []LoadError {
	return l.LoadErrors
}

// HasLoadErrors returns true if any errors occurred during loading
func (l *Library) HasLoadErrors() bool {
	return len(l.LoadErrors) > 0
}

// GetErrorCount returns the number of errors that occurred during loading
func (l *Library) GetErrorCount() int {
	return len(l.LoadErrors)
}

// GetErrorSummary returns a summary of loading errors
func (l *Library) GetErrorSummary() string {
	if len(l.LoadErrors) == 0 {
		return "No errors"
	}

	errorCount := 0
	warningCount := 0
	for _, err := range l.LoadErrors {
		if err.Severity == "error" {
			errorCount++
		} else {
			warningCount++
		}
	}

	if errorCount > 0 && warningCount > 0 {
		return fmt.Sprintf("%d errors, %d warnings", errorCount, warningCount)
	} else if errorCount > 0 {
		return fmt.Sprintf("%d errors", errorCount)
	} else {
		return fmt.Sprintf("%d warnings", warningCount)
	}
}

// ValidationSummary represents validation results summary
type ValidationSummary struct {
	ErrorCount   int
	WarningCount int
	TotalFiles   int
}

// GetValidationSummary returns a summary of validation results across all prompts
func (l *Library) GetValidationSummary() ValidationSummary {
	summary := ValidationSummary{
		TotalFiles: len(l.Prompts),
	}

	for _, p := range l.Prompts {
		if !p.ValidationStatus.IsValid {
			summary.ErrorCount += len(p.ValidationStatus.Errors)
		}
		summary.WarningCount += len(p.ValidationStatus.Warnings)
	}

	return summary
}

// HasValidationErrors returns true if any prompts have validation errors
func (l *Library) HasValidationErrors() bool {
	for _, p := range l.Prompts {
		if !p.ValidationStatus.IsValid {
			return true
		}
	}
	return false
}

// HasValidationWarnings returns true if any prompts have validation warnings
func (l *Library) HasValidationWarnings() bool {
	for _, p := range l.Prompts {
		if len(p.ValidationStatus.Warnings) > 0 {
			return true
		}
	}
	return false
}

// parseFrontmatter extracts YAML frontmatter from content
func parseFrontmatter(content string) (map[string]string, string, error) {
	// Check for frontmatter markers
	if !strings.HasPrefix(content, "---") {
		return nil, content, nil
	}

	// Find end of frontmatter
	endIdx := strings.Index(content[3:], "---")
	if endIdx == -1 {
		return nil, content, fmt.Errorf("unclosed frontmatter")
	}

	frontmatter := content[3 : endIdx+3]
	remainingContent := content[endIdx+6:] // Skip both --- markers

	// Parse YAML frontmatter
	metadata := make(map[string]string)
	lines := strings.Split(frontmatter, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first colon
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		metadata[key] = value
	}

	return metadata, remainingContent, nil
}

// parseTags parses comma-separated tags
func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}

	tags := strings.Split(tagsStr, ",")
	var result []string
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result
}

// AddPrompt adds a new prompt to the library
func (l *Library) AddPrompt(p *prompt.Prompt) {
	l.Prompts[p.FilePath] = p

	// Add to index
	indexed := &prompt.IndexedPrompt{
		PromptID:      p.ID,
		Title:         p.Title,
		Tags:          p.Tags,
		Category:      p.Category,
		WordFrequency: prompt.ExtractKeywords(p.Content),
		LastUsed:      p.UsageStats.LastUsed,
		UseCount:      p.UsageStats.UseCount,
	}

	l.Index.Prompts[p.FilePath] = indexed

	l.logger.Info("Prompt added to library",
		zap.String("title", p.Title),
		zap.String("path", p.FilePath))
}

// generateID generates a unique ID for a prompt
func generateID(filePath string) string {
	// Use file path as ID for now
	// Could use UUID or hash in the future
	return filePath
}
