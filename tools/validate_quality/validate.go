package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ValidationReport represents a generic validation report
type ValidationReport struct {
	Valid           bool    `json:"valid,omitempty"`
	OverallResult   string  `json:"overall_result,omitempty"`
	OverallScore    float64 `json:"overall_score,omitempty"`
	ValidationType  string  `json:"validation_type,omitempty"`
	ValidatedAt     string  `json:"validated_at,omitempty"`
	Validator       string  `json:"validator,omitempty"`
	Summary         any     `json:"summary,omitempty"`
	Recommendations []any   `json:"recommendations,omitempty"`
	Violations      []any   `json:"violations,omitempty"`
}

// QualityReport represents the final comprehensive quality report
type QualityReport struct {
	ReportType        string                    `json:"report_type"`
	Timestamp         string                    `json:"timestamp"`
	Milestone         string                    `json:"milestone"`
	RequirementsFile  string                    `json:"requirements_file"`
	GeneratedYAML     string                    `json:"generated_yaml"`
	QualityScore      float64                   `json:"quality_score"`
	ApprovalStatus    string                    `json:"approval_status"`
	ApprovalReason    string                    `json:"approval_reason"`
	Issues            []Issue                   `json:"issues"`
	Recommendations   []Recommendation          `json:"recommendations"`
	ValidationSummary map[string]string         `json:"validation_summary"`
	ComponentScores   map[string]ComponentScore `json:"component_scores"`
}

type Issue struct {
	ID             string `json:"id"`
	Severity       string `json:"severity"`
	Category       string `json:"category"`
	Description    string `json:"description"`
	Impact         string `json:"impact"`
	Recommendation string `json:"recommendation"`
	Status         string `json:"status"`
}

type Recommendation struct {
	Priority string `json:"priority"`
	Action   string `json:"action"`
	Benefit  string `json:"benefit"`
}

