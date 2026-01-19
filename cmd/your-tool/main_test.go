package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
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
				"your-tool",
				"Usage:",
				"flags",
			},
		},
		{
			name:         "help flag shows help",
			args:         []string{"--help"},
			wantExitCode: 0,
			outputContains: []string{
				"your-tool",
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

var osExit = os.Exit
