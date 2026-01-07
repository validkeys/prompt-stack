package app

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyledavis/prompt-stack/internal/ai"
	"github.com/kyledavis/prompt-stack/internal/commands"
	"github.com/kyledavis/prompt-stack/internal/history"
	"github.com/kyledavis/prompt-stack/internal/library"
	"github.com/kyledavis/prompt-stack/internal/prompt"
	"github.com/kyledavis/prompt-stack/ui/browser"
	"github.com/kyledavis/prompt-stack/ui/diffviewer"
	historyui "github.com/kyledavis/prompt-stack/ui/history"
	"github.com/kyledavis/prompt-stack/ui/palette"
	"github.com/kyledavis/prompt-stack/ui/suggestions"
	"github.com/kyledavis/prompt-stack/ui/workspace"
)

// Model represents the root application model
type Model struct {
	workspace       workspace.Model
	suggestions     suggestions.Model
	diffViewer      diffviewer.Model
	browser         browser.Model
	palette         palette.Model
	history         historyui.Model
	activePanel     string // "workspace", "suggestions", "diffviewer", "browser", "palette", "history"
	width           int
	height          int
	pendingEdits    []ai.Edit // Store edits for diff viewer
	pendingOriginal string    // Store original content for diff viewer
	library         *library.Library
	commands        *commands.Registry
	historyManager  *history.Manager
	vimMode         bool
	aiClient        *ai.Client
	contextSelector *ai.ContextSelector
}

// New creates a new application model
func New(workingDir string) Model {
	return Model{
		workspace:      workspace.New(workingDir),
		suggestions:    suggestions.NewModel(),
		diffViewer:     diffviewer.NewModel(),
		browser:        browser.New(make(map[string]*prompt.Prompt), false),
		palette:        palette.New(commands.NewRegistry(), false),
		history:        historyui.NewModel(),
		activePanel:    "workspace",
		width:          80,
		height:         24,
		library:        nil,
		commands:       nil,
		historyManager: nil,
		vimMode:        false,
	}
}

// NewWithDependencies creates a new application model with dependencies
func NewWithDependencies(workingDir string, lib *library.Library, registry *commands.Registry, historyMgr *history.Manager, vimMode bool, aiClient *ai.Client, contextSelector *ai.ContextSelector) Model {
	return Model{
		workspace:       workspace.New(workingDir),
		suggestions:     suggestions.NewModel(),
		diffViewer:      diffviewer.NewModel(),
		browser:         browser.New(lib.Prompts, vimMode),
		palette:         palette.New(registry, vimMode),
		history:         historyui.NewModel(),
		activePanel:     "workspace",
		width:           80,
		height:          24,
		library:         lib,
		commands:        registry,
		historyManager:  historyMgr,
		vimMode:         vimMode,
		aiClient:        aiClient,
		contextSelector: contextSelector,
	}
}

