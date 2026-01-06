# PromptStack - Requirements Document

## Overview
**PromptStack** is a Go-based CLI tool using Bubble Tea for interactive prompt composition, allowing users to build complex AI prompts by combining reusable prompt templates, adding file references, and receiving AI-powered suggestions.

### Terminology
- **Composition** - The working document where users build prompts (also "composition workspace" or "composition editor")
- **Library** - The global collection of reusable prompt templates stored in `~/.promptstack/data`
- **Prompt** - An individual template file in the library containing reusable text/instructions
- **Placeholder** - A variable in a prompt template that can be filled in (syntax: `{{type:name}}`)

## Core Concepts

### Prompt Library
- **Location**: Global library at `~/.promptstack/data`
- **Initial Setup**: Bundled starter prompts embedded in binary at build time
- **Structure**: Maintains existing folder organization:
  - `/workflows` - Multi-step processes
  - `/commands` - Reusable prompt templates
  - `/prompt-decorations` - Enhancement patterns
  - `/rules` - Guidelines and standards
- **Format**: Markdown files with YAML frontmatter

### Prompt Metadata
Stored in YAML frontmatter:
```yaml
---
title: "Prompt Title"
description: "Optional description"
tags: ["tag1", "tag2"]
---
```

**Required Fields:**
- `title`

**Optional Fields:**
- `description`
- `tags`

### Composition Model
The tool supports flexible composition through:

1. **Linear Concatenation** - Insert prompts sequentially to build up context (e.g., workflow → command → rules)
2. **Template-Based** - Insert prompts containing placeholders, then fill them in with specific values

Both modes can be used together in a single composition. Placeholders are automatically detected when prompts are inserted.

### Placeholder System

**Syntax:** `{{type:name}}`

**Supported Types:**
- `{{text:placeholder_name}}` - Simple text replacement
- `{{list:placeholder_name}}` - Multiple items in markdown bullet list format

**Example:**
```yaml
---
title: "API Documentation Generator"
---
Generate documentation for {{text:endpoint_name}} that handles {{list:http_methods}}.
```

**Validation Rules:**
- **Type:** Must be either `text` or `list` (case-sensitive)
- **Name:** Alphanumeric characters and underscores only (`a-zA-Z0-9_`)
- **Duplicate names:** Not allowed within same prompt - validation error
- **Unclosed/malformed placeholders:** Treated as literal text

**Behavior:**
- Text placeholders: Highlight on tab, allow direct typing to replace
- List placeholders: Support add/remove/edit operations for multiple items
- Tab/Shift+Tab navigation between placeholders

**Loading Prompts with Invalid Placeholders:**
- Validation runs when prompt is loaded from library
- Invalid placeholders treated as literal text (displayed as `{{...}}`)
- Prompt still loads and functions normally
- No error/warning shown during normal use

## Features

### 1. Composition Workspace

**Layout:**
- Split-pane interface:
  - Left: Composition editor
  - Right: AI suggestions panel (toggleable)

**Live Saving:**
- New composition creates timestamped file in `~/.promptstack/data/.history/` (format: `YYYY-MM-DD_HH-MM-SS.md`)
- Auto-saves changes as user composes

**Vim Bindings:**
- Optional vim keybindings support
- Normal mode and insert mode

**Undo/Redo:**
- Supports up to 100 levels of undo history
- Undo history persists in memory across auto-saves
- Cleared when composition is closed
- **Standard mode keybindings:**
  - Ctrl+Z - Undo
  - Ctrl+Y - Redo
- **Vim mode keybindings:**
  - u - Undo
  - Ctrl+R - Redo
- **Action granularity (smart batching):**
  - Continuous typing batched as one action (break on >1 second pause or cursor move)
  - Paste operations = one action each
  - Prompt insertions = one action each
  - Placeholder fills = one action each
  - Action boundaries: cursor movement, mode change, command execution
- **Visual feedback:**
  - Status message at undo/redo boundaries ("No more undo history")
  - Optional undo stack depth in status bar

### 2. Library Browser

**Trigger:** `/` key

**Display:**
- Modal overlay
- Flat list view (all prompts regardless of folder)
- Color-coded labels based on source folder
- Fuzzy search/filter capability

