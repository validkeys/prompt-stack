# Key Learnings - PromptStack Implementation

This document captures key learnings and insights from implementing PromptStack to assist future development.

## Go Embed Limitations

**Issue**: `go:embed` does not support parent directory references (`..`)

**Problem**: Initially tried to embed `starter-prompts` from `internal/bootstrap/starter.go` using `//go:embed ../../starter-prompts`, which resulted in "invalid pattern syntax" error.

**Solution**: Created the embed file at the root level (`starter.go`) where it can directly reference `starter-prompts` without parent directory traversal.

**Lesson**: Always place `go:embed` directives in files at the same directory level or higher than the target directory. Never use `..` in embed patterns.

## Zap Logger Structured Fields

**Issue**: Zap requires structured field objects, not string literals

**Problem**: Initially used `logger.Info("message", "key", value)` which caused compilation errors about untyped string constants.

**Solution**: Use zap field constructors: `logger.Info("message", zap.String("key", value))`

**Lesson**: Always use zap's field constructors (`zap.String()`, `zap.Int()`, `zap.Error()`, etc.) for structured logging. This provides type safety and better performance.

## Regex Matching in Go

**Issue**: Different regex methods return different data types

**Problem**: Used `FindAllStringSubmatchIndex()` which returns `[][]int` (positions), but needed actual string matches.

**Solution**: Switched to `FindAllStringSubmatch()` which returns `[][]string` with the actual matched text.

**Lesson**: Carefully choose the right regex method based on whether you need positions or actual matches. For placeholder parsing, we needed the actual text.

## SQLite Driver Selection

**Decision**: Chose `modernc.org/sqlite` over `github.com/mattn/go-sqlite3`

**Rationale**:
- Pure Go implementation (no CGO dependency)
- Simplifies cross-platform builds (macOS Intel/ARM)
- Adequate performance for personal-scale usage
- FTS5 support for full-text search

**Trade-off**: Slightly slower than CGO-based driver, but build simplicity outweighs performance for this use case.

**Lesson**: Consider build complexity vs. performance trade-offs when choosing dependencies. For CLI tools distributed as binaries, pure Go implementations often win.

## Go Version Requirements

**Issue**: Some packages require newer Go versions

**Problem**: `modernc.org/sqlite` required Go 1.24+, which was newer than the installed version (1.23.2).

**Solution**: Running `go get` automatically upgraded the Go toolchain to 1.24.11.

**Lesson**: Be aware of Go version requirements in dependencies. The Go toolchain can manage multiple versions, but this may surprise users with older Go installations.

## Project Structure Organization

**Pattern**: Standard Go project layout with feature-based internal packages

**Structure**:
```
cmd/promptstack/    # Main application entry point
internal/            # Private packages organized by feature
  ├── config/      # Configuration management
  ├── setup/        # First-run setup
  ├── bootstrap/    # Application initialization
  ├── library/      # Prompt library management
  ├── history/      # History and database
  ├── prompt/       # Prompt data models
  └── logging/      # Logging setup
ui/                 # TUI components (future)
starter.go           # Embedded resources at root
```

**Benefits**:
- Clear separation of concerns
- Easy to locate code by feature
- Standard Go conventions
- Internal packages are truly private

**Lesson**: Organize internal packages by domain/feature rather than technical layer. This makes the codebase more navigable and maintainable.

## Error Handling Patterns

**Pattern**: Use `fmt.Errorf` with `%w` for error wrapping

**Example**:
```go
if err != nil {
    return nil, fmt.Errorf("failed to load config: %w", err)
}
```

**Benefits**:
- Preserves original error for unwrapping
- Adds context at each layer
- Enables `errors.Is()` and `errors.As()` checks
- Clear error messages for debugging

**Lesson**: Always wrap errors with context using `%w`. Never discard the original error. This makes debugging and error handling much easier.

## Frontmatter Parsing Strategy

**Decision**: Simple string-based parser instead of full YAML library

**Rationale**:
- Only need to extract key-value pairs
- Simple format: `key: value`
- Avoids additional dependency
- Sufficient for current requirements

**Implementation**:
```go
func parseFrontmatter(content string) (map[string]string, string, error) {
    // Check for --- markers
    // Split by lines
    // Parse key: value pairs
    // Return metadata and remaining content
}
```

**Trade-off**: Less robust than full YAML parser, but adequate for simple frontmatter.

**Lesson**: Choose the simplest solution that meets requirements. Don't over-engineer for future needs that may never materialize.

## Placeholder Parsing

**Pattern**: Regex with position tracking

**Implementation**:
```go
re := regexp.MustCompile(`\{\{(\w+):(\w+)\}\}`)
matches := re.FindAllStringSubmatch(content, -1)

for _, match := range matches {
    // Find positions in content
    startPos := strings.Index(content, fullMatch)
    endPos := startPos + len(fullMatch)
    // Create placeholder with positions
}
```

**Lesson**: When parsing structured text, track both the content and its positions. This enables features like highlighting and navigation.

## Index Scoring Algorithm

**Design**: Multi-factor scoring for relevance

**Factors**:
1. Tag matches: +10 per matching tag
2. Keyword overlap: +1 per matching word (weighted by frequency)
3. Recently used: +3 if used in last 24 hours
4. Frequently used: +use_count

**Rationale**:
- Tags are strong signals of relevance
- Keywords provide content-based matching
- Usage patterns reflect user preferences
- Time decay ensures fresh content

**Lesson**: Relevance scoring should combine multiple signals. No single factor is sufficient for good recommendations.

## Validation Strategy

**Pattern**: Separate errors and warnings

**Implementation**:
```go
type ValidationResult struct {
    Errors   []ValidationError  // Block insertion
    Warnings []ValidationError  // Allow with indicator
    IsValid  bool
}
```

**Benefits**:
- Graceful degradation
- User can still use prompts with warnings
- Clear distinction between critical and minor issues

**Lesson**: Validation should have severity levels. Not all issues should block functionality. Provide users with information and let them decide.

## Database Schema Design

**Pattern**: Separate tables with FTS5 for search

**Schema**:
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

CREATE VIRTUAL TABLE compositions_fts USING fts5(
    content,
    working_directory,
    content='compositions',
    content_rowid='id'
);
```

**Benefits**:
- Efficient full-text search
- Separate concerns (storage vs. search)
- Fast queries on metadata

**Lesson**: Use SQLite's FTS5 for full-text search. It's purpose-built for this use case and much faster than LIKE queries.

## Configuration Management

**Pattern**: YAML with validation

**Implementation**:
```go
type Config struct {
    Version      string `yaml:"version"`
    ClaudeAPIKey string `yaml:"claude_api_key"`
    Model        string `yaml:"model"`
    VimMode      bool   `yaml:"vim_mode"`
}

