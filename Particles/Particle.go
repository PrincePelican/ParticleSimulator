package particles

type Point struct {
	X, Y int
}

type MovedParticle struct {
	PreviousPosition, CurrentPosition Point
	Type                              ParticleType
}

func NewMovedParticle(PreviousPosition Point, CurrentPosition Point, Type ParticleType) *MovedParticle {
	return &MovedParticle{
		PreviousPosition: PreviousPosition,
		CurrentPosition:  CurrentPosition,
		Type:             Type,
	}
}

func NewPoint(x int, y int) *Point {
	return &Point{
		X: x,
		Y: y,
	}
}

type Particle struct {
	Position Point
	Gravity  int
	Velocity int
	Type     ParticleType
}

func NewParticle(xPosition int, yPosition int, Type ParticleType) *Particle {
	position := NewPoint(xPosition, yPosition)
	return &Particle{
		Position: *position,
		Type:     Type,
		Velocity: 1,
		Gravity:  1,
	}
}

type ParticleType int

const (
	EMPTY ParticleType = 0
	SAND  ParticleType = 1
	WATER ParticleType = 2
)
