package errors

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyledavis/prompt-stack/internal/logging"
	"github.com/kyledavis/prompt-stack/ui/common"
	"go.uber.org/zap"
)

// Handler manages error display and recovery
type Handler struct {
	showModal bool
	modal     common.Modal
	logger    *zap.Logger
}

// NewHandler creates a new error handler
func NewHandler() *Handler {
	return &Handler{
		showModal: false,
	}
}

// NewHandlerWithLogger creates a new error handler with a logger
func NewHandlerWithLogger(logger *zap.Logger) *Handler {
	return &Handler{
		showModal: false,
		logger:    logger,
	}
}

// SetLogger sets the logger for the error handler
func (h *Handler) SetLogger(logger *zap.Logger) {
	h.logger = logger
}

// Handle processes an error and returns appropriate commands
func (h *Handler) Handle(err error) tea.Cmd {
	if err == nil {
		return nil
	}

	// Log the error
	h.logError(err)

	// Check if it's an AppError
	appErr, ok := err.(*AppError)
	if !ok {
		// Wrap generic errors
		appErr = New(ErrorTypeFile, err.Error())
	}

	// Determine display strategy based on severity
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

// logError logs an error to the debug log
func (h *Handler) logError(err error) {
	if err == nil {
		return
	}

	if h.logger == nil {
		// If no logger is set, use the global LogError function
		LogError(err)
		return
	}

	// Log using zap logger
	if appErr, ok := err.(*AppError); ok {
		// Determine log level based on severity
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
		// Log generic errors
		h.logger.Error("Error occurred",
			zap.Error(err))
	}
}

// handleError handles error-level errors
func (h *Handler) handleError(err *AppError) tea.Cmd {
	// Show persistent error in status bar
	cmd := h.createPersistentErrorMessage(err.Message)

	// For critical errors, also show modal
	if h.isCriticalError(err) {
		h.showModal = true
		h.modal = common.ErrorModal("Error", h.formatErrorDetails(err))
		return tea.Batch(cmd, h.showModalCmd())
	}

	return cmd
}

// isCriticalError determines if an error is critical enough to show a modal
func (h *Handler) isCriticalError(err *AppError) bool {
	// Critical errors that require user attention
	switch err.Type {
	case ErrorTypeConfig:
		return true
	case ErrorTypeDatabase:
		return true
	case ErrorTypeFile:
		// File errors are critical if they prevent core functionality
		return true
	default:
		return false
	}
}

// formatErrorDetails formats error details for modal display
func (h *Handler) formatErrorDetails(err *AppError) string {
	details := err.Message

	if err.Details != "" {
		details += "\n\n" + err.Details
	}

	if err.Err != nil {
		details += "\n\n" + fmt.Sprintf("Technical details: %v", err.Err)
	}

	if err.Retryable {
		details += "\n\nThis error can be retried."
	}

	return details
}

// showModalCmd returns a command to show a modal
func (h *Handler) showModalCmd() tea.Cmd {
	return func() tea.Msg {
		return ShowErrorModalMsg{
			Modal: h.modal,
		}
	}
}

// Update handles messages for the error handler
func (h *Handler) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ShowErrorModalMsg:
		h.showModal = true
		h.modal = msg.Modal
	case common.CloseModalMsg:
		h.showModal = false
	case common.ModalActionMsg:
		if msg.Action == "primary" {
			h.showModal = false
		}
	}
	return nil
}

// IsModalVisible returns whether an error modal is currently visible
func (h *Handler) IsModalVisible() bool {
	return h.showModal
}

// GetModal returns the current error modal
func (h *Handler) GetModal() common.Modal {
	return h.modal
}

// Message types for error handling

// ShowErrorModalMsg signals to show an error modal
type ShowErrorModalMsg struct {
	Modal common.Modal
}

// SetStatusMessageMsg is a message to set status bar content
type SetStatusMessageMsg struct {
	Message string
	Type    string
	Timeout time.Duration
}

// Helper functions to create status bar messages

func (h *Handler) createErrorMessage(message string) tea.Cmd {
	return func() tea.Msg {
		return SetStatusMessageMsg{
			Message: message,
			Type:    "error",
			Timeout: 5 * time.Second,
		}
	}
}

func (h *Handler) createWarningMessage(err *AppError) tea.Cmd {
	message := err.Message
	if err.Details != "" {
		message = fmt.Sprintf("%s: %s", err.Message, err.Details)
	}
	return func() tea.Msg {
		return SetStatusMessageMsg{
			Message: message,
			Type:    "warning",
			Timeout: 5 * time.Second,
		}
	}
}

func (h *Handler) createInfoMessage(message string) tea.Cmd {
	return func() tea.Msg {
		return SetStatusMessageMsg{
			Message: message,
			Type:    "info",
			Timeout: 3 * time.Second,
		}
	}
}

func (h *Handler) createSuccessMessage(message string) tea.Cmd {
	return func() tea.Msg {
		return SetStatusMessageMsg{
			Message: message,
			Type:    "success",
			Timeout: 2 * time.Second,
		}
	}
}

func (h *Handler) createPersistentErrorMessage(message string) tea.Cmd {
	return func() tea.Msg {
		return SetStatusMessageMsg{
			Message: message,
			Type:    "error",
			Timeout: 0, // No timeout
		}
	}
}

// Helper functions for common error scenarios