// Init initializes the application model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keybindings
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+p":
			// Show command palette
			m.palette.Show()
			m.activePanel = "palette"
			return m, nil
		case "ctrl+b":
			// Show library browser
			m.browser.Show()
			m.activePanel = "browser"
			return m, nil
		case "ctrl+h":
			// Show history browser
			m.showHistoryBrowser()
			m.activePanel = "history"
			return m, nil
		default:
			// Let the key fall through to the active panel
			// This allows normal typing and other keys to be handled by the workspace
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.workspace.SetSize(msg.Width, msg.Height)
		m.suggestions.SetSize(msg.Width, msg.Height)
		m.diffViewer.SetSize(msg.Width, msg.Height)
		m.history.SetSize(msg.Width, msg.Height)
		// Forward size updates to palette and browser
		updatedPalette, _ := m.palette.Update(msg)
		m.palette = updatedPalette.(palette.Model)
		updatedBrowser, _ := m.browser.Update(msg)
		m.browser = updatedBrowser.(browser.Model)

	case diffviewer.ShowDiffMsg:
		// Show diff viewer modal
		m.activePanel = "diffviewer"
		m.diffViewer.SetDiff(msg.Diff)
		m.pendingOriginal = msg.Original
		m.pendingEdits = msg.Edits

		m.diffViewer.SetOnAccept(func() tea.Cmd {
			return func() tea.Msg {
				return diffviewer.AcceptDiffMsg{}
			}
		})
		m.diffViewer.SetOnReject(func() tea.Cmd {
			return func() tea.Msg {
				return diffviewer.RejectDiffMsg{}
			}
		})

	case diffviewer.HideDiffMsg:
		// Hide diff viewer modal
		m.activePanel = "workspace"

	case diffviewer.AcceptDiffMsg:
		// User accepted the diff - apply edits
		if m.pendingOriginal != "" && len(m.pendingEdits) > 0 {
			err := m.workspace.ApplyEditsAsSingleAction(m.pendingEdits)
			if err != nil {
				// Handle error - show in status bar
				m.workspace.SetStatus(fmt.Sprintf("Failed to apply edits: %v", err))
			}
		}
		// Unlock editor
		m.workspace.SetAIApplying(false)
		m.workspace.SetReadOnly(false)
		// Hide diff viewer
		m.activePanel = "workspace"
		// Clear pending data
		m.pendingEdits = nil
		m.pendingOriginal = ""

	case diffviewer.RejectDiffMsg:
		// User rejected the diff - discard changes
		// Unlock editor
		m.workspace.SetAIApplying(false)
		m.workspace.SetReadOnly(false)
		// Hide diff viewer
		m.activePanel = "workspace"
		// Clear pending data
		m.pendingEdits = nil
		m.pendingOriginal = ""
	}

	// Update active panel
	switch m.activePanel {
	case "workspace":
		updatedModel, updateCmd := m.workspace.Update(msg)
		m.workspace = updatedModel.(workspace.Model)
		cmd = updateCmd
	case "suggestions":
		m.suggestions, cmd = m.suggestions.Update(msg)
	case "diffviewer":
		m.diffViewer, cmd = m.diffViewer.Update(msg)
	case "browser":
		updatedBrowser, browserCmd := m.browser.Update(msg)
		m.browser = updatedBrowser.(browser.Model)
		cmd = browserCmd
	case "palette":
		updatedPalette, paletteCmd := m.palette.Update(msg)
		m.palette = updatedPalette.(palette.Model)
		cmd = paletteCmd
	case "history":
		updatedHistory, historyCmd := m.history.Update(msg)
		m.history = updatedHistory.(historyui.Model)
		cmd = historyCmd
	}

	// Handle browser messages
	switch msg := msg.(type) {
	case browser.InsertPromptMsg:
		// Insert prompt into workspace
		if m.library != nil {
			if p, exists := m.library.Prompts[msg.FilePath]; exists {
				// Insert prompt content at cursor position
				m.insertPromptAtCursor(p.Content, msg.InsertMode)
				m.workspace.SetStatus(fmt.Sprintf("Inserted: %s", p.Title))
				m.browser.Hide()
				m.activePanel = "workspace"
			}
		}
	case browser.ValidationErrorMsg:
		// Show validation error
		m.workspace.SetStatus(fmt.Sprintf("Validation error: %s", msg.Errors[0].Message))
		m.browser.Hide()
		m.activePanel = "workspace"
	}

	// Handle history messages
	switch msg := msg.(type) {
	case LoadHistoryMsg:
		// Load history composition into workspace
		if m.historyManager != nil {
			content, err := m.historyManager.LoadComposition(msg.FilePath)
			if err != nil {
				m.workspace.SetStatus(fmt.Sprintf("Failed to load: %v", err))
			} else {
				m.workspace.SetContent(content)
				m.workspace.SetStatus(fmt.Sprintf("Loaded: %s", msg.FilePath))
			}
		}
		m.activePanel = "workspace"
	case DeleteHistoryMsg:
		// Delete history composition
		if m.historyManager != nil {
			err := m.historyManager.DeleteComposition(msg.FilePath)
			if err != nil {
				m.workspace.SetStatus(fmt.Sprintf("Failed to delete: %v", err))
			} else {
				m.workspace.SetStatus(fmt.Sprintf("Deleted: %s", msg.FilePath))
				// Refresh history browser
				m.showHistoryBrowser()
			}
		}
	}

	// Handle command palette execution messages
	switch msg := msg.(type) {
	case palette.ExecuteSuccessMsg:
		// Check if this is the AI suggestions command
		if msg.CommandID == "ai-suggestions" {
			// Trigger AI suggestions
			return m, func() tea.Msg {
				return TriggerAISuggestionsMsg{}
			}
		}
	case palette.ExecuteErrorMsg:
		// Show error message
		m.workspace.SetStatus(fmt.Sprintf("Command failed: %v", msg.Error))
	}

	// Handle AI suggestion messages
	switch msg := msg.(type) {
	case TriggerAISuggestionsMsg:
		// Check if AI client is available
		if m.aiClient == nil {
			m.workspace.SetStatus("AI client not initialized. Please configure API key in settings.")
			return m, nil
		}

		// Check token budget
		composition := m.workspace.GetContent()
		tokenBudget := ai.NewTokenBudget(m.aiClient.GetModelContextLimit())
		_, atWarning, atBlock, tokens := tokenBudget.CheckComposition(composition)

		if atBlock {
			m.workspace.SetStatus(fmt.Sprintf("Composition exceeds token budget (%s). Please reduce content.", ai.FormatTokenCount(tokens)))
			return m, nil
		}

		if atWarning {
			m.workspace.SetStatus(fmt.Sprintf("Warning: Composition approaching token budget (%s)", ai.FormatTokenCount(tokens)))
		}

		// Start suggestion generation in background
		return m, m.generateAISuggestions(composition)

	case AISuggestionsGeneratedMsg:
		// Handle generated suggestions
		if msg.Error != nil {
			m.workspace.SetStatus(fmt.Sprintf("Failed to generate suggestions: %v", msg.Error))
			return m, nil
		}

		// Set suggestions in panel
		m.SetSuggestions(msg.Suggestions)
		m.ShowSuggestions()
		m.workspace.SetStatus(fmt.Sprintf("Generated %d suggestion(s)", len(msg.Suggestions)))
		return m, nil

	case AISuggestionsErrorMsg:
		// Handle AI error
		m.workspace.SetStatus(fmt.Sprintf("AI error: %v", msg.Error))
		return m, nil
	}

	return m, cmd
}

