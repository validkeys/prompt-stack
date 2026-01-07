# PromptStack → Project OS: Architectural Evolution Strategy

## Overview

This document outlines how PromptStack's current architecture can naturally evolve toward the broader Project OS vision using idiomatic Go patterns. The goal is to build PromptStack today in a way that makes future extensibility straightforward without over-engineering.

**Key Principle:** Use interfaces and abstraction layers where they provide value, not everywhere. Implement patterns incrementally as needs arise.

---

## The Vision Connection

### PromptStack Today
- CLI tool for AI prompt composition
- Prompt library with workflows, commands, decorations, rules
- Claude API integration
- SQLite-based history
- TUI with Bubble Tea

### Project OS Tomorrow
- Knowledge & Context Department for organizations
- Queryable knowledge base of business and code
- MCP-based specialist servers (Tax, Security, Compliance, etc.)
- PostgreSQL + Neo4j for structured and graph data
- Implementation planning graphs
- Complete SDLC integration

### The Bridge
```mermaid
graph TD
    A[PromptStack Today] -->|Evolves Into| B[Project OS Tomorrow]
    
    A1[Prompt Library] -->|Becomes| B1[Domain Specialist Servers]
    A2[Claude Integration] -->|Becomes| B2[MCP Host/Orchestrator]
    A3[SQLite History] -->|Becomes| B3[PostgreSQL Knowledge Base]
    A4[Composition Workspace] -->|Becomes| B4[Implementation Planning Graphs]
    A5[File References] -->|Becomes| B5[AST-Based Code Indexing]
    A6[Placeholder System] -->|Becomes| B6[Context Variables & Templates]
    
    style A fill:#e1f5ff
    style B fill:#fff4e1
```

---

## Core Architectural Patterns

### 1. Interface-First Design Pattern

**Concept:** Define interfaces where they're **used**, not where they're **implemented**. This enables dependency inversion and easy swapping of implementations.

**Current PromptStack Structure:**
```
internal/
├── ai/
│   └── client.go          # Direct Claude API calls
├── library/
│   └── loader.go          # Loads prompts from filesystem
└── history/
    └── database.go        # SQLite operations
```

**Evolution-Ready Pattern:**

```go
// internal/ai/provider.go - Define interface where used
package ai

import "context"

type AIProvider interface {
    GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error)
    ApplyChanges(ctx context.Context, changes []Change) (string, error)
    EstimateTokens(ctx context.Context, text string) (int, error)
}

// internal/ai/claude.go - Current implementation
package ai

type ClaudeProvider struct {
    apiKey string
    model  string
    client *anthropic.Client
}

func NewClaudeProvider(apiKey, model string) (*ClaudeProvider, error) {
    client := anthropic.NewClient(apiKey)
    return &ClaudeProvider{
        apiKey: apiKey,
        model:  model,
        client: client,
    }, nil
}

func (c *ClaudeProvider) GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error) {
    // Direct Claude API implementation
    resp, err := c.client.CreateMessage(ctx, anthropic.MessageCreateParams{
        Model:     anthropic.F(c.model),
        MaxTokens: anthropic.F(4096),
        Messages:  []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewText(prompt))},
    })
    if err != nil {
        return nil, err
    }
    
    // Parse response into suggestions
    return parseSuggestions(resp.Content), nil
}

// Future: internal/ai/mcp.go - MCP-based implementation
package ai

type MCPProvider struct {
    orchestrator *MCPOrchestrator
}

func NewMCPProvider(host string) (*MCPProvider, error) {
    orchestrator, err := NewMCPOrchestrator(host)
    if err != nil {
        return nil, err
    }
    return &MCPProvider{orchestrator: orchestrator}, nil
}

func (m *MCPProvider) GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error) {
    // Route to specialist servers
    request := QueryRequest{
        Query:   prompt,
        Context: m.orchestrator.BuildContext(prompt),
    }
    
    responses, err := m.orchestrator.QuerySpecialists(ctx, request)
    if err != nil {
        return nil, err
    }
    
    return aggregateResponses(responses), nil
}
```

**Usage in Application:**

```go
// cmd/promptstack/main.go
func main() {
    cfg := config.Load()
    
    // Current: Claude provider
    var provider ai.AIProvider
    provider, err = ai.NewClaudeProvider(cfg.ClaudeAPIKey, cfg.Model)
    
    // Future: Switch to MCP provider
    // provider, err = ai.NewMCPProvider(cfg.MCPHost)
    
    if err != nil {
        log.Fatal(err)
    }
    
    app := NewApp(cfg, provider)
    app.Run()
}
```

**Benefits:**
- Easy to swap implementations for testing
- Future-proof for MCP integration
- Clear separation of concerns
- Dependency injection enables mocking

---

### 2. Repository Pattern for Data Layer

**Concept:** Abstract storage operations behind an interface, making it trivial to swap SQLite → PostgreSQL → Neo4j.

**Current Approach (from milestones):**
```go
// internal/history/database.go - Direct SQLite
func (db *Database) SaveComposition(comp Composition) error {
    // SQLite-specific code
}
```

**Evolution-Ready Pattern:**

