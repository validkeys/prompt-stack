package prompt

import (
	"strings"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		wantMetadata Metadata
		wantContent  string
		wantErr      bool
		description  string
	}{
		{
			name: "valid frontmatter with multiple keys",
			content: `---
title: "Test Prompt"
tags:
  - test
  - example
category: testing
description: A test prompt
---

This is the content.`,
			wantMetadata: Metadata{
				Title:       "Test Prompt",
				Tags:        []string{"test", "example"},
				Category:    "testing",
				Description: "A test prompt",
			},
			wantContent: "This is the content.",
		},
		{
			name:         "file without frontmatter",
			content:      `This is just markdown content.`,
			wantMetadata: Metadata{},
			wantContent:  "This is just markdown content.",
		},
		{
			name: "file with only frontmatter",
			content: `---
title: "Only Frontmatter"
tags:
  - tag1
---

`,
			wantMetadata: Metadata{
				Title: "Only Frontmatter",
				Tags:  []string{"tag1"},
			},
			wantContent: "",
		},
		{
			name: "malformed frontmatter missing closing delimiter",
			content: `---
title: "Missing Closing"
tags:
  - tag1

This is content that should be treated as plain markdown.`,
			wantMetadata: Metadata{},
			wantContent: `---
title: "Missing Closing"
tags:
  - tag1

This is content that should be treated as plain markdown.`,
		},
		{
			name: "frontmatter with quoted values",
			content: `---
title: "Quoted Title"
description: "Description with 'quotes' and \"double quotes\""
---

Content here.`,
			wantMetadata: Metadata{
				Title:       "Quoted Title",
				Description: `Description with 'quotes' and "double quotes"`,
			},
			wantContent: "Content here.",
		},
		{
			name: "empty frontmatter",
			content: `---
---

Content after empty frontmatter.`,
			wantMetadata: Metadata{},
			wantContent:  "Content after empty frontmatter.",
		},
		{
			name: "frontmatter with empty values",
			content: `---
title: ""
description: ""
tags: []
category: ""

---

Content.`,
			wantMetadata: Metadata{
				Title:       "",
				Description: "",
				Tags:        []string{},
				Category:    "",
			},
			wantContent: "Content.",
		},
		{
			name: "frontmatter with only title",
			content: `---
title: "Simple Prompt"
---

Simple content.`,
			wantMetadata: Metadata{
				Title: "Simple Prompt",
			},
			wantContent: "Simple content.",
		},
		{
			name: "frontmatter with empty lines",
			content: `---
title: "Prompt with Empty Lines"

tags:
  - tag1
  - tag2


description: "Test"


---

Content here.`,
			wantMetadata: Metadata{
				Title:       "Prompt with Empty Lines",
				Tags:        []string{"tag1", "tag2"},
				Description: "Test",
			},
			wantContent: "Content here.",
		},
		{
			name: "invalid YAML in frontmatter",
			content: `---
title: "Invalid YAML"
tags: [unclosed bracket
description: "Test"
---

This should be treated as plain markdown.`,
			wantMetadata: Metadata{},
			wantContent: `---
title: "Invalid YAML"
tags: [unclosed bracket
description: "Test"
---

This should be treated as plain markdown.`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMetadata, gotContent, err := ParseFrontmatter(tt.content)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFrontmatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotMetadata.Title != tt.wantMetadata.Title {
				t.Errorf("ParseFrontmatter() Title = %q, want %q", gotMetadata.Title, tt.wantMetadata.Title)
			}

			if gotMetadata.Category != tt.wantMetadata.Category {
				t.Errorf("ParseFrontmatter() Category = %q, want %q", gotMetadata.Category, tt.wantMetadata.Category)
			}

			if gotMetadata.Description != tt.wantMetadata.Description {
				t.Errorf("ParseFrontmatter() Description = %q, want %q", gotMetadata.Description, tt.wantMetadata.Description)
			}

			if !stringSlicesEqual(gotMetadata.Tags, tt.wantMetadata.Tags) {
				t.Errorf("ParseFrontmatter() Tags = %v, want %v", gotMetadata.Tags, tt.wantMetadata.Tags)
			}

			if gotContent != tt.wantContent {
				t.Errorf("ParseFrontmatter() Content = %q, want %q", gotContent, tt.wantContent)
			}
		})
	}
}

