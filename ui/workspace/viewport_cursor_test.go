package workspace

import (
	"fmt"
	"testing"
)

// TestAdjustViewport_MiddleThirdScrolling tests the middle-third scrolling strategy
func TestAdjustViewport_MiddleThirdScrolling(t *testing.T) {
	tests := []struct {
		name           string
		viewportHeight int
		totalLines     int
		cursorY        int
		initialOffset  int
		expectedOffset int
		expectedScroll bool
	}{
		{
			name:           "cursor in middle third, no scroll",
			viewportHeight: 12,
			totalLines:     100,
			cursorY:        6,
			initialOffset:  0,
			expectedOffset: 0,
			expectedScroll: false,
		},
		{
			name:           "cursor in top third, scroll up",
			viewportHeight: 12,
			totalLines:     100,
			cursorY:        2,
			initialOffset:  5,
			expectedOffset: 0,
			expectedScroll: true,
		},
		{
			name:           "cursor in bottom third, scroll down",
			viewportHeight: 12,
			totalLines:     100,
			cursorY:        10,
			initialOffset:  0,
			expectedOffset: 2,
			expectedScroll: true,
		},
		{
			name:           "cursor at document start",
			viewportHeight: 12,
			totalLines:     100,
			cursorY:        0,
			initialOffset:  5,
			expectedOffset: 0,
			expectedScroll: true,
		},
		{
			name:           "cursor at document end",
			viewportHeight: 12,
			totalLines:     100,
			cursorY:        99,
			initialOffset:  0,
			expectedOffset: 88,
			expectedScroll: true,
		},
		{
			name:           "viewport taller than document",
			viewportHeight: 50,
			totalLines:     10,
			cursorY:        5,
			initialOffset:  0,
			expectedOffset: 0,
			expectedScroll: false,
		},
		{
			name:           "one line viewport",
			viewportHeight: 1,
			totalLines:     100,
			cursorY:        50,
			initialOffset:  0,
			expectedOffset: 49,
			expectedScroll: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			model := New()
			model = model.SetSize(80, tt.viewportHeight+1)
			model = model.SetContent(generateLines(tt.totalLines))
			model.viewport.SetYOffset(tt.initialOffset)

			// Sync viewport content so it knows the total line count
			model = model.syncViewportLines()

			// Set cursor position directly
			if tt.cursorY < tt.totalLines {
				err := model.buffer.SetCursorPosition(0, tt.cursorY)
				if err != nil {
					t.Fatalf("Failed to set cursor position: %v", err)
				}
			}

			// Verify cursor position was set
			_, actualCursorY := model.buffer.CursorPosition()
			if actualCursorY != tt.cursorY {
				t.Fatalf("Expected cursor at line %d, got line %d", tt.cursorY, actualCursorY)
			}

			// Execute
			model = model.adjustViewport()

			// Assert
			actualOffset := model.viewport.YOffset
			if actualOffset != tt.expectedOffset {
				t.Errorf("Expected offset %d, got %d", tt.expectedOffset, actualOffset)
			}

			// Verify scroll occurred (or didn't)
			scrolled := (actualOffset != tt.initialOffset)
			if scrolled != tt.expectedScroll {
				t.Errorf("Expected scroll=%v, got scroll=%v", tt.expectedScroll, scrolled)
			}

			// Verify YOffset is never negative
			if model.viewport.YOffset < 0 {
				t.Errorf("YOffset should never be negative, got %d", model.viewport.YOffset)
			}

			// Verify YOffset doesn't exceed document bounds
			maxOffset := tt.totalLines - tt.viewportHeight
			if maxOffset < 0 {
				maxOffset = 0
			}
			if model.viewport.YOffset > maxOffset {
				t.Errorf("YOffset %d exceeds max offset %d", model.viewport.YOffset, maxOffset)
			}
		})
	}
}

