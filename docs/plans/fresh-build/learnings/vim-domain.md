# Vim Domain Key Learnings

**Purpose**: Key learnings and implementation patterns for vim mode functionality from previous PromptStack implementation.

**Related Milestones**: M34, M35

**Related Documents**: 
- [`project-structure.md`](../project-structure.md) - Vim domain structure
- [`go-style-guide.md`](../go-style-guide.md) - Go coding standards
- [`go-testing-guide.md`](../go-testing-guide.md) - Testing patterns

---

## Learning Categories

### Category 1: Vim Mode Transition Logic

**Learning**: Implement vim mode transitions as a state machine with proper cleanup and initialization

**Problem**: Need to manage vim mode transitions (Normal, Insert, Visual) with proper state management.

**Solution**: State machine with transition hooks and mode-specific state management

**Implementation Pattern**:
```go
// internal/vim/state.go
type Mode int

const (
    NormalMode Mode = iota
    InsertMode
    VisualMode
)

type State struct {
    CurrentMode     Mode
    PreviousMode    Mode
    transitionHook  TransitionHook
    // ... other state fields
}

type TransitionHook func(from, to Mode)

func (s *State) SetMode(mode Mode) {
    if s.CurrentMode == mode {
        return
    }

    // Validate transition
    if !s.isValidTransition(s.CurrentMode, mode) {
        return
    }

    from := s.CurrentMode
    s.PreviousMode = from
    s.CurrentMode = mode

    // Call transition hook if set
    if s.transitionHook != nil {
        s.transitionHook(from, mode)
    }

    // Reset mode-specific state
    s.cleanupModeState(from)
    s.initializeModeState(mode)
}

func (s *State) cleanupModeState(mode Mode) {
    switch mode {
    case InsertMode:
        // Leaving insert mode - no special cleanup needed
    case VisualMode:
        // Leaving visual mode - clear visual selection
        s.IsVisualBlock = false
        s.IsVisualLine = false
        s.VisualStart = struct{ Line, Column int }{Line: 0, Column: 0}
        s.VisualEnd = struct{ Line, Column int }{Line: 0, Column: 0}
    case NormalMode:
        // Leaving normal mode - no special cleanup needed
    }
}

func (s *State) initializeModeState(mode Mode) {
    switch mode {
    case NormalMode:
        // Reset count and pending command when entering normal mode
        s.Count = 0
        s.PendingCommand = ""
    case InsertMode:
        // No special initialization needed
    case VisualMode:
        // Visual mode state is set by EnterVisualMode methods
    }
}

// Helper methods for mode checking
func (s *State) IsInNormalMode() bool {
    return s.CurrentMode == NormalMode
}

func (s *State) IsInInsertMode() bool {
    return s.CurrentMode == InsertMode
}

func (s *State) IsInVisualMode() bool {
    return s.CurrentMode == VisualMode
}

func (s *State) IsInVisualBlockMode() bool {
    return s.CurrentMode == VisualMode && s.IsVisualBlock
}

func (s *State) IsInVisualLineMode() bool {
    return s.CurrentMode == VisualMode && s.IsVisualLine
}

// Mode transition helpers
func (s *State) ToggleMode() {
    if s.CurrentMode == NormalMode {
        s.EnterInsertMode()
    } else if s.CurrentMode == InsertMode {
        s.EnterNormalMode()
    }
}

func (s *State) ReturnToPreviousMode() {
    if s.PreviousMode != s.CurrentMode {
        s.SetMode(s.PreviousMode)
    }
}

func (s *State) CanTransitionTo(mode Mode) bool {
    return s.isValidTransition(s.CurrentMode, mode)
}
```

**Benefits**:
- Clean state machine pattern for mode transitions
- Transition hooks enable external code to react to mode changes
- Mode-specific state cleanup prevents stale data
- Mode-specific initialization ensures proper state setup
- Helper methods provide convenient mode checking
- Previous mode tracking enables return-to-previous functionality
- Transition validation prevents invalid state changes
- Extensible design for future mode types

**Lesson**: Implement vim mode transitions as a state machine with proper cleanup and initialization. Use transition hooks to decouple mode changes from side effects. Always clean up state when leaving a mode (e.g., clear visual selection when leaving Visual mode). Initialize state when entering a mode (e.g., reset count when entering Normal mode). Track previous mode to enable return-to-previous functionality. Provide helper methods for mode checking to simplify conditional logic. This provides a robust, maintainable vim mode system that handles all edge cases correctly.

**Related Milestones**: M34, M35

**When to Apply**: When implementing vim mode or similar state machine patterns

---

### Category 2: Read-Only Mode During Async Operations

**Learning**: Use read-only mode to prevent concurrent edits during async operations

**Problem**: Need to prevent user from editing while AI is applying changes.

**Solution**: State-based UI feedback with editing restrictions

**Implementation Pattern**:
```go
// ui/workspace/model.go
type Model struct {
    content              string
    cursor               cursor
    // ... other fields
    isReadOnly           bool // true when AI is applying suggestion (blocks editing)
    aiApplying           bool // true when AI is actively applying a suggestion
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Block all editing when in read-only mode (AI applying suggestion)
        if m.isReadOnly {
            // Only allow cursor navigation in read-only mode
            switch msg.Type {
            case tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
                // Allow cursor navigation
            default:
                // Block all other keys
                return m, nil
            }
        }
        // ... rest of key handling
    }
    return m, nil
}

func (m *Model) SetAIApplying(applying bool) {
    m.aiApplying = applying
    // When AI is applying, also set read-only mode
    m.isReadOnly = applying
}

func (m Model) renderStatusBar() string {
    // Build status message
    var parts []string
    
    // AI applying indicator (highest priority)
    if m.aiApplying {
        parts = append(parts, "✨ AI is applying...")
    }
    
    // ... other status indicators
    
    return statusStyle.Render(statusText)
}
```

**Benefits**:
- Clear visual feedback when AI is applying changes
- Prevents user from editing while AI is modifying content
- Allows cursor navigation for viewing during application
- Simple state management with two boolean flags
- Automatic read-only mode activation when AI applies
- Status bar indicator shows highest priority message

**Lesson**: When implementing async operations that modify content, use read-only mode to prevent concurrent edits. Provide clear visual feedback in status bar. Allow cursor navigation so users can view changes while they're being applied. Use separate flags for state (aiApplying) and behavior (isReadOnly) to enable flexible control. This prevents race conditions and provides good user experience during async operations.

**Related Milestones**: M34, M35

**When to Apply**: When implementing async operations that modify content

---

## Quick Reference

| Learning | Milestone | Priority |
|----------|-----------|----------|
| Vim Mode Transition Logic | M34, M35 | High |
| Read-Only Mode During Async Operations | M34, M35 | High |

---

**Last Updated**: 2026-01-07
**Source**: [`key-learnings.md`](../key-learnings.md)