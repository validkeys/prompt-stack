package prompt

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
)

func BenchmarkParseFrontmatter(b *testing.B) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "small frontmatter",
			content: `---
title: "Small Prompt"
tags:
  - test
---

Content here.`,
		},
		{
			name: "medium frontmatter",
			content: `---
title: "Medium Prompt"
description: "A medium prompt with more metadata"
tags:
  - test
  - example
  - benchmark
  - performance
category: testing
author: "Benchmark Author"
version: "1.0.0"

---

Content here.`,
		},
		{
			name: "large frontmatter",
			content: `---
title: "Large Prompt"
description: "A large prompt with many fields"
tags:
  - tag1
  - tag2
  - tag3
  - tag4
  - tag5
  - tag6
  - tag7
  - tag8
  - tag9
  - tag10
category: testing
author: "Benchmark Author"
version: "1.0.0"
custom_field1: "value1"
custom_field2: "value2"
custom_field3: "value3"
custom_field4: "value4"
custom_field5: "value5"

---

Content here.`,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _ = ParseFrontmatter(tt.content)
			}
		})
	}
}

func BenchmarkFormatFrontmatter(b *testing.B) {
	tests := []struct {
		name     string
		metadata Metadata
		content  string
	}{
		{
			name: "small",
			metadata: Metadata{
				Title: "Small",
			},
			content: "Content",
		},
		{
			name: "medium",
			metadata: Metadata{
				Title:       "Medium",
				Description: "Description",
				Category:    "testing",
				Tags:        []string{"test", "example"},
			},
			content: "Content",
		},
		{
			name: "large",
			metadata: Metadata{
				Title:       "Large",
				Description: "Description",
				Category:    "testing",
				Tags:        []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
				Author:      "Author",
				Version:     "1.0.0",
			},
			content: "Content",
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = FormatFrontmatter(tt.metadata, tt.content)
			}
		})
	}
}

func BenchmarkStorageLoad(b *testing.B) {
	tempDir := b.TempDir()

	// Create test files of different sizes
	testCases := []struct {
		name    string
		size    int
		content string
	}{
		{
			name:    "1KB file",
			size:    1024,
			content: generateTestContent(1 * 1024),
		},
		{
			name:    "10KB file",
			size:    10 * 1024,
			content: generateTestContent(10 * 1024),
		},
		{
			name:    "100KB file",
			size:    100 * 1024,
			content: generateTestContent(100 * 1024),
		},
		{
			name:    "1MB file",
			size:    1024 * 1024,
			content: generateTestContent(1024 * 1024),
		},
	}

	storage := NewStorage(zap.NewNop())

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			// Create test file
			testPath := filepath.Join(tempDir, tc.name+".md")
			if err := os.WriteFile(testPath, []byte(tc.content), 0644); err != nil {
				b.Fatalf("failed to create test file: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = storage.Load(testPath)
			}
		})
	}
}

func BenchmarkStorageSave(b *testing.B) {
	tempDir := b.TempDir()

	testCases := []struct {
		name    string
		size    int
		content string
	}{
		{
			name:    "1KB file",
			size:    1024,
			content: generateTestContent(1 * 1024),
		},
		{
			name:    "10KB file",
			size:    10 * 1024,
			content: generateTestContent(10 * 1024),
		},
		{
			name:    "100KB file",
			size:    100 * 1024,
			content: generateTestContent(100 * 1024),
		},
		{
			name:    "1MB file",
			size:    1024 * 1024,
			content: generateTestContent(1024 * 1024),
		},
	}

	storage := NewStorage(zap.NewNop())

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			prompt := &Prompt{
				Metadata: Metadata{
					Title:       tc.name,
					Description: "Benchmark test",
					Category:    "benchmark",
					Tags:        []string{"benchmark", "performance"},
				},
				Content: tc.content,
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				testPath := filepath.Join(tempDir, tc.name+"_"+string(rune('0'+i%10))+".md")
				_ = storage.Save(testPath, prompt)
			}
		})
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	tempDir := b.TempDir()

	testCases := []struct {
		name    string
		size    int
		content string
	}{
		{
			name:    "1KB round trip",
			size:    1024,
			content: generateTestContent(1 * 1024),
		},
		{
			name:    "10KB round trip",
			size:    10 * 1024,
			content: generateTestContent(10 * 1024),
		},
		{
			name:    "100KB round trip",
			size:    100 * 1024,
			content: generateTestContent(100 * 1024),
		},
	}

	storage := NewStorage(zap.NewNop())

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			prompt := &Prompt{
				Metadata: Metadata{
					Title:       tc.name,
					Description: "Round trip test",
					Category:    "benchmark",
					Tags:        []string{"benchmark", "round-trip"},
				},
				Content: tc.content,
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				testPath := filepath.Join(tempDir, tc.name+"_"+string(rune('0'+i%10))+".md")
				_ = storage.Save(testPath, prompt)
				_, _ = storage.Load(testPath)
			}
		})
	}
}

func BenchmarkBatchOperations(b *testing.B) {
	tempDir := b.TempDir()
	storage := NewStorage(zap.NewNop())

	// Create multiple prompts
	prompts := make([]*Prompt, 100)
	for i := 0; i < 100; i++ {
		prompts[i] = &Prompt{
			Metadata: Metadata{
				Title:       "Prompt " + string(rune('0'+i%10)),
				Description: "Batch test",
				Tags:        []string{"batch", "test"},
			},
			Content: generateTestContent(1024), // 1KB each
		}
	}

	paths := make([]string, 100)
	for i := 0; i < 100; i++ {
		paths[i] = filepath.Join(tempDir, "batch", "prompt"+string(rune('0'+i%10))+".md")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Save all
		for _, p := range prompts {
			for _, path := range paths {
				_ = storage.Save(path, p)
			}
		}
	}
}

func generateTestContent(size int) string {
	content := `---
title: "Test Prompt"
description: "Test content for benchmarks"
tags:
  - test
  - benchmark
---

`

	remaining := size - len(content)
	for remaining > 0 {
		line := "This is test line content for benchmarking purposes.\n"
		if len(line) > remaining {
			line = line[:remaining]
		}
		content += line
		remaining -= len(line)
	}

	return content
}
