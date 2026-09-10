package world

import "github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity"

type World struct {
	entities []entity.Entity
}

func New() *World {
	return &World{}
}

func (world *World) Add(entity entity.Entity) {
	world.entities = append(world.entities, entity)
}

func (world *World) Step() {
	for _, entity := range world.entities {
		entity.Step()
	}
}
