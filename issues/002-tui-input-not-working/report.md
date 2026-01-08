# Bug Report: TUI Keyboard Input Not Working

**Bug ID**: 002  
**Severity**: Critical  
**Status**: Root Cause Identified, Fix Proposed  
**Reported**: 2026-01-08  
**Milestone**: M2 - Basic TUI Shell

---

## Bug Description

The TUI launches successfully but does not capture or respond to keyboard input. Users can see the "Press 'q' or Ctrl+C to quit" message, but typing any characters (other than 'q') produces no response. The character count in the status bar remains at 0, indicating that keyboard events are not being processed.

---

## Steps to Reproduce

1. Build the application:
   ```bash
   go build -o promptstack ./cmd/promptstack
   ```

2. Launch the TUI:
   ```bash
   ./promptstack
   ```

3. Attempt to type characters (e.g., "Hello World")

4. Observe that:
   - No characters appear to be captured
   - Status bar character count remains at 0
   - Only 'q' key works to quit

---

## Expected Behavior

According to the M2 manual testing guide (Test Scenario 2: Keyboard Input Handling):

- Characters should be captured when typed
- Status bar character count should increment with each character
- Special characters should be captured
- Unicode characters should be captured
- Backspace should remove characters and decrement count

---

## Actual Behavior

- TUI launches successfully
- Status bar displays correctly (Chars: 0, Lines: 0)
- 'q' key works to quit the application
- Ctrl+C works to quit the application
- **All other keyboard input is ignored**
- Character count remains at 0 regardless of typing
- No visual feedback for any typed characters

---

## Environment

- **OS**: macOS Sonoma
- **Terminal**: [To be confirmed]
- **Go version**: [To be confirmed]
- **Application Version**: M2 - Basic TUI Shell

---

## Impact

This is a **critical** bug that completely blocks the core functionality of Milestone 2. The TUI is supposed to capture and track keyboard input, but it's not processing any keyboard events except for the quit command.

---

## Root Cause Analysis

After thorough code review, the root cause has been identified:

### Issue 1: Missing Character Input Handling in App Model
The `ui/app/model.go` Update() method only handles quit commands (Ctrl+C and 'q') but does not process regular character input. When a user types any character other than 'q', the code falls through to the default case and doesn't update the character count.

### Issue 2: Status Bar Not Updated with Character Counts
The status bar component (`ui/statusbar/model.go`) has `SetCharCount()` and `SetLineCount()` methods, but they are never called from the app model. The status bar's Update() method only handles window resize events, not keyboard input.

### Issue 3: No Text Buffer or Input Tracking
The app model doesn't maintain any text buffer or input tracking mechanism. According to the manual testing guide, characters should be "captured" (tracked for counting purposes), even if not displayed (display is M3 functionality).

### Code Evidence:
1. In `ui/app/model.go` lines 43-73, the Update() method:
   - Handles `tea.KeyCtrlC` and `tea.KeyRunes` only for 'q' key
   - Doesn't process regular character input
   - Doesn't update status bar character/line counts

2. In `ui/statusbar/model.go` lines 36-45, the Update() method:
   - Only handles `tea.WindowSizeMsg`
   - Doesn't process keyboard events

---

## Proposed Fix

### Overview
Modify `ui/app/model.go` to track keyboard input and update status bar counts. The solution follows Go style guide principles:
- Unexported fields for encapsulation
- Consistent pointer receiver usage
- Unicode-aware character counting
- Efficient string building

### Implementation Changes

#### 1. Update Model Struct
```go
// Model represents the root application model state.
// It manages the overall TUI state and coordinates child components.
type Model struct {
    statusBar  *statusbar.Model  // pointer for consistent receiver usage
    textBuffer strings.Builder   // efficient for append operations
    charCount  int               // tracks runes, not bytes (Unicode-aware)
    lineCount  int               // tracks newlines
    width      int
    height     int
    quitting   bool
}
```

**Key Changes**:
- Changed `statusBar` from value to pointer type for consistency with `SetCharCount()` and `SetLineCount()` pointer receivers
- Added `textBuffer` as `strings.Builder` for efficient string concatenation
- Added `charCount` to track Unicode runes (not bytes)
- Added `lineCount` to track newlines
- All new fields are unexported for proper encapsulation

