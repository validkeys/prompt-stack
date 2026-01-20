package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kyledavis/prompt-stack/internal/cli/prompt"
	"github.com/spf13/cobra"
)

var (
	requirementsOutput string
)

var requirementsCmd = &cobra.Command{
	Use:   "requirements",
	Short: "Interactive requirements gathering for planning input",
	Long:  `Run an interactive interview to gather planning input and save to YAML for downstream Plan Mode.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		fmt.Println("=== Planning Input Requirements Gathering ===")
		fmt.Println("This will ask you a series of questions to define planning input for the Plan Mode.")
		fmt.Println("Press Ctrl+C to cancel at any time.")
		fmt.Println()

		questions := PlanningQuestions()
		p := prompt.NewPrompt(questions)

		result, err := p.Run(ctx)
		if err != nil {
			return fmt.Errorf("interview failed: %w", err)
		}

		if err := savePlanningResult(result, requirementsOutput); err != nil {
			return fmt.Errorf("failed to save planning results: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(requirementsCmd)

	defaultDir := filepath.Join("docs", "implementation-plan", "m1")
	requirementsCmd.Flags().StringVarP(&requirementsOutput, "output", "o", defaultDir, "Directory to save planning input YAML")
}

func PlanningQuestions() []prompt.Question {
	return []prompt.Question{
		{
			ID:       "id",
			Text:     "What is the short slug for this plan? (required, e.g., m1, m2)",
			Required: true,
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("ID cannot be empty")
				}
				return nil
			},
		},
		{
			ID:       "title",
			Text:     "What is the one-line title for this plan? (required)",
			Required: true,
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("title cannot be empty")
				}
				return nil
			},
		},
		{
			ID:       "short_description",
			Text:     "Provide a 1-2 sentence goal description. (required)",
			Required: true,
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("short description cannot be empty")
				}
				return nil
			},
		},
		{
			ID:       "background",
			Text:     "Provide 1-3 sentences giving context. (optional)",
			Required: false,
		},
		{
			ID:       "objectives",
			Text:     "What are the measurable objectives? (max 5, one per line)",
			Required: true,
			Validate: func(s string) error {
				lines := strings.Split(strings.TrimSpace(s), "\n")
				validLines := 0
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						validLines++
					}
				}
				if validLines == 0 {
					return fmt.Errorf("at least one objective is required")
				}
				if validLines > 5 {
					return fmt.Errorf("please limit to 5 objectives maximum")
				}
				return nil
			},
		},
		{
			ID:       "success_metrics",
			Text:     "What are the success metrics? (format: metric: target, one per line)",
			Required: true,
		},
		{
			ID:       "requirements_file",
			Text:     "What is the path to the full requirements doc? (required, e.g., docs/implementation-plan/m1/requirements.md)",
			Required: true,
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("requirements file path cannot be empty")
				}
				return nil
			},
		},
		{
			ID:       "style_anchors",
			Text:     "What are the style anchor files or code references? (2-3 paths or URLs, one per line)",
			Required: true,
			Validate: func(s string) error {
				lines := strings.Split(strings.TrimSpace(s), "\n")
				validLines := 0
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						validLines++
					}
				}
				if validLines == 0 {
					return fmt.Errorf("at least one style anchor is required")
				}
				if validLines > 3 {
					return fmt.Errorf("please limit to 3 style anchors maximum")
				}
				return nil
			},
		},
		{
			ID:       "start_date",
			Text:     "What is the start date? (format: YYYY-MM-DD, required)",
			Required: true,
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("start date cannot be empty")
				}
				return nil
			},
		},
		{
			ID:       "target_completion",
			Text:     "What is the target completion date? (format: YYYY-MM-DD, required)",
			Required: true,
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("target completion date cannot be empty")
				}
				return nil
			},
		},
		{
			ID:       "scope_in",
			Text:     "What is in scope? (concise list, one per line)",
			Required: true,
		},
		{
			ID:       "scope_out",
			Text:     "What is out of scope? (concise list, one per line)",
			Required: true,
		},
		{
			ID:       "constraints",
			Text:     "What are the constraints? (e.g., browsers, platforms, budget, one per line)",
			Required: true,
		},
		{
			ID:       "assumptions",
			Text:     "What are the assumptions? (one per line)",
			Required: true,
		},
		{
			ID:       "deliverables",
			Text:     "What are the deliverables? (list each with name and description, one per line)",
			Required: true,
		},
		{
			ID:       "acceptance_criteria",
			Text:     "What are the acceptance criteria? (list each with id, title, expected outcome, one per line)",
			Required: true,
		},
		{
			ID:       "tech_stack_languages",
			Text:     "What programming languages are used? (comma-separated)",
			Required: true,
		},
		{
			ID:       "tech_stack_frameworks",
			Text:     "What frameworks are used? (comma-separated)",
			Required: false,
		},
		{
			ID:       "tech_stack_infra",
			Text:     "What infrastructure is used? (comma-separated)",
			Required: false,
		},
		{
			ID:       "integrations",
			Text:     "What integrations are needed? (system: notes, one per line)",
			Required: false,
		},
		{
			ID:       "attachments",
			Text:     "What attachment paths or URLs? (one per line, optional)",
			Required: false,
		},
		{
			ID:       "repo_access",
			Text:     "What is the repo path or URL for read access? (required)",
			Required: true,
		},
		{
			ID:       "require_unit_tests",
			Text:     "Are unit tests required? (yes/no, required)",
			Required: true,
			Validate: func(s string) error {
				s = strings.ToLower(strings.TrimSpace(s))
				if s != "yes" && s != "no" {
					return fmt.Errorf("please enter 'yes' or 'no'")
				}
				return nil
			},
		},
		{
			ID:       "require_integration_tests",
			Text:     "Are integration tests required? (yes/no, required)",
			Required: true,
			Validate: func(s string) error {
				s = strings.ToLower(strings.TrimSpace(s))
				if s != "yes" && s != "no" {
					return fmt.Errorf("please enter 'yes' or 'no'")
				}
				return nil
			},
		},
		{
			ID:       "require_e2e",
			Text:     "Are e2e tests required? (yes/no)",
			Required: false,
			Validate: func(s string) error {
				if s == "" {
					return nil
				}
				s = strings.ToLower(strings.TrimSpace(s))
				if s != "yes" && s != "no" {
					return fmt.Errorf("please enter 'yes' or 'no'")
				}
				return nil
			},
		},
		{
			ID:       "data_classification",
			Text:     "What is the data classification? (public/internal/confidential, required)",
			Required: true,
			Validate: func(s string) error {
				s = strings.ToLower(strings.TrimSpace(s))
				if s != "public" && s != "internal" && s != "confidential" {
					return fmt.Errorf("please enter 'public', 'internal', or 'confidential'")
				}
				return nil
			},
		},
		{
			ID:       "secrets_included",
			Text:     "Are secrets included in the documentation? (yes/no, required)",
			Required: true,
			Validate: func(s string) error {
				s = strings.ToLower(strings.TrimSpace(s))
				if s != "yes" && s != "no" {
					return fmt.Errorf("please enter 'yes' or 'no'")
				}
				return nil
			},
		},
	}
}

func savePlanningResult(result *prompt.InterviewResult, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	transcriptPath := ".prompt-stack/requirements-transcript.txt"
	if err := os.MkdirAll(filepath.Dir(transcriptPath), 0755); err != nil {
		return fmt.Errorf("failed to create transcript directory: %w", err)
	}
	if err := os.WriteFile(transcriptPath, []byte(result.Transcript), 0644); err != nil {
		return fmt.Errorf("failed to write transcript: %w", err)
	}
	fmt.Printf("✓ Saved transcript to %s\n", transcriptPath)

	yamlPath := filepath.Join(outputDir, "planning-input.yaml")
	yamlContent := generatePlanningYAML(result)
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		return fmt.Errorf("failed to write planning YAML: %w", err)
	}
	fmt.Printf("✓ Saved planning input to %s\n", yamlPath)

	fmt.Println("\n✓ Planning input generation complete!")
	fmt.Printf("  Transcript: %s\n", transcriptPath)
	fmt.Printf("  Planning input: %s\n", yamlPath)

	return nil
}

func generatePlanningYAML(result *prompt.InterviewResult) string {
	return fmt.Sprintf(`# Planning Input YAML

id: "%s"
title: "%s"
short_description: "%s"
background: "%s"

objectives:
%s

success_metrics:
%s

requirements_file: "%s"
style_anchors:
%s

timeline:
  start_date: "%s"
  target_completion: "%s"

scope:
  in_scope:
%s
  out_of_scope:
%s

constraints:
%s

assumptions:
%s

deliverables:
%s

acceptance_criteria:
%s

tech_stack:
  languages: [%s]
  frameworks: [%s]
  infra: [%s]

integrations:
%s

attachments:
%s

repo_access:
  repo: "%s"
  read_only: true

testing:
  require_unit_tests: %t
  require_integration_tests: %t
  require_e2e: %t

data_classification: "%s"
secrets_included: %t
`,
		result.Responses["id"],
		result.Responses["title"],
		result.Responses["short_description"],
		result.Responses["background"],
		formatYAMLList(result.Responses["objectives"]),
		formatSuccessMetrics(result.Responses["success_metrics"]),
		result.Responses["requirements_file"],
		formatYAMLList(result.Responses["style_anchors"]),
		result.Responses["start_date"],
		result.Responses["target_completion"],
		formatYAMLList(result.Responses["scope_in"]),
		formatYAMLList(result.Responses["scope_out"]),
		formatYAMLList(result.Responses["constraints"]),
		formatYAMLList(result.Responses["assumptions"]),
		formatDeliverables(result.Responses["deliverables"]),
		formatAcceptanceCriteria(result.Responses["acceptance_criteria"]),
		result.Responses["tech_stack_languages"],
		result.Responses["tech_stack_frameworks"],
		result.Responses["tech_stack_infra"],
		formatIntegrations(result.Responses["integrations"]),
		formatYAMLList(result.Responses["attachments"]),
		result.Responses["repo_access"],
		strings.ToLower(result.Responses["require_unit_tests"]) == "yes",
		strings.ToLower(result.Responses["require_integration_tests"]) == "yes",
		result.Responses["require_e2e"] != "" && strings.ToLower(result.Responses["require_e2e"]) == "yes",
		strings.ToLower(result.Responses["data_classification"]),
		strings.ToLower(result.Responses["secrets_included"]) == "yes",
	)
}

func formatYAMLList(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var result strings.Builder
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result.WriteString(fmt.Sprintf("  - \"%s\"\n", line))
		}
	}
	return result.String()
}

func formatSuccessMetrics(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var result strings.Builder
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				result.WriteString(fmt.Sprintf("  - metric: \"%s\"\n", strings.TrimSpace(parts[0])))
				result.WriteString(fmt.Sprintf("    target: \"%s\"\n", strings.TrimSpace(parts[1])))
			} else {
				result.WriteString(fmt.Sprintf("  - metric: \"%s\"\n", line))
				result.WriteString("    target: \"\"\n")
			}
		}
	}
	return result.String()
}

func formatDeliverables(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var result strings.Builder
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result.WriteString(fmt.Sprintf("  - name: \"deliverable-%d\"\n", i+1))
			result.WriteString(fmt.Sprintf("    description: \"%s\"\n", line))
			result.WriteString("    owner: \"\"\n")
			result.WriteString("    format: \"\"\n")
			result.WriteString("    due: \"\"\n")
		}
	}
	return result.String()
}

func formatAcceptanceCriteria(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var result strings.Builder
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result.WriteString(fmt.Sprintf("  - id: \"AC-%d\"\n", i+1))
			result.WriteString(fmt.Sprintf("    title: \"%s\"\n", line))
			result.WriteString("    scenario: \"\"\n")
			result.WriteString("    expected_outcome: \"\"\n")
			result.WriteString("    validation_method: \"\"\n")
			result.WriteString("    stakeholder_signoff: \"\"\n")
			result.WriteString("    related_deliverables: []\n")
		}
	}
	return result.String()
}

func formatIntegrations(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var result strings.Builder
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				result.WriteString(fmt.Sprintf("  - system: \"%s\"\n", strings.TrimSpace(parts[0])))
				result.WriteString(fmt.Sprintf("    notes: \"%s\"\n", strings.TrimSpace(parts[1])))
			} else {
				result.WriteString(fmt.Sprintf("  - system: \"%s\"\n", line))
				result.WriteString("    notes: \"\"\n")
			}
		}
	}
	return result.String()
}
