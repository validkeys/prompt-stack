package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"
)

const (
	ralphyScriptPath = ".your-tool/vendor/ralphy/ralphy.sh"
	reportFile       = ".your-tool/report.txt"
	auditLogFile     = ".your-tool/audit.log"
)

type Executor struct {
	workingDir string
	dryRun     bool
}

type ExecutionConfig struct {
	Task       string
	Prompt     string
	AIEngine   string
	SkipTests  bool
	SkipLint   bool
	MaxRetries int
	Timeout    time.Duration
	DryRun     bool
	WorkingDir string
}

type ExecutionResult struct {
	Success      bool
	ExitCode     int
	Stdout       string
	Stderr       string
	Duration     time.Duration
	Error        error
	DryRunOutput string
}

type AuditLogEntry struct {
	Timestamp  time.Time
	Task       string
	AIEngine   string
	Success    bool
	ExitCode   int
	Duration   time.Duration
	Error      string
	ReportPath string
}

func NewExecutor(workingDir string, dryRun bool) *Executor {
	return &Executor{
		workingDir: workingDir,
		dryRun:     dryRun,
	}
}

func (e *Executor) Execute(config ExecutionConfig) (*ExecutionResult, error) {
	startTime := time.Now()
	result := &ExecutionResult{
		Duration:     0,
		DryRunOutput: "",
	}

	if config.DryRun || e.dryRun {
		return e.executeDryRun(config, startTime)
	}

	scriptPath := filepath.Join(e.workingDir, ralphyScriptPath)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("ralphy.sh not found at %s", scriptPath)
	}

	args := e.buildCommandArgs(config)

	cmd := exec.Command("/bin/bash", scriptPath)
	cmd.Args = append([]string{"bash"}, args...)
	cmd.Dir = e.workingDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
		result.Stderr = stderr.String()
		result.Stdout = stdout.String()
		result.Error = fmt.Errorf("execution failed: %w", err)
		result.Success = false
	} else {
		result.ExitCode = 0
		result.Stdout = stdout.String()
		result.Stderr = stderr.String()
		result.Success = true
	}

	result.Duration = time.Since(startTime)

	if result.Success {
		if err := e.writeReport(config, result, startTime); err != nil {
			result.Stderr += fmt.Sprintf("\nWarning: failed to write report: %v", err)
		}
	}

	if err := e.logAudit(config, result, startTime); err != nil {
		result.Stderr += fmt.Sprintf("\nWarning: failed to write audit log: %v", err)
	}

	return result, nil
}

func (e *Executor) executeDryRun(config ExecutionConfig, startTime time.Time) (*ExecutionResult, error) {
	dryRunOutput, err := e.generateDryRunReport(config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate dry run report: %w", err)
	}

	reportPath := filepath.Join(e.workingDir, reportFile)
	if err := os.WriteFile(reportPath, []byte(dryRunOutput), 0644); err != nil {
		return nil, fmt.Errorf("failed to write dry run report: %w", err)
	}

	result := &ExecutionResult{
		Success:      true,
		ExitCode:     0,
		Stdout:       dryRunOutput,
		Duration:     time.Since(startTime),
		DryRunOutput: dryRunOutput,
	}

	if err := e.logAudit(config, result, startTime); err != nil {
		result.Stderr = fmt.Sprintf("Warning: failed to write audit log: %v", err)
	}

	return result, nil
}

func (e *Executor) buildCommandArgs(config ExecutionConfig) []string {
	args := []string{}

	if config.DryRun {
		args = append(args, "--dry-run")
	}

	if config.AIEngine != "" {
		args = append(args, "--"+config.AIEngine)
	}

	if config.SkipTests {
		args = append(args, "--no-tests")
	}

	if config.SkipLint {
		args = append(args, "--no-lint")
	}

	if config.Task != "" {
		args = append(args, config.Task)
	}

	return args
}

