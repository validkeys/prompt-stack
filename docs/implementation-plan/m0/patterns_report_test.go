package m0

import (
	"encoding/json"
	"os"
	"testing"
)

// PatternsReport represents the structure of patterns_report.json
type PatternsReport struct {
	Milestone       string           `json:"milestone"`
	Title           string           `json:"title"`
	GeneratedAt     string           `json:"generated_at"`
	BasedOn         string           `json:"based_on"`
	Patterns        []Pattern        `json:"patterns_identified"`
	StyleAnchors    StyleAnchorsMap  `json:"style_anchors_by_relevance"`
	CommonPitfalls  []Pitfall        `json:"common_pitfalls_to_avoid"`
	Verification    Verification     `json:"verification_criteria"`
	Recommendations []Recommendation `json:"recommendations"`
	NextSteps       NextSteps        `json:"next_steps"`
}

type Pattern struct {
	ID          string   `json:"id"`
	Category    string   `json:"category"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Confidence  float64  `json:"confidence"`
	Source      Source   `json:"source"`
	Relevance   []string `json:"relevance_to_requirements"`
	Applicable  []string `json:"applicable_to_tasks"`
}

type Source struct {
	File           string `json:"file"`
	Lines          string `json:"lines"`
	SnippetPattern string `json:"snippet_pattern"`
}

type StyleAnchorsMap struct {
	HighRelevance       []StyleAnchor `json:"high_relevance"`
	MediumRelevance     []StyleAnchor `json:"medium_relevance"`
	ContextualRelevance []StyleAnchor `json:"contextual_relevance"`
}

type StyleAnchor struct {
	File            string   `json:"file"`
	Reason          string   `json:"reason"`
	ApplicableTo    []string `json:"applicable_to"`
	ConfidenceScore float64  `json:"confidence_score"`
}

type Pitfall struct {
	Pitfall          string `json:"pitfall"`
	Prevention       string `json:"prevention"`
	ReferencePattern string `json:"reference_pattern"`
}

type Verification struct {
	MinimumPatterns          int     `json:"minimum_patterns_identified"`
	ActualPatterns           int     `json:"actual_patterns_identified"`
	MinimumConfidence        float64 `json:"minimum_confidence_score"`
	ActualMinConfidence      float64 `json:"actual_min_confidence_score"`
	StyleAnchorsCategorized  bool    `json:"style_anchors_categorized"`
	ConfidenceAboveThreshold bool    `json:"confidence_scores_above_threshold"`
	MeetsPreCommitCriteria   bool    `json:"meets_pre_commit_criteria"`
}

type Recommendation struct {
	Priority       string `json:"priority"`
	Recommendation string `json:"recommendation"`
	Rationale      string `json:"rationale"`
}

type NextSteps struct {
	ReadyForBreakdown bool     `json:"ready_for_task_breakdown"`
	RecommendedTasks  []string `json:"recommended_tasks_for_next_phase"`
	KeyPatterns       []string `json:"key_patterns_to_apply"`
	KeyStyleAnchors   []string `json:"key_style_anchors_to_include"`
}

// TestPatternsReportExists verifies that patterns_report.json file exists
func TestPatternsReportExists(t *testing.T) {
	path := "patterns_report.json"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("patterns_report.json does not exist at %s", path)
	}
}

// TestPatternsReportValidJSON verifies that patterns_report.json is valid JSON
func TestPatternsReportValidJSON(t *testing.T) {
	data, err := os.ReadFile("patterns_report.json")
	if err != nil {
		t.Fatalf("Failed to read patterns_report.json: %v", err)
	}

	var report PatternsReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse patterns_report.json as JSON: %v", err)
	}
}

// TestPatternsReportStructure verifies required fields in patterns_report.json
func TestPatternsReportStructure(t *testing.T) {
	data, err := os.ReadFile("patterns_report.json")
	if err != nil {
		t.Fatalf("Failed to read patterns_report.json: %v", err)
	}

	var report PatternsReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse patterns_report.json as JSON: %v", err)
	}

	// Check required top-level fields
	if report.Milestone == "" {
		t.Error("Missing required field: milestone")
	}
	if report.Title == "" {
		t.Error("Missing required field: title")
	}
	if report.GeneratedAt == "" {
		t.Error("Missing required field: generated_at")
	}
	if report.BasedOn == "" {
		t.Error("Missing required field: based_on")
	}
}

// TestPatternsReportMinimumPatterns verifies at least 5 patterns are identified
func TestPatternsReportMinimumPatterns(t *testing.T) {
	data, err := os.ReadFile("patterns_report.json")
	if err != nil {
		t.Fatalf("Failed to read patterns_report.json: %v", err)
	}

	var report PatternsReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse patterns_report.json as JSON: %v", err)
	}

	minPatterns := 5
	if len(report.Patterns) < minPatterns {
		t.Errorf("Expected at least %d patterns, got %d", minPatterns, len(report.Patterns))
	}

	// Verify internal count matches
	if report.Verification.MinimumPatterns != minPatterns {
		t.Errorf("Verification minimum_patterns_identified should be %d, got %d", minPatterns, report.Verification.MinimumPatterns)
	}
	if report.Verification.ActualPatterns != len(report.Patterns) {
		t.Errorf("Verification actual_patterns_identified should match patterns count, got %d vs %d", report.Verification.ActualPatterns, len(report.Patterns))
	}
}

// TestPatternsReportPatternFields verifies each pattern has required fields
func TestPatternsReportPatternFields(t *testing.T) {
	data, err := os.ReadFile("patterns_report.json")
	if err != nil {
		t.Fatalf("Failed to read patterns_report.json: %v", err)
	}

	var report PatternsReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse patterns_report.json as JSON: %v", err)
	}

	for i, pattern := range report.Patterns {
		if pattern.ID == "" {
			t.Errorf("Pattern[%d]: missing required field id", i)
		}
		if pattern.Category == "" {
			t.Errorf("Pattern[%d]: missing required field category", i)
		}
		if pattern.Name == "" {
			t.Errorf("Pattern[%d]: missing required field name", i)
		}
		if pattern.Description == "" {
			t.Errorf("Pattern[%d]: missing required field description", i)
		}
		if pattern.Confidence < 0 || pattern.Confidence > 1 {
			t.Errorf("Pattern[%d]: confidence must be between 0 and 1, got %f", i, pattern.Confidence)
		}
		if pattern.Source.File == "" {
			t.Errorf("Pattern[%d]: missing required field source.file", i)
		}
		if len(pattern.Relevance) == 0 {
			t.Errorf("Pattern[%d]: relevance_to_requirements must not be empty", i)
		}
		if len(pattern.Applicable) == 0 {
			t.Errorf("Pattern[%d]: applicable_to_tasks must not be empty", i)
		}
	}
}

// TestPatternsReportStyleAnchorsCategorized verifies style anchors are properly categorized
func TestPatternsReportStyleAnchorsCategorized(t *testing.T) {
	data, err := os.ReadFile("patterns_report.json")
	if err != nil {
		t.Fatalf("Failed to read patterns_report.json: %v", err)
	}

	var report PatternsReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse patterns_report.json as JSON: %v", err)
	}

	// Check that style anchors are categorized
	if len(report.StyleAnchors.HighRelevance) == 0 {
		t.Error("Expected at least one high-relevance style anchor")
	}
	if len(report.StyleAnchors.MediumRelevance) == 0 {
		t.Error("Expected at least one medium-relevance style anchor")
	}
	if len(report.StyleAnchors.ContextualRelevance) == 0 {
		t.Error("Expected at least one contextual-relevance style anchor")
	}

	// Verify each style anchor has required fields
	checkStyleAnchorFields := func(anchors []StyleAnchor, category string) {
		for i, anchor := range anchors {
			if anchor.File == "" {
				t.Errorf("Style anchor %s[%d]: missing required field file", category, i)
			}
			if anchor.Reason == "" {
				t.Errorf("Style anchor %s[%d]: missing required field reason", category, i)
			}
			if anchor.ConfidenceScore < 0 || anchor.ConfidenceScore > 1 {
				t.Errorf("Style anchor %s[%d]: confidence_score must be between 0 and 1, got %f", category, i, anchor.ConfidenceScore)
			}
		}
	}

	checkStyleAnchorFields(report.StyleAnchors.HighRelevance, "high_relevance")
	checkStyleAnchorFields(report.StyleAnchors.MediumRelevance, "medium_relevance")
	checkStyleAnchorFields(report.StyleAnchors.ContextualRelevance, "contextual_relevance")
}

// TestPatternsReportConfidenceThreshold verifies all patterns meet minimum confidence score
func TestPatternsReportConfidenceThreshold(t *testing.T) {
	data, err := os.ReadFile("patterns_report.json")
	if err != nil {
		t.Fatalf("Failed to read patterns_report.json: %v", err)
	}

	var report PatternsReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse patterns_report.json as JSON: %v", err)
	}

	minConfidence := 0.7
	for i, pattern := range report.Patterns {
		if pattern.Confidence < minConfidence {
			t.Errorf("Pattern[%d] (%s): confidence %f is below threshold %f", i, pattern.ID, pattern.Confidence, minConfidence)
		}
	}

	// Verify verification criteria
	if report.Verification.MinimumConfidence != minConfidence {
		t.Errorf("Verification minimum_confidence_score should be %f, got %f", minConfidence, report.Verification.MinimumConfidence)
	}
}

// TestPatternsReportVerificationCriteria verifies verification criteria are met
func TestPatternsReportVerificationCriteria(t *testing.T) {
	data, err := os.ReadFile("patterns_report.json")
	if err != nil {
		t.Fatalf("Failed to read patterns_report.json: %v", err)
	}

	var report PatternsReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse patterns_report.json as JSON: %v", err)
	}

	if !report.Verification.StyleAnchorsCategorized {
		t.Error("Verification criteria: style_anchors_categorized should be true")
	}
	if !report.Verification.ConfidenceAboveThreshold {
		t.Error("Verification criteria: confidence_scores_above_threshold should be true")
	}
	if !report.Verification.MeetsPreCommitCriteria {
		t.Error("Verification criteria: meets_pre_commit_criteria should be true")
	}
}

// TestPatternsReportCommonPitfalls verifies common pitfalls are documented
func TestPatternsReportCommonPitfalls(t *testing.T) {
	data, err := os.ReadFile("patterns_report.json")
	if err != nil {
		t.Fatalf("Failed to read patterns_report.json: %v", err)
	}

	var report PatternsReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse patterns_report.json as JSON: %v", err)
	}

	// Should have at least 3 documented pitfalls
	if len(report.CommonPitfalls) < 3 {
		t.Errorf("Expected at least 3 documented pitfalls, got %d", len(report.CommonPitfalls))
	}

	// Verify each pitfall has required fields
	for i, pitfall := range report.CommonPitfalls {
		if pitfall.Pitfall == "" {
			t.Errorf("Pitfall[%d]: missing required field pitfall", i)
		}
		if pitfall.Prevention == "" {
			t.Errorf("Pitfall[%d]: missing required field prevention", i)
		}
		if pitfall.ReferencePattern == "" {
			t.Errorf("Pitfall[%d]: missing required field reference_pattern", i)
		}
	}
}

// TestPatternsReportBasedOnStructuredAnalysis verifies it's based on structured_analysis.json
func TestPatternsReportBasedOnStructuredAnalysis(t *testing.T) {
	data, err := os.ReadFile("patterns_report.json")
	if err != nil {
		t.Fatalf("Failed to read patterns_report.json: %v", err)
	}

	var report PatternsReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse patterns_report.json as JSON: %v", err)
	}

	expectedBasedOn := "structured_analysis.json"
	if report.BasedOn != expectedBasedOn {
		t.Errorf("Expected based_on to be %s, got %s", expectedBasedOn, report.BasedOn)
	}
}
