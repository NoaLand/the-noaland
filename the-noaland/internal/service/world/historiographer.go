package world

import (
	"fmt"

	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity"
	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity/worldtree"
)

type record struct {
	time       uint64
	entityId   entity.EntityID
	entityName entity.EntityName
	message    string
}

type Annal struct {
	records []record
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
		record := record{
			time:       time,
			entityId:   tree.ID(),
			entityName: tree.Name(),
			message:    worldTreeMessage(previousExpression, current),
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
