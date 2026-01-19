// scan_secrets — Scans YAML files for embedded secrets and sensitive data.
//
// # Purpose
//
// This tool provides security scanning for Ralphy YAML inputs, detecting potential
// embedded secrets such as API keys, tokens, passwords, and other sensitive data.
// It ensures that secrets are not committed to version control and are instead
// referenced via environment variables or secure vault mechanisms.
//
// Usage
//
//	go run scan_secrets.go --file <input.yaml>
//
//	go run scan_secrets.go -f docs/implementation-plan/m0/ralphy_inputs.yaml
//
// Exit Codes
//
//	0 - No secrets found (scan passed)
//	1 - Secrets detected (scan failed)
//	2 - Error in command execution (file not found, invalid format, etc.)
//
// Features
//
//   - Pattern-based detection of common secret types (API keys, tokens, passwords)
//   - Line-by-line scanning with precise location reporting
//   - Configurable severity levels (critical, high, medium, low)
//   - Support for common secret formats (base64, UUID-like strings, hex strings)
//   - JSON report output for CI/CD integration
//   - Whitelist support for known-safe values
//
// Integration Points
//
//   - CI/CD pipelines: Security gate before merge
//   - Pre-commit hooks: Prevent secret commits
//   - Build Mode: Validate generated Ralphy YAML before execution
//   - Security audits: Regular scanning of YAML files
//
// Security Considerations
//
//   - Reads only local files specified via command-line arguments
//   - Does not execute external commands or evaluate arbitrary code
//   - Pattern matching is heuristics-based and may have false positives
//   - Does not connect to external services for validation
//
// Example Workflows
//
//	CI Pipeline:
//	  - name: Scan for secrets
//	    run: |
//	      cd tools
//	      go run scan_secrets.go \
//	        --file ../docs/implementation-plan/m0/ralphy_inputs.yaml
//
//	Makefile Target:
//	  scan-secrets:
//	      cd tools && go run scan_secrets.go --file $(FILE)
//
//	Pre-commit Hook:
//	  if git diff --cached --name-only | grep -q '\.yaml$'; then
//	    tools/scan_secrets.go --file docs/implementation-plan/m0/ralphy_inputs.yaml
//	  fi
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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

// scanFile scans a file for secrets and returns findings.
func scanFile(filePath string) ([]SecretFinding, error) {
	content, err := ioutil.ReadFile(filePath)
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

// generateReport creates a scan report from findings.
func generateReport(filePath string, findings []SecretFinding) ScanReport {
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
		ScannerTool:    "scan_secrets.go",
		FileScanned:    filePath,
		ScanTimestamp:  time.Now().UTC().Format(time.RFC3339),
		ScanStatus:     status,
		SecretsFound:   len(findings),
		Findings:       findings,
		Summary:        summary,
		Recommendation: recommendation,
	}
}

// parseFlags parses and validates command-line arguments.
func parseFlags() (string, string, error) {
	filePath := flag.String("file", "docs/implementation-plan/m0/ralphy_inputs.yaml", "Path to YAML file to scan")
	outputPath := flag.String("output", "", "Path to output JSON report (optional)")
	flag.Parse()

	if *filePath == "" {
		return "", "", fmt.Errorf("--file flag is required")
	}

	return *filePath, *outputPath, nil
}

// writeReport writes the scan report to a JSON file.
func writeReport(report ScanReport, outputPath string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write report to %q: %w", outputPath, err)
	}

	return nil
}

// main is the entry point for the scan_secrets tool.
//
// Execution flow:
//  1. Parse command-line flags
//  2. Scan the specified file for secrets
//  3. Generate scan report
//  4. Write report to output file if specified
//  5. Print summary and exit with appropriate status code
//
// Error handling:
//   - Errors during setup (file I/O) exit with code 2
//   - Secrets detected exit with code 1 and print detailed findings
//   - Successful scan (no secrets) exits with code 0
//
// See package-level documentation for usage examples and exit code details.
func main() {
	filePath, outputPath, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		fmt.Fprintf(os.Stderr, "Usage: go run scan_secrets.go --file <input.yaml> [--output <report.json>]\n")
		os.Exit(ExitExecution)
	}

	findings, err := scanFile(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ExitExecution)
	}

	report := generateReport(filePath, findings)

	if outputPath != "" {
		if err := writeReport(report, outputPath); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(ExitExecution)
		}
		fmt.Printf("Report written to %s\n", outputPath)
	}

	fmt.Printf("Secrets scan: %s\n", report.ScanStatus)
	fmt.Printf("Secrets found: %d\n", report.SecretsFound)

	if len(findings) > 0 {
		fmt.Printf("\nFindings:\n")
		for _, finding := range findings {
			fmt.Printf("  - [%s] %s (line %d)\n", finding.Severity, finding.Type, finding.Line)
			fmt.Printf("    Context: %s\n", finding.Context)
			fmt.Printf("    %s\n", finding.Description)
		}

		fmt.Printf("\nSummary:\n")
		for severity, count := range report.Summary {
			fmt.Printf("  %s: %d\n", severity, count)
		}

		fmt.Printf("\nRecommendation: %s\n", report.Recommendation)
		os.Exit(ExitFailed)
	}

	fmt.Printf("Recommendation: %s\n", report.Recommendation)
	os.Exit(ExitSuccess)
}
