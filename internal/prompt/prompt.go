// Package prompt provides prompt file I/O and frontmatter parsing for
// markdown-based prompt templates.
package prompt

// Metadata represents typed frontmatter fields for prompt files.
type Metadata struct {
	Title       string   `yaml:"title" json:"title"`
	Tags        []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	Category    string   `yaml:"category,omitempty" json:"category,omitempty"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Author      string   `yaml:"author,omitempty" json:"author,omitempty"`
	Version     string   `yaml:"version,omitempty" json:"version,omitempty"`
}

// Prompt represents a parsed prompt file with metadata and content.
type Prompt struct {
	Metadata Metadata `yaml:"-" json:"metadata"`
	Content  string   `yaml:"-" json:"content"`
}

const (
	// FrontmatterDelimiter is the marker that separates frontmatter from content
	FrontmatterDelimiter = "---"

	// MaxFileSize is the maximum allowed size for a prompt file in bytes
	MaxFileSize = 1 << 20 // 1MB
)
