package database

import (
	"fmt"
	"strings"

	"github.com/kyledavis/prompt-stack/internal/security"
)

var (
	sensitiveColumns = map[string]bool{
		"patterns.description":           true,
		"patterns.source":                true,
		"requirements.content":           true,
		"tasks.description":              true,
		"validation_reports.report_json": true,
	}
)

func ValidateNoSecrets(tableName string, data map[string]interface{}) error {
	for columnName := range sensitiveColumns {
		prefix := tableName + "."
		if strings.HasPrefix(columnName, prefix) {
			columnNameWithoutPrefix := strings.TrimPrefix(columnName, prefix)
			value, exists := data[columnNameWithoutPrefix]
			if exists {
				if strValue, ok := value.(string); ok {
					if security.ContainsSecrets(strValue) {
						return fmt.Errorf("potential secret detected in %s.%s", tableName, columnNameWithoutPrefix)
					}
				}
			}
		}
	}
	return nil
}

func ValidatePatternData(data map[string]interface{}) error {
	return ValidateNoSecrets("patterns", data)
}

func ValidateRequirementData(data map[string]interface{}) error {
	return ValidateNoSecrets("requirements", data)
}

func ValidateTaskData(data map[string]interface{}) error {
	return ValidateNoSecrets("tasks", data)
}

func ValidateValidationReportData(data map[string]interface{}) error {
	return ValidateNoSecrets("validation_reports", data)
}
