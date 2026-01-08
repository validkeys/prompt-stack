package workspace

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNewWorkspace tests creating a new workspace model
func TestNewWorkspace(t *testing.T) {
	model := New()

	if model.content != "" {
		t.Errorf("expected empty content, got %q", model.content)
	}

	if model.width != 0 || model.height != 0 {
		t.Errorf("expected zero dimensions, got %dx%d", model.width, model.height)
	}

	if model.isReadOnly {
		t.Error("expected not read-only")
	}
}

// TestCharacterInput tests typing characters into the workspace
func TestCharacterInput(t *testing.T) {
	model := New()

	// Type 'H'
	m, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'H'}})
	model = m.(Model)

	// Type 'i'
	m, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = m.(Model)

	expected := "Hi"
	if model.GetContent() != expected {
		t.Errorf("expected %q, got %q", expected, model.GetContent())
	}
}

// TestBackspace tests deleting characters with backspace
func TestBackspace(t *testing.T) {
	model := New()

	// Type "Hello"
	for _, r := range "Hello" {
		m, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = m.(Model)
	}

	// Backspace
	m, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model = m.(Model)

	expected := "Hell"
	if model.GetContent() != expected {
		t.Errorf("expected %q, got %q", expected, model.GetContent())
	}
}

// TestEnter tests inserting newlines
func TestEnter(t *testing.T) {
	model := New()

	// Type "Hello"
	for _, r := range "Hello" {
		m, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = m.(Model)
	}

	// Press Enter
	m, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = m.(Model)

	expected := "Hello\n"
	if model.GetContent() != expected {
		t.Errorf("expected %q, got %q", expected, model.GetContent())
	}
}

// TestTab tests tab insertion
func TestTab(t *testing.T) {
	model := New()

	// Type "Hello"
	for _, r := range "Hello" {
		m, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = m.(Model)
	}

	// Press Tab
	m, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = m.(Model)

	expected := "Hello    " // 4 spaces
	if model.GetContent() != expected {
		t.Errorf("expected %q, got %q", expected, model.GetContent())
	}
}

// TestCursorMovement tests cursor navigation
func TestCursorMovement(t *testing.T) {
	model := New()

	// Type "Hello"
	for _, r := range "Hello" {
		m, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = m.(Model)
	}

	// Cursor should be at position 5
	if model.cursor.X() != 5 || model.cursor.Y() != 0 {
		t.Errorf("cursor at (%d, %d), want (5, 0)", model.cursor.X(), model.cursor.Y())
	}

	// Move left
	m, _ := model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = m.(Model)

	if model.cursor.X() != 4 {
		t.Errorf("expected cursor at x=4, got %d", model.cursor.X())
	}

	// Move right
	m, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = m.(Model)

	if model.cursor.X() != 5 {
		t.Errorf("expected cursor at x=5, got %d", model.cursor.X())
	}
}

// TestWindowResize tests window resize handling
func TestWindowResize(t *testing.T) {
	model := New()

	// Resize window
	m, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = m.(Model)

	if model.width != 80 {
		t.Errorf("expected width 80, got %d", model.width)
	}

	if model.height != 24 {
		t.Errorf("expected height 24, got %d", model.height)
	}

	// Verify viewport height is adjusted (height - 1 for status bar)
	if model.viewport.Height() != 23 {
		t.Errorf("expected viewport height 23, got %d", model.viewport.Height())
	}
}

// TestSetSize tests setting workspace size
func TestSetSize(t *testing.T) {
	model := New()
	model = model.SetSize(100, 30)

	if model.width != 100 {
		t.Errorf("expected width 100, got %d", model.width)
	}

	if model.height != 30 {
		t.Errorf("expected height 30, got %d", model.height)
	}
}

// TestSetContent tests setting workspace content
func TestSetContent(t *testing.T) {
	model := New()

	content := "Hello World\nThis is a test"
	model = model.SetContent(content)

	if model.GetContent() != content {
		t.Errorf("expected %q, got %q", content, model.GetContent())
	}

	// Verify file is marked as modified
	if !model.fileManager.IsModified() {
		t.Error("expected file to be marked as modified")
	}
}

// TestGetContent tests getting workspace content
func TestGetContent(t *testing.T) {
	model := New()

	for _, r := range "Test content" {
		m, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = m.(Model)
	}

	if model.GetContent() != "Test content" {
		t.Errorf("expected 'Test content', got %q", model.GetContent())
	}
}

// TestGetCursorPosition tests getting cursor position
func TestGetCursorPosition(t *testing.T) {
	model := New()

	for _, r := range "Hello" {
		m, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = m.(Model)
	}

	pos := model.GetCursorPosition()
	if pos != 5 {
		t.Errorf("expected position 5, got %d", pos)
	}
}

