package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestStorageLoad(t *testing.T) {
	logger := zap.NewNop()
	storage := NewStorage(logger)

	tests := []struct {
		name        string
		setupFile   func(t *testing.T, path string)
		wantPrompt  *Prompt
		wantErr     bool
		errContains string
	}{
		{
			name: "successful read with frontmatter",
			setupFile: func(t *testing.T, path string) {
				content := `---
title: "Test Prompt"
tags:
  - test
  - example
category: testing
description: A test prompt
---

This is the content.`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			},
			wantPrompt: &Prompt{
				Metadata: Metadata{
					Title:       "Test Prompt",
					Tags:        []string{"test", "example"},
					Category:    "testing",
					Description: "A test prompt",
				},
				Content: "This is the content.",
			},
		},
		{
			name: "successful read without frontmatter",
			setupFile: func(t *testing.T, path string) {
				content := `This is just markdown content.`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			},
			wantPrompt: &Prompt{
				Metadata: Metadata{},
				Content:  "This is just markdown content.",
			},
		},
		{
			name:        "file not found error",
			wantErr:     true,
			errContains: "failed to stat file",
		},
		{
			name: "reject non-markdown file",
			setupFile: func(t *testing.T, path string) {
				txtPath := strings.TrimSuffix(path, ".md") + ".txt"
				content := `---
title: "Test"
---

Test content`
				if err := os.WriteFile(txtPath, []byte(content), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			},
			wantErr:     true,
			errContains: "must have .md extension",
		},
		{
			name: "empty file",
			setupFile: func(t *testing.T, path string) {
				if err := os.WriteFile(path, []byte(""), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			},
			wantPrompt: &Prompt{
				Metadata: Metadata{},
				Content:  "",
			},
		},
		{
			name: "file with only frontmatter",
			setupFile: func(t *testing.T, path string) {
				content := `---
title: "Only Frontmatter"
tags:
  - tag1
---`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			},
			wantPrompt: &Prompt{
				Metadata: Metadata{
					Title: "Only Frontmatter",
					Tags:  []string{"tag1"},
				},
				Content: "",
			},
		},
		{
			name: "file size exceeds limit",
			setupFile: func(t *testing.T, path string) {
				content := make([]byte, MaxFileSize+1)
				if err := os.WriteFile(path, content, 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			},
			wantErr:     true,
			errContains: "exceeds maximum size",
		},
		{
			name: "file at size limit boundary",
			setupFile: func(t *testing.T, path string) {
				content := make([]byte, MaxFileSize)
				if err := os.WriteFile(path, content, 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			},
			wantPrompt: &Prompt{
				Metadata: Metadata{},
				Content:  string(make([]byte, MaxFileSize)),
			},
		},
		{
			name: "unicode content",
			setupFile: func(t *testing.T, path string) {
				content := `---
title: "Test 日本語 🇯🇵"
---

Content with unicode: 你好世界 and émojis 😀`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			},
			wantPrompt: &Prompt{
				Metadata: Metadata{
					Title: "Test 日本語 🇯🇵",
				},
				Content: "Content with unicode: 你好世界 and émojis 😀",
			},
		},
		{
			name: "special characters",
			setupFile: func(t *testing.T, path string) {
				content := `---
title: "Test Special Chars <>&\"'"
---

Content with: <tag> &amp; "quotes" 'apostrophes'`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			},
			wantPrompt: &Prompt{
				Metadata: Metadata{
					Title: `Test Special Chars <>&"'`,
				},
				Content: `Content with: <tag> &amp; "quotes" 'apostrophes'`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testPath := filepath.Join(tempDir, "test.md")

			if tt.name == "reject non-markdown file" {
				testPath = strings.TrimSuffix(testPath, ".md") + ".txt"
			}

			if tt.setupFile != nil {
				tt.setupFile(t, testPath)
			}

			gotPrompt, err := storage.Load(testPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Storage.Load() error should contain %q, got %q", tt.errContains, err.Error())
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Storage.Load() unexpected error: %v", err)
				return
			}

			if gotPrompt.Metadata.Title != tt.wantPrompt.Metadata.Title {
				t.Errorf("Storage.Load() Title = %q, want %q", gotPrompt.Metadata.Title, tt.wantPrompt.Metadata.Title)
			}

			if gotPrompt.Metadata.Category != tt.wantPrompt.Metadata.Category {
				t.Errorf("Storage.Load() Category = %q, want %q", gotPrompt.Metadata.Category, tt.wantPrompt.Metadata.Category)
			}

			if gotPrompt.Metadata.Description != tt.wantPrompt.Metadata.Description {
				t.Errorf("Storage.Load() Description = %q, want %q", gotPrompt.Metadata.Description, tt.wantPrompt.Metadata.Description)
			}

			if !stringSlicesEqual(gotPrompt.Metadata.Tags, tt.wantPrompt.Metadata.Tags) {
				t.Errorf("Storage.Load() Tags = %v, want %v", gotPrompt.Metadata.Tags, tt.wantPrompt.Metadata.Tags)
			}

			if gotPrompt.Content != tt.wantPrompt.Content {
				t.Errorf("Storage.Load() Content = %q, want %q", gotPrompt.Content, tt.wantPrompt.Content)
			}
		})
	}
}

func TestStorageLoadEdgeCases(t *testing.T) {
	logger := zap.NewNop()
	storage := NewStorage(logger)

	t.Run("follow symlinks", func(t *testing.T) {
		tempDir := t.TempDir()
		targetPath := filepath.Join(tempDir, "target.md")
		symlinkPath := filepath.Join(tempDir, "symlink.md")

		content := `---
title: "Symlink Test"
---

Content from symlink`
		if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create target file: %v", err)
		}

		if err := os.Symlink(targetPath, symlinkPath); err != nil {
			t.Skip("symlinks not supported on this system")
		}

		prompt, err := storage.Load(symlinkPath)
		if err != nil {
			t.Errorf("Storage.Load() failed to load symlink: %v", err)
		}

		if prompt.Metadata.Title != "Symlink Test" {
			t.Errorf("Storage.Load() Title = %q, want %q", prompt.Metadata.Title, "Symlink Test")
		}
	})

	t.Run("file with BOM", func(t *testing.T) {
		tempDir := t.TempDir()
		testPath := filepath.Join(tempDir, "test.md")

		content := "\xEF\xBB\xBF" + `---
title: "BOM Test"
---

Content after BOM`
		if err := os.WriteFile(testPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		prompt, err := storage.Load(testPath)
		if err != nil {
			t.Errorf("Storage.Load() failed to load file with BOM: %v", err)
		}

		if prompt.Metadata.Title != "BOM Test" {
			t.Errorf("Storage.Load() should handle BOM correctly, got Title = %q", prompt.Metadata.Title)
		}
	})

	t.Run("different line endings - CRLF", func(t *testing.T) {
		tempDir := t.TempDir()
		testPath := filepath.Join(tempDir, "test.md")

		content := "---\r\ntitle: \"CRLF Test\"\r\n---\r\n\r\nContent with CRLF\r\nline endings"
		if err := os.WriteFile(testPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		prompt, err := storage.Load(testPath)
		if err != nil {
			t.Errorf("Storage.Load() failed with CRLF: %v", err)
		}

		if prompt.Metadata.Title != "CRLF Test" {
			t.Errorf("Storage.Load() Title = %q, want %q", prompt.Metadata.Title, "CRLF Test")
		}

		if !strings.Contains(prompt.Content, "Content with CRLF") {
			t.Errorf("Storage.Load() should preserve CRLF content")
		}
	})

	t.Run("mixed line endings", func(t *testing.T) {
		tempDir := t.TempDir()
		testPath := filepath.Join(tempDir, "test.md")

		content := "---\ntitle: \"Mixed Test\"\r\n---\n\nContent with mixed\r\nline\nendings"
		if err := os.WriteFile(testPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		prompt, err := storage.Load(testPath)
		if err != nil {
			t.Errorf("Storage.Load() failed with mixed line endings: %v", err)
		}

		if prompt.Metadata.Title != "Mixed Test" {
			t.Errorf("Storage.Load() Title = %q, want %q", prompt.Metadata.Title, "Mixed Test")
		}
	})
}

func TestStorageSave(t *testing.T) {
	logger := zap.NewNop()
	storage := NewStorage(logger)

	tests := []struct {
		name        string
		prompt      *Prompt
		setupFile   func(t *testing.T, path string)
		wantErr     bool
		errContains string
		verify      func(t *testing.T, path string)
	}{
		{
			name: "special_characters",
			prompt: &Prompt{
				Metadata: Metadata{
					Title:       "Special Characters",
					Tags:        []string{"test"},
					Category:    "testing",
					Description: "Content with special: <tag> &amp; \"quotes\" 'apostrophes'",
				},
				Content: "Special chars: <>&\"'",
			},
			verify: func(t *testing.T, path string) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("failed to read saved file: %v", err)
				}

				content := string(data)
				if !strings.Contains(content, "<>&") {
					t.Errorf("saved file should preserve special characters")
				}
			},
		},
		{
			name:        "reject non-markdown file",
			prompt:      &Prompt{Metadata: Metadata{Title: "Test"}, Content: "Content"},
			wantErr:     true,
			errContains: "must have .md extension",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testPath := filepath.Join(tempDir, "subdir", "test.md")

			if tt.name == "reject non-markdown file" {
				testPath = strings.TrimSuffix(testPath, ".md") + ".txt"
			}

			if tt.setupFile != nil {
				tt.setupFile(t, testPath)
			}

			err := storage.Save(testPath, tt.prompt)

			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Storage.Save() error should contain %q, got %q", tt.errContains, err.Error())
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Storage.Save() unexpected error: %v", err)
				return
			}

			if tt.verify != nil {
				tt.verify(t, testPath)
			}
		})
	}
}

func TestStorageSaveEdgeCases(t *testing.T) {
	logger := zap.NewNop()
	storage := NewStorage(logger)

	t.Run("nested directory creation", func(t *testing.T) {
		tempDir := t.TempDir()
		testPath := filepath.Join(tempDir, "a", "b", "c", "test.md")

		prompt := &Prompt{
			Metadata: Metadata{Title: "Nested"},
			Content:  "Content",
		}

		if err := storage.Save(testPath, prompt); err != nil {
			t.Errorf("Storage.Save() failed with nested directories: %v", err)
		}

		if _, err := os.Stat(testPath); err != nil {
			t.Errorf("nested directories should be created: %v", err)
		}
	})

	t.Run("very long content at limit boundary", func(t *testing.T) {
		tempDir := t.TempDir()
		testPath := filepath.Join(tempDir, "long.md")

		content := ""
		for i := 0; i < 10000; i++ {
			content += "Line " + string(rune('0'+i%10)) + "\n"
		}

		prompt := &Prompt{
			Metadata: Metadata{Title: "Long Content"},
			Content:  content,
		}

		if err := storage.Save(testPath, prompt); err != nil {
			t.Errorf("Storage.Save() failed with long content: %v", err)
		}

		loaded, err := storage.Load(testPath)
		if err != nil {
			t.Errorf("Storage.Load() failed after saving long content: %v", err)
		}

		if loaded.Content != content {
			t.Errorf("long content should be preserved")
		}
	})

	t.Run("content with different line endings", func(t *testing.T) {
		tempDir := t.TempDir()
		testPath := filepath.Join(tempDir, "lineendings.md")

		prompt := &Prompt{
			Metadata: Metadata{Title: "Line Endings"},
			Content:  "Line 1\r\nLine 2\nLine 3\r\nLine 4",
		}

		if err := storage.Save(testPath, prompt); err != nil {
			t.Errorf("Storage.Save() failed: %v", err)
		}

		data, err := os.ReadFile(testPath)
		if err != nil {
			t.Fatalf("failed to read saved file: %v", err)
		}

		if !strings.Contains(string(data), "Line 1") {
			t.Errorf("saved file should preserve content with mixed line endings")
		}
	})
}

func TestStorageRoundTrip(t *testing.T) {
	logger := zap.NewNop()
	storage := NewStorage(logger)

	tests := []struct {
		name   string
		prompt *Prompt
	}{
		{
			name: "full metadata",
			prompt: &Prompt{
				Metadata: Metadata{
					Title:       "Test Prompt",
					Tags:        []string{"test", "example"},
					Category:    "testing",
					Description: "A test prompt",
					Author:      "Test Author",
					Version:     "1.0.0",
				},
				Content: "This is the content.\nWith multiple lines.",
			},
		},
		{
			name: "minimal metadata",
			prompt: &Prompt{
				Metadata: Metadata{
					Title: "Simple",
				},
				Content: "Simple content",
			},
		},
		{
			name: "empty metadata",
			prompt: &Prompt{
				Metadata: Metadata{},
				Content:  "Just content",
			},
		},
		{
			name: "unicode",
			prompt: &Prompt{
				Metadata: Metadata{
					Title: "日本語 🇯🇵",
				},
				Content: "你好世界\némojis 😀",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testPath := filepath.Join(tempDir, "test.md")

			// Save
			if err := storage.Save(testPath, tt.prompt); err != nil {
				t.Fatalf("Storage.Save() failed: %v", err)
			}

			// Load
			loaded, err := storage.Load(testPath)
			if err != nil {
				t.Fatalf("Storage.Load() failed: %v", err)
			}

			// Verify
			if loaded.Metadata.Title != tt.prompt.Metadata.Title {
				t.Errorf("round trip Title changed: %q -> %q", tt.prompt.Metadata.Title, loaded.Metadata.Title)
			}

			if loaded.Metadata.Category != tt.prompt.Metadata.Category {
				t.Errorf("round trip Category changed: %q -> %q", tt.prompt.Metadata.Category, loaded.Metadata.Category)
			}

			if loaded.Metadata.Description != tt.prompt.Metadata.Description {
				t.Errorf("round trip Description changed: %q -> %q", tt.prompt.Metadata.Description, loaded.Metadata.Description)
			}

			if loaded.Metadata.Author != tt.prompt.Metadata.Author {
				t.Errorf("round trip Author changed: %q -> %q", tt.prompt.Metadata.Author, loaded.Metadata.Author)
			}

			if loaded.Metadata.Version != tt.prompt.Metadata.Version {
				t.Errorf("round trip Version changed: %q -> %q", tt.prompt.Metadata.Version, loaded.Metadata.Version)
			}

			if !stringSlicesEqual(loaded.Metadata.Tags, tt.prompt.Metadata.Tags) {
				t.Errorf("round trip Tags changed: %v -> %v", tt.prompt.Metadata.Tags, loaded.Metadata.Tags)
			}

			if loaded.Content != tt.prompt.Content {
				t.Errorf("round trip Content changed: %q -> %q", tt.prompt.Content, loaded.Content)
			}
		})
	}
}
