package main

import (
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
	gridWidth    = 640
	gridHeight   = 360
	cellSize     = 2
	clickSize    = 20
)

type Game struct {
	grid       [gridWidth][gridHeight]bool
	imageCache *ebiten.Image
}

func randomSign() int {
	return 2*rand.Intn(2) - 1
}

func (g *Game) Update() error {
	for y := gridHeight - 2; y >= 0; y-- {
		for x := 1; x < gridWidth-1; x++ {
			if g.grid[x][y] && !g.grid[x][y+1] {
				g.grid[x][y] = false
				g.grid[x][y+1] = true
			} else if g.grid[x][y] && g.grid[x][y+1] {
				if !g.grid[x+1][y+1] && !g.grid[x-1][y+1] {
					g.grid[x][y] = false
					g.grid[x+randomSign()][y+1] = true
				} else if !g.grid[x-1][y+1] {
					g.grid[x][y] = false
					g.grid[x-1][y+1] = true
				} else if !g.grid[x+1][y+1] {
					g.grid[x][y] = false
					g.grid[x+1][y+1] = true
				}
			}
		}

	}

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
	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
