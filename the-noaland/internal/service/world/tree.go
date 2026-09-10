package world

type Tree struct {
	identity entityID
	state    int
}

func NewTree(id entityID) *Tree {
	return &Tree{
		identity: id,
	}
}

func (t *Tree) id() entityID {
	return t.identity
}

func (t *Tree) step() {
	t.state++
}
