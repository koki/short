package pager

import (
	"bufio"
	"index/suffixarray"
	"io"
)

type ViewBuffer struct {
	buf         []string
	begin_index int
	max_height  int
	last_line   int
	scanner     *bufio.Scanner
	index       *suffixarray.Index
	run_length  []int
	prev_token  string
	index_index int
}

func NewViewBuffer(r io.Reader, index *suffixarray.Index) *ViewBuffer {
	return &ViewBuffer{
		buf:         []string{},
		begin_index: 0,
		last_line:   -1,
		max_height:  64,
		scanner:     bufio.NewScanner(r),
		index:       index,
		run_length:  []int{},
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
		return false
	}
	if !v.scanner.Scan() {
		//reached the end of input
		v.last_line = len(v.buf)
		return false
	}
	line := v.scanner.Text()
	v.buf = append(v.buf, line)

	cummulative_length := 0
	if len(v.run_length) != 0 {
		cummulative_length = v.run_length[len(v.run_length)-1]
	}
	v.run_length = append(v.run_length, cummulative_length+len(line))

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
		i := 0
		for {
			if i == n {
				break
			}
			if !v.scan() {
				break
			}
			i = i + 1
		}
		v.begin_index = v.begin_index + i
	}
}

func (v *ViewBuffer) ScrollBottom() {
	i := 0
	for {
		if !v.scan() {
			break
		}
		i = i + 1
	}
	v.begin_index = len(v.buf) - v.max_height
	if v.begin_index <= 0 {
		v.begin_index = 0
	}
}

func (v *ViewBuffer) Search(token string) bool {
	if token == "/" {
		return false
	}
	indices := v.index.Lookup([]byte(token[1:]), -1)
	if len(indices) == 0 {
		return false
	}

	if token == v.prev_token {
		v.index_index = v.index_index + 1
	} else {
		v.prev_token = token
		v.index_index = 0
	}

	if v.index_index >= len(indices) {
		v.index_index = 0
	}

	for i := range v.run_length {
		if v.run_length[i] >= indices[v.index_index] {
			if i > 0 {
				v.begin_index = i - 1
			} else {
				v.begin_index = 0
			}
			return true
		}
	}
	return false
}