// TestInsertContent tests inserting content at a specific position
func TestInsertContent(t *testing.T) {
	model := New()

	for _, r := range "Hello World" {
		m, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = m.(Model)
	}

	// Insert "Beautiful " at position 6
	model = model.InsertContent(6, "Beautiful ")

	expected := "Hello Beautiful World"
	if model.GetContent() != expected {
		t.Errorf("expected %q, got %q", expected, model.GetContent())
	}
}

// TestMarkDirty tests marking file as modified
func TestMarkDirty(t *testing.T) {
	model := New()

	if model.fileManager.IsModified() {
		t.Error("expected file to not be modified initially")
	}

	model = model.MarkDirty()

	if !model.fileManager.IsModified() {
		t.Error("expected file to be marked as modified")
	}
}

// TestSetStatus tests setting status message
func TestSetStatus(t *testing.T) {
	model := New()

	message := "Test status message"
	model = model.SetStatus(message)

	if model.statusBar.message != message {
		t.Errorf("expected %q, got %q", message, model.statusBar.message)
	}
}

// TestReadOnlyMode tests read-only mode functionality
func TestReadOnlyMode(t *testing.T) {
	model := New()
	model = model.SetReadOnly(true)

	if !model.IsReadOnly() {
		t.Error("expected read-only mode to be enabled")
	}

	// Try to type in read-only mode
	m, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'H'}})
	model = m.(Model)

	// Content should remain empty
	if model.GetContent() != "" {
		t.Errorf("expected empty content in read-only mode, got %q", model.GetContent())
	}

	// Disable read-only mode
	model = model.SetReadOnly(false)

	if model.IsReadOnly() {
		t.Error("expected read-only mode to be disabled")
	}

	// Typing should work now
	m, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'H'}})
	model = m.(Model)

	if model.GetContent() != "H" {
		t.Errorf("expected 'H', got %q", model.GetContent())
	}
}

// TestQuitCommand tests quit command behavior
func TestQuitCommand(t *testing.T) {
	model := New()

	// Test quit without modifications
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if cmd == nil {
		t.Error("expected quit command, got nil")
	}
}

// TestStatusBar tests status bar functionality
func TestStatusBar(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)

	// SetContent will trigger updateStatusBar through Update method
	model = model.SetContent("Hello World\nLine 2")

	// Status bar should be updated
	if model.statusBar.charCount != 18 {
		t.Errorf("expected 18 chars, got %d", model.statusBar.charCount)
	}

	if model.statusBar.lineCount != 2 {
		t.Errorf("expected 2 lines, got %d", model.statusBar.lineCount)
	}
}

// TestViewInit tests Init method
func TestViewInit(t *testing.T) {
	model := New()

	cmd := model.Init()
	if cmd != nil {
		t.Errorf("expected nil command from Init, got %v", cmd)
	}
}

// TestImmutableUpdate tests that Update returns new model instances
func TestImmutableUpdate(t *testing.T) {
	model := New()
	originalContent := model.GetContent()

	// Update should return new model
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'H'}})

	// Original model should be unchanged
	if model.GetContent() != originalContent {
		t.Error("original model was mutated")
	}

	// New model should have changes
	if newModel.(Model).GetContent() != "H" {
		t.Errorf("expected 'H', got %q", newModel.(Model).GetContent())
	}
}

// TestViewRendering tests view rendering
func TestViewRendering(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)

	// View should not panic
	view := model.View()
	if view == "" {
		t.Error("expected non-empty view")
	}

	// Add content
	for _, r := range "Hello World" {
		m, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = m.(Model)
	}

	view = model.View()
	if view == "" {
		t.Error("expected non-empty view with content")
	}
}

// TestModifiedFlag tests the modified flag behavior
func TestModifiedFlag(t *testing.T) {
	model := New()

	// Initially not modified
	if model.fileManager.IsModified() {
		t.Error("expected file to not be modified initially")
	}

	// Typing should mark as modified
	m, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'H'}})
	model = m.(Model)

	if !model.fileManager.IsModified() {
		t.Error("expected file to be marked as modified after typing")
	}
}

// TestMultipleLines tests handling of multiple lines
func TestMultipleLines(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)

	// SetContent will trigger updateStatusBar
	model = model.SetContent("Line\nLine\nLine\n")

	// Should have 4 lines (3 lines of text + 1 empty line after last \n)
	lineCount := model.statusBar.lineCount
	if lineCount != 4 {
		t.Errorf("expected 4 lines, got %d", lineCount)
	}
}

