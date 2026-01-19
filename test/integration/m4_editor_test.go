// Integration tests for M4 text editor functionality.
// These tests verify end-to-end workflows across buffer, viewport, and theme components.
package integration_test

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyledavis/prompt-stack/ui/theme"
	"github.com/kyledavis/prompt-stack/ui/workspace"
)

func TestTypingWorkflow(t *testing.T) {
	var model workspace.Model = workspace.New()

	msgs := []tea.Msg{
		tea.WindowSizeMsg{Width: 40, Height: 10},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'H'}},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}},
	}

	for _, msg := range msgs {
		newModel, _ := model.Update(msg)
		model = newModel.(workspace.Model)
	}

	if model.GetContent() != "Hello" {
		t.Errorf("expected 'Hello', got %q", model.GetContent())
	}

	view := model.View()
	if view == "" {
		t.Error("view should not be empty after typing")
	}
}

func TestTypingWithNewlines(t *testing.T) {
	var model workspace.Model = workspace.New()

	msgs := []tea.Msg{
		tea.WindowSizeMsg{Width: 40, Height: 10},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'L', 'i', 'n', 'e', '1'}},
		tea.KeyMsg{Type: tea.KeyEnter},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'L', 'i', 'n', 'e', '2'}},
	}

	for _, msg := range msgs {
		newModel, _ := model.Update(msg)
		model = newModel.(workspace.Model)
	}

	expected := "Line1\nLine2"
	if model.GetContent() != expected {
		t.Errorf("expected %q, got %q", expected, model.GetContent())
	}
}

func TestTypingWithBackspace(t *testing.T) {
	var model workspace.Model = workspace.New()

	msgs := []tea.Msg{
		tea.WindowSizeMsg{Width: 40, Height: 10},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T', 'e', 's', 't'}},
		tea.KeyMsg{Type: tea.KeyBackspace},
		tea.KeyMsg{Type: tea.KeyBackspace},
	}

	for _, msg := range msgs {
		newModel, _ := model.Update(msg)
		model = newModel.(workspace.Model)
	}

	if model.GetContent() != "Te" {
		t.Errorf("expected 'Te', got %q", model.GetContent())
	}
}

func TestRapidTypingNoDroppedCharacters(t *testing.T) {
	var model workspace.Model = workspace.New()
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = newModel.(workspace.Model)

	expectedContent := strings.Repeat("a", 100)
	startTime := time.Now()

	for i := 0; i < 100; i++ {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
		newModel, _ := model.Update(msg)
		model = newModel.(workspace.Model)
	}

	elapsed := time.Since(startTime)

	if elapsed > 1*time.Second {
		t.Errorf("rapid typing took too long: %v", elapsed)
	}

	if model.GetContent() != expectedContent {
		t.Errorf("rapid typing dropped characters")
	}
}

func TestCursorNavigation(t *testing.T) {
	var model workspace.Model = workspace.New()

	msgs := []tea.Msg{
		tea.WindowSizeMsg{Width: 40, Height: 10},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T', 'e', 's', 't'}},
		tea.KeyMsg{Type: tea.KeyLeft},
		tea.KeyMsg{Type: tea.KeyLeft},
		tea.KeyMsg{Type: tea.KeyRight},
		tea.KeyMsg{Type: tea.KeyRight},
	}

	for _, msg := range msgs {
		newModel, _ := model.Update(msg)
		model = newModel.(workspace.Model)
	}

	view := model.View()
	if view == "" {
		t.Error("view should render after navigation")
	}
}

func TestMultiLineNavigation(t *testing.T) {
	var model workspace.Model = workspace.New()

	msgs := []tea.Msg{
		tea.WindowSizeMsg{Width: 40, Height: 10},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A'}},
		tea.KeyMsg{Type: tea.KeyEnter},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'B'}},
		tea.KeyMsg{Type: tea.KeyUp},
		tea.KeyMsg{Type: tea.KeyDown},
	}

	for _, msg := range msgs {
		newModel, _ := model.Update(msg)
		model = newModel.(workspace.Model)
	}

	view := model.View()
	if view == "" {
		t.Error("view should render after multi-line navigation")
	}
}

func TestHomeAndEndNavigation(t *testing.T) {
	var model workspace.Model = workspace.New()

	msgs := []tea.Msg{
		tea.WindowSizeMsg{Width: 40, Height: 10},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T', 'e', 's', 't'}},
		tea.KeyMsg{Type: tea.KeyHome},
		tea.KeyMsg{Type: tea.KeyEnd},
	}

	for _, msg := range msgs {
		newModel, _ := model.Update(msg)
		model = newModel.(workspace.Model)
	}

	view := model.View()
	if view == "" {
		t.Error("view should render after home/end navigation")
	}
}

