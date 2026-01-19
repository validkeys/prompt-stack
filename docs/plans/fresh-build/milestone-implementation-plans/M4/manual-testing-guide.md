# Milestone 4 Manual Testing Guide

**Milestone**: 4 - Basic Text Editor
**Purpose**: Verify text editing functionality through manual testing
**Prerequisites**: Go 1.21+ installed, Milestones 1-3 completed, built application

---

## Test Environment Setup

### 1. Verify Previous Milestones **VERIFIED WORKING**

Before testing Milestone 4, ensure previous milestones are complete:

```bash
# Verify config exists (M1)
ls -la ~/.promptstack/config.yaml

# Verify logs are working (M1)
ls -la ~/.promptstack/debug.log

# Verify TUI launches (M2)
go build -o promptstack ./cmd/promptstack
./promptstack  # Should launch TUI, then press 'q' to quit
```

Expected: All previous milestones functioning correctly

### 2. Build the Application **VERIFIED WORKING**

```bash
# Build the application
go build -o promptstack ./cmd/promptstack

# Verify build
./promptstack --help  # Should show help or launch TUI
```

### 3. Verify Terminal Requirements

```bash
# Check terminal size (minimum 80x24 recommended)
echo $COLUMNS $LINES
```

Expected: Terminal should be at least 80 columns wide and 24 lines tall

---

## Test Scenario 1: TUI Launch with Workspace **VERIFIED WORKING**

### Objective

Verify the text editor workspace launches and renders correctly.

### Steps

1. **Launch the application**

   ```bash
   ./promptstack
   ```

2. **Expected Behavior**

   - Terminal should clear and enter alternate screen mode
   - Main workspace area should be visible
   - Status bar should appear at the bottom
   - Cursor should be visible in the workspace
   - No error messages should appear

3. **Verify Workspace Display**

   - Main editor area should occupy most of the screen
   - Status bar at bottom should show: "0 chars, 0 lines"
   - Editor background should use theme background color
   - Cursor should be at top-left position (0, 0)

4. **Verify Theme Colors**
   - Background should be: #1e1e2e (dark blue-gray)
   - Text foreground should be: #cdd6f4 (light gray)
   - Status bar should have: #181825 background (darker blue-gray)
   - Cursor should have: #cdd6f4 background, #11111b foreground

### Success Criteria

- [ ] TUI launches without errors
- [ ] Workspace area is visible
- [ ] Status bar displays initial counts (0 chars, 0 lines)
- [ ] Cursor is visible at position (0, 0)
- [ ] Colors match Catppuccin Mocha theme
- [ ] No visual artifacts or glitches

---

## Test Scenario 2: Character Input **NOT WORKING**

**REPORT** When I open the app, i see the footer and the main window with the label on how to close but i can not type anything

### Objective

Verify characters are typed, displayed, and counted correctly.

### Steps

1. **Launch the application**

   ```bash
   ./promptstack
   ```

2. **Test Basic Typing**

   - Type: `Hello World`
   - Observe: Each character should appear immediately
   - Cursor should move right with each character typed
   - Status bar should update: "11 chars, 1 lines"

3. **Test Character Counting**

   - Verify status bar shows correct count (11 chars)
   - Type additional characters: `! Testing`
   - Status bar should update to: "20 chars, 1 lines"
   - Backspace should decrement count

4. **Test Special Characters**

   - Type: `!@#$%^&*()_+-=[]{}|;':",./<>?`
   - Observe: All special characters should display correctly
   - Character count should increment

5. **Test Unicode Characters**

   - Type: `é ñ 中文 🎉`
   - Observe: All Unicode characters should display correctly
   - Character count should count runes (not bytes)

6. **Test Rapid Typing**
   - Type rapidly: `abcdefghijklmnopqrstuvwxyz`
   - Observe: All characters should appear
   - No characters should be dropped
   - Status bar should update smoothly

### Success Criteria

- [ ] All characters appear immediately when typed
- [ ] Cursor moves correctly with each character
- [ ] Character count updates in real-time
- [ ] Special characters display correctly
- [ ] Unicode characters display correctly
- [ ] No dropped characters during rapid typing
- [ ] Character count is accurate (rune-based)

---

## Test Scenario 3: Cursor Movement (Arrow Keys)

### Objective

Verify cursor navigation with arrow keys works correctly.

### Steps

1. **Launch the application**

   ```bash
   ./promptstack
   ```

