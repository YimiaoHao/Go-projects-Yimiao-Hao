package main

import (
	"image/color" // color utilities
	"log"         // logging
	"math/rand"   // random number generation
	"runtime"     // runtime control (e.g. GOMAXPROCS)
	"sync"        // WaitGroup for goroutine sync
	"time"        // time utilities

	"github.com/hajimehoshi/ebiten/v2" // game / rendering framework
)

const (
	w         = 300 // grid width (logical)
	h         = 300 // grid height (logical)
	scale     = 2   // each cell is drawn as 2x2 pixels
	tickEvery = 5   // update the life board every 5 frames
)

// ArticleStyleGame holds the whole game state
type ArticleStyleGame struct {
	grid     [w][h]uint8 // current generation: 0 = dead, 1 = alive
	buffer   [w][h]uint8 // next generation buffer, to avoid overwriting
	frameCnt int         // frame counter, used to throttle updates
}

// NewArticleStyleGame creates and seeds a new game
func NewArticleStyleGame() *ArticleStyleGame {
	g := &ArticleStyleGame{}
	// seed the grid with random cells (skip the outer border)
	for x := 1; x < w-1; x++ {
		for y := 1; y < h-1; y++ {
			if rand.Float32() < 0.5 { // 50% chance to start alive
				g.grid[x][y] = 1
			}
		}
	}
	return g
}

// parallelUpdate computes the next generation in parallel
func (g *ArticleStyleGame) parallelUpdate() {
	var wg sync.WaitGroup
	// let Go use all available CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// process each column in a separate goroutine
	for x := 1; x < w-1; x++ {
		wg.Add(1)
		// capture x
		go func(x int) {
			defer wg.Done()

			// iterate over all rows in this column (skip top/bottom border)
			for y := 1; y < h-1; y++ {
				// count 8 neighbours around (x,y)
				n := g.grid[x-1][y-1] + // top-left
					g.grid[x-1][y] + // left
					g.grid[x-1][y+1] + // bottom-left
					g.grid[x][y-1] + // top
					g.grid[x][y+1] + // bottom
					g.grid[x+1][y-1] + // top-right
					g.grid[x+1][y] + // right
					g.grid[x+1][y+1] // bottom-right

				// apply Conway's rules into the buffer
				if g.grid[x][y] == 0 && n == 3 {
					// dead cell with exactly 3 neighbours → birth
					g.buffer[x][y] = 1
				} else if n < 2 || n > 3 {
					// underpopulation or overpopulation → death
					g.buffer[x][y] = 0
				} else {
					// stays the same
					g.buffer[x][y] = g.grid[x][y]
				}
			}
		}(x)
	}
	// wait for all columns to finish
	wg.Wait()

	// swap current grid and buffer
	// (new generation becomes current, old current becomes next buffer)
	g.grid, g.buffer = g.buffer, g.grid
}

// Update is called every frame by Ebiten
func (g *ArticleStyleGame) Update() error {
	g.frameCnt++
	// update only every N frames so we can see the changes
	if g.frameCnt%tickEvery == 0 {
		g.parallelUpdate()
	}
	return nil
}

// Draw renders the current grid to the screen
func (g *ArticleStyleGame) Draw(screen *ebiten.Image) {
	// fill background
	screen.Fill(color.RGBA{69, 145, 196, 255})
	aliveColor := color.RGBA{255, 230, 120, 255}

	// draw all alive cells
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if g.grid[x][y] == 1 {
				// draw a scale x scale block for each alive cell
				for i := 0; i < scale; i++ {
					for j := 0; j < scale; j++ {
						screen.Set(x*scale+i, y*scale+j, aliveColor)
					}
				}
			}
		}
	}
}

// Layout defines the logical screen size
func (g *ArticleStyleGame) Layout(ow, oh int) (int, int) {
	return w * 1, h * 1
}

func main() {
	// seed RNG once
	rand.Seed(time.Now().UnixNano())

	game := NewArticleStyleGame()

	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Game of Life (parallel from article)")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
