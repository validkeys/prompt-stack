package database

import (
	"testing"
)

func TestValidateNoSecrets(t *testing.T) {
	tests := []struct {
		name      string
		tableName string
		data      map[string]interface{}
		wantErr   bool
	}{
		{
			name:      "clean data",
			tableName: "patterns",
			data: map[string]interface{}{
				"description": "This is a clean description",
				"source":      "github.com/user/repo",
			},
			wantErr: false,
		},
		{
			name:      "contains API key",
			tableName: "patterns",
			data: map[string]interface{}{
				"description": "Use api_key = my_secret_api_key_123456",
				"source":      "github.com/user/repo",
			},
			wantErr: true,
		},
		{
			name:      "contains password",
			tableName: "requirements",
			data: map[string]interface{}{
				"content": "The password is mypassword1234",
			},
			wantErr: true,
		},
		{
			name:      "contains AWS key",
			tableName: "tasks",
			data: map[string]interface{}{
				"description": "AKIAIOSFODNN7EXAMPLE",
			},
			wantErr: true,
		},
		{
			name:      "non-sensitive column",
			tableName: "patterns",
			data: map[string]interface{}{
				"name": "AKIAIOSFODNN7EXAMPLE",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNoSecrets(tt.tableName, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNoSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePatternData(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid pattern data",
			data: map[string]interface{}{
				"name":        "Test Pattern",
				"description": "This is a test pattern",
				"source":      "github.com/user/repo",
			},
			wantErr: false,
		},
		{
			name: "pattern with secret",
			data: map[string]interface{}{
				"name":        "Test Pattern",
				"description": "Use token = ghp_1234567890abcdefghijklmnopqrstuvwxyz",
				"source":      "github.com/user/repo",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePatternData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePatternData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRequirementData(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid requirement data",
			data: map[string]interface{}{
				"milestone_id":     "m1",
				"requirement_type": "functional",
				"content":          "The system should support feature X",
			},
			wantErr: false,
		},
		{
			name: "requirement with secret",
			data: map[string]interface{}{
				"milestone_id":     "m1",
				"requirement_type": "functional",
				"content":          "Use secret: very_secret_value_123456789",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequirementData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequirementData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTaskData(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid task data",
			data: map[string]interface{}{
				"milestone_id": "m1",
				"task_id":      "m1-001",
				"title":        "Implement feature",
				"description":  "Implement the feature using clean code",
			},
			wantErr: false,
		},
		{
			name: "task with secret",
			data: map[string]interface{}{
				"milestone_id": "m1",
				"task_id":      "m1-001",
				"title":        "Implement feature",
				"description":  "Use password: mypassword1234",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTaskData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTaskData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateValidationReportData(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid report data",
			data: map[string]interface{}{
				"plan_path":       "docs/plan.yaml",
				"overall_score":   0.95,
				"approval_status": "APPROVED",
				"report_json":     `{"score": 0.95}`,
			},
			wantErr: false,
		},
		{
			name: "report with secret in JSON",
			data: map[string]interface{}{
				"plan_path":       "docs/plan.yaml",
				"overall_score":   0.95,
				"approval_status": "APPROVED",
				"report_json":     `{"api_key": "my_secret_api_key_123456"}`,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateValidationReportData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateValidationReportData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
