---
title: "PromptStack - Complete Implementation Plan"
domain: "CLI Tools, TUI Applications, AI Integration"
keywords: ["go", "bubble-tea", "tui", "cli", "prompt-management", "claude-api", "markdown", "sqlite"]
related: ["requirements.md", "code-samples/"]
last_updated: "2026-01-06"
status: "draft"
---

# PromptStack - Complete Implementation Plan

## Overview and Objectives

### Project Summary
PromptStack is a Go-based TUI (Text User Interface) application that enables users to build complex AI prompts through interactive composition. The tool combines reusable prompt templates from a global library, supports placeholder-based templating, provides file references, and offers AI-powered suggestions through Claude API integration.

### Core Objectives
1. **Intuitive Prompt Composition** - Provide a vim-style editor for building prompts with auto-save and undo/redo
2. **Reusable Template Library** - Manage a global collection of prompt templates with metadata and validation
3. **Smart Templating** - Support text and list placeholders with vim-style editing interactions
4. **AI Integration** - Offer conservative, token-efficient AI suggestions via Claude API
5. **History Management** - Track composition history with SQLite-backed search and cleanup capabilities
6. **Developer Experience** - Minimal, modern UI with universal vim mode support and comprehensive logging

### Success Criteria
- Responsive TUI that works in terminals from 80 to 200+ columns
- Sub-100ms response time for library searches and navigation
- Conservative AI token usage (compositions limited to 25% of context window)
- Reliable auto-save with visual feedback
- Zero data loss through markdown-as-source-of-truth architecture

## Architecture and Approach

### High-Level Architecture

**Three-Layer Architecture:**

1. **Presentation Layer (TUI)**
   - Bubble Tea models for each screen/modal
   - Lipgloss styling for minimal, modern aesthetics
   - Event routing and keybinding management
   - Vim mode state machine

2. **Business Logic Layer**
   - Library management (CRUD, indexing, validation)
   - Composition management (editing, placeholders, undo/redo)
   - History management (storage, search, cleanup)
   - AI suggestion orchestration (context selection, diff application)
   - File reference system (fuzzy finding, title extraction)

3. **Data Layer**
   - Markdown file I/O with YAML frontmatter parsing
   - SQLite database for history indexing and full-text search
   - Embedded starter prompts via go:embed
   - Config file management (YAML)

See code-samples/001-architecture-diagram.sample.md for visual architecture diagram.

### Technology Stack Selection Rationale

**TUI Framework: Bubble Tea Ecosystem**
- Battle-tested framework with strong community
- Component-based architecture aligns with our modal/panel design
- Built-in support for complex event handling needed for vim mode
- Lipgloss provides declarative styling matching our minimal aesthetic

**Database: modernc.org/sqlite**
- Pure Go implementation eliminates CGO dependency
- Simplifies cross-platform builds (macOS Intel/ARM)
- FTS5 support for full-text search in history
- Adequate performance for personal-scale usage (hundreds of compositions)

**AI Integration: anthropic-sdk-go**
- Official SDK ensures API compatibility
- Structured response handling for suggestion parsing
- Built-in retry logic for transient failures

**Fuzzy Matching: sahilm/fuzzy**
- Lightweight algorithm suitable for in-memory collections
- Fast enough for real-time filtering (<10ms for 1000+ items)
- Simple API for integration with custom TUI modals

See code-samples/002-dependency-versions.sample.txt for pinned versions and go.mod structure.

### Design Philosophy

**Minimal and Functional**
- Color only for functional emphasis (categories, placeholders, errors)
- No decorative borders or boxes
- Generous whitespace for readability
- Typography hierarchy through weight and spacing

**Vim-First When Enabled**
- Universal vim keybindings across all components
- Modal editing paradigm for placeholders
- Consistent mode indicators

**Conservative AI Usage**
- Token limits enforced at composition level
- Top 3-5 relevant library prompts only
- Visible token estimates in status bar
- Block suggestions when composition exceeds 25% context budget

## Component Breakdown

### 1. Application Bootstrap & Configuration

**Responsibilities:**
- Initialize config at `~/.promptstack/config.yaml`
- Interactive setup wizard on first run
- Version-aware starter prompt extraction
- Database initialization and verification
- Logging setup with zap

**Key Components:**
- `cmd/promptstack/main.go` - Entry point
- `internal/config` - Config management
- `internal/setup` - First-run wizard
- `internal/bootstrap` - Initialization orchestration

See code-samples/003-config-structure.sample.go for config types and validation.
See code-samples/004-bootstrap-flow.sample.go for initialization sequence.

### 2. Library Management System

