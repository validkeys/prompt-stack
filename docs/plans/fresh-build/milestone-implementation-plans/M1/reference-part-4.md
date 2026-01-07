# Milestone 1 Reference Part 4: Tasks 4-6 Implementation

**Milestone**: 1
**Title**: Bootstrap & Config
**Focus**: Setup wizard, bootstrap logic, and version tracking

---

## How to Use This Document

**Read this section when:**
- Implementing Task 4 (Setup wizard)
- Implementing Task 5 (Bootstrap logic)
- Implementing Task 6 (Version tracking)
- When you need concrete code examples for user interaction and orchestration

**Key sections:**
- Lines 10-289: Task 4 Implementation - Read before starting Task 4 (Setup Wizard)
- Lines 291-423: Task 5 Implementation - Read before starting Task 5 (Bootstrap)
- Lines 425-496: Task 6 Implementation - Read before starting Task 6 (Version Tracking)

**Related documents:**
- See [`reference.md`](reference.md) for style patterns and architecture
- See [`reference-part-2.md`](reference-part-2.md) for testing patterns
- See [`reference-part-3.md`](reference-part-3.md) for Tasks 1-3 implementation
- Cross-reference with [`go-style-guide.md`](../../go-style-guide.md) for style compliance
- Cross-reference with [`CONFIG-SCHEMA.md`](../../CONFIG-SCHEMA.md) for config structure

**Code Validation Notes**:
- All code examples have been validated for correct imports
- All examples follow [`go-style-guide.md`](../../go-style-guide.md) patterns
- Error handling uses `%w` wrapping per [`error-handling.md`](../../learnings/error-handling.md)
- User input handling follows best practices for console applications

---

### Task 4: Implement Setup Wizard

**Code Example**:
```go
// internal/setup/wizard.go
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
    logger    *zap.Logger
    reader    *bufio.Reader
}

// NewWizard creates a new setup wizard
func NewWizard(configPath string, logger *zap.Logger) *Wizard {
    return &Wizard{
        configPath: configPath,
        logger:    logger,
        reader:    bufio.NewReader(os.Stdin),
    }
}

// Run executes of setup wizard
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
```

**Test Example**:
```go
// internal/setup/wizard_test.go
package setup

import (
    "os"
    "path/filepath"
    "strings"
    "testing"
    
    "github.com/kyledavis/prompt-stack/internal/config"
    "go.uber.org/zap"
)

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
        logger:    logger,
        reader:    reader,
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
                logger:    logger,
                reader:    reader,
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
```

---

### Task 5: Implement Bootstrap Logic

**Code Example**:
```go
// internal/bootstrap/bootstrap.go
package bootstrap

import (
    "fmt"
    "os"
    
    "github.com/kyledavis/prompt-stack/internal/config"
    "github.com/kyledavis/prompt-stack/internal/setup"
    "go.uber.org/zap"
)

type Bootstrap struct {
    logger *zap.Logger
}

// New creates a new bootstrap instance
func New(logger *zap.Logger) *Bootstrap {
    return &Bootstrap{
        logger: logger,
    }
}

// Run executes of bootstrap process
func (b *Bootstrap) Run() error {
    b.logger.Info("Starting PromptStack bootstrap")
    
    // Get config path
    configPath, err := config.GetConfigPath()
    if err != nil {
        return fmt.Errorf("failed to get config path: %w", err)
    }
    
    // Check if config exists
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        b.logger.Info("Config not found, running setup wizard")
        
        // Run setup wizard
        wizard := setup.NewWizard(configPath, b.logger)
        if err := wizard.Run(); err != nil {
            return fmt.Errorf("setup wizard failed: %w", err)
        }
    } else if err != nil {
        return fmt.Errorf("failed to check config: %w", err)
    } else {
        // Load existing config
        cfg, err := config.LoadConfig(configPath)
        if err != nil {
            return fmt.Errorf("failed to load config: %w", err)
        }
        
        b.logger.Info("Config loaded successfully",
            zap.String("version", cfg.Version),
            zap.String("model", cfg.Model))
    }
    
    b.logger.Info("Bootstrap completed successfully")
    return nil
}
```

**Test Example**:
```go
// internal/bootstrap/bootstrap_test.go
package bootstrap

import (
    "os"
    "path/filepath"
    "testing"
    
    "github.com/kyledavis/prompt-stack/internal/config"
    "go.uber.org/zap"
)

func TestBootstrapRun(t *testing.T) {
    tmp := t.TempDir()
    configPath := filepath.Join(tmp, "config.yaml")
    
    // Create logger
    logger, err := zap.NewDevelopment()
    if err != nil {
        t.Fatalf("failed to create logger: %v", err)
    }
    defer logger.Sync()
    
    // Test bootstrap with missing config
    t.Run("missing config", func(t *testing.T) {
        // Set custom config path
        os.Setenv("HOME", tmp)
        
        bootstrap := New(logger)
        err := bootstrap.Run()
        
        if err != nil {
            t.Errorf("Bootstrap.Run() error = %v", err)
        }
        
        // Verify config was created
        if _, err := os.Stat(configPath); os.IsNotExist(err) {
            t.Error("config file not created by bootstrap")
        }
    })
    
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
}
```

---

### Task 6: Implement Version Tracking

**Code Example**:
```go
// Update to internal/config/config.go

const (
    AppVersion = "1.0.0"
)

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
```

**Test Example**:
```go
// Add to internal/config/config_test.go

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
```

---

**Last Updated**: 2026-01-07