package theme

import "github.com/charmbracelet/lipgloss"

// Color Constants - Catppuccin Mocha Theme
// These constants provide a single source of truth for all UI colors

// Background Colors
const (
	BackgroundPrimary   = "#1e1e2e" // Main modal background
	BackgroundSecondary = "#181825" // Status bar background
	BackgroundTertiary  = "#313244" // Secondary button background
	BackgroundInput     = "236"     // Input field background (terminal color)
)

// Foreground Colors
const (
	ForegroundPrimary   = "#cdd6f4" // Main text
	ForegroundSecondary = "#a6adc8" // Secondary text
	ForegroundMuted     = "#6c7086" // Muted text
	ForegroundWhite     = "15"      // White text (terminal color)
)

// Accent Colors
const (
	AccentBlue   = "#89b4fa" // Primary accent
	AccentGreen  = "#a6e3a1" // Success
	AccentYellow = "#f9e2af" // Warning
	AccentRed    = "#f38ba8" // Error
	AccentCyan   = "39"      // Category labels (terminal color)
)

// Border Colors
const (
	BorderPrimary = "#45475a" // Modal borders
	BorderMuted   = "240"     // Secondary borders (terminal color)
)

// Special Colors
const (
	CursorBackground = "7"   // Cursor highlight
	CursorForeground = "0"   // Cursor text
	TextMuted        = "245" // Dimmed text (terminal color)
)

// Modal Styles

// ModalStyle returns the base style for modal containers
func ModalStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(BorderPrimary)).
		Padding(1, 2).
		Background(lipgloss.Color(BackgroundPrimary))
}

// ModalTitleStyle returns the style for modal titles
func ModalTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundPrimary)).
		Bold(true).
		MarginBottom(1)
}

// ModalContentStyle returns the style for modal content
func ModalContentStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundSecondary)).
		Width(60)
}

// ModalButtonStyle returns the style for primary buttons
func ModalButtonStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundPrimary)).
		Background(lipgloss.Color(BorderPrimary)).
		Padding(0, 2).
		MarginRight(1)
}

// ModalButtonFocusedStyle returns the style for focused primary buttons
func ModalButtonFocusedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(BackgroundPrimary)).
		Background(lipgloss.Color(AccentBlue)).
		Padding(0, 2).
		MarginRight(1)
}

// ModalButtonSecondaryStyle returns the style for secondary buttons
func ModalButtonSecondaryStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundPrimary)).
		Background(lipgloss.Color(BackgroundTertiary)).
		Padding(0, 2).
		MarginRight(1)
}

// ModalButtonSecondaryFocusedStyle returns the style for focused secondary buttons
func ModalButtonSecondaryFocusedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(BackgroundPrimary)).
		Background(lipgloss.Color(AccentBlue)).
		Padding(0, 2).
		MarginRight(1)
}

// Status Bar Styles

// StatusStyle returns the style for normal status messages
func StatusStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundMuted)).
		Background(lipgloss.Color(BackgroundSecondary)).
		Padding(0, 1)
}

// InfoStyle returns the style for info messages
func InfoStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentBlue)).
		Background(lipgloss.Color(BackgroundSecondary)).
		Padding(0, 1)
}

// SuccessStyle returns the style for success messages
func SuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentGreen)).
		Background(lipgloss.Color(BackgroundSecondary)).
		Padding(0, 1)
}

// WarningStyle returns the style for warning messages
func WarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentYellow)).
		Background(lipgloss.Color(BackgroundSecondary)).
		Padding(0, 1)
}

// ErrorStyle returns the style for error messages
func ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentRed)).
		Background(lipgloss.Color(BackgroundSecondary)).
		Padding(0, 1)
}

// SeparatorStyle returns the style for status bar separators
func SeparatorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(BorderPrimary)).
		Padding(0, 1)
}

// Input Styles

// InputStyle returns the base style for input fields
func InputStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundWhite)).
		Background(lipgloss.Color(BackgroundInput)).
		Padding(0, 1)
}

// SearchInputStyle returns the style for search input fields
func SearchInputStyle() lipgloss.Style {
	return InputStyle()
}

// List Styles

// ListItemStyle returns the style for unselected list items
func ListItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundWhite))
}

// ListItemSelectedStyle returns the style for selected list items
func ListItemSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(BorderMuted)).
		Foreground(lipgloss.Color(ForegroundWhite))
}

// ListCategoryStyle returns the style for category labels in lists
func ListCategoryStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentCyan)).
		Bold(true)
}

// ListDescriptionStyle returns the style for descriptions in lists
func ListDescriptionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(TextMuted)).
		Italic(true)
}

// ListEmptyStyle returns the style for empty list states
func ListEmptyStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(TextMuted)).
		Italic(true)
}

// Preview Styles

// PreviewStyle returns the base style for preview panes
func PreviewStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(BorderMuted))
}

// PreviewTitleStyle returns the style for preview titles
func PreviewTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundWhite)).
		Bold(true)
}

// PreviewDescriptionStyle returns the style for preview descriptions
func PreviewDescriptionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(TextMuted)).
		Italic(true)
}

// PreviewTagsStyle returns the style for preview tags
func PreviewTagsStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentCyan))
}

// PreviewContentStyle returns the style for preview content
func PreviewContentStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundWhite))
}

// Cursor Styles

// CursorStyle returns the style for cursor highlighting
func CursorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(CursorBackground)).
		Foreground(lipgloss.Color(CursorForeground))
}

// ActivePlaceholderStyle returns the style for active placeholder highlighting
func ActivePlaceholderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(AccentYellow)).
		Foreground(lipgloss.Color(BackgroundPrimary)).
		Bold(true)
}

// Header Styles

// HeaderStyle returns the style for section headers
func HeaderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundWhite)).
		Bold(true).
		MarginBottom(1)
}

// Base Styles for Composition

// BaseModal returns a base modal style that can be extended
func BaseModal() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(BorderPrimary)).
		Padding(0, 1)
}

// BaseInput returns a base input style that can be extended
func BaseInput() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundWhite)).
		Background(lipgloss.Color(BackgroundInput)).
		Padding(0, 1)
}

// BaseListItem returns a base list item style that can be extended
func BaseListItem() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundWhite))
}
