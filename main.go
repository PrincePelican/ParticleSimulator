package main

import (
	"FallingSand/sandsimulation"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	gridWidth    = 320
	gridHeight   = 180
	cellSize     = 4
	clickSize    = 10
)

type Game struct {
	grid       [gridWidth][gridHeight]bool
	imageCache *ebiten.Image
	gridInfo   chan [gridWidth][gridHeight]bool
}

func (g *Game) Update() error {

	g.grid = <-g.gridInfo

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		gridX, gridY := mx/cellSize, my/cellSize
		if gridX-clickSize >= 0 && gridX+clickSize < gridWidth && gridY-clickSize >= 0 && gridY+clickSize < gridHeight {
			for x := gridX; x < gridX+clickSize; x++ {
				for y := gridY; y < gridY+clickSize; y++ {
					g.grid[x][y] = true
				}
			}
		}
	}
	go sandsimulation.ContinueFall(g.grid, g.gridInfo)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.imageCache == nil {
		g.imageCache = ebiten.NewImage(screenWidth, screenHeight)
	}

	g.imageCache.Clear()

	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			if g.grid[x][y] {
				vector.DrawFilledRect(g.imageCache, float32(x*cellSize), float32(y*cellSize), cellSize, cellSize, color.RGBA{255, 215, 0, 1}, true)
			}
		}
	}
	screen.DrawImage(g.imageCache, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Falling Sand")
	gridChannel := make(chan [gridWidth][gridHeight]bool)
	game := &Game{gridInfo: gridChannel}
	go sandsimulation.ContinueFall(game.grid, game.gridInfo)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