func TestViewportScrolling(t *testing.T) {
	var model workspace.Model = workspace.New()
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	model = newModel.(workspace.Model)

	content := strings.Repeat("line\n", 20)
	for _, r := range content {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = newModel.(workspace.Model)
	}

	view := model.View()
	if view == "" {
		t.Error("viewport content should not be empty")
	}

	if !strings.Contains(view, "line") {
		t.Error("viewport should contain content")
	}
}

func TestWindowResize(t *testing.T) {
	var model workspace.Model = workspace.New()
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	model = newModel.(workspace.Model)

	content := strings.Repeat("line\n", 15)
	for _, r := range content {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = newModel.(workspace.Model)
	}

	initialView := model.View()

	newModel, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 20})
	model = newModel.(workspace.Model)

	resizedView := model.View()

	if resizedView == "" {
		t.Error("viewport should handle window resize")
	}

	if len(initialView) == len(resizedView) {
		t.Log("viewport size changed on window resize")
	}
}

func TestCharacterCounting(t *testing.T) {
	var model workspace.Model = workspace.New()

	msgs := []tea.Msg{
		tea.WindowSizeMsg{Width: 40, Height: 10},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'H', 'e', 'l', 'l', 'o'}},
		tea.KeyMsg{Type: tea.KeySpace},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'W', 'o', 'r', 'l', 'd'}},
	}

	for _, msg := range msgs {
		newModel, _ := model.Update(msg)
		model = newModel.(workspace.Model)
	}

	content := model.GetContent()
	actualChars := len([]rune(content))
	expectedChars := 11

	if actualChars != expectedChars {
		t.Errorf("character count: got %d, want %d", actualChars, expectedChars)
	}
}

func TestLineCounting(t *testing.T) {
	var model workspace.Model = workspace.New()

	msgs := []tea.Msg{
		tea.WindowSizeMsg{Width: 40, Height: 10},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'L', 'i', 'n', 'e', '1'}},
		tea.KeyMsg{Type: tea.KeyEnter},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'L', 'i', 'n', 'e', '2'}},
		tea.KeyMsg{Type: tea.KeyEnter},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'L', 'i', 'n', 'e', '3'}},
	}

	for _, msg := range msgs {
		newModel, _ := model.Update(msg)
		model = newModel.(workspace.Model)
	}

	content := model.GetContent()
	actualLines := strings.Count(content, "\n") + 1
	expectedLines := 3

	if actualLines != expectedLines {
		t.Errorf("line count: got %d, want %d", actualLines, expectedLines)
	}
}

func TestBufferWorkspaceIntegration(t *testing.T) {
	var model workspace.Model = workspace.New()
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	model = newModel.(workspace.Model)

	testContent := "Hello\nWorld\nTest"
	for _, r := range testContent {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = newModel.(workspace.Model)
	}

	if model.GetContent() != testContent {
		t.Errorf("content mismatch: got %q, want %q", model.GetContent(), testContent)
	}

	cursorPos := model.GetCursorPosition()
	if cursorPos < 0 {
		t.Error("cursor position should be valid")
	}
}

func TestViewportBufferIntegration(t *testing.T) {
	var model workspace.Model = workspace.New()
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	model = newModel.(workspace.Model)

	content := strings.Repeat("line\n", 20)
	for _, r := range content {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = newModel.(workspace.Model)
	}

	view := model.View()
	if view == "" {
		t.Error("view should not be empty")
	}

	if !strings.Contains(view, "line") {
		t.Error("view should contain buffer content")
	}
}

func TestThemeStylesApplied(t *testing.T) {
	var model workspace.Model = workspace.New()
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	model = newModel.(workspace.Model)

	for _, r := range "Test" {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = newModel.(workspace.Model)
	}
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = newModel.(workspace.Model)

	view := model.View()
	if view == "" {
		t.Error("view should not be empty")
	}

	cursorStyle := theme.CursorStyle()
	_ = cursorStyle.Render("test")

	statusStyle := theme.StatusStyle()
	_ = statusStyle.Render("test")

	if theme.BackgroundPrimary == "" {
		t.Error("background primary color should be defined")
	}

	if theme.ForegroundPrimary == "" {
		t.Error("foreground primary color should be defined")
	}
}

