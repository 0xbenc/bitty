// Package app owns bitty's interactive editor and simulation loop.
package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/0xbenc/bitty/internal/life"
	"github.com/0xbenc/bitty/internal/terminal"
	"github.com/0xbenc/bitty/internal/termstyle"
	"github.com/0xbenc/termchrome"
	"github.com/0xbenc/termnav"
)

type Options struct {
	Width, Height int
	Density       float64
	Speed         int
	Seed          uint64
	Wrap          bool
	Paused        bool
	Theme         termstyle.Theme
	Input         *os.File
	Output        *os.File
}

type Model struct {
	grid                 *life.Grid
	cursorX, cursorY     int
	scrollX, scrollPairs int
	generation           uint64
	density              float64
	seed                 uint64
	speed                int
	wrap, paused         bool
	theme                termstyle.Theme
}

func New(opts Options, termWidth, termHeight int) *Model {
	width, height := opts.Width, opts.Height
	if width <= 0 {
		width = termWidth
	}
	if height <= 0 {
		height = max(2, (termHeight-3)*2)
	}
	g := life.New(width, height)
	g.Randomize(opts.Density, opts.Seed)
	speed := opts.Speed
	if speed < 1 {
		speed = 1
	}
	return &Model{grid: g, density: opts.Density, seed: opts.Seed, speed: speed, wrap: opts.Wrap, paused: opts.Paused, theme: opts.Theme}
}

func Run(opts Options) error {
	in, out := opts.Input, opts.Output
	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}
	if !terminal.IsTTY(in) || !terminal.IsTTY(out) {
		return fmt.Errorf("interactive terminal required")
	}
	cols, rows, err := terminal.Size(out)
	if err != nil {
		return err
	}
	m := New(opts, cols, rows)
	old, err := terminal.MakeRaw(in)
	if err != nil {
		return err
	}
	fmt.Fprint(out, "\x1b[?1049h\x1b[?25l")
	defer func() {
		fmt.Fprint(out, "\x1b[?25h\x1b[?1049l")
		_ = terminal.Restore(in, old)
	}()

	keys := make(chan key, 8)
	go readKeys(in, keys)
	stepTimer := time.NewTimer(frameDuration(m.speed))
	defer stepTimer.Stop()
	resize := time.NewTicker(250 * time.Millisecond)
	defer resize.Stop()

	for {
		m.fit(cols, rows)
		if _, err := io.WriteString(out, m.Render(cols, rows)); err != nil {
			return err
		}
		select {
		case k, ok := <-keys:
			if !ok {
				return nil
			}
			if m.handle(k) {
				return nil
			}
		case <-stepTimer.C:
			if !m.paused {
				m.grid.Step(m.wrap)
				m.generation++
			}
			stepTimer.Reset(frameDuration(m.speed))
		case <-resize.C:
			if newCols, newRows, sizeErr := terminal.Size(out); sizeErr == nil {
				cols, rows = newCols, newRows
			}
		}
	}
}

func frameDuration(speed int) time.Duration { return time.Second / time.Duration(max(1, speed)) }

type key int

const (
	keyUnknown key = iota
	keyUp
	keyDown
	keyLeft
	keyRight
	keyToggle
	keyPause
	keyStep
	keyRandom
	keyClear
	keyWrap
	keyFaster
	keySlower
	keyQuit
)

func readKeys(r io.Reader, dst chan<- key) {
	reader := bufio.NewReader(r)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			close(dst)
			return
		}
		if b == 0x1b {
			second, err := reader.ReadByte()
			if err != nil || second != '[' {
				continue
			}
			third, err := reader.ReadByte()
			if err != nil {
				continue
			}
			switch third {
			case 'A':
				dst <- keyUp
			case 'B':
				dst <- keyDown
			case 'C':
				dst <- keyRight
			case 'D':
				dst <- keyLeft
			}
			continue
		}
		switch b {
		case 3, 'q':
			dst <- keyQuit
		case ' ', 'x':
			dst <- keyToggle
		case '\r', '\n', 'p':
			dst <- keyPause
		case 'n':
			dst <- keyStep
		case 'r':
			dst <- keyRandom
		case 'c':
			dst <- keyClear
		case 'w':
			dst <- keyWrap
		case '+', '=':
			dst <- keyFaster
		case '-', '_':
			dst <- keySlower
		case 'k':
			dst <- keyUp
		case 'j':
			dst <- keyDown
		case 'h':
			dst <- keyLeft
		case 'l':
			dst <- keyRight
		}
	}
}

