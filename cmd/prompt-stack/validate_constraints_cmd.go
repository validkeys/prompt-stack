package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kyledavis/prompt-stack/internal/validation/constraints"
	"github.com/spf13/cobra"
)

var validateConstraintsCmd = &cobra.Command{
	Use:   "validate-constraints",
	Short: "Validate affirmative constraints usage",
	Long:  `Validates that constraints in Ralphy YAML files use affirmative language (e.g., "do this" instead of "don't do that").`,
	Run: func(cmd *cobra.Command, args []string) {
		yamlPath, _ := cmd.Flags().GetString("file")

		if yamlPath == "" {
			fmt.Fprintln(os.Stderr, "Error: --file is required")
			_ = cmd.Help()
			os.Exit(2)
		}

		exitCode, result, err := constraints.ValidateConstraintsFromFile(yamlPath)
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
	rootCmd.AddCommand(validateConstraintsCmd)
	validateConstraintsCmd.Flags().String("file", "final_ralphy_inputs.yaml", "Path to YAML file to validate")
}
