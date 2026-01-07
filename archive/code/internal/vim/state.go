package vim

import (
	"fmt"
	"strings"
)

// Mode represents the current vim editing mode
type Mode int

const (
	// NormalMode is for navigation and commands
	NormalMode Mode = iota
	// InsertMode is for text editing
	InsertMode
	// VisualMode is for text selection
	VisualMode
)

// String returns the string representation of the mode
func (m Mode) String() string {
	switch m {
	case NormalMode:
		return "NORMAL"
	case InsertMode:
		return "INSERT"
	case VisualMode:
		return "VISUAL"
	default:
		return "UNKNOWN"
	}
}

// IsEditing returns true if the mode allows text editing
func (m Mode) IsEditing() bool {
	return m == InsertMode
}

// IsNavigation returns true if the mode is for navigation
func (m Mode) IsNavigation() bool {
	return m == NormalMode || m == VisualMode
}

// State represents the current vim mode state
type State struct {
	// CurrentMode is the active vim mode
	CurrentMode Mode
	// PreviousMode stores the mode before entering a temporary mode (e.g., from Normal to Insert)
	PreviousMode Mode
	// transitionHook is called during mode transitions
	transitionHook TransitionHook
	// VisualStart is the starting position for visual selection (line, column)
	VisualStart struct {
		Line   int
		Column int
	}
	// VisualEnd is the ending position for visual selection (line, column)
	VisualEnd struct {
		Line   int
		Column int
	}
	// IsVisualBlock indicates if visual block mode is active (Ctrl+V)
	IsVisualBlock bool
	// IsVisualLine indicates if visual line mode is active (Shift+V)
	IsVisualLine bool
	// Register stores the yanked or deleted text
	Register string
	// RegisterType stores the type of content in the register (charwise, linewise, blockwise)
	RegisterType RegisterType
	// Count stores the pending count for commands (e.g., 5j means move down 5 lines)
	Count int
	// PendingCommand stores the command being built (e.g., 'd' followed by 'w' for delete word)
	PendingCommand string
	// LastCommand stores the last executed command for repeat with '.'
	LastCommand string
	// LastCount stores the count used with the last command
	LastCount int
	// SearchPattern stores the current search pattern
	SearchPattern string
	// SearchDirection stores the direction of the search (forward or backward)
	SearchDirection SearchDirection
	// MacroRecording stores the macro being recorded (key name)
	MacroRecording string
	// MacroContent stores the recorded macro commands
	MacroContent string
	// Marks stores named marks for positions in the text
	Marks map[string]Mark
}

// RegisterType represents the type of content in a register
type RegisterType int

const (
	// RegisterCharwise is for character-wise operations
	RegisterCharwise RegisterType = iota
	// RegisterLinewise is for line-wise operations
	RegisterLinewise
	// RegisterBlockwise is for block-wise operations
	RegisterBlockwise
)

// SearchDirection represents the direction of a search
type SearchDirection int

const (
	// SearchForward searches forward in the text
	SearchForward SearchDirection = iota
	// SearchBackward searches backward in the text
	SearchBackward
)

// Mark represents a named mark in the text
type Mark struct {
	Line   int
	Column int
}

// NewState creates a new vim state with default values
func NewState() *State {
	return &State{
		CurrentMode:     NormalMode,
		PreviousMode:    NormalMode,
		VisualStart:     struct{ Line, Column int }{Line: 0, Column: 0},
		VisualEnd:       struct{ Line, Column int }{Line: 0, Column: 0},
		IsVisualBlock:   false,
		IsVisualLine:    false,
		Register:        "",
		RegisterType:    RegisterCharwise,
		Count:           0,
		PendingCommand:  "",
		LastCommand:     "",
		LastCount:       0,
		SearchPattern:   "",
		SearchDirection: SearchForward,
		MacroRecording:  "",
		MacroContent:    "",
		Marks:           make(map[string]Mark),
	}
}

// TransitionHook is a function called during mode transitions
type TransitionHook func(from, to Mode)

// SetMode changes the current vim mode
func (s *State) SetMode(mode Mode) {
	if s.CurrentMode == mode {
		return
	}

	// Validate transition
	if !s.isValidTransition(s.CurrentMode, mode) {
		return
	}

	from := s.CurrentMode
	s.PreviousMode = from
	s.CurrentMode = mode

	// Call transition hook if set
	if s.transitionHook != nil {
		s.transitionHook(from, mode)
	}

	// Reset mode-specific state
	s.cleanupModeState(from)
	s.initializeModeState(mode)
}

// isValidTransition checks if a mode transition is valid
func (s *State) isValidTransition(from, to Mode) bool {
	// All transitions are valid in basic vim
	// This can be extended for more complex validation
	return true
}

