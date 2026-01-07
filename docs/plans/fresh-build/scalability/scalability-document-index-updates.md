# Document Index Updates for Scalability

**Document**: [`DOCUMENT-INDEX.md`](../DOCUMENT-INDEX.md)  
**Priority**: **Low**  
**Status**: Should update after other documents

---

## Overview

This document details all changes required to [`DOCUMENT-INDEX.md`](../DOCUMENT-INDEX.md) to incorporate Phase 1 scalability abstractions. These changes update references to reflect new packages and milestones.

---

## Domain Updates

### AI Domain

```markdown
### AI Domain
**Milestones**: M27, M28, M29, M30, M31, M32, M33, M39, M41
**Key Documents**:
- [`project-structure.md`](../project-structure.md) (ai domain - updated with provider interface)
- [`requirements.md`](../requirements.md) (AI sections - updated with provider abstraction)
- [`DEPENDENCIES.md`](../DEPENDENCIES.md) (anthropic-sdk-go, tidwall/*, sergi/go-diff)
- [`go-style-guide.md`](../go-style-guide.md) (interface design, factory pattern, middleware pattern - NEW)
- [`scalability-review.md`](../scalability-review.md) (AI domain analysis - NEW)

**New Packages**:
- `internal/ai/provider.go` - AIProvider interface
- `internal/ai/claude.go` - Claude implementation
- `internal/ai/selector.go` - ContextSelector interface
- `internal/ai/middleware.go` - ProviderMiddleware type

**Updated Packages**:
- `internal/ai/context.go` - Refactored to implement ContextSelector
```

---

### Storage Domain (NEW)

```markdown
### Storage Domain (NEW)
**Milestones**: M15, M16, M17
**Key Documents**:
- [`project-structure.md`](../project-structure.md) (storage domain - NEW)
- [`DATABASE-SCHEMA.md`](../DATABASE-SCHEMA.md) (all sections - updated with migration strategy)
- [`requirements.md`](../requirements.md) (History section - updated with storage abstraction)
- [`DEPENDENCIES.md`](../DEPENDENCIES.md) (modernc.org/sqlite, lib/pq, neo4j-go-driver - updated)
- [`go-style-guide.md`](../go-style-guide.md) (factory pattern - NEW)
- [`scalability-review.md`](../scalability-review.md) (Storage domain analysis - NEW)

**New Packages**:
- `internal/storage/repository.go` - CompositionRepository interface
- `internal/storage/sqlite.go` - SQLite implementation
- `internal/storage/factory.go` - Repository factory
- `internal/storage/postgres.go` - PostgreSQL implementation (future)
- `internal/storage/graph.go` - Neo4j implementation (future)

**Refactored Packages**:
- `internal/history/manager.go` - Use repository pattern
- `internal/history/database.go` - Deprecated, moved to storage/sqlite.go
- `internal/history/storage.go` - Deprecated, moved to storage/sqlite.go
- `internal/history/sync.go` - Use repository pattern
- `internal/history/search.go` - Use repository pattern
- `internal/history/cleanup.go` - Use repository pattern
- `internal/history/listing.go` - Use repository pattern
```

---

### Library Domain

```markdown
### Library Domain
**Milestones**: M7, M8, M9, M10, M23, M24, M25, M26
**Key Documents**:
- [`project-structure.md`](../project-structure.md) (library domain - updated with source interface)
- [`requirements.md`](../requirements.md) (Library section - updated with source abstraction)
- [`DEPENDENCIES.md`](../DEPENDENCIES.md) (sahilm/fuzzy, yaml.v3)
- [`go-style-guide.md`](../go-style-guide.md) (interface design - NEW)
- [`scalability-review.md`](../scalability-review.md) (Library domain analysis - NEW)

**New Packages**:
- `internal/library/source.go` - PromptSource interface
- `internal/library/filesystem.go` - Filesystem implementation
- `internal/library/cache.go` - Prompt cache
- `internal/library/specialist.go` - MCP specialist source (future)
- `internal/library/remote.go` - Remote repository source (future)

**Refactored Packages**:
- `internal/library/library.go` - Use PromptSource
- `internal/library/loader.go` - Deprecated, moved to filesystem.go
```

---

### Events Domain (NEW)

```markdown
### Events Domain (NEW)
**Milestones**: M40
**Key Documents**:
- [`project-structure.md`](../project-structure.md) (events domain - NEW)
- [`go-style-guide.md`](../go-style-guide.md) (event patterns - NEW)
- [`scalability-review.md`](../scalability-review.md) (Domain Events section - NEW)

**New Packages**:
- `internal/events/events.go` - Event types
- `internal/events/dispatcher.go` - Event dispatcher
```

