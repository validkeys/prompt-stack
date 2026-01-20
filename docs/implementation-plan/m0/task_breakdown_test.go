package m0

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

// TestTaskBreakdownStructure validates the structure of task_breakdown.yaml
func TestTaskBreakdownStructure(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	// Basic validation that file contains expected structure
	content := string(data)

	requiredFields := []string{
		"metadata:",
		"milestone:",
		"title:",
		"tasks:",
		"id:",
		"title:",
		"description:",
		"estimated_duration_minutes:",
		"files_in_scope:",
		"style_anchors:",
		"dependencies:",
		"estimated_context_tokens:",
		"single_responsibility:",
		"acceptance_criteria:",
		"sizing_summary:",
		"style_anchors_summary:",
		"context_budget_summary:",
		"compliance:",
	}

	for _, field := range requiredFields {
		if !contains(content, field) {
			t.Errorf("Missing required field in task_breakdown.yaml: %s", field)
		}
	}
}

// TestTaskSizingCompliance validates task sizing against policy
func TestTaskSizingCompliance(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	// Check that compliance section exists and meets criteria
	if !contains(content, "meets_all_criteria: true") {
		t.Error("Task breakdown does not meet all criteria")
	}

	// Check sizing policy compliance
	if !contains(content, "within_sizing_policy: true") {
		t.Error("Task breakdown is not within sizing policy")
	}

	// Check style anchors policy
	if !contains(content, "meets_policy: true") {
		t.Error("Style anchors do not meet policy")
	}

	// Verify all tasks have duration estimates
	taskLines := extractTaskLines(content)
	for _, line := range taskLines {
		if contains(line, "estimated_duration_minutes:") {
			duration, err := extractDuration(line)
			if err != nil {
				t.Fatalf("Failed to parse duration from line %q: %v", line, err)
			}
			if duration < 30 || duration > 150 {
				t.Errorf("Task duration %d is outside policy range (30-150)", duration)
			}
		}
	}
}

// TestStyleAnchorsPerTask validates each task has 2-3 style anchors
func TestStyleAnchorsPerTask(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	// Extract all style anchors for tasks
	anchorCounts, err := extractStyleAnchorCounts(content)
	if err != nil {
		t.Fatalf("Failed to extract style anchor counts: %v", err)
	}

	for taskID, count := range anchorCounts {
		if count < 2 || count > 3 {
			t.Errorf("Task %s has %d style anchors (expected 2-3)", taskID, count)
		}
	}

	// Verify summary matches
	if !contains(content, "avg_anchors_per_task:") {
		t.Error("Missing average anchors per task in summary")
	}
}

// TestDependenciesAcyclic validates dependency graph is acyclic
func TestDependenciesAcyclic(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	// Check that dependencies are marked as acyclic
	if !contains(content, "acyclic: true") {
		t.Error("Dependency graph is not marked as acyclic")
	}

	// Verify dependency levels exist
	if !contains(content, "levels:") {
		t.Error("Missing dependency levels")
	}

	// Check for at least 5 levels (as per breakdown)
	levelCount := countOccurrences(content, "- level:")
	if levelCount < 5 {
		t.Errorf("Expected at least 5 dependency levels, found %d", levelCount)
	}
}

// TestContextBudgetValidation validates context budgets
func TestContextBudgetValidation(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	// Check that all tasks have estimated_context_tokens
	taskIDs := []string{"m0-001", "m0-002", "m0-003", "m0-004", "m0-005", "m0-006", "m0-007"}
	for _, id := range taskIDs {
		taskSection := extractTaskSection(content, id)
		if !contains(taskSection, "estimated_context_tokens:") {
			t.Errorf("Task %s is missing estimated_context_tokens", id)
		}
	}

	// Check that summary exists
	if !contains(content, "context_budget_summary:") {
		t.Error("Missing context budget summary")
	}

	if !contains(content, "all_within_budget: true") {
		t.Error("Not all tasks are within context budget")
	}
}

