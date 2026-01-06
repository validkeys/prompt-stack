package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// Default log file size before rotation (10MB)
	maxLogSize = 10
	// Number of old log files to keep
	maxBackups = 3
	// Maximum age of old log files to keep (days)
	maxAge = 30
)

var (
	// Global logger instance
	globalLogger *zap.Logger
	// Mutex for thread-safe access to global logger
	loggerMutex sync.RWMutex
)

// Initialize creates and configures the zap logger
func Initialize() (*zap.Logger, error) {
	// Get log level from environment variable
	logLevel := getLogLevel()

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create .promptstack directory if it doesn't exist
	promptStackDir := filepath.Join(homeDir, ".promptstack")
	if err := os.MkdirAll(promptStackDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .promptstack directory: %w", err)
	}

	// Configure log rotation
	logFile := filepath.Join(promptStackDir, "debug.log")
	writer := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxLogSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   true,
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		logLevel,
	)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Store global logger instance
	loggerMutex.Lock()
	globalLogger = logger
	loggerMutex.Unlock()

	return logger, nil
}

// GetLogger returns the global logger instance
func GetLogger() (*zap.Logger, error) {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()

	if globalLogger == nil {
		return nil, fmt.Errorf("logger not initialized")
	}

	return globalLogger, nil
}

// getLogLevel determines the log level from environment variable
func getLogLevel() zapcore.Level {
	levelStr := os.Getenv("PROMPTSTACK_LOG_LEVEL")
	if levelStr == "" {
		// Default to INFO level
		return zapcore.InfoLevel
	}

	// Check for debug mode
	if os.Getenv("PROMPTSTACK_DEBUG") == "1" {
		return zapcore.DebugLevel
	}

	// Parse the log level string
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		// If parsing fails, default to INFO
		return zapcore.InfoLevel
	}

	return level
}
