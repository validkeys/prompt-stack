# Project Structure Updates for Scalability

**Document**: [`project-structure.md`](../project-structure.md)  
**Priority**: **Critical**  
**Status**: Must update before any implementation begins

---

## Overview

This document details all changes required to [`project-structure.md`](../project-structure.md) to incorporate Phase 1 scalability abstractions. These changes add new packages and refactor existing ones to support multiple AI providers, storage backends, and prompt sources.

---

## A. Add New Packages to AI Domain

### Current Structure
```
internal/ai/
├── client.go              # Claude API client
├── context.go             # Context selection
├── tokens.go              # Token estimation
├── suggestions.go         # Suggestion parsing
└── diff.go                # Diff generation
```

### Updated Structure
```
internal/ai/
├── provider.go            # NEW: AIProvider interface
├── claude.go              # NEW: Claude implementation
├── selector.go            # NEW: ContextSelector interface
├── context.go             # REFACTOR: Implement ContextSelector
├── middleware.go          # NEW: ProviderMiddleware type
├── tokens.go              # Unchanged
├── suggestions.go         # Unchanged
└── diff.go                # Unchanged
```

### Code to Add

#### internal/ai/provider.go
```go
package ai

import "context"

// AIProvider defines the interface for AI providers
type AIProvider interface {
    // GetSuggestions returns AI suggestions for the given prompt
    GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error)
    
    // ApplyChanges applies suggested changes to the composition
    ApplyChanges(ctx context.Context, changes []Change) (string, error)
    
    // EstimateTokens estimates token count for text
    EstimateTokens(ctx context.Context, text string) (int, error)
}

// Change represents a suggested edit
type Change struct {
    StartLine    int
    EndLine      int
    NewContent   string
    Description  string
}

// Suggestion represents an AI suggestion
type Suggestion struct {
    Type        string
    Title       string
    Description string
    Changes     []Change
}
```

#### internal/ai/claude.go
```go
package ai

import (
    "context"
    "github.com/anthropics/anthropic-sdk-go"
)

// ClaudeProvider implements AIProvider for Claude API
type ClaudeProvider struct {
    client *anthropic.Client
    model  string
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(apiKey, model string) *ClaudeProvider {
    client := anthropic.NewClient(apiKey)
    return &ClaudeProvider{
        client: client,
        model:  model,
    }
}

// GetSuggestions implements AIProvider
func (c *ClaudeProvider) GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error) {
    // Implementation using Claude API
    // ...
}

// ApplyChanges implements AIProvider
func (c *ClaudeProvider) ApplyChanges(ctx context.Context, changes []Change) (string, error) {
    // Implementation
    // ...
}

// EstimateTokens implements AIProvider
func (c *ClaudeProvider) EstimateTokens(ctx context.Context, text string) (int, error) {
    // Implementation
    // ...
}
```

#### internal/ai/selector.go
```go
package ai

import "context"

// ContextSelector defines the interface for context selection algorithms
type ContextSelector interface {
    // Select returns relevant prompts from the library
    Select(ctx context.Context, composition string, library []Prompt) ([]Prompt, error)
    
    // Name returns the selector name
    Name() string
}

// DefaultSelector implements basic context selection
type DefaultSelector struct {
    tokenBudget int
}

// NewDefaultSelector creates a new default selector
func NewDefaultSelector(tokenBudget int) *DefaultSelector {
    return &DefaultSelector{
        tokenBudget: tokenBudget,
    }
}

// Select implements ContextSelector
func (s *DefaultSelector) Select(ctx context.Context, composition string, library []Prompt) ([]Prompt, error) {
    // Implementation using scoring algorithm from requirements.md
    // ...
}

// Name implements ContextSelector
func (s *DefaultSelector) Name() string {
    return "default"
}
```

