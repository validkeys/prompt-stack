package validation

import (
	"fmt"
	"os"
	"path/filepath"
)

type ProjectStructureResult struct {
	IsValid  bool
	Errors   []string
	Warnings []string
}

func ValidateProjectStructure(rootDir string) *ProjectStructureResult {
	result := &ProjectStructureResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	requiredDirs := []string{
		"cmd/prompt-stack",
		"internal",
		"pkg",
		"docs",
		".prompt-stack",
	}

	for _, dir := range requiredDirs {
		dirPath := filepath.Join(rootDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			result.IsValid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Missing required directory: %s", dir))
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
		filePath := filepath.Join(rootDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			result.IsValid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Missing required file: %s", file))
		}
	}

	requiredPackages := []string{
		"cmd/prompt-stack",
		"internal/config",
		"internal/database",
		"internal/validation",
		"pkg/executor",
		"pkg/prompt",
	}

	for _, pkg := range requiredPackages {
		pkgPath := filepath.Join(rootDir, pkg)
		if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
			result.IsValid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Missing required package: %s", pkg))
		}
	}

	return result
}
