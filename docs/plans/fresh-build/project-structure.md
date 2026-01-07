# PromptStack - Go-Idiomatic Project Structure

## Domain-Driven Design Analysis

Based on the implementation plan, PromptStack has **8 core domains**:

### 1. **Editor Domain**
Core text editing functionality independent of prompts or AI.
- Text buffer management
- Cursor positioning and navigation
- Undo/redo stack
- Selection handling

### 2. **Prompt Domain**
Everything related to prompt templates and placeholders.
- Prompt models and metadata
- Placeholder parsing and validation
- Prompt file I/O with YAML frontmatter

### 3. **Library Domain**
Managing collections of prompts.
- Prompt indexing and scoring
- Fuzzy search
- Library validation
- Usage tracking

### 4. **History Domain**
Composition persistence and retrieval.
- SQLite database operations
- Markdown file storage
- Full-text search
- Cleanup strategies

### 5. **AI Domain**
Claude API integration and suggestion system.
- API client wrapper
- Context selection algorithm
- Suggestion parsing
- Diff generation and application

### 6. **Config Domain**
Application configuration.
- Config loading/saving
- First-run setup
- Settings management

### 7. **UI Domain**
Bubble Tea TUI components.
- Screen models (workspace, history, settings)
- Modal models (palette, browser, diff viewer)
- Shared UI components

### 8. **Platform Domain**
Cross-cutting infrastructure.
- Logging
- Error handling
- File system utilities
- Bootstrap orchestration

---

## Recommended Project Structure

