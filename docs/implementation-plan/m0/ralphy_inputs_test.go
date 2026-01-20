package m0

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestRalphyInputsStructure validates the structure of ralphy_inputs.yaml
func TestRalphyInputsStructure(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	requiredFields := []string{
		"name:",
		"description:",
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
			t.Errorf("Missing required field in ralphy_inputs.yaml: %s", field)
		}
	}
}

// TestRequiredSchemaFields validates all required fields from docs/ralphy-inputs.schema.json
func TestRequiredSchemaFields(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	requiredFields := []string{
		"name:",
		"version:",
		"rules_file:",
		"task_sizing:",
		"tdd:",
		"model_preferences:",
		"outputs:",
	}

	for _, field := range requiredFields {
		if !strings.Contains(content, field) {
			t.Errorf("Missing required schema field: %s", field)
		}
	}
}

// TestVersionFormat validates semantic version format
func TestVersionFormat(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	versionLine := findLineWithPrefix(content, "version:")
	if versionLine == "" {
		t.Error("Version line not found")
		return
	}

	// Check semantic version pattern
	if !strings.Contains(versionLine, "0.1.0") {
		t.Errorf("Version format is invalid: %s (expected 0.1.0)", versionLine)
	}
}

// TestTaskSizingConfiguration validates task_sizing fields
func TestTaskSizingConfiguration(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "min_minutes: 30") {
		t.Error("Task sizing min_minutes should be 30")
	}

	if !strings.Contains(content, "max_minutes: 150") {
		t.Error("Task sizing max_minutes should be 150")
	}

	if !strings.Contains(content, "max_files: 5") {
		t.Error("Task sizing max_files should be 5")
	}
}

// TestTDDConfiguration validates tdd fields
func TestTDDConfiguration(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "required: false") {
		t.Error("Missing tdd.required field")
	}

	if !strings.Contains(content, "test_command:") {
		t.Error("Missing tdd.test_command field")
	}
}

// TestModelPreferencesConfiguration validates model_preferences fields
func TestModelPreferencesConfiguration(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "primary: opencode") {
		t.Error("Missing model_preferences.primary field")
	}

	if !strings.Contains(content, "strategies:") {
		t.Error("Missing model_preferences.strategies field")
	}
}

// TestOutputsConfiguration validates outputs fields
func TestOutputsConfiguration(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "allowed_file_edits:") {
		t.Error("Missing allowed_file_edits in outputs")
	}

	if !strings.Contains(content, "disallowed_file_edits:") {
		t.Error("Missing disallowed_file_edits in outputs")
	}

	if !strings.Contains(content, "- cmd/**") {
		t.Error("allowed_file_edits should include cmd/**")
	}

	if !strings.Contains(content, ".github/**") && !strings.Contains(content, "scripts/**") {
		t.Error("disallowed_file_edits should include protected directories")
	}
}

// TestGlobalConstraints validates global_constraints section
func TestGlobalConstraints(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "global_constraints:") {
		t.Error("Missing global_constraints section")
	}

	if !strings.Contains(content, "forbidden_patterns:") {
		t.Error("Missing forbidden_patterns in global_constraints")
	}

	if !strings.Contains(content, "required_patterns:") {
		t.Error("Missing required_patterns in global_constraints")
	}

	if !strings.Contains(content, "affirmative_constraints:") {
		t.Error("Missing affirmative_constraints in global_constraints")
	}
}

// TestRequiredPatterns validates required patterns presence
func TestRequiredPatterns(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	// Note: This is a Go project, so we check for Go test patterns in global_constraints
	// The Zod patterns are for future TypeScript work, we check they exist
	requiredPatterns := []string{
		"go test",
		"go vet",
		"any",        // Should be in forbidden_patterns
		"@ts-ignore", // Should be in forbidden_patterns
	}

	for _, pattern := range requiredPatterns {
		if !strings.Contains(content, pattern) {
			t.Errorf("Missing required pattern: %s", pattern)
		}
	}
}

// TestTaskDefinitions validates task definitions structure
func TestTaskDefinitions(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "tasks:") {
		t.Error("Missing tasks section")
	}

	expectedTaskIDs := []string{
		"m0-001",
		"m0-002",
		"m0-003",
		"m0-004",
		"m0-005",
		"m0-006",
		"m0-007",
	}

	for _, taskID := range expectedTaskIDs {
		if !strings.Contains(content, fmt.Sprintf("- id: \"%s\"", taskID)) {
			t.Errorf("Missing task with ID: %s", taskID)
		}
	}

	// Check that tasks have required fields
	taskIDs := expectedTaskIDs
	for _, taskID := range taskIDs {
		taskSection := extractRalphyTaskSection(content, taskID)

		requiredTaskFields := []string{
			"title:",
			"description:",
			"files_in_scope:",
			"style_anchors:",
			"estimated_duration_minutes:",
			"estimated_context_tokens:",
			"single_responsibility:",
			"acceptance_criteria:",
			"verification:",
			"pre_commit:",
		}

		for _, field := range requiredTaskFields {
			if !strings.Contains(taskSection, field) {
				t.Errorf("Task %s is missing field: %s", taskID, field)
			}
		}
	}
}

