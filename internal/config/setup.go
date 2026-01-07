package setup

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kyledavis/prompt-stack/internal/config"
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
	cfg := &config.Config{
		Version:      config.DefaultVersion,
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
		if !strings.HasPrefix(apiKey, "sk-ant-") {
			fmt.Println("Invalid API key format. API keys should start with 'sk-ant-'")
			continue
		}

		return apiKey, nil
	}
}

func (w *Wizard) promptModel() (string, error) {
	fmt.Println("\nSelect Claude model:")
	fmt.Println("1. claude-3-sonnet-20240229 (default)")
	fmt.Println("2. claude-3-opus-20240229")
	fmt.Println("3. claude-3-haiku-20240307")
	fmt.Print("Enter choice [1-3]: ")

	input, err := w.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	choice := strings.TrimSpace(input)

	switch choice {
	case "1", "":
		return "claude-3-sonnet-20240229", nil
	case "2":
		return "claude-3-opus-20240229", nil
	case "3":
		return "claude-3-haiku-20240307", nil
	default:
		fmt.Println("Invalid choice. Using default model.")
		return "claude-3-sonnet-20240229", nil
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

func (w *Wizard) showSummary(cfg *config.Config) {
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
