package editor

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	m := New()

	if len(m.placeholders) != 0 {
		t.Errorf("expected empty placeholders, got %d", len(m.placeholders))
	}

	if m.activeIndex != -1 {
		t.Errorf("expected activeIndex -1, got %d", m.activeIndex)
	}

	if m.isEditing {
		t.Error("expected isEditing false")
	}
}

func TestParsePlaceholders(t *testing.T) {
	content := "Hello {{text:name}}, here is a {{list:items}} list"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)

	placeholders := m.Placeholders()
	if len(placeholders) != 2 {
		t.Fatalf("expected 2 placeholders, got %d", len(placeholders))
	}

	// Check first placeholder
	if placeholders[0].Type != "text" {
		t.Errorf("expected type 'text', got '%s'", placeholders[0].Type)
	}
	if placeholders[0].Name != "name" {
		t.Errorf("expected name 'name', got '%s'", placeholders[0].Name)
	}
	if !placeholders[0].IsValid {
		t.Error("expected placeholder to be valid")
	}

	// Check second placeholder
	if placeholders[1].Type != "list" {
		t.Errorf("expected type 'list', got '%s'", placeholders[1].Type)
	}
	if placeholders[1].Name != "items" {
		t.Errorf("expected name 'items', got '%s'", placeholders[1].Name)
	}
}

func TestParsePlaceholdersEmpty(t *testing.T) {
	content := "No placeholders here"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)

	placeholders := m.Placeholders()
	if len(placeholders) != 0 {
		t.Errorf("expected 0 placeholders, got %d", len(placeholders))
	}
}

func TestParsePlaceholdersInvalid(t *testing.T) {
	content := "Invalid {{invalid:type}} and {{text:invalid-name}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)

	placeholders := m.Placeholders()
	if len(placeholders) != 2 {
		t.Fatalf("expected 2 placeholders, got %d", len(placeholders))
	}

	// First placeholder should be invalid (invalid type)
	if placeholders[0].IsValid {
		t.Error("expected first placeholder to be invalid")
	}

	// Second placeholder should be invalid (invalid name with hyphen)
	if placeholders[1].IsValid {
		t.Error("expected second placeholder to be invalid")
	}
}

func TestActivatePlaceholder(t *testing.T) {
	content := "{{text:first}} and {{text:second}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)

	// Activate first placeholder
	m = m.Update(ActivatePlaceholderMsg{Index: 0})

	if m.activeIndex != 0 {
		t.Errorf("expected activeIndex 0, got %d", m.activeIndex)
	}

	active := m.Active()
	if active == nil {
		t.Fatal("expected active placeholder, got nil")
	}

	if active.Name != "first" {
		t.Errorf("expected name 'first', got '%s'", active.Name)
	}

	if !active.IsActive {
		t.Error("expected IsActive true")
	}
}

func TestActivatePlaceholderInvalidIndex(t *testing.T) {
	content := "{{text:first}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)

	// Try to activate invalid index
	m = m.Update(ActivatePlaceholderMsg{Index: 5})

	if m.activeIndex != -1 {
		t.Errorf("expected activeIndex -1, got %d", m.activeIndex)
	}
}

func TestNextPlaceholder(t *testing.T) {
	content := "{{text:first}} {{text:second}} {{text:third}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)

	// Activate first
	m = m.Update(ActivatePlaceholderMsg{Index: 0})

	// Move to next
	m = m.Update(tea.KeyMsg{Type: tea.KeyTab})

	if m.activeIndex != 1 {
		t.Errorf("expected activeIndex 1, got %d", m.activeIndex)
	}

	active := m.Active()
	if active.Name != "second" {
		t.Errorf("expected name 'second', got '%s'", active.Name)
	}
}

func TestNextPlaceholderWrapAround(t *testing.T) {
	content := "{{text:first}} {{text:second}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)

	// Activate second
	m = m.Update(ActivatePlaceholderMsg{Index: 1})

	// Move to next (should wrap to first)
	m = m.Update(tea.KeyMsg{Type: tea.KeyTab})

	if m.activeIndex != 0 {
		t.Errorf("expected activeIndex 0, got %d", m.activeIndex)
	}
}

