# Theme System

The theme package provides a centralized styling system for PromptStack TUI using the Catppuccin Mocha color palette.

## Usage

### Using Style Helpers

```go
import "github.com/kyledavis/prompt-stack/ui/theme"

// Apply modal style
modal := theme.ModalStyle().Width(80).Height(40)

// Apply status bar style
status := theme.StatusStyle().Width(100)

// Apply active placeholder style
placeholder := theme.ActivePlaceholderStyle()
```

### Using Color Constants

```go
import (
    "github.com/charmbracelet/lipgloss"
    "github.com/kyledavis/prompt-stack/ui/theme"
)

// Custom style using theme colors
custom := lipgloss.NewStyle().
    Foreground(lipgloss.Color(theme.ForegroundPrimary)).
    Background(lipgloss.Color(theme.BackgroundPrimary))
```

## Color Palette

The theme uses the Catppuccin Mocha color palette:

### Background Colors
- `BackgroundPrimary` (#1e1e2e) - Main modal background
- `BackgroundSecondary` (#181825) - Status bar background
- `BackgroundTertiary` (#313244) - Secondary button background
- `BackgroundInput` (236) - Input field background (terminal color)

### Foreground Colors
- `ForegroundPrimary` (#cdd6f4) - Main text
- `ForegroundSecondary` (#a6adc8) - Secondary text
- `ForegroundMuted` (#6c7086) - Muted text
- `ForegroundWhite` (15) - White text (terminal color)

### Accent Colors
- `AccentBlue` (#89b4fa) - Primary accent, info states
- `AccentGreen` (#a6e3a1) - Success states
- `AccentYellow` (#f9e2af) - Warning states
- `AccentRed` (#f38ba8) - Error states
- `AccentCyan` (39) - Category labels (terminal color)

### Border Colors
- `BorderPrimary` (#45475a) - Modal borders
- `BorderMuted` (240) - Secondary borders (terminal color)

### Special Colors
- `CursorBackground` (7) - Cursor highlight
- `CursorForeground` (0) - Cursor text
- `TextMuted` (245) - Dimmed text (terminal color)

## Guidelines

- ✅ Always use theme helpers for standard UI patterns
- ✅ Use color constants for custom styles
- ✅ Extend base styles for component-specific needs
- ❌ Don't hard-code hex colors or terminal color numbers
- ❌ Don't create new color constants without consulting the design system

## Style Helpers

### ModalStyle
Returns the base style for modal containers with rounded borders and primary colors.

### StatusStyle
Returns the style for status bar components with secondary background.

### ActivePlaceholderStyle
Returns the style for active placeholder highlighting with cursor colors.