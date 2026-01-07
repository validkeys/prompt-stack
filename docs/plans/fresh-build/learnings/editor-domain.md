# Editor Domain Key Learnings

**Purpose**: Key learnings and implementation patterns for editor functionality from previous PromptStack implementation.

**Related Milestones**: M4, M5, M6, M11, M12, M13, M14

**Related Documents**: 
- [`project-structure.md`](../project-structure.md) - Editor domain structure
- [`go-style-guide.md`](../go-style-guide.md) - Go coding standards
- [`go-testing-guide.md`](../go-testing-guide.md) - Testing patterns

---

## Learning Categories

### Category 1: Placeholder Parsing

**Learning**: Regex with position tracking for structured text parsing

**Problem**: Need to parse placeholders from markdown content and track their positions for navigation and highlighting.

**Solution**: Use regex with position tracking to extract both content and positions

**Implementation Pattern**:
```go
import "regexp"

type Placeholder struct {
    Type         string   // "text" or "list"
    Name         string   // placeholder name
    StartPos     int      // position in content
    EndPos       int      // position in content
    CurrentValue string   // current filled value (for text)
    ListValues   []string // current filled values (for list)
    IsValid      bool     // whether syntax is valid
    IsActive     bool     // whether currently selected
}

func ParsePlaceholders(content string) []Placeholder {
    re := regexp.MustCompile(`\{\{(\w+):(\w+)\}\}`)
    matches := re.FindAllStringSubmatchIndex(content, -1)
    
    var placeholders []Placeholder
    for _, match := range matches {
        // match[0], match[1] = full match start/end
        // match[2], match[3] = type group start/end
        // match[4], match[5] = name group start/end
        
        phType := content[match[2]:match[3]]
        name := content[match[4]:match[5]]
        
        placeholders = append(placeholders, Placeholder{
            Type:     phType,
            Name:     name,
            StartPos: match[0],
            EndPos:   match[1],
            IsValid:  isValidPlaceholderType(phType) && isValidPlaceholderName(name),
        })
    }
    
    return placeholders
}
```

**Lesson**: When parsing structured text, track both the content and its positions. This enables features like highlighting and navigation.

**Related Milestones**: M4, M5

**When to Apply**: When parsing structured text with placeholders or similar patterns

---

### Category 2: Index Scoring Algorithm

**Learning**: Multi-factor scoring for relevance ranking

**Problem**: Need to rank prompts by relevance when searching or providing suggestions.

**Solution**: Combine multiple signals (tags, keywords, usage patterns) for scoring

**Implementation Pattern**:
```go
type PromptScore struct {
    Prompt    *Prompt
    Score     int
    Reasoning []string
}

func ScorePrompts(prompts []*Prompt, keywords map[string]int) []PromptScore {
    var scores []PromptScore
    
    for _, prompt := range prompts {
        score := 0
        var reasoning []string
        
        // Tag matches: +10 per matching tag
        tagMatches := countTagMatches(prompt, keywords)
        score += tagMatches * 10
        reasoning = append(reasoning, fmt.Sprintf("Tag matches: %d (+%d)", tagMatches, tagMatches*10))
        
        // Keyword overlap: +1 per matching word (weighted by frequency)
        keywordScore := calculateKeywordScore(prompt, keywords)
        score += keywordScore
        reasoning = append(reasoning, fmt.Sprintf("Keyword score: %d", keywordScore))
        
        // Recently used: +3 if used in last 24 hours
        if isRecentlyUsed(prompt) {
            score += 3
            reasoning = append(reasoning, "Recently used: +3")
        }
        
        // Frequently used: +use_count
        score += prompt.UseCount
        reasoning = append(reasoning, fmt.Sprintf("Use count: +%d", prompt.UseCount))
        
        scores = append(scores, PromptScore{
            Prompt:    prompt,
            Score:     score,
            Reasoning: reasoning,
        })
    }
    
    // Sort by score descending
    sort.Slice(scores, func(i, j int) bool {
        return scores[i].Score > scores[j].Score
    })
    
    return scores
}
```

**Rationale**:
- Tags are strong signals of relevance
- Keywords provide content-based matching
- Usage patterns reflect user preferences
- Time decay ensures fresh content

**Lesson**: Relevance scoring should combine multiple signals. No single factor is sufficient for good recommendations.

**Related Milestones**: M4, M5

**When to Apply**: When implementing search, ranking, or recommendation systems

---

### Category 3: Validation Strategy

**Learning**: Separate errors and warnings with severity levels

**Problem**: Need to validate prompts and provide feedback without blocking all functionality.