2. **Type multi-line content**

   ```
   Line 1: First line of text
   Line 2: Second line of text
   Line 3: Third line of text
   ```

   (Press Enter after each line)

3. **Test Left Arrow**

   - Press Left arrow repeatedly
   - Observe: Cursor moves left one character at a time
   - At start of line (column 0), cursor should stop
   - Should not wrap to previous line

4. **Test Right Arrow**

   - Press Right arrow repeatedly
   - Observe: Cursor moves right one character at a time
   - At end of line, cursor should stop
   - Should not wrap to next line

5. **Test Up Arrow**

   - Press Up arrow
   - Observe: Cursor moves to previous line
   - Cursor column should be preserved if possible
   - At first line (line 0), cursor should stop

6. **Test Down Arrow**

   - Press Down arrow
   - Observe: Cursor moves to next line
   - Cursor column should be preserved if possible
   - At last line, cursor should stop

7. **Test Column Preservation**
   - Move to line 1, column 10
   - Press Up arrow to line 0
   - If line 0 has fewer than 10 characters, cursor should go to end of line
   - Press Down arrow to line 1
   - Cursor should return to column 10

### Success Criteria

- [ ] Left arrow moves cursor left within line
- [ ] Right arrow moves cursor right within line
- [ ] Up arrow moves cursor up, preserving column
- [ ] Down arrow moves cursor down, preserving column
- [ ] Cursor stops at line boundaries
- [ ] Column preservation works correctly
- [ ] No cursor out of bounds errors

---

## Test Scenario 4: Cursor Movement (Home/End Keys)

### Objective

Verify Home and End keys navigate to line boundaries correctly.

### Steps

1. **Launch the application**

   ```bash
   ./promptstack
   ```

2. **Type multi-line content**

   ```
   Line 1: First line with some text
   Line 2: Second line with some text
   ```

3. **Test Home Key**

   - Move cursor to middle of line (e.g., column 10)
   - Press Home key
   - Observe: Cursor should jump to column 0 (start of line)
   - Cursor should always move to column 0 regardless of position

4. **Test End Key**

   - Move cursor to start of line (column 0)
   - Press End key
   - Observe: Cursor should jump to end of line
   - Cursor should be after the last character

5. **Test Home on Empty Line**

   - Create empty line (press Enter twice)
   - Move cursor to empty line
   - Press Home key
   - Observe: Cursor should stay at column 0

6. **Test End on Empty Line**

   - Move cursor to empty line
   - Press End key
   - Observe: Cursor should stay at column 0

7. **Test Home/End Navigation**
   - Type: `Test line for Home/End testing`
   - Press Home (cursor at column 0)
   - Press End (cursor at end of line)
   - Press Home again (cursor at column 0)
   - Repeat multiple times

### Success Criteria

- [ ] Home key moves cursor to column 0
- [ ] End key moves cursor to end of line
- [ ] Home works correctly on empty lines
- [ ] End works correctly on empty lines
- [ ] Repeated Home/End presses work correctly
- [ ] Navigation is immediate and smooth

---

## Test Scenario 5: Line Editing (Enter Key)

### Objective

Verify Enter key creates new lines correctly.

### Steps

1. **Launch the application**

   ```bash
   ./promptstack
   ```

2. **Test Single Line**

   - Type: `First line`
   - Press Enter
   - Observe: Cursor should move to next line (line 1)
   - Status bar should show: "10 chars, 2 lines"

3. **Test Multiple Lines**

   - Type: `Second line`
   - Press Enter
   - Observe: Cursor should move to line 2
   - Status bar should show: "21 chars, 3 lines"

4. **Test Empty Lines**

   - Press Enter without typing
   - Observe: Cursor should move to next line
   - Line count should increment, character count unchanged

5. **Test Line Continuation**
   - Type: `Line part 1`
   - Press Enter
   - Type: `Line part 2`
   - Observe: Two separate lines should be created
   - Both lines should be visible

### Success Criteria

- [ ] Enter creates new line
- [ ] Cursor moves to start of new line
- [ ] Line count increments correctly
- [ ] Character count unaffected
- [ ] Empty lines handled correctly
- [ ] Multi-line content displays correctly

---

## Test Scenario 6: Character Deletion (Backspace Key)

### Objective

Verify Backspace key deletes characters correctly.

### Steps

1. **Launch the application**

   ```bash
   ./promptstack
   ```

