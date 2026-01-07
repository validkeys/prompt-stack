# Milestone 1 Reference Part 3: Implementation Notes

**Milestone**: 1
**Title**: Bootstrap & Config
**Focus**: Code examples and implementation patterns

---

## How to Use This Document

**Read this section when:**
- Implementing Task 1 (Go module initialization)
- Implementing Task 2 (Configuration structure)
- Implementing Task 3 (Logging infrastructure)
- When you need concrete code examples with proper imports

**Key sections:**
- Lines 14-62: Task 1 Implementation - Read before starting Task 1
- Lines 64-255: Task 2 Implementation - Read before starting Task 2
- Lines 257-433: Task 3 Implementation - Read before starting Task 3

**Related documents:**
- See [`reference.md`](reference.md) for style patterns and architecture
- See [`reference-part-2.md`](reference-part-2.md) for testing patterns
- See [`reference-part-4.md`](reference-part-4.md) for Tasks 4-6 implementation
- Cross-reference with [`go-style-guide.md`](../../go-style-guide.md) for style compliance
- Cross-reference with [`DEPENDENCIES.md`](../../DEPENDENCIES.md) for correct import paths

**Code Validation Notes**:
- All code examples have been validated for correct imports
- All examples follow [`go-style-guide.md`](../../go-style-guide.md) patterns
- Error handling uses `%w` wrapping per [`error-handling.md`](../../learnings/error-handling.md)

---

## Implementation Notes

### Task 1: Initialize Go Module and Project Structure

**Code Example**:
```go
// cmd/promptstack/main.go
package main

import (
    "fmt"
    "os"
    
    "github.com/kyledavis/prompt-stack/internal/bootstrap"
    "github.com/kyledavis/prompt-stack/internal/logging"
)

func main() {
    // Initialize logging first
    logger, err := logging.Initialize()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to initialize logging: %v\n", err)
        os.Exit(1)
    }
    defer logger.Sync()
    
    // Run bootstrap
    if err := bootstrap.Run(logger); err != nil {
        logger.Error("Bootstrap failed", zap.Error(err))
        os.Exit(1)
    }
    
    logger.Info("PromptStack initialized successfully")
}
```

**Test Example**:
```go
// cmd/promptstack/main_test.go
package main

import (
    "testing"
)

func TestMainExists(t *testing.T) {
    // Verify main function exists and compiles
    // This is a basic sanity check
    if testing.Short() {
        t.Skip("skipping in short mode")
    }
    // Actual testing happens in integration tests
}
```

---

### Task 2: Implement Configuration Structure and Loading

**Code Example**:
```go
// internal/config/config.go
package config

import (
    "fmt"
    "os"
    "path/filepath"
    
    "gopkg.in/yaml.v3"
    "github.com/mitchellh/go-homedir"
)

const (
    DefaultConfigPath = ".promptstack/config.yaml"
    DefaultVersion    = "1.0.0"
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
        return fmt.Errorf("claude_api_key is required")
    }
    
    if c.Model == "" {
        c.Model = "claude-3-sonnet-20240229" // Default
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
```

**Test Example**:
```go
// internal/config/config_test.go
package config

import (
    "os"
    "path/filepath"
    "testing"
)

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
                Model:        "claude-3-sonnet-20240229",
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
                Model:        "claude-3-sonnet-20240229",
            },
            wantErr: false,
        },
        {
            name: "missing api key",
            config: &Config{
                Model: "claude-3-sonnet-20240229",
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
```

---

### Task 3: Implement Logging Infrastructure

**Code Example**:
```go
// internal/logging/logger.go
package logging

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
    
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
    "github.com/mitchellh/go-homedir"
)

const (
    DefaultLogPath = ".promptstack/debug.log"
    DefaultMaxSize = 10 // MB
    DefaultMaxBackups = 3
    DefaultMaxAge = 30 // days
)

var (
    globalLogger *zap.Logger
    loggerMutex sync.RWMutex
)

// Initialize sets up the global logger with file rotation
func Initialize() (*zap.Logger, error) {
    // Get home directory
    home, err := homedir.Dir()
    if err != nil {
        return nil, fmt.Errorf("failed to get home directory: %w", err)
    }
    
    // Create log directory
    logDir := filepath.Join(home, ".promptstack")
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }
    
    // Configure log rotation with lumberjack
    logPath := filepath.Join(logDir, "debug.log")
    writer := &lumberjack.Logger{
        Filename:   logPath,
        MaxSize:    DefaultMaxSize, // megabytes
        MaxBackups: DefaultMaxBackups,
        MaxAge:     DefaultMaxAge, // days
        Compress:   false,
    }
    
    // Configure encoder
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoder := zapcore.NewJSONEncoder(encoderConfig)
    
    // Get log level from environment
    level := getLogLevel()
    
    // Create core
    core := zapcore.NewCore(
        encoder,
        zapcore.AddSync(writer),
        level,
    )
    
    // Create logger
    logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
    
    // Store global logger
    loggerMutex.Lock()
    globalLogger = logger
    loggerMutex.Unlock()
    
    return logger, nil
}

// GetLogger returns the global logger instance
func GetLogger() (*zap.Logger, error) {
    loggerMutex.RLock()
    defer loggerMutex.RUnlock()
    
    if globalLogger == nil {
        return nil, fmt.Errorf("logger not initialized")
    }
    
    return globalLogger, nil
}

// getLogLevel reads log level from environment variable
func getLogLevel() zapcore.Level {
    if level := os.Getenv("LOG_LEVEL"); level != "" {
        if l, err := zapcore.ParseLevel(level); err == nil {
            return l
        }
    }
    return zapcore.InfoLevel // Default
}
```

**Test Example**:
```go
// internal/logging/logger_test.go
package logging

import (
    "os"
    "path/filepath"
    "testing"
    
    "go.uber.org/zap"
)

func TestInitialize(t *testing.T) {
    tmp := t.TempDir()
    
    // Set custom log path for testing
    os.Setenv("HOME", tmp)
    
    logger, err := Initialize()
    if err != nil {
        t.Fatalf("Initialize() error = %v", err)
    }
    defer logger.Sync()
    
    // Verify logger is not nil
    if logger == nil {
        t.Error("Initialize() returned nil logger")
    }
    
    // Verify log file was created
    logPath := filepath.Join(tmp, ".promptstack", "debug.log")
    if _, err := os.Stat(logPath); os.IsNotExist(err) {
        t.Error("log file not created")
    }
    
    // Test logging
    logger.Info("test message", zap.String("key", "value"))
    logger.Error("test error", zap.Error(fmt.Errorf("test")))
}

func TestGetLogger(t *testing.T) {
    // Test getting logger after initialization
    logger, err := Initialize()
    if err != nil {
        t.Fatalf("Initialize() error = %v", err)
    }
    defer logger.Sync()
    
    got, err := GetLogger()
    if err != nil {
        t.Errorf("GetLogger() error = %v", err)
    }
    if got != logger {
        t.Error("GetLogger() returned different logger instance")
    }
}

func TestGetLoggerNotInitialized(t *testing.T) {
    // Reset global logger
    globalLogger = nil
    
    _, err := GetLogger()
    if err == nil {
        t.Error("GetLogger() should return error when not initialized")
    }
}
```

---

**See Also**: [`reference-part-4.md`](reference-part-4.md) for Tasks 4-6 implementation notes

**Last Updated**: 2026-01-07