package world

type worldTreeState int

const (
	WorldTreeSeed worldTreeState = iota
	WorldTreeSprout
	WorldTreeYoung
	WorldTreeMature
)

type WorldTree struct {
	identity entityID
	state    worldTreeState
}

func NewWorldTree(id entityID) *WorldTree {
	return &WorldTree{
		identity: id,
		state:    WorldTreeSeed,
	}
}

func (t *WorldTree) id() entityID {
	return t.identity
}

func (t *WorldTree) step() {
	switch t.state {
	case WorldTreeSeed:
		t.state = WorldTreeSprout
	case WorldTreeSprout:
		t.state = WorldTreeYoung
	case WorldTreeYoung:
		t.state = WorldTreeMature
	case WorldTreeMature:
		t.state = WorldTreeMature
	}
}

func (t *WorldTree) State() any {
	return t.state
}

func (s worldTreeState) String() string {
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
