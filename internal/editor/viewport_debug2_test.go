package editor

import (
	"fmt"
	"testing"
)

func TestViewportEnsureVisible_Debug(t *testing.T) {
	testCases := []struct {
		name    string
		initial Viewport
		line    int
		wantTop int
	}{
		{
			name:    "line already visible in middle third",
			initial: NewViewport(10).ScrollTo(5),
			line:    7,
			wantTop: 5,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.initial
			result := v.EnsureVisible(tt.line)

			fmt.Printf("\n=== Debug ===\n")
			fmt.Printf("Initial: topLine=%d, height=%d\n", v.topLine, v.height)

			start, end := v.VisibleLines()
			third := v.height / 3
			fmt.Printf("Visible range: [%d, %d)\n", start, end)
			fmt.Printf("Third: %d\n", third)
			fmt.Printf("Top third: [%d, %d)\n", start, start+third)
			fmt.Printf("Middle third: [%d, %d)\n", start+third, end-third)
			fmt.Printf("Bottom third: [%d, %d)\n", end-third, end)

			fmt.Printf("Target line: %d\n", tt.line)
			fmt.Printf("Line %d in range? %v\n", tt.line, tt.line >= start && tt.line < end)
			fmt.Printf("Line %d in top third? %v\n", tt.line, tt.line >= start && tt.line < start+third)
			fmt.Printf("Line %d in middle third? %v\n", tt.line, tt.line >= start+third && tt.line < end-third)
			fmt.Printf("Line %d in bottom third? %v\n", tt.line, tt.line >= end-third && tt.line < end)

			fmt.Printf("Expected topLine: %d, Got: %d\n", tt.wantTop, result.topLine)
			fmt.Printf("Match? %v\n", tt.wantTop == result.topLine)

			if result.topLine != tt.wantTop {
				t.Errorf("expected topLine %d, got %d", tt.wantTop, result.topLine)
			}
		})
	}
}
