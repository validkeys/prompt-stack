# AI Domain Key Learnings

**Purpose**: Key learnings and implementation patterns for AI integration from previous PromptStack implementation.

**Related Milestones**: M27, M28, M29, M30, M31, M32, M33

**Related Documents**: 
- [`project-structure.md`](../project-structure.md) - AI domain structure
- [`go-style-guide.md`](../go-style-guide.md) - Go coding standards
- [`go-testing-guide.md`](../go-testing-guide.md) - Testing patterns

---

## Learning Categories

### Category 1: AI Applying Indicator and Read-Only Mode

**Learning**: Use read-only mode to prevent concurrent edits during async operations

**Problem**: Need to prevent user from editing while AI is applying changes.

**Solution**: State-based UI feedback with editing restrictions

**Implementation Pattern**:
```go
// ui/workspace/model.go
type Model struct {
    content              string
    cursor               cursor
    // ... other fields
    isReadOnly           bool // true when AI is applying suggestion (blocks editing)
    aiApplying           bool // true when AI is actively applying a suggestion
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Block all editing when in read-only mode (AI applying suggestion)
        if m.isReadOnly {
            // Only allow cursor navigation in read-only mode
            switch msg.Type {
            case tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
                // Allow cursor navigation
            default:
                // Block all other keys
                return m, nil
            }
        }
        // ... rest of key handling
    }
    return m, nil
}

func (m *Model) SetAIApplying(applying bool) {
    m.aiApplying = applying
    // When AI is applying, also set read-only mode
    m.isReadOnly = applying
}

func (m Model) renderStatusBar() string {
    // Build status message
    var parts []string
    
    // AI applying indicator (highest priority)
    if m.aiApplying {
        parts = append(parts, "✨ AI is applying...")
    }
    
    // ... other status indicators
    
    return statusStyle.Render(statusText)
}
```

**Benefits**:
- Clear visual feedback when AI is applying changes
- Prevents user from editing while AI is modifying content
- Allows cursor navigation for viewing during application
- Simple state management with two boolean flags
- Automatic read-only mode activation when AI applies
- Status bar indicator shows highest priority message

**Lesson**: When implementing async operations that modify content, use read-only mode to prevent concurrent edits. Provide clear visual feedback in status bar. Allow cursor navigation so users can view changes while they're being applied. Use separate flags for state (aiApplying) and behavior (isReadOnly) to enable flexible control. This prevents race conditions and provides good user experience during async operations.

**Related Milestones**: M32, M33

**When to Apply**: When implementing async operations that modify content

---

### Category 2: Diff Viewer Modal Implementation

**Learning**: Use Bubble Tea's viewport component for scrollable content

**Problem**: Need to display AI-generated diffs with scrolling and navigation.

**Solution**: Viewport-based modal with unified diff display and color-coded changes

**Implementation Pattern**:
```go
// ui/diffviewer/model.go
type Model struct {
    viewport    viewport.Model
    diff        *ai.UnifiedDiff
    original    string
    edits       []ai.Edit
    width       int
    height      int
    onAccept    func() tea.Cmd
    onReject    func() tea.Cmd
    scrollOffset int
}

func (m Model) renderDiff() {
    var builder strings.Builder
    
    // Write header
    builder.WriteString(m.diff.Header)
    builder.WriteString("\n\n")
    
    // Write hunks
    for _, hunk := range m.diff.Hunks {
        // Write hunk header
        hunkHeader := fmt.Sprintf("@@ -%d,%d +%d,%d @@",
            hunk.OldStart, hunk.OldLines, hunk.NewStart, hunk.NewLines)
        builder.WriteString(theme.DiffHunkHeaderStyle().Render(hunkHeader))
        builder.WriteString("\n")
        
        // Write lines
        for _, line := range hunk.Lines {
            var styledLine string
            switch line.Type {
            case ai.DiffLineContext:
                styledLine = theme.DiffContextStyle().Render(" " + line.Content)
            case ai.DiffLineAddition:
                styledLine = theme.DiffAdditionStyle().Render("+" + line.Content)
            case ai.DiffLineDeletion:
                styledLine = theme.DiffDeletionStyle().Render("-" + line.Content)
            }
            builder.WriteString(styledLine)
            builder.WriteString("\n")
        }
        
        builder.WriteString("\n")
    }
    
    m.viewport.SetContent(builder.String())
    m.viewport.GotoTop()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter", "y":
            // Accept diff
            if m.onAccept != nil {
                return m, m.onAccept()
            }
        case "esc", "n", "q":
            // Reject diff
            if m.onReject != nil {
                return m, m.onReject()
            }
        case "up", "k":
            m.viewport.LineUp(1)
        case "down", "j":
            m.viewport.LineDown(1)
        case "pgup":
            m.viewport.HalfViewUp()
        case "pgdown":
            m.viewport.HalfViewDown()
        case "home", "g":
            m.viewport.GotoTop()
        case "end", "G":
            m.viewport.GotoBottom()
        }
    }
    return m, nil
}
```