---

## Milestone Updates

### Updated Milestones

```markdown
### Milestone 7: Prompt Source Interface & Filesystem Implementation
**Status**: Updated
**Changes**: 
- Added PromptSource interface
- Added FilesystemSource implementation
- Added PromptCache interface
- Added MemoryCache implementation
- Updated test criteria

**Related Documents**:
- [`scalability-milestones-updates.md`](./scalability-milestones-updates.md)
- [`scalability-project-structure-updates.md`](./scalability-project-structure-updates.md)
```

```markdown
### Milestone 15: Repository Pattern & SQLite Implementation
**Status**: Updated
**Changes**: 
- Added CompositionRepository interface
- Added SQLiteRepository implementation
- Added repository factory
- Updated test criteria

**Related Documents**:
- [`scalability-milestones-updates.md`](./scalability-milestones-updates.md)
- [`scalability-project-structure-updates.md`](./scalability-project-structure-updates.md)
- [`scalability-database-schema-updates.md`](./scalability-database-schema-updates.md)
```

```markdown
### Milestone 27: AI Provider Interface & Claude Implementation
**Status**: Updated
**Changes**: 
- Added AIProvider interface
- Added ClaudeProvider implementation
- Added provider factory
- Updated test criteria

**Related Documents**:
- [`scalability-milestones-updates.md`](./scalability-milestones-updates.md)
- [`scalability-project-structure-updates.md`](./scalability-project-structure-updates.md)
```

---

### New Milestones

```markdown
### Milestone 39: Context Selector Interface
**Status**: NEW
**Description**: Implement pluggable context selection algorithm

**Deliverables**:
- ContextSelector interface in `internal/ai/selector.go`
- DefaultSelector implementation
- Scoring algorithm implementation
- Token budget enforcement

**Related Documents**:
- [`scalability-milestones-updates.md`](./scalability-milestones-updates.md)
- [`scalability-project-structure-updates.md`](./scalability-project-structure-updates.md)
```

```markdown
### Milestone 40: Domain Events System
**Status**: NEW
**Description**: Implement domain events for decoupling

**Deliverables**:
- Event interface in `internal/events/events.go`
- Event types (CompositionSaved, PromptUsed, SuggestionAccepted)
- Event dispatcher in `internal/events/dispatcher.go`
- Subscribe/Publish pattern
- Async event handling

**Related Documents**:
- [`scalability-milestones-updates.md`](./scalability-milestones-updates.md)
- [`scalability-project-structure-updates.md`](./scalability-project-structure-updates.md)
```

```markdown
### Milestone 41: AI Provider Middleware
**Status**: NEW
**Description**: Implement middleware pattern for cross-cutting concerns

**Deliverables**:
- ProviderMiddleware type in `internal/ai/middleware.go`
- WithLogging middleware
- WithCaching middleware
- WithMetrics middleware
- Middleware chaining support

**Related Documents**:
- [`scalability-milestones-updates.md`](./scalability-milestones-updates.md)
- [`scalability-project-structure-updates.md`](./scalability-project-structure-updates.md)
```

---

## Document Cross-References

### Scalability Implementation Documents

```markdown
### Scalability Implementation Plan
**Purpose**: Update existing planning documents to incorporate Phase 1 scalability abstractions

**Documents**:
- [`scalability-implementation-summary.md`](./scalability-implementation-summary.md) - Executive summary and overview
- [`scalability-project-structure-updates.md`](./scalability-project-structure-updates.md) - Project structure changes
- [`scalability-config-schema-updates.md`](./scalability-config-schema-updates.md) - Configuration schema changes
- [`scalability-milestones-updates.md`](./scalability-milestones-updates.md) - Milestone updates
- [`scalability-requirements-updates.md`](./scalability-requirements-updates.md) - Requirements updates
- [`scalability-go-style-guide-updates.md`](./scalability-go-style-guide-updates.md) - Style guide updates
- [`scalability-dependencies-updates.md`](./scalability-dependencies-updates.md) - Dependencies updates
- [`scalability-database-schema-updates.md`](./scalability-database-schema-updates.md) - Database schema updates
- [`scalability-document-index-updates.md`](./scalability-document-index-updates.md) - Document index updates
- [`scalability-implementation-order.md`](./scalability-implementation-order.md) - Implementation order
- [`scalability-architecture-evolution.md`](./scalability-architecture-evolution.md) - Architecture evolution
- [`scalability-implementation-plan-index.md`](./scalability-implementation-plan-index.md) - Master index
```

