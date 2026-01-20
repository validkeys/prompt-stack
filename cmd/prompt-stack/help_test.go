package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestHelpCommandExists(t *testing.T) {
	t.Run("help command has required fields", func(t *testing.T) {
		if helpCmd == nil {
			t.Fatal("helpCmd is nil")
		}

		if helpCmd.Use == "" {
			t.Error("help command has empty Use field")
		}

		if helpCmd.Short == "" {
			t.Error("help command has empty Short field")
		}

		if helpCmd.Long == "" {
			t.Error("help command has empty Long field")
		}

		if helpCmd.Run == nil {
			t.Error("help command has nil Run function")
		}
	})

	t.Run("help command follows Cobra conventions", func(t *testing.T) {
		if !strings.HasPrefix(helpCmd.Use, "help") {
			t.Errorf("help command Use should start with 'help', got %q", helpCmd.Use)
		}

		if helpCmd.DisableAutoGenTag {
			t.Error("help command should not disable auto-gen tag")
		}
	})
}

func TestHelpCommandListsCoreCommands(t *testing.T) {
	t.Run("help without args lists all commands", func(t *testing.T) {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"help"})

		osExit = func(code int) {}
		defer func() { osExit = os.Exit }()

		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("Execute() error = %v, want nil", err)
		}

		output := buf.String()

		coreCommands := []string{"init", "plan", "validate", "review", "build"}
		for _, cmd := range coreCommands {
			if !strings.Contains(output, cmd) {
				t.Errorf("help output should contain core command %q", cmd)
			}
		}
	})

	t.Run("help with specific command argument", func(t *testing.T) {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"help", "plan"})

		osExit = func(code int) {}
		defer func() { osExit = os.Exit }()

		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("Execute() error = %v, want nil", err)
		}

		output := buf.String()

		if !strings.Contains(output, "Generate implementation plans") {
			t.Error("help for plan command should contain its description")
		}

		if !strings.Contains(output, "plan") {
			t.Error("help for plan command should contain its name")
		}
	})

	t.Run("help with unknown command shows general help", func(t *testing.T) {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"help", "nonexistent"})

		osExit = func(code int) {}
		defer func() { osExit = os.Exit }()

		_ = rootCmd.Execute()
		output := buf.String()

		if !strings.Contains(output, "Available Commands:") {
			t.Error("help for unknown command should show available commands list")
		}
	})
}

func TestHelpCommandExitStatus(t *testing.T) {
	t.Run("help command exits with status 0", func(t *testing.T) {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"help"})

		err := rootCmd.Execute()

		if err != nil {
			t.Errorf("help command should execute without error, got %v", err)
		}
	})

	t.Run("help with command argument exits with status 0", func(t *testing.T) {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"help", "validate"})

		err := rootCmd.Execute()

		if err != nil {
			t.Errorf("help validate should execute without error, got %v", err)
		}
	})
}

func TestHelpCommandIntegration(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		outputContains []string
	}{
		{
			name: "help shows usage",
			args: []string{"help"},
			outputContains: []string{
				"Usage:",
				"Available Commands:",
			},
		},
		{
			name: "help init shows init command details",
			args: []string{"help", "init"},
			outputContains: []string{
				"Run an interactive interview",
				"milestone requirements",
			},
		},
		{
			name: "help validate shows validate command details",
			args: []string{"help", "validate"},
			outputContains: []string{
				"Validate implementation plans",
			},
		},
		{
			name: "help review shows review command details",
			args: []string{"help", "review"},
			outputContains: []string{
				"Review implementation progress",
			},
		},
		{
			name: "help build shows build command details",
			args: []string{"help", "build"},
			outputContains: []string{
				"Build project components",
				"implementation plan",
			},
		},
		{
			name: "help plan shows plan command details",
			args: []string{"help", "plan"},
			outputContains: []string{
				"Generate implementation plans",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tt.args)

			osExit = func(code int) {}
			defer func() { osExit = os.Exit }()

			err := rootCmd.Execute()
			if err != nil {
				t.Errorf("Execute() error = %v, want nil", err)
			}

			output := buf.String()
			for _, substring := range tt.outputContains {
				if !strings.Contains(output, substring) {
					t.Errorf("Output does not contain expected substring %q\nOutput: %s", substring, output)
				}
			}
		})
	}
}
