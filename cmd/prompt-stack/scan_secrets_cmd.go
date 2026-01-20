package main

import (
	"fmt"
	"os"

	"github.com/kyledavis/prompt-stack/internal/security"
	"github.com/spf13/cobra"
)

var scanSecretsCmd = &cobra.Command{
	Use:   "scan-secrets",
	Short: "Scan YAML files for embedded secrets",
	Long:  `Scans YAML files for potential embedded secrets such as API keys, tokens, passwords, and other sensitive data.`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")
		outputPath, _ := cmd.Flags().GetString("output")

		if filePath == "" {
			fmt.Fprintln(os.Stderr, "Error: --file is required")
			_ = cmd.Help()
			os.Exit(2)
		}

		exitCode, report, err := security.ScanSecrets(filePath, outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(exitCode)
		}

		fmt.Printf("Secrets scan: %s\n", report.ScanStatus)
		fmt.Printf("Secrets found: %d\n", report.SecretsFound)

		if len(report.Findings) > 0 {
			fmt.Printf("\nFindings:\n")
			for _, finding := range report.Findings {
				fmt.Printf("  - [%s] %s (line %d)\n", finding.Severity, finding.Type, finding.Line)
				fmt.Printf("    Context: %s\n", finding.Context)
				fmt.Printf("    %s\n", finding.Description)
			}

			fmt.Printf("\nSummary:\n")
			for severity, count := range report.Summary {
				fmt.Printf("  %s: %d\n", severity, count)
			}

			fmt.Printf("\nRecommendation: %s\n", report.Recommendation)
			os.Exit(exitCode)
		}

		fmt.Printf("Recommendation: %s\n", report.Recommendation)
		os.Exit(exitCode)
	},
}

func init() {
	rootCmd.AddCommand(scanSecretsCmd)
	scanSecretsCmd.Flags().String("file", "docs/implementation-plan/m0/ralphy_inputs.yaml", "Path to YAML file to scan")
	scanSecretsCmd.Flags().String("output", "", "Path to output JSON report (optional)")
}
