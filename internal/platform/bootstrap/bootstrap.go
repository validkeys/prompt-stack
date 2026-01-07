// Package bootstrap provides application initialization and startup logic.
package bootstrap

import (
	"fmt"
	"os"

	"github.com/kyledavis/prompt-stack/internal/config"
	"go.uber.org/zap"
)

type Bootstrap struct {
	logger *zap.Logger
}

// New creates a new bootstrap instance
func New(logger *zap.Logger) *Bootstrap {
	return &Bootstrap{
		logger: logger,
	}
}

// Run executes bootstrap process
func (b *Bootstrap) Run() error {
	b.logger.Info("Starting PromptStack bootstrap")

	// Get config path
	configPath, err := config.GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		b.logger.Info("Config not found, running setup wizard")

		// Run setup wizard
		wizard := config.NewWizard(configPath, b.logger)
		if err := wizard.Run(); err != nil {
			return fmt.Errorf("setup wizard failed: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check config: %w", err)
	} else {
		// Load existing config
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Check version
		if err := cfg.CheckVersion(); err != nil {
			b.logger.Warn("Version mismatch", zap.Error(err))
		}

		b.logger.Info("Config loaded successfully",
			zap.String("version", cfg.Version),
			zap.String("model", cfg.Model))
	}

	b.logger.Info("Bootstrap completed successfully")
	return nil
}

// Run is a convenience function that creates and runs bootstrap
func Run(logger *zap.Logger) error {
	bootstrap := New(logger)
	return bootstrap.Run()
}
