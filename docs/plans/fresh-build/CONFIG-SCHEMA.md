# Configuration Schema

**Version**: 1.0.0  
**Last Updated**: 2026-01-07  
**Status**: Ready for Implementation

---

## Overview

This document defines the complete configuration schema for PromptStack, including all available options, default values, validation rules, and environment variable overrides.

## Config File Location

**Default Path**: `~/.promptstack/config.yaml`

The configuration file is created automatically on first launch through an interactive setup wizard.

---

## Complete Schema

```yaml
# ============================================
# PromptStack Configuration
# ============================================

# Version tracking (auto-managed, do not edit manually)
version: "1.0.0"

# ============================================
# API Configuration
# ============================================
api:
  # Claude API key (required)
  # Can be set via environment variable: PROMPTSTACK_CLAUDE_API_KEY
  claude_api_key: ""
  
  # Claude model selection (required)
  # Valid values: claude-3-sonnet-20240229, claude-3-opus-20240229, claude-3-haiku-20240307
  # Can be set via environment variable: PROMPTSTACK_MODEL
  model: "claude-3-sonnet-20240229"
  
  # Maximum retry attempts for API requests (default: 3)
  max_retries: 3
  
  # Request timeout duration (default: 30s)
  # Format: Go duration string (e.g., 30s, 1m, 500ms)
  timeout: "30s"
  
  # AI provider selection (default: claude)
  # Valid values: claude, mcp, openai
  # Can be set via environment variable: PROMPTSTACK_AI_PROVIDER
  provider: "claude"
  
  # Context window management
  context:
    # Maximum percentage of context window for composition content (default: 25)
    # Range: 10-50
    composition_max_percent: 25
    
    # Maximum percentage of context window for library prompts (default: 15)
    # Range: 5-30
    library_max_percent: 15
    
    # Warning threshold for composition size (default: 15)
    # Show warning when composition exceeds this percentage
    warning_threshold: 15

# ============================================
# Editor Configuration
# ============================================
editor:
  # Enable vim keybindings (default: false)
  # Requires application restart to take effect
  vim_mode: false
  
  # Auto-save interval (default: 1s)
  # Format: Go duration string (e.g., 500ms, 1s, 2s)
  auto_save_interval: "1s"
  
  # Undo history size (default: 100)
  # Maximum number of undo actions to keep in memory
  undo_stack_size: 100
  
  # Smart batching timeout (default: 1s)
  # Time to wait before batching continuous typing as one undo action
  batch_timeout: "1s"

# ============================================
# Paths Configuration
# ============================================
paths:
  # Data directory (default: ~/.promptstack/data)
  # Contains library, history, and index files
  data_dir: "~/.promptstack/data"
  
  # Library directory (default: ~/.promptstack/data/library)
  # Contains prompt templates organized by category
  library_dir: "~/.promptstack/data/library"
  
  # History directory (default: ~/.promptstack/data/.history)
  # Contains composition history markdown files
  history_dir: "~/.promptstack/data/.history"
  
  # Database file (default: ~/.promptstack/data/history.db)
  # SQLite database for history indexing and search
  database_file: "~/.promptstack/data/history.db"
  
  # Index file (default: ~/.promptstack/data/.index.json)
  # Library index for AI context optimization
  index_file: "~/.promptstack/data/.index.json"
  
  # Plugin directory (default: ~/.promptstack/plugins)
  # Contains third-party plugin files
  plugin_dir: "~/.promptstack/plugins"

# ============================================
# Logging Configuration
# ============================================
logging:
  # Log level (default: info)
  # Valid values: debug, info, warn, error
  # Can be set via environment variable: PROMPTSTACK_LOG_LEVEL
  level: "info"
  
  # Log file path (default: ~/.promptstack/debug.log)
  file: "~/.promptstack/debug.log"
  
  # Maximum log file size before rotation (default: 10MB)
  # Format: size string (e.g., 10MB, 5MB)
  max_size: "10MB"
  
  # Number of rotated log files to keep (default: 3)
  max_backups: 3
  
  # Log format (default: json)
  # Valid values: json, text
  format: "json"

# ============================================
# UI Configuration
# ============================================
ui:
  # Theme name (default: catppuccin-mocha)
  # Valid values: catppuccin-mocha, dracula, nord, monokai
  theme: "catppuccin-mocha"
  
  # AI panel initial width as percentage of screen (default: 40)
  # Range: 20-60
  ai_panel_width_percent: 40
  
  # Minimum terminal width for split-pane mode (default: 100)
  # Below this width, AI panel becomes overlay modal
  min_split_width: 100
  
  # Status bar configuration
  status_bar:
    # Show character count (default: true)
    show_char_count: true
    
    # Show line count (default: true)
    show_line_count: true
    
    # Show token estimate (default: true)
    show_token_estimate: true
    
    # Show vim mode indicator (default: true)
    show_vim_mode: true
    
    # Auto-save notification duration (default: 2s)
    # Format: Go duration string
    notification_duration: "2s"

# ============================================
# Library Configuration
# ============================================
library:
  # Maximum file size for prompts (default: 1MB)
  # Format: size string (e.g., 1MB, 500KB)
  max_file_size: "1MB"
  
  # Auto-validate library on startup (default: true)
  validate_on_startup: true
  
  # Show validation warnings in library browser (default: true)
  show_warnings: true

# ============================================
# History Configuration
# ============================================
history:
  # Maximum history files to keep (default: unlimited)
  # Set to 0 for unlimited
  max_files: 0
  
  # Maximum age of history files (default: unlimited)
  # Format: Go duration string (e.g., 90d, 180d)
  # Set to 0 for unlimited
  max_age: "0"
  
  # Auto-rebuild index on startup if out of sync (default: false)
  # If false, prompts user to confirm rebuild
  auto_rebuild_index: false

# ============================================
# AI Suggestions Configuration
# ============================================
ai:
  # Enable AI suggestions (default: true)
  enabled: true
  
  # Maximum number of library prompts to include in context (default: 5)
  # Range: 1-10
  max_library_prompts: 5
  
  # Suggestion types to enable (default: all)
  # Valid types: recommendations, gap_analysis, formatting, contradictions, clarity, reformatting
  enabled_types:
    - recommendations
    - gap_analysis
    - formatting
    - contradictions
    - clarity
    - reformatting
  
  # Keyword extraction configuration
  keywords:
    # Minimum word length to consider (default: 3)
    min_length: 3
    
    # Maximum number of keywords to extract (default: 20)
    max_count: 20
    
    # Enable TF-IDF scoring (default: true)
    use_tfidf: true

# ============================================
# File Reference Configuration
# ============================================
files:
  # Maximum file size for references (default: 1MB)
  # Format: size string (e.g., 1MB, 500KB)
  max_file_size: "1MB"
  
  # Respect .gitignore (default: true)
  respect_gitignore: true
  
  # Default encoding for file reading (default: utf-8)
  encoding: "utf-8"
  
  # Preserve original line endings (default: true)
  preserve_line_endings: true

# ============================================
# Advanced Configuration
# ============================================
advanced:
  # Enable debug mode (default: false)
  # Can be set via environment variable: PROMPTSTACK_DEBUG
  debug: false
  
  # Enable verbose logging (default: false)
  verbose: false
  
  # Performance profiling (default: false)
  profile: false
  
  # Cache library in memory (default: true)
  cache_library: true

# ============================================
# Storage Configuration
# ============================================
storage:
  # Storage backend type (default: sqlite)
  # Valid values: sqlite, postgres, graph
  # Can be set via environment variable: PROMPTSTACK_STORAGE
  type: "sqlite"
  
  # SQLite database path (default: ~/.promptstack/data/history.db)
  # Used when storage.type == "sqlite"
  database_path: "~/.promptstack/data/history.db"
  
  # PostgreSQL connection string (default: "")
  # Used when storage.type == "postgres"
  # Format: postgres://user:password@host:port/dbname
  # Can be set via environment variable: PROMPTSTACK_POSTGRES_URL
  postgres_url: ""
  
  # Neo4j connection URL (default: "")
  # Used when storage.type == "graph"
  # Format: bolt://user:password@host:7687
  # Can be set via environment variable: PROMPTSTACK_NEO4J_URL
  neo4j_url: ""

# ============================================
# MCP Configuration (Future)
# ============================================
mcp:
  # MCP orchestrator address (default: "")
  # Format: host:port
  # Can be set via environment variable: PROMPTSTACK_MCP_HOST
  host: ""
  
  # Enabled specialist servers (default: [])
  # List of specialist server names to enable
  specialists: []

# ============================================
# Plugin Configuration (Future)
# ============================================
plugins:
  # Enable external plugin loading (default: false)
  # Can be set via environment variable: PROMPTSTACK_ENABLE_PLUGINS
  enabled: false
```

