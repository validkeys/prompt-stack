package ai

import (
	"fmt"
	"strings"
)

// DiffGenerator handles diff generation and application
type DiffGenerator struct {
	// Configuration
	contextLines int // Number of context lines to show in diff
}

// NewDiffGenerator creates a new diff generator
func NewDiffGenerator() *DiffGenerator {
	return &DiffGenerator{
		contextLines: 3, // Default: 3 lines of context
	}
}

// UnifiedDiff represents a unified diff format
type UnifiedDiff struct {
	// Header contains diff header information
	Header string

	// Hunks are the individual diff hunks
	Hunks []DiffHunk
}

// DiffHunk represents a single hunk in a unified diff
type DiffHunk struct {
	// OldStart is the starting line in the old content
	OldStart int

	// OldLines is the number of lines in the old content
	OldLines int

	// NewStart is the starting line in the new content
	NewStart int

	// NewLines is the number of lines in the new content
	NewLines int

	// Lines are the diff lines
	Lines []DiffLine
}

// DiffLine represents a single line in a diff
type DiffLine struct {
	// Type is the type of diff line
	Type DiffLineType

	// Content is the line content
	Content string

	// LineNumber is the line number (for display)
	LineNumber int
}

// DiffLineType represents the type of diff line
type DiffLineType string

const (
	// DiffLineContext is a context line (unchanged)
	DiffLineContext DiffLineType = "context"

	// DiffLineAddition is an added line
	DiffLineAddition DiffLineType = "addition"

	// DiffLineDeletion is a deleted line
	DiffLineDeletion DiffLineType = "deletion"
)

// GenerateUnifiedDiff generates a unified diff from edits
func (dg *DiffGenerator) GenerateUnifiedDiff(original string, edits []Edit) (*UnifiedDiff, error) {
	if original == "" {
		return nil, fmt.Errorf("original content is empty")
	}

	if len(edits) == 0 {
		return nil, fmt.Errorf("no edits to apply")
	}

	// Split original content into lines
	originalLines := strings.Split(original, "\n")

	// Apply edits to create new content
	newContent, err := dg.applyEdits(original, edits)
	if err != nil {
		return nil, fmt.Errorf("failed to apply edits: %w", err)
	}

	// Split new content into lines
	newLines := strings.Split(newContent, "\n")

	// Generate diff hunks
	hunks := dg.generateHunks(originalLines, newLines)

	return &UnifiedDiff{
		Header: dg.generateHeader(original, newContent),
		Hunks:  hunks,
	}, nil
}

// applyEdits applies edits to original content
func (dg *DiffGenerator) applyEdits(original string, edits []Edit) (string, error) {
	// Sort edits by position (line, column) to apply in order
	sortedEdits := make([]Edit, len(edits))
	copy(sortedEdits, edits)

	// Simple bubble sort by line, then column
	for i := 0; i < len(sortedEdits); i++ {
		for j := i + 1; j < len(sortedEdits); j++ {
			if sortedEdits[j].Line < sortedEdits[i].Line ||
				(sortedEdits[j].Line == sortedEdits[i].Line && sortedEdits[j].Column < sortedEdits[i].Column) {
				sortedEdits[i], sortedEdits[j] = sortedEdits[j], sortedEdits[i]
			}
		}
	}

	// Apply edits in order
	result := original
	offset := 0 // Track offset due to previous edits

	for _, edit := range sortedEdits {
		// Convert line/column to character position
		pos := dg.lineColumnToPosition(result, edit.Line, edit.Column)

		// Adjust position based on previous edits
		pos += offset

		// Validate position
		if pos < 0 || pos > len(result) {
			return "", fmt.Errorf("edit position out of bounds: line %d, column %d", edit.Line, edit.Column)
		}

		// Validate old content matches
		if pos+edit.Length > len(result) {
			return "", fmt.Errorf("edit length exceeds content: line %d, column %d", edit.Line, edit.Column)
		}

		oldContent := result[pos : pos+edit.Length]
		if oldContent != edit.OldContent {
			return "", fmt.Errorf("old content mismatch at line %d, column %d: expected %q, got %q",
				edit.Line, edit.Column, edit.OldContent, oldContent)
		}

		// Apply edit
		result = result[:pos] + edit.NewContent + result[pos+edit.Length:]

		// Update offset for next edit
		offset += len(edit.NewContent) - edit.Length
	}

	return result, nil
}

// lineColumnToPosition converts line and column to character position
func (dg *DiffGenerator) lineColumnToPosition(content string, line, column int) int {
	if line < 1 || column < 1 {
		return 0
	}

	lines := strings.Split(content, "\n")

	// Sum up lengths of previous lines
	pos := 0
	for i := 0; i < line-1 && i < len(lines); i++ {
		pos += len(lines[i]) + 1 // +1 for newline
	}

	// Add column offset (column is 1-indexed)
	pos += column - 1

	return pos
}

