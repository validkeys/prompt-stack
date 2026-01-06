# PromptStack Implementation Tracking

This document tracks implementation progress across all milestones. Mark tasks as complete using `[x]` as you finish them.

---

## Milestone 1: Foundation (MVP Core)

### Application Bootstrap
- [ ] Create project structure (`cmd/promptstack/`, `internal/`, `ui/`)
- [ ] Implement config structure and YAML parsing
- [ ] Build first-run interactive setup wizard
- [ ] Add version detection and comparison logic
- [ ] Implement starter prompt extraction (go:embed)
- [ ] Create database initialization with schema
- [ ] Set up zap logger with rotation (10MB, keep 3)
- [ ] Add environment variable log level configuration

### Library Management
- [ ] Define prompt model with metadata types
- [ ] Implement library loader (filesystem scanning)
- [ ] Add YAML frontmatter parser
- [ ] Build in-memory index structure
- [ ] Implement index scoring algorithm (tags, keywords, usage)
- [ ] Create basic validation checks (YAML, placeholders, file size)
- [ ] Add validation results storage

### Basic Composition Workspace
- [ ] Create workspace Bubble Tea model
- [ ] Implement basic text editor with cursor management
- [ ] Add status bar with character/line counts
- [ ] Implement auto-save with debouncing (500ms-1s)
- [ ] Create timestamped markdown file on new composition
- [ ] Add visual auto-save indicator in status bar

### Library Browser
- [ ] Create library browser modal (Bubble Tea model)
- [ ] Integrate sahilm/fuzzy for filtering
- [ ] Implement flat list view with category labels
- [ ] Add preview pane showing prompt content
- [ ] Add keyboard navigation (arrows/vim j/k)
- [ ] Implement insert at cursor functionality
- [ ] Add insert on new line option

### Command Palette
- [ ] Create command palette modal
- [ ] Build command registry structure
- [ ] Implement fuzzy command filtering
- [ ] Add command execution dispatcher
- [ ] Register core commands (toggle AI, copy, save, etc.)

### Basic Error Handling
- [ ] Create error display components (status bar, modal)
- [ ] Implement graceful file read error handling
- [ ] Add config error recovery (prompt for reset)
- [ ] Log all errors to debug.log

---

## Milestone 2: Advanced Editing

### Placeholder System
- [ ] Implement placeholder parser (regex: `{{type:name}}`)
- [ ] Add placeholder validation (type, name, duplicates)
- [ ] Track placeholder positions in composition
- [ ] Implement Tab/Shift+Tab navigation between placeholders
- [ ] Create text placeholder editing mode (vim-style)
- [ ] Build list placeholder editing UI
- [ ] Add list item CRUD operations (add, edit, delete, navigate)
- [ ] Highlight active placeholder in editor

### Undo/Redo System
- [ ] Design undo action data structure
- [ ] Implement undo/redo stacks (100 action limit)
- [ ] Add smart action batching logic (typing, pauses, mode changes)
- [ ] Create undo/redo commands (Ctrl+Z, Ctrl+Shift+Z)
- [ ] Integrate with placeholder editing actions
- [ ] Add visual feedback for undo/redo boundaries

### Composition State
- [ ] Persist composition state to markdown
- [ ] Track dirty state (unsaved changes)
- [ ] Prevent data loss on exit with unsaved work
- [ ] Auto-save integration with undo history

---

## Milestone 3: History & File References

### SQLite Database
- [ ] Create database schema (compositions table)
- [ ] Add FTS5 virtual table for full-text search
- [ ] Create indexes (created_at, working_directory)
- [ ] Implement prepared statements for queries
- [ ] Enable SQLite WAL mode for concurrency
- [ ] Add database connection pooling

### History Management
- [ ] Implement history file creation (timestamped markdown)
- [ ] Add auto-save to both markdown and SQLite
- [ ] Create history browser UI (Bubble Tea model)
- [ ] Implement FTS5 search with query highlighting
- [ ] Add history listing with sorting (recent, directory)
- [ ] Create sync verification logic (markdown ↔ SQLite)
- [ ] Implement database rebuild with backup creation

### History Cleanup
- [ ] Define cleanup strategies (age-based, count-based, directory-based)
- [ ] Create cleanup UI with preview
- [ ] Implement batch deletion (markdown + SQLite)
- [ ] Add confirmation dialog for destructive operations
- [ ] Log cleanup operations

### File Reference System
- [ ] Implement file traversal respecting .gitignore
- [ ] Create fuzzy file finder modal
- [ ] Add multiple file selection support
- [ ] Implement YAML frontmatter title extraction
- [ ] Build batch title editor UI
- [ ] Add markdown link insertion at cursor
- [ ] Handle relative path generation from working directory

