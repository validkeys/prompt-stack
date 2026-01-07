package statusbar

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/ui/theme"
)

// Model represents the status bar
type Model struct {
	width          int
	height         int
	message        string
	messageType    MessageType
	messageTimeout time.Time
	charCount      int
	lineCount      int
	tokenEstimate  int
	vimMode        string
	editMode       string
	showAutoSave   bool
	autoSaveStatus string
}

// MessageType represents the type of status message
type MessageType string

const (
	MessageTypeNormal  MessageType = "normal"
	MessageTypeInfo    MessageType = "info"
	MessageTypeSuccess MessageType = "success"
	MessageTypeWarning MessageType = "warning"
	MessageTypeError   MessageType = "error"
)

// Styles for the status bar
var (
	statusStyle = theme.StatusStyle()

	infoStyle = theme.InfoStyle()

	successStyle = theme.SuccessStyle()

	warningStyle = theme.WarningStyle()

	errorStyle = theme.ErrorStyle()

	separatorStyle = theme.SeparatorStyle()
)

// New creates a new status bar model
func New() Model {
	return Model{
		messageType: MessageTypeNormal,
	}
}

// Init initializes the status bar
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case SetMessageMsg:
		m.message = msg.Message
		m.messageType = msg.Type
		m.messageTimeout = time.Now().Add(msg.Timeout)
	case ClearMessageMsg:
		m.message = ""
		m.messageType = MessageTypeNormal
	case SetStatsMsg:
		m.charCount = msg.CharCount
		m.lineCount = msg.LineCount
	case SetTokenEstimateMsg:
		m.tokenEstimate = msg.TokenEstimate
	case SetVimModeMsg:
		m.vimMode = msg.Mode
	case SetEditModeMsg:
		m.editMode = msg.Mode
	case SetAutoSaveStatusMsg:
		m.showAutoSave = msg.Show
		m.autoSaveStatus = msg.Status
	}

	// Auto-clear message after timeout
	if !m.messageTimeout.IsZero() && time.Now().After(m.messageTimeout) {
		m.message = ""
		m.messageType = MessageTypeNormal
		m.messageTimeout = time.Time{}
	}

	return m, nil
}

// View renders the status bar
func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	// Left side: message or auto-save status
	leftContent := ""
	if m.showAutoSave && m.autoSaveStatus != "" {
		leftContent = m.autoSaveStatus
	} else if m.message != "" {
		leftContent = m.formatMessage()
	}

	// Right side: stats and modes
	rightContent := m.formatRightContent()

	// Combine left and right with proper spacing
	if leftContent != "" && rightContent != "" {
		availableWidth := m.width - lipgloss.Width(leftContent) - lipgloss.Width(rightContent) - 2
		if availableWidth > 0 {
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				leftContent,
				strings.Repeat(" ", availableWidth),
				rightContent,
			)
		}
	}

	// If only one side has content
	if leftContent != "" {
		return leftContent
	}
	if rightContent != "" {
		return lipgloss.PlaceHorizontal(m.width, lipgloss.Right, rightContent)
	}

	return statusStyle.Render("")
}

// formatMessage formats the status message with appropriate styling
func (m Model) formatMessage() string {
	switch m.messageType {
	case MessageTypeInfo:
		return infoStyle.Render(m.message)
	case MessageTypeSuccess:
		return successStyle.Render(m.message)
	case MessageTypeWarning:
		return warningStyle.Render(m.message)
	case MessageTypeError:
		return errorStyle.Render(m.message)
	default:
		return statusStyle.Render(m.message)
	}
}

// formatRightContent formats the right side of the status bar
func (m Model) formatRightContent() string {
	var parts []string

	// Token estimate (when AI panel is open)
	if m.tokenEstimate > 0 {
		parts = append(parts, fmt.Sprintf("~%dK tokens", m.tokenEstimate/1000))
	}

	// Edit mode indicator
	if m.editMode != "" {
		parts = append(parts, m.editMode)
	}

	// Vim mode indicator
	if m.vimMode != "" {
		parts = append(parts, m.vimMode)
	}

	// Line and character counts
	if m.lineCount > 0 || m.charCount > 0 {
		parts = append(parts, fmt.Sprintf("%d lines, %d chars", m.lineCount, m.charCount))
	}

	// Join with separators
	if len(parts) == 0 {
		return ""
	}

	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += separatorStyle.Render("•") + parts[i]
	}

	return statusStyle.Render(result)
}

