package m1

import (
	"encoding/json"
	"os"
	"testing"
)

type StructuredAnalysis struct {
	AnalysisType               string                   `json:"analysis_type"`
	Timestamp                  string                   `json:"timestamp"`
	RequirementsFile           string                   `json:"requirements_file"`
	MilestoneID                string                   `json:"milestone_id"`
	MilestoneTitle             string                   `json:"milestone_title"`
	CoreFunctionalRequirements []map[string]interface{} `json:"core_functional_requirements"`
	NonFunctionalRequirements  []map[string]interface{} `json:"non_functional_requirements"`
	IntegrationPoints          []map[string]interface{} `json:"integration_points"`
	Dependencies               []map[string]interface{} `json:"dependencies"`
	Constraints                []map[string]interface{} `json:"constraints"`
	AcceptanceCriteria         []map[string]interface{} `json:"acceptance_criteria"`
	Deliverables               []map[string]interface{} `json:"deliverables"`
	StyleAnchors               []map[string]interface{} `json:"style_anchors"`
	QualityTargets             []map[string]interface{} `json:"quality_targets"`
	AnalysisNotes              []string                 `json:"analysis_notes"`
	AmbiguitiesFlagged         []map[string]interface{} `json:"ambiguities_flagged"`
}

func TestStructuredAnalysisExists(t *testing.T) {
	data, err := os.ReadFile("structured_analysis.json")
	if err != nil {
		t.Fatalf("Failed to read structured_analysis.json: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("structured_analysis.json is empty")
	}
}

func TestStructuredAnalysisValidJSON(t *testing.T) {
	data, err := os.ReadFile("structured_analysis.json")
	if err != nil {
		t.Fatalf("Failed to read structured_analysis.json: %v", err)
	}

	var analysis StructuredAnalysis
	if err := json.Unmarshal(data, &analysis); err != nil {
		t.Fatalf("Failed to parse structured_analysis.json as JSON: %v", err)
	}
}

func TestStructuredAnalysisRequiredFields(t *testing.T) {
	data, err := os.ReadFile("structured_analysis.json")
	if err != nil {
		t.Fatalf("Failed to read structured_analysis.json: %v", err)
	}

	var analysis StructuredAnalysis
	if err := json.Unmarshal(data, &analysis); err != nil {
		t.Fatalf("Failed to parse structured_analysis.json as JSON: %v", err)
	}

	requiredFields := map[string]string{
		"AnalysisType":     analysis.AnalysisType,
		"Timestamp":        analysis.Timestamp,
		"RequirementsFile": analysis.RequirementsFile,
		"MilestoneID":      analysis.MilestoneID,
		"MilestoneTitle":   analysis.MilestoneTitle,
	}

	for fieldName, fieldValue := range requiredFields {
		if fieldValue == "" {
			t.Errorf("Required field %s is empty", fieldName)
		}
	}

	if len(analysis.CoreFunctionalRequirements) == 0 {
		t.Error("core_functional_requirements array is empty")
	}

	if len(analysis.NonFunctionalRequirements) == 0 {
		t.Error("non_functional_requirements array is empty")
	}

	if len(analysis.AcceptanceCriteria) == 0 {
		t.Error("acceptance_criteria array is empty")
	}

	if len(analysis.Deliverables) == 0 {
		t.Error("deliverables array is empty")
	}

	if len(analysis.StyleAnchors) == 0 {
		t.Error("style_anchors array is empty")
	}
}

func TestStructuredAnalysisContent(t *testing.T) {
	data, err := os.ReadFile("structured_analysis.json")
	if err != nil {
		t.Fatalf("Failed to read structured_analysis.json: %v", err)
	}

	var analysis StructuredAnalysis
	if err := json.Unmarshal(data, &analysis); err != nil {
		t.Fatalf("Failed to parse structured_analysis.json as JSON: %v", err)
	}

	// Check milestone ID matches m1
	if analysis.MilestoneID != "m1" {
		t.Errorf("Expected milestone_id to be 'm1', got %s", analysis.MilestoneID)
	}

	// Check requirements file path is correct
	expectedPath := "docs/implementation-plan/m1/requirements.md"
	if analysis.RequirementsFile != expectedPath {
		t.Errorf("Expected requirements_file to be %s, got %s", expectedPath, analysis.RequirementsFile)
	}

	// Check core functional requirements have expected structure
	for i, fr := range analysis.CoreFunctionalRequirements {
		if id, ok := fr["id"].(string); !ok || id == "" {
			t.Errorf("Core functional requirement %d missing 'id' field", i)
		}
		if title, ok := fr["title"].(string); !ok || title == "" {
			t.Errorf("Core functional requirement %d missing 'title' field", i)
		}
	}

	// Check acceptance criteria have IDs
	for i, ac := range analysis.AcceptanceCriteria {
		if id, ok := ac["id"].(string); !ok || id == "" {
			t.Errorf("Acceptance criterion %d missing 'id' field", i)
		}
	}

	// Check deliverables have names
	for i, deliverable := range analysis.Deliverables {
		if name, ok := deliverable["name"].(string); !ok || name == "" {
			t.Errorf("Deliverable %d missing 'name' field", i)
		}
	}
}

func TestStructuredAnalysisAmbiguities(t *testing.T) {
	data, err := os.ReadFile("structured_analysis.json")
	if err != nil {
		t.Fatalf("Failed to read structured_analysis.json: %v", err)
	}

	var analysis StructuredAnalysis
	if err := json.Unmarshal(data, &analysis); err != nil {
		t.Fatalf("Failed to parse structured_analysis.json as JSON: %v", err)
	}

	// Check that ambiguities are flagged if present
	// This is optional but good practice
	if len(analysis.AmbiguitiesFlagged) > 0 {
		t.Logf("Found %d ambiguities flagged in analysis", len(analysis.AmbiguitiesFlagged))
		for i, ambiguity := range analysis.AmbiguitiesFlagged {
			if issue, ok := ambiguity["issue"].(string); !ok || issue == "" {
				t.Errorf("Ambiguity %d missing 'issue' field", i)
			}
		}
	}
}
