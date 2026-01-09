package workspace

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyledavis/prompt-stack/internal/editor"
)

// handlePlaceholderEdit handles key events when in placeholder edit mode.
func (m Model) handlePlaceholderEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		// Exit placeholder edit mode and save the value
		newModel := m
		ph := newModel.placeholders.Active()
		if ph != nil {
			// Replace placeholder in content with the filled value
			newContent := editor.ReplacePlaceholder(newModel.buffer.Content(), *ph)
			newModel.buffer.SetContent(newContent)
			newModel = newModel.updatePlaceholders()
			newModel = newModel.markDirty()
			newModel = newModel.updateStatusBar()
		}
		newModel.placeholders = newModel.placeholders.Update(editor.ExitEditModeMsg{})
		return newModel, nil

	case tea.KeyBackspace:
		// Delete character from edit value
		newModel := m
		editValue := newModel.placeholders.EditValue()
		if len(editValue) > 0 {
			newModel.placeholders = newModel.placeholders.Update(editor.EditPlaceholderMsg{Value: editValue[:len(editValue)-1]})
		}
		return newModel, nil

	case tea.KeyEnter:
		// Exit placeholder edit mode and save the value
		newModel := m
		ph := newModel.placeholders.Active()
		if ph != nil {
			// Replace placeholder in content with the filled value
			newContent := editor.ReplacePlaceholder(newModel.buffer.Content(), *ph)
			newModel.buffer.SetContent(newContent)
			newModel = newModel.updatePlaceholders()
			newModel = newModel.markDirty()
			newModel = newModel.updateStatusBar()
		}
		newModel.placeholders = newModel.placeholders.Update(editor.ExitEditModeMsg{})
		return newModel, nil

	case tea.KeyRunes:
		// Append characters to edit value
		newModel := m
		editValue := newModel.placeholders.EditValue()
		newModel.placeholders = newModel.placeholders.Update(editor.EditPlaceholderMsg{Value: editValue + string(msg.Runes)})
		newModel = newModel.updateStatusBar()
		return newModel, nil
	}

	return m, nil
}
