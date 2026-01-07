package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DefaultLogPath    = ".promptstack/debug.log"
	DefaultMaxSize    = 10 // MB
	DefaultMaxBackups = 3
	DefaultMaxAge     = 30 // days
)

var (
	globalLogger *zap.Logger
	loggerMutex  sync.RWMutex
)

// Initialize sets up the global logger with file rotation
func Initialize() (*zap.Logger, error) {
	// Get home directory
	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create log directory
	logDir := filepath.Join(home, ".promptstack")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Configure log rotation with lumberjack
	logPath := filepath.Join(logDir, "debug.log")
	writer := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    DefaultMaxSize, // megabytes
		MaxBackups: DefaultMaxBackups,
		MaxAge:     DefaultMaxAge, // days
		Compress:   false,
	}

	// Configure encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Get log level from environment
	level := getLogLevel()

	// Create core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(writer),
		level,
	)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Store global logger
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

// getLogLevel reads log level from environment variable
func getLogLevel() zapcore.Level {
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		if l, err := zapcore.ParseLevel(level); err == nil {
			return l
		}
	}
	return zapcore.InfoLevel // Default
}
