# Config Schema Updates for Scalability

**Document**: [`CONFIG-SCHEMA.md`](../CONFIG-SCHEMA.md)  
**Priority**: **Critical**  
**Status**: Must update before M1 (Bootstrap & Config)

---

## Overview

This document details all changes required to [`CONFIG-SCHEMA.md`](../CONFIG-SCHEMA.md) to support Phase 1 scalability abstractions. These changes add new configuration fields for AI providers, storage backends, and future features.

---

## New Configuration Fields

Add new configuration fields to support Phase 1 abstractions:

```yaml
# ~/.promptstack/config.yaml

# Existing fields (unchanged)
claude_api_key: "sk-ant-..."
model: "claude-3-sonnet-20240229"
vim_mode: false
data_dir: "~/.promptstack/data"

# NEW: AI Provider Configuration
ai_provider: "claude"  # Options: "claude", "mcp", "openai"

# NEW: Storage Configuration
storage: "sqlite"  # Options: "sqlite", "postgres", "graph"
database_path: "~/.promptstack/data/history.db"  # SQLite path
postgres_url: ""  # PostgreSQL connection string (future)
neo4j_url: ""  # Neo4j connection string (future)

# NEW: MCP Configuration (future)
mcp_host: ""  # MCP orchestrator address
specialists: []  # Enabled specialist servers

# NEW: Plugin Configuration (future)
enable_plugins: false
plugin_dir: "~/.promptstack/plugins"

# Existing field (unchanged)
version: "1.0.0"
```

---

## Field Descriptions

### ai_provider
- **Type**: string
- **Optional**: Yes
- **Default**: "claude"
- **Description**: AI provider to use for suggestions
- **Options**: "claude", "mcp", "openai"
- **Validation**: Must be one of the supported providers
- **Future**: Add more providers as needed

### storage
- **Type**: string
- **Optional**: Yes
- **Default**: "sqlite"
- **Description**: Storage backend for composition history
- **Options**: "sqlite", "postgres", "graph"
- **Validation**: Must be one of the supported storage types
- **Future**: Add more storage backends as needed

### database_path
- **Type**: string
- **Optional**: Yes
- **Default**: "~/.promptstack/data/history.db"
- **Description**: Path to SQLite database file
- **Validation**: Must be a valid file path
- **Used when**: storage == "sqlite"

### postgres_url
- **Type**: string
- **Optional**: Yes
- **Default**: ""
- **Description**: PostgreSQL connection string
- **Validation**: Must be a valid connection string format
- **Used when**: storage == "postgres"
- **Future**: For team collaboration features

### neo4j_url
- **Type**: string
- **Optional**: Yes
- **Default**: ""
- **Description**: Neo4j connection URL
- **Validation**: Must be a valid URL format
- **Used when**: storage == "graph"
- **Future**: For graph-based knowledge management

### mcp_host
- **Type**: string
- **Optional**: Yes
- **Default**: ""
- **Description**: MCP orchestrator address
- **Validation**: Must be a valid host:port format
- **Future**: For MCP-based specialist servers

### specialists
- **Type**: array of strings
- **Optional**: Yes
- **Default**: []
- **Description**: List of enabled specialist servers
- **Validation**: Each must be a valid specialist name
- **Future**: For MCP-based specialist servers

### enable_plugins
- **Type**: boolean
- **Optional**: Yes
- **Default**: false
- **Description**: Enable external plugin loading
- **Future**: For third-party specialist servers

### plugin_dir
- **Type**: string
- **Optional**: Yes
- **Default**: "~/.promptstack/plugins"
- **Description**: Directory for plugin files
- **Validation**: Must be a valid directory path
- **Used when**: enable_plugins == true
- **Future**: For third-party specialist servers

---

## Update Validation Rules

Add validation for new fields:

```go
// Validation rules for new fields
func validateConfig(cfg *Config) error {
    // Validate ai_provider
    if cfg.AIProvider != "" {
        validProviders := []string{"claude", "mcp", "openai"}
        if !contains(validProviders, cfg.AIProvider) {
            return fmt.Errorf("invalid ai_provider: %s", cfg.AIProvider)
        }
    }
    
    // Validate storage
    if cfg.Storage != "" {
        validStorage := []string{"sqlite", "postgres", "graph"}
        if !contains(validStorage, cfg.Storage) {
            return fmt.Errorf("invalid storage: %s", cfg.Storage)
        }
    }
    
    // Validate storage-specific fields
    switch cfg.Storage {
    case "sqlite":
        if cfg.DatabasePath == "" {
            return fmt.Errorf("database_path required for sqlite storage")
        }
    case "postgres":
        if cfg.PostgresURL == "" {
            return fmt.Errorf("postgres_url required for postgres storage")
        }
    case "graph":
        if cfg.Neo4jURL == "" {
            return fmt.Errorf("neo4j_url required for graph storage")
        }
    }
    
    return nil
}
```

