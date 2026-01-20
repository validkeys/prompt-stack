package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version    string           `yaml:"version"`
	DefaultDir string           `yaml:"default_output_dir"`
	Database   DatabaseConfig   `yaml:"database"`
	Validation ValidationConfig `yaml:"validation"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type ValidationConfig struct {
	Strict bool `yaml:"strict"`
}

var DefaultConfig = Config{
	Version:    "0.1.0",
	DefaultDir: "docs/implementation-plan/m0",
	Database: DatabaseConfig{
		Path: ".prompt-stack/knowledge.db",
	},
	Validation: ValidationConfig{
		Strict: false,
	},
}

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

func Save(path string, cfg *Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func Init(path string) error {
	return Save(path, &DefaultConfig)
}