// View renders the application
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	switch m.activePanel {
	case "workspace":
		return m.renderWorkspace()
	case "suggestions":
		return m.renderSuggestions()
	case "diffviewer":
		return m.renderDiffViewer()
	case "browser":
		return m.renderBrowser()
	case "palette":
		return m.renderPalette()
	case "history":
		return m.renderHistory()
	default:
		return m.renderWorkspace()
	}
}

// renderWorkspace renders the workspace view
func (m Model) renderWorkspace() string {
	// Render workspace
	workspaceView := m.workspace.View()

	// If suggestions panel is active, render split view
	if m.activePanel == "workspace" && len(m.suggestions.GetSuggestions()) > 0 {
		// Calculate split (70% workspace, 30% suggestions)
		workspaceWidth := m.width * 7 / 10
		suggestionsWidth := m.width - workspaceWidth

		// Style workspace
		workspaceStyle := lipgloss.NewStyle().
			Width(workspaceWidth).
			Height(m.height)

		// Style suggestions
		suggestionsStyle := lipgloss.NewStyle().
			Width(suggestionsWidth).
			Height(m.height)

		// Render suggestions panel
		m.suggestions.SetSize(suggestionsWidth, m.height)
		suggestionsView := suggestionsStyle.Render(m.suggestions.View())

		// Combine views
		return lipgloss.JoinHorizontal(lipgloss.Top,
			workspaceStyle.Render(workspaceView),
			suggestionsView,
		)
	}

	return workspaceView
}

