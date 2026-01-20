package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kyledavis/prompt-stack/pkg/prompt"
)

func TestRequirementsCommandExists(t *testing.T) {
	t.Run("requirements_command_follows_Cobra_conventions", func(t *testing.T) {
		if requirementsCmd == nil {
			t.Fatal("requirementsCmd is not initialized")
		}

		if requirementsCmd.Use == "" {
			t.Error("requirementsCmd.Use is empty")
		}

		if requirementsCmd.Short == "" {
			t.Error("requirementsCmd.Short is empty")
		}

		if requirementsCmd.Long == "" {
			t.Error("requirementsCmd.Long is empty")
		}

		if requirementsCmd.RunE == nil {
			t.Error("requirementsCmd.RunE is nil")
		}
	})
}

func TestPlanningQuestions(t *testing.T) {
	questions := PlanningQuestions()

	t.Run("planning_questions_has_required_questions", func(t *testing.T) {
		requiredQuestions := []string{
			"id", "title", "short_description", "objectives", "success_metrics",
			"requirements_file", "style_anchors", "start_date", "target_completion",
			"scope_in", "scope_out", "constraints", "assumptions", "deliverables",
			"acceptance_criteria", "tech_stack_languages", "repo_access",
			"require_unit_tests", "require_integration_tests",
			"data_classification", "secrets_included",
		}

		questionIDs := make(map[string]bool)
		for _, q := range questions {
			questionIDs[q.ID] = true
		}

		for _, id := range requiredQuestions {
			if !questionIDs[id] {
				t.Errorf("Required question '%s' is missing", id)
			}
		}
	})

	t.Run("planning_questions_has_optional_questions", func(t *testing.T) {
		optionalQuestions := []string{"background", "tech_stack_frameworks", "tech_stack_infra", "integrations", "attachments", "require_e2e"}

		questionIDs := make(map[string]bool)
		for _, q := range questions {
			questionIDs[q.ID] = true
		}

		for _, id := range optionalQuestions {
			if !questionIDs[id] {
				t.Errorf("Optional question '%s' is missing", id)
			}
		}
	})

	t.Run("planning_questions_have_validation", func(t *testing.T) {
		questionsWithValidation := []string{"id", "title", "short_description", "objectives", "style_anchors", "start_date", "target_completion", "require_unit_tests", "require_integration_tests", "require_e2e", "data_classification", "secrets_included"}

		questionMap := make(map[string]prompt.Question)
		for _, q := range questions {
			questionMap[q.ID] = q
		}

		for _, id := range questionsWithValidation {
			q, exists := questionMap[id]
			if !exists {
				t.Errorf("Question '%s' not found", id)
				continue
			}

			if q.Validate == nil {
				t.Errorf("Question '%s' does not have validation function", id)
			}
		}
	})
}

