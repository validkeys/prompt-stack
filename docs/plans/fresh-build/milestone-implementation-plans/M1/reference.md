# Milestone 1 Reference: Bootstrap & Config

**Milestone**: 1
**Title**: Bootstrap & Config
**Focus**: Foundation infrastructure initialization

---

## How to Use This Document

**Read this section when:**
- Before implementing any Task 1-6 to understand overall architecture
- When you need to understand package structure and dependencies
- To review style guide patterns before writing code
- When troubleshooting integration between config, logging, and bootstrap

**Key sections:**
- Lines 24-48: Package Structure - Read before starting any task
- Lines 50-68: Package Naming - Reference when creating new packages
- Lines 70-98: Error Messages - Reference when writing error handling
- Lines 100-122: Exported vs Unexported - Reference when defining functions
- Lines 124-140: Method Receivers - Reference when adding methods to structs
- Lines 142-168: Dependency Injection - Reference when structuring components

**Related documents:**
- See [`reference-part-2.md`](reference-part-2.md) for testing patterns and key learnings
- See [`reference-part-3.md`](reference-part-3.md) for Tasks 1-3 implementation examples
- See [`reference-part-4.md`](reference-part-4.md) for Tasks 4-6 implementation examples
- Cross-reference with [`go-style-guide.md`](../../go-style-guide.md) for complete style rules
- Cross-reference with [`CONFIG-SCHEMA.md`](../../CONFIG-SCHEMA.md) for config structure

---

## Architecture Context

### Domain Overview

Milestone 1 establishes foundational infrastructure for PromptStack. This includes:

1. **Configuration Management** - Load, validate, and manage user configuration
2. **Setup Wizard** - Interactive first-run configuration
3. **Logging Infrastructure** - Structured logging with rotation
4. **Bootstrap Logic** - Application initialization orchestration
5. **Version Tracking** - Configuration version management

### Package Structure

```
cmd/promptstack/
└── main.go              # Application entry point

internal/
├── config/
│   └── config.go         # Configuration structure and loading
├── setup/
│   └── wizard.go         # Setup wizard implementation
├── logging/
│   └── logger.go         # Logging initialization
└── bootstrap/
    └── bootstrap.go      # Bootstrap orchestration
```

### Dependencies

**External Dependencies** (from [`DEPENDENCIES.md`](../../DEPENDENCIES.md)):
- `gopkg.in/yaml.v3` - YAML parsing
- `go.uber.org/zap` - Structured logging
- `gopkg.in/natefinch/lumberjack.v2` - Log rotation
- `github.com/mitchellh/go-homedir` - Home directory detection

**Internal Dependencies**:
- None (first milestone)

---

## Style Guide References

### Package Naming

From [`go-style-guide.md`](../../go-style-guide.md):

```go
// ✅ GOOD: Singular, lowercase
package config
package setup
package logging
package bootstrap

// ❌ BAD: Plural or uppercase
package configs
package Setup
package Logging
```

### Error Messages

From [`go-style-guide.md`](../../go-style-guide.md):

```go
// ✅ GOOD: Lowercase, no punctuation, %w for wrapping
return fmt.Errorf("failed to load config: %w", err)

// ❌ BAD: Uppercase, punctuation, no wrapping
return fmt.Errorf("Failed to load config: %v", err)
```

### Exported vs Unexported

From [`go-style-guide.md`](../../go-style-guide.md):

```go
// ✅ GOOD: Exported functions have comments
// LoadConfig loads configuration from specified path.
func LoadConfig(path string) (*Config, error) {
    // ...
}

// Unexported helper functions
func validateConfig(cfg *Config) error {
    // ...
}

// ❌ BAD: Missing comments on exported functions
func LoadConfig(path string) (*Config, error) {
    // ...
}
```

### Method Receivers

From [`go-style-guide.md`](../../go-style-guide.md):

```go
// ✅ GOOD: Consistent pointer receivers for mutable types
func (c *Config) Validate() error {
    // ...
}

func (c *Config) Save(path string) error {
    // ...
}

// ✅ GOOD: Value receivers for immutable types
func (c Config) Version() string {
    return c.Version
}
```

### No Global State

From [`go-style-guide.md`](../../go-style-guide.md):

```go
// ✅ GOOD: Dependency injection
func NewBootstrap(cfg *Config, logger *zap.Logger) *Bootstrap {
    return &Bootstrap{
        config: cfg,
        logger: logger,
    }
}

// ❌ BAD: Global variables
var globalConfig *Config
var globalLogger *zap.Logger
```

### Dependency Injection

From [`go-style-guide.md`](../../go-style-guide.md):

```go
// ✅ GOOD: Inject dependencies
type Bootstrap struct {
    config *Config
    logger *zap.Logger
}

func NewBootstrap(cfg *Config, logger *zap.Logger) *Bootstrap {
    return &Bootstrap{
        config: cfg,
        logger: logger,
    }
}

// ❌ BAD: Direct dependencies
type Bootstrap struct {}

func (b *Bootstrap) Run() error {
    cfg := LoadConfig() // Direct dependency
    logger := GetLogger() // Direct dependency
    // ...
}
```

---

## Integration Considerations

### Internal Integration
- **Config → Logging**: Config provides log level and file path
- **Setup → Config**: Setup wizard creates and validates config
- **Bootstrap → All**: Bootstrap orchestrates initialization order
- **Version → Config**: Version stored in config structure

### External Integration
- None (first milestone)

---

## Common Pitfalls to Avoid

### From [`go-fundamentals.md`](../../learnings/go-fundamentals.md)

1. **Don't ignore errors**: Always handle errors, don't use `_`
2. **Don't use string literals for zap fields**: Use `zap.String()`, `zap.Int()`, etc.
3. **Don't discard original errors**: Always wrap with `%w`
4. **Don't use global state**: Use dependency injection

### From [`architecture-patterns.md`](../../learnings/architecture-patterns.md)

1. **Don't skip validation**: Always validate configuration after loading
2. **Don't overwrite user data**: Check for existing files before writing
3. **Don't use unstructured logging**: Use zap with structured fields

### From [`error-handling.md`](../../learnings/error-handling.md)

1. **Don't create import cycles**: Keep message types in lower-level packages
2. **Don't lose error context**: Always wrap errors with context
3. **Don't ignore error severity**: Use severity levels for display strategy

---

**See Also**: [`reference-part-2.md`](reference-part-2.md) for testing guide references and key learnings

**Last Updated**: 2026-01-07