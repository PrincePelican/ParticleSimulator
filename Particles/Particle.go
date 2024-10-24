package partjicles

type Vector2D struct {
	X, Y int
}

type ChangeParticle struct {
	PreviousPosition, CurrentPosition Vector2D
	Type                              ParticleType
	ChangeType                        ChangeType
}

func NewChangeParticle(PreviousPosition Vector2D, CurrentPosition Vector2D, Type ParticleType, ChangeType ChangeType) *ChangeParticle {
	return &ChangeParticle{
		PreviousPosition: PreviousPosition,
		CurrentPosition:  CurrentPosition,
		Type:             Type,
		ChangeType:       ChangeType,
	}
}

func NewPoint(x int, y int) *Vector2D {
	return &Vector2D{
		X: x,
		Y: y,
	}
}

type Particle struct {
	Position  Vector2D
	Gravity   int
	Velocity  Vector2D
	Type      ParticleType
	FlameTime int
	FlameRate int
}

func NewParticle(xPosition int, yPosition int, Type ParticleType) *Particle {
	position := NewPoint(xPosition, yPosition)

	return &Particle{
		Position:  *position,
		Type:      Type,
		Velocity:  *NewPoint(10, 5),
		Gravity:   1,
		FlameTime: 30,
		FlameRate: getFlameRate(Type),
	}
}

func NewSmokeParticle(xPosition int, yPosition int, Type ParticleType) *Particle {
	position := NewPoint(xPosition, yPosition)

	return &Particle{
		Position:  *position,
		Type:      Type,
		Velocity:  *NewPoint(30, 2),
		Gravity:   0,
		FlameTime: 0,
		FlameRate: getFlameRate(Type),
	}
}

type ParticleType int

func getFlameRate(Type ParticleType) int {
	switch Type {
	case SAND:
		return 0
	case WATER:
		return 0
	case WOOD:
		return 1
	default:
		return 0
	}
}

const (
	EMPTY ParticleType = 0
	SAND  ParticleType = 1
	WATER ParticleType = 2
	WOOD  ParticleType = 3
	FLAME ParticleType = 4
	SMOKE ParticleType = 5
)

type ChangeType int

const (
	MOVE   ChangeType = 0
	CHANGE ChangeType = 1
	VANISH ChangeType = 2
)
