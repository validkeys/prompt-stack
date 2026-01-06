package diffviewer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyledavis/prompt-stack/internal/ai"
)

// ShowDiffMsg is sent to show the diff viewer modal
type ShowDiffMsg struct {
	Diff     *ai.UnifiedDiff
	Original string
	Edits    []ai.Edit
	OnAccept func() tea.Cmd
	OnReject func() tea.Cmd
}

// HideDiffMsg is sent to hide the diff viewer modal
type HideDiffMsg struct{}

// AcceptDiffMsg is sent when user accepts the diff
type AcceptDiffMsg struct{}

// RejectDiffMsg is sent when user rejects the diff
type RejectDiffMsg struct{}
