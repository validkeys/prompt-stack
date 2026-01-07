# UI/TUI Domain Key Learnings

**Purpose**: Key learnings and implementation patterns for UI/TUI functionality from previous PromptStack implementation.

**Related Milestones**: M2, M8, M19, M24, M25, M26, M31, M36, M37

**Related Documents**: 
- [`project-structure.md`](../project-structure.md) - UI domain structure
- [`go-style-guide.md`](../go-style-guide.md) - Go coding standards
- [`go-testing-guide.md`](../go-testing-guide.md) - Testing patterns

---

## Learning Categories

### Category 1: Bubble Tea Model Implementation

**Learning**: Follow Bubble Tea's model-view-update pattern strictly

**Problem**: Need to implement TUI components using Bubble Tea framework.

**Solution**: Standard Bubble Tea model structure with Init(), Update(), View()

**Implementation Pattern**:
```go
import "github.com/charmbracelet/bubbletea"

type Model struct {
    // State fields
    content string
    cursor  cursor
    width   int
    height  int
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle keyboard input
        return m, nil
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil
    }
    return m, nil
}

func (m Model) View() string {
    // Render UI
    return styledContent
}
```

**Lesson**: Follow Bubble Tea's model-view-update pattern strictly. Always return the updated model and any commands from Update(). This ensures proper state management and message flow.

**Related Milestones**: M2, M8, M19

**When to Apply**: When implementing any Bubble Tea TUI component

---

### Category 2: Cursor and Viewport Management

**Learning**: Always keep cursor visible by adjusting viewport

**Problem**: Need to scroll content while keeping cursor visible.

**Solution**: Track both cursor position and viewport offset with "middle third" strategy

**Implementation Pattern**:
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

func (m *Model) moveCursorDown() {
    lines := strings.Split(m.content, "\n")
    if m.cursor.y < len(lines)-1 {
        m.cursor.y++
        m.adjustViewport()
    }
}
```

**Lesson**: Always keep cursor visible by adjusting viewport. Use a "middle third" strategy to provide smooth scrolling. This prevents cursor from getting stuck at edges.

**Related Milestones**: M2, M8

**When to Apply**: When implementing scrollable content in TUI

---

### Category 3: Auto-save Debouncing with Bubble Tea

**Learning**: Use Bubble Tea's tea.Tick for timer-based operations

**Problem**: Need to auto-save content after user stops typing.

**Solution**: Use tea.Tick instead of time.AfterFunc for proper integration

**Implementation Pattern**:
```go
type autoSaveMsg struct{}

func (m Model) scheduleAutoSave() tea.Cmd {
    return tea.Tick(750*time.Millisecond, func(t time.Time) tea.Msg {
        return autoSaveMsg{}
    })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle key input
        m.insertRune(msg.Runes)
        // Schedule auto-save
        return m, m.scheduleAutoSave()
        
    case autoSaveMsg:
        m.saveStatus = "saving"
        return m, tea.Cmd(func() tea.Msg {
            err := m.saveToFile()
            if err != nil {
                return saveErrorMsg{err}
            }
            return saveSuccessMsg{}
        })
    }
    return m, nil
}
```

**Lesson**: Use Bubble Tea's tea.Tick for timer-based operations instead of time.AfterFunc. This ensures proper integration with the message system and allows for clean state management.

**Related Milestones**: M2, M8

**When to Apply**: When implementing timer-based operations in Bubble Tea

---

### Category 4: Custom Message Types

**Learning**: Define custom message types for async operations

**Problem**: Need to handle async operations like auto-save and AI suggestions.

**Solution**: Define custom message types for each async operation

**Implementation Pattern**:
```go
type autoSaveMsg struct{}
type saveSuccessMsg struct{}
type saveErrorMsg struct {
    err error
}
type clearSaveStatusMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case autoSaveMsg:
        // Handle auto-save trigger
        m.saveStatus = "saving"
        return m, tea.Cmd(func() tea.Msg {
            err := m.saveToFile()
            if err != nil {
                return saveErrorMsg{err}
            }
            return saveSuccessMsg{}
        })
        
    case saveSuccessMsg:
        m.saveStatus = "saved"
        m.lastSave = time.Now()
        m.isDirty = false
        // Clear saved status after 2 seconds
        return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
            return clearSaveStatusMsg{}
        })
        
    case saveErrorMsg:
        m.saveStatus = "error: " + msg.err.Error()
        return m, nil
        
    case clearSaveStatusMsg:
        m.saveStatus = ""
        return m, nil
    }
    return m, nil
}
```

**Lesson**: Define custom message types for each async operation. This makes the code more readable and easier to maintain. Use type assertions to handle different message types.

**Related Milestones**: M2, M8, M27

**When to Apply**: When implementing async operations in Bubble Tea

---

### Category 5: Status Bar State Management

**Learning**: Use explicit status states with auto-clear timers

**Problem**: Need to show transient status messages without cluttering UI.

**Solution**: Track status with explicit states and auto-clear timers

**Implementation Pattern**:
```go
type statusBar struct {
    charCount int
    lineCount int
    message   string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
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
    }
    return m, nil
}