**Actions:**
- `Enter` - Insert at cursor position
- `Ctrl+Enter` - Insert on new line
- `Esc` - Cancel/close modal

**Preview:**
- Shows prompt content before insertion

### 3. Command Palette

**Trigger:** `Space` key

**Features:**
- Fuzzy-filterable command list
- Arrow key navigation

**Commands:**
- Toggle AI panel
- Copy to clipboard (copies plain markdown as-is, including unfilled placeholders)
- Trigger AI suggestions
- Add file reference
- Create new prompt
- Save file
- View history
- Clean up history
- Validate library
- Settings

### 4. File Reference System

**Purpose:** Add references to files in current repository

**Trigger:** Via command palette

**Behavior:**
- Uses built-in Go fuzzy finder (primary method via library like `go-fuzzyfinder`)
- Falls back to system `fzf` if available and configured in preferences
- Searches from tool launch directory downward
- Respects `.gitignore` (walks up to find git root for rules)
- Supports multiple file selection

**Workflow:**
- User selects file(s) via fuzzy finder, provides title for each
- Inserts formatted markdown links: `[title](path/to/file)`

### 5. AI Suggestions

**Provider:** Claude API only (initially)

**Trigger:** Manual via command palette (not real-time)

**Context:**
- Current composition content
- Selected prompts from library (intelligently filtered - see "AI Context Window Management" in Technical Requirements)

**Suggestion Types:**
1. **Prompt Recommendations** - Relevant prompts from library to add
2. **Gap Analysis** - Missing context or information
3. **Formatting** - Better structure or organization
4. **Contradictions** - Conflicting instructions
5. **Clarity** - Unclear or ambiguous instructions
6. **Reformatting** - Alternative ways to structure content

**Presentation:**
- Scrollable list in right panel
- Navigate with arrow keys or vim keys (j/k)
- `a` - Accept/apply selected suggestion (triggers AI to apply changes)
- `d` - Dismiss selected suggestion
- Suggestions remain in list until explicitly dismissed

**Applying Suggestions:**

When user presses `a` to accept a suggestion:

1. **API Request:**
   - Send composition + selected suggestion to Claude
   - Request specific edits/patches (not full replacement)
   - Claude returns structured changes with location and new content

2. **Composition State:**
   - Composition enters read-only mode
   - Display status message: "✨ AI is applying suggestion..."
   - All editing controls disabled
   - Cursor navigation still works

3. **On Success:**
   - Display diff view modal showing changes
   - Side-by-side or inline diff highlighting additions/deletions
   - Options:
     - `Accept` - Apply changes, add to undo history as single action
     - `Reject` - Discard changes, return to normal editing
   - After accept/reject, unlock composition for editing
   - Suggestion remains in list (user can dismiss with `d`)

4. **On Failure:**
   - Keep composition unchanged
   - Unlock composition immediately
   - Display error within suggestion box in AI panel
   - Show error message and retry button
   - Format: "❌ Failed to apply: [error message] [Retry]"

5. **Diff View:**
   - Modal showing side-by-side diff with additions (green) and deletions (red)
   - Accept with Enter, reject with Esc

**Note:** See "AI Context Window Management" in Technical Requirements for details on how library prompts are intelligently selected when sending context to Claude.

### 6. Prompt Creation

**Trigger:** Via command palette

**Workflow:** Interactive guided entry using Bubble Tea components

**Steps:**
1. **Title** (required) - Text input
2. **Description** (optional) - Text input
3. **Tags** (optional) - Comma-delimited with autocomplete from existing tags
4. **Category** (required) - Select from existing folders
5. **Initial Content** - Multi-line text editor

**Result:** Creates new markdown file in selected category folder

### 7. Prompt Editing

**Behavior:**
- Edit prompts from library or history
- Changes require explicit save action (not auto-save)
- Updates original file in place

**Edit/Preview Mode:**
- Toggle between edit and preview modes with `Ctrl+P` (or `Cmd+P` on macOS)
- **Edit Mode:**
  - Raw markdown editor with syntax-aware editing
  - Shows YAML frontmatter and markdown syntax as-is
  - Placeholder highlighting active