---

## Configuration Updates

### New Configuration Fields

```markdown
### Configuration Schema Updates
**Document**: [`CONFIG-SCHEMA.md`](../CONFIG-SCHEMA.md)

**New Fields**:
- `ai_provider` - AI provider selection (claude, mcp, openai)
- `storage` - Storage backend selection (sqlite, postgres, graph)
- `database_path` - SQLite database path
- `postgres_url` - PostgreSQL connection string (future)
- `neo4j_url` - Neo4j connection URL (future)
- `mcp_host` - MCP orchestrator address (future)
- `specialists` - Enabled specialist servers (future)
- `enable_plugins` - Enable external plugins (future)
- `plugin_dir` - Plugin directory (future)

**Related Documents**:
- [`scalability-config-schema-updates.md`](./scalability-config-schema-updates.md)
```

---

## Pattern Updates

### New Design Patterns

```markdown
### Design Patterns
**Document**: [`go-style-guide.md`](../go-style-guide.md)

**New Patterns**:
- Interface Design - Define interfaces where used
- Factory Pattern - Create objects based on configuration
- Middleware Pattern - Wrap objects with cross-cutting concerns
- Event Pattern - Pub/sub pattern for domain events
- Repository Pattern - Abstract data access logic

**Related Documents**:
- [`scalability-go-style-guide-updates.md`](./scalability-go-style-guide-updates.md)
```

---

## Testing Updates

### Testing Guides

```markdown
### Testing Guides
**Updated for Scalability**:

- [`milestones/FOUNDATION-TESTING-GUIDE.md`](../milestones/FOUNDATION-TESTING-GUIDE.md) - Updated with interface testing
- [`milestones/AI-INTEGRATION-TESTING-GUIDE.md`](../milestones/AI-INTEGRATION-TESTING-GUIDE.md) - Updated with provider testing
- [`milestones/HISTORY-TESTING-GUIDE.md`](../milestones/HISTORY-TESTING-GUIDE.md) - Updated with repository testing
- [`milestones/LIBRARY-INTEGRATION-TESTING-GUIDE.md`](../milestones/LIBRARY-INTEGRATION-TESTING-GUIDE.md) - Updated with source testing

**New Testing Requirements**:
- Interface testing with mocks
- Factory testing with different configurations
- Middleware testing with chaining
- Event testing with async handlers
- Repository testing with different backends
```

---

## Migration Updates

### Database Migration

```markdown
### Database Migration
**Document**: [`DATABASE-SCHEMA.md`](../DATABASE-SCHEMA.md)

**New Sections**:
- Migration Strategy
- Version Management
- Migration Files
- Rollback Support
- Cross-Database Migration (SQLite → PostgreSQL)

**Related Documents**:
- [`scalability-database-schema-updates.md`](./scalability-database-schema-updates.md)
```

---

## Summary of Changes

### Documents Updated

1. **project-structure.md** - Added 6 new packages, refactored 3 existing
2. **CONFIG-SCHEMA.md** - Added 8 new config fields
3. **milestones.md** - Updated 3 milestones, added 3 new milestones
4. **requirements.md** - Updated AI, Storage, Library sections, added Domain Events
5. **go-style-guide.md** - Added interface, factory, middleware, event patterns
6. **DEPENDENCIES.md** - Added future dependencies (PostgreSQL, Neo4j, MCP)
7. **DATABASE-SCHEMA.md** - Added migration strategy section
8. **DOCUMENT-INDEX.md** - Updated references (this document)

### New Documents Created

1. **scalability-implementation-summary.md** - Executive summary
2. **scalability-project-structure-updates.md** - Project structure changes
3. **scalability-config-schema-updates.md** - Configuration changes
4. **scalability-milestones-updates.md** - Milestone updates
5. **scalability-requirements-updates.md** - Requirements updates
6. **scalability-go-style-guide-updates.md** - Style guide updates
7. **scalability-dependencies-updates.md** - Dependencies updates
8. **scalability-database-schema-updates.md** - Database schema updates
9. **scalability-document-index-updates.md** - Document index updates
10. **scalability-implementation-order.md** - Implementation order
11. **scalability-architecture-evolution.md** - Architecture evolution
12. **scalability-implementation-plan-index.md** - Master index

---

**Last Updated**: 2026-01-07  
**Related Documents**: 
- [`scalability-implementation-summary.md`](./scalability-implementation-summary.md)
- [`scalability-implementation-plan-index.md`](./scalability-implementation-plan-index.md)