#### 2. Update Constructor
```go
// New creates a new root app model with default values.
func New() Model {
    return Model{
        statusBar:  statusbar.New(),  // returns pointer now
        textBuffer: strings.Builder{},
        charCount:  0,
        lineCount:  0,
        width:      80,
        height:     24,
        quitting:   false,
    }
}
```

#### 3. Update the Update() Method
```go
// Update handles incoming messages for the root app model.
// It processes keyboard input, window resize events, and quit commands.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle quit commands
        switch msg.Type {
        case tea.KeyCtrlC:
            m.quitting = true
            return m, tea.Quit
            
        case tea.KeyRunes:
            // Check for 'q' key to quit
            if len(msg.Runes) == 1 && msg.Runes[0] == 'q' {
                m.quitting = true
                return m, tea.Quit
            }
            // Handle regular character input
            for _, r := range msg.Runes {
                m.textBuffer.WriteRune(r)
                m.charCount++
            }
            m.statusBar.SetCharCount(m.charCount)
            
        case tea.KeyEnter:
            // Add newline to buffer and increment line count
            m.textBuffer.WriteRune('\n')
            m.lineCount++
            m.statusBar.SetLineCount(m.lineCount)
            
        case tea.KeyBackspace:
            // Remove last character if buffer is not empty
            if m.charCount > 0 {
                str := m.textBuffer.String()
                if len(str) > 0 {
                    // Convert to runes for proper Unicode handling
                    runes := []rune(str)
                    m.textBuffer.Reset()
                    m.textBuffer.WriteString(string(runes[:len(runes)-1]))
                    m.charCount--
                    m.statusBar.SetCharCount(m.charCount)
                }
            }
            // Note: This implementation has O(n) complexity for each backspace.
            // For M2 scope, this is acceptable. Future optimizations could track
            // runes in a slice for O(1) backspace operations.
        }

    case tea.WindowSizeMsg:
        // Update window dimensions
        m.width = msg.Width
        m.height = msg.Height
    }

    // Update status bar with any message
    var statusCmd tea.Cmd
    statusModel, statusCmd := m.statusBar.Update(msg)
    m.statusBar = statusModel.(*statusbar.Model)
    return m, statusCmd
}
```

**Key Changes**:
- Added `tea.KeyRunes` case to handle regular character input
- Added `tea.KeyEnter` case to track newlines
- Added `tea.KeyBackspace` case to remove characters
- Uses rune-based counting for proper Unicode support
- Updates status bar counts after each input event
- Changed status bar type assertion to pointer type

#### 4. Update statusbar Package
```go
// In ui/statusbar/model.go, update New() to return pointer:
func New() *Model {
    return &Model{
        charCount: 0,
        lineCount: 0,
        width:     80,
    }
}
```

### Style Guide Compliance

This solution follows the Go style guide principles from `docs/plans/fresh-build/go-style-guide.md`:

1. **Encapsulation** (lines 74-95): All new fields (`textBuffer`, `charCount`, `lineCount`) are unexported
2. **Consistent Receivers** (lines 98-109): Changed `statusBar` to pointer type for consistency with `SetCharCount()` and `SetLineCount()` pointer receivers
3. **Unicode Support**: Uses rune-based counting instead of byte-based `len()` for proper Unicode character handling
4. **Efficient String Building**: Uses `strings.Builder` instead of string concatenation for better performance
5. **Clear Error Handling**: Validates buffer state before operations (e.g., checking `charCount > 0` before backspace)
6. **Explicit Behavior**: All state changes are explicit and traceable
7. **Bubble Tea Pattern Compliance**: Uses value receiver in Update() method following Bubble Tea's immutable model pattern
8. **Test Integration**: Tests should be added to existing `ui/app/model_test.go` file rather than as standalone functions

### Testing Requirements

#### Unit Tests to Add/Update:

1. **Character Input Tests** (update existing in `ui/app/model_test.go`):
   ```go
   func TestUpdateHandlesCharacterInputWithCounting(t *testing.T) {
       model := New()
       
       // Type "Hello"
       chars := []rune{'H', 'e', 'l', 'l', 'o'}
       for i, char := range chars {
           msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}}
           newModel, _ := model.Update(msg)
           model = newModel.(Model)
           
           if model.statusBar.GetCharCount() != i+1 {
               t.Errorf("Expected char count %d, got %d", i+1, model.statusBar.GetCharCount())
           }
       }
   }
   ```