2. **Type content**

   - Type: `Hello World`
   - Press Backspace once
   - Observe: Last character should be deleted
   - Cursor should move left one position
   - Status bar should update to: "10 chars, 1 lines"

3. **Test Multiple Backspaces**

   - Press Backspace 5 times
   - Observe: 5 characters should be deleted
   - Cursor should move left 5 positions
   - Status bar should update to: "5 chars, 1 lines"

4. **Test Backspace at Line Start**

   - Type content on multiple lines:
     ```
     Line 1
     Line 2
     ```
   - Move cursor to start of line 2
   - Press Backspace
   - Observe: Cursor should move to end of line 1
   - Lines should merge (no newline)

5. **Test Backspace at Document Start**

   - Delete all content
   - Ensure cursor is at position (0, 0)
   - Press Backspace
   - Observe: Nothing should be deleted
   - No error should occur

6. **Test Backspace with Empty Buffer**
   - Press Backspace on empty buffer
   - Observe: Nothing should happen
   - No error should occur

### Success Criteria

- [ ] Backspace deletes character at cursor position
- [ ] Cursor moves left after deletion
- [ ] Character count decrements
- [ ] Line count decrements when joining lines
- [ ] No deletion at document start
- [ ] No errors on empty buffer

---

## Test Scenario 7: Character Deletion (Delete Key)

### Objective

Verify Delete key deletes characters correctly.

### Steps

1. **Launch the application**

   ```bash
   ./promptstack
   ```

2. **Type content**

   - Type: `Hello World`
   - Move cursor to position after "Hello" (before space)
   - Press Delete
   - Observe: Space should be deleted
   - Cursor should not move
   - Content should become: `HelloWorld`

3. **Test Multiple Deletes**

   - Press Delete 5 times
   - Observe: 5 characters should be deleted
   - Cursor should stay in place
   - Character count should decrease

4. **Test Delete at Line End**

   - Type multi-line content:
     ```
     Line 1
     Line 2
     ```
   - Move cursor to end of line 1
   - Press Delete
   - Observe: Newline should be deleted
   - Lines should merge: `Line 1Line 2`

5. **Test Delete at Document End**

   - Move cursor to end of document
   - Press Delete
   - Observe: Nothing should be deleted
   - No error should occur

6. **Test Delete with Empty Buffer**
   - Press Delete on empty buffer
   - Observe: Nothing should happen
   - No error should occur

### Success Criteria

- [ ] Delete deletes character after cursor
- [ ] Cursor does not move after deletion
- [ ] Character count decrements
- [ ] Line count decrements when joining lines
- [ ] No deletion at document end
- [ ] No errors on empty buffer

---

## Test Scenario 8: Viewport Scrolling

### Objective

Verify viewport scrolls to keep cursor visible.

### Steps

1. **Launch the application**

   ```bash
   ./promptstack
   ```

2. **Resize terminal to small size**

   - Resize to approximately 80x20
   - This makes viewport visible (20 lines tall, minus 1 for status bar)

3. **Create content larger than viewport**

   - Type 30 lines of content:
     ```
     Line 1
     Line 2
     ...
     Line 30
     ```
     (Use copy-paste or rapid typing)

4. **Test Downward Scrolling**

   - Move cursor down with Down arrow
   - Observe: When cursor reaches middle of viewport, scrolling should start
   - Viewport should adjust to keep cursor visible
   - Cursor should always be visible

5. **Test Upward Scrolling**

   - Move to bottom of document
   - Move cursor up with Up arrow
   - Observe: Viewport should scroll up when cursor reaches middle
   - Cursor should always remain visible

6. **Test Middle-Third Scrolling**

   - Move cursor to top of viewport
   - Press Down arrow repeatedly
   - Observe: Scrolling should start when cursor passes 1/3 mark
   - Cursor should be kept in middle 1/3 of viewport
   - Scrolling should be smooth

7. **Test Rapid Scrolling**

   - Hold Down arrow key
   - Observe: Rapid scrolling should work
   - No visual glitches
   - Cursor remains visible

8. **Test Scroll to Top**

   - Press Home key
   - Observe: Viewport should scroll to top
   - First line should be visible

9. **Test Scroll to Bottom**
   - Press End key
   - Observe: Viewport should scroll to bottom
   - Last line should be visible

### Success Criteria

