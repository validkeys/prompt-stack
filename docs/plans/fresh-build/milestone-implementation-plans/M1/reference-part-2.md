# Milestone 1 Reference Part 2: Testing & Key Learnings

**Milestone**: 1
**Title**: Bootstrap & Config
**Focus**: Testing patterns and key learnings application

---

## How to Use This Document

**Read this section when:**
- Before writing tests for any Task 1-6
- When you need concrete test examples for config, logging, or setup
- To understand acceptance criteria verification approach
- When applying key learnings from previous implementations

**Key sections:**
- Lines 25-60: Table-Driven Tests - Read before writing any test
- Lines 62-95: Test Helpers - Reference when creating test utilities
- Lines 97-118: Mock File System - Reference for file operation tests
- Lines 120-153: Acceptance Criteria Verification - Read before Task completion
- Lines 155-230: Error Handling Patterns - Reference throughout implementation
- Lines 232-290: Logging Strategy - Read before Task 3 (logging)
- Lines 292-386: Error Handling Architecture - Reference for error design

**Related documents:**
- See [`reference.md`](reference.md) for architecture and style patterns
- See [`reference-part-3.md`](reference-part-3.md) for Tasks 1-3 code examples
- See [`reference-part-4.md`](reference-part-4.md) for Tasks 4-6 code examples
- Cross-reference with [`go-testing-guide.md`](../../go-testing-guide.md) for complete testing patterns
- Cross-reference with [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md) for milestone-specific testing

---

## Testing Guide References

### Testing Guide for Foundation Milestones

From [`FOUNDATION-TESTING-GUIDE.md`](../../milestones/FOUNDATION-TESTING-GUIDE.md):

**Key Testing Patterns**:
- Table-driven tests for multiple scenarios
- Mock file system for file operations
- Test error paths and edge cases
- Integration tests for bootstrap flow

### Test Patterns from [`go-testing-guide.md`](../../go-testing-guide.md)

#### Table-Driven Tests

```go
func TestLoadConfig(t *testing.T) {
    tests := []struct {
        name    string
        content string
        want    *Config
        wantErr bool
    }{
        {
            name:    "valid config",
            content: "version: \"1.0.0\"\napi:\n  claude_api_key: \"sk-ant-test\"",
            want:    &Config{Version: "1.0.0", ClaudeAPIKey: "sk-ant-test"},
            wantErr: false,
        },
        {
            name:    "missing api key",
            content: "version: \"1.0.0\"",
            want:    nil,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := LoadConfig(tt.content)
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
```

#### Test Helpers

```go
// testutil/helpers.go
package testutil

import (
    "testing"
    "os"
    "path/filepath"
)

// TempConfig creates a temporary config file
func TempConfig(t *testing.T, content string) string {
    t.Helper()
    
    tmp := t.TempDir()
    path := filepath.Join(tmp, "config.yaml")
    
    if err := os.WriteFile(path, []byte(content), 0600); err != nil {
        t.Fatalf("failed to create temp config: %v", err)
    }
    
    return path
}

// AssertEqual marks function as test helper
func AssertEqual(t *testing.T, got, want interface{}) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

#### Mock File System

```go
func TestSetupWizard(t *testing.T) {
    // Use t.TempDir() for temporary files
    tmp := t.TempDir()
    configPath := filepath.Join(tmp, "config.yaml")
    
    // Test wizard with temp directory
    wizard := NewWizard(configPath)
    err := wizard.Run()
    
    if err != nil {
        t.Errorf("wizard.Run() error = %v", err)
    }
    
    // Verify config was created
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        t.Error("config file not created")
    }
}
```

### Acceptance Criteria Verification

From [`ENHANCED-TEST-CRITERIA-TEMPLATE.md`](../../milestones/ENHANCED-TEST-CRITERIA-TEMPLATE.md):

**Functional Requirements**:
- App launches without errors on fresh install
- Config file created at `~/.promptstack/config.yaml` with correct structure
- Setup wizard prompts for API key and preferences
- Logs written to `~/.promptstack/debug.log`
- Version stored in config and compared on startup

**Integration Requirements**:
- Config loading integrates with logging system
- Setup wizard integrates with config persistence
- Version tracking integrates with starter prompt extraction

**Edge Cases & Error Handling**:
- Handle missing config directory (create automatically)
- Handle corrupted config file (show error, offer reset)
- Handle invalid API key format (reject with clear message)
- Handle read-only filesystem (show error, exit gracefully)
- Handle interrupted setup wizard (resume on next launch)

**Performance Requirements**:
- App startup time <500ms on fresh install
- Config file read/write <50ms
- Log file write <10ms per entry

**User Experience Requirements**:
- Setup wizard provides clear instructions
- Error messages are actionable and specific
- Progress indicators shown during initialization
- Keyboard navigation works in setup wizard

---

## Key Learnings References

### From [`go-fundamentals.md`](../../learnings/go-fundamentals.md)

#### Error Handling Patterns

**Learning**: Use `fmt.Errorf` with `%w` for error wrapping

```go
// ✅ GOOD: Wrap errors with context
if err != nil {
    return nil, fmt.Errorf("failed to load config: %w", err)
}

