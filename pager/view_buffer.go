package pager

import (
	"bufio"
	"io"
	"strings"
)

type ViewBuffer struct {
	// Lines
	buf []string

	beginIndex int
	maxHeight  int

	// Line index of the last line (-1 if we haven't gotten there yet.)
	lastLine int

	scanner *bufio.Scanner

	// Search
	lastSearchToken      string
	lastSearchResultLine int
}

func NewViewBuffer(r io.Reader) *ViewBuffer {
	return &ViewBuffer{
		buf:        []string{},
		beginIndex: 0,
		lastLine:   -1,
		maxHeight:  64,
		scanner:    bufio.NewScanner(r),
	}
}

func (v *ViewBuffer) SetMaxHeight(height int) {
	v.maxHeight = height
}

func (v *ViewBuffer) CurrentView() []string {
	for len(v.buf)-v.beginIndex < v.maxHeight {
		if !v.scan() {
			break
		}
	}
	if len(v.buf)-v.beginIndex < v.maxHeight {
		v.beginIndex = len(v.buf) - v.maxHeight
	}

	if v.beginIndex < 0 {
		v.beginIndex = 0
	}

	start := v.beginIndex
	end := v.beginIndex + v.maxHeight
	if end > len(v.buf) {
		end = len(v.buf)
	}

	return v.buf[start:end]
}

func (v *ViewBuffer) scan() bool {
	if v.lastLine != -1 {
		// Reached the end of input earlier.
		return false
	}
	if !v.scanner.Scan() {
		// Just reached the end of input.
		v.lastLine = len(v.buf)
		return false
	}

	// Grab the next line.
	line := v.scanner.Text()
	v.buf = append(v.buf, line)

	return true
}

func (v *ViewBuffer) ScrollUp() {
	if v.beginIndex > 0 {
		v.beginIndex = v.beginIndex - 1
	}
}

func (v *ViewBuffer) ScrollTop() {
	v.beginIndex = 0
}

func (v *ViewBuffer) ScrollDown() {
	if len(v.buf)-v.beginIndex > v.maxHeight {
		v.beginIndex = v.beginIndex + 1
	} else if v.lastLine == -1 {
		if v.scan() {
			v.beginIndex = v.beginIndex + 1
		}
	}
}

func (v *ViewBuffer) ScrollDownN(n int) {
	if len(v.buf)-v.beginIndex > v.maxHeight {
		if len(v.buf)-v.beginIndex >= v.maxHeight+n {
			v.beginIndex = v.beginIndex + n
		} else {
			remaining := len(v.buf) - v.beginIndex - v.maxHeight
			v.beginIndex = v.beginIndex + remaining
		}
	} else if v.lastLine == -1 {
		var i int
		for i = 0; i < n; i++ {
			if !v.scan() {
				break
			}
		}
		v.beginIndex = v.beginIndex + i
	}
}

func (v *ViewBuffer) ScrollToLine(projectedBeginIndex int) {
	projectedEndIndex := projectedBeginIndex + v.maxHeight
	if projectedBeginIndex <= 0 {
		// Scroll to the top.
		v.beginIndex = 0
		return
	}

	// Make sure we've read enough lines to scroll to the right spot.
	for len(v.buf) <= projectedEndIndex {
		if !v.scan() {
			break
		}
	}

	// Don't actually scroll past the end of the buffer.
	if len(v.buf) >= projectedEndIndex {
		v.beginIndex = projectedBeginIndex
		return
	}

	// Scroll to the bottom.
	v.beginIndex = len(v.buf) - v.maxHeight
}

func (v *ViewBuffer) ScrollBottom() {
	for v.scan() {
	}
	v.beginIndex = len(v.buf) - v.maxHeight
	if v.beginIndex <= 0 {
		v.beginIndex = 0
	}
}

func (v *ViewBuffer) SearchFromLine(lineIndex int, searchToken string, stopAtLineIndex int) (int, bool) {
	for {
		if lineIndex == stopAtLineIndex {
			// Don't keep looping through the lines forever.
			return -1, false
		}

		// Make sure we know when to give up the search.
		if stopAtLineIndex == -1 {
			// Don't search lineIndex again.
			stopAtLineIndex = lineIndex
		}

		// Make sure we've read as far as the line we're searching.
		for lineIndex >= len(v.buf) && v.scan() {
		}

		if lineIndex >= len(v.buf) {
			// We passed the end of the file. Jump back to the start of the file.
			lineIndex = 0
			continue
		}

		line := v.buf[lineIndex]
		if strings.Contains(line, searchToken) {
			// Found it!
			return lineIndex, true
		}

		// Keep looking.
		lineIndex++
	}
}

func (v *ViewBuffer) Search(token string) bool {
	token = token[1:]
	if len(token) == 0 {
		return false
	}

	if v.lastSearchToken != token {
		v.lastSearchToken = token
		v.lastSearchResultLine = -1
	}

	if lineIndex, ok := v.SearchFromLine(v.lastSearchResultLine+1, token, v.lastSearchResultLine); ok {
		v.ScrollToLine(lineIndex - 1)
		v.lastSearchResultLine = lineIndex
		return true
	}

	return false
}
