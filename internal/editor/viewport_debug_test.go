package editor

import (
	"fmt"
	"testing"
)

func TestMiddleThirdLogic(t *testing.T) {
	testCases := []struct {
		height  int
		topLine int
		line    int
	}{
		{10, 5, 7},   // line 7 in viewport with topLine=5
		{10, 5, 8},   // line 8 in viewport with topLine=5
		{10, 5, 5},   // line at top edge
		{10, 0, 9},   // line at bottom edge
		{30, 10, 15}, // line in top third
		{30, 10, 25}, // line in bottom third
		{30, 10, 20}, // line in middle third
	}

	for _, tc := range testCases {
		v := NewViewport(tc.height).ScrollTo(tc.topLine)
		result := v.EnsureVisible(tc.line)

		fmt.Printf("\n=== Test Case ===\n")
		fmt.Printf("Height: %d, TopLine: %d, TargetLine: %d\n", tc.height, tc.topLine, tc.line)
		fmt.Printf("Before: topLine=%d, height=%d\n", v.topLine, v.height)

		start, end := v.VisibleLines()
		third := v.height / 3
		fmt.Printf("Visible range: [%d, %d)\n", start, end)
		fmt.Printf("Third: %d\n", third)
		fmt.Printf("Top third: [%d, %d)\n", start, start+third)
		fmt.Printf("Middle third: [%d, %d)\n", start+third, end-third)
		fmt.Printf("Bottom third: [%d, %d)\n", end-third, end)
		fmt.Printf("Line %d in range? %v\n", tc.line, tc.line >= start && tc.line < end)
		fmt.Printf("Line %d in top third? %v\n", tc.line, tc.line >= start && tc.line < start+third)
		fmt.Printf("Line %d in middle third? %v\n", tc.line, tc.line >= start+third && tc.line < end-third)
		fmt.Printf("Line %d in bottom third? %v\n", tc.line, tc.line >= end-third && tc.line < end)

		fmt.Printf("After: topLine=%d\n", result.topLine)
	}
}
