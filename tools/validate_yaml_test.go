package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestYAMLValidationSuccess(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		schema  string
		wantErr bool
	}{
		{
			name:    "Valid simple YAML",
			yaml:    "test_data/valid_simple.yaml",
			schema:  "../docs/ralphy-inputs.schema.json",
			wantErr: false,
		},
		{
			name:    "Invalid YAML - missing required field",
			yaml:    "test_data/invalid_missing_required.yaml",
			schema:  "../docs/ralphy-inputs.schema.json",
			wantErr: true,
		},
		{
			name:    "Invalid YAML - wrong type",
			yaml:    "test_data/invalid_wrong_type.yaml",
			schema:  "../docs/ralphy-inputs.schema.json",
			wantErr: true,
		},
		{
			name:    "Valid full Ralphy inputs",
			yaml:    "../docs/implementation-plan/m0/ralphy_inputs.yaml",
			schema:  "../docs/ralphy-inputs.schema.json",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, err := loadAndCompileSchema(tt.schema)
			if err != nil {
				t.Fatalf("Failed to load schema: %v", err)
			}

			document, err := loadAndConvertYAML(tt.yaml)
			if err != nil {
				t.Fatalf("Failed to load YAML: %v", err)
			}

			validationErr := validateAgainstSchema(schema, document)
			if (validationErr != nil) != tt.wantErr {
				t.Errorf("validateAgainstSchema() error = %v, wantErr %v", validationErr, tt.wantErr)
			}
		})
	}
}

func TestLoadAndCompileSchema(t *testing.T) {
	tests := []struct {
		name      string
		schema    string
		wantErr   bool
		errString string
	}{
		{
			name:    "Valid schema",
			schema:  "../docs/ralphy-inputs.schema.json",
			wantErr: false,
		},
		{
			name:      "Non-existent schema",
			schema:    "../nonexistent/schema.json",
			wantErr:   true,
			errString: "failed to compile schema",
		},
		{
			name:      "Invalid schema JSON",
			schema:    "test_data/invalid_schema.json",
			wantErr:   true,
			errString: "failed to compile schema",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadAndCompileSchema(tt.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadAndCompileSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errString != "" {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.errString)
					return
				}
				if !contains(err.Error(), tt.errString) {
					t.Errorf("Expected error containing %q, got %q", tt.errString, err.Error())
				}
			}
		})
	}
}

func TestLoadAndConvertYAML(t *testing.T) {
	tests := []struct {
		name      string
		yaml      string
		wantErr   bool
		errString string
	}{
		{
			name:    "Valid YAML",
			yaml:    "../docs/implementation-plan/m0/ralphy_inputs.yaml",
			wantErr: false,
		},
		{
			name:      "Non-existent file",
			yaml:      "../nonexistent/file.yaml",
			wantErr:   true,
			errString: "failed to read YAML file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadAndConvertYAML(tt.yaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadAndConvertYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errString != "" {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.errString)
					return
				}
				if !contains(err.Error(), tt.errString) {
					t.Errorf("Expected error containing %q, got %q", tt.errString, err.Error())
				}
			}
		})
	}
}

