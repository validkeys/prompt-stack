package setup

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kyledavis/prompt-stack/internal/config"
	"go.uber.org/zap"
)

// RunWizard runs the interactive first-run setup wizard
func RunWizard(version string, logger *zap.Logger) (*config.Config, error) {
	fmt.Println("Welcome to PromptStack!")
	fmt.Println("Let's set up your configuration.")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// Get Claude API key
	apiKey, err := prompt(reader, "Enter your Claude API key: ", true)
	if err != nil {
		return nil, fmt.Errorf("failed to read API key: %w", err)
	}

	// Select model
	fmt.Println("\nAvailable Claude models:")
	fmt.Println("1. claude-3-sonnet-20240229 (default)")
	fmt.Println("2. claude-3-opus-20240229")
	fmt.Println("3. claude-3-haiku-20240307")

	modelChoice, err := prompt(reader, "Select model (1-3, default: 1): ", false)
	if err != nil {
		return nil, fmt.Errorf("failed to read model choice: %w", err)
	}

	var model string
	switch strings.TrimSpace(modelChoice) {
	case "2":
		model = "claude-3-opus-20240229"
	case "3":
		model = "claude-3-haiku-20240307"
	default:
		model = "claude-3-sonnet-20240229"
	}

	// Ask about vim mode
	vimModeStr, err := prompt(reader, "Enable vim mode? (y/N): ", false)
	if err != nil {
		return nil, fmt.Errorf("failed to read vim mode choice: %w", err)
	}

	vimMode := strings.ToLower(strings.TrimSpace(vimModeStr)) == "y" || strings.ToLower(strings.TrimSpace(vimModeStr)) == "yes"

	// Create config
	cfg := &config.Config{
		Version:      version,
		ClaudeAPIKey: strings.TrimSpace(apiKey),
		Model:        model,
		VimMode:      vimMode,
	}

	logger.Info("Setup wizard completed",
		zap.String("model", model),
		zap.Bool("vim_mode", vimMode),
	)

	return cfg, nil
}

// prompt prompts the user for input
func prompt(reader *bufio.Reader, message string, required bool) (string, error) {
	for {
		fmt.Print(message)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		input = strings.TrimSpace(input)

		if required && input == "" {
			fmt.Println("This field is required. Please try again.")
			continue
		}

		return input, nil
	}
}
