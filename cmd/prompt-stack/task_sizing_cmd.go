package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kyledavis/prompt-stack/internal/build"
	"github.com/spf13/cobra"
)

var validateTaskSizingCmd = &cobra.Command{
	Use:   "validate-task-sizing",
	Short: "Validate task sizing compliance",
	Long:  `Validates that tasks in Ralphy YAML files comply with task sizing guidelines (30-150 minute range).`,
	Run: func(cmd *cobra.Command, args []string) {
		yamlPath, _ := cmd.Flags().GetString("file")

		if yamlPath == "" {
			fmt.Fprintln(os.Stderr, "Error: --file is required")
			_ = cmd.Help()
			os.Exit(2)
		}

		exitCode, result, err := build.ValidateTaskSizing(yamlPath)
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
	rootCmd.AddCommand(validateTaskSizingCmd)
	validateTaskSizingCmd.Flags().String("file", "final_ralphy_inputs.yaml", "Path to YAML file to validate")
}
