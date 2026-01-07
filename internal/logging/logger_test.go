package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
)

func TestInitialize(t *testing.T) {
	tmp := t.TempDir()

	// Set custom log path for testing
	os.Setenv("HOME", tmp)

	logger, err := Initialize()
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	defer logger.Sync()

	// Verify logger is not nil
	if logger == nil {
		t.Error("Initialize() returned nil logger")
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

func TestGetLogger(t *testing.T) {
	// Test getting logger after initialization
	logger, err := Initialize()
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	defer logger.Sync()

	got, err := GetLogger()
	if err != nil {
		t.Errorf("GetLogger() error = %v", err)
	}
	if got != logger {
		t.Error("GetLogger() returned different logger instance")
	}
}

func TestGetLoggerNotInitialized(t *testing.T) {
	// Reset global logger
	globalLogger = nil

	_, err := GetLogger()
	if err == nil {
		t.Error("GetLogger() should return error when not initialized")
	}
}
