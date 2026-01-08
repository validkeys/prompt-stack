package testutil

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TypeText simulates typing text into a Bubble Tea model
// It sends each character as a separate KeyMsg with runes
func TypeText(model tea.Model, text string) tea.Model {
	for _, r := range text {
		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{r},
		}
		model, _ = model.Update(msg)
	}
	return model
}

// PressKey simulates pressing a key on a Bubble Tea model
// It sends a single KeyMsg with the specified key type
func PressKey(model tea.Model, keyType tea.KeyType) tea.Model {
	msg := tea.KeyMsg{
		Type: keyType,
	}
	model, _ = model.Update(msg)
	return model
}

// AssertState compares two states and reports an error if they differ
// This is a generic assertion helper that can be used with any state type
func AssertState(t *testing.T, got, want interface{}) {
	t.Helper()

	if got != want {
		t.Errorf("state mismatch:\ngot:  %v\nwant: %v", got, want)
	}
}

// AssertStatef compares two states and reports a formatted error if they differ
func AssertStatef(t *testing.T, got, want interface{}, format string, args ...interface{}) {
	t.Helper()

	if got != want {
		t.Errorf(format, args...)
	}
}

// WorkflowStep represents a single step in a workflow simulation
type WorkflowStep struct {
	Name string
	Msg  tea.Msg
}

// SimulateWorkflow simulates a sequence of messages sent to a model
// This is useful for testing complex user interactions
func SimulateWorkflow(model tea.Model, steps []WorkflowStep) tea.Model {
	for _, step := range steps {
		model, _ = model.Update(step.Msg)
	}
	return model
}

// SimulateWorkflowWithCallback simulates a workflow and calls a callback after each step
// This allows for intermediate state verification
func SimulateWorkflowWithCallback(model tea.Model, steps []WorkflowStep, callback func(step WorkflowStep, model tea.Model)) tea.Model {
	for _, step := range steps {
		model, _ = model.Update(step.Msg)
		if callback != nil {
			callback(step, model)
		}
	}
	return model
}

// TypeTextWithCursor simulates typing text and returns the model with cursor position
// This is useful for models that track cursor position
func TypeTextWithCursor(model tea.Model, text string) tea.Model {
	for _, r := range text {
		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{r},
		}
		model, _ = model.Update(msg)
	}
	return model
}

// PressKeyWithModifiers simulates pressing a key with modifiers (Alt)
// Note: Ctrl and Shift are handled via specific key types (e.g., tea.KeyCtrlC)
func PressKeyWithModifiers(model tea.Model, keyType tea.KeyType, alt bool) tea.Model {
	msg := tea.KeyMsg{
		Type: keyType,
		Alt:  alt,
	}
	model, _ = model.Update(msg)
	return model
}

// ResizeWindow simulates a window resize event
func ResizeWindow(model tea.Model, width, height int) tea.Model {
	msg := tea.WindowSizeMsg{
		Width:  width,
		Height: height,
	}
	model, _ = model.Update(msg)
	return model
}

// ClickMouse simulates a mouse click event
func ClickMouse(model tea.Model, x, y int, button tea.MouseButton) tea.Model {
	msg := tea.MouseMsg{
		X:      x,
		Y:      y,
		Button: button,
		Type:   tea.MouseLeft,
	}
	model, _ = model.Update(msg)
	return model
}

// AssertNoCommand verifies that no command was returned from an Update
func AssertNoCommand(t *testing.T, cmd tea.Cmd) {
	t.Helper()

	if cmd != nil {
		t.Errorf("expected no command, got %v", cmd)
	}
}

// AssertCommand verifies that a command was returned from an Update
func AssertCommand(t *testing.T, cmd tea.Cmd) {
	t.Helper()

	if cmd == nil {
		t.Error("expected command, got nil")
	}
}

// AssertCommandType verifies that a command was returned
// Note: tea.Cmd is a function type, so we can only verify it's not nil
// For specific command type checking, you'll need to execute the command and check its behavior
func AssertCommandType(t *testing.T, cmd tea.Cmd) {
	t.Helper()

	if cmd == nil {
		t.Error("expected command, got nil")
	}
}

// ModelWithContent is an interface for models that have content
type ModelWithContent interface {
	GetContent() string
}

// ModelWithCursor is an interface for models that track cursor position
type ModelWithCursor interface {
	GetCursorX() int
	GetCursorY() int
}

// AssertContent asserts that a model's content matches the expected value
func AssertContent(t *testing.T, model ModelWithContent, expected string) {
	t.Helper()

	got := model.GetContent()
	if got != expected {
		t.Errorf("content mismatch:\ngot:  %q\nwant: %q", got, expected)
	}
}

// AssertCursorPosition asserts that a model's cursor position matches the expected values
func AssertCursorPosition(t *testing.T, model ModelWithCursor, expectedX, expectedY int) {
	t.Helper()

	gotX := model.GetCursorX()
	gotY := model.GetCursorY()

	if gotX != expectedX || gotY != expectedY {
		t.Errorf("cursor position mismatch:\ngot:  (%d, %d)\nwant: (%d, %d)", gotX, gotY, expectedX, expectedY)
	}
}

// RunUpdateSequence runs a sequence of updates on a model and returns the final model
func RunUpdateSequence(model tea.Model, msgs []tea.Msg) tea.Model {
	for _, msg := range msgs {
		model, _ = model.Update(msg)
	}
	return model
}

// CollectCommands runs a sequence of updates and collects all returned commands
func CollectCommands(model tea.Model, msgs []tea.Msg) []tea.Cmd {
	commands := make([]tea.Cmd, 0, len(msgs))

	for _, msg := range msgs {
		var cmd tea.Cmd
		model, cmd = model.Update(msg)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}

	return commands
}

// AssertNoStateMutation verifies that a model's state hasn't changed after an update
// This is useful for testing that certain messages don't modify state
func AssertNoStateMutation(t *testing.T, original, updated tea.Model) {
	t.Helper()

	if original != updated {
		t.Errorf("unexpected state mutation:\noriginal: %v\nupdated:  %v", original, updated)
	}
}

// FormatModelDiff creates a formatted diff between two models for debugging
func FormatModelDiff(original, updated interface{}) string {
	return fmt.Sprintf("Original: %+v\nUpdated:  %+v", original, updated)
}
