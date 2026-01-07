# Milestones Updates for Scalability

**Document**: [`milestones.md`](../milestones.md)  
**Priority**: **Critical**  
**Status**: Must update before implementation begins

---

## Overview

This document details all changes required to [`milestones.md`](../milestones.md) to incorporate Phase 1 scalability abstractions. These changes update existing milestones and add new milestones for scalability features.

---

## A. Update M27: Claude API Client

### Current M27
```
## Milestone 27: Claude API Client
**Goal:** Send requests to Claude API

**Deliverables:**
- API client wrapper
- Authentication with API key
- Send message, receive response
- Error handling (rate limit, auth, network, timeout)
- Retry logic (max 3 attempts)
```

### Updated M27
```
## Milestone 27: AI Provider Interface & Claude Implementation
**Goal:** Implement AI provider abstraction with Claude as first provider

**Deliverables:**
- AIProvider interface in `internal/ai/provider.go`
- ClaudeProvider implementation in `internal/ai/claude.go`
- Provider factory in `internal/ai/factory.go`
- Authentication with API key
- Send message, receive response
- Error handling (rate limit, auth, network, timeout)
- Retry logic (max 3 attempts)
- Unit tests for interface and implementation
- Integration tests with Claude API

**Test Criteria:**

### Functional Requirements
- [ ] AIProvider interface defined with required methods
- [ ] ClaudeProvider implements AIProvider interface
- [ ] Provider factory creates correct provider based on config
- [ ] Send test request to Claude API
- [ ] Receive valid response
- [ ] Authentication with API key works
- [ ] Send message, receive response
- [ ] Handle 401 auth error
- [ ] Handle 429 rate limit
- [ ] Handle network timeout
- [ ] Retry on transient failures
- [ ] Stop after 3 failed retries
- [ ] Exponential backoff between retries

### Integration Requirements
- [ ] Provider factory integrates with config system
- [ ] Provider factory reads ai_provider field from config
- [ ] API client integrates with logging system
- [ ] Error handling integrates with error handler

### Edge Cases & Error Handling
- [ ] Handle missing API key (show error, prompt to configure)
- [ ] Handle invalid API key (show error, prompt to reconfigure)
- [ ] Handle API service unavailable (retry with backoff)
- [ ] Handle malformed API response (show error, log details)
- [ ] Handle concurrent API requests (queue properly)
- [ ] Handle very long responses (>100KB)
- [ ] Handle rate limit exceeded (wait and retry)
- [ ] Handle network connection refused
- [ ] Handle DNS resolution failures

### Performance Requirements
- [ ] API request latency <2s for typical response
- [ ] Retry backoff: 1s, 2s, 4s (exponential)
- [ ] Connection timeout: 30s
- [ ] Read timeout: 60s
- [ ] Handle 10 concurrent requests

### User Experience Requirements
- [ ] Clear error messages for API failures
- [ ] Progress indicator during API requests
- [ ] Retry status shown in status bar
- [ ] No UI blocking during API requests
- [ ] Graceful degradation if API unavailable

**Files:** [`internal/ai/provider.go`](../../archive/code/internal/ai/provider.go), [`internal/ai/claude.go`](../../archive/code/internal/ai/claude.go), [`internal/ai/factory.go`](../../archive/code/internal/ai/factory.go)
```

---

## B. Update M15: SQLite Setup

### Current M15
```
## Milestone 15: SQLite Setup
**Goal:** Initialize history database

**Deliverables:**
- Database schema with FTS5
- Create at `~/.promptstack/data/history.db`
- Basic CRUD operations
- Index on created_at and working_directory
```