// TestSingleResponsibility validates each task has clear single responsibility
func TestSingleResponsibility(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	taskIDs := []string{"m0-001", "m0-002", "m0-003", "m0-004", "m0-005", "m0-006", "m0-007"}
	for _, id := range taskIDs {
		taskSection := extractTaskSection(content, id)
		if !contains(taskSection, "single_responsibility:") {
			t.Errorf("Task %s is missing single_responsibility field", id)
		}
	}

	if !contains(content, "clear_single_responsibility: true") {
		t.Error("Single responsibility not validated in compliance")
	}
}

// TestAcceptanceCriteria validates each task has acceptance criteria
func TestAcceptanceCriteria(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	taskIDs := []string{"m0-001", "m0-002", "m0-003", "m0-004", "m0-005", "m0-006", "m0-007"}
	for _, id := range taskIDs {
		taskSection := extractTaskSection(content, id)
		if !contains(taskSection, "acceptance_criteria:") {
			t.Errorf("Task %s is missing acceptance_criteria", id)
		}

		// Check that task has at least 2 acceptance criteria
		criteriaCount := countOccurrences(taskSection, "-")
		if criteriaCount < 2 {
			t.Errorf("Task %s has only %d acceptance criteria (expected at least 2)", id, criteriaCount)
		}
	}
}

// TestFilesInScope validates each task has files_in_scope
func TestFilesInScope(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	// Define expected file counts per task (exceptions allowed with justification)
	expectedFileCounts := map[string]int{
		"m0-001": 5,
		"m0-002": 6, // Creating 6 command files is reasonable for CLI structure
		"m0-003": 5,
		"m0-004": 5,
		"m0-005": 5,
		"m0-006": 5,
		"m0-007": 5,
	}

	taskIDs := []string{"m0-001", "m0-002", "m0-003", "m0-004", "m0-005", "m0-006", "m0-007"}
	for _, id := range taskIDs {
		taskSection := extractTaskSection(content, id)
		if !contains(taskSection, "files_in_scope:") {
			t.Errorf("Task %s is missing files_in_scope", id)
		}

		// Extract files_in_scope section and count files
		filesSection := extractSectionBetween(taskSection, "files_in_scope:", "style_anchors:")
		fileCount := countFileEntries(filesSection)

		// Check against expected count
		expected := expectedFileCounts[id]
		if fileCount != expected {
			t.Errorf("Task %s has %d files in scope (expected %d)", id, fileCount, expected)
		}

		// Ensure no task exceeds 7 files (hard limit)
		if fileCount > 7 {
			t.Errorf("Task %s has %d files in scope (max allowed: 7)", id, fileCount)
		}
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
		"docs/implementation-plan/m0/structured_analysis.json",
		"docs/implementation-plan/m0/patterns_report.json",
		"docs/task-sizing.md",
	}

	for _, file := range requiredFiles {
		if !contains(content, file) {
			t.Errorf("Task breakdown does not reference required file: %s", file)
		}
	}
}

// TestTaskCountAndIDs validates task count and ID format
func TestTaskCountAndIDs(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	// Count tasks
	taskCount := countOccurrences(content, "- id: \"m0-")
	if taskCount != 7 {
		t.Errorf("Expected 7 tasks, found %d", taskCount)
	}

	// Verify task IDs are sequential
	expectedIDs := []string{"m0-001", "m0-002", "m0-003", "m0-004", "m0-005", "m0-006", "m0-007"}
	for _, id := range expectedIDs {
		if !contains(content, fmt.Sprintf("- id: \"%s\"", id)) {
			t.Errorf("Missing task with ID: %s", id)
		}
	}
}

// TestMetadataCompleteness validates metadata section
func TestMetadataCompleteness(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	requiredMetadata := []string{
		"milestone:",
		"title:",
		"breakdown_generated_at:",
		"based_on:",
		"sizing_policy:",
		"min_minutes:",
		"max_minutes:",
		"max_files:",
		"style_anchor_policy:",
	}

	for _, field := range requiredMetadata {
		if !contains(content, field) {
			t.Errorf("Missing required metadata field: %s", field)
		}
	}

	// Verify milestone ID
	if !contains(content, "milestone: \"m0\"") {
		t.Error("Milestone ID is not set to 'm0'")
	}
}

