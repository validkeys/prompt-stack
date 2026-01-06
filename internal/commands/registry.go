package commands

import (
	"fmt"
	"strings"

	"github.com/sahilm/fuzzy"
)

// Command represents a command that can be executed
type Command struct {
	ID          string
	Name        string
	Description string
	Category    string
	Handler     func() error
}

// Registry manages all available commands
type Registry struct {
	commands map[string]*Command
}

// NewRegistry creates a new command registry
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]*Command),
	}
}

// Register adds a command to the registry
func (r *Registry) Register(cmd *Command) error {
	if cmd.ID == "" {
		return fmt.Errorf("command ID cannot be empty")
	}
	if cmd.Name == "" {
		return fmt.Errorf("command name cannot be empty")
	}
	if cmd.Handler == nil {
		return fmt.Errorf("command handler cannot be nil")
	}

	if _, exists := r.commands[cmd.ID]; exists {
		return fmt.Errorf("command with ID %s already registered", cmd.ID)
	}

	r.commands[cmd.ID] = cmd
	return nil
}

// Get retrieves a command by ID
func (r *Registry) Get(id string) (*Command, bool) {
	cmd, exists := r.commands[id]
	return cmd, exists
}

// GetAll returns all commands
func (r *Registry) GetAll() []*Command {
	commands := make([]*Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		commands = append(commands, cmd)
	}
	return commands
}

// Search performs fuzzy search on commands
func (r *Registry) Search(query string) []*Command {
	if query == "" {
		return r.GetAll()
	}

	// Build searchable strings
	var stringsToMatch []string
	var commands []*Command
	for _, cmd := range r.commands {
		// Combine name, description, and category for search
		searchable := fmt.Sprintf("%s %s %s", cmd.Name, cmd.Description, cmd.Category)
		stringsToMatch = append(stringsToMatch, searchable)
		commands = append(commands, cmd)
	}

	// Apply fuzzy matching
	matches := fuzzy.Find(query, stringsToMatch)

	// Return matched commands
	result := make([]*Command, 0, len(matches))
	for _, match := range matches {
		result = append(result, commands[match.Index])
	}

	return result
}

// GetByCategory returns all commands in a category
func (r *Registry) GetByCategory(category string) []*Command {
	var commands []*Command
	for _, cmd := range r.commands {
		if cmd.Category == category {
			commands = append(commands, cmd)
		}
	}
	return commands
}

// GetCategories returns all unique categories
func (r *Registry) GetCategories() []string {
	categories := make(map[string]bool)
	for _, cmd := range r.commands {
		if cmd.Category != "" {
			categories[cmd.Category] = true
		}
	}

	result := make([]string, 0, len(categories))
	for cat := range categories {
		result = append(result, cat)
	}
	return result
}

// Count returns the total number of commands
func (r *Registry) Count() int {
	return len(r.commands)
}

// Unregister removes a command from the registry
func (r *Registry) Unregister(id string) bool {
	if _, exists := r.commands[id]; exists {
		delete(r.commands, id)
		return true
	}
	return false
}

// Clear removes all commands from the registry
func (r *Registry) Clear() {
	r.commands = make(map[string]*Command)
}

// String returns a string representation of the command
func (c *Command) String() string {
	var parts []string
	if c.Category != "" {
		parts = append(parts, fmt.Sprintf("[%s]", c.Category))
	}
	parts = append(parts, c.Name)
	if c.Description != "" {
		parts = append(parts, fmt.Sprintf("- %s", c.Description))
	}
	return strings.Join(parts, " ")
}
