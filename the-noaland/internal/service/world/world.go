package world

type World struct {
	entities []*Entity
}

func New() *World {
	return &World{}
}

func (world *World) Add(entity *Entity) {
	world.entities = append(world.entities, entity)
}

func (world *World) Step() {
	for _, entity := range world.entities {
		entity.Step()
	}
}
