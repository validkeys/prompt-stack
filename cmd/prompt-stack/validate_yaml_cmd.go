package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kyledavis/prompt-stack/internal/validation"
	"github.com/spf13/cobra"
)

var validateYamlCmd = &cobra.Command{
	Use:   "validate-yaml",
	Short: "Validate YAML files against JSON Schema",
	Long:  `Validates YAML files against a JSON Schema specification, ensuring generated implementation plans conform to required structure.`,
	Run: func(cmd *cobra.Command, args []string) {
		schemaPath, _ := cmd.Flags().GetString("schema")
		yamlPath, _ := cmd.Flags().GetString("file")

		if schemaPath == "" || yamlPath == "" {
			fmt.Fprintln(os.Stderr, "Error: both --schema and --file are required")
			_ = cmd.Help()
			os.Exit(2)
		}

		exitCode, err := validation.ValidateYAML(schemaPath, yamlPath)
		report := map[string]interface{}{
			"summary": map[string]interface{}{
				"validation_type": "yaml_syntax",
				"status":          "passed",
				"validator_tool":  "prompt-stack",
				"file_validated":  yamlPath,
				"schema_used":     schemaPath,
				"validation_details": map[string]string{
					"yaml_conversion":   "success",
					"json_conversion":   "success",
					"schema_validation": "success",
				},
				"exit_code": float64(exitCode),
				"message":   "Validation passed",
			},
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(exitCode)
		}

		jsonReport, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(jsonReport))
		os.Exit(exitCode)
	},
}

func init() {
	rootCmd.AddCommand(validateYamlCmd)
	validateYamlCmd.Flags().String("schema", "docs/ralphy-inputs.schema.json", "Path to JSON Schema file")
	validateYamlCmd.Flags().String("file", "final_ralphy_inputs.yaml", "Path to YAML file to validate")
}