// Check for specific errors
if errors.Is(err, os.ErrNotExist) {
    // Handle file not found
}
```

**When to Apply**: When handling errors throughout application

#### Project Structure Organization

**Learning**: Organize internal packages by domain/feature

```
internal/
├── config/      # Configuration management
├── setup/        # First-run setup
├── bootstrap/    # Application initialization
└── logging/      # Logging setup
```

**When to Apply**: When structuring new Go projects

### From [`architecture-patterns.md`](../../learnings/architecture-patterns.md)

#### Configuration Management

**Learning**: Always validate configuration after loading

```go
type Config struct {
    Version      string `yaml:"version"`
    ClaudeAPIKey string `yaml:"claude_api_key"`
    Model        string `yaml:"model"`
    VimMode      bool   `yaml:"vim_mode"`
}

func (c *Config) Validate() error {
    if c.ClaudeAPIKey == "" {
        return fmt.Errorf("claude_api_key is required")
    }
    return nil
}

func LoadConfig(configPath string) (*Config, error) {
    data, err := os.ReadFile(configPath)
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
```

**When to Apply**: When implementing configuration management

#### Logging Strategy

**Learning**: Use structured logging from start with rotation

```go
func Initialize() (*zap.Logger, error) {
    // Create log directory
    logDir, err := config.GetLogPath()
    if err != nil {
        return nil, fmt.Errorf("failed to get log path: %w", err)
    }

    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }

    // Configure rotation
    logPath := filepath.Join(logDir, "debug.log")
    writer, err := rotatelogs.New(
        logPath,
        rotatelogs.WithMaxAge(30*24*time.Hour), // 30 days
        rotatelogs.WithRotationTime(24*time.Hour), // Daily
        rotatelogs.WithMaxBackups(3), // Keep last 3
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create log writer: %w", err)
    }

    // Configure encoder
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoder := zapcore.NewJSONEncoder(encoderConfig)

    // Create core
    core := zapcore.NewCore(
        encoder,
        zapcore.AddSync(writer),
        getLogLevel(),
    )

    // Create logger
    logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

    return logger, nil
}

func getLogLevel() zapcore.Level {
    // Read from environment variable
    if level := os.Getenv("LOG_LEVEL"); level != "" {
        if l, err := zapcore.ParseLevel(level); err == nil {
            return l
        }
    }
    return zapcore.InfoLevel // Default
}
```

**When to Apply**: When implementing logging in Go applications

### From [`error-handling.md`](../../learnings/error-handling.md)

#### Error Handling Architecture

**Learning**: Create structured error types with severity levels

```go
type ErrorType string
const (
    ErrorTypeConfig ErrorType = "config"
    ErrorTypeFile   ErrorType = "file"
)

type Severity string
const (
    SeverityError   Severity = "error"
    SeverityWarning Severity = "warning"
)

type AppError struct {
    Type      ErrorType
    Severity  Severity
    Message   string
    Details   string
    Timestamp time.Time
    Err       error
}

func New(errorType ErrorType, message string) *AppError {
    return &AppError{
        Type:      errorType,
        Severity:  SeverityError,
        Message:   message,
        Timestamp: time.Now(),
    }
}

func (e *AppError) WithDetails(details string) *AppError {
    e.Details = details
    return e
}

func (e *AppError) WithError(err error) *AppError {
    e.Err = err
    return e
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}
```

**When to Apply**: When implementing error handling across application

#### Error Logging Integration

**Learning**: Integrate error logging throughout application using global logger pattern

```go
var (
    globalLogger *zap.Logger
    loggerMutex sync.RWMutex
)

func Initialize() (*zap.Logger, error) {
    // ... logger setup ...
    logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
    
    // Store global logger instance
    loggerMutex.Lock()
    globalLogger = logger
    loggerMutex.Unlock()
    
    return logger, nil
}

func GetLogger() (*zap.Logger, error) {
    loggerMutex.RLock()
    defer loggerMutex.RUnlock()
    
    if globalLogger == nil {
        return nil, fmt.Errorf("logger not initialized")
    }
    
    return globalLogger, nil
}
```

**When to Apply**: When implementing error handling and logging across application

---

**See Also**: [`reference-part-3.md`](reference-part-3.md) for implementation notes

**Last Updated**: 2026-01-07