package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kyledavis/prompt-stack/internal/validation/enforcement"
	"github.com/spf13/cobra"
)

var validateEnforcementCmd = &cobra.Command{
	Use:   "validate-enforcement",
	Short: "Validate multi-layer enforcement and commit/scope policies",
	Long:  `Validates that Ralphy YAML files include comprehensive multi-layer enforcement (prompt-level, IDE, pre-commit, CI, runtime) and commit/scope policies.`,
	Run: func(cmd *cobra.Command, args []string) {
		yamlPath, _ := cmd.Flags().GetString("file")

		if yamlPath == "" {
			fmt.Fprintln(os.Stderr, "Error: --file is required")
			_ = cmd.Help()
			os.Exit(2)
		}

		exitCode, result, err := enforcement.ValidateEnforcementFromFile(yamlPath)
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
	rootCmd.AddCommand(validateEnforcementCmd)
	validateEnforcementCmd.Flags().String("file", "final_ralphy_inputs.yaml", "Path to YAML file to validate")
}
