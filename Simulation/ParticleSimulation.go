package simulation

import (
	particles "FallingSand/Particles"
	flameparticle "FallingSand/Simulation/FlameParticle"
)

const (
	gridWidth  = 1280
	gridHeight = 720
)

type Simulation struct {
	grid           [][]particles.Particle
	buffer         [][]particles.Particle
	addedParticles []particles.MovedParticle
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

		addedParticle := &particles.MovedParticle{
			PreviousPosition: *particles.NewPoint(particle.Position.X, particle.Position.Y),
			CurrentPosition:  *particles.NewPoint(particle.Position.X, particle.Position.Y),
			Type:             particle.Type,
		}
		s.addedParticles = append(s.addedParticles, *addedParticle)
	}
}

func (s *Simulation) NextStep() []particles.MovedParticle {
	var movedParticles []particles.MovedParticle = s.addedParticles
	s.addedParticles = make([]particles.MovedParticle, 0)

	for y := gridHeight - 2; y >= 0; y-- {
		for x := 1; x < gridWidth-1; x++ {
			var movedParticle *particles.MovedParticle
			var Type particles.ParticleType = particles.EMPTY
			switch s.grid[x][y].Type {
			case particles.SAND:
				movedParticle = s.handleSandParticle(x, y)
				Type = particles.SAND
			case particles.WATER:
				movedParticle = s.handleWaterParticle(x, y)
				Type = particles.WATER
			case particles.FLAME:
				movedParticle = s.handleFlameParticle(x, y)
				if movedParticle != nil {
					movedParticle.Type = Type
					movedParticles = append(movedParticles, *movedParticle)
				}
				movedParticle = s.handleFlameParticle(x, y)
				if movedParticle != nil {
					movedParticle.Type = Type
					movedParticles = append(movedParticles, *movedParticle)
				}
				movedParticle = s.handleFlameParticle(x, y)
				if movedParticle != nil {
					movedParticle.Type = Type
					movedParticles = append(movedParticles, *movedParticle)
				}
				movedParticle = s.handleFlameParticle(x, y)
				Type = particles.FLAME
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
	if s.grid[x][y+1].Type == particles.EMPTY {
		return s.handleDownMove(x, y)
	} else if s.grid[x-1][y+1].Type == particles.EMPTY {
		return s.handleDiagonalLeftDown(x, y)
	} else if s.grid[x+1][y+1].Type == particles.EMPTY {
		return s.handleDiagonalRigthDown(x, y)
	} else {
		s.grid[x][y].Velocity.Y = 1
	}

	return nil
}

func (s *Simulation) handleWaterParticle(x int, y int) *particles.MovedParticle {
	if s.grid[x][y+1].Type == particles.EMPTY {
		return s.handleDownMove(x, y)
	} else if s.grid[x-1][y+1].Type == particles.EMPTY {
		return s.handleDiagonalLeftDown(x, y)
	} else if s.grid[x+1][y+1].Type == particles.EMPTY {
		return s.handleDiagonalRigthDown(x, y)
	} else if s.grid[x-1][y].Type == particles.EMPTY {
		return s.handleMoveLeft(x, y)
	} else if s.grid[x+1][y].Type == particles.EMPTY {
		return s.handleMoveRight(x, y)
	}
	return nil
}

func (s *Simulation) handleFlameParticle(x int, y int) *particles.MovedParticle {
	if s.grid[x][y-1].FlameRate != 0 {
		return flameparticle.HandleFlameUp(x, y, s.grid)
	}
	if s.grid[x][y+1].FlameRate != 0 {
		return flameparticle.HandleFlameDown(x, y, s.grid)
	}
	if s.grid[x-1][y].FlameRate != 0 {
		return flameparticle.HandleFlameLeft(x, y, s.grid)
	}
	if s.grid[x+1][y].FlameRate != 0 {
		return flameparticle.HandleFlameRight(x, y, s.grid)
	}
	return nil
}

func (s *Simulation) handleDownMove(x int, y int) *particles.MovedParticle {
	startY := y
	particle := s.grid[x][y]
	for ; y < startY+particle.Velocity.Y && y < gridHeight-1; y++ {
		if s.grid[x][y+1].Type != particles.EMPTY {
			break
		}
	}
	s.grid[x][startY], s.grid[x][y] = s.grid[x][y], s.grid[x][startY]
	s.grid[x][y].Velocity.Y += particle.Gravity
	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(x, startY), CurrentPosition: *particles.NewPoint(x, y)}
}

func (s *Simulation) handleMoveRight(x int, y int) *particles.MovedParticle {
	startX := x
	particle := s.grid[x][y]
	for ; x < startX+particle.Velocity.X && x < gridWidth-1; x++ {
		if s.grid[x+1][y].Type != particles.EMPTY {
			break
		}
	}
	s.grid[startX][y], s.grid[x][y] = s.grid[x][y], s.grid[startX][y]
	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(startX, y), CurrentPosition: *particles.NewPoint(x, y)}
}

func (s *Simulation) handleMoveLeft(x int, y int) *particles.MovedParticle {
	startX := x
	particle := s.grid[x][y]
	for ; x > startX-particle.Velocity.X && x > 1; x-- {
		if s.grid[x-1][y].Type != particles.EMPTY {
			break
		}
	}
	s.grid[startX][y], s.grid[x][y] = s.grid[x][y], s.grid[startX][y]
	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(startX, y), CurrentPosition: *particles.NewPoint(x, y)}
}

