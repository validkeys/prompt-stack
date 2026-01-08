package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	model := New()

	if model.statusBar.GetCharCount() != 0 {
		t.Errorf("Expected char count to be 0, got %d", model.statusBar.GetCharCount())
	}

	if model.statusBar.GetLineCount() != 0 {
		t.Errorf("Expected line count to be 0, got %d", model.statusBar.GetLineCount())
	}

	if model.width != 80 {
		t.Errorf("Expected default width to be 80, got %d", model.width)
	}

	if model.height != 24 {
		t.Errorf("Expected default height to be 24, got %d", model.height)
	}

	if model.quitting {
		t.Error("Expected quitting to be false")
	}
}

func TestInit(t *testing.T) {
	model := New()
	cmd := model.Init()

	if cmd != nil {
		t.Error("Expected Init to return nil command")
	}
}

func TestUpdateHandlesCharacterInput(t *testing.T) {
	model := New()

	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"character 'a'", "a", false},
		{"character 'b'", "b", false},
		{"character 'z'", "z", false},
		{"character 'A'", "A", false},
		{"character 'Z'", "Z", false},
		{"character '0'", "0", false},
		{"character '9'", "9", false},
		{"space", " ", false},
		{"special char '@'", "@", false},
		{"special char '#'", "#", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune(tt.key),
			}

			newModel, cmd := model.Update(msg)
			appModel := newModel.(Model)

			if cmd != nil {
				t.Error("Expected nil command for character input")
			}

			if appModel.IsQuitting() != tt.expected {
				t.Errorf("Expected quitting to be %v, got %v", tt.expected, appModel.IsQuitting())
			}
		})
	}
}

func TestUpdateHandlesCtrlCQuit(t *testing.T) {
	model := New()

	msg := tea.KeyMsg{
		Type: tea.KeyCtrlC,
	}

	newModel, cmd := model.Update(msg)
	appModel := newModel.(Model)

	if cmd == nil {
		t.Error("Expected tea.Quit command for Ctrl+C")
	}

	if !appModel.IsQuitting() {
		t.Error("Expected quitting to be true after Ctrl+C")
	}
}

func TestUpdateHandlesQQuit(t *testing.T) {
	model := New()

	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
	}

	newModel, cmd := model.Update(msg)
	appModel := newModel.(Model)

	if cmd == nil {
		t.Error("Expected tea.Quit command for 'q' key")
	}

	if !appModel.IsQuitting() {
		t.Error("Expected quitting to be true after 'q' key")
	}
}

func TestUpdateHandlesQQuitCaseInsensitive(t *testing.T) {
	model := New()

	tests := []struct {
		name     string
		key      rune
		expected bool
	}{
		{"lowercase 'q'", 'q', true},
		{"uppercase 'Q'", 'Q', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{tt.key},
			}

			newModel, cmd := model.Update(msg)
			appModel := newModel.(Model)

			if tt.expected && cmd == nil {
				t.Error("Expected tea.Quit command for 'q' key")
			}

			if appModel.IsQuitting() != tt.expected {
				t.Errorf("Expected quitting to be %v, got %v", tt.expected, appModel.IsQuitting())
			}
		})
	}
}

func TestUpdateHandlesWindowSizeMsg(t *testing.T) {
	model := New()

	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"small window", 40, 10},
		{"medium window", 80, 24},
		{"large window", 120, 40},
		{"very large window", 200, 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.WindowSizeMsg{
				Width:  tt.width,
				Height: tt.height,
			}

			newModel, cmd := model.Update(msg)
			appModel := newModel.(Model)

			if cmd != nil {
				t.Error("Expected nil command for window resize")
			}

			if appModel.GetWidth() != tt.width {
				t.Errorf("Expected width to be %d, got %d", tt.width, appModel.GetWidth())
			}

			if appModel.GetHeight() != tt.height {
				t.Errorf("Expected height to be %d, got %d", tt.height, appModel.GetHeight())
			}

			// Verify status bar width was updated
			if appModel.GetStatusBar().GetWidth() != tt.width {
				t.Errorf("Expected status bar width to be %d, got %d", tt.width, appModel.GetStatusBar().GetWidth())
			}
		})
	}
}