---

## Milestone 4: Prompt Management

### Prompt Creation
- [ ] Design guided creation wizard flow
- [ ] Create title input step (required validation)
- [ ] Add description input step (optional)
- [ ] Implement tags input with autocomplete
- [ ] Add category selection from existing folders
- [ ] Create initial content editor
- [ ] Generate filename from title (kebab-case)
- [ ] Save prompt with YAML frontmatter

### Prompt Editing
- [ ] Create prompt editor Bubble Tea model
- [ ] Implement edit mode (raw markdown with frontmatter)
- [ ] Integrate Glamour for preview mode
- [ ] Add Ctrl+P toggle between edit/preview
- [ ] Show mode indicator in status bar
- [ ] Implement explicit save (no auto-save)
- [ ] Handle save errors gracefully

### Library Validation
- [ ] Run validation on startup (silent unless errors)
- [ ] Add manual validation trigger via command palette
- [ ] Check for duplicate placeholder names
- [ ] Validate YAML frontmatter syntax
- [ ] Check for required fields (title)
- [ ] Enforce 1MB file size limit
- [ ] Detect invalid placeholder types/names
- [ ] Create validation results modal UI
- [ ] Show ⚠️ icon for invalid prompts in browser
- [ ] Block insertion of error-level prompts
- [ ] Allow insertion of warning-level prompts with indicator

---

## Milestone 5: AI Integration

### Claude API Integration
- [ ] Create Claude API client wrapper
- [ ] Implement retry logic for transient failures
- [ ] Add timeout configuration
- [ ] Handle rate limiting (429) with retry-after
- [ ] Implement auth error detection (401)
- [ ] Log all API requests/responses

### Context Selection
- [ ] Implement keyword extraction from composition
- [ ] Create library prompt scoring algorithm
- [ ] Calculate tag matches (+10 per tag)
- [ ] Add category bonus (+5)
- [ ] Weight keyword overlaps by frequency
- [ ] Bonus for recently used prompts (+3)
- [ ] Bonus for frequently used prompts (use_count)
- [ ] Sort and select top 3-5 within token budget

### Token Budget Management
- [ ] Detect Claude model context limit dynamically
- [ ] Implement token estimation for text
- [ ] Enforce 25% limit for composition content
- [ ] Enforce 15% limit for library context
- [ ] Show warning at 15% threshold
- [ ] Block suggestions at 25% threshold
- [ ] Display token estimate in status bar

### Suggestion System
- [ ] Define 6 suggestion types (recommendation, gap, formatting, contradiction, clarity, reformatting)
- [ ] Parse structured suggestion responses from Claude
- [ ] Create suggestions panel UI (scrollable list)
- [ ] Implement suggestion selection navigation
- [ ] Add "Apply" action (press 'a')
- [ ] Show "✨ AI is applying..." indicator
- [ ] Enter read-only mode during application

### Diff Application
- [ ] Generate unified diff using go-diff library
- [ ] Create diff viewer modal
- [ ] Parse structured edits from Claude
- [ ] Apply changes line-by-line
- [ ] Integrate with undo system (single action)
- [ ] Add accept/reject flow
- [ ] Unlock editor after completion
- [ ] Handle application errors with retry

---

## Milestone 6: Vim Mode

### Vim State Machine
- [ ] Design vim mode state structure (Normal/Insert/Visual)
- [ ] Implement mode transition logic
- [ ] Create mode indicator for status bar
- [ ] Add global vim mode toggle in config

### Keybinding System
- [ ] Define keybinding maps for Normal mode
- [ ] Define keybinding maps for Insert mode
- [ ] Define keybinding maps for Visual mode
- [ ] Create context-aware event router
- [ ] Map component-specific keybindings

### Universal Vim Support
- [ ] Add vim support to composition workspace
- [ ] Add vim navigation to library browser (j/k, /)
- [ ] Add vim navigation to command palette (j/k)
- [ ] Add vim navigation to AI suggestions panel (j/k)
- [ ] Add vim navigation to history browser (j/k)
- [ ] Add vim navigation to settings panel (j/k)
- [ ] Add vim editing to text inputs where applicable

---

## Milestone 7: Polish & Settings

### Settings Panel
- [ ] Create settings modal UI
- [ ] Add vim mode toggle (with restart warning)
- [ ] Add Claude API key input (masked)
- [ ] Add Claude model selection dropdown
- [ ] Implement immediate persistence to config.yaml
- [ ] Validate settings before saving
- [ ] Show success/error feedback

