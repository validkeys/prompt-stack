package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyledavis/prompt-stack/internal/platform/bootstrap"
	"github.com/kyledavis/prompt-stack/internal/platform/logging"
	"github.com/kyledavis/prompt-stack/ui/app"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// run initializes the application and performs bootstrap, then launches the TUI
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

	// Create the app model
	appModel := app.New()

	// Log TUI startup
	logger.Info("Starting TUI")

	// Create and run the Bubble Tea program with alt screen
	program := tea.NewProgram(appModel, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		logger.Error("TUI error", zap.Error(err))
		return fmt.Errorf("tui error: %w", err)
	}

	// Log TUI shutdown
	logger.Info("TUI shutdown complete")

	return nil
}
