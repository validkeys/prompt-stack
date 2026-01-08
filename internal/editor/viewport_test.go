package editor

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewViewport(t *testing.T) {
	v := NewViewport(24)

	if v.topLine != 0 {
		t.Errorf("expected topLine to be 0, got %d", v.topLine)
	}
	if v.height != 24 {
		t.Errorf("expected height to be 24, got %d", v.height)
	}
	if v.totalLines != 0 {
		t.Errorf("expected totalLines to be 0, got %d", v.totalLines)
	}
}

func TestViewportUpdate(t *testing.T) {
	tests := []struct {
		name       string
		initial    Viewport
		msg        tea.Msg
		wantTop    int
		wantHeight int
	}{
		{
			name:       "window size message updates height",
			initial:    NewViewport(10),
			msg:        tea.WindowSizeMsg{Height: 20},
			wantTop:    0,
			wantHeight: 20,
		},
		{
			name:       "scroll message updates top line",
			initial:    NewViewport(10),
			msg:        ScrollMsg{Line: 5},
			wantTop:    5,
			wantHeight: 10,
		},
		{
			name:       "unhandled message returns unchanged",
			initial:    NewViewport(10),
			msg:        tea.KeyMsg{Type: tea.KeyEnter},
			wantTop:    0,
			wantHeight: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.Update(tt.msg)
			if result.topLine != tt.wantTop {
				t.Errorf("expected topLine %d, got %d", tt.wantTop, result.topLine)
			}
			if result.height != tt.wantHeight {
				t.Errorf("expected height %d, got %d", tt.wantHeight, result.height)
			}
		})
	}
}

func TestViewportScrollTo(t *testing.T) {
	tests := []struct {
		name    string
		initial Viewport
		line    int
		wantTop int
	}{
		{
			name:    "scroll to line 5",
			initial: NewViewport(10),
			line:    5,
			wantTop: 5,
		},
		{
			name:    "scroll to line 0",
			initial: NewViewport(10),
			line:    0,
			wantTop: 0,
		},
		{
			name:    "scroll to line 100",
			initial: NewViewport(10),
			line:    100,
			wantTop: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.ScrollTo(tt.line)
			if result.topLine != tt.wantTop {
				t.Errorf("expected topLine %d, got %d", tt.wantTop, result.topLine)
			}
		})
	}
}

func TestViewportScrollUp(t *testing.T) {
	tests := []struct {
		name    string
		initial Viewport
		lines   int
		wantTop int
	}{
		{
			name:    "scroll up 5 lines",
			initial: NewViewport(10).ScrollTo(10),
			lines:   5,
			wantTop: 5,
		},
		{
			name:    "scroll up past top",
			initial: NewViewport(10).ScrollTo(3),
			lines:   5,
			wantTop: 0,
		},
		{
			name:    "scroll up 0 lines",
			initial: NewViewport(10).ScrollTo(10),
			lines:   0,
			wantTop: 10,
		},
		{
			name:    "scroll up from top",
			initial: NewViewport(10),
			lines:   5,
			wantTop: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.ScrollUp(tt.lines)
			if result.topLine != tt.wantTop {
				t.Errorf("expected topLine %d, got %d", tt.wantTop, result.topLine)
			}
		})
	}
}

func TestViewportScrollDown(t *testing.T) {
	tests := []struct {
		name       string
		initial    Viewport
		lines      int
		wantTop    int
		totalLines int
	}{
		{
			name:       "scroll down 5 lines",
			initial:    NewViewport(10).SetTotalLines(100),
			lines:      5,
			wantTop:    5,
			totalLines: 100,
		},
		{
			name:       "scroll down past end",
			initial:    NewViewport(10).ScrollTo(85).SetTotalLines(100),
			lines:      10,
			wantTop:    90,
			totalLines: 100,
		},
		{
			name:       "scroll down 0 lines",
			initial:    NewViewport(10).SetTotalLines(100),
			lines:      0,
			wantTop:    0,
			totalLines: 100,
		},
		{
			name:       "scroll down with short document",
			initial:    NewViewport(10).SetTotalLines(5),
			lines:      5,
			wantTop:    0,
			totalLines: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.ScrollDown(tt.lines)
			if result.topLine != tt.wantTop {
				t.Errorf("expected topLine %d, got %d", tt.wantTop, result.topLine)
			}
		})
	}
}

func TestViewportVisibleLines(t *testing.T) {
	tests := []struct {
		name      string
		viewport  Viewport
		wantStart int
		wantEnd   int
	}{
		{
			name:      "viewport at top",
			viewport:  NewViewport(10),
			wantStart: 0,
			wantEnd:   10,
		},
		{
			name:      "viewport scrolled down",
			viewport:  NewViewport(10).ScrollTo(5),
			wantStart: 5,
			wantEnd:   15,
		},
		{
			name:      "viewport with height 0",
			viewport:  NewViewport(0),
			wantStart: 0,
			wantEnd:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := tt.viewport.VisibleLines()
			if start != tt.wantStart {
				t.Errorf("expected start %d, got %d", tt.wantStart, start)
			}
			if end != tt.wantEnd {
				t.Errorf("expected end %d, got %d", tt.wantEnd, end)
			}
		})
	}
}