- **Preview Mode:**
  - Rendered markdown view using `glow` or similar markdown renderer
  - Shows formatted output: headings, lists, code blocks, emphasis
  - YAML frontmatter hidden or displayed as styled metadata card
  - Read-only view (press `e` or `Ctrl+P` to return to edit mode)
- Mode indicator displayed in status bar ("EDIT" / "PREVIEW")

**Preview Styling:**
- Follows minimal aesthetic: clean typography, subtle styling
- Syntax highlighting for code blocks
- Proper spacing and hierarchy for headings
- Matches overall TUI color scheme

### 8. History Browser

**Purpose:** View and load previous compositions from history

**Trigger:** "View history" command in command palette

**Display:**
- Modal overlay similar to library browser
- List view of all history items
- Fuzzy search/filter capability
- Each item shows:
  - Timestamp (formatted: "Jan 5, 2026 2:30 PM")
  - Working directory where created
  - First 100 characters of content (preview)
  - Character count and line count

**Search:**
- Fuzzy search across: timestamp, working directory, content
- Real-time filtering as user types

**Actions:**
- `Enter` - Load selected composition
- `Delete` - Delete selected history item (with confirmation)
- `Esc` - Close modal

**Loading Behavior:**
- If current composition has unsaved changes: Prompt to save/discard/cancel
- Loaded composition creates new history entry (with same working directory)
- Original history file remains unchanged
- New composition starts auto-saving as new timestamped file

**Data Storage:**
- History metadata stored in SQLite database (`~/.promptstack/data/history.db`)
- Markdown files still created in `.history/` folder (source of truth)
- SQLite acts as index for fast searching and metadata queries

### 9. History Cleanup

**Purpose:** Clean up composition history files in `~/.promptstack/data/.history/`

**Trigger:** "Clean up history" command in command palette

**Workflow:**
1. Display modal showing history statistics:
   - Total number of history files
   - Total size on disk
   - Date range (oldest to newest)
2. Present cleanup options:
   - Delete all history
   - Delete history older than [30/60/90] days
   - Delete all except last [10/25/50] files
   - Custom date range
3. Show preview of files to be deleted
4. Require confirmation before deletion
5. Display result: "Deleted X files, freed Y MB"

**Behavior:**
- No automatic deletion
- Manual cleanup only via command palette
- Deleted files are permanently removed (no trash/recycle bin)
- Corresponding SQLite entries also deleted

### 10. Settings Panel

**Purpose:** Configure application settings

**Trigger:** "Settings" command in command palette

**Display:**
- Modal overlay with form-style interface
- Toggle switches and input fields for settings
- Changes persist to `~/.promptstack/config.yaml` immediately

**Available Settings:**
- **Vim Mode:** Toggle vim keybindings on/off (requires restart to take effect)
- **Use System fzf:** Prefer system fzf over built-in fuzzy finder if available
- **Claude API Key:** Edit API key (masked input)
- **Claude Model:** Select model from dropdown

**Behavior:**
- All changes auto-save to config file
- Settings that require restart show indicator: "(restart required)"
- Validation on inputs (e.g., test API key before saving)
- Keyboard navigation: Tab/Shift+Tab between fields, Space to toggle switches

**Actions:**
- Esc: Close settings panel
- No explicit "Save" button - changes persist immediately

### 11. Library Validation

**Purpose:** Validate all prompt files in library for errors and issues

**Trigger:** "Validate library" command in command palette

**Validation Checks:**
1. **Placeholder syntax:**
   - Invalid placeholder types (not `text` or `list`)
   - Invalid placeholder names (non-alphanumeric/underscore characters)
   - Duplicate placeholder names within same prompt
   - Malformed/unclosed placeholders
2. **YAML frontmatter:**
   - Missing frontmatter entirely
   - Invalid YAML syntax
   - Missing required field: `title`
3. **File size:**
   - Files exceeding 1MB limit
4. **File readability:**
   - Permission errors
   - Encoding issues

**Output:**
- Modal overlay showing validation results
- Grouped by error type
- Each error shows:
  - File path (relative to library root)
  - Line number (if applicable)
  - Error description
  - Severity: Error (red) or Warning (yellow)

