package editor

import (
	"time"

	"github.com/kyledavis/prompt-stack/internal/ai"
)

// ActionType represents the type of undo action
type ActionType string

const (
	ActionTypeInsert          ActionType = "insert"
	ActionTypeDelete          ActionType = "delete"
	ActionTypePaste           ActionType = "paste"
	ActionTypePromptInsert    ActionType = "prompt_insert"
	ActionTypePlaceholderFill ActionType = "placeholder_fill"
	ActionTypeNewline         ActionType = "newline"
	ActionTypeBackspace       ActionType = "backspace"
	ActionTypeBatchEdit       ActionType = "batch_edit"
)

// UndoAction represents a single undoable action
type UndoAction struct {
	Type       ActionType
	Content    string    // The content that was inserted or deleted
	Position   int       // Absolute position in the content
	CursorPos  cursorPos // Cursor position after the action
	Timestamp  time.Time // When the action was created
	BatchID    string    // ID for batching related actions
	IsBatchEnd bool      // Marks the end of a batch
	Edits      []ai.Edit // Edits for batch edit actions (for redo)
}

// cursorPos represents cursor position
type cursorPos struct {
	X int
	Y int
}

// UndoStack manages the undo and redo history
type UndoStack struct {
	undoActions    []UndoAction
	redoActions    []UndoAction
	maxSize        int // Maximum number of actions to keep (default: 100)
	currentBatch   string
	lastActionTime time.Time
}

// NewUndoStack creates a new undo stack with default settings
func NewUndoStack() *UndoStack {
	return &UndoStack{
		undoActions: make([]UndoAction, 0),
		redoActions: make([]UndoAction, 0),
		maxSize:     100,
	}
}

// PushAction adds a new action to the undo stack
func (s *UndoStack) PushAction(action UndoAction) {
	// Clear redo stack when new action is added
	s.redoActions = make([]UndoAction, 0)

	// Check if we should batch this action with the previous one
	if s.shouldBatch(action) {
		// Add to current batch
		action.BatchID = s.currentBatch
	} else {
		// Start a new batch
		s.currentBatch = generateBatchID()
		action.BatchID = s.currentBatch
		action.IsBatchEnd = true
	}

	// Add to undo stack
	s.undoActions = append(s.undoActions, action)

	// Enforce size limit
	if len(s.undoActions) > s.maxSize {
		// Remove oldest batch
		s.removeOldestBatch()
	}

	// Update last action time
	s.lastActionTime = time.Now()
}

// shouldBatch determines if an action should be batched with the previous one
func (s *UndoStack) shouldBatch(action UndoAction) bool {
	if len(s.undoActions) == 0 {
		return false
	}

	lastAction := s.undoActions[len(s.undoActions)-1]

	// Different action types cannot be batched
	if lastAction.Type != action.Type {
		return false
	}

	// Only certain action types can be batched
	switch action.Type {
	case ActionTypeInsert:
		// Batch continuous typing (break on >1 second pause)
		if time.Since(s.lastActionTime) > time.Second {
			return false
		}
		// Batch if position is contiguous
		return lastAction.Position+len(lastAction.Content) == action.Position

	case ActionTypeBackspace:
		// Batch continuous backspace (break on >1 second pause)
		if time.Since(s.lastActionTime) > time.Second {
			return false
		}
		// Batch if position is contiguous
		return lastAction.Position == action.Position

	default:
		// Other action types are not batched
		return false
	}
}

// Undo performs an undo operation
func (s *UndoStack) Undo() (UndoAction, bool) {
	if len(s.undoActions) == 0 {
		return UndoAction{}, false
	}

	// Find the end of the current batch
	batchEnd := len(s.undoActions) - 1
	batchStart := batchEnd

	// Find the start of the batch
	for i := batchEnd; i >= 0; i-- {
		if s.undoActions[i].IsBatchEnd {
			batchStart = i
			break
		}
	}

	// Extract the batch
	batch := s.undoActions[batchStart : batchEnd+1]

	// Remove from undo stack
	s.undoActions = s.undoActions[:batchStart]

	// Add to redo stack (in reverse order)
	for i := len(batch) - 1; i >= 0; i-- {
		s.redoActions = append(s.redoActions, batch[i])
	}

	// Return the first action of the batch (for cursor positioning)
	return batch[0], true
}

// Redo performs a redo operation
func (s *UndoStack) Redo() (UndoAction, bool) {
	if len(s.redoActions) == 0 {
		return UndoAction{}, false
	}

	// Find the end of the current batch in redo stack
	batchEnd := len(s.redoActions) - 1
	batchStart := batchEnd

	// Find the start of the batch
	for i := batchEnd; i >= 0; i-- {
		if s.redoActions[i].IsBatchEnd {
			batchStart = i
			break
		}
	}

	// Extract the batch
	batch := s.redoActions[batchStart : batchEnd+1]

	// Remove from redo stack
	s.redoActions = s.redoActions[:batchStart]

	// Add to undo stack
	for _, action := range batch {
		s.undoActions = append(s.undoActions, action)
	}

	// Return the first action of the batch (for cursor positioning)
	return batch[0], true
}

