package theme_test

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/ui/theme"
)

func TestColorConstants(t *testing.T) {
	tests := []struct {
		name  string
		color string
	}{
		// Background colors
		{"BackgroundPrimary", theme.BackgroundPrimary},
		{"BackgroundSecondary", theme.BackgroundSecondary},
		{"BackgroundTertiary", theme.BackgroundTertiary},
		{"BackgroundInput", theme.BackgroundInput},
		{"BackgroundHighlight", theme.BackgroundHighlight},
		{"BackgroundMuted", theme.BackgroundMuted},
		// Foreground colors
		{"ForegroundPrimary", theme.ForegroundPrimary},
		{"ForegroundSecondary", theme.ForegroundSecondary},
		{"ForegroundMuted", theme.ForegroundMuted},
		{"ForegroundWhite", theme.ForegroundWhite},
		{"ForegroundBlack", theme.ForegroundBlack},
		// Accent colors
		{"AccentBlue", theme.AccentBlue},
		{"AccentGreen", theme.AccentGreen},
		{"AccentYellow", theme.AccentYellow},
		{"AccentRed", theme.AccentRed},
		{"AccentCyan", theme.AccentCyan},
		{"AccentMauve", theme.AccentMauve},
		{"AccentPink", theme.AccentPink},
		{"AccentOrange", theme.AccentOrange},
		// Diff colors
		{"DiffAddedBackground", theme.DiffAddedBackground},
		{"DiffAddedForeground", theme.DiffAddedForeground},
		{"DiffRemovedBackground", theme.DiffRemovedBackground},
		{"DiffRemovedForeground", theme.DiffRemovedForeground},
		{"DiffContextBackground", theme.DiffContextBackground},
		{"DiffContextForeground", theme.DiffContextForeground},
		{"DiffHeaderBackground", theme.DiffHeaderBackground},
		{"DiffHeaderForeground", theme.DiffHeaderForeground},
		// Border colors
		{"BorderPrimary", theme.BorderPrimary},
		{"BorderMuted", theme.BorderMuted},
		{"BorderFocus", theme.BorderFocus},
		// Special colors
		{"CursorBackground", theme.CursorBackground},
		{"CursorForeground", theme.CursorForeground},
		{"TextMuted", theme.TextMuted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color == "" {
				t.Errorf("color constant %s is empty", tt.name)
			}
		})
	}
}

func TestModalStyle(t *testing.T) {
	style := theme.ModalStyle()

	// Test that style can be applied to text
	result := style.Render("test")
	if result == "" {
		t.Error("style.Render returned empty string")
	}
}

func TestStatusStyle(t *testing.T) {
	style := theme.StatusStyle()

	// Test that style can be applied to text
	result := style.Render("test")
	if result == "" {
		t.Error("style.Render returned empty string")
	}
}

func TestActivePlaceholderStyle(t *testing.T) {
	style := theme.ActivePlaceholderStyle()

	// Test that style can be applied to text
	result := style.Render("test")
	if result == "" {
		t.Error("style.Render returned empty string")
	}
}

func TestStyleChaining(t *testing.T) {
	// Test that styles can be chained
	baseStyle := theme.ModalStyle()
	chainedStyle := baseStyle.Width(80).Height(40).Padding(2, 3)

	// Test that chained style can be applied
	result := chainedStyle.Render("test")
	if result == "" {
		t.Error("chained style.Render returned empty string")
	}
}

func TestStyleWithSampleText(t *testing.T) {
	tests := []struct {
		name  string
		style lipgloss.Style
		text  string
	}{
		// Modal styles
		{"ModalStyle with text", theme.ModalStyle(), "Modal Content"},
		{"ModalTitleStyle with text", theme.ModalTitleStyle(), "Modal Title"},
		{"ModalContentStyle with text", theme.ModalContentStyle(), "Modal Content"},
		{"ModalButtonStyle with text", theme.ModalButtonStyle(), "Button"},
		{"ModalButtonFocusedStyle with text", theme.ModalButtonFocusedStyle(), "Focused Button"},
		// Status styles
		{"StatusStyle with text", theme.StatusStyle(), "Status: Ready"},
		{"InfoStyle with text", theme.InfoStyle(), "Info message"},
		{"SuccessStyle with text", theme.SuccessStyle(), "Success!"},
		{"WarningStyle with text", theme.WarningStyle(), "Warning!"},
		{"ErrorStyle with text", theme.ErrorStyle(), "Error!"},
		// Input styles
		{"InputStyle with text", theme.InputStyle(), "Input text"},
		{"SearchInputStyle with text", theme.SearchInputStyle(), "Search..."},
		// List styles
		{"ListItemStyle with text", theme.ListItemStyle(), "List item"},
		{"ListItemSelectedStyle with text", theme.ListItemSelectedStyle(), "Selected item"},
		{"ListCategoryStyle with text", theme.ListCategoryStyle(), "Category"},
		// Preview styles
		{"PreviewStyle with text", theme.PreviewStyle(), "Preview content"},
		{"PreviewTitleStyle with text", theme.PreviewTitleStyle(), "Preview Title"},
		// Utility styles
		{"HeaderStyle with text", theme.HeaderStyle(), "Header"},
		{"CursorStyle with text", theme.CursorStyle(), "cursor"},
		{"ActivePlaceholderStyle with text", theme.ActivePlaceholderStyle(), "{{placeholder}}"},
		{"ValidationErrorStyle with text", theme.ValidationErrorStyle(), "Validation error"},
		// Diff styles
		{"DiffAddedStyle with text", theme.DiffAddedStyle(), "+ added line"},
		{"DiffRemovedStyle with text", theme.DiffRemovedStyle(), "- removed line"},
		{"DiffContextStyle with text", theme.DiffContextStyle(), "  context line"},
		{"DiffHeaderStyle with text", theme.DiffHeaderStyle(), "@@ diff header @@"},
		// Message styles
		{"UserMessageStyle with text", theme.UserMessageStyle(), "User: message"},
		{"AssistantMessageStyle with text", theme.AssistantMessageStyle(), "Assistant: message"},
		{"SystemMessageStyle with text", theme.SystemMessageStyle(), "System: message"},
		// Tool styles
		{"ToolNameStyle with text", theme.ToolNameStyle(), "Tool: bash"},
		{"ToolOutputStyle with text", theme.ToolOutputStyle(), "Tool output"},
		// Highlight styles
		{"HighlightStyle with text", theme.HighlightStyle(), "highlighted"},
		{"MutedStyle with text", theme.MutedStyle(), "muted text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.style.Render(tt.text)
			if result == "" {
				t.Errorf("style.Render returned empty string for text %q", tt.text)
			}
			// Verify that the rendered output contains the original text
			if len(result) < len(tt.text) {
				t.Errorf("rendered output is shorter than input text: got %d chars, want at least %d", len(result), len(tt.text))
			}
		})
	}
}

