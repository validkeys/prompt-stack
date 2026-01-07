# PromptStack - Granular Implementation Milestones

**Strategy:** Build incrementally with testing at each milestone. Order prioritizes seeing the application come to life with integrated features early, following tighter integration loops.

**Testing Philosophy:** Each milestone ends with explicit verification that the feature works in isolation and integrates with previous milestones.

---

## Milestone 1: Bootstrap & Config ✅ **COMPLETED**
**Goal:** Initialize application foundation

**Deliverables:**
- Config structure at `~/.promptstack/config.yaml` ✅
- First-run interactive setup wizard ✅
- Logging setup with zap ✅
- Version tracking ✅

**Test Criteria:**

### Functional Requirements
- [x] App launches without errors on fresh install
- [x] Config file created at `~/.promptstack/config.yaml` with correct structure
- [x] Setup wizard prompts for API key and preferences
  - [x] Validates API key format before accepting
  - [x] Shows error message for invalid keys
  - [x] Allows re-entry of invalid fields
- [x] Logs written to `~/.promptstack/debug.log`
  - [x] Log file created with correct permissions (0600)
  - [x] Log entries include timestamp and level
  - [ ] Log rotation works at 10MB limit (not manually tested)
- [x] Version stored in config and compared on startup

### Integration Requirements
- [x] Config loading integrates with logging system
- [x] Setup wizard integrates with config persistence
- [ ] Version tracking integrates with starter prompt extraction (N/A - no starter prompts yet)

### Edge Cases & Error Handling
- [x] Handle missing config directory (create automatically)
- [ ] Handle corrupted config file (show error, offer reset) (not manually tested)
- [x] Handle invalid API key format (reject with clear message)
- [ ] Handle read-only filesystem (show error, exit gracefully) (not manually tested)
- [ ] Handle interrupted setup wizard (resume on next launch) (not manually tested)

### Performance Requirements
- [x] App startup time <500ms on fresh install
- [x] Config file read/write <50ms
- [x] Log file write <10ms per entry

### User Experience Requirements
- [x] Setup wizard provides clear instructions
- [x] Error messages are actionable and specific
- [ ] Progress indicators shown during initialization (not manually tested)
- [ ] Keyboard navigation works in setup wizard (not manually tested)

**Files:** [`internal/config/config.go`], [`internal/setup/wizard.go`], [`internal/logging/logger.go`], [`cmd/promptstack/main.go`]

---

## Milestone 2: Basic TUI Shell
**Goal:** Render functional TUI with quit handling

**Deliverables:**
- Root Bubble Tea model
- Basic status bar
- Keyboard input handling
- Clean quit (Ctrl+C, q)

**Test Criteria:**

### Functional Requirements
- [ ] App renders in terminal without errors
- [ ] Status bar visible at bottom of screen
- [ ] Responds to keyboard input
- [ ] Quits cleanly with Ctrl+C
- [ ] Quits cleanly with 'q' key
- [ ] Terminal restored to normal state on exit

### Integration Requirements
- [ ] TUI shell integrates with config system
- [ ] Status bar integrates with logging system
- [ ] Keyboard input integrates with event system

### Edge Cases & Error Handling
- [ ] Handle terminal resize during runtime
- [ ] Handle terminal too small (<80 columns)
- [ ] Handle rapid quit commands
- [ ] Handle interrupted quit sequence
- [ ] Handle display errors (graceful degradation)

### Performance Requirements
- [ ] Initial render <100ms
- [ ] Frame rate >=30fps during input
- [ ] Memory footprint <50MB

### User Experience Requirements
- [ ] Status bar shows app name and version
- [ ] Clear visual feedback on keypress
- [ ] Smooth quit animation (if applicable)
- [ ] No flickering or visual artifacts

**Files:** [`ui/app/model.go`], [`ui/statusbar/model.go`]

---

## Milestone 3: File I/O Foundation
**Goal:** Read/write markdown with YAML frontmatter

**Deliverables:**
- Markdown file reader
- YAML frontmatter parser
- Markdown file writer
- Error handling for file operations

**Test Criteria:**

### Functional Requirements
- [ ] Read markdown file with frontmatter
- [ ] Parse frontmatter into struct
- [ ] Write markdown with frontmatter
- [ ] Preserve markdown content exactly
- [ ] Handle files without frontmatter

### Integration Requirements
- [ ] File I/O integrates with logging system
- [ ] Markdown parser integrates with prompt structure
- [ ] Error handling integrates with error handler

