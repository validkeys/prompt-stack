package editor

import (
	"fmt"
	"testing"
)

func TestBoundaryConditions_Debug(t *testing.T) {
	testCases := []struct {
		name    string
		height  int
		topLine int
		line    int
	}{
		{"cursor at viewport top line", 30, 10, 10},
		{"cursor at viewport bottom-1 line", 30, 10, 39},
		{"cursor enters top third", 30, 0, 5},
		{"cursor enters bottom third", 30, 0, 28},
		{"cursor moves rapidly upward", 30, 100, 80},
	}

	for _, tc := range testCases {
		v := NewViewport(tc.height).ScrollTo(tc.topLine)
		result := v.EnsureVisible(tc.line)

		fmt.Printf("\n=== %s ===\n", tc.name)
		fmt.Printf("Height: %d, TopLine: %d, TargetLine: %d\n", tc.height, tc.topLine, tc.line)

		start, end := v.VisibleLines()
		third := v.height / 3
		fmt.Printf("Visible range: [%d, %d)\n", start, end)
		fmt.Printf("Third: %d\n", third)
		fmt.Printf("Top third: [%d, %d)\n", start, start+third)
		fmt.Printf("Middle third: [%d, %d)\n", start+third, end-third)
		fmt.Printf("Bottom third: [%d, %d)\n", end-third, end)

		fmt.Printf("Result topLine: %d\n", result.topLine)
	}
}
