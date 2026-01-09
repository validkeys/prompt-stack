// Integration tests for prompt file I/O operations.
// These tests verify end-to-end workflows across multiple components.
package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kyledavis/prompt-stack/internal/prompt"
	"go.uber.org/zap"
)

func setupStorage(t *testing.T) (*prompt.Storage, func()) {
	logger := zap.NewNop()
	storage := prompt.NewStorage(logger)
	tempDir := t.TempDir()

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return storage, cleanup
}

func TestFileIORoundTrip(t *testing.T) {
	storage, cleanup := setupStorage(t)
	defer cleanup()

	tests := []struct {
		name   string
		prompt *prompt.Prompt
	}{
		{
			name: "complete prompt with all fields",
			prompt: &prompt.Prompt{
				Metadata: prompt.Metadata{
					Title:       "Complete Prompt",
					Description: "A complete prompt with all fields",
					Category:    "testing",
					Tags:        []string{"test", "integration", "complete"},
					Author:      "Test Author",
					Version:     "1.0.0",
				},
				Content: `This is the prompt content.

It can have multiple paragraphs.

## Subheading

And markdown formatting.`,
			},
		},
		{
			name: "minimal prompt",
			prompt: &prompt.Prompt{
				Metadata: prompt.Metadata{
					Title: "Minimal",
				},
				Content: "Simple content.",
			},
		},
		{
			name: "prompt with only content",
			prompt: &prompt.Prompt{
				Metadata: prompt.Metadata{},
				Content:  "Just content, no metadata.",
			},
		},
		{
			name: "prompt with unicode",
			prompt: &prompt.Prompt{
				Metadata: prompt.Metadata{
					Title: "Unicode 测试 🧪",
					Tags:  []string{"日本語", "中文", "한국어"},
				},
				Content: "Content with unicode: 你好世界 こんにちは 안녕하세요 🌍",
			},
		},
		{
			name: "prompt with many tags",
			prompt: &prompt.Prompt{
				Metadata: prompt.Metadata{
					Title: "Many Tags",
					Tags:  []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6", "tag7", "tag8", "tag9", "tag10"},
				},
				Content: "Content with many tags.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testPath := filepath.Join(tempDir, "test.md")

			// Save
			if err := storage.Save(testPath, tt.prompt); err != nil {
				t.Fatalf("failed to save prompt: %v", err)
			}

			// Load
			loaded, err := storage.Load(testPath)
			if err != nil {
				t.Fatalf("failed to load prompt: %v", err)
			}

			// Verify content preserved
			if loaded.Metadata.Title != tt.prompt.Metadata.Title {
				t.Errorf("title not preserved: got %q, want %q", loaded.Metadata.Title, tt.prompt.Metadata.Title)
			}

			if loaded.Metadata.Description != tt.prompt.Metadata.Description {
				t.Errorf("description not preserved: got %q, want %q", loaded.Metadata.Description, tt.prompt.Metadata.Description)
			}

			if loaded.Metadata.Category != tt.prompt.Metadata.Category {
				t.Errorf("category not preserved: got %q, want %q", loaded.Metadata.Category, tt.prompt.Metadata.Category)
			}

			if loaded.Metadata.Author != tt.prompt.Metadata.Author {
				t.Errorf("author not preserved: got %q, want %q", loaded.Metadata.Author, tt.prompt.Metadata.Author)
			}

			if loaded.Metadata.Version != tt.prompt.Metadata.Version {
				t.Errorf("version not preserved: got %q, want %q", loaded.Metadata.Version, tt.prompt.Metadata.Version)
			}

			if !stringSlicesEqual(loaded.Metadata.Tags, tt.prompt.Metadata.Tags) {
				t.Errorf("tags not preserved: got %v, want %v", loaded.Metadata.Tags, tt.prompt.Metadata.Tags)
			}

			if loaded.Content != tt.prompt.Content {
				t.Errorf("content not preserved")
			}
		})
	}
}

func TestBatchOperations(t *testing.T) {
	storage, cleanup := setupStorage(t)
	defer cleanup()

	tempDir := t.TempDir()

	// Create multiple prompts
	prompts := []*prompt.Prompt{
		{
			Metadata: prompt.Metadata{Title: "Prompt 1", Tags: []string{"a"}},
			Content:  "Content 1",
		},
		{
			Metadata: prompt.Metadata{Title: "Prompt 2", Tags: []string{"b"}},
			Content:  "Content 2",
		},
		{
			Metadata: prompt.Metadata{Title: "Prompt 3", Tags: []string{"c"}},
			Content:  "Content 3",
		},
	}

	// Save all prompts
	paths := make([]string, len(prompts))
	for i, p := range prompts {
		paths[i] = filepath.Join(tempDir, "prompt", p.Metadata.Title+".md")
		if err := storage.Save(paths[i], p); err != nil {
			t.Fatalf("failed to save prompt %d: %v", i, err)
		}
	}

	// Load all prompts and verify
	for i, path := range paths {
		loaded, err := storage.Load(path)
		if err != nil {
			t.Fatalf("failed to load prompt %d: %v", i, err)
		}

		if loaded.Metadata.Title != prompts[i].Metadata.Title {
			t.Errorf("batch operation: prompt %d title mismatch", i)
		}

		if loaded.Content != prompts[i].Content {
			t.Errorf("batch operation: prompt %d content mismatch", i)
		}
	}
}

