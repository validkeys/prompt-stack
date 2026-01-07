// Package config provides application configuration management including
// loading, validation, and persistence of user settings.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

const (
	// Application metadata
	DefaultVersion = "1.0.0"
	AppVersion     = "1.0.0"

	// File paths
	DefaultConfigPath = ".promptstack/config.yaml"

	// API key validation
	APIKeyPrefix    = "sk-ant-"
	APIKeyMinLength = 20

	// Model identifiers
	ModelSonnet = "claude-3-sonnet-20240229"
	ModelOpus   = "claude-3-opus-20240229"
	ModelHaiku  = "claude-3-haiku-20240307"
)

type Config struct {
	Version      string `yaml:"version"`
	ClaudeAPIKey string `yaml:"claude_api_key"`
	Model        string `yaml:"model"`
	VimMode      bool   `yaml:"vim_mode"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.ClaudeAPIKey == "" {
		return fmt.Errorf("claude api key is required")
	}

	if c.Model == "" {
		c.Model = ModelSonnet // Default
	}

	return nil
}

// CheckVersion compares config version with app version
func (c *Config) CheckVersion() error {
	if c.Version == "" {
		c.Version = AppVersion
		return nil
	}

	if c.Version != AppVersion {
		// Log warning but continue
		// Future: Implement migration logic
		return fmt.Errorf("config version %s does not match app version %s", c.Version, AppVersion)
	}

	return nil
}

// GetConfigPath returns the default config file path
func GetConfigPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, DefaultConfigPath), nil
}

// LoadConfig loads configuration from the specified path
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to the specified path
func (c *Config) SaveConfig(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Create directory if needed
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
