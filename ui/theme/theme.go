// Package theme provides a centralized theme system for PromptStack TUI.
// It uses the Catppuccin Mocha color palette and provides style helper functions
// for consistent UI styling across all components.
//
// Style categories:
//   - Modal styles: ModalStyle, ModalTitleStyle, ModalContentStyle, ModalButtonStyle, ModalButtonFocusedStyle
//   - Status styles: StatusStyle, InfoStyle, SuccessStyle, WarningStyle, ErrorStyle
//   - Input styles: InputStyle, SearchInputStyle
//   - List styles: ListItemStyle, ListItemSelectedStyle, ListCategoryStyle
//   - Preview styles: PreviewStyle, PreviewTitleStyle
//   - Utility styles: HeaderStyle, CursorStyle, ActivePlaceholderStyle, ValidationErrorStyle, HighlightStyle, MutedStyle
//   - Diff styles: DiffAddedStyle, DiffRemovedStyle, DiffContextStyle, DiffHeaderStyle
//   - Message styles: UserMessageStyle, AssistantMessageStyle, SystemMessageStyle
//   - Tool styles: ToolNameStyle, ToolOutputStyle
package theme

import "github.com/charmbracelet/lipgloss"

// Spacing system constants
const (
	Unit          = 1      // 1 character width/height (monospace grid)
	DoubleNewline = "\n\n" // Two newlines for vertical spacing
)

// Keyboard shortcut constants
const (
	LeaderKey         = "Ctrl+X"
	CommandPaletteKey = "Ctrl+P"
	SwitchAgentsKey   = "Tab"
	CancelKey         = "Esc"
	InterruptKey      = "Ctrl+C"
	SuspendKey        = "Ctrl+Z"
	AbortKey          = "Ctrl+G"
)

// Color constants using Catppuccin Mocha palette
const (
	// Background colors
	BackgroundPrimary   = "#1e1e2e" // Main modal background
	BackgroundSecondary = "#181825" // Status bar background
	BackgroundTertiary  = "#313244" // Secondary button background
	BackgroundInput     = "#313244" // Input field background
	BackgroundHighlight = "#45475a" // Highlight/background selection
	BackgroundMuted     = "#1e1e2e" // Muted background

	// Foreground colors
	ForegroundPrimary   = "#cdd6f4" // Main text
	ForegroundSecondary = "#a6adc8" // Secondary text
	ForegroundMuted     = "#6c7086" // Muted text
	ForegroundWhite     = "#ffffff" // White text
	ForegroundBlack     = "#11111b" // Black text

	// Accent colors
	AccentBlue   = "#89b4fa" // Primary accent, info states
	AccentGreen  = "#a6e3a1" // Success states
	AccentYellow = "#f9e2af" // Warning states
	AccentRed    = "#f38ba8" // Error states
	AccentCyan   = "#94e2d5" // Category labels
	AccentMauve  = "#cba6f7" // User messages
	AccentPink   = "#f5c2e7" // Special highlights
	AccentOrange = "#fab387" // Tool names

	// Diff colors
	DiffAddedBackground   = "#a6e3a1" // Added line background
	DiffAddedForeground   = "#1e1e2e" // Added line text
	DiffRemovedBackground = "#f38ba8" // Removed line background
	DiffRemovedForeground = "#1e1e2e" // Removed line text
	DiffContextBackground = "#313244" // Context line background
	DiffContextForeground = "#a6adc8" // Context line text
	DiffHeaderBackground  = "#45475a" // Diff header background
	DiffHeaderForeground  = "#89b4fa" // Diff header text

	// Border colors
	BorderPrimary = "#45475a" // Modal borders
	BorderMuted   = "#585b70" // Secondary borders
	BorderFocus   = "#89b4fa" // Focused element border

	// Special colors
	CursorBackground = "#cdd6f4" // Cursor highlight (light gray)
	CursorForeground = "#11111b" // Cursor text (dark)
	TextMuted        = "#6c7086" // Dimmed text
)

// Modal styles
func ModalStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(BackgroundPrimary)).
		Foreground(lipgloss.Color(ForegroundPrimary)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(BorderPrimary))
}

func ModalTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundPrimary)).
		Bold(true)
}

func ModalContentStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundPrimary))
}

func ModalButtonStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundPrimary)).
		Background(lipgloss.Color(BackgroundTertiary)).
		Padding(0, 2)
}

func ModalButtonFocusedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(BackgroundPrimary)).
		Background(lipgloss.Color(AccentBlue)).
		Padding(0, 2)
}

// Status styles
func StatusStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(BackgroundSecondary)).
		Foreground(lipgloss.Color(ForegroundPrimary))
}

func InfoStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentBlue))
}

func SuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentGreen))
}

func WarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentYellow))
}

func ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentRed))
}

// Input styles
func InputStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(BackgroundInput)).
		Foreground(lipgloss.Color(ForegroundPrimary)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(BorderMuted))
}

func SearchInputStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(BackgroundInput)).
		Foreground(lipgloss.Color(ForegroundPrimary)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(AccentBlue))
}

// List styles
func ListItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundPrimary))
}

func ListItemSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(BackgroundPrimary)).
		Background(lipgloss.Color(AccentBlue))
}

func ListCategoryStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentCyan)).
		Bold(true)
}

// Preview styles
func PreviewStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(BackgroundPrimary)).
		Foreground(lipgloss.Color(ForegroundPrimary)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(BorderMuted))
}

func PreviewTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentMauve)).
		Bold(true)
}

// Utility styles
func HeaderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundSecondary)).
		Bold(true)
}

func CursorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(CursorBackground)).
		Foreground(lipgloss.Color(CursorForeground))
}

func ActivePlaceholderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(CursorBackground)).
		Foreground(lipgloss.Color(CursorForeground))
}

func ValidationErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentRed)).
		Background(lipgloss.Color(BackgroundPrimary)).
		Bold(true)
}

// Diff styles
func DiffAddedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(DiffAddedBackground)).
		Foreground(lipgloss.Color(DiffAddedForeground))
}

func DiffRemovedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(DiffRemovedBackground)).
		Foreground(lipgloss.Color(DiffRemovedForeground))
}

func DiffContextStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(DiffContextBackground)).
		Foreground(lipgloss.Color(DiffContextForeground))
}

func DiffHeaderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(DiffHeaderBackground)).
		Foreground(lipgloss.Color(DiffHeaderForeground)).
		Bold(true)
}

// Message styles (for chat interface)
func UserMessageStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentBlue)).
		Bold(true)
}

func AssistantMessageStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentCyan)).
		Bold(true)
}

func SystemMessageStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(TextMuted)).
		Italic(true)
}

// Tool styles
func ToolNameStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(AccentOrange)).
		Bold(true)
}

func ToolOutputStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ForegroundPrimary))
}

// Highlight styles
func HighlightStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(AccentYellow)).
		Foreground(lipgloss.Color(BackgroundPrimary))
}

func MutedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(TextMuted))
}

// KeyboardHelp returns a formatted help text string for keyboard shortcuts.
// This centralizes keyboard shortcut documentation for consistency across the UI.
func KeyboardHelp() string {
	return LeaderKey + " (command), " + CommandPaletteKey + " (palette), " + CancelKey + " (cancel)"
}
