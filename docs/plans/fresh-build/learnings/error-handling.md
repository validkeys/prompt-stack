# Error Handling Key Learnings

**Purpose**: Key learnings and implementation patterns for error handling from previous PromptStack implementation.

**Related Milestones**: M1-M38 (All milestones)

**Related Documents**: 
- [`project-structure.md`](../project-structure.md) - Error handling structure
- [`go-style-guide.md`](../go-style-guide.md) - Go error handling patterns
- [`go-testing-guide.md`](../go-testing-guide.md) - Testing error scenarios

---

## Learning Categories

### Category 1: Error Handling Architecture

**Learning**: Create structured error types with severity levels and display strategies

**Problem**: Need consistent error handling across the application with different display strategies.

**Solution**: Structured error types with severity levels

**Implementation Pattern**:
```go
type ErrorType string
const (
    ErrorTypeFile      ErrorType = "file"
    ErrorTypeDatabase  ErrorType = "database"
    ErrorTypeAPI       ErrorType = "api"
    ErrorTypeParsing   ErrorType = "parsing"
    ErrorTypeConfig    ErrorType = "config"
    ErrorTypeValidation ErrorType = "validation"
)

type Severity string
const (
    SeverityError   Severity = "error"
    SeverityWarning Severity = "warning"
    SeverityInfo    Severity = "info"
)

type AppError struct {
    Type      ErrorType
    Severity  Severity
    Message   string
    Details   string
    Timestamp time.Time
    Retryable bool
    Err       error
}

func New(errorType ErrorType, message string) *AppError {
    return &AppError{
        Type:      errorType,
        Severity:  SeverityError,
        Message:   message,
        Timestamp: time.Now(),
        Retryable: false,
    }
}

func (e *AppError) WithDetails(details string) *AppError {
    e.Details = details
    return e
}

func (e *AppError) WithSeverity(severity Severity) *AppError {
    e.Severity = severity
    return e
}

func (e *AppError) WithRetryable(retryable bool) *AppError {
    e.Retryable = retryable
    return e
}

func (e *AppError) WithError(err error) *AppError {
    e.Err = err
    return e
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}
```

**Benefits**:
- Clear categorization of error types
- Severity-based display strategy (modal vs. status bar)
- Retryable flag for transient failures
- Preserves original error for debugging
- Structured logging support

**Lesson**: Create structured error types with severity levels. Use severity to determine display strategy (modal for critical errors, status bar for warnings). Mark retryable errors to enable automatic retry logic.

**Related Milestones**: M1-M38 (All milestones)

**When to Apply**: When implementing error handling across the application

---

### Category 2: Status Bar Component Design

**Learning**: Use message-based updates with auto-dismiss and persistent modes

**Problem**: Need to show transient and persistent status messages.

**Solution**: Message-based updates with timeout support

**Implementation Pattern**:
```go
type Model struct {
    message        string
    messageType    MessageType
    messageTimeout time.Time
    charCount      int
    lineCount      int
    tokenEstimate  int
    vimMode        string
    editMode       string
    showAutoSave   bool
    autoSaveStatus string
}

type SetMessageMsg struct {
    Message string
    Type    MessageType
    Timeout time.Duration
}

func SetErrorMessage(message string) tea.Cmd {
    return func() tea.Msg {
        return SetMessageMsg{
            Message: message,
            Type:    MessageTypeError,
            Timeout: 5 * time.Second,
        }
    }
}

func SetPersistentErrorMessage(message string) tea.Cmd {
    return func() tea.Msg {
        return SetMessageMsg{
            Message: message,
            Type:    MessageTypeError,
            Timeout: 0, // No timeout
        }
    }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case SetMessageMsg:
        m.message = msg.Message
        m.messageType = msg.Type
        if msg.Timeout > 0 {
            m.messageTimeout = time.Now().Add(msg.Timeout)
            return m, tea.Tick(msg.Timeout, func(t time.Time) tea.Msg {
                return ClearMessageMsg{}
            })
        }
        return m, nil
        
    case ClearMessageMsg:
        m.message = ""
        m.messageType = MessageTypeInfo
        return m, nil
    }
    return m, nil
}
```

