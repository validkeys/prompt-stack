package main

import (
	"fmt"
	"os"

	"github.com/kyledavis/prompt-stack/internal/bootstrap"
	"github.com/kyledavis/prompt-stack/internal/logging"
	"go.uber.org/zap"
)

const (
	version = "0.1.0"
)

func main() {
	// Initialize logger
	logger, err := logging.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting PromptStack", zap.String("version", version))

	// Bootstrap's application
	app, err := bootstrap.Initialize(version, StarterPrompts, logger)
	if err != nil {
		logger.Error("Failed to initialize application", zap.Error(err))
		fmt.Fprintf(os.Stderr, "Failed to initialize: %v\n", err)
		os.Exit(1)
	}

	// Run application
	if err := app.Run(); err != nil {
		logger.Error("Application error", zap.Error(err))
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	logger.Info("PromptStack exited successfully")
}
