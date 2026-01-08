package statusbar

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kyledavis/prompt-stack/ui/theme"
)

// TestNew verifies that New creates a status bar model with default values.
func TestNew(t *testing.T) {
	model := New()

	if model.charCount != 0 {
		t.Errorf("Expected charCount to be 0, got %d", model.charCount)
	}
	if model.lineCount != 0 {
		t.Errorf("Expected lineCount to be 0, got %d", model.lineCount)
	}
	if model.width != 80 {
		t.Errorf("Expected width to be 80, got %d", model.width)
	}
}

// TestInit verifies that Init returns nil command.
func TestInit(t *testing.T) {
	model := New()
	cmd := model.Init()

	if cmd != nil {
		t.Errorf("Expected Init to return nil command, got %v", cmd)
	}
}

// TestUpdateWindowSizeMsg verifies that Update handles WindowSizeMsg correctly.
func TestUpdateWindowSizeMsg(t *testing.T) {
	tests := []struct {
		name     string
		initial  Model
		msg      tea.WindowSizeMsg
		expected int
	}{
		{
			name:     "standard width",
			initial:  New(),
			msg:      tea.WindowSizeMsg{Width: 100, Height: 24},
			expected: 100,
		},
		{
			name:     "zero width",
			initial:  New(),
			msg:      tea.WindowSizeMsg{Width: 0, Height: 24},
			expected: 0,
		},
		{
			name:     "large width",
			initial:  New(),
			msg:      tea.WindowSizeMsg{Width: 1000, Height: 24},
			expected: 1000,
		},
		{
			name:     "small width",
			initial:  New(),
			msg:      tea.WindowSizeMsg{Width: 20, Height: 24},
			expected: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newModel, cmd := tt.initial.Update(tt.msg)

			if cmd != nil {
				t.Errorf("Expected Update to return nil command, got %v", cmd)
			}

			statusBar, ok := newModel.(Model)
			if !ok {
				t.Fatalf("Expected Model type, got %T", newModel)
			}

			if statusBar.width != tt.expected {
				t.Errorf("Expected width to be %d, got %d", tt.expected, statusBar.width)
			}
		})
	}
}

// TestUpdateOtherMessages verifies that Update ignores non-WindowSizeMsg messages.
func TestUpdateOtherMessages(t *testing.T) {
	model := New()
	initialWidth := model.width

	// Test with KeyMsg
	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Errorf("Expected Update to return nil command for KeyMsg, got %v", cmd)
	}
	statusBar, ok := newModel.(Model)
	if !ok {
		t.Fatalf("Expected Model type, got %T", newModel)
	}
	if statusBar.width != initialWidth {
		t.Errorf("Expected width to remain %d, got %d", initialWidth, statusBar.width)
	}

	// Test with string message
	newModel, cmd = model.Update("test message")
	if cmd != nil {
		t.Errorf("Expected Update to return nil command for string message, got %v", cmd)
	}
	statusBar, ok = newModel.(Model)
	if !ok {
		t.Fatalf("Expected Model type, got %T", newModel)
	}
	if statusBar.width != initialWidth {
		t.Errorf("Expected width to remain %d, got %d", initialWidth, statusBar.width)
	}
}

// TestView verifies that View renders correct content.
func TestView(t *testing.T) {
	tests := []struct {
		name      string
		charCount int
		lineCount int
		width     int
		contains  string
	}{
		{
			name:      "zero counts",
			charCount: 0,
			lineCount: 0,
			width:     80,
			contains:  "Chars: 0 | Lines: 0",
		},
		{
			name:      "positive counts",
			charCount: 100,
			lineCount: 5,
			width:     80,
			contains:  "Chars: 100 | Lines: 5",
		},
		{
			name:      "large counts",
			charCount: 999999,
			lineCount: 1000,
			width:     80,
			contains:  "Chars: 999999 | Lines: 1000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New()
			model.SetCharCount(tt.charCount)
			model.SetLineCount(tt.lineCount)
			model.width = tt.width

			view := model.View()

			if !strings.Contains(view, tt.contains) {
				t.Errorf("Expected view to contain %q, got %q", tt.contains, view)
			}
		})
	}
}

// TestViewUsesThemeStyles verifies that View uses theme styles correctly.
func TestViewUsesThemeStyles(t *testing.T) {
	model := New()
	model.SetCharCount(50)
	model.SetLineCount(3)

	view := model.View()

	// Get the expected style from theme
	expectedStyle := theme.StatusStyle().Width(model.width)

	// Render the same content with the expected style
	expectedView := expectedStyle.Render("Chars: 50 | Lines: 3")

	if view != expectedView {
		t.Errorf("View does not match theme style output.\nExpected: %q\nGot:      %q", expectedView, view)
	}
}

// TestViewWithZeroWidth verifies that View handles zero width gracefully.
func TestViewWithZeroWidth(t *testing.T) {
	model := New()
	model.width = 0
	model.SetCharCount(10)
	model.SetLineCount(1)

	view := model.View()

	// Should not panic and should still render something
	if view == "" {
		t.Error("Expected View to render something even with zero width")
	}
}

// TestViewWithLargeWidth verifies that View handles large width correctly.
func TestViewWithLargeWidth(t *testing.T) {
	model := New()
	model.width = 1000
	model.SetCharCount(100)
	model.SetLineCount(10)

	view := model.View()

	// Should render without issues
	if !strings.Contains(view, "Chars: 100 | Lines: 10") {
		t.Errorf("Expected view to contain status text, got %q", view)
	}
}

