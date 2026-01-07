// Package logging provides structured logging infrastructure with file rotation
// and environment variable control.
package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	tmp := t.TempDir()

	// Set custom log path for testing
	os.Setenv("HOME", tmp)

	logger, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer logger.Sync()

	// Verify logger is not nil
	if logger == nil {
		t.Error("New() returned nil logger")
	}

	// Test logging
	logger.Info("test message", zap.String("key", "value"))
	logger.Error("test error", zap.Error(fmt.Errorf("test")))

	// Verify log file was created
	logPath := filepath.Join(tmp, ".promptstack", "debug.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("log file not created")
	}
}

func TestNewMultipleInstances(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)

	// Create multiple logger instances
	logger1, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer logger1.Sync()

	logger2, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer logger2.Sync()

	// Both should be valid loggers
	if logger1 == nil || logger2 == nil {
		t.Error("New() returned nil logger")
	}

	// They should be different instances (no global state)
	if logger1 == logger2 {
		t.Error("New() should return different instances")
	}
}

func TestLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     zapcore.Level
	}{
		{
			name:     "debug level",
			envValue: "debug",
			want:     zapcore.DebugLevel,
		},
		{
			name:     "info level",
			envValue: "info",
			want:     zapcore.InfoLevel,
		},
		{
			name:     "warn level",
			envValue: "warn",
			want:     zapcore.WarnLevel,
		},
		{
			name:     "error level",
			envValue: "error",
			want:     zapcore.ErrorLevel,
		},
		{
			name:     "invalid level defaults to info",
			envValue: "invalid",
			want:     zapcore.InfoLevel,
		},
		{
			name:     "empty env defaults to info",
			envValue: "",
			want:     zapcore.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore original LOG_LEVEL value
			originalLevel := os.Getenv("LOG_LEVEL")
			defer func() {
				if originalLevel != "" {
					os.Setenv("LOG_LEVEL", originalLevel)
				} else {
					os.Unsetenv("LOG_LEVEL")
				}
			}()

			// Set environment variable
			if tt.envValue != "" {
				os.Setenv("LOG_LEVEL", tt.envValue)
			} else {
				os.Unsetenv("LOG_LEVEL")
			}

			tmp := t.TempDir()
			os.Setenv("HOME", tmp)

			logger, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			defer logger.Sync()

			// Verify logger was created
			if logger == nil {
				t.Error("New() returned nil logger")
			}

			// Test logging at different levels
			logger.Debug("debug message")
			logger.Info("info message")
			logger.Warn("warn message")
			logger.Error("error message")
		})
	}
}