func TestViewportEnsureVisible(t *testing.T) {
	tests := []struct {
		name    string
		initial Viewport
		line    int
		wantTop int
	}{
		{
			name:    "line already visible",
			initial: NewViewport(10).ScrollTo(5),
			line:    7,
			wantTop: 5,
		},
		{
			name:    "line above viewport",
			initial: NewViewport(10).ScrollTo(10),
			line:    5,
			wantTop: 5,
		},
		{
			name:    "line below viewport",
			initial: NewViewport(10).ScrollTo(0),
			line:    15,
			wantTop: 6,
		},
		{
			name:    "line at top edge",
			initial: NewViewport(10).ScrollTo(5),
			line:    5,
			wantTop: 5,
		},
		{
			name:    "line at bottom edge",
			initial: NewViewport(10).ScrollTo(0),
			line:    9,
			wantTop: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.EnsureVisible(tt.line)
			if result.topLine != tt.wantTop {
				t.Errorf("expected topLine %d, got %d", tt.wantTop, result.topLine)
			}
		})
	}
}

func TestViewportSetHeight(t *testing.T) {
	tests := []struct {
		name       string
		initial    Viewport
		height     int
		wantHeight int
	}{
		{
			name:       "set height to 20",
			initial:    NewViewport(10),
			height:     20,
			wantHeight: 20,
		},
		{
			name:       "set height to 0",
			initial:    NewViewport(10),
			height:     0,
			wantHeight: 0,
		},
		{
			name:       "set height to 100",
			initial:    NewViewport(10),
			height:     100,
			wantHeight: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.SetHeight(tt.height)
			if result.height != tt.wantHeight {
				t.Errorf("expected height %d, got %d", tt.wantHeight, result.height)
			}
		})
	}
}

func TestViewportSetTotalLines(t *testing.T) {
	tests := []struct {
		name       string
		initial    Viewport
		totalLines int
		wantTop    int
		wantTotal  int
	}{
		{
			name:       "set total lines to 100",
			initial:    NewViewport(10),
			totalLines: 100,
			wantTop:    0,
			wantTotal:  100,
		},
		{
			name:       "set total lines less than height",
			initial:    NewViewport(10).ScrollTo(5),
			totalLines: 5,
			wantTop:    0,
			wantTotal:  5,
		},
		{
			name:       "set total lines with viewport at end",
			initial:    NewViewport(10).ScrollTo(95),
			totalLines: 100,
			wantTop:    90,
			wantTotal:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.SetTotalLines(tt.totalLines)
			if result.topLine != tt.wantTop {
				t.Errorf("expected topLine %d, got %d", tt.wantTop, result.topLine)
			}
			if result.totalLines != tt.wantTotal {
				t.Errorf("expected totalLines %d, got %d", tt.wantTotal, result.totalLines)
			}
		})
	}
}

func TestViewportGetters(t *testing.T) {
	v := NewViewport(24).ScrollTo(10).SetTotalLines(100)

	if v.TopLine() != 10 {
		t.Errorf("expected TopLine() to return 10, got %d", v.TopLine())
	}
	if v.Height() != 24 {
		t.Errorf("expected Height() to return 24, got %d", v.Height())
	}
	if v.TotalLines() != 100 {
		t.Errorf("expected TotalLines() to return 100, got %d", v.TotalLines())
	}
}

func TestViewportImmutability(t *testing.T) {
	original := NewViewport(10).ScrollTo(5)

	// All operations should return new instances
	scrolled := original.ScrollTo(10)
	if original.topLine == scrolled.topLine {
		t.Error("ScrollTo should return a new instance")
	}

	up := original.ScrollUp(2)
	if original.topLine == up.topLine {
		t.Error("ScrollUp should return a new instance")
	}

	down := original.ScrollDown(2)
	if original.topLine == down.topLine {
		t.Error("ScrollDown should return a new instance")
	}

	height := original.SetHeight(20)
	if original.height == height.height {
		t.Error("SetHeight should return a new instance")
	}

	total := original.SetTotalLines(100)
	if original.totalLines == total.totalLines {
		t.Error("SetTotalLines should return a new instance")
	}

	// Original should remain unchanged
	if original.topLine != 5 {
		t.Errorf("original topLine should remain 5, got %d", original.topLine)
	}
	if original.height != 10 {
		t.Errorf("original height should remain 10, got %d", original.height)
	}
}

func TestViewportEdgeCases(t *testing.T) {
	t.Run("empty document", func(t *testing.T) {
		v := NewViewport(10).SetTotalLines(0)
		if v.topLine != 0 {
			t.Errorf("expected topLine to be 0 for empty document, got %d", v.topLine)
		}
	})

	t.Run("document shorter than viewport", func(t *testing.T) {
		v := NewViewport(10).SetTotalLines(5)
		if v.topLine != 0 {
			t.Errorf("expected topLine to be 0 when document shorter than viewport, got %d", v.topLine)
		}
	})

	t.Run("document exactly viewport height", func(t *testing.T) {
		v := NewViewport(10).SetTotalLines(10)
		if v.topLine != 0 {
			t.Errorf("expected topLine to be 0 when document equals viewport height, got %d", v.topLine)
		}
	})

	t.Run("scroll to negative line", func(t *testing.T) {
		v := NewViewport(10).ScrollTo(-5)
		if v.topLine != -5 {
			t.Errorf("expected topLine to be -5, got %d", v.topLine)
		}
	})

	t.Run("ensure visible with negative line", func(t *testing.T) {
		v := NewViewport(10).ScrollTo(5)
		result := v.EnsureVisible(-1)
		if result.topLine != -1 {
			t.Errorf("expected topLine to be -1, got %d", result.topLine)
		}
	})
}
