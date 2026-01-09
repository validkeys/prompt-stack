// Package examples provides reference implementations demonstrating OpenCode design system patterns.
// These are NOT production code - they are design validation and developer reference material.
package examples

import (
	"fmt"
	"strings"

	"github.com/kyledavis/prompt-stack/ui/theme"
)

// ModalExample demonstrates the OpenCode modal dialog pattern.
func ModalExample() string {
	dialogStyle := theme.ModalStyle().
		Width(60).
		Height(10).
		Padding(theme.Unit*2, theme.Unit*3)

	titleStyle := theme.ModalTitleStyle()
	contentStyle := theme.ModalContentStyle()
	buttonStyle := theme.ModalButtonStyle()
	buttonFocusedStyle := theme.ModalButtonFocusedStyle()

	return dialogStyle.Render(
		titleStyle.Render("Select Theme") + "\n\n" +
			contentStyle.Render("  • opencode (current)\n    tokyonight\n    catppuccin\n    nord") + "\n\n" +
			buttonFocusedStyle.Render("[ OK ]") + " " +
			buttonStyle.Render("[ Cancel ]") + "\n\n" +
			theme.MutedStyle().Render("Press Enter to select, Esc to cancel"),
	)
}

// StatusBarExample demonstrates the OpenCode status bar pattern.
func StatusBarExample() string {
	return theme.StatusStyle().
		Width(80).
		Height(1).
		Padding(0, theme.Unit).
		Render("Chars: 125 | Lines: 8 | Modified | " + theme.InfoStyle().Render("[PLACEHOLDER EDIT]"))
}

// InputExample demonstrates the OpenCode input field pattern.
func InputExample() string {
	inputStyle := theme.InputStyle().
		Width(50).
		Height(3)

	return inputStyle.Render("Search files...")
}

// ListExample demonstrates the OpenCode list pattern.
func ListExample() string {
	itemStyle := theme.ListItemStyle()
	selectedStyle := theme.ListItemSelectedStyle()
	categoryStyle := theme.ListCategoryStyle()

	return theme.ModalStyle().
		Width(60).
		Height(12).
		Padding(theme.Unit, theme.Unit*2).
		Render(
			categoryStyle.Render("Files") + "\n" +
				selectedStyle.Render("> src/main.go") + "\n" +
				itemStyle.Render("  src/utils/file.go") + "\n" +
				itemStyle.Render("  internal/config/config.go") + "\n" +
				itemStyle.Render("  cmd/server/main.go") + "\n\n" +
				theme.MutedStyle().Render("Press Up/Down to navigate, Enter to select"),
		)
}

// PreviewExample demonstrates the OpenCode preview pane pattern.
func PreviewExample() string {
	previewStyle := theme.PreviewStyle().
		Width(60).
		Height(15)

	titleStyle := theme.PreviewTitleStyle()
	contentStyle := theme.ModalContentStyle()

	return previewStyle.Render(
		titleStyle.Render("Preview: src/main.go") + "\n\n" +
			contentStyle.Render("package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}") + "\n\n" +
			theme.MutedStyle().Render("Lines: 6 | Size: 85 bytes"),
	)
}

// DiffExample demonstrates the OpenCode diff display pattern.
func DiffExample() string {
	return theme.ModalStyle().
		Width(80).
		Height(10).
		Padding(theme.Unit, theme.Unit*2).
		Render(
			theme.DiffHeaderStyle().Render("@@ -1,5 +1,5 @@") + "\n" +
				theme.DiffRemovedStyle().Render("-func oldFunction() {") + "\n" +
				theme.DiffAddedStyle().Render("+func newFunction() {") + "\n" +
				theme.DiffContextStyle().Render("     // Old comment") + "\n" +
				theme.DiffRemovedStyle().Render("-    return oldValue") + "\n" +
				theme.DiffAddedStyle().Render("+    return newValue") + "\n" +
				theme.DiffContextStyle().Render(" }") + "\n\n" +
				theme.SuccessStyle().Render("Green:") + " Added lines    " +
				theme.ErrorStyle().Render("Red:") + " Removed lines    " +
				theme.MutedStyle().Render("Gray:") + " Context lines",
		)
}

// ChatExample demonstrates the OpenCode chat interface pattern.
func ChatExample() string {
	userStyle := theme.UserMessageStyle()
	assistantStyle := theme.AssistantMessageStyle()
	toolStyle := theme.ToolNameStyle()

	return theme.ModalStyle().
		Width(80).
		Height(15).
		Padding(theme.Unit, theme.Unit*2).
		Render(
			userStyle.Render("User:") + " Can you help me with this code?" + "\n\n" +
				assistantStyle.Render("Assistant:") + " Sure! I'd be happy to help. Let me look at the code..." + "\n\n" +
				toolStyle.Render("[Tool: bash]") + " Running: ls -la" + "\n" +
				theme.ToolOutputStyle().Render("total 48\ndrwxr-xr-x  11 user  staff   352 Jan  8 12:00 .\ndrwxr-xr-x   5 user  staff   160 Jan  8 11:30 ..\n-rw-r--r--   1 user  staff   287 Jan  8 11:45 .git") + "\n\n" +
				toolStyle.Render("[Tool: bash]") + " " + theme.SuccessStyle().Render("Completed successfully") + "\n\n" +
				assistantStyle.Render("Assistant:") + " I found the issue. Here's the fix..." + "\n\n" +
				theme.MutedStyle().Render("Press Ctrl+C to interrupt, Ctrl+G to abort"),
		)
}

