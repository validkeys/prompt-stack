package errors

import (
	"fmt"
	"time"
)

// ErrorType represents the category of error
type ErrorType string

const (
	ErrorTypeFile       ErrorType = "file"
	ErrorTypeDatabase   ErrorType = "database"
	ErrorTypeAPI        ErrorType = "api"
	ErrorTypeParsing    ErrorType = "parsing"
	ErrorTypeConfig     ErrorType = "config"
	ErrorTypeValidation ErrorType = "validation"
)

// Severity represents the severity level of an error
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// AppError represents an application error with context
type AppError struct {
	Type      ErrorType
	Severity  Severity
	Message   string
	Details   string
	Timestamp time.Time
	Retryable bool
	Err       error // Original error if any
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:      errorType,
		Severity:  SeverityError,
		Message:   message,
		Timestamp: time.Now(),
		Retryable: false,
	}
}

// WithSeverity sets the severity level
func (e *AppError) WithSeverity(severity Severity) *AppError {
	e.Severity = severity
	return e
}

// WithDetails adds additional details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithRetryable marks the error as retryable
func (e *AppError) WithRetryable(retryable bool) *AppError {
	e.Retryable = retryable
	return e
}

// Wrap wraps an existing error
func (e *AppError) Wrap(err error) *AppError {
	e.Err = err
	return e
}

// FileError creates a file-related error
func FileError(message string, err error) *AppError {
	return New(ErrorTypeFile, message).Wrap(err)
}

// DatabaseError creates a database-related error
func DatabaseError(message string, err error) *AppError {
	return New(ErrorTypeDatabase, message).Wrap(err).WithRetryable(true)
}

// APIError creates an API-related error
func APIError(message string, err error) *AppError {
	return New(ErrorTypeAPI, message).Wrap(err).WithRetryable(true)
}

// ParsingError creates a parsing-related error
func ParsingError(message string, err error) *AppError {
	return New(ErrorTypeParsing, message).Wrap(err).WithSeverity(SeverityWarning)
}

// ConfigError creates a config-related error
func ConfigError(message string, err error) *AppError {
	return New(ErrorTypeConfig, message).Wrap(err)
}

// ValidationError creates a validation-related error
func ValidationError(message string) *AppError {
	return New(ErrorTypeValidation, message).WithSeverity(SeverityWarning)
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Retryable
	}
	return false
}

// GetSeverity returns the severity of an error
func GetSeverity(err error) Severity {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Severity
	}
	return SeverityError
}
