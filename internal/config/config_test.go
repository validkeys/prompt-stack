package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestInit(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	t.Run("creates config file with defaults", func(t *testing.T) {
		err := Init(configPath)
		if err != nil {
			t.Fatalf("Init() error = %v", err)
		}

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Fatal("config file was not created")
		}

		data, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("failed to read config file: %v", err)
		}

		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			t.Fatalf("failed to parse config: %v", err)
		}

		if cfg.Version != DefaultConfig.Version {
			t.Errorf("version = %q, want %q", cfg.Version, DefaultConfig.Version)
		}

		if cfg.Database.Path != DefaultConfig.Database.Path {
			t.Errorf("database path = %q, want %q", cfg.Database.Path, DefaultConfig.Database.Path)
		}

		if cfg.DefaultDir != DefaultConfig.DefaultDir {
			t.Errorf("default_dir = %q, want %q", cfg.DefaultDir, DefaultConfig.DefaultDir)
		}
	})

	t.Run("creates parent directories if needed", func(t *testing.T) {
		nestedPath := filepath.Join(tmpDir, "nested", "dir", "config.yaml")
		err := Init(nestedPath)
		if err != nil {
			t.Fatalf("Init() error = %v", err)
		}

		if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
			t.Fatal("config file was not created in nested directory")
		}
	})
}

func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	t.Run("loads existing config", func(t *testing.T) {
		if err := Init(configPath); err != nil {
			t.Fatalf("failed to create test config: %v", err)
		}

		cfg, err := Load(configPath)
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.Version != DefaultConfig.Version {
			t.Errorf("version = %q, want %q", cfg.Version, DefaultConfig.Version)
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		nonExistentPath := filepath.Join(tmpDir, "nonexistent.yaml")
		_, err := Load(nonExistentPath)
		if err == nil {
			t.Fatal("Load() should return error for non-existent file")
		}
	})

	t.Run("returns error for invalid YAML", func(t *testing.T) {
		invalidPath := filepath.Join(tmpDir, "invalid.yaml")
		if err := os.WriteFile(invalidPath, []byte("invalid: yaml: content:"), 0644); err != nil {
			t.Fatalf("failed to write invalid YAML: %v", err)
		}

		_, err := Load(invalidPath)
		if err == nil {
			t.Fatal("Load() should return error for invalid YAML")
		}
	})
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	t.Run("saves config to file", func(t *testing.T) {
		customConfig := Config{
			Version:    "2.0.0",
			DefaultDir: "/custom/path",
			Database: DatabaseConfig{
				Path: "/custom/db.sqlite",
			},
		}

		err := Save(configPath, &customConfig)
		if err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		cfg, err := Load(configPath)
		if err != nil {
			t.Fatalf("failed to load saved config: %v", err)
		}

		if cfg.Version != "2.0.0" {
			t.Errorf("version = %q, want 2.0.0", cfg.Version)
		}

		if cfg.DefaultDir != "/custom/path" {
			t.Errorf("default_dir = %q, want /custom/path", cfg.DefaultDir)
		}

		if cfg.Database.Path != "/custom/db.sqlite" {
			t.Errorf("database path = %q, want /custom/db.sqlite", cfg.Database.Path)
		}
	})

	t.Run("creates parent directories", func(t *testing.T) {
		nestedPath := filepath.Join(tmpDir, "nested", "dir", "config.yaml")
		customConfig := Config{Version: "1.0.0"}

		err := Save(nestedPath, &customConfig)
		if err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
			t.Fatal("config file was not created in nested directory")
		}
	})
}

func TestDefaultConfig(t *testing.T) {
	t.Run("has required fields", func(t *testing.T) {
		if DefaultConfig.Version == "" {
			t.Error("default config version should not be empty")
		}

		if DefaultConfig.DefaultDir == "" {
			t.Error("default config default_dir should not be empty")
		}

		if DefaultConfig.Database.Path == "" {
			t.Error("default config database path should not be empty")
		}
	})
}
