package bootstrap

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kyledavis/prompt-stack/internal/config"
	"go.uber.org/zap"
)

// ExtractStarterPrompts extracts starter prompts to the data directory
// Only extracts prompts that don't already exist (version-aware)
func ExtractStarterPrompts(dataPath string, starterFS fs.FS, logger *zap.Logger) error {
	logger.Info("Extracting starter prompts", zap.String("data_path", dataPath))

	// Walk through embedded starter prompts
	err := fs.WalkDir(starterFS, "starter-prompts", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Get relative path from starter-prompts/
		relPath, err := filepath.Rel("starter-prompts", path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Destination path in data directory
		destPath := filepath.Join(dataPath, relPath)

		// Check if file already exists
		if _, err := os.Stat(destPath); err == nil {
			logger.Debug("Prompt already exists, skipping", zap.String("path", relPath))
			return nil
		}

		// Create destination directory if needed
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destDir, err)
		}

		// Read embedded file
		content, err := fs.ReadFile(starterFS, path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Write to destination
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}

		logger.Info("Extracted starter prompt", zap.String("path", relPath))
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to extract starter prompts: %w", err)
	}

	return nil
}

// InitializeDataDirectory creates the data directory structure
func InitializeDataDirectory(logger *zap.Logger) error {
	dataPath, err := config.GetDataPath()
	if err != nil {
		return fmt.Errorf("failed to get data path: %w", err)
	}

	// Create data directory
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create history directory
	historyPath, err := config.GetHistoryPath()
	if err != nil {
		return fmt.Errorf("failed to get history path: %w", err)
	}

	if err := os.MkdirAll(historyPath, 0755); err != nil {
		return fmt.Errorf("failed to create history directory: %w", err)
	}

	logger.Info("Data directory initialized", zap.String("data_path", dataPath))
	return nil
}
