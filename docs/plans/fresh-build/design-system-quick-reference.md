# OpenCode Design System - Quick Reference

## How to Use This Document

**Read this section when:**
- Before implementing any UI component
- When styling TUI elements in PromptStack
- When choosing colors or spacing for components
- When implementing keyboard shortcuts
- During code reviews for UI components

**Key sections:**
- Lines 7-30: Color palette - Reference for all color constants
- Lines 32-36: Spacing system - Reference for layout spacing
- Lines 38-45: Keyboard shortcuts - Default keybindings to support
- Lines 47-100: UI patterns - Common component implementations
- Lines 102-130: Usage examples - How to use theme package

**Related documents:**
- See [opencode-design-system.md](opencode-design-system.md) for complete design specification
- See [milestone-execution-prompt.md](milestone-execution-prompt.md) for design system compliance checks
- See [ui/theme/theme.go](../../ui/theme/theme.go) for theme package implementation

---

## Color Palette (Catppuccin Mocha)

### Background Colors
- `BackgroundPrimary` = "#1e1e2e"   // Main modal background
- `BackgroundSecondary` = "#181825" // Status bar background
- `BackgroundTertiary` = "#313244"  // Secondary button background
- `BackgroundInput` = "#313244"      // Input field background
- `BackgroundHighlight` = "#45475a" // Highlight/background selection
- `BackgroundMuted` = "#1e1e2e"     // Muted background

### Foreground Colors
- `ForegroundPrimary` = "#cdd6f4"   // Main text
- `ForegroundSecondary` = "#a6adc8" // Secondary text
- `ForegroundMuted` = "#6c7086"     // Muted text
- `ForegroundWhite` = "#ffffff"     // White text
- `ForegroundBlack` = "#11111b"     // Black text

### Accent Colors
- `AccentBlue` = "#89b4fa"   // Primary accent, info states
- `AccentGreen` = "#a6e3a1"  // Success states
- `AccentYellow` = "#f9e2af" // Warning states
- `AccentRed` = "#f38ba8"    // Error states
- `AccentCyan` = "#94e2d5"  // Category labels
- `AccentMauve` = "#cba6f7"  // User messages
- `AccentPink` = "#f5c2e7"   // Special highlights
- `AccentOrange` = "#fab387" // Tool names

### Diff Colors
- `DiffAddedBackground` = "#a6e3a1" // Added line background
- `DiffAddedForeground` = "#1e1e2e" // Added line text
- `DiffRemovedBackground` = "#f38ba8" // Removed line background
- `DiffRemovedForeground` = "#1e1e2e" // Removed line text
- `DiffContextBackground` = "#313244" // Context line background
- `DiffContextForeground` = "#a6adc8" // Context line text
- `DiffHeaderBackground` = "#45475a" // Diff header background
- `DiffHeaderForeground` = "#89b4fa" // Diff header text

### Border Colors
- `BorderPrimary` = "#45475a" // Modal borders
- `BorderMuted` = "#585b70"   // Secondary borders
- `BorderFocus` = "#89b4fa"   // Focused element border

### Special Colors
- `CursorBackground` = "7"   // Cursor highlight
- `CursorForeground` = "0"   // Cursor text
- `TextMuted` = "#6c7086"   // Dimmed text

---

## Spacing System

**Unit**: 1 character width/height (monospace grid)

- **Padding**: 1-2 units for content padding
- **Margin**: 1 unit between related elements, 2 units between sections
- **Gutter**: 1 unit for column spacing

---

## Keyboard Shortcuts (OpenCode Defaults)

- `Ctrl+X` - Leader key for commands
- `Ctrl+P` - Command palette
- `Tab` - Switch between agents (build/plan)
- `Esc` - Cancel/close dialogs
- `Ctrl+C` - Interrupt running operations
- `Ctrl+Z` - Suspend to background
- `Ctrl+G` - Cancel popovers and abort responses

---

## Common UI Patterns

### Modal Dialogs

```go
dialogStyle := theme.ModalStyle().
    Width(60).
    Height(20).
    Padding(1, 2)

titleStyle := theme.ModalTitleStyle()
contentStyle := theme.ModalContentStyle()
```

