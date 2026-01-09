package editor

import tea "github.com/charmbracelet/bubbletea"

// Viewport manages viewport scrolling and visible area
type Viewport struct {
	topLine    int // First visible line
	height     int // Viewport height in lines
	totalLines int // Total lines in document
}

// NewViewport creates a new viewport with specified height
func NewViewport(height int) Viewport {
	return Viewport{
		topLine:    0,
		height:     height,
		totalLines: 0,
	}
}

// Update handles viewport-related messages
func (v Viewport) Update(msg tea.Msg) Viewport {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return v.SetHeight(msg.Height)
	case ScrollMsg:
		return v.ScrollTo(msg.Line)
	}
	return v
}

// ScrollTo moves viewport to show specified line
func (v Viewport) ScrollTo(line int) Viewport {
	newViewport := v
	newViewport.topLine = line
	return newViewport
}

// ScrollUp moves viewport up by specified number of lines
func (v Viewport) ScrollUp(lines int) Viewport {
	newViewport := v
	newViewport.topLine -= lines
	if newViewport.topLine < 0 {
		newViewport.topLine = 0
	}
	return newViewport
}

// ScrollDown moves viewport down by specified number of lines
func (v Viewport) ScrollDown(lines int) Viewport {
	newViewport := v
	newViewport.topLine += lines
	// Don't scroll past the end of the document
	maxTop := v.totalLines - v.height
	if maxTop < 0 {
		maxTop = 0
	}
	if newViewport.topLine > maxTop {
		newViewport.topLine = maxTop
	}
	return newViewport
}

// VisibleLines returns range of visible line numbers
func (v Viewport) VisibleLines() (start, end int) {
	return v.topLine, v.topLine + v.height
}

// EnsureVisible ensures line is visible in viewport using middle-third scrolling strategy
// Keeps cursor in middle third of viewport for better visibility in large documents
func (v Viewport) EnsureVisible(line int) Viewport {
	if v.height <= 0 {
		return v
	}

	start, end := v.VisibleLines()
	third := v.height / 3

	if line < start {
		return v.ScrollTo(line)
	}

	if line >= end {
		return v.ScrollTo(line - v.height + 1)
	}

	if line < start+third {
		return v.ScrollTo(line - third)
	}

	if line >= end-third {
		return v.ScrollTo(line - v.height + third)
	}

	return v
}

// SetHeight updates the viewport height
func (v Viewport) SetHeight(height int) Viewport {
	newViewport := v
	newViewport.height = height
	return newViewport
}

// SetTotalLines updates the total number of lines in the document
func (v Viewport) SetTotalLines(total int) Viewport {
	newViewport := v
	newViewport.totalLines = total
	// Adjust topLine if needed
	maxTop := total - v.height
	if maxTop < 0 {
		maxTop = 0
	}
	if newViewport.topLine > maxTop {
		newViewport.topLine = maxTop
	}
	return newViewport
}

// TopLine returns the first visible line number
func (v Viewport) TopLine() int {
	return v.topLine
}

// Height returns the viewport height in lines
func (v Viewport) Height() int {
	return v.height
}

// TotalLines returns the total number of lines in the document
func (v Viewport) TotalLines() int {
	return v.totalLines
}

// ScrollMsg is a message to scroll the viewport
type ScrollMsg struct {
	Line int
}
