package security

import (
	"testing"
)

func TestScanForSecrets(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectMatches int
	}{
		{
			name:          "no secrets",
			content:       "This is just some normal text",
			expectMatches: 0,
		},
		{
			name:          "AWS access key",
			content:       "My key is AKIAIOSFODNN7EXAMPLE",
			expectMatches: 1,
		},
		{
			name:          "AWS secret key",
			content:       "The secret is wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			expectMatches: 1,
		},
		{
			name:          "GitHub token",
			content:       "Token: ghp_1234567890abcdefghijklmnopqrstuvwxyz",
			expectMatches: 1,
		},
		{
			name:          "API key",
			content:       "api_key = my_secret_api_key_123456",
			expectMatches: 1,
		},
		{
			name:          "Secret variable",
			content:       "secret: very_secret_value_123456789",
			expectMatches: 1,
		},
		{
			name:          "Password",
			content:       "password: mypassword1234",
			expectMatches: 1,
		},
		{
			name:          "Private key",
			content:       "-----BEGIN RSA PRIVATE KEY-----",
			expectMatches: 1,
		},
		{
			name:          "Bearer token",
			content:       "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectMatches: 1,
		},
		{
			name:          "multiple secrets",
			content:       "AKIAIOSFODNN7EXAMPLE and secret: very_secret_value_12345678",
			expectMatches: 2,
		},
		{
			name:          "case insensitive secret",
			content:       "SECRET: very_secret_value_12345678",
			expectMatches: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := ScanForSecrets(tt.content)
			if len(matches) != tt.expectMatches {
				t.Errorf("ScanForSecrets() returned %d matches, expected %d", len(matches), tt.expectMatches)
			}
		})
	}
}

func TestContainsSecrets(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "normal text",
			content:  "This is just some normal text",
			expected: false,
		},
		{
			name:     "contains API key",
			content:  "api_key = my_secret_api_key_123456",
			expected: true,
		},
		{
			name:     "contains password",
			content:  "password: mypassword1234",
			expected: true,
		},
		{
			name:     "contains AWS key",
			content:  "AKIAIOSFODNN7EXAMPLE",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsSecrets(tt.content)
			if result != tt.expected {
				t.Errorf("ContainsSecrets() returned %v, expected %v", result, tt.expected)
			}
		})
	}
}
