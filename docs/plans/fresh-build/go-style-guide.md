# Go Code Style Guide for AI Assistants

**Target**: AI code generation for PromptStack
**Length**: Concise, actionable rules
**Focus**: Go idioms, clean architecture, testability

---

## Core Principles

1. **Clarity over cleverness** - Code should be obvious, not clever
2. **Explicit over implicit** - Make dependencies and behavior clear
3. **Simple over complex** - Prefer straightforward solutions
4. **Testable by design** - Write code that's easy to test

---

## Package Organization

### Naming
```go
// ✅ GOOD: Singular, lowercase, no underscores
package editor
package prompt
package library

// ❌ BAD: Plural, mixed case, underscores
package editors
package promptUtils
package prompt_manager
```

### Structure
```go
// Standard package layout:
package-name/
├── entity.go           // Core types, constructors
├── entity_test.go      // Tests
├── operation.go        // Operations on entities
├── operation_test.go   // Operation tests
```

### Package Comments
```go
// Package editor provides text editing functionality including
// cursor management, buffer operations, and undo/redo support.
package editor
```

---

## Type Design

### Constructors
```go
// ✅ GOOD: New() for single type, NewType() for multiple
func New(content string) *Buffer {
    return &Buffer{
        content: content,
        cursor:  0,
    }
}

func NewClient(cfg Config) (*Client, error) {
    if cfg.APIKey == "" {
        return nil, fmt.Errorf("API key required")
    }
    return &Client{apiKey: cfg.APIKey}, nil
}

// ❌ BAD: Don't use CreateBuffer, MakeBuffer, etc.
```

### Struct Design
```go
// ✅ GOOD: Exported type, unexported fields
type Buffer struct {
    content string  // unexported implementation
    cursor  int     // unexported state
}

// ✅ GOOD: Config structs for complex initialization
type Config struct {
    APIKey     string
    Model      string
    MaxRetries int
    Timeout    time.Duration
}

// ❌ BAD: Exported fields expose implementation
type Buffer struct {
    Content string  // allows direct manipulation
    Cursor  int     // breaks encapsulation
}
```

### Method Receivers
```go
// ✅ GOOD: Pointer for mutation, value for read-only
func (b *Buffer) Insert(text string) error { ... }  // mutates
func (b Buffer) Content() string { ... }            // reads

// ✅ GOOD: Be consistent - if any method uses pointer, all should
type Buffer struct { ... }
func (b *Buffer) Insert(text string) error { ... }
func (b *Buffer) Content() string { ... }  // pointer for consistency

// ❌ BAD: Mixing pointer and value receivers inconsistently
```

---

## Error Handling

### Error Creation
```go
// ✅ GOOD: Lowercase, no punctuation, wrap with %w
return fmt.Errorf("failed to load prompt: %w", err)
return fmt.Errorf("invalid placeholder type %q", typ)

// ✅ GOOD: Include context
return fmt.Errorf("failed to parse frontmatter in %s: %w", path, err)

// ❌ BAD: Uppercase, punctuation, no context
return fmt.Errorf("Error loading prompt!")
return fmt.Errorf("Failed: %v", err)  // use %w not %v
```

### Error Checking
```go
// ✅ GOOD: Check immediately, handle explicitly
prompts, err := library.Load(path)
if err != nil {
    return nil, fmt.Errorf("failed to load library: %w", err)
}

// ❌ BAD: Ignoring errors
prompts, _ := library.Load(path)  // never ignore errors
```

### Custom Errors
```go
// ✅ GOOD: Define error types for specific cases
type ValidationError struct {
    Type    string
    Message string
    Line    int
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s at line %d: %s", e.Type, e.Line, e.Message)
}
```

### Error Wrapping with Context
```go
// ✅ GOOD: Wrap errors with context
func (r *SQLiteRepository) Load(ctx context.Context, id string) (Composition, error) {
    var comp Composition
    err := r.db.QueryRowContext(ctx,
        "SELECT id, file_path, content FROM compositions WHERE id = ?", id).
        Scan(&comp.ID, &comp.FilePath, &comp.Content)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return Composition{}, fmt.Errorf("composition not found: %w", err)
        }
        return Composition{}, fmt.Errorf("failed to load composition: %w", err)
    }
    return comp, nil
}
```