**Solution**: Separate validation errors (block insertion) from warnings (allow with indicator)

**Implementation Pattern**:
```go
type ValidationError struct {
    Type    string // "error" or "warning"
    Message string // human-readable message
    Line    int    // line number
    Column  int    // column number
}

type ValidationResult struct {
    Errors   []ValidationError  // Block insertion
    Warnings []ValidationError  // Allow with indicator
    IsValid  bool
}

func ValidatePlaceholders(placeholders []Placeholder) []ValidationError {
    var errors []ValidationError
    nameMap := make(map[string]int)
    
    // Check for duplicate names
    for i, ph := range placeholders {
        if !ph.IsValid {
            continue
        }
        if prevIndex, exists := nameMap[ph.Name]; exists {
            errors = append(errors, ValidationError{
                Type:    "error",
                Message: "Duplicate placeholder name: " + ph.Name,
                Line:    getLineNumber(placeholders, i),
                Column:  ph.StartPos,
            })
        } else {
            nameMap[ph.Name] = i
        }
    }
    
    return errors
}
```

**Benefits**:
- Graceful degradation
- User can still use prompts with warnings
- Clear distinction between critical and minor issues

**Lesson**: Validation should have severity levels. Not all issues should block functionality. Provide users with information and let them decide.

**Related Milestones**: M4, M5

**When to Apply**: When implementing validation for user input or data

---

### Category 4: Placeholder System Implementation

**Learning**: Re-parse placeholders on every content change to keep state synchronized

**Problem**: Need to maintain placeholder state as user edits content.

**Solution**: Re-parse placeholders on every content change

**Implementation Pattern**:
```go
type Model struct {
    content           string
    placeholders      []editor.Placeholder
    activePlaceholder int // -1 if none active
    // ... other fields
}

func (m *Model) updatePlaceholders() {
    m.placeholders = editor.ParsePlaceholders(m.content)
    
    // Maintain active placeholder if still valid
    if m.activePlaceholder >= 0 && m.activePlaceholder < len(m.placeholders) {
        activeName := m.placeholders[m.activePlaceholder].Name
        found := false
        for i, ph := range m.placeholders {
            if ph.Name == activeName {
                m.activePlaceholder = i
                found = true
                break
            }
        }
        if !found {
            m.activePlaceholder = -1
        }
    } else {
        m.activePlaceholder = -1
    }
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle key input
        m.insertRune(msg.Runes)
        m.updatePlaceholders() // Re-parse after content change
    }
    return m, nil
}
```

**Benefits**:
- Automatic placeholder detection on content changes
- State stays synchronized with content
- Active placeholder tracking maintained across edits

**Lesson**: Re-parse placeholders on every content change to keep state synchronized. This ensures placeholders are always up-to-date with the current content.

**Related Milestones**: M4, M5

**When to Apply**: When implementing systems that parse structured content

---

### Category 5: Placeholder Validation Strategy

**Learning**: Validate for both syntax and semantic correctness

**Problem**: Need to ensure placeholders are valid before they're used.

**Solution**: Check for duplicate names, validate types and names

**Implementation Pattern**:
```go
func ValidatePlaceholders(placeholders []Placeholder) []ValidationError {
    var errors []ValidationError
    nameMap := make(map[string]int)
    
    // Check for duplicate names
    for i, ph := range placeholders {
        if !ph.IsValid {
            continue
        }
        if prevIndex, exists := nameMap[ph.Name]; exists {
            errors = append(errors, ValidationError{
                Type:    "error",
                Message: "Duplicate placeholder name: " + ph.Name,
                Line:    getLineNumber(placeholders, i),
                Column:  ph.StartPos,
            })
            errors = append(errors, ValidationError{
                Type:    "error",
                Message: "Duplicate placeholder name: " + ph.Name,
                Line:    getLineNumber(placeholders, prevIndex),
                Column:  placeholders[prevIndex].StartPos,
            })
        } else {
            nameMap[ph.Name] = i
        }
    }
    
    // Validate each placeholder
    for i, ph := range placeholders {
        if !ph.IsValid {
            if !isValidPlaceholderType(ph.Type) {
                errors = append(errors, ValidationError{
                    Type:    "error",
                    Message: "Invalid placeholder type: " + ph.Type + " (must be 'text' or 'list')",
                    Line:    getLineNumber(placeholders, i),
                    Column:  ph.StartPos,
                })
            }
            if !isValidPlaceholderName(ph.Name) {
                errors = append(errors, ValidationError{
                    Type:    "error",
                    Message: "Invalid placeholder name: " + ph.Name + " (must be alphanumeric and underscores only)",
                    Line:    getLineNumber(placeholders, i),
                    Column:  ph.StartPos,
                })
            }
        }
    }
    
    return errors
}
```