**Benefits**:
- Viewport-based scrolling handles large diffs efficiently
- Color-coded changes (green for additions, red for deletions, cyan for hunk headers)
- Multiple accept/reject keybindings (Enter/y, Esc/n/q) for accessibility
- Statistics display (+X/-Y) shows change magnitude
- Help text in footer shows all available keybindings
- Message-based callbacks decouple UI from business logic
- Empty state handling when no diff is available
- Responsive to window size changes

**Lesson**: Use Bubble Tea's viewport component for scrollable content. Implement multiple keybindings for the same action (Enter/y for accept, Esc/n/q for reject) to accommodate different user preferences. Color-code diff lines by type (additions, deletions, context) for immediate visual recognition. Show statistics in header to provide context about change magnitude. Use message-based callbacks for accept/reject actions to decouple UI from business logic. This provides a clean, user-friendly diff review experience.

**Related Milestones**: M32, M33

**When to Apply**: When implementing diff viewers or scrollable content in TUI

---

### Category 3: AI Message-Based Workflow

**Learning**: Use custom message types for async AI operations

**Problem**: Need to handle async AI operations without blocking UI.

**Solution**: Custom message types for trigger, success, and error states

**Implementation Pattern**:
```go
// ui/app/model.go
type TriggerAISuggestionsMsg struct{}
type AISuggestionsGeneratedMsg struct {
    Suggestions []ai.Suggestion
    Error       error
}
type AISuggestionsErrorMsg struct {
    Error error
}

func (m Model) generateAISuggestions(composition string) tea.Cmd {
    return func() tea.Msg {
        // Extract keywords
        keywords := m.contextSelector.KeywordExtraction(composition)
        
        // Score and select prompts
        indexedPrompts := m.convertLibraryToIndexedPrompts()
        scoredPrompts := m.contextSelector.ScorePrompts(indexedPrompts, keywords)
        selectedPrompts := m.contextSelector.SelectTopPrompts(scoredPrompts, 5)
        
        // Build request
        request := ai.SuggestionRequest{
            Composition: composition,
            Library:     selectedPrompts,
        }
        
        // Send to AI API
        response, err := m.aiClient.SendMessage(request.GetSystemPrompt(), request)
        if err != nil {
            return AISuggestionsErrorMsg{Error: err}
        }
        
        // Parse response
        suggestionsResp, err := ai.ParseSuggestionsResponse(response.Content)
        if err != nil {
            return AISuggestionsErrorMsg{Error: err}
        }
        
        return AISuggestionsGeneratedMsg{Suggestions: suggestionsResp.Suggestions}
    }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case TriggerAISuggestionsMsg:
        if m.aiClient == nil {
            m.workspace.SetStatus("AI client not initialized. Please configure API key in settings.")
            return m, nil
        }
        
        // Check token budget
        tokenBudget := ai.NewTokenBudget(m.aiClient.GetModelContextLimit())
        _, atWarning, atBlock, tokens := tokenBudget.CheckComposition(m.composition)
        
        if atBlock {
            m.workspace.SetStatus(fmt.Sprintf("Composition exceeds token budget (%s). Please reduce content.", ai.FormatTokenCount(tokens)))
            return m, nil
        }
        
        if atWarning {
            m.workspace.SetStatus(fmt.Sprintf("Warning: Composition approaching token limit (%s)", ai.FormatTokenCount(tokens)))
        }
        
        // Trigger suggestion generation
        return m, m.generateAISuggestions(m.composition)
        
    case AISuggestionsGeneratedMsg:
        m.suggestions = msg.Suggestions
        m.activePanel = "suggestions"
        m.workspace.SetStatus(fmt.Sprintf("Generated %d AI suggestions", len(msg.Suggestions)))
        return m, nil
        
    case AISuggestionsErrorMsg:
        m.workspace.SetStatus(fmt.Sprintf("AI suggestions failed: %v", msg.Error))
        return m, nil
    }
    return m, nil
}
```

