# OpenCode Design System for PromptStack TUI

## Overview

This document outlines the design principles, patterns, and implementation guidelines for building a TUI application that follows OpenCode's beautiful design aesthetic. OpenCode is known for its clean, modern, and highly functional terminal interface that balances aesthetics with usability.

## Core Design Principles

### 1. Minimalist Aesthetic
- Clean, uncluttered interfaces with ample whitespace
- Consistent spacing and alignment
- Subtle visual hierarchy through color and typography
- No unnecessary decorations or visual noise

### 2. Functional Beauty
- Every visual element serves a purpose
- Colors indicate state and function (not just decoration)
- Typography enhances readability and information hierarchy
- Interactive elements are clearly distinguishable

### 3. Terminal-First Design
- Respect terminal constraints while pushing boundaries
- Optimize for monospace fonts and limited color palettes
- Ensure compatibility across different terminal emulators
- Leverage truecolor support for richer visuals

### 4. Consistency
- Uniform spacing and padding throughout
- Consistent color usage patterns
- Predictable interaction patterns
- Cohesive visual language across all components

## Color System

### Base Color Palette (Catppuccin Mocha)
OpenCode uses the Catppuccin Mocha palette as its foundation:

```go
// Background Colors
BackgroundPrimary   = "#1e1e2e"  // Main modal background
BackgroundSecondary = "#181825"  // Status bar background  
BackgroundTertiary  = "#313244"  // Secondary button background
BackgroundInput     = "236"      // Input field background (terminal color)

// Foreground Colors
ForegroundPrimary   = "#cdd6f4"  // Main text
ForegroundSecondary = "#a6adc8"  // Secondary text
ForegroundMuted     = "#6c7086"  // Muted text
ForegroundWhite     = "15"       // White text (terminal color)

// Accent Colors
AccentBlue   = "#89b4fa"  // Primary accent, info states
AccentGreen  = "#a6e3a1"  // Success states
AccentYellow = "#f9e2af"  // Warning states
AccentRed    = "#f38ba8"  // Error states
AccentCyan   = "39"       // Category labels (terminal color)

// Border Colors
BorderPrimary = "#45475a"  // Modal borders
BorderMuted   = "240"      // Secondary borders (terminal color)

// Special Colors
CursorBackground = "7"   // Cursor highlight
CursorForeground = "0"   // Cursor text
TextMuted        = "245" // Dimmed text (terminal color)
```

### Color Usage Guidelines

1. **Primary Background**: Use for main content areas and modal dialogs
2. **Secondary Background**: Use for status bars, footers, and secondary panels
3. **Primary Text**: Default text color for most content
4. **Secondary Text**: Labels, metadata, less important information
5. **Accent Colors**: Use consistently for specific states:
   - Blue: Primary actions, links, information, user messages
   - Green: Success states, positive actions, added diff lines
   - Yellow: Warnings, caution states, highlighted search terms
   - Red: Errors, destructive actions, removed diff lines
   - Cyan: Categories, tags, metadata, assistant messages

### Specific Color Applications

#### Diff Display Colors
- **Added lines**: Green background with darker green text
- **Removed lines**: Red background with darker red text  
- **Context lines**: Muted background with secondary text
- **Line numbers**: Border color with muted text
- **Diff headers**: Accent cyan with bold text

#### Tool Execution Colors
- **Tool name**: Accent blue with brackets `[Tool: name]`
- **Command**: Primary text with subtle highlighting
- **Output**: Default terminal colors respecting syntax
- **Success indicator**: Green text with checkmark
- **Error indicator**: Red text with cross mark

#### Status Indicators
- **Modified files**: Yellow warning color
- **Read-only mode**: Cyan accent color
- **Placeholder edit**: Inverse cursor colors
- **Loading state**: Pulsing accent color
- **Error state**: Red with clear icon

#### Search & Highlighting
- **Matched text**: Yellow background with dark text
- **Current selection**: Blue background with white text
- **File matches**: Cyan accent in file browser
- **Command matches**: Blue accent in palette

#### Message Types
- **User messages**: Blue accent with clear label
- **Assistant messages**: Cyan accent with clear label
- **System messages**: Muted text with brackets
- **Tool results**: Border color with collapsible header

## Typography & Layout

### Spacing System
- **Unit**: 1 character width/height (monospace grid)
- **Padding**: 1-2 units for content padding
- **Margin**: 1 unit between related elements, 2 units between sections
- **Gutter**: 1 unit for column spacing