```go
// internal/storage/repository.go - Interface
package storage

import "context"

type CompositionRepository interface {
    Save(ctx context.Context, comp Composition) error
    Load(ctx context.Context, id string) (Composition, error)
    Search(ctx context.Context, query string) ([]Composition, error)
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, opts ListOptions) ([]Composition, error)
}

type ListOptions struct {
    Limit      int
    Offset     int
    SortBy     string
    SortOrder  string
    WorkingDir string
}

// internal/storage/sqlite.go - Current implementation
package storage

import (
    "database/sql"
    "context"
    _ "modernc.org/sqlite"
)

type SQLiteRepository struct {
    db *sql.DB
}

func NewSQLiteRepository(dataDir string) (*SQLiteRepository, error) {
    dbPath := filepath.Join(dataDir, "history.db")
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }
    
    // Initialize schema
    if err := initSchema(db); err != nil {
        return nil, err
    }
    
    return &SQLiteRepository{db: db}, nil
}

func (r *SQLiteRepository) Save(ctx context.Context, comp Composition) error {
    query := `
        INSERT INTO compositions (file_path, created_at, working_directory, content, character_count, line_count, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(file_path) DO UPDATE SET
            content = excluded.content,
            character_count = excluded.character_count,
            line_count = excluded.line_count,
            updated_at = excluded.updated_at
    `
    _, err := r.db.ExecContext(ctx, query,
        comp.FilePath,
        comp.CreatedAt,
        comp.WorkingDirectory,
        comp.Content,
        comp.CharacterCount,
        comp.LineCount,
        comp.UpdatedAt,
    )
    return err
}

func (r *SQLiteRepository) Search(ctx context.Context, query string) ([]Composition, error) {
    // FTS5 full-text search
    sqlQuery := `
        SELECT id, file_path, created_at, working_directory, content, character_count, line_count, updated_at
        FROM compositions
        WHERE content MATCH ?
        ORDER BY rank
        LIMIT 100
    `
    rows, err := r.db.QueryContext(ctx, sqlQuery, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var results []Composition
    for rows.Next() {
        var comp Composition
        if err := rows.Scan(&comp.ID, &comp.FilePath, &comp.CreatedAt, &comp.WorkingDirectory,
            &comp.Content, &comp.CharacterCount, &comp.LineCount, &comp.UpdatedAt); err != nil {
            return nil, err
        }
        results = append(results, comp)
    }
    
    return results, nil
}

// Future: internal/storage/postgres.go
package storage

import (
    "database/sql"
    "context"
    _ "github.com/lib/pq"
)

type PostgresRepository struct {
    db *sql.DB
}

func NewPostgresRepository(connString string) (*PostgresRepository, error) {
    db, err := sql.Open("postgres", connString)
    if err != nil {
        return nil, err
    }
    
    // Initialize schema with pgvector for semantic search
    if err := initSchema(db); err != nil {
        return nil, err
    }
    
    return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Save(ctx context.Context, comp Composition) error {
    // PostgreSQL implementation with vector embeddings
    query := `
        INSERT INTO compositions (file_path, created_at, working_directory, content, character_count, line_count, updated_at, embedding)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (file_path) DO UPDATE SET
            content = EXCLUDED.content,
            character_count = EXCLUDED.character_count,
            line_count = EXCLUDED.line_count,
            updated_at = EXCLUDED.updated_at,
            embedding = EXCLUDED.embedding
    `
    
    // Generate embedding for semantic search
    embedding, err := generateEmbedding(comp.Content)
    if err != nil {
        return err
    }
    
    _, err = r.db.ExecContext(ctx, query,
        comp.FilePath,
        comp.CreatedAt,
        comp.WorkingDirectory,
        comp.Content,
        comp.CharacterCount,
        comp.LineCount,
        comp.UpdatedAt,
        embedding,
    )
    return err
}

// Future: internal/storage/graph.go - For Neo4j
package storage

import (
    "context"
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphRepository struct {
    driver neo4j.Driver
}

func NewGraphRepository(uri, username, password string) (*GraphRepository, error) {
    driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
    if err != nil {
        return nil, err
    }
    
    return &GraphRepository{driver: driver}, nil
}

func (g *GraphRepository) Save(ctx context.Context, comp Composition) error {
    session := g.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
    defer session.Close(ctx)
    
    _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        // Create composition node with relationships
        query := `
            MERGE (c:Composition {id: $id})
            SET c.file_path = $file_path,
                c.content = $content,
                c.created_at = $created_at,
                c.working_directory = $working_directory
            
            // Create relationships to referenced files
            UNWIND $file_refs AS ref
            MERGE (f:File {path: ref.path})
            MERGE (c)-[:REFERENCES]->(f)
            
            // Create relationships to used prompts
            UNWIND $prompt_refs AS ref
            MERGE (p:Prompt {id: ref.id})
            MERGE (c)-[:USES]->(p)
        `
        
        params := map[string]interface{}{
            "id":              comp.ID,
            "file_path":       comp.FilePath,
            "content":         comp.Content,
            "created_at":      comp.CreatedAt,
            "working_directory": comp.WorkingDirectory,
            "file_refs":       extractFileReferences(comp.Content),
            "prompt_refs":     extractPromptReferences(comp.Content),
        }
        
        return tx.Run(ctx, query, params)
    })
    
    return err
}

func (g *GraphRepository) QueryRelated(ctx context.Context, compID string) ([]Composition, error) {
    // Find compositions that share files or prompts
    session := g.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
    defer session.Close(ctx)
    
    result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        query := `
            MATCH (c:Composition {id: $id})-[:REFERENCES|USES]->(x)<-[:REFERENCES|USES]-(other:Composition)
            WHERE other.id <> $id
            RETURN other
            ORDER BY COUNT(x) DESC
            LIMIT 10
        `
        
        result, err := tx.Run(ctx, query, map[string]interface{}{"id": compID})
        if err != nil {
            return nil, err
        }
        
        var related []Composition
        for result.Next(ctx) {
            var comp Composition
            if err := result.Scan(&comp); err != nil {
                return nil, err
            }
            related = append(related, comp)
        }
        
        return related, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    return result.([]Composition), nil
}
```

**Factory Pattern for Configuration:**

```go
// internal/storage/factory.go
package storage

import "github.com/yourorg/promptstack/internal/config"

func NewRepository(cfg *config.Config) (CompositionRepository, error) {
    switch cfg.Storage {
    case "sqlite", "":
        return NewSQLiteRepository(cfg.DataDir)
    case "postgres":
        return NewPostgresRepository(cfg.PostgresURL)
    case "graph":
        return NewGraphRepository(cfg.Neo4jURL)
    default:
        return nil, fmt.Errorf("unknown storage type: %s", cfg.Storage)
    }
}
```

**Benefits:**
- Easy migration between storage backends
- Testable with in-memory implementations
- Clear separation of storage logic
- Future-proof for PostgreSQL + Neo4j

---

### 3. Strategy Pattern for Context Selection

**Concept:** Make context selection algorithm pluggable, enabling evolution from simple keyword matching to sophisticated specialist routing.

**Current Approach (from milestones):**
```go
// internal/ai/context.go - Hardcoded algorithm
func SelectContext(composition string, library []Prompt) []Prompt {
    // Fixed scoring algorithm
}
```

**Evolution-Ready Pattern:**

```go
// internal/ai/selector.go - Interface
package ai

import "context"

type ContextSelector interface {
    Select(ctx context.Context, composition string, library []Prompt) ([]Prompt, error)
    Name() string
}

// internal/ai/keyword_selector.go - Current implementation
package ai

type KeywordSelector struct {
    scorer *Scorer
}

func NewKeywordSelector() *KeywordSelector {
    return &KeywordSelector{
        scorer: NewScorer(),
    }
}

func (k *KeywordSelector) Select(ctx context.Context, composition string, library []Prompt) ([]Prompt, error) {
    // Extract keywords from composition
    keywords := extractKeywords(composition)
    
    // Score each prompt
    scored := make([]ScoredPrompt, 0, len(library))
    for _, prompt := range library {
        score := k.scorer.Score(keywords, prompt)
        scored = append(scored, ScoredPrompt{
            Prompt: prompt,
            Score:  score,
        })
    }
    
    // Sort by score (highest first)
    sort.Slice(scored, func(i, j int) bool {
        return scored[i].Score > scored[j].Score
    })
    
    // Select top prompts within token budget
    selected := make([]Prompt, 0)
    totalTokens := 0
    for _, sp := range scored {
        tokens := estimateTokens(sp.Prompt.Content)
        if totalTokens+tokens > maxContextTokens {
            break
        }
        selected = append(selected, sp.Prompt)
        totalTokens += tokens
    }
    
    return selected, nil
}

func (k *KeywordSelector) Name() string {
    return "keyword"
}

// Future: internal/ai/specialist_selector.go
package ai

type SpecialistSelector struct {
    specialists map[string]SpecialistServer
    router      *ContextRouter
}

func NewSpecialistSelector(specialists map[string]SpecialistServer) *SpecialistSelector {
    return &SpecialistSelector{
        specialists: specialists,
        router:      NewContextRouter(),
    }
}

func (s *SpecialistSelector) Select(ctx context.Context, composition string, library []Prompt) ([]Prompt, error) {
    // Analyze composition to determine relevant specialists
    relevantSpecialists := s.router.Route(ctx, composition)
    
    // Query each specialist for relevant prompts
    var selected []Prompt
    for _, specialist := range relevantSpecialists {
        request := QueryRequest{
            Query:   composition,
            Library: library,
        }
        
        response, err := specialist.Query(ctx, request)
        if err != nil {
            log.Printf("Specialist %s query failed: %v", specialist.Name(), err)
            continue
        }
        
        selected = append(selected, response.RelevantPrompts...)
    }
    
    // Deduplicate and rank
    return deduplicateAndRank(selected), nil
}

func (s *SpecialistSelector) Name() string {
    return "specialist"
}

// internal/ai/scorer.go
package ai

type Scorer struct {
    // Scoring weights
    tagMatchWeight      float64
    categoryMatchWeight float64
    keywordWeight       float64
    recentWeight        float64
    frequencyWeight     float64
}

func NewScorer() *Scorer {
    return &Scorer{
        tagMatchWeight:      10.0,
        categoryMatchWeight: 5.0,
        keywordWeight:       1.0,
        recentWeight:        3.0,
        frequencyWeight:     1.0,
    }
}

func (s *Scorer) Score(keywords []string, prompt Prompt) float64 {
    var score float64
    
    // Tag match: +10 per matching tag
    for _, tag := range prompt.Tags {
        for _, keyword := range keywords {
            if strings.EqualFold(tag, keyword) {
                score += s.tagMatchWeight
            }
        }
    }
    
    // Category match: +5 if same category as keywords
    if s.categoryMatches(keywords, prompt.Category) {
        score += s.categoryMatchWeight
    }
    
    // Keyword overlap: +1 per matching word
    for _, keyword := range keywords {
        if strings.Contains(strings.ToLower(prompt.Content), strings.ToLower(keyword)) {
            score += s.keywordWeight
        }
    }
    
    // Recently used: +3 if used in last session
    if s.isRecentlyUsed(prompt) {
        score += s.recentWeight
    }
    
    // Frequently used: +use_count
    score += float64(prompt.UseCount) * s.frequencyWeight
    
    return score
}
```

**Benefits:**
- Easy to experiment with different selection algorithms
- Future-proof for specialist-based routing
- Testable with mock selectors
- Clear separation of selection logic

---

### 4. Plugin Architecture with Go Plugins

**Concept:** Enable loading specialist servers at runtime for true extensibility.

**Folder Structure:**
```
internal/
├── specialist/
│   ├── registry.go          # Plugin registry
│   ├── interface.go         # Specialist interface
│   └── builtin/             # Built-in specialists
│       ├── workflow.go
│       ├── code.go
│       └── compliance.go
└── plugins/                # External plugins (future)
    └── tax_specialist.so
```

**Interface Definition:**

```go
// internal/specialist/interface.go
package specialist

import "context"

type Specialist interface {
    Name() string
    Category() string
    Description() string
    Query(ctx context.Context, request QueryRequest) (*QueryResponse, error)
    Capabilities() []Capability
}

type QueryRequest struct {
    Query   string
    Context map[string]interface{}
    Library []Prompt
}

type QueryResponse struct {
    RelevantPrompts []Prompt
    Suggestions     []Suggestion
    Confidence     float64
    Metadata       map[string]interface{}
}

type Capability struct {
    Name        string
    Description string
    InputType   string
    OutputType  string
}

// internal/specialist/registry.go
package specialist

import (
    "sync"
)

type Registry struct {
    mu          sync.RWMutex
    specialists map[string]Specialist
}

func NewRegistry() *Registry {
    return &Registry{
        specialists: make(map[string]Specialist),
    }
}

func (r *Registry) Register(s Specialist) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.specialists[s.Name()] = s
}

func (r *Registry) Get(name string) (Specialist, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    s, ok := r.specialists[name]
    return s, ok
}

func (r *Registry) List() []Specialist {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    list := make([]Specialist, 0, len(r.specialists))
    for _, s := range r.specialists {
        list = append(list, s)
    }
    return list
}

func (r *Registry) GetByCategory(category string) []Specialist {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var filtered []Specialist
    for _, s := range r.specialists {
        if s.Category() == category {
            filtered = append(filtered, s)
        }
    }
    return filtered
}

// internal/specialist/builtin/workflow.go
package builtin

import (
    "context"
    "github.com/yourorg/promptstack/internal/specialist"
)

type WorkflowSpecialist struct {
    workflows []Prompt
}

func NewWorkflowSpecialist(workflows []Prompt) *WorkflowSpecialist {
    return &WorkflowSpecialist{workflows: workflows}
}

func (w *WorkflowSpecialist) Name() string {
    return "workflow"
}

func (w *WorkflowSpecialist) Category() string {
    return "workflows"
}

func (w *WorkflowSpecialist) Description() string {
    return "Specialist for multi-step workflow prompts"
}

func (w *WorkflowSpecialist) Query(ctx context.Context, request specialist.QueryRequest) (*specialist.QueryResponse, error) {
    // Analyze query to find relevant workflows
    relevant := w.findRelevantWorkflows(request.Query)
    
    return &specialist.QueryResponse{
        RelevantPrompts: relevant,
        Confidence:     0.85,
        Metadata: map[string]interface{}{
            "matched_count": len(relevant),
        },
    }, nil
}

func (w *WorkflowSpecialist) Capabilities() []specialist.Capability {
    return []specialist.Capability{
        {
            Name:        "workflow_recommendation",
            Description: "Suggest relevant workflows based on context",
            InputType:   "text",
            OutputType:  "prompt_list",
        },
    }
}

func (w *WorkflowSpecialist) findRelevantWorkflows(query string) []Prompt {
    // Simple keyword matching for now
    // Future: Use semantic search with embeddings
    var relevant []Prompt
    for _, workflow := range w.workflows {
        if strings.Contains(strings.ToLower(workflow.Content), strings.ToLower(query)) {
            relevant = append(relevant, workflow)
        }
    }
    return relevant
}
```

**Future Plugin Loading:**

```go
// internal/specialist/plugin.go
package specialist

import (
    "fmt"
    "plugin"
)

// LoadPlugin loads a specialist from a compiled Go plugin
func LoadPlugin(path string) (Specialist, error) {
    p, err := plugin.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open plugin: %w", err)
    }
    
    // Look for the NewSpecialist symbol
    sym, err := p.Lookup("NewSpecialist")
    if err != nil {
        return nil, fmt.Errorf("plugin missing NewSpecialist symbol: %w", err)
    }
    
    // Type assert to function
    newFunc, ok := sym.(func() Specialist)
    if !ok {
        return nil, fmt.Errorf("plugin has unexpected type for NewSpecialist")
    }
    
    // Create specialist instance
    specialist := newFunc()
    return specialist, nil
}

// Example plugin implementation (in separate package)
/*
package main

import "github.com/yourorg/promptstack/internal/specialist"

type TaxSpecialist struct {
    // Specialist implementation
}

func NewSpecialist() specialist.Specialist {
    return &TaxSpecialist{}
}

func (t *TaxSpecialist) Name() string {
    return "tax"
}

func (t *TaxSpecialist) Category() string {
    return "business"
}

// ... implement other interface methods

func main() {
    // Plugin entry point (required for Go plugins)
}
*/
```

**Benefits:**
- Extensible without modifying core code
- Third-party specialists can be added
- Clear plugin contract via interface
- Built-in specialists always available

---

### 5. Configuration-Driven Behavior

**Concept:** Use configuration to enable/disable features and switch implementations without code changes.

**Config Structure:**

```go
// internal/config/config.go
package config

import "time"

type Config struct {
    // Current fields
    ClaudeAPIKey string `yaml:"claude_api_key"`
    Model        string `yaml:"model"`
    VimMode      bool   `yaml:"vim_mode"`
    DataDir      string `yaml:"data_dir"`
    
    // Future fields (add when needed)
    AIProvider   string   `yaml:"ai_provider"`   // "claude" | "mcp"
    Storage      string   `yaml:"storage"`       // "sqlite" | "postgres" | "graph"
    Specialists  []string `yaml:"specialists"`   // Enabled specialists
    MCPHost      string   `yaml:"mcp_host"`      // MCP orchestrator address
    PostgresURL  string   `yaml:"postgres_url"`   // PostgreSQL connection string
    Neo4jURL     string   `yaml:"neo4j_url"`     // Neo4j connection string
    EnablePlugins bool     `yaml:"enable_plugins"` // Enable external plugins
    PluginDir    string   `yaml:"plugin_dir"`    // Directory for plugins
    
    // Performance tuning
    MaxContextTokens int           `yaml:"max_context_tokens"`
    CacheEnabled    bool          `yaml:"cache_enabled"`
    CacheTTL       time.Duration `yaml:"cache_ttl"`
    
    // Logging
    LogLevel string `yaml:"log_level"` // "debug" | "info" | "warn" | "error"
    LogFile  string `yaml:"log_file"`
}
```

**Factory Pattern for Configuration:**

```go
// internal/ai/factory.go
package ai

import (
    "fmt"
    "github.com/yourorg/promptstack/internal/config"
)

func NewProvider(cfg *config.Config) (AIProvider, error) {
    switch cfg.AIProvider {
    case "claude", "":
        return NewClaudeProvider(cfg.ClaudeAPIKey, cfg.Model)
    case "mcp":
        return NewMCPProvider(cfg.MCPHost)
    default:
        return nil, fmt.Errorf("unknown AI provider: %s", cfg.AIProvider)
    }
}

// internal/storage/factory.go
package storage

import (
    "fmt"
    "github.com/yourorg/promptstack/internal/config"
)

func NewRepository(cfg *config.Config) (CompositionRepository, error) {
    switch cfg.Storage {
    case "sqlite", "":
        return NewSQLiteRepository(cfg.DataDir)
    case "postgres":
        return NewPostgresRepository(cfg.PostgresURL)
    case "graph":
        return NewGraphRepository(cfg.Neo4jURL)
    default:
        return nil, fmt.Errorf("unknown storage type: %s", cfg.Storage)
    }
}

// internal/ai/selector_factory.go
package ai

import (
    "fmt"
    "github.com/yourorg/promptstack/internal/config"
)

func NewSelector(cfg *config.Config, specialists map[string]SpecialistServer) (ContextSelector, error) {
    switch cfg.ContextSelector {
    case "keyword", "":
        return NewKeywordSelector(), nil
    case "specialist":
        return NewSpecialistSelector(specialists), nil
    default:
        return nil, fmt.Errorf("unknown context selector: %s", cfg.ContextSelector)
    }
}
```

**Example Config File:**

```yaml
# ~/.promptstack/config.yaml

# Current settings
claude_api_key: "sk-ant-..."
model: "claude-3-sonnet-20240229"
vim_mode: true
data_dir: "~/.promptstack/data"

# Future settings (add when needed)
ai_provider: "claude"  # or "mcp"
storage: "sqlite"        # or "postgres", "graph"
context_selector: "keyword"  # or "specialist"

# MCP settings (future)
mcp_host: "localhost:8080"
specialists:
  - "workflow"
  - "code"
  - "compliance"

# PostgreSQL settings (future)
postgres_url: "postgres://user:pass@localhost/promptstack"

# Neo4j settings (future)
neo4j_url: "bolt://localhost:7687"

# Plugin settings (future)
enable_plugins: true
plugin_dir: "~/.promptstack/plugins"

# Performance tuning
max_context_tokens: 100000
cache_enabled: true
cache_ttl: 5m

# Logging
log_level: "info"
log_file: "~/.promptstack/debug.log"
```

**Benefits:**
- Switch implementations without code changes
- Easy A/B testing of different algorithms
- Clear configuration documentation
- Future-proof for new features

---

### 6. Middleware Pattern for Cross-Cutting Concerns

**Concept:** Compose behavior through middleware layers for logging, caching, metrics, etc.

```go
// internal/ai/middleware.go
package ai

import (
    "context"
    "time"
    "go.uber.org/zap"
)

type ProviderMiddleware func(AIProvider) AIProvider

// WithLogging adds logging to an AIProvider
func WithLogging(logger *zap.Logger) ProviderMiddleware {
    return func(next AIProvider) AIProvider {
        return &loggingProvider{
            next:   next,
            logger: logger,
        }
    }
}

type loggingProvider struct {
    next   AIProvider
    logger *zap.Logger
}

func (l *loggingProvider) GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error) {
    start := time.Now()
    l.logger.Info("Getting suggestions",
        zap.String("prompt", prompt[:min(100, len(prompt))]),
    )
    
    suggestions, err := l.next.GetSuggestions(ctx, prompt)
    
    duration := time.Since(start)
    if err != nil {
        l.logger.Error("Failed to get suggestions",
            zap.Error(err),
            zap.Duration("duration", duration),
        )
    } else {
        l.logger.Info("Got suggestions",
            zap.Int("count", len(suggestions)),
            zap.Duration("duration", duration),
        )
    }
    
    return suggestions, err
}

func (l *loggingProvider) ApplyChanges(ctx context.Context, changes []Change) (string, error) {
    start := time.Now()
    l.logger.Info("Applying changes",
        zap.Int("change_count", len(changes)),
    )
    
    result, err := l.next.ApplyChanges(ctx, changes)
    
    duration := time.Since(start)
    if err != nil {
        l.logger.Error("Failed to apply changes",
            zap.Error(err),
            zap.Duration("duration", duration),
        )
    } else {
        l.logger.Info("Applied changes",
            zap.Duration("duration", duration),
        )
    }
    
    return result, err
}

// WithCaching adds caching to an AIProvider
func WithCaching(cache Cache) ProviderMiddleware {
    return func(next AIProvider) AIProvider {
        return &cachingProvider{
            next:  next,
            cache: cache,
        }
    }
}

type cachingProvider struct {
    next  AIProvider
    cache Cache
}

func (c *cachingProvider) GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("suggestions:%s", hash(prompt))
    if cached, ok := c.cache.Get(cacheKey); ok {
        return cached.([]Suggestion), nil
    }
    
    // Cache miss - call next provider
    suggestions, err := c.next.GetSuggestions(ctx, prompt)
    if err != nil {
        return nil, err
    }
    
    // Cache the result
    c.cache.Set(cacheKey, suggestions, 5*time.Minute)
    
    return suggestions, nil
}

// WithMetrics adds metrics collection to an AIProvider
func WithMetrics(metrics Metrics) ProviderMiddleware {
    return func(next AIProvider) AIProvider {
        return &metricsProvider{
            next:    next,
            metrics: metrics,
        }
    }
}

type metricsProvider struct {
    next    AIProvider
    metrics Metrics
}

func (m *metricsProvider) GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error) {
    start := time.Now()
    
    suggestions, err := m.next.GetSuggestions(ctx, prompt)
    
    duration := time.Since(start)
    m.metrics.RecordDuration("ai.suggestions.duration", duration)
    m.metrics.RecordCount("ai.suggestions.count", len(suggestions))
    
    if err != nil {
        m.metrics.RecordError("ai.suggestions.errors")
    }
    
    return suggestions, err
}

// Cache interface
type Cache interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration)
    Delete(key string)
}

// Metrics interface
type Metrics interface {
    RecordDuration(name string, duration time.Duration)
    RecordCount(name string, count int)
    RecordError(name string)
}
```

**Usage in main.go:**

```go
func main() {
    cfg := config.Load()
    logger := zap.NewExample()
    
    // Create base provider
    provider, err := ai.NewClaudeProvider(cfg.ClaudeAPIKey, cfg.Model)
    if err != nil {
        log.Fatal(err)
    }
    
    // Apply middleware (order matters!)
    if cfg.CacheEnabled {
        cache := NewMemoryCache()
        provider = ai.WithCaching(cache)(provider)
    }
    
    provider = ai.WithLogging(logger)(provider)
    provider = ai.WithMetrics(metrics)(provider)
    
    app := NewApp(cfg, provider)
    app.Run()
}
```

**Benefits:**
- Composable behavior
- Easy to add/remove cross-cutting concerns
- Clear separation of concerns
- Testable middleware layers

---

### 7. Domain Events for Decoupling

**Concept:** Use domain events to enable future features like audit logging, analytics, or real-time notifications without tight coupling.

```go
// internal/events/events.go
package events

import "time"

type Event interface {
    Type() string
    Timestamp() time.Time
    Payload() map[string]interface{}
}

type CompositionSavedEvent struct {
    ID        string
    Content   string
    Timestamp time.Time
}

func (e *CompositionSavedEvent) Type() string {
    return "composition.saved"
}

func (e *CompositionSavedEvent) Timestamp() time.Time {
    return e.Timestamp
}

func (e *CompositionSavedEvent) Payload() map[string]interface{} {
    return map[string]interface{}{
        "id":      e.ID,
        "content": e.Content,
        "length":  len(e.Content),
    }
}

type PromptUsedEvent struct {
    PromptID   string
    PromptName string
    Context    string
    Timestamp  time.Time
}

func (e *PromptUsedEvent) Type() string {
    return "prompt.used"
}

func (e *PromptUsedEvent) Timestamp() time.Time {
    return e.Timestamp
}

func (e *PromptUsedEvent) Payload() map[string]interface{} {
    return map[string]interface{}{
        "prompt_id":   e.PromptID,
        "prompt_name": e.PromptName,
        "context":     e.Context,
    }
}

type SuggestionAcceptedEvent struct {
    SuggestionID string
    Suggestion  Suggestion
    Timestamp   time.Time
}

func (e *SuggestionAcceptedEvent) Type() string {
    return "suggestion.accepted"
}

func (e *SuggestionAcceptedEvent) Timestamp() time.Time {
    return e.Timestamp
}

func (e *SuggestionAcceptedEvent) Payload() map[string]interface{} {
    return map[string]interface{}{
        "suggestion_id": e.SuggestionID,
        "suggestion":    e.Suggestion,
    }
}

// internal/events/dispatcher.go
package events

import (
    "sync"
)

type EventHandler func(Event) error

type Dispatcher struct {
    mu       sync.RWMutex
    handlers map[string][]EventHandler
}

func NewDispatcher() *Dispatcher {
    return &Dispatcher{
        handlers: make(map[string][]EventHandler),
    }
}

func (d *Dispatcher) Subscribe(eventType string, handler EventHandler) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.handlers[eventType] = append(d.handlers[eventType], handler)
}

func (d *Dispatcher) Publish(event Event) {
    d.mu.RLock()
    handlers := d.handlers[event.Type()]
    d.mu.RUnlock()
    
    // Dispatch to all handlers asynchronously
    for _, handler := range handlers {
        go func(h EventHandler) {
            if err := h(event); err != nil {
                log.Printf("Event handler error: %v", err)
            }
        }(handler)
    }
}

// Example: Analytics handler
type AnalyticsHandler struct {
    client *AnalyticsClient
}

func (a *AnalyticsHandler) Handle(event Event) error {
    switch e := event.(type) {
    case *CompositionSavedEvent:
        return a.client.Track("composition_saved", e.Payload())
    case *PromptUsedEvent:
        return a.client.Track("prompt_used", e.Payload())
    case *SuggestionAcceptedEvent:
        return a.client.Track("suggestion_accepted", e.Payload())
    default:
        return nil
    }
}

// Example: Audit log handler
type AuditLogHandler struct {
    logger *zap.Logger
}

func (a *AuditLogHandler) Handle(event Event) error {
    a.logger.Info("Audit event",
        zap.String("type", event.Type()),
        zap.Time("timestamp", event.Timestamp()),
        zap.Any("payload", event.Payload()),
    )
    return nil
}

// Example: Notification handler
type NotificationHandler struct {
    notifier Notifier
}

func (n *NotificationHandler) Handle(event Event) error {
    switch e := event.(type) {
    case *SuggestionAcceptedEvent:
        return n.notifier.Notify("Suggestion accepted", e.Suggestion.Title)
    default:
        return nil
    }
}
```

**Usage in Application:**

```go
func main() {
    cfg := config.Load()
    
    // Create event dispatcher
    dispatcher := events.NewDispatcher()
    
    // Subscribe handlers
    if cfg.AnalyticsEnabled {
        analytics := NewAnalyticsHandler(cfg.AnalyticsAPIKey)
        dispatcher.Subscribe("composition.saved", analytics.Handle)
        dispatcher.Subscribe("prompt.used", analytics.Handle)
        dispatcher.Subscribe("suggestion.accepted", analytics.Handle)
    }
    
    if cfg.AuditLoggingEnabled {
        audit := NewAuditLogHandler(logger)
        dispatcher.Subscribe("*", audit.Handle) // Subscribe to all events
    }
    
    if cfg.NotificationsEnabled {
        notifier := NewNotifier(cfg.NotificationConfig)
        dispatcher.Subscribe("suggestion.accepted", notifier.Handle)
    }
    
    // Pass dispatcher to application
    app := NewApp(cfg, dispatcher)
    app.Run()
}

// In composition save logic
func (c *Composition) Save() error {
    // Save to storage
    if err := c.repo.Save(c); err != nil {
        return err
    }
    
    // Publish event
    event := &events.CompositionSavedEvent{
        ID:        c.ID,
        Content:   c.Content,
        Timestamp: time.Now(),
    }
    c.dispatcher.Publish(event)
    
    return nil
}
```

**Benefits:**
- Decoupled architecture
- Easy to add new event handlers
- Future-proof for analytics, audit logging, notifications
- Clear event flow

---

## Folder Structure Evolution

The current structure in [`project-structure.md`](../plans/fresh-build/project-structure.md) is well-designed and can evolve gracefully:

### Current Structure
```
promptstack/
├── cmd/promptstack/
│   └── main.go
│
├── internal/
│   ├── ai/
│   │   ├── client.go
│   │   ├── context.go
│   │   ├── tokens.go
│   │   ├── suggestions.go
│   │   └── diff.go
│   ├── library/
│   │   ├── library.go
│   │   ├── loader.go
│   │   ├── index.go
│   │   ├── scorer.go
│   │   └── search.go
│   ├── history/
│   │   ├── manager.go
│   │   ├── database.go
│   │   ├── storage.go
│   │   ├── sync.go
│   │   ├── search.go
│   │   └── cleanup.go
│   ├── editor/
│   ├── prompt/
│   ├── config/
│   ├── platform/
│   ├── vim/
│   └── commands/
│
├── ui/
└── test/
```

### Evolution-Ready Structure
```
promptstack/
├── cmd/promptstack/
│   └── main.go
│
├── internal/
│   ├── ai/
│   │   ├── provider.go           # NEW: Interface (add now)
│   │   ├── claude.go             # Current implementation
│   │   ├── context.go            # Context selection
│   │   ├── selector.go           # NEW: Selector interface (add now)
│   │   ├── keyword_selector.go   # Current implementation
│   │   ├── tokens.go
│   │   ├── suggestions.go
│   │   ├── diff.go
│   │   ├── middleware.go         # NEW: Middleware support (add now)
│   │   └── factory.go           # NEW: Factory for providers (add now)
│   │
│   ├── specialist/               # NEW: For future specialists
│   │   ├── interface.go          # Specialist interface
│   │   ├── registry.go           # Plugin registry
│   │   ├── plugin.go             # Plugin loading
│   │   └── builtin/              # Built-in specialists
│   │       ├── workflow.go
│   │       ├── code.go
│   │       └── compliance.go
│   │
│   ├── storage/                  # RENAME from history/
│   │   ├── repository.go         # NEW: Repository interface (add now)
│   │   ├── sqlite.go             # Current implementation
│   │   ├── postgres.go           # Future implementation
│   │   ├── graph.go              # Future implementation
│   │   ├── factory.go            # NEW: Factory for repositories (add now)
│   │   ├── manager.go            # Composition manager
│   │   ├── sync.go              # Sync verification
│   │   ├── search.go            # Search operations
│   │   └── cleanup.go           # Cleanup operations
│   │
│   ├── library/
│   │   ├── source.go             # NEW: Source interface (add now)
│   │   ├── filesystem.go         # Current implementation
│   │   ├── specialist.go         # Future implementation
│   │   ├── library.go
│   │   ├── loader.go
│   │   ├── index.go
│   │   ├── scorer.go
│   │   └── search.go
│   │
│   ├── events/                   # NEW: Domain events
│   │   ├── events.go
│   │   └── dispatcher.go
│   │
│   ├── editor/                   # Unchanged
│   ├── prompt/                   # Unchanged
│   ├── config/                   # Unchanged
│   ├── platform/                 # Unchanged
│   ├── vim/                      # Unchanged
│   └── commands/                 # Unchanged
│
├── ui/                           # Unchanged
└── test/                         # Unchanged
```

### Migration Strategy

**Phase 1: Add Abstraction Layers (Now)**
- Add [`internal/ai/provider.go`](internal/ai/provider.go) - AI provider interface
- Add [`internal/ai/selector.go`](internal/ai/selector.go) - Context selector interface
- Add [`internal/storage/repository.go`](internal/storage/repository.go) - Repository interface
- Add [`internal/library/source.go`](internal/library/source.go) - Source interface
- Add [`internal/events/`](internal/events/) - Domain events system

**Phase 2: Implement Current Logic with Interfaces (Now)**
- Refactor [`internal/ai/client.go`](internal/ai/client.go) → [`internal/ai/claude.go`](internal/ai/claude.go)
- Refactor [`internal/history/database.go`](internal/history/database.go) → [`internal/storage/sqlite.go`](internal/storage/sqlite.go)
- Refactor [`internal/library/loader.go`](internal/library/loader.go) → [`internal/library/filesystem.go`](internal/library/filesystem.go)

**Phase 3: Add Factory Patterns (When Needed)**
- Add [`internal/ai/factory.go`](internal/ai/factory.go) - Provider factory
- Add [`internal/storage/factory.go`](internal/storage/factory.go) - Repository factory
- Add [`internal/ai/selector_factory.go`](internal/ai/selector_factory.go) - Selector factory

**Phase 4: Add Middleware (When Needed)**
- Add [`internal/ai/middleware.go`](internal/ai/middleware.go) - Middleware support
- Implement logging middleware
- Implement caching middleware (when performance is an issue)
- Implement metrics middleware (when monitoring is needed)

**Phase 5: Add Specialist Architecture (When Needed)**
- Add [`internal/specialist/`](internal/specialist/) package
- Implement built-in specialists
- Add plugin loading support
- Implement MCP integration

**Phase 6: Migrate Storage (When Needed)**
- Add [`internal/storage/postgres.go`](internal/storage/postgres.go) - PostgreSQL implementation
- Add [`internal/storage/graph.go`](internal/storage/graph.go) - Neo4j implementation
- Update factory to support new storage types
- Provide migration tools

---

## Implementation Guidelines

### When to Add Interfaces

**Add interfaces when:**
- You need to swap implementations for testing
- You anticipate multiple implementations in the future
- You want to decouple components
- You need to mock dependencies

**Don't add interfaces when:**
- There's only one implementation and no need for testing
- The interface adds complexity without benefit
- The implementation is unlikely to change

### When to Use Patterns

**Use repository pattern when:**
- You have multiple storage backends (SQLite, PostgreSQL, Neo4j)
- You need to test business logic without real storage
- You want to abstract storage operations

**Use strategy pattern when:**
- You have multiple algorithms for the same operation
- You want to switch algorithms at runtime
- You need to test different approaches

**Use middleware pattern when:**
- You have cross-cutting concerns (logging, caching, metrics)
- You want to compose behavior
- You need to add/remove features without modifying core logic

**Use plugin architecture when:**
- You want third-party extensibility
- You have clear plugin contracts
- You need to load code at runtime

### Incremental Implementation

**Start simple:**
1. Implement current functionality without abstractions
2. Add interfaces when you need to swap implementations
3. Add patterns when they solve a real problem
4. Refactor incrementally, not all at once

**Example progression:**

```go
// Step 1: Direct implementation
func GetSuggestions(prompt string) ([]Suggestion, error) {
    // Direct Claude API call
}

// Step 2: Add interface when needed for testing
type AIProvider interface {
    GetSuggestions(prompt string) ([]Suggestion, error)
}

type ClaudeProvider struct {
    apiKey string
}

func (c *ClaudeProvider) GetSuggestions(prompt string) ([]Suggestion, error) {
    // Claude API implementation
}

// Step 3: Add middleware when needed for logging
type LoggingProvider struct {
    next   AIProvider
    logger *zap.Logger
}

func (l *LoggingProvider) GetSuggestions(prompt string) ([]Suggestion, error) {
    l.logger.Info("Getting suggestions", zap.String("prompt", prompt))
    return l.next.GetSuggestions(prompt)
}

// Step 4: Add factory when needed for configuration
func NewProvider(cfg *Config) (AIProvider, error) {
    provider, err := NewClaudeProvider(cfg.APIKey)
    if err != nil {
        return nil, err
    }
    return WithLogging(logger)(provider), nil
}
```

---

## Key Takeaways

1. **Define interfaces where used**, not where implemented
2. **Use repository pattern** for data layer abstraction
3. **Factory pattern** for configuration-driven instantiation
4. **Middleware pattern** for cross-cutting concerns
5. **Domain events** for decoupling future features
6. **Plugin architecture** for extensibility (add when needed)
7. **Keep current structure** - just add abstraction layers
8. **Implement incrementally** - add patterns when they solve real problems
9. **Don't over-engineer** - YAGNI (You Aren't Gonna Need It)
10. **Test-driven approach** - patterns emerge from test needs

The beauty of these patterns is that you can implement them **incrementally**:
- Start with current implementation
- Add interfaces when you need to swap implementations
- Add middleware when you need logging/caching
- Add plugins when you need external specialists

This keeps the codebase maintainable while enabling future extensibility toward Project OS.

---

## Related Documents

- [`../plans/fresh-build/what-we-are-building.md`](../plans/fresh-build/what-we-are-building.md) - PromptStack requirements
- [`../plans/fresh-build/project-structure.md`](../plans/fresh-build/project-structure.md) - Current project structure
- [`../plans/fresh-build/milestones.md`](../plans/fresh-build/milestones.md) - Implementation milestones
- [`/Users/kyledavis/Sites/go-cli-ai/README.md`](/Users/kyledavis/Sites/go-cli-ai/README.md) - Project OS vision

---

**Last Updated:** 2026-01-07
**Status:** Architectural Discussion - Ready for Implementation