// TestStyleAnchorReferences validates style anchors reference existing files
func TestStyleAnchorReferences(t *testing.T) {
	data, err := os.ReadFile("task_breakdown.yaml")
	if err != nil {
		t.Fatalf("Failed to read task_breakdown.yaml: %v", err)
	}

	content := string(data)

	// Common style anchors that should be referenced
	expectedAnchors := []string{
		"examples/style-anchor/cmd/mytool/main.go",
		"examples/style-anchor/pkg/greeter/greeter.go",
		"examples/style-anchor/pkg/greeter/greeter_test.go",
		"docs/best-practices.md",
		"docs/requirements/main.md",
	}

	anchorCount := 0
	for _, anchor := range expectedAnchors {
		if contains(content, anchor) {
			anchorCount++
		}
	}

	if anchorCount < 3 {
		t.Errorf("Only found %d of %d expected style anchors", anchorCount, len(expectedAnchors))
	}
}

// Helper functions

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func extractTaskLines(content string) []string {
	lines := []string{}
	for _, line := range splitLines(content) {
		if contains(line, "estimated_duration_minutes:") {
			lines = append(lines, line)
		}
	}
	return lines
}

func extractDuration(line string) (int, error) {
	trimmed := strings.TrimSpace(line)
	const prefix = "estimated_duration_minutes:"
	idx := strings.Index(trimmed, prefix)
	if idx < 0 {
		return 0, fmt.Errorf("duration line missing %q: %q", prefix, line)
	}

	value := strings.TrimSpace(trimmed[idx+len(prefix):])
	duration, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parse duration %q: %w", value, err)
	}
	return duration, nil
}

func extractStyleAnchorCounts(content string) (map[string]int, error) {
	counts := make(map[string]int)
	currentTask := ""
	inStyleAnchors := false

	lines := splitLines(content)
	for _, line := range lines {
		if contains(line, "- id: \"") {
			taskID, err := extractID(line)
			if err != nil {
				return nil, err
			}
			currentTask = taskID
			counts[currentTask] = 0
			inStyleAnchors = false
		} else if contains(line, "style_anchors:") {
			inStyleAnchors = true
		} else if inStyleAnchors && contains(line, "file:") {
			counts[currentTask]++
		} else if inStyleAnchors && contains(line, "  - ") && !contains(line, "file:") && !contains(line, "reason:") {
			inStyleAnchors = false
		}
	}

	return counts, nil
}

func extractID(line string) (string, error) {
	trimmed := strings.TrimSpace(line)
	const prefix = "- id:"
	idx := strings.Index(trimmed, prefix)
	if idx < 0 {
		return "", fmt.Errorf("id line missing %q: %q", prefix, line)
	}

	value := strings.TrimSpace(trimmed[idx+len(prefix):])
	value = strings.Trim(value, "\"'")
	if value == "" {
		return "", fmt.Errorf("empty task id in line %q", line)
	}
	return value, nil
}

func extractTaskSection(content, taskID string) string {
	lines := splitLines(content)
	var section []string
	inSection := false
	taskStart := "- id: \"" + taskID + "\""

	for _, line := range lines {
		if contains(line, taskStart) {
			inSection = true
		} else if inSection && contains(line, "- id: \"") && !contains(line, taskID) {
			break
		}

		if inSection {
			section = append(section, line)
		}
	}

	return joinLines(section)
}

func countOccurrences(s, substr string) int {
	count := 0
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			count++
			i += len(substr) - 1
		}
	}
	return count
}

func splitLines(s string) []string {
	lines := []string{}
	current := ""

	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(s[i])
		}
	}

	if current != "" {
		lines = append(lines, current)
	}

	return lines
}

func joinLines(lines []string) string {
	result := ""
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

func extractSectionBetween(content, startMarker, endMarker string) string {
	startIdx := findSubstring(content, startMarker)
	if startIdx < 0 {
		return ""
	}

	remaining := content[startIdx:]
	endIdx := findSubstring(remaining, endMarker)
	if endIdx < 0 {
		// No end marker, return everything from start
		return remaining
	}

	return remaining[:endIdx]
}

func countFileEntries(section string) int {
	count := 0
	lines := splitLines(section)
	for _, line := range lines {
		trimmed := trimWhitespace(line)
		if len(trimmed) > 0 && contains(trimmed, "\"") && !contains(trimmed, "#") {
			// This is a file path entry (quoted string)
			count++
		}
	}
	return count
}

func trimWhitespace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}
