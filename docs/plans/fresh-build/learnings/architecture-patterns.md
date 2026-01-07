# Architecture Patterns Key Learnings

**Purpose**: Key learnings and implementation patterns for architecture and design from previous PromptStack implementation.

**Related Milestones**: M1-M38 (All milestones)

**Related Documents**: 
- [`project-structure.md`](../project-structure.md) - Architecture structure
- [`go-style-guide.md`](../go-style-guide.md) - Go coding standards
- [`go-testing-guide.md`](../go-testing-guide.md) - Testing patterns

---

## Learning Categories

### Category 1: Database Schema Design

**Learning**: Use SQLite's FTS5 for full-text search

**Problem**: Need efficient full-text search for history compositions.

**Solution**: Separate tables with FTS5 virtual table for search

**Implementation Pattern**:
```go
// Schema design
CREATE TABLE compositions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_path TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    working_directory TEXT NOT NULL,
    content TEXT NOT NULL,
    character_count INTEGER NOT NULL,
    line_count INTEGER NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE VIRTUAL TABLE compositions_fts USING fts5(
    content,
    working_directory,
    content='compositions',
    content_rowid='id'
);

// Query with full-text search
SELECT * FROM compositions
WHERE id IN (
    SELECT rowid FROM compositions_fts
    WHERE compositions_fts MATCH ?
)
ORDER BY created_at DESC
LIMIT ?;
```

**Benefits**:
- Efficient full-text search
- Separate concerns (storage vs. search)
- Fast queries on metadata
- Built-in ranking and relevance

**Lesson**: Use SQLite's FTS5 for full-text search. It's purpose-built for this use case and much faster than LIKE queries. Separate the FTS table from the main table to keep concerns clean.

**Related Milestones**: M15, M16, M17

**When to Apply**: When implementing full-text search with SQLite

---

### Category 2: Configuration Management

**Learning**: Always validate configuration after loading

**Problem**: Need to load and validate user configuration.

**Solution**: YAML with validation

**Implementation Pattern**:
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

**Benefits**:
- Human-readable configuration
- Type-safe with struct tags
- Validation ensures correctness
- Easy to extend

**Lesson**: Always validate configuration after loading. Catch errors early and provide clear error messages. Use struct tags for YAML mapping and validation methods for business rules.

**Related Milestones**: M1, M2

**When to Apply**: When implementing configuration management

---

### Category 3: Starter Prompt Extraction

**Learning**: Never overwrite user data when extracting bundled resources

**Problem**: Need to extract starter prompts on first run and upgrades.

**Solution**: Version-aware extraction with existence checks

**Implementation Pattern**:
```go
func ExtractStarterPrompts(dataPath string, logger *zap.Logger) error {
    // Walk embedded filesystem
    err := fs.WalkDir(starterFS, "starter-prompts", func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            return nil
        }

        // Calculate destination path
        relPath, err := filepath.Rel("starter-prompts", path)
        if err != nil {
            return err
        }
        destPath := filepath.Join(dataPath, relPath)

        // Check if file already exists
        if _, err := os.Stat(destPath); err == nil {
            logger.Info("Skipping existing file", zap.String("path", destPath))
            return nil
        }

        // Create directory if needed
        if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
            return fmt.Errorf("failed to create directory: %w", err)
        }

        // Extract file
        content, err := fs.ReadFile(starterFS, path)
        if err != nil {
            return fmt.Errorf("failed to read embedded file: %w", err)
        }

        if err := os.WriteFile(destPath, content, 0644); err != nil {
            return fmt.Errorf("failed to write file: %w", err)
        }

        logger.Info("Extracted starter prompt", zap.String("path", destPath))
        return nil
    })

    return err
}
```

**Benefits**:
- Preserves user modifications
- Adds new prompts on upgrade
- Never overwrites user changes
- Idempotent operation
- Clear logging of extraction process

**Lesson**: When extracting bundled resources, always check for existing files. Never overwrite user data. This enables safe upgrades and preserves user customizations.

**Related Milestones**: M1, M2

**When to Apply**: When extracting bundled resources or templates

---

### Category 4: Logging Strategy

**Learning**: Use structured logging from the start with rotation

**Problem**: Need persistent logs for debugging without disk bloat.

**Solution**: Structured logging with rotation

**Implementation Pattern**:
```go
// Configuration
- File-based logging to `~/.promptstack/debug.log`
- JSON format for easy parsing
- Rotation at 10MB, keep last 3
- Level control via environment variable

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

**Benefits**:
- Persistent logs for debugging
- Automatic cleanup prevents disk bloat
- Structured logs enable filtering and analysis
- Environment variable control for different environments
- JSON format for easy parsing

**Lesson**: Use structured logging from the start. It pays dividends in debugging and monitoring. Configure rotation to prevent disk issues. Use environment variables for level control in different environments.

**Related Milestones**: M1, M2

**When to Apply**: When implementing logging in Go applications

---

### Category 5: Command Registry Pattern

**Learning**: Use a registry pattern for commands with validation

**Problem**: Need to manage and execute commands throughout the application.

**Solution**: Centralized command registration with handler functions

**Implementation Pattern**:
```go
type Command struct {
    ID          string
    Name        string
    Description string
    Category    string
    Handler     func() error
}

type Registry struct {
    commands map[string]*Command
}

func NewRegistry() *Registry {
    return &Registry{
        commands: make(map[string]*Command),
    }
}