// HandleFileError handles file-related errors
func HandleFileError(operation string, err error) tea.Cmd {
	if err == nil {
		return nil
	}

	appErr := FileError(fmt.Sprintf("Failed to %s", operation), err)
	LogError(appErr)
	return NewHandler().Handle(appErr)
}

// HandleDatabaseError handles database-related errors
func HandleDatabaseError(operation string, err error) tea.Cmd {
	if err == nil {
		return nil
	}

	appErr := DatabaseError(fmt.Sprintf("Failed to %s", operation), err)
	LogError(appErr)
	return NewHandler().Handle(appErr)
}

// HandleAPIError handles API-related errors
func HandleAPIError(operation string, err error) tea.Cmd {
	if err == nil {
		return nil
	}

	appErr := APIError(fmt.Sprintf("Failed to %s", operation), err)
	LogError(appErr)
	return NewHandler().Handle(appErr)
}

// HandleParsingError handles parsing-related errors
func HandleParsingError(operation string, err error) tea.Cmd {
	if err == nil {
		return nil
	}

	appErr := ParsingError(fmt.Sprintf("Failed to %s", operation), err)
	LogError(appErr)
	return NewHandler().Handle(appErr)
}

// HandleConfigError handles config-related errors
func HandleConfigError(operation string, err error) tea.Cmd {
	if err == nil {
		return nil
	}

	appErr := ConfigError(fmt.Sprintf("Failed to %s", operation), err)
	LogError(appErr)
	return NewHandler().Handle(appErr)
}

// HandleValidationError handles validation-related errors
func HandleValidationError(message string) tea.Cmd {
	appErr := ValidationError(message)
	LogError(appErr)
	return NewHandler().Handle(appErr)
}

// HandleAutoSaveError handles auto-save failures with retry logic
func HandleAutoSaveError(err error, retryCount int) tea.Cmd {
	if err == nil {
		return NewHandler().createSuccessMessage("Saved")
	}

	if retryCount < 3 {
		// Retry silently
		return nil
	}

	// Show persistent error after max retries
	appErr := FileError("Auto-save failed after multiple attempts", err).
		WithDetails("Your work is preserved in memory. Please save manually.")
	LogError(appErr)
	return NewHandler().Handle(appErr)
}

// ShowConfigResetPrompt shows a modal prompting for config reset
func ShowConfigResetPrompt(reason string) tea.Cmd {
	message := fmt.Sprintf(
		"Configuration issue detected: %s\n\nWould you like to reset to defaults?\n\nYour current configuration will be backed up.",
		reason,
	)
	modal := common.ConfirmModal("Configuration Error", message)

	return func() tea.Msg {
		return ShowErrorModalMsg{Modal: modal}
	}
}

// ConfigResetMsg is a message to trigger config reset
type ConfigResetMsg struct {
	Reason string
}

// HandleConfigReset handles config reset during runtime
func HandleConfigReset(reason string) tea.Cmd {
	return func() tea.Msg {
		return ConfigResetMsg{Reason: reason}
	}
}

// ShowDatabaseRebuildPrompt shows a modal prompting for database rebuild
func ShowDatabaseRebuildPrompt(reason string) tea.Cmd {
	message := fmt.Sprintf(
		"Database issue detected: %s\n\nWould you like to rebuild the history index?\n\nA backup will be created before rebuilding.",
		reason,
	)
	modal := common.ConfirmModal("Database Error", message)

	return func() tea.Msg {
		return ShowErrorModalMsg{Modal: modal}
	}
}

// ShowRetryPrompt shows a modal with retry option
func ShowRetryPrompt(title, message string) tea.Cmd {
	modal := common.ConfirmModal(title, message)

	return func() tea.Msg {
		return ShowErrorModalMsg{Modal: modal}
	}
}

// LogError logs an error to the debug log (global function for use without handler)
func LogError(err error) {
	if err == nil {
		return
	}

	// Get logger instance
	logger, err := logging.GetLogger()
	if err != nil || logger == nil {
		// If logger is not available, we can't log
		return
	}

	// Log using zap logger
	if appErr, ok := err.(*AppError); ok {
		// Determine log level based on severity
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
		// Log generic errors
		logger.Error("Error occurred",
			zap.Error(err))
	}
}

// LogErrorWithRetry logs an error with retry information
func LogErrorWithRetry(err error, retryCount int, maxRetries int) {
	if err == nil {
		return
	}

	retryInfo := ""
	if retryCount > 0 {
		retryInfo = fmt.Sprintf(" (retry %d/%d)", retryCount, maxRetries)
	}

	LogError(fmt.Errorf("%w%s", err, retryInfo))
}

// CreateTimeoutError creates a timeout error
func CreateTimeoutError(operation string, timeout time.Duration) *AppError {
	return New(ErrorTypeAPI, fmt.Sprintf("%s timed out after %v", operation, timeout)).
		WithRetryable(true)
}

// CreateRateLimitError creates a rate limit error
func CreateRateLimitError(message string, retryAfter time.Duration) *AppError {
	details := fmt.Sprintf("Retry after: %v", retryAfter)
	return New(ErrorTypeAPI, message).
		WithDetails(details).
		WithRetryable(true)
}

// CreateAuthError creates an authentication error
func CreateAuthError(message string) *AppError {
	return New(ErrorTypeAPI, message).
		WithDetails("Please check your API key in settings").
		WithRetryable(false)
}
