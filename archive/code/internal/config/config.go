package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Version      string `yaml:"version"`
	ClaudeAPIKey string `yaml:"claude_api_key"`
	Model        string `yaml:"model"`
	VimMode      bool   `yaml:"vim_mode"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Version: "",
		Model:   "claude-3-sonnet-20240229",
		VimMode: false,
	}
}

// Load reads the configuration from the specified path
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to the specified path
func (c *Config) Save(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with restricted permissions (0600 - user read/write only)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.ClaudeAPIKey == "" {
		return fmt.Errorf("claude_api_key is required")
	}
	if c.Model == "" {
		return fmt.Errorf("model is required")
	}
	return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".promptstack", "config.yaml"), nil
}

// GetDataPath returns the path to the data directory
func GetDataPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".promptstack", "data"), nil
}

// GetHistoryPath returns the path to the history directory
func GetHistoryPath() (string, error) {
	dataPath, err := GetDataPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataPath, ".history"), nil
}

// GetDatabasePath returns the path to the SQLite database
func GetDatabasePath() (string, error) {
	dataPath, err := GetDataPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataPath, "history.db"), nil
}

// ResetConfig resets the configuration to defaults
// It backs up the existing config file before resetting
func ResetConfig(configPath string) error {
	// Check if config file exists
	if _, err := os.Stat(configPath); err == nil {
		// Create backup
		backupPath := configPath + ".backup"
		if err := os.Rename(configPath, backupPath); err != nil {
			return fmt.Errorf("failed to backup config file: %w", err)
		}
	}

	// Create default config
	cfg := DefaultConfig()
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save default config: %w", err)
	}

	return nil
}

// BackupConfig creates a backup of the existing config file
func BackupConfig(configPath string) (string, error) {
	backupPath := configPath + ".backup"

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %w", err)
	}

	// Write backup
	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return "", fmt.Errorf("failed to write backup: %w", err)
	}

	return backupPath, nil
}
