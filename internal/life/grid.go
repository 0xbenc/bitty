// Package life implements Conway's Game of Life independently of the terminal UI.
package life

import "math/rand/v2"

// Grid is a rectangular collection of cells stored row-major.
type Grid struct {
	width, height int
	cells         []bool
}

// New returns an empty grid. Width and height are clamped to at least one.
func New(width, height int) *Grid {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	return &Grid{width: width, height: height, cells: make([]bool, width*height)}
}

func (g *Grid) Width() int  { return g.width }
func (g *Grid) Height() int { return g.height }

func (g *Grid) Alive(x, y int) bool {
	return x >= 0 && x < g.width && y >= 0 && y < g.height && g.cells[y*g.width+x]
}

func (g *Grid) Set(x, y int, alive bool) {
	if x >= 0 && x < g.width && y >= 0 && y < g.height {
		g.cells[y*g.width+x] = alive
	}
}

func (g *Grid) Toggle(x, y int) { g.Set(x, y, !g.Alive(x, y)) }

func (g *Grid) Clear() { clear(g.cells) }

func (g *Grid) Randomize(density float64, seed uint64) {
	if density < 0 {
		density = 0
	}
	if density > 1 {
		density = 1
	}
	rng := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))
	for i := range g.cells {
		g.cells[i] = rng.Float64() < density
	}
}

func (g *Grid) Population() int {
	n := 0
	for _, alive := range g.cells {
		if alive {
			n++
		}
	}
	return n
}

// Step advances the grid by one generation. With wrap enabled, opposite edges
// touch; otherwise cells outside the grid are always dead.
func (g *Grid) Step(wrap bool) {
	next := make([]bool, len(g.cells))
	for y := range g.height {
		for x := range g.width {
			neighbors := g.neighbors(x, y, wrap)
			next[y*g.width+x] = neighbors == 3 || (g.Alive(x, y) && neighbors == 2)
		}
	}
	g.cells = next
}

func (g *Grid) neighbors(x, y int, wrap bool) int {
	n := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if wrap {
				nx = (nx + g.width) % g.width
				ny = (ny + g.height) % g.height
			}
			if g.Alive(nx, ny) {
				n++
			}
		}
	}
	return n
}
