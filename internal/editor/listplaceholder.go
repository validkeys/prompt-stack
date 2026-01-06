package editor

import (
	"strings"
)

// ListEditState represents the state when editing a list placeholder
type ListEditState struct {
	PlaceholderIndex int      // Index of the placeholder being edited
	Items            []string // Current list items
	SelectedItem     int      // Index of currently selected item (-1 if none)
	EditMode         bool     // true when editing a specific item
	EditValue        string   // Value being edited for the selected item
}

// NewListEditState creates a new list edit state from a placeholder
func NewListEditState(placeholderIndex int, ph Placeholder) ListEditState {
	return ListEditState{
		PlaceholderIndex: placeholderIndex,
		Items:            ph.ListValues,
		SelectedItem:     0,
		EditMode:         false,
		EditValue:        "",
	}
}

// AddItem adds a new item to the list
func (s *ListEditState) AddItem() {
	s.Items = append(s.Items, "")
	s.SelectedItem = len(s.Items) - 1
	s.EditMode = true
	s.EditValue = ""
}

// DeleteItem deletes the currently selected item
func (s *ListEditState) DeleteItem() bool {
	if s.SelectedItem < 0 || s.SelectedItem >= len(s.Items) {
		return false
	}

	s.Items = append(s.Items[:s.SelectedItem], s.Items[s.SelectedItem+1:]...)

	// Adjust selection
	if s.SelectedItem >= len(s.Items) {
		s.SelectedItem = len(s.Items) - 1
	}

	return true
}

// EditItem enters edit mode for the selected item
func (s *ListEditState) EditItem() bool {
	if s.SelectedItem < 0 || s.SelectedItem >= len(s.Items) {
		return false
	}

	s.EditMode = true
	s.EditValue = s.Items[s.SelectedItem]
	return true
}

// SaveEdit saves the current edit value to the selected item
func (s *ListEditState) SaveEdit() bool {
	if s.SelectedItem < 0 || s.SelectedItem >= len(s.Items) {
		return false
	}

	s.Items[s.SelectedItem] = s.EditValue
	s.EditMode = false
	s.EditValue = ""
	return true
}

// CancelEdit cancels the current edit without saving
func (s *ListEditState) CancelEdit() {
	s.EditMode = false
	s.EditValue = ""
}

// MoveUp moves selection to the previous item
func (s *ListEditState) MoveUp() {
	if s.SelectedItem > 0 {
		s.SelectedItem--
	}
}

// MoveDown moves selection to the next item
func (s *ListEditState) MoveDown() {
	if s.SelectedItem < len(s.Items)-1 {
		s.SelectedItem++
	}
}

// GetPlaceholder returns the updated placeholder with current list values
func (s *ListEditState) GetPlaceholder(original Placeholder) Placeholder {
	ph := original
	ph.ListValues = s.Items
	return ph
}

// FormatAsMarkdown formats the list items as markdown bullet list
func (s *ListEditState) FormatAsMarkdown() string {
	if len(s.Items) == 0 {
		return ""
	}

	items := make([]string, len(s.Items))
	for i, val := range s.Items {
		items[i] = "- " + val
	}
	return strings.Join(items, "\n")
}