func TestValidateAgainstSchema(t *testing.T) {
	validSchemaPath := "../docs/ralphy-inputs.schema.json"
	schema, err := loadAndCompileSchema(validSchemaPath)
	if err != nil {
		t.Fatalf("Failed to load test schema: %v", err)
	}

	tests := []struct {
		name    string
		doc     interface{}
		wantErr bool
	}{
		{
			name: "Valid minimal document",
			doc: map[string]interface{}{
				"name":        "test",
				"version":     "1.0.0",
				"rules_file":  "rules.md",
				"task_sizing": map[string]interface{}{"min_minutes": float64(30), "max_minutes": float64(150)},
				"tdd":         map[string]interface{}{"required": false, "test_command": "go test"},
				"model_preferences": map[string]interface{}{
					"primary": "opencode",
				},
				"outputs": map[string]interface{}{
					"allowed_file_edits":    []interface{}{"cmd/**"},
					"disallowed_file_edits": []interface{}{".github/**"},
				},
				"tasks": []interface{}{},
			},
			wantErr: false,
		},
		{
			name: "Missing required field",
			doc: map[string]interface{}{
				"name":        "test",
				"version":     "1.0.0",
				"rules_file":  "rules.md",
				"task_sizing": map[string]interface{}{"min_minutes": float64(30), "max_minutes": float64(150)},
				"tdd":         map[string]interface{}{"required": false, "test_command": "go test"},
				"model_preferences": map[string]interface{}{
					"primary": "opencode",
				},
				"outputs": map[string]interface{}{
					"allowed_file_edits":    []interface{}{"cmd/**"},
					"disallowed_file_edits": []interface{}{".github/**"},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid version format",
			doc: map[string]interface{}{
				"name":        "test",
				"version":     "v1.0.0",
				"rules_file":  "rules.md",
				"task_sizing": map[string]interface{}{"min_minutes": float64(30), "max_minutes": float64(150)},
				"tdd":         map[string]interface{}{"required": false, "test_command": "go test"},
				"model_preferences": map[string]interface{}{
					"primary": "opencode",
				},
				"outputs": map[string]interface{}{
					"allowed_file_edits":    []interface{}{"cmd/**"},
					"disallowed_file_edits": []interface{}{".github/**"},
				},
				"tasks": []interface{}{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAgainstSchema(schema, tt.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAgainstSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExitCodes(t *testing.T) {
	tests := []struct {
		name string
		code int
		desc string
	}{
		{
			name: "ExitSuccess",
			code: ExitSuccess,
			desc: "Validation passed",
		},
		{
			name: "ExitFailed",
			code: ExitFailed,
			desc: "Validation failed",
		},
		{
			name: "ExitExecution",
			code: ExitExecution,
			desc: "Execution error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code < 0 || tt.code > 2 {
				t.Errorf("Exit code %d out of range [0-2]", tt.code)
			}
		})
	}
}

func TestValidationReportStructure(t *testing.T) {
	reportPath := "../docs/implementation-plan/m0/yaml_validation_report.json"

	data, err := ioutil.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("Failed to read validation report: %v", err)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse validation report: %v", err)
	}

	requiredFields := []string{
		"validation_type",
		"status",
		"timestamp",
		"validator_tool",
		"file_validated",
		"schema_used",
		"validation_details",
		"exit_code",
		"message",
	}

	for _, field := range requiredFields {
		if _, ok := report[field]; !ok {
			t.Errorf("Validation report missing required field: %s", field)
		}
	}

	if report["validation_type"] != "yaml_syntax" {
		t.Errorf("Expected validation_type 'yaml_syntax', got %v", report["validation_type"])
	}

	if report["exit_code"].(float64) != 0 {
		t.Errorf("Expected exit_code 0, got %v", report["exit_code"])
	}
}

func TestYAMLValidationReportContent(t *testing.T) {
	reportPath := "../docs/implementation-plan/m0/yaml_validation_report.json"

	data, err := ioutil.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("Failed to read validation report: %v", err)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to parse validation report: %v", err)
	}

	details := report["validation_details"].(map[string]interface{})

	if details["yaml_conversion"] != "success" {
		t.Errorf("Expected yaml_conversion 'success', got %v", details["yaml_conversion"])
	}

	if details["json_conversion"] != "success" {
		t.Errorf("Expected json_conversion 'success', got %v", details["json_conversion"])
	}

	if details["schema_validation"] != "success" {
		t.Errorf("Expected schema_validation 'success', got %v", details["schema_validation"])
	}

	fixes := details["fixes_applied"].([]interface{})
	if len(fixes) == 0 {
		t.Error("Expected at least one fix to be documented")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