2. **Enter Key Tests** (new):
   ```go
   func TestUpdateHandlesEnterKey(t *testing.T) {
       model := New()
       
       msg := tea.KeyMsg{Type: tea.KeyEnter}
       newModel, _ := model.Update(msg)
       model = newModel.(Model)
       
       if model.statusBar.GetLineCount() != 1 {
           t.Errorf("Expected line count 1, got %d", model.statusBar.GetLineCount())
       }
   }
   ```

3. **Backspace Tests** (new):
   ```go
   func TestUpdateHandlesBackspace(t *testing.T) {
       model := New()
       
       // Type "abc"
       for _, char := range []rune{'a', 'b', 'c'} {
           msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}}
           newModel, _ := model.Update(msg)
           model = newModel.(Model)
       }
       
       // Backspace once
       msg := tea.KeyMsg{Type: tea.KeyBackspace}
       newModel, _ := model.Update(msg)
       model = newModel.(Model)
       
       if model.statusBar.GetCharCount() != 2 {
           t.Errorf("Expected char count 2 after backspace, got %d", model.statusBar.GetCharCount())
       }
   }
   ```

4. **Unicode Tests** (new):
   ```go
   func TestUpdateHandlesUnicodeCharacters(t *testing.T) {
       model := New()
       
       // Type emoji and multi-byte characters
       chars := []string{"😀", "中", "é"}
       for i, char := range chars {
           msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(char)}
           newModel, _ := model.Update(msg)
           model = newModel.(Model)
           
           if model.statusBar.GetCharCount() != i+1 {
               t.Errorf("Expected char count %d for Unicode, got %d", i+1, model.statusBar.GetCharCount())
           }
       }
   }
   ```

5. **Backspace on Empty Buffer** (new):
   ```go
   func TestUpdateHandlesBackspaceOnEmptyBuffer(t *testing.T) {
       model := New()
       
       // Backspace on empty buffer should not panic
       msg := tea.KeyMsg{Type: tea.KeyBackspace}
       newModel, _ := model.Update(msg)
       model = newModel.(Model)
       
       if model.statusBar.GetCharCount() != 0 {
           t.Errorf("Expected char count 0, got %d", model.statusBar.GetCharCount())
       }
   }
   ```

**Note**: All new test functions should be added to the existing `ui/app/model_test.go` file following the established test patterns. Consider integrating new test cases into existing table-driven tests where appropriate (e.g., adding Enter and Backspace cases to `TestUpdateHandlesSpecialKeys`).

#### Manual Testing:
Run all tests from `docs/plans/fresh-build/milestone-implementation-plans/M2/manual-testing-guide.md`:
- **Test Scenario 2**: Keyboard Input Handling (lines 86-129)
  - Verify character counting for ASCII characters
  - Verify special character handling
  - Verify Unicode character handling
  - Verify backspace functionality
- **Test Scenario 3**: Line Count Tracking (lines 132-171)
  - Verify Enter key increments line count
  - Verify multiple Enter keys
  - Verify line count persists across other operations

#### Files to Modify:
1. `ui/app/model.go` - Add fields and update Update() method
2. `ui/statusbar/model.go` - Change New() to return pointer
3. `ui/app/model_test.go` - Add new test functions and update existing ones
4. `ui/statusbar/model_test.go` - Update tests for pointer return type

**Note**: When adding tests, follow the existing table-driven test patterns in `ui/app/model_test.go` for consistency.

---

## Related Documentation

- Manual Testing Guide: `docs/plans/fresh-build/milestone-implementation-plans/M2/manual-testing-guide.md`
- Test Scenario 2: Keyboard Input Handling (lines 86-129)
- Test Scenario 3: Line Count Tracking (lines 132-171)
- Source Code: `ui/app/model.go`, `ui/statusbar/model.go`

---

## Next Steps

1. **Immediate**: Implement the proposed fix in `ui/app/model.go`
2. **Testing**: Update and run tests to verify character and line count updates
3. **Validation**: Perform manual testing per M2 testing guide
4. **Documentation**: Update any documentation if behavior changes

---

## Attachments

- Debug logs: `~/.promptstack/debug.log` (to be collected)
- Screenshots: (to be collected if needed)

---

**Last Updated**: 2026-01-08  
**Assigned To**: Development Team  
**Priority**: P0 - Critical  
**Fix Complexity**: Low (requires ~20-30 lines of code changes)