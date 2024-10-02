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
	imageCache          *ebiten.Image
	simulation          simulation.Simulation
	movedParticles      []particles.MovedParticle
	currentParticleType particles.ParticleType
}

func (g *Game) Update() error {
	g.movedParticles = g.simulation.NextStep()

	g.HandleInput()

	return nil
}

func (g *Game) HandleInput() {

	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.currentParticleType = particles.SAND
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		g.currentParticleType = particles.WATER
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		g.currentParticleType = particles.WOOD
	} else if ebiten.IsKeyPressed(ebiten.Key4) {
		g.currentParticleType = particles.EMPTY
	} else if ebiten.IsKeyPressed(ebiten.Key5) {
		g.currentParticleType = particles.FLAME
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		var createdParticles []particles.Particle
		mx, my := ebiten.CursorPosition()
		gridX, gridY := mx/cellSize, my/cellSize
		if gridX >= 0 && gridX+clickSize < gridWidth && gridY >= 0 && gridY+clickSize < gridHeight {
			if g.currentParticleType == particles.SAND || g.currentParticleType == particles.WATER {
				for x := gridX + 1; x < gridX+clickSize; x += rand.Intn(10) + 1 {
					for y := gridY; y < gridY+clickSize; y += rand.Intn(10) + 1 {
						createdParticles = append(createdParticles, *particles.NewParticle(x, y, g.currentParticleType))
					}
				}
			} else if g.currentParticleType == particles.FLAME {
				createdParticles = append(createdParticles, *particles.NewParticle(gridX, gridY, g.currentParticleType))
			} else {
				for x := gridX + 1; x < gridX+clickSize/2; x++ {
					for y := gridY; y < gridY+clickSize/2; y++ {
						createdParticles = append(createdParticles, *particles.NewParticle(x, y, g.currentParticleType))
					}
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
		switch g.movedParticles[x].Type {
		case particles.SAND:
			particleColor = color.RGBA{255, 215, 0, 255} // Yellow for Sand
		case particles.WATER:
			particleColor = color.RGBA{0, 0, 255, 255} // Blue for Water
		case particles.WOOD:
			particleColor = color.RGBA{100, 25, 0, 255} // Brown for Wood
		case particles.FLAME:
			particleColor = color.RGBA{255, 165, 0, 255} // Brown for Wood
		default:
			particleColor = color.RGBA{0, 0, 0, 255} // Default to black
		}
		if g.movedParticles[x].Type != particles.FLAME {
			g.imageCache.Set(g.movedParticles[x].PreviousPosition.X*cellSize, g.movedParticles[x].PreviousPosition.Y*cellSize, color.RGBA{0, 0, 0, 255})
			g.imageCache.Set(g.movedParticles[x].CurrentPosition.X*cellSize, g.movedParticles[x].CurrentPosition.Y*cellSize, particleColor)
		} else {
			g.imageCache.Set(g.movedParticles[x].PreviousPosition.X*cellSize, g.movedParticles[x].PreviousPosition.Y*cellSize, particleColor)
			g.imageCache.Set(g.movedParticles[x].CurrentPosition.X*cellSize, g.movedParticles[x].CurrentPosition.Y*cellSize, particleColor)
		}

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