func (c *Config) Validate() error {
    if c.ClaudeAPIKey == "" {
        return fmt.Errorf("claude_api_key is required")
    }
    return nil
}
```

**Benefits**:
- Human-readable configuration
- Type-safe with struct tags
- Validation ensures correctness
- Easy to extend

**Lesson**: Always validate configuration after loading. Catch errors early and provide clear error messages.

## Starter Prompt Extraction

**Pattern**: Version-aware extraction

**Implementation**:
```go
func ExtractStarterPrompts(dataPath string, logger *zap.Logger) error {
    // Walk embedded filesystem
    // Check if file already exists
    // Only extract if not present
    // Log each extraction
}
```

**Benefits**:
- Preserves user modifications
- Adds new prompts on upgrade
- Never overwrites user changes
- Idempotent operation

**Lesson**: When extracting bundled resources, always check for existing files. Never overwrite user data. This enables safe upgrades.

## Logging Strategy

**Pattern**: Structured logging with rotation

**Configuration**:
- File-based logging to `~/.promptstack/debug.log`
- JSON format for easy parsing
- Rotation at 10MB, keep last 3
- Level control via environment variable

**Benefits**:
- Persistent logs for debugging
- Automatic cleanup prevents disk bloat
- Structured logs enable filtering and analysis
- Environment variable control for different environments

**Lesson**: Use structured logging from the start. It pays dividends in debugging and monitoring. Configure rotation to prevent disk issues.

## Bubble Tea Model Implementation

**Pattern**: Standard Bubble Tea model structure with Init(), Update(), View()

**Implementation**:
```go
type Model struct {
    // State fields
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle messages
    return m, cmd
}

func (m Model) View() string {
    // Render UI
    return styledContent
}
```

**Lesson**: Follow Bubble Tea's model-view-update pattern strictly. Always return the updated model and any commands from Update(). This ensures proper state management and message flow.

## Cursor and Viewport Management

**Pattern**: Track both cursor position and viewport offset

**Implementation**:
```go
type cursor struct {
    x int
    y int
}

type viewport struct {
    x int
    y int
}

func (m *Model) adjustViewport() {
    availableHeight := m.height - 1
    third := availableHeight / 3
    
    if m.cursor.y < m.viewport.y+third {
        m.viewport.y = max(0, m.cursor.y-third)
    } else if m.cursor.y > m.viewport.y+availableHeight-third {
        m.viewport.y = max(0, m.cursor.y-availableHeight+third)
    }
}
```

**Lesson**: Always keep cursor visible by adjusting viewport. Use a "middle third" strategy to provide smooth scrolling. This prevents cursor from getting stuck at edges.

## Auto-save Debouncing with Bubble Tea

**Pattern**: Use tea.Tick for timer-based operations

**Initial Approach** (problematic):
```go
func (m *Model) scheduleAutoSave() {
    if m.saveTimer != nil {
        m.saveTimer.Stop()
    }
    m.saveTimer = time.AfterFunc(750*time.Millisecond, func() {
        // Doesn't integrate with Bubble Tea message system
    })
}
```

**Better Approach** (Bubble Tea native):
```go
func (m Model) scheduleAutoSave() tea.Cmd {
    return tea.Tick(750*time.Millisecond, func(t time.Time) tea.Msg {
        return autoSaveMsg{}
    })
}

// In Update:
case autoSaveMsg:
    m.saveStatus = "saving"
    return m, tea.Cmd(func() tea.Msg {
        err := m.saveToFile()
        if err != nil {
            return saveErrorMsg{err}
        }
        return saveSuccessMsg{}
    })
```

**Lesson**: Use Bubble Tea's tea.Tick for timer-based operations instead of time.AfterFunc. This ensures proper integration with the message system and allows for clean state management.

## Custom Message Types

**Pattern**: Define custom message types for async operations

**Implementation**:
```go
type autoSaveMsg struct{}
type saveSuccessMsg struct{}
type saveErrorMsg struct {
    err error
}
type clearSaveStatusMsg struct{}

// In Update:
switch msg := msg.(type) {
case autoSaveMsg:
    // Handle auto-save trigger
case saveSuccessMsg:
    // Handle successful save
case saveErrorMsg:
    // Handle save error
case clearSaveStatusMsg:
    // Clear status after timeout
}
```

**Lesson**: Define custom message types for each async operation. This makes the code more readable and easier to maintain. Use type assertions to handle different message types.

## Status Bar State Management

**Pattern**: Track status with explicit states and auto-clear

**Implementation**:
```go
type statusBar struct {
    charCount int
    lineCount int
    message   string
}

// In Update:
case saveSuccessMsg:
    m.saveStatus = "saved"
    m.lastSave = time.Now()
    m.isDirty = false
    // Clear saved status after 2 seconds
    return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
        return clearSaveStatusMsg{}
    })

case clearSaveStatusMsg:
    m.saveStatus = ""
    m.statusBar.message = ""
    return m, nil
```

**Lesson**: Use explicit status states with auto-clear timers. This provides good user feedback without cluttering the UI. Always clear transient status messages after a reasonable timeout.

## Text Editor Cursor Positioning

**Pattern**: Handle cursor movement across line boundaries

**Implementation**:
```go
func (m *Model) moveCursorLeft() {
    if m.cursor.x > 0 {
        m.cursor.x--
    } else if m.cursor.y > 0 {
        // Move to end of previous line
        m.cursor.y--
        lines := strings.Split(m.content, "\n")
        if m.cursor.y < len(lines) {
            lineLen := lines[m.cursor.y]
            m.cursor.x = len(lineLen)
        }
    }
}
```

**Lesson**: Always handle edge cases when moving cursor. When moving left at column 0, move to end of previous line. When moving right at end of line, move to start of next line. This provides natural text editing behavior.

## File Path Management for History

**Pattern**: Timestamp-based file naming with directory creation

**Implementation**:
```go
func (m *Model) saveToFile() error {
    if m.filePath == "" {
        timestamp := time.Now().Format("2006-01-02_15-04-05")
        m.filePath = filepath.Join(m.workingDir, ".promptstack", ".history", timestamp+".md")
    }
    
    dir := filepath.Dir(m.filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }
    
    return os.WriteFile(m.filePath, []byte(m.content), 0644)
}
```

**Lesson**: Use timestamp-based naming for history files to avoid conflicts. Always create directories with MkdirAll before writing files. Use filepath.Join for cross-platform path construction.

## Lipgloss Styling

**Pattern**: Define reusable styles and compose them

**Implementation**:
```go
editorStyle := lipgloss.NewStyle().
    Width(m.width).
    Height(availableHeight).
    Padding(0, 1)

statusStyle := lipgloss.NewStyle().
    Width(m.width).
    Height(1).
    Background(lipgloss.Color("240")).
    Foreground(lipgloss.Color("15")).
    Padding(0, 1)

cursorStyle := lipgloss.NewStyle().
    Background(lipgloss.Color("7")).
    Foreground(lipgloss.Color("0"))
```

**Lesson**: Define styles as reusable variables or functions. Use Lipgloss's fluent API for clean style definitions. Use color codes (240 for gray, 7 for white, etc.) for consistent theming.

## Centralized Theme System

**Pattern**: Single source of truth for all UI colors and styles

**Implementation**:
```go
// ui/theme/theme.go
package theme

import "github.com/charmbracelet/lipgloss"

// Color Constants
const (
    BackgroundPrimary   = "#1e1e2e"
    BackgroundSecondary = "#181825"
    ForegroundPrimary   = "#cdd6f4"
    AccentBlue   = "#89b4fa"
    AccentGreen  = "#a6e3a1"
    // ... more constants
)