**Benefits**:
- Async operations don't block UI
- Clear separation between trigger, success, and error states
- Token budget checking before expensive API calls
- User feedback at each stage
- Error handling doesn't crash application

**Lesson**: Use custom message types for async operations. Define separate message types for trigger, success, and error states. Check prerequisites (like token budget) before starting expensive operations. Provide user feedback at each stage. This creates a smooth, responsive user experience even for slow operations.

**Related Milestones**: M27, M28, M29, M30, M31

**When to Apply**: When implementing async operations in Bubble Tea

---

### Category 4: Token Budget Enforcement

**Learning**: Implement conservative token budgeting for AI features

**Problem**: Need to prevent token waste and API cost overruns.

**Solution**: Conservative token allocation with proactive warnings

**Implementation Pattern**:
```go
// internal/ai/tokens.go
type TokenBudget struct {
    contextLimit      int
    compositionLimit  int // 25% of context
    libraryLimit     int // 15% of context
    warningThreshold int // 15% of context
    blockThreshold   int // 25% of context
}

func NewTokenBudget(contextLimit int) *TokenBudget {
    return &TokenBudget{
        contextLimit:      contextLimit,
        compositionLimit:  contextLimit / 4,  // 25%
        libraryLimit:     contextLimit * 15 / 100,  // 15%
        warningThreshold: contextLimit * 15 / 100,  // 15%
        blockThreshold:   contextLimit / 4,  // 25%
    }
}

func (tb *TokenBudget) CheckComposition(content string) (withinBudget, atWarning, atBlock bool, tokenCount int) {
    tokens := tb.EstimateTokensDetailed(content)
    
    atWarning = tokens >= tb.warningThreshold
    atBlock = tokens >= tb.blockThreshold
    withinBudget = tokens < tb.compositionLimit
    
    return withinBudget, atWarning, atBlock, tokens
}

func (tb *TokenBudget) EstimateTokensDetailed(content string) int {
    // Weighted estimation using multiple factors
    words := len(strings.Fields(content))
    chars := len(content)
    lines := strings.Count(content, "\n")
    
    // Weighted formula: (words * 1.3) + (chars / 4) + (lines * 0.5)
    return int(float64(words)*1.3 + float64(chars)/4.0 + float64(lines)*0.5)
}

func FormatTokenCount(tokens int) string {
    if tokens < 1000 {
        return fmt.Sprintf("%d tokens", tokens)
    }
    return fmt.Sprintf("%.1fK tokens", float64(tokens)/1000.0)
}
```

**Benefits**:
- Prevents token waste and API cost overruns
- Proactive warnings before blocking
- Conservative allocation ensures reliability
- Weighted token estimation more accurate than simple character count
- Clear status messages for users

**Lesson**: Implement conservative token budgeting for AI features. Use weighted estimation (words, characters, lines) rather than simple character counts. Provide warnings before blocking operations. Enforce strict limits to prevent cost overruns. This ensures reliable, cost-effective AI usage.