// Message types for tea.Msg

// SetMessageMsg sets a status message
type SetMessageMsg struct {
	Message string
	Type    MessageType
	Timeout time.Duration
}

// ClearMessageMsg clears the current status message
type ClearMessageMsg struct{}

// SetStatsMsg sets character and line counts
type SetStatsMsg struct {
	CharCount int
	LineCount int
}

// SetTokenEstimateMsg sets the token estimate
type SetTokenEstimateMsg struct {
	TokenEstimate int
}

// SetVimModeMsg sets the vim mode indicator
type SetVimModeMsg struct {
	Mode string
}

// SetEditModeMsg sets the edit mode indicator
type SetEditModeMsg struct {
	Mode string
}

// SetAutoSaveStatusMsg sets the auto-save status
type SetAutoSaveStatusMsg struct {
	Show   bool
	Status string
}

// Helper functions to create messages

// SetInfoMessage creates an info message
func SetInfoMessage(message string) tea.Cmd {
	return func() tea.Msg {
		return SetMessageMsg{
			Message: message,
			Type:    MessageTypeInfo,
			Timeout: 3 * time.Second,
		}
	}
}

// SetSuccessMessage creates a success message
func SetSuccessMessage(message string) tea.Cmd {
	return func() tea.Msg {
		return SetMessageMsg{
			Message: message,
			Type:    MessageTypeSuccess,
			Timeout: 2 * time.Second,
		}
	}
}

// SetWarningMessage creates a warning message
func SetWarningMessage(message string) tea.Cmd {
	return func() tea.Msg {
		return SetMessageMsg{
			Message: message,
			Type:    MessageTypeWarning,
			Timeout: 5 * time.Second,
		}
	}
}

// SetErrorMessage creates an error message
func SetErrorMessage(message string) tea.Cmd {
	return func() tea.Msg {
		return SetMessageMsg{
			Message: message,
			Type:    MessageTypeError,
			Timeout: 5 * time.Second,
		}
	}
}

// SetPersistentErrorMessage creates a persistent error message (no timeout)
func SetPersistentErrorMessage(message string) tea.Cmd {
	return func() tea.Msg {
		return SetMessageMsg{
			Message: message,
			Type:    MessageTypeError,
			Timeout: 0, // No timeout
		}
	}
}

// SetStats creates a stats update message
func SetStats(charCount, lineCount int) tea.Cmd {
	return func() tea.Msg {
		return SetStatsMsg{
			CharCount: charCount,
			LineCount: lineCount,
		}
	}
}

// SetTokenEstimate creates a token estimate update message
func SetTokenEstimate(tokenEstimate int) tea.Cmd {
	return func() tea.Msg {
		return SetTokenEstimateMsg{
			TokenEstimate: tokenEstimate,
		}
	}
}

// SetVimMode creates a vim mode update message
func SetVimMode(mode string) tea.Cmd {
	return func() tea.Msg {
		return SetVimModeMsg{
			Mode: mode,
		}
	}
}

// SetEditMode creates an edit mode update message
func SetEditMode(mode string) tea.Cmd {
	return func() tea.Msg {
		return SetEditModeMsg{
			Mode: mode,
		}
	}
}

// SetAutoSaveStatus creates an auto-save status update message
func SetAutoSaveStatus(show bool, status string) tea.Cmd {
	return func() tea.Msg {
		return SetAutoSaveStatusMsg{
			Show:   show,
			Status: status,
		}
	}
}

// ClearMessage creates a clear message command
func ClearMessage() tea.Cmd {
	return func() tea.Msg {
		return ClearMessageMsg{}
	}
}
