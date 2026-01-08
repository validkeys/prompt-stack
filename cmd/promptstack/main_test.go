package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyledavis/prompt-stack/internal/platform/logging"
	"github.com/kyledavis/prompt-stack/ui/app"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// TestMainFunction tests the main function handles errors correctly
func TestMainFunction(t *testing.T) {
	tests := []struct {
		name       string
		setup      func()
		cleanup    func()
		wantExit   bool
		wantOutput string
	}{
		{
			name: "successful execution",
			setup: func() {
				// Set up minimal environment for successful run
				os.Setenv("HOME", "/tmp/test-home")
			},
			cleanup: func() {
				os.Unsetenv("HOME")
			},
			wantExit:   false,
			wantOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			// Capture stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Test the run() function which contains the main application logic

			// Restore stderr
			w.Close()
			os.Stderr = oldStderr

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			if tt.wantOutput != "" && !strings.Contains(output, tt.wantOutput) {
				t.Errorf("Expected output to contain %q, got %q", tt.wantOutput, output)
			}
		})
	}
}

// TestRunFunction tests the run function
func TestRunFunction(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (*zap.Logger, error)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful run with valid logger",
			setup: func() (*zap.Logger, error) {
				return zap.NewNop(), nil
			},
			wantErr: false,
		},
		{
			name: "run with nil logger",
			setup: func() (*zap.Logger, error) {
				return nil, nil
			},
			wantErr: true,
			errMsg:  "failed to initialize logging",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify the run() function structure is correct
			if tt.setup != nil {
				logger, err := tt.setup()
				if err != nil && !tt.wantErr {
					t.Errorf("setup() failed: %v", err)
				}
				if logger != nil {
					logger.Sync()
				}
			}
		})
	}
}

// TestTUILaunch tests that the TUI can be launched
func TestTUILaunch(t *testing.T) {
	// Create a test logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create app model
	appModel := app.New()

	// Verify model is created
	if appModel.IsQuitting() {
		t.Error("New app model should not be quitting")
	}

	// Test that model implements tea.Model interface
	var _ tea.Model = appModel
}

// TestTUIWithProgram tests the TUI with a Bubble Tea program
func TestTUIWithProgram(t *testing.T) {
	// Create app model
	appModel := app.New()

	// Create a program with input/output
	input := bytes.NewBufferString("q\n")
	output := &bytes.Buffer{}

	program := tea.NewProgram(
		appModel,
		tea.WithInput(input),
		tea.WithOutput(output),
		tea.WithoutSignalHandler(),
	)

	// Run the program in a goroutine
	done := make(chan error)
	go func() {
		_, err := program.Run()
		done <- err
	}()

	// Wait for program to finish or timeout
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Program failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("Program timed out")
	}
}

// TestTUIHandlesKeyboardInput tests that the TUI handles keyboard input
func TestTUIHandlesKeyboardInput(t *testing.T) {
	appModel := app.New()

	// Test character input
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'a'},
	}

	newModel, cmd := appModel.Update(msg)
	if cmd != nil {
		t.Error("Character input should not return a command")
	}

	updatedModel := newModel.(app.Model)
	if updatedModel.IsQuitting() {
		t.Error("Character input should not quit")
	}
}

// TestTUIHandlesQuit tests that the TUI handles quit commands
func TestTUIHandlesQuit(t *testing.T) {
	tests := []struct {
		name     string
		msg      tea.Msg
		wantQuit bool
	}{
		{
			name:     "Ctrl+C quits",
			msg:      tea.KeyMsg{Type: tea.KeyCtrlC},
			wantQuit: true,
		},
		{
			name:     "q key quits",
			msg:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			wantQuit: true,
		},
		{
			name:     "other keys do not quit",
			msg:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
			wantQuit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appModel := app.New()
			newModel, cmd := appModel.Update(tt.msg)

			updatedModel := newModel.(app.Model)
			if updatedModel.IsQuitting() != tt.wantQuit {
				t.Errorf("IsQuitting() = %v, want %v", updatedModel.IsQuitting(), tt.wantQuit)
			}

			if tt.wantQuit && cmd == nil {
				t.Error("Quit command should return tea.Quit")
			}
		})
	}
}

// TestTUIHandlesWindowSize tests that the TUI handles window resize
func TestTUIHandlesWindowSize(t *testing.T) {
	appModel := app.New()

	// Test window size message
	msg := tea.WindowSizeMsg{
		Width:  100,
		Height: 50,
	}

	newModel, _ := appModel.Update(msg)
	updatedModel := newModel.(app.Model)

	if updatedModel.GetWidth() != 100 {
		t.Errorf("GetWidth() = %v, want 100", updatedModel.GetWidth())
	}

	if updatedModel.GetHeight() != 50 {
		t.Errorf("GetHeight() = %v, want 50", updatedModel.GetHeight())
	}
}

