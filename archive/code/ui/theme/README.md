# Theme System

A centralized theme system for managing lipgloss styling across the prompt-stack UI.

## Overview

The theme system provides a single source of truth for all UI colors and styles, making it easy to maintain consistency and update the visual appearance of the application.

## Features

- **Centralized color constants** - All colors defined in one place
- **Pre-built style helpers** - Common UI patterns ready to use
- **Catppuccin Mocha theme** - Beautiful, modern color palette
- **Type-safe** - Constants prevent typos in color values
- **Easy to extend** - Base styles can be customized for specific needs

## Usage

### Import the Theme Package

```go
import "github.com/kyledavis/prompt-stack/ui/theme"
```

### Using Style Helpers

Replace hard-coded lipgloss styles with theme helper functions:

**Before:**
```go
modalStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#45475a")).
    Background(lipgloss.Color("#1e1e2e"))
```

**After:**
```go
modalStyle := theme.ModalStyle()
```

### Available Style Helpers

#### Modal Styles
- `ModalStyle()` - Base modal container
- `ModalTitleStyle()` - Modal titles
- `ModalContentStyle()` - Modal content text
- `ModalButtonStyle()` - Primary buttons
- `ModalButtonFocusedStyle()` - Focused primary buttons
- `ModalButtonSecondaryStyle()` - Secondary buttons
- `ModalButtonSecondaryFocusedStyle()` - Focused secondary buttons

#### Status Bar Styles
- `StatusStyle()` - Normal status messages
- `InfoStyle()` - Info messages (blue)
- `SuccessStyle()` - Success messages (green)
- `WarningStyle()` - Warning messages (yellow)
- `ErrorStyle()` - Error messages (red)
- `SeparatorStyle()` - Status bar separators

#### Input Styles
- `InputStyle()` - Base input field
- `SearchInputStyle()` - Search input field

#### List Styles
- `ListItemStyle()` - Unselected list items
- `ListItemSelectedStyle()` - Selected list items
- `ListCategoryStyle()` - Category labels
- `ListDescriptionStyle()` - Item descriptions
- `ListEmptyStyle()` - Empty state message

#### Preview Styles
- `PreviewStyle()` - Preview pane container
- `PreviewTitleStyle()` - Preview titles
- `PreviewDescriptionStyle()` - Preview descriptions
- `PreviewTagsStyle()` - Preview tags
- `PreviewContentStyle()` - Preview content

#### Cursor Styles
- `CursorStyle()` - Cursor highlighting

#### Header Styles
- `HeaderStyle()` - Section headers

### Using Color Constants

You can also use color constants directly for custom styles:

```go
import "github.com/charmbracelet/lipgloss"
import "github.com/kyledavis/prompt-stack/ui/theme"

customStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color(theme.ForegroundPrimary)).
    Background(lipgloss.Color(theme.BackgroundPrimary)).
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color(theme.BorderPrimary))
```

### Available Color Constants

#### Background Colors
- `BackgroundPrimary` - Main modal background (#1e1e2e)
- `BackgroundSecondary` - Status bar background (#181825)
- `BackgroundTertiary` - Secondary button background (#313244)
- `BackgroundInput` - Input field background (terminal color 236)

#### Foreground Colors
- `ForegroundPrimary` - Main text (#cdd6f4)
- `ForegroundSecondary` - Secondary text (#a6adc8)
- `ForegroundMuted` - Muted text (#6c7086)
- `ForegroundWhite` - White text (terminal color 15)

#### Accent Colors
- `AccentBlue` - Primary accent (#89b4fa)
- `AccentGreen` - Success (#a6e3a1)
- `AccentYellow` - Warning (#f9e2af)
- `AccentRed` - Error (#f38ba8)
- `AccentCyan` - Category labels (terminal color 39)

#### Border Colors
- `BorderPrimary` - Modal borders (#45475a)
- `BorderMuted` - Secondary borders (terminal color 240)

#### Special Colors
- `CursorBackground` - Cursor highlight (terminal color 7)
- `CursorForeground` - Cursor text (terminal color 0)
- `TextMuted` - Dimmed text (terminal color 245)

### Extending Base Styles

For custom components, extend the base styles:

```go
customModal := theme.BaseModal().
    Width(80).
    Height(40).
    Padding(2, 3)

customInput := theme.BaseInput().
    Width(60).
    Placeholder("Enter text...")
```

## Color Palette

The theme uses the **Catppuccin Mocha** color scheme, a popular dark theme with excellent contrast and readability.

### Theme Colors

| Category | Color | Hex | Usage |
|----------|-------|-----|-------|
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

## Migration Guide

To migrate existing code to use the theme system:

1. Import the theme package:
   ```go
   import "github.com/kyledavis/prompt-stack/ui/theme"
   ```

2. Find hard-coded lipgloss styles:
   ```go
   var style = lipgloss.NewStyle().
       Foreground(lipgloss.Color("#cdd6f4")).
       Background(lipgloss.Color("#1e1e2e"))
   ```

3. Replace with theme helper:
   ```go
   var style = theme.ModalStyle()
   ```

4. For custom styles, use color constants:
   ```go
   var style = lipgloss.NewStyle().
       Foreground(lipgloss.Color(theme.ForegroundPrimary)).
       Background(lipgloss.Color(theme.BackgroundPrimary))
   ```

## Benefits

- **Consistency** - All components use the same color palette
- **Maintainability** - Update colors in one place
- **Type Safety** - Constants prevent typos
- **Code Reduction** - Helper functions eliminate duplication
- **Easy Theming** - Simple to switch color schemes
- **Better Organization** - Clear separation of concerns

## Future Enhancements

Potential improvements to the theme system:

- **Multiple Themes** - Support for light/dark mode
- **Theme Configuration** - Load theme from config file
- **Dynamic Theming** - Runtime theme switching
- **Color Validation** - Ensure colors are valid before use
- **Accessibility** - High-contrast mode options
- **Custom Themes** - Allow users to define custom color schemes

## Contributing

When adding new UI components:

1. Check if existing theme helpers can be used
2. If not, add new helpers to `theme.go`
3. Follow the naming convention: `ComponentStyle()` or `ComponentStateStyle()`
4. Document the new helpers in this README
5. Update the implementation plan if needed

## See Also

- [Implementation Plan](../../docs/plans/theme-system/implementation-plan.md) - Detailed design document
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss) - Styling library
- [Catppuccin Theme](https://catppuccin.com/) - Color palette inspiration