**Benefits**:
- Duplicate detection prevents confusion
- Type validation ensures only supported types are used
- Name validation prevents syntax errors
- Line and column information for precise error reporting

**Lesson**: Validate placeholders for both syntax and semantic correctness. Check for duplicate names to prevent runtime confusion. Separate errors (block insertion) from warnings (allow with indicator).

**Related Milestones**: M4, M5

**When to Apply**: When implementing validation for structured data

---

### Category 6: Cursor Position Management for Placeholders

**Learning**: Implement bidirectional conversion between cursor coordinates and absolute positions

**Problem**: Need to navigate to specific placeholder positions in content.

**Solution**: Convert between cursor coordinates (x, y) and absolute positions

**Implementation Pattern**:
```go
func (m *Model) getCursorPosition() int {
    lines := strings.Split(m.content, "\n")
    pos := 0
    
    for i := 0; i < m.cursor.y && i < len(lines); i++ {
        pos += len(lines[i]) + 1 // +1 for newline
    }
    
    if m.cursor.y < len(lines) {
        pos += m.cursor.x
    }
    
    return pos
}

func (m *Model) setCursorToPosition(pos int) {
    lines := strings.Split(m.content, "\n")
    currentPos := 0
    
    for i, line := range lines {
        lineEnd := currentPos + len(line)
        
        if pos <= lineEnd {
            m.cursor.y = i
            m.cursor.x = pos - currentPos
            return
        }
        
        currentPos = lineEnd + 1 // +1 for newline
    }
    
    // If position is beyond content, set to end
    m.cursor.y = len(lines) - 1
    m.cursor.x = len(lines[len(lines)-1])
}

func (m *Model) navigateToNextPlaceholder() bool {
    cursorPos := m.getCursorPosition()
    nextIndex := editor.GetNextPlaceholder(m.placeholders, cursorPos)
    if nextIndex >= 0 {
        m.activePlaceholder = nextIndex
        ph := m.placeholders[nextIndex]
        m.setCursorToPosition(ph.StartPos)
        return true
    }
    return false
}
```

**Benefits**:
- Accurate navigation to placeholder positions
- Handles multi-line content correctly
- Edge case handling (position beyond content)
- Enables precise placeholder selection

**Lesson**: Implement bidirectional conversion between cursor coordinates (x, y) and absolute positions. This is essential for features like placeholder navigation where you need to jump to specific positions.

**Related Milestones**: M4, M5

**When to Apply**: When implementing text editing with navigation features

---

### Category 7: Placeholder Highlighting in TUI

**Learning**: Render placeholders line-by-line with position-based highlighting

**Problem**: Need to visually highlight active placeholder in TUI.

**Solution**: Calculate line positions and apply styles to active placeholder

**Implementation Pattern**:
```go
func (m *Model) renderLineWithPlaceholders(line string, lineIndex int) string {
    if len(m.placeholders) == 0 {
        return line
    }
    
    // Calculate line start position in content
    lineStartPos := 0
    lines := strings.Split(m.content, "\n")
    for i := 0; i < lineIndex && i < len(lines); i++ {
        lineStartPos += len(lines[i]) + 1 // +1 for newline
    }
    
    // Find placeholders on this line
    result := line
    offset := 0
    
    for _, ph := range m.placeholders {
        // Check if placeholder is on this line
        if ph.StartPos >= lineStartPos && ph.EndPos <= lineStartPos+len(line) {
            // Calculate position within line
            phStart := ph.StartPos - lineStartPos
            phEnd := ph.EndPos - lineStartPos
            
            // Apply highlighting if active
            if m.activePlaceholder >= 0 && m.placeholders[m.activePlaceholder].Name == ph.Name {
                placeholderStyle := theme.ActivePlaceholderStyle()
                placeholderText := line[phStart+offset : phEnd+offset]
                result = result[:phStart+offset] + placeholderStyle.Render(placeholderText) + result[phEnd+offset:]
                offset += len(placeholderStyle.Render(placeholderText)) - (phEnd - phStart)
            }
        }
    }
    
    return result
}
```

**Benefits**:
- Visual feedback for active placeholder
- Line-by-line rendering for performance
- Offset tracking for styled text length
- Only highlights active placeholder, not all placeholders

**Lesson**: Render placeholders line-by-line with position-based highlighting. Calculate line start position in content to determine which placeholders are on each line. Track offset when applying styles because styled text has different length than plain text.