**Related Milestones**: M27, M28, M29, M30, M31

**When to Apply**: When implementing AI features with token limits

---

### Category 5: Context Selection Algorithm

**Learning**: Use multi-factor scoring for intelligent prompt selection

**Problem**: Need to select relevant prompts for AI context.

**Solution**: Multi-factor scoring with token budget constraints

**Implementation Pattern**:
```go
// internal/ai/context.go
type PromptScore struct {
    Prompt    IndexedPrompt
    Score      int
    Reasoning []string
}

func (cs *ContextSelector) ScorePrompts(prompts []IndexedPrompt, keywords map[string]int) []PromptScore {
    var scores []PromptScore
    
    for _, prompt := range prompts {
        score := 0
        var reasoning []string
        
        // Tag matches: +10 per matching tag
        tagMatches := cs.countTagMatches(prompt, keywords)
        score += tagMatches * 10
        reasoning = append(reasoning, fmt.Sprintf("Tag matches: %d (+%d)", tagMatches, tagMatches*10))
        
        // Category bonus: +5 if same category
        if cs.isSameCategory(prompt) {
            score += 5
            reasoning = append(reasoning, "Category match: +5")
        }
        
        // Keyword overlap: +1 per matching word (weighted by frequency)
        keywordScore := cs.calculateKeywordScore(prompt, keywords)
        score += keywordScore
        reasoning = append(reasoning, fmt.Sprintf("Keyword score: %d", keywordScore))
        
        // Recently used: +3 if used in last 24 hours
        if cs.isRecentlyUsed(prompt) {
            score += 3
            reasoning = append(reasoning, "Recently used: +3")
        }
        
        // Frequently used: +use_count
        score += prompt.UseCount
        reasoning = append(reasoning, fmt.Sprintf("Use count: +%d", prompt.UseCount))
        
        scores = append(scores, PromptScore{
            Prompt:    prompt,
            Score:      score,
            Reasoning: reasoning,
        })
    }
    
    // Sort by score descending
    cs.sortByScore(scores)
    return scores
}

func (cs *ContextSelector) SelectTopPrompts(scores []PromptScore, maxPrompts int) []IndexedPrompt {
    var selected []IndexedPrompt
    totalTokens := 0
    tokenBudget := ai.NewTokenBudget(cs.contextLimit)
    
    for _, scored := range scores {
        promptTokens := cs.estimateTokens(scored.Prompt.Content)
        
        // Check if adding this prompt would exceed budget
        if totalTokens + promptTokens > tokenBudget.GetLibraryLimit() {
            break
        }
        
        selected = append(selected, scored.Prompt)
        totalTokens += promptTokens
        
        if len(selected) >= maxPrompts {
            break
        }
    }
    
    return selected
}
```

**Benefits**:
- Intelligent prompt selection based on multiple signals
- Token budget constraints prevent context overflow
- Transparent scoring with reasoning array
- Usage patterns (recent, frequent) improve relevance
- Tag and category matching provide strong signals

**Lesson**: Use multi-factor scoring for AI context selection. Combine strong signals (tags, category) with content-based signals (keywords). Incorporate usage patterns (recent, frequent) for personalization. Enforce token budget constraints to prevent context overflow. Document scoring reasoning for transparency and debugging. This provides high-quality, relevant AI suggestions.

**Related Milestones**: M28, M29, M30, M31

**When to Apply**: When implementing AI context selection or recommendation systems

---

### Category 6: Command Palette Integration for AI Features

**Learning**: Command handler returns success, TUI handles execution

**Problem**: Need to integrate AI features into command palette.

**Solution**: Command handler returns success, TUI handles actual execution