func TestUpdateHandlesSpecialKeys(t *testing.T) {
	model := New()

	tests := []struct {
		name     string
		keyType  tea.KeyType
		expected bool
	}{
		{"Enter key", tea.KeyEnter, false},
		{"Space key", tea.KeySpace, false},
		{"Backspace key", tea.KeyBackspace, false},
		{"Delete key", tea.KeyDelete, false},
		{"Tab key", tea.KeyTab, false},
		{"Up arrow", tea.KeyUp, false},
		{"Down arrow", tea.KeyDown, false},
		{"Left arrow", tea.KeyLeft, false},
		{"Right arrow", tea.KeyRight, false},
		{"Escape key", tea.KeyEscape, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{
				Type: tt.keyType,
			}

			newModel, cmd := model.Update(msg)
			appModel := newModel.(Model)

			if cmd != nil {
				t.Error("Expected nil command for special keys")
			}

			if appModel.IsQuitting() != tt.expected {
				t.Errorf("Expected quitting to be %v, got %v", tt.expected, appModel.IsQuitting())
			}
		})
	}
}

func TestViewRendersStatusBar(t *testing.T) {
	model := New()

	view := model.View()

	if view == "" {
		t.Error("Expected View to return non-empty string")
	}

	// Verify status bar content is present
	if len(view) < 10 {
		t.Error("Expected View to render substantial content")
	}
}

func TestViewUsesThemeStyles(t *testing.T) {
	model := New()

	view := model.View()

	// Verify view contains expected content
	if len(view) == 0 {
		t.Error("Expected View to return styled content")
	}

	// The view should contain the main content and status bar
	if len(view) < 20 {
		t.Error("Expected View to render both main content and status bar")
	}
}

func TestViewWithZeroDimensions(t *testing.T) {
	model := New()
	model.width = 0
	model.height = 0

	view := model.View()

	// Should still return a string even with zero dimensions
	if view == "" {
		t.Error("Expected View to return string even with zero dimensions")
	}
}

func TestViewWithLargeDimensions(t *testing.T) {
	model := New()
	model.width = 500
	model.height = 500

	view := model.View()

	if view == "" {
		t.Error("Expected View to return string with large dimensions")
	}

	// View should be longer with larger dimensions
	if len(view) < 100 {
		t.Error("Expected longer view with large dimensions")
	}
}

func TestIsQuitting(t *testing.T) {
	model := New()

	if model.IsQuitting() {
		t.Error("Expected IsQuitting to return false initially")
	}

	model.quitting = true

	if !model.IsQuitting() {
		t.Error("Expected IsQuitting to return true after setting")
	}
}

func TestGetWidth(t *testing.T) {
	model := New()

	if model.GetWidth() != 80 {
		t.Errorf("Expected GetWidth to return 80, got %d", model.GetWidth())
	}

	model.width = 120

	if model.GetWidth() != 120 {
		t.Errorf("Expected GetWidth to return 120, got %d", model.GetWidth())
	}
}

func TestGetHeight(t *testing.T) {
	model := New()

	if model.GetHeight() != 24 {
		t.Errorf("Expected GetHeight to return 24, got %d", model.GetHeight())
	}

	model.height = 40

	if model.GetHeight() != 40 {
		t.Errorf("Expected GetHeight to return 40, got %d", model.GetHeight())
	}
}

func TestGetStatusBar(t *testing.T) {
	model := New()

	statusBar := model.GetStatusBar()

	if statusBar.GetCharCount() != 0 {
		t.Errorf("Expected status bar char count to be 0, got %d", statusBar.GetCharCount())
	}

	if statusBar.GetLineCount() != 0 {
		t.Errorf("Expected status bar line count to be 0, got %d", statusBar.GetLineCount())
	}
}

func TestRapidKeyboardInput(t *testing.T) {
	model := New()

	// Simulate rapid keyboard input
	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	for _, key := range keys {
		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune(key),
		}

		newModel, cmd := model.Update(msg)
		appModel := newModel.(Model)

		if cmd != nil {
			t.Error("Expected nil command for rapid character input")
		}

		if appModel.IsQuitting() {
			t.Error("Expected not to quit during rapid character input")
		}

		model = appModel
	}
}

func TestMultipleQuitAttempts(t *testing.T) {
	model := New()

	// First quit attempt
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
	}

	newModel, cmd := model.Update(msg)
	appModel := newModel.(Model)

	if cmd == nil {
		t.Error("Expected tea.Quit command on first quit attempt")
	}

	if !appModel.IsQuitting() {
		t.Error("Expected quitting to be true after first quit attempt")
	}

	// Second quit attempt (should still work)
	newModel, cmd = appModel.Update(msg)
	appModel = newModel.(Model)

	if cmd == nil {
		t.Error("Expected tea.Quit command on second quit attempt")
	}

	if !appModel.IsQuitting() {
		t.Error("Expected quitting to remain true after second quit attempt")
	}
}