// Style Helper Functions
func ModalStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color(BorderPrimary)).
        Padding(1, 2).
        Background(lipgloss.Color(BackgroundPrimary))
}

func StatusStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color(ForegroundMuted)).
        Background(lipgloss.Color(BackgroundSecondary)).
        Padding(0, 1)
}

// ... more helper functions
```

**Usage in Components**:
```go
import "github.com/kyledavis/prompt-stack/ui/theme"

// Before: Hard-coded styles
modalStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#45475a")).
    Background(lipgloss.Color("#1e1e2e"))

// After: Theme helper
modalStyle := theme.ModalStyle()
```

**Benefits**:
- Single source of truth for all colors
- Easy to update entire UI theme
- Consistent color palette across components
- Type-safe constants prevent typos
- Code reduction through helper functions
- Simple to extend for new components

**Lesson**: Create a centralized theme package with color constants and style helper functions. Replace hard-coded lipgloss styles with theme helpers. This provides maintainability, consistency, and makes theme updates trivial. Organize colors by semantic purpose (backgrounds, foregrounds, accents, borders) rather than by component.

## Library Browser Implementation

**Pattern**: Modal overlay with fuzzy search and preview pane

**Implementation**:
```go
type Model struct {
    prompts      map[string]*prompt.Prompt
    filtered     []string
    selected     int
    searchInput  string
    width        int
    height       int
    visible      bool
    insertMode   InsertMode
    vimMode      bool
}

func (m *Model) applyFilter() {
    if m.searchInput == "" {
        // Show all prompts
        m.filtered = make([]string, 0, len(m.prompts))
        for path := range m.prompts {
            m.filtered = append(m.filtered, path)
        }
        return
    }

    // Build searchable strings (title + tags + category)
    var stringsToMatch []string
    var paths []string
    for path, p := range m.prompts {
        searchable := fmt.Sprintf("%s %s %s", p.Title, strings.Join(p.Tags, " "), p.Category)
        stringsToMatch = append(stringsToMatch, searchable)
        paths = append(paths, path)
    }

    // Apply fuzzy matching
    matches := fuzzy.Find(m.searchInput, stringsToMatch)
    
    // Update filtered list
    m.filtered = make([]string, 0, len(matches))
    for _, match := range matches {
        m.filtered = append(m.filtered, paths[match.Index])
    }
}
```

**Benefits**:
- Real-time filtering as user types
- Combines multiple fields for better search relevance
- Sahilm/fuzzy provides fast, simple fuzzy matching
- Modal overlay doesn't disrupt main workspace

**Lesson**: When implementing fuzzy search, combine multiple searchable fields (title, tags, category) into a single string. This provides better relevance than searching fields separately. Use sahilm/fuzzy for simple, fast fuzzy matching in Go.

## Modal Overlay Pattern

**Pattern**: Visibility flag with Show()/Hide() methods

**Implementation**:
```go
type Model struct {
    visible bool
    // ... other fields
}

func (m *Model) Show() {
    m.visible = true
    m.searchInput = ""
    m.applyFilter()
}

func (m *Model) Hide() {
    m.visible = false
}

func (m Model) View() string {
    if !m.visible {
        return ""
    }
    // Render modal content
}
```

**Benefits**:
- Clean separation of modal and main UI
- Easy to toggle visibility
- Modal state is self-contained

**Lesson**: Use a visibility flag for modals. Return empty string from View() when not visible. This keeps the main UI clean and makes modals easy to manage.

## Fuzzy Matching Integration

**Library**: sahilm/fuzzy

**Usage**:
```go
import "github.com/sahilm/fuzzy"

// Build list of strings to search
stringsToMatch := []string{"Command Palette", "Library Browser", "Settings"}

// Find matches
matches := fuzzy.Find("lib", stringsToMatch)

// Access results
for _, match := range matches {
    fmt.Printf("Match: %s (score: %d)\n", match.Str, match.Score)
}
```

**Benefits**:
- Simple API
- Fast performance (<10ms for 1000+ items)
- Returns ranked results with scores
- No complex configuration needed

**Lesson**: sahilm/fuzzy is ideal for in-memory fuzzy matching. It's fast, simple, and returns ranked results. Perfect for command palettes, library browsers, and similar UI components.

## Split-Pane Layout with Lipgloss

**Pattern**: Use lipgloss.JoinHorizontal() for side-by-side panels

**Implementation**:
```go
// Render prompt list
promptList := m.renderPromptList(modalWidth, modalHeight-4)

// Render preview
preview := m.renderPreview(modalWidth, modalHeight-4)

// Combine horizontally
content := lipgloss.JoinHorizontal(lipgloss.Top, promptList, preview)
```

**Benefits**:
- Clean separation of list and preview
- Automatic width calculation
- Responsive to modal size changes

**Lesson**: Use lipgloss.JoinHorizontal() for split-pane layouts. It handles width distribution automatically and keeps panels aligned. Add borders between panels for visual separation.

## Keyboard Navigation with Vim Mode

**Pattern**: Conditional keybinding based on vim mode flag

**Implementation**:
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyUp:
            m.moveSelection(-1)
        case tea.KeyDown:
            m.moveSelection(1)
        default:
            if msg.Type == tea.KeyRunes {
                m.searchInput += string(msg.Runes)
                m.applyFilter()
            }
        }

        // Vim mode keybindings
        if m.vimMode {
            switch msg.String() {
            case "j":
                m.moveSelection(1)
            case "k":
                m.moveSelection(-1)
            }
        }
    }
    return m, nil
}
```

**Benefits**:
- Universal vim support when enabled
- Falls back to arrow keys when disabled
- Consistent with user expectations

**Lesson**: Support both arrow keys and vim keybindings. Check vim mode flag and apply vim keys (j/k) when enabled. This provides a familiar experience for vim users while remaining accessible to everyone.

## Message-Based Command Execution

**Pattern**: Return custom message from Update() for async operations

**Implementation**:
```go
type InsertPromptMsg struct {
    FilePath   string
    InsertMode InsertMode
}

// In Update:
case tea.KeyEnter:
    if len(m.filtered) > 0 {
        return m, func() tea.Msg {
            return InsertPromptMsg{
                FilePath:   m.filtered[m.selected],
                InsertMode: InsertAtCursor,
            }
        }
    }
```

**Benefits**:
- Decouples UI from business logic
- Parent model handles command execution
- Clean separation of concerns

**Lesson**: Use custom message types for commands. Return them as tea.Cmd from Update(). This allows parent models to handle the actual execution, keeping the modal focused on UI logic.

## Command Registry Pattern

**Pattern**: Centralized command registration with handler functions

**Implementation**:
```go
type Command struct {
    ID          string
    Name        string
    Description string
    Category    string
    Handler     func() error
}

type Registry struct {
    commands map[string]*Command
}

func (r *Registry) Register(cmd *Command) error {
    if cmd.ID == "" {
        return fmt.Errorf("command ID cannot be empty")
    }
    if cmd.Handler == nil {
        return fmt.Errorf("command handler cannot be nil")
    }
    if _, exists := r.commands[cmd.ID]; exists {
        return fmt.Errorf("command with ID %s already registered", cmd.ID)
    }
    r.commands[cmd.ID] = cmd
    return nil
}
```

