package worldtree

type worldTreeAppearance int

const (
	Seed worldTreeAppearance = iota
	Sprout
	Young
	Mature
)

func (s worldTreeAppearance) String() string {
	switch s {
	case Seed:
		return "Seed"
	case Sprout:
		return "Sprout"
	case Young:
		return "Young"
	case Mature:
		return "Mature"
	default:
		return "Unknown"
	}
}