type ComponentScore struct {
	Score   float64 `json:"score"`
	Reason  string  `json:"reason"`
	Details string  `json:"details"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run validate.go <validation_reports_dir>")
		os.Exit(1)
	}

	reportsDir := os.Args[1]
	qualityReport, err := generateQualityReport(reportsDir)
	if err != nil {
		fmt.Printf("Error generating quality report: %v\n", err)
		os.Exit(1)
	}

	outputJSON, err := json.MarshalIndent(qualityReport, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling quality report: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(outputJSON))
}

func generateQualityReport(reportsDir string) (*QualityReport, error) {
	// Read all validation reports
	reports, err := readValidationReports(reportsDir)
	if err != nil {
		return nil, err
	}

	// Calculate component scores
	componentScores := calculateComponentScores(reports)

	// Calculate overall quality score
	qualityScore := calculateOverallScore(componentScores)

	// Determine approval status
	approvalStatus, approvalReason := determineApprovalStatus(qualityScore, reports)

	// Generate issues and recommendations
	issues := generateIssues(reports)
	recommendations := generateRecommendations(reports)

	// Create validation summary
	validationSummary := generateValidationSummary(reports)

	return &QualityReport{
		ReportType:        "final_quality_report",
		Timestamp:         time.Now().UTC().Format(time.RFC3339),
		Milestone:         "m0",
		RequirementsFile:  "docs/implementation-plan/m0/requirements.md",
		GeneratedYAML:     "docs/implementation-plan/m0/ralphy_inputs.yaml",
		QualityScore:      qualityScore,
		ApprovalStatus:    approvalStatus,
		ApprovalReason:    approvalReason,
		Issues:            issues,
		Recommendations:   recommendations,
		ValidationSummary: validationSummary,
		ComponentScores:   componentScores,
	}, nil
}

func readValidationReports(reportsDir string) (map[string]ValidationReport, error) {
	reports := make(map[string]ValidationReport)

	// List of expected validation reports
	expectedReports := []string{
		"validation_report_anchors.json",
		"validation_report_constraints.json",
		"validation_report_implementation_guidelines.json",
		"validation_report_sizing.json",
		"validation_report_enforcement.json",
		"yaml_validation_report.json",
		"schema_validation_report.json",
		"secrets_scan_report.json",
	}

	for _, reportName := range expectedReports {
		reportPath := filepath.Join(reportsDir, reportName)
		if _, err := os.Stat(reportPath); os.IsNotExist(err) {
			// Report doesn't exist, skip it
			continue
		}

		data, err := os.ReadFile(reportPath)
		if err != nil {
			return nil, fmt.Errorf("error reading %s: %v", reportName, err)
		}

		var report ValidationReport
		if err := json.Unmarshal(data, &report); err != nil {
			return nil, fmt.Errorf("error unmarshaling %s: %v", reportName, err)
		}

		reports[reportName] = report
	}

	return reports, nil
}

func calculateComponentScores(reports map[string]ValidationReport) map[string]ComponentScore {
	scores := make(map[string]ComponentScore)

	// Style anchors score
	if report, ok := reports["validation_report_anchors.json"]; ok {
		scores["style_anchors"] = ComponentScore{
			Score:   report.OverallScore,
			Reason:  "Style anchors compliance validation",
			Details: fmt.Sprintf("Overall score: %.2f, Result: %s", report.OverallScore, report.OverallResult),
		}
	}

	// Constraints score
	if report, ok := reports["validation_report_constraints.json"]; ok {
		score := 0.7 // Default score for constraints
		if report.Valid {
			score = 0.85
		}
		scores["affirmative_constraints"] = ComponentScore{
			Score:   score,
			Reason:  "Affirmative constraints validation",
			Details: fmt.Sprintf("Valid: %v", report.Valid),
		}
	}

	// Implementation guidelines score
	if report, ok := reports["validation_report_implementation_guidelines.json"]; ok {
		score := 0.6 // Default score
		if report.Valid {
			score = 0.8
		}
		scores["implementation_guidelines"] = ComponentScore{
			Score:   score,
			Reason:  "Implementation guidelines validation",
			Details: fmt.Sprintf("Valid: %v", report.Valid),
		}
	}

	// Task sizing score
	if report, ok := reports["validation_report_sizing.json"]; ok {
		scores["task_sizing"] = ComponentScore{
			Score:   report.OverallScore,
			Reason:  "Task sizing compliance validation",
			Details: fmt.Sprintf("Overall score: %.2f, Result: %s", report.OverallScore, report.OverallResult),
		}
	}

	// Enforcement score
	if report, ok := reports["validation_report_enforcement.json"]; ok {
		scores["multi_layer_enforcement"] = ComponentScore{
			Score:   report.OverallScore,
			Reason:  "Multi-layer enforcement validation",
			Details: fmt.Sprintf("Overall score: %.2f, Result: %s", report.OverallScore, report.OverallResult),
		}
	}

	// YAML syntax score
	if report, ok := reports["yaml_validation_report.json"]; ok {
		score := 0.9
		if report.Valid {
			score = 1.0
		}
		scores["yaml_syntax"] = ComponentScore{
			Score:   score,
			Reason:  "YAML syntax validation",
			Details: fmt.Sprintf("Valid: %v", report.Valid),
		}
	}

	// Schema validation score
	if report, ok := reports["schema_validation_report.json"]; ok {
		score := 0.9
		if report.Valid {
			score = 1.0
		}
		scores["schema_compliance"] = ComponentScore{
			Score:   score,
			Reason:  "JSON Schema compliance validation",
			Details: fmt.Sprintf("Valid: %v", report.Valid),
		}
	}

	// Secrets scan score
	if report, ok := reports["secrets_scan_report.json"]; ok {
		score := 0.9
		if report.Valid {
			score = 1.0
		}
		scores["secrets_scan"] = ComponentScore{
			Score:   score,
			Reason:  "Secrets scan validation",
			Details: fmt.Sprintf("Valid: %v", report.Valid),
		}
	}

	return scores
}

func calculateOverallScore(componentScores map[string]ComponentScore) float64 {
	if len(componentScores) == 0 {
		return 0.0
	}

	// Define weights for each component
	weights := map[string]float64{
		"style_anchors":           0.30,
		"task_sizing":             0.25,
		"schema_compliance":       0.20,
		"secrets_scan":            0.15,
		"affirmative_constraints": 0.10,
	}

	totalWeight := 0.0
	weightedSum := 0.0

	for component, score := range componentScores {
		if weight, ok := weights[component]; ok {
			totalWeight += weight
			weightedSum += score.Score * weight
		}
	}

	if totalWeight == 0 {
		return 0.0
	}

	return weightedSum / totalWeight
}

func determineApprovalStatus(qualityScore float64, reports map[string]ValidationReport) (string, string) {
	minThreshold := 0.95

	if qualityScore >= minThreshold {
		return "APPROVED", fmt.Sprintf("Quality score %.2f ≥ minimum threshold %.2f", qualityScore, minThreshold)
	}

	// Check for critical failures
	for name, report := range reports {
		if report.OverallResult == "FAIL" || (report.Valid == false && isCriticalReport(name)) {
			return "REJECTED", fmt.Sprintf("Critical validation failure in %s", name)
		}
	}

	return "NEEDS_REVISION", fmt.Sprintf("Quality score %.2f < minimum threshold %.2f", qualityScore, minThreshold)
}

func isCriticalReport(reportName string) bool {
	criticalReports := map[string]bool{
		"schema_validation_report.json": true,
		"secrets_scan_report.json":      true,
		"yaml_validation_report.json":   true,
	}
	return criticalReports[reportName]
}

func generateIssues(reports map[string]ValidationReport) []Issue {
	var issues []Issue

	issueID := 1
	for name, report := range reports {
		if !report.Valid || report.OverallResult == "FAIL" {
			severity := "MEDIUM"
			if isCriticalReport(name) {
				severity = "HIGH"
			}

			issues = append(issues, Issue{
				ID:             fmt.Sprintf("VALIDATION-%03d", issueID),
				Severity:       severity,
				Category:       "validation",
				Description:    fmt.Sprintf("Validation failure in %s", name),
				Impact:         "Reduces overall quality score and may block approval",
				Recommendation: "Fix the validation issues identified in the report",
				Status:         "OPEN",
			})
			issueID++
		}
	}

	return issues
}

func generateRecommendations(reports map[string]ValidationReport) []Recommendation {
	var recommendations []Recommendation

	// Check for missing reports
	expectedReports := []string{
		"validation_report_sizing.json",
		"validation_report_enforcement.json",
	}

	for _, reportName := range expectedReports {
		if _, ok := reports[reportName]; !ok {
			recommendations = append(recommendations, Recommendation{
				Priority: "HIGH",
				Action:   fmt.Sprintf("Generate missing validation report: %s", reportName),
				Benefit:  "Complete validation coverage for comprehensive quality assessment",
			})
		}
	}

	// Add general recommendations
	recommendations = append(recommendations,
		Recommendation{
			Priority: "MEDIUM",
			Action:   "Review and address all validation issues",
			Benefit:  "Improve overall quality score and ensure approval",
		},
		Recommendation{
			Priority: "LOW",
			Action:   "Consider adding more project-specific style anchors",
			Benefit:  "Further reduce architectural drift risk",
		},
	)

	return recommendations
}

func generateValidationSummary(reports map[string]ValidationReport) map[string]string {
	summary := make(map[string]string)

	for name, report := range reports {
		status := "PASS"
		if !report.Valid || report.OverallResult == "FAIL" {
			status = "FAIL"
		} else if report.OverallResult == "" && report.Valid {
			status = "PASS"
		} else {
			status = report.OverallResult
		}

		// Map report names to summary keys
		switch name {
		case "validation_report_anchors.json":
			summary["style_anchors_compliance"] = status
		case "validation_report_constraints.json":
			summary["affirmative_constraints"] = status
		case "validation_report_implementation_guidelines.json":
			summary["implementation_guidelines"] = status
		case "validation_report_sizing.json":
			summary["task_sizing_compliance"] = status
		case "validation_report_enforcement.json":
			summary["multi_layer_enforcement"] = status
		case "yaml_validation_report.json":
			summary["yaml_syntax"] = status
		case "schema_validation_report.json":
			summary["schema_validation"] = status
		case "secrets_scan_report.json":
			summary["secrets_scan"] = status
		}
	}

	return summary
}
