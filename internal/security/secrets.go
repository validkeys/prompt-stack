// secrets — Scans YAML files for embedded secrets and sensitive data.
//
// # Purpose
//
// This package provides security scanning for Ralphy YAML inputs, detecting potential
// embedded secrets such as API keys, tokens, passwords, and other sensitive data.
// It ensures that secrets are not committed to version control and are instead
// referenced via environment variables or secure vault mechanisms.
//
// Features
//
//   - Pattern-based detection of common secret types (API keys, tokens, passwords)
//   - Line-by-line scanning with precise location reporting
//   - Configurable severity levels (critical, high, medium, low)
//   - Support for common secret formats (base64, UUID-like strings, hex strings)
//   - JSON report output for CI/CD integration
//   - Whitelist support for known-safe values
package security

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// Exit codes for predictable script behavior
const (
	ExitSuccess   = 0 // No secrets found
	ExitFailed    = 1 // Secrets detected
	ExitExecution = 2 // Execution error (file I/O, invalid args, etc.)
)

// SecretFinding represents a single detected secret.
type SecretFinding struct {
	Type        string `json:"type"`
	Pattern     string `json:"pattern"`
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	Severity    string `json:"severity"`
	Context     string `json:"context"`
	Description string `json:"description"`
}

// ScanReport represents the complete scan results.
type ScanReport struct {
	ScanType       string          `json:"scan_type"`
	ScannerTool    string          `json:"scanner_tool"`
	FileScanned    string          `json:"file_scanned"`
	ScanTimestamp  string          `json:"scan_timestamp"`
	ScanStatus     string          `json:"scan_status"`
	SecretsFound   int             `json:"secrets_found"`
	Findings       []SecretFinding `json:"findings"`
	Summary        map[string]int  `json:"summary"`
	Recommendation string          `json:"recommendation"`
}

// SecretPattern defines a pattern for detecting secrets.
type SecretPattern struct {
	Type        string
	Regex       *regexp.Regexp
	Severity    string
	Description string
}