**Severity Levels:**
- **Error:** Duplicate placeholders, invalid YAML, missing title, file size exceeded
- **Warning:** Missing optional metadata (description/tags), malformed placeholders

**Example Output:**
```
Library Validation Results
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

❌ ERRORS (2)
commands/api-prompt.md:15
  Duplicate placeholder name: {{text:endpoint}}

workflows/setup.md
  Invalid YAML frontmatter

⚠️  WARNINGS (1)
prompt-decorations/enhance.md:8
  Invalid placeholder type: {{data:values}}

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Total: 2 errors, 1 warning in 3 files
[Dismiss]
```

**Actions:**
- Arrow keys to navigate between errors
- Esc: Close modal

## Configuration

### Global Config
**Location:** `~/.promptstack/config.yaml`

**First Run:**
- Interactive setup if config doesn't exist
- Guided prompts for required fields

**Fields:**
- `claude_api_key` (required) - Claude API key
- `model` (required) - Claude model selection (e.g., "claude-3-sonnet")
- `vim_mode` (default: `false`) - Enable vim keybindings
- `use_system_fzf` (default: `false`) - Prefer system fzf over built-in fuzzy finder

### Initialization
On first launch:
1. Check for global config
2. If missing, run interactive setup
3. Create `~/.promptstack/data/` directory
4. Extract bundled starter prompts from binary to data directory
5. Create `.history/` subfolder
6. Initialize SQLite database (`history.db`) with schema

**Subsequent Launches:**
- Existing library files are never overwritten by bundled prompts
- User modifications and additions persist across updates
- Only extract bundled prompts if data directory is empty or missing

## User Interface

### Vim Support
- Optional vim keybindings (configured via `vim_mode` in config file or Settings panel)
- Requires restart to take effect after toggling
- **Scope:** Applies to composition editor only (not modals or other inputs)
- **Mode indicator:** Display current mode in status bar (INSERT/NORMAL/VISUAL)
- **Supported modes:**
  - Normal mode: Navigation and commands
  - Insert mode: Text editing
  - Visual mode: Text selection
- Vim-style navigation with standard keybindings (h/j/k/l, w/b, etc.)

### Hotkeys

**Global Hotkeys:**
- `/` - Open library browser
- `Space` - Open command palette
- `Tab` / `Shift+Tab` - Navigate placeholders
- `Ctrl+P` (or `Cmd+P` on macOS) - Toggle edit/preview mode (when editing prompts)
- `Esc` - Close modals/cancel

### Visual Design

**Design Philosophy:**
- **Modern and Minimal** - Clean interface with intentional use of whitespace
- **Color Usage** - Color applied sparingly, only for functional emphasis:
  - Category/folder labels in library browser
  - Placeholder highlighting in composition editor
  - Error states (red) and success states (green)
  - Syntax highlighting in preview mode
- **Typography** - Clear hierarchy through font weight and spacing, not decoration
- **No Visual Clutter** - Avoid borders, boxes, and unnecessary dividers

**Layout:**
- Split-pane layout with toggleable AI panel
- Modal overlays for browser and dialogs
- Generous padding and line spacing for readability

### Status Bar
Located at bottom of screen, displays:
- **Character count** (e.g., "1,234 chars")
- **Line count** (e.g., "45 lines")
- **Edit/Preview mode indicator** (when editing prompts: "EDIT" or "PREVIEW")
- **Vim mode indicator** (when vim mode enabled: "INSERT", "NORMAL", "VISUAL")
- **Notifications/warnings** (temporary messages, dismissible)

## Technical Requirements

### Technology Stack
- **Language:** Go
- **TUI Framework:** Bubble Tea
- **Markdown Rendering:** Glamour (library used by `glow`) or similar for rich markdown rendering in TUI
- **File Selection:** Go fuzzy finder library (e.g., `go-fuzzyfinder`) with optional system fzf integration
- **AI Integration:** Claude API
- **Database:** SQLite for history metadata and indexing
- **Embedding:** Go embed (`//go:embed`) for bundling starter prompts in binary
- **Platform Support:** macOS only (primary development and testing platform)

### File Operations
- Read/write markdown files
- Parse YAML frontmatter
- Handle timestamps for history
- Git integration (find root, parse .gitignore)
- Maintain prompt library index (`.index.json`) for AI context optimization