**Related Milestones**: M4, M5

**When to Apply**: When implementing visual highlighting in TUI applications

---

### Category 8: Theme Integration for Placeholders

**Learning**: Create dedicated style functions for specific UI elements

**Problem**: Need consistent styling for placeholder highlighting across the application.

**Solution**: Centralized theme package with style helper functions

**Implementation Pattern**:
```go
// ui/theme/theme.go
func ActivePlaceholderStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Background(lipgloss.Color(AccentYellow)).
        Foreground(lipgloss.Color(BackgroundPrimary)).
        Bold(true)
}

// Usage in workspace
import "github.com/kyledavis/prompt-stack/ui/theme"

func (m *Model) renderLineWithPlaceholders(line string, lineIndex int) string {
    // ...
    if m.activePlaceholder >= 0 && m.placeholders[m.activePlaceholder].Name == ph.Name {
        placeholderStyle := theme.ActivePlaceholderStyle()
        placeholderText := line[phStart+offset : phEnd+offset]
        result = result[:phStart+offset] + placeholderStyle.Render(placeholderText) + result[phEnd+offset:]
    }
    // ...
}
```

**Benefits**:
- Consistent styling across application
- Easy to update placeholder appearance
- Uses existing color palette
- High contrast for visibility

**Lesson**: Create dedicated style functions for specific UI elements like active placeholders. Use high-contrast colors (yellow background with dark foreground) for visibility. Keep all styles in the centralized theme package for consistency.

**Related Milestones**: M4, M5

**When to Apply**: When implementing visual styling in TUI applications

---

### Category 9: Text Placeholder Editing Mode

**Learning**: Implement placeholder editing as a separate mode with its own key handling

**Problem**: Need to allow users to edit placeholder values directly.

**Solution**: Separate edit mode with state management and value replacement

**Implementation Pattern**:
```go
type Model struct {
    content              string
    cursor               cursor
    placeholders         []editor.Placeholder
    activePlaceholder     int // -1 if none active
    placeholderEditMode   bool // true when editing a placeholder
    placeholderEditValue  string // current value being edited
}

func (m *Model) enterPlaceholderEditMode() {
    if m.activePlaceholder < 0 || m.activePlaceholder >= len(m.placeholders) {
        return
    }
    
    ph := &m.placeholders[m.activePlaceholder]
    
    // Only text placeholders can be edited in this mode
    if ph.Type != "text" {
        return
    }
    
    // Initialize edit value with current value or empty string
    m.placeholderEditMode = true
    m.placeholderEditValue = ph.CurrentValue
    
    // Move cursor to placeholder position
    m.setCursorToPosition(ph.StartPos)
}

func (m Model) handlePlaceholderEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.Type {
    case tea.KeyEsc:
        // Exit placeholder edit mode and save the value
        m.exitPlaceholderEditMode()
        return m, nil
        
    case tea.KeyBackspace:
        // Delete character from edit value
        if len(m.placeholderEditValue) > 0 {
            m.placeholderEditValue = m.placeholderEditValue[:len(m.placeholderEditValue)-1]
        }
        return m, nil
        
    case tea.KeyEnter:
        // Exit placeholder edit mode and save the value
        m.exitPlaceholderEditMode()
        return m, nil
        
    case tea.KeyRunes:
        // Append characters to edit value
        m.placeholderEditValue += string(msg.Runes)
        return m, nil
    }
    
    return m, nil
}

func (m *Model) exitPlaceholderEditMode() {
    if m.activePlaceholder < 0 || m.activePlaceholder >= len(m.placeholders) {
        m.placeholderEditMode = false
        m.placeholderEditValue = ""
        return
    }
    
    ph := &m.placeholders[m.activePlaceholder]
    
    // Update placeholder's current value
    ph.CurrentValue = m.placeholderEditValue
    
    // Replace placeholder in content with filled value
    m.content = editor.ReplacePlaceholder(m.content, *ph)
    
    // Re-parse placeholders after replacement
    m.updatePlaceholders()
    
    // Mark as dirty and schedule auto-save
    m.markDirty()
    m.scheduleAutoSave()
    
    // Exit edit mode
    m.placeholderEditMode = false
    m.placeholderEditValue = ""
}
```

**Benefits**:
- Vim-style editing workflow familiar to developers
- Clear visual feedback with status bar indicator
- Type to replace placeholder content directly
- Esc or Enter to save and exit edit mode
- Automatic placeholder replacement with value
- Content re-parsing after editing
- Auto-save triggered after placeholder is filled

