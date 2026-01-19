# Milestone 5: Auto-save - Technical Reference

## Message Flow Diagram

```
User Types → markDirty() → scheduleAutoSave()
    ↓
  (500ms timer)
    ↓
autoSaveMsg{timerID: N} → Update() → saveStatus = "saving"
    ↓
tea.Cmd → saveToHistory() → (async file write)
    ↓
saveSuccessMsg/saveErrorMsg → Update()
    ↓
saveStatus = "saved"/"error" → (auto-clear timer)
    ↓
clearSaveStatusMsg → Update() → saveStatus = ""
```

## Message Types

```go
// Sent after debounce timer expires
type autoSaveMsg struct {
    timerID int  // Matches model.autoSaveTimerID to ignore stale timers
}

// Sent by async save command on success
type saveSuccessMsg struct{}

// Sent by async save command on error  
type saveErrorMsg struct {
    err error
}

// Sent by auto-clear timer to reset status
type clearSaveStatusMsg struct{}
```

## Model State Additions

```go
type Model struct {
    // ... existing fields
    
    // Auto-save state
    autoSaveTimerID int     // Incremented on each new timer
    saveStatus      string  // "saving", "saved", "error", or ""
    saveError       string  // Error message when saveStatus == "error"
}
```

## Key Methods

### `scheduleAutoSave() (Model, tea.Cmd)`
- Increments `autoSaveTimerID`
- Returns `tea.Tick(500ms)` command that sends `autoSaveMsg` with current timer ID
- Called from all edit handlers (Backspace, Enter, Space, Tab, Runes)

### `saveToHistory() error`
1. Gets history directory: `~/.promptstack/data/.history/`
2. Creates directory if needed with `os.MkdirAll(0755)`
3. Generates filename: `time.Now().Format("2006-01-02_15-04-05.md")`
4. Writes buffer content with `os.WriteFile(0644)`
5. Returns error on failure

### Status Bar Integration
- `saveStatus` displayed in status bar via `view.go`
- Error messages truncated to 40 characters
- "Saved" clears after 2 seconds via `tea.Tick(2s)`
- "Error" clears after 5 seconds via `tea.Tick(5s)`

## File System Structure

```
~/.promptstack/
├── config.yaml
├── debug.log
└── data/
    └── .history/
        ├── 2025-01-11_14-30-05.md
        ├── 2025-01-11_14-30-10.md
        └── 2025-01-11_14-30-15.md
```

## Error Handling Scenarios

### 1. Permission Denied
- User lacks write permission to history directory
- `saveToHistory()` returns `os.Permission` error
- Status bar shows "Error: permission denied" for 5 seconds
- Editor continues working (failed save doesn't affect editing)

### 2. Disk Full
- `os.WriteFile()` returns disk full error
- Status bar shows truncated error message
- User can continue editing (content not lost from memory)

### 3. Directory Creation Failed
- `os.MkdirAll()` fails (e.g., read-only filesystem)
- Error shown in status bar
- Auto-save attempts continue on next timer

### 4. Concurrent Saves
- Timer ID system prevents race conditions
- Only latest timer triggers save
- Multiple rapid edits result in single save

## Performance Characteristics

### Timing
- **Debounce delay**: 500ms (configurable constant)
- **Save operation**: <1ms for typical prompts (<10KB)
- **Status display**: "Saved" for 2000ms, "Error" for 5000ms
- **Timer resolution**: System-dependent, typically ~10ms

### Memory
- **Additional model state**: ~24 bytes (int + two strings)
- **No persistent buffers**: Content written directly to disk
- **Timer commands**: Minimal overhead, garbage collected after use

### Disk Usage
- **File size**: Exactly buffer content size (UTF-8 encoded)
- **Growth rate**: One file per save session (after 500ms inactivity)
- **No compression**: Raw markdown text

## Testing Strategy

### Unit Tests
1. **Timer scheduling**: Verify ID increments and command returned
2. **Message flow**: Simulate complete save cycle
3. **Stale timer**: Ensure old messages ignored
4. **Error handling**: Test error state transitions

### Integration Tests Needed
1. **File system**: Actual file creation in temp directory
2. **Concurrent edits**: Rapid typing produces single file
3. **Error simulation**: Mock filesystem failures

### Manual Tests
1. **Visual feedback**: Verify status bar updates
2. **File verification**: Check created files contain correct content
3. **Error scenarios**: Simulate disk full, permission denied

## Dependencies

### Internal
- `internal/editor.FileManager` (for modified flag clearing)
- `ui/theme` (for status bar styling)

### External
- `github.com/charmbracelet/bubbletea` (timer commands)
- `os`, `path/filepath`, `time` (standard library)

## Configuration Constants

```go
const (
    autoSaveDelay = 500 * time.Millisecond
    savedStatusDuration = 2 * time.Second
    errorStatusDuration = 5 * time.Second
    maxErrorDisplayLength = 40
)
```

## Integration with Other Milestones

### M4 (Basic Text Editor)
- Builds on existing workspace model
- Uses existing buffer for content access
- Integrates with existing status bar

### M6 (Undo/Redo) - FUTURE
- Auto-save could trigger undo snapshots
- History files could serve as persistent undo points
- Consider saving cursor position in metadata

### M15 (SQLite History) - FUTURE
- Could replace file-based history with database
- Same API, different storage backend
- Migration path: copy existing history files to database

## Security Considerations

### File Permissions
- History directory: `0755` (rwxr-xr-x)
- History files: `0644` (rw-r--r--)
- No sensitive data in prompts (user responsibility)

### Path Safety
- All paths constructed with `filepath.Join()`
- No user input in filenames (uses timestamp)
- Directory traversal not possible

### Memory Safety
- No buffers retained after write
- Error messages sanitized (truncated, not executed)

## Monitoring and Debugging

### Logging
- File operations not currently logged (could add debug logging)
- Errors captured in model state for UI display

### Diagnostics
- Status bar shows current save state
- No performance metrics collected (could add)

## Known Limitations

1. **No file locking**: Concurrent editor instances could overwrite
2. **No compression**: History files accumulate indefinitely
3. **No metadata**: Only content saved, no cursor position or viewport state
4. **Simple error handling**: Retry not implemented, just display error

## Future Enhancements

1. **Configurable delay**: User setting for auto-save frequency
2. **Save retention policy**: Auto-delete files older than N days
3. **Save confirmation sound**: Audio feedback option
4. **Cloud backup integration**: Sync history to remote storage
5. **Diff viewing**: Compare between auto-save points