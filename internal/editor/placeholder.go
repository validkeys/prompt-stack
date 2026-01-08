package editor

import (
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Placeholder represents a template variable in a composition
type Placeholder struct {
	Type         string   `json:"type"`          // "text" or "list"
	Name         string   `json:"name"`          // placeholder name
	StartPos     int      `json:"start_pos"`     // position in content
	EndPos       int      `json:"end_pos"`       // position in content
	CurrentValue string   `json:"current_value"` // current filled value (for text)
	ListValues   []string `json:"list_values"`   // current filled values (for list)
	IsValid      bool     `json:"is_valid"`      // whether syntax is valid
	IsActive     bool     `json:"is_active"`     // whether currently selected
}

// Manager handles placeholder detection and editing
type Manager struct {
	placeholders []Placeholder
	activeIndex  int
	isEditing    bool
	editValue    string // Current value being edited
}

// New creates a new placeholder manager
func New() Manager {
	return Manager{
		placeholders: []Placeholder{},
		activeIndex:  -1,
		isEditing:    false,
		editValue:    "",
	}
}

// Update handles placeholder-related messages
func (m Manager) Update(msg tea.Msg) Manager {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case ParsePlaceholdersMsg:
		return m.parsePlaceholders(msg.Content)
	case ActivatePlaceholderMsg:
		return m.activatePlaceholder(msg.Index)
	case EditPlaceholderMsg:
		return m.editPlaceholder(msg.Value)
	case ExitEditModeMsg:
		return m.exitEditMode()
	}
	return m
}

// handleKey processes keyboard input for placeholder navigation and editing
func (m Manager) handleKey(msg tea.KeyMsg) Manager {
	switch msg.Type {
	case tea.KeyTab:
		return m.nextPlaceholder()
	case tea.KeyShiftTab:
		return m.previousPlaceholder()
	case tea.KeyEnter:
		if m.isEditing {
			return m.saveEdit()
		}
	case tea.KeyEscape:
		if m.isEditing {
			return m.exitEditMode()
		}
	case tea.KeyRunes:
		if m.isEditing {
			newValue := m.editValue + string(msg.Runes)
			return m.editPlaceholder(newValue)
		}
	case tea.KeyBackspace:
		if m.isEditing {
			if len(m.editValue) > 0 {
				newValue := m.editValue[:len(m.editValue)-1]
				return m.editPlaceholder(newValue)
			}
		}
	}
	return m
}

// parsePlaceholders extracts all placeholders from content
func (m Manager) parsePlaceholders(content string) Manager {
	newModel := m
	newModel.placeholders = extractPlaceholders(content)

	// Reset active index if out of bounds
	if newModel.activeIndex >= len(newModel.placeholders) {
		newModel.activeIndex = -1
	}

	return newModel
}

// activatePlaceholder activates a placeholder by index
func (m Manager) activatePlaceholder(index int) Manager {
	newModel := m

	// Deactivate current
	if newModel.activeIndex >= 0 && newModel.activeIndex < len(newModel.placeholders) {
		ph := &newModel.placeholders[newModel.activeIndex]
		ph.IsActive = false
	}

	// Activate new
	if index >= 0 && index < len(newModel.placeholders) {
		newModel.activeIndex = index
		ph := &newModel.placeholders[index]
		ph.IsActive = true

		// Initialize edit value for text placeholders and enter edit mode
		if ph.Type == "text" {
			newModel.editValue = ph.CurrentValue
			newModel.isEditing = true
		}
	} else {
		newModel.activeIndex = -1
	}

	return newModel
}

// editPlaceholder updates the edit value for the active placeholder
func (m Manager) editPlaceholder(value string) Manager {
	if m.activeIndex < 0 || m.activeIndex >= len(m.placeholders) {
		return m
	}

	ph := m.placeholders[m.activeIndex]
	if ph.Type != "text" {
		return m
	}

	newModel := m
	newModel.editValue = value
	newModel.isEditing = true
	return newModel
}

// saveEdit saves the current edit value to the active placeholder
func (m Manager) saveEdit() Manager {
	if m.activeIndex < 0 || m.activeIndex >= len(m.placeholders) {
		return m
	}

	ph := m.placeholders[m.activeIndex]
	if ph.Type != "text" {
		return m
	}

	newModel := m
	newModel.placeholders[m.activeIndex].CurrentValue = newModel.editValue
	newModel.isEditing = false
	return newModel
}

// exitEditMode exits placeholder edit mode without saving
func (m Manager) exitEditMode() Manager {
	newModel := m
	newModel.isEditing = false
	newModel.editValue = ""
	return newModel
}

// nextPlaceholder moves to the next placeholder
func (m Manager) nextPlaceholder() Manager {
	if len(m.placeholders) == 0 {
		return m
	}

	newModel := m

	// Save current edit if editing
	if newModel.isEditing && newModel.activeIndex >= 0 && newModel.activeIndex < len(newModel.placeholders) {
		ph := newModel.placeholders[newModel.activeIndex]
		if ph.Type == "text" {
			newModel.placeholders[newModel.activeIndex].CurrentValue = newModel.editValue
		}
		newModel.isEditing = false
	}

	newModel.activeIndex = (newModel.activeIndex + 1) % len(newModel.placeholders)

	// Update active states
	for i := range newModel.placeholders {
		newModel.placeholders[i].IsActive = (i == newModel.activeIndex)
	}

	// Initialize edit value for text placeholders and enter edit mode
	if newModel.activeIndex >= 0 {
		ph := newModel.placeholders[newModel.activeIndex]
		if ph.Type == "text" {
			newModel.editValue = ph.CurrentValue
			newModel.isEditing = true
		}
	}

	return newModel
}