### Enhanced Status Bar
- [ ] Add auto-save indicator (transient)
- [ ] Display character and line counts
- [ ] Show token estimate when AI panel open
- [ ] Add vim mode indicator
- [ ] Add edit/preview mode indicator
- [ ] Display notifications and warnings
- [ ] Add visual styling with Lipgloss

### Responsive Layout
- [ ] Implement split-pane layout (workspace + AI panel)
- [ ] Add resizable divider
- [ ] Handle narrow terminal gracefully (80 columns)
- [ ] Test layout from 80 to 200+ columns
- [ ] Ensure text wrapping works correctly
- [ ] Adapt panel sizes based on terminal width

### Log Viewer
- [ ] Create log viewer modal (Ctrl+L)
- [ ] Tail last N lines from debug.log
- [ ] Add auto-refresh option
- [ ] Implement log level filtering
- [ ] Add search/highlighting in logs

### Comprehensive Error Messages
- [ ] Improve file system error messages
- [ ] Enhance database error messages
- [ ] Clarify API error messages
- [ ] Add contextual help for common errors
- [ ] Provide actionable next steps in errors

### Startup Validation
- [ ] Verify config file integrity on startup
- [ ] Check database accessibility
- [ ] Validate library directory structure
- [ ] Detect and report configuration issues
- [ ] Offer fixes for common problems

---

## Milestone 8: Testing & Documentation

### Unit Tests
- [ ] Set up test framework structure
- [ ] Create test fixtures (prompts, compositions)
- [ ] Write tests for library loader
- [ ] Write tests for placeholder parser
- [ ] Write tests for undo/redo stack
- [ ] Write tests for context selection algorithm
- [ ] Write tests for validation checks
- [ ] Write tests for token estimation
- [ ] Achieve 80%+ coverage on critical paths

### Integration Tests
- [ ] Test library loading and indexing (100+ prompts)
- [ ] Test history database operations
- [ ] Test composition save/load cycle
- [ ] Test placeholder editing flow
- [ ] Test AI context selection with real library
- [ ] Test database sync verification and rebuild

### TUI Tests
- [ ] Use Bubble Tea testing utilities
- [ ] Test workspace model updates
- [ ] Test modal state transitions
- [ ] Test keyboard navigation
- [ ] Test vim mode transitions
- [ ] Verify state consistency

### End-to-End Scenarios
- [ ] First launch initialization
- [ ] Create composition, insert prompts, fill placeholders
- [ ] AI suggestion request and acceptance
- [ ] History search and load
- [ ] Prompt creation and validation
- [ ] Settings update and persistence
- [ ] Database corruption recovery

### User Documentation
- [ ] Write README with installation instructions
- [ ] Document first-run setup process
- [ ] Create quick start guide
- [ ] Document all keyboard shortcuts
- [ ] Explain vim mode usage
- [ ] Document placeholder syntax
- [ ] Explain AI suggestion workflow
- [ ] Provide troubleshooting guide

### Developer Documentation
- [ ] Document architecture and design decisions
- [ ] Create contribution guidelines
- [ ] Document code organization
- [ ] Add inline code comments for complex logic
- [ ] Document build process
- [ ] Explain testing approach

---

## Milestone 9: Build & Distribution

### Build Process
- [ ] Create build script for macOS Intel
- [ ] Create build script for macOS ARM
- [ ] Embed starter prompts using go:embed
- [ ] Test builds on both architectures
- [ ] Verify binary size is reasonable (<50MB)
- [ ] Test cross-compilation

### Version Management
- [ ] Implement version detection in config
- [ ] Create upgrade detection logic
- [ ] Implement new starter prompt extraction on upgrade
- [ ] Skip existing files during upgrade
- [ ] Update version in config after upgrade
- [ ] Test upgrade from v1.0 → v1.1

### Release Preparation
- [ ] Create CHANGELOG.md
- [ ] Write release notes
- [ ] Tag version in git
- [ ] Build release binaries
- [ ] Test release binaries on clean systems
- [ ] Create GitHub Release with artifacts

### Distribution
- [ ] Upload binaries to GitHub Releases
- [ ] Write installation instructions
- [ ] Create homebrew formula (optional)
- [ ] Announce release
- [ ] Monitor for bug reports

---

## Progress Summary

- **Milestone 1:** 0/XX tasks complete
- **Milestone 2:** 0/XX tasks complete
- **Milestone 3:** 0/XX tasks complete
- **Milestone 4:** 0/XX tasks complete
- **Milestone 5:** 0/XX tasks complete
- **Milestone 6:** 0/XX tasks complete
- **Milestone 7:** 0/XX tasks complete
- **Milestone 8:** 0/XX tasks complete
- **Milestone 9:** 0/XX tasks complete

---

**Last Updated:** 2026-01-06
**Status:** Ready for implementation