**Lesson**: Implement placeholder editing as a separate mode with its own key handling. Use a boolean flag to track edit mode state. Store edit value separately from placeholder until exit. Show edit value instead of placeholder syntax during editing.

**Related Milestones**: M4, M5

**When to Apply**: When implementing text editing with special modes

---

### Category 10: Text Editor Cursor Positioning

**Learning**: Handle cursor movement across line boundaries

**Problem**: Need natural cursor movement when editing text.

**Solution**: Handle edge cases when moving cursor across lines

**Implementation Pattern**:
```go
func (m *Model) moveCursorLeft() {
    if m.cursor.x > 0 {
        m.cursor.x--
    } else if m.cursor.y > 0 {
        // Move to end of previous line
        m.cursor.y--
        lines := strings.Split(m.content, "\n")
        if m.cursor.y < len(lines) {
            lineLen := len(lines[m.cursor.y])
            m.cursor.x = lineLen
        }
    }
}

func (m *Model) moveCursorRight() {
    lines := strings.Split(m.content, "\n")
    if m.cursor.y < len(lines) {
        lineLen := len(lines[m.cursor.y])
        if m.cursor.x < lineLen {
            m.cursor.x++
        } else if m.cursor.y < len(lines)-1 {
            // Move to start of next line
            m.cursor.y++
            m.cursor.x = 0
        }
    }
}
```

**Lesson**: Always handle edge cases when moving cursor. When moving left at column 0, move to end of previous line. When moving right at end of line, move to start of next line. This provides natural text editing behavior.

**Related Milestones**: M4, M5

**When to Apply**: When implementing text editing with cursor movement

---

### Category 11: File Path Management for History

**Learning**: Use timestamp-based file naming with directory creation

**Problem**: Need to save history files without conflicts.

**Solution**: Timestamp-based naming with automatic directory creation

**Implementation Pattern**:
```go
import (
    "os"
    "path/filepath"
    "time"
)

func (m *Model) saveToFile() error {
    if m.filePath == "" {
        timestamp := time.Now().Format("2006-01-02_15-04-05")
        m.filePath = filepath.Join(m.workingDir, ".promptstack", ".history", timestamp+".md")
    }
    
    dir := filepath.Dir(m.filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }
    
    return os.WriteFile(m.filePath, []byte(m.content), 0644)
}
```

**Lesson**: Use timestamp-based naming for history files to avoid conflicts. Always create directories with MkdirAll before writing files. Use filepath.Join for cross-platform path construction.

**Related Milestones**: M15, M16

**When to Apply**: When saving files with unique names

---

### Category 12: Lipgloss Styling

**Learning**: Define reusable styles and compose them

**Problem**: Need consistent styling across TUI components.

**Solution**: Define styles as reusable variables or functions

**Implementation Pattern**:
```go
import "github.com/charmbracelet/lipgloss"

editorStyle := lipgloss.NewStyle().
    Width(m.width).
    Height(availableHeight).
    Padding(0, 1)

statusStyle := lipgloss.NewStyle().
    Width(m.width).
    Height(1).
    Background(lipgloss.Color("240")).
    Foreground(lipgloss.Color("15")).
    Padding(0, 1)

cursorStyle := lipgloss.NewStyle().
    Background(lipgloss.Color("7")).
    Foreground(lipgloss.Color("0"))
```

**Lesson**: Define styles as reusable variables or functions. Use Lipgloss's fluent API for clean style definitions. Use color codes (240 for gray, 7 for white, etc.) for consistent theming.

**Related Milestones**: M4, M5, M6

**When to Apply**: When styling TUI components with Lipgloss

---

## Quick Reference

| Learning | Milestone | Priority |
|----------|-----------|----------|
| Placeholder Parsing | M4, M5 | High |
| Index Scoring Algorithm | M4, M5 | High |
| Validation Strategy | M4, M5 | High |
| Placeholder System Implementation | M4, M5 | High |
| Placeholder Validation Strategy | M4, M5 | High |
| Cursor Position Management for Placeholders | M4, M5 | High |
| Placeholder Highlighting in TUI | M4, M5 | High |
| Theme Integration for Placeholders | M4, M5 | High |
| Text Placeholder Editing Mode | M4, M5 | High |
| Text Editor Cursor Positioning | M4, M5 | High |
| File Path Management for History | M15, M16 | Medium |
| Lipgloss Styling | M4, M5, M6 | High |

---

**Last Updated**: 2026-01-07
**Source**: [`key-learnings.md`](../key-learnings.md)