// renderSuggestions renders the suggestions panel view
func (m Model) renderSuggestions() string {
	return m.suggestions.View()
}

// renderDiffViewer renders the diff viewer modal
func (m Model) renderDiffViewer() string {
	// Render diff viewer as modal overlay
	diffView := m.diffViewer.View()

	// Create modal style
	modalStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	return modalStyle.Render(diffView)
}

// renderBrowser renders the library browser modal
func (m Model) renderBrowser() string {
	// Render browser as modal overlay
	browserView := m.browser.View()

	// Create modal style
	modalStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	return modalStyle.Render(browserView)
}

// renderPalette renders the command palette modal
func (m Model) renderPalette() string {
	// Render palette as modal overlay
	paletteView := m.palette.View()

	// Create modal style
	modalStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	return modalStyle.Render(paletteView)
}

// SetSuggestions sets the suggestions to display
func (m *Model) SetSuggestions(suggestionList []ai.Suggestion) {
	m.suggestions.SetSuggestions(suggestionList)

	// Set up callbacks for applying/dismissing suggestions
	m.suggestions.SetOnApply(func(suggestion *ai.Suggestion) tea.Cmd {
		return func() tea.Msg {
			// When user presses 'a' to apply suggestion:
			// 1. Lock editor (read-only mode)
			m.workspace.SetAIApplying(true)
			m.workspace.SetReadOnly(true)

			// 2. Generate diff from suggestion edits
			diff, err := ai.NewDiffGenerator().GenerateUnifiedDiff(m.workspace.GetContent(), suggestion.ProposedChanges)
			if err != nil {
				// Handle error - show in status bar and unlock editor
				m.workspace.SetStatus(fmt.Sprintf("Failed to generate diff: %v", err))
				m.workspace.SetAIApplying(false)
				m.workspace.SetReadOnly(false)
				return nil
			}

			// 3. Store pending data
			m.pendingOriginal = m.workspace.GetContent()
			m.pendingEdits = suggestion.ProposedChanges

			// 4. Show diff viewer modal
			return diffviewer.ShowDiffMsg{
				Diff:     diff,
				Original: m.workspace.GetContent(),
				Edits:    suggestion.ProposedChanges,
			}
		}
	})

	m.suggestions.SetOnDismiss(func(suggestion *ai.Suggestion) tea.Cmd {
		return func() tea.Msg {
			// Mark suggestion as dismissed
			suggestion.MarkAsDismissed()
			return nil
		}
	})
}

// GetSuggestions returns the current suggestions
func (m *Model) GetSuggestions() []ai.Suggestion {
	return m.suggestions.GetSuggestions()
}

// ShowSuggestions shows the suggestions panel
func (m *Model) ShowSuggestions() {
	m.activePanel = "workspace" // Show split view
}

// HideSuggestions hides the suggestions panel
func (m *Model) HideSuggestions() {
	m.activePanel = "workspace"
}

// showHistoryBrowser shows the history browser with current history items
func (m *Model) showHistoryBrowser() {
	if m.historyManager == nil {
		m.workspace.SetStatus("History manager not initialized")
		return
	}

	// Get all compositions from history
	compositions, err := m.historyManager.GetAllCompositions()
	if err != nil {
		m.workspace.SetStatus(fmt.Sprintf("Failed to load history: %v", err))
		return
	}

	// Convert to history items
	var items []list.Item
	for _, comp := range compositions {
		// Create preview (first 100 chars)
		preview := comp.Content
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}

		item := historyui.NewItem(
			comp.FilePath,
			comp.CreatedAt,
			comp.WorkingDirectory,
			preview,
			comp.CharacterCount,
			comp.LineCount,
		)
		items = append(items, item)
	}

	// Set items in history browser
	m.history.SetItems(items)

	// Set up callbacks
	m.history.SetOnSelect(func(filePath string) tea.Cmd {
		return func() tea.Msg {
			return LoadHistoryMsg{FilePath: filePath}
		}
	})

	m.history.SetOnDelete(func(filePath string) tea.Cmd {
		return func() tea.Msg {
			return DeleteHistoryMsg{FilePath: filePath}
		}
	})
}

