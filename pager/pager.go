package pager

import (
	"index/suffixarray"
	"io"
	"strings"

	"github.com/nsf/termbox-go"
)

type Pager struct {
	buf *ViewBuffer
}

func NewPager(r io.Reader, index *suffixarray.Index) *Pager {
	return &Pager{
		buf: NewViewBuffer(r, index),
	}
}

// Renders the input buffer until the user exits
func (p *Pager) Render() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()
	max_width, max_height := termbox.Size()

	termbox.SetInputMode(termbox.InputEsc)

	//very important to set height
	p.buf.SetMaxHeight(max_height - 1)

	g := false
	brk := false
	search := false
	search_result := false
	token := ""
	back := false
	p.draw(max_width, max_height, "")
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventResize {
			// TODO: handle resize
		}
		if ev.Type == termbox.EventKey {
			switch ev.Key {
			case termbox.KeyEnter:
				if search || search_result {
					p.buf.Search(token)
					search_result = true
				} else {
					p.buf.ScrollDown()
				}
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				p.buf.ScrollDown()
			case termbox.KeySpace:
				if !search && !search_result {
					p.buf.ScrollDownN(30)
				} else {
					token = token + " "
				}
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				p.buf.ScrollUp()
			case termbox.KeyEsc:
				if !search_result && !search {
					brk = true
				}
				if search {
					search = false
					token = ""
				}
				if search_result {
					search_result = false
					token = ""
				}
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if len(token) > 0 {
					token = token[:len(token)-1]
				}
				if len(token) == 0 {
					search = false
					search_result = false
				}
				back = true
			case termbox.KeyCtrlC, termbox.KeyCtrlD:
				brk = true
			default:
			}

			if !search && !search_result {
				switch ev.Ch {
				case 'G':
					p.buf.ScrollBottom()
				case 'g':
					if g {
						p.buf.ScrollTop()
						g = false
					}
					g = true
				case 'q':
					brk = true
				case '/':
					search = true
					token = "/"
				default:
				}
			} else {
				if !back {
					if int32(ev.Ch) != 0 {
						token = token + string(ev.Ch)
					}
				}
			}
		}

		if brk {
			break
		}

		p.draw(max_width, max_height, token)
		back = false
	}
	termbox.Sync()
	return nil
}

// Draws the current view onto the screen
func (p *Pager) draw(max_width, max_height int, token string) error {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	lines := p.buf.CurrentView()

	width := 0
	height := 0

	for _, line := range lines {
		copyOfLine := line
		toSearch := ""
		if len(token) > 0 {
			toSearch = token[1:]
		}
		juju := strings.Map(func(r rune) rune {
			return r + 1
		}, toSearch)
		copyOfLine = strings.Replace(copyOfLine, toSearch, juju, -1)
		for i, ch := range line {
			fg := termbox.ColorDefault
			bg := termbox.ColorDefault
			if copyOfLine[i] != line[i] {
				fg = termbox.ColorWhite
				bg = termbox.ColorBlack
			}
			termbox.SetCell(width, height, ch, fg, bg)
			if width < max_width {
				width = width + 1
			} else if height < max_height {
				height = height + 1
			}
		}
		width = 0
		if height < max_height {
			height = height + 1
		} else {
			break
		}
	}

	tail_token := token
	if token == "" {
		tail_token = `Press '/' for search`
	}

	for i, ch := range tail_token {
		termbox.SetCell(i, max_height-1, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()
	return nil
}
