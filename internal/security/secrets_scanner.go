package security

import (
	"regexp"
)

var (
	awsAccessKeyPattern = regexp.MustCompile(`AKIA[0-9A-Z]{16}`)
	awsSecretKeyPattern = regexp.MustCompile(`[0-9a-zA-Z/+]{40}`)
	githubTokenPattern  = regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}`)
	apiKeyPattern       = regexp.MustCompile(`(?i)api[_-]?key['":\s]*[=:]['"\s]*[a-zA-Z0-9_\-]{16,}`)
	secretPattern       = regexp.MustCompile(`(?i)secret['":\s]*[=:]['"\s]*[a-zA-Z0-9_\-]{16,}`)
	passwordPattern     = regexp.MustCompile(`(?i)password['":\s]*(?:is\s+|[=:]['"\s]+)[a-zA-Z0-9_\-]{4,}`)
	privateKeyPattern   = regexp.MustCompile(`-----BEGIN (RSA )?PRIVATE KEY-----`)
	bearerTokenPattern  = regexp.MustCompile(`Bearer [a-zA-Z0-9_\-\.]{20,}`)
)

type SecretMatch struct {
	PatternName string
	Match       string
	LineNumber  int
}

func ScanForSecrets(content string) []SecretMatch {
	matches := []SecretMatch{}

	patterns := map[string]*regexp.Regexp{
		"AWS Access Key": awsAccessKeyPattern,
		"AWS Secret Key": awsSecretKeyPattern,
		"GitHub Token":   githubTokenPattern,
		"API Key":        apiKeyPattern,
		"Secret":         secretPattern,
		"Password":       passwordPattern,
		"Private Key":    privateKeyPattern,
		"Bearer Token":   bearerTokenPattern,
	}

	for name, pattern := range patterns {
		if found := pattern.FindAllString(content, -1); len(found) > 0 {
			for _, match := range found {
				matches = append(matches, SecretMatch{
					PatternName: name,
					Match:       match,
					LineNumber:  0,
				})
			}
		}
	}

	return matches
}

func ContainsSecrets(content string) bool {
	matches := ScanForSecrets(content)
	return len(matches) > 0
}
