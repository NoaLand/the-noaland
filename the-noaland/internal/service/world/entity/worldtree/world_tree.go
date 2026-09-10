package worldtree

import "github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity"

type worldTreeState struct {
	growth int
}

type WorldTree struct {
	id    entity.EntityID
	state worldTreeState
}

func New(id entity.EntityID) *WorldTree {
	return &WorldTree{
		id: id,
		state: worldTreeState{
			growth: 0,
		},
	}
}

func (t *WorldTree) ID() entity.EntityID {
	return t.id
}

func (t *WorldTree) Step() {
	t.state.growth++
}

type Expressions struct {
	Appearance worldTreeAppearance
}

func (t *WorldTree) Express() *Expressions {
	return &Expressions{
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