#### internal/ai/middleware.go
```go
package ai

// ProviderMiddleware wraps an AIProvider with additional behavior
type ProviderMiddleware func(AIProvider) AIProvider

// WithLogging adds logging to provider calls
func WithLogging(logger Logger) ProviderMiddleware {
    return func(provider AIProvider) AIProvider {
        return &loggingProvider{
            provider: provider,
            logger:   logger,
        }
    }
}

// WithCaching adds caching to provider calls
func WithCaching(cache Cache) ProviderMiddleware {
    return func(provider AIProvider) AIProvider {
        return &cachingProvider{
            provider: provider,
            cache:    cache,
        }
    }
}

// WithMetrics adds metrics collection to provider calls
func WithMetrics(metrics Metrics) ProviderMiddleware {
    return func(provider AIProvider) AIProvider {
        return &metricsProvider{
            provider: provider,
            metrics:  metrics,
        }
    }
}

// loggingProvider wraps provider with logging
type loggingProvider struct {
    provider AIProvider
    logger   Logger
}

// cachingProvider wraps provider with caching
type cachingProvider struct {
    provider AIProvider
    cache    Cache
}

// metricsProvider wraps provider with metrics
type metricsProvider struct {
    provider AIProvider
    metrics  Metrics
}
```

---

## B. Add New Packages to Storage Domain

### Current Structure
```
internal/history/
├── manager.go
├── database.go
├── storage.go
├── sync.go
├── search.go
├── cleanup.go
└── listing.go
```

### Updated Structure
```
internal/storage/              # NEW: Storage abstraction layer
├── repository.go              # NEW: CompositionRepository interface
├── sqlite.go                  # NEW: SQLite implementation
├── factory.go                 # NEW: Repository factory
├── postgres.go                # FUTURE: PostgreSQL implementation
└── graph.go                   # FUTURE: Neo4j implementation

internal/history/              # REFACTOR: Use repository pattern
├── manager.go                # REFACTOR: Use repository
├── database.go                # DEPRECATE: Move to storage/sqlite.go
├── storage.go                # DEPRECATE: Move to storage/sqlite.go
├── sync.go                   # REFACTOR: Use repository
├── search.go                 # REFACTOR: Use repository
├── cleanup.go                # REFACTOR: Use repository
└── listing.go                # REFACTOR: Use repository
```

### Code to Add

#### internal/storage/repository.go
```go
package storage

import "context"

// CompositionRepository defines the interface for composition storage
type CompositionRepository interface {
    // Save saves a composition
    Save(ctx context.Context, comp Composition) error
    
    // Load loads a composition by ID
    Load(ctx context.Context, id string) (Composition, error)
    
    // Search searches compositions by query
    Search(ctx context.Context, query string) ([]Composition, error)
    
    // Delete deletes a composition
    Delete(ctx context.Context, id string) error
    
    // List lists compositions with options
    List(ctx context.Context, opts ListOptions) ([]Composition, error)
}

// Composition represents a saved composition
type Composition struct {
    ID               string
    FilePath         string
    CreatedAt        time.Time
    WorkingDirectory string
    Content          string
    CharacterCount   int
    LineCount        int
    UpdatedAt        time.Time
}

// ListOptions defines options for listing compositions
type ListOptions struct {
    Limit      int
    Offset     int
    SortBy     string
    SortOrder  string
    WorkingDir string
}
```

#### internal/storage/sqlite.go
```go
package storage

import (
    "context"
    "database/sql"
    "modernc.org/sqlite"
)

// SQLiteRepository implements CompositionRepository for SQLite
type SQLiteRepository struct {
    db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }
    
    // Create tables
    if err := createSchema(db); err != nil {
        return nil, err
    }
    
    return &SQLiteRepository{db: db}, nil
}

// Save implements CompositionRepository
func (r *SQLiteRepository) Save(ctx context.Context, comp Composition) error {
    // Implementation
    // ...
}

// Load implements CompositionRepository
func (r *SQLiteRepository) Load(ctx context.Context, id string) (Composition, error) {
    // Implementation
    // ...
}

// Search implements CompositionRepository
func (r *SQLiteRepository) Search(ctx context.Context, query string) ([]Composition, error) {
    // Implementation using FTS5
    // ...
}

// Delete implements CompositionRepository
func (r *SQLiteRepository) Delete(ctx context.Context, id string) error {
    // Implementation
    // ...
}

// List implements CompositionRepository
func (r *SQLiteRepository) List(ctx context.Context, opts ListOptions) ([]Composition, error) {
    // Implementation
    // ...
}

func createSchema(db *sql.DB) error {
    // Schema creation from DATABASE-SCHEMA.md
    // ...
}
```