### Status Bars

```go
statusStyle := theme.StatusStyle().
    Width(screenWidth).
    Height(1).
    Padding(0, 1)
```

### Input Fields

```go
inputStyle := theme.InputStyle().
    Width(50).
    Placeholder("Search...")
```

### Lists

```go
itemStyle := theme.ListItemStyle()
selectedStyle := theme.ListItemSelectedStyle()
categoryStyle := theme.ListCategoryStyle()
```

---

## Usage Examples

### Using Style Helpers

```go
import "github.com/kyledavis/prompt-stack/ui/theme"

// Simple usage
modal := theme.ModalStyle()
status := theme.StatusStyle()
input := theme.InputStyle()

// Customizing
customModal := theme.ModalStyle().
    Width(80).
    Height(40).
    Padding(2, 3)
```

### Using Color Constants

```go
import (
    "github.com/charmbracelet/lipgloss"
    "github.com/kyledavis/prompt-stack/ui/theme"
)

customStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color(theme.ForegroundPrimary)).
    Background(lipgloss.Color(theme.BackgroundPrimary)).
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color(theme.BorderPrimary))
```

### Using Spacing Constants

```go
import "github.com/kyledavis/prompt-stack/ui/theme"

// 1 unit spacing between elements
margin := theme.Unit
padding := theme.Unit * 2
gutter := theme.Unit
```

### Using Keyboard Shortcuts

```go
import "github.com/kyledavis/prompt-stack/ui/theme"

// Show in help text
helpText := fmt.Sprintf("Press %s for commands", theme.LeaderKey)
helpText += fmt.Sprintf("Press %s to close", theme.CancelKey)
```

---

## Style Functions Available

### Modal Styles
- `ModalStyle()` - Base modal container style
- `ModalTitleStyle()` - Modal title text style
- `ModalContentStyle()` - Modal content text style
- `ModalButtonStyle()` - Modal button style
- `ModalButtonFocusedStyle()` - Focused button style

### Status Styles
- `StatusStyle()` - Base status bar style
- `InfoStyle()` - Info status style
- `SuccessStyle()` - Success status style
- `WarningStyle()` - Warning status style
- `ErrorStyle()` - Error status style

### Input Styles
- `InputStyle()` - Base input field style
- `SearchInputStyle()` - Search input field style

### List Styles
- `ListItemStyle()` - List item style
- `ListItemSelectedStyle()` - Selected item style
- `ListCategoryStyle()` - Category label style

### Preview Styles
- `PreviewStyle()` - Base preview pane style
- `PreviewTitleStyle()` - Preview title style

### Utility Styles
- `HeaderStyle()` - Section header style
- `CursorStyle()` - Cursor highlight style
- `ActivePlaceholderStyle()` - Active placeholder style
- `ValidationErrorStyle()` - Validation error style
- `MutedStyle()` - Dimmed text style
- `HighlightStyle()` - Highlighted text style

### Diff Styles
- `DiffAddedStyle()` - Added line style (green background)
- `DiffRemovedStyle()` - Removed line style (red background)
- `DiffContextStyle()` - Context line style
- `DiffHeaderStyle()` - Diff header style

### Message Styles
- `UserMessageStyle()` - User message style (blue)
- `AssistantMessageStyle()` - Assistant message style (cyan)
- `SystemMessageStyle()` - System message style (italic, muted)

### Tool Styles
- `ToolNameStyle()` - Tool name style (orange)
- `ToolOutputStyle()` - Tool output style

---

## Design Compliance Checklist

Before committing UI code, verify:

- [ ] Colors used from OpenCode palette only
- [ ] Spacing uses 1-character unit system
- [ ] Style functions used instead of inline styles
- [ ] Keyboard shortcuts match OpenCode defaults
- [ ] Visual hierarchy follows OpenCode patterns
- [ ] Interactive elements have clear visual feedback
- [ ] Error states use AccentRed color
- [ ] Success states use AccentGreen color
- [ ] Warnings use AccentYellow color
- [ ] Info states use AccentBlue color

---

**Last Updated**: 2026-01-08  
**Status**: Active - Use this reference for all UI development
