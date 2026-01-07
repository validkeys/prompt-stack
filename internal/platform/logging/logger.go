// Package logging provides structured logging infrastructure with file rotation
// and environment variable control.
package logging

import (
	"fmt"
	"os"
	"path/filepath"

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

// New creates a new logger instance with file rotation
func New() (*zap.Logger, error) {
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

	return logger, nil
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
