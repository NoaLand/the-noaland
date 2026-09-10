package world

type WorldTree struct {
	identity entityID
	state    int
}

func NewWorldTree(id entityID) *WorldTree {
	return &WorldTree{
		identity: id,
	}
}

func (t *WorldTree) id() entityID {
	return t.identity
}

func (t *WorldTree) step() {
	t.state++
}