// previousPlaceholder moves to the previous placeholder
func (m Manager) previousPlaceholder() Manager {
	if len(m.placeholders) == 0 {
		return m
	}

	newModel := m

	// Save current edit if editing
	if newModel.isEditing && newModel.activeIndex >= 0 && newModel.activeIndex < len(newModel.placeholders) {
		ph := newModel.placeholders[newModel.activeIndex]
		if ph.Type == "text" {
			newModel.placeholders[newModel.activeIndex].CurrentValue = newModel.editValue
		}
		newModel.isEditing = false
	}

	if newModel.activeIndex <= 0 {
		newModel.activeIndex = len(newModel.placeholders) - 1
	} else {
		newModel.activeIndex--
	}

	// Update active states
	for i := range newModel.placeholders {
		newModel.placeholders[i].IsActive = (i == newModel.activeIndex)
	}

	// Initialize edit value for text placeholders and enter edit mode
	if newModel.activeIndex >= 0 {
		ph := newModel.placeholders[newModel.activeIndex]
		if ph.Type == "text" {
			newModel.editValue = ph.CurrentValue
			newModel.isEditing = true
		}
	}

	return newModel
}

// Active returns the currently active placeholder
func (m Manager) Active() *Placeholder {
	if m.activeIndex < 0 || m.activeIndex >= len(m.placeholders) {
		return nil
	}
	return &m.placeholders[m.activeIndex]
}

// IsEditing returns true if editing a placeholder
func (m Manager) IsEditing() bool {
	return m.isEditing
}

// Placeholders returns all placeholders
func (m Manager) Placeholders() []Placeholder {
	return m.placeholders
}

// EditValue returns the current edit value
func (m Manager) EditValue() string {
	return m.editValue
}

// extractPlaceholders extracts all placeholders from content
func extractPlaceholders(content string) []Placeholder {
	// Regex to match {{type:name}} pattern - more permissive to catch invalid ones too
	re := regexp.MustCompile(`\{\{([^}:]+):([^}]+)\}\}`)
	matches := re.FindAllStringSubmatchIndex(content, -1)

	var placeholders []Placeholder
	for _, match := range matches {
		if len(match) >= 6 {
			fullMatchStart := match[0]
			fullMatchEnd := match[1]
			typeStart := match[2]
			typeEnd := match[3]
			nameStart := match[4]
			nameEnd := match[5]

			// Extract type and name
			placeholderType := content[typeStart:typeEnd]
			placeholderName := content[nameStart:nameEnd]

			placeholder := Placeholder{
				Type:     placeholderType,
				Name:     placeholderName,
				StartPos: fullMatchStart,
				EndPos:   fullMatchEnd,
				IsValid:  isValidPlaceholderType(placeholderType) && isValidPlaceholderName(placeholderName),
				IsActive: false,
			}

			// Initialize current value based on type
			if placeholderType == "list" {
				placeholder.ListValues = []string{}
			} else {
				placeholder.CurrentValue = ""
			}

			placeholders = append(placeholders, placeholder)
		}
	}

	return placeholders
}

// isValidPlaceholderType checks if placeholder type is valid
func isValidPlaceholderType(typ string) bool {
	return typ == "text" || typ == "list"
}

// isValidPlaceholderName checks if placeholder name is valid
func isValidPlaceholderName(name string) bool {
	if name == "" {
		return false
	}
	// Must be alphanumeric and underscores only
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return true
}

// ReplacePlaceholder replaces a placeholder with its filled value in content
func ReplacePlaceholder(content string, ph Placeholder) string {
	var replacement string
	if ph.Type == "list" {
		// Format as markdown bullet list
		if len(ph.ListValues) > 0 {
			items := make([]string, len(ph.ListValues))
			for i, val := range ph.ListValues {
				items[i] = "- " + val
			}
			replacement = strings.Join(items, "\n")
		} else {
			replacement = ""
		}
	} else {
		replacement = ph.CurrentValue
	}

	// Replace placeholder with value
	before := content[:ph.StartPos]
	after := content[ph.EndPos:]
	return before + replacement + after
}

// FindPlaceholderAtPosition finds the placeholder at the given cursor position
func FindPlaceholderAtPosition(placeholders []Placeholder, pos int) int {
	for i, ph := range placeholders {
		if pos >= ph.StartPos && pos <= ph.EndPos {
			return i
		}
	}
	return -1
}

// Messages

// ParsePlaceholdersMsg triggers placeholder parsing
type ParsePlaceholdersMsg struct {
	Content string
}

// ActivatePlaceholderMsg activates a placeholder by index
type ActivatePlaceholderMsg struct {
	Index int
}

// EditPlaceholderMsg updates the edit value
type EditPlaceholderMsg struct {
	Value string
}

// ExitEditModeMsg exits edit mode
type ExitEditModeMsg struct{}
