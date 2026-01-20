package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRalphyCommandExists(t *testing.T) {
	if ralphyCmd == nil {
		t.Error("ralphy command is nil")
	}

	if ralphyCmd.Use == "" {
		t.Error("ralphy command has empty Use field")
	}

	if ralphyCmd.Short == "" {
		t.Error("ralphy command has empty Short field")
	}

	if ralphyCmd.RunE == nil {
		t.Error("ralphy command has nil RunE function")
	}
}

func TestRalphyCommandDryRunFlag(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		expectFlag bool
		flagValue  bool
	}{
		{
			name:       "dry-run flag present",
			args:       []string{"--dry-run"},
			expectFlag: true,
			flagValue:  true,
		},
		{
			name:       "dry-run flag absent",
			args:       []string{},
			expectFlag: true,
			flagValue:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := ralphyCmd.Flags().Lookup("dry-run")
			if flag == nil {
				t.Error("dry-run flag not found")
				return
			}

			if !tt.expectFlag {
				t.Error("expected dry-run flag to be present")
				return
			}

			if tt.flagValue {
				if err := ralphyCmd.Flags().Set("dry-run", "true"); err != nil {
					t.Errorf("failed to set dry-run flag: %v", err)
				}
			}
		})
	}
}

func TestRalphyDryRunCreatesReportFiles(t *testing.T) {
	tmpDir := t.TempDir()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Errorf("failed to restore working directory to %q: %v", originalWd, err)
		}
	})

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	if err := os.MkdirAll(".prompt-stack/vendor/ralphy", 0755); err != nil {
		t.Fatalf("failed to create vendor directory: %v", err)
	}

	rootCmd.SetArgs([]string{"ralphy", "--dry-run"})
	rootCmd.SetOut(nil)
	rootCmd.SetErr(nil)

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("ralphy dry-run failed: %v", err)
	}

	reportPath := filepath.Join(tmpDir, ".prompt-stack", "report.txt")
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Error("report.txt was not created")
	}

	auditPath := filepath.Join(tmpDir, ".prompt-stack", "audit.log")
	if _, err := os.Stat(auditPath); os.IsNotExist(err) {
		t.Error("audit.log was not created")
	}

	if content, err := os.ReadFile(reportPath); err != nil {
		t.Errorf("failed to read report.txt: %v", err)
	} else if len(content) == 0 {
		t.Error("report.txt is empty")
	}

	if content, err := os.ReadFile(auditPath); err != nil {
		t.Errorf("failed to read audit.log: %v", err)
	} else if len(content) == 0 {
		t.Error("audit.log is empty")
	}
}

func TestRalphyCommandIntegration(t *testing.T) {
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Errorf("failed to restore working directory to %q: %v", originalWd, err)
		}
	})

	t.Run("ralphy command with --help", func(t *testing.T) {
		tmpDir := t.TempDir()
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to change to temp directory: %v", err)
		}
		t.Cleanup(func() {
			if err := os.Chdir(originalWd); err != nil {
				t.Errorf("failed to restore working directory to %q: %v", originalWd, err)
			}
		})

		if err := os.MkdirAll(".prompt-stack/vendor/ralphy", 0755); err != nil {
			t.Fatalf("failed to create vendor directory: %v", err)
		}

		buf := new(bytes.Buffer)
		rootCmd.SetArgs([]string{"ralphy", "--help"})
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("ralphy --help failed: %v", err)
		}
	})
}