func TestParseFrontmatterEdgeCases(t *testing.T) {
	t.Run("very long frontmatter", func(t *testing.T) {
		lines := []string{FrontmatterDelimiter}
		lines = append(lines, "title: \"Long Frontmatter\"")
		for i := 0; i < 100; i++ {
			lines = append(lines, "field"+string(rune('0'+i%10))+string(rune('A'+i%26))+": \"value\"")
		}
		lines = append(lines, FrontmatterDelimiter)
		lines = append(lines, "Content")

		content := strings.Join(lines, "\n")
		_, gotContent, err := ParseFrontmatter(content)

		if err != nil {
			t.Errorf("ParseFrontmatter() unexpected error: %v", err)
		}

		if gotContent != "Content" {
			t.Errorf("ParseFrontmatter() Content = %q, want %q", gotContent, "Content")
		}
	})

	t.Run("unicode characters in values", func(t *testing.T) {
		content := "---\ntitle: \"Test 日本語 🇯🇵\"\ndescription: \"Description with émojis 😀 and spëcial çhars\"\ntags:\n  - français\n  - español\n---\n\nContent with unicode: 你好世界"
		_, gotContent, err := ParseFrontmatter(content)

		if err != nil {
			t.Errorf("ParseFrontmatter() unexpected error: %v", err)
		}

		if !strings.Contains(gotContent, "你好世界") {
			t.Errorf("ParseFrontmatter() should preserve unicode content")
		}
	})

	t.Run("special characters in keys", func(t *testing.T) {
		content := "---\ntitle: \"Test\"\ncustom_field: \"value\"\nanother.field: \"value2\"\n---\n\nContent"
		metadata, gotContent, err := ParseFrontmatter(content)

		if err != nil {
			t.Errorf("ParseFrontmatter() unexpected error: %v", err)
		}

		if gotContent != "Content" {
			t.Errorf("ParseFrontmatter() Content = %q, want %q", gotContent, "Content")
		}

		if metadata.Title != "Test" {
			t.Errorf("ParseFrontmatter() Title = %q, want %q", metadata.Title, "Test")
		}
	})

	t.Run("multiple --- markers in file", func(t *testing.T) {
		content := "---\ntitle: \"Test\"\n---\n\nContent with --- in the middle\nand another --- here"
		metadata, gotContent, err := ParseFrontmatter(content)

		if err != nil {
			t.Errorf("ParseFrontmatter() unexpected error: %v", err)
		}

		if metadata.Title != "Test" {
			t.Errorf("ParseFrontmatter() Title = %q, want %q", metadata.Title, "Test")
		}

		if !strings.Contains(gotContent, "Content with --- in the middle") {
			t.Errorf("ParseFrontmatter() should preserve --- in content")
		}
	})
}

func TestFormatFrontmatter(t *testing.T) {
	tests := []struct {
		name         string
		metadata     Metadata
		content      string
		wantContains []string
		description  string
	}{
		{
			name: "full metadata",
			metadata: Metadata{
				Title:       "Test Prompt",
				Tags:        []string{"test", "example"},
				Category:    "testing",
				Description: "A test prompt",
			},
			content: "This is the content.",
			wantContains: []string{
				"---",
				"title: Test Prompt",
				"- test",
				"- example",
				"category: testing",
				"description: A test prompt",
				"This is the content.",
			},
		},
		{
			name: "minimal metadata",
			metadata: Metadata{
				Title: "Simple",
			},
			content: "Content",
			wantContains: []string{
				"---",
				"title: Simple",
				"Content",
			},
		},
		{
			name: "empty content",
			metadata: Metadata{
				Title: "No Content",
			},
			content: "",
			wantContains: []string{
				"---",
				"title: No Content",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatFrontmatter(tt.metadata, tt.content)

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("FormatFrontmatter() missing expected content: %q\nGot:\n%s", want, got)
				}
			}
		})
	}
}

func TestFrontmatterRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "full round trip",
			content: `---
title: "Test Prompt"
tags:
  - test
  - example
category: testing
description: A test prompt
---

This is the content.`,
		},
		{
			name:    "no frontmatter",
			content: `Just markdown content.`,
		},
		{
			name: "empty frontmatter",
			content: `---
---

Content.`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata1, content1, err := ParseFrontmatter(tt.content)
			if err != nil {
				t.Fatalf("ParseFrontmatter() failed: %v", err)
			}

			formatted := FormatFrontmatter(metadata1, content1)

			metadata2, content2, err := ParseFrontmatter(formatted)
			if err != nil {
				t.Fatalf("ParseFrontmatter() on formatted content failed: %v", err)
			}

			if metadata1.Title != metadata2.Title {
				t.Errorf("Round trip Title changed: %q -> %q", metadata1.Title, metadata2.Title)
			}

			if metadata1.Category != metadata2.Category {
				t.Errorf("Round trip Category changed: %q -> %q", metadata1.Category, metadata2.Category)
			}

			if metadata1.Description != metadata2.Description {
				t.Errorf("Round trip Description changed: %q -> %q", metadata1.Description, metadata2.Description)
			}

			if !stringSlicesEqual(metadata1.Tags, metadata2.Tags) {
				t.Errorf("Round trip Tags changed: %v -> %v", metadata1.Tags, metadata2.Tags)
			}

			if content1 != content2 {
				t.Errorf("Round trip Content changed: %q -> %q", content1, content2)
			}
		})
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
