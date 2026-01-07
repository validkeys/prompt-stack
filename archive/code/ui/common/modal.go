package common

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/ui/theme"
)

// Modal represents a generic modal dialog
type Modal struct {
	title        string
	content      string
	width        int
	height       int
	showButtons  bool
	primaryBtn   string
	secondaryBtn string
	focused      bool
}

// Styles for modal
var (
	modalStyle = theme.ModalStyle()

	titleStyle = theme.ModalTitleStyle()

	contentStyle = theme.ModalContentStyle()

	buttonStyle = theme.ModalButtonStyle()

	buttonFocusedStyle = theme.ModalButtonFocusedStyle()

	buttonSecondaryStyle = theme.ModalButtonSecondaryStyle()

	buttonSecondaryFocusedStyle = theme.ModalButtonSecondaryFocusedStyle()
)

// NewModal creates a new modal
func NewModal(title, content string) Modal {
	return Modal{
		title:       title,
		content:     content,
		width:       64,
		height:      20,
		showButtons: false,
		focused:     true,
	}
}

// WithButtons adds buttons to the modal
func (m Modal) WithButtons(primary, secondary string) Modal {
	m.showButtons = true
	m.primaryBtn = primary
	m.secondaryBtn = secondary
	return m
}

// WithSize sets the modal size
func (m Modal) WithSize(width, height int) Modal {
	m.width = width
	m.height = height
	return m
}

// Init initializes the modal
func (m Modal) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Modal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.showButtons {
			switch msg.String() {
			case "esc", "enter", "q":
				return m, tea.Quit
			}
		} else {
			switch msg.String() {
			case "esc":
				return m, func() tea.Msg { return CloseModalMsg{} }
			case "enter":
				return m, func() tea.Msg { return ModalActionMsg{Action: "primary"} }
			case "tab":
				// Toggle button focus (for future implementation)
			}
		}
	}
	return m, nil
}

// View renders the modal
func (m Modal) View() string {
	// Build modal content
	var content strings.Builder

	// Title
	content.WriteString(titleStyle.Render(m.title))
	content.WriteString("\n\n")

	// Content (wrap text)
	wrappedContent := contentStyle.Width(m.width - 4).Render(m.content)
	content.WriteString(wrappedContent)

	// Buttons
	if m.showButtons {
		content.WriteString("\n\n")
		buttons := m.renderButtons()
		content.WriteString(buttons)
	}

	// Apply modal style
	modalContent := modalStyle.Render(content.String())

	// Center modal in viewport
	return lipgloss.Place(
		m.height,
		m.width,
		lipgloss.Center,
		lipgloss.Center,
		modalContent,
	)
}

// renderButtons renders the modal buttons
func (m Modal) renderButtons() string {
	var buttons strings.Builder

	if m.secondaryBtn != "" {
		buttons.WriteString(buttonSecondaryStyle.Render(m.secondaryBtn))
	}

	if m.primaryBtn != "" {
		buttons.WriteString(buttonStyle.Render(m.primaryBtn))
	}

	return buttons.String()
}

// Message types for modal

// CloseModalMsg signals to close the modal
type CloseModalMsg struct{}

// ModalActionMsg signals a button action
type ModalActionMsg struct {
	Action string
}

// ErrorModal creates an error modal
func ErrorModal(title, message string) Modal {
	return NewModal(title, message).WithButtons("OK", "")
}

// WarningModal creates a warning modal
func WarningModal(title, message string) Modal {
	return NewModal(title, message).WithButtons("OK", "")
}

// ConfirmModal creates a confirmation modal
func ConfirmModal(title, message string) Modal {
	return NewModal(title, message).WithButtons("Confirm", "Cancel")
}

// InfoModal creates an info modal
func InfoModal(title, message string) Modal {
	return NewModal(title, message).WithButtons("OK", "")
}
