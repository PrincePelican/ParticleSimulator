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
	for x := 1; x < gridWidth-1; x++ {
		for y := gridHeight - 2; y >= 0; y-- {
			var movedParticle *particles.MovedParticle
			var Type particles.ParticleType = particles.EMPTY
			switch s.grid[x][y].Type {
			case particles.SAND:
				movedParticle = s.handleSandParticle(x, y)
				Type = particles.SAND
			case particles.WATER:
				movedParticle = s.handleWaterParticle(x, y)
				Type = particles.WATER
			}

			if movedParticle != nil {
				movedParticle.Type = Type
				movedParticles = append(movedParticles, *movedParticle)
			}
		}
	}
	return movedParticles
}

func (s *Simulation) handleSandParticle(x int, y int) *particles.MovedParticle {
	if s.grid[x][y].Type != particles.EMPTY {
		if s.grid[x][y+1].Type == particles.EMPTY {
			return s.handleDownMove(x, y)
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

func (s *Simulation) handleWaterParticle(x int, y int) *particles.MovedParticle {
	if s.grid[x][y].Type != particles.EMPTY {
		if s.grid[x][y+1].Type == particles.EMPTY {
			return s.handleDownMove(x, y)
		} else if s.grid[x-1][y+1].Type == particles.EMPTY {
			s.grid[x][y], s.grid[x-1][y+1] = s.grid[x-1][y+1], s.grid[x][y]
			return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, y), CurrentPosition: *particles.NewPoint(x-1, y+1)}
		} else if s.grid[x+1][y+1].Type == particles.EMPTY {
			s.grid[x][y], s.grid[x+1][y+1] = s.grid[x+1][y+1], s.grid[x][y]
			return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, y), CurrentPosition: *particles.NewPoint(x+1, y+1)}
		} else if s.grid[x+1][y].Type == particles.EMPTY {
			return s.handleMoveRight(x, y)
		} else if s.grid[x-1][y].Type == particles.EMPTY {
			return s.handleMoveLeft(x, y)
		}

	}
	return nil
}

func (s *Simulation) handleDownMove(x int, y int) *particles.MovedParticle {
	startY := y
	particle := s.grid[x][y]
	for ; y < startY+particle.Velocity && y < gridHeight-1; y++ {
		if s.grid[x][y+1].Type != particles.EMPTY {
			break
		}
	}
	s.grid[x][startY], s.grid[x][y] = s.grid[x][y], s.grid[x][startY]
	s.grid[x][y].Velocity += particle.Gravity
	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, startY), CurrentPosition: *particles.NewPoint(x, y)}
}

func (s *Simulation) handleMoveRight(x int, y int) *particles.MovedParticle {
	startX := x
	particle := s.grid[x][y]
	for ; x < startX+particle.Velocity && x < gridWidth-1; x++ {
		if s.grid[x+1][y].Type != particles.EMPTY {
			break
		}
	}
	s.grid[startX][y], s.grid[x][y] = s.grid[x][y], s.grid[startX][y]
	//s.grid[x][y].Velocity += particle.Gravity
	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(startX, y), CurrentPosition: *particles.NewPoint(x, y)}
}

func (s *Simulation) handleMoveLeft(x int, y int) *particles.MovedParticle {
	startX := x
	particle := s.grid[x][y]
	for ; x > startX-particle.Velocity && x > 1; x-- {
		if s.grid[x-1][y].Type != particles.EMPTY {
			break
		}
	}
	s.grid[startX][y], s.grid[x][y] = s.grid[x][y], s.grid[startX][y]
	//s.grid[x][y].Velocity += particle.Gravity
	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(startX, y), CurrentPosition: *particles.NewPoint(x, y)}
}
