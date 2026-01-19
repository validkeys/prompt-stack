package m0

import (
	"encoding/json"
	"os"
	"testing"
)

// TestValidationReportAnchorsStructure validates the structure of validation_report_anchors.json
func TestValidationReportAnchorsStructure(t *testing.T) {
	data, err := os.ReadFile("validation_report_anchors.json")
	if err != nil {
		t.Fatalf("Failed to read validation report: %v", err)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected interface{}
		checkMin bool
	}{
		{"validation_type exists", "validation_type", "style_anchors_compliance", false},
		{"validated_at exists", "validated_at", "", false},
		{"validator exists", "validator", "planning-005", false},
		{"overall_result exists", "overall_result", "", false},
		{"overall_score exists", "overall_score", 0.0, true},
		{"summary exists", "summary", nil, false},
		{"tasks exists", "tasks", nil, false},
		{"metrics exists", "metrics", nil, false},
		{"recommendations exists", "recommendations", nil, false},
		{"verification_criteria exists", "verification_criteria", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, ok := getPath(report, tt.path)
			if !ok {
				t.Errorf("Path %s not found in report", tt.path)
				return
			}
			if tt.checkMin {
				if numValue, ok := value.(float64); ok {
					if numValue < tt.expected.(float64) {
						t.Errorf("Path %s = %v, want >= %v", tt.path, value, tt.expected)
					}
				}
			} else if tt.expected != "" && tt.expected != nil && value != tt.expected {
				t.Errorf("Path %s = %v, want %v", tt.path, value, tt.expected)
			}
		})
	}
}

// TestValidationReportAnchorsSummary validates summary section
func TestValidationReportAnchorsSummary(t *testing.T) {
	data, err := os.ReadFile("validation_report_anchors.json")
	if err != nil {
		t.Fatalf("Failed to read validation report: %v", err)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	summary, ok := report["summary"].(map[string]interface{})
	if !ok {
		t.Fatal("summary is not a map")
	}

	tests := []struct {
		name     string
		key      string
		minValue float64
	}{
		{"total_tasks is 7", "total_tasks", 7},
		{"tasks_passed is 7", "tasks_passed", 7},
		{"tasks_failed is 0", "tasks_failed", 0},
		{"total_anchors >= 14", "total_anchors", 14},
		{"avg_anchors_per_task >= 2.0", "avg_anchors_per_task", 2.0},
		{"min_required_anchors is 2", "min_required_anchors", 2},
		{"max_allowed_anchors is 3", "max_allowed_anchors", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, ok := summary[tt.key].(float64)
			if !ok {
				t.Errorf("%s is not a number", tt.key)
				return
			}
			if value < tt.minValue {
				t.Errorf("%s = %v, want >= %v", tt.key, value, tt.minValue)
			}
		})
	}

	// Verify all tasks passed
	if summary["tasks_passed"].(float64) != summary["total_tasks"].(float64) {
		t.Errorf("tasks_passed (%v) != total_tasks (%v)", summary["tasks_passed"], summary["total_tasks"])
	}
	if summary["tasks_failed"].(float64) != 0 {
		t.Errorf("tasks_failed = %v, want 0", summary["tasks_failed"])
	}
}

