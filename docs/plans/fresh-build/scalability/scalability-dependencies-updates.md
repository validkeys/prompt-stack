# Dependencies Updates for Scalability

**Document**: [`DEPENDENCIES.md`](../DEPENDENCIES.md)  
**Priority**: **Medium**  
**Status**: Should update before implementation begins

---

## Overview

This document details all changes required to [`DEPENDENCIES.md`](../DEPENDENCIES.md) to incorporate Phase 1 scalability abstractions. These changes add future dependencies for PostgreSQL, Neo4j, and MCP integration.

---

## Future Dependencies

### PostgreSQL Support

```markdown
**PostgreSQL Support:**
- **Package**: `github.com/lib/pq`
- **Purpose**: PostgreSQL driver for team collaboration features
- **Status**: Future (Phase 4)
- **When to Add**: When implementing PostgreSQL repository backend
- **Version**: Latest stable
- **License**: MIT

**Usage:**
```go
import (
    "database/sql"
    _ "github.com/lib/pq"
)

func NewPostgreSQLRepository(connStr string) (*PostgreSQLRepository, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    return &PostgreSQLRepository{db: db}, nil
}
```

**Configuration:**
```yaml
storage: "postgres"
postgres_url: "postgres://user:password@localhost:5432/promptstack?sslmode=disable"
```
```

---

### Neo4j Support

```markdown
**Neo4j Support:**
- **Package**: `github.com/neo4j/neo4j-go-driver`
- **Purpose**: Neo4j driver for graph-based knowledge management
- **Status**: Future (Phase 4)
- **When to Add**: When implementing Neo4j repository backend
- **Version**: Latest stable
- **License**: Apache 2.0

**Usage:**
```go
import (
    "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func NewGraphRepository(uri, username, password string) (*GraphRepository, error) {
    driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
    if err != nil {
        return nil, err
    }
    return &GraphRepository{driver: driver}, nil
}
```

**Configuration:**
```yaml
storage: "graph"
neo4j_url: "bolt://localhost:7687"
```
```

---

### MCP Integration

```markdown
**MCP Integration:**
- **Package**: `github.com/modelcontextprotocol/sdk-go` (hypothetical)
- **Purpose**: Model Context Protocol integration for specialist servers
- **Status**: Future (Phase 5)
- **When to Add**: When implementing MCP provider
- **Version**: Latest stable
- **License**: TBD

**Usage:**
```go
import (
    "github.com/modelcontextprotocol/sdk-go"
)

func NewMCPProvider(host string) (*MCPProvider, error) {
    client := mcp.NewClient(host)
    return &MCPProvider{client: client}, nil
}
```

**Configuration:**
```yaml
ai_provider: "mcp"
mcp_host: "localhost:8080"
specialists:
  - "code-review"
  - "documentation"
```
```

---

### Plugin System

```markdown
**Plugin System:**
- **Package**: `github.com/hashicorp/go-plugin`
- **Purpose**: Plugin system for third-party specialist servers
- **Status**: Future (Phase 3)
- **When to Add**: When implementing plugin system
- **Version**: Latest stable
- **License**: MPL-2.0

**Usage:**
```go
import (
    "github.com/hashicorp/go-plugin"
)

