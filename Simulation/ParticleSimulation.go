package simulation

import (
	particles "FallingSand/Particles"
	flameparticle "FallingSand/Simulation/FlameParticle"
	"math/rand"
)

const (
	gridWidth  = 720
	gridHeight = 360
)

type Simulation struct {
	grid           [][]particles.Particle
	addedParticles []particles.ChangeParticle
	indexes        []particles.Vector2D
}

func NewSimulation() *Simulation {
	return &Simulation{
		grid: make([][]particles.Particle, gridWidth),
	}
}

func (s *Simulation) Initialize() {
	for i := range s.grid {
		s.grid[i] = make([]particles.Particle, gridHeight)
		for j := range s.grid[i] {
			if j < gridHeight-1 && i < gridWidth-1 {
				s.indexes = append(s.indexes, particles.Vector2D{X: i, Y: j})
			}

		}
	}
}

func (s *Simulation) AddParticles(newParticlesList []particles.Particle) {
	for x := 0; x < len(newParticlesList); x++ {
		particle := newParticlesList[x]
		s.grid[particle.Position.X][particle.Position.Y] = particle

		addedParticle := &particles.ChangeParticle{
			PreviousPosition: *particles.NewPoint(particle.Position.X, particle.Position.Y),
			CurrentPosition:  *particles.NewPoint(particle.Position.X, particle.Position.Y),
			Type:             particle.Type,
		}
		s.addedParticles = append(s.addedParticles, *addedParticle)
	}
}

func (s *Simulation) randomizeListOfIndexes() []particles.Vector2D {
	var indexes []particles.Vector2D = s.indexes
	rand.Shuffle(len(indexes), func(x, y int) {
		indexes[x], indexes[y] = indexes[y], indexes[x]
	})
	return indexes
}

func (s *Simulation) NextStep() []particles.ChangeParticle {
	var movedParticles []particles.ChangeParticle = s.addedParticles
	s.addedParticles = make([]particles.ChangeParticle, 0)
	var indexTmp = s.randomizeListOfIndexes()

	for x := 0; x < len(indexTmp); x++ {
		var changesParticle []particles.ChangeParticle
		switch s.grid[indexTmp[x].X][indexTmp[x].Y].Type {
		case particles.SAND:
			changesParticle = s.handleSandParticle(indexTmp[x].X, indexTmp[x].Y)
		case particles.WATER:
			changesParticle = s.handleWaterParticle(indexTmp[x].X, indexTmp[x].Y)
		case particles.FLAME:
			changesParticle = s.handleFlameParticle(indexTmp[x].X, indexTmp[x].Y)
		case particles.SMOKE:
			changesParticle = s.handleSmokeParticle(indexTmp[x].X, indexTmp[x].Y)
		}

		if changesParticle != nil {
			movedParticles = append(movedParticles, changesParticle...)
		}
	}

	return movedParticles
}

func (s *Simulation) handleSandParticle(x int, y int) []particles.ChangeParticle {
	var changes []particles.ChangeParticle
	if s.grid[x][y+1].Type == particles.EMPTY {
		change := s.handleMoveDown(x, y)
		change.Type = particles.SAND
		return append(changes, *change)
	} else if s.grid[x-1][y+1].Type == particles.EMPTY {
		change := s.handleDiagonalLeftDown(x, y)
		change.Type = particles.SAND
		return append(changes, *change)
	} else if s.grid[x+1][y+1].Type == particles.EMPTY {
		change := s.handleDiagonalRigthDown(x, y)
		change.Type = particles.SAND
		return append(changes, *change)
	} else {
		s.grid[x][y].Velocity.Y = 1
	}

	return nil
}

func (s *Simulation) handleWaterParticle(x int, y int) []particles.ChangeParticle {
	var changes []particles.ChangeParticle
	if s.grid[x][y+1].Type == particles.EMPTY {
		change := s.handleMoveDown(x, y)
		change.Type = particles.WATER
		return append(changes, *change)
	} else if s.grid[x-1][y+1].Type == particles.EMPTY {
		change := s.handleDiagonalLeftDown(x, y)
		change.Type = particles.WATER
		return append(changes, *change)
	} else if s.grid[x+1][y+1].Type == particles.EMPTY {
		change := s.handleDiagonalRigthDown(x, y)
		change.Type = particles.WATER
		return append(changes, *change)
	} else if s.grid[x-1][y].Type == particles.EMPTY {
		change := s.handleMoveLeft(x, y)
		change.Type = particles.WATER
		return append(changes, *change)
	} else if s.grid[x+1][y].Type == particles.EMPTY {
		change := s.handleMoveRight(x, y)
		change.Type = particles.WATER
		return append(changes, *change)
	}
	return nil
}