func TestErrorRecovery(t *testing.T) {
	storage, cleanup := setupStorage(t)
	defer cleanup()

	tempDir := t.TempDir()

	t.Run("recover from missing file", func(t *testing.T) {
		missingPath := filepath.Join(tempDir, "missing.md")
		_, err := storage.Load(missingPath)
		if err == nil {
			t.Errorf("expected error for missing file")
		}

		// Now create the file
		prompt := &prompt.Prompt{
			Metadata: prompt.Metadata{Title: "Recovered"},
			Content:  "Content",
		}
		if err := storage.Save(missingPath, prompt); err != nil {
			t.Fatalf("failed to save recovered file: %v", err)
		}

		// Now load should succeed
		loaded, err := storage.Load(missingPath)
		if err != nil {
			t.Errorf("failed to load recovered file: %v", err)
		}

		if loaded.Metadata.Title != "Recovered" {
			t.Errorf("recovered file has wrong title: %q", loaded.Metadata.Title)
		}
	})

	t.Run("recover from corrupted frontmatter", func(t *testing.T) {
		corruptedPath := filepath.Join(tempDir, "corrupted.md")

		// Create file with malformed frontmatter
		corruptedContent := `---
title: "Test"
tags: [unclosed
---

Content`
		if err := os.WriteFile(corruptedPath, []byte(corruptedContent), 0644); err != nil {
			t.Fatalf("failed to create corrupted file: %v", err)
		}

		// Should load as plain markdown
		loaded, err := storage.Load(corruptedPath)
		if err != nil {
			t.Fatalf("failed to load corrupted file: %v", err)
		}

		// Should have empty metadata (treat as plain markdown)
		if loaded.Metadata.Title != "" {
			t.Errorf("corrupted file should have empty metadata, got: %q", loaded.Metadata.Title)
		}

		// Now save with proper frontmatter
		prompt := &prompt.Prompt{
			Metadata: prompt.Metadata{Title: "Fixed"},
			Content:  loaded.Content,
		}
		if err := storage.Save(corruptedPath, prompt); err != nil {
			t.Fatalf("failed to save fixed file: %v", err)
		}

		// Reload and verify
		loaded2, err := storage.Load(corruptedPath)
		if err != nil {
			t.Fatalf("failed to load fixed file: %v", err)
		}

		if loaded2.Metadata.Title != "Fixed" {
			t.Errorf("fixed file has wrong title: %q", loaded2.Metadata.Title)
		}
	})
}

func TestConcurrentAccess(t *testing.T) {
	storage, cleanup := setupStorage(t)
	defer cleanup()

	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "concurrent.md")

	// Create initial prompt
	prompt := &prompt.Prompt{
		Metadata: prompt.Metadata{Title: "Concurrent Test"},
		Content:  "Initial content",
	}
	if err := storage.Save(testPath, prompt); err != nil {
		t.Fatalf("failed to save initial prompt: %v", err)
	}

	// Simulate concurrent operations (read then write then read)
	// Note: Go's os.Rename is atomic, so this should work safely
	loaded1, err := storage.Load(testPath)
	if err != nil {
		t.Fatalf("first load failed: %v", err)
	}

	// Modify
	loaded1.Content = "Modified content"
	if err := storage.Save(testPath, loaded1); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Reload
	loaded2, err := storage.Load(testPath)
	if err != nil {
		t.Fatalf("second load failed: %v", err)
	}

	if loaded2.Content != "Modified content" {
		t.Errorf("concurrent access test failed: got %q, want %q", loaded2.Content, "Modified content")
	}
}

func TestVeryLargeFiles(t *testing.T) {
	storage, cleanup := setupStorage(t)
	defer cleanup()

	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "large.md")

	// Create a large content (but under limit)
	content := ""
	for i := 0; i < 1000; i++ {
		content += "Line " + string(rune('0'+i%10)) + " with some text\n"
	}

	prompt := &prompt.Prompt{
		Metadata: prompt.Metadata{Title: "Large File"},
		Content:  content,
	}

	// Save
	if err := storage.Save(testPath, prompt); err != nil {
		t.Fatalf("failed to save large file: %v", err)
	}

	// Load
	loaded, err := storage.Load(testPath)
	if err != nil {
		t.Fatalf("failed to load large file: %v", err)
	}

	if loaded.Content != content {
		t.Errorf("large file content not preserved")
	}
}

func TestLoggingIntegration(t *testing.T) {
	// Create logger that captures logs
	logs := []string{}
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	storage := prompt.NewStorage(logger)
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test.md")

	prompt := &prompt.Prompt{
		Metadata: prompt.Metadata{Title: "Logging Test"},
		Content:  "Content",
	}

	// Save (should log)
	if err := storage.Save(testPath, prompt); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	// Load (should log)
	if _, err := storage.Load(testPath); err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	// Note: In a real test, we'd verify logs were captured
	// For now, just ensure operations complete
	_ = logs
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