func (m Model) renderStatusBar() string {
    statusStyle := theme.StatusStyle().
        Width(m.width).
        Height(1)
    
    // Build status message
    var parts []string
    
    if m.saveStatus != "" {
        parts = append(parts, m.saveStatus)
    }
    
    parts = append(parts, fmt.Sprintf("%d chars, %d lines", 
        m.statusBar.charCount, m.statusBar.lineCount))
    
    statusText := strings.Join(parts, " | ")
    return statusStyle.Render(statusText)
}
```

**Lesson**: Use explicit status states with auto-clear timers. This provides good user feedback without cluttering the UI. Always clear transient status messages after a reasonable timeout.

**Related Milestones**: M2, M8

**When to Apply**: When implementing status indicators in TUI

---

### Category 6: Centralized Theme System

**Learning**: Create a centralized theme package with color constants and style helper functions

**Problem**: Need consistent styling across all TUI components.

**Solution**: Single source of truth for all UI colors and styles

**Implementation Pattern**:
```go
// ui/theme/theme.go
package theme

import "github.com/charmbracelet/lipgloss"

// Color Constants
const (
    BackgroundPrimary   = "#1e1e2e"
    BackgroundSecondary = "#181825"
    ForegroundPrimary   = "#cdd6f4"
    ForegroundMuted    = "#a6adc8"
    AccentBlue   = "#89b4fa"
    AccentGreen  = "#a6e3a1"
    AccentYellow = "#f9e2af"
    AccentRed    = "#f38ba8"
    BorderPrimary = "#45475a"
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

func ActivePlaceholderStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Background(lipgloss.Color(AccentYellow)).
        Foreground(lipgloss.Color(BackgroundPrimary)).
        Bold(true)
}
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

**Lesson**: Create a centralized theme package with color constants and style helper functions. Replace hard-coded lipgloss styles with theme helpers. Organize colors by semantic purpose (backgrounds, foregrounds, accents, borders) rather than by component.

**Related Milestones**: M2, M8, M19

**When to Apply**: When implementing TUI with multiple components

---

### Category 7: Library Browser Implementation

**Learning**: Combine multiple searchable fields for better fuzzy matching

**Problem**: Need to search prompts by title, tags, and category.

**Solution**: Build searchable strings combining multiple fields

**Implementation Pattern**:
```go
import "github.com/sahilm/fuzzy"

type Model struct {
    prompts      map[string]*prompt.Prompt
    filtered     []string
    selected     int
    searchInput  string
    width        int
    height       int
    visible      bool
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
        searchable := fmt.Sprintf("%s %s %s", 
            p.Title, strings.Join(p.Tags, " "), p.Category)
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

**Related Milestones**: M19, M24

**When to Apply**: When implementing search/filter functionality

---

### Category 8: Modal Overlay Pattern

**Learning**: Use a visibility flag for modals

**Problem**: Need to show/hide modals without disrupting main UI.

**Solution**: Visibility flag with Show()/Hide() methods

**Implementation Pattern**:
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
    return styledContent
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if !m.visible {
        return m, nil
    }
    // Handle messages when visible
    return m, nil
}
```

**Benefits**:
- Clean separation of modal and main UI
- Easy to toggle visibility
- Modal state is self-contained

**Lesson**: Use a visibility flag for modals. Return empty string from View() when not visible. This keeps the main UI clean and makes modals easy to manage.

**Related Milestones**: M19, M24, M25

**When to Apply**: When implementing modal overlays in TUI

---

### Category 9: Fuzzy Matching Integration

**Learning**: Use sahilm/fuzzy for simple, fast fuzzy matching

**Problem**: Need fuzzy search for command palette and library browser.

**Solution**: Use sahilm/fuzzy library

**Implementation Pattern**:
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

**Related Milestones**: M19, M24, M25

**When to Apply**: When implementing fuzzy search in Go applications

---

### Category 10: Split-Pane Layout with Lipgloss

**Learning**: Use lipgloss.JoinHorizontal() for side-by-side panels

**Problem**: Need to display list and preview side-by-side in modal.

**Solution**: Use lipgloss.JoinHorizontal() for automatic width distribution

**Implementation Pattern**:
```go
func (m Model) View() string {
    if !m.visible {
        return ""
    }
    
    modalWidth := min(m.width-4, 80)
    modalHeight := min(m.height-4, 20)
    
    // Render prompt list
    promptList := m.renderPromptList(modalWidth, modalHeight-4)
    
    // Render preview
    preview := m.renderPreview(modalWidth, modalHeight-4)
    
    // Combine horizontally
    content := lipgloss.JoinHorizontal(lipgloss.Top, promptList, preview)
    
    // Wrap in modal border
    return theme.ModalStyle().
        Width(modalWidth).
        Height(modalHeight).
        Render(content)
}
```

**Benefits**:
- Clean separation of list and preview
- Automatic width calculation
- Responsive to modal size changes

**Lesson**: Use lipgloss.JoinHorizontal() for split-pane layouts. It handles width distribution automatically and keeps panels aligned. Add borders between panels for visual separation.