- [ ] Viewport scrolls when cursor reaches boundaries
- [ ] Cursor always remains visible
- [ ] Middle-third scrolling strategy works
- [ ] Scrolling is smooth
- [ ] No visual artifacts during scroll
- [ ] Rapid scrolling works
- [ ] Scroll to top/bottom works
- [ ] No negative YOffset errors

---

## Test Scenario 9: Window Resize Handling

### Objective

Verify editor handles terminal window resizing correctly.

### Steps

1. **Launch the application**

   ```bash
   ./promptstack
   ```

2. **Type multi-line content**

   - Type 15-20 lines of content

3. **Test Smaller Resize**

   - Resize terminal to 60x20
   - Observe: Workspace should adjust
   - Status bar should adjust width
   - No visual glitches
   - Content should remain accessible

4. **Test Larger Resize**

   - Resize terminal to 120x40
   - Observe: Workspace should expand
   - More content should be visible
   - Status bar should adjust width
   - No visual glitches

5. **Test Minimum Size**

   - Resize terminal to minimum (e.g., 40x10)
   - Observe: Editor should still render
   - Status bar should be visible
   - Content should scroll if needed

6. **Test Rapid Resizing**

   - Quickly resize terminal multiple times
   - Observe: No crashes
   - No visual artifacts
   - Editor remains responsive

7. **Test Resize During Editing**
   - While typing content, resize terminal
   - Observe: Typing should continue
   - No dropped characters
   - Editor should adjust in real-time

### Success Criteria

- [ ] Workspace adjusts to terminal size
- [ ] Status bar adjusts width
- [ ] No visual glitches on resize
- [ ] Works at minimum terminal size
- [ ] Works at maximum terminal size
- [ ] Rapid resizing doesn't cause crashes
- [ ] Resize during editing works
- [ ] No YOffset errors

---

## Test Scenario 10: Character and Line Counting

### Objective

Verify character and line counts are accurate and update in real-time.

### Steps

1. **Launch the application**

   ```bash
   ./promptstack
   ```

2. **Verify Initial Counts**

   - Status bar should show: "0 chars, 0 lines"

3. **Test Character Counting**

   - Type: `Hello` (5 characters)
   - Status bar should show: "5 chars, 0 lines"
   - Type: ` World` (6 characters including space)
   - Status bar should show: "11 chars, 0 lines"
   - Verify count is accurate

4. **Test Line Counting**

   - Press Enter
   - Status bar should show: "11 chars, 1 lines"
   - Type: `Second line`
   - Status bar should show: "22 chars, 2 lines"
   - Press Enter
   - Status bar should show: "22 chars, 3 lines"

5. **Test Count Updates on Delete**

   - Press Backspace
   - Character count should decrement
   - Move cursor to line boundary and delete newline
   - Line count should decrement

6. **Test Unicode Counting**

   - Type: `é ñ 中文` (7 runes)
   - Character count should show 7
   - Type emoji: `🎉` (1 rune)
   - Character count should show 8
   - Verify rune-based counting (not byte-based)

7. **Test Large Document**
   - Create document with 100+ lines
   - Verify line count is accurate
   - Type 1000+ characters
   - Verify character count is accurate

### Success Criteria

- [ ] Initial counts are correct (0, 0)
- [ ] Character count increments with each character
- [ ] Line count increments with each Enter
- [ ] Character count decrements on delete
- [ ] Line count decrements on line join
- [ ] Unicode characters counted as runes
- [ ] Counts accurate for large documents
- [ ] Counts update in real-time

---

## Test Scenario 11: Theme Application

### Objective

Verify Catppuccin Mocha theme is applied correctly to all components.

### Steps

1. **Launch the application**

   ```bash
   ./promptstack
   ```

2. **Verify Background Colors**

   - Main editor background: #1e1e2e (dark blue-gray)
   - Status bar background: #181825 (darker blue-gray)

3. **Verify Foreground Colors**

   - Editor text: #cdd6f4 (light gray)
   - Status bar text: #a6adc8 (medium gray)

4. **Verify Cursor Styling**

   - Cursor background: #cdd6f4 (light gray)
   - Cursor text: #11111b (dark)
   - Cursor should have high contrast (easily visible)

5. **Verify Status Bar Styling**

   - Status bar should be visually distinct from editor
   - Background color should be different
   - Text should be readable
   - Spacing should follow 1-unit system

6. **Verify Color Consistency**

   - All components should use consistent palette
   - No hardcoded color values
   - Colors should match design system

