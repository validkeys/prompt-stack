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
		{"BackgroundPrimary", theme.BackgroundPrimary},
		{"BackgroundSecondary", theme.BackgroundSecondary},
		{"BackgroundTertiary", theme.BackgroundTertiary},
		{"BackgroundInput", theme.BackgroundInput},
		{"ForegroundPrimary", theme.ForegroundPrimary},
		{"ForegroundSecondary", theme.ForegroundSecondary},
		{"ForegroundMuted", theme.ForegroundMuted},
		{"ForegroundWhite", theme.ForegroundWhite},
		{"AccentBlue", theme.AccentBlue},
		{"AccentGreen", theme.AccentGreen},
		{"AccentYellow", theme.AccentYellow},
		{"AccentRed", theme.AccentRed},
		{"AccentCyan", theme.AccentCyan},
		{"BorderPrimary", theme.BorderPrimary},
		{"BorderMuted", theme.BorderMuted},
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
		{"ModalStyle with text", theme.ModalStyle(), "Modal Content"},
		{"StatusStyle with text", theme.StatusStyle(), "Status: Ready"},
		{"ActivePlaceholderStyle with text", theme.ActivePlaceholderStyle(), "{{placeholder}}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.style.Render(tt.text)
			if result == "" {
				t.Errorf("style.Render returned empty string for text %q", tt.text)
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
		{"BackgroundPrimary", theme.BackgroundPrimary, "#1e1e2e"},
		{"BackgroundSecondary", theme.BackgroundSecondary, "#181825"},
		{"BackgroundInput", theme.BackgroundInput, "236"},
		{"ForegroundPrimary", theme.ForegroundPrimary, "#cdd6f4"},
		{"ForegroundMuted", theme.ForegroundMuted, "#6c7086"},
		{"ForegroundWhite", theme.ForegroundWhite, "15"},
		{"AccentBlue", theme.AccentBlue, "#89b4fa"},
		{"AccentGreen", theme.AccentGreen, "#a6e3a1"},
		{"AccentYellow", theme.AccentYellow, "#f9e2af"},
		{"AccentRed", theme.AccentRed, "#f38ba8"},
		{"AccentCyan", theme.AccentCyan, "39"},
		{"BorderPrimary", theme.BorderPrimary, "#45475a"},
		{"BorderMuted", theme.BorderMuted, "240"},
		{"CursorBackground", theme.CursorBackground, "7"},
		{"CursorForeground", theme.CursorForeground, "0"},
		{"TextMuted", theme.TextMuted, "245"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color != tt.expected {
				t.Errorf("color %s = %q, want %q", tt.name, tt.color, tt.expected)
			}
		})
	}
}
