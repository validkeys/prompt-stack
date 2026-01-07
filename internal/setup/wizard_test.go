package setup

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kyledavis/prompt-stack/internal/config"
	"go.uber.org/zap"
)

func TestNewWizard(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	configPath := "/tmp/config.yaml"
	wizard := NewWizard(configPath, logger)

	if wizard == nil {
		t.Error("NewWizard() returned nil")
	}
	if wizard.configPath != configPath {
		t.Errorf("configPath = %s, want %s", wizard.configPath, configPath)
	}
	if wizard.logger != logger {
		t.Error("logger not set correctly")
	}
	if wizard.reader == nil {
		t.Error("reader not initialized")
	}
}

func TestWizardRun(t *testing.T) {
	tmp := t.TempDir()
	configPath := filepath.Join(tmp, "config.yaml")

	// Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Simulate user input
	input := "sk-ant-test123456789\n1\ny\ny\n"
	reader := bufio.NewReader(strings.NewReader(input))

	wizard := &Wizard{
		configPath: configPath,
		logger:     logger,
		reader:     reader,
	}

	// Run wizard
	err = wizard.Run()
	if err != nil {
		t.Errorf("Wizard.Run() error = %v", err)
	}

	// Verify config was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file not created")
	}

	// Load and verify config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Errorf("failed to load config: %v", err)
	}

	if cfg.ClaudeAPIKey != "sk-ant-test123456789" {
		t.Errorf("API key = %s, want sk-ant-test123456789", cfg.ClaudeAPIKey)
	}

	if cfg.Model != "claude-3-sonnet-20240229" {
		t.Errorf("Model = %s, want claude-3-sonnet-20240229", cfg.Model)
	}

	if !cfg.VimMode {
		t.Error("Vim mode should be true")
	}
}

func TestPromptAPIKey(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid api key",
			input:   "sk-ant-test123\n",
			want:    "sk-ant-test123",
			wantErr: false,
		},
		{
			name:    "invalid api key format",
			input:   "invalid-key\nsk-ant-test123\n",
			want:    "sk-ant-test123",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, _ := zap.NewDevelopment()
			defer logger.Sync()

			reader := bufio.NewReader(strings.NewReader(tt.input))
			wizard := &Wizard{
				configPath: "/tmp/config.yaml",
				logger:     logger,
				reader:     reader,
			}

			got, err := wizard.promptAPIKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("promptAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("promptAPIKey() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestPromptModel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "default model",
			input: "\n",
			want:  "claude-3-sonnet-20240229",
		},
		{
			name:  "sonnet model",
			input: "1\n",
			want:  "claude-3-sonnet-20240229",
		},
		{
			name:  "opus model",
			input: "2\n",
			want:  "claude-3-opus-20240229",
		},
		{
			name:  "haiku model",
			input: "3\n",
			want:  "claude-3-haiku-20240307",
		},
		{
			name:  "invalid choice defaults",
			input: "99\n",
			want:  "claude-3-sonnet-20240229",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, _ := zap.NewDevelopment()
			defer logger.Sync()

			reader := bufio.NewReader(strings.NewReader(tt.input))
			wizard := &Wizard{
				configPath: "/tmp/config.yaml",
				logger:     logger,
				reader:     reader,
			}

			got, err := wizard.promptModel()
			if err != nil {
				t.Errorf("promptModel() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("promptModel() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestPromptVimMode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "enable vim mode",
			input: "y\n",
			want:  true,
		},
		{
			name:  "enable vim mode yes",
			input: "yes\n",
			want:  true,
		},
		{
			name:  "disable vim mode",
			input: "n\n",
			want:  false,
		},
		{
			name:  "disable vim mode no",
			input: "no\n",
			want:  false,
		},
		{
			name:  "default disable",
			input: "\n",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, _ := zap.NewDevelopment()
			defer logger.Sync()

			reader := bufio.NewReader(strings.NewReader(tt.input))
			wizard := &Wizard{
				configPath: "/tmp/config.yaml",
				logger:     logger,
				reader:     reader,
			}

			got, err := wizard.promptVimMode()
			if err != nil {
				t.Errorf("promptVimMode() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("promptVimMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskAPIKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			name: "short key",
			key:  "sk-ant-",
			want: "***",
		},
		{
			name: "normal key",
			key:  "sk-ant-apikey123456789",
			want: "sk-ant-api...",
		},
		{
			name: "exact length",
			key:  "sk-ant-123",
			want: "***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maskAPIKey(tt.key)
			if got != tt.want {
				t.Errorf("maskAPIKey() = %s, want %s", got, tt.want)
			}
		})
	}
}
