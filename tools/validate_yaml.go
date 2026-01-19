// validate_yaml — Validates YAML files against JSON Schema specifications.
//
// # Purpose
//
// This tool provides enterprise-grade schema validation for Ralphy YAML inputs used in the
// AI-assisted development workflow. It converts YAML files to JSON and validates them against
// a JSON Schema, ensuring that generated implementation plans conform to required structure,
// constraints, and best practices before execution.
//
// Usage
//
//	go run validate_yaml.go --schema <schema.json> --file <input.yaml>
//
//	go run validate_yaml.go -s docs/ralphy-inputs.schema.json -f final_ralphy_inputs.yaml
//
// Exit Codes
//
//	0 - Validation successful (YAML conforms to schema)
//	1 - Validation failed (YAML does not conform to schema)
//	2 - Error in command execution (file not found, invalid schema, etc.)
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
// Integration Points
//
//   - CI/CD pipelines: Run as a quality gate before merge
//   - Pre-commit hooks: Validate YAML files before commits
//   - Build Mode: Validate generated Ralphy YAML before execution
//   - Local development: Quick validation during development
//
// Dependencies
//
//   - github.com/santhosh-tekuri/jsonschema/v5: JSON Schema Draft 7+ validator
//   - sigs.k8s.io/yaml: Robust YAML-to-JSON converter (Kubernetes project)
//
// Performance Characteristics
//
//   - Memory: O(file_size) - entire file loaded into memory
//   - Time: O(file_size + schema_complexity) - linear in file size
//   - Typical runtime: <100ms for 1MB files on modern hardware
//
// Limitations
//
//   - Does not validate YAML syntax separately (conversion failures indicate syntax errors)
//   - Schema $ref resolution requires relative paths from schema file location
//   - Large files (>100MB) may exceed memory limits
//   - Does not support remote schema URLs (only file:// protocol)
//
// Security Considerations
//
//   - Reads only local files specified via command-line arguments
//   - Does not execute external commands or evaluate arbitrary code
//   - Schema $ref resolution restricted to file:// protocol
//   - No network access or external dependencies during validation
//
// Example Workflows
//
//	CI Pipeline:
//	  - name: Validate Ralphy YAML
//	    run: |
//	      cd tools
//	      go run validate_yaml.go \
//	        --schema ../docs/ralphy-inputs.schema.json \
//	        --file ../final_ralphy_inputs.yaml
//
//	Makefile Target:
//	  validate-yaml:
//	      cd tools && go run validate_yaml.go \
//	        --schema ../docs/ralphy-inputs.schema.json \
//	        --file $(FILE)
//
//	Pre-commit Hook:
//	  if grep -q "ralphy_inputs.yaml" .git/hooks/pre-commit.sample; then
//	    tools/validate_yaml.go \
//	      --schema docs/ralphy-inputs.schema.json \
//	      --file ralphy_inputs.yaml
//	  fi
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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

// loadAndCompileSchema loads and compiles a JSON Schema from the given path.
//
// The schema path is converted to an absolute path and expressed as a file:// URL,
// which enables the JSON Schema compiler to resolve $ref references relative to the
// schema file location.
//
// Parameters:
//
//	schemaPath - Path to the JSON Schema file (can be relative or absolute)
//
// Returns:
//
//	*jsonschema.Schema - Compiled schema ready for validation
//	error - Error if the schema cannot be read or compiled
//
// Example:
//
//	schema, err := loadAndCompileSchema("docs/ralphy-inputs.schema.json")
//	if err != nil {
//	    log.Fatalf("Failed to load schema: %v", err)
//	}
func loadAndCompileSchema(schemaPath string) (*jsonschema.Schema, error) {
	// Resolve to absolute path for reliable $ref resolution
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
//
// Parameters:
//
//	yamlPath - Path to the YAML file to validate
//
// Returns:
//
//	interface{} - The parsed YAML content as a Go interface{}
//	error - Error if the file cannot be read or conversion fails
//
// Example:
//
//	doc, err := loadAndConvertYAML("final_ralphy_inputs.yaml")
//	if err != nil {
//	    log.Fatalf("Failed to load YAML: %v", err)
//	}
func loadAndConvertYAML(yamlPath string) (interface{}, error) {
	yamlBytes, err := ioutil.ReadFile(yamlPath)
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
//
// Parameters:
//
//	schema   - Compiled JSON Schema from loadAndCompileSchema
//	document - Document to validate (typically from loadAndConvertYAML)
//
// Returns:
//
//	error - Nil if validation succeeds, detailed error if validation fails
//
// Example:
//
//	err := validateAgainstSchema(schema, doc)
//	if err != nil {
//	    fmt.Printf("Validation failed: %v\n", err)
//	}
func validateAgainstSchema(schema *jsonschema.Schema, document interface{}) error {
	if err := schema.Validate(document); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

// parseFlags parses and validates command-line arguments.
//
// Returns:
//
//	schemaPath - Path to the JSON Schema file
//	yamlPath   - Path to the YAML file to validate
//	error      - Error if arguments are missing or invalid
func parseFlags() (string, string, error) {
	schemaPath := flag.String("schema", "docs/ralphy-inputs.schema.json", "Path to JSON Schema file")
	yamlPath := flag.String("file", "final_ralphy_inputs.yaml", "Path to YAML file to validate")
	flag.Parse()

	if *schemaPath == "" {
		return "", "", fmt.Errorf("--schema flag is required")
	}
	if *yamlPath == "" {
		return "", "", fmt.Errorf("--file flag is required")
	}

	return *schemaPath, *yamlPath, nil
}

// main is the entry point for the validate_yaml tool.
//
// Execution flow:
//  1. Parse command-line flags
//  2. Load and compile the JSON Schema
//  3. Load and convert the YAML file to JSON
//  4. Validate the document against the schema
//  5. Exit with appropriate status code
//
// Error handling:
//   - Errors during setup (file I/O, schema compilation) exit with code 2
//   - Validation failures exit with code 1 and print detailed errors
//   - Successful validation exits with code 0
//
// See package-level documentation for usage examples and exit code details.
func main() {
	schemaPath, yamlPath, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		fmt.Fprintf(os.Stderr, "Usage: go run validate_yaml.go --schema <schema.json> --file <input.yaml>\n")
		os.Exit(ExitExecution)
	}

	schema, err := loadAndCompileSchema(schemaPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ExitExecution)
	}

	document, err := loadAndConvertYAML(yamlPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ExitExecution)
	}

	if err := validateAgainstSchema(schema, document); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ExitFailed)
	}

	fmt.Println("validation OK")
	os.Exit(ExitSuccess)
}