### Layout Patterns

#### 1. Modal Dialogs
- Rounded borders with subtle shadows (via border styling)
- Centered or positioned based on context
- Minimum padding of 1 unit on all sides
- Clear visual separation from background

#### 2. Status Bars
- Full width at bottom of screen
- Secondary background color
- Compact information display
- Right-aligned status indicators
- Multiple information sections separated by `|`
- Dynamic content based on context

#### 3. Input Fields
- Distinct background color (BackgroundInput)
- Clear visual focus state
- Placeholder text in muted color
- Consistent padding and borders
- Multi-line support with `Shift+Enter`
- Command history navigation with arrow keys

#### 4. Lists & Menus
- Clear visual hierarchy
- Hover/focus states for interactive items
- Consistent indentation and alignment
- Visual indicators for selected items
- Fuzzy search filtering
- Virtual scrolling for long lists

#### 5. Chat Interface Layout
- Message bubbles with sender identification
- Tool results in collapsible sections
- Code blocks with syntax highlighting
- Inline file references with preview
- Timestamps for message history
- User/assistant role indicators

#### 6. Multi-pane Layouts
- Sidebar navigation (toggle with `Ctrl+X B`)
- Main content area with scrolling
- Status bar at bottom
- Adjustable pane widths
- Responsive collapse on small terminals

#### 7. Diff View Layouts
- Side-by-side comparison (wide terminals)
- Inline diff display (narrow terminals)
- Collapsible file sections
- Line number gutter
- Change statistics summary

#### 8. Search Interface Layout
- Overlay or inline search panels
- Real-time results filtering
- Category filters and tags
- Recent searches history
- Search syntax help

## Component Design Patterns

### 1. Buttons & Interactive Elements
- Clear visual feedback on hover/focus
- Consistent padding and borders
- Color coding based on action type
- Disabled state with reduced opacity

### 2. Forms & Inputs
- Clear labels using secondary text color
- Validation states with appropriate colors
- Help text in muted color
- Consistent spacing between form elements

### 3. Status Indicators
- Color-coded based on status type
- Compact, information-dense display
- Consistent positioning
- Clear visual hierarchy

### 4. Navigation Elements
- Clear visual hierarchy
- Consistent spacing
- Active state clearly indicated
- Breadcrumb-style navigation where appropriate

### 5. File Reference Components
- `@` prefix for file references
- Fuzzy search with real-time filtering
- Visual indication of matched characters
- Tab cycling through results
- Enter to select and insert reference

### 6. Tool Execution Components
- Clear tool name and command display
- Collapsible output sections
- Success/error state indicators
- Execution time display
- Copy-to-clipboard functionality for outputs

### 7. Diff Display Components
- Side-by-side or inline diff views
- Color coding: green for additions, red for deletions
- Line numbers for context
- Collapsible diff sections
- File path and change summary headers

### 8. Loading & Progress Indicators
- Subtle spinners for short operations
- Progress bars for longer operations
- Status messages with estimated time
- Cancellable operations with clear feedback
- Skeleton loading states for content

### 9. Error States & Recovery
- Clear error messages with context
- Suggested recovery actions
- Error codes for debugging
- Retry mechanisms where appropriate
- Fallback content when available

### 10. Search & Filter Interfaces
- Command palette (Ctrl+P) for actions
- Fuzzy search with highlighting
- Filter chips for categories
- Clear search state indicators
- Recent searches history

## UI Patterns from OpenCode

### 1. Chat Interface Pattern
```
┌─────────────────────────────────────────────────────┐
│                                                     │
│  User: Can you help me with this code?              │
│                                                     │
│  Assistant: Sure! I'd be happy to help.             │
│  Let me look at the code...                         │
│                                                     │
│  [Tool Result: File read successfully]              │
│                                                     │
│  Assistant: I found the issue. Here's the fix...    │
│                                                     │
└─────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────┐
│ Type your message...                                │
└─────────────────────────────────────────────────────┘
```

### 2. Status Bar Pattern
```
┌─────────────────────────────────────────────────────┐
│ Chars: 125 | Lines: 8 | Modified | [PLACEHOLDER EDIT]│
└─────────────────────────────────────────────────────┘
```

