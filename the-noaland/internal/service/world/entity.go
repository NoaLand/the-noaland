package world

type EntityID string

type Entity struct {
	ID    EntityID
	State int
}

func NewEntity(id EntityID) *Entity {
	return &Entity{
		ID:    id,
		State: 0,
	}
}

func (entity *Entity) Step() {
	entity.State++
}
