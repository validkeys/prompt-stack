// Package theme provides a centralized theme system for the PromptStack TUI.
// It uses the Catppuccin Mocha color palette and provides style helper functions
// for consistent UI styling across all components.
package theme

import "github.com/charmbracelet/lipgloss"

// Color constants using Catppuccin Mocha palette
const (
	// Background colors
	BackgroundPrimary   = "#1e1e2e" // Main modal background
	BackgroundSecondary = "#181825" // Status bar background
	BackgroundTertiary  = "#313244" // Secondary button background
	BackgroundInput     = "236"     // Input field background (terminal color)

	// Foreground colors
	ForegroundPrimary   = "#cdd6f4" // Main text
	ForegroundSecondary = "#a6adc8" // Secondary text
	ForegroundMuted     = "#6c7086" // Muted text
	ForegroundWhite     = "15"      // White text (terminal color)

	// Accent colors
	AccentBlue   = "#89b4fa" // Primary accent, info states
	AccentGreen  = "#a6e3a1" // Success states
	AccentYellow = "#f9e2af" // Warning states
	AccentRed    = "#f38ba8" // Error states
	AccentCyan   = "39"      // Category labels (terminal color)

	// Border colors
	BorderPrimary = "#45475a" // Modal borders
	BorderMuted   = "240"     // Secondary borders (terminal color)

	// Special colors
	CursorBackground = "7"   // Cursor highlight
	CursorForeground = "0"   // Cursor text
	TextMuted        = "245" // Dimmed text (terminal color)
)

// ModalStyle returns the base style for modal containers.
func ModalStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(BackgroundPrimary)).
		Foreground(lipgloss.Color(ForegroundPrimary)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(BorderPrimary))
}

// StatusStyle returns the style for status bar components.
func StatusStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(BackgroundSecondary)).
		Foreground(lipgloss.Color(ForegroundPrimary))
}

// ActivePlaceholderStyle returns the style for active placeholder highlighting.
func ActivePlaceholderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(CursorBackground)).
		Foreground(lipgloss.Color(CursorForeground))
}