**Responsibilities:**
- Load prompts from `~/.promptstack/data/` on startup
- Parse YAML frontmatter and validate metadata
- Build in-memory index for AI context selection
- Support CRUD operations on prompts
- Run validation checks (startup + manual)

**Key Components:**
- `internal/library` - Core library management
- `internal/library/loader.go` - Filesystem scanning and loading
- `internal/library/indexer.go` - Build usage index for AI context
- `internal/library/validator.go` - Validation checks
- `internal/prompt` - Prompt model and metadata types

**Index Structure:**
- Metadata: title, description, tags, category, file_path
- Usage tracking: last_used, use_count
- Content analysis: word_frequency_map (for keyword matching)
- Validation status: errors, warnings

See code-samples/005-prompt-model.sample.go for prompt types and metadata.
See code-samples/006-library-indexer.sample.go for index structure and scoring algorithm.
See code-samples/007-validator.sample.go for validation checks implementation.

### 3. Composition Workspace

**Responsibilities:**
- Primary text editor for composition
- Placeholder detection and navigation
- Undo/redo history (100 levels)
- Auto-save with debouncing (500ms-1s)
- Vim mode support (Normal/Insert/Visual)
- Split-pane layout with resizable divider

**Key Components:**
- `internal/editor` - Core text editing logic
- `internal/editor/placeholder.go` - Placeholder parsing and state
- `internal/editor/undo.go` - Undo/redo stack with smart batching
- `internal/editor/vim.go` - Vim mode state machine
- `ui/workspace` - Bubble Tea model for workspace

**Undo/Redo Smart Batching:**
- Continuous typing batched as one action (break on >1 second pause)
- Cursor movement breaks batching
- Paste = one action
- Prompt insertion = one action
- Placeholder fill = one action
- Mode change breaks batching

See code-samples/008-editor-model.sample.go for editor state structure.
See code-samples/009-placeholder-parser.sample.go for placeholder detection logic.
See code-samples/010-undo-stack.sample.go for undo/redo implementation.
See code-samples/011-vim-state-machine.sample.go for vim mode handling.

### 4. Placeholder System

**Responsibilities:**
- Parse and validate placeholders in composition
- Track placeholder positions and state
- Handle Tab/Shift+Tab navigation between placeholders
- Provide vim-style editing for text placeholders
- Provide list editing mode for list placeholders

**Placeholder Types:**
- `{{text:name}}` - Single text value
- `{{list:name}}` - Multiple items (markdown bullet list)

**Validation Rules:**
- Type must be "text" or "list" (case-sensitive)
- Name must be alphanumeric + underscores only
- No duplicate names within same prompt
- Malformed placeholders treated as literal text

**Editing Interaction:**

**Text Placeholders:**
1. Tab to highlight placeholder (enters "placeholder normal mode")
2. Press 'i' or 'a' to enter placeholder edit mode
3. Type to replace content
4. Esc to exit edit mode, returns to composition

**List Placeholders:**
1. Tab to highlight placeholder (enters "list normal mode")
2. Press 'i' to enter list edit mode
3. Up/Down to navigate items
4. 'e' to edit highlighted item (inline editor)
5. 'd' to delete highlighted item
6. 'n' or 'o' to add new item
7. Esc to exit list edit mode

See code-samples/012-placeholder-types.sample.go for placeholder data structures.
See code-samples/013-placeholder-editing.sample.go for placeholder editing state machine.
See code-samples/014-list-placeholder-ui.sample.go for list editing UI component.

### 5. History Management

**Responsibilities:**
- Create timestamped markdown files on new composition
- Auto-save composition changes to markdown + SQLite
- Provide searchable history browser
- Support history cleanup with multiple strategies
- Maintain sync between markdown files (source of truth) and SQLite index

**Storage Strategy:**
- Markdown files: `~/.promptstack/data/.history/YYYY-MM-DD_HH-MM-SS.md`
- SQLite database: `~/.promptstack/data/history.db`
- Verify sync on startup, offer rebuild if mismatch detected

**Key Components:**
- `internal/history` - History management
- `internal/history/storage.go` - Markdown file operations
- `internal/history/database.go` - SQLite operations with FTS5
- `internal/history/sync.go` - Sync verification and rebuild
- `ui/history` - History browser Bubble Tea model

**Database Schema:**
```sql
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

CREATE INDEX idx_created_at ON compositions(created_at);
CREATE INDEX idx_working_directory ON compositions(working_directory);

CREATE VIRTUAL TABLE compositions_fts USING fts5(
  content,
  working_directory,
  content='compositions',
  content_rowid='id'
);
```

See code-samples/015-history-model.sample.go for history data structures.
See code-samples/016-history-database.sample.go for SQLite operations with FTS5.
See code-samples/017-history-sync.sample.go for sync verification logic.
See code-samples/018-history-cleanup.sample.go for cleanup strategies.

