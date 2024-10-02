package flameparticle

import (
	particles "FallingSand/Particles"
	"math/rand"
)

func HandleFlameDown(x int, y int, grid [][]particles.Particle) *particles.MovedParticle {
	grid[x][y].Type = particles.FLAME
	if tryToFlame(x, y+1, grid) {
		return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, y), CurrentPosition: *particles.NewPoint(x, y+1)}
	}

	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, y), CurrentPosition: *particles.NewPoint(x, y)}
}

func HandleFlameUp(x int, y int, grid [][]particles.Particle) *particles.MovedParticle {
	grid[x][y].Type = particles.FLAME
	if tryToFlame(x, y-1, grid) {
		return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, y), CurrentPosition: *particles.NewPoint(x, y-1)}
	}
	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, y), CurrentPosition: *particles.NewPoint(x, y)}
}

func HandleFlameLeft(x int, y int, grid [][]particles.Particle) *particles.MovedParticle {
	grid[x][y].Type = particles.FLAME
	if tryToFlame(x-1, y, grid) {
		return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, y), CurrentPosition: *particles.NewPoint(x-1, y)}
	}
	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, y), CurrentPosition: *particles.NewPoint(x, y)}
}

func HandleFlameRight(x int, y int, grid [][]particles.Particle) *particles.MovedParticle {
	grid[x][y].Type = particles.FLAME
	if tryToFlame(x+1, y, grid) {
		return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, y), CurrentPosition: *particles.NewPoint(x+1, y)}
	}
	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, y), CurrentPosition: *particles.NewPoint(x, y)}
}

func tryToFlame(x int, y int, grid [][]particles.Particle) bool {
	flameRoll := rand.Intn(9)
	if flameRoll >= grid[x][y].FlameRate {
		grid[x][y].Type = particles.FLAME
		grid[x][y].FlameRate = 0
		return true
	}
	return false
}
