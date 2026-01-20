// schema_validator — Validates YAML files against JSON Schema specifications.
//
// # Purpose
//
// This package provides enterprise-grade schema validation for Ralphy YAML inputs used in the
// AI-assisted development workflow. It converts YAML files to JSON and validates them against
// a JSON Schema, ensuring that generated implementation plans conform to required structure,
// constraints, and best practices before execution.
//
// Features
//
//   - YAML to JSON conversion using sigs.k8s.io/yaml (handles complex YAML types)
//   - Strict JSON Schema validation via santhosh-tekuri/jsonschema/v5
//   - Support for $ref resolution (schema can reference external files)
//   - Detailed error reporting with path and keyword information
//   - File:// URL schema resolution for local schema files
//   - Absolute path computation for reliable file references
//
// Dependencies
//
//   - github.com/santhosh-tekuri/jsonschema/v5: JSON Schema Draft 7+ validator
//   - sigs.k8s.io/yaml: Robust YAML-to-JSON converter (Kubernetes project)
//
// Example
//
//	import "github.com/kyledavis/prompt-stack/internal/validation"
//
//	err := validation.ValidateYAMLAgainstSchema("docs/ralphy-inputs.schema.json", "final_ralphy_inputs.yaml")
//	if err != nil {
//	    log.Fatalf("Validation failed: %v", err)
//	}
package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"sigs.k8s.io/yaml"
)

// Exit codes for predictable script behavior
const (
	ExitSuccess   = 0 // Validation passed
	ExitFailed    = 1 // Validation failed (schema violations)
	ExitExecution = 2 // Execution error (file I/O, invalid args, etc.)
)

// ValidateYAML validates a YAML file against a JSON Schema and returns exit code.
//
// This is the main entrypoint for the validate-yaml command.
//
// Parameters:
//
//	schemaPath - Path to the JSON Schema file
//	yamlPath   - Path to the YAML file to validate
//
// Returns:
//
//	int - Exit code (0=success, 1=validation failed, 2=execution error)
//	error - Details about execution error (nil on validation failure)
func ValidateYAML(schemaPath, yamlPath string) (int, error) {
	schema, err := loadAndCompileSchema(schemaPath)
	if err != nil {
		return ExitExecution, err
	}

	document, err := loadAndConvertYAML(yamlPath)
	if err != nil {
		return ExitExecution, err
	}

	if err := validateAgainstSchema(schema, document); err != nil {
		return ExitFailed, err
	}

	return ExitSuccess, nil
}

// loadAndCompileSchema loads and compiles a JSON Schema from the given path.
//
// The schema path is converted to an absolute path and expressed as a file:// URL,
// which enables the JSON Schema compiler to resolve $ref references relative to the
// schema file location.
func loadAndCompileSchema(schemaPath string) (*jsonschema.Schema, error) {
	absSchema, err := filepath.Abs(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve schema path %q: %w", schemaPath, err)
	}

	schemaURL := "file://" + absSchema
	compiler := jsonschema.NewCompiler()

	schema, err := compiler.Compile(schemaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema %q: %w", schemaPath, err)
	}

	return schema, nil
}

// loadAndConvertYAML reads a YAML file, converts it to JSON, and unmarshals it.
//
// This function handles the three-step process of YAML validation:
// 1. Read the YAML file from disk
// 2. Convert YAML to JSON (sigs.k8s.io/yaml handles complex types)
// 3. Unmarshal JSON into a generic interface{} for schema validation
func loadAndConvertYAML(yamlPath string) (interface{}, error) {
	yamlBytes, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file %q: %w", yamlPath, err)
	}

	jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML to JSON: %w", err)
	}

	var doc interface{}
	if err := json.Unmarshal(jsonBytes, &doc); err != nil {
		return nil, fmt.Errorf("invalid JSON after conversion: %w", err)
	}

	return doc, nil
}

// validateAgainstSchema validates a document against a compiled JSON Schema.
//
// The validation uses the jsonschema library's Validate method, which
// provides detailed error information including the location in the document that
// failed validation and the specific schema constraint that was violated.
func validateAgainstSchema(schema *jsonschema.Schema, document interface{}) error {
	if err := schema.Validate(document); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

// ValidateYAMLAgainstSchema validates a YAML file against a JSON Schema.
//
// Parameters:
//
//	schemaPath - Path to the JSON Schema file
//	yamlPath   - Path to the YAML file to validate
//
// Returns:
//
//	error - Nil if validation succeeds, detailed error if validation fails
//
// Example:
//
//	err := ValidateYAMLAgainstSchema("docs/ralphy-inputs.schema.json", "final_ralphy_inputs.yaml")
//	if err != nil {
//	    fmt.Printf("Validation failed: %v\n", err)
//	}
func ValidateYAMLAgainstSchema(schemaPath, yamlPath string) error {
	schema, err := loadAndCompileSchema(schemaPath)
	if err != nil {
		return err
	}

	document, err := loadAndConvertYAML(yamlPath)
	if err != nil {
		return err
	}

	return validateAgainstSchema(schema, document)
}