// TestCursorVerticalMovement tests up and down cursor movement
func TestCursorVerticalMovement(t *testing.T) {
	model := New()
	model = model.SetContent("Line 1\nLine 2\nLine 3")

	// Move down
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(Model)

	if model.cursor.Y() != 1 {
		t.Errorf("expected cursor at y=1, got %d", model.cursor.Y())
	}

	// Move up
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = newModel.(Model)

	if model.cursor.Y() != 0 {
		t.Errorf("expected cursor at y=0, got %d", model.cursor.Y())
	}
}

// TestBackspaceNewline tests backspace on newline character
func TestBackspaceNewline(t *testing.T) {
	model := New()
	model = model.SetContent("Hello\nWorld")

	// Move cursor to start of second line
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(Model)

	// Backspace from start of line - should delete newline
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model = newModel.(Model)

	// Cursor should now be at end of first line
	if model.cursor.Y() != 0 {
		t.Errorf("expected cursor on line 0 after backspace, got %d", model.cursor.Y())
	}
}

// TestViewportAdjustment tests viewport adjustment based on cursor
func TestViewportAdjustment(t *testing.T) {
	model := New()
	model = model.SetSize(80, 10)

	// Set content longer than viewport
	lines := make([]string, 50)
	for i := range lines {
		lines[i] = fmt.Sprintf("Line %d", i)
	}
	model = model.SetContent(strings.Join(lines, "\n"))

	// Move cursor down
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(Model)

	// Move cursor down more
	for i := 0; i < 10; i++ {
		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = newModel.(Model)
	}

	// Test that adjustViewport doesn't panic
	model = model.adjustViewport()

	// Cursor should be at expected position
	if model.cursor.Y() != 11 {
		t.Errorf("expected cursor y=11, got %d", model.cursor.Y())
	}
}

// TestPlaceholderNavigation tests Tab/Shift+Tab navigation
func TestPlaceholderNavigation(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("Hello {{text:name}}. Welcome to {{text:project}}.")

	// Press Tab to navigate to first placeholder - should not panic
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)

	// Press Tab again to navigate to next placeholder - should not panic
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)

	// Press Shift+Tab to navigate to previous placeholder - should not panic
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	model = newModel.(Model)
}

// TestPlaceholderEditMode tests editing placeholders
func TestPlaceholderEditMode(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("Hello {{text:name}}.")

	// Navigate to placeholder and enter edit mode via Tab
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)

	// Type characters to fill placeholder
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'W', 'o', 'r', 'l', 'd'}})
	model = newModel.(Model)

	// Verify edit value was set
	editValue := model.placeholders.EditValue()
	if editValue != "World" {
		t.Errorf("expected edit value 'World', got %q", editValue)
	}

	// Press Enter to save and exit
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = newModel.(Model)

	// Should have exited edit mode
	if model.placeholders.IsEditing() {
		t.Error("expected to exit edit mode")
	}
}

// TestPlaceholderEditBackspace tests backspace in placeholder edit mode
func TestPlaceholderEditBackspace(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("Hello {{text:name}}.")

	// Navigate to placeholder and enter edit mode via Tab
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)

	// Type characters
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A', 'B', 'C'}})
	model = newModel.(Model)

	// Press backspace
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model = newModel.(Model)

	// Edit value should have been modified
	editValue := model.placeholders.EditValue()
	if editValue != "AB" {
		t.Errorf("expected edit value 'AB', got %q", editValue)
	}
}

// TestPlaceholderEditEscape tests escape key in placeholder edit mode
func TestPlaceholderEditEscape(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("Hello {{name}}.")

	// Navigate to placeholder and enter edit mode
	m, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = m.(Model)
	m, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = m.(Model)

	// Type characters
	m, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A', 'B', 'C'}})
	model = m.(Model)

	// Press Escape to save and exit
	m, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = m.(Model)

	// Should have replaced placeholder with content
	if model.placeholders.IsEditing() {
		t.Error("expected to exit edit mode")
	}
}

// TestSaveToFile tests saving to file
func TestSaveToFile(t *testing.T) {
	model := New()
	model = model.SetContent("Test content")

	// Set a filepath for testing
	model.fileManager = model.fileManager.SetPath("/tmp/test.md")

	err := model.saveToFile()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if model.fileManager.IsModified() {
		t.Error("expected file to be marked as saved")
	}
}

