package bootstrap

import (
	"fmt"
	"io/fs"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyledavis/prompt-stack/internal/ai"
	"github.com/kyledavis/prompt-stack/internal/commands"
	"github.com/kyledavis/prompt-stack/internal/config"
	"github.com/kyledavis/prompt-stack/internal/history"
	"github.com/kyledavis/prompt-stack/internal/library"
	"github.com/kyledavis/prompt-stack/internal/prompt"
	"github.com/kyledavis/prompt-stack/internal/setup"
	"github.com/kyledavis/prompt-stack/ui/app"
	"go.uber.org/zap"
)

// App represents the initialized application
type App struct {
	Config *config.Config
	Logger *zap.Logger
}

// Initialize sets up application by loading config, running setup if needed, etc.
func Initialize(version string, starterFS fs.FS, logger *zap.Logger) (*App, error) {
	// Get config path
	configPath, err := config.GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// Check if config exists
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		// First run - run setup wizard
		logger.Info("First run detected, starting setup wizard")
		cfg, err := setup.RunWizard(version, logger)
		if err != nil {
			return nil, fmt.Errorf("setup wizard failed: %w", err)
		}

		// Initialize data directory
		if err := InitializeDataDirectory(logger); err != nil {
			return nil, fmt.Errorf("failed to initialize data directory: %w", err)
		}

		// Extract starter prompts
		dataPath, err := config.GetDataPath()
		if err != nil {
			return nil, fmt.Errorf("failed to get data path: %w", err)
		}

		if err := ExtractStarterPrompts(dataPath, starterFS, logger); err != nil {
			return nil, fmt.Errorf("failed to extract starter prompts: %w", err)
		}

		// Initialize database
		dbPath, err := config.GetDatabasePath()
		if err != nil {
			return nil, fmt.Errorf("failed to get database path: %w", err)
		}

		if _, err := history.Initialize(dbPath, logger); err != nil {
			return nil, fmt.Errorf("failed to initialize database: %w", err)
		}

		// Save config
		if err := cfg.Save(configPath); err != nil {
			return nil, fmt.Errorf("failed to save config: %w", err)
		}

		logger.Info("Setup completed successfully")
		return &App{
			Config: cfg,
			Logger: logger,
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to check config file: %w", err)
	}

	// Load existing config with error recovery
	cfg, err := loadConfigWithRecovery(configPath, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Check for version upgrade
	if cfg.Version != version {
		logger.Info("Version upgrade detected", zap.String("old_version", cfg.Version), zap.String("new_version", version))
		// TODO: Handle version upgrade (extract new starter prompts)
		cfg.Version = version
		if err := cfg.Save(configPath); err != nil {
			logger.Warn("Failed to update version in config", zap.Error(err))
		}
	}

	return &App{
		Config: cfg,
		Logger: logger,
	}, nil
}

// loadConfigWithRecovery loads config with error recovery
// If config loading or validation fails, it offers to reset's config
func loadConfigWithRecovery(configPath string, logger *zap.Logger) (*config.Config, error) {
	// Try to load config
	cfg, err := config.Load(configPath)
	if err != nil {
		// Config file exists but cannot be loaded
		logger.Warn("Failed to load config file", zap.Error(err))

		// Attempt to reset config
		if resetErr := attemptConfigReset(configPath, logger); resetErr != nil {
			return nil, fmt.Errorf("config load failed and reset failed: %w (reset error: %v)", err, resetErr)
		}

		// Load of newly created default config
		cfg, err = config.Load(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load reset config: %w", err)
		}

		logger.Info("Config reset to defaults")
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		logger.Warn("Config validation failed", zap.Error(err))

		// Attempt to reset config
		if resetErr := attemptConfigReset(configPath, logger); resetErr != nil {
			return nil, fmt.Errorf("config validation failed and reset failed: %w (reset error: %v)", err, resetErr)
		}

		// Load of newly created default config
		cfg, err = config.Load(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load reset config: %w", err)
		}

		logger.Info("Config reset to defaults due to validation failure")
	}

	return cfg, nil
}

// attemptConfigReset attempts to reset's config file
// It creates a backup of existing config before resetting
func attemptConfigReset(configPath string, logger *zap.Logger) error {
	logger.Info("Attempting to reset config", zap.String("path", configPath))

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// No config file exists, just create default
		cfg := config.DefaultConfig()
		if err := cfg.Save(configPath); err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
		return nil
	}

	// Config file exists, create backup first
	backupPath, err := config.BackupConfig(configPath)
	if err != nil {
		logger.Warn("Failed to backup config before reset", zap.Error(err))
		// Continue with reset even if backup fails
	} else {
		logger.Info("Config backed up", zap.String("backup_path", backupPath))
	}

	// Reset config to defaults
	if err := config.ResetConfig(configPath); err != nil {
		return fmt.Errorf("failed to reset config: %w", err)
	}

	return nil
}

// Run starts the application
func (a *App) Run() error {
	// Get working directory
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Get data path for library
	dataPath, err := config.GetDataPath()
	if err != nil {
		return fmt.Errorf("failed to get data path: %w", err)
	}

	// Load library
	lib, err := library.Load(dataPath, a.Logger)
	if err != nil {
		a.Logger.Warn("Failed to load library, continuing with empty library", zap.Error(err))
		lib = &library.Library{
			Prompts: make(map[string]*prompt.Prompt),
		}
	}

	// Create command registry
	registry := commands.NewRegistry()
	commands.RegisterCoreCommands(registry)

	// Initialize history manager
	dbPath, err := config.GetDatabasePath()
	if err != nil {
		return fmt.Errorf("failed to get database path: %w", err)
	}

	db, err := history.Initialize(dbPath, a.Logger)
	if err != nil {
		return fmt.Errorf("failed to initialize history database: %w", err)
	}

	// Get history directory path
	historyDir, err := config.GetHistoryPath()
	if err != nil {
		return fmt.Errorf("failed to get history path: %w", err)
	}

	storage, err := history.NewStorage(historyDir, a.Logger)
	if err != nil {
		return fmt.Errorf("failed to initialize history storage: %w", err)
	}

	historyMgr := history.NewManager(db, storage, a.Logger)

	// Initialize AI client if API key is configured
	var aiClient *ai.Client
	if a.Config.ClaudeAPIKey != "" {
		aiConfig := ai.Config{
			APIKey:     a.Config.ClaudeAPIKey,
			Model:      a.Config.Model,
			MaxRetries: 3,
			Timeout:    60 * 1000000000, // 60 seconds in nanoseconds
			Logger:     a.Logger,
		}

		aiClient, err = ai.NewClient(aiConfig)
		if err != nil {
			a.Logger.Warn("Failed to initialize AI client, AI features will be disabled", zap.Error(err))
			aiClient = nil
		} else {
			a.Logger.Info("AI client initialized successfully", zap.String("model", a.Config.Model))
		}
	} else {
		a.Logger.Info("No API key configured, AI features will be disabled")
	}

	// Initialize context selector
	contextSelector := ai.NewContextSelector()

	// Initialize TUI model with library, commands, history manager, and AI client
	tuiModel := app.NewWithDependencies(workingDir, lib, registry, historyMgr, a.Config.VimMode, aiClient, contextSelector)

	// Create Bubble Tea program with options
	program := tea.NewProgram(
		tuiModel,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	return nil
}