### Updated M15
```
## Milestone 15: Repository Pattern & SQLite Implementation
**Goal:** Implement repository abstraction with SQLite as first backend

**Deliverables:**
- CompositionRepository interface in `internal/storage/repository.go`
- SQLiteRepository implementation in `internal/storage/sqlite.go`
- Repository factory in `internal/storage/factory.go`
- Database schema with FTS5
- Create at `~/.promptstack/data/history.db`
- Basic CRUD operations
- Index on created_at and working_directory
- Unit tests for interface and implementation
- Integration tests with SQLite

**Test Criteria:**

### Functional Requirements
- [ ] CompositionRepository interface defined with required methods
- [ ] SQLiteRepository implements CompositionRepository interface
- [ ] Repository factory creates correct repository based on config
- [ ] Database file created at `~/.promptstack/data/history.db` on first run
- [ ] Schema matches specification (compositions table with FTS5)
- [ ] Insert composition record with all fields
- [ ] Query composition by ID
- [ ] Full-text search works on content field
- [ ] Update existing composition record
- [ ] Delete composition record
- [ ] Index created on created_at column
- [ ] Index created on working_directory column

### Integration Requirements
- [ ] Repository factory integrates with config system
- [ ] Repository factory reads storage field from config
- [ ] Database initialization integrates with config system
- [ ] Database operations integrate with logging system
- [ ] Error handling integrates with error handler

### Edge Cases & Error Handling
- [ ] Handle missing data directory (create automatically)
- [ ] Handle corrupted database file (show error, offer rebuild)
- [ ] Handle database locked (retry with backoff)
- [ ] Handle disk full errors
- [ ] Handle permission denied errors
- [ ] Handle concurrent database access
- [ ] Handle very large records (>1MB content)

### Performance Requirements
- [ ] Database initialization <100ms
- [ ] Insert operation <10ms
- [ ] Query by ID <5ms
- [ ] Full-text search <50ms for 1000 records
- [ ] Update operation <10ms
- [ ] Delete operation <10ms

### User Experience Requirements
- [ ] Clear error messages for database failures
- [ ] Progress indicators for long operations
- [ ] No UI blocking during database operations
- [ ] Graceful degradation if database unavailable

**Files:** [`internal/storage/repository.go`](../../archive/code/internal/storage/repository.go), [`internal/storage/sqlite.go`](../../archive/code/internal/storage/sqlite.go), [`internal/storage/factory.go`](../../archive/code/internal/storage/factory.go)
```

---

## C. Update M7: Library Loader

### Current M7
```
## Milestone 7: Library Loader
**Goal:** Load prompts from filesystem into memory

**Deliverables:**
- Scan `~/.promptstack/data/` directory
- Parse YAML frontmatter from each prompt
- Load into in-memory collection
- Extract metadata (title, tags, category, description)
```

### Updated M7
```
## Milestone 7: Prompt Source Interface & Filesystem Implementation
**Goal:** Implement prompt source abstraction with filesystem as first source

**Deliverables:**
- PromptSource interface in `internal/library/source.go`
- FilesystemSource implementation in `internal/library/filesystem.go`
- PromptCache interface in `internal/library/cache.go`
- MemoryCache implementation in `internal/library/cache.go`
- Scan `~/.promptstack/data/` directory
- Parse YAML frontmatter from each prompt
- Load into in-memory collection
- Extract metadata (title, tags, category, description)
- Unit tests for interfaces and implementations
- Integration tests with filesystem

**Test Criteria:**

### Functional Requirements
- [ ] PromptSource interface defined with required methods
- [ ] FilesystemSource implements PromptSource interface
- [ ] PromptCache interface defined with required methods
- [ ] MemoryCache implements PromptCache interface
- [ ] Load sample prompts from directory
- [ ] Verify prompt count matches files
- [ ] Metadata parsed correctly (title, tags, category, description)
- [ ] Category derived from folder structure
- [ ] Invalid prompts loaded (not excluded)
- [ ] Empty directories handled gracefully
- [ ] Nested directories scanned recursively
- [ ] Cache stores and retrieves prompts correctly
- [ ] Cache invalidation works correctly

### Integration Requirements
- [ ] FilesystemSource integrates with config system
- [ ] Metadata extraction integrates with prompt structure
- [ ] Error handling integrates with logging system
- [ ] Cache integrates with FilesystemSource

### Edge Cases & Error Handling
- [ ] Handle missing library directory (create automatically)
- [ ] Handle empty library directory
- [ ] Handle files without frontmatter
- [ ] Handle malformed YAML frontmatter
- [ ] Handle duplicate filenames
- [ ] Handle very large libraries (1000+ prompts)
- [ ] Handle permission denied errors
- [ ] Handle circular symlinks (if applicable)
- [ ] Handle cache misses
- [ ] Handle cache invalidation during load

### Performance Requirements
- [ ] Startup time <500ms for 100 prompts
- [ ] Startup time <2s for 1000 prompts
- [ ] Memory usage scales linearly with prompt count
- [ ] File I/O operations are concurrent where possible
- [ ] Cache hit time <1ms
- [ ] Cache miss time <10ms

### User Experience Requirements
- [ ] Loading indicator shown during startup
- [ ] Error messages are specific to file/issue
- [ ] Progress indicator for large libraries
- [ ] No UI blocking during load

**Files:** [`internal/library/source.go`](../../archive/code/internal/library/source.go), [`internal/library/filesystem.go`](../../archive/code/internal/library/filesystem.go), [`internal/library/cache.go`](../../archive/code/internal/library/cache.go)
```