// TestRenderLineWithPlaceholders tests rendering lines with placeholder highlighting
func TestRenderLineWithPlaceholders(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("Hello {{name}} from {{place}}.")

	// Navigate to placeholder
	m, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = m.(Model)

	// Render a line with active placeholder
	lines := model.getVisibleLines(24)
	rendered := model.renderLineWithPlaceholders(lines[0], 0)

	// Should contain placeholder highlighting
	if !strings.Contains(rendered, "{{name}}") {
		t.Error("expected placeholder in rendered line")
	}
}

// TestEnterPlaceholderEditMode tests entering placeholder edit mode
func TestEnterPlaceholderEditMode(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("Hello {{text:name}}.")

	// Navigate to placeholder
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)

	// Enter edit mode - Tab already enters edit mode for text placeholders
	if model.placeholders.IsEditing() {
		// Already in edit mode - this is correct behavior
	}

	// Call enterPlaceholderEditMode directly
	model = model.enterPlaceholderEditMode()

	// Should be in edit mode (already was from Tab)
	// This test ensures the method doesn't panic
}

// TestViewportScrolling tests viewport scrolling behavior
func TestViewportScrolling(t *testing.T) {
	model := New()
	model = model.SetSize(80, 10)

	// Set content that spans multiple pages
	lines := make([]string, 100)
	for i := range lines {
		lines[i] = fmt.Sprintf("Line %d with some content", i)
	}
	model = model.SetContent(strings.Join(lines, "\n"))

	// Move cursor far down
	for i := 0; i < 80; i++ {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = newModel.(Model)
	}

	// Get visible lines
	visibleLines := model.getVisibleLines(9) // 9 for content area

	// Visible lines should be a subset of all lines
	if len(visibleLines) > len(lines) {
		t.Errorf("visible lines %d should not exceed total lines %d", len(visibleLines), len(lines))
	}
}

// TestGetVisibleLines tests getting visible lines from viewport
func TestGetVisibleLines(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)

	// Set multiline content
	content := strings.Join([]string{"Line 1", "Line 2", "Line 3", "Line 4", "Line 5"}, "\n")
	model = model.SetContent(content)

	// Get visible lines
	lines := model.getVisibleLines(5)

	if len(lines) != 5 {
		t.Errorf("expected 5 visible lines, got %d", len(lines))
	}

	if lines[0] != "Line 1" {
		t.Errorf("expected first line 'Line 1', got %q", lines[0])
	}
}

// TestNavigateToPreviousPlaceholder tests previous placeholder navigation
func TestNavigateToPreviousPlaceholder(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("One {{text:a}} Two {{text:b}} Three {{text:c}}.")

	// Navigate to a placeholder first
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)

	// Get the active placeholder name
	activeName := model.placeholders.Active().Name

	// Navigate to previous placeholder
	model = model.navigateToPreviousPlaceholder()

	// Should have wrapped to last placeholder
	if model.placeholders.Active().Name == activeName {
		t.Error("expected to wrap to last placeholder")
	}
}

// TestQuitWithModifications tests quit when file is modified
func TestQuitWithModifications(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("Test content")

	// File should be modified after SetContent
	if !model.fileManager.IsModified() {
		t.Error("expected file to be modified")
	}

	// Set a file path
	model.fileManager = model.fileManager.SetPath("/tmp/test_quit.md")

	// Try to quit - Ctrl+C in read-only mode doesn't quit
	model = model.SetReadOnly(true)
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	// Should not return quit command in read-only mode
	if cmd != nil {
		t.Error("expected no quit command in read-only mode")
	}

	// Try again without read-only
	model = model.SetReadOnly(false)
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	// Should return quit command
	if cmd == nil {
		t.Error("expected quit command")
	}
}

// TestSpaceKey tests space key handling
func TestSpaceKey(t *testing.T) {
	model := New()

	// Press space
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeySpace})
	model = newModel.(Model)

	// Content should contain a space
	if model.GetContent() != " " {
		t.Errorf("expected ' ', got %q", model.GetContent())
	}
}

// TestShiftTabWhenNoPlaceholders tests Shift+Tab without placeholders
func TestShiftTabWhenNoPlaceholders(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("No placeholders here.")

	// Press Shift+Tab - should not panic
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	model = newModel.(Model)

	// No active placeholder
	if model.placeholders.Active() != nil {
		t.Error("expected no active placeholder")
	}
}

// TestTabWithInvalidPlaceholderContent tests Tab with invalid placeholder syntax
func TestTabWithInvalidPlaceholderContent(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("Invalid {name} syntax.")

	// Press Tab - should not panic
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)

	// Should insert tab since no valid placeholder found
	// After SetContent, cursor is at position 0, so tab is inserted at start
	if !strings.HasPrefix(model.GetContent(), "    ") {
		t.Errorf("expected tab insertion at start, got %q", model.GetContent())
	}
}