7. **Test in Different Terminals**
   - Test in 256-color terminal
   - Test in truecolor terminal
   - Colors should work correctly in both

### Success Criteria

- [ ] Editor background is #1e1e2e
- [ ] Status bar background is #181825
- [ ] Text foreground is #cdd6f4
- [ ] Cursor has high contrast
- [ ] Status bar is visually distinct
- [ ] Colors are consistent
- [ ] Works in 256-color terminal
- [ ] Works in truecolor terminal

---

## Test Scenario 12: Quit Functionality

### Objective

Verify editor exits cleanly.

### Steps

1. **Launch the application**

   ```bash
   ./promptstack
   ```

2. **Type some content**

   - Type: `Testing quit functionality`

3. **Press 'q' to quit**

   - Press lowercase 'q'
   - Observe: TUI should exit immediately
   - Terminal should return to normal mode
   - No error messages

4. **Verify Clean Exit**

   - Terminal prompt should be visible
   - No error messages in terminal
   - Check logs:
     ```bash
     cat ~/.promptstack/debug.log | grep -E "TUI shutdown|workspace"
     ```

5. **Test Ctrl+C quit**
   - Launch application again
   - Type content
   - Press Ctrl+C
   - Observe: TUI should exit immediately

### Success Criteria

- [ ] 'q' key exits TUI immediately
- [ ] Ctrl+C exits TUI immediately
- [ ] Terminal returns to normal mode
- [ ] No error messages
- [ ] Exit is logged
- [ ] Clean exit with no state corruption

---

## Test Scenario 13: Performance Testing

### Objective

Verify editor performance meets requirements.

### Steps

1. **Measure Startup Time**

   ```bash
   time ./promptstack
   ```

   Expected: Startup time < 100ms

2. **Test Input Response**

   - Type rapidly: `abcdefghijklmnopqrstuvwxyz`
   - Observe: Each character should appear immediately
   - Perceived latency < 50ms
   - No lag or delay

3. **Test Cursor Movement Performance**

   - Press arrow keys rapidly
   - Observe: Cursor movement should be smooth
   - No lag
   - Movement should complete in < 10ms per operation

4. **Test Large Document**

   - Create document with 10,000 characters
   - Type and move cursor
   - Observe: No lag
   - Operations should complete quickly

5. **Test Long Session**

   - Keep editor open for 5+ minutes
   - Type periodically
   - Observe: No memory leaks
   - No performance degradation

6. **Test Memory Usage**
   - Launch editor
   - Check memory in another terminal:
     ```bash
     ps aux | grep promptstack
     ```
     Expected: Memory usage < 50MB

### Success Criteria

- [ ] Startup time < 100ms
- [ ] Input latency < 50ms
- [ ] Cursor movement < 10ms
- [ ] Large documents handled smoothly
- [ ] No memory leaks
- [ ] No performance degradation
- [ ] Memory usage reasonable

---

## Test Scenario 14: Edge Cases

### Objective

Verify editor handles edge cases correctly.

### Steps

1. **Test Empty Buffer**

   - Launch editor
   - Try all operations (typing, deletion, cursor movement)
   - Observe: No crashes or errors

2. **Test Single Character**

   - Type single character: `a`
   - Test all operations
   - Observe: Works correctly

3. **Test Very Long Lines**

   - Create line with 1000+ characters
   - Type more content
   - Test cursor movement
   - Observe: No performance issues

4. **Test Many Empty Lines**

   - Press Enter 100 times
   - Test cursor movement
   - Observe: Works correctly

5. **Test Unicode Edge Cases**

   - Type: `🎉` (emoji)
   - Type: `中文` (multi-byte characters)
   - Type: `é` (accented character)
   - Test all operations
   - Observe: Works correctly

6. **Test Rapid Key Presses**
   - Press arrow keys rapidly
   - Type rapidly
   - Delete rapidly
   - Observe: No dropped operations
   - No crashes

### Success Criteria

- [ ] Empty buffer handled correctly
- [ ] Single character works
- [ ] Very long lines work
- [ ] Many empty lines work
- [ ] Unicode edge cases work
- [ ] Rapid input works
- [ ] No crashes in edge cases

---

## Test Scenario 15: Integration Testing

### Objective

Verify all components work together correctly.

### Steps