#### internal/storage/factory.go
```go
package storage

import (
    "github.com/yourorg/promptstack/internal/config"
)

// NewRepository creates a repository based on config
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

## C. Add New Packages to Library Domain

### Current Structure
```
internal/library/
├── library.go
├── loader.go
├── index.go
├── scorer.go
├── search.go
└── validator.go
```

### Updated Structure
```
internal/library/
├── source.go                 # NEW: PromptSource interface
├── filesystem.go             # NEW: Filesystem implementation
├── cache.go                 # NEW: Prompt cache
├── library.go               # REFACTOR: Use PromptSource
├── loader.go                # DEPRECATE: Move to filesystem.go
├── index.go                 # Unchanged
├── scorer.go                # Unchanged
├── search.go                # Unchanged
├── validator.go             # Unchanged
├── specialist.go            # FUTURE: MCP specialist source
└── remote.go                # FUTURE: Remote repository source
```

### Code to Add

#### internal/library/source.go
```go
package library

import "context"

// PromptSource defines the interface for prompt sources
type PromptSource interface {
    // Load loads all prompts from the source
    Load(ctx context.Context) ([]Prompt, error)
    
    // Get loads a specific prompt by ID
    Get(ctx context.Context, id string) (Prompt, error)
    
    // Watch watches for prompt changes
    Watch(ctx context.Context) (<-chan PromptEvent, error)
}

// PromptEvent represents a prompt change event
type PromptEvent struct {
    Type    EventType
    Prompt  Prompt
}

// EventType represents the type of event
type EventType string

const (
    EventCreated  EventType = "created"
    EventUpdated  EventType = "updated"
    EventDeleted  EventType = "deleted"
)
```

#### internal/library/filesystem.go
```go
package library

import (
    "context"
    "os"
    "path/filepath"
)

// FilesystemSource implements PromptSource for filesystem
type FilesystemSource struct {
    rootDir string
    cache   PromptCache
}

// NewFilesystemSource creates a new filesystem source
func NewFilesystemSource(rootDir string, cache PromptCache) *FilesystemSource {
    return &FilesystemSource{
        rootDir: rootDir,
        cache:   cache,
    }
}

// Load implements PromptSource
func (s *FilesystemSource) Load(ctx context.Context) ([]Prompt, error) {
    // Implementation scanning filesystem
    // ...
}

// Get implements PromptSource
func (s *FilesystemSource) Get(ctx context.Context, id string) (Prompt, error) {
    // Implementation
    // ...
}

// Watch implements PromptSource
func (s *FilesystemSource) Watch(ctx context.Context) (<-chan PromptEvent, error) {
    // Implementation using fsnotify
    // ...
}
```

#### internal/library/cache.go
```go
package library

// PromptCache defines the interface for prompt caching
type PromptCache interface {
    // Get retrieves a prompt from cache
    Get(id string) (Prompt, bool)
    
    // Set stores a prompt in cache
    Set(id string, prompt Prompt)
    
    // Invalidate removes a prompt from cache
    Invalidate(id string)
    
    // Clear clears all cached prompts
    Clear()
}

// MemoryCache implements PromptCache with in-memory storage
type MemoryCache struct {
    prompts map[string]Prompt
}

// NewMemoryCache creates a new memory cache
func NewMemoryCache() *MemoryCache {
    return &MemoryCache{
        prompts: make(map[string]Prompt),
    }
}

// Get implements PromptCache
func (c *MemoryCache) Get(id string) (Prompt, bool) {
    prompt, ok := c.prompts[id]
    return prompt, ok
}

// Set implements PromptCache
func (c *MemoryCache) Set(id string, prompt Prompt) {
    c.prompts[id] = prompt
}

// Invalidate implements PromptCache
func (c *MemoryCache) Invalidate(id string) {
    delete(c.prompts, id)
}