### 6. Library Browser & Fuzzy Search

**Responsibilities:**
- Display flat list of all prompts with category labels
- Real-time fuzzy filtering as user types
- Preview pane showing prompt content
- Keyboard navigation (arrows or vim j/k)
- Insert at cursor or on new line

**Key Components:**
- `ui/browser` - Library browser Bubble Tea model
- `internal/fuzzy` - Fuzzy matching integration
- `ui/common/fuzzylist.go` - Reusable fuzzy list component

**Fuzzy Matching:**
- Search across: title, tags, description, category
- Rank by match quality (exact > prefix > fuzzy)
- Real-time filtering (<10ms target)

See code-samples/019-library-browser.sample.go for browser UI implementation.
See code-samples/020-fuzzy-list-component.sample.go for reusable fuzzy list.

### 7. Command Palette

**Responsibilities:**
- Modal overlay with fuzzy-filterable command list
- Execute commands by name
- Keyboard-driven navigation

**Commands:**
- Toggle AI panel
- Copy to clipboard
- Trigger AI suggestions
- Add file reference
- Create new prompt
- Edit prompt
- Save composition
- View history
- Clean up history
- Validate library
- Settings
- View logs

**Key Components:**
- `ui/palette` - Command palette Bubble Tea model
- `internal/commands` - Command registry and execution

See code-samples/021-command-palette.sample.go for palette implementation.
See code-samples/022-command-registry.sample.go for command system.

### 8. File Reference System

**Responsibilities:**
- Fuzzy find files from working directory
- Respect .gitignore rules
- Extract title from YAML frontmatter or use filename
- Support multiple file selection
- Batch title editor for selected files
- Insert markdown links: `[title](path/to/file)`

**Key Components:**
- `internal/files` - File traversal and .gitignore handling
- `internal/files/finder.go` - Fuzzy file finder
- `internal/files/title_extractor.go` - YAML frontmatter parsing
- `ui/filereference` - File reference UI modal

**Workflow:**
1. User triggers "Add file reference" from command palette
2. Fuzzy finder modal shows all non-ignored files
3. User selects one or more files
4. For each file, extract title (frontmatter → filename)
5. Batch editor shows all files with editable titles
6. User can edit titles or quick-accept all
7. Insert markdown links at cursor

See code-samples/023-file-finder.sample.go for file traversal and fuzzy finding.
See code-samples/024-title-extractor.sample.go for frontmatter title extraction.
See code-samples/025-batch-title-editor.sample.go for batch title editing UI.

### 9. AI Suggestions System

**Responsibilities:**
- Manual trigger via command palette
- Select top 3-5 relevant library prompts using index scoring
- Send composition + context to Claude API
- Parse structured suggestions (6 types)
- Display suggestions in scrollable right panel
- Handle suggestion acceptance with diff preview
- Apply changes as single undo action

**Suggestion Types:**
1. Prompt Recommendations - Suggest relevant library prompts
2. Gap Analysis - Missing context or information
3. Formatting - Better structure or organization
4. Contradictions - Conflicting instructions
5. Clarity - Unclear or ambiguous instructions
6. Reformatting - Alternative ways to structure

**Context Selection Algorithm:**
1. Extract keywords from composition (word frequency)
2. Score each library prompt:
   - Tag match: +10 per matching tag
   - Category match: +5 if same category
   - Keyword overlap: +1 per matching word (weighted by frequency)
   - Recently used: +3 if used in last session
   - Frequently used: +use_count
3. Sort by score, select top 3-5 that fit in token budget
4. Always include explicitly referenced prompts

**Token Budget:**
- Detect model context limit dynamically (200K for Claude 3)
- Composition content: Max 25% of context window
- Library prompts: Max 15% of context window
- Reserve 60% for system prompt, response, safety margin
- Show warning at 15%, block at 25%

**Applying Suggestions:**
1. User presses 'a' on selected suggestion
2. Composition enters read-only mode, show "✨ AI is applying suggestion..."
3. Send request to Claude for structured edits
4. On success: Show unified diff modal
5. User accepts → Apply as single undo action, unlock editor
6. User rejects → Discard changes, unlock editor
7. On failure: Show error in suggestion box with retry button

**Diff Format:**
- Use `github.com/sergi/go-diff` to generate unified diff
- Request from Claude: Structured edits with line ranges and new content
- Parse into diff format for display
- Apply accepted changes line-by-line

**Key Components:**
- `internal/ai` - AI integration
- `internal/ai/client.go` - Claude API client wrapper
- `internal/ai/context.go` - Context selection logic
- `internal/ai/suggestions.go` - Suggestion parsing
- `internal/ai/diff.go` - Diff generation and application
- `ui/suggestions` - Suggestions panel Bubble Tea model