```
promptstack/
├── cmd/
│   └── promptstack/
│       ├── main.go                    # Entry point
│       ├── embed.go                   # Embedded starter prompts
│       └── starter-prompts/           # Embedded prompt files
│           ├── commands/
│           ├── rules/
│           └── workflows/
│
├── internal/
│   ├── editor/                        # EDITOR DOMAIN
│   │   ├── buffer.go                  # Text buffer with cursor
│   │   ├── buffer_test.go
│   │   ├── selection.go               # Text selection logic
│   │   ├── selection_test.go
│   │   ├── undo.go                    # Undo/redo stack
│   │   └── undo_test.go
│   │
│   ├── prompt/                        # PROMPT DOMAIN
│   │   ├── prompt.go                  # Prompt model
│   │   ├── prompt_test.go
│   │   ├── metadata.go                # YAML frontmatter types
│   │   ├── metadata_test.go
│   │   ├── placeholder.go             # Placeholder parsing
│   │   ├── placeholder_test.go
│   │   ├── validator.go               # Prompt validation
│   │   ├── validator_test.go
│   │   ├── storage.go                 # Prompt file I/O
│   │   └── storage_test.go
│   │
│   ├── library/                       # LIBRARY DOMAIN
│   │   ├── library.go                 # Library manager
│   │   ├── library_test.go
│   │   ├── loader.go                  # Load prompts from filesystem
│   │   ├── loader_test.go
│   │   ├── index.go                   # In-memory indexing
│   │   ├── index_test.go
│   │   ├── scorer.go                  # Context relevance scoring
│   │   ├── scorer_test.go
│   │   ├── search.go                  # Fuzzy search
│   │   └── search_test.go
│   │
│   ├── history/                       # HISTORY DOMAIN
│   │   ├── manager.go                 # History manager
│   │   ├── manager_test.go
│   │   ├── database.go                # SQLite operations
│   │   ├── database_test.go
│   │   ├── storage.go                 # Markdown file operations
│   │   ├── storage_test.go
│   │   ├── sync.go                    # DB/file sync verification
│   │   ├── sync_test.go
│   │   ├── search.go                  # Full-text search
│   │   ├── search_test.go
│   │   ├── cleanup.go                 # History cleanup strategies
│   │   └── cleanup_test.go
│   │
│   ├── ai/                            # AI DOMAIN
│   │   ├── client.go                  # Claude API client
│   │   ├── client_test.go
│   │   ├── context.go                 # Context selection
│   │   ├── context_test.go
│   │   ├── tokens.go                  # Token estimation & budgeting
│   │   ├── tokens_test.go
│   │   ├── suggestions.go             # Suggestion types & parsing
│   │   ├── suggestions_test.go
│   │   ├── diff.go                    # Diff generation & application
│   │   └── diff_test.go
│   │
│   ├── config/                        # CONFIG DOMAIN
│   │   ├── config.go                  # Config types & loading
│   │   ├── config_test.go
│   │   ├── setup.go                   # First-run wizard
│   │   ├── setup_test.go
│   │   ├── settings.go                # Settings updates
│   │   └── settings_test.go
│   │
│   ├── platform/                      # PLATFORM DOMAIN
│   │   ├── bootstrap/
│   │   │   ├── bootstrap.go           # App initialization
│   │   │   ├── bootstrap_test.go
│   │   │   ├── starter.go             # Extract starter prompts
│   │   │   └── starter_test.go
│   │   ├── logging/
│   │   │   ├── logger.go              # Zap logger setup
│   │   │   └── logger_test.go
│   │   ├── errors/
│   │   │   ├── errors.go              # Error types
│   │   │   ├── handler.go             # Error handling utilities
│   │   │   └── errors_test.go
│   │   └── files/
│   │       ├── traversal.go           # File system traversal
│   │       ├── traversal_test.go
│   │       ├── gitignore.go           # .gitignore handling
│   │       └── gitignore_test.go
│   │
│   ├── vim/                           # VIM MODE (cross-cutting)
│   │   ├── state.go                   # Vim state machine
│   │   ├── state_test.go
│   │   ├── keymaps.go                 # Mode keybindings
│   │   └── keymaps_test.go
│   │
│   └── commands/                      # COMMAND SYSTEM (cross-cutting)
│       ├── registry.go                # Command registry
│       ├── registry_test.go
│       ├── core.go                    # Core commands
│       └── core_test.go
│
├── ui/                                # UI DOMAIN (Bubble Tea)
│   ├── app/
│   │   ├── model.go                   # Root application model
│   │   ├── update.go                  # Root update handler
│   │   ├── view.go                    # Root view renderer
│   │   └── messages.go                # Global message types
│   │
│   ├── workspace/                     # Primary editing screen
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   ├── messages.go
│   │   └── keybindings.go
│   │
│   ├── browser/                       # Library browser modal
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   └── messages.go
│   │
│   ├── palette/                       # Command palette modal
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   └── messages.go
│   │
│   ├── history/                       # History browser screen
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   └── messages.go
│   │
│   ├── prompteditor/                  # Prompt creation/editing screen
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   └── messages.go
│   │
│   ├── settings/                      # Settings modal
│   │   ├── model.go
│   │   ├── update.go
│   │   └── view.go
│   │
│   ├── suggestions/                   # AI suggestions panel
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   └── messages.go
│   │
│   ├── diffviewer/                    # Diff viewer modal
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   └── messages.go
│   │
│   ├── filereference/                 # File reference picker modal
│   │   ├── model.go
│   │   ├── update.go
│   │   └── view.go
│   │
│   ├── validation/                    # Validation results modal
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   └── messages.go
│   │
│   ├── cleanup/                       # History cleanup modal
│   │   ├── model.go
│   │   ├── update.go
│   │   └── view.go
│   │
│   ├── logviewer/                     # Log viewer modal
│   │   ├── model.go
│   │   ├── update.go
│   │   └── view.go
│   │
│   ├── statusbar/                     # Status bar component
│   │   ├── model.go
│   │   └── view.go
│   │
│   ├── theme/                         # Theme system
│   │   ├── theme.go                   # Lipgloss styles
│   │   └── README.md
│   │
│   └── common/                        # Shared UI components
│       ├── modal.go                   # Modal overlay wrapper
│       ├── confirmation.go            # Confirmation dialog
│       ├── fuzzylist.go               # Reusable fuzzy list
│       └── markdownpreview.go         # Glamour markdown renderer
│
├── test/
│   ├── fixtures/                      # Test data
│   │   ├── prompts/                   # Sample prompt files
│   │   ├── compositions/              # Sample compositions
│   │   └── configs/                   # Sample configs
│   ├── integration/                   # Integration tests
│   │   ├── library_test.go
│   │   ├── history_test.go
│   │   └── ai_test.go
│   └── testutil/                      # Test utilities
│       ├── helpers.go
│       └── mocks.go
│
├── docs/
│   ├── plans/
│   │   ├── fresh-build/
│   │   │   ├── tracking.md
│   │   │   ├── project-structure.md   # This file
│   │   │   └── implementation-plan.md # To be created
│   │   └── initial-build-archive/     # Archived initial plan
│   └── architecture/
│       ├── domain-model.md
│       └── data-flow.md
│
├── archive/
│   └── code/                          # Previous implementation
│
├── issues/                            # Issue tracking
│   └── 001-spacebar-issue/
│
├── go.mod
├── go.sum
├── README.md
├── .gitignore
└── Makefile                           # Build automation
```

