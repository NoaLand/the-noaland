package world

type entityID string

type Entity interface {
	id() entityID
	step()
	State() any
}
