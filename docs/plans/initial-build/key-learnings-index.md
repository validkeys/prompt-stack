# Key Learnings Index - AI Optimized

This index provides quick navigation to all issues, patterns, and learnings documented in key-learnings.md with specific line numbers.

## Core Go & Language Issues

- **Go Embed Limitations** - Line 5: `go:embed` does not support parent directory references (`..`)
- **Zap Logger Structured Fields** - Line 15: Zap requires structured field objects, not string literals
- **Regex Matching in Go** - Line 25: Different regex methods return different data types
- **Go Version Requirements** - Line 49: Some packages require newer Go versions

## Architecture & Design Decisions

- **SQLite Driver Selection** - Line 35: Chose `modernc.org/sqlite` over `github.com/mattn/go-sqlite3`
- **Project Structure Organization** - Line 59: Standard Go project layout with feature-based internal packages
- **Error Handling Patterns** - Line 86: Use `fmt.Errorf` with `%w` for error wrapping
- **Frontmatter Parsing Strategy** - Line 105: Simple string-based parser instead of full YAML library
- **Placeholder Parsing** - Line 129: Regex with position tracking
- **Index Scoring Algorithm** - Line 148: Multi-factor scoring for relevance
- **Validation Strategy** - Line 166: Separate errors and warnings
- **Database Schema Design** - Line 186: Separate tables with FTS5 for search
- **Configuration Management** - Line 218: YAML with validation
- **Starter Prompt Extraction** - Line 247: Version-aware extraction
- **Logging Strategy** - Line 269: Structured logging with rotation

## Bubble Tea & TUI Patterns

- **Bubble Tea Model Implementation** - Line 287: Standard Bubble Tea model structure with Init(), Update(), View()
- **Cursor and Viewport Management** - Line 314: Track both cursor position and viewport offset
- **Auto-save Debouncing with Bubble Tea** - Line 344: Use tea.Tick for timer-based operations
- **Custom Message Types** - Line 382: Define custom message types for async operations
- **Status Bar State Management** - Line 410: Track status with explicit states and auto-clear
- **Text Editor Cursor Positioning** - Line 440: Handle cursor movement across line boundaries
- **File Path Management for History** - Line 463: Timestamp-based file naming with directory creation
- **Lipgloss Styling** - Line 486: Define reusable styles and compose them
- **Centralized Theme System** - Line 511: Single source of truth for all UI colors and styles

## UI Component Patterns

- **Library Browser Implementation** - Line 575: Modal overlay with fuzzy search and preview pane
- **Modal Overlay Pattern** - Line 631: Visibility flag with Show()/Hide() methods
- **Fuzzy Matching Integration** - Line 667: sahilm/fuzzy library usage
- **Split-Pane Layout with Lipgloss** - Line 695: Use lipgloss.JoinHorizontal() for side-by-side panels
- **Keyboard Navigation with Vim Mode** - Line 718: Conditional keybinding based on vim mode flag
- **Message-Based Command Execution** - Line 760: Return custom message from Update() for async operations
- **Command Registry Pattern** - Line 790: Centralized command registration with handler functions
- **Command Palette Implementation** - Line 832: Modal overlay with fuzzy search and message-based execution
- **Command Categorization** - Line 877: Group commands by category for better organization
- **Placeholder Command Handlers** - Line 907: Register commands with placeholder handlers for future implementation

## Error Handling Architecture

- **Error Handling Architecture** - Line 936: Structured error types with severity levels and display strategies
- **Status Bar Component Design** - Line 979: Message-based updates with auto-dismiss and persistent modes
- **Modal Component Pattern** - Line 1034: Reusable modal with visibility flag and message-based control
- **Error Handler Integration** - Line 1089: Centralized error handler with display strategy routing
- **Import Cycle Prevention** - Line 1144: UI components importing internal packages can create circular dependencies
- **Error Recovery Strategies** - Line 1178: Graceful degradation with user-friendly messages
- **Graceful File Read Error Handling** - Line 1227: Comprehensive error handling with graceful degradation
- **Error Logging Integration** - Line 1366: Centralized error logging with global logger access

## Placeholder System

- **Placeholder System Implementation** - Line 1542: Regex-based parsing with position tracking and navigation
- **Placeholder Validation Strategy** - Line 1629: Separate errors and warnings with severity levels
- **Cursor Position Management for Placeholders** - Line 1704: Convert between cursor coordinates and absolute positions
- **Placeholder Highlighting in TUI** - Line 1755: Line-by-line rendering with position-based highlighting
- **Theme Integration for Placeholders** - Line 1806: Centralized style for active placeholder highlighting
- **Text Placeholder Editing Mode** - Line 1830: Vim-style editing mode with state management and value replacement
- **Confirmation Dialog Integration for Destructive Operations** - Line 2043: Multi-step confirmation workflow with type assertion for Bubble Tea models
- **Type Assertion for Bubble Tea Model Updates** - Line 2120: Handling tea.Model interface return type with type assertions
- **AI Applying Indicator and Read-Only Mode** - Line 2143: State-based UI feedback with editing restrictions during async operations
- **Diff Viewer Modal Implementation** - Line 2207: Viewport-based modal with unified diff display and color-coded changes

## Future Considerations

- **Future Considerations** - Line 1523: Potential improvements, technical debt, and architecture decisions to revisit

---

**Total Sections**: 52
**Document Length**: 2195 lines
**Last Updated**: 2026-01-06