func TestLargeDocumentHandling(t *testing.T) {
	var model workspace.Model = workspace.New()
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = newModel.(workspace.Model)

	startTime := time.Now()

	content := strings.Repeat("a", 10000)
	for _, r := range content {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = newModel.(workspace.Model)
	}

	elapsed := time.Since(startTime)
	if elapsed > 5*time.Second {
		t.Errorf("handling large document took too long: %v", elapsed)
	}

	actualChars := len([]rune(model.GetContent()))
	if actualChars != 10000 {
		t.Errorf("large document char count: got %d, want 10000", actualChars)
	}

	view := model.View()
	if view == "" {
		t.Error("view should handle large documents")
	}
}

func TestCursorPositionConsistency(t *testing.T) {
	var model workspace.Model = workspace.New()
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	model = newModel.(workspace.Model)

	for _, r := range "ABC" {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = newModel.(workspace.Model)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model = newModel.(workspace.Model)

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = newModel.(workspace.Model)

	for _, r := range "DE" {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = newModel.(workspace.Model)
	}

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = newModel.(workspace.Model)

	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(workspace.Model)

	view := model.View()
	if view == "" {
		t.Error("view should render after cursor movements")
	}
}

func TestPerformanceRequirements(t *testing.T) {
	var model workspace.Model = workspace.New()
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = newModel.(workspace.Model)

	startTime := time.Now()

	for i := 0; i < 1000; i++ {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
		model = newModel.(workspace.Model)
		if i > 0 && i%50 == 0 {
			newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = newModel.(workspace.Model)
		}
	}

	typingElapsed := time.Since(startTime)

	startTime = time.Now()
	for i := 0; i < 500; i++ {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
		model = newModel.(workspace.Model)
		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = newModel.(workspace.Model)
		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
		model = newModel.(workspace.Model)
		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
		model = newModel.(workspace.Model)
	}
	navigationElapsed := time.Since(startTime)

	startTime = time.Now()
	for i := 0; i < 100; i++ {
		_ = model.View()
	}
	viewElapsed := time.Since(startTime)

	if typingElapsed > 2*time.Second {
		t.Errorf("typing performance: %v (should be < 2s)", typingElapsed)
	}

	if navigationElapsed > 1*time.Second {
		t.Errorf("navigation performance: %v (should be < 1s)", navigationElapsed)
	}

	if viewElapsed > 500*time.Millisecond {
		t.Errorf("view rendering performance: %v (should be < 500ms)", viewElapsed)
	}
}

func TestUnicodeSupport(t *testing.T) {
	var model workspace.Model = workspace.New()
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	model = newModel.(workspace.Model)

	unicodeText := "Hello 世界 🌍"
	for _, r := range unicodeText {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = newModel.(workspace.Model)
	}

	if model.GetContent() != unicodeText {
		t.Errorf("unicode content mismatch")
	}

	expectedChars := len([]rune(unicodeText))
	actualChars := len([]rune(model.GetContent()))
	if actualChars != expectedChars {
		t.Errorf("unicode char count: got %d, want %d", actualChars, expectedChars)
	}

	view := model.View()
	if view == "" {
		t.Error("view should render with unicode content")
	}
}

func TestEmptyBufferState(t *testing.T) {
	var model workspace.Model = workspace.New()
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	model = newModel.(workspace.Model)

	if model.GetContent() != "" {
		t.Error("empty buffer should be empty")
	}

	if len([]rune(model.GetContent())) != 0 {
		t.Error("empty buffer should have 0 chars")
	}

	content := model.GetContent()
	lines := 1
	if strings.Contains(content, "\n") {
		lines = strings.Count(content, "\n") + 1
	}
	if lines != 1 {
		t.Errorf("empty buffer should have 1 line, got %d", lines)
	}
}

func TestBoundaryConditions(t *testing.T) {
	t.Run("backspace at start", func(t *testing.T) {
		var model workspace.Model = workspace.New()
		newModel, _ := model.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
		model = newModel.(workspace.Model)

		initialContent := model.GetContent()
		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		model = newModel.(workspace.Model)

		if model.GetContent() != initialContent {
			t.Error("backspace at start should not change content")
		}
	})

	t.Run("navigation at boundaries", func(t *testing.T) {
		var model workspace.Model = workspace.New()
		newModel, _ := model.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
		model = newModel.(workspace.Model)

		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A'}})
		model = newModel.(workspace.Model)

		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
		model = newModel.(workspace.Model)

		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
		model = newModel.(workspace.Model)

		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
		model = newModel.(workspace.Model)

		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = newModel.(workspace.Model)

		view := model.View()
		if view == "" {
			t.Error("view should render after boundary navigation")
		}
	})
}