### Character Encoding & Line Endings
- **Encoding:** Auto-detect file encoding on read, convert to UTF-8 internally
- **Line Endings:** Preserve original line endings (CRLF or LF) when writing files
- **New Files:** Create with LF line endings (macOS standard)

### Search & Navigation
- Fuzzy search implementation
- File tree traversal with ignore patterns
- Tag aggregation across library

### AI Context Window Management

To prevent exceeding Claude's context limits, the tool maintains an index and intelligently selects which library prompts to include in AI suggestion requests.

**Index Structure (per prompt):**
- Metadata: title, description, tags, category, file path
- Usage tracking: last_used timestamp, use_count
- Content analysis: word_frequency_map (for keyword matching)

**Selection Algorithm:**
1. Extract keywords from current composition (word frequency analysis)
2. Score each prompt in library:
   - Tag match: +10 points per matching tag
   - Category match: +5 points if same category as prompts in composition
   - Keyword overlap: +1 point per matching word (weighted by frequency)
   - Recently used: +3 points if used in last session
   - Frequently used: +use_count points
3. Sort prompts by score (highest first)
4. Select top N prompts that fit within context budget
5. Always include prompts explicitly referenced in current composition

**Context Budget:**
- Reserve 50% of context window for composition content
- Use remaining 50% for highest-scoring library prompts
- Token estimation: ~4 characters ≈ 1 token
- Dynamically adjust based on composition size

**Index Storage & Updates:**
- Built on library load (startup)
- Incrementally updated when prompts are used/created/edited
- Usage stats persisted to `~/.promptstack/data/.index.json`

### Database Schema

**Location:** `~/.promptstack/data/history.db`

**Purpose:** Index composition history for fast searching and metadata queries

**Schema:**
```sql
CREATE TABLE compositions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  file_path TEXT NOT NULL UNIQUE,           -- Path to markdown file in .history/
  created_at TIMESTAMP NOT NULL,            -- When composition was created
  working_directory TEXT NOT NULL,          -- CWD when composition started
  content TEXT NOT NULL,                    -- Full markdown content
  character_count INTEGER NOT NULL,         -- Total characters
  line_count INTEGER NOT NULL,              -- Total lines
  updated_at TIMESTAMP NOT NULL             -- Last modification time
);

CREATE INDEX idx_created_at ON compositions(created_at);
CREATE INDEX idx_working_directory ON compositions(working_directory);
CREATE INDEX idx_content_fts ON compositions(content);  -- For full-text search
```

**Data Flow:**
1. New composition created → Insert into SQLite + create markdown file
2. Auto-save triggered → Update `content`, `character_count`, `line_count`, `updated_at`
3. History browser opened → Query SQLite for list with metadata
4. Search in history → Query SQLite with full-text search
5. History cleanup → Delete from both SQLite and filesystem

**Sync Strategy:**
- Markdown files are source of truth
- On startup, verify SQLite matches filesystem
- If mismatch detected, rebuild SQLite from markdown files
- Support manual "Rebuild history index" command

## Error Handling

### General Principles
- All errors must be caught to prevent crashes
- Display clear error messages with context
- Preserve user work in memory when operations fail
- Offer retry for transient failures (network, rate limit, 5xx errors)
- Never fail silently

### Error Handling by Type

| Error Type | Display Location | Retryable | Specific Behavior |
|------------|------------------|-----------|-------------------|
| **File Parsing (Invalid YAML/Markdown)** | Status bar warning at startup | No | Load as plain markdown without metadata. Show details modal on 'i' press. Continue loading other files. |
| **Claude API - Setup** | Inline in setup wizard | Yes | Allow re-entry or skip (disables AI features) |
| **Claude API - Rate Limit (429)** | AI suggestions panel | Yes | Show retry button with API message |
| **Claude API - Auth (401)** | AI suggestions panel | No | Display error, suggest checking config |
| **Claude API - Network/5xx** | AI suggestions panel | Yes | Show retry button, track retry count |
| **Auto-Save Failure** | Status bar notification | Yes (auto) | Retry on next change, offer manual save |
| **Prompt Save Failure** | Error modal | Yes | Keep editor open, allow clipboard copy |
| **Database Corruption** | Modal on startup | Yes | Offer "Rebuild index" with progress indicator |
| **Database Query Failure** | Status bar notification | No | Continue with in-memory data, markdown files as fallback |

