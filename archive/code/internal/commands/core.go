package commands

import (
	"fmt"
)

// RegisterCoreCommands registers all core application commands
func RegisterCoreCommands(registry *Registry) error {
	// General Commands
	if err := registry.Register(&Command{
		ID:          "toggle-ai",
		Name:        "Toggle AI Panel",
		Description: "Show or hide the AI suggestions panel",
		Category:    "General",
		Handler: func() error {
			// TODO: Implement AI panel toggle
			return fmt.Errorf("AI panel toggle not yet implemented")
		},
	}); err != nil {
		return err
	}

	if err := registry.Register(&Command{
		ID:          "copy",
		Name:        "Copy to Clipboard",
		Description: "Copy current composition to clipboard",
		Category:    "General",
		Handler: func() error {
			// TODO: Implement copy to clipboard
			return fmt.Errorf("copy to clipboard not yet implemented")
		},
	}); err != nil {
		return err
	}

	if err := registry.Register(&Command{
		ID:          "save",
		Name:        "Save Composition",
		Description: "Manually save the current composition",
		Category:    "General",
		Handler: func() error {
			// TODO: Implement manual save
			return fmt.Errorf("manual save not yet implemented")
		},
	}); err != nil {
		return err
	}

	// AI Commands
	if err := registry.Register(&Command{
		ID:          "ai-suggestions",
		Name:        "Get AI Suggestions",
		Description: "Request AI suggestions for current composition",
		Category:    "AI",
		Handler: func() error {
			// This command will be handled by the TUI to trigger AI suggestions
			// The handler returns nil to indicate success, and the TUI will
			// handle sending the TriggerAISuggestionsMsg
			return nil
		},
	}); err != nil {
		return err
	}

	// File Commands
	if err := registry.Register(&Command{
		ID:          "add-file-ref",
		Name:        "Add File Reference",
		Description: "Insert a reference to a file in the working directory",
		Category:    "Files",
		Handler: func() error {
			// TODO: Implement file reference
			return fmt.Errorf("file reference not yet implemented")
		},
	}); err != nil {
		return err
	}

	// Prompt Commands
	if err := registry.Register(&Command{
		ID:          "create-prompt",
		Name:        "Create New Prompt",
		Description: "Create a new prompt template",
		Category:    "Prompts",
		Handler: func() error {
			// TODO: Implement prompt creation
			return fmt.Errorf("prompt creation not yet implemented")
		},
	}); err != nil {
		return err
	}

	if err := registry.Register(&Command{
		ID:          "edit-prompt",
		Name:        "Edit Prompt",
		Description: "Edit an existing prompt template",
		Category:    "Prompts",
		Handler: func() error {
			// TODO: Implement prompt editing
			return fmt.Errorf("prompt editing not yet implemented")
		},
	}); err != nil {
		return err
	}

	if err := registry.Register(&Command{
		ID:          "validate-library",
		Name:        "Validate Library",
		Description: "Run validation checks on all prompts in the library",
		Category:    "Prompts",
		Handler: func() error {
			// This command will be handled by the TUI to show validation modal
			// The handler returns nil to indicate success, and the TUI will
			// handle showing the validation results
			return nil
		},
	}); err != nil {
		return err
	}

	// History Commands
	if err := registry.Register(&Command{
		ID:          "view-history",
		Name:        "View History",
		Description: "Browse and load previous compositions",
		Category:    "History",
		Handler: func() error {
			// TODO: Implement history browser
			return fmt.Errorf("history browser not yet implemented")
		},
	}); err != nil {
		return err
	}

	if err := registry.Register(&Command{
		ID:          "cleanup-history",
		Name:        "Clean Up History",
		Description: "Remove old or unwanted compositions from history",
		Category:    "History",
		Handler: func() error {
			// TODO: Implement history cleanup
			return fmt.Errorf("history cleanup not yet implemented")
		},
	}); err != nil {
		return err
	}

	// Settings Commands
	if err := registry.Register(&Command{
		ID:          "settings",
		Name:        "Settings",
		Description: "Open settings panel",
		Category:    "Settings",
		Handler: func() error {
			// TODO: Implement settings panel
			return fmt.Errorf("settings panel not yet implemented")
		},
	}); err != nil {
		return err
	}

	// Debug Commands
	if err := registry.Register(&Command{
		ID:          "view-logs",
		Name:        "View Logs",
		Description: "Open log viewer to see recent application logs",
		Category:    "Debug",
		Handler: func() error {
			// TODO: Implement log viewer
			return fmt.Errorf("log viewer not yet implemented")
		},
	}); err != nil {
		return err
	}

	return nil
}