### Error Types for Validation
```go
// ✅ GOOD: Use custom error types for validation
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// Usage
if cfg.AIProvider == "" {
    return &ValidationError{Field: "ai_provider", Message: "required"}
}
```

### Error Wrapping in Factories
```go
// ✅ GOOD: Wrap errors in factory functions
func NewProvider(cfg *config.Config) (AIProvider, error) {
    if cfg.ClaudeAPIKey == "" {
        return nil, fmt.Errorf("missing API key: %w", ErrMissingConfig)
    }
    // ...
}
```

---

## Interfaces

### Definition Location
```go
// ✅ GOOD: Define interfaces where they're USED, not implemented
// ui/workspace/model.go
type LibrarySearcher interface {
    Search(query string) ([]prompt.Prompt, error)
}

type Model struct {
    library LibrarySearcher  // interface, not concrete type
}

// internal/library/library.go
type Library struct { ... }
func (l *Library) Search(query string) ([]prompt.Prompt, error) {
    // Library implicitly implements LibrarySearcher
}

// ❌ BAD: Defining interface in implementation package
// internal/library/library.go
type LibrarySearcher interface { ... }  // wrong location
```

### Interface Size
```go
// ✅ GOOD: Small, focused interfaces (3-5 methods maximum)
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Searcher interface {
    Search(query string) ([]Prompt, error)
}

// ❌ BAD: Large, kitchen-sink interfaces
type LibraryManager interface {
    Load() error
    Search(query string) ([]Prompt, error)
    Validate() error
    Index() error
    Save() error
    // ... 10 more methods
}
```

### Interface Naming
```go
// ✅ GOOD: Descriptive, domain-specific names
type AIProvider interface { ... }
type PromptSource interface { ... }
type CompositionRepository interface { ... }

// ❌ BAD: Generic names
type Provider interface { ... }
type Source interface { ... }
type Repository interface { ... }
```

### Interface Composition
```go
// ✅ GOOD: Compose small interfaces into larger ones
type Reader interface {
    Read(ctx context.Context, id string) (Data, error)
}

type Writer interface {
    Write(ctx context.Context, data Data) error
}

// Composed interface
type Store interface {
    Reader
    Writer
}

// Usage
func NewService(store Store) *Service {
    // Can use any type that implements both Reader and Writer
}
```

### Interface Testing
```go
// ✅ GOOD: Use interfaces for testing with mocks
type MockProvider struct {
    suggestions []Suggestion
    err         error
}

func (m *MockProvider) GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error) {
    return m.suggestions, m.err
}

// In tests
func TestService(t *testing.T) {
    mock := &MockProvider{
        suggestions: []Suggestion{{Text: "test"}},
    }
    service := NewService(mock)
    // test with mock
}
```

---

## Dependency Management

### Injection Pattern
```go
// ✅ GOOD: Pass dependencies explicitly
func main() {
    cfg := config.Load()
    logger := logging.New(cfg.LogLevel)
    
    lib := library.New(library.Config{
        DataDir: cfg.DataDir,
        Logger:  logger,
    })
    
    app := ui.NewApp(cfg, lib, logger)
    app.Run()
}

// ❌ BAD: Global state
var GlobalConfig *Config  // avoid global state
var GlobalLogger *Logger  // makes testing hard

func (l *Loader) LoadPrompts() {
    dir := GlobalConfig.DataDir  // hidden dependency
}
```

### Dependency Direction
```
✅ GOOD Flow:
ui/ → internal/domain/ → internal/platform/

❌ BAD Flow:
internal/domain/ → ui/  (domain shouldn't know about UI)
internal/platform/ → internal/domain/  (infra shouldn't know business logic)
```

---

## Concurrency

### Bubble Tea Pattern
```go
// ✅ GOOD: Use tea.Cmd for async operations
func loadLibrary(path string) tea.Cmd {
    return func() tea.Msg {
        prompts, err := library.Load(path)
        if err != nil {
            return libraryLoadedMsg{err: err}
        }
        return libraryLoadedMsg{prompts: prompts}
    }
}

// ❌ BAD: Goroutines in domain logic
func (l *Library) Load(path string) ([]Prompt, error) {
    go func() {  // don't do this in domain packages
        // async work
    }()
}
```

