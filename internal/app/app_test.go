package app

import (
	"strings"
	"testing"

	"github.com/0xbenc/bitty/internal/termstyle"
)

func testModel() *Model {
	theme := termstyle.TerminalTheme()
	theme.NoColor = true
	return New(Options{Width: 8, Height: 8, Density: 0, Speed: 12, Seed: 42, Wrap: true, Paused: true, Theme: theme}, 8, 6)
}

func TestEditingAndStepping(t *testing.T) {
	m := testModel()
	m.handle(keyToggle)
	m.handle(keyRight)
	m.handle(keyToggle)
	m.handle(keyDown)
	m.handle(keyToggle)
	if got := m.grid.Population(); got != 3 {
		t.Fatalf("population = %d, want 3", got)
	}
	m.handle(keyStep)
	if m.generation != 1 {
		t.Fatalf("generation = %d, want 1", m.generation)
	}
	m.handle(keyClear)
	if !m.paused || m.grid.Population() != 0 || m.generation != 0 {
		t.Fatalf("clear left model in state %+v", m)
	}
}

func TestFitPansOversizedWorld(t *testing.T) {
	m := testModel()
	m.cursorX, m.cursorY = 7, 7
	m.fit(4, 4)
	if m.scrollX != 4 || m.scrollPairs != 2 {
		t.Fatalf("scroll = (%d,%d), want (4,2)", m.scrollX, m.scrollPairs)
	}
}

func collectKeys(t *testing.T, input string) []key {
	t.Helper()
	ch := make(chan key)
	go readKeys(strings.NewReader(input), ch)
	var got []key
	for k := range ch {
		got = append(got, k)
	}
	return got
}

func TestReadKeysQuitGrammar(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []key
	}{
		{"ctrl+c quits", "\x03", []key{keyQuit}},
		{"ctrl+q quits", "\x11", []key{keyQuit}},
		{"bare esc quits", "\x1b", []key{keyQuit}},
		{"up arrow unaffected", "\x1b[A", []key{keyUp}},
		{"q is inert", "q", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectKeys(t, tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("readKeys(%q) = %v, want %v", tt.input, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("readKeys(%q)[%d] = %v, want %v", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

// Chosen semantics for esc followed by a buffered byte: with in-memory input
// the whole burst is already buffered when the 0x1b is read, so "\x1bp" is
// treated like an alt+key chord and degrades to the plain key — the 'p'
// dispatches as keyPause and no keyQuit is emitted. A standalone esc (nothing
// buffered, as when the terminal delivers it alone) quits instead; see the
// "bare esc quits" case above.
func TestReadKeysEscThenBufferedKeyIsNotSwallowed(t *testing.T) {
	got := collectKeys(t, "\x1bp")
	if len(got) != 1 || got[0] != keyPause {
		t.Fatalf("readKeys(\"\\x1bp\") = %v, want [keyPause] (p must not be swallowed)", got)
	}
}

func TestRenderUsesHalfBlocksAndFooter(t *testing.T) {
	m := testModel()
	m.grid.Set(0, 0, true)
	m.grid.Set(1, 1, true)
	got := m.Render(80, 6)
	if !strings.Contains(got, "▀▄") {
		t.Fatalf("render did not pack cells into half blocks: %q", got)
	}
	if !strings.Contains(got, "arrows move / space cell") {
		t.Fatalf("render did not include canonical footer: %q", got)
	}
}
