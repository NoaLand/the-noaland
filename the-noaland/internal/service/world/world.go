package world

import (
	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity"
)

type World struct {
	entities        []entity.Entity
	seed            uint64
	historiographer *historiographer
	timer           uint64
}

func New(seed uint64) *World {
	return &World{
		seed:            seed,
		historiographer: NewHistoriographer(),
		timer:           0,
	}
}

func (world *World) Add(entity entity.Entity) {
	world.entities = append(world.entities, entity)
	world.historiographer.record(world.timer, entity)
}

func (world *World) Step() {
	world.timer++

	for _, entity := range world.entities {
		entity.Step()
		world.historiographer.record(world.timer, entity)
	}
}

func (world *World) Seed() uint64 {
	return world.seed
}
