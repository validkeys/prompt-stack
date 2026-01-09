package prompt

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseFrontmatter parses YAML frontmatter from markdown content.
// Returns metadata, content after frontmatter, and error if any.
func ParseFrontmatter(content string) (Metadata, string, error) {
	var metadata Metadata
	var parsedContent string

	// Strip UTF-8 BOM if present
	const utf8BOM = "\xEF\xBB\xBF"
	if strings.HasPrefix(content, utf8BOM) {
		content = strings.TrimPrefix(content, utf8BOM)
	}

	// Normalize line endings to LF
	content = strings.ReplaceAll(content, "\r\n", "\n")

	// Split content into lines
	lines := strings.Split(content, "\n")

	// Check if content starts with frontmatter delimiter
	if len(lines) < 2 || lines[0] != FrontmatterDelimiter {
		// No frontmatter, return empty metadata and full content
		return Metadata{}, content, nil
	}

	// Find the closing delimiter
	var endIdx int
	found := false
	for i := 1; i < len(lines); i++ {
		if lines[i] == FrontmatterDelimiter {
			endIdx = i
			found = true
			break
		}
	}

	if !found {
		// Malformed frontmatter (no closing delimiter), treat as plain markdown
		return Metadata{}, content, nil
	}

	// Extract frontmatter content (between delimiters)
	frontmatterLines := lines[1:endIdx]
	frontmatterContent := strings.Join(frontmatterLines, "\n")

	// Parse YAML
	if err := yaml.Unmarshal([]byte(frontmatterContent), &metadata); err != nil {
		// Invalid YAML, treat as plain markdown with empty metadata
		return Metadata{}, content, nil
	}

	// Extract content after frontmatter
	if endIdx+1 < len(lines) {
		// Skip the first line after delimiter (empty line)
		// and only include remaining content
		contentLines := lines[endIdx+1:]
		if len(contentLines) > 0 && contentLines[0] == "" {
			contentLines = contentLines[1:]
		}
		parsedContent = strings.Join(contentLines, "\n")
	} else {
		parsedContent = ""
	}

	return metadata, parsedContent, nil
}

// FormatFrontmatter formats metadata and content into a markdown file with frontmatter.
func FormatFrontmatter(metadata Metadata, content string) string {
	var buf bytes.Buffer

	buf.WriteString(FrontmatterDelimiter)
	buf.WriteString("\n")

	yamlData, err := yaml.Marshal(metadata)
	if err != nil {
		// If marshaling fails, write minimal frontmatter
		buf.WriteString("title: \"")
		buf.WriteString(metadata.Title)
		buf.WriteString("\"\n")
	} else {
		buf.Write(yamlData)
	}

	buf.WriteString(FrontmatterDelimiter)
	buf.WriteString("\n")

	if content != "" {
		buf.WriteString(content)
	}

	return buf.String()
}
