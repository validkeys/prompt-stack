package executor

import (
	"fmt"
	"os"
	"path/filepath"
)

type DryRunReport struct {
	Timestamp   string
	Task        string
	AIEngine    string
	ScriptPath  string
	Args        []string
	Environment map[string]string
	WorkingDir  string
}

type DryRunValidator struct {
	executor *Executor
}

func (e *Executor) NewDryRunValidator() *DryRunValidator {
	return &DryRunValidator{
		executor: e,
	}
}

func (dv *DryRunValidator) ValidateConfig(config ExecutionConfig) error {
	if config.Task == "" {
		return fmt.Errorf("dry run validation failed: task is required")
	}

	scriptPath := filepath.Join(dv.executor.workingDir, ralphyScriptPath)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("dry run validation failed: ralphy.sh not found at %s", scriptPath)
	}

	if config.AIEngine != "" {
		validEngines := map[string]bool{
			"claude":   true,
			"opencode": true,
			"cursor":   true,
			"codex":    true,
			"qwen":     true,
			"droid":    true,
		}
		if !validEngines[config.AIEngine] {
			return fmt.Errorf("dry run validation failed: invalid AI engine '%s'", config.AIEngine)
		}
	}

	if config.MaxRetries < 0 {
		return fmt.Errorf("dry run validation failed: max retries cannot be negative (got %d)", config.MaxRetries)
	}

	if config.Timeout < 0 {
		return fmt.Errorf("dry run validation failed: timeout cannot be negative (got %v)", config.Timeout)
	}

	return nil
}

func (dv *DryRunValidator) ValidateScriptMaterialization() error {
	scriptPath := filepath.Join(dv.executor.workingDir, ralphyScriptPath)

	fileInfo, err := os.Stat(scriptPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("ralphy.sh not found at %s", scriptPath)
		}
		return fmt.Errorf("failed to stat ralphy.sh: %w", err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("ralphy.sh is a directory, expected a file")
	}

	if fileInfo.Mode().Perm()&0111 == 0 {
		return fmt.Errorf("ralphy.sh is not executable")
	}

	return nil
}

func (dv *DryRunValidator) ValidateWorkingDirectory() error {
	if dv.executor.workingDir == "" {
		return fmt.Errorf("working directory is empty")
	}

	fileInfo, err := os.Stat(dv.executor.workingDir)
	if err != nil {
		return fmt.Errorf("failed to stat working directory: %w", err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("working directory is not a directory")
	}

	return nil
}

func (dv *DryRunValidator) ValidateOutputPaths() error {
	reportPath := filepath.Join(dv.executor.workingDir, reportFile)
	auditPath := filepath.Join(dv.executor.workingDir, auditLogFile)

	reportDir := filepath.Dir(reportPath)
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return fmt.Errorf("failed to create report directory: %w", err)
	}

	auditDir := filepath.Dir(auditPath)
	if err := os.MkdirAll(auditDir, 0755); err != nil {
		return fmt.Errorf("failed to create audit log directory: %w", err)
	}

	return nil
}

func (dv *DryRunValidator) GenerateDryRunReport(config ExecutionConfig) (*DryRunReport, error) {
	report := &DryRunReport{
		Timestamp:   fmt.Sprintf("%v", config.Timeout),
		Task:        config.Task,
		AIEngine:    config.AIEngine,
		ScriptPath:  filepath.Join(dv.executor.workingDir, ralphyScriptPath),
		Args:        dv.executor.buildCommandArgs(config),
		WorkingDir:  config.WorkingDir,
		Environment: make(map[string]string),
	}

	report.Environment["DRY_RUN"] = "true"
	report.Environment["WORKING_DIR"] = dv.executor.workingDir

	return report, nil
}

func (dv *DryRunValidator) ValidateAll(config ExecutionConfig) error {
	if err := dv.ValidateWorkingDirectory(); err != nil {
		return fmt.Errorf("working directory validation failed: %w", err)
	}

	if err := dv.ValidateScriptMaterialization(); err != nil {
		return fmt.Errorf("script materialization validation failed: %w", err)
	}

	if err := dv.ValidateConfig(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	if err := dv.ValidateOutputPaths(); err != nil {
		return fmt.Errorf("output paths validation failed: %w", err)
	}

	return nil
}

func (e *Executor) ValidateDryRun(config ExecutionConfig) error {
	validator := e.NewDryRunValidator()
	return validator.ValidateAll(config)
}