### 3. Modal Dialog Pattern
```
┌─────────────────────────────────────────────────────┐
│                                                     │
│  ╭─────────────────────────────────────────────╮    │
│  │          Select Theme                      │    │
│  │                                             │    │
│  │  • opencode (current)                      │    │
│  │    tokyonight                              │    │
│  │    catppuccin                              │    │
│  │    nord                                    │    │
│  │                                             │    │
│  │  Press Enter to select, Esc to cancel      │    │
│  ╰─────────────────────────────────────────────╯    │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### 4. File Reference Pattern (@ fuzzy search)
```
┌─────────────────────────────────────────────────────┐
│ Type: @                                             │
│                                                     │
│  Search files...                                    │
│  > src/main.go                                      │
│    src/utils/file.go                                │
│    internal/config/config.go                        │
│    cmd/server/main.go                               │
│                                                     │
│  Press Tab to cycle, Enter to select                │
└─────────────────────────────────────────────────────┘
```

### 5. Tool Execution Pattern
```
┌─────────────────────────────────────────────────────┐
│ [Tool: bash] Running: ls -la                        │
│                                                     │
│ total 48                                            │
│ drwxr-xr-x  11 user  staff   352 Jan  8 12:00 .     │
│ drwxr-xr-x   5 user  staff   160 Jan  8 11:30 ..    │
│ -rw-r--r--   1 user  staff   287 Jan  8 11:45 .git  │
│ -rw-r--r--   1 user  staff    45 Jan  8 11:45 .git  │
│                                                     │
│ [Tool: bash] Completed successfully                 │
└─────────────────────────────────────────────────────┘
```

### 6. Diff Display Pattern
```
┌─────────────────────────────────────────────────────┐
│ @@ -1,5 +1,5 @@                                     │
│ -func oldFunction() {                               │
│ +func newFunction() {                               │
│     // Old comment                                  │
│ -    return oldValue                                │
│ +    return newValue                                │
│ }                                                   │
│                                                     │
│ Green: Added lines                                  │
│ Red: Removed lines                                  │
│ Gray: Context lines                                 │
└─────────────────────────────────────────────────────┘
```

### 7. Command Palette Pattern (Ctrl+P)
```
┌─────────────────────────────────────────────────────┐
│ >                                                   │
│                                                     │
│  /connect    Add a provider                         │
│  /theme      Change theme                           │
│  /models     List available models                  │
│  /help       Show help                              │
│  /export     Export conversation                    │
│  /share      Share current session                  │
│                                                     │
│  Type to filter commands...                         │
└─────────────────────────────────────────────────────┘
```

### 8. Agent Switching Indicator
```
┌─────────────────────────────────────────────────────┐
│                                                     │
│  [BUILD] Ready to make changes                      │
│                                                     │
│  or                                                 │
│                                                     │
│  [PLAN] Read-only analysis mode                     │
│                                                     │
│  Press Tab to switch agents                         │
└─────────────────────────────────────────────────────┘
```

## Implementation Guidelines

### 1. Using Lipgloss for Styling
```go
import "github.com/charmbracelet/lipgloss"

// Base styles should be defined in theme package
func ModalStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Background(lipgloss.Color(BackgroundPrimary)).
        Foreground(lipgloss.Color(ForegroundPrimary)).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color(BorderPrimary))
}

// Component-specific styles extend base styles
func ButtonStyle() lipgloss.Style {
    return ModalStyle().
        Padding(0, 2).
        Background(lipgloss.Color(BackgroundTertiary)).
        BorderStyle(lipgloss.NormalBorder())
}
```

### 2. Consistent Component Structure
```go
type Component struct {
    width  int
    height int
    // Component-specific state
}