**Benefits**:
- Centralized command management
- Type-safe command registration
- Easy to add new commands
- Validation prevents duplicate IDs
- Handler functions encapsulate command logic

**Lesson**: Use a registry pattern for commands. Validate command registration (ID, handler, duplicates). This prevents runtime errors and makes command management maintainable.

## Command Palette Implementation

**Pattern**: Modal overlay with fuzzy search and message-based execution

**Implementation**:
```go
type Model struct {
    registry    *commands.Registry
    filtered    []*commands.Command
    selected    int
    searchInput string
    width       int
    height      int
    visible     bool
    vimMode     bool
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyEnter:
        if len(m.filtered) > 0 {
            cmd := m.filtered[m.selected]
            m.visible = false
            return m, func() tea.Msg {
                err := cmd.Handler()
                if err != nil {
                    return ExecuteErrorMsg{CommandID: cmd.ID, Error: err}
                }
                return ExecuteSuccessMsg{CommandID: cmd.ID}
            }
        }
    }
    return m, nil
}
```

**Benefits**:
- Fast command discovery with fuzzy search
- Keyboard-driven workflow
- Message-based execution decouples UI from logic
- Success/error feedback for user
- Vim mode support for consistency

**Lesson**: Implement command palette with fuzzy search across command name, description, and category. Return execution results as messages for proper error handling. This provides a fast, keyboard-driven command interface.

## Command Categorization

**Pattern**: Group commands by category for better organization

**Implementation**:
```go
// Register commands with categories
registry.Register(&Command{
    ID:          "toggle-ai",
    Name:        "Toggle AI Panel",
    Description: "Show or hide the AI suggestions panel",
    Category:    "General",
    Handler:     toggleAIPanel,
})

// Get commands by category
commands := registry.GetByCategory("General")

// Get all unique categories
categories := registry.GetCategories()
```

**Benefits**:
- Logical grouping of related commands
- Easier to find commands
- Can filter by category
- Better organization in UI

**Lesson**: Categorize commands to improve discoverability. Show category labels in UI (e.g., "[General] Toggle AI Panel"). This helps users quickly find relevant commands.

## Placeholder Command Handlers

**Pattern**: Register commands with placeholder handlers for future implementation

**Implementation**:
```go
func RegisterCoreCommands(registry *Registry) error {
    registry.Register(&Command{
        ID:          "toggle-ai",
        Name:        "Toggle AI Panel",
        Description: "Show or hide the AI suggestions panel",
        Category:    "General",
        Handler: func() error {
            // TODO: Implement AI panel toggle
            return fmt.Errorf("AI panel toggle not yet implemented")
        },
    })
    return nil
}
```

**Benefits**:
- Commands appear in palette immediately
- Clear indication of unimplemented features
- Easy to implement handlers later
- UI can show error feedback

**Lesson**: Register commands with placeholder handlers that return descriptive errors. This allows the UI to be complete while features are being implemented incrementally. Users see what's available and get clear feedback when selecting unimplemented commands.

## Error Handling Architecture

**Pattern**: Structured error types with severity levels and display strategies

**Implementation**:
```go
type ErrorType string
const (
    ErrorTypeFile      ErrorType = "file"
    ErrorTypeDatabase  ErrorType = "database"
    ErrorTypeAPI       ErrorType = "api"
    ErrorTypeParsing   ErrorType = "parsing"
    ErrorTypeConfig    ErrorType = "config"
    ErrorTypeValidation ErrorType = "validation"
)

type Severity string
const (
    SeverityError   Severity = "error"
    SeverityWarning Severity = "warning"
    SeverityInfo    Severity = "info"
)

type AppError struct {
    Type      ErrorType
    Severity  Severity
    Message   string
    Details   string
    Timestamp time.Time
    Retryable bool
    Err       error
}
```

**Benefits**:
- Clear categorization of error types
- Severity-based display strategy (modal vs. status bar)
- Retryable flag for transient failures
- Preserves original error for debugging
- Structured logging support

**Lesson**: Create structured error types with severity levels. Use severity to determine display strategy (modal for critical errors, status bar for warnings). Mark retryable errors to enable automatic retry logic.

## Status Bar Component Design

**Pattern**: Message-based updates with auto-dismiss and persistent modes

**Implementation**:
```go
type Model struct {
    message        string
    messageType    MessageType
    messageTimeout time.Time
    charCount      int
    lineCount      int
    tokenEstimate  int
    vimMode        string
    editMode       string
    showAutoSave   bool
    autoSaveStatus string
}

type SetMessageMsg struct {
    Message string
    Type    MessageType
    Timeout time.Duration
}

func SetErrorMessage(message string) tea.Cmd {
    return func() tea.Msg {
        return SetMessageMsg{
            Message: message,
            Type:    MessageTypeError,
            Timeout: 5 * time.Second,
        }
    }
}

func SetPersistentErrorMessage(message string) tea.Cmd {
    return func() tea.Msg {
        return SetMessageMsg{
            Message: message,
            Type:    MessageTypeError,
            Timeout: 0, // No timeout
        }
    }
}
```

**Benefits**:
- Message-based updates integrate with Bubble Tea
- Auto-dismiss for transient messages
- Persistent mode for critical errors
- Multiple message types (info, success, warning, error)
- Displays stats and mode indicators

**Lesson**: Use message-based updates for status bar. Support both auto-dismissing (with timeout) and persistent (no timeout) messages. Display contextual information (stats, modes) alongside messages.

## Modal Component Pattern

**Pattern**: Reusable modal with visibility flag and message-based control

**Implementation**:
```go
type Modal struct {
    title       string
    content     string
    width       int
    height      int
    showButtons bool
    primaryBtn  string
    secondaryBtn string
    focused     bool
}

func (m Modal) Update(msg tea.Msg) (Modal, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if !m.showButtons {
            switch msg.String() {
            case "esc", "enter", "q":
                return m, tea.Quit
            }
        } else {
            switch msg.String() {
            case "esc":
                return m, func() tea.Msg { return CloseModalMsg{} }
            case "enter":
                return m, func() tea.Msg { return ModalActionMsg{Action: "primary"} }
            }
        }
    }
    return m, nil
}

func ErrorModal(title, message string) Modal {
    return NewModal(title, message).WithButtons("OK", "")
}

func ConfirmModal(title, message string) Modal {
    return NewModal(title, message).WithButtons("Confirm", "Cancel")
}
```

**Benefits**:
- Reusable across different error types
- Visibility flag for clean UI integration
- Message-based control (CloseModalMsg, ModalActionMsg)
- Helper functions for common modal types
- Centered layout with proper sizing

**Lesson**: Create reusable modal components with visibility flags. Return empty string from View() when not visible. Use helper functions for common modal types (error, warning, confirm). This provides consistent error UI across the application.

## Error Handler Integration

**Pattern**: Centralized error handler with display strategy routing

