# Milestone 2 Manual Testing Guide

**Milestone**: 2 - Basic TUI Shell
**Purpose**: Verify TUI functionality through manual testing
**Prerequisites**: Go 1.21+ installed, Milestone 1 completed, built application

---

## Test Environment Setup

### 1. Verify Milestone 1 Completion

Before testing Milestone 2, ensure Milestone 1 is complete:

```bash
# Verify config exists
ls -la ~/.promptstack/config.yaml

# Verify config is valid
cat ~/.promptstack/config.yaml
```

Expected: Config file exists with valid YAML structure

### 2. Build the Application

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

## Test Scenario 1: TUI Launch

### Objective
Verify the TUI launches successfully with proper rendering.

### Steps

1. **Launch the application**
   ```bash
   ./promptstack
   ```

2. **Expected Behavior**
   - Terminal should clear and enter alternate screen mode
   - Status bar should appear at the bottom of the screen
   - No error messages should appear
   - Application should be responsive to keyboard input

3. **Verify Status Bar Display**
   - Status bar should be at the bottom of the screen
   - Should display character count (initially 0)
   - Should display line count (initially 0)
   - Should use Catppuccin Mocha color scheme

4. **Verify Theme Colors**
   - Background should be dark (#1e1e2e)
   - Text should be light (#cdd6f4)
   - Status bar should have accent colors (blue/green/yellow/red)
   - Borders should be visible (#45475a)

### Success Criteria
- [ ] TUI launches without errors
- [ ] Terminal enters alternate screen mode
- [ ] Status bar displays at bottom of screen
- [ ] Character count shows 0
- [ ] Line count shows 0
- [ ] Colors match Catppuccin Mocha theme

---

## Test Scenario 2: Keyboard Input Handling

### Objective
Verify the TUI captures and processes keyboard input correctly.

### Steps

1. **Launch the application**
   ```bash
   ./promptstack
   ```
   STATUS: PASS

2. **Test Character Input**
   - Type: `Hello World`
   - Observe: Characters should be captured
   - Status bar character count should increment with each character
   STATUS: FAIL. I AM UNABLE TO TYPE ANYHING. BUG REPORT: /Users/kyledavis/Sites/prompt-stack/issues/002-tui-input-not-working/report.md

3. **Verify Character Count**
   - After typing "Hello World" (11 characters including space)
   - Status bar should show: `Chars: 11`

4. **Test Special Characters**
   - Type: `!@#$%^&*()`
   - Observe: Special characters should be captured
   - Character count should increment

5. **Test Unicode Characters**
   - Type: `é ñ 中文`
   - Observe: Unicode characters should be captured
   - Character count should increment correctly

6. **Test Backspace**
   - Press Backspace key
   - Observe: Last character should be removed
   - Character count should decrement

### Success Criteria
- [ ] All characters are captured
- [ ] Character count updates in real-time
- [ ] Special characters work correctly
- [ ] Unicode characters work correctly
- [ ] Backspace removes characters
- [ ] Character count decrements on backspace

---

## Test Scenario 3: Line Count Tracking

### Objective
Verify line count updates correctly when Enter is pressed.

### Steps

1. **Launch the application**
   ```bash
   ./promptstack
   ```

2. **Test Single Line**
   - Type: `First line`
   - Press Enter
   - Observe: Line count should increment to 1

3. **Test Multiple Lines**
   - Type: `Second line`
   - Press Enter
   - Observe: Line count should increment to 2
   - Type: `Third line`
   - Press Enter
   - Observe: Line count should increment to 3

4. **Verify Status Bar**
   - Status bar should show: `Lines: 3`
   - Character count should reflect total characters typed

5. **Test Empty Lines**
   - Press Enter without typing
   - Observe: Line count should increment
   - Character count should remain unchanged

### Success Criteria
- [ ] Line count increments on Enter
- [ ] Line count displays correctly in status bar
- [ ] Empty lines increment line count
- [ ] Line count and character count are independent

---

## Test Scenario 4: Quit Functionality (q key)

### Objective
Verify the TUI exits cleanly when 'q' is pressed.

### Steps

1. **Launch the application**
   ```bash
   ./promptstack
   ```

2. **Type some text**
   - Type: `Testing quit functionality`

3. **Press 'q' key**
   - Press lowercase 'q'
   - Observe: TUI should exit immediately

4. **Verify Clean Exit**
   - Terminal should return to normal mode
   - No error messages should appear
   - Terminal prompt should be visible

5. **Check Logs**
   ```bash
   cat ~/.promptstack/debug.log | grep "TUI shutdown"
   ```
   Expected: Log entry showing TUI shutdown

### Success Criteria
- [ ] 'q' key exits TUI immediately
- [ ] Terminal returns to normal mode
- [ ] No error messages on exit
- [ ] Shutdown is logged
- [ ] Exit is clean and graceful

---

## Test Scenario 5: Quit Functionality (Ctrl+C)

### Objective
Verify the TUI exits cleanly when Ctrl+C is pressed.

### Steps

1. **Launch the application**
   ```bash
   ./promptstack
   ```

2. **Type some text**
   - Type: `Testing Ctrl+C quit`

3. **Press Ctrl+C**
   - Hold Ctrl and press C
   - Observe: TUI should exit immediately

4. **Verify Clean Exit**
   - Terminal should return to normal mode
   - No error messages should appear
   - Terminal prompt should be visible

5. **Check Logs**
   ```bash
   cat ~/.promptstack/debug.log | grep "TUI shutdown"
   ```
   Expected: Log entry showing TUI shutdown

### Success Criteria
- [ ] Ctrl+C exits TUI immediately
- [ ] Terminal returns to normal mode
- [ ] No error messages on exit
- [ ] Shutdown is logged
- [ ] Exit is clean and graceful

---

## Test Scenario 6: Window Resize Handling

### Objective
Verify the TUI handles terminal window resizing correctly.

### Steps

1. **Launch the application**
   ```bash
   ./promptstack
   ```

2. **Resize Terminal Smaller**
   - Drag terminal corner to make it smaller (e.g., 60x20)
   - Observe: Status bar should adjust to new width
   - No visual glitches should occur

3. **Resize Terminal Larger**
   - Drag terminal corner to make it larger (e.g., 120x40)
   - Observe: Status bar should adjust to new width
   - No visual glitches should occur

4. **Test Extreme Sizes**
   - Resize to minimum (e.g., 40x10)
   - Observe: TUI should still render
   - Status bar should be visible
   - Resize back to normal size

5. **Test Rapid Resizing**
   - Quickly resize terminal multiple times
   - Observe: No crashes or visual artifacts
   - Status bar should adjust smoothly

### Success Criteria
- [ ] Status bar adjusts to window width
- [ ] No visual glitches on resize
- [ ] TUI handles small terminals
- [ ] TUI handles large terminals
- [ ] Rapid resizing doesn't cause crashes
- [ ] Rendering remains smooth during resize

---

## Test Scenario 7: Theme System Verification

### Objective
Verify the Catppuccin Mocha theme is applied correctly.

### Steps

1. **Launch the application**
   ```bash
   ./promptstack
   ```

2. **Verify Background Colors**
   - Main background should be: #1e1e2e (dark blue-gray)
   - Secondary background should be: #181825 (darker blue-gray)

3. **Verify Foreground Colors**
   - Primary text should be: #cdd6f4 (light gray)
   - Muted text should be: #a6adc8 (medium gray)

4. **Verify Accent Colors**
   - Blue accent should be: #89b4fa
   - Green accent should be: #a6e3a1
   - Yellow accent should be: #f9e2af
   - Red accent should be: #f38ba8

5. **Verify Border Colors**
   - Borders should be: #45475a (medium gray)

6. **Take Screenshot for Reference**
   - Use terminal screenshot tool
   - Compare with Catppuccin Mocha reference

### Success Criteria
- [ ] Background colors match Catppuccin Mocha
- [ ] Foreground colors match Catppuccin Mocha
- [ ] Accent colors match Catppuccin Mocha
- [ ] Border colors match Catppuccin Mocha
- [ ] Colors are consistent across components
- [ ] No color bleeding or artifacts

---

## Test Scenario 8: Status Bar Component

### Objective
Verify the status bar component functions correctly.

### Steps

1. **Launch the application**
   ```bash
   ./promptstack
   ```

2. **Verify Status Bar Position**
   - Status bar should be at the bottom of the screen
   - Should span full width of terminal
   - Should be clearly visible

3. **Verify Status Bar Content**
   - Should display character count label: "Chars:"
   - Should display line count label: "Lines:"
   - Should display current counts

4. **Test Status Bar Updates**
   - Type characters: `Test`
   - Observe: Character count updates to 4
   - Press Enter
   - Observe: Line count updates to 1

5. **Test Status Bar Styling**
   - Status bar should have distinct background color
   - Text should be readable
   - Should use theme accent colors

6. **Test Status Bar on Resize**
   - Resize terminal
   - Observe: Status bar adjusts width
   - Content remains visible

### Success Criteria
- [ ] Status bar is at bottom of screen
- [ ] Status bar spans full width
- [ ] Character count displays correctly
- [ ] Line count displays correctly
- [ ] Status bar updates in real-time
- [ ] Status bar styling is consistent
- [ ] Status bar handles resize

---

## Test Scenario 9: Root App Model

### Objective
Verify the root app model coordinates components correctly.

### Steps

1. **Launch the application**
   ```bash
   ./promptstack
   ```

2. **Verify Model Initialization**
   - TUI should start without errors
   - Status bar should be initialized
   - Initial state should be clean

3. **Test Message Handling**
   - Type characters: Verify Update() handles KeyMsg
   - Press Enter: Verify Update() handles KeyMsg
   - Resize terminal: Verify Update() handles WindowSizeMsg

4. **Test View Rendering**
   - Verify View() renders status bar
   - Verify View() uses theme styles
   - Verify View() handles different terminal sizes

5. **Test State Management**
   - Type text, verify state updates
   - Resize terminal, verify state updates
   - Quit, verify clean state cleanup

### Success Criteria
- [ ] Model initializes correctly
- [ ] Update() handles all message types
- [ ] View() renders correctly
- [ ] State is managed properly
- [ ] No state corruption occurs
- [ ] Clean state cleanup on exit

---

## Test Scenario 10: Integration with Main Application

### Objective
Verify TUI integrates correctly with main application.

### Steps

1. **Launch the application**
   ```bash
   ./promptstack
   ```

2. **Verify Bootstrap Integration**
   - Config should be loaded from M1
   - Logger should be initialized from M1
   - No bootstrap errors should occur

3. **Verify TUI Startup**
   - TUI should start after bootstrap
   - Startup should be logged
   - No startup errors should occur

4. **Verify TUI Shutdown**
   - Quit TUI (press 'q')
   - Shutdown should be logged
   - No shutdown errors should occur

5. **Check Logs**
   ```bash
   cat ~/.promptstack/debug.log | grep -E "TUI startup|TUI shutdown"
   ```
   Expected: Both startup and shutdown log entries

6. **Verify Error Handling**
   - Test with invalid config (if possible)
   - Verify errors are logged
   - Verify graceful error handling

### Success Criteria
- [ ] Bootstrap completes before TUI starts
- [ ] Config is loaded correctly
- [ ] Logger is initialized correctly
- [ ] TUI startup is logged
- [ ] TUI shutdown is logged
- [ ] Errors are handled gracefully
- [ ] No crashes occur

---

## Test Scenario 11: Performance Testing

### Objective
Verify TUI performance meets requirements.

### Steps

1. **Measure Startup Time**
   ```bash
   time ./promptstack
   ```
   Expected: Startup time < 100ms

2. **Test Input Response**
   - Type rapidly: `abcdefghijklmnopqrstuvwxyz`
   - Observe: Each character should appear immediately
   - No lag or delay should be noticeable

3. **Test Render Performance**
   - Resize terminal rapidly
   - Observe: Rendering should be smooth
   - No flickering or stuttering

4. **Test Memory Usage**
   - Launch TUI
   - Check memory usage in another terminal:
     ```bash
     ps aux | grep promptstack
     ```
   Expected: Memory usage should be reasonable (< 50MB)

5. **Test Long Running Session**
   - Keep TUI open for 5 minutes
   - Type periodically
   - Observe: No memory leaks or performance degradation

### Success Criteria
- [ ] Startup time < 100ms
- [ ] Input response is immediate
- [ ] Rendering is smooth (60 FPS)
- [ ] Memory usage is reasonable
- [ ] No memory leaks
- [ ] No performance degradation over time

---

## Test Scenario 12: Edge Cases

### Objective
Verify TUI handles edge cases correctly.

### Steps

1. **Test Zero Terminal Size**
   - Resize terminal to minimum possible
   - Observe: TUI should still render
   - No crashes should occur

2. **Test Very Large Terminal**
   - Resize terminal to maximum possible
   - Observe: TUI should still render
   - No performance issues should occur

3. **Test Rapid Keyboard Input**
   - Type as fast as possible
   - Observe: All characters should be captured
   - No characters should be lost

4. **Test Special Key Combinations**
   - Press Ctrl+Z (suspend)
   - Resume with `fg`
   - Observe: TUI should resume correctly

5. **Test Unicode Input**
   - Type various Unicode characters: `é ñ 中文 🎉`
   - Observe: All characters should display correctly
   - Character count should be accurate

6. **Test Multiple Quit Attempts**
   - Press 'q' multiple times rapidly
   - Observe: First 'q' should exit
   - No errors should occur

### Success Criteria
- [ ] Handles minimal terminal size
- [ ] Handles maximal terminal size
- [ ] Captures all rapid input
- [ ] Handles suspend/resume
- [ ] Handles Unicode correctly
- [ ] Handles multiple quit attempts
- [ ] No crashes in edge cases

---

## Test Scenario 13: Logging Verification

### Objective
Verify TUI operations are logged correctly.

### Steps

1. **Clear existing logs**
   ```bash
   rm ~/.promptstack/debug.log
   ```

2. **Launch and quit TUI**
   ```bash
   ./promptstack
   # Type some text, then press 'q'
   ```

3. **Check log file**
   ```bash
   cat ~/.promptstack/debug.log
   ```

4. **Verify Log Entries**
   Expected log entries:
   - TUI startup message
   - TUI shutdown message
   - Any errors (if they occurred)

5. **Verify Log Format**
   - Logs should be in JSON format
   - Should have: level, timestamp, msg fields
   - Should be structured and parseable

6. **Test Debug Logging**
   ```bash
   LOG_LEVEL=debug ./promptstack
   cat ~/.promptstack/debug.log | grep "level\":\"debug"
   ```
   Expected: Debug messages should appear

### Success Criteria
- [ ] TUI startup is logged
- [ ] TUI shutdown is logged
- [ ] Logs are in JSON format
- [ ] Logs have required fields
- [ ] Debug logging works
- [ ] No log corruption

---

## Cleanup After Testing

```bash
# Remove built binary
rm -f promptstack

# Optional: Remove test logs (keep config for next milestone)
# rm ~/.promptstack/debug.log
```

---

## Test Results Checklist

After completing all scenarios, verify:

- [ ] All 13 test scenarios completed
- [ ] All success criteria met
- [ ] No crashes or panics observed
- [ ] All error messages are clear and helpful
- [ ] TUI behavior is consistent
- [ ] Performance meets requirements
- [ ] Theme colors are correct
- [ ] Documentation matches actual behavior

---

## Known Limitations

1. **TUI Functionality**: Only basic shell implemented; advanced features coming in later milestones
2. **Input Display**: Characters are captured but not displayed in main area (coming in M3)
3. **Status Bar**: Only shows counts; more info coming in later milestones
4. **Theme**: Only Catppuccin Mocha implemented; other themes coming later

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

**Last Updated**: 2026-01-08
**Tested By**: _______________
**Test Date**: _______________
**Overall Result**: [ ] Pass [ ] Fail