func (s *Simulation) handleDiagonalLeftDown(x int, y int) *particles.MovedParticle {
	startX := x
	startY := y
	particle := s.grid[x][y]

	for i := 1; i < particle.Velocity.Y && y+i < gridHeight-1 && x-i > 1; i++ {
		x -= 1
		y += 1
		if s.grid[x-1][y+1].Type != particles.EMPTY {
			break
		}
	}

	s.grid[startX][startY], s.grid[x][y] = s.grid[x][y], s.grid[startX][startY]
	s.grid[x][y].Velocity.Y += particle.Gravity
	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(startX, startY), CurrentPosition: *particles.NewPoint(x, y)}
}

func (s *Simulation) handleDiagonalRigthDown(x int, y int) *particles.MovedParticle {
	startX := x
	startY := y
	particle := s.grid[x][y]

	for i := 1; i < particle.Velocity.Y && y+i < gridHeight-1 && x+i < gridWidth-1; i++ {
		x += 1
		y += 1
		if s.grid[x+1][y+1].Type != particles.EMPTY {
			break
		}
	}

	s.grid[startX][startY], s.grid[x][y] = s.grid[x][y], s.grid[startX][startY]
	s.grid[x][y].Velocity.Y += particle.Gravity
	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(startX, startY), CurrentPosition: *particles.NewPoint(x, y)}
}

func (s *Simulation) handleDiagonalLeftUp(x int, y int) *particles.MovedParticle {
	startX := x
	startY := y
	particle := s.grid[x][y]

	for i := 1; i < particle.Velocity.Y && y-i > 1 && x-i > 1; i++ {
		x -= 1
		y -= 1
		if s.grid[x-1][y-1].Type != particles.EMPTY {
			break
		}
	}

	s.grid[startX][startY], s.grid[x][y] = s.grid[x][y], s.grid[startX][startY]
	s.grid[x][y].Velocity.Y += particle.Gravity
	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(startX, startY), CurrentPosition: *particles.NewPoint(x, y)}
}

func (s *Simulation) handleDiagonalRightUp(x int, y int) *particles.MovedParticle {
	startX := x
	startY := y
	particle := s.grid[x][y]

	for i := 1; i < particle.Velocity.Y && y-i > 1 && x+i < gridWidth-1; i++ {
		x += 1
		y -= 1
		if s.grid[x+1][y-1].Type != particles.EMPTY {
			break
		}
	}

	s.grid[startX][startY], s.grid[x][y] = s.grid[x][y], s.grid[startX][startY]
	s.grid[x][y].Velocity.Y += particle.Gravity
	return &particles.MovedParticle{PreviousPosition: *particles.NewPoint(startX, startY), CurrentPosition: *particles.NewPoint(x, y)}
}