// cleanupModeState cleans up state when leaving a mode
func (s *State) cleanupModeState(mode Mode) {
	switch mode {
	case InsertMode:
		// Leaving insert mode - no special cleanup needed
	case VisualMode:
		// Leaving visual mode - clear visual selection
		s.IsVisualBlock = false
		s.IsVisualLine = false
		s.VisualStart = struct{ Line, Column int }{Line: 0, Column: 0}
		s.VisualEnd = struct{ Line, Column int }{Line: 0, Column: 0}
	case NormalMode:
		// Leaving normal mode - no special cleanup needed
	}
}

// initializeModeState initializes state when entering a mode
func (s *State) initializeModeState(mode Mode) {
	switch mode {
	case NormalMode:
		// Reset count and pending command when entering normal mode
		s.Count = 0
		s.PendingCommand = ""
	case InsertMode:
		// No special initialization needed
	case VisualMode:
		// Visual mode state is set by EnterVisualMode methods
	}
}

// EnterInsertMode switches to insert mode
func (s *State) EnterInsertMode() {
	s.SetMode(InsertMode)
}

// EnterNormalMode switches to normal mode
func (s *State) EnterNormalMode() {
	s.SetMode(NormalMode)
}

// EnterVisualMode switches to visual mode
func (s *State) EnterVisualMode() {
	s.SetMode(VisualMode)
	s.IsVisualBlock = false
	s.IsVisualLine = false
}

// EnterVisualBlockMode switches to visual block mode
func (s *State) EnterVisualBlockMode() {
	s.SetMode(VisualMode)
	s.IsVisualBlock = true
	s.IsVisualLine = false
}

// EnterVisualLineMode switches to visual line mode
func (s *State) EnterVisualLineMode() {
	s.SetMode(VisualMode)
	s.IsVisualBlock = false
	s.IsVisualLine = true
}

// SetTransitionHook sets a function to be called during mode transitions
func (s *State) SetTransitionHook(hook TransitionHook) {
	s.transitionHook = hook
}

// GetTransitionHook returns the current transition hook
func (s *State) GetTransitionHook() TransitionHook {
	return s.transitionHook
}

// CanTransitionTo checks if a transition to the given mode is allowed
func (s *State) CanTransitionTo(mode Mode) bool {
	return s.isValidTransition(s.CurrentMode, mode)
}

// GetModeHistory returns the mode transition history
func (s *State) GetModeHistory() (current, previous Mode) {
	return s.CurrentMode, s.PreviousMode
}

// ToggleMode toggles between Normal and Insert modes
func (s *State) ToggleMode() {
	if s.CurrentMode == NormalMode {
		s.EnterInsertMode()
	} else if s.CurrentMode == InsertMode {
		s.EnterNormalMode()
	}
}

// ReturnToPreviousMode returns to the previous mode
func (s *State) ReturnToPreviousMode() {
	if s.PreviousMode != s.CurrentMode {
		s.SetMode(s.PreviousMode)
	}
}

// IsInVisualMode returns true if in any visual mode (visual, visual block, or visual line)
func (s *State) IsInVisualMode() bool {
	return s.CurrentMode == VisualMode
}

// IsInVisualBlockMode returns true if in visual block mode
func (s *State) IsInVisualBlockMode() bool {
	return s.CurrentMode == VisualMode && s.IsVisualBlock
}

// IsInVisualLineMode returns true if in visual line mode
func (s *State) IsInVisualLineMode() bool {
	return s.CurrentMode == VisualMode && s.IsVisualLine
}

// IsInNormalMode returns true if in normal mode
func (s *State) IsInNormalMode() bool {
	return s.CurrentMode == NormalMode
}

// IsInInsertMode returns true if in insert mode
func (s *State) IsInInsertMode() bool {
	return s.CurrentMode == InsertMode
}

// SetVisualStart sets the starting position for visual selection
func (s *State) SetVisualStart(line, column int) {
	s.VisualStart.Line = line
	s.VisualStart.Column = column
	s.VisualEnd.Line = line
	s.VisualEnd.Column = column
}

// SetVisualEnd sets the ending position for visual selection
func (s *State) SetVisualEnd(line, column int) {
	s.VisualEnd.Line = line
	s.VisualEnd.Column = column
}

// GetVisualRange returns the visual selection range
func (s *State) GetVisualRange() (startLine, startCol, endLine, endCol int) {
	// Ensure start is before end
	if s.VisualStart.Line < s.VisualEnd.Line ||
		(s.VisualStart.Line == s.VisualEnd.Line && s.VisualStart.Column < s.VisualEnd.Column) {
		return s.VisualStart.Line, s.VisualStart.Column, s.VisualEnd.Line, s.VisualEnd.Column
	}
	return s.VisualEnd.Line, s.VisualEnd.Column, s.VisualStart.Line, s.VisualStart.Column
}

// SetRegister stores text in the register
func (s *State) SetRegister(text string, regType RegisterType) {
	s.Register = text
	s.RegisterType = regType
}