**Benefits**:
- Message-based updates integrate with Bubble Tea
- Auto-dismiss for transient messages
- Persistent mode for critical errors
- Multiple message types (info, success, warning, error)
- Displays stats and mode indicators

**Lesson**: Use message-based updates for status bar. Support both auto-dismissing (with timeout) and persistent (no timeout) messages. Display contextual information (stats, modes) alongside messages.

**Related Milestones**: M2, M8, M38

**When to Apply**: When implementing status bar components in TUI

---

### Category 3: Modal Component Pattern

**Learning**: Create reusable modal with visibility flag and message-based control

**Problem**: Need consistent error modals across the application.

**Solution**: Reusable modal component with visibility management

**Implementation Pattern**:
```go
type Modal struct {
    title       string
    content     string
    width       int
    height      int
    showButtons bool
    primaryBtn  string
    secondaryBtn string
    focused     bool
}

func NewModal(title, content string) Modal {
    return Modal{
        title:   title,
        content: content,
        focused:  true,
    }
}

func (m Modal) WithButtons(primary, secondary string) Modal {
    m.showButtons = true
    m.primaryBtn = primary
    m.secondaryBtn = secondary
    return m
}

func (m Modal) WithSize(width, height int) Modal {
    m.width = width
    m.height = height
    return m
}

func (m Modal) Update(msg tea.Msg) (Modal, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if !m.showButtons {
            switch msg.String() {
            case "esc", "enter", "q":
                return m, tea.Quit
            }
        } else {
            switch msg.String() {
            case "esc":
                return m, func() tea.Msg { return CloseModalMsg{} }
            case "enter":
                return m, func() tea.Msg { return ModalActionMsg{Action: "primary"} }
            }
        }
    }
    return m, nil
}

func (m Modal) View() string {
    if m.width == 0 || m.height == 0 {
        return ""
    }
    
    content := theme.ModalStyle().
        Width(m.width).
        Height(m.height).
        Render(m.content)
    
    return content
}

func ErrorModal(title, message string) Modal {
    return NewModal(title, message).WithButtons("OK", "")
}

func ConfirmModal(title, message string) Modal {
    return NewModal(title, message).WithButtons("Confirm", "Cancel")
}
```

**Benefits**:
- Reusable across different error types
- Visibility flag for clean UI integration
- Message-based control (CloseModalMsg, ModalActionMsg)
- Helper functions for common modal types
- Centered layout with proper sizing

**Lesson**: Create reusable modal components with visibility flags. Return empty string from View() when not visible. Use helper functions for common modal types (error, warning, confirm). This provides consistent error UI across the application.

**Related Milestones**: M38

**When to Apply**: When implementing modal dialogs in TUI

---

### Category 4: Error Handler Integration

**Learning**: Centralized error handler with display strategy routing

**Problem**: Need consistent error handling with severity-based display.

**Solution**: Centralized handler that routes errors based on severity

**Implementation Pattern**:
```go
type Handler struct {
    showModal bool
    modal    common.Modal
    logger    *zap.Logger
}

func NewHandler() *Handler {
    return &Handler{
        showModal: false,
    }
}

func NewHandlerWithLogger(logger *zap.Logger) *Handler {
    return &Handler{
        showModal: false,
        logger:    logger,
    }
}

func (h *Handler) Handle(err error) tea.Cmd {
    if err == nil {
        return nil
    }

    appErr, ok := err.(*AppError)
    if !ok {
        appErr = New(ErrorTypeFile, err.Error())
    }

    switch appErr.Severity {
    case SeverityError:
        return h.handleError(appErr)
    case SeverityWarning:
        return h.createWarningMessage(appErr)
    case SeverityInfo:
        return h.createInfoMessage(appErr.Message)
    default:
        return h.createErrorMessage(appErr.Message)
    }
}

func (h *Handler) handleError(err *AppError) tea.Cmd {
    cmd := h.createPersistentErrorMessage(err.Message)
    
    if h.isCriticalError(err) {
        h.showModal = true
        h.modal = common.ErrorModal("Error", h.formatErrorDetails(err))
        return tea.Batch(cmd, h.showModalCmd())
    }
    
    return cmd
}

func (h *Handler) isCriticalError(err *AppError) bool {
    // Define what constitutes a critical error
    return err.Type == ErrorTypeConfig || err.Type == ErrorTypeDatabase
}

func (h *Handler) formatErrorDetails(err *AppError) string {
    var builder strings.Builder
    builder.WriteString(err.Message)
    
    if err.Details != "" {
        builder.WriteString("\n\n")
        builder.WriteString(err.Details)
    }
    
    if err.Err != nil {
        builder.WriteString("\n\n")
        builder.WriteString(fmt.Sprintf("Details: %v", err.Err))
    }
    
    return builder.String()
}
```

