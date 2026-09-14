package worldtree

import (
	"math/rand/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity"
)

type worldTreeState struct {
	growth int
}

type WorldTree struct {
	id    entity.EntityID
	name  entity.EntityName
	state worldTreeState
	rng   *rand.Rand
}

func New(id entity.EntityID, name entity.EntityName, worldSeed uint64) *WorldTree {
	seed := entity.DeriveSeed(worldSeed, id)

	return &WorldTree{
		id:   id,
		name: name,
		state: worldTreeState{
			growth: 0,
		},
		rng: rand.New(rand.NewPCG(seed, 0)),
	}
}

func (t *WorldTree) ID() entity.EntityID {
	return t.id
}

func (t *WorldTree) Name() entity.EntityName {
	return t.name
}

func (t *WorldTree) Step() {
	if t.rng.Float64() < 0.7 {
		t.state.growth++
	}
}

type Expressions struct {
	Name       entity.EntityName
	Appearance worldTreeAppearance
}

func (t *WorldTree) Express() *Expressions {
	return &Expressions{
		Name:       t.name,
		Appearance: t.appearance(),
	}
}

func (t *WorldTree) appearance() worldTreeAppearance {
	switch {
	case t.state.growth < 3:
		return Seed
	case t.state.growth < 6:
		return Sprout
	case t.state.growth < 10:
		return Young
	default:
		return Mature
	}
}