func LoadPlugin(pluginPath string) (Specialist, error) {
    client := plugin.NewClient(&plugin.ClientConfig{
        HandshakeConfig: plugin.HandshakeConfig{
            ProtocolVersion:  1,
            MagicCookieKey:   "BASIC_PLUGIN",
            MagicCookieValue: "hello",
        },
        Plugins: map[string]plugin.Plugin{
            "specialist": &SpecialistPlugin{},
        },
        Cmd: exec.Command(pluginPath),
    }
    
    rpcClient, err := client.Client()
    if err != nil {
        return nil, err
    }
    
    raw, err := rpcClient.Dispense("specialist")
    if err != nil {
        return nil, err
    }
    
    return raw.(Specialist), nil
}
```

**Configuration:**
```yaml
enable_plugins: true
plugin_dir: "~/.promptstack/plugins"
```
```

---

## Dependency Categories

### Current Dependencies (Phase 1)

```markdown
### Core Dependencies
- `github.com/anthropics/anthropic-sdk-go` - Claude API client
- `modernc.org/sqlite` - SQLite database driver
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling for TUI
- `github.com/sahilm/fuzzy` - Fuzzy search
- `gopkg.in/yaml.v3` - YAML parsing
- `github.com/sergi/go-diff` - Diff generation
- `github.com/fsnotify/fsnotify` - File system notifications
```

### Future Dependencies (Phase 2+)

```markdown
### Storage Backends
- `github.com/lib/pq` - PostgreSQL driver (Phase 4)
- `github.com/neo4j/neo4j-go-driver` - Neo4j driver (Phase 4)

### AI Providers
- `github.com/modelcontextprotocol/sdk-go` - MCP SDK (Phase 5)
- `github.com/sashabaranov/go-openai` - OpenAI client (Phase 5)

### Plugin System
- `github.com/hashicorp/go-plugin` - Plugin system (Phase 3)

### Testing
- `github.com/stretchr/testify` - Testing utilities
- `github.com/golang/mock` - Mock generation
```

---

## Dependency Management

### go.mod Updates

When adding new dependencies, update [`go.mod`](../../archive/code/go.mod):

```go
module github.com/yourorg/promptstack

go 1.21

require (
    github.com/anthropics/anthropic-sdk-go v0.0.0
    modernc.org/sqlite v1.28.0
    github.com/charmbracelet/bubbletea v0.24.2
    // ... existing dependencies
    
    // Future dependencies (add when needed)
    github.com/lib/pq v1.10.9 // Phase 4
    github.com/neo4j/neo4j-go-driver/v4 v4.4.1 // Phase 4
    github.com/modelcontextprotocol/sdk-go v0.0.0 // Phase 5
    github.com/hashicorp/go-plugin v1.5.5 // Phase 3
)
```

### Dependency Versioning

```markdown
**Version Strategy:**
- Use semantic versioning (semver)
- Pin to specific versions in go.mod
- Update dependencies regularly for security patches
- Test thoroughly before upgrading major versions

**Version Ranges:**
- Use `vX.Y.Z` for specific versions
- Use `vX.Y` for compatible minor versions
- Avoid using `latest` or master branches

**Security:**
- Monitor for security advisories
- Update vulnerable dependencies promptly
- Use `go get -u` to check for updates
- Use `go mod tidy` to clean up dependencies
```

---

## Dependency Testing

### Integration Testing

```markdown
**Testing with External Dependencies:**

1. **SQLite Testing:**
   - Use in-memory databases for unit tests
   - Use test databases for integration tests
   - Clean up test data after each test

2. **PostgreSQL Testing (Future):**
   - Use Docker containers for integration tests
   - Use test databases separate from production
   - Mock PostgreSQL for unit tests

3. **Neo4j Testing (Future):**
   - Use Docker containers for integration tests
   - Use test databases separate from production
   - Mock Neo4j for unit tests

4. **AI Provider Testing:**
   - Mock AI providers for unit tests
   - Use test API keys for integration tests
   - Record and replay API responses for tests
```

### Mock Dependencies

```markdown
**Mocking External Dependencies:**

```go
// Mock AI provider for testing
type MockProvider struct {
    suggestions []Suggestion
    err         error
}

func (m *MockProvider) GetSuggestions(ctx context.Context, prompt string) ([]Suggestion, error) {
    return m.suggestions, m.err
}

// Mock repository for testing
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
```
```

---

## Dependency Migration

### Migration Strategy

```markdown
**When Adding New Dependencies:**

1. **Phase 1 (Current):**
   - SQLite only
   - Claude API only
   - Filesystem library only

2. **Phase 2 (Future):**
   - Add plugin system
   - Add OpenAI provider

3. **Phase 3 (Future):**
   - Add MCP integration
   - Add specialist servers

4. **Phase 4 (Future):**
   - Add PostgreSQL backend
   - Add Neo4j backend

**Migration Path:**
- Keep existing dependencies
- Add new dependencies incrementally
- Maintain backward compatibility
- Provide migration tools when needed
```

---

**Last Updated**: 2026-01-07  
**Related Documents**: 
- [`scalability-implementation-summary.md`](./scalability-implementation-summary.md)
- [`scalability-project-structure-updates.md`](./scalability-project-structure-updates.md)
- [`scalability-config-schema-updates.md`](./scalability-config-schema-updates.md)