// Common secret patterns to detect
var secretPatterns = []SecretPattern{
	{
		Type:        "api_key",
		Regex:       regexp.MustCompile(`(?i)(api[_-]?key|apikey)\s*[:=]\s*['"]?([^'"\s]+)['"]?`),
		Severity:    "critical",
		Description: "API key detected. Use environment variable or vault reference instead.",
	},
	{
		Type:        "secret",
		Regex:       regexp.MustCompile(`(?i)(secret|client[_-]?secret)\s*[:=]\s*['"]?([^'"\s]{16,})['"]?`),
		Severity:    "critical",
		Description: "Secret detected. Use environment variable or vault reference instead.",
	},
	{
		Type:        "password",
		Regex:       regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[:=]\s*['"]?([^'"\s]{8,})['"]?`),
		Severity:    "critical",
		Description: "Password detected. Use environment variable or vault reference instead.",
	},
	{
		Type:        "token",
		Regex:       regexp.MustCompile(`(?i)(token|access[_-]?token)\s*[:=]\s*['"]?([^'"\s]{20,})['"]?`),
		Severity:    "critical",
		Description: "Token detected. Use environment variable or vault reference instead.",
	},
	{
		Type:        "private_key",
		Regex:       regexp.MustCompile(`(?i)(private[_-]?key)\s*[:=]\s*['"]?([^'"\s]{50,})['"]?`),
		Severity:    "critical",
		Description: "Private key detected. Use environment variable or vault reference instead.",
	},
	{
		Type:        "aws_key",
		Regex:       regexp.MustCompile(`(?i)(aws[_-]?(access|secret)[_-]?key)\s*[:=]\s*['"]?(AKIA[A-Z0-9]{16})['"]?`),
		Severity:    "critical",
		Description: "AWS access key detected. Use IAM roles or environment variable instead.",
	},
	{
		Type:        "github_token",
		Regex:       regexp.MustCompile(`(?i)(gh[pou]_|github[_-]?token)\s*[:=]\s*['"]?([^'"\s]{36,})['"]?`),
		Severity:    "critical",
		Description: "GitHub token detected. Use environment variable or secret store instead.",
	},
	{
		Type:        "base64_secret",
		Regex:       regexp.MustCompile(`(?i)(secret|key|token|password)\s*[:=]\s*['"]?([A-Za-z0-9+/]{40,}={0,2})['"]?`),
		Severity:    "high",
		Description: "Potential base64-encoded secret detected.",
	},
	{
		Type:        "uuid_like",
		Regex:       regexp.MustCompile(`(?i)(api[_-]?key|secret|token)\s*[:=]\s*['"]?([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})['"]?`),
		Severity:    "medium",
		Description: "UUID-like string in secret context detected.",
	},
	{
		Type:        "hex_secret",
		Regex:       regexp.MustCompile(`(?i)(secret|key|token|password)\s*[:=]\s*['"]?([0-9a-f]{32,})['"]?`),
		Severity:    "medium",
		Description: "Potential hex-encoded secret detected.",
	},
}

// Safe patterns that should be ignored (known placeholders)
var safePatterns = []*regexp.Regexp{
	regexp.MustCompile(`\{\{\s*\w+\s*\}\}`),          // Template placeholders like {{ ENV_VAR }}
	regexp.MustCompile(`\$[A-Z_]+`),                  // Environment variable references like $API_KEY
	regexp.MustCompile(`%[A-Z_]+%`),                  // Windows-style env vars like %API_KEY%
	regexp.MustCompile(`<\s*[A-Z_]+\s*>`),            // XML-style placeholders like <API_KEY>
	regexp.MustCompile(`\b(vault://|env://)`),        // Vault/env references
	regexp.MustCompile(`\b(secret://|key://)`),       // Secret manager references
	regexp.MustCompile(`\b(NONE|TODO|FIXME)\b`),      // Placeholders
	regexp.MustCompile(`\bxxx+\b`),                   // Redacted values
	regexp.MustCompile(`\b[A-Z_]+(_PLACEHOLDER)?\b`), // Generic placeholders
}

// isSafeValue checks if a matched value is a known-safe placeholder.
func isSafeValue(value string) bool {
	for _, pattern := range safePatterns {
		if pattern.MatchString(value) {
			return true
		}
	}
	return false
}

// ScanFile scans a file for secrets and returns findings.
func ScanFile(filePath string) ([]SecretFinding, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", filePath, err)
	}

	lines := strings.Split(string(content), "\n")
	var findings []SecretFinding

	for lineNum, line := range lines {
		for _, pattern := range secretPatterns {
			matches := pattern.Regex.FindAllStringSubmatchIndex(line, -1)
			for _, match := range matches {
				if len(match) < 4 {
					continue
				}

				matchedText := line[match[0]:match[1]]
				capturedValue := line[match[2]:match[3]]

				if isSafeValue(capturedValue) {
					continue
				}

				finding := SecretFinding{
					Type:        pattern.Type,
					Pattern:     pattern.Regex.String(),
					Line:        lineNum + 1,
					Column:      match[0] + 1,
					Severity:    pattern.Severity,
					Context:     strings.TrimSpace(matchedText),
					Description: pattern.Description,
				}
				findings = append(findings, finding)
			}
		}
	}

	return findings, nil
}

// GenerateReport creates a scan report from findings.
func GenerateReport(filePath string, findings []SecretFinding) ScanReport {
	summary := make(map[string]int)
	for _, finding := range findings {
		summary[finding.Severity]++
		summary[finding.Type]++
	}

	status := "passed"
	recommendation := "No secrets detected. Safe to commit."
	if len(findings) > 0 {
		status = "failed"
		recommendation = "Secrets detected. Use environment variables or vault references instead of embedding secrets."
	}

	return ScanReport{
		ScanType:       "secrets_scan",
		ScannerTool:    "secrets",
		FileScanned:    filePath,
		ScanTimestamp:  time.Now().UTC().Format(time.RFC3339),
		ScanStatus:     status,
		SecretsFound:   len(findings),
		Findings:       findings,
		Summary:        summary,
		Recommendation: recommendation,
	}
}

// ScanSecrets scans a YAML file for embedded secrets and returns exit code and report.
//
// Parameters:
//
//	filePath   - Path to the YAML file to scan
//	outputPath - Path to output JSON report (optional, empty string for no output)
//
// Returns:
//
//	int - Exit code (0=success, 1=secrets detected, 2=execution error)
//	ScanReport - Scan report containing findings
//	error - Details about execution error (nil on scan failure)
func ScanSecrets(filePath, outputPath string) (int, ScanReport, error) {
	findings, err := ScanFile(filePath)
	if err != nil {
		return ExitExecution, ScanReport{}, err
	}

	report := GenerateReport(filePath, findings)

	if outputPath != "" {
		if err := WriteReport(report, outputPath); err != nil {
			return ExitExecution, ScanReport{}, err
		}
	}

	if len(findings) > 0 {
		return ExitFailed, report, nil
	}

	return ExitSuccess, report, nil
}

// WriteReport writes the scan report to a JSON file.
func WriteReport(report ScanReport, outputPath string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write report to %q: %w", outputPath, err)
	}

	return nil
}
