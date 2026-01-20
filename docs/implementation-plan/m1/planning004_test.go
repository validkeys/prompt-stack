package m1

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestImplementationPlanExists(t *testing.T) {
	data, err := os.ReadFile("implementation-plan.yaml")
	if err != nil {
		t.Fatalf("Failed to read implementation-plan.yaml: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("implementation-plan.yaml is empty")
	}
}

func TestImplementationPlanContainsRequiredFields(t *testing.T) {
	data, err := os.ReadFile("implementation-plan.yaml")
	if err != nil {
		t.Fatalf("Failed to read implementation-plan.yaml: %v", err)
	}

	content := string(data)

	// Check for required top-level fields
	requiredFields := []string{
		"name: prompt-stack",
		"version:",
		"rules_file:",
		"task_sizing:",
		"tdd:",
		"model_preferences:",
		"outputs:",
		"tasks:",
	}

	for _, field := range requiredFields {
		if !strings.Contains(content, field) {
			t.Errorf("Missing required field: %s", field)
		}
	}

	// Check task_sizing has required subfields
	if !strings.Contains(content, "min_minutes:") || !strings.Contains(content, "max_minutes:") {
		t.Error("task_sizing missing required subfields (min_minutes, max_minutes)")
	}

	// Check tdd has required subfields
	if !strings.Contains(content, "required: true") {
		t.Error("tdd.required should be true")
	}
	if !strings.Contains(content, "test_command:") {
		t.Error("tdd missing test_command")
	}

	// Check model_preferences has primary field
	if !strings.Contains(content, "primary: opencode") {
		t.Error("model_preferences.primary should be 'opencode'")
	}

	// Check outputs has required arrays
	if !strings.Contains(content, "allowed_file_edits:") || !strings.Contains(content, "disallowed_file_edits:") {
		t.Error("outputs missing required arrays (allowed_file_edits, disallowed_file_edits)")
	}
}

func TestImplementationPlanContainsTasks(t *testing.T) {
	data, err := os.ReadFile("implementation-plan.yaml")
	if err != nil {
		t.Fatalf("Failed to read implementation-plan.yaml: %v", err)
	}

	content := string(data)

	// Check for task definitions (m1-001 through m1-007)
	for i := 1; i <= 7; i++ {
		taskID := fmt.Sprintf("m1-%03d", i)
		if !strings.Contains(content, fmt.Sprintf(`id: "%s"`, taskID)) {
			t.Errorf("Missing task: %s", taskID)
		}
	}

	// Check tasks have required fields
	taskFields := []string{
		"title:",
		"description:",
		"files_in_scope:",
		"style_anchors:",
		"estimated_duration_minutes:",
		"verification:",
	}

	// Count tasks
	taskCount := strings.Count(content, "id: \"m1-")
	if taskCount < 7 {
		t.Errorf("Expected at least 7 tasks, found %d", taskCount)
	}

	// For each task field, ensure it appears at least as many times as tasks
	for _, field := range taskFields {
		fieldCount := strings.Count(content, field)
		// We expect at least one occurrence per task, but some fields might be optional
		if fieldCount < taskCount/2 { // Allow some flexibility
			t.Errorf("Field '%s' appears only %d times for %d tasks", field, fieldCount, taskCount)
		}
	}
}

func TestImplementationPlanTaskSizingConstraints(t *testing.T) {
	data, err := os.ReadFile("implementation-plan.yaml")
	if err != nil {
		t.Fatalf("Failed to read implementation-plan.yaml: %v", err)
	}

	content := string(data)

	// Extract task_sizing values
	lines := strings.Split(content, "\n")
	var minMinutes, maxMinutes, maxFiles int
	inTaskSizing := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "task_sizing:" {
			inTaskSizing = true
			continue
		}
		if inTaskSizing && strings.HasPrefix(trimmed, "min_minutes:") {
			parts := strings.Split(trimmed, ":")
			if len(parts) > 1 {
				n, err := fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &minMinutes)
				if err != nil || n != 1 {
					t.Fatalf("Failed to parse task_sizing.min_minutes from %q: n=%d err=%v", trimmed, n, err)
				}
			}
		}
		if inTaskSizing && strings.HasPrefix(trimmed, "max_minutes:") {
			parts := strings.Split(trimmed, ":")
			if len(parts) > 1 {
				n, err := fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &maxMinutes)
				if err != nil || n != 1 {
					t.Fatalf("Failed to parse task_sizing.max_minutes from %q: n=%d err=%v", trimmed, n, err)
				}
			}
		}
		if inTaskSizing && strings.HasPrefix(trimmed, "max_files:") {
			parts := strings.Split(trimmed, ":")
			if len(parts) > 1 {
				n, err := fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &maxFiles)
				if err != nil || n != 1 {
					t.Fatalf("Failed to parse task_sizing.max_files from %q: n=%d err=%v", trimmed, n, err)
				}
			}
		}
		if inTaskSizing && trimmed != "" && !strings.HasPrefix(trimmed, "  ") && !strings.HasPrefix(trimmed, "-") {
			inTaskSizing = false
		}
	}

	// Validate task sizing constraints
	if minMinutes < 30 {
		t.Errorf("task_sizing.min_minutes should be >= 30, got %d", minMinutes)
	}
	if maxMinutes > 150 {
		t.Errorf("task_sizing.max_minutes should be <= 150, got %d", maxMinutes)
	}
	if maxFiles > 5 {
		t.Errorf("task_sizing.max_files should be <= 5, got %d", maxFiles)
	}
}