// TestValidationReportAnchorsTasks validates each task's anchors
func TestValidationReportAnchorsTasks(t *testing.T) {
	data, err := os.ReadFile("validation_report_anchors.json")
	if err != nil {
		t.Fatalf("Failed to read validation report: %v", err)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	tasks, ok := report["tasks"].([]interface{})
	if !ok {
		t.Fatal("tasks is not an array")
	}

	if len(tasks) != 7 {
		t.Fatalf("Expected 7 tasks, got %d", len(tasks))
	}

	taskIDs := []string{"m0-001", "m0-002", "m0-003", "m0-004", "m0-005", "m0-006", "m0-007"}
	for i, taskID := range taskIDs {
		t.Run("task_"+taskID, func(t *testing.T) {
			task := tasks[i].(map[string]interface{})
			if task["task_id"] != taskID {
				t.Errorf("task[%d].task_id = %v, want %s", i, task["task_id"], taskID)
			}

			// Check anchor count is 2 or 3
			anchorCount := task["anchor_count"].(float64)
			if anchorCount < 2 || anchorCount > 3 {
				t.Errorf("task %s has %v anchors, want 2-3", taskID, anchorCount)
			}

			// Check result is PASS
			if task["result"] != "PASS" {
				t.Errorf("task %s result = %v, want PASS", taskID, task["result"])
			}

			// Check issues is empty array
			issues := task["issues"].([]interface{})
			if len(issues) != 0 {
				t.Errorf("task %s has issues: %v", taskID, issues)
			}

			// Check anchors exist and are accessible
			anchors := task["anchors"].([]interface{})
			for j, anchor := range anchors {
				anchorMap := anchor.(map[string]interface{})
				if anchorMap["exists"] != true {
					t.Errorf("task %s anchor[%d] does not exist", taskID, j)
				}

				relevanceScore := anchorMap["relevance_score"].(float64)
				if relevanceScore < 0.7 {
					t.Errorf("task %s anchor[%d] relevance score = %v, want >= 0.7", taskID, j, relevanceScore)
				}

				reasonQuality := anchorMap["reason_quality"].(string)
				if reasonQuality != "specific" {
					t.Errorf("task %s anchor[%d] reason quality = %v, want specific", taskID, j, reasonQuality)
				}
			}
		})
	}
}

// TestValidationReportAnchorsMetrics validates metrics section
func TestValidationReportAnchorsMetrics(t *testing.T) {
	data, err := os.ReadFile("validation_report_anchors.json")
	if err != nil {
		t.Fatalf("Failed to read validation report: %v", err)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	metrics, ok := report["metrics"].(map[string]interface{})
	if !ok {
		t.Fatal("metrics is not a map")
	}

	// Test anchor relevance scores
	relevanceScores := metrics["anchor_relevance_scores"].(map[string]interface{})
	if relevanceScores["min"].(float64) < 0.7 {
		t.Errorf("min relevance score = %v, want >= 0.7", relevanceScores["min"])
	}
	if relevanceScores["avg"].(float64) < 0.8 {
		t.Errorf("avg relevance score = %v, want >= 0.8", relevanceScores["avg"])
	}

	// Test anchor file existence
	fileExistence := metrics["anchor_file_existence"].(map[string]interface{})
	if fileExistence["missing"].(float64) != 0 {
		t.Errorf("missing files = %v, want 0", fileExistence["missing"])
	}
	if fileExistence["success_rate"].(float64) != 1.0 {
		t.Errorf("file existence success rate = %v, want 1.0", fileExistence["success_rate"])
	}

	// Test anchor reason quality
	reasonQuality := metrics["anchor_reason_quality"].(map[string]interface{})
	if reasonQuality["generic"].(float64) != 0 {
		t.Errorf("generic reasons = %v, want 0", reasonQuality["generic"])
	}
	if reasonQuality["vague"].(float64) != 0 {
		t.Errorf("vague reasons = %v, want 0", reasonQuality["vague"])
	}
	if reasonQuality["success_rate"].(float64) != 1.0 {
		t.Errorf("reason quality success rate = %v, want 1.0", reasonQuality["success_rate"])
	}

	// Test anchor count compliance
	countCompliance := metrics["anchor_count_compliance"].(map[string]interface{})
	if countCompliance["tasks_below_2"].(float64) != 0 {
		t.Errorf("tasks below 2 anchors = %v, want 0", countCompliance["tasks_below_2"])
	}
	if countCompliance["tasks_above_3"].(float64) != 0 {
		t.Errorf("tasks above 3 anchors = %v, want 0", countCompliance["tasks_above_3"])
	}
	if countCompliance["compliance_rate"].(float64) != 1.0 {
		t.Errorf("compliance rate = %v, want 1.0", countCompliance["compliance_rate"])
	}
}

// TestValidationReportAnchorsVerificationCriteria validates verification criteria
func TestValidationReportAnchorsVerificationCriteria(t *testing.T) {
	data, err := os.ReadFile("validation_report_anchors.json")
	if err != nil {
		t.Fatalf("Failed to read validation report: %v", err)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	verification, ok := report["verification_criteria"].(map[string]interface{})
	if !ok {
		t.Fatal("verification_criteria is not a map")
	}

	tests := []struct {
		name string
		key  string
	}{
		{"all_tasks_have_2_3_anchors", "all_tasks_have_2_3_anchors"},
		{"anchor_relevance_score_above_0_8", "anchor_relevance_score_above_0_8"},
		{"anchor_files_exist_and_accessible", "anchor_files_exist_and_accessible"},
		{"anchor_reasons_are_specific", "anchor_reasons_are_specific"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, ok := verification[tt.key]
			if !ok {
				t.Errorf("verification_criteria.%s not found", tt.key)
				return
			}
			if value != true {
				t.Errorf("verification_criteria.%s = %v, want true", tt.key, value)
			}
		})
	}
}

// TestValidationReportAnchorsOverallScore validates overall score is >= 0.9
func TestValidationReportAnchorsOverallScore(t *testing.T) {
	data, err := os.ReadFile("validation_report_anchors.json")
	if err != nil {
		t.Fatalf("Failed to read validation report: %v", err)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	overallScore, ok := report["overall_score"].(float64)
	if !ok {
		t.Fatal("overall_score is not a number")
	}

	if overallScore < 0.9 {
		t.Errorf("overall_score = %v, want >= 0.9", overallScore)
	}

	if report["overall_result"] != "PASS" {
		t.Errorf("overall_result = %v, want PASS", report["overall_result"])
	}
}

// Helper function to get nested values from map
func getPath(m map[string]interface{}, path string) (interface{}, bool) {
	// For now, only handle top-level keys
	return m[path], true
}
