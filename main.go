package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game implements ebiten.Game interface.
type Game struct {
	grid [][]bool
	size int
}

func NewGrid(W, H, size int, probability float32) [][]bool {
	grid := [][]bool{}//make([][]bool, W/size)
	for y := size; y < H; y += size {
		row := []bool{}
		for x := size; x < W; x += size {
			row = append(row, rand.Float32() > probability)
		}
		grid = append(grid, row)
	}
	return grid
}

func countNeighbours(grid [][]bool, x, y int) int {
	N := 0
	if y > 0 && grid[y-1][x] {N += 1}
	if y+1 < len(grid) && grid[y+1][x] {N += 1}
	if x > 0 && grid[y][x-1] {N += 1}
	if x+1 < len(grid[y]) && grid[y][x+1] {N += 1}
	return N
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
    // Write your game's logical update.	
	for y, row := range g.grid {
		for x, _ := range row {
			N := countNeighbours(g.grid, x, y)
			alive := g.grid[y][x]
			// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
			// Any live cell with two or three live neighbours lives on to the next generation.
			// Any live cell with more than three live neighbours dies, as if by overpopulation.
			// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
			if N < 2 || N > 3 {
				alive = false
			}
			if N == 2 || N == 3 {
				alive = true
			}
			g.grid[y][x] = alive
		}
	}
    return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
    // Write your game's rendering.
	block := ebiten.NewImage(g.size, g.size)
	block.Fill(color.White)
	position := ebiten.GeoM{}
	for i, row := range g.grid {
		for _, cell := range row {
			if cell {
				screen.DrawImage(block, &ebiten.DrawImageOptions{
					GeoM: position,
				})
			}
			position.Translate(float64(g.size), 0)
			} 
		position = ebiten.GeoM{}
		position.Translate(0, float64(i * g.size))
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 400, 300
}

func main() {
	W := 800
	H := 600
	SIZE := 5
    game := &Game{
		grid: NewGrid(W, H, SIZE, .75),
		size: SIZE,
	}
    // Specify the window size as you like. Here, a doubled size is specified.
    ebiten.SetWindowSize(W, H)
    ebiten.SetWindowTitle("Conway's Game of Life")
	ebiten.SetMaxTPS(5)

	rand.Seed(time.Now().UnixNano())

    // Call ebiten.RunGame to start your game loop.
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}