---

## Design Principles

### 1. **Domain Separation**

Each domain package is **self-contained** and has clear responsibilities:

- **`internal/editor/`** - Pure text editing logic, no knowledge of prompts or AI
- **`internal/prompt/`** - Prompt models and operations, no knowledge of library or history
- **`internal/library/`** - Library management, depends on `prompt` package
- **`internal/history/`** - Composition storage, minimal dependencies
- **`internal/ai/`** - AI integration, depends on `prompt` for context selection
- **`internal/config/`** - Configuration, no business logic dependencies
- **`internal/platform/`** - Infrastructure, no business logic dependencies

### 2. **Dependency Direction**

Dependencies flow **downward** and **inward**:

```
ui/           →  internal/{domain}  →  internal/platform
(presentation)   (business logic)      (infrastructure)
```

**Rules:**
- UI packages can import domain packages
- Domain packages can import platform packages
- Domain packages **should not** import UI packages
- Platform packages **should not** import domain or UI packages

**Example Valid Imports:**
- `ui/workspace` → `internal/editor`, `internal/prompt`
- `internal/library` → `internal/prompt`, `internal/platform/files`
- `internal/ai` → `internal/prompt`, `internal/library`

**Example Invalid Imports:**
- `internal/editor` → `ui/workspace` ❌
- `internal/platform/logging` → `internal/library` ❌

### 3. **Package Naming**

Follow Go conventions:
- **Singular names**: `editor`, `prompt`, `library` (not `editors`, `prompts`, `libraries`)
- **No stuttering**: Types are `editor.Buffer`, not `editor.EditorBuffer`
- **Clear purpose**: Package name should describe what it does

### 4. **File Organization Within Packages**

**Standard Pattern:**
```
package-name/
├── {entity}.go           # Core types and constructors
├── {entity}_test.go      # Unit tests
├── {operation}.go        # Operations on entities
├── {operation}_test.go   # Operation tests
└── messages.go           # Package-specific message types (UI only)
```

**Example: `internal/prompt/`**
```
prompt/
├── prompt.go             # Prompt type, NewPrompt()
├── prompt_test.go
├── metadata.go           # Metadata types
├── metadata_test.go
├── placeholder.go        # Placeholder parsing
├── placeholder_test.go
├── validator.go          # Validation logic
├── validator_test.go
├── storage.go            # File I/O
└── storage_test.go
```

**UI Package Pattern:**
```
workspace/
├── model.go              # Bubble Tea model type
├── update.go             # Update() function
├── view.go               # View() function
├── messages.go           # Bubble Tea messages
└── keybindings.go        # Key handlers (optional)
```

### 5. **Testing Structure**

**Unit Tests:**
- Co-located with source: `buffer.go` → `buffer_test.go`
- Test package: `package editor_test` (black-box testing preferred)
- Focus on **public API**, not internal implementation

**Integration Tests:**
- Located in `test/integration/`
- Test interactions between domains
- Use real dependencies where possible

**Test Fixtures:**
- Shared test data in `test/fixtures/`
- Helper functions in `test/testutil/`

### 6. **Interface Design**

Define interfaces **where they're used**, not where they're implemented:

**Good:**
```go
// ui/workspace/model.go
type EditorService interface {
    GetContent() string
    SetContent(string) error
    Undo() error
    Redo() error
}

// internal/editor/buffer.go
type Buffer struct { ... }

func (b *Buffer) GetContent() string { ... }
// Buffer implements EditorService implicitly
```

