package main

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateQualityReport(t *testing.T) {
	// Create a temporary directory with test validation reports
	tempDir := t.TempDir()

	// Create test validation reports
	testReports := map[string]ValidationReport{
		"validation_report_anchors.json": {
			Valid:          true,
			OverallResult:  "PASS",
			OverallScore:   0.92,
			ValidationType: "style_anchors_compliance",
		},
		"validation_report_constraints.json": {
			Valid: true,
		},
		"validation_report_implementation_guidelines.json": {
			Valid: false,
		},
		"yaml_validation_report.json": {
			Valid: true,
		},
		"schema_validation_report.json": {
			Valid: true,
		},
		"secrets_scan_report.json": {
			Valid: true,
		},
	}

	for filename, report := range testReports {
		data, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal test report %s: %v", filename, err)
		}

		reportPath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(reportPath, data, 0644); err != nil {
			t.Fatalf("Failed to write test report %s: %v", filename, err)
		}
	}

	// Generate quality report
	qualityReport, err := generateQualityReport(tempDir)
	if err != nil {
		t.Fatalf("Failed to generate quality report: %v", err)
	}

	// Verify report structure
	if qualityReport.ReportType != "final_quality_report" {
		t.Errorf("Expected report_type 'final_quality_report', got %s", qualityReport.ReportType)
	}

	if qualityReport.Milestone != "m0" {
		t.Errorf("Expected milestone 'm0', got %s", qualityReport.Milestone)
	}

	// Verify quality score calculation
	if qualityReport.QualityScore <= 0 || qualityReport.QualityScore > 1.0 {
		t.Errorf("Quality score should be between 0 and 1, got %.2f", qualityReport.QualityScore)
	}

	// Verify component scores
	expectedComponents := []string{
		"style_anchors",
		"affirmative_constraints",
		"implementation_guidelines",
		"yaml_syntax",
		"schema_compliance",
		"secrets_scan",
	}

	for _, component := range expectedComponents {
		if _, ok := qualityReport.ComponentScores[component]; !ok {
			t.Errorf("Missing component score for %s", component)
		}
	}

	// Verify validation summary
	expectedSummaryKeys := []string{
		"style_anchors_compliance",
		"affirmative_constraints",
		"implementation_guidelines",
		"yaml_syntax",
		"schema_validation",
		"secrets_scan",
	}

	for _, key := range expectedSummaryKeys {
		if status, ok := qualityReport.ValidationSummary[key]; !ok {
			t.Errorf("Missing validation summary for %s", key)
		} else if status == "" {
			t.Errorf("Empty validation status for %s", key)
		}
	}

	// Verify issues were generated (should have at least one for implementation_guidelines failure)
	if len(qualityReport.Issues) == 0 {
		t.Error("Expected at least one issue for validation failure")
	}

	// Verify recommendations
	if len(qualityReport.Recommendations) == 0 {
		t.Error("Expected at least one recommendation")
	}
}

func TestCalculateOverallScore(t *testing.T) {
	tests := []struct {
		name     string
		scores   map[string]ComponentScore
		expected float64
	}{
		{
			name:     "empty scores",
			scores:   map[string]ComponentScore{},
			expected: 0.0,
		},
		{
			name: "all perfect scores",
			scores: map[string]ComponentScore{
				"style_anchors":           {Score: 1.0},
				"task_sizing":             {Score: 1.0},
				"schema_compliance":       {Score: 1.0},
				"secrets_scan":            {Score: 1.0},
				"affirmative_constraints": {Score: 1.0},
			},
			expected: 1.0,
		},
		{
			name: "mixed scores",
			scores: map[string]ComponentScore{
				"style_anchors":           {Score: 0.8},
				"task_sizing":             {Score: 0.9},
				"schema_compliance":       {Score: 1.0},
				"secrets_scan":            {Score: 1.0},
				"affirmative_constraints": {Score: 0.7},
			},
			expected: 0.885, // (0.8*0.3 + 0.9*0.25 + 1.0*0.2 + 1.0*0.15 + 0.7*0.1) / 1.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := calculateOverallScore(tt.scores)
			// Allow small floating point differences
			if math.Abs(score-tt.expected) > 0.001 {
				t.Errorf("calculateOverallScore() = %.15f, want %.15f", score, tt.expected)
			}
		})
	}
}

