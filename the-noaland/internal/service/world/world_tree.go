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