1. **Test Complete Editing Workflow**

   - Launch editor
   - Type multi-line content
   - Navigate with arrow keys
   - Edit with Backspace/Delete
   - Scroll with viewport
   - Verify counts update
   - Quit cleanly

2. **Test Real-World Editing**

   - Type a short paragraph of text
   - Make edits throughout
   - Navigate to different positions
   - Delete and add content
   - Verify everything works together

3. **Test Window Resize During Editing**

   - Type content
   - Resize terminal while typing
   - Continue editing
   - Verify no issues

4. **Test Large Document Editing**
   - Create large document (100+ lines)
   - Navigate throughout
   - Make edits at various positions
   - Verify smooth performance

### Success Criteria

- [ ] Complete editing workflow works
- [ ] Real-world editing is smooth
- [ ] Resize during editing works
- [ ] Large document editing works
- [ ] All components integrate correctly

---

## Cleanup After Testing

```bash
# Remove built binary
rm -f promptstack

# Optional: Clear logs for clean state
# rm ~/.promptstack/debug.log
```

---

## Test Results Checklist

After completing all scenarios, verify:

- [ ] All 15 test scenarios completed
- [ ] All success criteria met
- [ ] No crashes or panics observed
- [ ] All error messages are clear and helpful
- [ ] Editor behavior is consistent
- [ ] Performance meets requirements
- [ ] Theme colors are correct
- [ ] Documentation matches actual behavior

---

## Known Limitations

1. **File Operations**: Not yet implemented (coming in M5)
2. **Undo/Redo**: Not yet implemented (coming in M6)
3. **Search/Replace**: Not yet implemented (coming in later milestones)
4. **Syntax Highlighting**: Not yet implemented (coming in later milestones)

---

## Notes

- Record any unexpected behavior or bugs found
- Note any areas where UX could be improved
- Document any deviations from expected behavior
- Take screenshots of any issues for reference
- Note any performance issues or lag

---

## Bug Report Template

If you find a bug, use this template:

```
**Bug Description**: [Brief description]

**Steps to Reproduce**:
1. [Step 1]
2. [Step 2]
3. [Step 3]

**Expected Behavior**: [What should happen]

**Actual Behavior**: [What actually happened]

**Environment**:
- OS: [macOS/Linux/Windows]
- Terminal: [Terminal name and version]
- Go version: [go version]

**Screenshots**: [Attach if applicable]

**Logs**: [Attach ~/.promptstack/debug.log]
```

---

## What to Look For

### Visual Quality

- ✓ Cursor is always visible and has high contrast
- ✓ Text is readable with proper contrast
- ✓ No visual flickering or artifacts
- ✓ Smooth scrolling without jumping
- ✓ Colors are consistent and match theme

### Performance

- ✓ Typing appears immediately (< 50ms latency)
- ✓ Cursor movement is smooth (< 10ms per operation)
- ✓ No lag when editing large documents
- ✓ No memory leaks during long sessions
- ✓ Responsive to rapid input

### Functionality

- ✓ All editing operations work correctly
- ✓ Character and line counts are accurate
- ✓ Cursor movement is accurate and smooth
- ✓ Viewport scrolling works correctly
- ✓ Window resizing handled gracefully

### User Experience

- ✓ Editor feels natural and responsive
- ✓ No confusing behavior
- ✓ Error messages are clear (if any)
- ✓ Smooth and predictable interactions
- ✓ No crashes or panics

---

## Feedback Points

Please provide feedback on:

1. **Typing Experience**

   - Does typing feel smooth and responsive?
   - Any noticeable lag or delay?
   - Are characters appearing immediately?

2. **Cursor Movement**

   - Is cursor movement smooth?
   - Does it feel natural?
   - Any unexpected behavior at boundaries?

3. **Viewport Scrolling**

   - Is scrolling smooth?
   - Does the middle-third strategy feel right?
   - Any jumping or jarring movements?

4. **Theme and Visuals**

   - Are colors readable and pleasant?
   - Is cursor visibility good?
   - Any visual artifacts or glitches?

5. **Performance**

   - Any performance issues?
   - Lag on large documents?
   - Memory usage concerns?

6. **Overall UX**
   - Does the editor feel natural?
   - Any confusing interactions?
   - Suggestions for improvement?

---

**Last Updated**: 2026-01-09
**Tested By**: ******\_\_\_******
**Test Date**: ******\_\_\_******
**Overall Result**: [ ] Pass [ ] Fail