// ValidationExample demonstrates validation error styling.
func ValidationExample() string {
	return theme.ModalStyle().
		Width(60).
		Height(8).
		Padding(theme.Unit, theme.Unit*2).
		Render(
			theme.ValidationErrorStyle().Render("Error: Invalid frontmatter") + "\n\n" +
				theme.MutedStyle().Render("Missing required field: 'name'") + "\n\n" +
				theme.WarningStyle().Render("Press Enter to fix, Esc to cancel"),
		)
}

// HighlightExample demonstrates text highlighting.
func HighlightExample() string {
	content := "The quick brown fox jumps over the lazy dog."
	highlighted := "quick brown fox"

	return theme.ModalStyle().
		Width(70).
		Height(6).
		Padding(theme.Unit, theme.Unit*2).
		Render(
			theme.HeaderStyle().Render("Search Results") + "\n\n" +
				strings.Replace(content, highlighted, theme.HighlightStyle().Render(highlighted), 1) + "\n\n" +
				theme.MutedStyle().Render("3 matches found in document"),
		)
}

// KeyboardExample demonstrates keyboard shortcut help.
func KeyboardExample() string {
	return theme.ModalStyle().
		Width(60).
		Height(10).
		Padding(theme.Unit, theme.Unit*2).
		Render(
			theme.HeaderStyle().Render("Keyboard Shortcuts") + "\n\n" +
				theme.InfoStyle().Render("Ctrl+X") + " - Command menu\n" +
				theme.InfoStyle().Render("Ctrl+P") + " - Command palette\n" +
				theme.InfoStyle().Render("Tab") + " - Switch agents\n" +
				theme.InfoStyle().Render("Esc") + " - Cancel/close\n" +
				theme.InfoStyle().Render("Ctrl+C") + " - Interrupt\n" +
				theme.InfoStyle().Render("q") + " - Quit\n\n" +
				theme.MutedStyle().Render("Press any key to close"),
		)
}

// ToolCallExample demonstrates tool invocation and output.
func ToolCallExample() string {
	return theme.ModalStyle().
		Width(70).
		Height(12).
		Padding(theme.Unit, theme.Unit*2).
		Render(
			theme.HeaderStyle().Render("Tool Execution") + "\n\n" +
				theme.ToolNameStyle().Render("bash") + " Running: go test ./...\n\n" +
				theme.ToolOutputStyle().Render("=== RUN   TestExample\n    --- PASS: TestExample (0.02s)\n    PASS\n    ok      github.com/example/pkg  0.123s\n\n") +
				theme.SuccessStyle().Render("✓ All tests passed") + "\n\n" +
				theme.MutedStyle().Render("Duration: 0.123s | Exit code: 0"),
		)
}

// StatusMessageExample demonstrates different status message types.
func StatusMessageExample() string {
	return theme.ModalStyle().
		Width(60).
		Height(12).
		Padding(theme.Unit, theme.Unit*2).
		Render(
			theme.HeaderStyle().Render("Status Messages") + "\n\n" +
				theme.InfoStyle().Render("ℹ Info:") + "  Loading configuration...\n" +
				theme.SuccessStyle().Render("✓ Success:") + " File saved successfully\n" +
				theme.WarningStyle().Render("⚠ Warning:") + " File may be corrupted\n" +
				theme.ErrorStyle().Render("✗ Error:") + " Failed to write file\n\n" +
				theme.MutedStyle().Render("Status messages provide feedback to users"),
		)
}

// RunExamples runs all example components and prints them to stdout.
func RunExamples() {
	fmt.Println("\n" + theme.HeaderStyle().Render("=== OpenCode Design System Examples ===") + "\n")

	fmt.Println("\n" + theme.HeaderStyle().Render("1. Modal Dialog"))
	fmt.Println(ModalExample())

	fmt.Println("\n" + theme.HeaderStyle().Render("2. Status Bar"))
	fmt.Println(StatusBarExample())

	fmt.Println("\n" + theme.HeaderStyle().Render("3. Input Field"))
	fmt.Println(InputExample())

	fmt.Println("\n" + theme.HeaderStyle().Render("4. List View"))
	fmt.Println(ListExample())

	fmt.Println("\n" + theme.HeaderStyle().Render("5. Preview Pane"))
	fmt.Println(PreviewExample())

	fmt.Println("\n" + theme.HeaderStyle().Render("6. Diff Display"))
	fmt.Println(DiffExample())

	fmt.Println("\n" + theme.HeaderStyle().Render("7. Chat Interface"))
	fmt.Println(ChatExample())

	fmt.Println("\n" + theme.HeaderStyle().Render("8. Validation Errors"))
	fmt.Println(ValidationExample())

	fmt.Println("\n" + theme.HeaderStyle().Render("9. Text Highlighting"))
	fmt.Println(HighlightExample())

	fmt.Println("\n" + theme.HeaderStyle().Render("10. Keyboard Shortcuts"))
	fmt.Println(KeyboardExample())

	fmt.Println("\n" + theme.HeaderStyle().Render("11. Tool Execution"))
	fmt.Println(ToolCallExample())

	fmt.Println("\n" + theme.HeaderStyle().Render("12. Status Messages"))
	fmt.Println(StatusMessageExample())

	fmt.Println("\n" + theme.MutedStyle().Render("\n=== Examples Complete ==="))
}
