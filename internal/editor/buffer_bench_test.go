package editor_test

import (
	"strings"
	"testing"

	"github.com/kyledavis/prompt-stack/internal/editor"
)

func BenchmarkBufferMoveUp(b *testing.B) {
	buf := editor.NewBufferWithContent(strings.Repeat("line\n", 1000))
	buf.SetCursorPosition(10, 500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.MoveUp()
	}
}

func BenchmarkBufferMoveDown(b *testing.B) {
	buf := editor.NewBufferWithContent(strings.Repeat("line\n", 1000))
	buf.SetCursorPosition(10, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.MoveDown()
	}
}

func BenchmarkBufferMoveLeft(b *testing.B) {
	buf := editor.NewBufferWithContent(strings.Repeat("a", 10000))
	buf.SetCursorPosition(5000, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.MoveLeft()
	}
}

func BenchmarkBufferMoveRight(b *testing.B) {
	buf := editor.NewBufferWithContent(strings.Repeat("a", 10000))
	buf.SetCursorPosition(5000, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.MoveRight()
	}
}

func BenchmarkBufferMoveToLineStart(b *testing.B) {
	buf := editor.NewBufferWithContent(strings.Repeat("a", 10000))
	buf.SetCursorPosition(5000, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.MoveToLineStart()
	}
}

func BenchmarkBufferMoveToLineEnd(b *testing.B) {
	buf := editor.NewBufferWithContent(strings.Repeat("a", 10000))
	buf.SetCursorPosition(5000, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.MoveToLineEnd()
	}
}

func BenchmarkBufferRapidMovement(b *testing.B) {
	buf := editor.NewBufferWithContent(strings.Repeat("line\n", 1000))
	buf.SetCursorPosition(10, 500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.MoveUp()
		buf.MoveDown()
		buf.MoveLeft()
		buf.MoveRight()
	}
}