func TestPreviousPlaceholder(t *testing.T) {
	content := "{{text:first}} {{text:second}} {{text:third}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)

	// Activate second
	m = m.Update(ActivatePlaceholderMsg{Index: 1})

	// Move to previous
	m = m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})

	if m.activeIndex != 0 {
		t.Errorf("expected activeIndex 0, got %d", m.activeIndex)
	}

	active := m.Active()
	if active.Name != "first" {
		t.Errorf("expected name 'first', got '%s'", active.Name)
	}
}

func TestPreviousPlaceholderWrapAround(t *testing.T) {
	content := "{{text:first}} {{text:second}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)

	// Activate first
	m = m.Update(ActivatePlaceholderMsg{Index: 0})

	// Move to previous (should wrap to last)
	m = m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})

	if m.activeIndex != 1 {
		t.Errorf("expected activeIndex 1, got %d", m.activeIndex)
	}
}

func TestEditPlaceholder(t *testing.T) {
	content := "{{text:name}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)
	m = m.Update(ActivatePlaceholderMsg{Index: 0})

	// Type some text
	m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("H")})
	m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")})
	m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("o")})

	if !m.isEditing {
		t.Error("expected isEditing true")
	}

	if m.editValue != "Hello" {
		t.Errorf("expected editValue 'Hello', got '%s'", m.editValue)
	}
}

func TestEditPlaceholderBackspace(t *testing.T) {
	content := "{{text:name}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)
	m = m.Update(ActivatePlaceholderMsg{Index: 0})

	// Type text
	m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Hello")})

	// Delete last character
	m = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	if m.editValue != "Hell" {
		t.Errorf("expected editValue 'Hell', got '%s'", m.editValue)
	}
}

func TestSaveEdit(t *testing.T) {
	content := "{{text:name}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)
	m = m.Update(ActivatePlaceholderMsg{Index: 0})

	// Type text
	m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("World")})

	// Save with Enter
	m = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.isEditing {
		t.Error("expected isEditing false after save")
	}

	active := m.Active()
	if active.CurrentValue != "World" {
		t.Errorf("expected CurrentValue 'World', got '%s'", active.CurrentValue)
	}
}

func TestExitEditMode(t *testing.T) {
	content := "{{text:name}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)
	m = m.Update(ActivatePlaceholderMsg{Index: 0})

	// Type text
	m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Test")})

	// Exit with Escape
	m = m.Update(tea.KeyMsg{Type: tea.KeyEscape})

	if m.isEditing {
		t.Error("expected isEditing false after exit")
	}

	if m.editValue != "" {
		t.Errorf("expected editValue empty, got '%s'", m.editValue)
	}

	// Original value should not be saved
	active := m.Active()
	if active.CurrentValue != "" {
		t.Errorf("expected CurrentValue empty, got '%s'", active.CurrentValue)
	}
}

func TestEditPlaceholderListType(t *testing.T) {
	content := "{{list:items}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)
	m = m.Update(ActivatePlaceholderMsg{Index: 0})

	// Try to type (should not work for list type)
	m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Test")})

	if m.isEditing {
		t.Error("expected isEditing false for list type")
	}
}

