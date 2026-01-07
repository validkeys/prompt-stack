package main

import (
	"fmt"
	"os"

	"github.com/kyledavis/prompt-stack/internal/bootstrap"
	"github.com/kyledavis/prompt-stack/internal/logging"
	"go.uber.org/zap"
)

func main() {
	// Initialize logging first
	logger, err := logging.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logging: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Run bootstrap
	if err := bootstrap.Run(logger); err != nil {
		logger.Error("Bootstrap failed", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("PromptStack initialized successfully")
}
