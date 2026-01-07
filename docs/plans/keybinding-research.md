# Keybinding Management Research for Prompt Stack

## Executive Summary

This document researches Go libraries for managing keybindings in terminal UI applications, specifically for use with Bubble Tea and the vimtea-based Prompt Stack application.

## Current Implementation Analysis

### Existing Keybinding Architecture

The Prompt Stack application currently uses **Bubble Tea's native key handling** with manual switch statements across all UI components:

#### Key Components:
1. **Global Keybindings** ([`ui/app/model.go`](archive/code/ui/app/model.go:95-118))
   - Ctrl+C: Quit
   - Ctrl+P: Command palette
   - Ctrl+B: Library browser
   - Ctrl+H: History browser

2. **Workspace Keybindings** ([`ui/workspace/model.go`](archive/code/ui/workspace/model.go:103-182))
   - Ctrl+Z/Y: Undo/Redo
   - Arrow keys: Navigation
   - Tab/Shift+Tab: Navigate placeholders
   - Space, Enter, Backspace: Text editing

3. **Modal Keybindings** (browser, palette, validation, history)
   - Esc: Close modal
   - Enter: Select/Execute
   - Arrow keys: Navigate
   - Vim mode: j/k for navigation

4. **Vim State Machine** ([`internal/vim/state.go`](archive/code/internal/vim/state.go))
   - Sophisticated mode management (Normal, Insert, Visual)
   - Visual selection tracking
   - Register management
   - Command building with counts
   - Macro recording support

### Current Challenges

1. **Code Duplication**: Each component implements similar keybinding patterns
2. **Maintenance Burden**: Adding new keybindings requires changes in multiple files
3. **No Central Registry**: Keybindings are scattered across components
4. **Limited Customization**: Users cannot customize keybindings
5. **Conflict Detection**: No mechanism to detect keybinding conflicts
6. **Context Awareness**: Limited context-aware keybinding routing

## Go Keybinding Libraries Research

### 1. Bubble Tea Native (Current Approach)

**Library**: Built into Bubble Tea framework

**Pros**:
- ✅ Zero dependencies
- ✅ Full control over key handling
- ✅ Already integrated
- ✅ Well-documented in Bubble Tea ecosystem
- ✅ Type-safe with `tea.KeyMsg`

**Cons**:
- ❌ No built-in keybinding registry
- ❌ Manual switch statements required
- ❌ No conflict detection
- ❌ No user customization support
- ❌ Code duplication across components

**Verdict**: Good for simple apps, but doesn't scale well for complex applications with many keybindings.

---

### 2. github.com/charmbracelet/bubbletea (with custom keybinding layer)

**Approach**: Build a custom keybinding management layer on top of Bubble Tea

**Pros**:
- ✅ Tailored to specific needs
- ✅ Full control over implementation
- ✅ Can integrate with existing vim state machine
- ✅ No external dependencies
- ✅ Can support user customization

**Cons**:
- ❌ Requires significant development effort
- ❌ Maintenance burden on the team
- ❌ Need to implement conflict detection, routing, etc.

**Verdict**: Best option if no suitable library exists, but requires substantial investment.

---

### 3. github.com/charmbracelet/bubbles (list, textinput, etc.)

**Library**: Bubble Tea component library

**Pros**:
- ✅ Official Bubble Tea components
- ✅ Built-in keybinding handling for common patterns
- ✅ Well-tested and maintained
- ✅ Consistent with Bubble Tea patterns

**Cons**:
- ❌ Not a keybinding management library per se
- ❌ Limited to specific components (list, textinput, etc.)
- ❌ Doesn't solve global keybinding management

**Verdict**: Good for individual components, but not a solution for global keybinding management.

---

### 4. github.com/maaslalani/gambit

**Library**: Game engine for terminal UI with keybinding support

**Pros**:
- ✅ Keybinding management built-in
- ✅ Supports complex key sequences
- ✅ Context-aware keybinding routing

**Cons**:
- ❌ Designed for games, not TUI applications
- ❌ Heavy dependency for just keybindings
- ❌ May not integrate well with Bubble Tea
- ❌ Less active maintenance

**Verdict**: Not suitable for this use case.

---

### 5. github.com/rivo/tview

**Library**: Rich TUI framework with keybinding support

**Pros**:
- ✅ Comprehensive keybinding system
- ✅ Context-aware routing
- ✅ Supports key chords and sequences
- ✅ Well-maintained

**Cons**:
- ❌ Competing framework to Bubble Tea
- ❌ Would require major rewrite
- ❌ Not compatible with existing Bubble Tea code
- ❌ Heavy dependency

**Verdict**: Not compatible with Bubble Tea architecture.

---

### 6. github.com/gdamore/tcell

**Library**: Low-level terminal interface

**Pros**:
- ✅ Low-level control over input
- ✅ Supports all key combinations
- ✅ Cross-platform

**Cons**:
- ❌ Too low-level for keybinding management
- ❌ Would require building abstraction layer
- ❌ Not designed for keybinding management

**Verdict**: Useful as a foundation, but not a keybinding management solution.

---

### 7. github.com/nsf/termbox-go

**Library**: Terminal interface library

**Pros**:
- ✅ Simple API
- ✅ Cross-platform

**Cons**:
- ❌ Low-level input handling
- ❌ No keybinding management features
- ❌ Less active maintenance

**Verdict**: Not suitable for keybinding management.

---

## Recommended Approach: Custom Keybinding Management Layer

Based on the research, **the best approach is to build a custom keybinding management layer** on top of Bubble Tea. Here's why:

### Why Custom?