func (r *Registry) Register(cmd *Command) error {
    if cmd.ID == "" {
        return fmt.Errorf("command ID cannot be empty")
    }
    if cmd.Handler == nil {
        return fmt.Errorf("command handler cannot be nil")
    }
    if _, exists := r.commands[cmd.ID]; exists {
        return fmt.Errorf("command with ID %s already registered", cmd.ID)
    }
    r.commands[cmd.ID] = cmd
    return nil
}

func (r *Registry) GetCommand(id string) *Command {
    return r.commands[id]
}

func (r *Registry) GetAll() []*Command {
    commands := make([]*Command, 0, len(r.commands))
    for _, cmd := range r.commands {
        commands = append(commands, cmd)
    }
    return commands
}

func (r *Registry) GetByCategory(category string) []*Command {
    var commands []*Command
    for _, cmd := range r.commands {
        if cmd.Category == category {
            commands = append(commands, cmd)
        }
    }
    return commands
}

func (r *Registry) GetCategories() []string {
    categories := make(map[string]bool)
    for _, cmd := range r.commands {
        categories[cmd.Category] = true
    }
    
    result := make([]string, 0, len(categories))
    for category := range categories {
        result = append(result, category)
    }
    sort.Strings(result)
    return result
}
```

**Benefits**:
- Centralized command management
- Type-safe command registration
- Easy to add new commands
- Validation prevents duplicate IDs
- Handler functions encapsulate command logic
- Category-based filtering

**Lesson**: Use a registry pattern for commands. Validate command registration (ID, handler, duplicates). This prevents runtime errors and makes command management maintainable.

**Related Milestones**: M24, M25

**When to Apply**: When implementing command systems

---

### Category 6: Confirmation Dialog Integration for Destructive Operations

**Learning**: Implement multi-step confirmation workflow with type assertion for Bubble Tea models

**Problem**: Need to prevent accidental destructive operations.

**Solution**: Multi-step confirmation workflow with type assertion

**Implementation Pattern**:
```go
// ui/cleanup/model.go
type Model struct {
    showConfirmation  bool
    confirmation      common.ConfirmationModel
    // ... other fields
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            if m.showPreview && !m.showConfirmation {
                // Show confirmation dialog before executing
                m.showConfirmation = true
                m.confirmation = common.DestructiveConfirmation(
                    "Confirm Cleanup",
                    fmt.Sprintf("Are you sure you want to delete %d files? This action cannot be undone.", m.previewResult.FileCount),
                )
                m.confirmation.SetOnConfirm(func() tea.Cmd {
                    if m.onExecute != nil {
                        return m.onExecute()
                    }
                    return nil
                })
                m.confirmation.SetOnCancel(func() tea.Cmd {
                    m.showConfirmation = false
                    return nil
                })
                m.confirmation.SetSize(m.width, m.height)
                return m, nil
            }
        }
    }
    
    // Update confirmation dialog if active
    if m.showConfirmation {
        var confirmCmd tea.Cmd
        model, confirmCmd := m.confirmation.Update(msg)
        m.confirmation = model.(common.ConfirmationModel)
        
        // Check if confirmation was handled
        if m.confirmation.IsConfirmed() || m.confirmation.IsCancelled() {
            m.showConfirmation = false
        }
        
        return m, confirmCmd
    }
    
    // ... rest of update logic
}

func (m Model) View() string {
    // Show confirmation dialog if active
    if m.showConfirmation {
        return m.confirmation.View()
    }
    
    // ... render normal UI
}
```

**Benefits**:
- Prevents accidental destructive operations
- User must type explicit confirmation text ("DELETE")
- Confirmation dialog overlays entire UI when active
- Clean separation between confirmation and main UI
- Type assertion handles Bubble Tea model interface
- Callback-based execution decouples UI from business logic

**Lesson**: For destructive operations, implement a multi-step confirmation workflow. Show confirmation dialog before executing. Use type assertion when updating Bubble Tea models that implement the tea.Model interface. Overlay confirmation dialog by returning its View() when active. This prevents accidental data loss and provides clear user feedback.

**Related Milestones**: M38

**When to Apply**: When implementing destructive operations

---

### Category 7: Type Assertion for Bubble Tea Model Updates

**Learning**: Use type assertion to convert interface back to concrete type

**Problem**: Bubble Tea's Update() method returns tea.Model interface, not concrete type.

**Solution**: Use type assertion to convert interface back to concrete type

**Implementation Pattern**:
```go
// Incorrect approach (compilation error):
m.confirmation, cmd = m.confirmation.Update(msg)
// Error: cannot use m.confirmation.Update(msg) (value of interface type tea.Model)
//        as common.ConfirmationModel value in assignment

// Correct approach with type assertion:
model, cmd := m.confirmation.Update(msg)
m.confirmation = model.(common.ConfirmationModel)
```

**Lesson**: When updating Bubble Tea models that implement the tea.Model interface, use type assertion to convert the returned interface back to the concrete type. This is necessary because Update() returns tea.Model interface for flexibility, but you need to concrete type to access its methods and fields.

**Related Milestones**: M38

**When to Apply**: When updating Bubble Tea child models from parent models

---

## Quick Reference

| Learning | Milestone | Priority |
|----------|-----------|----------|
| Database Schema Design | M15, M16, M17 | High |
| Configuration Management | M1, M2 | High |
| Starter Prompt Extraction | M1, M2 | Medium |
| Logging Strategy | M1, M2 | High |
| Command Registry Pattern | M24, M25 | High |
| Confirmation Dialog Integration for Destructive Operations | M38 | High |
| Type Assertion for Bubble Tea Model Updates | M38 | High |

---

**Last Updated**: 2026-01-07
**Source**: [`key-learnings.md`](../key-learnings.md)