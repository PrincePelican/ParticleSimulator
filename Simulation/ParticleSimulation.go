package simulation

import (
	particles "FallingSand/Particles"
	flameparticle "FallingSand/Simulation/FlameParticle"
	"math/rand"
)

const (
	gridWidth  = 2560
	gridHeight = 1440
)

type Simulation struct {
	grid            [][]particles.Particle
	addedParticles  []particles.ChangeParticle
	indexes         []particles.Vector2D
	movedParticles  []particles.ChangeParticle
	activeParticles []particles.Vector2D
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
	s.movedParticles = make([]particles.ChangeParticle, 0, gridHeight*gridWidth)

	go s.randomizeListOfIndexes()
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
		s.activeParticles = append(s.activeParticles, addedParticle.CurrentPosition)
	}
}

func (s *Simulation) randomizeListOfIndexes() {
	rand.Shuffle(len(s.indexes), func(x, y int) {
		s.indexes[x], s.indexes[y] = s.indexes[y], s.indexes[x]
	})
}

func (s *Simulation) NextStep() []particles.ChangeParticle {
	s.movedParticles = s.movedParticles[:0]
	s.movedParticles = append(s.movedParticles, s.addedParticles...)
	s.addedParticles = s.addedParticles[:0]
	var currentActive = s.activeParticles
	s.activeParticles = s.activeParticles[:0]

	for x := 0; x < len(currentActive); x++ {
		var changesParticles []particles.ChangeParticle
		var newActiveParticles []particles.Vector2D
		switch s.grid[currentActive[x].X][currentActive[x].Y].Type {
		case particles.SAND:
			changesParticles, newActiveParticles = s.handleSandParticle(currentActive[x].X, currentActive[x].Y)
		case particles.WATER:
			changesParticles, newActiveParticles = s.handleWaterParticle(currentActive[x].X, currentActive[x].Y)
		case particles.FLAME:
			changesParticles, newActiveParticles = s.handleFlameParticle(currentActive[x].X, currentActive[x].Y)
		case particles.SMOKE:
			changesParticles, newActiveParticles = s.handleSmokeParticle(currentActive[x].X, currentActive[x].Y)
		}

		if changesParticles != nil {
			s.movedParticles = append(s.movedParticles, changesParticles...)
		}

		if (newActiveParticles) != nil {
			s.activeParticles = append(s.activeParticles, newActiveParticles...)
		}
	}
	return s.movedParticles
}

func (s *Simulation) handleSandParticle(x int, y int) ([]particles.ChangeParticle, []particles.Vector2D) {
	var changes []particles.ChangeParticle
	var newActives []particles.Vector2D
	if s.grid[x][y+1].Type == particles.EMPTY {
		change, newActive := s.handleMoveDown(x, y)
		change.Type = particles.SAND
		return append(changes, *change), append(newActives, *newActive)
	} else if s.grid[x-1][y+1].Type == particles.EMPTY {
		change, newActive := s.handleDiagonalLeftDown(x, y)
		change.Type = particles.SAND
		return append(changes, *change), append(newActives, *newActive)
	} else if s.grid[x+1][y+1].Type == particles.EMPTY {
		change, newActive := s.handleDiagonalRigthDown(x, y)
		change.Type = particles.SAND
		return append(changes, *change), append(newActives, *newActive)
	} else {
		s.grid[x][y].Velocity.Y = 1
	}

	return nil, nil
}

func (s *Simulation) handleWaterParticle(x int, y int) ([]particles.ChangeParticle, []particles.Vector2D) {
	var changes []particles.ChangeParticle
	var newActives []particles.Vector2D
	if s.grid[x][y+1].Type == particles.EMPTY {
		change, newActive := s.handleMoveDown(x, y)
		change.Type = particles.WATER
		return append(changes, *change), append(newActives, *newActive)
	} else if s.grid[x-1][y+1].Type == particles.EMPTY {
		change, newActive := s.handleDiagonalLeftDown(x, y)
		change.Type = particles.WATER
		return append(changes, *change), append(newActives, *newActive)
	} else if s.grid[x+1][y+1].Type == particles.EMPTY {
		change, newActive := s.handleDiagonalRigthDown(x, y)
		change.Type = particles.WATER
		return append(changes, *change), append(newActives, *newActive)
	} else if s.grid[x-1][y].Type == particles.EMPTY {
		change, newActive := s.handleMoveLeft(x, y)
		change.Type = particles.WATER
		return append(changes, *change), append(newActives, *newActive)
	} else if s.grid[x+1][y].Type == particles.EMPTY {
		change, newActive := s.handleMoveRight(x, y)
		change.Type = particles.WATER
		return append(changes, *change), append(newActives, *newActive)
	}
	return nil, nil
}

