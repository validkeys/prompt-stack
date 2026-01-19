package main

import (
	"fmt"
	"os"

	"github.com/kyledavis/prompt-stack/internal/validation"
	"github.com/spf13/cobra"
)

var (
	validateInput  string
	validateOutput string
	validateStrict bool
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate implementation plans",
	Long:  `Validate implementation plans against schema and quality standards.`,
	Run: func(cmd *cobra.Command, args []string) {
		if validateInput == "" {
			fmt.Println("Error: --input is required")
			cmd.Help()
			os.Exit(1)
		}

		config := validation.Config{
			InputPath:  validateInput,
			OutputPath: validateOutput,
			Strict:     validateStrict,
		}

		result, err := validation.Validate(config)
		if err != nil {
			fmt.Printf("Validation error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Validation result: %s (score: %.2f)\n", result.OverallResult, result.OverallScore)
		for name, score := range result.ComponentScores {
			fmt.Printf("  %s: %.2f\n", name, score.Score)
		}
		for _, issue := range result.Issues {
			fmt.Printf("    [%s] %s\n", issue.Severity, issue.Message)
		}

		if result.ApprovalStatus != "" {
			fmt.Printf("\nApproval status: %s\n", result.ApprovalStatus)
			if result.ApprovalReason != "" {
				fmt.Printf("Reason: %s\n", result.ApprovalReason)
			}
		}

		if result.OverallResult == "FAIL" {
			os.Exit(1)
		}
	},
}

func init() {
	validateCmd.Flags().StringVarP(&validateInput, "input", "i", "", "Input file to validate (required)")
	validateCmd.Flags().StringVarP(&validateOutput, "output", "o", ".prompt-stack/reports/validation_report.json", "Output report path")
	validateCmd.Flags().BoolVar(&validateStrict, "strict", false, "Fail validation on any issue")
	rootCmd.AddCommand(validateCmd)
}
