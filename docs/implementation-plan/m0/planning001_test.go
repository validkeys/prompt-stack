package m0

import (
	"encoding/json"
	"os"
	"testing"
)

type StructuredAnalysis struct {
	Milestone                  string   `json:"milestone"`
	Title                      string   `json:"title"`
	RequirementsFile           string   `json:"requirements_file"`
	AnalyzedAt                 string   `json:"analyzed_at"`
	CoreFunctionalRequirements []any    `json:"core_functional_requirements"`
	NonFunctionalRequirements  []any    `json:"non_functional_requirements"`
	IntegrationPoints          []any    `json:"integration_points"`
	Dependencies               []any    `json:"dependencies"`
	Constraints                []any    `json:"constraints"`
	AcceptanceCriteria         []any    `json:"acceptance_criteria"`
	Deliverables               []string `json:"deliverables"`
	TechStack                  any      `json:"tech_stack"`
	Stakeholders               any      `json:"stakeholders"`
	Timeline                   any      `json:"timeline"`
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
		"Milestone":        analysis.Milestone,
		"Title":            analysis.Title,
		"RequirementsFile": analysis.RequirementsFile,
		"AnalyzedAt":       analysis.AnalyzedAt,
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

	if len(analysis.Deliverables) == 0 {
		t.Error("deliverables array is empty")
	}
}

func TestStructuredAnalysisAcceptanceCriteria(t *testing.T) {
	data, err := os.ReadFile("structured_analysis.json")
	if err != nil {
		t.Fatalf("Failed to read structured_analysis.json: %v", err)
	}

	var analysis StructuredAnalysis
	if err := json.Unmarshal(data, &analysis); err != nil {
		t.Fatalf("Failed to parse structured_analysis.json as JSON: %v", err)
	}

	if len(analysis.AcceptanceCriteria) == 0 {
		t.Error("acceptance_criteria array is empty")
	}
}

func TestStructuredAnalysisConstraints(t *testing.T) {
	data, err := os.ReadFile("structured_analysis.json")
	if err != nil {
		t.Fatalf("Failed to read structured_analysis.json: %v", err)
	}

	var analysis StructuredAnalysis
	if err := json.Unmarshal(data, &analysis); err != nil {
		t.Fatalf("Failed to parse structured_analysis.json as JSON: %v", err)
	}

	if len(analysis.Constraints) == 0 {
		t.Error("constraints array is empty")
	}
}
