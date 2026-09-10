package world

type entityID string

type Entity struct {
	id    entityID
	state int
}

func NewEntity(id entityID) *Entity {
	return &Entity{
		id:    id,
		state: 0,
	}
}

func (entity *Entity) Step() {
	entity.state++
}