---

## Update Setup Wizard

Add setup wizard steps for new fields:

```go
// Setup wizard additions
func runSetupWizard(cfg *Config) error {
    // Existing steps...
    
    // NEW: AI Provider Selection
    fmt.Println("\nSelect AI Provider:")
    fmt.Println("1. Claude (default)")
    fmt.Println("2. MCP (future)")
    fmt.Println("3. OpenAI (future)")
    choice := promptUser("Choice [1-3]", "1")
    
    switch choice {
    case "1":
        cfg.AIProvider = "claude"
    case "2":
        cfg.AIProvider = "mcp"
        cfg.MCPHost = promptUser("MCP Host Address", "")
    case "3":
        cfg.AIProvider = "openai"
    }
    
    // NEW: Storage Selection
    fmt.Println("\nSelect Storage Backend:")
    fmt.Println("1. SQLite (default)")
    fmt.Println("2. PostgreSQL (future)")
    fmt.Println("3. Neo4j (future)")
    choice = promptUser("Choice [1-3]", "1")
    
    switch choice {
    case "1":
        cfg.Storage = "sqlite"
    case "2":
        cfg.Storage = "postgres"
        cfg.PostgresURL = promptUser("PostgreSQL Connection String", "")
    case "3":
        cfg.Storage = "graph"
        cfg.Neo4jURL = promptUser("Neo4j Connection URL", "")
    }
    
    // Save config...
}
```

---

## Config Structure Updates

Update the Config struct in [`internal/config/config.go`](../../archive/code/internal/config/config.go):

```go
type Config struct {
    // Existing fields
    ClaudeAPIKey string `yaml:"claude_api_key"`
    Model        string `yaml:"model"`
    VimMode      bool   `yaml:"vim_mode"`
    DataDir      string `yaml:"data_dir"`
    Version      string `yaml:"version"`
    
    // NEW: AI Provider Configuration
    AIProvider string `yaml:"ai_provider"`
    
    // NEW: Storage Configuration
    Storage     string `yaml:"storage"`
    DatabasePath string `yaml:"database_path"`
    PostgresURL  string `yaml:"postgres_url"`
    Neo4jURL     string `yaml:"neo4j_url"`
    
    // NEW: MCP Configuration
    MCPHost     string   `yaml:"mcp_host"`
    Specialists []string `yaml:"specialists"`
    
    // NEW: Plugin Configuration
    EnablePlugins bool   `yaml:"enable_plugins"`
    PluginDir     string `yaml:"plugin_dir"`
}
```

---

## Migration Strategy

### Existing Users

When existing users upgrade, the config should be migrated automatically:

1. Load existing config
2. Add new fields with default values
3. Validate config
4. Save updated config

### Example Migration Code

```go
func migrateConfig(cfg *Config) error {
    // Set defaults for new fields
    if cfg.AIProvider == "" {
        cfg.AIProvider = "claude"
    }
    
    if cfg.Storage == "" {
        cfg.Storage = "sqlite"
    }
    
    if cfg.DatabasePath == "" {
        cfg.DatabasePath = filepath.Join(cfg.DataDir, "history.db")
    }
    
    if cfg.PluginDir == "" {
        cfg.PluginDir = filepath.Join(cfg.DataDir, "plugins")
    }
    
    return nil
}
```

---

## Testing Requirements

### Unit Tests

- [ ] Test validation of ai_provider field
- [ ] Test validation of storage field
- [ ] Test storage-specific field validation
- [ ] Test config migration
- [ ] Test default values

### Integration Tests

- [ ] Test config loading with new fields
- [ ] Test config saving with new fields
- [ ] Test setup wizard with new fields
- [ ] Test config validation with invalid values

---

**Last Updated**: 2026-01-07  
**Related Documents**: 
- [`scalability-implementation-summary.md`](./scalability-implementation-summary.md)
- [`scalability-project-structure-updates.md`](./scalability-project-structure-updates.md)
- [`scalability-milestones-updates.md`](./scalability-milestones-updates.md)