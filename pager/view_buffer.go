package pager

import (
	"bufio"
	"io"
	"strings"
)

type ViewBuffer struct {
	// Lines
	buf []string

	begin_index int
	max_height  int

	// Line index of the last line (-1 if we haven't gotten there yet.)
	last_line int

	scanner *bufio.Scanner

	// Search
	lastSearchToken      string
	lastSearchResultLine int
}

func NewViewBuffer(r io.Reader) *ViewBuffer {
	return &ViewBuffer{
		buf:         []string{},
		begin_index: 0,
		last_line:   -1,
		max_height:  64,
		scanner:     bufio.NewScanner(r),
	}
}

func (v *ViewBuffer) SetMaxHeight(height int) {
	v.max_height = height
}

func (v *ViewBuffer) CurrentView() []string {
	for len(v.buf)-v.begin_index < v.max_height {
		if !v.scan() {
			break
		}
	}
	if len(v.buf)-v.begin_index < v.max_height {
		v.begin_index = len(v.buf) - v.max_height
	}

	if v.begin_index < 0 {
		v.begin_index = 0
	}

	start := v.begin_index
	end := v.begin_index + v.max_height
	if end > len(v.buf) {
		end = len(v.buf)
	}

	return v.buf[start:end]
}

func (v *ViewBuffer) scan() bool {
	if v.last_line != -1 {
		// Reached the end of input earlier.
		return false
	}
	if !v.scanner.Scan() {
		// Just reached the end of input.
		v.last_line = len(v.buf)
		return false
	}

	// Grab the next line.
	line := v.scanner.Text()
	v.buf = append(v.buf, line)

	return true
}

func (v *ViewBuffer) ScrollUp() {
	if v.begin_index > 0 {
		v.begin_index = v.begin_index - 1
	}
}

func (v *ViewBuffer) ScrollTop() {
	v.begin_index = 0
}

func (v *ViewBuffer) ScrollDown() {
	if len(v.buf)-v.begin_index > v.max_height {
		v.begin_index = v.begin_index + 1
	} else if v.last_line == -1 {
		if v.scan() {
			v.begin_index = v.begin_index + 1
		}
	}
}

func (v *ViewBuffer) ScrollDownN(n int) {
	if len(v.buf)-v.begin_index > v.max_height {
		if len(v.buf)-v.begin_index >= v.max_height+n {
			v.begin_index = v.begin_index + n
		} else {
			remaining := len(v.buf) - v.begin_index - v.max_height
			v.begin_index = v.begin_index + remaining
		}
	} else if v.last_line == -1 {
		var i int
		for i = 0; i < n; i++ {
			if !v.scan() {
				break
			}
		}
		v.begin_index = v.begin_index + i
	}
}

func (v *ViewBuffer) ScrollToLine(projectedBeginIndex int) {
	projectedEndIndex := projectedBeginIndex + v.max_height
	if projectedBeginIndex <= 0 {
		// Scroll to the top.
		v.begin_index = 0
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
		v.begin_index = projectedBeginIndex
		return
	}

	// Scroll to the bottom.
	v.begin_index = len(v.buf) - v.max_height
}

func (v *ViewBuffer) ScrollBottom() {
	for v.scan() {
	}
	v.begin_index = len(v.buf) - v.max_height
	if v.begin_index <= 0 {
		v.begin_index = 0
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
