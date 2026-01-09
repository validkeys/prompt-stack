package theme_test

import (
	"strings"
	"testing"

	"github.com/kyledavis/prompt-stack/ui/theme"
)

func BenchmarkModalStyleRender(b *testing.B) {
	style := theme.ModalStyle()
	text := "Test content for modal style"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkStatusStyleRender(b *testing.B) {
	style := theme.StatusStyle()
	text := "Status: Ready"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkInfoStyleRender(b *testing.B) {
	style := theme.InfoStyle()
	text := "Info message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkSuccessStyleRender(b *testing.B) {
	style := theme.SuccessStyle()
	text := "Success!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkWarningStyleRender(b *testing.B) {
	style := theme.WarningStyle()
	text := "Warning!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkErrorStyleRender(b *testing.B) {
	style := theme.ErrorStyle()
	text := "Error!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkInputStyleRender(b *testing.B) {
	style := theme.InputStyle()
	text := "Input text"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkSearchInputStyleRender(b *testing.B) {
	style := theme.SearchInputStyle()
	text := "Search..."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkListItemStyleRender(b *testing.B) {
	style := theme.ListItemStyle()
	text := "List item"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkListItemSelectedStyleRender(b *testing.B) {
	style := theme.ListItemSelectedStyle()
	text := "Selected item"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkListCategoryStyleRender(b *testing.B) {
	style := theme.ListCategoryStyle()
	text := "Category"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkDiffAddedStyleRender(b *testing.B) {
	style := theme.DiffAddedStyle()
	text := "+ added line"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkDiffRemovedStyleRender(b *testing.B) {
	style := theme.DiffRemovedStyle()
	text := "- removed line"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkDiffContextStyleRender(b *testing.B) {
	style := theme.DiffContextStyle()
	text := "  context line"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkChainedStylesRender(b *testing.B) {
	style := theme.ModalStyle().Width(80).Height(40).Padding(2, 3)
	text := "Test content"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(text)
	}
}

func BenchmarkComplexStyledContent(b *testing.B) {
	tests := []struct {
		name  string
		lines int
	}{
		{"10 lines", 10},
		{"50 lines", 50},
		{"100 lines", 100},
		{"500 lines", 500},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			content := generateStyledContent(tt.lines)
			style := theme.ModalContentStyle()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = style.Render(content)
			}
		})
	}
}

func BenchmarkMultipleStyles(b *testing.B) {
	tests := []struct {
		name    string
		content string
	}{
		{"single line", "Single line of text"},
		{"short paragraph", "This is a short paragraph of text that demonstrates basic rendering performance."},
		{"medium paragraph", "This is a medium paragraph of text. It contains multiple sentences. The purpose is to test rendering performance with a reasonable amount of content. It should be representative of typical use cases."},
		{"long paragraph", strings.Repeat("This is a long paragraph. ", 20)},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = theme.ModalStyle().Render(tt.content)
				_ = theme.InfoStyle().Render(tt.content)
				_ = theme.SuccessStyle().Render(tt.content)
				_ = theme.WarningStyle().Render(tt.content)
				_ = theme.ErrorStyle().Render(tt.content)
				_ = theme.InputStyle().Render(tt.content)
				_ = theme.ListItemStyle().Render(tt.content)
			}
		})
	}
}

func BenchmarkStyleWithModifiers(b *testing.B) {
	text := "Test content with modifiers"

	tests := []struct {
		name  string
		style func(string) string
	}{
		{"width only", func(s string) string { return theme.ModalStyle().Width(80).Render(s) }},
		{"height only", func(s string) string { return theme.ModalStyle().Height(40).Render(s) }},
		{"padding only", func(s string) string { return theme.ModalStyle().Padding(2, 3).Render(s) }},
		{"width and height", func(s string) string { return theme.ModalStyle().Width(80).Height(40).Render(s) }},
		{"padding and margin", func(s string) string { return theme.ModalStyle().Padding(2, 3).Margin(1, 2).Render(s) }},
		{"all modifiers", func(s string) string {
			return theme.ModalStyle().Width(80).Height(40).Padding(2, 3).Margin(1, 2).Render(s)
		}},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = tt.style(text)
			}
		})
	}
}

func generateStyledContent(lines int) string {
	builder := strings.Builder{}
	for i := 0; i < lines; i++ {
		builder.WriteString("Line ")
		builder.WriteString(string(rune('0' + i%10)))
		builder.WriteString(": This is styled content for benchmarking purposes.\n")
	}
	return builder.String()
}
