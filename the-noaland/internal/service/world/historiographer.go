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

type annal struct {
	records []record
}

type historiographer struct {
	last  map[entity.EntityID]any
	annal annal
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
			time:     time,
			entityId: tree.ID(),
			entityName: tree.Name(),
			message:  worldTreeMessage(previousExpression, current),
		}

		h.annal.records = append(h.annal.records, record)

		fmt.Printf(
			"Time #%d - %s: %s\n",
			record.time,
			record.entityId,
			record.message,
		)
	}

	h.last[tree.ID()] = current
}

func worldTreeMessage(previous, current *worldtree.Expressions) string {
	return fmt.Sprintf(
		"The World Tree changed from %s to %s.",
		previous.Appearance,
		current.Appearance,
	)
}
