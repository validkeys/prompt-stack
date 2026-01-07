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
// ✅ GOOD: Small, focused interfaces
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