func (m *Model) handle(k key) bool {
	switch k {
	case keyQuit:
		return true
	case keyUp:
		m.cursorY = max(0, m.cursorY-1)
	case keyDown:
		m.cursorY = min(m.grid.Height()-1, m.cursorY+1)
	case keyLeft:
		m.cursorX = max(0, m.cursorX-1)
	case keyRight:
		m.cursorX = min(m.grid.Width()-1, m.cursorX+1)
	case keyToggle:
		m.grid.Toggle(m.cursorX, m.cursorY)
	case keyPause:
		m.paused = !m.paused
	case keyStep:
		m.grid.Step(m.wrap)
		m.generation++
	case keyRandom:
		m.seed++
		m.grid.Randomize(m.density, m.seed)
		m.generation = 0
	case keyClear:
		m.grid.Clear()
		m.generation = 0
		m.paused = true
	case keyWrap:
		m.wrap = !m.wrap
	case keyFaster:
		m.speed = min(60, m.speed+1)
	case keySlower:
		m.speed = max(1, m.speed-1)
	}
	return false
}

func (m *Model) fit(cols, rows int) {
	boardRows := max(1, rows-2)
	m.cursorX, m.scrollX = termnav.ClampWindow(m.grid.Width(), m.cursorX, m.scrollX, func(start, cursor int) bool {
		return cursor >= start && cursor < start+max(1, cols)
	})
	pairs := (m.grid.Height() + 1) / 2
	cursorPair := m.cursorY / 2
	cursorPair, m.scrollPairs = termnav.ClampWindow(pairs, cursorPair, m.scrollPairs, func(start, cursor int) bool {
		return cursor >= start && cursor < start+boardRows
	})
	m.cursorY = min(m.grid.Height()-1, cursorPair*2+m.cursorY%2)
}

func (m *Model) Render(cols, rows int) string {
	if cols < 1 || rows < 1 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\x1b[H")
	status := "RUN"
	statusRole := termstyle.RoleSuccess
	if m.paused {
		status, statusRole = "PAUSE", termstyle.RoleWarning
	}
	header := fmt.Sprintf("bitty  %s  gen %d  pop %d  %d/s  wrap %s  cell %d,%d",
		m.theme.Style(statusRole, status), m.generation, m.grid.Population(), m.speed, onOff(m.wrap), m.cursorX, m.cursorY)
	b.WriteString(termstyle.Truncate(termstyle.Sanitize(header), cols))
	b.WriteString("\x1b[K\r\n")
	boardRows := max(0, rows-2)
	for screenY := 0; screenY < boardRows; screenY++ {
		pair := m.scrollPairs + screenY
		y := pair * 2
		for screenX := 0; screenX < cols; screenX++ {
			x := m.scrollX + screenX
			if x >= m.grid.Width() || y >= m.grid.Height() {
				b.WriteByte(' ')
				continue
			}
			top, bottom := m.grid.Alive(x, y), m.grid.Alive(x, y+1)
			cell := " "
			switch {
			case top && bottom:
				cell = "█"
			case top:
				cell = "▀"
			case bottom:
				cell = "▄"
			}
			if x == m.cursorX && pair == m.cursorY/2 {
				cell = m.theme.Style(termstyle.RoleSelected, cell)
			} else if top || bottom {
				cell = m.theme.Style(termstyle.RolePrimary, cell)
			}
			b.WriteString(cell)
		}
		b.WriteString("\x1b[K")
		if screenY < boardRows-1 || rows > 1 {
			b.WriteString("\r\n")
		}
	}
	if rows > 1 {
		hints := []termchrome.KeyHint{
			{Key: "arrows", Label: "move"}, {Key: "space", Label: "cell"},
			{Key: "p", Label: "pause"}, {Key: "n", Label: "step"},
			{Key: "r", Label: "random"}, {Key: "c", Label: "clear"},
			{Key: "+/-", Label: "speed"}, {Key: "w", Label: "wrap"}, {Key: "q", Label: "quit"},
		}
		footer := termchrome.Footer(hints, cols)
		b.WriteString(termstyle.Truncate(m.theme.Style(termstyle.RoleMuted, footer), cols))
		b.WriteString("\x1b[K")
	}
	return b.String()
}

func onOff(v bool) string {
	if v {
		return "on"
	}
	return "off"
}
