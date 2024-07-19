package particles

type Vector2D struct {
	X, Y int
}

type MovedParticle struct {
	PreviousPosition, CurrentPosition Vector2D
	Type                              ParticleType
}

func NewMovedParticle(PreviousPosition Vector2D, CurrentPosition Vector2D, Type ParticleType) *MovedParticle {
	return &MovedParticle{
		PreviousPosition: PreviousPosition,
		CurrentPosition:  CurrentPosition,
		Type:             Type,
	}
}

func NewPoint(x int, y int) *Vector2D {
	return &Vector2D{
		X: x,
		Y: y,
	}
}

type Particle struct {
	Position Vector2D
	Gravity  int
	Velocity Vector2D
	Type     ParticleType
}

func NewParticle(xPosition int, yPosition int, Type ParticleType) *Particle {
	position := NewPoint(xPosition, yPosition)
	return &Particle{
		Position: *position,
		Type:     Type,
		Velocity: *NewPoint(10, 1),
		Gravity:  1,
	}
}

type ParticleType int

const (
	EMPTY ParticleType = 0
	SAND  ParticleType = 1
	WATER ParticleType = 2
)
