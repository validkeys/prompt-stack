package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCheckVerificationLayers tests the verification layers detection
func TestCheckVerificationLayers(t *testing.T) {
	tests := []struct {
		name     string
		config   *RalphyYAML
		expected VerificationLayers
	}{
		{
			name: "complete verification layers",
			config: &RalphyYAML{
				RulesFile: "docs/opencode-rules.md",
				ValidationSchemas: []string{
					"docs/ralphy-inputs.schema.json",
				},
				GlobalConstraints: GlobalConstraints{
					ForbiddenPatterns: []PatternConstraint{
						{Pattern: "\\bany\\b", Message: "Use unknown with type guards instead"},
					},
					RequiredPatterns: []PatternConstraint{
						{Pattern: "import.*zod", When: "validating_external_data"},
					},
					AffirmativeConstraints: []string{
						"Use affirmative language",
					},
				},
				CI: CI{
					Precommit: []string{
						"go test ./...",
						"go vet ./...",
					},
					CIChecks: []string{
						"build-and-test",
					},
				},
				DriftPolicyRef: "docs/drift-policy.md",
			},
			expected: VerificationLayers{
				PromptLevel:    true,
				IDEIntegration: true,
				PreCommit:      true,
				CIChecks:       true,
				Runtime:        true,
				TotalLayers:    5,
			},
		},
		{
			name: "minimal verification layers",
			config: &RalphyYAML{
				RulesFile: "",
				GlobalConstraints: GlobalConstraints{
					ForbiddenPatterns: []PatternConstraint{
						{Pattern: "\\bany\\b", Message: "Use unknown with type guards instead"},
					},
				},
				CI: CI{
					Precommit: []string{
						"go test ./...",
					},
				},
			},
			expected: VerificationLayers{
				PromptLevel:    true,
				IDEIntegration: false,
				PreCommit:      true,
				CIChecks:       false,
				Runtime:        false,
				TotalLayers:    2,
			},
		},
		{
			name: "no verification layers",
			config: &RalphyYAML{
				RulesFile:         "",
				GlobalConstraints: GlobalConstraints{},
				CI:                CI{},
			},
			expected: VerificationLayers{
				PromptLevel:    false,
				IDEIntegration: false,
				PreCommit:      false,
				CIChecks:       false,
				Runtime:        false,
				TotalLayers:    0,
			},
		},
		{
			name: "only prompt-level constraints",
			config: &RalphyYAML{
				GlobalConstraints: GlobalConstraints{
					ForbiddenPatterns: []PatternConstraint{
						{Pattern: "@ts-ignore", Message: "Fix type errors properly"},
					},
				},
			},
			expected: VerificationLayers{
				PromptLevel:    true,
				IDEIntegration: false,
				PreCommit:      false,
				CIChecks:       false,
				Runtime:        false,
				TotalLayers:    1,
			},
		},
		{
			name: "only IDE integration",
			config: &RalphyYAML{
				RulesFile: "docs/opencode-rules.md",
			},
			expected: VerificationLayers{
				PromptLevel:    false,
				IDEIntegration: true,
				PreCommit:      false,
				CIChecks:       false,
				Runtime:        false,
				TotalLayers:    1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkVerificationLayers(tt.config)
			if result != tt.expected {
				t.Errorf("checkVerificationLayers() = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}

// TestCheckCommitPolicy tests the commit policy validation
func TestCheckCommitPolicy(t *testing.T) {
	tests := []struct {
		name     string
		config   *RalphyYAML
		expected CommitPolicyStatus
	}{
		{
			name: "complete commit policy",
			config: &RalphyYAML{
				Outputs: Outputs{
					CommitPolicy: CommitPolicy{
						PrefixRules: []string{
							"feat:",
							"fix:",
							"docs:",
						},
						RequireScope:               true,
						RequireConventionalCommits: true,
					},
				},
			},
			expected: CommitPolicyStatus{
				HasPrefixRules:         true,
				HasScopeRequirement:    true,
				HasConventionalCommits: true,
				Complete:               true,
			},
		},
		{
			name: "minimal commit policy",
			config: &RalphyYAML{
				Outputs: Outputs{
					CommitPolicy: CommitPolicy{
						PrefixRules: []string{
							"feat:",
						},
					},
				},
			},
			expected: CommitPolicyStatus{
				HasPrefixRules:         true,
				HasScopeRequirement:    false,
				HasConventionalCommits: false,
				Complete:               true,
			},
		},
		{
			name: "no commit policy",
			config: &RalphyYAML{
				Outputs: Outputs{},
			},
			expected: CommitPolicyStatus{
				HasPrefixRules:         false,
				HasScopeRequirement:    false,
				HasConventionalCommits: false,
				Complete:               false,
			},
		},
		{
			name: "empty prefix rules",
			config: &RalphyYAML{
				Outputs: Outputs{
					CommitPolicy: CommitPolicy{
						PrefixRules: []string{},
					},
				},
			},
			expected: CommitPolicyStatus{
				HasPrefixRules:         false,
				HasScopeRequirement:    false,
				HasConventionalCommits: false,
				Complete:               false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkCommitPolicy(tt.config)
			if result != tt.expected {
				t.Errorf("checkCommitPolicy() = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}

// TestCheckScopeEnforcement tests the scope enforcement validation
func TestCheckScopeEnforcement(t *testing.T) {
	tests := []struct {
		name     string
		config   *RalphyYAML
		expected ScopeEnforcement
	}{
		{
			name: "complete scope enforcement",
			config: &RalphyYAML{
				Outputs: Outputs{
					AllowedFileEdits: []string{
						"docs/**",
						"src/**/*.ts",
					},
					DisallowedFileEdits: []string{
						".github/**",
						"scripts/**",
					},
				},
			},
			expected: ScopeEnforcement{
				HasAllowedFileEdits:      true,
				HasDisallowedFileEdits:   true,
				AllTasksHaveFilesInScope: true,
				Complete:                 true,
			},
		},
		{
			name: "only allowed file edits",
			config: &RalphyYAML{
				Outputs: Outputs{
					AllowedFileEdits: []string{
						"docs/**",
					},
				},
			},
			expected: ScopeEnforcement{
				HasAllowedFileEdits:      true,
				HasDisallowedFileEdits:   false,
				AllTasksHaveFilesInScope: true,
				Complete:                 false,
			},
		},
		{
			name: "only disallowed file edits",
			config: &RalphyYAML{
				Outputs: Outputs{
					DisallowedFileEdits: []string{
						".github/**",
					},
				},
			},
			expected: ScopeEnforcement{
				HasAllowedFileEdits:      false,
				HasDisallowedFileEdits:   true,
				AllTasksHaveFilesInScope: true,
				Complete:                 false,
			},
		},
		{
			name: "no scope enforcement",
			config: &RalphyYAML{
				Outputs: Outputs{},
			},
			expected: ScopeEnforcement{
				HasAllowedFileEdits:      false,
				HasDisallowedFileEdits:   false,
				AllTasksHaveFilesInScope: true,
				Complete:                 false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkScopeEnforcement(tt.config)
			if result != tt.expected {
				t.Errorf("checkScopeEnforcement() = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}

// TestCheckTasks tests task-level validation
func TestCheckTasks(t *testing.T) {
	tests := []struct {
		name           string
		config         *RalphyYAML
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "all tasks complete",
			config: &RalphyYAML{
				Tasks: []Task{
					{
						ID:                   "task-001",
						Title:                "Test task",
						Description:          "A test task",
						FilesInScope:         []string{"file1.go", "file2.go"},
						SingleResponsibility: "Test implementation",
						Verification: Verification{
							PreCommit: []string{"go test"},
						},
					},
					{
						ID:                   "task-002",
						Title:                "Another task",
						Description:          "Another test task",
						FilesInScope:         []string{"file3.go"},
						SingleResponsibility: "Documentation update",
						Verification: Verification{
							PreCommit: []string{"markdownlint"},
						},
					},
				},
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "tasks missing files_in_scope",
			config: &RalphyYAML{
				Tasks: []Task{
					{
						ID:                   "task-001",
						Title:                "Test task",
						Description:          "A test task",
						FilesInScope:         []string{},
						SingleResponsibility: "Test implementation",
						Verification: Verification{
							PreCommit: []string{"go test"},
						},
					},
					{
						ID:                   "task-002",
						Title:                "Another task",
						Description:          "Another test task",
						FilesInScope:         []string{"file3.go"},
						SingleResponsibility: "Documentation update",
						Verification: Verification{
							PreCommit: []string{"markdownlint"},
						},
					},
				},
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "tasks missing verification",
			config: &RalphyYAML{
				Tasks: []Task{
					{
						ID:                   "task-001",
						Title:                "Test task",
						Description:          "A test task",
						FilesInScope:         []string{"file1.go"},
						SingleResponsibility: "Test implementation",
						Verification:         Verification{},
					},
				},
			},
			expectedValid:  true, // Missing verification is a warning, not an error
			expectedErrors: 1,
		},
		{
			name: "tasks missing single_responsibility",
			config: &RalphyYAML{
				Tasks: []Task{
					{
						ID:                   "task-001",
						Title:                "Test task",
						Description:          "A test task",
						FilesInScope:         []string{"file1.go"},
						SingleResponsibility: "",
						Verification: Verification{
							PreCommit: []string{"go test"},
						},
					},
				},
			},
			expectedValid:  true, // Missing single responsibility is a warning, not an error
			expectedErrors: 1,
		},
		{
			name: "multiple task violations",
			config: &RalphyYAML{
				Tasks: []Task{
					{
						ID:                   "task-001",
						Title:                "Test task",
						Description:          "A test task",
						FilesInScope:         []string{},
						SingleResponsibility: "",
						Verification:         Verification{},
					},
				},
			},
			expectedValid:  false,
			expectedErrors: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidationResult{
				Valid:              true,
				TotalTasks:         len(tt.config.Tasks),
				Violations:         []Violation{},
				Recommendations:    []string{},
				VerificationLayers: VerificationLayers{},
				CommitPolicy:       CommitPolicyStatus{},
				ScopeEnforcement:   ScopeEnforcement{AllTasksHaveFilesInScope: true},
			}

			result = checkTasks(tt.config, result)

			if result.Valid != tt.expectedValid {
				t.Errorf("checkTasks().Valid = %v, expected %v", result.Valid, tt.expectedValid)
			}

			if len(result.Violations) != tt.expectedErrors {
				t.Errorf("checkTasks() violations = %d, expected %d", len(result.Violations), tt.expectedErrors)
			}
		})
	}
}

// TestValidateEnforcement tests the complete validation function
func TestValidateEnforcement(t *testing.T) {
	tests := []struct {
		name     string
		config   *RalphyYAML
		expected bool
	}{
		{
			name: "complete valid configuration",
			config: &RalphyYAML{
				RulesFile: "docs/opencode-rules.md",
				ValidationSchemas: []string{
					"docs/ralphy-inputs.schema.json",
				},
				GlobalConstraints: GlobalConstraints{
					ForbiddenPatterns: []PatternConstraint{
						{Pattern: "\\bany\\b", Message: "Use unknown with type guards instead"},
					},
					RequiredPatterns: []PatternConstraint{
						{Pattern: "import.*zod", When: "validating_external_data"},
					},
					AffirmativeConstraints: []string{
						"Use affirmative language",
					},
				},
				CI: CI{
					Precommit: []string{
						"go test ./...",
						"go vet ./...",
					},
					CIChecks: []string{
						"build-and-test",
					},
				},
				DriftPolicyRef: "docs/drift-policy.md",
				Outputs: Outputs{
					AllowedFileEdits: []string{
						"docs/**",
						"src/**/*.ts",
					},
					DisallowedFileEdits: []string{
						".github/**",
						"scripts/**",
					},
					CommitPolicy: CommitPolicy{
						PrefixRules: []string{
							"feat:",
							"fix:",
							"docs:",
						},
					},
				},
				Tasks: []Task{
					{
						ID:                   "task-001",
						Title:                "Test task",
						Description:          "A test task",
						FilesInScope:         []string{"file1.go", "file2.go"},
						SingleResponsibility: "Test implementation",
						Verification: Verification{
							PreCommit: []string{"go test"},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "insufficient verification layers",
			config: &RalphyYAML{
				GlobalConstraints: GlobalConstraints{
					ForbiddenPatterns: []PatternConstraint{
						{Pattern: "\\bany\\b", Message: "Use unknown with type guards instead"},
					},
				},
				Outputs: Outputs{
					AllowedFileEdits: []string{
						"docs/**",
					},
					DisallowedFileEdits: []string{
						".github/**",
					},
					CommitPolicy: CommitPolicy{
						PrefixRules: []string{
							"feat:",
						},
					},
				},
				Tasks: []Task{
					{
						ID:                   "task-001",
						Title:                "Test task",
						Description:          "A test task",
						FilesInScope:         []string{"file1.go"},
						SingleResponsibility: "Test implementation",
						Verification: Verification{
							PreCommit: []string{"go test"},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "incomplete commit policy",
			config: &RalphyYAML{
				RulesFile: "docs/opencode-rules.md",
				GlobalConstraints: GlobalConstraints{
					ForbiddenPatterns: []PatternConstraint{
						{Pattern: "\\bany\\b", Message: "Use unknown with type guards instead"},
					},
					RequiredPatterns: []PatternConstraint{
						{Pattern: "import.*zod", When: "validating_external_data"},
					},
				},
				CI: CI{
					Precommit: []string{
						"go test ./...",
					},
					CIChecks: []string{
						"build-and-test",
					},
				},
				Outputs: Outputs{
					AllowedFileEdits: []string{
						"docs/**",
					},
					DisallowedFileEdits: []string{
						".github/**",
					},
					CommitPolicy: CommitPolicy{},
				},
				Tasks: []Task{
					{
						ID:                   "task-001",
						Title:                "Test task",
						Description:          "A test task",
						FilesInScope:         []string{"file1.go"},
						SingleResponsibility: "Test implementation",
						Verification: Verification{
							PreCommit: []string{"go test"},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "incomplete scope enforcement",
			config: &RalphyYAML{
				RulesFile: "docs/opencode-rules.md",
				GlobalConstraints: GlobalConstraints{
					ForbiddenPatterns: []PatternConstraint{
						{Pattern: "\\bany\\b", Message: "Use unknown with type guards instead"},
					},
					RequiredPatterns: []PatternConstraint{
						{Pattern: "import.*zod", When: "validating_external_data"},
					},
				},
				CI: CI{
					Precommit: []string{
						"go test ./...",
					},
					CIChecks: []string{
						"build-and-test",
					},
				},
				Outputs: Outputs{
					AllowedFileEdits: []string{
						"docs/**",
					},
					CommitPolicy: CommitPolicy{
						PrefixRules: []string{
							"feat:",
						},
					},
				},
				Tasks: []Task{
					{
						ID:                   "task-001",
						Title:                "Test task",
						Description:          "A test task",
						FilesInScope:         []string{"file1.go"},
						SingleResponsibility: "Test implementation",
						Verification: Verification{
							PreCommit: []string{"go test"},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "tasks missing files_in_scope",
			config: &RalphyYAML{
				RulesFile: "docs/opencode-rules.md",
				GlobalConstraints: GlobalConstraints{
					ForbiddenPatterns: []PatternConstraint{
						{Pattern: "\\bany\\b", Message: "Use unknown with type guards instead"},
					},
					RequiredPatterns: []PatternConstraint{
						{Pattern: "import.*zod", When: "validating_external_data"},
					},
				},
				CI: CI{
					Precommit: []string{
						"go test ./...",
					},
					CIChecks: []string{
						"build-and-test",
					},
				},
				Outputs: Outputs{
					AllowedFileEdits: []string{
						"docs/**",
					},
					DisallowedFileEdits: []string{
						".github/**",
					},
					CommitPolicy: CommitPolicy{
						PrefixRules: []string{
							"feat:",
						},
					},
				},
				Tasks: []Task{
					{
						ID:                   "task-001",
						Title:                "Test task",
						Description:          "A test task",
						FilesInScope:         []string{},
						SingleResponsibility: "Test implementation",
						Verification: Verification{
							PreCommit: []string{"go test"},
						},
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateEnforcement(tt.config)
			if result.Valid != tt.expected {
				t.Errorf("validateEnforcement().Valid = %v, expected %v", result.Valid, tt.expected)
			}
		})
	}
}

// TestValidationResultStructure tests that validation result has correct structure
func TestValidationResultStructure(t *testing.T) {
	config := &RalphyYAML{
		RulesFile: "docs/opencode-rules.md",
		GlobalConstraints: GlobalConstraints{
			ForbiddenPatterns: []PatternConstraint{
				{Pattern: "\\bany\\b", Message: "Use unknown with type guards instead"},
			},
		},
		CI: CI{
			Precommit: []string{"go test"},
		},
		Outputs: Outputs{
			AllowedFileEdits:    []string{"docs/**"},
			DisallowedFileEdits: []string{".github/**"},
			CommitPolicy: CommitPolicy{
				PrefixRules: []string{"feat:"},
			},
		},
		Tasks: []Task{
			{
				ID:                   "task-001",
				Title:                "Test task",
				Description:          "A test task",
				FilesInScope:         []string{"file1.go"},
				SingleResponsibility: "Test implementation",
				Verification: Verification{
					PreCommit: []string{"go test"},
				},
			},
			{
				ID:                   "task-002",
				Title:                "Another task",
				Description:          "Another test task",
				FilesInScope:         []string{"file2.go"},
				SingleResponsibility: "Documentation update",
				Verification: Verification{
					PreCommit: []string{"markdownlint"},
				},
			},
		},
	}

	result := validateEnforcement(config)

	// Check basic structure
	if result.TotalTasks != 2 {
		t.Errorf("TotalTasks = %d, expected 2", result.TotalTasks)
	}

	if result.TasksWithFilesInScope != 2 {
		t.Errorf("TasksWithFilesInScope = %d, expected 2", result.TasksWithFilesInScope)
	}

	if result.TasksWithVerification != 2 {
		t.Errorf("TasksWithVerification = %d, expected 2", result.TasksWithVerification)
	}

	// Check verification layers
	if result.VerificationLayers.TotalLayers < 3 {
		t.Errorf("VerificationLayers.TotalLayers = %d, expected at least 3", result.VerificationLayers.TotalLayers)
	}

	// Check commit policy
	if !result.CommitPolicy.Complete {
		t.Error("CommitPolicy.Complete = false, expected true")
	}

	// Check scope enforcement
	if !result.ScopeEnforcement.Complete {
		t.Error("ScopeEnforcement.Complete = false, expected true")
	}

	if !result.ScopeEnforcement.AllTasksHaveFilesInScope {
		t.Error("ScopeEnforcement.AllTasksHaveFilesInScope = false, expected true")
	}
}

// TestCommandLineInterface tests the command-line interface
func TestCommandLineInterface(t *testing.T) {
	// Create a temporary test YAML file
	tempDir := t.TempDir()
	testYAML := `name: test-project
description: Test project
version: 1.0.0
rules_file: docs/opencode-rules.md
global_constraints:
  forbidden_patterns:
    - pattern: \bany\b
      message: Use unknown with type guards instead
  required_patterns:
    - pattern: import.*zod
      when: validating_external_data
ci:
  precommit:
    - go test ./...
  ci_checks:
    - build-and-test
outputs:
  allowed_file_edits:
    - docs/**
  disallowed_file_edits:
    - .github/**
  commit_policy:
    prefix_rules:
      - "feat:"
      - "fix:"
tasks:
  - id: task-001
    title: Test task
    description: A test task
    files_in_scope:
      - file1.go
    single_responsibility: Test implementation
    verification:
      pre_commit:
        - go test
`

	yamlPath := filepath.Join(tempDir, "test.yaml")
	if err := os.WriteFile(yamlPath, []byte(testYAML), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Test with valid file
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"validate", "--file", yamlPath}

	// We can't easily test the main function exit codes in a unit test
	// without complex mocking, so we'll just verify the file can be loaded
	config, err := loadYAML(yamlPath)
	if err != nil {
		t.Fatalf("Failed to load test YAML: %v", err)
	}

	if config.Name != "test-project" {
		t.Errorf("Config.Name = %s, expected test-project", config.Name)
	}

	if len(config.Tasks) != 1 {
		t.Errorf("Config.Tasks length = %d, expected 1", len(config.Tasks))
	}
}

// TestEmptyConfig tests validation with empty configuration
func TestEmptyConfig(t *testing.T) {
	config := &RalphyYAML{
		Tasks: []Task{},
	}

	result := validateEnforcement(config)

	if result.TotalTasks != 0 {
		t.Errorf("TotalTasks = %d, expected 0", result.TotalTasks)
	}

	if result.Valid {
		t.Error("Empty config should not be valid (missing verification layers)")
	}

	if len(result.Violations) == 0 {
		t.Error("Empty config should have violations")
	}
}

// TestMinimumVerificationLayers tests the minimum verification layers constant
func TestMinimumVerificationLayers(t *testing.T) {
	if minVerificationLayers != 3 {
		t.Errorf("minVerificationLayers = %d, expected 3", minVerificationLayers)
	}
}