---

## Validation Rules

### Required Fields

The following fields must be non-empty and valid:

- [`api.claude_api_key`](#api-configuration) - Must be a valid Claude API key
- [`api.model`](#api-configuration) - Must be a valid Claude model name

### Field Validation

| Field | Type | Validation | Default |
|-------|------|------------|---------|
| `version` | string | Semantic version format | Auto-managed |
| `api.claude_api_key` | string | Non-empty string | Required |
| `api.model` | string | Valid model name | claude-3-sonnet-20240229 |
| `api.provider` | string | claude\|mcp\|openai | claude |
| `api.max_retries` | int | Range: 0-10 | 3 |
| `api.timeout` | duration | Valid Go duration | 30s |
| `api.context.composition_max_percent` | int | Range: 10-50 | 25 |
| `api.context.library_max_percent` | int | Range: 5-30 | 15 |
| `api.context.warning_threshold` | int | Range: 5-40 | 15 |
| `editor.vim_mode` | bool | true or false | false |
| `editor.auto_save_interval` | duration | Valid Go duration, min: 100ms | 1s |
| `editor.undo_stack_size` | int | Range: 10-1000 | 100 |
| `editor.batch_timeout` | duration | Valid Go duration, min: 100ms | 1s |
| `paths.data_dir` | string | Valid directory path | ~/.promptstack/data |
| `paths.library_dir` | string | Valid directory path | ~/.promptstack/data/library |
| `paths.history_dir` | string | Valid directory path | ~/.promptstack/data/.history |
| `paths.database_file` | string | Valid file path | ~/.promptstack/data/history.db |
| `paths.index_file` | string | Valid file path | ~/.promptstack/data/.index.json |
| `paths.plugin_dir` | string | Valid directory path | ~/.promptstack/plugins |
| `logging.level` | string | debug\|info\|warn\|error | info |
| `logging.file` | string | Valid file path | ~/.promptstack/debug.log |
| `logging.max_size` | size | Valid size string | 10MB |
| `logging.max_backups` | int | Range: 0-10 | 3 |
| `logging.format` | string | json\|text | json |
| `ui.theme` | string | Valid theme name | catppuccin-mocha |
| `ui.ai_panel_width_percent` | int | Range: 20-60 | 40 |
| `ui.min_split_width` | int | Range: 80-120 | 100 |
| `ui.status_bar.show_char_count` | bool | true or false | true |
| `ui.status_bar.show_line_count` | bool | true or false | true |
| `ui.status_bar.show_token_estimate` | bool | true or false | true |
| `ui.status_bar.show_vim_mode` | bool | true or false | true |
| `ui.status_bar.notification_duration` | duration | Valid Go duration | 2s |
| `library.max_file_size` | size | Valid size string, max: 10MB | 1MB |
| `library.validate_on_startup` | bool | true or false | true |
| `library.show_warnings` | bool | true or false | true |
| `history.max_files` | int | Range: 0 or 10-10000 | 0 (unlimited) |
| `history.max_age` | duration | Valid Go duration, 0 for unlimited | 0 (unlimited) |
| `history.auto_rebuild_index` | bool | true or false | false |
| `ai.enabled` | bool | true or false | true |
| `ai.max_library_prompts` | int | Range: 1-10 | 5 |
| `ai.enabled_types` | list | Valid suggestion types | All types |
| `ai.keywords.min_length` | int | Range: 2-10 | 3 |
| `ai.keywords.max_count` | int | Range: 5-50 | 20 |
| `ai.keywords.use_tfidf` | bool | true or false | true |
| `files.max_file_size` | size | Valid size string, max: 10MB | 1MB |
| `files.respect_gitignore` | bool | true or false | true |
| `files.encoding` | string | Valid encoding name | utf-8 |
| `files.preserve_line_endings` | bool | true or false | true |
| `advanced.debug` | bool | true or false | false |
| `advanced.verbose` | bool | true or false | false |
| `advanced.profile` | bool | true or false | false |
| `advanced.cache_library` | bool | true or false | true |
| `storage.type` | string | sqlite\|postgres\|graph | sqlite |
| `storage.database_path` | string | Valid file path | ~/.promptstack/data/history.db |
| `storage.postgres_url` | string | Valid connection string | "" |
| `storage.neo4j_url` | string | Valid URL format | "" |
| `mcp.host` | string | Valid host:port format | "" |
| `mcp.specialists` | list | Valid specialist names | [] |
| `plugins.enabled` | bool | true or false | false |

### Validation Error Messages

- **Missing required field**: "Configuration error: Field '{field}' is required but not set"
- **Invalid value type**: "Configuration error: Field '{field}' must be of type {type}"
- **Value out of range**: "Configuration error: Field '{field}' value {value} is out of range [{min}, {max}]"
- **Invalid enum value**: "Configuration error: Field '{field}' value '{value}' is not valid. Valid values: {valid_values}"
- **Invalid duration**: "Configuration error: Field '{field}' value '{value}' is not a valid duration format"
- **Invalid size**: "Configuration error: Field '{field}' value '{value}' is not a valid size format"
- **Invalid path**: "Configuration error: Field '{field}' path '{value}' is not a valid path"

---

## Environment Variable Overrides

Configuration values can be overridden using environment variables. Environment variables take precedence over values in the config file.

### Supported Environment Variables

| Environment Variable | Config Field | Type | Notes |
|---------------------|--------------|------|-------|
| `PROMPTSTACK_CLAUDE_API_KEY` | `api.claude_api_key` | string | Required for AI features |
| `PROMPTSTACK_MODEL` | `api.model` | string | Claude model selection |
| `PROMPTSTACK_AI_PROVIDER` | `api.provider` | string | AI provider selection |
| `PROMPTSTACK_LOG_LEVEL` | `logging.level` | string | debug\|info\|warn\|error |
| `PROMPTSTACK_DEBUG` | `advanced.debug` | bool | Set to "1" or "true" to enable |
| `PROMPTSTACK_VERBOSE` | `advanced.verbose` | bool | Set to "1" or "true" to enable |
| `PROMPTSTACK_STORAGE` | `storage.type` | string | Storage backend selection |
| `PROMPTSTACK_POSTGRES_URL` | `storage.postgres_url` | string | PostgreSQL connection string |
| `PROMPTSTACK_NEO4J_URL` | `storage.neo4j_url` | string | Neo4j connection URL |
| `PROMPTSTACK_MCP_HOST` | `mcp.host` | string | MCP orchestrator address |
| `PROMPTSTACK_ENABLE_PLUGINS` | `plugins.enabled` | bool | Set to "1" or "true" to enable |

### Usage Examples

```bash
# Set API key via environment variable
export PROMPTSTACK_CLAUDE_API_KEY="sk-ant-api03-..."

# Set model via environment variable
export PROMPTSTACK_MODEL="claude-3-opus-20240229"

# Enable debug mode
export PROMPTSTACK_DEBUG=1

# Set log level
export PROMPTSTACK_LOG_LEVEL=debug

# Run promptstack with environment overrides
promptstack
```

### Environment Variable Precedence

1. Environment variables (highest priority)
2. Config file values
3. Default values (lowest priority)

---

## Configuration Migration

### Version Tracking

The `version` field in the config file tracks the last application version that wrote the configuration. This is used for:

- Detecting application upgrades
- Migrating configuration between versions
- Managing starter prompt updates

### Migration Process

When the application detects a version mismatch:

1. **Config version < App version (upgrade)**:
   - Run migration scripts for each intermediate version
   - Add new fields with default values
   - Remove deprecated fields (with warning)
   - Update `version` field to current app version

2. **Config version > App version (downgrade)**:
   - Warn user about potential incompatibility
   - Continue with best-effort compatibility
   - Do not modify config file

3. **Config version == App version**:
   - No migration needed
   - Continue normally

### Migration Examples

#### Example: Adding a New Field

```go
// Migration from 1.0.0 to 1.1.0
func migrate_1_0_0_to_1_1_0(cfg *Config) {
    // Add new field with default value
    if cfg.Editor.BatchTimeout == "" {
        cfg.Editor.BatchTimeout = "1s"
    }
}
```

#### Example: Renaming a Field

```go
// Migration from 1.0.0 to 1.2.0
func migrate_1_1_0_to_1_2_0(cfg *Config) {
    // Rename field: auto_save -> auto_save_interval
    if cfg.Editor.AutoSave > 0 {
        cfg.Editor.AutoSaveInterval = fmt.Sprintf("%dms", cfg.Editor.AutoSave)
        cfg.Editor.AutoSave = 0 // Clear old field
    }
}
```

#### Example: Removing a Deprecated Field

```go
// Migration from 1.2.0 to 1.3.0
func migrate_1_2_0_to_1_3_0(cfg *Config) {
    // Remove deprecated field with warning
    if cfg.Editor.OldSetting != "" {
        log.Warn("Deprecated field 'editor.old_setting' removed. Use 'editor.new_setting' instead.")
        cfg.Editor.OldSetting = ""
    }
}
```

### Breaking Changes

Breaking changes in configuration format require a major version bump (e.g., 1.x.x → 2.0.0). Examples:

- Removing required fields
- Changing field types
- Restructuring configuration hierarchy
- Changing validation rules in incompatible ways

### Non-Breaking Changes

Non-breaking changes require a minor or patch version bump:

- Adding new optional fields
- Adding new enum values
- Relaxing validation rules
- Deprecating fields (with migration)

---

## Configuration Loading

### Load Order

1. Check for config file at `~/.promptstack/config.yaml`
2. If missing, run interactive setup wizard
3. Parse YAML configuration
4. Apply environment variable overrides
5. Validate all fields
6. Run migrations if needed
7. Return validated configuration

### Error Handling

- **Config file not found**: Run setup wizard
- **Invalid YAML**: Show error, suggest manual fix or recreate
- **Validation errors**: Show specific field errors, suggest fixes
- **Migration errors**: Show error, offer to continue with defaults or abort

### Setup Wizard

On first launch (or when config is missing), the application runs an interactive setup wizard:

1. Welcome message
2. Prompt for Claude API key (required)
3. Select Claude model (required)
4. Enable vim mode? (optional, default: false)
5. Review configuration summary
6. Create config file
7. Initialize data directory
8. Extract starter prompts
9. Initialize database
10. Launch application

---

## Configuration Examples

### Minimal Configuration

```yaml
version: "1.0.0"
api:
  claude_api_key: "sk-ant-api03-..."
  model: "claude-3-sonnet-20240229"
```

### Development Configuration

```yaml
version: "1.0.0"
api:
  claude_api_key: "sk-ant-api03-..."
  model: "claude-3-sonnet-20240229"
  max_retries: 5
  timeout: "60s"
editor:
  vim_mode: true
  auto_save_interval: "500ms"
  undo_stack_size: 200
logging:
  level: "debug"
  format: "text"
advanced:
  debug: true
  verbose: true
```

### Production Configuration

```yaml
version: "1.0.0"
api:
  claude_api_key: "sk-ant-api03-..."
  model: "claude-3-sonnet-20240229"
  max_retries: 3
  timeout: "30s"
  context:
    composition_max_percent: 25
    library_max_percent: 15
    warning_threshold: 15
editor:
  vim_mode: false
  auto_save_interval: "1s"
  undo_stack_size: 100
logging:
  level: "info"
  format: "json"
ui:
  theme: "catppuccin-mocha"
  ai_panel_width_percent: 40
history:
  max_files: 0
  max_age: "0"
```

### Custom Paths Configuration

```yaml
version: "1.0.0"
api:
  claude_api_key: "sk-ant-api03-..."
  model: "claude-3-sonnet-20240229"
paths:
  data_dir: "/custom/path/promptstack/data"
  library_dir: "/custom/path/promptstack/library"
  history_dir: "/custom/path/promptstack/history"
  database_file: "/custom/path/promptstack/history.db"
  index_file: "/custom/path/promptstack/index.json"
```

---

## Configuration Best Practices

### Security

- **Never commit config files** to version control
- **Use environment variables** for sensitive data (API keys)
- **Set appropriate file permissions** (0600 for config file)
- **Rotate API keys** regularly

### Performance

- **Adjust undo_stack_size** based on available memory
- **Enable cache_library** for faster startup (default: true)
- **Tune auto_save_interval** based on usage patterns
- **Limit max_library_prompts** to reduce API costs

### Usability

- **Enable vim_mode** if you're comfortable with vim
- **Choose appropriate theme** for your terminal
- **Adjust ai_panel_width_percent** based on screen size
- **Enable show_token_estimate** to monitor API usage

### Debugging

- **Set logging.level to debug** for troubleshooting
- **Enable advanced.debug** for verbose output
- **Check debug.log** for detailed error information
- **Use PROMPTSTACK_DEBUG=1** for temporary debugging

---

## Troubleshooting

### Common Issues

#### Issue: "Configuration file not found"

**Solution**: Run the application to trigger the setup wizard, or manually create the config file at `~/.promptstack/config.yaml`.

#### Issue: "Invalid API key"

**Solution**: Verify your Claude API key is correct and active. You can set it via environment variable: `export PROMPTSTACK_CLAUDE_API_KEY="your-key"`

#### Issue: "Invalid model name"

**Solution**: Check that the model name is valid. Valid models: `claude-3-sonnet-20240229`, `claude-3-opus-20240229`, `claude-3-haiku-20240307`

#### Issue: "Configuration validation failed"

**Solution**: Check the specific validation error message and fix the corresponding field in your config file.

#### Issue: "Migration failed"

**Solution**: Try manually updating your config file to match the current schema, or delete the config file and run the setup wizard again.

### Resetting Configuration

To reset to default configuration:

```bash
# Backup existing config
cp ~/.promptstack/config.yaml ~/.promptstack/config.yaml.backup

# Delete config file
rm ~/.promptstack/config.yaml

# Run promptstack to trigger setup wizard
promptstack
```

---

## Future Enhancements

Potential future configuration options:

- **Custom themes**: User-defined color schemes
- **Keybinding customization**: Custom keybinding configurations
- **Workspace-specific configs**: Per-project configuration files
- **Cloud sync**: Configuration synchronization across devices
- **Profiles**: Multiple configuration profiles for different use cases
- **Additional AI providers**: Support for more AI providers beyond Claude, MCP, and OpenAI
- **Additional storage backends**: Support for more storage backends beyond SQLite, PostgreSQL, and Neo4j

---

## References

- [Requirements Document](requirements.md)
- [Project Structure](project-structure.md)
- [Go Style Guide](go-style-guide.md)
- [Milestones](milestones.md)

---

**Document Status**: ✅ Complete  
**Implementation Status**: ⏳ Not Started  
**Next Steps**: Implement config loading and validation in Milestone 1