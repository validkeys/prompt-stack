package m1

import (
	"os"
	"testing"
)

// TestTaskBreakdownExists validates the task breakdown file exists
func TestTaskBreakdownExists(t *testing.T) {
	if _, err := os.Stat("task_breakdown.yaml"); os.IsNotExist(err) {
		t.Fatal("task_breakdown.yaml does not exist")
	}
}

// TestTaskBreakdownStructure validates basic structure of task breakdown
func TestTaskBreakdownStructure(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	// Check for required top-level fields
	requiredFields := []string{
		"milestone:",
		"milestone_title:",
		"generated_at:",
		"based_on:",
		"tasks:",
		"dependency_graph:",
		"parallel_execution_opportunities:",
		"validation_summary:",
	}

	for _, field := range requiredFields {
		if !contains(content, field) {
			t.Errorf("Missing required field in task_breakdown.yaml: %s", field)
		}
	}

	// Check milestone is m1
	if !contains(content, "milestone: m1") {
		t.Error("Milestone should be 'm1'")
	}
}

// TestTaskCount validates we have the expected number of tasks
func TestTaskCount(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	// Count tasks by looking for task IDs
	taskCount := 0
	expectedIDs := []string{"m1-001", "m1-002", "m1-003", "m1-004", "m1-005", "m1-006", "m1-007"}
	for _, id := range expectedIDs {
		if contains(content, "- id: \""+id+"\"") {
			taskCount++
		}
	}

	if taskCount != 7 {
		t.Errorf("Expected 7 tasks, found %d", taskCount)
	}
}

// TestTaskSizing validates task sizing compliance
func TestTaskSizing(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	// Check validation summary indicates compliance
	if !contains(content, "task_sizing_compliance: true") {
		t.Error("Task sizing compliance should be true")
	}

	// Check task durations are within policy (30-150 minutes)
	taskDurations := map[string]int{
		"m1-001": 45,
		"m1-002": 90,
		"m1-003": 60,
		"m1-004": 75,
		"m1-005": 60,
		"m1-006": 90,
		"m1-007": 45,
	}

	for taskID, expectedDuration := range taskDurations {
		durationPattern := "- id: \"" + taskID + "\"\n"
		durationIdx := index(content, durationPattern)
		if durationIdx == -1 {
			t.Errorf("Task %s not found", taskID)
			continue
		}

		// Look for estimated_duration_minutes after the task ID
		remaining := content[durationIdx:]
		durationLine := "    estimated_duration_minutes: " + string(rune('0'+expectedDuration/10)) + string(rune('0'+expectedDuration%10))
		if !contains(remaining, durationLine) {
			t.Errorf("Task %s should have estimated_duration_minutes: %d", taskID, expectedDuration)
		}
	}
}

// TestStyleAnchors validates each task has 2-3 style anchors
func TestStyleAnchors(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	// Check validation summary indicates compliance
	if !contains(content, "style_anchors_compliance: true") {
		t.Error("Style anchors compliance should be true")
	}

	// Check each task has style_anchors section
	taskIDs := []string{"m1-001", "m1-002", "m1-003", "m1-004", "m1-005", "m1-006", "m1-007"}
	for _, taskID := range taskIDs {
		taskSection := extractTaskSection(content, taskID)
		if !contains(taskSection, "style_anchors:") {
			t.Errorf("Task %s missing style_anchors section", taskID)
		}
	}
}

// TestDependencies validates dependency structure
func TestDependencies(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	// Check validation summary indicates acyclic dependencies
	if !contains(content, "dependencies_acyclic: true") {
		t.Error("Dependencies should be marked as acyclic")
	}

	// Check dependency graph exists
	if !contains(content, "dependency_graph:") {
		t.Error("Missing dependency_graph section")
	}
}

// TestBasedOnFiles validates that breakdown references correct input files
func TestBasedOnFiles(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	requiredFiles := []string{
		"docs/implementation-plan/m1/structured_analysis.json",
		"docs/implementation-plan/m0/patterns_report.json",
	}

	for _, file := range requiredFiles {
		if !contains(content, file) {
			t.Errorf("Task breakdown does not reference required file: %s", file)
		}
	}
}

// TestAcceptanceCriteria validates each task has acceptance criteria
func TestAcceptanceCriteria(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	taskIDs := []string{"m1-001", "m1-002", "m1-003", "m1-004", "m1-005", "m1-006", "m1-007"}
	for _, taskID := range taskIDs {
		taskSection := extractTaskSection(content, taskID)
		if !contains(taskSection, "acceptance_criteria:") {
			t.Errorf("Task %s missing acceptance_criteria section", taskID)
		}
	}
}

// Helper functions
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func index(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func extractTaskSection(content, taskID string) string {
	taskStart := "- id: \"" + taskID + "\""
	startIdx := index(content, taskStart)
	if startIdx == -1 {
		return ""
	}

	// Find the end of this task (start of next task or end of tasks section)
	remaining := content[startIdx:]

	// Look for next task
	nextTaskIdx := -1
	for i := 1; i < 8; i++ {
		nextID := taskID[:3] + string(rune('0'+i/10)) + string(rune('0'+i%10))
		if nextID != taskID {
			nextTaskPattern := "- id: \"" + nextID + "\""
			if idx := index(remaining, nextTaskPattern); idx != -1 && (nextTaskIdx == -1 || idx < nextTaskIdx) {
				nextTaskIdx = idx
			}
		}
	}

	if nextTaskIdx != -1 {
		return remaining[:nextTaskIdx]
	}
	return remaining
}
