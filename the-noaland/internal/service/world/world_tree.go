package world

type worldTreeAppearance int

const (
	WorldTreeSeed worldTreeAppearance = iota
	WorldTreeSprout
	WorldTreeYoung
	WorldTreeMature
)

type WorldTree struct {
	identity   entityID
	appearance worldTreeAppearance
}

func NewWorldTree(id entityID) *WorldTree {
	return &WorldTree{
		identity:   id,
		appearance: WorldTreeSeed,
	}
}

func (t *WorldTree) id() entityID {
	return t.identity
}

func (t *WorldTree) step() {
	switch t.appearance {
	case WorldTreeSeed:
		t.appearance = WorldTreeSprout
	case WorldTreeSprout:
		t.appearance = WorldTreeYoung
	case WorldTreeYoung:
		t.appearance = WorldTreeMature
	}
}

func (t *WorldTree) Appearance() worldTreeAppearance {
	return t.appearance
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