**Bad:**
```go
// internal/editor/buffer.go
type EditorService interface { ... } // ❌ Defining interface in implementation package
```

This enables **dependency inversion** and easier testing.

### 7. **Error Handling**

**Domain packages** return errors:
```go
// internal/library/loader.go
func (l *Loader) LoadPrompts(dir string) ([]prompt.Prompt, error) {
    // Return descriptive errors
    return nil, fmt.Errorf("failed to load prompts from %s: %w", dir, err)
}
```

**Platform error package** provides utilities:
```go
// internal/platform/errors/errors.go
type Error struct {
    Op      string // Operation that failed
    Kind    Kind   // Error category
    Err     error  // Underlying error
}

func (e *Error) Error() string
func Is(err error, kind Kind) bool
```

**UI packages** handle errors for display:
```go
// ui/workspace/update.go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case errMsg:
        m.err = msg.error
        // Display error in status bar or modal
    }
}
```

### 8. **Concurrency Patterns**

**Use Bubble Tea Cmd for async operations:**
```go
// ui/workspace/update.go
func loadLibrary(path string) tea.Cmd {
    return func() tea.Msg {
        prompts, err := library.Load(path)
        if err != nil {
            return libraryLoadedMsg{err: err}
        }
        return libraryLoadedMsg{prompts: prompts}
    }
}
```

**Domain packages are synchronous** (no goroutines in business logic):
```go
// internal/library/loader.go
func (l *Loader) LoadPrompts(dir string) ([]prompt.Prompt, error) {
    // Synchronous operation
    // Let caller decide concurrency model
}
```

### 9. **Configuration Injection**

**Pass config down, don't access globally:**

**Good:**
```go
// cmd/promptstack/main.go
func main() {
    cfg := config.Load()
    
    lib := library.New(library.Config{
        DataDir: cfg.DataDir,
    })
    
    app := ui.NewApp(cfg, lib)
    app.Run()
}
```

**Bad:**
```go
// internal/library/loader.go
func (l *Loader) LoadPrompts() {
    dir := config.Global.DataDir // ❌ Global state
}
```

### 10. **Mocking for Tests**

**Define minimal interfaces:**
```go
// ui/workspace/model_test.go
type mockLibrary struct {
    prompts []prompt.Prompt
    err     error
}

func (m *mockLibrary) Search(query string) ([]prompt.Prompt, error) {
    return m.prompts, m.err
}
```

**Use dependency injection:**
```go
// ui/workspace/model.go
type Model struct {
    library LibrarySearcher // Interface, not concrete type
}

type LibrarySearcher interface {
    Search(query string) ([]prompt.Prompt, error)
}
```

---

## Migration Strategy from Archive

### Phase 1: Platform Setup
1. Create `internal/platform/` packages (bootstrap, logging, errors, files)
2. Extract and adapt from `archive/code/internal/bootstrap/`, `archive/code/internal/logging/`, etc.
3. Write tests for each platform package

### Phase 2: Core Domains
1. Create `internal/editor/` with clean text editing logic
2. Create `internal/prompt/` with prompt models and parsing
3. Create `internal/config/` with configuration management
4. Write comprehensive tests for each domain

### Phase 3: Secondary Domains
1. Create `internal/library/` using `internal/prompt/`
2. Create `internal/history/` for persistence
3. Create `internal/ai/` for Claude integration
4. Write comprehensive tests for each domain

### Phase 4: UI Layer
1. Create `ui/app/` root model
2. Create `ui/workspace/` primary screen
3. Create other UI components as needed
4. Write UI tests using Bubble Tea test utilities

### Phase 5: Integration
1. Wire up all components in `cmd/promptstack/main.go`
2. Write integration tests in `test/integration/`
3. Manual testing with user feedback checkpoints

---

## Key Differences from Archive

### What's Better

1. **Clear Domain Boundaries**: Each domain is self-contained with explicit dependencies
2. **Platform Layer**: Infrastructure code separated from business logic
3. **Testable Design**: Interfaces defined at usage site, dependency injection
4. **Standard Patterns**: Consistent file naming and organization
5. **UI Separation**: Business logic never in UI components

### What's Removed