---

## D. Add New Milestones

### Milestone 39: Context Selector Interface

```
## Milestone 39: Context Selector Interface
**Goal:** Implement pluggable context selection algorithm

**Deliverables:**
- ContextSelector interface in `internal/ai/selector.go`
- DefaultSelector implementation
- Scoring algorithm implementation
- Token budget enforcement
- Unit tests for interface and implementation
- Integration tests with library

**Test Criteria:**

### Functional Requirements
- [ ] ContextSelector interface defined with required methods
- [ ] DefaultSelector implements ContextSelector interface
- [ ] Given composition, extract keywords
- [ ] Score all library prompts
- [ ] Top 3-5 prompts selected
- [ ] Relevant prompts ranked higher
- [ ] Token budget respected
- [ ] Explicitly referenced prompts always included
- [ ] Scoring algorithm applied correctly
  - [ ] Tag match: +10 per matching tag
  - [ ] Category match: +5 if same category
  - [ ] Keyword overlap: +1 per matching word
  - [ ] Recently used: +3 if used in last session
  - [ ] Frequently used: +use_count

### Integration Requirements
- [ ] Context selection integrates with library loader
- [ ] Context selection integrates with placeholder parser
- [ ] Context selection integrates with token estimator
- [ ] Context selection integrates with history tracking

### Edge Cases & Error Handling
- [ ] Handle empty composition (no keywords extracted)
- [ ] Handle composition with no matching prompts
- [ ] Handle library with no prompts
- [ ] Handle very large library (1000+ prompts)
- [ ] Handle very long composition (>10000 words)
- [ ] Handle prompts with no metadata
- [ ] Handle prompts with duplicate scores
- [ ] Handle token budget exceeded (select fewer prompts)
- [ ] Handle explicitly referenced prompts exceeding budget

### Performance Requirements
- [ ] Keyword extraction <50ms for 1000-word composition
- [ ] Scoring <100ms for 1000 prompts
- [ ] Selection <10ms
- [ ] Total context selection <200ms
- [ ] Memory usage scales linearly with prompt count

### User Experience Requirements
- [ ] Selected prompts shown in status bar
- [ ] Token usage displayed
- [ ] Clear indication of context limit
- [ ] No UI blocking during selection
- [ ] Progress indicator for large libraries

**Files:** [`internal/ai/selector.go`](../../archive/code/internal/ai/selector.go)
```

---

### Milestone 40: Domain Events System

```
## Milestone 40: Domain Events System
**Goal:** Implement domain events for decoupling

**Deliverables:**
- Event interface in `internal/events/events.go`
- Event types (CompositionSaved, PromptUsed, SuggestionAccepted)
- Event dispatcher in `internal/events/dispatcher.go`
- Subscribe/Publish pattern
- Async event handling
- Unit tests for events and dispatcher
- Integration tests with components

**Test Criteria:**

### Functional Requirements
- [ ] Event interface defined with required methods
- [ ] BaseEvent implements Event interface
- [ ] CompositionSavedEvent defined
- [ ] PromptUsedEvent defined
- [ ] SuggestionAcceptedEvent defined
- [ ] Dispatcher manages event subscriptions
- [ ] Subscribe registers handler for event type
- [ ] Publish sends event to all handlers
- [ ] Events handled asynchronously
- [ ] Multiple handlers can subscribe to same event
- [ ] Handlers receive correct event data

### Integration Requirements
- [ ] Events integrate with composition save
- [ ] Events integrate with prompt use
- [ ] Events integrate with suggestion accept
- [ ] Dispatcher integrates with logging system
- [ ] Event handlers integrate with analytics (future)

### Edge Cases & Error Handling
- [ ] Handle handler panics (recover and log)
- [ ] Handle no handlers subscribed (no-op)
- [ ] Handle rapid event publishing (queue properly)
- [ ] Handle handler errors (log and continue)
- [ ] Handle concurrent subscriptions
- [ ] Handle concurrent publishing
- [ ] Handle very large event payloads

### Performance Requirements
- [ ] Subscribe operation <1ms
- [ ] Publish operation <1ms
- [ ] Event handler execution <10ms (typical)
- [ ] No blocking on publish
- [ ] Memory usage scales with handler count

### User Experience Requirements
- [ ] No UI blocking during event handling
- [ ] Errors in handlers don't crash app
- [ ] Events processed in reasonable time
- [ ] No event loss under normal load

**Files:** [`internal/events/events.go`](../../archive/code/internal/events/events.go), [`internal/events/dispatcher.go`](../../archive/code/internal/events/dispatcher.go)
```