See code-samples/026-ai-context-selection.sample.go for context scoring algorithm.
See code-samples/027-ai-suggestion-types.sample.go for suggestion data structures.
See code-samples/028-ai-diff-application.sample.go for diff parsing and application.
See code-samples/029-suggestions-panel.sample.go for suggestions UI.

### 10. Prompt Creation & Editing

**Responsibilities:**
- Interactive guided prompt creation
- Edit existing prompts from library or history
- Edit/Preview mode toggle (Ctrl+P)
- Preview using Glamour markdown rendering
- Explicit save required (no auto-save for prompts)

**Creation Workflow:**
1. Title (required) - Text input
2. Description (optional) - Text input
3. Tags (optional) - Comma-delimited with autocomplete
4. Category (required) - Select from existing folders
5. Initial Content - Multi-line text editor

**Edit/Preview Mode:**
- Edit Mode: Raw markdown with YAML frontmatter visible, placeholder highlighting
- Preview Mode: Rendered markdown via Glamour, frontmatter hidden or styled
- Toggle: Ctrl+P (or Cmd+P on macOS)
- Mode indicator in status bar

**Key Components:**
- `ui/prompteditor` - Prompt creation/editing Bubble Tea model
- `internal/prompt/creator.go` - Prompt creation logic
- `ui/common/markdownpreview.go` - Glamour-based preview component

See code-samples/030-prompt-creator.sample.go for creation workflow.
See code-samples/031-prompt-editor.sample.go for edit/preview toggle.
See code-samples/032-markdown-preview.sample.go for Glamour integration.

### 11. Library Validation

**Responsibilities:**
- Run automatically on startup (silent unless errors found)
- Manual trigger via command palette
- Check placeholder syntax, YAML frontmatter, file size, readability
- Display errors/warnings grouped by type
- Prevent insertion of prompts with errors

**Validation Checks:**

**Errors (block insertion):**
- Duplicate placeholder names in same prompt
- Invalid YAML frontmatter syntax
- Missing required field: title
- File size exceeds 1MB

**Warnings (allow insertion with indicator):**
- Invalid placeholder types (not text/list)
- Invalid placeholder names (non-alphanumeric/underscore)
- Malformed/unclosed placeholders
- Missing optional metadata (description/tags)