### Synchronous Domain Logic
```go
// ✅ GOOD: Domain packages are synchronous
func (l *Loader) LoadPrompts(dir string) ([]Prompt, error) {
    // Synchronous operation
    // Let caller decide concurrency model
    return prompts, nil
}
```

---

## Testing

### Test Package Names
```go
// ✅ GOOD: Use _test suffix for black-box testing
package editor_test

import (
    "testing"
    "github.com/yourorg/promptstack/internal/editor"
)

func TestBuffer(t *testing.T) {
    buf := editor.New("content")
    // Test public API only
}

// ✅ ACCEPTABLE: Same package for white-box testing
package editor

func TestInternalHelper(t *testing.T) {
    // Test internal functions when necessary
}
```

### Table-Driven Tests
```go
// ✅ GOOD: Use table-driven tests for multiple cases
func TestParsePlaceholders(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    []Placeholder
        wantErr bool
    }{
        {
            name:  "valid text placeholder",
            input: "{{text:name}}",
            want:  []Placeholder{{Type: "text", Name: "name"}},
        },
        {
            name:    "invalid syntax",
            input:   "{{invalid",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParsePlaceholders(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Mocking
```go
// ✅ GOOD: Simple interface mocks
type mockLibrary struct {
    prompts []prompt.Prompt
    err     error
}

func (m *mockLibrary) Search(query string) ([]prompt.Prompt, error) {
    return m.prompts, m.err
}

// Use in tests
func TestWorkspace(t *testing.T) {
    mock := &mockLibrary{
        prompts: []prompt.Prompt{{Title: "test"}},
    }
    
    model := NewModel(mock)
    // test with mock
}
```

### Mock Interfaces for Testing
```go
// ✅ GOOD: Mock repository for testing
type MockRepository struct {
    compositions map[string]Composition
    err          error
}

func (m *MockRepository) Save(ctx context.Context, comp Composition) error {
    if m.err != nil {
        return m.err
    }
    m.compositions[comp.ID] = comp
    return nil
}

func (m *MockRepository) Load(ctx context.Context, id string) (Composition, error) {
    if m.err != nil {
        return Composition{}, m.err
    }
    comp, ok := m.compositions[id]
    if !ok {
        return Composition{}, sql.ErrNoRows
    }
    return comp, nil
}
```

### Test Helpers
```go
// ✅ GOOD: Reusable test helpers
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatal(err)
    }
    
    // Create schema
    _, err = db.Exec(`
        CREATE TABLE compositions (
            id TEXT PRIMARY KEY,
            file_path TEXT NOT NULL,
            content TEXT NOT NULL
        )
    `)
    if err != nil {
        t.Fatal(err)
    }
    
    return db
}
```

---

## Code Organization

### File Naming
```go
// ✅ GOOD: Descriptive, lowercase, underscores for multi-word
buffer.go
buffer_test.go
placeholder.go
placeholder_test.go
http_client.go
http_client_test.go

// ❌ BAD: CamelCase, unclear names
Buffer.go
bufferStuff.go
utils.go  // too generic
```

### Function Length
```go
// ✅ GOOD: Short, focused functions (< 50 lines ideal)
func (b *Buffer) Insert(text string) error {
    if err := b.validateCursor(); err != nil {
        return err
    }
    b.content = b.content[:b.cursor] + text + b.content[b.cursor:]
    b.cursor += len(text)
    return nil
}

// ❌ BAD: Long, multi-purpose functions (> 100 lines)
func (b *Buffer) DoEverything() error {
    // 200 lines of mixed concerns
}
```

### Comments
```go
// ✅ GOOD: Explain WHY, not WHAT
// Use exponential backoff to avoid overwhelming the API
// during temporary outages or rate limiting.
backoff := time.Duration(1<<uint(retry)) * time.Second

