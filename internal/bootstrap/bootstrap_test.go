package bootstrap

import (
	"path/filepath"
	"testing"

	"github.com/kyledavis/prompt-stack/internal/config"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	bootstrap := New(logger)
	if bootstrap == nil {
		t.Error("New() returned nil")
	}
	if bootstrap.logger != logger {
		t.Error("New() did not set logger")
	}
}

func TestBootstrapRun(t *testing.T) {
	tmp := t.TempDir()
	configPath := filepath.Join(tmp, "config.yaml")

	// Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Test bootstrap with existing config
	t.Run("existing config", func(t *testing.T) {
		// Create valid config
		cfg := &config.Config{
			Version:      "1.0.0",
			ClaudeAPIKey: "sk-ant-test",
			Model:        "claude-3-sonnet-20240229",
			VimMode:      false,
		}

		if err := cfg.SaveConfig(configPath); err != nil {
			t.Fatalf("failed to create test config: %v", err)
		}

		bootstrap := New(logger)
		err := bootstrap.Run()

		if err != nil {
			t.Errorf("Bootstrap.Run() error = %v", err)
		}
	})

	// Test bootstrap with version mismatch
	t.Run("version mismatch", func(t *testing.T) {
		// Create config with old version
		cfg := &config.Config{
			Version:      "0.9.0",
			ClaudeAPIKey: "sk-ant-test",
			Model:        "claude-3-sonnet-20240229",
			VimMode:      false,
		}

		if err := cfg.SaveConfig(configPath); err != nil {
			t.Fatalf("failed to create test config: %v", err)
		}

		bootstrap := New(logger)
		err := bootstrap.Run()

		// Should succeed but log warning
		if err != nil {
			t.Errorf("Bootstrap.Run() error = %v", err)
		}
	})

	// Test Run convenience function
	t.Run("Run convenience function", func(t *testing.T) {
		cfg := &config.Config{
			Version:      "1.0.0",
			ClaudeAPIKey: "sk-ant-test",
			Model:        "claude-3-sonnet-20240229",
			VimMode:      false,
		}

		if err := cfg.SaveConfig(configPath); err != nil {
			t.Fatalf("failed to create test config: %v", err)
		}

		err := Run(logger)
		if err != nil {
			t.Errorf("Run() error = %v", err)
		}
	})
}