**Implementation**:
```go
type Handler struct {
    showModal bool
    modal    common.Modal
}

func (h *Handler) Handle(err error) tea.Cmd {
    if err == nil {
        return nil
    }

    appErr, ok := err.(*AppError)
    if !ok {
        appErr = New(ErrorTypeFile, err.Error())
    }

    switch appErr.Severity {
    case SeverityError:
        return h.handleError(appErr)
    case SeverityWarning:
        return h.createWarningMessage(appErr)
    case SeverityInfo:
        return h.createInfoMessage(appErr.Message)
    default:
        return h.createErrorMessage(appErr.Message)
    }
}

func (h *Handler) handleError(err *AppError) tea.Cmd {
    cmd := h.createPersistentErrorMessage(err.Message)
    
    if h.isCriticalError(err) {
        h.showModal = true
        h.modal = common.ErrorModal("Error", h.formatErrorDetails(err))
        return tea.Batch(cmd, h.showModalCmd())
    }
    
    return cmd
}
```

**Benefits**:
- Centralized error handling logic
- Severity-based display routing
- Critical errors show both status bar and modal
- Helper functions for common error scenarios
- Consistent error display across application

**Lesson**: Create centralized error handler that routes errors based on severity. Critical errors show both status bar message and modal. Helper functions (HandleFileError, HandleDatabaseError, etc.) provide convenient error handling for common scenarios.

## Import Cycle Prevention

**Issue**: UI components importing internal packages can create circular dependencies

**Problem**: Status bar needed to import internal/errors for SetErrorFromAppError, but error handler needed to import ui/statusbar for message creation.

**Solution**: Create message types in internal/errors that status bar can handle, avoiding direct dependency

**Implementation**:
```go
// In internal/errors/handler.go
type SetStatusMessageMsg struct {
    Message string
    Type    string
    Timeout time.Duration
}

func (h *Handler) createErrorMessage(message string) tea.Cmd {
    return func() tea.Msg {
        return SetStatusMessageMsg{
            Message: message,
            Type:    "error",
            Timeout: 5 * time.Second,
        }
    }
}

// In ui/statusbar/model.go
case SetStatusMessageMsg:
    // Handle message from error handler
```

**Lesson**: Avoid import cycles by creating message types in lower-level packages. UI components handle messages from internal packages without importing them directly. This maintains clean architecture and prevents circular dependencies.

## Error Recovery Strategies

**Pattern**: Graceful degradation with user-friendly messages

**Implementation**:
```go
// File read errors: Load as plain text, warn in validation
func HandleFileError(operation string, err error) tea.Cmd {
    appErr := FileError(fmt.Sprintf("Failed to %s", operation), err)
    return NewHandler().Handle(appErr)
}

// Auto-save errors: Retry silently, show persistent error after max retries
func HandleAutoSaveError(err error, retryCount int) tea.Cmd {
    if err == nil {
        return NewHandler().createSuccessMessage("Saved")
    }
    
    if retryCount < 3 {
        return nil // Retry silently
    }
    
    appErr := FileError("Auto-save failed after multiple attempts", err).
        WithDetails("Your work is preserved in memory. Please save manually.")
    return NewHandler().Handle(appErr)
}

// Config errors: Prompt for reset
func ShowConfigResetPrompt(reason string) tea.Cmd {
    message := fmt.Sprintf(
        "Configuration issue detected: %s\n\nWould you like to reset to configuration?",
        reason,
    )
    modal := common.ConfirmModal("Configuration Error", message)
    return func() tea.Msg {
        return ShowErrorModalMsg{Modal: modal}
    }
}
```

**Benefits**:
- Preserves user work in memory
- Provides actionable next steps
- Retry logic for transient failures
- User-friendly error messages
- Graceful degradation

**Lesson**: Implement error recovery strategies that preserve user work. Retry transient failures silently. Show persistent errors with actionable next steps. Use modals for critical errors requiring user action. This provides a robust user experience even when things go wrong.

## Graceful File Read Error Handling

**Pattern**: Comprehensive error handling with graceful degradation

**Implementation**:
```go
type LoadError struct {
    FilePath string
    Error    error
    Severity string // "error" or "warning"
}

type Library struct {
    Prompts          map[string]*prompt.Prompt
    Index            *prompt.LibraryIndex
    logger           *zap.Logger
    LoadErrors       []LoadError // Track errors during loading
    ValidationErrors []errors.AppError
}

func readFileGracefully(filePath string, logger *zap.Logger) ([]byte, error) {
    // Check if file exists
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, errors.FileError("File not found", err).
                WithDetails(fmt.Sprintf("The file %s does not exist", filePath))
        }
        return nil, errors.FileError("Failed to access file", err).
            WithDetails(fmt.Sprintf("Could not access file: %s", filePath))
    }

    // Check file size (1MB limit)
    const maxFileSize = 1 << 20 // 1MB
    if fileInfo.Size() > maxFileSize {
        err := errors.FileError("File size exceeds limit", nil).
            WithDetails(fmt.Sprintf("File %s is %.2f MB (max: 1MB)",
                filePath, float64(fileInfo.Size())/(1024*1024)))
        logger.Warn("File size exceeds limit",
            zap.String("path", filePath),
            zap.Int64("size", fileInfo.Size()))
        return nil, err
    }

    // Check file permissions
    if fileInfo.Mode().Perm()&0400 == 0 {
        err := errors.FileError("Permission denied", nil).
            WithDetails(fmt.Sprintf("No read permission for file: %s", filePath))
        logger.Warn("Permission denied", zap.String("path", filePath))
        return nil, err
    }

    // Read file content
    content, err := os.ReadFile(filePath)
    if err != nil {
        // Handle specific error types
        if os.IsPermission(err) {
            return nil, errors.FileError("Permission denied", err).
                WithDetails(fmt.Sprintf("Cannot read file: %s", filePath))
        }
        if errors.Is(err, os.ErrClosed) {
            return nil, errors.FileError("File closed unexpectedly", err).
                WithDetails(fmt.Sprintf("File handle closed: %s", filePath))
        }
        
        // Generic file read error
        return nil, errors.FileError("Failed to read file", err).
            WithDetails(fmt.Sprintf("Error reading file: %s", filePath))
    }

    return content, nil
}
```

**Benefits**:
- Comprehensive error detection (not found, size, permissions, read errors)
- Detailed error messages with context
- Graceful degradation (continues loading other files)
- Error tracking for reporting
- Severity-based handling (error vs warning)
- All errors logged appropriately

**Lesson**: Implement comprehensive file read error handling that checks for multiple failure modes before attempting to read. Use structured errors with details for debugging. Track all errors during batch operations and continue processing other items. This provides robust error handling without stopping the entire operation.

**Error Tracking and Reporting**:
```go
// Helper methods for error reporting
func (l *Library) GetLoadErrors() []LoadError {
    return l.LoadErrors
}

func (l *Library) HasLoadErrors() bool {
    return len(l.LoadErrors) > 0
}

func (l *Library) GetErrorCount() int {
    return len(l.LoadErrors)
}

func (l *Library) GetErrorSummary() string {
    if len(l.LoadErrors) == 0 {
        return "No errors"
    }
    
    errorCount := 0
    warningCount := 0
    for _, err := range l.LoadErrors {
        if err.Severity == "error" {
            errorCount++
        } else {
            warningCount++
        }
    }
    
    if errorCount > 0 && warningCount > 0 {
        return fmt.Sprintf("%d errors, %d warnings", errorCount, warningCount)
    } else if errorCount > 0 {
        return fmt.Sprintf("%d errors", errorCount)
    } else {
        return fmt.Sprintf("%d warnings", warningCount)
    }
}
```

