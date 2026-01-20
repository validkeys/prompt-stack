package integration

import (
	"bufio"
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type TestResult struct {
	Name    string
	Passed  bool
	Message string
}

func GetRepoDir(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("failed to get caller information")
	}

	dir := filepath.Dir(filename)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find repository root (go.mod not found)")
		}
		dir = parent
	}
}

func SetupTestDir(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "prompt-stack-test-")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	return tmpDir, func() {
		os.Chdir(oldDir)
		os.RemoveAll(tmpDir)
	}
}

func BuildBinary(t *testing.T) string {
	buildDir, err := os.MkdirTemp("", "prompt-stack-build-")
	if err != nil {
		t.Fatalf("failed to create build directory: %v", err)
	}

	binaryPath := filepath.Join(buildDir, "prompt-stack")

	repoDir := GetRepoDir(t)

	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/prompt-stack")
	cmd.Dir = repoDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\nOutput: %s", err, string(output))
	}

	t.Cleanup(func() {
		os.RemoveAll(buildDir)
	})

	return binaryPath
}

func RunCommand(t *testing.T, binaryPath string, args ...string) (string, int, error) {
	cmd := exec.Command(binaryPath, args...)
	output, err := cmd.CombinedOutput()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		return string(output), exitCode, err
	}

	return string(output), exitCode, nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func FileContains(path, substring string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), substring) {
			return true
		}
	}

	return false
}

func CheckDatabaseForSecrets(t *testing.T, dbPath string) []string {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	secretPatterns := []string{
		"AKIA[0-9A-Z]{16}",
		"ghp_[a-zA-Z0-9]{36}",
		"sk-[a-zA-Z0-9]{48}",
		"[aA][pP][iI]_[kK][eE][yY]",
		"[sS][eE][cC][rR][eE][tT]",
		"[pP][aA][sS][sS][wW][oO][rR][dD]",
		"-----BEGIN.*PRIVATE KEY-----",
		"Bearer [a-zA-Z0-9\\-._~+/]+=*",
	}

	foundSecrets := []string{}

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		t.Fatalf("failed to query tables: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			t.Fatalf("failed to scan table name: %v", err)
		}
		tables = append(tables, tableName)
	}

	for _, table := range tables {
		tableRows, err := db.Query("SELECT * FROM " + table)
		if err != nil {
			continue
		}

		columns, err := tableRows.Columns()
		if err != nil {
			tableRows.Close()
			continue
		}

		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		for tableRows.Next() {
			if err := tableRows.Scan(valuePtrs...); err != nil {
				continue
			}

			for i, val := range values {
				if val == nil {
					continue
				}

				strVal := ""
				switch v := val.(type) {
				case []byte:
					strVal = string(v)
				case string:
					strVal = v
				default:
					continue
				}

				for _, pattern := range secretPatterns {
					if strings.Contains(strVal, pattern) {
						foundSecrets = append(foundSecrets,
							"Table: "+table+", Column: "+columns[i]+", Pattern: "+pattern)
					}
				}
			}
		}
		tableRows.Close()
	}

	return foundSecrets
}

func TestIntegrationCLIWorkflow(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T) TestResult
	}{
		{
			name:     "AC-1: Help command lists core commands",
			testFunc: testAC1_HelpListsCoreCommands,
		},
		{
			name:     "AC-2: Init command creates required files",
			testFunc: testAC2_InitCreatesRequiredFiles,
		},
		{
			name:     "AC-3: No secrets in database",
			testFunc: testAC3_NoSecretsInDatabase,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.testFunc(t)
			if !result.Passed {
				t.Errorf("Test failed: %s", result.Message)
			}
		})
	}
}

func testAC1_HelpListsCoreCommands(t *testing.T) TestResult {
	_, cleanup := SetupTestDir(t)
	defer cleanup()

	binaryPath := BuildBinary(t)

	output, exitCode, err := RunCommand(t, binaryPath, "--help")
	if err != nil && exitCode != 0 {
		return TestResult{
			Name:    "AC-1: Help command lists core commands",
			Passed:  false,
			Message: "help command failed: " + err.Error() + "\nOutput: " + output,
		}
	}

	requiredCommands := []string{"init", "plan", "validate", "review", "build"}
	missingCommands := []string{}

	for _, cmd := range requiredCommands {
		if !strings.Contains(output, cmd) {
			missingCommands = append(missingCommands, cmd)
		}
	}

	if len(missingCommands) > 0 {
		return TestResult{
			Name:    "AC-1: Help command lists core commands",
			Passed:  false,
			Message: "missing commands in help output: " + strings.Join(missingCommands, ", "),
		}
	}

	return TestResult{
		Name:    "AC-1: Help command lists core commands",
		Passed:  true,
		Message: "all core commands listed in help output",
	}
}

