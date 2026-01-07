// Package config provides application configuration management including
// loading, validation, and persistence of user settings.
package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

type Wizard struct {
	configPath string
	logger     *zap.Logger
	reader     *bufio.Reader
}

// NewWizard creates a new setup wizard
func NewWizard(configPath string, logger *zap.Logger) *Wizard {
	return &Wizard{
		configPath: configPath,
		logger:     logger,
		reader:     bufio.NewReader(os.Stdin),
	}
}

// Run executes setup wizard
func (w *Wizard) Run() error {
	fmt.Println("Welcome to PromptStack!")
	fmt.Println("Let's configure your application.")
	fmt.Println()

	// Prompt for API key
	apiKey, err := w.promptAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	// Prompt for model selection
	model, err := w.promptModel()
	if err != nil {
		return fmt.Errorf("failed to get model: %w", err)
	}

	// Prompt for vim mode
	vimMode, err := w.promptVimMode()
	if err != nil {
		return fmt.Errorf("failed to get vim mode preference: %w", err)
	}

	// Create config
	cfg := &Config{
		Version:      DefaultVersion,
		ClaudeAPIKey: apiKey,
		Model:        model,
		VimMode:      vimMode,
	}

	// Show summary
	w.showSummary(cfg)

	// Confirm and save
	if !w.confirm() {
		fmt.Println("Setup cancelled.")
		return nil
	}

	// Save config
	if err := cfg.SaveConfig(w.configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	w.logger.Info("Configuration saved", zap.String("path", w.configPath))
	fmt.Println("Configuration saved successfully!")

	return nil
}

func (w *Wizard) promptAPIKey() (string, error) {
	for {
		fmt.Print("Enter your Claude API key: ")
		input, err := w.reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		apiKey := strings.TrimSpace(input)

		// Validate API key format
		if err := validateAPIKey(apiKey); err != nil {
			fmt.Printf("Invalid API key: %v\n", err)
			continue
		}

		return apiKey, nil
	}
}

// validateAPIKey checks if the API key meets format requirements
func validateAPIKey(apiKey string) error {
	if !strings.HasPrefix(apiKey, APIKeyPrefix) {
		return fmt.Errorf("API key must start with '%s'", APIKeyPrefix)
	}

	if len(apiKey) < APIKeyMinLength {
		return fmt.Errorf("API key must be at least %d characters", APIKeyMinLength)
	}

	// Check for valid characters (alphanumeric and hyphens)
	for _, r := range apiKey {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-') {
			return fmt.Errorf("API key contains invalid character '%c'", r)
		}
	}

	return nil
}

func (w *Wizard) promptModel() (string, error) {
	fmt.Println("\nSelect Claude model:")
	fmt.Printf("1. %s (default)\n", ModelSonnet)
	fmt.Printf("2. %s\n", ModelOpus)
	fmt.Printf("3. %s\n", ModelHaiku)
	fmt.Print("Enter choice [1-3]: ")

	input, err := w.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	choice := strings.TrimSpace(input)

	switch choice {
	case "1", "":
		return ModelSonnet, nil
	case "2":
		return ModelOpus, nil
	case "3":
		return ModelHaiku, nil
	default:
		fmt.Println("Invalid choice. Using default model.")
		return ModelSonnet, nil
	}
}

func (w *Wizard) promptVimMode() (bool, error) {
	fmt.Print("\nEnable vim mode? [y/N]: ")

	input, err := w.reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	choice := strings.TrimSpace(strings.ToLower(input))
	return choice == "y" || choice == "yes", nil
}

func (w *Wizard) showSummary(cfg *Config) {
	fmt.Println("\nConfiguration Summary:")
	fmt.Printf("  API Key: %s\n", maskAPIKey(cfg.ClaudeAPIKey))
	fmt.Printf("  Model: %s\n", cfg.Model)
	fmt.Printf("  Vim Mode: %v\n", cfg.VimMode)
}

func (w *Wizard) confirm() bool {
	fmt.Print("\nSave this configuration? [Y/n]: ")

	input, err := w.reader.ReadString('\n')
	if err != nil {
		return false
	}

	choice := strings.TrimSpace(strings.ToLower(input))
	return choice == "" || choice == "y" || choice == "yes"
}

func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 10 {
		return "***"
	}
	return apiKey[:10] + "..."
}