func (e *Executor) generateDryRunReport(config ExecutionConfig) (string, error) {
	const reportTemplate = `RALPHY EXECUTION - DRY RUN REPORT
========================================
Timestamp: {{.Timestamp}}
Task: {{.Task}}
AI Engine: {{.AIEngine}}
Mode: Dry Run

Configuration:
  Working Directory: {{.WorkingDir}}
  Skip Tests: {{.SkipTests}}
  Skip Lint: {{.SkipLint}}
  Max Retries: {{.MaxRetries}}

Script to Execute:
  Path: {{.ScriptPath}}
  Arguments:
{{range $arg := .Arguments}}    - {{$arg}}
{{end}}

Validation: PASSED
Note: This is a dry run. No actual execution was performed.
`

	type ReportData struct {
		Timestamp  string
		Task       string
		AIEngine   string
		WorkingDir string
		SkipTests  bool
		SkipLint   bool
		MaxRetries int
		ScriptPath string
		Arguments  []string
	}

	data := ReportData{
		Timestamp:  time.Now().Format(time.RFC3339),
		Task:       config.Task,
		AIEngine:   config.AIEngine,
		WorkingDir: config.WorkingDir,
		SkipTests:  config.SkipTests,
		SkipLint:   config.SkipLint,
		MaxRetries: config.MaxRetries,
		ScriptPath: filepath.Join(e.workingDir, ralphyScriptPath),
		Arguments:  e.buildCommandArgs(config),
	}

	tmpl, err := template.New("dryrun-report").Parse(reportTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse report template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (e *Executor) writeReport(config ExecutionConfig, result *ExecutionResult, startTime time.Time) error {
	reportPath := filepath.Join(e.workingDir, reportFile)

	const reportTemplate = `RALPHY EXECUTION REPORT
========================
Timestamp: {{.Timestamp}}
Task: {{.Task}}
AI Engine: {{.AIEngine}}
Duration: {{.Duration}}

Result:
  Status: {{if .Success}}SUCCESS{{else}}FAILED{{end}}
  Exit Code: {{.ExitCode}}

Output:
{{.Output}}

Errors:
{{.Error}}
`

	type ReportData struct {
		Timestamp string
		Task      string
		AIEngine  string
		Duration  string
		Success   bool
		ExitCode  int
		Output    string
		Error     string
	}

	data := ReportData{
		Timestamp: startTime.Format(time.RFC3339),
		Task:      config.Task,
		AIEngine:  config.AIEngine,
		Duration:  result.Duration.String(),
		Success:   result.Success,
		ExitCode:  result.ExitCode,
		Output:    result.Stdout,
		Error:     result.Stderr,
	}

	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse report template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return os.WriteFile(reportPath, buf.Bytes(), 0644)
}

func (e *Executor) logAudit(config ExecutionConfig, result *ExecutionResult, startTime time.Time) error {
	auditPath := filepath.Join(e.workingDir, auditLogFile)

	entry := AuditLogEntry{
		Timestamp:  startTime,
		Task:       config.Task,
		AIEngine:   config.AIEngine,
		Success:    result.Success,
		ExitCode:   result.ExitCode,
		Duration:   result.Duration,
		ReportPath: filepath.Join(e.workingDir, reportFile),
	}

	if result.Error != nil {
		entry.Error = result.Error.Error()
	}

	const auditTemplate = `{{.Timestamp}} | {{.Task}} | {{.AIEngine}} | {{if .Success}}SUCCESS{{else}}FAILED{{end}} | ExitCode: {{.ExitCode}} | Duration: {{.Duration}} | Error: {{.Error}}
`

	tmpl, err := template.New("audit").Parse(auditTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse audit template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, entry); err != nil {
		return fmt.Errorf("failed to execute audit template: %w", err)
	}

	file, err := os.OpenFile(auditPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open audit log: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(buf.String())
	return err
}

func (e *Executor) ValidateInputs(config ExecutionConfig) error {
	if config.Task == "" {
		return fmt.Errorf("task cannot be empty")
	}

	scriptPath := filepath.Join(e.workingDir, ralphyScriptPath)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("ralphy.sh not found at %s", scriptPath)
	}

	validEngines := map[string]bool{
		"claude":   true,
		"opencode": true,
		"cursor":   true,
		"codex":    true,
		"qwen":     true,
		"droid":    true,
	}

	if config.AIEngine != "" && !validEngines[config.AIEngine] {
		return fmt.Errorf("invalid AI engine: %s", config.AIEngine)
	}

	if config.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}

	return nil
}
