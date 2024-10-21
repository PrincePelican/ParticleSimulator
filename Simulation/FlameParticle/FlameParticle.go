package flameparticle

import (
	particles "FallingSand/Particles"
	"math/rand"
)

func HandleFlameDown(x int, y int, grid [][]particles.Particle) *particles.ChangeParticle {
	grid[x][y].Type = particles.FLAME
	if tryToFlame(x, y+1, grid) {
		return &particles.ChangeParticle{CurrentPosition: *particles.NewPoint(x, y+1), Type: particles.FLAME, ChangeType: particles.CHANGE}
	}
	return nil
}

func HandleFlameUp(x int, y int, grid [][]particles.Particle) *particles.ChangeParticle {
	grid[x][y].Type = particles.FLAME
	if tryToFlame(x, y-1, grid) {
		return &particles.ChangeParticle{CurrentPosition: *particles.NewPoint(x, y-1), Type: particles.FLAME, ChangeType: particles.CHANGE}
	}
	return nil
}

func HandleFlameLeft(x int, y int, grid [][]particles.Particle) *particles.ChangeParticle {
	grid[x][y].Type = particles.FLAME
	if tryToFlame(x-1, y, grid) {
		return &particles.ChangeParticle{CurrentPosition: *particles.NewPoint(x-1, y), Type: particles.FLAME, ChangeType: particles.CHANGE}
	}
	return nil
}

func HandleFlameRight(x int, y int, grid [][]particles.Particle) *particles.ChangeParticle {
	grid[x][y].Type = particles.FLAME
	if tryToFlame(x+1, y, grid) {
		return &particles.ChangeParticle{CurrentPosition: *particles.NewPoint(x+1, y), Type: particles.FLAME, ChangeType: particles.CHANGE}
	}
	return nil
}

func TickFlameTime(flameParticle *particles.Particle, grid [][]particles.Particle) *particles.ChangeParticle {
	if flameParticle.FlameTime == 0 {
		flameParticle.Type = particles.EMPTY
		grid[flameParticle.Position.X][flameParticle.Position.Y] = *particles.NewSmokeParticle(flameParticle.Position.X, flameParticle.Position.Y, particles.SMOKE)
		return particles.NewChangeParticle(flameParticle.Position, flameParticle.Position, particles.SMOKE, particles.CHANGE)
	}
	flameParticle.FlameTime--
	return nil
}

func tryToFlame(x int, y int, grid [][]particles.Particle) bool {
	flameRoll := rand.Intn(9)
	if flameRoll <= grid[x][y].FlameRate {
		grid[x][y].Type = particles.FLAME
		grid[x][y].FlameRate = 0
		return true
	}
	return false
}