1. **No Suitable Library Exists**: No Go library provides keybinding management specifically for Bubble Tea
2. **Existing Investment**: The project already has a sophisticated vim state machine
3. **Full Control**: Custom implementation allows perfect integration with existing architecture
4. **Future Flexibility**: Can add features like user customization, conflict detection, etc.

### Proposed Architecture

```go
// internal/keybindings/registry.go
package keybindings

import (
    tea "github.com/charmbracelet/bubbletea"
)

// Key represents a key combination
type Key struct {
    Type  tea.KeyType
    Alt   bool
    Ctrl  bool
    Runes []rune
}

// Action represents a keybinding action
type Action struct {
    ID          string
    Description string
    Handler     func() tea.Cmd
    Context     string // "global", "workspace", "modal", etc.
    Mode        string // "normal", "insert", "visual" for vim
}

// Registry manages keybindings
type Registry struct {
    bindings map[Key]Action
    contexts map[string]map[Key]Action
    modes    map[string]map[Key]Action
}

// Register adds a keybinding
func (r *Registry) Register(key Key, action Action) error

// Lookup finds an action for a key in a context
func (r *Registry) Lookup(key Key, context string, mode string) (Action, bool)

// DetectConflicts finds conflicting keybindings
func (r *Registry) DetectConflicts() []Conflict

// Export/Import for user customization
func (r *Registry) Export() map[string]string
func (r *Registry) Import(bindings map[string]string) error
```

### Key Features

1. **Centralized Registry**: All keybindings in one place
2. **Context-Aware Routing**: Different keybindings for different contexts
3. **Vim Mode Support**: Integration with existing vim state machine
4. **Conflict Detection**: Automatically detect conflicting keybindings
5. **User Customization**: Export/import keybindings for user customization
6. **Documentation**: Auto-generate keybinding documentation

### Implementation Plan

#### Phase 1: Core Registry
- [ ] Create `internal/keybindings/registry.go`
- [ ] Implement key representation
- [ ] Implement action representation
- [ ] Implement basic registration and lookup
- [ ] Add tests

#### Phase 2: Context-Aware Routing
- [ ] Add context support (global, workspace, modal, etc.)
- [ ] Add vim mode support (normal, insert, visual)
- [ ] Implement priority-based lookup
- [ ] Add tests

#### Phase 3: Integration
- [ ] Integrate with app model for global keybindings
- [ ] Integrate with workspace model
- [ ] Integrate with modal components
- [ ] Update existing keybinding code to use registry
- [ ] Add tests

#### Phase 4: Advanced Features
- [ ] Implement conflict detection
- [ ] Add keybinding export/import
- [ ] Add user customization support
- [ ] Auto-generate documentation
- [ ] Add tests

#### Phase 5: Migration
- [ ] Migrate all existing keybindings to registry
- [ ] Remove old switch statements
- [ ] Update documentation
- [ ] Add migration guide

### Example Usage

```go
// Initialize registry
registry := keybindings.NewRegistry()

// Register global keybindings
registry.Register(
    keybindings.Key{Type: tea.KeyCtrlC},
    keybindings.Action{
        ID:          "quit",
        Description:  "Quit the application",
        Handler:     func() tea.Cmd { return tea.Quit },
        Context:     "global",
    },
)

// Register workspace keybindings
registry.Register(
    keybindings.Key{Type: tea.KeyRunes, Runes: []rune{'j'}},
    keybindings.Action{
        ID:          "move_down",
        Description:  "Move cursor down",
        Handler:     func() tea.Cmd { return MoveCursorDownMsg{} },
        Context:     "workspace",
        Mode:        "normal",
    },
)

// In Update method
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        key := keybindings.FromKeyMsg(msg)
        action, found := registry.Lookup(key, m.context, m.vimState.CurrentMode)
        if found {
            return m, action.Handler()
        }
    }
    // ... rest of update logic
}
```

## Comparison Summary

| Approach | Pros | Cons | Recommendation |
|----------|------|-------|--------------|
| **Bubble Tea Native** | Zero dependencies, already integrated | No registry, code duplication | ❌ Not scalable |
| **Custom Layer** | Full control, integrates with vim state machine | Requires development effort | ✅ **Recommended** |
| **Bubbles Components** | Official, well-tested | Limited to specific components | ⚠️ Use for components only |
| **Gambit** | Keybinding management built-in | Game-focused, heavy dependency | ❌ Not suitable |
| **Tview** | Comprehensive system | Competing framework | ❌ Incompatible |
| **Tcell** | Low-level control | Too low-level | ❌ Not a solution |
| **Termbox** | Simple API | No keybinding features | ❌ Not suitable |

## Conclusion

After thorough research, **building a custom keybinding management layer on top of Bubble Tea is the best approach** for the Prompt Stack application. This provides:

1. **Perfect Integration**: Works seamlessly with existing Bubble Tea code
2. **Vim Mode Support**: Integrates with the sophisticated vim state machine
3. **Future Flexibility**: Can add user customization, conflict detection, etc.
4. **Maintainability**: Centralized keybinding registry reduces code duplication
5. **No External Dependencies**: Keeps the dependency tree small

The investment in building this layer will pay off in:
- Easier maintenance
- Better user experience (customizable keybindings)
- Reduced bugs (conflict detection)
- Better documentation (auto-generated)

## Next Steps

1. Review and approve this research
2. Create detailed implementation plan
3. Begin Phase 1: Core Registry implementation
4. Iterate through phases with testing at each step
5. Migrate existing keybindings gradually
6. Update documentation

## References

- Bubble Tea Documentation: https://github.com/charmbracelet/bubbletea
- Bubbles Components: https://github.com/charmbracelet/bubbles
- Current Implementation: [`archive/code/ui/`](archive/code/ui/)
- Vim State Machine: [`archive/code/internal/vim/state.go`](archive/code/internal/vim/state.go)