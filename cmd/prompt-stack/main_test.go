package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCmd(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantOutput     string
		wantExitCode   int
		outputContains []string
	}{
		{
			name:         "root command shows help",
			args:         []string{},
			wantExitCode: 0,
			outputContains: []string{
				"prompt-stack",
				"Usage:",
				"flags",
			},
		},
		{
			name:         "help flag shows help",
			args:         []string{"--help"},
			wantExitCode: 0,
			outputContains: []string{
				"prompt-stack",
				"generating and validating Ralphy YAML files",
				"flags",
			},
		},
		{
			name:         "short help flag",
			args:         []string{"-h"},
			wantExitCode: 0,
			outputContains: []string{
				"Usage:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)

			osExit = func(code int) {
				if code != tt.wantExitCode {
					t.Errorf("ExitCode() got = %d, want %d", code, tt.wantExitCode)
				}
			}
			defer func() { osExit = os.Exit }()

			rootCmd.SetArgs(tt.args)
			err := rootCmd.Execute()

			output := buf.String()

			if err != nil && tt.wantExitCode == 0 {
				t.Errorf("Execute() error = %v, want nil", err)
			}

			for _, substring := range tt.outputContains {
				if !strings.Contains(output, substring) {
					t.Errorf("Output does not contain expected substring %q\nOutput: %s", substring, output)
				}
			}
		})
	}
}

func TestCommandsExist(t *testing.T) {
	commands := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"plan", planCmd},
		{"validate", validateCmd},
		{"review", reviewCmd},
		{"build", buildCmd},
	}

	for _, tc := range commands {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cmd == nil {
				t.Errorf("%s command is nil", tc.name)
			}

			if tc.cmd.Use == "" {
				t.Errorf("%s command has empty Use field", tc.name)
			}

			if tc.cmd.Short == "" {
				t.Errorf("%s command has empty Short field", tc.name)
			}

			if tc.cmd.Run == nil {
				t.Errorf("%s command has nil Run function", tc.name)
			}
		})
	}
}

func TestAllCommandsAvailable(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	osExit = func(code int) {}
	defer func() { osExit = os.Exit }()

	_ = rootCmd.Execute()
	output := buf.String()

	requiredCommands := []string{"plan", "validate", "review", "build"}
	for _, cmd := range requiredCommands {
		if !strings.Contains(output, cmd) {
			t.Errorf("help output should contain command %q", cmd)
		}
	}
}

func TestCommandsCompile(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		setup func(t *testing.T) func()
	}{
		{"plan command compiles", []string{"plan"}, nil},
		{"validate command compiles", []string{"validate", "--input", "test.yaml"}, func(t *testing.T) func() {
			// Create a temporary YAML file for validation
			tmpDir := t.TempDir()
			yamlPath := filepath.Join(tmpDir, "test.yaml")
			yamlContent := `name: test
description: Test YAML file
tasks:
  - name: test-task
    description: Test task
    implementation: echo "test"`

			if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
				t.Fatalf("failed to create test YAML file: %v", err)
			}

			// Change to temp directory so test.yaml is found
			oldDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("failed to get current directory: %v", err)
			}

			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("failed to change to temp directory: %v", err)
			}

			return func() {
				if err := os.Chdir(oldDir); err != nil {
					t.Errorf("failed to restore working directory to %q: %v", oldDir, err)
				}
			}
		}},
		{"review command compiles", []string{"review"}, nil},
		{"build command compiles", []string{"build"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cleanup func()
			if tt.setup != nil {
				cleanup = tt.setup(t)
				if cleanup != nil {
					defer cleanup()
				}
			}

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tt.args)

			osExit = func(code int) {}
			defer func() { osExit = os.Exit }()

			err := rootCmd.Execute()
			if err != nil {
				t.Errorf("command should compile without error, got %v", err)
			}
		})
	}
}