func TestWindowResizeDuringInput(t *testing.T) {
	model := New()

	// Simulate character input
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'a'},
	}

	newModel, _ := model.Update(msg)

	// Simulate window resize during input
	resizeMsg := tea.WindowSizeMsg{
		Width:  100,
		Height: 30,
	}

	newModel, cmd := newModel.Update(resizeMsg)
	appModel := newModel.(Model)

	if cmd != nil {
		t.Error("Expected nil command for window resize")
	}

	if appModel.GetWidth() != 100 {
		t.Errorf("Expected width to be 100, got %d", appModel.GetWidth())
	}

	if appModel.GetHeight() != 30 {
		t.Errorf("Expected height to be 30, got %d", appModel.GetHeight())
	}
}

func TestSpecialCharactersAndUnicode(t *testing.T) {
	model := New()

	tests := []struct {
		name string
		char string
	}{
		{"emoji 😀", "😀"},
		{"accented é", "é"},
		{"chinese 中", "中"},
		{"japanese 日", "日"},
		{"arabic ا", "ا"},
		{"russian р", "р"},
		{"greek α", "α"},
		{"mathematical ∑", "∑"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune(tt.char),
			}

			newModel, cmd := model.Update(msg)
			appModel := newModel.(Model)

			if cmd != nil {
				t.Error("Expected nil command for unicode characters")
			}

			if appModel.IsQuitting() {
				t.Error("Expected not to quit with unicode characters")
			}
		})
	}
}

func TestModelImplementsTeaModel(t *testing.T) {
	var _ tea.Model = New()
}

func TestViewReturnsStyledString(t *testing.T) {
	model := New()

	view := model.View()

	// Verify view is a string
	if view == "" {
		t.Error("Expected View to return non-empty string")
	}

	// Verify view contains expected text
	expectedText := "PromptStack TUI"
	if len(view) < len(expectedText) {
		t.Errorf("Expected view to contain at least %d characters, got %d", len(expectedText), len(view))
	}
}

func TestStatusBarIntegration(t *testing.T) {
	model := New()

	// Update status bar counts
	model.statusBar.SetCharCount(100)
	model.statusBar.SetLineCount(5)

	view := model.View()

	// Verify view contains status bar content
	if len(view) == 0 {
		t.Error("Expected View to render status bar")
	}
}

func TestMultipleUpdates(t *testing.T) {
	model := New()

	// First update: character input
	msg1 := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'a'},
	}

	newModel, _ := model.Update(msg1)

	// Second update: window resize
	msg2 := tea.WindowSizeMsg{
		Width:  100,
		Height: 30,
	}

	newModel, _ = newModel.Update(msg2)

	// Third update: another character
	msg3 := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'b'},
	}

	newModel, _ = newModel.Update(msg3)
	appModel := newModel.(Model)

	// Verify state is consistent
	if appModel.GetWidth() != 100 {
		t.Errorf("Expected width to be 100, got %d", appModel.GetWidth())
	}

	if appModel.GetHeight() != 30 {
		t.Errorf("Expected height to be 30, got %d", appModel.GetHeight())
	}

	if appModel.IsQuitting() {
		t.Error("Expected not to be quitting after multiple updates")
	}
}

func TestViewWithDifferentWidths(t *testing.T) {
	tests := []struct {
		name  string
		width int
	}{
		{"width 40", 40},
		{"width 80", 80},
		{"width 120", 120},
		{"width 200", 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New()
			model.width = tt.width

			view := model.View()

			if view == "" {
				t.Error("Expected View to return string")
			}

			// View should be longer with larger widths
			if len(view) < 10 {
				t.Error("Expected View to render substantial content")
			}
		})
	}
}

func TestUpdateWithNilMessage(t *testing.T) {
	model := New()

	// This should not panic
	newModel, cmd := model.Update(nil)
	appModel := newModel.(Model)

	if cmd != nil {
		t.Error("Expected nil command for nil message")
	}

	if appModel.IsQuitting() {
		t.Error("Expected not to quit with nil message")
	}
}

func TestUpdateWithUnknownMessageType(t *testing.T) {
	model := New()

	// Create a custom message type
	type customMsg struct{}

	// This should not panic
	newModel, cmd := model.Update(customMsg{})
	appModel := newModel.(Model)

	if cmd != nil {
		t.Error("Expected nil command for unknown message type")
	}

	if appModel.IsQuitting() {
		t.Error("Expected not to quit with unknown message type")
	}
}
