package world

import (
	"fmt"
	"math/rand/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity"
)

var worldtreePrefixes = []string{
	"Ancient",
	"Silver",
	"Silent",
	"Golden",
	"Green",
}

var worldtreeNames = []string{
	"Oak",
	"Ash",
	"Willow",
	"Thorn",
	"Cedar",
}

func generateWorldTreeName(rng *rand.Rand) entity.EntityName {
	return (entity.EntityName)(fmt.Sprintf(
		"%s-%s",
		worldtreePrefixes[rng.IntN(len(worldtreePrefixes))],
		worldtreeNames[rng.IntN(len(worldtreeNames))],
	))
}
