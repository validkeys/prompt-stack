package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyledavis/prompt-stack/internal/cli/prompt"
	"github.com/kyledavis/prompt-stack/internal/config"
	"github.com/kyledavis/prompt-stack/internal/knowledge/database"
	"github.com/spf13/cobra"
)

var (
	interactiveOutputDir string
	noInteractive        bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new milestone with interactive requirements gathering",
	Long:  `Run an interactive interview to gather milestone requirements and save them to YAML and markdown files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		fmt.Println("=== Prompt Stack Initialization ===")

		configPath := ".prompt-stack/config.yaml"
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Printf("Creating config file at %s\n", configPath)
			if err := config.Init(configPath); err != nil {
				return fmt.Errorf("failed to initialize config: %w", err)
			}
			fmt.Printf("✓ Created %s\n", configPath)
		} else {
			fmt.Printf("✓ Config already exists at %s\n", configPath)
		}

		dbPath := ".prompt-stack/knowledge.db"
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			fmt.Printf("Creating database at %s\n", dbPath)
			if err := database.Init(dbPath); err != nil {
				return fmt.Errorf("failed to initialize database: %w", err)
			}
			fmt.Printf("✓ Created %s\n", dbPath)
		} else {
			fmt.Printf("✓ Database already exists at %s\n", dbPath)
		}

		if !noInteractive {
			fmt.Println("\n=== Requirements Gathering Interview ===")
			fmt.Println("This will ask you a series of questions to define your milestone requirements.")
			fmt.Println("Press Ctrl+C to cancel at any time.")
			fmt.Println()

			questions := prompt.DefaultQuestions()
			p := prompt.NewPrompt(questions)

			result, err := p.Run(ctx)
			if err != nil {
				return fmt.Errorf("interview failed: %w", err)
			}

			if err := saveInterviewResult(result, interactiveOutputDir); err != nil {
				return fmt.Errorf("failed to save interview results: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	defaultDir := filepath.Join("docs", "implementation-plan", "m0")
	initCmd.Flags().StringVarP(&interactiveOutputDir, "output-dir", "o", defaultDir, "Directory to save output files")
	initCmd.Flags().BoolVar(&noInteractive, "no-interactive", false, "Skip interactive requirements gathering")
}

func saveInterviewResult(result *prompt.InterviewResult, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	transcriptPath := filepath.Join(outputDir, "requirements-interview.md")
	if err := os.WriteFile(transcriptPath, []byte(result.Transcript), 0644); err != nil {
		return fmt.Errorf("failed to write transcript: %w", err)
	}
	fmt.Printf("✓ Saved transcript to %s\n", transcriptPath)

	yamlPath := filepath.Join(outputDir, "requirements.md")
	yamlContent := generateYAML(result)
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		return fmt.Errorf("failed to write YAML: %w", err)
	}
	fmt.Printf("✓ Saved requirements to %s\n", yamlPath)

	fmt.Println("\n✓ Requirements gathering complete!")
	fmt.Printf("  Transcript: %s\n", transcriptPath)
	fmt.Printf("  Requirements: %s\n", yamlPath)

	return nil
}

func generateYAML(result *prompt.InterviewResult) string {
	return fmt.Sprintf(`# Milestone Requirements

## Milestone Information

**ID:** %s
**Title:** %s
**Description:** %s

## Stakeholder

%s

## Objectives

%s

## Success Metrics

%s

## Style Anchors

%s

## Scope

%s

## Out of Scope

%s

## Constraints & Assumptions

%s

## Deliverables

%s

## Timeline

%s

## Testing Requirements

%s

## Privacy & Security

%s
`,
		result.Responses["milestone_id"],
		result.Responses["milestone_title"],
		result.Responses["milestone_description"],
		result.Responses["stakeholder"],
		result.Responses["objectives"],
		result.Responses["success_metrics"],
		result.Responses["style_anchors"],
		result.Responses["scope"],
		result.Responses["out_of_scope"],
		result.Responses["constraints"],
		result.Responses["deliverables"],
		result.Responses["timeline"],
		result.Responses["testing"],
		result.Responses["privacy"],
	)
}