// ✅ GOOD: Document exported functions
// ParsePlaceholders extracts all placeholder definitions from content.
// It returns an error if any placeholder has invalid syntax.
func ParsePlaceholders(content string) ([]Placeholder, error) { ... }

// ❌ BAD: Obvious comments
// Set x to 5
x := 5

// ❌ BAD: Commented-out code
// func oldImplementation() { ... }
```

---

## Common Patterns

### Options Pattern
```go
// ✅ GOOD: For optional configuration
type Option func(*Client)

func WithTimeout(d time.Duration) Option {
    return func(c *Client) {
        c.timeout = d
    }
}

func WithRetries(n int) Option {
    return func(c *Client) {
        c.maxRetries = n
    }
}

func NewClient(apiKey string, opts ...Option) *Client {
    c := &Client{
        apiKey:     apiKey,
        timeout:    30 * time.Second,  // defaults
        maxRetries: 3,
    }
    for _, opt := range opts {
        opt(c)
    }
    return c
}

// Usage
client := NewClient(key, WithTimeout(60*time.Second), WithRetries(5))
```

### Builder Pattern (Avoid)
```go
// ❌ BAD: Unnecessary in Go, use Config struct instead
client := NewClientBuilder().
    WithAPIKey(key).
    WithTimeout(60).
    Build()

// ✅ GOOD: Use Config struct
client := NewClient(Config{
    APIKey:  key,
    Timeout: 60 * time.Second,
})
```

### Context Usage
```go
// ✅ GOOD: First parameter, named ctx
func (c *Client) SendMessage(ctx context.Context, req Request) (*Response, error) {
    // Check context before expensive operations
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    // Use context in HTTP requests
    httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
}

// ❌ BAD: Context not first parameter
func (c *Client) SendMessage(req Request, ctx context.Context) (*Response, error) {
```

---

## Factory Pattern

### Purpose
Create objects based on configuration without exposing creation logic.

### Basic Factory
```go
// ✅ GOOD: Simple factory for provider selection
// internal/ai/factory.go
package ai

import "github.com/yourorg/promptstack/internal/config"

func NewProvider(cfg *config.Config) (AIProvider, error) {
    switch cfg.AIProvider {
    case "claude":
        return NewClaudeProvider(cfg.ClaudeAPIKey, cfg.Model)
    case "mcp":
        return NewMCPProvider(cfg.MCPHost)
    case "openai":
        return NewOpenAIProvider(cfg.OpenAIAPIKey, cfg.Model)
    default:
        return NewClaudeProvider(cfg.ClaudeAPIKey, cfg.Model)
    }
}
```

### Benefits
- Decouples creation from usage
- Easy to add new providers
- Configuration-driven instantiation
- Testable with mocks

### Usage
```go
// cmd/promptstack/main.go
provider, err := ai.NewProvider(cfg)
if err != nil {
    log.Fatal(err)
}
```

### Factory with Options
```go
// ✅ GOOD: Functional options for complex factories
type ProviderOption func(*providerConfig)

type providerConfig struct {
    timeout time.Duration
    retries int
}

func WithTimeout(timeout time.Duration) ProviderOption {
    return func(cfg *providerConfig) {
        cfg.timeout = timeout
    }
}

func WithRetries(retries int) ProviderOption {
    return func(cfg *providerConfig) {
        cfg.retries = retries
    }
}

func NewProvider(cfg *config.Config, opts ...ProviderOption) (AIProvider, error) {
    pc := &providerConfig{
        timeout: 30 * time.Second,
        retries: 3,
    }
    
    for _, opt := range opts {
        opt(pc)
    }
    
    // Create provider with config...
}
```

### Repository Factory
```go
// ✅ GOOD: Factory for storage backends
// internal/storage/factory.go
package storage

import "github.com/yourorg/promptstack/internal/config"

func NewRepository(cfg *config.Config) (CompositionRepository, error) {
    switch cfg.Storage {
    case "sqlite":
        return NewSQLiteRepository(cfg.DatabasePath)
    case "postgres":
        return NewPostgreSQLRepository(cfg.PostgresURL)
    case "graph":
        return NewGraphRepository(cfg.Neo4jURL)
    default:
        return NewSQLiteRepository(cfg.DatabasePath)
    }
}
```

---

## Middleware Pattern

### Purpose
Wrap objects with cross-cutting concerns (logging, caching, metrics).

### Middleware Definition
```go
// ✅ GOOD: Middleware as a function type
// internal/ai/middleware.go
package ai

type ProviderMiddleware func(AIProvider) AIProvider

func WithLogging(logger Logger) ProviderMiddleware {
    return func(provider AIProvider) AIProvider {
        return &loggingProvider{
            provider: provider,
            logger:   logger,
        }
    }
}

func WithCaching(cache Cache) ProviderMiddleware {
    return func(provider AIProvider) AIProvider {
        return &cachingProvider{
            provider: provider,
            cache:    cache,
        }
    }
}
```

### Usage
```go
// cmd/promptstack/main.go
provider, _ := ai.NewProvider(cfg)

// Apply middleware
provider = ai.WithLogging(logger)(provider)
provider = ai.WithCaching(cache)(provider)
provider = ai.WithMetrics(metrics)(provider)
```

### Benefits
- Composable behavior
- Separation of concerns
- Easy to add/remove middleware
- Testable in isolation

### Middleware Implementation
```go
// ✅ GOOD: Implement middleware with struct embedding
type loggingProvider struct {
    provider AIProvider
    logger   Logger
}

func (l *loggingProvider) GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error) {
    l.logger.Printf("Getting suggestions for prompt: %s", prompt)
    suggestions, err := l.provider.GetSuggestions(ctx, prompt)
    if err != nil {
        l.logger.Printf("Error getting suggestions: %v", err)
    } else {
        l.logger.Printf("Got %d suggestions", len(suggestions))
    }
    return suggestions, err
}