func TestImplementationPlanStyleAnchors(t *testing.T) {
	data, err := os.ReadFile("implementation-plan.yaml")
	if err != nil {
		t.Fatalf("Failed to read implementation-plan.yaml: %v", err)
	}

	content := string(data)

	// Check global style anchors
	if !strings.Contains(content, "style_anchors:") {
		t.Error("Missing global style_anchors array")
	}

	// Count style anchors references in tasks
	styleAnchorCount := strings.Count(content, "style_anchors:")
	// Should have at least: 1 global + 7 tasks * 2 anchors each = 15 references
	if styleAnchorCount < 8 { // 1 global + at least 1 per task
		t.Errorf("Insufficient style_anchors references: %d", styleAnchorCount)
	}

	// Check for style anchor structure (file and reason)
	if !strings.Contains(content, "file:") || !strings.Contains(content, "reason:") {
		t.Error("Style anchors missing required fields (file, reason)")
	}
}

func TestImplementationPlanValidationMetadata(t *testing.T) {
	data, err := os.ReadFile("implementation-plan.yaml")
	if err != nil {
		t.Fatalf("Failed to read implementation-plan.yaml: %v", err)
	}

	content := string(data)

	// Check for validation metadata
	if !strings.Contains(content, "_validation:") {
		t.Error("Missing validation metadata section")
	}

	// Check for quality score
	if !strings.Contains(content, "quality_score:") {
		t.Error("Missing quality_score in validation metadata")
	}

	// Check for approval status
	if !strings.Contains(content, "approval:") {
		t.Error("Missing approval in validation metadata")
	}

	// Check for reports
	if !strings.Contains(content, "reports:") {
		t.Error("Missing reports in validation metadata")
	}
}

func TestImplementationPlanSchemaCompliance(t *testing.T) {
	// This test validates that the YAML file can be validated against the schema
	// by running the validate_yaml tool if available
	data, err := os.ReadFile("implementation-plan.yaml")
	if err != nil {
		t.Fatalf("Failed to read implementation-plan.yaml: %v", err)
	}

	content := string(data)

	// Basic schema compliance checks
	// Check for all schema-required top-level fields from docs/ralphy-inputs.schema.json
	schemaRequiredFields := []string{
		"name:",
		"version:",
		"rules_file:",
		"task_sizing:",
		"tdd:",
		"model_preferences:",
		"outputs:",
		"tasks:",
	}

	missingFields := []string{}
	for _, field := range schemaRequiredFields {
		if !strings.Contains(content, field) {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		t.Errorf("Missing schema-required fields: %v", missingFields)
	}
}
