package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kyledavis/prompt-stack/internal/validation/quality"
	"github.com/spf13/cobra"
)

var validateQualityCmd = &cobra.Command{
	Use:   "validate-quality",
	Short: "Generate comprehensive quality report",
	Long:  `Generates a comprehensive quality report by aggregating all validation reports from a directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		reportsDir, _ := cmd.Flags().GetString("reports-dir")

		if reportsDir == "" {
			fmt.Fprintln(os.Stderr, "Error: --reports-dir is required")
			_ = cmd.Help()
			os.Exit(2)
		}

		report, err := quality.GenerateQualityReport(reportsDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}

		jsonResult, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to marshal report: %v\n", err)
			os.Exit(2)
		}

		fmt.Println(string(jsonResult))
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(validateQualityCmd)
	validateQualityCmd.Flags().String("reports-dir", "", "Path to directory containing validation reports")
}