**Related Milestones**: M19, M24

**When to Apply**: When implementing split-pane layouts in TUI

---

### Category 11: Keyboard Navigation with Vim Mode

**Learning**: Support both arrow keys and vim keybindings

**Problem**: Need to support vim-style navigation when enabled.

**Solution**: Conditional keybinding based on vim mode flag

**Implementation Pattern**:
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

**Related Milestones**: M34, M35, M37

**When to Apply**: When implementing keyboard navigation in TUI

---

### Category 12: Message-Based Command Execution

**Learning**: Return custom message from Update() for async operations

**Problem**: Need to execute commands from command palette.

**Solution**: Return custom message types as tea.Cmd

**Implementation Pattern**:
```go
type InsertPromptMsg struct {
    FilePath   string
    InsertMode InsertMode
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyEnter:
        if len(m.filtered) > 0 {
            return m, func() tea.Msg {
                return InsertPromptMsg{
                    FilePath:   m.filtered[m.selected],
                    InsertMode: InsertAtCursor,
                }
            }
        }
    }
    return m, nil
}

// In parent model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case InsertPromptMsg:
        m.insertPromptAtCursor(msg.FilePath, msg.InsertMode)
        return m, nil
    }
    return m, nil
}
```

**Benefits**:
- Decouples UI from business logic
- Parent model handles command execution
- Clean separation of concerns

**Lesson**: Use custom message types for commands. Return them as tea.Cmd from Update(). This allows parent models to handle the actual execution, keeping the modal focused on UI logic.

**Related Milestones**: M24, M25

**When to Apply**: When implementing command execution in Bubble Tea

---

### Category 13: Command Palette Implementation

**Learning**: Implement command palette with fuzzy search across multiple fields

**Problem**: Need fast command discovery and execution.

**Solution**: Modal overlay with fuzzy search and message-based execution

**Implementation Pattern**:
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

func (m Model) applyFilter() {
    if m.searchInput == "" {
        // Show all commands
        m.filtered = m.registry.GetAll()
        return
    }

    // Build searchable strings (name + description + category)
    var stringsToMatch []string
    for _, cmd := range m.registry.GetAll() {
        searchable := fmt.Sprintf("%s %s %s", 
            cmd.Name, cmd.Description, cmd.Category)
        stringsToMatch = append(stringsToMatch, searchable)
    }

    // Apply fuzzy matching
    matches := fuzzy.Find(m.searchInput, stringsToMatch)
    
    // Update filtered list
    m.filtered = make([]*commands.Command, 0, len(matches))
    for _, match := range matches {
        m.filtered = append(m.filtered, m.registry.GetAll()[match.Index])
    }
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

**Related Milestones**: M24, M25

**When to Apply**: When implementing command palette in TUI applications

---

### Category 14: Command Categorization

**Learning**: Categorize commands to improve discoverability

**Problem**: Need to organize many commands for easy discovery.

**Solution**: Group commands by category

**Implementation Pattern**:
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

// Display with category labels
func (m Model) renderCommandList() string {
    var builder strings.Builder
    for _, cmd := range m.filtered {
        builder.WriteString(fmt.Sprintf("[%s] %s - %s\n", 
            cmd.Category, cmd.Name, cmd.Description))
    }
    return builder.String()
}
```

**Benefits**:
- Logical grouping of related commands
- Easier to find commands
- Can filter by category
- Better organization in UI

**Lesson**: Categorize commands to improve discoverability. Show category labels in UI (e.g., "[General] Toggle AI Panel"). This helps users quickly find relevant commands.

**Related Milestones**: M24, M25

**When to Apply**: When implementing command systems with many commands

---

### Category 15: Placeholder Command Handlers

**Learning**: Register commands with placeholder handlers for future implementation

**Problem**: Need to show all commands even if some aren't implemented yet.

**Solution**: Register commands with placeholder handlers that return descriptive errors

**Implementation Pattern**:
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

**Related Milestones**: M24, M25

**When to Apply**: When implementing command systems with incremental feature rollout

---

## Quick Reference

| Learning | Milestone | Priority |
|----------|-----------|----------|
| Bubble Tea Model Implementation | M2, M8, M19 | High |
| Cursor and Viewport Management | M2, M8 | High |
| Auto-save Debouncing with Bubble Tea | M2, M8 | High |
| Custom Message Types | M2, M8, M27 | High |
| Status Bar State Management | M2, M8 | High |
| Centralized Theme System | M2, M8, M19 | High |
| Library Browser Implementation | M19, M24 | High |
| Modal Overlay Pattern | M19, M24, M25 | High |
| Fuzzy Matching Integration | M19, M24, M25 | High |
| Split-Pane Layout with Lipgloss | M19, M24 | High |
| Keyboard Navigation with Vim Mode | M34, M35, M37 | High |
| Message-Based Command Execution | M24, M25 | High |
| Command Palette Implementation | M24, M25 | High |
| Command Categorization | M24, M25 | Medium |
| Placeholder Command Handlers | M24, M25 | Medium |

---

**Last Updated**: 2026-01-07
**Source**: [`key-learnings.md`](../key-learnings.md)