### Error Message Formats

**File Parsing Warning:**
```
⚠ Warning: 2 files loaded without metadata due to parsing errors. Press 'i' for details.
```

**API Error with Retry:**
```
❌ API Error: Rate limit exceeded
Claude API returned: "You have exceeded your rate limit. Please try again in 30 seconds."

[Retry]  [Dismiss]
```

**Database Rebuild Progress:**
```
Rebuilding history index from markdown files...
Progress: 45/120 files processed

[Cancel]
```

## Security

### API Key Storage
- Stored in plain text in `~/.promptstack/config.yaml`
- Config file created with user-only read/write permissions (0600)
- Note: Consider more secure storage mechanisms in future versions

### Input Validation
- Validate file paths to prevent directory traversal
- Sanitize user input in file references
- Escape special characters in markdown generation

### Prompt Injection
- AI suggestions treated as untrusted content
- User must explicitly accept suggestions (no auto-apply)
- Display suggestions in read-only view before acceptance

## Performance & Limits

### File Size Limits
- **Composition files:** 1MB maximum
- **Referenced files:** 1MB maximum
- **Behavior on limit exceeded:**
  - Display error message with file size
  - Prevent loading/saving
  - Suggest breaking into smaller files

### Library Scale
- No hard limit on number of prompts
- Library browser loads all prompts into memory
- Fuzzy search operates on full library
- Designed for personal use (hundreds of prompts, not thousands)

### Memory Management
- Undo history: ~100 actions in memory
- AI suggestions: Keep in memory until dismissed
- Library: Full library cached after initial load

## Installation & Distribution

### Build Process
- **Source Repository:** Contains Go source code and `starter-prompts/` folder
- **Starter Prompts Location:** `starter-prompts/` directory in repository root
- **Embedding:** Use Go embed to bundle starter prompts into binary at compile time
- **Build Tool:** Standard Go build toolchain

**Starter Prompts Directory Structure:**
```
starter-prompts/
├── workflows/
├── commands/
├── prompt-decorations/
└── rules/
```

### Distribution
- **Method:** Pre-built binaries distributed via GitHub Releases
- **Platform:** macOS binaries (Intel and Apple Silicon)
- **Release Artifacts:**
  - `promptstack-darwin-amd64` (Intel Macs)
  - `promptstack-darwin-arm64` (Apple Silicon Macs)
  - README with installation instructions
  - CHANGELOG

### Installation Instructions
1. Download appropriate binary for architecture from GitHub Releases
2. Make executable: `chmod +x promptstack-darwin-*`
3. Move to PATH: `mv promptstack-darwin-* /usr/local/bin/promptstack`
4. Run: `promptstack`
5. Complete interactive setup on first launch

### Updates
- User downloads new binary and replaces old one
- Existing `~/.promptstack/` directory and user data preserved
- Config file format version checked on startup (future-proofing)

## Feature Priority

### MVP (Minimum Viable Product)
Must-have features for initial release:
- Features 1-4: Composition Workspace, Library Browser, Command Palette, File Reference System
- Feature 6: Prompt Creation
- Feature 7: Prompt Editing
- Feature 8: History Browser
- Settings Panel (Feature 10)
- Basic error handling

### Post-MVP
Features that can be added after initial release:
- Feature 5: AI Suggestions (full implementation)
- Feature 9: History Cleanup
- Feature 11: Library Validation
- Advanced vim mode support
- Comprehensive error recovery

## Future Considerations
Not in scope for initial versions:
- Multiple library locations
- Real-time AI monitoring
- Additional placeholder types beyond text/list
- Multiple AI provider support (beyond Claude)
- Local per-repo libraries
- Export functionality (PDF, HTML, etc.)

## Open Questions
1. Should there be additional keyboard shortcuts for common commands beyond Space for palette?
2. How should the tool handle very large files (approaching 1MB limit) - show warning at 75%?
3. What happens if user edits a library file externally while tool is running?
