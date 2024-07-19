package main

import (
	particles "FallingSand/Particles"
	simulation "FallingSand/Simulation"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	gridWidth    = 1280
	gridHeight   = 720
	cellSize     = 1
	clickSize    = 50
)

type Game struct {
	imageCache     *ebiten.Image
	simulation     simulation.Simulation
	movedParticles []particles.MovedParticle
}

func (g *Game) Update() error {
	g.movedParticles = g.simulation.NextStep()

	g.HandleInput()

	return nil
}

func (g *Game) HandleInput() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		var createdParticles []particles.Particle
		mx, my := ebiten.CursorPosition()
		gridX, gridY := mx/cellSize, my/cellSize
		if gridX >= 0 && gridX+clickSize < gridWidth && gridY >= 0 && gridY+clickSize < gridHeight {
			for x := gridX + 1; x < gridX+clickSize; x += rand.Intn(5) + 1 {
				for y := gridY; y < gridY+clickSize; y += rand.Intn(5) + 1 {
					createdParticles = append(createdParticles, *particles.NewParticle(x, y, particles.SAND))
				}
			}
		}
		g.simulation.AddParticles(createdParticles)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		var createdParticles []particles.Particle
		mx, my := ebiten.CursorPosition()
		gridX, gridY := mx/cellSize, my/cellSize
		if gridX >= 0 && gridX+clickSize < gridWidth && gridY >= 0 && gridY+clickSize < gridHeight {
			for x := gridX + 1; x < gridX+clickSize; x += rand.Intn(5) + 1 {
				for y := gridY; y < gridY+clickSize; y += rand.Intn(5) + 1 {
					createdParticles = append(createdParticles, *particles.NewParticle(x, y, particles.WATER))
				}
			}
		}
		g.simulation.AddParticles(createdParticles)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.imageCache == nil {
		g.imageCache = ebiten.NewImage(screenWidth, screenHeight)
	}

	g.DrawGrid()

	screen.DrawImage(g.imageCache, nil)
}

func (g *Game) DrawGrid() {
	for x := 0; x < len(g.movedParticles); x++ {
		particleColor := color.RGBA{28, 107, 160, 1}
		if g.movedParticles[x].Type == particles.SAND {
			particleColor = color.RGBA{255, 215, 0, 255}
		}
		g.imageCache.Set(g.movedParticles[x].CurrentPosition.X*cellSize, g.movedParticles[x].CurrentPosition.Y*cellSize, particleColor)
		g.imageCache.Set(g.movedParticles[x].PreviousPosition.X*cellSize, g.movedParticles[x].PreviousPosition.Y*cellSize, color.RGBA{0, 0, 0, 255})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Falling Sand")
	simulation := simulation.NewSimulation()
	simulation.Initialize()
	game := &Game{simulation: *simulation}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