// TestSetCharCount verifies that SetCharCount updates the character count.
func TestSetCharCount(t *testing.T) {
	tests := []struct {
		name          string
		count         int
		expectedCount int
	}{
		{"zero", 0, 0},
		{"positive", 100, 100},
		{"large", 999999, 999999},
		{"negative is clamped to 0", -1, 0},
		{"large negative is clamped to 0", -1000, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New()
			model.SetCharCount(tt.count)

			if model.charCount != tt.expectedCount {
				t.Errorf("Expected charCount to be %d, got %d", tt.expectedCount, model.charCount)
			}
		})
	}
}

// TestSetLineCount verifies that SetLineCount updates the line count.
func TestSetLineCount(t *testing.T) {
	tests := []struct {
		name          string
		count         int
		expectedCount int
	}{
		{"zero", 0, 0},
		{"positive", 5, 5},
		{"large", 1000, 1000},
		{"negative is clamped to 0", -1, 0},
		{"large negative is clamped to 0", -500, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New()
			model.SetLineCount(tt.count)

			if model.lineCount != tt.expectedCount {
				t.Errorf("Expected lineCount to be %d, got %d", tt.expectedCount, model.lineCount)
			}
		})
	}
}

// TestGetCharCount verifies that GetCharCount returns the correct character count.
func TestGetCharCount(t *testing.T) {
	model := New()
	model.SetCharCount(42)

	if model.GetCharCount() != 42 {
		t.Errorf("Expected GetCharCount to return 42, got %d", model.GetCharCount())
	}
}

// TestGetLineCount verifies that GetLineCount returns the correct line count.
func TestGetLineCount(t *testing.T) {
	model := New()
	model.SetLineCount(7)

	if model.GetLineCount() != 7 {
		t.Errorf("Expected GetLineCount to return 7, got %d", model.GetLineCount())
	}
}

// TestGetWidth verifies that GetWidth returns the correct width.
func TestGetWidth(t *testing.T) {
	model := New()
	model.width = 120

	if model.GetWidth() != 120 {
		t.Errorf("Expected GetWidth to return 120, got %d", model.GetWidth())
	}
}

// TestRapidWindowResize verifies that the model handles rapid window resize events.
func TestRapidWindowResize(t *testing.T) {
	model := New()

	// Simulate rapid resize events
	widths := []int{80, 100, 120, 80, 60, 100, 80}
	for _, width := range widths {
		msg := tea.WindowSizeMsg{Width: width, Height: 24}
		newModel, cmd := model.Update(msg)

		if cmd != nil {
			t.Errorf("Expected Update to return nil command, got %v", cmd)
		}

		statusBar, ok := newModel.(Model)
		if !ok {
			t.Fatalf("Expected Model type, got %T", newModel)
		}

		if statusBar.width != width {
			t.Errorf("Expected width to be %d, got %d", width, statusBar.width)
		}

		model = statusBar
	}
}

// TestModelImplementsTeaModel verifies that Model implements the tea.Model interface.
func TestModelImplementsTeaModel(t *testing.T) {
	var _ tea.Model = New()
}

// TestViewReturnsStyledString verifies that View returns a properly styled string.
func TestViewReturnsStyledString(t *testing.T) {
	model := New()
	model.SetCharCount(25)
	model.SetLineCount(2)

	view := model.View()

	// Verify the view is not empty
	if view == "" {
		t.Error("Expected View to return a non-empty string")
	}

	// Verify the view contains the expected content
	if !strings.Contains(view, "Chars: 25") {
		t.Errorf("Expected view to contain 'Chars: 25', got %q", view)
	}
	if !strings.Contains(view, "Lines: 2") {
		t.Errorf("Expected view to contain 'Lines: 2', got %q", view)
	}
}

// TestStyleIntegration verifies that the status bar integrates correctly with the theme system.
func TestStyleIntegration(t *testing.T) {
	model := New()
	model.SetCharCount(10)
	model.SetLineCount(1)

	view := model.View()

	// The view should be styled with the theme's colors
	// We can't directly test the styling, but we can verify the view is not empty
	if view == "" {
		t.Error("Expected styled view to be non-empty")
	}
}

// TestMultipleUpdates verifies that multiple Update calls work correctly.
func TestMultipleUpdates(t *testing.T) {
	model := New()

	// First update
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 24})
	model = newModel.(Model)
	if model.width != 100 {
		t.Errorf("Expected width to be 100 after first update, got %d", model.width)
	}

	// Second update
	newModel, _ = model.Update(tea.WindowSizeMsg{Width: 150, Height: 30})
	model = newModel.(Model)
	if model.width != 150 {
		t.Errorf("Expected width to be 150 after second update, got %d", model.width)
	}

	// Third update
	newModel, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = newModel.(Model)
	if model.width != 80 {
		t.Errorf("Expected width to be 80 after third update, got %d", model.width)
	}
}

// TestViewWithDifferentWidths verifies that View renders correctly with different widths.
func TestViewWithDifferentWidths(t *testing.T) {
	model := New()
	model.SetCharCount(50)
	model.SetLineCount(5)

	widths := []int{20, 40, 80, 120, 200}
	for _, width := range widths {
		t.Run(fmt.Sprintf("width_%d", width), func(t *testing.T) {
			model.width = width
			view := model.View()

			// Should render without errors
			if view == "" {
				t.Errorf("Expected View to render something with width %d", width)
			}

			// Should contain the status text
			if !strings.Contains(view, "Chars: 50 | Lines: 5") {
				t.Errorf("Expected view to contain status text with width %d, got %q", width, view)
			}
		})
	}
}