1. **No `internal/embed/`**: Move to `cmd/promptstack/embed.go` (closer to `go:embed` directive)
2. **No `internal/setup/`**: Move to `internal/config/setup.go` (setup is configuration)
3. **No scattered files**: Everything has a clear domain home

### What's New

1. **`internal/platform/`**: Infrastructure package for cross-cutting concerns
2. **`test/` directory**: Centralized integration tests and fixtures
3. **Consistent UI patterns**: `model.go`, `update.go`, `view.go`, `messages.go`

---

## Go-Idiomatic Conventions

### 1. **Package Comments**
```go
// Package editor provides core text editing functionality including
// cursor management, text buffer operations, and undo/redo support.
package editor
```

### 2. **Type Constructors**
```go
// New creates a new Buffer with the given content.
func New(content string) *Buffer {
    return &Buffer{
        content: content,
        cursor:  0,
    }
}
```

### 3. **Method Receivers**
```go
// Use pointer receivers for mutating methods
func (b *Buffer) Insert(text string) error { ... }

// Use value receivers for read-only methods
func (b Buffer) Content() string { ... }

// Be consistent: if any method uses pointer, all should
```

### 4. **Error Messages**
```go
// Lower case, no punctuation
return fmt.Errorf("failed to load prompt: %w", err)

// Include context
return fmt.Errorf("failed to parse frontmatter in %s: %w", path, err)
```

### 5. **Exported vs Unexported**
```go
// Exported (public API)
type Buffer struct {
    // Unexported fields (implementation details)
    content string
    cursor  int
}

// Exported methods
func (b *Buffer) Content() string { ... }

// Unexported helpers
func (b *Buffer) validateCursor() error { ... }
```

---

## Theme System

### Overview

The theme system provides a **centralized source of truth** for all UI styling using the **Catppuccin Mocha** color palette. This ensures visual consistency across all components and makes it easy to update the entire UI appearance.

### Location

```
ui/theme/
├── theme.go      # Color constants and style helper functions
└── README.md     # Usage documentation and examples
```

### Design Philosophy

1. **Single Source of Truth**: All colors defined once in [`theme.go`](ui/theme/theme.go)
2. **Semantic Naming**: Colors named by purpose, not appearance
3. **Helper Functions**: Pre-built styles for common UI patterns
4. **Composable**: Base styles can be extended for custom needs
5. **Type Safe**: Constants prevent typos in color values

### Color Constants

**Background Colors:**
```go
const (
    BackgroundPrimary   = "#1e1e2e"  // Main modal background
    BackgroundSecondary = "#181825"  // Status bar background
    BackgroundTertiary  = "#313244"  // Secondary button background
    BackgroundInput     = "236"      // Input field background (terminal color)
)
```

**Foreground Colors:**
```go
const (
    ForegroundPrimary   = "#cdd6f4"  // Main text
    ForegroundSecondary = "#a6adc8"  // Secondary text
    ForegroundMuted     = "#6c7086"  // Muted text
    ForegroundWhite     = "15"       // White text (terminal color)
)
```

**Accent Colors:**
```go
const (
    AccentBlue   = "#89b4fa"  // Primary accent, info states
    AccentGreen  = "#a6e3a1"  // Success states
    AccentYellow = "#f9e2af"  // Warning states
    AccentRed    = "#f38ba8"  // Error states
    AccentCyan   = "39"       // Category labels (terminal color)
)
```

**Border Colors:**
```go
const (
    BorderPrimary = "#45475a"  // Modal borders
    BorderMuted   = "240"      // Secondary borders (terminal color)
)
```

**Special Colors:**
```go
const (
    CursorBackground = "7"    // Cursor highlight
    CursorForeground = "0"    // Cursor text
    TextMuted        = "245"  // Dimmed text (terminal color)
)
```

### Style Helper Functions

The theme provides pre-built style functions for common UI patterns:

**Modal Styles:**
- [`ModalStyle()`](ui/theme/theme.go:49) - Base modal container with rounded border
- [`ModalTitleStyle()`](ui/theme/theme.go:58) - Bold modal titles
- [`ModalContentStyle()`](ui/theme/theme.go:66) - Modal content text
- [`ModalButtonStyle()`](ui/theme/theme.go:73) - Primary buttons (normal state)
- [`ModalButtonFocusedStyle()`](ui/theme/theme.go:82) - Primary buttons (focused state)
- [`ModalButtonSecondaryStyle()`](ui/theme/theme.go:91) - Secondary buttons (normal state)
- [`ModalButtonSecondaryFocusedStyle()`](ui/theme/theme.go:100) - Secondary buttons (focused state)

**Status Bar Styles:**
- [`StatusStyle()`](ui/theme/theme.go:111) - Normal status messages
- [`InfoStyle()`](ui/theme/theme.go:119) - Info messages (blue)
- [`SuccessStyle()`](ui/theme/theme.go:127) - Success messages (green)
- [`WarningStyle()`](ui/theme/theme.go:135) - Warning messages (yellow)
- [`ErrorStyle()`](ui/theme/theme.go:143) - Error messages (red)
- [`SeparatorStyle()`](ui/theme/theme.go:151) - Status bar separators

**Input Styles:**
- [`InputStyle()`](ui/theme/theme.go:160) - Base input field
- [`SearchInputStyle()`](ui/theme/theme.go:168) - Search input field

**List Styles:**
- [`ListItemStyle()`](ui/theme/theme.go:175) - Unselected list items
- [`ListItemSelectedStyle()`](ui/theme/theme.go:181) - Selected list items
- [`ListCategoryStyle()`](ui/theme/theme.go:188) - Category labels (cyan, bold)
- [`ListDescriptionStyle()`](ui/theme/theme.go:195) - Item descriptions (muted, italic)
- [`ListEmptyStyle()`](ui/theme/theme.go:202) - Empty state message

**Preview Styles:**
- [`PreviewStyle()`](ui/theme/theme.go:211) - Preview pane container
- [`PreviewTitleStyle()`](ui/theme/theme.go:219) - Preview titles
- [`PreviewDescriptionStyle()`](ui/theme/theme.go:226) - Preview descriptions
- [`PreviewTagsStyle()`](ui/theme/theme.go:233) - Preview tags
- [`PreviewContentStyle()`](ui/theme/theme.go:239) - Preview content

**Cursor Styles:**
- [`CursorStyle()`](ui/theme/theme.go:247) - Cursor highlighting
- [`ActivePlaceholderStyle()`](ui/theme/theme.go:254) - Active placeholder highlighting

**Validation Styles:**
- [`ValidationErrorStyle()`](ui/theme/theme.go:298) - Validation errors (red, bold)
- [`ValidationWarningStyle()`](ui/theme/theme.go:305) - Validation warnings (yellow)
- [`ValidationIconStyle()`](ui/theme/theme.go:311) - Validation icons
- [`ValidationPromptStyle()`](ui/theme/theme.go:317) - Prompt titles in validation list
- [`ValidationPromptErrorStyle()`](ui/theme/theme.go:323) - Prompts with errors
- [`ValidationPromptWarningStyle()`](ui/theme/theme.go:330) - Prompts with warnings
- [`ValidationDetailStyle()`](ui/theme/theme.go:336) - Validation detail messages
- [`ValidationSummaryStyle()`](ui/theme/theme.go:343) - Validation summary

**Diff Styles:**
- [`DiffHunkHeaderStyle()`](ui/theme/theme.go:353) - Diff hunk headers (cyan, bold)
- [`DiffContextStyle()`](ui/theme/theme.go:360) - Diff context lines (muted)
- [`DiffAdditionStyle()`](ui/theme/theme.go:366) - Diff addition lines (green, bold)
- [`DiffDeletionStyle()`](ui/theme/theme.go:373) - Diff deletion lines (red, bold)

**Other Styles:**
- [`HeaderStyle()`](ui/theme/theme.go:264) - Section headers (bold)

### Base Styles for Composition

For custom components, extend base styles:

```go
// Base styles that can be customized
func BaseModal() lipgloss.Style
func BaseInput() lipgloss.Style
func BaseListItem() lipgloss.Style
```

### Usage Examples

**Using Style Helpers:**
```go
import "github.com/yourorg/promptstack/ui/theme"

// Simple usage - just call the helper
modalStyle := theme.ModalStyle()
titleStyle := theme.ModalTitleStyle()
buttonStyle := theme.ModalButtonFocusedStyle()
```