**Lesson**: Provide helper methods for error reporting and summary. This makes it easy for UI components to display error information to users. Track both errors and warnings separately to provide accurate summaries.

**Package Naming Conflicts**:
```go
// Use alias to avoid conflict with standard errors package
import (
    "errors"
    apperrors "github.com/kyledavis/prompt-stack/internal/errors"
)

// Now use apperrors.FileError() instead of errors.FileError()
```

**Lesson**: When importing custom error packages that conflict with standard library packages, use import aliases. This prevents naming conflicts and makes code clearer by explicitly showing which package is being used.

## Error Logging Integration

**Pattern**: Centralized error logging with global logger access

**Implementation**:
```go
// In internal/logging/logger.go
var (
    globalLogger *zap.Logger
    loggerMutex sync.RWMutex
)

func Initialize() (*zap.Logger, error) {
    // ... logger setup ...
    logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
    
    // Store global logger instance
    loggerMutex.Lock()
    globalLogger = logger
    loggerMutex.Unlock()
    
    return logger, nil
}

func GetLogger() (*zap.Logger, error) {
    loggerMutex.RLock()
    defer loggerMutex.RUnlock()
    
    if globalLogger == nil {
        return nil, fmt.Errorf("logger not initialized")
    }
    
    return globalLogger, nil
}

// In internal/errors/handler.go
type Handler struct {
    showModal bool
    modal     common.Modal
    logger    *zap.Logger
}

func NewHandlerWithLogger(logger *zap.Logger) *Handler {
    return &Handler{
        showModal: false,
        logger:    logger,
    }
}

func (h *Handler) logError(err error) {
    if err == nil {
        return
    }
    
    if h.logger == nil {
        // Fall back to global LogError function
        LogError(err)
        return
    }
    
    // Log using zap logger with appropriate severity
    if appErr, ok := err.(*AppError); ok {
        switch appErr.Severity {
        case SeverityError:
            h.logger.Error(appErr.Message,
                zap.String("type", string(appErr.Type)),
                zap.String("severity", string(appErr.Severity)),
                zap.String("details", appErr.Details),
                zap.Bool("retryable", appErr.Retryable),
                zap.Error(appErr.Err))
        case SeverityWarning:
            h.logger.Warn(appErr.Message,
                zap.String("type", string(appErr.Type)),
                zap.String("details", appErr.Details),
                zap.Bool("retryable", appErr.Retryable),
                zap.Error(appErr.Err))
        case SeverityInfo:
            h.logger.Info(appErr.Message,
                zap.String("type", string(appErr.Type)),
                zap.String("details", appErr.Details))
        }
    } else {
        h.logger.Error("Error occurred", zap.Error(err))
    }
}

func (h *Handler) Handle(err error) tea.Cmd {
    if err == nil {
        return nil
    }
    
    // Log the error
    h.logError(err)
    
    // ... rest of error handling ...
}

// Global LogError function for use without handler
func LogError(err error) {
    if err == nil {
        return
    }
    
    logger, err := logging.GetLogger()
    if err != nil || logger == nil {
        return
    }
    
    // Log using zap logger with appropriate severity
    if appErr, ok := err.(*AppError); ok {
        switch appErr.Severity {
        case SeverityError:
            logger.Error(appErr.Message,
                zap.String("type", string(appErr.Type)),
                zap.String("severity", string(appErr.Severity)),
                zap.String("details", appErr.Details),
                zap.Bool("retryable", appErr.Retryable),
                zap.Error(appErr.Err))
        case SeverityWarning:
            logger.Warn(appErr.Message,
                zap.String("type", string(appErr.Type)),
                zap.String("details", appErr.Details),
                zap.Bool("retryable", appErr.Retryable),
                zap.Error(appErr.Err))
        case SeverityInfo:
            logger.Info(appErr.Message,
                zap.String("type", string(appErr.Type)),
                zap.String("details", appErr.Details))
        }
    } else {
        logger.Error("Error occurred", zap.Error(err))
    }
}

// Update all helper functions to log errors
func HandleFileError(operation string, err error) tea.Cmd {
    if err == nil {
        return nil
    }
    
    appErr := FileError(fmt.Sprintf("Failed to %s", operation), err)
    LogError(appErr) // Log before handling
    return NewHandler().Handle(appErr)
}
```

**Benefits**:
- All errors automatically logged to debug.log
- Structured logging with severity levels (error, warning, info)
- Thread-safe global logger access
- Both Handler instances and global function can log errors
- Detailed error context (type, severity, details, retryable, stack trace)
- Automatic logging in all error helper functions
- Easy to debug issues with comprehensive error logs

**Lesson**: Integrate error logging throughout the application using a global logger pattern. Store logger instance in logging package with thread-safe access. Add logger field to error handler for instance-based logging. Create global LogError function for use without handler instances. Log all errors with appropriate severity levels and structured fields. This provides comprehensive error tracking for debugging and monitoring.

## Future Considerations

### Potential Improvements
1. **Frontmatter parsing**: Consider full YAML parser if frontmatter becomes more complex
2. **Index scoring**: Add machine learning for better relevance
3. **Validation**: Add more sophisticated checks (e.g., placeholder usage)
4. **Database**: Consider connection pooling for concurrent access
5. **Error handling**: Add retry logic for transient failures

### Technical Debt
1. **Sorting algorithm**: Bubble sort in indexer.go should be replaced with more efficient algorithm
2. **Placeholder parsing**: Could be more robust with better error handling
3. **Frontmatter parsing**: Simple parser may fail on complex YAML

### Architecture Decisions to Revisit
1. **Library loading**: Currently loads all prompts into memory. Consider lazy loading for large libraries.
2. **Index building**: Rebuilds entire index on load. Consider incremental updates.
3. **Validation**: Runs on all prompts. Consider caching results.

## Placeholder System Implementation

**Pattern**: Regex-based parsing with position tracking and navigation

**Implementation**:
```go
// internal/editor/placeholder.go
type Placeholder struct {
    Type         string   // "text" or "list"
    Name         string   // placeholder name
    StartPos     int      // position in content
    EndPos       int      // position in content
    CurrentValue string   // current filled value (for text)
    ListValues   []string // current filled values (for list)
    IsValid      bool     // whether syntax is valid
    IsActive     bool     // whether currently selected
}

func ParsePlaceholders(content string) []Placeholder {
    re := regexp.MustCompile(`\{\{(\w+):(\w+)\}\}`)
    matches := re.FindAllStringSubmatchIndex(content, -1)
    
    var placeholders []Placeholder
    for _, match := range matches {
        // Extract type, name, and positions
        // Create placeholder with validation
    }
    return placeholders
}

func ValidatePlaceholders(placeholders []Placeholder) []ValidationError {
    // Check for duplicate names
    // Validate types and names
    // Return errors and warnings
}

// Navigation helpers
func GetNextPlaceholder(placeholders []Placeholder, currentPos int) int
func GetPreviousPlaceholder(placeholders []Placeholder, currentPos int) int
func FindPlaceholderAtPosition(placeholders []Placeholder, pos int) int
```

