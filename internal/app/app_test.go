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