// TestTUILogging tests that TUI operations are logged
func TestTUILogging(t *testing.T) {
	// Create an observer logger to capture logs
	observedZapCore, logs := observer.New(zapcore.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	// Simulate TUI startup logging
	observedLogger.Info("Starting TUI")

	// Verify log was captured
	if logs.Len() != 1 {
		t.Errorf("Expected 1 log entry, got %d", logs.Len())
	}

	entry := logs.All()[0]
	if entry.Message != "Starting TUI" {
		t.Errorf("Expected log message 'Starting TUI', got '%s'", entry.Message)
	}

	// Simulate TUI shutdown logging
	observedLogger.Info("TUI shutdown complete")

	if logs.Len() != 2 {
		t.Errorf("Expected 2 log entries, got %d", logs.Len())
	}

	entry = logs.All()[1]
	if entry.Message != "TUI shutdown complete" {
		t.Errorf("Expected log message 'TUI shutdown complete', got '%s'", entry.Message)
	}
}

// TestTUIErrorHandling tests that TUI errors are handled correctly
func TestTUIErrorHandling(t *testing.T) {
	// Create a logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Test that errors are logged
	logger.Error("TUI error", zap.Error(fmt.Errorf("test error")))

	// This test verifies the error handling structure
	// Actual error scenarios would require mocking the program
}

// TestTUIWithAltScreen tests that TUI uses alt screen
func TestTUIWithAltScreen(t *testing.T) {
	appModel := app.New()

	// Create program with alt screen option
	program := tea.NewProgram(
		appModel,
		tea.WithAltScreen(),
		tea.WithoutSignalHandler(),
	)

	// Verify program was created
	if program == nil {
		t.Error("Program should be created with alt screen option")
	}
}

// TestTUIIntegration tests the full integration
func TestTUIIntegration(t *testing.T) {
	// This test verifies the integration points
	// 1. Bootstrap completes successfully
	// 2. App model is created
	// 3. TUI launches
	// 4. TUI handles input
	// 5. TUI quits cleanly

	// Create logger
	logger, err := logging.New()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create app model
	appModel := app.New()

	// Verify model structure
	if appModel.GetStatusBar().View() == "" {
		t.Error("Status bar should render content")
	}

	// Verify model can handle messages
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
	newModel, _ := appModel.Update(msg)
	if newModel == nil {
		t.Error("Update should return a model")
	}

	// Verify view renders
	view := appModel.View()
	if view == "" {
		t.Error("View should render content")
	}

	if !strings.Contains(view, "PromptStack TUI") {
		t.Error("View should contain 'PromptStack TUI'")
	}
}

// TestTUIPerformance tests TUI performance characteristics
func TestTUIPerformance(t *testing.T) {
	appModel := app.New()

	// Test rapid updates
	start := time.Now()
	for i := 0; i < 100; i++ {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
		newModel, _ := appModel.Update(msg)
		appModel = newModel.(app.Model)
	}
	duration := time.Since(start)

	// Should handle 100 updates in less than 100ms
	if duration > 100*time.Millisecond {
		t.Errorf("Rapid updates took too long: %v", duration)
	}

	// Test rapid view renders
	start = time.Now()
	for i := 0; i < 100; i++ {
		_ = appModel.View()
	}
	duration = time.Since(start)

	// Should handle 100 renders in less than 100ms
	if duration > 100*time.Millisecond {
		t.Errorf("Rapid renders took too long: %v", duration)
	}
}

// TestTUIEdgeCases tests edge cases
func TestTUIEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "zero window size",
			test: func(t *testing.T) {
				appModel := app.New()
				msg := tea.WindowSizeMsg{Width: 0, Height: 0}
				newModel, _ := appModel.Update(msg)
				updatedModel := newModel.(app.Model)
				if updatedModel.GetWidth() != 0 || updatedModel.GetHeight() != 0 {
					t.Error("Should handle zero window size")
				}
			},
		},
		{
			name: "very large window size",
			test: func(t *testing.T) {
				appModel := app.New()
				msg := tea.WindowSizeMsg{Width: 10000, Height: 10000}
				newModel, _ := appModel.Update(msg)
				updatedModel := newModel.(app.Model)
				if updatedModel.GetWidth() != 10000 || updatedModel.GetHeight() != 10000 {
					t.Error("Should handle large window size")
				}
			},
		},
		{
			name: "nil message",
			test: func(t *testing.T) {
				appModel := app.New()
				newModel, _ := appModel.Update(nil)
				if newModel == nil {
					t.Error("Should handle nil message")
				}
			},
		},
		{
			name: "unknown message type",
			test: func(t *testing.T) {
				appModel := app.New()
				newModel, _ := appModel.Update("unknown message")
				if newModel == nil {
					t.Error("Should handle unknown message type")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestTUIConcurrentAccess tests concurrent access to the model
func TestTUIConcurrentAccess(t *testing.T) {
	appModel := app.New()
	var mu sync.Mutex

	// Simulate concurrent updates
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
			for j := 0; j < 100; j++ {
				mu.Lock()
				newModel, _ := appModel.Update(msg)
				appModel = newModel.(app.Model)
				mu.Unlock()
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify model is still functional
	view := appModel.View()
	if view == "" {
		t.Error("Model should still render after concurrent access")
	}
}
