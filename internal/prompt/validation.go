package prompt

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

const (
	MaxTitleLength       = 200
	MaxDescriptionLength = 1000
	MaxCategoryLength    = 100
	MaxAuthorLength      = 100
	MaxTagsCount         = 20
	MaxTagLength         = 50
)

var (
	semverRegex = regexp.MustCompile(`^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
)

// ValidateMetadata validates frontmatter metadata fields.
// Returns error if any validation fails.
func ValidateMetadata(metadata Metadata) error {
	if utf8.RuneCountInString(metadata.Title) > MaxTitleLength {
		return fmt.Errorf("title exceeds maximum length of %d characters", MaxTitleLength)
	}

	if utf8.RuneCountInString(metadata.Description) > MaxDescriptionLength {
		return fmt.Errorf("description exceeds maximum length of %d characters", MaxDescriptionLength)
	}

	if utf8.RuneCountInString(metadata.Category) > MaxCategoryLength {
		return fmt.Errorf("category exceeds maximum length of %d characters", MaxCategoryLength)
	}

	if utf8.RuneCountInString(metadata.Author) > MaxAuthorLength {
		return fmt.Errorf("author exceeds maximum length of %d characters", MaxAuthorLength)
	}

	if len(metadata.Tags) > MaxTagsCount {
		return fmt.Errorf("tags count exceeds maximum of %d", MaxTagsCount)
	}

	for i, tag := range metadata.Tags {
		if utf8.RuneCountInString(tag) > MaxTagLength {
			return fmt.Errorf("tag at index %d exceeds maximum length of %d characters", i, MaxTagLength)
		}
	}

	if metadata.Version != "" && !isValidSemVer(metadata.Version) {
		return fmt.Errorf("version must follow semantic versioning format (e.g., 1.0.0)")
	}

	return nil
}

// isValidSemVer checks if a string is a valid semantic version.
func isValidSemVer(version string) bool {
	return semverRegex.MatchString(version)
}
