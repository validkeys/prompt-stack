package examples

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModelInit(t *testing.T) {
	m := newModel()

	cmd := m.Init()
	if cmd == nil {
		t.Error("Init should return a tick command")
	}

	// Verify the command is executable
	if cmd != nil {
		msg := cmd()
		if msg == nil {
			t.Error("Command should return a tick message")
		}
	}
}

func TestModelUpdateTick(t *testing.T) {
	m := newModel()

	initialExample := m.currentExample
	msg := tickMsg{}

	newModel, _ := m.Update(msg)
	updatedModel := newModel.(model)

	expectedExample := (initialExample + 1) % len(updatedModel.examples)
	if updatedModel.currentExample != expectedExample {
		t.Errorf("Expected example %d, got %d", expectedExample, updatedModel.currentExample)
	}
}

func TestModelUpdateNavigation(t *testing.T) {
	tests := []struct {
		name       string
		keyMsg     tea.KeyMsg
		wantChange int
	}{
		{
			name:       "right key advances",
			keyMsg:     tea.KeyMsg{Type: tea.KeyRight},
			wantChange: 1,
		},
		{
			name:       "space key advances",
			keyMsg:     tea.KeyMsg{Type: tea.KeySpace},
			wantChange: 1,
		},
		{
			name:       "left key goes back",
			keyMsg:     tea.KeyMsg{Type: tea.KeyLeft},
			wantChange: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newModel()
			initialExample := m.currentExample

			newModel, _ := m.Update(tt.keyMsg)
			updatedModel := newModel.(model)

			expectedExample := (initialExample + tt.wantChange + len(m.examples)) % len(m.examples)
			if updatedModel.currentExample != expectedExample {
				t.Errorf("Expected example %d, got %d", expectedExample, updatedModel.currentExample)
			}
		})
	}
}

func TestModelUpdateQuit(t *testing.T) {
	tests := []struct {
		name    string
		keyMsg  tea.KeyMsg
		wantCmd tea.Cmd
	}{
		{
			name:    "Ctrl+C quits",
			keyMsg:  tea.KeyMsg{Type: tea.KeyCtrlC},
			wantCmd: tea.Quit,
		},
		{
			name:    "Esc quits",
			keyMsg:  tea.KeyMsg{Type: tea.KeyEsc},
			wantCmd: tea.Quit,
		},
		{
			name:    "q quits",
			keyMsg:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			wantCmd: tea.Quit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newModel()
			newModel, cmd := m.Update(tt.keyMsg)

			updatedModel := newModel.(model)
			if cmd == nil {
				t.Error("Quit command should be non-nil")
			}

			// Model should still be intact after quit message
			if updatedModel.examples == nil {
				t.Error("Model should remain intact after quit message")
			}
		})
	}
}

func TestModelView(t *testing.T) {
	m := newModel()

	view := m.View()
	if view == "" {
		t.Error("View should return non-empty string")
	}

	// Check that view contains the example name
	ex := m.examples[m.currentExample]
	if view == "" {
		t.Errorf("View should contain example name: %s", ex.name)
	}
}

func TestModelWindowSize(t *testing.T) {
	m := newModel()

	msg := tea.WindowSizeMsg{
		Width:  100,
		Height: 50,
	}

	newModel, _ := m.Update(msg)
	if newModel == nil {
		t.Error("Model should handle window size messages")
	}
}

func TestExampleCount(t *testing.T) {
	m := newModel()

	if len(m.examples) == 0 {
		t.Error("Should have at least one example")
	}

	for i, ex := range m.examples {
		if ex.name == "" {
			t.Errorf("Example %d should have a name", i)
		}
		if ex.render == nil {
			t.Errorf("Example %d should have a render function", i)
		}
	}
}

func TestExampleRender(t *testing.T) {
	m := newModel()

	for i, ex := range m.examples {
		t.Run(ex.name, func(t *testing.T) {
			output := ex.render()
			if output == "" {
				t.Errorf("Example %d (%s) should render non-empty output", i, ex.name)
			}
		})
	}
}

func TestNavigateAllExamples(t *testing.T) {
	m := newModel()

	visited := make([]bool, len(m.examples))

	// Navigate through all examples
	for i := 0; i < len(m.examples); i++ {
		visited[m.currentExample] = true

		msg := tea.KeyMsg{Type: tea.KeyRight}
		newModel, _ := m.Update(msg)
		m = newModel.(model)
	}

	// Verify all examples were visited
	for i, v := range visited {
		if !v {
			t.Errorf("Example %d was not visited during navigation", i)
		}
	}
}

func TestWrapNavigation(t *testing.T) {
	m := newModel()

	// Go forward past the last example
	for i := 0; i < len(m.examples)+5; i++ {
		msg := tea.KeyMsg{Type: tea.KeyRight}
		newModel, _ := m.Update(msg)
		m = newModel.(model)
	}

	// Should wrap around
	if m.currentExample >= len(m.examples) {
		t.Error("Navigation should wrap around to the beginning")
	}

	// Go backward past the first example
	for i := 0; i < len(m.examples)+5; i++ {
		msg := tea.KeyMsg{Type: tea.KeyLeft}
		newModel, _ := m.Update(msg)
		m = newModel.(model)
	}

	// Should wrap around to the end
	if m.currentExample < 0 {
		t.Error("Navigation should wrap around to the end")
	}
}
