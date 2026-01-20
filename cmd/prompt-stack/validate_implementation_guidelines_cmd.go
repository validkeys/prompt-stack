package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kyledavis/prompt-stack/internal/validation/implementationguidelines"
	"github.com/spf13/cobra"
)

var validateImplementationGuidelinesCmd = &cobra.Command{
	Use:   "validate-implementation-guidelines",
	Short: "Validate implementation guidelines inclusion",
	Long:  `Validates that implementation-phase guidelines are present where applicable, including test-first workflow identification and testable acceptance criteria.`,
	Run: func(cmd *cobra.Command, args []string) {
		yamlPath, _ := cmd.Flags().GetString("file")

		if yamlPath == "" {
			fmt.Fprintln(os.Stderr, "Error: --file is required")
			_ = cmd.Help()
			os.Exit(2)
		}

		exitCode, result, err := implementationguidelines.ValidateImplementationGuidelinesFromFile(yamlPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(exitCode)
		}

		jsonResult, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to marshal result: %v\n", err)
			os.Exit(2)
		}

		fmt.Println(string(jsonResult))
		os.Exit(exitCode)
	},
}

func init() {
	rootCmd.AddCommand(validateImplementationGuidelinesCmd)
	validateImplementationGuidelinesCmd.Flags().String("file", "final_ralphy_inputs.yaml", "Path to YAML file to validate")
}