type cachingProvider struct {
    provider AIProvider
    cache    Cache
}

func (c *cachingProvider) GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error) {
    // Check cache
    if cached, ok := c.cache.Get(prompt); ok {
        return cached, nil
    }
    
    // Call provider
    suggestions, err := c.provider.GetSuggestions(ctx, prompt)
    if err != nil {
        return nil, err
    }
    
    // Cache result
    c.cache.Set(prompt, suggestions)
    return suggestions, nil
}
```

### Middleware Chaining
```go
// ✅ GOOD: Helper function for chaining middleware
func ApplyMiddleware(provider AIProvider, middlewares ...ProviderMiddleware) AIProvider {
    for _, mw := range middlewares {
        provider = mw(provider)
    }
    return provider
}

// Usage
provider := ApplyMiddleware(
    ai.NewProvider(cfg),
    ai.WithLogging(logger),
    ai.WithCaching(cache),
    ai.WithMetrics(metrics),
)
```

---

## Event Pattern

### Purpose
Decouple components using publish-subscribe pattern for domain events.

### Event Interface
```go
// ✅ GOOD: Define event interface
type Event interface {
    Type() string
    Timestamp() time.Time
    Payload() map[string]interface{}
}
```

### Event Dispatcher
```go
// ✅ GOOD: Simple event dispatcher
type Dispatcher struct {
    mu       sync.RWMutex
    handlers map[string][]EventHandler
}

type EventHandler func(Event)

func (d *Dispatcher) Subscribe(eventType string, handler EventHandler) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.handlers[eventType] = append(d.handlers[eventType], handler)
}

func (d *Dispatcher) Publish(event Event) {
    d.mu.RLock()
    handlers := d.handlers[event.Type()]
    d.mu.RUnlock()
    
    for _, handler := range handlers {
        go handler(event) // Async handling
    }
}
```

### Usage
```go
// ✅ GOOD: Subscribe to events
dispatcher.Subscribe("composition.saved", func(e Event) {
    log.Printf("Composition saved: %v", e.Payload())
})