func (s *Simulation) handleFlameParticle(x int, y int) ([]particles.ChangeParticle, []particles.Vector2D) {
	var changes []particles.ChangeParticle
	var newActives []particles.Vector2D

	if s.grid[x][y-1].FlameRate != 0 {
		change, newActive := flameparticle.HandleFlameUp(x, y, s.grid)
		if change != nil {
			changes, newActives = append(changes, *change), append(newActives, *newActive)
		}
	}
	if s.grid[x][y+1].FlameRate != 0 {
		change, newActive := flameparticle.HandleFlameDown(x, y, s.grid)
		if change != nil {
			changes, newActives = append(changes, *change), append(newActives, *newActive)
		}
	}
	if s.grid[x-1][y].FlameRate != 0 {
		change, newActive := flameparticle.HandleFlameLeft(x, y, s.grid)
		if change != nil {
			changes, newActives = append(changes, *change), append(newActives, *newActive)
		}
	}
	if s.grid[x+1][y].FlameRate != 0 {
		change, newActive := flameparticle.HandleFlameRight(x, y, s.grid)
		if change != nil {
			changes, newActives = append(changes, *change), append(newActives, *newActive)
		}
	}
	change, newActive := flameparticle.TickFlameTime(&s.grid[x][y], s.grid)
	if change != nil {
		changes, newActives = append(changes, *change), append(newActives, *newActive)
	}

	return changes, newActives
}

func (s *Simulation) handleSmokeParticle(x int, y int) ([]particles.ChangeParticle, []particles.Vector2D) {
	var changes []particles.ChangeParticle
	var newActives []particles.Vector2D
	if s.grid[x][y-1].Type == particles.EMPTY {
		change, newActive := s.handleMoveUp(x, y)
		change.Type = particles.SMOKE
		return append(changes, *change), append(newActives, *newActive)
	} else if s.grid[x-1][y-1].Type == particles.EMPTY {
		change, newActive := s.handleDiagonalLeftUp(x, y)
		change.Type = particles.SMOKE
		return append(changes, *change), append(newActives, *newActive)
	} else if s.grid[x+1][y-1].Type == particles.EMPTY {
		change, newActive := s.handleDiagonalRightUp(x, y)
		change.Type = particles.SMOKE
		return append(changes, *change), append(newActives, *newActive)
	} else if s.grid[x-1][y].Type == particles.EMPTY {
		change, newActive := s.handleMoveLeft(x, y)
		change.Type = particles.SMOKE
		return append(changes, *change), append(newActives, *newActive)
	} else if s.grid[x+1][y].Type == particles.EMPTY {
		change, newActive := s.handleMoveRight(x, y)
		change.Type = particles.SMOKE
		return append(changes, *change), append(newActives, *newActive)
	}
	return nil, nil
}

func (s *Simulation) handleMoveDown(x int, y int) (*particles.ChangeParticle, *particles.Vector2D) {
	startY := y
	particle := s.grid[x][y]
	for ; y < startY+particle.Velocity.Y && y < gridHeight-2; y++ {
		if s.grid[x][y+1].Type != particles.EMPTY {
			break
		}
	}
	s.grid[x][startY], s.grid[x][y] = s.grid[x][y], s.grid[x][startY]
	s.grid[x][y].Velocity.Y += particle.Gravity
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(x, startY), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}, particles.NewPoint(x, y)
}

func (s *Simulation) handleMoveUp(x int, y int) (*particles.ChangeParticle, *particles.Vector2D) {
	startY := y
	particle := s.grid[x][y]
	for ; y > startY-particle.Velocity.Y && y > 1; y-- {
		if s.grid[x][y-1].Type != particles.EMPTY {
			break
		}
	}
	s.grid[x][startY], s.grid[x][y] = s.grid[x][y], s.grid[x][startY]
	s.grid[x][y].Velocity.Y += particle.Gravity
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(x, startY), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}, particles.NewPoint(x, y)
}

func (s *Simulation) handleMoveRight(x int, y int) (*particles.ChangeParticle, *particles.Vector2D) {
	startX := x
	particle := s.grid[x][y]
	for ; x < startX+particle.Velocity.X && x < gridWidth-2; x++ {
		if s.grid[x+1][y].Type != particles.EMPTY {
			break
		}
	}
	s.grid[startX][y], s.grid[x][y] = s.grid[x][y], s.grid[startX][y]
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(startX, y), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}, particles.NewPoint(x, y)
}

func (s *Simulation) handleMoveLeft(x int, y int) (*particles.ChangeParticle, *particles.Vector2D) {
	startX := x
	particle := s.grid[x][y]
	for ; x > startX-particle.Velocity.X && x > 1; x-- {
		if s.grid[x-1][y].Type != particles.EMPTY {
			break
		}
	}
	s.grid[startX][y], s.grid[x][y] = s.grid[x][y], s.grid[startX][y]
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(startX, y), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}, particles.NewPoint(x, y)

}

func (s *Simulation) handleDiagonalLeftDown(x int, y int) (*particles.ChangeParticle, *particles.Vector2D) {
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
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(startX, startY), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}, particles.NewPoint(x, y)
}

func (s *Simulation) handleDiagonalRigthDown(x int, y int) (*particles.ChangeParticle, *particles.Vector2D) {
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
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(startX, startY), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}, particles.NewPoint(x, y)
}

func (s *Simulation) handleDiagonalLeftUp(x int, y int) (*particles.ChangeParticle, *particles.Vector2D) {
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
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(startX, startY), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}, particles.NewPoint(x, y)
}

func (s *Simulation) handleDiagonalRightUp(x int, y int) (*particles.ChangeParticle, *particles.Vector2D) {
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
	return &particles.ChangeParticle{PreviousPosition: *particles.NewPoint(startX, startY), CurrentPosition: *particles.NewPoint(x, y), ChangeType: particles.MOVE}, particles.NewPoint(x, y)
}