// GetRegister retrieves text from the register
func (s *State) GetRegister() (string, RegisterType) {
	return s.Register, s.RegisterType
}

// SetCount sets the pending count for commands
func (s *State) SetCount(count int) {
	s.Count = count
}

// GetCount returns the pending count (defaults to 1 if not set)
func (s *State) GetCount() int {
	if s.Count == 0 {
		return 1
	}
	return s.Count
}

// ResetCount resets the pending count
func (s *State) ResetCount() {
	s.Count = 0
}

// SetPendingCommand sets the command being built
func (s *State) SetPendingCommand(cmd string) {
	s.PendingCommand = cmd
}

// GetPendingCommand returns the pending command
func (s *State) GetPendingCommand() string {
	return s.PendingCommand
}

// ResetPendingCommand resets the pending command
func (s *State) ResetPendingCommand() {
	s.PendingCommand = ""
}

// SetLastCommand stores the last executed command for repeat
func (s *State) SetLastCommand(cmd string, count int) {
	s.LastCommand = cmd
	s.LastCount = count
}

// GetLastCommand returns the last executed command
func (s *State) GetLastCommand() (string, int) {
	return s.LastCommand, s.LastCount
}

// SetSearchPattern sets the current search pattern
func (s *State) SetSearchPattern(pattern string, direction SearchDirection) {
	s.SearchPattern = pattern
	s.SearchDirection = direction
}

// GetSearchPattern returns the current search pattern
func (s *State) GetSearchPattern() (string, SearchDirection) {
	return s.SearchPattern, s.SearchDirection
}

// StartMacroRecording starts recording a macro
func (s *State) StartMacroRecording(key string) {
	s.MacroRecording = key
	s.MacroContent = ""
}

// StopMacroRecording stops recording a macro
func (s *State) StopMacroRecording() {
	s.MacroRecording = ""
}

// IsRecordingMacro returns true if a macro is being recorded
func (s *State) IsRecordingMacro() bool {
	return s.MacroRecording != ""
}

// AppendToMacro appends a command to the macro being recorded
func (s *State) AppendToMacro(cmd string) {
	if s.IsRecordingMacro() {
		s.MacroContent += cmd
	}
}

// GetMacroContent returns the content of the recorded macro
func (s *State) GetMacroContent() string {
	return s.MacroContent
}

// SetMark sets a named mark at a position
func (s *State) SetMark(name string, line, column int) {
	s.Marks[name] = Mark{
		Line:   line,
		Column: column,
	}
}

// GetMark retrieves a named mark
func (s *State) GetMark(name string) (Mark, bool) {
	mark, exists := s.Marks[name]
	return mark, exists
}

// DeleteMark removes a named mark
func (s *State) DeleteMark(name string) {
	delete(s.Marks, name)
}

// ClearMarks removes all marks
func (s *State) ClearMarks() {
	s.Marks = make(map[string]Mark)
}

// String returns a string representation of the vim state
func (s *State) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Mode: %s", s.CurrentMode))

	if s.CurrentMode == VisualMode {
		if s.IsVisualBlock {
			sb.WriteString(" (BLOCK)")
		} else if s.IsVisualLine {
			sb.WriteString(" (LINE)")
		}
	}

	if s.Count > 0 {
		sb.WriteString(fmt.Sprintf(" | Count: %d", s.Count))
	}

	if s.PendingCommand != "" {
		sb.WriteString(fmt.Sprintf(" | Pending: %s", s.PendingCommand))
	}

	if s.SearchPattern != "" {
		direction := "→"
		if s.SearchDirection == SearchBackward {
			direction = "←"
		}
		sb.WriteString(fmt.Sprintf(" | Search: %s %s", direction, s.SearchPattern))
	}

	if s.IsRecordingMacro() {
		sb.WriteString(fmt.Sprintf(" | Recording: @%s", s.MacroRecording))
	}

	return sb.String()
}

// Clone creates a deep copy of the vim state
func (s *State) Clone() *State {
	clone := NewState()
	clone.CurrentMode = s.CurrentMode
	clone.PreviousMode = s.PreviousMode
	clone.transitionHook = s.transitionHook
	clone.VisualStart = s.VisualStart
	clone.VisualEnd = s.VisualEnd
	clone.IsVisualBlock = s.IsVisualBlock
	clone.IsVisualLine = s.IsVisualLine
	clone.Register = s.Register
	clone.RegisterType = s.RegisterType
	clone.Count = s.Count
	clone.PendingCommand = s.PendingCommand
	clone.LastCommand = s.LastCommand
	clone.LastCount = s.LastCount
	clone.SearchPattern = s.SearchPattern
	clone.SearchDirection = s.SearchDirection
	clone.MacroRecording = s.MacroRecording
	clone.MacroContent = s.MacroContent

	// Deep copy marks
	clone.Marks = make(map[string]Mark)
	for k, v := range s.Marks {
		clone.Marks[k] = v
	}

	return clone
}