func TestReplacePlaceholderText(t *testing.T) {
	content := "Hello {{text:name}}!"
	ph := Placeholder{
		Type:         "text",
		Name:         "name",
		StartPos:     6,
		EndPos:       19,
		CurrentValue: "World",
		IsValid:      true,
	}

	result := ReplacePlaceholder(content, ph)
	expected := "Hello World!"

	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestReplacePlaceholderList(t *testing.T) {
	content := "Items:\n{{list:items}}"
	ph := Placeholder{
		Type:       "list",
		Name:       "items",
		StartPos:   7,
		EndPos:     21,
		ListValues: []string{"First", "Second", "Third"},
		IsValid:    true,
	}

	result := ReplacePlaceholder(content, ph)
	expected := "Items:\n- First\n- Second\n- Third"

	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestReplacePlaceholderListEmpty(t *testing.T) {
	content := "Items:\n{{list:items}}"
	ph := Placeholder{
		Type:       "list",
		Name:       "items",
		StartPos:   7,
		EndPos:     21,
		ListValues: []string{},
		IsValid:    true,
	}

	result := ReplacePlaceholder(content, ph)
	expected := "Items:\n"

	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestFindPlaceholderAtPosition(t *testing.T) {
	placeholders := []Placeholder{
		{StartPos: 0, EndPos: 10, Name: "first"},
		{StartPos: 15, EndPos: 25, Name: "second"},
		{StartPos: 30, EndPos: 40, Name: "third"},
	}

	// Find at position 5 (in first placeholder)
	idx := FindPlaceholderAtPosition(placeholders, 5)
	if idx != 0 {
		t.Errorf("expected index 0, got %d", idx)
	}

	// Find at position 20 (in second placeholder)
	idx = FindPlaceholderAtPosition(placeholders, 20)
	if idx != 1 {
		t.Errorf("expected index 1, got %d", idx)
	}

	// Find at position 35 (in third placeholder)
	idx = FindPlaceholderAtPosition(placeholders, 35)
	if idx != 2 {
		t.Errorf("expected index 2, got %d", idx)
	}

	// Find at position 12 (not in any placeholder)
	idx = FindPlaceholderAtPosition(placeholders, 12)
	if idx != -1 {
		t.Errorf("expected index -1, got %d", idx)
	}
}

func TestIsValidPlaceholderType(t *testing.T) {
	tests := []struct {
		typ  string
		want bool
	}{
		{"text", true},
		{"list", true},
		{"invalid", false},
		{"", false},
		{"TEXT", false},
		{"LIST", false},
	}

	for _, tt := range tests {
		t.Run(tt.typ, func(t *testing.T) {
			got := isValidPlaceholderType(tt.typ)
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestIsValidPlaceholderName(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"valid_name", true},
		{"ValidName123", true},
		{"_underscore", true},
		{"", false},
		{"invalid-name", false},
		{"invalid name", false},
		{"invalid.name", false},
		{"invalid/name", false},
		{"invalid@name", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidPlaceholderName(tt.name)
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestImmutability(t *testing.T) {
	content := "{{text:first}} {{text:second}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m1 := New()
	m2 := m1.Update(msg)

	// m1 should not be modified
	if len(m1.placeholders) != 0 {
		t.Error("expected m1 to remain unchanged")
	}

	// m2 should have placeholders
	if len(m2.placeholders) != 2 {
		t.Errorf("expected m2 to have 2 placeholders, got %d", len(m2.placeholders))
	}
}

func TestNavigationWithNoPlaceholders(t *testing.T) {
	m := New()

	// Try to navigate with no placeholders
	m = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	if m.activeIndex != -1 {
		t.Errorf("expected activeIndex -1, got %d", m.activeIndex)
	}

	m = m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if m.activeIndex != -1 {
		t.Errorf("expected activeIndex -1, got %d", m.activeIndex)
	}
}

func TestEditWithoutActivePlaceholder(t *testing.T) {
	m := New()

	// Try to type without active placeholder
	m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Test")})

	if m.isEditing {
		t.Error("expected isEditing false")
	}

	if m.editValue != "" {
		t.Errorf("expected editValue empty, got '%s'", m.editValue)
	}
}

func TestMultiplePlaceholdersSameType(t *testing.T) {
	content := "{{text:first}} and {{text:second}} and {{text:third}}"
	msg := ParsePlaceholdersMsg{Content: content}

	m := New()
	m = m.Update(msg)

	// Activate first
	m = m.Update(ActivatePlaceholderMsg{Index: 0})
	m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("One")})

	// Move to second (auto-saves first)
	m = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Two")})

	// Move to third (auto-saves second)
	m = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Three")})

	// Save third placeholder
	m = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Check all values
	placeholders := m.Placeholders()
	if placeholders[0].CurrentValue != "One" {
		t.Errorf("expected 'One', got '%s'", placeholders[0].CurrentValue)
	}
	if placeholders[1].CurrentValue != "Two" {
		t.Errorf("expected 'Two', got '%s'", placeholders[1].CurrentValue)
	}
	if placeholders[2].CurrentValue != "Three" {
		t.Errorf("expected 'Three', got '%s'", placeholders[2].CurrentValue)
	}
}