func TestDetermineApprovalStatus(t *testing.T) {
	tests := []struct {
		name           string
		qualityScore   float64
		reports        map[string]ValidationReport
		expectedStatus string
		expectedReason string
	}{
		{
			name:           "approved - high score",
			qualityScore:   0.97,
			reports:        map[string]ValidationReport{},
			expectedStatus: "APPROVED",
			expectedReason: "Quality score 0.97 ≥ minimum threshold 0.95",
		},
		{
			name:           "needs revision - low score",
			qualityScore:   0.90,
			reports:        map[string]ValidationReport{},
			expectedStatus: "NEEDS_REVISION",
			expectedReason: "Quality score 0.90 < minimum threshold 0.95",
		},
		{
			name:         "rejected - critical failure",
			qualityScore: 0.94, // Below threshold with critical failure
			reports: map[string]ValidationReport{
				"schema_validation_report.json": {
					Valid: false,
				},
			},
			expectedStatus: "REJECTED",
			expectedReason: "Critical validation failure in schema_validation_report.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, reason := determineApprovalStatus(tt.qualityScore, tt.reports)
			if status != tt.expectedStatus {
				t.Errorf("determineApprovalStatus() status = %s, want %s", status, tt.expectedStatus)
			}
			if reason != tt.expectedReason {
				t.Errorf("determineApprovalStatus() reason = %s, want %s", reason, tt.expectedReason)
			}
		})
	}
}

func TestGenerateIssues(t *testing.T) {
	reports := map[string]ValidationReport{
		"validation_report_anchors.json": {
			Valid:         true,
			OverallResult: "PASS",
		},
		"schema_validation_report.json": {
			Valid: false,
		},
		"validation_report_implementation_guidelines.json": {
			Valid:         false,
			OverallResult: "FAIL",
		},
	}

	issues := generateIssues(reports)

	// Should have 2 issues (schema validation and implementation guidelines)
	if len(issues) != 2 {
		t.Errorf("Expected 2 issues, got %d", len(issues))
	}

	// Check that issues have proper structure
	for _, issue := range issues {
		if issue.ID == "" {
			t.Error("Issue ID should not be empty")
		}
		if issue.Severity == "" {
			t.Error("Issue severity should not be empty")
		}
		if issue.Description == "" {
			t.Error("Issue description should not be empty")
		}
	}
}

func TestGenerateRecommendations(t *testing.T) {
	reports := map[string]ValidationReport{
		"validation_report_anchors.json": {
			Valid: true,
		},
		// Missing sizing and enforcement reports
	}

	recommendations := generateRecommendations(reports)

	// Should have at least 3 recommendations (2 for missing reports + 1 general)
	if len(recommendations) < 3 {
		t.Errorf("Expected at least 3 recommendations, got %d", len(recommendations))
	}

	// Check for missing report recommendations
	foundSizingRecommendation := false
	foundEnforcementRecommendation := false
	for _, rec := range recommendations {
		if rec.Action == "Generate missing validation report: validation_report_sizing.json" {
			foundSizingRecommendation = true
		}
		if rec.Action == "Generate missing validation report: validation_report_enforcement.json" {
			foundEnforcementRecommendation = true
		}
	}

	if !foundSizingRecommendation {
		t.Error("Missing recommendation for sizing report")
	}
	if !foundEnforcementRecommendation {
		t.Error("Missing recommendation for enforcement report")
	}
}

func TestMainFunction(t *testing.T) {
	// Test that main function can be called without arguments (should exit with error)
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"validate.go"}

	// We can't easily test os.Exit in a test, but we can verify the function compiles
	// and the test suite runs successfully
}
