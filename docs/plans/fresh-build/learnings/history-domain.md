# History Domain Key Learnings

**Purpose**: Key learnings and implementation patterns for history management from previous PromptStack implementation.

**Related Milestones**: M15, M16, M17

**Related Documents**: 
- [`project-structure.md`](../project-structure.md) - History domain structure
- [`go-style-guide.md`](../go-style-guide.md) - Go coding standards
- [`go-testing-guide.md`](../go-testing-guide.md) - Testing patterns

---

## Learning Categories

### Category 1: History Browser Integration

**Learning**: Use message-based operations for history browser actions

**Problem**: Need to display and interact with history compositions.

**Solution**: Modal overlay with list-based navigation and message-based operations

**Implementation Pattern**:
```go
// ui/app/model.go
type Model struct {
    history         historyui.Model
    historyManager  *history.Manager
    // ... other fields
}

func (m *Model) showHistoryBrowser() {
    if m.historyManager == nil {
        m.workspace.SetStatus("History manager not initialized")
        return
    }

    // Get all compositions from history
    compositions, err := m.historyManager.GetAllCompositions()
    if err != nil {
        m.workspace.SetStatus(fmt.Sprintf("Failed to load history: %v", err))
        return
    }

    // Convert to history items
    var items []list.Item
    for _, comp := range compositions {
        preview := comp.Content
        if len(preview) > 100 {
            preview = preview[:100] + "..."
        }
        
        item := historyui.NewItem(
            comp.FilePath,
            comp.CreatedAt,
            comp.WorkingDirectory,
            preview,
            comp.CharacterCount,
            comp.LineCount,
        )
        items = append(items, item)
    }

    // Set items in history browser
    m.history.SetItems(items)

    // Set up callbacks
    m.history.SetOnSelect(func(filePath string) tea.Cmd {
        return func() tea.Msg {
            return LoadHistoryMsg{FilePath: filePath}
        }
    })

    m.history.SetOnDelete(func(filePath string) tea.Cmd {
        return func() tea.Msg {
            return DeleteHistoryMsg{FilePath: filePath}
        }
    })
}

// Message Handling
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case LoadHistoryMsg:
        // Load history composition into workspace
        if m.historyManager != nil {
            content, err := m.historyManager.LoadComposition(msg.FilePath)
            if err != nil {
                m.workspace.SetStatus(fmt.Sprintf("Failed to load: %v", err))
            } else {
                m.workspace.SetContent(content)
                m.workspace.SetStatus(fmt.Sprintf("Loaded: %s", msg.FilePath))
            }
        }
        m.activePanel = "workspace"
        
    case DeleteHistoryMsg:
        // Delete history composition
        if m.historyManager != nil {
            err := m.historyManager.DeleteComposition(msg.FilePath)
            if err != nil {
                m.workspace.SetStatus(fmt.Sprintf("Failed to delete: %v", err))
            } else {
                m.workspace.SetStatus(fmt.Sprintf("Deleted: %s", msg.FilePath))
                // Refresh history browser
                m.showHistoryBrowser()
            }
        }
    }
    return m, nil
}
```

**Benefits**:
- History browser displays all saved compositions with timestamps and previews
- Search functionality with "/" key for filtering
- Load composition with Enter key
- Delete composition with Delete key
- Message-based operations decouple UI from business logic
- Automatic refresh after delete operations

**Lesson**: Use message-based operations for history browser actions. Define custom message types (LoadHistoryMsg, DeleteHistoryMsg) that are handled in the parent app model. This keeps the history browser focused on UI logic while the parent handles the actual business operations (loading, deleting). Always refresh the browser after destructive operations to keep the UI in sync.

**Related Milestones**: M15, M16, M17

**When to Apply**: When implementing history or similar data management features

---

### Category 2: History Manager Initialization in Bootstrap

**Learning**: Initialize all dependencies in bootstrap before creating TUI model

**Problem**: Need to properly initialize history manager with database and storage.

**Solution**: Dependency injection with proper initialization order

**Implementation Pattern**:
```go
// internal/bootstrap/bootstrap.go
func (a *App) Run() error {
    // ... existing code ...
    
    // Initialize history manager
    dbPath, err := config.GetDatabasePath()
    if err != nil {
        return fmt.Errorf("failed to get database path: %w", err)
    }

    db, err := history.Initialize(dbPath, a.Logger)
    if err != nil {
        return fmt.Errorf("failed to initialize history database: %w", err)
    }

    // Get history directory path
    historyDir, err := config.GetHistoryPath()
    if err != nil {
        return fmt.Errorf("failed to get history path: %w", err)
    }

    storage, err := history.NewStorage(historyDir, a.Logger)
    if err != nil {
        return fmt.Errorf("failed to initialize history storage: %w", err)
    }

    historyMgr := history.NewManager(db, storage, a.Logger)

    // Initialize TUI model with library, commands, and history manager
    tuiModel := app.NewWithDependencies(workingDir, lib, registry, historyMgr, a.Config.VimMode)
    
    // ... rest of code ...
}
```

**Benefits**:
- History manager properly initialized before TUI starts
- Database and storage configured with correct paths
- Dependency injection enables testing
- Clear initialization order (db → storage → manager → app)
- Error handling at each initialization step

**Lesson**: Initialize all dependencies in bootstrap before creating the TUI model. Use dependency injection to pass initialized components to the app model. Follow a clear initialization order (database first, then storage, then manager). Handle errors at each step to provide clear feedback about what failed. This ensures all components are ready when the TUI starts.

**Related Milestones**: M15, M16, M17

**When to Apply**: When initializing complex dependencies in application bootstrap

---

## Quick Reference

| Learning | Milestone | Priority |
|----------|-----------|----------|
| History Browser Integration | M15, M16, M17 | High |
| History Manager Initialization in Bootstrap | M15, M16, M17 | High |

---

**Last Updated**: 2026-01-07
**Source**: [`key-learnings.md`](../key-learnings.md)