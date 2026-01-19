package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewExecutor(t *testing.T) {
	tests := []struct {
		name       string
		workingDir string
		dryRun     bool
	}{
		{
			name:       "standard executor",
			workingDir: "/tmp/test",
			dryRun:     false,
		},
		{
			name:       "dry run executor",
			workingDir: "/tmp/test",
			dryRun:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := NewExecutor(tt.workingDir, tt.dryRun)
			if exec == nil {
				t.Fatal("expected non-nil executor")
			}
			if exec.workingDir != tt.workingDir {
				t.Errorf("workingDir = %v, want %v", exec.workingDir, tt.workingDir)
			}
			if exec.dryRun != tt.dryRun {
				t.Errorf("dryRun = %v, want %v", exec.dryRun, tt.dryRun)
			}
		})
	}
}

func TestValidateInputs(t *testing.T) {
	tests := []struct {
		name    string
		config  ExecutionConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: ExecutionConfig{
				Task:       "test task",
				AIEngine:   "claude",
				MaxRetries: 3,
			},
			wantErr: false,
		},
		{
			name: "empty task",
			config: ExecutionConfig{
				Task: "",
			},
			wantErr: true,
			errMsg:  "task cannot be empty",
		},
		{
			name: "invalid AI engine",
			config: ExecutionConfig{
				Task:     "test task",
				AIEngine: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid AI engine",
		},
		{
			name: "negative max retries",
			config: ExecutionConfig{
				Task:       "test task",
				MaxRetries: -1,
			},
			wantErr: true,
			errMsg:  "max retries cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			vendorDir := filepath.Join(tmpDir, ".your-tool/vendor/ralphy")
			if err := os.MkdirAll(vendorDir, 0755); err != nil {
				t.Fatalf("failed to create vendor directory: %v", err)
			}

			scriptPath := filepath.Join(vendorDir, "ralphy.sh")
			scriptContent := []byte("#!/bin/bash\necho test\n")
			if err := os.WriteFile(scriptPath, scriptContent, 0755); err != nil {
				t.Fatalf("failed to create test script: %v", err)
			}

			tt.config.WorkingDir = tmpDir
			exec := NewExecutor(tmpDir, false)
			err := exec.ValidateInputs(tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInputs() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != nil && tt.errMsg != "" {
				if len(err.Error()) < len(tt.errMsg) || err.Error()[:len(tt.errMsg)] != tt.errMsg {
					t.Errorf("error message = %v, want prefix %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestExecuteDryRun(t *testing.T) {
	tests := []struct {
		name       string
		config     ExecutionConfig
		wantReport bool
		wantErr    bool
	}{
		{
			name: "successful dry run",
			config: ExecutionConfig{
				Task:     "test task",
				AIEngine: "claude",
				DryRun:   true,
			},
			wantReport: true,
			wantErr:    false,
		},
		{
			name: "dry run with all options",
			config: ExecutionConfig{
				Task:       "test task",
				AIEngine:   "opencode",
				DryRun:     true,
				SkipTests:  true,
				SkipLint:   true,
				MaxRetries: 5,
			},
			wantReport: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			yourToolDir := filepath.Join(tmpDir, ".your-tool")
			if err := os.MkdirAll(yourToolDir, 0755); err != nil {
				t.Fatalf("failed to create .your-tool directory: %v", err)
			}

			exec := NewExecutor(tmpDir, true)
			result, err := exec.Execute(tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if !result.Success {
				t.Errorf("expected success = true, got %v", result.Success)
			}

			if tt.wantReport && result.DryRunOutput == "" {
				t.Error("expected dry run output to be non-empty")
			}

			if tt.wantReport {
				reportPath := filepath.Join(tmpDir, reportFile)
				if _, err := os.Stat(reportPath); os.IsNotExist(err) {
					t.Errorf("expected report file at %s", reportPath)
				}
			}
		})
	}
}

func TestBuildCommandArgs(t *testing.T) {
	tests := []struct {
		name     string
		config   ExecutionConfig
		expected []string
	}{
		{
			name: "minimal config",
			config: ExecutionConfig{
				Task: "test task",
			},
			expected: []string{"test task"},
		},
		{
			name: "config with dry run",
			config: ExecutionConfig{
				Task:   "test task",
				DryRun: true,
			},
			expected: []string{"--dry-run", "test task"},
		},
		{
			name: "config with AI engine",
			config: ExecutionConfig{
				Task:     "test task",
				AIEngine: "claude",
			},
			expected: []string{"--claude", "test task"},
		},
		{
			name: "config with skip flags",
			config: ExecutionConfig{
				Task:      "test task",
				SkipTests: true,
				SkipLint:  true,
			},
			expected: []string{"--no-tests", "--no-lint", "test task"},
		},
		{
			name: "full config",
			config: ExecutionConfig{
				Task:      "test task",
				AIEngine:  "opencode",
				SkipTests: true,
				SkipLint:  true,
				DryRun:    true,
			},
			expected: []string{"--dry-run", "--opencode", "--no-tests", "--no-lint", "test task"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := NewExecutor(testWorkDir(), false)
			args := exec.buildCommandArgs(tt.config)

			if len(args) != len(tt.expected) {
				t.Errorf("buildCommandArgs() returned %d args, want %d", len(args), len(tt.expected))
			}

			for i, arg := range args {
				if i < len(tt.expected) && arg != tt.expected[i] {
					t.Errorf("arg[%d] = %v, want %v", i, arg, tt.expected[i])
				}
			}
		})
	}
}

func TestDryRunValidator(t *testing.T) {
	tests := []struct {
		name    string
		config  ExecutionConfig
		wantErr bool
	}{
		{
			name: "valid dry run config",
			config: ExecutionConfig{
				Task:       "test task",
				AIEngine:   "claude",
				MaxRetries: 3,
			},
			wantErr: false,
		},
		{
			name: "empty task",
			config: ExecutionConfig{
				Task: "",
			},
			wantErr: true,
		},
		{
			name: "invalid AI engine",
			config: ExecutionConfig{
				Task:     "test task",
				AIEngine: "invalid",
			},
			wantErr: true,
		},
		{
			name: "negative max retries",
			config: ExecutionConfig{
				Task:       "test task",
				MaxRetries: -1,
			},
			wantErr: true,
		},
		{
			name: "negative timeout",
			config: ExecutionConfig{
				Task:    "test task",
				Timeout: -1 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			if !tt.wantErr {
				vendorDir := filepath.Join(tmpDir, ".your-tool/vendor/ralphy")
				if err := os.MkdirAll(vendorDir, 0755); err != nil {
					t.Fatalf("failed to create vendor directory: %v", err)
				}

				scriptPath := filepath.Join(vendorDir, "ralphy.sh")
				scriptContent := []byte("#!/bin/bash\necho test\n")
				if err := os.WriteFile(scriptPath, scriptContent, 0755); err != nil {
					t.Fatalf("failed to create test script: %v", err)
				}
			}

			tt.config.WorkingDir = tmpDir
			exec := NewExecutor(tmpDir, false)
			validator := exec.NewDryRunValidator()
			err := validator.ValidateConfig(tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDryRunValidatorValidateScriptMaterialization(t *testing.T) {
	tests := []struct {
		name       string
		scriptFile string
		executable bool
		wantErr    bool
	}{
		{
			name:       "executable script exists",
			scriptFile: "ralphy.sh",
			executable: true,
			wantErr:    false,
		},
		{
			name:       "non-executable script",
			scriptFile: "ralphy.sh",
			executable: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			scriptPath := filepath.Join(tmpDir, tt.scriptFile)

			if err := os.WriteFile(scriptPath, []byte("#!/bin/bash\necho test\n"), 0644); err != nil {
				t.Fatalf("failed to create test script: %v", err)
			}

			if tt.executable {
				if err := os.Chmod(scriptPath, 0755); err != nil {
					t.Fatalf("failed to make script executable: %v", err)
				}
			}

			vendorDir := filepath.Join(tmpDir, ".your-tool/vendor/ralphy")
			if err := os.MkdirAll(vendorDir, 0755); err != nil {
				t.Fatalf("failed to create vendor directory: %v", err)
			}

			targetScriptPath := filepath.Join(vendorDir, "ralphy.sh")
			if err := os.Rename(scriptPath, targetScriptPath); err != nil {
				t.Fatalf("failed to move script to vendor dir: %v", err)
			}

			exec := NewExecutor(tmpDir, false)
			validator := exec.NewDryRunValidator()

			err := validator.ValidateScriptMaterialization()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateScriptMaterialization() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDryRun(t *testing.T) {
	tests := []struct {
		name    string
		config  ExecutionConfig
		wantErr bool
	}{
		{
			name: "all validations pass",
			config: ExecutionConfig{
				Task:       "test task",
				AIEngine:   "claude",
				MaxRetries: 3,
				Timeout:    30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "empty task fails validation",
			config: ExecutionConfig{
				Task: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			vendorDir := filepath.Join(tmpDir, ".your-tool/vendor/ralphy")
			if err := os.MkdirAll(vendorDir, 0755); err != nil {
				t.Fatalf("failed to create vendor directory: %v", err)
			}

			scriptPath := filepath.Join(vendorDir, "ralphy.sh")
			scriptContent := []byte("#!/bin/bash\necho test\n")
			if err := os.WriteFile(scriptPath, scriptContent, 0755); err != nil {
				t.Fatalf("failed to create test script: %v", err)
			}

			tt.config.WorkingDir = tmpDir
			exec := NewExecutor(tmpDir, false)
			err := exec.ValidateDryRun(tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDryRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateDryRunReport(t *testing.T) {
	tests := []struct {
		name    string
		config  ExecutionConfig
		wantErr bool
	}{
		{
			name: "generate report successfully",
			config: ExecutionConfig{
				Task:       "test task",
				AIEngine:   "claude",
				SkipTests:  true,
				SkipLint:   true,
				MaxRetries: 3,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tt.config.WorkingDir = tmpDir
			exec := NewExecutor(tmpDir, false)
			validator := exec.NewDryRunValidator()

			report, err := validator.GenerateDryRunReport(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateDryRunReport() error = %v, wantErr %v", err, tt.wantErr)
			}

			if report == nil && !tt.wantErr {
				t.Error("expected non-nil report")
			}

			if report != nil {
				if report.Task != tt.config.Task {
					t.Errorf("report.Task = %v, want %v", report.Task, tt.config.Task)
				}
				if report.AIEngine != tt.config.AIEngine {
					t.Errorf("report.AIEngine = %v, want %v", report.AIEngine, tt.config.AIEngine)
				}
			}
		})
	}
}

func TestExecutionResult(t *testing.T) {
	tests := []struct {
		name   string
		result *ExecutionResult
	}{
		{
			name: "successful result",
			result: &ExecutionResult{
				Success:  true,
				ExitCode: 0,
				Stdout:   "output",
				Stderr:   "",
			},
		},
		{
			name: "failed result",
			result: &ExecutionResult{
				Success:  false,
				ExitCode: 1,
				Stdout:   "partial output",
				Stderr:   "error message",
				Error:    fmt.Errorf("execution failed"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.result
			if result.Success != (result.ExitCode == 0) {
				t.Errorf("Success = %v, but ExitCode = %d (expected 0 for success)", result.Success, result.ExitCode)
			}
		})
	}
}

func testWorkDir() string {
	return ".your-tool"
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
