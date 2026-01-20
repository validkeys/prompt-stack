package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateProjectStructure(t *testing.T) {
	t.Run("validates correct project structure", func(t *testing.T) {
		tmpDir := t.TempDir()

		requiredDirs := []string{
			"cmd/prompt-stack",
			"internal",
			"docs",
			".prompt-stack",
		}
		for _, dir := range requiredDirs {
			if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
				t.Fatalf("failed to create test directory: %v", err)
			}
		}

		requiredFiles := []string{
			"cmd/prompt-stack/main.go",
			"README.md",
			"Makefile",
			"go.mod",
			"go.sum",
		}
		for _, file := range requiredFiles {
			if err := os.WriteFile(filepath.Join(tmpDir, file), []byte("test"), 0644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}
		}

		requiredPackages := []string{
			"cmd/prompt-stack",
			"internal/config",
			"internal/validation",
			"internal/knowledge/database",
			"internal/cli/prompt",
			"internal/executor",
		}
		for _, pkg := range requiredPackages {
			if err := os.MkdirAll(filepath.Join(tmpDir, pkg), 0755); err != nil {
				t.Fatalf("failed to create test package: %v", err)
			}
		}

		result := ValidateProjectStructure(tmpDir)

		if !result.IsValid {
			t.Errorf("expected valid structure, got errors: %v", result.Errors)
		}

		if len(result.Errors) > 0 {
			t.Errorf("expected no errors, got: %v", result.Errors)
		}
	})

	t.Run("reports missing directories", func(t *testing.T) {
		tmpDir := t.TempDir()

		result := ValidateProjectStructure(tmpDir)

		if result.IsValid {
			t.Error("expected invalid structure for missing directories")
		}

		expectedMissingDirs := []string{
			"cmd/prompt-stack",
			"internal",
			"docs",
			".prompt-stack",
		}

		for _, expectedDir := range expectedMissingDirs {
			found := false
			for _, err := range result.Errors {
				if err == fmt.Sprintf("Missing required directory: %s", expectedDir) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected error for missing directory %s", expectedDir)
			}
		}
	})

	t.Run("reports missing files", func(t *testing.T) {
		tmpDir := t.TempDir()

		requiredDirs := []string{
			"cmd/prompt-stack",
			"internal",
			"docs",
			".prompt-stack",
		}
		for _, dir := range requiredDirs {
			if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
				t.Fatalf("failed to create test directory: %v", err)
			}
		}

		result := ValidateProjectStructure(tmpDir)

		if result.IsValid {
			t.Error("expected invalid structure for missing files")
		}

		expectedMissingFiles := []string{
			"cmd/prompt-stack/main.go",
			"README.md",
			"Makefile",
			"go.mod",
			"go.sum",
		}

		for _, expectedFile := range expectedMissingFiles {
			found := false
			for _, err := range result.Errors {
				if err == fmt.Sprintf("Missing required file: %s", expectedFile) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected error for missing file %s", expectedFile)
			}
		}
	})

	t.Run("reports missing packages", func(t *testing.T) {
		tmpDir := t.TempDir()

		requiredDirs := []string{
			"cmd/prompt-stack",
			"internal",
			"docs",
			".prompt-stack",
		}
		for _, dir := range requiredDirs {
			if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
				t.Fatalf("failed to create test directory: %v", err)
			}
		}

		requiredFiles := []string{
			"cmd/prompt-stack/main.go",
			"README.md",
			"Makefile",
			"go.mod",
			"go.sum",
		}
		for _, file := range requiredFiles {
			if err := os.WriteFile(filepath.Join(tmpDir, file), []byte("test"), 0644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}
		}

		result := ValidateProjectStructure(tmpDir)

		if result.IsValid {
			t.Error("expected invalid structure for missing packages")
		}

		expectedMissingPackages := []string{
			"internal/config",
			"internal/validation",
			"internal/knowledge/database",
			"internal/cli/prompt",
			"internal/executor",
		}

		for _, expectedPkg := range expectedMissingPackages {
			found := false
			for _, err := range result.Errors {
				if err == fmt.Sprintf("Missing required package: %s", expectedPkg) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected error for missing package %s", expectedPkg)
			}
		}
	})

	t.Run("provides clear error messages", func(t *testing.T) {
		tmpDir := t.TempDir()

		result := ValidateProjectStructure(tmpDir)

		for _, errMsg := range result.Errors {
			if errMsg == "" {
				t.Error("error message should not be empty")
			}
			if !strings.Contains(errMsg, "Missing required") {
				t.Errorf("error message should indicate missing item, got: %s", errMsg)
			}
		}
	})
}