**Benefits**:
- Centralized error handling logic
- Severity-based display routing
- Critical errors show both status bar and modal
- Helper functions for common error scenarios
- Consistent error display across application

**Lesson**: Create centralized error handler that routes errors based on severity. Critical errors show both status bar message and modal. Helper functions (HandleFileError, HandleDatabaseError, etc.) provide convenient error handling for common scenarios.

**Related Milestones**: M38

**When to Apply**: When implementing error handling across the application

---

### Category 5: Import Cycle Prevention

**Learning**: Avoid import cycles by creating message types in lower-level packages

**Problem**: UI components importing internal packages can create circular dependencies.

**Solution**: Create message types in internal/errors that status bar can handle

**Implementation Pattern**:
```go
// In internal/errors/handler.go
type SetStatusMessageMsg struct {
    Message string
    Type    string
    Timeout time.Duration
}

func (h *Handler) createErrorMessage(message string) tea.Cmd {
    return func() tea.Msg {
        return SetStatusMessageMsg{
            Message: message,
            Type:    "error",
            Timeout: 5 * time.Second,
        }
    }
}

// In ui/statusbar/model.go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case SetStatusMessageMsg:
        // Handle message from error handler
        m.message = msg.Message
        m.messageType = msg.Type
        if msg.Timeout > 0 {
            return m, tea.Tick(msg.Timeout, func(t time.Time) tea.Msg {
                return ClearMessageMsg{}
            })
        }
        return m, nil
    }
    return m, nil
}
```

**Lesson**: Avoid import cycles by creating message types in lower-level packages. UI components handle messages from internal packages without importing them directly. This maintains clean architecture and prevents circular dependencies.

**Related Milestones**: M38

**When to Apply**: When designing package architecture to avoid circular dependencies

---

### Category 6: Error Recovery Strategies

**Learning**: Implement graceful degradation with user-friendly messages

**Problem**: Need to handle errors without losing user work.

**Solution**: Graceful degradation with retry logic and user-friendly messages

**Implementation Pattern**:
```go
// File read errors: Load as plain text, warn in validation
func HandleFileError(operation string, err error) tea.Cmd {
    appErr := FileError(fmt.Sprintf("Failed to %s", operation), err)
    return NewHandler().Handle(appErr)
}

// Auto-save errors: Retry silently, show persistent error after max retries
func HandleAutoSaveError(err error, retryCount int) tea.Cmd {
    if err == nil {
        return NewHandler().createSuccessMessage("Saved")
    }
    
    if retryCount < 3 {
        return nil // Retry silently
    }
    
    appErr := FileError("Auto-save failed after multiple attempts", err).
        WithDetails("Your work is preserved in memory. Please save manually.")
    return NewHandler().Handle(appErr)
}

// Config errors: Prompt for reset
func ShowConfigResetPrompt(reason string) tea.Cmd {
    message := fmt.Sprintf(
        "Configuration issue detected: %s\n\nWould you like to reset to defaults?",
        reason,
    )
    modal := common.ConfirmModal("Configuration Error", message)
    return func() tea.Msg {
        return ShowErrorModalMsg{Modal: modal}
    }
}
```

**Benefits**:
- Preserves user work in memory
- Provides actionable next steps
- Retry logic for transient failures
- User-friendly error messages
- Graceful degradation

**Lesson**: Implement error recovery strategies that preserve user work. Retry transient failures silently. Show persistent errors with actionable next steps. Use modals for critical errors requiring user action. This provides a robust user experience even when things go wrong.

**Related Milestones**: M38

