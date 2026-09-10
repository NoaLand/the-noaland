package entity

type EntityID string

type Entity interface {
	ID() EntityID
	Step()
}
