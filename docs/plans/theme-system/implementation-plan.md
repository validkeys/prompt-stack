# Theme System Implementation Plan

## Overview

Create a centralized theme system to manage all lipgloss styling across the UI components, replacing scattered hard-coded hex codes and terminal color numbers with a single source of truth.

## Current State Analysis

### Files Using lipgloss
- [`ui/common/modal.go`](../../ui/common/modal.go) - Modal dialogs with buttons
- [`ui/statusbar/model.go`](../../ui/statusbar/model.go) - Status bar with message types
- [`ui/workspace/model.go`](../../ui/workspace/model.go) - Workspace editor
- [`ui/palette/model.go`](../../ui/palette/model.go) - Command palette
- [`ui/browser/model.go`](../../ui/browser/model.go) - Library browser

### Color Patterns Identified

**Catppuccin Mocha Theme Colors (Hex Codes):**
- Backgrounds: `#1e1e2e`, `#181825`, `#313244`
- Foregrounds: `#cdd6f4`, `#a6adc8`, `#6c7086`
- Accents: `#89b4fa` (blue), `#a6e3a1` (green), `#f9e2af` (yellow), `#f38ba8` (red)
- Borders: `#45475a`

**Terminal Color Numbers:**
- White: `15`
- Gray: `240`, `245`
- Dark Gray: `236`
- Blue: `39`
- Light Gray: `7`
- Black: `0`

### Usage Categories
1. **Modal Components**: Borders, backgrounds, titles, buttons (normal/focused)
2. **Status Bar**: Normal, info, success, warning, error states
3. **Workspace**: Cursor highlighting, status bar
4. **Command Palette**: Modal, headers, search input, list items, categories
5. **Library Browser**: Modal, headers, search input, list items, preview pane

## Proposed Solution

### File Structure
```
ui/
└── theme/
    └── theme.go          # Centralized theme definitions
```

### Theme Design

#### 1. Color Constants
Organize colors by semantic purpose:

```go
// Background Colors
const (
    BackgroundPrimary   = "#1e1e2e"  // Main modal background
    BackgroundSecondary = "#181825"  // Status bar background
    BackgroundTertiary  = "#313244"  // Secondary button background
    BackgroundInput     = "#236"     // Input field background (terminal color)
)

// Foreground Colors
const (
    ForegroundPrimary   = "#cdd6f4"  // Main text
    ForegroundSecondary = "#a6adc8"  // Secondary text
    ForegroundMuted     = "#6c7086"  // Muted text
    ForegroundWhite     = "15"       // White text (terminal color)
)

// Accent Colors
const (
    AccentBlue   = "#89b4fa"  // Primary accent
    AccentGreen  = "#a6e3a1"  // Success
    AccentYellow = "#f9e2af"  // Warning
    AccentRed    = "#f38ba8"  // Error
    AccentCyan   = "39"       // Category labels (terminal color)
)

// Border Colors
const (
    BorderPrimary = "#45475a"  // Modal borders
    BorderMuted   = "240"      // Secondary borders (terminal color)
)

// Special Colors
const (
    CursorBackground = "7"   // Cursor highlight
    CursorForeground  = "0"   // Cursor text
    TextMuted         = "245" // Dimmed text (terminal color)
)
```

#### 2. Helper Functions
Create reusable style builders:

```go
// Modal Styles
func ModalStyle() lipgloss.Style
func ModalTitleStyle() lipgloss.Style
func ModalContentStyle() lipgloss.Style
func ModalButtonStyle() lipgloss.Style
func ModalButtonFocusedStyle() lipgloss.Style
func ModalButtonSecondaryStyle() lipgloss.Style
func ModalButtonSecondaryFocusedStyle() lipgloss.Style

// Status Bar Styles
func StatusStyle() lipgloss.Style
func InfoStyle() lipgloss.Style
func SuccessStyle() lipgloss.Style
func WarningStyle() lipgloss.Style
func ErrorStyle() lipgloss.Style
func SeparatorStyle() lipgloss.Style

// Input Styles
func InputStyle() lipgloss.Style
func SearchInputStyle() lipgloss.Style

// List Styles
func ListItemStyle() lipgloss.Style
func ListItemSelectedStyle() lipgloss.Style
func ListCategoryStyle() lipgloss.Style
func ListDescriptionStyle() lipgloss.Style
func ListEmptyStyle() lipgloss.Style

// Preview Styles
func PreviewStyle() lipgloss.Style
func PreviewTitleStyle() lipgloss.Style
func PreviewDescriptionStyle() lipgloss.Style
func PreviewTagsStyle() lipgloss.Style
func PreviewContentStyle() lipgloss.Style

// Cursor Styles
func CursorStyle() lipgloss.Style

// Header Styles
func HeaderStyle() lipgloss.Style
```

#### 3. Style Composition
Allow flexible style composition:

```go
// Base styles that can be extended
func BaseModal() lipgloss.Style
func BaseInput() lipgloss.Style
func BaseListItem() lipgloss.Style
```

### Implementation Steps

1. ✅ **Create theme.go file** with all color constants and helper functions
2. ✅ **Update ui/common/modal.go** to use theme functions
3. ✅ **Update ui/statusbar/model.go** to use theme functions
4. ✅ **Update ui/workspace/model.go** to use theme functions
5. ✅ **Update ui/palette/model.go** to use theme functions
6. ✅ **Update ui/browser/model.go** to use theme functions
7. ✅ **Test all UI components** to ensure visual consistency
8. ✅ **Document usage** in README and key learnings

### Migration Strategy

For each file:
1. Import the theme package
2. Replace hard-coded `lipgloss.Color()` calls with theme constants
3. Replace style variable declarations with theme helper function calls
4. Test the component to verify visual appearance

Example migration:

**Before:**
```go
var (
    modalStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#45475a")).
        Padding(1, 2).
        Background(lipgloss.Color("#1e1e2e"))
)
```

**After:**
```go
import "github.com/kyledavis/prompt-stack/ui/theme"

var (
    modalStyle = theme.ModalStyle()
)
```

### Benefits

1. **Single Source of Truth**: All colors defined in one place
2. **Easy Theme Switching**: Change colors in theme.go to update entire UI
3. **Consistency**: Ensures all components use the same color palette
4. **Maintainability**: Easier to update and refactor styling
5. **Type Safety**: Constants prevent typos in color values
6. **Reusability**: Helper functions reduce code duplication

### Future Enhancements

1. **Multiple Themes**: Support light/dark mode or different color schemes
2. **Theme Configuration**: Load theme from config file
3. **Dynamic Theming**: Allow runtime theme switching
4. **Color Validation**: Ensure colors are valid before use
5. **Accessibility**: Add high-contrast mode options

## Testing Checklist

- [x] Modal dialogs display correctly with all button states
- [x] Status bar shows all message types (normal, info, success, warning, error)
- [x] Workspace cursor highlighting works properly
- [x] Command palette displays correctly with search and selection
- [x] Library browser shows list items and preview pane correctly
- [x] All borders and backgrounds render as expected
- [x] Text colors are readable and consistent
- [x] No visual regressions compared to current implementation
- [x] All UI components compile successfully (`go build ./ui/...`)

## Documentation

✅ **Created** `ui/theme/README.md` covering:
- Overview of the theme system
- Available color constants
- Helper function reference
- Usage examples
- Migration guide for existing code
- How to create custom themes

✅ **Updated** `docs/plans/initial-build/key-learnings.md` with:
- Centralized Theme System section documenting the pattern
- Implementation examples
- Benefits and lessons learned