---

### Milestone 41: AI Provider Middleware

```
## Milestone 41: AI Provider Middleware
**Goal:** Implement middleware pattern for cross-cutting concerns

**Deliverables:**
- ProviderMiddleware type in `internal/ai/middleware.go`
- WithLogging middleware
- WithCaching middleware
- WithMetrics middleware
- Middleware chaining support
- Unit tests for all middleware
- Integration tests with providers

**Test Criteria:**

### Functional Requirements
- [ ] ProviderMiddleware type defined
- [ ] WithLogging wraps provider with logging
- [ ] WithCaching wraps provider with caching
- [ ] WithMetrics wraps provider with metrics
- [ ] Middleware can be chained
- [ ] Middleware preserves provider interface
- [ ] Logging middleware logs all calls
- [ ] Caching middleware caches results
- [ ] Metrics middleware collects metrics

### Integration Requirements
- [ ] Middleware integrates with provider factory
- [ ] Logging middleware integrates with logging system
- [ ] Caching middleware integrates with cache
- [ ] Metrics middleware integrates with metrics system
- [ ] Middleware chain applies in correct order

### Edge Cases & Error Handling
- [ ] Handle middleware errors (propagate to caller)
- [ ] Handle cache misses (call provider)
- [ ] Handle cache invalidation
- [ ] Handle logging failures (don't block)
- [ ] Handle metrics collection failures (don't block)
- [ ] Handle empty middleware chain
- [ ] Handle very long middleware chains

### Performance Requirements
- [ ] Middleware overhead <1ms per call
- [ ] Logging middleware <5ms per call
- [ ] Caching middleware <1ms (hit), <provider time (miss)
- [ ] Metrics middleware <1ms per call
- [ ] No blocking in middleware

### User Experience Requirements
- [ ] No UI blocking from middleware
- [ ] Logging visible in debug mode
- [ ] Caching improves performance
- [ ] Metrics available for monitoring
- [ ] Middleware errors don't crash app

**Files:** [`internal/ai/middleware.go`](../../archive/code/internal/ai/middleware.go)
```

---

## E. Update Milestone Summary

Update the milestone summary at the end of milestones.md:

```
## Summary

**Total Milestones:** 41 (updated from 38)

**Milestone Groups:**
1. **Foundation** (1-6): Bootstrap, TUI shell, file I/O, basic editor, auto-save, undo/redo
2. **Library Integration** (7-10): Load library, browse, search, insert prompts
3. **Placeholders** (11-14): Parse, navigate, edit text/list placeholders
4. **History** (15-17): Repository pattern, sync, browser
5. **Commands & Files** (18-22): Command system, palette, file references
6. **Prompt Management** (23-26): Validation, results, creator, editor
7. **AI Integration** (27-33): Provider interface, context selection, tokens, suggestions, diff
8. **Vim Mode** (34-35): State machine, keybindings
9. **Polish** (36-38): Settings, responsive layout, error handling
10. **Scalability** (39-41): Context selector, domain events, middleware (NEW)

**Key Principles:**
- Each milestone is independently testable
- Clear pass/fail criteria for every test
- Incremental complexity with tight integration loops
- Library + placeholders working together early (Milestone 10-14)
- Scalability abstractions integrated early (Milestone 27, 15, 7, 39-41)
- No time estimates - focus on deliverables
- Test-driven approach ensures stability
```

---

**Last Updated**: 2026-01-07  
**Related Documents**: 
- [`scalability-implementation-summary.md`](./scalability-implementation-summary.md)
- [`scalability-project-structure-updates.md`](./scalability-project-structure-updates.md)
- [`scalability-config-schema-updates.md`](./scalability-config-schema-updates.md)