// Clear implements PromptCache
func (c *MemoryCache) Clear() {
    c.prompts = make(map[string]Prompt)
}
```

---

## D. Add New Packages for Domain Events

### New Structure
```
internal/events/               # NEW: Domain events system
├── events.go                # Event types
└── dispatcher.go            # Event dispatcher
```

### Code to Add

#### internal/events/events.go
```go
package events

import "time"

// Event defines the interface for domain events
type Event interface {
    // Type returns the event type
    Type() string
    
    // Timestamp returns when the event occurred
    Timestamp() time.Time
    
    // Payload returns the event payload
    Payload() map[string]interface{}
}

// BaseEvent provides common event functionality
type BaseEvent struct {
    eventType  string
    timestamp  time.Time
    payload    map[string]interface{}
}

// Type implements Event
func (e *BaseEvent) Type() string {
    return e.eventType
}

// Timestamp implements Event
func (e *BaseEvent) Timestamp() time.Time {
    return e.timestamp
}

// Payload implements Event
func (e *BaseEvent) Payload() map[string]interface{} {
    return e.payload
}

// CompositionSavedEvent is emitted when a composition is saved
type CompositionSavedEvent struct {
    BaseEvent
    CompositionID string
    FilePath      string
}

// PromptUsedEvent is emitted when a prompt is used
type PromptUsedEvent struct {
    BaseEvent
    PromptID string
    Context  string
}

// SuggestionAcceptedEvent is emitted when a suggestion is accepted
type SuggestionAcceptedEvent struct {
    BaseEvent
    SuggestionType string
    CompositionID  string
}
```

#### internal/events/dispatcher.go
```go
package events

import (
    "sync"
)

// EventHandler handles events
type EventHandler func(Event)

// Dispatcher manages event subscriptions and publishing
type Dispatcher struct {
    mu       sync.RWMutex
    handlers map[string][]EventHandler
}

// NewDispatcher creates a new event dispatcher
func NewDispatcher() *Dispatcher {
    return &Dispatcher{
        handlers: make(map[string][]EventHandler),
    }
}

// Subscribe registers a handler for an event type
func (d *Dispatcher) Subscribe(eventType string, handler EventHandler) {
    d.mu.Lock()
    defer d.mu.Unlock()
    
    d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Publish publishes an event to all registered handlers
func (d *Dispatcher) Publish(event Event) {
    d.mu.RLock()
    handlers := d.handlers[event.Type()]
    d.mu.RUnlock()
    
    for _, handler := range handlers {
        go handler(event) // Async handling
    }
}
```

---

## E. Update Domain Descriptions

Update the domain descriptions in project-structure.md to reflect new abstractions:

### AI Domain (updated)
```
### 5. AI Domain
Claude API integration and suggestion system with provider abstraction.
- AIProvider interface for multiple AI providers
- ContextSelector interface for pluggable context selection
- ProviderMiddleware for cross-cutting concerns
- Context selection algorithm
- Token estimation and budgeting
- Suggestion parsing and diff generation
```

### Storage Domain (new)
```
### 9. Storage Domain (NEW)
Composition persistence with repository pattern.
- CompositionRepository interface for multiple storage backends
- SQLite implementation (current)
- PostgreSQL implementation (future)
- Neo4j implementation (future)
- Factory pattern for repository instantiation
```

### Library Domain (updated)
```
### 3. Library Domain
Managing collections of prompts with source abstraction.
- PromptSource interface for multiple prompt sources
- Filesystem implementation (current)
- MCP specialist source (future)
- Remote repository source (future)
- Prompt caching layer
- Prompt indexing and scoring
- Fuzzy search
- Library validation
- Usage tracking
```

### Events Domain (new)
```
### 10. Events Domain (NEW)
Domain events for decoupling components.
- Event types (CompositionSaved, PromptUsed, SuggestionAccepted)
- Event dispatcher for pub/sub pattern
- Async event handling
```

---

**Last Updated**: 2026-01-07  
**Related Documents**: 
- [`scalability-implementation-summary.md`](./scalability-implementation-summary.md)
- [`scalability-config-schema-updates.md`](./scalability-config-schema-updates.md)
- [`scalability-milestones-updates.md`](./scalability-milestones-updates.md)