func (s *Simulation) handleFlameParticle(x int, y int) []particles.ChangeParticle {
	var changes []particles.ChangeParticle
	if s.grid[x][y-1].FlameRate != 0 {
		change := flameparticle.HandleFlameUp(x, y, s.grid)
		if change != nil {
			changes = append(changes, *change)
		}
	}
	if s.grid[x][y+1].FlameRate != 0 {
		change := flameparticle.HandleFlameDown(x, y, s.grid)
		if change != nil {
			changes = append(changes, *change)
		}
	}
	if s.grid[x-1][y].FlameRate != 0 {
		change := flameparticle.HandleFlameLeft(x, y, s.grid)
		if change != nil {
			changes = append(changes, *change)
		}
	}
	if s.grid[x+1][y].FlameRate != 0 {
		change := flameparticle.HandleFlameRight(x, y, s.grid)
		if change != nil {
			changes = append(changes, *change)
		}
	}
	change := flameparticle.TickFlameTime(&s.grid[x][y], s.grid)
	if change != nil {
		changes = append(changes, *change)
	}

	return changes
}

func (s *Simulation) handleSmokeParticle(x int, y int) []particles.ChangeParticle {
	var changes []particles.ChangeParticle
	if s.grid[x][y-1].Type == particles.EMPTY {
		change := s.handleMoveUp(x, y)
		change.Type = particles.SMOKE
		return append(changes, *change)
	} else if s.grid[x-1][y-1].Type == particles.EMPTY {
		change := s.handleDiagonalLeftUp(x, y)
		change.Type = particles.SMOKE
		return append(changes, *change)
	} else if s.grid[x+1][y-1].Type == particles.EMPTY {
		change := s.handleDiagonalRightUp(x, y)
		change.Type = particles.SMOKE
		return append(changes, *change)
	} else if s.grid[x-1][y].Type == particles.EMPTY {
		change := s.handleMoveLeft(x, y)
		change.Type = particles.SMOKE
		return append(changes, *change)
	} else if s.grid[x+1][y].Type == particles.EMPTY {
		change := s.handleMoveRight(x, y)
		change.Type = particles.SMOKE
		return append(changes, *change)
	}
	return nil
}

func (s *Simulation) handleMoveDown(x int, y int) *particles.ChangeParticle {
	startY := y
	particle := s.grid[x][y]
	for ; y < startY+particle.Velocity.Y && y < gridHeight-1; y++ {
		if s.grid[x][y+1].Type != particles.EMPTY {
			break
		}
	}
	s.grid[x][startY], s.grid[x][y] = s.grid[x][y], s.grid[x][startY]
	s.grid[x][y].Velocity.Y += particle.Gravity
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(x, startY), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}
}

func (s *Simulation) handleMoveUp(x int, y int) *particles.ChangeParticle {
	startY := y
	particle := s.grid[x][y]
	for ; y > startY-particle.Velocity.Y && y > 1; y-- {
		if s.grid[x][y-1].Type != particles.EMPTY {
			break
		}
	}
	s.grid[x][startY], s.grid[x][y] = s.grid[x][y], s.grid[x][startY]
	s.grid[x][y].Velocity.Y += particle.Gravity
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(x, startY), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}
}

func (s *Simulation) handleMoveRight(x int, y int) *particles.ChangeParticle {
	startX := x
	particle := s.grid[x][y]
	for ; x < startX+particle.Velocity.X && x < gridWidth-1; x++ {
		if s.grid[x+1][y].Type != particles.EMPTY {
			break
		}
	}
	s.grid[startX][y], s.grid[x][y] = s.grid[x][y], s.grid[startX][y]
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(startX, y), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}
}

func (s *Simulation) handleMoveLeft(x int, y int) *particles.ChangeParticle {
	startX := x
	particle := s.grid[x][y]
	for ; x > startX-particle.Velocity.X && x > 1; x-- {
		if s.grid[x-1][y].Type != particles.EMPTY {
			break
		}
	}
	s.grid[startX][y], s.grid[x][y] = s.grid[x][y], s.grid[startX][y]
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(startX, y), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}
}

func (s *Simulation) handleDiagonalLeftDown(x int, y int) *particles.ChangeParticle {
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
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(startX, startY), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}
}

func (s *Simulation) handleDiagonalRigthDown(x int, y int) *particles.ChangeParticle {
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
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(startX, startY), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}
}

func (s *Simulation) handleDiagonalLeftUp(x int, y int) *particles.ChangeParticle {
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
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(startX, startY), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}
}

func (s *Simulation) handleDiagonalRightUp(x int, y int) *particles.ChangeParticle {
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
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(startX, startY), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}
}