### Edge Cases & Error Handling
- [ ] Handle missing files (return error, don't crash)
- [ ] Handle malformed YAML (return error, don't crash)
- [ ] Handle empty files
- [ ] Handle files with only frontmatter
- [ ] Handle files with only markdown (no frontmatter)
- [ ] Handle permission denied errors
- [ ] Handle disk full errors
- [ ] Handle invalid UTF-8 encoding

### Performance Requirements
- [ ] Read 1MB file <100ms
- [ ] Write 1MB file <100ms
- [ ] Parse frontmatter <10ms
- [ ] Handle 1000+ files efficiently

### User Experience Requirements
- [ ] Clear error messages for file operations
- [ ] Progress indicators for large files
- [ ] No data loss on write errors
- [ ] Preserve file permissions

**Files:** [`internal/prompt/storage.go`], [`internal/files/markdown.go`]

---

## Milestone 4: Basic Text Editor
**Goal:** Type and edit text in workspace

**Deliverables:**
- Text input component
- Cursor movement (arrows, home, end)
- Character and line counting
- Display in main workspace area

**Test Criteria:**

### Functional Requirements
- [ ] Type text and see it appear immediately
- [ ] Move cursor with arrow keys (up, down, left, right)
- [ ] Home key moves to start of line
- [ ] End key moves to end of line
- [ ] Character count updates in status bar
- [ ] Line count updates in status bar
- [ ] Backspace deletes character before cursor
- [ ] Delete key deletes character at cursor

### Integration Requirements
- [ ] Editor integrates with TUI shell
- [ ] Status bar integrates with editor state
- [ ] Keyboard input integrates with editor model

### Edge Cases & Error Handling
- [ ] Handle empty editor
- [ ] Handle very long lines (>1000 characters)
- [ ] Handle very long documents (>10000 lines)
- [ ] Handle rapid typing (no dropped characters)
- [ ] Handle cursor at boundaries (start/end of file)
- [ ] Handle special characters (tabs, newlines, unicode)

### Performance Requirements
- [ ] Typing latency <10ms
- [ ] Cursor movement <5ms
- [ ] Render 10000 lines <50ms
- [ ] Memory usage scales linearly with document size

### User Experience Requirements
- [ ] Smooth cursor movement
- [ ] No lag during typing
- [ ] Clear visual feedback for cursor position
- [ ] Line numbers visible (if applicable)
- [ ] Word wrap (if applicable)

**Files:** [`ui/workspace/model.go`], [`internal/editor/editor.go`]

---

## Milestone 5: Auto-save
**Goal:** Automatically save composition to timestamped files

**Deliverables:**
- Debounced auto-save (500ms-1s)
- Timestamped filename generation
- Save to `~/.promptstack/data/.history/`
- Visual feedback in status bar

**Test Criteria:**

### Functional Requirements
- [ ] Type text, wait 500ms-1s, verify file created
- [ ] Filename format: `YYYY-MM-DD_HH-MM-SS.md`
- [ ] Status bar shows "Saving..." during save
- [ ] Status bar shows "Saved" after save completes
- [ ] Multiple edits batch into single save
- [ ] File contains current editor content
- [ ] Auto-save triggers on any edit

### Integration Requirements
- [ ] Auto-save integrates with editor model
- [ ] Status bar integrates with auto-save state
- [ ] File I/O integrates with auto-save system

### Edge Cases & Error Handling
- [ ] Handle empty editor (save empty file)
- [ ] Handle very large documents (>1MB)
- [ ] Handle rapid edits (debounce correctly)
- [ ] Handle save errors (show error, retry)
- [ ] Handle disk full errors
- [ ] Handle permission denied errors
- [ ] Handle concurrent saves (queue properly)

### Performance Requirements
- [ ] Auto-save trigger delay: 500ms-1s
- [ ] Save operation <100ms for 1MB file
- [ ] Debounce doesn't miss edits
- [ ] No UI lag during save

### User Experience Requirements
- [ ] Clear "Saving..." indicator
- [ ] Clear "Saved" confirmation
- [ ] Error messages are actionable
- [ ] No interruption to typing during save
- [ ] Visual feedback for save status

**Files:** [`internal/editor/autosave.go`], [`internal/history/storage.go`]

---

## Milestone 6: Undo/Redo
**Goal:** Undo and redo editing actions

**Deliverables:**
- Undo stack (100 levels)
- Redo stack
- Smart batching (continuous typing = one action)
- Keyboard shortcuts (Ctrl+Z, Ctrl+Y)

**Test Criteria:**

### Functional Requirements
- [ ] Type text, press Ctrl+Z, text removed
- [ ] Press Ctrl+Y, text restored
- [ ] Continuous typing batched as one undo action
- [ ] Cursor movement breaks batching
- [ ] Undo/redo at boundaries shows feedback
- [ ] Stack limited to 100 actions
- [ ] Redo stack cleared on new edit
- [ ] Undo/redo works with all edit types

### Integration Requirements
- [ ] Undo/redo integrates with editor model
- [ ] Undo/redo integrates with auto-save
- [ ] Keyboard shortcuts integrate with input system

### Edge Cases & Error Handling
- [ ] Handle undo at start of document (show feedback)
- [ ] Handle redo at end of history (show feedback)
- [ ] Handle undo after auto-save
- [ ] Handle undo with very large documents
- [ ] Handle rapid undo/redo commands
- [ ] Handle undo with special characters
- [ ] Handle undo with cursor at different positions

### Performance Requirements
- [ ] Undo operation <10ms
- [ ] Redo operation <10ms
- [ ] Memory usage scales with stack size
- [ ] No UI lag during undo/redo

### User Experience Requirements
- [ ] Clear feedback when undo/redo not available
- [ ] Smooth undo/redo transitions
- [ ] No visual artifacts during undo/redo
- [ ] Cursor position restored correctly
- [ ] Selection restored correctly (if applicable)

**Files:** [`internal/editor/undo.go`]

---

## Milestone 7: Prompt Source Interface & Filesystem Implementation
**Goal:** Implement prompt source abstraction with filesystem as first source

**Deliverables:**
- PromptSource interface in `internal/library/source.go`
- FilesystemSource implementation in `internal/library/filesystem.go`
- PromptCache interface in `internal/library/cache.go`
- MemoryCache implementation in `internal/library/cache.go`
- Scan `~/.promptstack/data/` directory
- Parse YAML frontmatter from each prompt
- Load into in-memory collection
- Extract metadata (title, tags, category, description)
- Unit tests for interfaces and implementations
- Integration tests with filesystem

**Test Criteria:**

### Functional Requirements
- [ ] PromptSource interface defined with required methods
- [ ] FilesystemSource implements PromptSource interface
- [ ] PromptCache interface defined with required methods
- [ ] MemoryCache implements PromptCache interface
- [ ] Load sample prompts from directory
- [ ] Verify prompt count matches files
- [ ] Metadata parsed correctly (title, tags, category, description)
- [ ] Category derived from folder structure
- [ ] Invalid prompts loaded (not excluded)
- [ ] Empty directories handled gracefully
- [ ] Nested directories scanned recursively
- [ ] Cache stores and retrieves prompts correctly
- [ ] Cache invalidation works correctly

### Integration Requirements
- [ ] FilesystemSource integrates with config system
- [ ] Metadata extraction integrates with prompt structure
- [ ] Error handling integrates with logging system
- [ ] Cache integrates with FilesystemSource

### Edge Cases & Error Handling
- [ ] Handle missing library directory (create automatically)
- [ ] Handle empty library directory
- [ ] Handle files without frontmatter
- [ ] Handle malformed YAML frontmatter
- [ ] Handle duplicate filenames
- [ ] Handle very large libraries (1000+ prompts)
- [ ] Handle permission denied errors
- [ ] Handle circular symlinks (if applicable)
- [ ] Handle cache misses
- [ ] Handle cache invalidation during load

### Performance Requirements
- [ ] Startup time <500ms for 100 prompts
- [ ] Startup time <2s for 1000 prompts
- [ ] Memory usage scales linearly with prompt count
- [ ] File I/O operations are concurrent where possible
- [ ] Cache hit time <1ms
- [ ] Cache miss time <10ms

### User Experience Requirements
- [ ] Loading indicator shown during startup
- [ ] Error messages are specific to file/issue
- [ ] Progress indicator for large libraries
- [ ] No UI blocking during load

**Files:** [`internal/library/source.go`], [`internal/library/filesystem.go`], [`internal/library/cache.go`]

---

## Milestone 8: Library Browser UI
**Goal:** Browse and preview prompts

**Deliverables:**
- Modal overlay with prompt list
- Keyboard navigation (arrows or j/k)
- Preview pane showing prompt content
- Category labels in list
- Open with keyboard shortcut

**Test Criteria:**

### Functional Requirements
- [ ] Press shortcut, browser opens
- [ ] See list of all prompts with categories
- [ ] Navigate with arrow keys (up/down)
- [ ] Navigate with j/k keys
- [ ] Preview pane shows selected prompt
- [ ] Esc closes browser
- [ ] Returns to editor on close
- [ ] Category labels displayed in list

### Integration Requirements
- [ ] Browser integrates with library loader
- [ ] Preview integrates with prompt structure
- [ ] Keyboard shortcuts integrate with input system

### Edge Cases & Error Handling
- [ ] Handle empty library (show empty state)
- [ ] Handle very long prompt titles (truncate)
- [ ] Handle very long descriptions (truncate in preview)
- [ ] Handle prompts with missing metadata
- [ ] Handle rapid navigation (no lag)
- [ ] Handle browser open when already open

### Performance Requirements
- [ ] Browser open <100ms
- [ ] Navigation latency <10ms
- [ ] Preview render <50ms
- [ ] Handle 1000+ prompts smoothly

### User Experience Requirements
- [ ] Smooth scrolling through list
- [ ] Clear visual indication of selected item
- [ ] Preview pane readable and well-formatted
- [ ] Category labels visually distinct
- [ ] Empty state message is helpful
- [ ] Keyboard shortcuts discoverable

**Files:** [`ui/browser/model.go`]

---

## Milestone 9: Fuzzy Search in Library
**Goal:** Filter library as user types

**Deliverables:**
- Text input at top of browser
- Real-time fuzzy filtering
- Search across title, tags, description, category
- Rank by match quality

**Test Criteria:**

### Functional Requirements
- [ ] Type in search box
- [ ] List filters in real-time (<10ms)
- [ ] Exact matches ranked higher
- [ ] Partial matches shown
- [ ] Clear search restores full list
- [ ] Search works with 100+ prompts
- [ ] Case-insensitive search
- [ ] Search across title, tags, description, category

### Integration Requirements
- [ ] Search integrates with library loader
- [ ] Search results integrate with browser UI
- [ ] Search input integrates with input system

### Edge Cases & Error Handling
- [ ] Handle empty search (show all)
- [ ] Handle no results (show empty state)
- [ ] Handle special characters in search
- [ ] Handle very long search queries
- [ ] Handle rapid typing (debounce properly)
- [ ] Handle unicode characters in search

### Performance Requirements
- [ ] Search latency <10ms for 100 prompts
- [ ] Search latency <50ms for 1000 prompts
- [ ] Debounce delay: 100-200ms
- [ ] No UI lag during search

### User Experience Requirements
- [ ] Search results update smoothly
- [ ] Clear indication of search active
- [ ] Empty search results message is helpful
- [ ] Search input clearly visible
- [ ] Clear search button (Esc or Ctrl+C)

**Files:** [`internal/library/search.go`], [`ui/browser/model.go`]

---

## Milestone 10: Prompt Insertion
**Goal:** Insert selected prompt into editor

**Deliverables:**
- Insert at cursor position
- Insert on new line option
- Close browser after insertion
- Update editor content

**Test Criteria:**

### Functional Requirements
- [ ] Select prompt, press Enter
- [ ] Prompt content appears at cursor
- [ ] Browser closes automatically
- [ ] Cursor positioned after inserted content
- [ ] Undo removes entire insertion
- [ ] Works with empty editor
- [ ] Works with existing content
- [ ] Insert on new line option works

### Integration Requirements
- [ ] Insertion integrates with editor model
- [ ] Browser close integrates with TUI shell
- [ ] Undo integration with undo stack

### Edge Cases & Error Handling
- [ ] Handle insertion at start of file
- [ ] Handle insertion at end of file
- [ ] Handle insertion in middle of line
- [ ] Handle very long prompts (>10000 characters)
- [ ] Handle prompts with special characters
- [ ] Handle insertion with selection (replace or insert)
- [ ] Handle insertion with undo stack full

### Performance Requirements
- [ ] Insertion operation <50ms
- [ ] Browser close <50ms
- [ ] Handle 10000+ character prompts smoothly

### User Experience Requirements
- [ ] Smooth insertion animation (if applicable)
- [ ] Clear visual feedback on insertion
- [ ] Cursor position intuitive after insertion
- [ ] Browser close is seamless
- [ ] No visual artifacts after insertion

**Files:** [`ui/browser/model.go`], [`ui/workspace/model.go`]

---

## Milestone 11: Placeholder Parser
**Goal:** Detect placeholders in composition

**Deliverables:**
- Regex-based parser for `{{text:name}}` and `{{list:name}}`
- Placeholder validation
- Position tracking
- Real-time parsing as user types

**Test Criteria:**

### Functional Requirements
- [ ] Parse text with `{{text:example}}`
- [ ] Parse text with `{{list:items}}`
- [ ] Detect multiple placeholders
- [ ] Track start/end positions
- [ ] Validate placeholder syntax
- [ ] Ignore malformed placeholders
- [ ] Handle duplicate names (validation error)
- [ ] Parse in real-time as user types

### Integration Requirements
- [ ] Parser integrates with editor model
- [ ] Validation integrates with error handler
- [ ] Position tracking integrates with cursor

### Edge Cases & Error Handling
- [ ] Handle nested placeholders (reject or handle)
- [ ] Handle unclosed placeholders
- [ ] Handle empty placeholder names
- [ ] Handle special characters in names
- [ ] Handle very long placeholder names
- [ ] Handle placeholders at document boundaries
- [ ] Handle rapid typing (parse efficiently)

### Performance Requirements
- [ ] Parse 1000-line document <50ms
- [ ] Real-time parsing <10ms per keystroke
- [ ] Memory usage scales with document size
- [ ] No UI lag during typing

### User Experience Requirements
- [ ] Visual highlighting of placeholders
- [ ] Clear indication of invalid placeholders
- [ ] Smooth parsing during typing
- [ ] No flickering or visual artifacts

**Files:** [`internal/editor/placeholder.go`]

---

## Milestone 12: Placeholder Navigation
**Goal:** Tab between placeholders

**Deliverables:**
- Tab key moves to next placeholder
- Shift+Tab moves to previous placeholder
- Visual highlighting of active placeholder
- Wrap around at boundaries

**Test Criteria:**

### Functional Requirements
- [ ] Insert prompt with placeholders
- [ ] Press Tab, first placeholder highlighted
- [ ] Press Tab again, moves to next
- [ ] Press Shift+Tab, moves to previous
- [ ] At last placeholder, Tab wraps to first
- [ ] At first placeholder, Shift+Tab wraps to last
- [ ] Visual indicator shows active placeholder
- [ ] Cursor moves to placeholder position

### Integration Requirements
- [ ] Navigation integrates with parser
- [ ] Highlighting integrates with editor UI
- [ ] Keyboard shortcuts integrate with input system

### Edge Cases & Error Handling
- [ ] Handle document with no placeholders (no action)
- [ ] Handle document with single placeholder
- [ ] Handle rapid Tab/Shift+Tab presses
- [ ] Handle navigation during editing
- [ ] Handle placeholders added/removed during navigation

### Performance Requirements
- [ ] Navigation latency <10ms
- [ ] Highlighting update <10ms
- [ ] No UI lag during navigation

### User Experience Requirements
- [ ] Smooth transition between placeholders
- [ ] Clear visual indication of active placeholder
- [ ] Cursor position intuitive
- [ ] No visual artifacts during navigation

**Files:** [`internal/editor/placeholder.go`], [`ui/workspace/model.go`]

---

## Milestone 13: Text Placeholder Editing
**Goal:** Edit text placeholder values

**Deliverables:**
- Enter edit mode with 'i' or 'a'
- Type to replace placeholder content
- Esc to exit edit mode
- Update composition with new value

**Test Criteria:**

### Functional Requirements
- [ ] Tab to text placeholder
- [ ] Press 'i', enter edit mode
- [ ] Type new value
- [ ] Press Esc, exit edit mode
- [ ] Placeholder replaced with typed value
- [ ] Undo restores original placeholder
- [ ] Works with multiple text placeholders
- [ ] 'a' key also enters edit mode

### Integration Requirements
- [ ] Edit mode integrates with navigation
- [ ] Value replacement integrates with editor
- [ ] Undo integration with undo stack

### Edge Cases & Error Handling
- [ ] Handle empty value (replace with empty string)
- [ ] Handle very long values (>1000 characters)
- [ ] Handle special characters in value
- [ ] Handle rapid typing
- [ ] Handle Esc during typing (cancel edit)
- [ ] Handle navigation during edit (cancel or save)

### Performance Requirements
- [ ] Enter edit mode <10ms
- [ ] Exit edit mode <10ms
- [ ] Value replacement <10ms
- [ ] No UI lag during typing

### User Experience Requirements
- [ ] Clear indication of edit mode
- [ ] Smooth transition to/from edit mode
- [ ] Cursor position intuitive in edit mode
- [ ] Visual feedback during editing

**Files:** [`internal/editor/placeholder.go`], [`ui/workspace/model.go`]

---

## Milestone 14: List Placeholder Editing
**Goal:** Manage list items in list placeholders

**Deliverables:**
- List editing UI modal
- Navigate items with Up/Down
- Edit item with 'e'
- Delete item with 'd'
- Add item with 'n' or 'o'
- Convert to markdown bullet list

**Test Criteria:**

### Functional Requirements
- [ ] Tab to list placeholder
- [ ] Press 'i', list editor opens
- [ ] See existing items (if any)
- [ ] Press 'n', add new item
- [ ] Type item text
- [ ] Press 'e' on item, edit inline
- [ ] Press 'd' on item, item deleted
- [ ] Press Esc, list converted to markdown bullets
- [ ] Undo restores original placeholder
- [ ] 'o' key also adds new item

### Integration Requirements
- [ ] List editor integrates with navigation
- [ ] Markdown conversion integrates with editor
- [ ] Undo integration with undo stack

### Edge Cases & Error Handling
- [ ] Handle empty list (no items)
- [ ] Handle very long item text (>500 characters)
- [ ] Handle very long lists (>100 items)
- [ ] Handle empty item text
- [ ] Handle special characters in items
- [ ] Handle rapid item operations
- [ ] Handle Esc during editing (cancel)

### Performance Requirements
- [ ] Open list editor <50ms
- [ ] Add item <10ms
- [ ] Edit item <10ms
- [ ] Delete item <10ms
- [ ] Convert to markdown <50ms
- [ ] Handle 100+ items smoothly

### User Experience Requirements
- [ ] Clear list editor UI
- [ ] Smooth item navigation
- [ ] Clear visual feedback for selected item
- [ ] Intuitive keyboard shortcuts
- [ ] Smooth markdown conversion

**Files:** [`internal/editor/listplaceholder.go`], [`ui/workspace/model.go`]

---

## Milestone 15: Repository Pattern & SQLite Implementation
**Goal:** Implement repository abstraction with SQLite as first backend

**Deliverables:**
- CompositionRepository interface in `internal/storage/repository.go`
- SQLiteRepository implementation in `internal/storage/sqlite.go`
- Repository factory in `internal/storage/factory.go`
- Database schema with FTS5
- Create at `~/.promptstack/data/history.db`
- Basic CRUD operations
- Index on created_at and working_directory
- Unit tests for interface and implementation
- Integration tests with SQLite

**Test Criteria:**

### Functional Requirements
- [ ] CompositionRepository interface defined with required methods
- [ ] SQLiteRepository implements CompositionRepository interface
- [ ] Repository factory creates correct repository based on config
- [ ] Database file created at `~/.promptstack/data/history.db` on first run
- [ ] Schema matches specification (compositions table with FTS5)
- [ ] Insert composition record with all fields
- [ ] Query composition by ID
- [ ] Full-text search works on content field
- [ ] Update existing composition record
- [ ] Delete composition record
- [ ] Index created on created_at column
- [ ] Index created on working_directory column

### Integration Requirements
- [ ] Repository factory integrates with config system
- [ ] Repository factory reads storage field from config
- [ ] Database initialization integrates with config system
- [ ] Database operations integrate with logging system
- [ ] Error handling integrates with error handler

### Edge Cases & Error Handling
- [ ] Handle missing data directory (create automatically)
- [ ] Handle corrupted database file (show error, offer rebuild)
- [ ] Handle database locked (retry with backoff)
- [ ] Handle disk full errors
- [ ] Handle permission denied errors
- [ ] Handle concurrent database access
- [ ] Handle very large records (>1MB content)

### Performance Requirements
- [ ] Database initialization <100ms
- [ ] Insert operation <10ms
- [ ] Query by ID <5ms
- [ ] Full-text search <50ms for 1000 records
- [ ] Update operation <10ms
- [ ] Delete operation <10ms

### User Experience Requirements
- [ ] Clear error messages for database failures
- [ ] Progress indicators for long operations
- [ ] No UI blocking during database operations
- [ ] Graceful degradation if database unavailable

**Files:** [`internal/storage/repository.go`], [`internal/storage/sqlite.go`], [`internal/storage/factory.go`]

---

## Milestone 16: History Sync
**Goal:** Save to both markdown and SQLite

**Deliverables:**
- Auto-save writes to markdown file
- Auto-save inserts/updates SQLite record
- Sync verification on startup
- Offer rebuild if mismatch detected

**Test Criteria:**

### Functional Requirements
- [ ] Edit composition, auto-save triggers
- [ ] Markdown file updated with current content
- [ ] SQLite record updated with current content
- [ ] Both markdown and SQLite contain identical content
- [ ] Restart app, sync verification runs automatically
- [ ] Sync verification passes for matching records
- [ ] Delete markdown file, sync detects mismatch
- [ ] Delete SQLite record, sync detects mismatch
- [ ] Rebuild option offered when mismatch detected
- [ ] Rebuild recreates missing records from markdown files
- [ ] Rebuild recreates missing markdown files from SQLite records

### Integration Requirements
- [ ] Sync integrates with auto-save system
- [ ] Sync integrates with database operations
- [ ] Sync integrates with file I/O operations
- [ ] Sync verification integrates with startup sequence
- [ ] Rebuild integrates with library loader

### Edge Cases & Error Handling
- [ ] Handle markdown file missing (rebuild from SQLite)
- [ ] Handle SQLite record missing (rebuild from markdown)
- [ ] Handle content mismatch (offer rebuild)
- [ ] Handle corrupted markdown file (show error, skip)
- [ ] Handle corrupted SQLite record (show error, skip)
- [ ] Handle concurrent sync operations (queue properly)
- [ ] Handle sync during rebuild (block or queue)
- [ ] Handle very large compositions (>1MB)

### Performance Requirements
- [ ] Sync operation <50ms per composition
- [ ] Sync verification <100ms for 100 compositions
- [ ] Rebuild operation <1s for 100 compositions
- [ ] No UI lag during sync operations
- [ ] Debounce sync operations (batch multiple edits)

### User Experience Requirements
- [ ] Clear indication of sync status in status bar
- [ ] Clear error messages for sync failures
- [ ] Clear rebuild option with explanation
- [ ] Progress indicator for rebuild operations
- [ ] No interruption to editing during sync

**Files:** [`internal/history/sync.go`], [`internal/history/storage.go`]

---

## Milestone 17: History Browser
**Goal:** Browse and load past compositions

**Deliverables:**
- Modal with list of compositions
- Sort by date (newest first)
- Show preview of content
- Load selected composition into editor
- Search by content (FTS5)

**Test Criteria:**

### Functional Requirements
- [ ] Open history browser with keyboard shortcut
- [ ] See list of all past compositions
- [ ] List sorted by date, newest first
- [ ] Preview shows first few lines of content
- [ ] Select composition, press Enter to load
- [ ] Editor content replaced with selected composition
- [ ] Search finds compositions by keyword (FTS5)
- [ ] Search works with 100+ history entries
- [ ] Search results ranked by relevance
- [ ] Clear search restores full list
- [ ] Esc closes history browser

### Integration Requirements
- [ ] History browser integrates with database operations
- [ ] History browser integrates with editor model
- [ ] Search integrates with FTS5 full-text search
- [ ] Load operation integrates with undo stack
- [ ] Keyboard shortcuts integrate with input system

### Edge Cases & Error Handling
- [ ] Handle empty history (show empty state)
- [ ] Handle very long composition titles (truncate)
- [ ] Handle very long previews (truncate)
- [ ] Handle very long history (>1000 entries)
- [ ] Handle search with no results (show empty state)
- [ ] Handle rapid navigation (no lag)
- [ ] Handle browser open when already open
- [ ] Handle database errors (show error, close browser)

### Performance Requirements
- [ ] Browser open <100ms
- [ ] List load <50ms for 100 entries
- [ ] List load <200ms for 1000 entries
- [ ] Navigation latency <10ms
- [ ] Search latency <50ms for 1000 entries
- [ ] Load composition <100ms
- [ ] Preview render <50ms

### User Experience Requirements
- [ ] Smooth scrolling through list
- [ ] Clear visual indication of selected item
- [ ] Preview pane readable and well-formatted
- [ ] Empty state message is helpful
- [ ] Search input clearly visible
- [ ] Keyboard shortcuts discoverable
- [ ] No visual artifacts during navigation

**Files:** [`ui/history/model.go`], [`internal/history/listing.go`]

---

## Milestone 18: Command Registry
**Goal:** Define and execute commands

**Deliverables:**
- Command registry structure
- Register commands with names and handlers
- Execute command by name
- Return results or errors

**Test Criteria:**

### Functional Requirements
- [ ] Register test command with name and handler
- [ ] Execute command by name
- [ ] Handler called with correct arguments
- [ ] Return value captured and returned
- [ ] Error handling works (errors propagated)
- [ ] List all registered commands
- [ ] Check if command exists
- [ ] Unregister command (if needed)
- [ ] Handle command with no return value

### Integration Requirements
- [ ] Command registry integrates with logging system
- [ ] Command execution integrates with error handler
- [ ] Command listing integrates with command palette

### Edge Cases & Error Handling
- [ ] Handle duplicate command names (reject or override)
- [ ] Handle command not found (return error)
- [ ] Handle handler panic (recover and return error)
- [ ] Handle very long command names
- [ ] Handle special characters in command names
- [ ] Handle concurrent command execution
- [ ] Handle command with invalid arguments

### Performance Requirements
- [ ] Register command <1ms
- [ ] Execute command <5ms
- [ ] List commands <10ms for 100 commands
- [ ] Command lookup <1ms

### User Experience Requirements
- [ ] Clear error messages for command failures
- [ ] Command names are descriptive and consistent
- [ ] Command documentation accessible
- [ ] No UI blocking during command execution

**Files:** [`internal/commands/registry.go`], [`internal/commands/core.go`]

---

## Milestone 19: Command Palette UI
**Goal:** Fuzzy-searchable command launcher

**Deliverables:**
- Modal overlay
- List all commands
- Fuzzy filter as user types
- Execute selected command
- Keyboard shortcut to open (Ctrl+P)

**Test Criteria:**

### Functional Requirements
- [ ] Press Ctrl+P, palette opens
- [ ] See list of all registered commands
- [ ] Type to filter commands (fuzzy search)
- [ ] Select command, press Enter to execute
- [ ] Command executes successfully
- [ ] Palette closes after execution
- [ ] Esc closes palette without executing
- [ ] Works with 20+ commands
- [ ] Works with 100+ commands
- [ ] Clear search restores full list

### Integration Requirements
- [ ] Palette integrates with command registry
- [ ] Command execution integrates with command system
- [ ] Keyboard shortcuts integrate with input system
- [ ] Palette close integrates with TUI shell

### Edge Cases & Error Handling
- [ ] Handle empty command list (show empty state)
- [ ] Handle no search results (show empty state)
- [ ] Handle very long command names (truncate)
- [ ] Handle very long command descriptions (truncate)
- [ ] Handle command execution errors (show error)
- [ ] Handle rapid typing (debounce properly)
- [ ] Handle palette open when already open
- [ ] Handle special characters in search

### Performance Requirements
- [ ] Palette open <50ms
- [ ] Search latency <10ms for 100 commands
- [ ] Command execution <50ms
- [ ] Palette close <50ms
- [ ] No UI lag during search

### User Experience Requirements
- [ ] Smooth palette open/close animation
- [ ] Clear visual indication of selected command
- [ ] Command descriptions visible and readable
- [ ] Empty state message is helpful
- [ ] Search input clearly visible
- [ ] Keyboard shortcuts discoverable
- [ ] No visual artifacts during navigation

**Files:** [`ui/palette/model.go`]

---

## Milestone 20: File Finder
**Goal:** Fuzzy find files in working directory

**Deliverables:**
- Traverse working directory
- Respect .gitignore rules
- Fuzzy search by filename
- Multiple file selection
- Modal UI

**Test Criteria:**

### Functional Requirements
- [ ] Open file finder with keyboard shortcut
- [ ] See list of all files in working directory
- [ ] .gitignore rules respected (ignored files not shown)
- [ ] Type to filter files (fuzzy search)
- [ ] Select single file with Enter
- [ ] Select multiple files with Space
- [ ] Confirm selection with Enter
- [ ] Cancel selection with Esc
- [ ] Works with 1000+ files
- [ ] Recursive directory traversal
- [ ] Show file extensions
- [ ] Show relative paths

### Integration Requirements
- [ ] File finder integrates with file system
- [ ] .gitignore parsing integrates with gitignore library
- [ ] File selection integrates with batch editor
- [ ] Keyboard shortcuts integrate with input system

### Edge Cases & Error Handling
- [ ] Handle empty directory (show empty state)
- [ ] Handle no search results (show empty state)
- [ ] Handle very long file paths (truncate)
- [ ] Handle very long filenames (truncate)
- [ ] Handle permission denied errors (skip file)
- [ ] Handle circular symlinks (detect and skip)
- [ ] Handle very deep directory structures
- [ ] Handle special characters in filenames
- [ ] Handle rapid navigation (no lag)

### Performance Requirements
- [ ] File finder open <100ms
- [ ] Directory scan <200ms for 1000 files
- [ ] Search latency <50ms for 1000 files
- [ ] File selection <10ms
- [ ] No UI lag during navigation

### User Experience Requirements
- [ ] Smooth scrolling through file list
- [ ] Clear visual indication of selected files
- [ ] Empty state message is helpful
- [ ] Search input clearly visible
- [ ] Keyboard shortcuts discoverable
- [ ] Progress indicator for large directories
- [ ] No visual artifacts during navigation

**Files:** [`internal/files/finder.go`], [`ui/filereference/model.go`]

---

## Milestone 21: Title Extraction
**Goal:** Extract titles from file frontmatter

**Deliverables:**
- Parse YAML frontmatter from files
- Extract title field
- Fallback to filename if no title
- Handle various file types

**Test Criteria:**

### Functional Requirements
- [ ] Extract title from markdown with frontmatter
- [ ] Use filename if no frontmatter present
- [ ] Use filename if no title field in frontmatter
- [ ] Handle .txt files
- [ ] Handle .md files
- [ ] Handle .js files
- [ ] Handle other text-based file types
- [ ] Strip file extension from filename
- [ ] Handle empty frontmatter
- [ ] Handle malformed frontmatter gracefully

### Integration Requirements
- [ ] Title extraction integrates with file I/O
- [ ] Title extraction integrates with YAML parser
- [ ] Title extraction integrates with batch editor

### Edge Cases & Error Handling
- [ ] Handle missing files (return error)
- [ ] Handle empty files (use filename)
- [ ] Handle files with only frontmatter
- [ ] Handle files with only content (no frontmatter)
- [ ] Handle invalid YAML frontmatter (use filename)
- [ ] Handle very long titles (truncate)
- [ ] Handle special characters in titles
- [ ] Handle unicode characters in titles
- [ ] Handle binary files (return error)

### Performance Requirements
- [ ] Extract title <10ms per file
- [ ] Handle 1000 files in <5s
- [ ] Memory usage scales linearly with file count

### User Experience Requirements
- [ ] Clear error messages for file read failures
- [ ] Progress indicators for batch operations
- [ ] No UI blocking during extraction

**Files:** [`internal/files/title_extractor.go`]

---

## Milestone 22: Batch Title Editor & Link Insertion
**Goal:** Edit titles and insert markdown links

**Deliverables:**
- Batch editor showing selected files
- Editable title for each file
- Quick-accept all option
- Insert markdown links at cursor: `[title](path)`

**Test Criteria:**

### Functional Requirements
- [ ] Select multiple files from file finder
- [ ] Batch editor shows all selected files
- [ ] Edit individual titles inline
- [ ] Accept all titles with single action
- [ ] Markdown links inserted at cursor position
- [ ] Link format: `[title](relative/path)`
- [ ] Multiple links on separate lines
- [ ] Cancel selection with Esc
- [ ] Remove individual files from selection
- [ ] Clear all selection

### Integration Requirements
- [ ] Batch editor integrates with file finder
- [ ] Batch editor integrates with title extractor
- [ ] Link insertion integrates with editor model
- [ ] Keyboard shortcuts integrate with input system

### Edge Cases & Error Handling
- [ ] Handle empty selection (show empty state)
- [ ] Handle very long titles (truncate in display)
- [ ] Handle very long file paths (truncate in display)
- [ ] Handle special characters in titles (escape in markdown)
- [ ] Handle special characters in paths (escape in markdown)
- [ ] Handle rapid title edits (debounce properly)
- [ ] Handle insertion at start of file
- [ ] Handle insertion at end of file
- [ ] Handle insertion in middle of line

### Performance Requirements
- [ ] Batch editor open <50ms
- [ ] Title edit <10ms
- [ ] Accept all <50ms for 100 files
- [ ] Link insertion <100ms for 100 links
- [ ] No UI lag during editing

### User Experience Requirements
- [ ] Clear batch editor UI
- [ ] Smooth navigation through file list
- [ ] Clear visual feedback for selected files
- [ ] Intuitive keyboard shortcuts
- [ ] Clear indication of link format
- [ ] Smooth link insertion
- [ ] No visual artifacts during editing

**Files:** [`ui/filereference/model.go`]

---

## Milestone 23: Prompt Validation
**Goal:** Check prompt syntax and metadata

**Deliverables:**
- Validation checks (errors and warnings)
- Run on startup (silent unless errors)
- Manual trigger via command
- Store validation status per prompt

**Validation Checks:**
- **Errors:** Duplicate placeholder names, invalid YAML, missing title, file >1MB
- **Warnings:** Invalid placeholder types, invalid names, malformed placeholders, missing optional metadata

**Test Criteria:**

### Functional Requirements
- [ ] Validate good prompt (no errors or warnings)
- [ ] Detect duplicate placeholder names (error)
- [ ] Detect invalid YAML frontmatter (error)
- [ ] Detect missing title field (error)
- [ ] Detect file size >1MB (error)
- [ ] Detect invalid placeholder types (warning)
- [ ] Detect invalid placeholder names (warning)
- [ ] Detect malformed placeholders (warning)
- [ ] Detect missing optional metadata (warning)
- [ ] Validation runs automatically on startup
- [ ] Manual validation via command
- [ ] Results stored per prompt
- [ ] Validation status persists across sessions

### Integration Requirements
- [ ] Validation integrates with library loader
- [ ] Validation integrates with YAML parser
- [ ] Validation integrates with placeholder parser
- [ ] Validation results integrate with UI display
- [ ] Validation status integrates with library browser

### Edge Cases & Error Handling
- [ ] Handle prompt with no frontmatter (error)
- [ ] Handle prompt with empty frontmatter (warning)
- [ ] Handle prompt with very long content (>1MB)
- [ ] Handle prompt with 100+ placeholders
- [ ] Handle prompt with nested placeholders (error)
- [ ] Handle prompt with unicode characters
- [ ] Handle corrupted prompt files (error)
- [ ] Handle validation of 1000+ prompts

### Performance Requirements
- [ ] Validate single prompt <10ms
- [ ] Validate 100 prompts <500ms
- [ ] Validate 1000 prompts <5s
- [ ] Startup validation <2s for 1000 prompts
- [ ] Manual validation <5s for 1000 prompts
- [ ] Memory usage scales linearly with prompt count

### User Experience Requirements
- [ ] Clear indication of validation status
- [ ] Silent validation on startup (no UI unless errors)
- [ ] Clear error messages for validation failures
- [ ] Progress indicator for large validation operations
- [ ] No UI blocking during validation
- [ ] Validation results easily accessible

**Files:** [`internal/library/validator.go`]

---

## Milestone 24: Validation Results Display
**Goal:** Show validation errors and warnings

**Deliverables:**
- Modal showing validation results
- Group by error/warning type
- Show affected prompts
- ⚠️ icon in library browser for invalid prompts
- Block insertion of error-level prompts

**Test Criteria:**

### Functional Requirements
- [ ] Run validation, open results modal
- [ ] Errors grouped separately from warnings
- [ ] Each result shows prompt name and issue
- [ ] Results sorted by severity (errors first)
- [ ] Navigate through results with arrow keys
- [ ] Select result to see details
- [ ] Close modal with Esc
- [ ] Library browser shows ⚠️ for invalid prompts
- [ ] Cannot insert error-level prompts
- [ ] Can insert warning-level prompts
- [ ] Clear validation results on re-validation

### Integration Requirements
- [ ] Results modal integrates with validation system
- [ ] Results display integrates with library browser
- [ ] Error blocking integrates with prompt insertion
- [ ] Keyboard shortcuts integrate with input system

### Edge Cases & Error Handling
- [ ] Handle no validation results (show empty state)
- [ ] Handle very long prompt names (truncate)
- [ ] Handle very long error messages (truncate or wrap)
- [ ] Handle 100+ validation results
- [ ] Handle rapid navigation through results
- [ ] Handle modal open when already open
- [ ] Handle validation during modal display

### Performance Requirements
- [ ] Results modal open <50ms
- [ ] Navigation latency <10ms
- [ ] Display 100 results <100ms
- [ ] No UI lag during navigation

### User Experience Requirements
- [ ] Clear visual distinction between errors and warnings
- [ ] Smooth scrolling through results
- [ ] Clear visual indication of selected result
- [ ] Empty state message is helpful
- [ ] Keyboard shortcuts discoverable
- [ ] No visual artifacts during navigation

**Files:** [`ui/validation/model.go`]

---

## Milestone 25: Prompt Creator
**Goal:** Create new prompts with guided wizard

**Deliverables:**
- Step-by-step wizard
- Title input (required)
- Description input (optional)
- Tags input with autocomplete
- Category selection
- Initial content editor
- Save to library

**Test Criteria:**

### Functional Requirements
- [ ] Open prompt creator with keyboard shortcut
- [ ] Enter title (required field)
- [ ] Enter description (optional field)
- [ ] Add tags (comma-separated, autocomplete)
- [ ] Select category from dropdown
- [ ] Write initial content in editor
- [ ] Navigate between wizard steps
- [ ] Save prompt to library
- [ ] Prompt appears in library browser
- [ ] File created in correct category folder
- [ ] Cancel creation with Esc
- [ ] Validate before saving

### Integration Requirements
- [ ] Prompt creator integrates with library loader
- [ ] Tag autocomplete integrates with existing tags
- [ ] Category selection integrates with category structure
- [ ] Content editor integrates with editor model
- [ ] Save operation integrates with file I/O
- [ ] Validation integrates with validator

### Edge Cases & Error Handling
- [ ] Handle missing title (show error, block save)
- [ ] Handle very long title (>200 characters)
- [ ] Handle very long description (>1000 characters)
- [ ] Handle very long content (>1MB)
- [ ] Handle very long tag names
- [ ] Handle duplicate tags (deduplicate)
- [ ] Handle special characters in title
- [ ] Handle invalid category (show error)
- [ ] Handle save errors (show error, retry)
- [ ] Handle disk full errors
- [ ] Handle permission denied errors

### Performance Requirements
- [ ] Prompt creator open <50ms
- [ ] Tag autocomplete <10ms
- [ ] Save operation <100ms
- [ ] Library reload <100ms
- [ ] No UI lag during typing

### User Experience Requirements
- [ ] Clear wizard step indicators
- [ ] Clear required vs optional fields
- [ ] Helpful placeholder text
- [ ] Real-time validation feedback
- [ ] Clear error messages for invalid input
- [ ] Smooth navigation between steps
- [ ] Clear save confirmation
- [ ] Keyboard shortcuts discoverable
- [ ] No visual artifacts during navigation

**Files:** [`ui/promptcreator/model.go`], [`internal/prompt/creator.go`]

---

## Milestone 26: Prompt Editor
**Goal:** Edit existing prompts with preview

**Deliverables:**
- Edit mode (raw markdown + frontmatter)
- Preview mode (rendered markdown)
- Toggle with Ctrl+P
- Mode indicator in status bar
- Explicit save (no auto-save)

**Test Criteria:**

### Functional Requirements
- [ ] Open existing prompt for editing
- [ ] See raw markdown with frontmatter
- [ ] Press Ctrl+P, switch to preview mode
- [ ] See rendered markdown in preview
- [ ] Press Ctrl+P, switch back to edit mode
- [ ] Make changes to content
- [ ] Make changes to frontmatter
- [ ] Save explicitly with keyboard shortcut
- [ ] Changes persisted to file
- [ ] Library reloaded with changes
- [ ] Cancel changes with Esc
- [ ] Mode indicator shows current mode

### Integration Requirements
- [ ] Prompt editor integrates with file I/O
- [ ] Preview mode integrates with markdown renderer
- [ ] Save operation integrates with library loader
- [ ] Mode toggle integrates with input system
- [ ] Validation integrates with validator

### Edge Cases & Error Handling
- [ ] Handle very long content (>1MB)
- [ ] Handle very long frontmatter
- [ ] Handle malformed frontmatter (show error)
- [ ] Handle save errors (show error, retry)
- [ ] Handle disk full errors
- [ ] Handle permission denied errors
- [ ] Handle unsaved changes on close (prompt to save)
- [ ] Handle rapid mode switching
- [ ] Handle special characters in content

### Performance Requirements
- [ ] Prompt editor open <50ms
- [ ] Mode switch <50ms
- [ ] Preview render <100ms for 1MB content
- [ ] Save operation <100ms
- [ ] Library reload <100ms
- [ ] No UI lag during typing

### User Experience Requirements
- [ ] Clear mode indicator in status bar
- [ ] Smooth transition between edit and preview
- [ ] Clear visual feedback for unsaved changes
- [ ] Clear save confirmation
- [ ] Helpful error messages for save failures
- [ ] Keyboard shortcuts discoverable
- [ ] No visual artifacts during mode switch
- [ ] Preview renders markdown correctly

**Files:** [`ui/prompteditor/model.go`]

---

## Milestone 27: AI Provider Interface & Claude Implementation
**Goal:** Implement AI provider abstraction with Claude as first provider

**Deliverables:**
- AIProvider interface in `internal/ai/provider.go`
- ClaudeProvider implementation in `internal/ai/claude.go`
- Provider factory in `internal/ai/factory.go`
- Authentication with API key
- Send message, receive response
- Error handling (rate limit, auth, network, timeout)
- Retry logic (max 3 attempts)
- Unit tests for interface and implementation
- Integration tests with Claude API

**Test Criteria:**

### Functional Requirements
- [ ] AIProvider interface defined with required methods
- [ ] ClaudeProvider implements AIProvider interface
- [ ] Provider factory creates correct provider based on config
- [ ] Send test request to Claude API
- [ ] Receive valid response
- [ ] Authentication with API key works
- [ ] Send message, receive response
- [ ] Handle 401 auth error
- [ ] Handle 429 rate limit
- [ ] Handle network timeout
- [ ] Retry on transient failures
- [ ] Stop after 3 failed retries
- [ ] Exponential backoff between retries

### Integration Requirements
- [ ] Provider factory integrates with config system
- [ ] Provider factory reads ai_provider field from config
- [ ] API client integrates with logging system
- [ ] Error handling integrates with error handler

### Edge Cases & Error Handling
- [ ] Handle missing API key (show error, prompt to configure)
- [ ] Handle invalid API key (show error, prompt to reconfigure)
- [ ] Handle API service unavailable (retry with backoff)
- [ ] Handle malformed API response (show error, log details)
- [ ] Handle concurrent API requests (queue properly)
- [ ] Handle very long responses (>100KB)
- [ ] Handle rate limit exceeded (wait and retry)
- [ ] Handle network connection refused
- [ ] Handle DNS resolution failures

### Performance Requirements
- [ ] API request latency <2s for typical response
- [ ] Retry backoff: 1s, 2s, 4s (exponential)
- [ ] Connection timeout: 30s
- [ ] Read timeout: 60s
- [ ] Handle 10 concurrent requests

### User Experience Requirements
- [ ] Clear error messages for API failures
- [ ] Progress indicator during API requests
- [ ] Retry status shown in status bar
- [ ] No UI blocking during API requests
- [ ] Graceful degradation if API unavailable

**Files:** [`internal/ai/provider.go`], [`internal/ai/claude.go`], [`internal/ai/factory.go`]

---

## Milestone 28: Context Selection Algorithm
**Goal:** Select relevant library prompts for AI context

**Deliverables:**
- Extract keywords from composition
- Score each library prompt
- Select top 3-5 prompts
- Fit within token budget (15% of context)

**Scoring Algorithm:**
- Tag match: +10 per matching tag
- Category match: +5 if same category
- Keyword overlap: +1 per matching word
- Recently used: +3 if used in last session
- Frequently used: +use_count

**Test Criteria:**

### Functional Requirements
- [ ] Given composition, extract keywords
- [ ] Score all library prompts
- [ ] Top 3-5 prompts selected
- [ ] Relevant prompts ranked higher
- [ ] Token budget respected
- [ ] Explicitly referenced prompts always included
- [ ] Scoring algorithm applied correctly
  - [ ] Tag match: +10 per matching tag
  - [ ] Category match: +5 if same category
  - [ ] Keyword overlap: +1 per matching word
  - [ ] Recently used: +3 if used in last session
  - [ ] Frequently used: +use_count

### Integration Requirements
- [ ] Context selection integrates with library loader
- [ ] Context selection integrates with placeholder parser
- [ ] Context selection integrates with token estimator
- [ ] Context selection integrates with history tracking

### Edge Cases & Error Handling
- [ ] Handle empty composition (no keywords extracted)
- [ ] Handle composition with no matching prompts
- [ ] Handle library with no prompts
- [ ] Handle very large library (1000+ prompts)
- [ ] Handle very long composition (>10000 words)
- [ ] Handle prompts with no metadata
- [ ] Handle prompts with duplicate scores
- [ ] Handle token budget exceeded (select fewer prompts)
- [ ] Handle explicitly referenced prompts exceeding budget

### Performance Requirements
- [ ] Keyword extraction <50ms for 1000-word composition
- [ ] Scoring <100ms for 1000 prompts
- [ ] Selection <10ms
- [ ] Total context selection <200ms
- [ ] Memory usage scales linearly with prompt count

### User Experience Requirements
- [ ] Selected prompts shown in status bar
- [ ] Token usage displayed
- [ ] Clear indication of context limit
- [ ] No UI blocking during selection
- [ ] Progress indicator for large libraries

**Files:** [`internal/ai/context.go`], [`internal/library/indexer.go`]

---

## Milestone 29: Token Estimation & Budget
**Goal:** Estimate tokens and enforce limits

**Deliverables:**
- Token counting function
- Estimate composition tokens
- Estimate library prompt tokens
- Show token count in status bar
- Warn at 15% of context window
- Block suggestions at 25%

**Test Criteria:**

### Functional Requirements
- [ ] Estimate tokens for sample text
- [ ] Show token count in status bar
- [ ] Warning shown at 15% threshold
- [ ] Suggestions blocked at 25% threshold
- [ ] Token count updates as user types
- [ ] Accurate within 10% of actual tokens
- [ ] Estimate composition tokens
- [ ] Estimate library prompt tokens
- [ ] Calculate total context tokens
- [ ] Calculate remaining tokens

### Integration Requirements
- [ ] Token estimation integrates with editor model
- [ ] Token estimation integrates with context selection
- [ ] Token estimation integrates with status bar
- [ ] Token estimation integrates with suggestions system

### Edge Cases & Error Handling
- [ ] Handle empty composition (0 tokens)
- [ ] Handle very long composition (>100000 tokens)
- [ ] Handle unicode characters correctly
- [ ] Handle special characters correctly
- [ ] Handle mixed language content
- [ ] Handle code blocks correctly
- [ ] Handle markdown formatting correctly
- [ ] Handle rapid typing (debounce updates)

### Performance Requirements
- [ ] Token estimation <10ms per keystroke
- [ ] Token estimation <100ms for 10000-word composition
- [ ] Debounce delay: 200-300ms
- [ ] No UI lag during typing
- [ ] Memory usage scales with document size

### User Experience Requirements
- [ ] Token count visible in status bar
- [ ] Clear warning at 15% threshold
- [ ] Clear blocking message at 25% threshold
- [ ] Smooth token count updates
- [ ] Visual indicator of token usage
- [ ] No interruption to typing

**Files:** [`internal/ai/tokens.go`]

---

## Milestone 30: Suggestion Parsing
**Goal:** Parse AI response into structured suggestions

**Deliverables:**
- Parse 6 suggestion types from API response
- Extract title, description, proposed changes
- Handle malformed responses gracefully

**Suggestion Types:**
1. Prompt Recommendations
2. Gap Analysis
3. Formatting
4. Contradictions
5. Clarity
6. Reformatting

**Test Criteria:**

### Functional Requirements
- [ ] Parse valid API response
- [ ] Extract all suggestion types
- [ ] Each suggestion has title and description
- [ ] Proposed changes extracted
- [ ] Handle missing fields gracefully
- [ ] Handle malformed JSON gracefully
- [ ] Parse 6 suggestion types correctly
  - [ ] Prompt Recommendations
  - [ ] Gap Analysis
  - [ ] Formatting
  - [ ] Contradictions
  - [ ] Clarity
  - [ ] Reformatting

### Integration Requirements
- [ ] Suggestion parsing integrates with API client
- [ ] Suggestion parsing integrates with suggestions panel
- [ ] Suggestion parsing integrates with logging system
- [ ] Error handling integrates with error handler

### Edge Cases & Error Handling
- [ ] Handle empty API response
- [ ] Handle API response with no suggestions
- [ ] Handle API response with partial suggestions
- [ ] Handle API response with duplicate suggestions
- [ ] Handle API response with invalid suggestion types
- [ ] Handle API response with very long suggestions
- [ ] Handle API response with unicode characters
- [ ] Handle API response with special characters
- [ ] Handle API response with nested JSON structures

### Performance Requirements
- [ ] Parse API response <100ms
- [ ] Parse 100 suggestions <200ms
- [ ] Handle very long responses (>10KB)
- [ ] Memory usage scales with suggestion count

### User Experience Requirements
- [ ] Clear error messages for parsing failures
- [ ] Graceful degradation for partial responses
- [ ] No UI blocking during parsing
- [ ] Progress indicator for long responses
- [ ] Log parsing errors for debugging

**Files:** [`internal/ai/suggestions.go`]

---

## Milestone 31: Suggestions Panel
**Goal:** Display AI suggestions in right panel

**Deliverables:**
- Split-pane layout (editor left, suggestions right)
- Scrollable suggestion list
- Navigate with j/k or arrows
- Show suggestion details
- Accept ('a') or dismiss ('d') actions
- Status indicators (pending, applying, applied, error)

**Test Criteria:**

### Functional Requirements
- [ ] Trigger AI suggestions
- [ ] Right panel opens with suggestions
- [ ] Navigate through suggestions
- [ ] See title and description
- [ ] Press 'a' to accept suggestion
- [ ] Press 'd' to dismiss suggestion
- [ ] Status updates shown
- [ ] Error messages displayed if API fails
- [ ] Split-pane layout (editor left, suggestions right)
- [ ] Scrollable suggestion list
- [ ] Navigate with j/k or arrows
- [ ] Show suggestion details
- [ ] Status indicators (pending, applying, applied, error)

### Integration Requirements
- [ ] Suggestions panel integrates with API client
- [ ] Suggestions panel integrates with suggestion parser
- [ ] Suggestions panel integrates with editor model
- [ ] Suggestions panel integrates with diff system
- [ ] Keyboard shortcuts integrate with input system

### Edge Cases & Error Handling
- [ ] Handle no suggestions (show empty state)
- [ ] Handle API failure (show error, retry option)
- [ ] Handle very long suggestion titles (truncate)
- [ ] Handle very long suggestion descriptions (truncate or wrap)
- [ ] Handle very long suggestion lists (100+ suggestions)
- [ ] Handle rapid navigation (no lag)
- [ ] Handle panel open when already open
- [ ] Handle suggestion acceptance during API request
- [ ] Handle suggestion dismissal during API request
- [ ] Handle panel close during API request

### Performance Requirements
- [ ] Panel open <100ms
- [ ] Navigation latency <10ms
- [ ] Display 100 suggestions <100ms
- [ ] Accept action <50ms
- [ ] Dismiss action <50ms
- [ ] No UI lag during navigation
- [ ] Handle 1000+ suggestions smoothly

### User Experience Requirements
- [ ] Smooth panel open/close animation
- [ ] Clear visual indication of selected suggestion
- [ ] Clear status indicators (pending, applying, applied, error)
- [ ] Empty state message is helpful
- [ ] Error messages are actionable
- [ ] Keyboard shortcuts discoverable
- [ ] Smooth scrolling through suggestions
- [ ] No visual artifacts during navigation
- [ ] Clear visual feedback for accept/dismiss actions

**Files:** [`ui/suggestions/model.go`]

---

## Milestone 32: Diff Generation
**Goal:** Create unified diff from AI edits

**Deliverables:**
- Request structured edits from Claude
- Parse line ranges and new content
- Generate unified diff format
- Display in modal

**Test Criteria:**

### Functional Requirements
- [ ] Accept suggestion
- [ ] Request edits from API
- [ ] Parse edit instructions
- [ ] Generate unified diff
- [ ] Diff modal opens
- [ ] Shows old vs new content
- [ ] Unified diff format correct
- [ ] Line numbers accurate
- [ ] Request structured edits from Claude
- [ ] Parse line ranges and new content
- [ ] Display in modal

### Integration Requirements
- [ ] Diff generation integrates with API client
- [ ] Diff generation integrates with suggestion parser
- [ ] Diff generation integrates with editor model
- [ ] Diff viewer integrates with TUI shell
- [ ] Diff application integrates with undo stack

### Edge Cases & Error Handling
- [ ] Handle API failure (show error, retry option)
- [ ] Handle malformed edit instructions (show error)
- [ ] Handle invalid line ranges (show error)
- [ ] Handle very large diffs (>1000 lines)
- [ ] Handle diffs with many changes (100+ hunks)
- [ ] Handle diffs with unicode characters
- [ ] Handle diffs with special characters
- [ ] Handle diff generation during editing
- [ ] Handle diff modal open when already open

### Performance Requirements
- [ ] API request <2s for typical response
- [ ] Parse edit instructions <50ms
- [ ] Generate diff <100ms for 1000 lines
- [ ] Diff modal open <100ms
- [ ] Render diff <100ms for 1000 lines
- [ ] Handle very large diffs (>10000 lines)

### User Experience Requirements
- [ ] Clear diff visualization
- [ ] Clear old vs new content distinction
- [ ] Accurate line numbers
- [ ] Smooth diff rendering
- [ ] Clear error messages for failures
- [ ] No UI blocking during generation
- [ ] Progress indicator for long operations
- [ ] Keyboard shortcuts discoverable
- [ ] No visual artifacts during rendering

**Files:** [`internal/ai/diff.go`], [`ui/diffviewer/model.go`]

---

## Milestone 33: Diff Application
**Goal:** Apply accepted changes as single undo action

**Deliverables:**
- Apply diff to composition
- Update editor content
- Add to undo stack as single action
- Unlock editor after application

**Test Criteria:**

### Functional Requirements
- [ ] Accept diff in modal
- [ ] Changes applied to editor
- [ ] Editor content updated correctly
- [ ] Press Ctrl+Z, entire change undone
- [ ] Press Ctrl+Y, change reapplied
- [ ] Editor unlocked after application
- [ ] Reject diff, no changes applied
- [ ] Apply diff to composition
- [ ] Update editor content
- [ ] Add to undo stack as single action
- [ ] Unlock editor after application

### Integration Requirements
- [ ] Diff application integrates with editor model
- [ ] Diff application integrates with undo stack
- [ ] Diff application integrates with auto-save
- [ ] Diff application integrates with diff viewer

### Edge Cases & Error Handling
- [ ] Handle diff application during editing (block or queue)
- [ ] Handle very large diffs (>1000 lines)
- [ ] Handle diffs with many changes (100+ hunks)
- [ ] Handle diffs with unicode characters
- [ ] Handle diffs with special characters
- [ ] Handle diff application errors (show error, rollback)
- [ ] Handle undo stack full (show error)
- [ ] Handle concurrent diff applications (queue properly)
- [ ] Handle diff rejection during application

### Performance Requirements
- [ ] Apply diff <100ms for 1000 lines
- [ ] Apply diff <1s for 10000 lines
- [ ] Undo operation <10ms
- [ ] Redo operation <10ms
- [ ] No UI lag during application
- [ ] Handle very large diffs smoothly

### User Experience Requirements
- [ ] Clear visual feedback during application
- [ ] Clear confirmation after application
- [ ] Clear error messages for failures
- [ ] Smooth undo/redo transitions
- [ ] No visual artifacts during application
- [ ] Editor locked during application (show indicator)
- [ ] Clear indication of undo/redo availability
- [ ] No interruption to editing after application

**Files:** [`internal/ai/diff.go`], [`ui/workspace/model.go`]

---

## Milestone 34: Vim State Machine
**Goal:** Implement vim mode state management

**Deliverables:**
- State machine for Normal/Insert/Visual modes
- Mode transitions
- Mode indicator in status bar
- Global vim mode toggle in config

**Test Criteria:**

### Functional Requirements
- [ ] Enable vim mode in config
- [ ] Start in Normal mode
- [ ] Press 'i', enter Insert mode
- [ ] Press Esc, return to Normal mode
- [ ] Press 'v', enter Visual mode
- [ ] Mode indicator shows current mode
- [ ] Disable vim mode, standard keybindings work
- [ ] State machine for Normal/Insert/Visual modes
- [ ] Mode transitions work correctly
- [ ] Global vim mode toggle in config

### Integration Requirements
- [ ] Vim state integrates with editor model
- [ ] Vim state integrates with keyboard input
- [ ] Mode indicator integrates with status bar
- [ ] Vim state integrates with config system

### Edge Cases & Error Handling
- [ ] Handle rapid mode switches
- [ ] Handle mode switch during editing
- [ ] Handle mode switch in different components
- [ ] Handle invalid mode transitions
- [ ] Handle mode switch with unsaved changes
- [ ] Handle mode switch during API requests
- [ ] Handle mode switch with active selections

### Performance Requirements
- [ ] Mode switch <10ms
- [ ] Mode indicator update <10ms
- [ ] No UI lag during mode switches
- [ ] State machine operations <1ms

### User Experience Requirements
- [ ] Clear mode indicator in status bar
- [ ] Smooth mode transitions
- [ ] Clear visual feedback for mode changes
- [ ] No visual artifacts during mode switches
- [ ] Intuitive mode transitions
- [ ] Clear indication of current mode

**Files:** [`internal/vim/state.go`]

---

## Milestone 35: Vim Keybindings
**Goal:** Context-aware vim keybindings

**Deliverables:**
- Keybinding maps per mode
- Context-aware routing (editor, browser, palette, etc.)
- j/k navigation in lists
- h/l for horizontal movement
- / for search
- Full vim editing in composition

**Test Criteria:**

### Functional Requirements
- [ ] In editor Normal mode: h/j/k/l moves cursor
- [ ] In editor Insert mode: type normally
- [ ] In library browser: j/k navigates list
- [ ] In command palette: j/k navigates commands
- [ ] In suggestions panel: j/k navigates suggestions
- [ ] / opens search in browser
- [ ] Vim commands work consistently across components
- [ ] Keybinding maps per mode
- [ ] Context-aware routing (editor, browser, palette, etc.)
- [ ] j/k navigation in lists
- [ ] h/l for horizontal movement
- [ ] / for search
- [ ] Full vim editing in composition

### Integration Requirements
- [ ] Vim keybindings integrate with state machine
- [ ] Vim keybindings integrate with keyboard input
- [ ] Context-aware routing integrates with component system
- [ ] Vim keybindings integrate with editor model
- [ ] Vim keybindings integrate with browser UI
- [ ] Vim keybindings integrate with command palette

### Edge Cases & Error Handling
- [ ] Handle keybindings in different modes
- [ ] Handle keybindings in different components
- [ ] Handle conflicting keybindings
- [ ] Handle rapid key presses
- [ ] Handle keybindings during API requests
- [ ] Handle keybindings with unsaved changes
- [ ] Handle keybindings with active selections
- [ ] Handle keybindings during modal display

### Performance Requirements
- [ ] Keybinding lookup <1ms
- [ ] Keybinding execution <5ms
- [ ] Context routing <1ms
- [ ] No UI lag during key presses
- [ ] Handle rapid key presses smoothly

### User Experience Requirements
- [ ] Keybindings work consistently across components
- [ ] Clear visual feedback for key presses
- [ ] Intuitive keybinding behavior
- [ ] No visual artifacts during key presses
- [ ] Smooth cursor movement
- [ ] Keyboard shortcuts discoverable
- [ ] Clear indication of current mode

**Files:** [`internal/vim/keymaps.go`], [`internal/vim/context.go`]

---

## Milestone 36: Settings Panel
**Goal:** Edit configuration in UI

**Deliverables:**
- Modal with form-style interface
- Toggle vim mode (requires restart)
- Edit API key (masked input)
- Select Claude model from dropdown
- Immediate persistence to config.yaml

**Test Criteria:**

### Functional Requirements
- [ ] Open settings panel
- [ ] Toggle vim mode
- [ ] Edit API key (shows masked)
- [ ] Select different model
- [ ] Save settings
- [ ] Config file updated
- [ ] Restart required message shown for vim toggle
- [ ] Other settings apply immediately
- [ ] Modal with form-style interface
- [ ] Immediate persistence to config.yaml

### Integration Requirements
- [ ] Settings panel integrates with config system
- [ ] Settings panel integrates with config updater
- [ ] Settings panel integrates with TUI shell
- [ ] Settings panel integrates with keyboard input

### Edge Cases & Error Handling
- [ ] Handle invalid API key format (show error)
- [ ] Handle invalid model selection (show error)
- [ ] Handle save errors (show error, retry)
- [ ] Handle disk full errors
- [ ] Handle permission denied errors
- [ ] Handle unsaved changes on close (prompt to save)
- [ ] Handle rapid setting changes
- [ ] Handle settings panel open when already open

### Performance Requirements
- [ ] Settings panel open <50ms
- [ ] Setting change <10ms
- [ ] Save operation <50ms
- [ ] Config file update <50ms
- [ ] No UI lag during editing

### User Experience Requirements
- [ ] Clear form-style interface
- [ ] Clear required vs optional fields
- [ ] Helpful placeholder text
- [ ] Clear save confirmation
- [ ] Clear error messages for invalid input
- [ ] Clear restart required message
- [ ] Keyboard shortcuts discoverable
- [ ] No visual artifacts during editing

**Files:** [`ui/settings/model.go`], [`internal/config/updater.go`]

---

## Milestone 37: Responsive Layout
**Goal:** Handle terminal resize and narrow terminals

**Deliverables:**
- Detect terminal size changes
- Adjust layout dynamically
- Split-pane resizing with divider
- Minimum width handling (80 columns)
- Hide panels if too narrow

**Test Criteria:**

### Functional Requirements
- [ ] Resize terminal, layout adjusts
- [ ] Drag divider, panes resize
- [ ] Narrow terminal (<100 cols), suggestions panel hidden
- [ ] Very narrow (<80 cols), warning shown
- [ ] Wide terminal (>200 cols), layout uses space well
- [ ] No visual artifacts on resize
- [ ] Detect terminal size changes
- [ ] Adjust layout dynamically
- [ ] Split-pane resizing with divider
- [ ] Minimum width handling (80 columns)
- [ ] Hide panels if too narrow

### Integration Requirements
- [ ] Responsive layout integrates with TUI shell
- [ ] Responsive layout integrates with editor model
- [ ] Responsive layout integrates with suggestions panel
- [ ] Responsive layout integrates with status bar

### Edge Cases & Error Handling
- [ ] Handle very narrow terminals (<80 columns)
- [ ] Handle very wide terminals (>300 columns)
- [ ] Handle rapid resize events
- [ ] Handle resize during editing
- [ ] Handle resize during API requests
- [ ] Handle resize with active modals
- [ ] Handle resize with active selections
- [ ] Handle divider drag at boundaries

### Performance Requirements
- [ ] Layout adjustment <50ms on resize
- [ ] Divider drag <10ms latency
- [ ] Panel hide/show <50ms
- [ ] No UI lag during resize
- [ ] Handle rapid resize events smoothly

### User Experience Requirements
- [ ] Smooth layout transitions
- [ ] Clear visual feedback during resize
- [ ] Clear warning for very narrow terminals
- [ ] Intuitive divider dragging
- [ ] No visual artifacts during resize
- [ ] Layout uses space efficiently
- [ ] Clear indication of hidden panels

**Files:** [`ui/app/model.go`], [`ui/workspace/model.go`]

---

## Milestone 38: Error Handling & Log Viewer
**Goal:** Graceful error handling and log access

**Deliverables:**
- Error display components (status bar, modals)
- Retry mechanisms for transient failures
- Graceful degradation (e.g., DB failure → markdown only)
- Log viewer modal (Ctrl+L)
- Show recent logs with filtering

**Test Criteria:**

### Functional Requirements
- [ ] Trigger file write error, see error message
- [ ] Retry button shown for transient errors
- [ ] Database failure, app continues with markdown only
- [ ] Press Ctrl+L, log viewer opens
- [ ] See recent log entries
- [ ] Filter by level (DEBUG, INFO, WARN, ERROR)
- [ ] Close log viewer
- [ ] No crashes from any error condition
- [ ] Error display components (status bar, modals)
- [ ] Retry mechanisms for transient failures
- [ ] Graceful degradation (e.g., DB failure → markdown only)
- [ ] Log viewer modal (Ctrl+L)
- [ ] Show recent logs with filtering

### Integration Requirements
- [ ] Error handling integrates with logging system
- [ ] Error handling integrates with config system
- [ ] Error handling integrates with TUI shell
- [ ] Error handling integrates with database operations
- [ ] Error handling integrates with file I/O operations
- [ ] Log viewer integrates with logging system
- [ ] Log viewer integrates with keyboard input

### Edge Cases & Error Handling
- [ ] Handle multiple concurrent errors
- [ ] Handle errors during error display
- [ ] Handle errors during log viewer display
- [ ] Handle very long error messages (truncate or wrap)
- [ ] Handle very long log entries (truncate or wrap)
- [ ] Handle rapid error events
- [ ] Handle errors during API requests
- [ ] Handle errors during editing
- [ ] Handle errors during modal display
- [ ] Handle log viewer open when already open

### Performance Requirements
- [ ] Error display <50ms
- [ ] Retry mechanism <100ms
- [ ] Graceful degradation <100ms
- [ ] Log viewer open <100ms
- [ ] Log filtering <50ms
- [ ] No UI lag during error handling

### User Experience Requirements
- [ ] Clear error messages in status bar
- [ ] Clear error messages in modals
- [ ] Retry button clearly visible
- [ ] Clear indication of degraded functionality
- [ ] Log viewer easy to navigate
- [ ] Clear filter options
- [ ] Keyboard shortcuts discoverable
- [ ] No visual artifacts during error display
- [ ] Graceful degradation is transparent to user

**Files:** [`internal/errors/handler.go`], [`ui/logviewer/model.go`]

---

## Milestone 39: Context Selector Interface
**Goal:** Implement pluggable context selection algorithm

**Deliverables:**
- ContextSelector interface in `internal/ai/selector.go`
- DefaultSelector implementation
- Scoring algorithm implementation
- Token budget enforcement
- Unit tests for interface and implementation
- Integration tests with library

**Test Criteria:**

### Functional Requirements
- [ ] ContextSelector interface defined with required methods
- [ ] DefaultSelector implements ContextSelector interface
- [ ] Given composition, extract keywords
- [ ] Score all library prompts
- [ ] Top 3-5 prompts selected
- [ ] Relevant prompts ranked higher
- [ ] Token budget respected
- [ ] Explicitly referenced prompts always included
- [ ] Scoring algorithm applied correctly
  - [ ] Tag match: +10 per matching tag
  - [ ] Category match: +5 if same category
  - [ ] Keyword overlap: +1 per matching word
  - [ ] Recently used: +3 if used in last session
  - [ ] Frequently used: +use_count

### Integration Requirements
- [ ] Context selection integrates with library loader
- [ ] Context selection integrates with placeholder parser
- [ ] Context selection integrates with token estimator
- [ ] Context selection integrates with history tracking

### Edge Cases & Error Handling
- [ ] Handle empty composition (no keywords extracted)
- [ ] Handle composition with no matching prompts
- [ ] Handle library with no prompts
- [ ] Handle very large library (1000+ prompts)
- [ ] Handle very long composition (>10000 words)
- [ ] Handle prompts with no metadata
- [ ] Handle prompts with duplicate scores
- [ ] Handle token budget exceeded (select fewer prompts)
- [ ] Handle explicitly referenced prompts exceeding budget

### Performance Requirements
- [ ] Keyword extraction <50ms for 1000-word composition
- [ ] Scoring <100ms for 1000 prompts
- [ ] Selection <10ms
- [ ] Total context selection <200ms
- [ ] Memory usage scales linearly with prompt count

### User Experience Requirements
- [ ] Selected prompts shown in status bar
- [ ] Token usage displayed
- [ ] Clear indication of context limit
- [ ] No UI blocking during selection
- [ ] Progress indicator for large libraries

**Files:** [`internal/ai/selector.go`]

---

## Milestone 40: Domain Events System
**Goal:** Implement domain events for decoupling

**Deliverables:**
- Event interface in `internal/events/events.go`
- Event types (CompositionSaved, PromptUsed, SuggestionAccepted)
- Event dispatcher in `internal/events/dispatcher.go`
- Subscribe/Publish pattern
- Async event handling
- Unit tests for events and dispatcher
- Integration tests with components

**Test Criteria:**

### Functional Requirements
- [ ] Event interface defined with required methods
- [ ] BaseEvent implements Event interface
- [ ] CompositionSavedEvent defined
- [ ] PromptUsedEvent defined
- [ ] SuggestionAcceptedEvent defined
- [ ] Dispatcher manages event subscriptions
- [ ] Subscribe registers handler for event type
- [ ] Publish sends event to all handlers
- [ ] Events handled asynchronously
- [ ] Multiple handlers can subscribe to same event
- [ ] Handlers receive correct event data

### Integration Requirements
- [ ] Events integrate with composition save
- [ ] Events integrate with prompt use
- [ ] Events integrate with suggestion accept
- [ ] Dispatcher integrates with logging system
- [ ] Event handlers integrate with analytics (future)

### Edge Cases & Error Handling
- [ ] Handle handler panics (recover and log)
- [ ] Handle no handlers subscribed (no-op)
- [ ] Handle rapid event publishing (queue properly)
- [ ] Handle handler errors (log and continue)
- [ ] Handle concurrent subscriptions
- [ ] Handle concurrent publishing
- [ ] Handle very large event payloads

### Performance Requirements
- [ ] Subscribe operation <1ms
- [ ] Publish operation <1ms
- [ ] Event handler execution <10ms (typical)
- [ ] No blocking on publish
- [ ] Memory usage scales with handler count

### User Experience Requirements
- [ ] No UI blocking during event handling
- [ ] Errors in handlers don't crash app
- [ ] Events processed in reasonable time
- [ ] No event loss under normal load

**Files:** [`internal/events/events.go`], [`internal/events/dispatcher.go`]

---

## Milestone 41: AI Provider Middleware
**Goal:** Implement middleware pattern for cross-cutting concerns

**Deliverables:**
- ProviderMiddleware type in `internal/ai/middleware.go`
- WithLogging middleware
- WithCaching middleware
- WithMetrics middleware
- Middleware chaining support
- Unit tests for all middleware
- Integration tests with providers

**Test Criteria:**

### Functional Requirements
- [ ] ProviderMiddleware type defined
- [ ] WithLogging wraps provider with logging
- [ ] WithCaching wraps provider with caching
- [ ] WithMetrics wraps provider with metrics
- [ ] Middleware can be chained
- [ ] Middleware preserves provider interface
- [ ] Logging middleware logs all calls
- [ ] Caching middleware caches results
- [ ] Metrics middleware collects metrics

### Integration Requirements
- [ ] Middleware integrates with provider factory
- [ ] Logging middleware integrates with logging system
- [ ] Caching middleware integrates with cache
- [ ] Metrics middleware integrates with metrics system
- [ ] Middleware chain applies in correct order

### Edge Cases & Error Handling
- [ ] Handle middleware errors (propagate to caller)
- [ ] Handle cache misses (call provider)
- [ ] Handle cache invalidation
- [ ] Handle logging failures (don't block)
- [ ] Handle metrics collection failures (don't block)
- [ ] Handle empty middleware chain
- [ ] Handle very long middleware chains

### Performance Requirements
- [ ] Middleware overhead <1ms per call
- [ ] Logging middleware <5ms per call
- [ ] Caching middleware <1ms (hit), <provider time (miss)
- [ ] Metrics middleware <1ms per call
- [ ] No blocking in middleware

### User Experience Requirements
- [ ] No UI blocking from middleware
- [ ] Logging visible in debug mode
- [ ] Caching improves performance
- [ ] Metrics available for monitoring
- [ ] Middleware errors don't crash app

**Files:** [`internal/ai/middleware.go`]

---

## Summary

**Total Milestones:** 41 (updated from 38)

**Milestone Groups:**
1. **Foundation** (1-6): Bootstrap, TUI shell, file I/O, basic editor, auto-save, undo/redo
2. **Library Integration** (7-10): Load library, browse, search, insert prompts
3. **Placeholders** (11-14): Parse, navigate, edit text/list placeholders
4. **History** (15-17): Repository pattern, sync, browser
5. **Commands & Files** (18-22): Command system, palette, file references
6. **Prompt Management** (23-26): Validation, results, creator, editor
7. **AI Integration** (27-33): Provider interface, context selection, tokens, suggestions, diff
8. **Vim Mode** (34-35): State machine, keybindings
9. **Polish** (36-38): Settings, responsive layout, error handling
10. **Scalability** (39-41): Context selector, domain events, middleware (NEW)

**Key Principles:**
- Each milestone is independently testable
- Clear pass/fail criteria for every test
- Incremental complexity with tight integration loops
- Library + placeholders working together early (Milestone 10-14)
- Scalability abstractions integrated early (Milestone 27, 15, 7, 39-41)
- No time estimates - focus on deliverables
- Test-driven approach ensures stability