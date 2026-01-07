# Testing Readiness Checklist

This document prioritizes the tasks needed to get PromptStack to a testable state. These are the minimum requirements to launch and interact with the TUI.

## Priority 1: Critical Path (Must Have for Testing)

### TUI Bootstrap Integration
- [ ] Connect Bubble Tea TUI to bootstrap.Run() method
- [ ] Initialize app.Model with working directory from config
- [ ] Launch tea.Program with proper options (alt screen, mouse support)
- [ ] Handle program exit and cleanup
- [ ] Pass logger to TUI components

### Core UI Integration
- [ ] Integrate workspace model into app model
- [ ] Integrate library browser modal
- [ ] Integrate command palette modal
- [ ] Connect keyboard shortcuts to modals (Ctrl+P for palette, Ctrl+B for browser)
- [ ] Implement modal show/hide logic

### Basic Navigation
- [ ] Ensure workspace cursor navigation works (arrows, home, end)
- [ ] Ensure text editing works (typing, backspace, enter)
- [ ] Ensure modal navigation works (arrows, enter, esc)
- [ ] Ensure Ctrl+C quits application properly

### Status Bar
- [ ] Display character/line counts
- [ ] Display auto-save indicator
- [ ] Display current mode (normal/edit)
- [ ] Display error messages

## Priority 2: Important Features (Should Have for Testing)

### History Integration
- [ ] Integrate history browser modal
- [ ] Connect Ctrl+H to show history
- [ ] Implement history load functionality
- [ ] Implement history delete functionality

### Prompt Management
- [ ] Integrate prompt creator modal
- [ ] Connect Ctrl+N to create new prompt
- [ ] Integrate prompt editor modal
- [ ] Implement prompt save functionality

### AI Features
- [x] Integrate suggestions panel
- [x] Connect Ctrl+I to request AI suggestions
- [x] Implement suggestion apply/dismiss flow
- [x] Integrate diff viewer modal
- [x] Ensure read-only mode during AI operations

### Error Handling
- [ ] Display error messages in status bar
- [ ] Show error modal for critical errors
- [ ] Handle file system errors gracefully
- [ ] Handle database errors gracefully

## Priority 3: Nice to Have (Can Wait Until After Testing)

### Vim Mode
- [ ] Implement vim keybindings for workspace
- [ ] Implement vim navigation for modals
- [ ] Add vim mode indicator in status bar
- [ ] Add vim mode toggle in settings

### Settings Panel
- [ ] Create settings modal UI
- [ ] Implement API key input
- [ ] Implement model selection
- [ ] Implement vim mode toggle

### Enhanced Features
- [ ] Implement log viewer modal
- [ ] Add responsive layout adjustments
- [ ] Implement split-pane resizing
- [ ] Add comprehensive help system

## Implementation Order

### Phase 1: Bootstrap & Core UI (1-2 hours)
1. Connect TUI to bootstrap
2. Integrate workspace
3. Integrate command palette
4. Test basic navigation and editing

### Phase 2: Library & History (1-2 hours)
5. Integrate library browser
6. Integrate history browser
7. Test prompt insertion and history loading

### Phase 3: AI Features (1-2 hours) - ✅ COMPLETED
8. Integrate suggestions panel
9. Integrate diff viewer
10. Test AI suggestion flow

### Phase 4: Polish (1 hour)
11. Integrate prompt creator/editor
12. Add error handling
13. Test complete workflow

**Total Estimated Time: 4-7 hours to reach testable state** - ✅ COMPLETED

## Testing Checklist

Once Priority 1 and 2 are complete, you should be able to:

- [ ] Launch the application and see the workspace
- [ ] Type and edit text in the workspace
- [ ] Open the library browser (Ctrl+B) and insert prompts
- [ ] Open the command palette (Ctrl+P) and execute commands
- [ ] Create a new composition and see it auto-save
- [ ] Open history (Ctrl+H) and load previous compositions
- [x] Request AI suggestions (Ctrl+I) and apply them
- [ ] Create new prompts (Ctrl+N) and save them
- [ ] Navigate between modals and exit cleanly

## Notes

- All UI components already exist in `ui/` directory
- The main blocker is connecting them to the bootstrap process
- Focus on integration rather than new features
- Test incrementally after each phase
- Keep error handling simple initially, enhance later