// insertPromptAtCursor inserts prompt content at cursor position
func (m *Model) insertPromptAtCursor(content string, mode browser.InsertMode) {
	// Get current cursor position
	cursorPos := m.workspace.GetCursorPosition()

	// Insert content at cursor position
	m.workspace.InsertContent(cursorPos, content)

	// If inserting on new line, add newline before
	if mode == browser.InsertOnNewLine {
		m.workspace.InsertContent(cursorPos, "\n")
	}

	// Mark as dirty and schedule auto-save
	m.workspace.MarkDirty()
}

// renderHistory renders the history browser modal
func (m Model) renderHistory() string {
	// Render history as modal overlay
	historyView := m.history.View()

	// Create modal style
	modalStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	return modalStyle.Render(historyView)
}

// generateAISuggestions generates AI suggestions for the current composition
func (m Model) generateAISuggestions(composition string) tea.Cmd {
	return func() tea.Msg {
		// Extract keywords from composition
		keywords := m.contextSelector.KeywordExtraction(composition)

		// Convert library prompts to indexed prompts
		var indexedPrompts []*ai.IndexedPrompt
		for _, p := range m.library.Prompts {
			indexedPrompts = append(indexedPrompts, &ai.IndexedPrompt{
				PromptID:      p.ID,
				Title:         p.Title,
				Description:   p.Description,
				Tags:          p.Tags,
				Category:      p.Category,
				WordFrequency: make(map[string]int), // TODO: Build from content
				LastUsed:      time.Time{},          // TODO: Track usage
				UseCount:      0,                    // TODO: Track usage
				Content:       p.Content,
			})
		}

		// Score prompts based on composition
		scored := m.contextSelector.ScorePrompts(
			indexedPrompts,
			keywords,
			[]string{}, // No tags from composition yet
			"",         // No category from composition yet
		)

		// Select top prompts within token budget
		tokenBudget := ai.NewTokenBudget(m.aiClient.GetModelContextLimit())
		selectedPrompts := m.contextSelector.SelectTopPrompts(scored, tokenBudget.GetLibraryLimit())

		// Build context for AI request
		var contextPrompts []ai.IndexedPrompt
		for _, p := range selectedPrompts {
			contextPrompts = append(contextPrompts, *p)
		}

		// Build messages for AI
		messages := []ai.Message{
			{Role: "user", Content: fmt.Sprintf("Composition:\n%s\n\nContext Prompts:\n%v", composition, contextPrompts)},
		}

		// Send request to AI
		response, err := m.aiClient.SendMessage(context.Background(), ai.MessageRequest{
			SystemPrompt: ai.GetSystemPrompt(),
			Messages:     messages,
			MaxTokens:    4000,
			Temperature:  0.7,
		})

		if err != nil {
			return AISuggestionsErrorMsg{Error: err}
		}

		// Parse suggestions from response
		suggestionsResp, err := ai.ParseSuggestionsResponse(response.Content)
		if err != nil {
			return AISuggestionsErrorMsg{Error: fmt.Errorf("failed to parse suggestions: %w", err)}
		}

		return AISuggestionsGeneratedMsg{Suggestions: suggestionsResp.Suggestions}
	}
}

// Messages for history
type LoadHistoryMsg struct {
	FilePath string
}

type DeleteHistoryMsg struct {
	FilePath string
}

// Messages for AI suggestions
type TriggerAISuggestionsMsg struct{}

type AISuggestionsGeneratedMsg struct {
	Suggestions []ai.Suggestion
	Error       error
}

type AISuggestionsErrorMsg struct {
	Error error
}
