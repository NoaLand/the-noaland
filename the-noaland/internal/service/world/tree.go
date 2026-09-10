package world

type Tree struct {
	identity entityID
	state    int
}

func NewTree(id string) *Tree {
	return &Tree{
		identity: entityID(id),
	}
}

func (t *Tree) id() entityID {
	return t.identity
}

func (t *Tree) step() {
	t.state++
}