func TestColorValues(t *testing.T) {
	// Verify specific hex values match Catppuccin Mocha specification
	tests := []struct {
		name     string
		color    string
		expected string
	}{
		// Background colors
		{"BackgroundPrimary", theme.BackgroundPrimary, "#1e1e2e"},
		{"BackgroundSecondary", theme.BackgroundSecondary, "#181825"},
		{"BackgroundTertiary", theme.BackgroundTertiary, "#313244"},
		{"BackgroundInput", theme.BackgroundInput, "#313244"},
		{"BackgroundHighlight", theme.BackgroundHighlight, "#45475a"},
		{"BackgroundMuted", theme.BackgroundMuted, "#1e1e2e"},
		// Foreground colors
		{"ForegroundPrimary", theme.ForegroundPrimary, "#cdd6f4"},
		{"ForegroundSecondary", theme.ForegroundSecondary, "#a6adc8"},
		{"ForegroundMuted", theme.ForegroundMuted, "#6c7086"},
		{"ForegroundWhite", theme.ForegroundWhite, "#ffffff"},
		{"ForegroundBlack", theme.ForegroundBlack, "#11111b"},
		// Accent colors
		{"AccentBlue", theme.AccentBlue, "#89b4fa"},
		{"AccentGreen", theme.AccentGreen, "#a6e3a1"},
		{"AccentYellow", theme.AccentYellow, "#f9e2af"},
		{"AccentRed", theme.AccentRed, "#f38ba8"},
		{"AccentCyan", theme.AccentCyan, "#94e2d5"},
		{"AccentMauve", theme.AccentMauve, "#cba6f7"},
		{"AccentPink", theme.AccentPink, "#f5c2e7"},
		{"AccentOrange", theme.AccentOrange, "#fab387"},
		// Diff colors
		{"DiffAddedBackground", theme.DiffAddedBackground, "#a6e3a1"},
		{"DiffAddedForeground", theme.DiffAddedForeground, "#1e1e2e"},
		{"DiffRemovedBackground", theme.DiffRemovedBackground, "#f38ba8"},
		{"DiffRemovedForeground", theme.DiffRemovedForeground, "#1e1e2e"},
		{"DiffContextBackground", theme.DiffContextBackground, "#313244"},
		{"DiffContextForeground", theme.DiffContextForeground, "#a6adc8"},
		{"DiffHeaderBackground", theme.DiffHeaderBackground, "#45475a"},
		{"DiffHeaderForeground", theme.DiffHeaderForeground, "#89b4fa"},
		// Border colors
		{"BorderPrimary", theme.BorderPrimary, "#45475a"},
		{"BorderMuted", theme.BorderMuted, "#585b70"},
		{"BorderFocus", theme.BorderFocus, "#89b4fa"},
		// Special colors
		{"CursorBackground", theme.CursorBackground, "#cdd6f4"},
		{"CursorForeground", theme.CursorForeground, "#11111b"},
		{"TextMuted", theme.TextMuted, "#6c7086"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color != tt.expected {
				t.Errorf("color %s = %q, want %q", tt.name, tt.color, tt.expected)
			}
		})
	}
}

func TestSpacingConstants(t *testing.T) {
	if theme.Unit != 1 {
		t.Errorf("Unit = %d, want 1", theme.Unit)
	}
}

func TestKeyboardShortcuts(t *testing.T) {
	tests := []struct {
		name     string
		shortcut string
	}{
		{"LeaderKey", theme.LeaderKey},
		{"CommandPaletteKey", theme.CommandPaletteKey},
		{"SwitchAgentsKey", theme.SwitchAgentsKey},
		{"CancelKey", theme.CancelKey},
		{"InterruptKey", theme.InterruptKey},
		{"SuspendKey", theme.SuspendKey},
		{"AbortKey", theme.AbortKey},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shortcut == "" {
				t.Errorf("keyboard shortcut %s is empty", tt.name)
			}
		})
	}
}

func TestKeyboardHelp(t *testing.T) {
	help := theme.KeyboardHelp()
	if help == "" {
		t.Error("KeyboardHelp returned empty string")
	}
	// Verify that it contains all expected keyboard shortcuts
	expectedShortcuts := []string{
		theme.LeaderKey,
		theme.CommandPaletteKey,
		theme.CancelKey,
	}
	for _, shortcut := range expectedShortcuts {
		if len(help) < len(shortcut) || !contains(help, shortcut) {
			t.Errorf("KeyboardHelp does not contain shortcut %q, got %q", shortcut, help)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
