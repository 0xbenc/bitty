package life

import "testing"

func TestBlinker(t *testing.T) {
	g := New(5, 5)
	g.Set(1, 2, true)
	g.Set(2, 2, true)
	g.Set(3, 2, true)
	g.Step(false)
	for y := range 5 {
		for x := range 5 {
			want := x == 2 && y >= 1 && y <= 3
			if got := g.Alive(x, y); got != want {
				t.Fatalf("cell (%d,%d) = %v, want %v", x, y, got, want)
			}
		}
	}
}

func TestWrap(t *testing.T) {
	g := New(3, 3)
	g.Set(0, 0, true)
	g.Set(0, 2, true)
	g.Set(2, 0, true)
	g.Step(true)
	if !g.Alive(2, 2) {
		t.Fatal("corner should be born from wrapped neighbors")
	}
}

func TestRandomizeIsRepeatable(t *testing.T) {
	a, b := New(20, 20), New(20, 20)
	a.Randomize(.3, 42)
	b.Randomize(.3, 42)
	for y := range a.Height() {
		for x := range a.Width() {
			if a.Alive(x, y) != b.Alive(x, y) {
				t.Fatal("same seed produced different grids")
			}
		}
	}
}
