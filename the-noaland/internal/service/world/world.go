package world

import (
	"fmt"
	"math/rand/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity"
	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity/worldtree"
)

type World struct {
	entities        []entity.Entity
	seed            uint64
	historiographer *historiographer
	timer           uint64
	rng             *rand.Rand
}

func New(seed uint64) *World {
	return &World{
		seed:            seed,
		historiographer: NewHistoriographer(),
		timer:           0,
		rng:             rand.New(rand.NewPCG(seed, 0)),
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

func (world *World) Generate() {
	rng := rand.New(rand.NewPCG(world.seed, 0))

	count := 1 + rng.IntN(5)

	for i := range count {
		id := entity.EntityID(fmt.Sprintf("world-tree-%d", i))

		tree := worldtree.New(id, generateWorldTreeName(world.rng), world.seed)
		world.Add(tree)
	}
}