**Integration with Workspace**:
```go
// ui/workspace/model.go
type Model struct {
    content           string
    placeholders      []editor.Placeholder
    activePlaceholder int // -1 if none active
    // ... other fields
}

func (m *Model) updatePlaceholders() {
    m.placeholders = editor.ParsePlaceholders(m.content)
    // Maintain active placeholder if still valid
}

func (m *Model) navigateToNextPlaceholder() bool {
    cursorPos := m.getCursorPosition()
    nextIndex := editor.GetNextPlaceholder(m.placeholders, cursorPos)
    if nextIndex >= 0 {
        m.activePlaceholder = nextIndex
        ph := m.placeholders[nextIndex]
        m.setCursorToPosition(ph.StartPos)
        return true
    }
    return false
}

func (m *Model) renderLineWithPlaceholders(line string, lineIndex int) string {
    // Calculate line start position in content
    // Find placeholders on this line
    // Apply highlighting if active
    return renderedLine
}
```

**Benefits**:
- Automatic placeholder detection on content changes
- Tab/Shift+Tab navigation between placeholders
- Visual highlighting of active placeholder
- Validation for duplicate names and invalid types
- Position tracking enables precise navigation and highlighting
- Separation of concerns (editor package handles parsing, workspace handles UI)

**Lesson**: Use regex with position tracking for placeholder parsing. Track both content and positions to enable features like navigation and highlighting. Re-parse placeholders on every content change to keep state synchronized. Use Tab/Shift+Tab for intuitive navigation between placeholders. Highlight active placeholders visually to provide clear feedback. Separate parsing logic from UI logic for better code organization.

**Placeholder Validation Strategy**:

**Pattern**: Separate errors and warnings with severity levels

**Implementation**:
```go
type ValidationError struct {
    Type    string // "error" or "warning"
    Message string // human-readable message
    Line    int    // line number
    Column  int    // column number
}

func ValidatePlaceholders(placeholders []Placeholder) []ValidationError {
    var errors []ValidationError
    nameMap := make(map[string]int)
    
    // Check for duplicate names
    for i, ph := range placeholders {
        if !ph.IsValid {
            continue
        }
        if prevIndex, exists := nameMap[ph.Name]; exists {
            errors = append(errors, ValidationError{
                Type:    "error",
                Message: "Duplicate placeholder name: " + ph.Name,
                Line:    getLineNumber(placeholders, i),
                Column:  ph.StartPos,
            })
            errors = append(errors, ValidationError{
                Type:    "error",
                Message: "Duplicate placeholder name: " + ph.Name,
                Line:    getLineNumber(placeholders, prevIndex),
                Column:  placeholders[prevIndex].StartPos,
            })
        } else {
            nameMap[ph.Name] = i
        }
    }
    
    // Validate each placeholder
    for i, ph := range placeholders {
        if !ph.IsValid {
            if !isValidPlaceholderType(ph.Type) {
                errors = append(errors, ValidationError{
                    Type:    "error",
                    Message: "Invalid placeholder type: " + ph.Type + " (must be 'text' or 'list')",
                    Line:    getLineNumber(placeholders, i),
                    Column:  ph.StartPos,
                })
            }
            if !isValidPlaceholderName(ph.Name) {
                errors = append(errors, ValidationError{
                    Type:    "error",
                    Message: "Invalid placeholder name: " + ph.Name + " (must be alphanumeric and underscores only)",
                    Line:    getLineNumber(placeholders, i),
                    Column:  ph.StartPos,
                })
            }
        }
    }
    
    return errors
}
```

**Benefits**:
- Clear distinction between errors (block insertion) and warnings (allow with indicator)
- Duplicate detection prevents confusion
- Type validation ensures only supported types are used
- Name validation prevents syntax errors
- Line and column information for precise error reporting

**Lesson**: Validate placeholders for both syntax and semantic correctness. Check for duplicate names to prevent runtime confusion. Separate errors (block insertion) from warnings (allow with indicator). Provide line and column information for precise error reporting. This ensures placeholders are valid before they're used.

**Cursor Position Management for Placeholders**:

**Pattern**: Convert between cursor coordinates and absolute positions

**Implementation**:
```go
func (m *Model) getCursorPosition() int {
    lines := strings.Split(m.content, "\n")
    pos := 0
    
    for i := 0; i < m.cursor.y && i < len(lines); i++ {
        pos += len(lines[i]) + 1 // +1 for newline
    }
    
    if m.cursor.y < len(lines) {
        pos += m.cursor.x
    }
    
    return pos
}

func (m *Model) setCursorToPosition(pos int) {
    lines := strings.Split(m.content, "\n")
    currentPos := 0
    
    for i, line := range lines {
        lineEnd := currentPos + len(line)
        
        if pos <= lineEnd {
            m.cursor.y = i
            m.cursor.x = pos - currentPos
            return
        }
        
        currentPos = lineEnd + 1 // +1 for newline
    }
    
    // If position is beyond content, set to end
    m.cursor.y = len(lines) - 1
    m.cursor.x = len(lines[len(lines)-1])
}
```

**Benefits**:
- Accurate navigation to placeholder positions
- Handles multi-line content correctly
- Edge case handling (position beyond content)
- Enables precise placeholder selection

**Lesson**: Implement bidirectional conversion between cursor coordinates (x, y) and absolute positions. This is essential for features like placeholder navigation where you need to jump to specific positions. Handle edge cases (position beyond content) gracefully. Account for newlines when calculating positions.

**Placeholder Highlighting in TUI**:

**Pattern**: Line-by-line rendering with position-based highlighting

**Implementation**:
```go
func (m *Model) renderLineWithPlaceholders(line string, lineIndex int) string {
    if len(m.placeholders) == 0 {
        return line
    }
    
    // Calculate line start position in content
    lineStartPos := 0
    lines := strings.Split(m.content, "\n")
    for i := 0; i < lineIndex && i < len(lines); i++ {
        lineStartPos += len(lines[i]) + 1 // +1 for newline
    }
    
    // Find placeholders on this line
    result := line
    offset := 0
    
    for _, ph := range m.placeholders {
        // Check if placeholder is on this line
        if ph.StartPos >= lineStartPos && ph.EndPos <= lineStartPos+len(line) {
            // Calculate position within line
            phStart := ph.StartPos - lineStartPos
            phEnd := ph.EndPos - lineStartPos
            
            // Apply highlighting if active
            if m.activePlaceholder >= 0 && m.placeholders[m.activePlaceholder].Name == ph.Name {
                placeholderStyle := theme.ActivePlaceholderStyle()
                placeholderText := line[phStart+offset : phEnd+offset]
                result = result[:phStart+offset] + placeholderStyle.Render(placeholderText) + result[phEnd+offset:]
                offset += len(placeholderStyle.Render(placeholderText)) - (phEnd - phStart)
            }
        }
    }
    
    return result
}
```

**Benefits**:
- Visual feedback for active placeholder
- Line-by-line rendering for performance
- Offset tracking for styled text length
- Only highlights active placeholder, not all placeholders

**Lesson**: Render placeholders line-by-line with position-based highlighting. Calculate line start position in content to determine which placeholders are on each line. Track offset when applying styles because styled text has different length than plain text. Only highlight the active placeholder to avoid visual clutter. This provides clear user feedback without overwhelming the UI.

