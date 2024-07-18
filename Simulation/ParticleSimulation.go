package simulation

import particles "FallingSand/Particles"

const (
	gridWidth  = 1280
	gridHeight = 720
)

type Simulation struct {
	grid [][]particles.Particle
}

func NewSimulation() *Simulation {
	return &Simulation{
		grid: make([][]particles.Particle, gridWidth),
	}
}

func (s *Simulation) Initialize() {
	for i := range s.grid {
		s.grid[i] = make([]particles.Particle, gridHeight)
	}
}

func (s *Simulation) AddParticles(newParticlesList []particles.Particle) {
	for x := 0; x < len(newParticlesList); x++ {
		particle := newParticlesList[x]
		s.grid[particle.Position.X][particle.Position.Y] = particle
	}
}

func (s *Simulation) NextStep() []particles.MovedParticle {
	var movedParticles []particles.MovedParticle
	for y := gridHeight - 2; y >= 0; y-- {
		for x := 1; x < gridWidth-1; x++ {
			movedParticle := s.handleCell(x, y)
			if movedParticle != nil {
				movedParticles = append(movedParticles, *movedParticle)
			}
		}
	}
	return movedParticles
}

func (s *Simulation) handleCell(x int, y int) *particles.MovedParticle {
	if s.grid[x][y].Type != particles.EMPTY {
		if s.grid[x][y+1].Type == particles.EMPTY {
			s.grid[x][y], s.grid[x][y+1] = s.grid[x][y+1], s.grid[x][y]
			return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, y), CurrentPosition: *particles.NewPoint(x, y+1)}
		} else if s.grid[x-1][y+1].Type == particles.EMPTY {
			s.grid[x][y], s.grid[x-1][y+1] = s.grid[x-1][y+1], s.grid[x][y]
			return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, y), CurrentPosition: *particles.NewPoint(x-1, y+1)}
		} else if s.grid[x+1][y+1].Type == particles.EMPTY {
			s.grid[x][y], s.grid[x+1][y+1] = s.grid[x+1][y+1], s.grid[x][y]
			return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, y), CurrentPosition: *particles.NewPoint(x+1, y+1)}
		}

	}
	return nil
}
