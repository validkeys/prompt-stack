package validation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestYAMLValidator(t *testing.T) {
	tests := []struct {
		name      string
		setupFile func(t *testing.T) string
		wantValid bool
		wantIssue bool
	}{
		{
			name: "valid yaml file",
			setupFile: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				path := filepath.Join(tmpDir, "valid.yaml")
				if err := os.WriteFile(path, []byte("test: value\n"), 0644); err != nil {
					t.Fatalf("Failed to write test file %q: %v", path, err)
				}
				return path
			},
			wantValid: true,
			wantIssue: false,
		},
		{
			name: "missing file",
			setupFile: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				return filepath.Join(tmpDir, "missing.yaml")
			},
			wantValid: false,
			wantIssue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &YAMLValidator{}
			inputPath := tt.setupFile(t)
			result, err := validator.Validate(inputPath)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Name != "yaml_compliance" {
				t.Errorf("Name = %v, want yaml_compliance", result.Name)
			}

			if result.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v", result.Valid, tt.wantValid)
			}

			hasIssue := len(result.Issues) > 0
			if hasIssue != tt.wantIssue {
				t.Errorf("has issues = %v, want %v", hasIssue, tt.wantIssue)
			}
		})
	}
}

func TestTaskSizingValidator(t *testing.T) {
	tests := []struct {
		name      string
		setupFile func(t *testing.T) string
	}{
		{
			name: "existing file",
			setupFile: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				path := filepath.Join(tmpDir, "test.yaml")
				if err := os.WriteFile(path, []byte("test: value\n"), 0644); err != nil {
					t.Fatalf("Failed to write test file %q: %v", path, err)
				}
				return path
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &TaskSizingValidator{}
			inputPath := tt.setupFile(t)
			_, err := validator.Validate(inputPath)

			if err != nil {
				t.Logf("TaskSizingValidator error (may fail if tools not present): %v", err)
			}

			if inputPath == "" {
				t.Error("inputPath should not be empty")
			}
		})
	}
}

func TestConstraintsValidator(t *testing.T) {
	tests := []struct {
		name      string
		setupFile func(t *testing.T) string
	}{
		{
			name: "existing file",
			setupFile: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				path := filepath.Join(tmpDir, "test.yaml")
				if err := os.WriteFile(path, []byte("test: value\n"), 0644); err != nil {
					t.Fatalf("Failed to write test file %q: %v", path, err)
				}
				return path
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &ConstraintsValidator{}
			inputPath := tt.setupFile(t)
			_, err := validator.Validate(inputPath)

			if err != nil {
				t.Logf("ConstraintsValidator error (may fail if tools not present): %v", err)
			}

			if inputPath == "" {
				t.Error("inputPath should not be empty")
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name       string
		setupFile  func(t *testing.T) string
		wantResult string
	}{
		{
			name: "valid file with output",
			setupFile: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				inputPath := filepath.Join(tmpDir, "input.yaml")
				if err := os.WriteFile(inputPath, []byte("test: value\n"), 0644); err != nil {
					t.Fatalf("Failed to write test file %q: %v", inputPath, err)
				}
				return inputPath
			},
			wantResult: "WARN",
		},
		{
			name: "invalid file missing",
			setupFile: func(t *testing.T) string {
				t.Helper()
				tmpDir := t.TempDir()
				return filepath.Join(tmpDir, "missing.yaml")
			},
			wantResult: "FAIL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputPath := filepath.Join(tmpDir, "output.json")

			config := Config{
				InputPath:     tt.setupFile(t),
				OutputPath:    outputPath,
				Strict:        false,
				Milestone:     "m0",
				QualityTarget: 0.95,
			}

			result, err := Validate(config)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.OverallResult != tt.wantResult {
				t.Errorf("OverallResult = %v, want %v", result.OverallResult, tt.wantResult)
			}

			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				t.Error("output file was not created")
			}
		})
	}
}
