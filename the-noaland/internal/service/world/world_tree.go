package world

type worldTreeState struct {
	growth int
}

type WorldTree struct {
	identity entityID
	state    worldTreeState
}

func NewWorldTree(id entityID) *WorldTree {
	return &WorldTree{
		identity: id,
		state: worldTreeState{
			growth: 0,
		},
	}
}

func (t *WorldTree) id() entityID {
	return t.identity
}

func (t *WorldTree) step() {
	t.state.growth++
}

type worldTreeAppearance int

const (
	WorldTreeSeed worldTreeAppearance = iota
	WorldTreeSprout
	WorldTreeYoung
	WorldTreeMature
)

func (t *WorldTree) Appearance() worldTreeAppearance {
	switch {
	case t.state.growth < 3:
		return WorldTreeSeed
	case t.state.growth < 6:
		return WorldTreeSprout
	case t.state.growth < 10:
		return WorldTreeYoung
	default:
		return WorldTreeMature
	}
}

func (s worldTreeAppearance) String() string {
	switch s {
	case WorldTreeSeed:
		return "Seed"
	case WorldTreeSprout:
		return "Sprout"
	case WorldTreeYoung:
		return "Young"
	case WorldTreeMature:
		return "Mature"
	default:
		return "Unknown"
	}
}
