# Go Style Guide Updates for Scalability

**Document**: [`go-style-guide.md`](../go-style-guide.md)  
**Priority**: **Medium**  
**Status**: Should update before implementation begins

---

## Overview

This document details all changes required to [`go-style-guide.md`](../go-style-guide.md) to incorporate Phase 1 scalability patterns. These changes add sections on interface design, factory patterns, and middleware patterns.

---

## A. Add Interface Design Section

```markdown
### Interface Design

**Define Interfaces Where Used:**
Define interfaces in the package that uses them, not where they're implemented.

**Good:**
```go
// internal/ai/context.go (uses AIProvider)
type AIProvider interface {
    GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error)
    ApplyChanges(ctx context.Context, changes []Change) (string, error)
    EstimateTokens(ctx context.Context, text string) (int, error)
}

// internal/ai/claude.go (implements AIProvider)
type ClaudeProvider struct { ... }

func (c *ClaudeProvider) GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error) {
    // Implementation
}
```

**Bad:**
```go
// internal/ai/claude.go (defines interface in implementation package)
type AIProvider interface { ... } // ❌
```

**Interface Size:**
Keep interfaces small and focused:
- 3-5 methods maximum
- Single responsibility
- Clear purpose

**Interface Naming:**
- Use descriptive names: `AIProvider`, `PromptSource`, `CompositionRepository`
- Avoid generic names: `Provider`, `Source`, `Repository`

**Interface Composition:**
Compose small interfaces into larger ones when needed:
```go
// Small, focused interfaces
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
```

**Interface Testing:**
Use interfaces for testing with mocks:
```go
// In tests
type MockProvider struct {
    suggestions []Suggestion
    err         error
}

func (m *MockProvider) GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error) {
    return m.suggestions, m.err
}
```
```

---

## B. Add Factory Pattern Section

```markdown
### Factory Pattern

**Purpose:**
Create objects based on configuration without exposing creation logic.

**Example:**
```go
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

**Benefits:**
- Decouples creation from usage
- Easy to add new providers
- Configuration-driven instantiation
- Testable with mocks

**Usage:**
```go
// cmd/promptstack/main.go
provider, err := ai.NewProvider(cfg)
if err != nil {
    log.Fatal(err)
}
```

**Factory with Options:**
Use functional options for complex factories:
```go
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

**Repository Factory:**
```go
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
```

---

## C. Add Middleware Pattern Section

```markdown
### Middleware Pattern

**Purpose:**
Wrap objects with cross-cutting concerns (logging, caching, metrics).

**Example:**
```go
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

**Usage:**
```go
// cmd/promptstack/main.go
provider, _ := ai.NewProvider(cfg)

// Apply middleware
provider = ai.WithLogging(logger)(provider)
provider = ai.WithCaching(cache)(provider)
provider = ai.WithMetrics(metrics)(provider)
```

**Benefits:**
- Composable behavior
- Separation of concerns
- Easy to add/remove middleware
- Testable in isolation

**Middleware Implementation:**
```go
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

**Middleware Chaining:**
```go
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
```

---

## D. Add Event Pattern Section

```markdown
### Event Pattern

**Purpose:**
Decouple components using publish-subscribe pattern for domain events.

**Event Interface:**
```go
type Event interface {
    Type() string
    Timestamp() time.Time
    Payload() map[string]interface{}
}
```

**Event Dispatcher:**
```go
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

**Usage:**
```go
// Subscribe to events
dispatcher.Subscribe("composition.saved", func(e Event) {
    log.Printf("Composition saved: %v", e.Payload())
})

// Publish events
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
```

---

## E. Add Repository Pattern Section

```markdown
### Repository Pattern

**Purpose:**
Abstract data access logic behind an interface.

**Repository Interface:**
```go
type CompositionRepository interface {
    Save(ctx context.Context, comp Composition) error
    Load(ctx context.Context, id string) (Composition, error)
    Search(ctx context.Context, query string) ([]Composition, error)
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, opts ListOptions) ([]Composition, error)
}
```

**Implementation:**
```go
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

**Usage:**
```go
repo, _ := storage.NewRepository(cfg)
err := repo.Save(ctx, composition)
```

**Benefits:**
- Swappable implementations
- Testable with mocks
- Centralized data access logic
- Easy to add caching, logging, etc.
```

---

## F. Add Error Handling Patterns

```markdown
### Error Handling

**Wrap Errors with Context:**
```go
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

**Error Types:**
Define custom error types for better error handling:
```go
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

**Error Wrapping:**
Use `fmt.Errorf` with `%w` to wrap errors:
```go
func NewProvider(cfg *config.Config) (AIProvider, error) {
    if cfg.ClaudeAPIKey == "" {
        return nil, fmt.Errorf("missing API key: %w", ErrMissingConfig)
    }
    // ...
}
```
```

---

## G. Add Testing Patterns

```markdown
### Testing Patterns

**Table-Driven Tests:**
```go
func TestContextSelector_Select(t *testing.T) {
    tests := []struct {
        name        string
        composition string
        library     []Prompt
        want        []Prompt
        wantErr     bool
    }{
        {
            name:        "selects matching prompts",
            composition: "write code",
            library:     []Prompt{codePrompt, textPrompt},
            want:        []Prompt{codePrompt},
        },
        {
            name:        "empty composition",
            composition: "",
            library:     []Prompt{codePrompt},
            want:        nil,
            wantErr:     true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            selector := NewDefaultSelector(1000)
            got, err := selector.Select(context.Background(), tt.composition, tt.library)
            if (err != nil) != tt.wantErr {
                t.Errorf("Select() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Select() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Mock Interfaces:**
```go
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

**Test Helpers:**
```go
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
```

---

**Last Updated**: 2026-01-07  
**Related Documents**: 
- [`scalability-implementation-summary.md`](./scalability-implementation-summary.md)
- [`scalability-project-structure-updates.md`](./scalability-project-structure-updates.md)
- [`scalability-milestones-updates.md`](./scalability-milestones-updates.md)