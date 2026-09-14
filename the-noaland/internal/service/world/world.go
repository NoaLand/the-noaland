package world

import (
	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity"
)

type World struct {
	entities []entity.Entity
	seed     uint64
}

func New(seed uint64) *World {
	return &World{
		seed: seed,
	}
}

func (world *World) Add(entity entity.Entity) {
	world.entities = append(world.entities, entity)
}

func (world *World) Step() {
	for _, entity := range world.entities {
		entity.Step()
	}
}

func (world *World) Seed() uint64 {
	return world.seed
}