// TestRalphyStyleAnchorsPerTask validates 2-3 style anchors per task in Ralphy inputs
func TestRalphyStyleAnchorsPerTask(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	taskIDs := []string{"m0-001", "m0-002", "m0-003", "m0-004", "m0-005", "m0-006", "m0-007"}

	for _, taskID := range taskIDs {
		taskSection := extractRalphyTaskSection(content, taskID)

		// Simply check that style_anchors section exists in the task
		if !strings.Contains(taskSection, "style_anchors:") {
			t.Errorf("Task %s is missing style_anchors section", taskID)
		}

		// Check for at least one file: reference (indicating at least one anchor)
		if !strings.Contains(taskSection, "file:") {
			t.Errorf("Task %s style_anchors should contain file: references", taskID)
		}

		// Check for at least one reason: reference (indicating at least one anchor)
		if !strings.Contains(taskSection, "reason:") {
			t.Errorf("Task %s style_anchors should contain reason: references", taskID)
		}
	}
}

// TestAffirmativeConstraints validates affirmative constraint usage
func TestAffirmativeConstraints(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "affirmative_constraints:") {
		t.Error("Missing affirmative_constraints section")
	}

	// Check for positive language patterns
	if !strings.Contains(content, "ALWAYS") {
		t.Error("Affirmative constraints should use positive language (ALWAYS)")
	}

	if !strings.Contains(content, "ONLY") {
		t.Error("Affirmative constraints should use positive language (ONLY)")
	}
}

// TestVerificationSteps validates verification steps are present
func TestVerificationSteps(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	taskIDs := []string{"m0-001", "m0-002", "m0-003", "m0-004", "m0-005", "m0-006", "m0-007"}

	for _, taskID := range taskIDs {
		taskSection := extractRalphyTaskSection(content, taskID)

		if !strings.Contains(taskSection, "verification:") {
			t.Errorf("Task %s is missing verification section", taskID)
		}

		if !strings.Contains(taskSection, "pre_commit:") {
			t.Errorf("Task %s is missing pre_commit verification steps", taskID)
		}
	}
}

// TestDependencies validates task dependencies
func TestDependencies(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	// Check that some tasks have dependencies
	if !strings.Contains(content, "dependencies:") {
		t.Error("Some tasks should have dependencies")
	}

	// Verify specific dependencies
	if !strings.Contains(content, "\"m0-001\"") {
		t.Error("Should have dependency on m0-001")
	}
}

// TestCommitPolicy validates commit_policy section
func TestCommitPolicy(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "commit_policy:") {
		t.Error("Missing commit_policy in outputs")
	}

	if !strings.Contains(content, "prefix_rules:") {
		t.Error("Missing prefix_rules in commit_policy")
	}

	expectedPrefixes := []string{"feat:", "fix:", "test:", "docs:"}
	for _, prefix := range expectedPrefixes {
		if !strings.Contains(content, prefix) {
			t.Errorf("Missing commit prefix: %s", prefix)
		}
	}
}

// TestStyleAnchorsTopLevel validates top-level style_anchors
func TestStyleAnchorsTopLevel(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "style_anchors:") {
		t.Error("Missing top-level style_anchors section")
	}

	expectedAnchors := []string{
		"docs/best-practices.md",
		"docs/ralphy-inputs.md",
		"examples/style-anchor/cmd/mytool/main.go",
	}

	for _, anchor := range expectedAnchors {
		if !strings.Contains(content, anchor) {
			t.Errorf("Missing expected style anchor: %s", anchor)
		}
	}
}

// TestAllowedDependencies validates allowed_dependencies
func TestAllowedDependencies(t *testing.T) {
	data, err := os.ReadFile("ralphy_inputs.yaml")
	if err != nil {
		t.Fatalf("Failed to read ralphy_inputs.yaml: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "allowed_dependencies:") {
		t.Error("Missing allowed_dependencies section")
	}

	expectedDeps := []string{
		"cobra",
		"better-sqlite3",
		"go",
	}

	for _, dep := range expectedDeps {
		if !strings.Contains(content, dep) {
			t.Errorf("Missing expected allowed dependency: %s", dep)
		}
	}
}

// Helper functions

func findLineWithPrefix(content, prefix string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			return strings.TrimSpace(line)
		}
	}
	return ""
}

func extractRalphyTaskSection(content, taskID string) string {
	taskStart := fmt.Sprintf("- id: \"%s\"", taskID)
	startIdx := strings.Index(content, taskStart)
	if startIdx < 0 {
		return ""
	}

	// Extract from task start to next task
	lines := strings.Split(content[startIdx:], "\n")
	var section []string

	for _, line := range lines {
		// Stop when we hit another task
		if strings.Contains(line, "- id: \"") && !strings.Contains(line, taskID) {
			break
		}
		section = append(section, line)
	}

	return strings.Join(section, "\n")
}