// CanUndo returns true if undo is available
func (s *UndoStack) CanUndo() bool {
	return len(s.undoActions) > 0
}

// CanRedo returns true if redo is available
func (s *UndoStack) CanRedo() bool {
	return len(s.redoActions) > 0
}

// Clear clears both undo and redo stacks
func (s *UndoStack) Clear() {
	s.undoActions = make([]UndoAction, 0)
	s.redoActions = make([]UndoAction, 0)
	s.currentBatch = ""
	s.lastActionTime = time.Time{}
}

// BreakBatch forces the current batch to end
func (s *UndoStack) BreakBatch() {
	if len(s.undoActions) > 0 {
		s.undoActions[len(s.undoActions)-1].IsBatchEnd = true
	}
	s.currentBatch = ""
}

// removeOldestBatch removes the oldest batch from the undo stack
func (s *UndoStack) removeOldestBatch() {
	if len(s.undoActions) == 0 {
		return
	}

	// Find the end of the oldest batch
	batchEnd := 0
	for i := 0; i < len(s.undoActions); i++ {
		if s.undoActions[i].IsBatchEnd {
			batchEnd = i
			break
		}
	}

	// Remove the oldest batch
	s.undoActions = s.undoActions[batchEnd+1:]
}

// generateBatchID generates a unique batch ID
func generateBatchID() string {
	return time.Now().Format("20060102150405.999999999")
}

// CreateInsertAction creates an insert action
func CreateInsertAction(content string, position int, cursorX, cursorY int) UndoAction {
	return UndoAction{
		Type:      ActionTypeInsert,
		Content:   content,
		Position:  position,
		CursorPos: cursorPos{X: cursorX, Y: cursorY},
		Timestamp: time.Now(),
	}
}

// CreateDeleteAction creates a delete action
func CreateDeleteAction(content string, position int, cursorX, cursorY int) UndoAction {
	return UndoAction{
		Type:      ActionTypeDelete,
		Content:   content,
		Position:  position,
		CursorPos: cursorPos{X: cursorX, Y: cursorY},
		Timestamp: time.Now(),
	}
}

// CreatePasteAction creates a paste action
func CreatePasteAction(content string, position int, cursorX, cursorY int) UndoAction {
	return UndoAction{
		Type:      ActionTypePaste,
		Content:   content,
		Position:  position,
		CursorPos: cursorPos{X: cursorX, Y: cursorY},
		Timestamp: time.Now(),
	}
}

// CreatePromptInsertAction creates a prompt insert action
func CreatePromptInsertAction(content string, position int, cursorX, cursorY int) UndoAction {
	return UndoAction{
		Type:      ActionTypePromptInsert,
		Content:   content,
		Position:  position,
		CursorPos: cursorPos{X: cursorX, Y: cursorY},
		Timestamp: time.Now(),
	}
}

// CreatePlaceholderFillAction creates a placeholder fill action
func CreatePlaceholderFillAction(content string, position int, cursorX, cursorY int) UndoAction {
	return UndoAction{
		Type:      ActionTypePlaceholderFill,
		Content:   content,
		Position:  position,
		CursorPos: cursorPos{X: cursorX, Y: cursorY},
		Timestamp: time.Now(),
	}
}

// CreateNewlineAction creates a newline action
func CreateNewlineAction(position int, cursorX, cursorY int) UndoAction {
	return UndoAction{
		Type:      ActionTypeNewline,
		Content:   "\n",
		Position:  position,
		CursorPos: cursorPos{X: cursorX, Y: cursorY},
		Timestamp: time.Now(),
	}
}

// CreateBackspaceAction creates a backspace action
func CreateBackspaceAction(content string, position int, cursorX, cursorY int) UndoAction {
	return UndoAction{
		Type:      ActionTypeBackspace,
		Content:   content,
		Position:  position,
		CursorPos: cursorPos{X: cursorX, Y: cursorY},
		Timestamp: time.Now(),
	}
}

// CreateBatchEditAction creates a batch edit action for applying multiple edits as one undo action
// This is used when applying AI suggestions via diff viewer
func CreateBatchEditAction(originalContent string, edits []ai.Edit, position int, cursorX, cursorY int) UndoAction {
	return UndoAction{
		Type:       ActionTypeBatchEdit,
		Content:    originalContent,
		Edits:      edits,
		Position:   position,
		CursorPos:  cursorPos{X: cursorX, Y: cursorY},
		Timestamp:  time.Now(),
		IsBatchEnd: true, // Always mark as batch end
	}
}