// generateHunks generates diff hunks from old and new lines
func (dg *DiffGenerator) generateHunks(oldLines, newLines []string) []DiffHunk {
	var hunks []DiffHunk

	// Simple line-by-line comparison
	// In a real implementation, you'd use a proper diff algorithm
	// For now, we'll do a basic comparison

	oldLine := 1
	newLine := 1

	for oldLine <= len(oldLines) || newLine <= len(newLines) {
		var hunk DiffHunk
		var lines []DiffLine

		// Find next difference
		for oldLine <= len(oldLines) && newLine <= len(newLines) {
			if oldLine > len(oldLines) {
				// Only additions remain
				lines = append(lines, DiffLine{
					Type:       DiffLineAddition,
					Content:    newLines[newLine-1],
					LineNumber: newLine,
				})
				newLine++
				break
			}

			if newLine > len(newLines) {
				// Only deletions remain
				lines = append(lines, DiffLine{
					Type:       DiffLineDeletion,
					Content:    oldLines[oldLine-1],
					LineNumber: oldLine,
				})
				oldLine++
				break
			}

			if oldLines[oldLine-1] == newLines[newLine-1] {
				// Lines match
				lines = append(lines, DiffLine{
					Type:       DiffLineContext,
					Content:    oldLines[oldLine-1],
					LineNumber: oldLine,
				})
				oldLine++
				newLine++
			} else {
				// Lines differ
				lines = append(lines, DiffLine{
					Type:       DiffLineDeletion,
					Content:    oldLines[oldLine-1],
					LineNumber: oldLine,
				})
				oldLine++

				if newLine <= len(newLines) {
					lines = append(lines, DiffLine{
						Type:       DiffLineAddition,
						Content:    newLines[newLine-1],
						LineNumber: newLine,
					})
					newLine++
				}
				break
			}
		}

		if len(lines) > 0 {
			hunk = DiffHunk{
				OldStart: oldLine - len(lines),
				OldLines: 0,
				NewStart: newLine - len(lines),
				NewLines: 0,
				Lines:    lines,
			}

			// Count old and new lines in hunk
			for _, line := range lines {
				if line.Type == DiffLineDeletion {
					hunk.OldLines++
				} else if line.Type == DiffLineAddition {
					hunk.NewLines++
				} else {
					hunk.OldLines++
					hunk.NewLines++
				}
			}

			hunks = append(hunks, hunk)
		}
	}

	return hunks
}

// generateHeader generates a diff header
func (dg *DiffGenerator) generateHeader(original, newContent string) string {
	originalLines := strings.Count(original, "\n") + 1
	newLines := strings.Count(newContent, "\n") + 1

	return fmt.Sprintf("--- original (%d lines)\n+++ new (%d lines)", originalLines, newLines)
}

// FormatUnifiedDiff formats a unified diff as a string
func (dg *DiffGenerator) FormatUnifiedDiff(diff *UnifiedDiff) string {
	var builder strings.Builder

	// Write header
	builder.WriteString(diff.Header)
	builder.WriteString("\n\n")

	// Write hunks
	for _, hunk := range diff.Hunks {
		// Write hunk header
		builder.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n",
			hunk.OldStart, hunk.OldLines, hunk.NewStart, hunk.NewLines))

		// Write lines
		for _, line := range hunk.Lines {
			prefix := " "
			switch line.Type {
			case DiffLineAddition:
				prefix = "+"
			case DiffLineDeletion:
				prefix = "-"
			}

			builder.WriteString(fmt.Sprintf("%s%s\n", prefix, line.Content))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

// ApplyEdits applies edits to content and returns new content
func (dg *DiffGenerator) ApplyEdits(original string, edits []Edit) (string, error) {
	return dg.applyEdits(original, edits)
}

// ValidateEdits validates that edits can be applied to content
func (dg *DiffGenerator) ValidateEdits(content string, edits []Edit) error {
	// Check each edit
	for _, edit := range edits {
		// Convert line/column to position
		pos := dg.lineColumnToPosition(content, edit.Line, edit.Column)

		// Validate position
		if pos < 0 || pos > len(content) {
			return fmt.Errorf("edit position out of bounds: line %d, column %d", edit.Line, edit.Column)
		}

		// Validate length
		if pos+edit.Length > len(content) {
			return fmt.Errorf("edit length exceeds content: line %d, column %d", edit.Line, edit.Column)
		}

		// Validate old content matches
		oldContent := content[pos : pos+edit.Length]
		if oldContent != edit.OldContent {
			return fmt.Errorf("old content mismatch at line %d, column %d: expected %q, got %q",
				edit.Line, edit.Column, edit.OldContent, oldContent)
		}
	}

	return nil
}