// ✅ GOOD: Publish events
dispatcher.Publish(&CompositionSavedEvent{
    BaseEvent: BaseEvent{
        eventType: "composition.saved",
        timestamp: time.Now(),
        payload:   map[string]interface{}{"id": "123"},
    },
    CompositionID: "123",
    FilePath:     "/path/to/file.md",
})
```

---

## Repository Pattern

### Purpose
Abstract data access logic behind an interface.

### Repository Interface
```go
// ✅ GOOD: Define repository interface
type CompositionRepository interface {
    Save(ctx context.Context, comp Composition) error
    Load(ctx context.Context, id string) (Composition, error)
    Search(ctx context.Context, query string) ([]Composition, error)
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, opts ListOptions) ([]Composition, error)
}
```

### Implementation
```go
// ✅ GOOD: SQLite implementation
type SQLiteRepository struct {
    db *sql.DB
}

func (r *SQLiteRepository) Save(ctx context.Context, comp Composition) error {
    _, err := r.db.ExecContext(ctx,
        "INSERT INTO compositions (id, file_path, content) VALUES (?, ?, ?)",
        comp.ID, comp.FilePath, comp.Content)
    return err
}
```

### Usage
```go
// ✅ GOOD: Use repository through interface
repo, _ := storage.NewRepository(cfg)
err := repo.Save(ctx, composition)
```

### Benefits
- Swappable implementations
- Testable with mocks
- Centralized data access logic
- Easy to add caching, logging, etc.

---

## Anti-Patterns to Avoid

### Don't
```go
// ❌ init() functions (except for registration)
func init() {
    config = loadConfig()  // side effects in init
}

// ❌ Panic in library code
func LoadConfig() Config {
    cfg, err := load()
    if err != nil {
        panic(err)  // let caller handle errors
    }
    return cfg
}

// ❌ Naked returns
func Calculate() (result int, err error) {
    result = 42
    return  // unclear what's being returned
}

// ❌ Generic variable names in large scopes
func Process() {
    d := getData()  // what is 'd'?
    // 50 lines later...
    use(d)  // hard to remember what 'd' is
}

// ❌ Else after return
func Check(x int) string {
    if x > 0 {
        return "positive"
    } else {  // unnecessary else
        return "non-positive"
    }
}
```

### Do
```go
// ✅ Explicit initialization
func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }
}

// ✅ Return errors
func LoadConfig() (Config, error) {
    cfg, err := load()
    if err != nil {
        return Config{}, fmt.Errorf("failed to load config: %w", err)
    }
    return cfg, nil
}

// ✅ Named returns only for documentation
func Calculate() (result int, err error) {
    result = 42
    err = nil
    return result, err  // explicit
}

// ✅ Descriptive names
func Process() {
    userData := getData()
    // 50 lines later...
    use(userData)  // clear what this is
}

// ✅ Early return
func Check(x int) string {
    if x > 0 {
        return "positive"
    }
    return "non-positive"
}
```

---

## Project-Specific Rules

### UI Components (Bubble Tea)
```go
// ✅ GOOD: Separate files for model, update, view
// workspace/model.go
type Model struct { ... }

// workspace/update.go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { ... }

// workspace/view.go
func (m Model) View() string { ... }

// workspace/messages.go
type libraryLoadedMsg struct { ... }
```

### Theme Usage
```go
// ✅ GOOD: Always use theme helpers
import "github.com/yourorg/promptstack/ui/theme"

modalStyle := theme.ModalStyle()
titleStyle := theme.ModalTitleStyle()

// ❌ BAD: Hard-coded colors
modalStyle := lipgloss.NewStyle().
    Background(lipgloss.Color("#1e1e2e"))  // use theme constant
```

### Logging
```go
// ✅ GOOD: Structured logging with zap
logger.Info("loading library",
    zap.String("path", path),
    zap.Int("prompt_count", len(prompts)),
)

// ❌ BAD: Printf-style logging
log.Printf("Loading library from %s with %d prompts", path, len(prompts))
```

---

## Quick Reference

**Naming**: `camelCase` (unexported), `PascalCase` (exported)
**Packages**: singular, lowercase, no underscores
**Errors**: lowercase, no punctuation, wrap with `%w`
**Interfaces**: define where used, keep small
**Tests**: table-driven, black-box when possible
**Comments**: explain WHY, document exported APIs
**Files**: one primary type per file, co-locate tests

---

**Remember**: Write code for humans first, computers second. When in doubt, choose clarity.