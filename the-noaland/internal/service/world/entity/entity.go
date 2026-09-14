package entity

import (
	"encoding/binary"
	"hash/fnv"
)

type EntityID string

type Entity interface {
	ID() EntityID
	Step()
}

func DeriveSeed(worldSeed uint64, id EntityID) uint64 {
	h := fnv.New64a()

	var buffer [8]byte
	binary.LittleEndian.PutUint64(buffer[:], worldSeed)

	_, _ = h.Write(buffer[:])
	_, _ = h.Write([]byte(id))

	return h.Sum64()
}
