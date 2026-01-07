// Package config provides application configuration management including
// loading, validation, and persistence of user settings.
package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() error = %v", err)
	}

	if path == "" {
		t.Error("GetConfigPath() returned empty path")
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    *Config
		wantErr bool
	}{
		{
			name: "valid config",
			content: `version: "1.0.0"
claude_api_key: "sk-ant-test"
model: "claude-3-sonnet-20240229"
vim_mode: false`,
			want: &Config{
				Version:      "1.0.0",
				ClaudeAPIKey: "sk-ant-test",
				Model:        ModelSonnet,
				VimMode:      false,
			},
			wantErr: false,
		},
		{
			name:    "missing api key",
			content: `version: "1.0.0"`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid yaml",
			content: `invalid: yaml: content: [`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			path := filepath.Join(tmp, "config.yaml")

			if err := os.WriteFile(path, []byte(tt.content), 0600); err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			got, err := LoadConfig(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				ClaudeAPIKey: "sk-ant-test",
				Model:        ModelSonnet,
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			config: &Config{
				Model: ModelSonnet,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigCheckVersion(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "matching version",
			config: &Config{
				Version: "1.0.0",
			},
			wantErr: false,
		},
		{
			name: "empty version",
			config: &Config{
				Version: "",
			},
			wantErr: false,
		},
		{
			name: "mismatched version",
			config: &Config{
				Version: "0.9.0",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.CheckVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.CheckVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveConfig(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.yaml")

	cfg := &Config{
		Version:      "1.0.0",
		ClaudeAPIKey: "sk-ant-test",
		Model:        ModelSonnet,
		VimMode:      false,
	}

	err := cfg.SaveConfig(path)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("config file not created")
	}

	// Load and verify content
	loaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}

	if loaded.Version != cfg.Version {
		t.Errorf("Version = %s, want %s", loaded.Version, cfg.Version)
	}

	if loaded.ClaudeAPIKey != cfg.ClaudeAPIKey {
		t.Errorf("ClaudeAPIKey = %s, want %s", loaded.ClaudeAPIKey, cfg.ClaudeAPIKey)
	}
}

func TestConfigValidateDefaults(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   string
	}{
		{
			name: "default model applied",
			config: &Config{
				ClaudeAPIKey: "sk-ant-test",
				Model:        "",
			},
			want: ModelSonnet,
		},
		{
			name: "model preserved",
			config: &Config{
				ClaudeAPIKey: "sk-ant-test",
				Model:        ModelOpus,
			},
			want: ModelOpus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err != nil {
				t.Errorf("Config.Validate() error = %v", err)
			}
			if tt.config.Model != tt.want {
				t.Errorf("Model = %s, want %s", tt.config.Model, tt.want)
			}
		})
	}
}