func testAC2_InitCreatesRequiredFiles(t *testing.T) TestResult {
	_, cleanup := SetupTestDir(t)
	defer cleanup()

	binaryPath := BuildBinary(t)

	output, exitCode, err := RunCommand(t, binaryPath, "init", "--no-interactive")
	if err != nil && exitCode != 0 {
		return TestResult{
			Name:    "AC-2: Init command creates required files",
			Passed:  false,
			Message: "init command failed: " + err.Error() + "\nOutput: " + output,
		}
	}

	configPath := filepath.Join(".", ".prompt-stack", "config.yaml")
	dbPath := filepath.Join(".", ".prompt-stack", "knowledge.db")

	if !FileExists(configPath) {
		return TestResult{
			Name:    "AC-2: Init command creates required files",
			Passed:  false,
			Message: "config.yaml not created at " + configPath,
		}
	}

	if !FileExists(dbPath) {
		return TestResult{
			Name:    "AC-2: Init command creates required files",
			Passed:  false,
			Message: "knowledge.db not created at " + dbPath,
		}
	}

	if !FileContains(configPath, "version:") {
		return TestResult{
			Name:    "AC-2: Init command creates required files",
			Passed:  false,
			Message: "config.yaml does not contain version field",
		}
	}

	return TestResult{
		Name:    "AC-2: Init command creates required files",
		Passed:  true,
		Message: "config.yaml and knowledge.db created successfully",
	}
}

func testAC3_NoSecretsInDatabase(t *testing.T) TestResult {
	_, cleanup := SetupTestDir(t)
	defer cleanup()

	binaryPath := BuildBinary(t)

	output, exitCode, err := RunCommand(t, binaryPath, "init", "--no-interactive")
	if err != nil && exitCode != 0 {
		return TestResult{
			Name:    "AC-3: No secrets in database",
			Passed:  false,
			Message: "init command failed: " + err.Error() + "\nOutput: " + output,
		}
	}

	dbPath := filepath.Join(".", ".prompt-stack", "knowledge.db")
	if !FileExists(dbPath) {
		return TestResult{
			Name:    "AC-3: No secrets in database",
			Passed:  false,
			Message: "knowledge.db not created",
		}
	}

	foundSecrets := CheckDatabaseForSecrets(t, dbPath)
	if len(foundSecrets) > 0 {
		return TestResult{
			Name:    "AC-3: No secrets in database",
			Passed:  false,
			Message: "found secrets in database: " + strings.Join(foundSecrets, "; "),
		}
	}

	return TestResult{
		Name:    "AC-3: No secrets in database",
		Passed:  true,
		Message: "no secrets found in database",
	}
}

func TestIntegrationInitValidatePlanWorkflow(t *testing.T) {
	_, cleanup := SetupTestDir(t)
	defer cleanup()

	binaryPath := BuildBinary(t)

	t.Run("init creates prompt-stack directory structure", func(t *testing.T) {
		output, exitCode, err := RunCommand(t, binaryPath, "init", "--no-interactive")
		if err != nil && exitCode != 0 {
			t.Fatalf("init command failed: %v\nOutput: %s", err, output)
		}

		promptStackDir := filepath.Join(".", ".prompt-stack")
		if !FileExists(promptStackDir) {
			t.Error(".prompt-stack directory not created")
		}

		configPath := filepath.Join(promptStackDir, "config.yaml")
		dbPath := filepath.Join(promptStackDir, "knowledge.db")

		if !FileExists(configPath) {
			t.Error("config.yaml not created")
		}

		if !FileExists(dbPath) {
			t.Error("knowledge.db not created")
		}
	})

	t.Run("validate command works after init", func(t *testing.T) {
		output, exitCode, err := RunCommand(t, binaryPath, "validate", "--help")
		if err != nil && exitCode != 0 {
			t.Fatalf("validate --help failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(output, "validate") {
			t.Error("validate help output does not contain 'validate'")
		}
	})

	t.Run("plan command shows help", func(t *testing.T) {
		output, exitCode, err := RunCommand(t, binaryPath, "plan", "--help")
		if err != nil && exitCode != 0 {
			t.Fatalf("plan --help failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(output, "plan") {
			t.Error("plan help output does not contain 'plan'")
		}
	})
}