func (c Component) View() string {
    // Use theme styles consistently
    style := theme.ComponentStyle().
        Width(c.width).
        Height(c.height)
    
    // Build content with proper spacing
    content := lipgloss.JoinVertical(
        lipgloss.Left,
        c.renderHeader(),
        c.renderContent(),
        c.renderFooter(),
    )
    
    return style.Render(content)
}
```

### 3. Responsive Layout
```go
func (m Model) View() string {
    // Calculate available space
    availableHeight := m.height - statusBarHeight
    
    // Adjust components based on available space
    if availableHeight < minHeight {
        return m.renderCompactView()
    }
    
    return m.renderFullView()
}
```

## Interactive Patterns

### 1. Focus Management
- Clear visual indication of focused element
- Tab-based navigation between interactive elements
- Escape key to cancel/close dialogs
- Enter key to activate primary action

### 2. Keyboard Shortcuts (OpenCode Defaults)
- `Ctrl+X` as leader key for commands
- `Ctrl+P` for command palette
- `Tab` to switch between agents (build/plan)
- `Esc` to cancel operations and close dialogs
- `Ctrl+C` to interrupt running operations
- `Ctrl+Z` to suspend to background
- `Ctrl+G` to cancel popovers and abort responses

### 3. Text Editing Shortcuts (Readline/Emacs-style)
- `Ctrl+A` - Move to start of line
- `Ctrl+E` - Move to end of line
- `Ctrl+B` - Move cursor back one character
- `Ctrl+F` - Move cursor forward one character
- `Alt+B` - Move cursor back one word
- `Alt+F` - Move cursor forward one word
- `Ctrl+D` - Delete character under cursor
- `Ctrl+K` - Kill to end of line
- `Ctrl+U` - Kill to start of line
- `Ctrl+W` - Kill previous word
- `Alt+D` - Kill next word
- `Ctrl+T` - Transpose characters

### 4. Visual Feedback
- Immediate visual response to user actions
- Loading states with subtle animations
- Success/error states with appropriate colors
- Hover states for interactive elements
- Scroll acceleration for smooth scrolling (macOS-style)
- Tool execution details toggle (`Ctrl+X D`)

### 5. File Operations
- `@` for fuzzy file search and reference
- Drag and drop images into terminal
- Auto-formatting after file edits
- Git integration for undo/redo operations
- File watcher for external changes

### 6. Session Management
- Multiple concurrent sessions
- Session switching with `Ctrl+X L`
- Session export to markdown
- Session sharing via links
- Session compaction to save context

## Accessibility Considerations

### 1. Color Contrast
- Ensure sufficient contrast between text and background
- Use color combinations that work in grayscale
- Provide alternative indicators beyond color alone

### 2. Keyboard Navigation
- All functionality accessible via keyboard
- Logical tab order
- Clear focus indicators
- Escape routes from all dialogs

### 3. Screen Reader Compatibility
- Semantic structure in rendered output
- Clear status messages
- Announce important state changes

## Testing & Quality Assurance

### 1. Visual Consistency
- Check spacing and alignment across components
- Verify color usage follows guidelines
- Ensure typography hierarchy is maintained

### 2. Responsive Behavior
- Test at different terminal sizes
- Verify layout adjustments work correctly
- Check for content truncation issues

### 3. Interactive Behavior
- Test keyboard navigation flows
- Verify focus management
- Check visual feedback for all interactions

## Theme System & Customization

### Built-in Themes
OpenCode supports multiple built-in themes that can be switched dynamically:
- `opencode` - Default OpenCode theme (Catppuccin Mocha)
- `system` - Adapts to terminal's color scheme
- `tokyonight` - Based on Tokyo Night theme
- `catppuccin` - Catppuccin variants
- `nord` - Nord theme
- `gruvbox` - Gruvbox theme
- `ayu` - Ayu dark theme
- `everforest` - Everforest theme
- `kanagawa` - Kanagawa theme
- `matrix` - Hacker-style green on black
- `one-dark` - Atom One Dark theme

### Theme Configuration
Themes support JSON configuration with:
- **Hex colors**: `"#ffffff"`
- **ANSI colors**: `3` (0-255)
- **Color references**: `"primary"` or custom definitions
- **Dark/light variants**: `{"dark": "#000", "light": "#fff"}`
- **No color**: `"none"` - Uses terminal's default

### Custom Theme Structure
```json
{
  "$schema": "https://opencode.ai/theme.json",
  "defs": {
    "customColor": "#ff0000"
  },
  "theme": {
    "primary": {
      "dark": "customColor",
      "light": "#0000ff"
    },
    "text": {
      "dark": "#cdd6f4",
      "light": "#1e1e2e"
    },
    "background": {
      "dark": "#1e1e2e",
      "light": "#ffffff"
    }
  }
}
```

### Theme Loading Hierarchy
1. **Built-in themes** - Embedded in binary
2. **User config directory** - `~/.config/opencode/themes/*.json`
3. **Project root directory** - `.opencode/themes/*.json`
4. **Current working directory** - `./.opencode/themes/*.json`

## Implementation Roadmap

### Phase 1: Foundation
1. Extend theme system with OpenCode color palette
2. Implement base component styles (modal, status bar, input)
3. Create consistent spacing utilities
4. Add theme switching support

### Phase 2: Core Components
1. Chat interface component with message bubbles
2. Enhanced status bar with multiple indicators
3. Modal dialog system with command palette
4. Form input components with validation
5. File reference (`@`) fuzzy search

### Phase 3: Advanced Features
1. Tabbed interfaces and multi-pane layouts
2. File tree navigation and browser
3. Search/filter components with highlighting
4. Progress indicators and loading states
5. Diff display with side-by-side/inline views
6. Tool execution results display

### Phase 4: Polish
1. Animations and transitions (spinners, fades)
2. Enhanced keyboard navigation and shortcuts
3. Scroll acceleration and smooth scrolling
4. Accessibility improvements
5. Error states and recovery patterns
6. Performance optimizations for large outputs

## Example Implementations

### 1. OpenCode-style Modal Dialog
```go
package dialog

import (
    "github.com/charmbracelet/lipgloss"
    "github.com/kyledavis/prompt-stack/ui/theme"
)

type Dialog struct {
    title   string
    content string
    width   int
    height  int
}

func New(title, content string) Dialog {
    return Dialog{
        title:   title,
        content: content,
        width:   60,
        height:  20,
    }
}

func (d Dialog) View() string {
    // Create dialog style
    dialogStyle := lipgloss.NewStyle().
        Width(d.width).
        Height(d.height).
        Background(lipgloss.Color(theme.BackgroundPrimary)).
        Foreground(lipgloss.Color(theme.ForegroundPrimary)).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color(theme.BorderPrimary)).
        Padding(1, 2)
    
    // Create title style
    titleStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color(theme.AccentBlue)).
        Bold(true).
        PaddingBottom(1)
    
    // Build content
    content := lipgloss.JoinVertical(
        lipgloss.Left,
        titleStyle.Render(d.title),
        d.content,
    )
    
    return dialogStyle.Render(content)
}
```

### 2. Status Bar Component
```go
package statusbar

import (
    "strings"
    "github.com/charmbracelet/lipgloss"
    "github.com/kyledavis/prompt-stack/ui/theme"
)

type StatusBar struct {
    width      int
    indicators []string
}

func New(width int) *StatusBar {
    return &StatusBar{
        width:      width,
        indicators: []string{},
    }
}

func (s *StatusBar) AddIndicator(text string) {
    s.indicators = append(s.indicators, text)
}

func (s *StatusBar) Clear() {
    s.indicators = []string{}
}

func (s *StatusBar) View() string {
    if len(s.indicators) == 0 {
        return ""
    }
    
    // Join indicators with separator
    statusText := strings.Join(s.indicators, " | ")
    
    // Apply status bar styling
    style := lipgloss.NewStyle().
        Width(s.width).
        Height(1).
        Background(lipgloss.Color(theme.BackgroundSecondary)).
        Foreground(lipgloss.Color(theme.ForegroundPrimary)).
        Padding(0, 1)
    
    return style.Render(statusText)
}
```

### 3. File Reference Search Component
```go
package filesearch

import (
    "strings"
    "github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/kyledavis/prompt-stack/ui/theme"
)

type FileSearch struct {
    input     textinput.Model
    files     []string
    matches   []string
    selected  int
    visible   bool
}

func New() FileSearch {
    ti := textinput.New()
    ti.Placeholder = "Search files..."
    ti.Prompt = "@"
    ti.CharLimit = 256
    ti.Width = 50
    
    return FileSearch{
        input:    ti,
        files:    []string{},
        matches:  []string{},
        selected: 0,
        visible:  false,
    }
}

func (fs FileSearch) Update(msg tea.Msg) (FileSearch, tea.Cmd) {
    var cmd tea.Cmd
    
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEsc:
            fs.visible = false
            fs.input.SetValue("")
            return fs, nil
        case tea.KeyEnter:
            if len(fs.matches) > 0 && fs.selected < len(fs.matches) {
                // Return selected file
                return fs, SelectFileCmd(fs.matches[fs.selected])
            }
        case tea.KeyTab:
            if len(fs.matches) > 0 {
                fs.selected = (fs.selected + 1) % len(fs.matches)
            }
        case tea.KeyShiftTab:
            if len(fs.matches) > 0 {
                fs.selected = (fs.selected - 1 + len(fs.matches)) % len(fs.matches)
            }
        case tea.KeyUp:
            if fs.selected > 0 {
                fs.selected--
            }
        case tea.KeyDown:
            if fs.selected < len(fs.matches)-1 {
                fs.selected++
            }
        }
    }
    
    fs.input, cmd = fs.input.Update(msg)
    
    // Update matches based on input
    if fs.input.Value() != "" {
        fs.matches = fs.filterFiles(fs.input.Value())
    } else {
        fs.matches = fs.files
    }
    
    return fs, cmd
}

func (fs FileSearch) View() string {
    if !fs.visible {
        return ""
    }
    
    var lines []string
    lines = append(lines, fs.input.View())
    lines = append(lines, "")
    
    for i, file := range fs.matches {
        if i >= 10 { // Limit display
            lines = append(lines, "...")
            break
        }
        
        line := file
        if i == fs.selected {
            line = "> " + lipgloss.NewStyle().
                Foreground(lipgloss.Color(theme.AccentBlue)).
                Render(file)
        } else {
            line = "  " + file
        }
        lines = append(lines, line)
    }
    
    if len(fs.matches) == 0 {
        lines = append(lines, "No matches found")
    }
    
    lines = append(lines, "")
    lines = append(lines, "Press Tab to cycle, Enter to select, Esc to cancel")
    
    content := strings.Join(lines, "\n")
    
    // Style the search box
    style := lipgloss.NewStyle().
        Width(60).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color(theme.BorderPrimary)).
        Padding(1).
        Background(lipgloss.Color(theme.BackgroundPrimary))
    
    return style.Render(content)
}

func (fs FileSearch) filterFiles(query string) []string {
    query = strings.ToLower(query)
    var matches []string
    
    for _, file := range fs.files {
        if strings.Contains(strings.ToLower(file), query) {
            matches = append(matches, file)
        }
    }
    
    return matches
}
```

## Key Design Insights from OpenCode

### 1. Terminal-First Philosophy
OpenCode demonstrates that terminal applications can be both beautiful and functional. The design respects terminal constraints while pushing boundaries with:
- Truecolor support for rich visuals
- Smooth scrolling and animations
- Responsive layouts that adapt to terminal size
- Keyboard-first interaction patterns

### 2. Information Density & Hierarchy
OpenCode balances information density with clarity through:
- Consistent color coding for different information types
- Collapsible sections for detailed content
- Progressive disclosure of complex information
- Clear visual hierarchy with spacing and typography

### 3. Interactive Patterns
The user experience is built around familiar patterns:
- Command palette (Ctrl+P) for discoverability
- Fuzzy search for file references (@)
- Readline/Emacs-style text editing
- Modal interactions with clear escape routes
- Multi-session management with easy switching

### 4. Developer-Centric Design
Every design decision serves developer workflows:
- Integrated diff views for code changes
- Tool execution with detailed outputs
- File reference system for context
- Agent switching for different tasks (build/plan)
- Git integration for undo/redo operations

## Implementation Priorities

### High Priority (Core UX)
1. **Theme system** with OpenCode color palette
2. **Chat interface** with message bubbles and tool results
3. **Status bar** with dynamic indicators
4. **File reference** (@) fuzzy search
5. **Keyboard shortcuts** with leader key (Ctrl+X)

### Medium Priority (Enhanced Features)
1. **Command palette** (Ctrl+P) for actions
2. **Diff display** with color-coded changes
3. **Multi-pane layouts** with sidebar
4. **Loading states** and progress indicators
5. **Error handling** with recovery options

### Low Priority (Polish)
1. **Smooth scrolling** with acceleration
2. **Animations** for state transitions
3. **Advanced theming** with custom themes
4. **Accessibility** enhancements
5. **Performance optimizations**

## Testing & Validation

### Visual Consistency Checks
- Verify color usage follows guidelines across all components
- Test spacing and alignment at different terminal sizes
- Ensure typography hierarchy is maintained
- Check contrast ratios for accessibility

### Interactive Behavior Tests
- Verify all keyboard shortcuts work correctly
- Test focus management and tab order
- Validate error states and recovery flows
- Check responsive behavior on resize

### Performance Considerations
- Optimize rendering for large outputs
- Implement virtual scrolling for long lists
- Cache frequently accessed data
- Monitor memory usage with large sessions

## Conclusion

By following these design principles and implementation guidelines, we can create a TUI application that matches OpenCode's beautiful aesthetic while maintaining consistency, usability, and accessibility. The key is to balance visual appeal with functional design, ensuring every element serves a purpose and contributes to a cohesive user experience.

Remember: Good design is invisible. Users should feel the beauty and functionality without noticing the individual design decisions that make it work.

The OpenCode design system demonstrates that terminal applications don't have to sacrifice aesthetics for functionality. By implementing these patterns, PromptStack can achieve the same level of polish and usability that makes OpenCode such a pleasure to use.