// TestAdjustViewport_BoundsChecking tests that bounds checking prevents invalid offsets
func TestAdjustViewport_BoundsChecking(t *testing.T) {
	tests := []struct {
		name           string
		viewportHeight int
		totalLines     int
		cursorY        int
		initialOffset  int
	}{
		{
			name:           "zero height viewport",
			viewportHeight: 0,
			totalLines:     100,
			cursorY:        50,
			initialOffset:  0,
		},
		{
			name:           "cursor beyond document",
			viewportHeight: 12,
			totalLines:     10,
			cursorY:        100,
			initialOffset:  0,
		},
		{
			name:           "negative initial offset (should be corrected)",
			viewportHeight: 12,
			totalLines:     100,
			cursorY:        50,
			initialOffset:  -10,
		},
		{
			name:           "offset beyond document end (should be corrected)",
			viewportHeight: 12,
			totalLines:     50,
			cursorY:        25,
			initialOffset:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			model := New()
			model = model.SetSize(80, tt.viewportHeight+1)
			model = model.SetContent(generateLines(tt.totalLines))
			model.viewport.SetYOffset(tt.initialOffset)

			// Sync viewport content before adjusting
			model = model.syncViewportLines()

			// Set cursor position
			if tt.cursorY < tt.totalLines {
				err := model.buffer.SetCursorPosition(0, tt.cursorY)
				if err != nil {
					t.Fatalf("Failed to set cursor position: %v", err)
				}
			}

			// Execute - should not panic
			model = model.adjustViewport()

			// Assert - YOffset should be valid (no panic means success)
			// Additional bounds checking
			if model.viewport.Height > 0 && model.viewport.YOffset < 0 {
				t.Errorf("YOffset should not be negative after adjustment, got %d", model.viewport.YOffset)
			}

			maxOffset := tt.totalLines - tt.viewportHeight
			if maxOffset < 0 {
				maxOffset = 0
			}
			if model.viewport.Height > 0 && model.viewport.YOffset > maxOffset {
				t.Errorf("YOffset %d exceeds max offset %d", model.viewport.YOffset, maxOffset)
			}
		})
	}
}

// TestAdjustViewport_RapidCursorMovement tests rapid cursor movement maintains smooth scrolling
func TestAdjustViewport_RapidCursorMovement(t *testing.T) {
	// Setup
	model := New()
	model = model.SetSize(80, 25)
	totalLines := 1000
	model = model.SetContent(generateLines(totalLines))

	// Move cursor rapidly down 100 times
	for i := 0; i < 100; i++ {
		// Set cursor position directly
		err := model.buffer.SetCursorPosition(0, i+1)
		if err != nil {
			t.Fatalf("Failed to set cursor position at line %d: %v", i+1, err)
		}
		// Sync viewport content before adjusting
		model = model.syncViewportLines()

		// Verify cursor position was set
		_, cursorY := model.buffer.CursorPosition()
		if cursorY != i+1 {
			t.Fatalf("Expected cursor at line %d, got line %d", i+1, cursorY)
		}

		model = model.adjustViewport()

		// Verify cursor stays visible
		_, cursorY = model.buffer.CursorPosition()
		viewportBottom := model.viewport.YOffset + model.viewport.Height

		if cursorY < model.viewport.YOffset || cursorY >= viewportBottom {
			t.Errorf("Cursor not visible at iteration %d: cursor=%d, offset=%d, bottom=%d",
				i, cursorY, model.viewport.YOffset, viewportBottom)
		}

		// Verify no negative offsets
		if model.viewport.YOffset < 0 {
			t.Errorf("Negative YOffset at iteration %d: %d", i, model.viewport.YOffset)
			break
		}
	}

	// Verify final state
	_, finalCursorY := model.buffer.CursorPosition()
	if finalCursorY != 100 {
		t.Errorf("Expected cursor at line 100, got %d", finalCursorY)
	}
}

// generateLines generates test content with specified number of lines
func generateLines(count int) string {
	if count <= 0 {
		return ""
	}

	result := ""
	for i := 0; i < count; i++ {
		if i > 0 {
			result += "\n"
		}
		result += fmt.Sprintf("Line %d", i+1)
	}
	return result
}
