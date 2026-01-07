package main

import (
	"fmt"
	"os"

	"github.com/kyledavis/prompt-stack/internal/platform/bootstrap"
	"github.com/kyledavis/prompt-stack/internal/platform/logging"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// run initializes the application and performs bootstrap
func run() error {
	// Initialize logging first
	logger, err := logging.New()
	if err != nil {
		return fmt.Errorf("failed to initialize logging: %w", err)
	}
	defer logger.Sync()

	// Run bootstrap
	if err := bootstrap.Run(logger); err != nil {
		logger.Error("Bootstrap failed", zap.Error(err))
		return fmt.Errorf("bootstrap failed: %w", err)
	}

	logger.Info("PromptStack initialized successfully")
	return nil
}
