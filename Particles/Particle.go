package particles

import "image/color"

type Point struct {
	X, Y int
}

type MovedParticle struct {
	PreviousPosition, CurrentPosition Point
}

func NewPoint(x int, y int) *Point {
	return &Point{
		X: x,
		Y: y,
	}
}

type Particle struct {
	Position Point
	Color    color.RGBA
	Gravity  float32
	Velocity float32
	Type     ParticleType
}

func NewParticle(xPosition int, yPosition int) *Particle {
	position := NewPoint(xPosition, yPosition)
	return &Particle{
		Position: *position,
		Type:     SAND,
	}
}

type ParticleType int

const (
	EMPTY ParticleType = 0
	SAND  ParticleType = 1
)
