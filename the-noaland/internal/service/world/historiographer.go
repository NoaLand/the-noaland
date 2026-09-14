package world

import (
	"fmt"
	"slices"

	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity"
	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity/worldtree"
)

// Record describes an observed change at a world time step.
type Record struct {
	Time       uint64
	EntityID   entity.EntityID
	EntityName entity.EntityName
	Message    string
}

type Annal struct {
	records []Record
}

// Records returns a chronological snapshot. Changes to it do not affect the annal.
func (a *Annal) Records() []Record {
	return slices.Clone(a.records)
}

type historiographer struct {
	last  map[entity.EntityID]any
	annal Annal
}

func NewHistoriographer() *historiographer {
	return &historiographer{
		last: make(map[entity.EntityID]any),
	}
}

func (h *historiographer) record(time uint64, entity entity.Entity) {
	switch v := entity.(type) {
	case *worldtree.WorldTree:
		h.recordWorldTree(time, v)
	}
}

func (h *historiographer) recordWorldTree(time uint64, tree *worldtree.WorldTree) {
	current := tree.Express()

	previous, ok := h.last[tree.ID()]
	if !ok {
		h.last[tree.ID()] = current
		return
	}

	previousExpression, ok := previous.(*worldtree.Expressions)
	if !ok {
		return
	}

	if previousExpression.Appearance != current.Appearance {
		record := Record{
			Time:       time,
			EntityID:   tree.ID(),
			EntityName: tree.Name(),
			Message:    worldTreeMessage(previousExpression, current),
		}

		h.annal.records = append(h.annal.records, record)
	}

	h.last[tree.ID()] = current
}

func worldTreeMessage(previous, current *worldtree.Expressions) string {
	return fmt.Sprintf(
		"%s changed from %s to %s.",
		current.Name,
		previous.Appearance,
		current.Appearance,
	)
}