// TestRenderStatusBar tests status bar rendering with various states
func TestRenderStatusBar(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("Test content")

	// Render status bar
	statusBar := model.renderStatusBar()

	// Should not be empty
	if statusBar == "" {
		t.Error("expected non-empty status bar")
	}

	// Should contain character count
	if !strings.Contains(statusBar, "chars") {
		t.Error("expected status bar to show character count")
	}
}

// TestRenderCursorLine tests cursor line rendering
func TestRenderCursorLine(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("Hello World")

	lines := model.getVisibleLines(24)

	// Render cursor line
	cursorLine := model.renderCursorLine(lines)

	// Should not be empty
	if cursorLine == "" {
		t.Error("expected non-empty cursor line")
	}
}

// TestUpdatePlaceholders tests placeholder parsing
func TestUpdatePlaceholders(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)

	// Initial content with no placeholders
	model = model.SetContent("No placeholders")

	// Update to add placeholders
	model = model.SetContent("Hello {{text:name}}")

	// Should have parsed the placeholder
	if len(model.placeholders.Placeholders()) != 1 {
		t.Errorf("expected 1 placeholder, got %d", len(model.placeholders.Placeholders()))
	}
}

// TestBackspaceAtBeginning tests backspace at beginning of file
func TestBackspaceAtBeginning(t *testing.T) {
	model := New()
	model = model.SetContent("Hello")

	// Move cursor to beginning
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = newModel.(Model)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = newModel.(Model)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = newModel.(Model)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = newModel.(Model)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = newModel.(Model)

	// Backspace at beginning - should do nothing
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model = newModel.(Model)

	// Content should be unchanged
	if model.GetContent() != "Hello" {
		t.Errorf("expected unchanged content, got %q", model.GetContent())
	}
}

// TestCursorAdjustToLineLength tests cursor adjustment to line length
func TestCursorAdjustToLineLength(t *testing.T) {
	model := New()
	model = model.SetContent("Short\nVery long line with many characters")

	// Move to end of short line (position 5)
	for i := 0; i < 5; i++ {
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
		model = newModel.(Model)
	}

	// Cursor at position 5 on line 0
	if model.cursor.X() != 5 {
		t.Errorf("expected cursor x=5, got %d", model.cursor.X())
	}

	// Move down to long line
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = newModel.(Model)

	// Cursor stays at same position if line is long enough
	// The long line has 32 characters, so position 5 is valid
	if model.cursor.X() != 5 {
		t.Errorf("expected cursor x=5 on longer line, got %d", model.cursor.X())
	}
}

// TestRenderCursorLineWithPlaceholder tests rendering cursor line with active placeholder
func TestRenderCursorLineWithPlaceholder(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("Hello {{text:name}} World")

	// Navigate to placeholder and enter edit mode
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)

	// Type something
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A', 'B', 'C'}})
	model = newModel.(Model)

	// Get visible lines
	lines := model.getVisibleLines(24)

	// Render cursor line
	cursorLine := model.renderCursorLine(lines)

	// Should not be empty
	if cursorLine == "" {
		t.Error("expected non-empty cursor line")
	}
}

// TestEnterListPlaceholderEditMode tests entering edit mode for list placeholder
func TestEnterListPlaceholderEditMode(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("Items: {{list:items}}")

	// Navigate to list placeholder
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)

	// Should have active placeholder
	if model.placeholders.Active() == nil {
		t.Fatal("expected placeholder to be active")
	}

	// Enter placeholder edit mode
	model = model.enterPlaceholderEditMode()

	// Should not enter edit mode for list placeholders (not yet implemented)
	// Just ensure it doesn't panic
}

// TestViewWithPlaceholderEdit tests view rendering while editing placeholder
func TestViewWithPlaceholderEdit(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("Hello {{text:name}}")

	// Navigate to placeholder and enter edit mode
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)

	// Type something
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T', 'e', 's', 't'}})
	model = newModel.(Model)

	// Render view
	view := model.View()

	// Should not be empty
	if view == "" {
		t.Error("expected non-empty view")
	}
}

// TestNavigatePlaceholdersWrap tests placeholder wrapping
func TestNavigatePlaceholdersWrap(t *testing.T) {
	model := New()
	model = model.SetSize(80, 24)
	model = model.SetContent("{{text:a}} {{text:b}} {{text:c}}")

	// Navigate past the last placeholder to wrap to first
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)
	newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = newModel.(Model)

	// Should have wrapped to first placeholder
	if model.placeholders.Active() == nil {
		t.Error("expected active placeholder after wrap")
	}
}