**When to Apply**: When implementing error handling for critical operations

---

### Category 7: Graceful File Read Error Handling

**Learning**: Implement comprehensive error handling with graceful degradation

**Problem**: Need to handle multiple file read failure modes gracefully.

**Solution**: Comprehensive error checking before attempting to read

**Implementation Pattern**:
```go
type LoadError struct {
    FilePath string
    Error    error
    Severity string // "error" or "warning"
}

type Library struct {
    Prompts          map[string]*prompt.Prompt
    Index            *prompt.LibraryIndex
    logger           *zap.Logger
    LoadErrors       []LoadError // Track errors during loading
    ValidationErrors []errors.AppError
}

func readFileGracefully(filePath string, logger *zap.Logger) ([]byte, error) {
    // Check if file exists
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, errors.FileError("File not found", err).
                WithDetails(fmt.Sprintf("The file %s does not exist", filePath))
        }
        return nil, errors.FileError("Failed to access file", err).
            WithDetails(fmt.Sprintf("Could not access file: %s", filePath))
    }

    // Check file size (1MB limit)
    const maxFileSize = 1 << 20 // 1MB
    if fileInfo.Size() > maxFileSize {
        err := errors.FileError("File size exceeds limit", nil).
            WithDetails(fmt.Sprintf("File %s is %.2f MB (max: 1MB)",
                filePath, float64(fileInfo.Size())/(1024*1024)))
        logger.Warn("File size exceeds limit",
            zap.String("path", filePath),
            zap.Int64("size", fileInfo.Size()))
        return nil, err
    }

    // Check file permissions
    if fileInfo.Mode().Perm()&0400 == 0 {
        err := errors.FileError("Permission denied", nil).
            WithDetails(fmt.Sprintf("No read permission for file: %s", filePath))
        logger.Warn("Permission denied", zap.String("path", filePath))
        return nil, err
    }

    // Read file content
    content, err := os.ReadFile(filePath)
    if err != nil {
        // Handle specific error types
        if os.IsPermission(err) {
            return nil, errors.FileError("Permission denied", err).
                WithDetails(fmt.Sprintf("Cannot read file: %s", filePath))
        }
        if errors.Is(err, os.ErrClosed) {
            return nil, errors.FileError("File closed unexpectedly", err).
                WithDetails(fmt.Sprintf("File handle closed: %s", filePath))
        }
        
        // Generic file read error
        return nil, errors.FileError("Failed to read file", err).
            WithDetails(fmt.Sprintf("Error reading file: %s", filePath))
    }

    return content, nil
}
```

**Benefits**:
- Comprehensive error detection (not found, size, permissions, read errors)
- Detailed error messages with context
- Graceful degradation (continues loading other files)
- Error tracking for reporting
- Severity-based handling (error vs warning)
- All errors logged appropriately

**Lesson**: Implement comprehensive file read error handling that checks for multiple failure modes before attempting to read. Use structured errors with details for debugging. Track all errors during batch operations and continue processing other items. This provides robust error handling without stopping the entire operation.

**Related Milestones**: M4, M5, M15

**When to Apply**: When implementing file I/O operations

---

### Category 8: Error Logging Integration

**Learning**: Integrate error logging throughout the application using a global logger pattern

**Problem**: Need consistent error logging across all components.

**Solution**: Global logger with thread-safe access