**Implementation Pattern**:
```go
// internal/commands/core.go
func RegisterCoreCommands(registry *Registry) error {
    if err := registry.Register(&Command{
        ID:          "ai-suggestions",
        Name:        "Get AI Suggestions",
        Description: "Request AI suggestions for current composition",
        Category:    "AI",
        Handler: func() error {
            // This command will be handled by TUI to trigger AI suggestions
            return nil
        },
    }); err != nil {
        return err
    }
    return nil
}

// ui/app/model.go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case ExecuteCommandMsg:
        // Execute command from palette
        cmd := m.registry.GetCommand(msg.CommandID)
        if cmd != nil {
            err := cmd.Handler()
            if err != nil {
                return m, SetErrorMessage(err.Error())
            }
            
            // Handle AI suggestions command
            if cmd.ID == "ai-suggestions" {
                return m, func() tea.Msg {
                    return TriggerAISuggestionsMsg{}
                }
            }
        }
    }
    return m, nil
}
```

**Benefits**:
- Command palette shows all available commands
- AI features discoverable through command palette
- Command handler returns success, TUI handles actual execution
- Clean separation between command registry and TUI logic
- Consistent with other command implementations

**Lesson**: For commands that trigger TUI-specific operations, have the handler return success and let the TUI handle the actual execution via message passing. This keeps the command registry simple and TUI-specific logic in the TUI layer. Use message-based execution to maintain clean separation of concerns.

**Related Milestones**: M27, M28, M29, M30, M31

**When to Apply**: When integrating TUI-specific features into command palette

---

### Category 7: Read-Only Mode During Async Operations

**Learning**: Use read-only mode to prevent concurrent edits during async operations

**Problem**: Need to prevent user from editing while AI is applying changes.

**Solution**: State-based UI feedback with editing restrictions

**Implementation Pattern**:
```go
// ui/workspace/model.go
type Model struct {
    content              string
    cursor               cursor
    // ... other fields
    isReadOnly           bool // true when AI is applying suggestion (blocks editing)
    aiApplying           bool // true when AI is actively applying a suggestion
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Block all editing when in read-only mode (AI applying suggestion)
        if m.isReadOnly {
            // Only allow cursor navigation in read-only mode
            switch msg.Type {
            case tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
                // Allow cursor navigation
            default:
                // Block all other keys
                return m, nil
            }
        }
        // ... rest of key handling
    }
    return m, nil
}

func (m *Model) SetAIApplying(applying bool) {
    m.aiApplying = applying
    // When AI is applying, also set read-only mode
    m.isReadOnly = applying
}

func (m Model) renderStatusBar() string {
    // Build status message
    var parts []string
    
    // AI applying indicator (highest priority)
    if m.aiApplying {
        parts = append(parts, "✨ AI is applying...")
    }
    
    // ... other status indicators
    
    return statusStyle.Render(statusText)
}
```

**Benefits**:
- Clear visual feedback when AI is applying changes
- Prevents user from editing while AI is modifying content
- Allows cursor navigation for viewing during application
- Simple state management with two boolean flags
- Automatic read-only mode activation when AI applies
- Status bar indicator shows highest priority message

**Lesson**: When implementing async operations that modify content, use read-only mode to prevent concurrent edits. Provide clear visual feedback in status bar. Allow cursor navigation so users can view changes while they're being applied. Use separate flags for state (aiApplying) and behavior (isReadOnly) to enable flexible control. This prevents race conditions and provides good user experience during async operations.

**Related Milestones**: M32, M33

**When to Apply**: When implementing async operations that modify content

---

## Quick Reference

| Learning | Milestone | Priority |
|----------|-----------|----------|
| AI Applying Indicator and Read-Only Mode | M32, M33 | High |
| Diff Viewer Modal Implementation | M32, M33 | High |
| AI Message-Based Workflow | M27, M28, M29, M30, M31 | High |
| Token Budget Enforcement | M27, M28, M29, M30, M31 | High |
| Context Selection Algorithm | M28, M29, M30, M31 | High |
| Command Palette Integration for AI Features | M27, M28, M29, M30, M31 | High |
| Read-Only Mode During Async Operations | M32, M33 | High |

---

**Last Updated**: 2026-01-07
**Source**: [`key-learnings.md`](../key-learnings.md)