**Behavior:**
- Load invalid prompts into library (don't exclude)
- Show ⚠️ icon in library browser for invalid prompts
- Block insertion of error-level prompts
- Allow insertion of warning-level prompts

**Key Components:**
- `internal/library/validator.go` - Validation logic
- `ui/validation` - Validation results modal

See code-samples/033-validation-checks.sample.go for validation implementation.
See code-samples/034-validation-results-ui.sample.go for results display.

### 12. Settings Panel

**Responsibilities:**
- Modal overlay with form-style interface
- Toggle vim mode (requires restart)
- Edit Claude API key (masked input)
- Select Claude model from dropdown
- Immediate persistence to config.yaml

**Key Components:**
- `ui/settings` - Settings panel Bubble Tea model
- `internal/config/updater.go` - Config update and validation

See code-samples/035-settings-panel.sample.go for settings UI.

### 13. Vim Mode System

**Responsibilities:**
- Universal vim keybindings across all components when enabled
- Mode state machine (Normal/Insert/Visual)
- Context-aware keybinding application
- Mode indicator in status bar

**Scope:**
- Composition workspace: Full Normal/Insert/Visual support
- Library browser: j/k navigation, / for search
- Command palette: j/k navigation
- AI suggestions panel: j/k navigation
- History browser: j/k navigation
- Settings panel: j/k navigation between fields
- Text inputs: Basic vim editing where applicable

**Implementation Strategy:**
- Central vim mode state manager
- Context-aware event routing
- Component-specific keybinding maps
- Mode transition rules

**Key Components:**
- `internal/vim` - Vim mode state machine
- `internal/vim/keymaps.go` - Keybinding definitions per mode
- `internal/vim/context.go` - Context-aware routing

See code-samples/036-vim-state-machine.sample.go for vim mode implementation.
See code-samples/037-vim-keybindings.sample.go for keybinding maps.

### 14. Status Bar

**Responsibilities:**
- Display contextual information at bottom of screen
- Auto-save indicator (transient)
- Character/line counts
- Token estimate (when AI panel open)
- Mode indicators (vim, edit/preview)
- Notifications and warnings

**Key Components:**
- `ui/statusbar` - Status bar Bubble Tea component

See code-samples/038-status-bar.sample.go for status bar implementation.

### 15. Logging System

**Responsibilities:**
- Structured logging with zap
- File-based output to `~/.promptstack/debug.log`
- Log rotation (10MB, keep last 3)
- Level control via environment variable
- Modal viewer for recent logs (Ctrl+L)

**Configuration:**
- `PROMPTSTACK_LOG_LEVEL`: DEBUG, INFO, WARN, ERROR (default: INFO)
- `PROMPTSTACK_DEBUG=1`: Enable DEBUG level

**Key Components:**
- `internal/logging` - Logger initialization and configuration
- `ui/logviewer` - Log viewer modal

See code-samples/039-logging-setup.sample.go for zap configuration.
See code-samples/040-log-viewer.sample.go for log viewer UI.

## Data Models and Relationships

### Core Data Models

**Prompt:**
```
- ID: string (generated)
- Title: string (required)
- Description: string (optional)
- Tags: []string (optional)
- Category: string (derived from folder)
- FilePath: string
- Content: string
- Metadata: map[string]interface{} (parsed frontmatter)
- Placeholders: []Placeholder
- ValidationStatus: ValidationResult
- UsageStats: UsageMetadata
```

**Composition:**
```
- ID: string (timestamp-based)
- FilePath: string
- Content: string
- CreatedAt: time.Time
- UpdatedAt: time.Time
- WorkingDirectory: string
- CharacterCount: int
- LineCount: int
- Placeholders: []Placeholder
- UndoStack: []UndoAction
- RedoStack: []UndoAction
- IsDirty: bool
```

**Placeholder:**
```
- Type: string ("text" | "list")
- Name: string
- StartPos: int
- EndPos: int
- CurrentValue: string | []string
- IsValid: bool
```

**Suggestion:**
```
- ID: string
- Type: SuggestionType (enum)
- Title: string
- Description: string
- ProposedChanges: []Edit
- Status: string ("pending" | "applying" | "applied" | "dismissed" | "error")
- Error: string (if status == error)
```

**LibraryIndex:**
```
- Prompts: map[string]IndexedPrompt
- LastBuilt: time.Time
- Version: string
```

**IndexedPrompt:**
```
- PromptID: string
- Title: string
- Tags: []string
- Category: string
- WordFrequency: map[string]int
- LastUsed: time.Time
- UseCount: int
```

See code-samples/041-data-models.sample.go for complete model definitions.

### Relationships

- **Library → Prompts** (1:N) - Library contains multiple prompts organized by folder
- **Composition → Placeholders** (1:N) - Composition contains zero or more placeholders
- **Composition → History** (1:1) - Each composition backed by one history entry
- **LibraryIndex → IndexedPrompts** (1:N) - Index tracks metadata for all prompts
- **Suggestions → Composition** (N:1) - Multiple suggestions may apply to one composition

See code-samples/042-relationship-diagram.sample.md for entity relationship diagram.

## State Management Approach

### Bubble Tea State Architecture

**Root Model:**
- Manages global application state
- Routes events to active screen
- Handles global keybindings (Esc, Ctrl+C, etc.)
- Manages modal stack (library browser, command palette, etc.)

**Screen Models:**
- Workspace (primary screen)
- Prompt Editor
- History Browser
- Settings Panel

**Modal Models:**
- Library Browser
- Command Palette
- File Reference Finder
- Validation Results
- Diff Viewer
- Log Viewer

**State Transitions:**
- Screens replace entire view
- Modals overlay current view
- Modal stack supports multiple overlays
- Esc dismisses top modal

**Shared State:**
- Library (immutable after load, reloaded on changes)
- Config (reloaded on settings changes)
- Vim Mode (global state)
- Active Composition (mutable)

See code-samples/043-root-model.sample.go for root model structure.
See code-samples/044-state-transitions.sample.go for screen/modal transitions.

### Concurrency Strategy

**Background Operations:**
- Library validation (startup)
- History database sync verification
- Auto-save debouncing
- AI API calls

**Concurrency Patterns:**
- Use channels for async operations
- Bubble Tea Cmd for background tasks
- Mutex protection for shared state (library index, undo stack)

See code-samples/045-concurrency-patterns.sample.go for async operation handling.

## Error Handling Strategy

### Error Handling Principles

1. **Never crash** - Catch all errors at appropriate boundaries
2. **Preserve user work** - Keep composition in memory on failures
3. **Clear context** - Show what operation failed and why
4. **Actionable** - Provide retry or alternative actions when possible
5. **Appropriate visibility** - Status bar for minor issues, modals for critical errors

### Error Handling by Category

**File System Errors:**
- Cannot read prompt: Load as plain text, warn in validation
- Cannot write composition: Retry auto-save, offer manual save, prevent data loss
- Cannot access config: Prompt for manual config creation or reset

**Database Errors:**
- Corruption detected: Prompt for rebuild with backup creation
- Query failures: Fall back to markdown file operations, suggest rebuild
- Write failures during auto-save: Retry silently, show persistent warning if retry fails

**API Errors:**
- Rate limit (429): Show retry button with API message in suggestions panel
- Auth failure (401): Show error, suggest checking settings
- Network/5xx: Show retry button, track retry count (max 3)
- Timeout: Show error with option to increase timeout in settings

**Parsing Errors:**
- Invalid YAML: Load prompt without metadata, show in validation results
- Malformed placeholders: Treat as literal text, continue normally
- Invalid markdown: Display as-is in preview mode

**Validation Errors:**
- Error-level: Block prompt insertion, show ⚠️ icon
- Warning-level: Allow insertion, show ⚠️ icon
- Show details on demand via validation results modal

See code-samples/046-error-handling-patterns.sample.go for error handling utilities.
See code-samples/047-error-display-components.sample.go for error UI components.

## Testing Strategy

### Testing Levels

**Unit Tests:**
- All business logic packages (`internal/*`)
- Pure functions (parsers, validators, algorithms)
- Target: 80%+ coverage for critical paths

**Integration Tests:**
- Library loading and indexing
- History database operations
- AI context selection
- Placeholder parsing and editing

**TUI Tests:**
- Use Bubble Tea's testing utilities
- Test model updates for user interactions
- Verify state transitions

**End-to-End Tests:**
- Manual testing scenarios
- Critical user workflows (create composition, insert prompt, apply suggestion)

### Test Organization

```
promptstack/
├── internal/
│   ├── library/
│   │   ├── loader.go
│   │   └── loader_test.go
│   ├── editor/
│   │   ├── undo.go
│   │   └── undo_test.go
│   └── ...
└── test/
    ├── fixtures/ (test prompts, compositions)
    ├── integration/ (integration tests)
    └── e2e/ (end-to-end scenarios)
```

### Key Test Scenarios

**Critical Paths:**
1. First launch initialization (config, starter prompts, database)
2. Load library with 100+ prompts
3. Create composition, insert prompts, fill placeholders
4. Auto-save and history persistence
5. AI suggestion request and diff application
6. Undo/redo operations
7. Database sync verification and rebuild

**Edge Cases:**
1. Large composition (approaching 1MB limit)
2. Prompt with duplicate placeholder names (validation)
3. Database corruption recovery
4. API rate limiting and retry
5. Terminal resize during operation
6. Concurrent auto-save and user editing

See code-samples/048-test-examples.sample.go for unit test examples.
See code-samples/049-integration-test-setup.sample.go for integration test utilities.

## Security Considerations

### API Key Storage

**Current Approach:**
- Plain text in `~/.promptstack/config.yaml`
- File permissions: 0600 (user read/write only)

**Risks:**
- API key visible to anyone with filesystem access
- No encryption at rest

**Mitigations:**
- Clear documentation about storage method
- Warn users not to commit config to version control
- Future: Consider OS keychain integration (macOS Keychain)

### Input Validation

**Prevent Directory Traversal:**
- Validate all file paths before reading/writing
- Ensure paths are within expected directories (`~/.promptstack/data/`)
- Use `filepath.Clean()` to normalize paths

**Sanitize User Input:**
- Escape special characters when generating markdown
- Validate frontmatter before parsing YAML
- Limit file sizes (1MB max for prompts and compositions)

### Prompt Injection

**AI Suggestions:**
- Treat all AI responses as untrusted content
- Never auto-apply suggestions without user review
- Display suggestions in read-only view before acceptance
- Diff view allows user to inspect exact changes

**User-Generated Prompts:**
- No execution of user prompts as code
- Library prompts are pure text/markdown
- No system command execution from prompt content

### Dependency Security

**Supply Chain:**
- Pin dependency versions in go.mod
- Review dependencies for known vulnerabilities
- Use only well-maintained libraries with active communities
- Minimize dependency count

See code-samples/050-input-validation.sample.go for validation utilities.

## Performance Considerations

### Performance Targets

**Startup:**
- Cold start: <500ms for 100 prompts
- Database verification: <200ms for 1000 history entries
- Library validation: Background, non-blocking

**Runtime:**
- Library search: <10ms for fuzzy filtering
- Auto-save: <50ms for typical composition
- History search: <100ms with SQLite FTS5
- UI rendering: 60fps (16ms per frame)

**AI Operations:**
- Context selection: <100ms for scoring 100+ prompts
- Suggestion display: Show immediately when received
- Diff application: <50ms for typical change set

### Optimization Strategies

**Library Loading:**
- Lazy load prompt content (load metadata first, content on demand)
- Cache parsed prompts in memory after first access
- Build index incrementally as prompts are loaded

**Fuzzy Search:**
- Use sahilm/fuzzy with pre-computed lowercase strings
- Limit result set to top 50 matches for display
- Debounce search input (50ms) to avoid excessive filtering

**Database Operations:**
- Use prepared statements for repeated queries
- Batch inserts/updates when possible
- FTS5 index for full-text search (optimized for read-heavy workload)
- SQLite WAL mode for better concurrent performance

**Memory Management:**
- Limit undo history to 100 actions
- Clear dismissed suggestions after session
- Release large buffers after operations complete

**TUI Rendering:**
- Use Bubble Tea's efficient diff-based rendering
- Minimize re-renders by memoizing expensive calculations
- Lazy render off-screen content (viewport pattern)

### Resource Limits

**File Sizes:**
- Prompts: 1MB max
- Compositions: 1MB max
- Referenced files: 1MB max (warning, not hard limit)

**Memory:**
- Target: <100MB resident for typical usage
- Library: ~10KB per prompt in memory
- Undo stack: ~1KB per action

**Database:**
- History database: No hard limit, user-managed via cleanup
- FTS5 index: Scales linearly with content size

See code-samples/051-performance-benchmarks.sample.go for benchmark tests.

## Implementation Phases

### Phase 1: Foundation (MVP Core)

**Goal:** Basic composition with library management

**Deliverables:**
1. Application bootstrap and configuration
2. Library loader and in-memory index
3. Composition workspace with basic editing
4. Auto-save to history
5. Library browser with fuzzy search
6. Command palette
7. Basic error handling

**Duration Estimate:** Foundation for all subsequent features

### Phase 2: Advanced Editing

**Goal:** Placeholder system and undo/redo

**Deliverables:**
1. Placeholder parser and validator
2. Placeholder navigation (Tab/Shift+Tab)
3. Text placeholder editing (vim-style)
4. List placeholder editing (vim-style)
5. Undo/redo with smart batching
6. Composition state persistence

**Duration Estimate:** Core editing experience

### Phase 3: History & File References

**Goal:** History management and file linking

**Deliverables:**
1. SQLite database setup with FTS5
2. History browser with search
3. History cleanup functionality
4. File reference system (fuzzy finder)
5. Title extraction from frontmatter
6. Batch title editor
7. Database sync verification

**Duration Estimate:** Complete history feature set

### Phase 4: Prompt Management

**Goal:** Create, edit, and validate prompts

**Deliverables:**
1. Prompt creation wizard
2. Prompt editor with edit/preview mode
3. Glamour integration for preview
4. Library validation checks
5. Validation results display
6. Invalid prompt handling in browser

**Duration Estimate:** Full prompt lifecycle

### Phase 5: AI Integration

**Goal:** AI suggestions with conservative token usage

**Deliverables:**
1. Claude API client integration
2. Context selection algorithm
3. Token estimation and budget enforcement
4. Suggestion parsing and display
5. Diff generation and preview
6. Diff application with undo integration
7. Error handling and retry logic

**Duration Estimate:** Complete AI feature set

### Phase 6: Vim Mode

**Goal:** Universal vim keybindings

**Deliverables:**
1. Vim state machine
2. Keybinding maps for each context
3. Mode indicators
4. Context-aware event routing
5. Vim support in all components

**Duration Estimate:** Full vim mode support

### Phase 7: Polish & Settings

**Goal:** Final UX improvements

**Deliverables:**
1. Settings panel
2. Status bar with all indicators
3. Responsive layout (narrow terminal support)
4. Split-pane resizing
5. Log viewer modal
6. Comprehensive error messages
7. Startup validation

**Duration Estimate:** Production-ready polish

### Phase 8: Testing & Documentation

**Goal:** Comprehensive tests and docs

**Deliverables:**
1. Unit tests for critical paths
2. Integration tests
3. E2E test scenarios
4. User documentation
5. Developer documentation
6. Build and release process

**Duration Estimate:** Quality assurance and launch prep

## Build and Distribution

### Build Process

**Build Tool:** Standard Go toolchain

**Build Steps:**
1. Embed starter prompts: `//go:embed starter-prompts`
2. Compile for macOS Intel: `GOOS=darwin GOARCH=amd64 go build`
3. Compile for macOS ARM: `GOOS=darwin GOARCH=arm64 go build`
4. Version tagging: Use git tags for version management

**Binary Naming:**
- `promptstack-darwin-amd64` (Intel Macs)
- `promptstack-darwin-arm64` (Apple Silicon)

### Distribution

**Method:** GitHub Releases

**Release Artifacts:**
1. Pre-built binaries (Intel + ARM)
2. README.md with installation instructions
3. CHANGELOG.md
4. LICENSE

**Installation:**
```bash
# Download from GitHub Releases
chmod +x promptstack-darwin-arm64
mv promptstack-darwin-arm64 /usr/local/bin/promptstack
promptstack
```

### Version Management

**Config Version Tracking:**
- `version` field in config.yaml
- Compare on startup to detect upgrades
- Extract new starter prompts only on version mismatch

**Upgrade Process:**
1. User downloads new binary
2. Replaces old binary
3. Runs application
4. Tool detects version mismatch
5. Extracts new starter prompts (skip existing)
6. Updates version in config
7. Normal operation

See code-samples/052-build-script.sample.sh for build automation.
See code-samples/053-version-upgrade.sample.go for upgrade handling.

## Appendix: Implementation Checklist

### Configuration & Bootstrap
- [ ] Config structure and YAML parsing
- [ ] First-run interactive setup wizard
- [ ] Version-aware starter prompt extraction
- [ ] Database initialization
- [ ] Logging setup with zap

### Library Management
- [ ] Prompt model and metadata types
- [ ] Library loader (filesystem scanning)
- [ ] YAML frontmatter parsing
- [ ] In-memory index builder
- [ ] Index scoring algorithm
- [ ] Library validation checks
- [ ] Validation results display

### Composition Workspace
- [ ] Basic text editor model
- [ ] Cursor management
- [ ] Auto-save with debouncing
- [ ] Split-pane layout
- [ ] Divider resizing
- [ ] Responsive layout (narrow terminal)
- [ ] Status bar with indicators

### Placeholder System
- [ ] Placeholder parser (regex-based)
- [ ] Placeholder validator
- [ ] Tab/Shift+Tab navigation
- [ ] Text placeholder editing (vim-style)
- [ ] List placeholder editing UI
- [ ] List item CRUD operations

### Undo/Redo
- [ ] Undo stack data structure
- [ ] Smart action batching
- [ ] Undo/redo commands
- [ ] Visual feedback for boundaries
- [ ] Integration with placeholder edits

### History Management
- [ ] Timestamped markdown file creation
- [ ] SQLite database schema with FTS5
- [ ] Auto-save to markdown + SQLite
- [ ] History browser UI
- [ ] Fuzzy search across history
- [ ] Database sync verification
- [ ] Rebuild with backup
- [ ] History cleanup strategies

### Library Browser
- [ ] Flat list view with category labels
- [ ] Fuzzy search integration
- [ ] Preview pane
- [ ] Keyboard navigation
- [ ] Insert at cursor / new line

### Command Palette
- [ ] Modal overlay UI
- [ ] Command registry
- [ ] Fuzzy command filtering
- [ ] Command execution

### File Reference System
- [ ] File traversal with .gitignore respect
- [ ] Fuzzy file finder UI
- [ ] Multiple file selection
- [ ] Title extraction from frontmatter
- [ ] Batch title editor
- [ ] Markdown link insertion

### AI Suggestions
- [ ] Claude API client wrapper
- [ ] Context selection algorithm
- [ ] Token estimation
- [ ] Token budget enforcement
- [ ] Suggestion parsing (6 types)
- [ ] Suggestions panel UI
- [ ] Diff generation
- [ ] Diff viewer modal (unified format)
- [ ] Diff application
- [ ] Integration with undo system
- [ ] Error handling and retry

### Prompt Creation & Editing
- [ ] Creation wizard (guided steps)
- [ ] Prompt editor UI
- [ ] Edit/preview mode toggle
- [ ] Glamour markdown preview
- [ ] Explicit save (no auto-save)

### Settings Panel
- [ ] Settings modal UI
- [ ] Vim mode toggle
- [ ] API key input (masked)
- [ ] Model selection
- [ ] Immediate persistence

### Vim Mode
- [ ] Vim state machine (Normal/Insert/Visual)
- [ ] Keybinding maps per mode
- [ ] Context-aware event routing
- [ ] Mode indicators
- [ ] Universal vim support (all components)

### Logging
- [ ] Zap logger initialization
- [ ] Log rotation setup
- [ ] Environment variable configuration
- [ ] Log viewer modal (Ctrl+L)

### Error Handling
- [ ] Error handling utilities
- [ ] Error display components
- [ ] Retry mechanisms
- [ ] Graceful degradation

### Testing
- [ ] Unit test framework setup
- [ ] Test fixtures (prompts, compositions)
- [ ] Unit tests for critical paths
- [ ] Integration tests
- [ ] E2E test scenarios

### Build & Distribution
- [ ] Build script (Intel + ARM)
- [ ] Starter prompts embedding
- [ ] Version management
- [ ] Upgrade handling
- [ ] GitHub release automation

---

**Last Updated:** 2026-01-06
**Status:** Draft - Ready for review and code sample creation
