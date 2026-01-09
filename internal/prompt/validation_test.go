package prompt

import (
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestValidateMetadata(t *testing.T) {
	tests := []struct {
		name        string
		metadata    Metadata
		wantErr     bool
		errContains string
	}{
		{
			name: "valid minimal metadata",
			metadata: Metadata{
				Title: "Test",
			},
			wantErr: false,
		},
		{
			name: "valid full metadata",
			metadata: Metadata{
				Title:       "Valid Title",
				Description: "A valid description",
				Category:    "testing",
				Author:      "Test Author",
				Version:     "1.0.0",
				Tags:        []string{"test", "validation"},
			},
			wantErr: false,
		},
		{
			name: "title at max length",
			metadata: Metadata{
				Title: strings.Repeat("a", MaxTitleLength),
			},
			wantErr: false,
		},
		{
			name: "title exceeds max length",
			metadata: Metadata{
				Title: strings.Repeat("a", MaxTitleLength+1),
			},
			wantErr:     true,
			errContains: "title exceeds maximum length",
		},
		{
			name: "description at max length",
			metadata: Metadata{
				Title:       "Test",
				Description: strings.Repeat("a", MaxDescriptionLength),
			},
			wantErr: false,
		},
		{
			name: "description exceeds max length",
			metadata: Metadata{
				Title:       "Test",
				Description: strings.Repeat("a", MaxDescriptionLength+1),
			},
			wantErr:     true,
			errContains: "description exceeds maximum length",
		},
		{
			name: "category at max length",
			metadata: Metadata{
				Title:    "Test",
				Category: strings.Repeat("a", MaxCategoryLength),
			},
			wantErr: false,
		},
		{
			name: "category exceeds max length",
			metadata: Metadata{
				Title:    "Test",
				Category: strings.Repeat("a", MaxCategoryLength+1),
			},
			wantErr:     true,
			errContains: "category exceeds maximum length",
		},
		{
			name: "author at max length",
			metadata: Metadata{
				Title:  "Test",
				Author: strings.Repeat("a", MaxAuthorLength),
			},
			wantErr: false,
		},
		{
			name: "author exceeds max length",
			metadata: Metadata{
				Title:  "Test",
				Author: strings.Repeat("a", MaxAuthorLength+1),
			},
			wantErr:     true,
			errContains: "author exceeds maximum length",
		},
		{
			name: "tags at max count",
			metadata: Metadata{
				Title: "Test",
				Tags:  make([]string, MaxTagsCount),
			},
			wantErr: false,
		},
		{
			name: "tags exceed max count",
			metadata: Metadata{
				Title: "Test",
				Tags:  make([]string, MaxTagsCount+1),
			},
			wantErr:     true,
			errContains: "tags count exceeds maximum",
		},
		{
			name: "tag at max length",
			metadata: Metadata{
				Title: "Test",
				Tags:  []string{strings.Repeat("a", MaxTagLength)},
			},
			wantErr: false,
		},
		{
			name: "tag exceeds max length",
			metadata: Metadata{
				Title: "Test",
				Tags:  []string{strings.Repeat("a", MaxTagLength+1)},
			},
			wantErr:     true,
			errContains: "tag at index 0 exceeds maximum length",
		},
		{
			name: "valid semver",
			metadata: Metadata{
				Title:   "Test",
				Version: "1.0.0",
			},
			wantErr: false,
		},
		{
			name: "valid semver with pre-release",
			metadata: Metadata{
				Title:   "Test",
				Version: "1.0.0-alpha",
			},
			wantErr: false,
		},
		{
			name: "valid semver with build metadata",
			metadata: Metadata{
				Title:   "Test",
				Version: "1.0.0+build.1",
			},
			wantErr: false,
		},
		{
			name: "valid semver with v prefix",
			metadata: Metadata{
				Title:   "Test",
				Version: "v1.0.0",
			},
			wantErr: false,
		},
		{
			name: "invalid semver",
			metadata: Metadata{
				Title:   "Test",
				Version: "invalid",
			},
			wantErr:     true,
			errContains: "version must follow semantic versioning format",
		},
		{
			name: "empty version is allowed",
			metadata: Metadata{
				Title:   "Test",
				Version: "",
			},
			wantErr: false,
		},
		{
			name: "unicode characters in title",
			metadata: Metadata{
				Title: "日本語 🇯🇵 中文",
			},
			wantErr: false,
		},
		{
			name: "multiple tags with various lengths",
			metadata: Metadata{
				Title: "Test",
				Tags:  []string{"short", strings.Repeat("b", MaxTagLength), "medium-length-tag"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetadata(tt.metadata)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateMetadata() error should contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

func TestIsValidSemVer(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{"1.0.0", true},
		{"v1.0.0", true},
		{"1.2.3", true},
		{"10.20.30", true},
		{"1.0.0-alpha", true},
		{"1.0.0-alpha.1", true},
		{"1.0.0-0.3.7", true},
		{"1.0.0-x.7.z.92", true},
		{"1.0.0-alpha+001", true},
		{"1.0.0+20130313144700", true},
		{"1.0.0-beta+exp.sha.5114f85", true},
		{"1.0.0-beta.2", true},
		{"invalid", false},
		{"1.0", false},
		{"1", false},
		{"v1", false},
		{"1.0.0-alpha_beta", false},
		{"1.0.0..alpha", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			got := isValidSemVer(tt.version)
			if got != tt.want {
				t.Errorf("isValidSemVer(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestStorageSaveWithValidation(t *testing.T) {
	logger := zap.NewNop()
	storage := NewStorage(logger)

	t.Run("reject invalid metadata", func(t *testing.T) {
		tempDir := t.TempDir()
		testPath := filepath.Join(tempDir, "test.md")

		prompt := &Prompt{
			Metadata: Metadata{
				Title: strings.Repeat("a", MaxTitleLength+1),
			},
			Content: "Content",
		}

		err := storage.Save(testPath, prompt)
		if err == nil {
			t.Error("Storage.Save() should reject invalid metadata")
		}

		if !strings.Contains(err.Error(), "invalid metadata") {
			t.Errorf("error should mention invalid metadata, got: %v", err)
		}
	})

	t.Run("reject invalid version format", func(t *testing.T) {
		tempDir := t.TempDir()
		testPath := filepath.Join(tempDir, "test.md")

		prompt := &Prompt{
			Metadata: Metadata{
				Title:   "Test",
				Version: "invalid-version",
			},
			Content: "Content",
		}

		err := storage.Save(testPath, prompt)
		if err == nil {
			t.Error("Storage.Save() should reject invalid version format")
		}

		if !strings.Contains(err.Error(), "version must follow semantic versioning") {
			t.Errorf("error should mention version format, got: %v", err)
		}
	})

	t.Run("accept valid metadata", func(t *testing.T) {
		tempDir := t.TempDir()
		testPath := filepath.Join(tempDir, "test.md")

		prompt := &Prompt{
			Metadata: Metadata{
				Title:       "Valid Title",
				Description: "Valid description",
				Version:     "1.0.0",
				Tags:        []string{"valid", "metadata"},
			},
			Content: "Content",
		}

		err := storage.Save(testPath, prompt)
		if err != nil {
			t.Errorf("Storage.Save() should accept valid metadata, got error: %v", err)
		}
	})
}