func TestPlanningQuestionsValidation(t *testing.T) {
	questions := PlanningQuestions()
	questionMap := make(map[string]prompt.Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	t.Run("id_validation_rejects_empty", func(t *testing.T) {
		q := questionMap["id"]
		err := q.Validate("")
		if err == nil {
			t.Error("Expected error for empty ID, got nil")
		}
	})

	t.Run("id_validation_accepts_valid", func(t *testing.T) {
		q := questionMap["id"]
		err := q.Validate("m1")
		if err != nil {
			t.Errorf("Expected no error for valid ID, got: %v", err)
		}
	})

	t.Run("title_validation_rejects_empty", func(t *testing.T) {
		q := questionMap["title"]
		err := q.Validate("")
		if err == nil {
			t.Error("Expected error for empty title, got nil")
		}
	})

	t.Run("title_validation_accepts_valid", func(t *testing.T) {
		q := questionMap["title"]
		err := q.Validate("CLI scaffold implementation")
		if err != nil {
			t.Errorf("Expected no error for valid title, got: %v", err)
		}
	})

	t.Run("short_description_validation_rejects_empty", func(t *testing.T) {
		q := questionMap["short_description"]
		err := q.Validate("")
		if err == nil {
			t.Error("Expected error for empty short description, got nil")
		}
	})

	t.Run("objectives_validation_rejects_empty", func(t *testing.T) {
		q := questionMap["objectives"]
		err := q.Validate("")
		if err == nil {
			t.Error("Expected error for empty objectives, got nil")
		}
	})

	t.Run("objectives_validation_rejects_too_many", func(t *testing.T) {
		q := questionMap["objectives"]
		objs := "obj1\nobj2\nobj3\nobj4\nobj6\nobj7"
		err := q.Validate(objs)
		if err == nil {
			t.Error("Expected error for too many objectives, got nil")
		}
	})

	t.Run("objectives_validation_accepts_valid", func(t *testing.T) {
		q := questionMap["objectives"]
		objs := "obj1\nobj2\nobj3"
		err := q.Validate(objs)
		if err != nil {
			t.Errorf("Expected no error for valid objectives, got: %v", err)
		}
	})

	t.Run("style_anchors_validation_rejects_empty", func(t *testing.T) {
		q := questionMap["style_anchors"]
		err := q.Validate("")
		if err == nil {
			t.Error("Expected error for empty style anchors, got nil")
		}
	})

	t.Run("style_anchors_validation_rejects_too_many", func(t *testing.T) {
		q := questionMap["style_anchors"]
		anchors := "anchor1\nanchor2\nanchor3\nanchor4\nanchor5"
		err := q.Validate(anchors)
		if err == nil {
			t.Error("Expected error for too many style anchors, got nil")
		}
	})

	t.Run("style_anchors_validation_accepts_valid", func(t *testing.T) {
		q := questionMap["style_anchors"]
		anchors := "docs/style-markers.md\nexamples/style-anchor/"
		err := q.Validate(anchors)
		if err != nil {
			t.Errorf("Expected no error for valid style anchors, got: %v", err)
		}
	})

	t.Run("start_date_validation_rejects_empty", func(t *testing.T) {
		q := questionMap["start_date"]
		err := q.Validate("")
		if err == nil {
			t.Error("Expected error for empty start date, got nil")
		}
	})

	t.Run("target_completion_validation_rejects_empty", func(t *testing.T) {
		q := questionMap["target_completion"]
		err := q.Validate("")
		if err == nil {
			t.Error("Expected error for empty target completion, got nil")
		}
	})

	t.Run("require_unit_tests_validation_rejects_invalid", func(t *testing.T) {
		q := questionMap["require_unit_tests"]
		err := q.Validate("maybe")
		if err == nil {
			t.Error("Expected error for invalid require_unit_tests value, got nil")
		}
	})

	t.Run("require_unit_tests_validation_accepts_yes", func(t *testing.T) {
		q := questionMap["require_unit_tests"]
		err := q.Validate("yes")
		if err != nil {
			t.Errorf("Expected no error for valid require_unit_tests 'yes', got: %v", err)
		}
	})

	t.Run("require_unit_tests_validation_accepts_no", func(t *testing.T) {
		q := questionMap["require_unit_tests"]
		err := q.Validate("no")
		if err != nil {
			t.Errorf("Expected no error for valid require_unit_tests 'no', got: %v", err)
		}
	})

	t.Run("require_integration_tests_validation_rejects_invalid", func(t *testing.T) {
		q := questionMap["require_integration_tests"]
		err := q.Validate("maybe")
		if err == nil {
			t.Error("Expected error for invalid require_integration_tests value, got nil")
		}
	})

	t.Run("require_integration_tests_validation_accepts_yes", func(t *testing.T) {
		q := questionMap["require_integration_tests"]
		err := q.Validate("yes")
		if err != nil {
			t.Errorf("Expected no error for valid require_integration_tests 'yes', got: %v", err)
		}
	})

	t.Run("data_classification_validation_rejects_invalid", func(t *testing.T) {
		q := questionMap["data_classification"]
		err := q.Validate("secret")
		if err == nil {
			t.Error("Expected error for invalid data_classification value, got nil")
		}
	})

	t.Run("data_classification_validation_accepts_public", func(t *testing.T) {
		q := questionMap["data_classification"]
		err := q.Validate("public")
		if err != nil {
			t.Errorf("Expected no error for valid data_classification 'public', got: %v", err)
		}
	})

	t.Run("data_classification_validation_accepts_internal", func(t *testing.T) {
		q := questionMap["data_classification"]
		err := q.Validate("internal")
		if err != nil {
			t.Errorf("Expected no error for valid data_classification 'internal', got: %v", err)
		}
	})

	t.Run("data_classification_validation_accepts_confidential", func(t *testing.T) {
		q := questionMap["data_classification"]
		err := q.Validate("confidential")
		if err != nil {
			t.Errorf("Expected no error for valid data_classification 'confidential', got: %v", err)
		}
	})

	t.Run("secrets_included_validation_rejects_invalid", func(t *testing.T) {
		q := questionMap["secrets_included"]
		err := q.Validate("maybe")
		if err == nil {
			t.Error("Expected error for invalid secrets_included value, got nil")
		}
	})

	t.Run("secrets_included_validation_accepts_yes", func(t *testing.T) {
		q := questionMap["secrets_included"]
		err := q.Validate("yes")
		if err != nil {
			t.Errorf("Expected no error for valid secrets_included 'yes', got: %v", err)
		}
	})

	t.Run("secrets_included_validation_accepts_no", func(t *testing.T) {
		q := questionMap["secrets_included"]
		err := q.Validate("no")
		if err != nil {
			t.Errorf("Expected no error for valid secrets_included 'no', got: %v", err)
		}
	})
}

func TestGeneratePlanningYAML(t *testing.T) {
	result := &prompt.InterviewResult{
		Responses: map[string]string{
			"id":                        "m1",
			"title":                     "CLI scaffold implementation",
			"short_description":         "Implement Go/Cobra CLI scaffold for milestone M1",
			"background":                "Initial implementation of CLI tool",
			"objectives":                "Implement CLI structure\nAdd init command\nAdd validation",
			"success_metrics":           "quality score: 0.95\ndelivery: on time",
			"requirements_file":         "docs/implementation-plan/m1/requirements.md",
			"style_anchors":             "docs/style-markers.md\nexamples/style-anchor/",
			"start_date":                "2026-01-21",
			"target_completion":         "2026-01-25",
			"scope_in":                  "CLI commands\nConfig management",
			"scope_out":                 "Web interface",
			"constraints":               "Go 1.21\nCobra framework",
			"assumptions":               "User has Go installed",
			"deliverables":              "Binary\nDocumentation",
			"acceptance_criteria":       "CLI runs\nTests pass",
			"tech_stack_languages":      "Go, Bash",
			"tech_stack_frameworks":     "Cobra",
			"tech_stack_infra":          "GitHub Actions",
			"integrations":              "OpenCode: for AI assistance",
			"attachments":               "docs/architecture.md",
			"repo_access":               "github.com/kyledavis/prompt-stack",
			"require_unit_tests":        "yes",
			"require_integration_tests": "yes",
			"require_e2e":               "no",
			"data_classification":       "internal",
			"secrets_included":          "no",
		},
	}

	yaml := generatePlanningYAML(result)

	t.Run("generate_planning_yaml_contains_id", func(t *testing.T) {
		if yaml == "" {
			t.Fatal("Generated YAML is empty")
		}

		if !contains(yaml, `id: "m1"`) {
			t.Error("Generated YAML does not contain id")
		}
	})

	t.Run("generate_planning_yaml_contains_title", func(t *testing.T) {
		if !contains(yaml, `title: "CLI scaffold implementation"`) {
			t.Error("Generated YAML does not contain title")
		}
	})

	t.Run("generate_planning_yaml_contains_short_description", func(t *testing.T) {
		if !contains(yaml, `short_description: "Implement Go/Cobra CLI scaffold for milestone M1"`) {
			t.Error("Generated YAML does not contain short_description")
		}
	})

	t.Run("generate_planning_yaml_contains_objectives", func(t *testing.T) {
		if !contains(yaml, "- \"Implement CLI structure\"") {
			t.Error("Generated YAML does not contain objective 1")
		}
		if !contains(yaml, "- \"Add init command\"") {
			t.Error("Generated YAML does not contain objective 2")
		}
		if !contains(yaml, "- \"Add validation\"") {
			t.Error("Generated YAML does not contain objective 3")
		}
	})

	t.Run("generate_planning_yaml_contains_success_metrics", func(t *testing.T) {
		if !contains(yaml, `metric: "quality score"`) {
			t.Error("Generated YAML does not contain success metric 1")
		}
		if !contains(yaml, `target: "0.95"`) {
			t.Error("Generated YAML does not contain success metric target 1")
		}
	})

	t.Run("generate_planning_yaml_contains_requirements_file", func(t *testing.T) {
		if !contains(yaml, `requirements_file: "docs/implementation-plan/m1/requirements.md"`) {
			t.Error("Generated YAML does not contain requirements_file")
		}
	})

	t.Run("generate_planning_yaml_contains_style_anchors", func(t *testing.T) {
		if !contains(yaml, `- "docs/style-markers.md"`) {
			t.Error("Generated YAML does not contain style anchor 1")
		}
		if !contains(yaml, `- "examples/style-anchor/"`) {
			t.Error("Generated YAML does not contain style anchor 2")
		}
	})

	t.Run("generate_planning_yaml_contains_timeline", func(t *testing.T) {
		if !contains(yaml, `start_date: "2026-01-21"`) {
			t.Error("Generated YAML does not contain start_date")
		}
		if !contains(yaml, `target_completion: "2026-01-25"`) {
			t.Error("Generated YAML does not contain target_completion")
		}
	})

	t.Run("generate_planning_yaml_contains_scope", func(t *testing.T) {
		if !contains(yaml, "in_scope:") {
			t.Error("Generated YAML does not contain in_scope")
		}
		if !contains(yaml, "out_of_scope:") {
			t.Error("Generated YAML does not contain out_of_scope")
		}
	})

	t.Run("generate_planning_yaml_contains_tech_stack", func(t *testing.T) {
		if !contains(yaml, "languages: [Go, Bash]") {
			t.Error("Generated YAML does not contain languages")
		}
		if !contains(yaml, "frameworks: [Cobra]") {
			t.Error("Generated YAML does not contain frameworks")
		}
		if !contains(yaml, "infra: [GitHub Actions]") {
			t.Error("Generated YAML does not contain infra")
		}
	})

	t.Run("generate_planning_yaml_contains_testing", func(t *testing.T) {
		if !contains(yaml, "require_unit_tests: true") {
			t.Error("Generated YAML does not contain require_unit_tests")
		}
		if !contains(yaml, "require_integration_tests: true") {
			t.Error("Generated YAML does not contain require_integration_tests")
		}
		if !contains(yaml, "require_e2e: false") {
			t.Error("Generated YAML does not contain require_e2e")
		}
	})

	t.Run("generate_planning_yaml_contains_data_classification", func(t *testing.T) {
		if !contains(yaml, `data_classification: "internal"`) {
			t.Error("Generated YAML does not contain data_classification")
		}
	})

	t.Run("generate_planning_yaml_contains_secrets_included", func(t *testing.T) {
		if !contains(yaml, "secrets_included: false") {
			t.Error("Generated YAML does not contain secrets_included")
		}
	})
}

func TestSavePlanningResult(t *testing.T) {
	result := &prompt.InterviewResult{
		Responses: map[string]string{
			"id":                        "m1",
			"title":                     "Test milestone",
			"short_description":         "Test description",
			"background":                "Test background",
			"objectives":                "Test objective 1\nTest objective 2",
			"success_metrics":           "metric: target",
			"requirements_file":         "docs/requirements.md",
			"style_anchors":             "docs/style.md",
			"start_date":                "2026-01-21",
			"target_completion":         "2026-01-25",
			"scope_in":                  "Test in scope",
			"scope_out":                 "Test out of scope",
			"constraints":               "Test constraint",
			"assumptions":               "Test assumption",
			"deliverables":              "Test deliverable",
			"acceptance_criteria":       "Test criteria",
			"tech_stack_languages":      "Go",
			"tech_stack_frameworks":     "",
			"tech_stack_infra":          "",
			"integrations":              "",
			"attachments":               "",
			"repo_access":               "github.com/test/repo",
			"require_unit_tests":        "yes",
			"require_integration_tests": "yes",
			"require_e2e":               "",
			"data_classification":       "internal",
			"secrets_included":          "no",
		},
		Transcript: "Q: Test question\nA: Test answer\n",
	}

	t.Run("save_planning_result_creates_directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputDir := filepath.Join(tmpDir, "test-output")

		defer func() {
			os.RemoveAll(".prompt-stack")
		}()

		err := savePlanningResult(result, outputDir)
		if err != nil {
			t.Fatalf("Failed to save planning result: %v", err)
		}

		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			t.Error("Output directory was not created")
		}
	})

	t.Run("save_planning_result_creates_transcript", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputDir := filepath.Join(tmpDir, "test-output-2")

		transcriptPath := ".prompt-stack/requirements-transcript.txt"
		defer func() {
			os.RemoveAll(".prompt-stack")
		}()

		err := savePlanningResult(result, outputDir)
		if err != nil {
			t.Fatalf("Failed to save planning result: %v", err)
		}

		content, err := os.ReadFile(transcriptPath)
		if err != nil {
			t.Fatalf("Failed to read transcript: %v", err)
		}

		if string(content) != result.Transcript {
			t.Error("Transcript content does not match")
		}
	})

	t.Run("save_planning_result_creates_yaml", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputDir := filepath.Join(tmpDir, "test-output-3")

		defer func() {
			os.RemoveAll(".prompt-stack")
		}()

		err := savePlanningResult(result, outputDir)
		if err != nil {
			t.Fatalf("Failed to save planning result: %v", err)
		}

		yamlPath := filepath.Join(outputDir, "planning-input.yaml")
		content, err := os.ReadFile(yamlPath)
		if err != nil {
			t.Fatalf("Failed to read YAML: %v", err)
		}

		if !contains(string(content), `id: "m1"`) {
			t.Error("YAML does not contain id")
		}
	})
}

func TestRequirementsCommandIntegration(t *testing.T) {
	t.Run("requirements_command_has_help", func(t *testing.T) {
		if requirementsCmd.Short == "" {
			t.Error("requirements command does not have Short description")
		}

		if requirementsCmd.Long == "" {
			t.Error("requirements command does not have Long description")
		}
	})

	t.Run("requirements_command_has_output_flag", func(t *testing.T) {
		flag := requirementsCmd.Flags().Lookup("output")
		if flag == nil {
			t.Error("requirements command does not have --output flag")
		}

		if flag.Name != "output" {
			t.Errorf("Expected flag name 'output', got '%s'", flag.Name)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
