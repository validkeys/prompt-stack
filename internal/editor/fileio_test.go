package editor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileManager(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantPath string
	}{
		{
			name:     "with path",
			path:     "/test/file.txt",
			wantPath: "/test/file.txt",
		},
		{
			name:     "empty path",
			path:     "",
			wantPath: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm := NewFileManager(tt.path)
			if fm.Path() != tt.wantPath {
				t.Errorf("NewFileManager() path = %v, want %v", fm.Path(), tt.wantPath)
			}
			if fm.IsModified() {
				t.Errorf("NewFileManager() should not be modified initially")
			}
		})
	}
}

func TestFileManager_Load(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		setup   func() FileManager
		want    string
		wantErr bool
	}{
		{
			name: "load existing file",
			setup: func() FileManager {
				path := filepath.Join(tmpDir, "test.txt")
				content := "Hello, World!"
				os.WriteFile(path, []byte(content), 0644)
				return NewFileManager(path)
			},
			want:    "Hello, World!",
			wantErr: false,
		},
		{
			name: "load empty file",
			setup: func() FileManager {
				path := filepath.Join(tmpDir, "empty.txt")
				os.WriteFile(path, []byte(""), 0644)
				return NewFileManager(path)
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "load non-existent file",
			setup: func() FileManager {
				path := filepath.Join(tmpDir, "nonexistent.txt")
				return NewFileManager(path)
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "load with empty path",
			setup: func() FileManager {
				return NewFileManager("")
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "load file with newlines",
			setup: func() FileManager {
				path := filepath.Join(tmpDir, "newlines.txt")
				content := "Line 1\nLine 2\nLine 3"
				os.WriteFile(path, []byte(content), 0644)
				return NewFileManager(path)
			},
			want:    "Line 1\nLine 2\nLine 3",
			wantErr: false,
		},
		{
			name: "load file with unicode",
			setup: func() FileManager {
				path := filepath.Join(tmpDir, "unicode.txt")
				content := "Hello 世界 🌍"
				os.WriteFile(path, []byte(content), 0644)
				return NewFileManager(path)
			},
			want:    "Hello 世界 🌍",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm := tt.setup()
			got, err := fm.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("FileManager.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FileManager.Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileManager_Save(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		setup   func() FileManager
		content string
		wantErr bool
		verify  func(path string) error
	}{
		{
			name: "save to new file",
			setup: func() FileManager {
				path := filepath.Join(tmpDir, "new.txt")
				return NewFileManager(path)
			},
			content: "New content",
			wantErr: false,
			verify: func(path string) error {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				if string(content) != "New content" {
					t.Errorf("saved content = %v, want %v", string(content), "New content")
				}
				return nil
			},
		},
		{
			name: "save to existing file",
			setup: func() FileManager {
				path := filepath.Join(tmpDir, "existing.txt")
				os.WriteFile(path, []byte("Old content"), 0644)
				return NewFileManager(path)
			},
			content: "Updated content",
			wantErr: false,
			verify: func(path string) error {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				if string(content) != "Updated content" {
					t.Errorf("saved content = %v, want %v", string(content), "Updated content")
				}
				return nil
			},
		},
		{
			name: "save empty content",
			setup: func() FileManager {
				path := filepath.Join(tmpDir, "empty.txt")
				return NewFileManager(path)
			},
			content: "",
			wantErr: false,
			verify: func(path string) error {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				if string(content) != "" {
					t.Errorf("saved content = %v, want empty", string(content))
				}
				return nil
			},
		},
		{
			name: "save with empty path",
			setup: func() FileManager {
				return NewFileManager("")
			},
			content: "Should not save",
			wantErr: false,
			verify:  nil,
		},
		{
			name: "save to nested directory",
			setup: func() FileManager {
				path := filepath.Join(tmpDir, "nested", "dir", "file.txt")
				return NewFileManager(path)
			},
			content: "Nested content",
			wantErr: false,
			verify: func(path string) error {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				if string(content) != "Nested content" {
					t.Errorf("saved content = %v, want %v", string(content), "Nested content")
				}
				return nil
			},
		},
		{
			name: "save with unicode",
			setup: func() FileManager {
				path := filepath.Join(tmpDir, "unicode.txt")
				return NewFileManager(path)
			},
			content: "Hello 世界 🌍",
			wantErr: false,
			verify: func(path string) error {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				if string(content) != "Hello 世界 🌍" {
					t.Errorf("saved content = %v, want %v", string(content), "Hello 世界 🌍")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm := tt.setup()
			err := fm.Save(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileManager.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.verify != nil && fm.Path() != "" {
				if err := tt.verify(fm.Path()); err != nil {
					t.Errorf("verification failed: %v", err)
				}
			}
		})
	}
}

func TestFileManager_MarkModified(t *testing.T) {
	fm := NewFileManager("/test/file.txt")
	if fm.IsModified() {
		t.Errorf("NewFileManager should not be modified initially")
	}

	fm = fm.MarkModified()
	if !fm.IsModified() {
		t.Errorf("MarkModified() should set modified flag to true")
	}

	// Verify immutability - original should still be unmodified
	original := NewFileManager("/test/file.txt")
	if original.IsModified() {
		t.Errorf("Original FileManager should remain unmodified")
	}
}

func TestFileManager_ClearModified(t *testing.T) {
	fm := NewFileManager("/test/file.txt")
	fm = fm.MarkModified()
	if !fm.IsModified() {
		t.Errorf("MarkModified() should set modified flag to true")
	}

	fm = fm.ClearModified()
	if fm.IsModified() {
		t.Errorf("ClearModified() should set modified flag to false")
	}
}

func TestFileManager_SetPath(t *testing.T) {
	fm := NewFileManager("/old/path.txt")
	if fm.Path() != "/old/path.txt" {
		t.Errorf("initial path = %v, want /old/path.txt", fm.Path())
	}

	fm = fm.SetPath("/new/path.txt")
	if fm.Path() != "/new/path.txt" {
		t.Errorf("SetPath() path = %v, want /new/path.txt", fm.Path())
	}

	// Verify that setting path clears modified flag
	fm = fm.MarkModified()
	fm = fm.SetPath("/another/path.txt")
	if fm.IsModified() {
		t.Errorf("SetPath() should clear modified flag")
	}
}

func TestFileManager_HasPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "with path",
			path: "/test/file.txt",
			want: true,
		},
		{
			name: "empty path",
			path: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm := NewFileManager(tt.path)
			if got := fm.HasPath(); got != tt.want {
				t.Errorf("FileManager.HasPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileManager_Immutability(t *testing.T) {
	// Test that all methods return new instances
	original := NewFileManager("/test/file.txt")

	// MarkModified
	modified := original.MarkModified()
	if original.IsModified() {
		t.Errorf("original should not be modified after MarkModified()")
	}
	if !modified.IsModified() {
		t.Errorf("returned instance should be modified")
	}

	// ClearModified
	cleared := modified.ClearModified()
	if !modified.IsModified() {
		t.Errorf("modified should still be modified after ClearModified()")
	}
	if cleared.IsModified() {
		t.Errorf("returned instance should not be modified")
	}

	// SetPath
	newPath := original.SetPath("/new/path.txt")
	if original.Path() != "/test/file.txt" {
		t.Errorf("original path should not change after SetPath()")
	}
	if newPath.Path() != "/new/path.txt" {
		t.Errorf("returned instance should have new path")
	}
}

func TestFileManager_Workflow(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "workflow.txt")

	// Create new file manager
	fm := NewFileManager(path)
	if !fm.HasPath() {
		t.Errorf("should have path")
	}
	if fm.IsModified() {
		t.Errorf("should not be modified initially")
	}

	// Save content
	content := "Initial content"
	err := fm.Save(content)
	if err != nil {
		t.Errorf("Save() error = %v", err)
	}

	// Mark as modified
	fm = fm.MarkModified()
	if !fm.IsModified() {
		t.Errorf("should be modified after MarkModified()")
	}

	// Load content
	loaded, err := fm.Load()
	if err != nil {
		t.Errorf("Load() error = %v", err)
	}
	if loaded != content {
		t.Errorf("loaded content = %v, want %v", loaded, content)
	}

	// Clear modified flag
	fm = fm.ClearModified()
	if fm.IsModified() {
		t.Errorf("should not be modified after ClearModified()")
	}

	// Update content
	newContent := "Updated content"
	err = fm.Save(newContent)
	if err != nil {
		t.Errorf("Save() error = %v", err)
	}

	// Verify file was updated
	loaded, err = fm.Load()
	if err != nil {
		t.Errorf("Load() error = %v", err)
	}
	if loaded != newContent {
		t.Errorf("loaded content = %v, want %v", loaded, newContent)
	}
}

func TestFileManager_EdgeCases(t *testing.T) {
	t.Run("very long content", func(t *testing.T) {
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "long.txt")
		fm := NewFileManager(path)

		longContent := string(make([]byte, 10000))
		for i := range longContent {
			longContent = longContent[:i] + "a" + longContent[i+1:]
		}

		err := fm.Save(longContent)
		if err != nil {
			t.Errorf("Save() error = %v", err)
		}

		loaded, err := fm.Load()
		if err != nil {
			t.Errorf("Load() error = %v", err)
		}
		if loaded != longContent {
			t.Errorf("loaded content length = %v, want %v", len(loaded), len(longContent))
		}
	})

	t.Run("special characters in path", func(t *testing.T) {
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "file with spaces.txt")
		fm := NewFileManager(path)

		content := "Special path content"
		err := fm.Save(content)
		if err != nil {
			t.Errorf("Save() error = %v", err)
		}

		loaded, err := fm.Load()
		if err != nil {
			t.Errorf("Load() error = %v", err)
		}
		if loaded != content {
			t.Errorf("loaded content = %v, want %v", loaded, content)
		}
	})
}