**Using Color Constants:**
```go
import (
    "github.com/charmbracelet/lipgloss"
    "github.com/yourorg/promptstack/ui/theme"
)

// Custom style using theme colors
customStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color(theme.ForegroundPrimary)).
    Background(lipgloss.Color(theme.BackgroundPrimary)).
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color(theme.BorderPrimary))
```

**Extending Base Styles:**
```go
// Start with base and customize
customModal := theme.BaseModal().
    Width(80).
    Height(40).
    Padding(2, 3)

customInput := theme.BaseInput().
    Width(60).
    Placeholder("Enter text...")
```

### Theme Guidelines

**DO:**
- ✅ Always use theme helpers for standard UI patterns
- ✅ Use color constants for custom styles
- ✅ Extend base styles for component-specific needs
- ✅ Document new style helpers if you add them
- ✅ Test visual appearance after theme changes

**DON'T:**
- ❌ Hard-code hex colors or terminal color numbers
- ❌ Create duplicate style definitions
- ❌ Mix theme colors with hard-coded colors
- ❌ Use colors for non-semantic purposes (e.g., AccentRed for non-errors)

### Adding New Styles

When creating new UI components:

1. **Check existing helpers** - Can you use an existing style?
2. **Use color constants** - Build custom styles with theme colors
3. **Add new helpers** - If pattern is reusable, add to [`theme.go`](ui/theme/theme.go):
   ```go
   // NewComponentStyle returns the style for new component
   func NewComponentStyle() lipgloss.Style {
       return lipgloss.NewStyle().
           Foreground(lipgloss.Color(ForegroundPrimary)).
           Background(lipgloss.Color(BackgroundPrimary))
   }
   ```
4. **Document it** - Update [`ui/theme/README.md`](ui/theme/README.md) with usage examples
5. **Follow naming** - Use pattern: `ComponentStyle()` or `ComponentStateStyle()`

### Testing Theme Changes

When modifying the theme:

1. **Visual regression check** - Test all UI components
2. **Contrast verification** - Ensure text is readable
3. **State testing** - Verify focused/selected/error states
4. **Terminal compatibility** - Test in different terminal emulators
5. **Color blindness** - Consider accessibility

### Future Enhancements

Potential improvements:

- **Multiple Themes**: Support light/dark mode switching
- **Theme Configuration**: Load theme from config file
- **Dynamic Theming**: Runtime theme switching
- **Custom Themes**: User-defined color schemes
- **Accessibility**: High-contrast mode options
- **Theme Validation**: Ensure colors meet contrast ratios

### Catppuccin Mocha Palette

The theme uses **Catppuccin Mocha**, a popular dark theme with excellent contrast:

| Purpose | Color | Hex | Usage |
|---------|-------|-----|-------|
| Background Primary | Dark Blue | #1e1e2e | Modal backgrounds |
| Background Secondary | Darker Blue | #181825 | Status bar |
| Background Tertiary | Gray Blue | #313244 | Secondary buttons |
| Foreground Primary | White | #cdd6f4 | Main text |
| Foreground Secondary | Light Gray | #a6adc8 | Secondary text |
| Foreground Muted | Muted Gray | #6c7086 | Dimmed text |
| Accent Blue | Blue | #89b4fa | Primary actions, info |
| Accent Green | Green | #a6e3a1 | Success states |
| Accent Yellow | Yellow | #f9e2af | Warning states |
| Accent Red | Red | #f38ba8 | Error states |
| Border Primary | Gray | #45475a | Modal borders |

**Why Catppuccin Mocha?**
- Excellent contrast and readability
- Consistent color relationships
- Popular and well-tested
- Beautiful, modern aesthetic
- Good terminal compatibility

---

## Next Steps

1. **Review this structure** with the implementation team
2. **Create implementation plan** based on this structure
3. **Set up project scaffolding** (directories, go.mod, Makefile)
4. **Begin Phase 1** (Platform packages) with TDD approach

---

**Last Updated**: 2026-01-07  
**Status**: Proposed - Ready for review and refinement