**Implementation Pattern**:
```go
// In internal/logging/logger.go
var (
    globalLogger *zap.Logger
    loggerMutex sync.RWMutex
)

func Initialize() (*zap.Logger, error) {
    // ... logger setup ...
    logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
    
    // Store global logger instance
    loggerMutex.Lock()
    globalLogger = logger
    loggerMutex.Unlock()
    
    return logger, nil
}

func GetLogger() (*zap.Logger, error) {
    loggerMutex.RLock()
    defer loggerMutex.RUnlock()
    
    if globalLogger == nil {
        return nil, fmt.Errorf("logger not initialized")
    }
    
    return globalLogger, nil
}

// In internal/errors/handler.go
type Handler struct {
    showModal bool
    modal    common.Modal
    logger    *zap.Logger
}

func NewHandlerWithLogger(logger *zap.Logger) *Handler {
    return &Handler{
        showModal: false,
        logger:    logger,
    }
}

func (h *Handler) logError(err error) {
    if err == nil {
        return
    }
    
    if h.logger == nil {
        // Fall back to global LogError function
        LogError(err)
        return
    }
    
    // Log using zap logger with appropriate severity
    if appErr, ok := err.(*AppError); ok {
        switch appErr.Severity {
        case SeverityError:
            h.logger.Error(appErr.Message,
                zap.String("type", string(appErr.Type)),
                zap.String("severity", string(appErr.Severity)),
                zap.String("details", appErr.Details),
                zap.Bool("retryable", appErr.Retryable),
                zap.Error(appErr.Err))
        case SeverityWarning:
            h.logger.Warn(appErr.Message,
                zap.String("type", string(appErr.Type)),
                zap.String("details", appErr.Details),
                zap.Bool("retryable", appErr.Retryable),
                zap.Error(appErr.Err))
        case SeverityInfo:
            h.logger.Info(appErr.Message,
                zap.String("type", string(appErr.Type)),
                zap.String("details", appErr.Details))
        }
    } else {
        h.logger.Error("Error occurred", zap.Error(err))
    }
}

func (h *Handler) Handle(err error) tea.Cmd {
    if err == nil {
        return nil
    }
    
    // Log the error
    h.logError(err)
    
    // ... rest of error handling ...
}

// Global LogError function for use without handler
func LogError(err error) {
    if err == nil {
        return
    }
    
    logger, err := logging.GetLogger()
    if err != nil || logger == nil {
        return
    }
    
    // Log using zap logger with appropriate severity
    if appErr, ok := err.(*AppError); ok {
        switch appErr.Severity {
        case SeverityError:
            logger.Error(appErr.Message,
                zap.String("type", string(appErr.Type)),
                zap.String("severity", string(appErr.Severity)),
                zap.String("details", appErr.Details),
                zap.Bool("retryable", appErr.Retryable),
                zap.Error(appErr.Err))
        case SeverityWarning:
            logger.Warn(appErr.Message,
                zap.String("type", string(appErr.Type)),
                zap.String("details", appErr.Details),
                zap.Bool("retryable", appErr.Retryable),
                zap.Error(appErr.Err))
        case SeverityInfo:
            logger.Info(appErr.Message,
                zap.String("type", string(appErr.Type)),
                zap.String("details", appErr.Details))
        }
    } else {
        logger.Error("Error occurred", zap.Error(err))
    }
}

// Update all helper functions to log errors
func HandleFileError(operation string, err error) tea.Cmd {
    if err == nil {
        return nil
    }
    
    appErr := FileError(fmt.Sprintf("Failed to %s", operation), err)
    LogError(appErr) // Log before handling
    return NewHandler().Handle(appErr)
}
```

**Benefits**:
- All errors automatically logged to debug.log
- Structured logging with severity levels (error, warning, info)
- Thread-safe global logger access
- Both Handler instances and global function can log errors
- Detailed error context (type, severity, details, retryable, stack trace)
- Automatic logging in all error helper functions
- Easy to debug issues with comprehensive error logs

**Lesson**: Integrate error logging throughout the application using a global logger pattern. Store logger instance in logging package with thread-safe access. Add logger field to error handler for instance-based logging. Create global LogError function for use without handler instances. Log all errors with appropriate severity levels and structured fields. This provides comprehensive error tracking for debugging and monitoring.

**Related Milestones**: M1-M38 (All milestones)

**When to Apply**: When implementing error handling and logging across the application

---

## Quick Reference

| Learning | Milestone | Priority |
|----------|-----------|----------|
| Error Handling Architecture | M1-M38 | High |
| Status Bar Component Design | M2, M8, M38 | High |
| Modal Component Pattern | M38 | High |
| Error Handler Integration | M38 | High |
| Import Cycle Prevention | M38 | High |
| Error Recovery Strategies | M38 | High |
| Graceful File Read Error Handling | M4, M5, M15 | High |
| Error Logging Integration | M1-M38 | High |

---

**Last Updated**: 2026-01-07
**Source**: [`key-learnings.md`](../key-learnings.md)