**Theme Integration for Placeholders**:

**Pattern**: Centralized style for active placeholder highlighting

**Implementation**:
```go
// ui/theme/theme.go
func ActivePlaceholderStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Background(lipgloss.Color(AccentYellow)).
        Foreground(lipgloss.Color(BackgroundPrimary)).
        Bold(true)
}
```

**Benefits**:
- Consistent styling across application
- Easy to update placeholder appearance
- Uses existing color palette
- High contrast for visibility

**Lesson**: Create dedicated style functions for specific UI elements like active placeholders. Use high-contrast colors (yellow background with dark foreground) for visibility. Keep all styles in the centralized theme package for consistency. This makes it easy to update the entire visual appearance of the application.

## Text Placeholder Editing Mode

**Pattern**: Vim-style editing mode with state management and value replacement

**Implementation**:
```go
// ui/workspace/model.go
type Model struct {
    content              string
    cursor               cursor
    // ... other fields
    placeholders         []editor.Placeholder
    activePlaceholder     int // -1 if none active
    placeholderEditMode   bool // true when editing a placeholder
    placeholderEditValue  string // current value being edited
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle placeholder edit mode
        if m.placeholderEditMode {
            return m.handlePlaceholderEdit(msg)
        }
        
        switch msg.Type {
        case tea.KeyRunes:
            // Check for 'i' or 'a' to enter placeholder edit mode
            if m.activePlaceholder >= 0 && len(msg.Runes) == 1 {
                r := msg.Runes[0]
                if r == 'i' || r == 'a' {
                    m.enterPlaceholderEditMode()
                    return m, nil
                }
            }
            // Normal typing
            m.insertRune(msg.Runes)
        }
    }
    return m, nil
}

func (m *Model) enterPlaceholderEditMode() {
    if m.activePlaceholder < 0 || m.activePlaceholder >= len(m.placeholders) {
        return
    }
    
    ph := &m.placeholders[m.activePlaceholder]
    
    // Only text placeholders can be edited in this mode
    if ph.Type != "text" {
        return
    }
    
    // Initialize edit value with current value or empty string
    m.placeholderEditMode = true
    m.placeholderEditValue = ph.CurrentValue
    
    // Move cursor to placeholder position
    m.setCursorToPosition(ph.StartPos)
}

func (m Model) handlePlaceholderEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.Type {
    case tea.KeyEsc:
        // Exit placeholder edit mode and save the value
        m.exitPlaceholderEditMode()
        return m, nil
        
    case tea.KeyBackspace:
        // Delete character from edit value
        if len(m.placeholderEditValue) > 0 {
            m.placeholderEditValue = m.placeholderEditValue[:len(m.placeholderEditValue)-1]
        }
        return m, nil
        
    case tea.KeyEnter:
        // Exit placeholder edit mode and save the value
        m.exitPlaceholderEditMode()
        return m, nil
        
    case tea.KeyRunes:
        // Append characters to edit value
        m.placeholderEditValue += string(msg.Runes)
        return m, nil
    }
    
    return m, nil
}

func (m *Model) exitPlaceholderEditMode() {
    if m.activePlaceholder < 0 || m.activePlaceholder >= len(m.placeholders) {
        m.placeholderEditMode = false
        m.placeholderEditValue = ""
        return
    }
    
    ph := &m.placeholders[m.activePlaceholder]
    
    // Update placeholder's current value
    ph.CurrentValue = m.placeholderEditValue
    
    // Replace placeholder in content with filled value
    m.content = editor.ReplacePlaceholder(m.content, *ph)
    
    // Re-parse placeholders after replacement
    m.updatePlaceholders()
    
    // Mark as dirty and schedule auto-save
    m.markDirty()
    m.scheduleAutoSave()
    
    // Exit edit mode
    m.placeholderEditMode = false
    m.placeholderEditValue = ""
}
```

**Rendering Edit Mode**:
```go
func (m *Model) renderCursorLine(lines []string) string {
    if m.cursor.y >= len(lines) {
        return ""
    }
    
    line := lines[m.cursor.y]
    
    // If in placeholder edit mode, show edit value instead of placeholder
    if m.placeholderEditMode && m.activePlaceholder >= 0 {
        ph := m.placeholders[m.activePlaceholder]
        if ph.Type == "text" {
            // Calculate line start position
            lineStartPos := 0
            allLines := strings.Split(m.content, "\n")
            for i := 0; i < m.cursor.y && i < len(allLines); i++ {
                lineStartPos += len(allLines[i]) + 1
            }
            
            // Check if placeholder is on this line
            if ph.StartPos >= lineStartPos && ph.EndPos <= lineStartPos+len(line) {
                phStart := ph.StartPos - lineStartPos
                phEnd := ph.EndPos - lineStartPos
                
                // Replace placeholder with edit value
                before := line[:phStart]
                after := line[phEnd:]
                editValue := m.placeholderEditValue
                
                // Adjust cursor position to be within edit value
                cursorPosInEdit := m.cursor.x - phStart
                if cursorPosInEdit < 0 {
                    cursorPosInEdit = 0
                } else if cursorPosInEdit > len(editValue) {
                    cursorPosInEdit = len(editValue)
                }
                
                // Style cursor
                cursorStyle := theme.CursorStyle()
                
                if cursorPosInEdit < len(editValue) {
                    return before + cursorStyle.Render(string(editValue[cursorPosInEdit])) + editValue[cursorPosInEdit+1:] + after
                }
                
                return before + editValue + cursorStyle.Render(" ") + after
            }
        }
    }
    
    // Normal rendering
    // ...
}
```

**Status Bar Indicator**:
```go
func (m *Model) renderStatusBar() string {
    statusStyle := theme.StatusStyle().
        Width(m.width).
        Height(1)
    
    // Build status message
    var parts []string
    
    // Placeholder edit mode indicator
    if m.placeholderEditMode {
        parts = append(parts, "[PLACEHOLDER EDIT MODE]")
    }
    
    // Auto-save indicator
    if m.saveStatus == "saving" {
        parts = append(parts, "Saving...")
    }
    
    // Join with separator
    statusText := strings.Join(parts, " | ")
    
    return statusStyle.Render(statusText)
}
```

**Benefits**:
- Vim-style editing workflow familiar to developers
- Clear visual feedback with status bar indicator
- Type to replace placeholder content directly
- Esc or Enter to save and exit edit mode
- Automatic placeholder replacement with value
- Content re-parsing after editing
- Auto-save triggered after placeholder is filled
- Only text placeholders can be edited (list placeholders require separate UI)
- Cursor positioning maintained within edit value

**Lesson**: Implement placeholder editing as a separate mode with its own key handling. Use a boolean flag to track edit mode state. Store edit value separately from placeholder until exit. Show edit value instead of placeholder syntax during editing. Display clear mode indicator in status bar. On exit, replace placeholder with filled value and re-parse content. This provides a clean, intuitive editing experience that feels natural to vim users while remaining accessible to everyone.

---

**Last Updated**: 2026-01-06
**Implementation Phase**: Milestone 2 - Advanced Editing - In Progress (6/13 tasks complete)